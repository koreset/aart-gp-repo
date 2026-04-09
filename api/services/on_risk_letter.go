package services

import (
	appLog "api/log"
	"api/models"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// OnRiskLetterData is the composite response sent to the frontend so it can
// render the On Risk letter document (DOCX / PDF) client-side.
type OnRiskLetterData struct {
	Letter          models.OnRiskLetter              `json:"letter"`
	Quote           models.GroupPricingQuote          `json:"quote"`
	Scheme          models.GroupScheme                `json:"scheme"`
	Insurer         models.GroupPricingInsurerDetail  `json:"insurer"`
	BenefitSummary  []BenefitPremiumLine             `json:"benefit_summary"`
}

// BenefitPremiumLine represents one row in the benefits summary table.
type BenefitPremiumLine struct {
	Benefit       string  `json:"benefit"`
	AnnualPremium float64 `json:"annual_premium"`
}

// CreateOnRiskLetter persists an OnRiskLetter record for the given quote.
// It is called automatically inside AcceptGroupPricingQuote and can also be
// called independently to re-issue a letter.
func CreateOnRiskLetter(quoteId int, userName string) (*models.OnRiskLetter, error) {
	var quote models.GroupPricingQuote
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("id = ?", quoteId).First(&quote).Error
	})
	if err != nil {
		return nil, fmt.Errorf("quote not found: %v", err)
	}

	if quote.Status != models.StatusAccepted && quote.Status != models.StatusInForce {
		return nil, fmt.Errorf("quote must be in accepted or in_force status to issue an On Risk letter")
	}

	letterRef := fmt.Sprintf("ORL-%d-%d-%d", quote.SchemeID, quoteId, time.Now().Unix())

	letter := models.OnRiskLetter{
		QuoteID:          quoteId,
		SchemeID:         quote.SchemeID,
		LetterDate:       time.Now(),
		CommencementDate: quote.CommencementDate,
		CoverEndDate:     quote.CoverEndDate,
		GeneratedBy:      userName,
		LetterReference:  letterRef,
	}

	if err := DB.Create(&letter).Error; err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to create On Risk letter record")
		return nil, err
	}

	appLog.WithFields(map[string]interface{}{
		"letter_id":  letter.ID,
		"quote_id":   quoteId,
		"reference":  letterRef,
	}).Info("On Risk letter record created")

	return &letter, nil
}

// GetOnRiskLetterData fetches all data the frontend needs to render the On
// Risk letter document for a given quote.
func GetOnRiskLetterData(quoteId int) (*OnRiskLetterData, error) {
	var quote models.GroupPricingQuote
	var scheme models.GroupScheme
	var insurer models.GroupPricingInsurerDetail
	var letter models.OnRiskLetter
	var summaries []models.MemberRatingResultSummary

	// Fetch quote
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Preload("SchemeCategories").Where("id = ?", quoteId).First(&quote).Error
	})
	if err != nil {
		return nil, fmt.Errorf("quote not found: %v", err)
	}

	// Fetch scheme
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("id = ?", quote.SchemeID).First(&scheme).Error
	})
	if err != nil {
		appLog.WithField("error", err.Error()).Warn("Scheme not found for On Risk letter")
	}

	// Fetch insurer details (single record)
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Limit(1).Find(&insurer).Error
	})
	if err != nil {
		appLog.WithField("error", err.Error()).Warn("Insurer details not found for On Risk letter")
	}

	// Fetch most recent On Risk letter for this quote
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("quote_id = ?", quoteId).Order("id desc").First(&letter).Error
	})
	if err != nil {
		appLog.WithField("error", err.Error()).Warn("No existing On Risk letter found for quote")
	}

	// Fetch member rating summaries to build benefit premium breakdown
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("quote_id = ?", quoteId).Find(&summaries).Error
	})
	if err != nil {
		appLog.WithField("error", err.Error()).Warn("Member rating summaries not found")
	}

	benefitSummary := buildBenefitSummary(summaries)

	return &OnRiskLetterData{
		Letter:         letter,
		Quote:          quote,
		Scheme:         scheme,
		Insurer:        insurer,
		BenefitSummary: benefitSummary,
	}, nil
}

// buildBenefitSummary aggregates premium totals across all member rating
// summaries and returns only benefits with non-zero premiums.
func buildBenefitSummary(summaries []models.MemberRatingResultSummary) []BenefitPremiumLine {
	var gla, ptd, ci, sgla, ttd, phi, funeral float64

	for _, s := range summaries {
		gla += s.ExpTotalGlaAnnualOfficePremium
		ptd += s.ExpTotalPtdAnnualOfficePremium
		ci += s.ExpTotalCiAnnualOfficePremium
		sgla += s.ExpTotalSglaAnnualOfficePremium
		ttd += s.ExpTotalTtdAnnualOfficePremium
		phi += s.ExpTotalPhiAnnualOfficePremium
		funeral += s.ExpTotalFunAnnualOfficePremium
	}

	var lines []BenefitPremiumLine

	add := func(name string, amount float64) {
		if amount > 0 {
			lines = append(lines, BenefitPremiumLine{Benefit: name, AnnualPremium: amount})
		}
	}

	add("Group Life Assurance (GLA)", gla)
	add("Permanent Total Disability (PTD)", ptd)
	add("Critical Illness (CI)", ci)
	add("Spouse Group Life Assurance (SGLA)", sgla)
	add("Temporary Total Disability (TTD)", ttd)
	add("Permanent Health Insurance (PHI)", phi)
	add("Group Funeral", funeral)

	return lines
}

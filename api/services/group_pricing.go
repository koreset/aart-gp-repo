package services

import (
	appLog "api/log"
	"api/models"
	"api/utils"
	"compress/gzip"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/gammazero/workerpool"
	"github.com/jszwec/csvutil"
	"github.com/montanaflynn/stats"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func validateFamilyFuneralSumAssureds(quote models.GroupPricingQuote) error {
	for _, sc := range quote.SchemeCategories {
		if !sc.FamilyFuneralBenefit {
			continue
		}
		main := sc.FamilyFuneralMainMemberFuneralSumAssured
		label := sc.SchemeCategory
		if sc.FamilyFuneralSpouseFuneralSumAssured > main {
			return fmt.Errorf("scheme category '%s': spouse funeral sum assured cannot exceed main member's sum assured", label)
		}
		if sc.FamilyFuneralChildrenFuneralSumAssured > main {
			return fmt.Errorf("scheme category '%s': children funeral sum assured cannot exceed main member's sum assured", label)
		}
		if sc.FamilyFuneralAdultDependantSumAssured > main {
			return fmt.Errorf("scheme category '%s': adult dependant sum assured cannot exceed main member's sum assured", label)
		}
		if sc.FamilyFuneralParentFuneralSumAssured > main {
			return fmt.Errorf("scheme category '%s': parent funeral sum assured cannot exceed main member's sum assured", label)
		}
	}
	return nil
}

func GenerateGroupPricingQuote(quote models.GroupPricingQuote, user models.AppUser) error {
	logger := appLog.WithFields(map[string]interface{}{
		"user_email": user.UserEmail,
		"user_name":  user.UserName,
		"action":     "GenerateGroupPricingQuote",
	})

	logger.Info("Starting group pricing quote generation")

	if err := validateFamilyFuneralSumAssureds(quote); err != nil {
		return err
	}

	// Capture before state if this is an update (quote.ID > 0)
	var __gpqBefore models.GroupPricingQuote
	if quote.ID > 0 {
		_ = DB.Where("id = ?", quote.ID).First(&__gpqBefore).Error
	}

	var err error

	// we need the group benefit mappers here
	logger.Debug("Retrieving group benefit mappers")
	var benefitMappers []models.GroupBenefitMapper
	err = DB.Find(&benefitMappers).Error
	if err != nil {
		//just log the error and continue
		logger.WithField("error", err.Error()).Error("Failed to retrieve benefit mappers, continuing anyway")
	}

	// Check if the quote is in edit mode
	if quote.EditMode {
		logger.WithField("quote_id", quote.ID).Debug("Edit mode enabled for existing quote")
	} else {
		logger.Debug("Creating new group pricing quote")
	}

	//// get a count of quotes associated with the scheme name
	//var existingQuote models.GroupPricingQuote
	//err = DB.Where("name = ?", quote.SchemeName).First(&existingQuote).Error

	// This might be an edit so ID might be set. If not, we will create a new one. Otherwise, we will update the existing one.
	if !quote.EditMode {
		// discovered the v-date-input component has a bug that sends the date as the day before the actual date.
		// This is a workaround to fix that.
		var quoteCount int64
		_ = DB.Model(models.GroupPricingQuote{}).Where("scheme_name = ? ", quote.SchemeName).Count(&quoteCount).Error

		if quoteCount == 0 {
			quoteCount = 1
		} else {
			quoteCount = quoteCount + 1
		}
		quote.QuoteName = strings.ReplaceAll(quote.SchemeName, " ", "_") + "_" + strconv.FormatInt(quoteCount, 10)

		quote.CommencementDate = quote.CommencementDate.AddDate(0, 0, 0)
		quote.CoverEndDate = quote.CommencementDate.AddDate(0, 12, 0)
		quote.CreationDate = time.Now()
		quote.ModificationDate = time.Now()
		quote.CreatedBy = user.UserName
		quote.Status = models.StatusInProgress
		quote.SchemeQuoteStatus = models.StatusNotInEffect
		// Prevent duplicate quotes: same scheme name + same quote type
		var dupCount int64
		_ = DB.Model(&models.GroupPricingQuote{}).
			Where("scheme_name = ? AND quote_type = ?", quote.SchemeName, quote.QuoteType).
			Count(&dupCount).Error
		if dupCount > 0 {
			return fmt.Errorf("a %s quote for scheme '%s' already exists", quote.QuoteType, quote.SchemeName)
		}

		addMappers(&quote, benefitMappers)
		for i := range quote.SchemeCategories {
			quote.SchemeCategories[i].ID = 0
		}
		err = DB.Create(&quote).Error
		if err != nil {
			return err
		}

		if quote.QuoteType == "New Business" {
			// Create a new scheme
			var scheme models.GroupScheme
			scheme.Name = quote.SchemeName
			scheme.ContactPerson = quote.SchemeContact
			scheme.ContactEmail = quote.SchemeEmail
			scheme.DistributionChannel = quote.DistributionChannel
			if quote.DistributionChannel != models.ChannelDirect {
				scheme.BrokerId = quote.QuoteBroker.ID
			}
			scheme.CreationDate = time.Now()
			scheme.CreatedBy = user.UserName
			scheme.QuoteId = quote.ID
			scheme.Status = models.StatusQuoted
			scheme.QuoteInForce = quote.QuoteName
			scheme.RenewalDate = quote.CoverEndDate
			scheme.CoverStartDate = quote.CommencementDate // Default date
			scheme.CoverEndDate = quote.CoverEndDate
			scheme.CommencementDate = quote.CommencementDate

			err = DB.Create(&scheme).Error
			if err != nil {
				return err
			}

			quote.SchemeID = scheme.ID
			err = DB.Save(&quote).Error
		}

	} else {

		addMappers(&quote, benefitMappers)
		quote.CommencementDate = quote.CommencementDate.AddDate(0, 0, 0)
		quote.CoverEndDate = quote.CommencementDate.AddDate(0, 12, 0)
		quote.ModificationDate = time.Now()
		quote.ModifiedBy = user.UserName
		updatedSchemeCategories := quote.SchemeCategories
		err = DB.Save(&quote).Error
		if err != nil {
			return err
		}

		for _, schemeCategory := range updatedSchemeCategories {
			// Update the scheme categories for the quote
			var existingSchemeCategory models.SchemeCategory
			err = DB.Where("quote_id = ? AND scheme_category = ?", quote.ID, schemeCategory.SchemeCategory).First(&existingSchemeCategory).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// If not found, create a new one
					schemeCategory.QuoteId = quote.ID
					err = DB.Create(&schemeCategory).Error
				} else {
					return err
				}
			} else {
				// If found, update it
				schemeCategory.ID = existingSchemeCategory.ID
				err = DB.Save(&schemeCategory).Error
			}
			if err != nil {
				return err
			}
		}
	}

	//if quote.QuoteType == "Renewal" && !quote.EditMode {
	//	// Get the existing scheme
	//	var scheme models.GroupScheme
	//	err = DB.Where("name = ?", quote.SchemeName).First(&scheme).Error
	//	if err != nil {
	//		return err
	//	}
	//	addMappers(&quote, benefitMappers)
	//	quote.SchemeID = scheme.ID
	//	quote.SchemeName = scheme.Name
	//	quote.Status = models.StatusQuoted
	//	quote.SchemeQuoteStatus = models.StatusNotInEffect
	//
	//	err = DB.Create(&quote).Error
	//}

	// Load after state and write generic audit (CREATE vs UPDATE)
	var __gpqAfter models.GroupPricingQuote
	if err2 := DB.Where("id = ?", quote.ID).First(&__gpqAfter).Error; err2 == nil && __gpqAfter.ID > 0 {
		action := "UPDATE"
		if __gpqBefore.ID == 0 {
			action = "CREATE"
		}
		_ = writeAudit(DB, AuditContext{
			Area:      "group-pricing",
			Entity:    "group_pricing_quotes",
			EntityID:  strconv.Itoa(__gpqAfter.ID),
			Action:    action,
			ChangedBy: user.UserName,
		}, __gpqBefore, __gpqAfter)
	}

	return nil
}

func UpdateGroupPricingQuote(quote models.GroupPricingQuote, user models.AppUser) error {
	if err := validateFamilyFuneralSumAssureds(quote); err != nil {
		return err
	}

	// Capture incoming status before transaction so we can fire notifications after commit.
	incomingStatus := quote.Status

	err := DB.Transaction(func(tx *gorm.DB) error {
		// Load before state
		var before models.GroupPricingQuote
		if err := tx.Where("id = ?", quote.ID).First(&before).Error; err != nil {
			return err
		}

		// Apply updates
		quote.ModifiedBy = user.UserName
		quote.ModificationDate = time.Now()
		if err := tx.Save(&quote).Error; err != nil {
			return err
		}

		// Load after state
		var after models.GroupPricingQuote
		if err := tx.Where("id = ?", quote.ID).First(&after).Error; err != nil {
			return err
		}

		// Write generic audit
		if err := writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "group_pricing_quotes",
			EntityID:  strconv.Itoa(quote.ID),
			Action:    "UPDATE",
			ChangedBy: user.UserName,
		}, before, after); err != nil {
			return err
		}
		go func(id int) { ScoreQuote(id) }(quote.ID)
		return nil
	})
	if err != nil {
		return err
	}

	// Fire notification hooks after successful commit
	switch incomingStatus {
	case models.StatusSubmitted:
		go NotifyQuoteSubmitted(quote, user)
	case models.StatusRejected:
		go NotifyQuoteRejected(quote, user, "")
	case models.StatusPendingReview, "Pending Review":
		go NotifyQuoteSubmitted(quote, user) // Reuse submitted notification for pending review status
	}

	return nil
}

func ApproveGroupPricingQuote(quoteId string, user models.AppUser) error {
	var quote models.GroupPricingQuote

	err := DB.Transaction(func(tx *gorm.DB) error {
		var memberRatingResultSummary models.MemberRatingResultSummary

		// Load before state
		if err := tx.Where("id = ?", quoteId).First(&quote).Error; err != nil {
			return err
		}
		before := quote

		if err := tx.Where("quote_id = ?", quoteId).First(&memberRatingResultSummary).Error; err != nil {
			return err
		}

		// Apply updates
		quote.Status = models.StatusApproved
		quote.ModifiedBy = user.UserName
		quote.ModificationDate = time.Now()
		if err := tx.Save(&quote).Error; err != nil {
			return err
		}
		memberRatingResultSummary.IfStatus = "Approved"
		if err := tx.Save(&memberRatingResultSummary).Error; err != nil {
			return err
		}

		// Load after state
		var after models.GroupPricingQuote
		if err := tx.Where("id = ?", quoteId).First(&after).Error; err != nil {
			return err
		}
		// Write generic audit
		if err := writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "group_pricing_quotes",
			EntityID:  quoteId,
			Action:    "UPDATE",
			ChangedBy: user.UserName,
		}, before, after); err != nil {
			return err
		}
		go func(id int) { ScoreQuote(id) }(intQuoteID(quoteId))
		return nil
	})
	if err != nil {
		return err
	}

	// Fire notification after successful commit
	go NotifyQuoteApproved(quote, user)

	return nil
}

func AcceptGroupPricingQuote(quoteId string, commencementDate string, term string, user models.AppUser) error {
	var quote models.GroupPricingQuote

	txErr := DB.Transaction(func(tx *gorm.DB) error {
		var memberRatingSummary []models.MemberRatingResultSummary
		var err error
		var groupAgeBands []models.GroupPricingAgeBands
		var groupBenefits []models.GroupBusinessBenefits
		var exposureData []models.GroupSchemeExposure
		var GroupPricingInsurerDetail models.GroupPricingInsurerDetail
		var insurerYearEndMonth int

		termMonths, _ := strconv.Atoi(term)

		err = tx.Where("id = ?", quoteId).First(&quote).Error
		if err != nil {
			return err
		}
		// Snapshot before state for auditing
		before := quote

		// Parse commencement date
		parsedCommencementDate, err := time.Parse("2006-01-02", commencementDate)
		if err != nil {
			return fmt.Errorf("invalid commencement date format: %v", err)
		}

		// Check the status of the quote
		if quote.Status != models.StatusApproved {
			return fmt.Errorf("the quote has not yet been approved. Please approve the quote before accepting it")
		}

		if quote.Status == models.StatusInForce || quote.Status == models.StatusAccepted {
			return fmt.Errorf("the quote is already accepted or in force")
		}

		groupAgeBands, err = GetGroupPricingAgeBands(context.Background())
		groupBenefits, err = GetGroupPricingBenefits(context.Background())

		// Update quote with the provided commencement date
		quote.CommencementDate = parsedCommencementDate
		quote.CoverEndDate = parsedCommencementDate.AddDate(0, termMonths, 0)
		quote.SchemeQuoteStatus = models.StatusNotInEffect
		quote.Status = models.StatusAccepted
		quote.ModifiedBy = user.UserName
		quote.ModificationDate = time.Now()

		err = tx.Save(&quote).Error
		if err != nil {
			return err
		}

		if err := tx.Limit(1).Find(&GroupPricingInsurerDetail).Error; err == nil {
			insurerYearEndMonth = GroupPricingInsurerDetail.YearEndMonth
		}

		err = tx.Where("quote_id = ?", quoteId).Find(&memberRatingSummary).Error
		if err != nil {
			return err
		}

		if len(memberRatingSummary) > 0 {
			for i := range memberRatingSummary {
				_, memberRatingSummary[i].FinancialYear = getGroupRiskQuotingFinancialYear(quote.CommencementDate, insurerYearEndMonth)
				memberRatingSummary[i].IfStatus = models.StatusAccepted
			}
			err = tx.Save(&memberRatingSummary).Error
			if err != nil {
				return err
			}
		}

		// we need to copy member data to member data in force
		// first check if there is existing data in member data in force for the quote_id and delete it
		tx.Model(&models.GPricingMemberDataInForce{}).Where("quote_id = ?", quoteId).Delete(&models.GPricingMemberDataInForce{})

		var gmdif []models.GPricingMemberDataInForce
		var selectedCategories []models.SchemeCategory

		err = tx.Model(&models.SchemeCategory{}).Where("quote_id = ?", quoteId).Find(&selectedCategories).Error
		if err != nil {
			return err
		}

		if quote.QuoteType == "New Business" {
			err = tx.Model(&models.GPricingMemberData{}).Where("quote_id = ?", quoteId).Find(&gmdif).Error
			if err != nil {
				return err
			}

			// Collect RSA IDs for bulk validation via CheckID API
			var rsaIDsToValidate []string
			for _, m := range gmdif {
				idType := strings.ToUpper(strings.TrimSpace(m.MemberIdType))
				if (idType == "RSA_ID" || idType == "ID" || idType == "RSA_ISD") && strings.TrimSpace(m.MemberIdNumber) != "" {
					rsaIDsToValidate = append(rsaIDsToValidate, strings.TrimSpace(m.MemberIdNumber))
				}
			}
			rsaIDResults, rsaIDErr := utils.ValidateRSAIDsBulk(rsaIDsToValidate)
			if rsaIDErr != nil {
				return fmt.Errorf("ID validation service error: %v", rsaIDErr)
			}

			for i := range gmdif {
				if strings.TrimSpace(gmdif[i].EmployeeNumber) == "" {
					return fmt.Errorf("employee number for member %s is empty. All members must have an employee number for New Business quotes", gmdif[i].MemberName)
				}

				if quote.ObligationType == "Compulsory" {
					if strings.TrimSpace(gmdif[i].MemberIdNumber) == "" {
						return fmt.Errorf("member id number for member %s is empty. All members must have an ID or Passport number for Compulsory quotes", gmdif[i].MemberName)
					}
				}

				idType := strings.ToUpper(strings.TrimSpace(gmdif[i].MemberIdType))
				if (idType == "RSA_ID" || idType == "ID" || idType == "RSA_ISD") && strings.TrimSpace(gmdif[i].MemberIdNumber) != "" {
					if valid, ok := rsaIDResults[strings.TrimSpace(gmdif[i].MemberIdNumber)]; ok && !valid {
						return fmt.Errorf("invalid RSA ID '%s' for member %s", gmdif[i].MemberIdNumber, gmdif[i].MemberName)
					}
				}

				gmdif[i].ID = 0
				gmdif[i].IsOriginalMember = true
				gmdif[i].Status = "Active"
				gmdif[i].Year = 2025
				//gmdif[i].ExitDate =0 time.Now() // reset exit date
				for _, category := range selectedCategories {
					if category.SchemeCategory == gmdif[i].SchemeCategory {
						if quote.UseGlobalSalaryMultiple {
							if category.GlaBenefit {
								gmdif[i].Benefits.GlaEnabled = true
								gmdif[i].Benefits.GlaMultiple = category.GlaSalaryMultiple
							}
							if category.SglaBenefit {
								gmdif[i].Benefits.SglaEnabled = true
								gmdif[i].Benefits.SglaMultiple = category.SglaSalaryMultiple
							}
							if category.PtdBenefit {
								gmdif[i].Benefits.PtdEnabled = true
								gmdif[i].Benefits.PtdMultiple = category.PtdSalaryMultiple
							}
							if category.CiBenefit {
								gmdif[i].Benefits.CiEnabled = true
								gmdif[i].Benefits.CiMultiple = category.CiCriticalIllnessSalaryMultiple
							}
							if category.PhiBenefit {
								gmdif[i].Benefits.PhiEnabled = true
								gmdif[i].Benefits.PhiMultiple = category.PhiIncomeReplacementPercentage / 100
							}
							if category.TtdBenefit {
								gmdif[i].Benefits.TtdEnabled = true
								gmdif[i].Benefits.TtdMultiple = category.TtdIncomeReplacementPercentage / 100
							}
							if category.FamilyFuneralBenefit {
								gmdif[i].Benefits.GffEnabled = true
							}
						} else {
							if category.GlaBenefit {
								gmdif[i].Benefits.GlaEnabled = true
							}
							if category.SglaBenefit {
								gmdif[i].Benefits.SglaEnabled = true
							}
							if category.PtdBenefit {
								gmdif[i].Benefits.PtdEnabled = true
							}
							if category.CiBenefit {
								gmdif[i].Benefits.CiEnabled = true
							}
							if category.PhiBenefit {
								gmdif[i].Benefits.PhiEnabled = true
								gmdif[i].Benefits.PhiMultiple = category.PhiIncomeReplacementPercentage / 100
							}
							if category.TtdBenefit {
								gmdif[i].Benefits.TtdEnabled = true
								gmdif[i].Benefits.TtdMultiple = category.TtdIncomeReplacementPercentage / 100
							}
							if category.FamilyFuneralBenefit {
								gmdif[i].Benefits.GffEnabled = true
							}
						}
					}
				}
			}
			if len(gmdif) > 0 {
				err = tx.CreateInBatches(&gmdif, 100).Error
				if err != nil {
					return err
				}
			}
		}

		var quoteStats models.GroupRiskQuoteStats
		var scheme models.GroupScheme
		err = tx.Where("name = ?", quote.SchemeName).First(&scheme).Error
		if err != nil {
			return err
		}

		if quote.QuoteType == "New Business" {
			scheme.CoverStartDate, _ = utils.ParseDateString(commencementDate)
			termMonths, _ = strconv.Atoi(term)
			scheme.CoverEndDate = scheme.CoverStartDate.AddDate(0, termMonths, 0)
		}

		for _, mrs := range memberRatingSummary {
			quoteStats.AnnualPremium += mrs.ExpTotalAnnualPremiumExclFuneral + models.ComputeOfficePremium(mrs.ExpTotalFunAnnualRiskPremium, &mrs)
			quoteStats.Commission += mrs.TotalCommission
			quoteStats.ExpectedExpenses += mrs.TotalExpenses
			quoteStats.ExpectedClaims += mrs.TotalExpectedClaims
			quoteStats.ExpectedGlaClaims += mrs.ExpTotalGlaAnnualRiskPremium
			quoteStats.ExpectedPtdClaims += mrs.ExpTotalPtdAnnualRiskPremium
			quoteStats.ExpectedCiClaims += mrs.ExpTotalCiAnnualRiskPremium
			quoteStats.ExpectedSglaClaims += mrs.ExpTotalSglaAnnualRiskPremium
			quoteStats.ExpectedTtdClaims += mrs.ExpTotalTtdAnnualRiskPremium
			quoteStats.ExpectedPhiClaims += mrs.ExpTotalPhiAnnualRiskPremium
			quoteStats.ExpectedFunClaims += mrs.ExpTotalFunAnnualRiskPremium
			quoteStats.MemberCount += int(mrs.MemberCount)
			quoteStats.GlaAnnualPremium += models.ComputeOfficePremium(mrs.ExpTotalGlaAnnualRiskPremium, &mrs)
			quoteStats.PtdAnnualPremium += models.ComputeOfficePremium(mrs.ExpTotalPtdAnnualRiskPremium, &mrs)
			quoteStats.CiAnnualPremium += models.ComputeOfficePremium(mrs.ExpTotalCiAnnualRiskPremium, &mrs)
			quoteStats.SglaAnnualPremium += models.ComputeOfficePremium(mrs.ExpTotalSglaAnnualRiskPremium, &mrs)
			quoteStats.TtdAnnualPremium += models.ComputeOfficePremium(mrs.ExpTotalTtdAnnualRiskPremium, &mrs)
			quoteStats.PhiAnnualPremium += models.ComputeOfficePremium(mrs.ExpTotalPhiAnnualRiskPremium, &mrs)
			quoteStats.FuneralAnnualPremium += models.ComputeOfficePremium(mrs.ExpTotalFunAnnualRiskPremium, &mrs)
		}

		if quote.QuoteType == "New Business" {
			scheme.CoverStartDate = quote.CommencementDate
			scheme.CommencementDate = quote.CommencementDate
			scheme.RenewalDate = scheme.CoverEndDate.AddDate(0, 0, 1)
			scheme.InForce = true
			scheme.Status = models.StatusAccepted
			scheme.NewBusiness = true
		}

		// Optimized Exposure Data Calculation
		type BenefitSum struct {
			AgeBand                   string
			Gender                    string
			GlaCappedSumAssured       float64
			PtdCappedSumAssured       float64
			CiCappedSumAssured        float64
			SpouseGlaCappedSumAssured float64
			TtdCappedIncome           float64
			PhiCappedIncome           float64
		}

		var benefitSums []BenefitSum
		err = tx.Table("member_rating_results").
			Select("age_band, gender, SUM(gla_capped_sum_assured) as gla_capped_sum_assured, SUM(ptd_capped_sum_assured) as ptd_capped_sum_assured, SUM(ci_capped_sum_assured) as ci_capped_sum_assured, SUM(spouse_gla_capped_sum_assured) as spouse_gla_capped_sum_assured, SUM(ttd_capped_income) as ttd_capped_income, SUM(phi_capped_income) as phi_capped_income").
			Where("quote_id = ?", quoteId).
			Group("age_band, gender").
			Scan(&benefitSums).Error

		if err != nil {
			return err
		}

		// Organize benefitSums for quick lookup: map[ageBand]map[gender]BenefitSum
		sumMap := make(map[string]map[string]BenefitSum)
		for _, bs := range benefitSums {
			if _, ok := sumMap[bs.AgeBand]; !ok {
				sumMap[bs.AgeBand] = make(map[string]BenefitSum)
			}
			sumMap[bs.AgeBand][bs.Gender] = bs
		}

		financialYear := 0
		if len(memberRatingSummary) > 0 {
			financialYear = memberRatingSummary[0].FinancialYear
		}

		for _, ageBand := range groupAgeBands {
			for _, groupBenefit := range groupBenefits {
				exposureDataPoint := models.GroupSchemeExposure{
					QuoteId:       quote.ID,
					SchemeName:    scheme.Name,
					Industry:      quote.Industry,
					AgeBand:       ageBand.Name,
					MinAge:        ageBand.MinAge,
					MaxAge:        ageBand.MaxAge,
					Benefit:       groupBenefit.Name,
					FinancialYear: financialYear,
					QuoteStatus:   string(quote.Status),
				}

				maleBS := sumMap[ageBand.Name]["M"]
				femaleBS := sumMap[ageBand.Name]["F"]

				var maleVal, femaleVal float64
				switch groupBenefit.Name {
				case "GLA":
					maleVal = maleBS.GlaCappedSumAssured
					femaleVal = femaleBS.GlaCappedSumAssured
				case "PTD":
					maleVal = maleBS.PtdCappedSumAssured
					femaleVal = femaleBS.PtdCappedSumAssured
				case "CI":
					maleVal = maleBS.CiCappedSumAssured
					femaleVal = femaleBS.CiCappedSumAssured
				case "SGLA":
					maleVal = maleBS.SpouseGlaCappedSumAssured
					femaleVal = femaleBS.SpouseGlaCappedSumAssured
				case "TTD":
					maleVal = maleBS.TtdCappedIncome
					femaleVal = femaleBS.TtdCappedIncome
				case "PHI":
					maleVal = maleBS.PhiCappedIncome
					femaleVal = femaleBS.PhiCappedIncome
				}

				exposureDataPoint.MaleSumAssured = maleVal
				exposureDataPoint.FemaleSumAssured = femaleVal
				exposureDataPoint.TotalSumAssured = maleVal + femaleVal

				exposureData = append(exposureData, exposureDataPoint)
			}
		}

		tx.Where("quote_id = ?", quoteId).Delete(&models.GroupSchemeExposure{})
		if len(exposureData) > 0 {
			err = tx.CreateInBatches(&exposureData, 100).Error
			if err != nil {
				return err
			}
		}

		tx.Where("quote_id = ?", quoteId).Delete(&models.GroupRiskQuoteStats{})
		quoteStats.QuoteID = quote.ID
		if quoteStats.AnnualPremium > 0 {
			quoteStats.ExpectedClaimsRatio = (quoteStats.ExpectedClaims / quoteStats.AnnualPremium)
		}
		quoteStats.CoverStartDate, _ = utils.ParseDateString(commencementDate)
		quoteStats.CoverEndDate = quote.CommencementDate.AddDate(0, termMonths, 0)
		quoteStats.Creator = user.UserName
		quoteStats.CreationDate = time.Now()
		err = tx.Create(&quoteStats).Error
		if err != nil {
			return err
		}

		err = tx.Save(&scheme).Error
		if err != nil {
			return err
		}

		// Audit logging inside transaction or after successful commit?
		// Usually better after commit, but since it's an internal function,
		// we can do it after the transaction block or use a hook.
		// Original code did it before returning.

		var after models.GroupPricingQuote
		if err := tx.Where("id = ?", quoteId).First(&after).Error; err == nil {
			_ = writeAudit(tx, AuditContext{
				Area:      "group-pricing",
				Entity:    "group_pricing_quotes",
				EntityID:  quoteId,
				Action:    "UPDATE",
				ChangedBy: user.UserName,
			}, before, after)
		}

		go func(id int) { ScoreQuote(id) }(intQuoteID(quoteId))
		return nil
	})
	if txErr != nil {
		return txErr
	}

	// Fire notification after successful commit
	go NotifyQuoteAccepted(quote, user)

	return nil
}

func GetGroupPricingQuotes(filter string) ([]models.GroupPricingQuote, error) {
	var quotes []models.GroupPricingQuote
	var err error

	if filter == "" {
		err = DB.Preload("SchemeCategories").Order("id desc").Find(&quotes).Error
	} else {
		err = DB.Preload("SchemeCategories").Where("status = ?", filter).Order("id desc").Find(&quotes).Error
	}
	return quotes, err
}

// GetGroupPricingQuotesBySchemeID returns all quotes linked to a given group scheme via scheme_id
func GetGroupPricingQuotesBySchemeID(schemeId string) ([]models.GroupPricingQuote, error) {
	var quotes []models.GroupPricingQuote
	err := DB.Preload("SchemeCategories").Where("scheme_id = ?", schemeId).Order("id desc").Find(&quotes).Error
	return quotes, err
}

func GetGroupPricingQuote(id string) (models.GroupPricingQuote, error) {
	var quote models.GroupPricingQuote
	var err error

	// Use QueryWithContextAndCache to cache the results for better performance
	cacheKey := fmt.Sprintf("group_pricing_quote_%s", id)
	err = QueryWithContextAndCache(cacheKey, &quote, 5*time.Second, func(ctx context.Context) error {
		// First get the quote
		if dbErr := DB.WithContext(ctx).Where("id = ?", id).Preload("SchemeCategories").First(&quote).Error; err != nil {
			return dbErr
		}

		// Use a single query with subqueries to get all counts in one database call
		type Counts struct {
			MemberDataCount              int64
			ClaimsExperienceCount        int64
			ExperienceRateOverridesCount int64
			MemberRatingResultCount      int64
			MemberPremiumScheduleCount   int64
			BordereauxCount              int64
		}

		var counts Counts
		// bordereaux_count tracks the number of rows that would be projected
		// for the quote's Bordereaux preview/download. Each MemberRatingResult
		// projects to exactly one bordereaux row, so we count that table.
		countQuery := DB.WithContext(ctx).Raw(`
			SELECT
				(SELECT COUNT(*) FROM g_pricing_member_data WHERE quote_id = ?) as member_data_count,
				(SELECT COUNT(*) FROM group_pricing_claims_experiences WHERE quote_id = ?) as claims_experience_count,
				(SELECT COUNT(*) FROM group_pricing_experience_rate_overrides WHERE quote_id = ?) as experience_rate_overrides_count,
				(SELECT COUNT(*) FROM member_rating_results WHERE quote_id = ?) as member_rating_result_count,
				(SELECT COUNT(*) FROM member_premium_schedules WHERE quote_id = ?) as member_premium_schedule_count,
				(SELECT COUNT(*) FROM member_rating_results WHERE quote_id = ?) as bordereaux_count
		`, id, id, id, id, id, id)

		if err = countQuery.Scan(&counts).Error; err != nil {
			// If the optimized query fails, fall back to individual counts
			appLog.WithField("error", err.Error()).Warn("Optimized count query failed, falling back to individual counts")

			// Get counts individually but still using the context
			DB.WithContext(ctx).Model(&models.GPricingMemberData{}).Where("quote_id = ?", id).Count(&counts.MemberDataCount)
			DB.WithContext(ctx).Model(&models.GroupPricingClaimsExperience{}).Where("quote_id = ?", id).Count(&counts.ClaimsExperienceCount)
			DB.WithContext(ctx).Model(&models.GroupPricingExperienceRateOverride{}).Where("quote_id = ?", id).Count(&counts.ExperienceRateOverridesCount)
			DB.WithContext(ctx).Model(&models.MemberRatingResult{}).Where("quote_id = ?", id).Count(&counts.MemberRatingResultCount)
			DB.WithContext(ctx).Model(&models.MemberPremiumSchedule{}).Where("quote_id = ?", id).Count(&counts.MemberPremiumScheduleCount)
			// Bordereaux is projected from MemberRatingResult on-the-fly, so
			// the bordereaux row count tracks the rating-result row count.
			DB.WithContext(ctx).Model(&models.MemberRatingResult{}).Where("quote_id = ?", id).Count(&counts.BordereauxCount)
		}

		// Update the quote with the counts
		if quote.QuoteType == "Renewal" {
			// For Renewal quotes, MemberDataCount must come from GPMemberDataInforce table using SchemeName
			var inforceCount int64
			DB.WithContext(ctx).Model(&models.GPricingMemberDataInForce{}).Where("scheme_name = ?", quote.SchemeName).Count(&inforceCount)
			quote.MemberDataCount = int(inforceCount)
		} else {
			quote.MemberDataCount = int(counts.MemberDataCount)
		}

		quote.ClaimsExperienceCount = int(counts.ClaimsExperienceCount)
		quote.ExperienceRateOverridesCount = int(counts.ExperienceRateOverridesCount)
		quote.MemberRatingResultCount = int(counts.MemberRatingResultCount)
		quote.MemberPremiumScheduleCount = int(counts.MemberPremiumScheduleCount)
		quote.BordereauxCount = int(counts.BordereauxCount)

		var memberIndicativeDataSet []models.MemberIndicativeDataSet
		DB.Where("quote_id = ?", id).Find(&memberIndicativeDataSet)
		quote.MemberIndicativeDataSet = memberIndicativeDataSet

		return nil
	})

	return quote, err
}

// GetGroupPricingQuoteBySchemeName finds a scheme by its name, then retrieves
// the associated group pricing quote using the scheme's quote name (if present)
// or falls back to the most recent quote for the scheme. It returns the
// enriched quote using the existing GetGroupPricingQuote logic to ensure
// consistent preload and counters.
func GetGroupPricingQuoteBySchemeName(schemeName string) (models.GroupPricingQuote, error) {
	var empty models.GroupPricingQuote

	name := strings.TrimSpace(schemeName)
	if name == "" {
		return empty, fmt.Errorf("scheme name is required")
	}

	var scheme models.GroupScheme
	// Case-insensitive exact match on scheme name
	if err := DB.Where("LOWER(name) = LOWER(?)", name).First(&scheme).Error; err != nil {
		return empty, err
	}

	// Prefer explicit quote name recorded on the scheme when available
	quoteName := strings.TrimSpace(scheme.QuoteInForce)
	if quoteName != "" {
		var q models.GroupPricingQuote
		if err := DB.Where("quote_name = ?", quoteName).First(&q).Error; err != nil {
			return empty, err
		}
		return GetGroupPricingQuote(strconv.Itoa(q.ID))
	}

	// Fallback: get latest quote for the scheme (by id desc)
	var latest models.GroupPricingQuote
	if err := DB.Where("scheme_id = ?", scheme.ID).Order("id desc").First(&latest).Error; err != nil {
		return empty, err
	}
	return GetGroupPricingQuote(strconv.Itoa(latest.ID))
}

func DeleteGroupPricingQuote(id string, user models.AppUser) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var quote models.GroupPricingQuote
		// Load before snapshot
		if err := tx.Where("id = ?", id).First(&quote).Error; err != nil {
			return err
		}

		// delete associated schemes first
		var schemes []models.GroupScheme
		if err := tx.Where("quote_id = ?", id).Find(&schemes).Error; err != nil {
			return err
		}
		for _, scheme := range schemes {
			if scheme.InForce {
				return fmt.Errorf("cannot delete quote as it has schemes that are in force")
			}
		}
		for _, scheme := range schemes {
			if err := tx.Where("id = ?", scheme.ID).Delete(&models.GroupScheme{}).Error; err != nil {
				return err
			}
		}

		// delete all tables with quote_id
		if err := tx.Where("quote_id = ?", id).Delete(&models.SchemeCategory{}).Error; err != nil {
			return err
		}
		if err := tx.Where("quote_id = ?", id).Delete(&models.GroupRiskQuoteStats{}).Error; err != nil {
			return err
		}
		if err := tx.Where("quote_id = ?", id).Delete(&models.QuoteWinProbability{}).Error; err != nil {
			return err
		}
		if err := tx.Where("quote_id = ?", id).Delete(&models.MemberIndicativeDataSet{}).Error; err != nil {
			return err
		}
		if err := tx.Where("quote_id = ?", id).Delete(&models.OnRiskLetter{}).Error; err != nil {
			return err
		}
		if err := tx.Where("quote_id = ?", id).Delete(&models.GPricingMemberDataInForce{}).Error; err != nil {
			return err
		}
		if err := tx.Where("quote_id = ?", id).Delete(&models.MovementMemberRatingResult{}).Error; err != nil {
			return err
		}
		if err := tx.Where("quote_id = ?", id).Delete(&models.GroupSchemeExposure{}).Error; err != nil {
			return err
		}
		if err := tx.Where("quote_id = ?", id).Delete(&models.GPricingMemberData{}).Error; err != nil {
			return err
		}
		if err := tx.Where("quote_id = ?", id).Delete(&models.GroupPricingClaimsExperience{}).Error; err != nil {
			return err
		}
		if err := tx.Where("quote_id = ?", id).Delete(&models.MemberRatingResult{}).Error; err != nil {
			return err
		}
		if err := tx.Where("quote_id = ?", id).Delete(&models.MemberPremiumSchedule{}).Error; err != nil {
			return err
		}
		if err := tx.Where("quote_id = ?", id).Delete(&models.Bordereaux{}).Error; err != nil {
			return err
		}
		if err := tx.Where("quote_id = ?", id).Delete(&models.MemberRatingResultSummary{}).Error; err != nil {
			return err
		}
		if err := tx.Where("quote_id = ?", id).Delete(&models.HistoricalCredibilityData{}).Error; err != nil {
			return err
		}

		// delete the quote
		if err := tx.Where("id = ? ", id).Delete(&models.GroupPricingQuote{}).Error; err != nil {
			return err
		}

		// Write generic DELETE audit with snapshot in prev_values
		if err := writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "group_pricing_quotes",
			EntityID:  id,
			Action:    "DELETE",
			ChangedBy: user.UserName,
		}, quote, struct{}{}); err != nil {
			return err
		}
		return nil
	})
}

// MemberRateResult bundles the per-member outputs produced by PopulateRatesPerMember
// so that the caller can collect them without shared-state mutation.
type MemberRateResult struct {
	Rating          models.MemberRatingResult
	PremiumSchedule models.MemberPremiumSchedule
	Bordereaux      models.Bordereaux
}

type TheoreticalRiskTotal struct {
	GlaSumRiskPremium float64
	PtdSumRiskPremium float64
	TtdSumRiskPremium float64
	PhiSumRiskPremium float64
	CiSumRiskPremium  float64
	GlaSumAssured     float64
	PtdSumAssured     float64
	TtdSumAssured     float64
	PhiSumAssured     float64
	CiSumAssured      float64
}

// sendCalculationProgress pushes a progress update to the user via WebSocket.
// It is a no-op if the hub is not initialized (e.g. in tests).
func sendCalculationProgress(userEmail string, progress CalculationProgress) {
	hub := GetHub()
	if hub == nil {
		return
	}
	hub.SendToUser(userEmail, WSEnvelope{
		Type:    WSCalculationProgress,
		Payload: progress,
	})
}

func CalculateGroupPricingQuote(quoteId string, basis string, credibility float64, user models.AppUser) error {
	// TODO: Implement the logic to calculate the group pricing quote based on the provided quoteId and basis and optionally, credibility.
	startTime := time.Now()

	logger := appLog.WithFields(map[string]interface{}{
		"user_email": user.UserEmail,
		"user_name":  user.UserName,
		"quote_id":   quoteId,
		"basis":      basis,
		"action":     "CalculateGroupPricingQuote",
	})

	logger.Info("Starting group pricing quote calculation")

	GroupPricingCache.Clear()
	logger.Debug("Group pricing cache cleared")

	var groupQuote models.GroupPricingQuote

	dbStartTime := time.Now()
	DB.Where("id = ?", quoteId).Preload("SchemeCategories").First(&groupQuote)
	dbElapsed := time.Since(dbStartTime)
	logger.WithFields(map[string]interface{}{
		"scheme_name": groupQuote.SchemeName,
		"elapsed_ms":  dbElapsed.Milliseconds(),
	}).Debug("Retrieved group quote")

	// Fetch tax table and tiered income replacement tiers once for this quote.
	taxTable, err := GetTaxTableByRiskRateCode(groupQuote.RiskRateCode)
	if err != nil {
		logger.WithField("error", err.Error()).Warn("Failed to load tax table; take-home pay will equal gross salary")
	}
	tieredIncomeTiers, err := GetTieredIncomeReplacementTiers(groupQuote.RiskRateCode)
	if err != nil {
		logger.WithField("error", err.Error()).Warn("Failed to load tiered income replacement tiers")
	}

	// Fetch the retirement-tax table only when the TaxSaver benefit is selected
	// on at least one scheme category — otherwise the rows aren't consulted.
	var taxRetirementBands []models.TaxRetirementTable
	for _, sc := range groupQuote.SchemeCategories {
		if sc.TaxSaverBenefit {
			taxRetirementBands, err = GetTaxRetirementTableByRiskRateCode(groupQuote.RiskRateCode)
			if err != nil {
				logger.WithField("error", err.Error()).Warn("Failed to load tax retirement table; TaxSaverSumAssured will be zero")
			}
			break
		}
	}

	// Fetch custom tiered income replacement tiers if any scheme category uses "custom" mode.
	var customTieredIncomeTiers []models.TieredIncomeReplacement
	needsCustomTiers := false
	for _, sc := range groupQuote.SchemeCategories {
		if (sc.PhiUseTieredIncomeReplacementRatio && sc.PhiTieredIncomeReplacementType == "custom") ||
			(sc.TtdUseTieredIncomeReplacementRatio && sc.TtdTieredIncomeReplacementType == "custom") {
			needsCustomTiers = true
			break
		}
	}
	if needsCustomTiers {
		customTieredIncomeTiers, err = GetCustomTieredIncomeReplacementTiers(groupQuote.SchemeName, groupQuote.RiskRateCode)
		if err != nil {
			return fmt.Errorf("failed to load custom tiered income replacement tiers: %v", err)
		}
		if len(customTieredIncomeTiers) == 0 {
			return fmt.Errorf("missing custom tiered income replacement table for scheme '%s' — this needs super admin attention", groupQuote.SchemeName)
		}
	}

	// Determine scheme size level once for this quote before spawning category workers.
	tempParam := models.GroupPricingParameters{RiskRateCode: groupQuote.RiskRateCode}
	schemeSizeRow := GetSchemeSizeLoading(tempParam, groupQuote.MemberDataCount)
	schemeSizeLevel := schemeSizeRow.SizeLevel
	logger.WithField("scheme_size_level", schemeSizeLevel).Debug("Determined scheme size level")

	totalCategories := len(groupQuote.SelectedSchemeCategories)
	var completedCategories int64

	// Progress callback invoked by each category worker at key checkpoints.
	emitProgress := func(category, phase string) {
		completed := int(atomic.LoadInt64(&completedCategories))
		// Each category has 3 phases; compute overall progress.
		phaseWeight := map[string]float64{"loading_data": 0.0, "rating_members": 0.33, "saving_results": 0.66, "category_done": 1.0}
		w := phaseWeight[phase]
		progress := (float64(completed) + w) / float64(totalCategories) * 100
		sendCalculationProgress(user.UserEmail, CalculationProgress{
			QuoteID:             quoteId,
			TotalCategories:     totalCategories,
			CompletedCategories: completed,
			CurrentCategory:     category,
			Phase:               phase,
			Progress:            math.Min(progress, 100),
		})
	}

	categoryWorkerPool := workerpool.New(runtime.NumCPU())
	for _, selectedSchemeCategory := range groupQuote.SelectedSchemeCategories {
		selectedSchemeCategory := selectedSchemeCategory // capture range variable
		categoryWorkerPool.Submit(func() {
			err2 := calculateForCategory(quoteId, basis, credibility, user, logger, groupQuote, dbStartTime, dbElapsed, selectedSchemeCategory, taxTable, taxRetirementBands, tieredIncomeTiers, customTieredIncomeTiers, schemeSizeLevel, schemeSizeRow, emitProgress)
			if err2 != nil {
				logger.WithField("error", err2.Error()).Error("Error calculating for scheme category")
			}
			atomic.AddInt64(&completedCategories, 1)
			emitProgress(selectedSchemeCategory, "category_done")
		})
	}

	categoryWorkerPool.StopWait()
	GroupPricingCache.Clear()

	// Scheme-wide commission pass — applied once all category workers have
	// written their summaries. Commission runs on the aggregate total premium
	// (scheme-level) via the tiered CommissionStructure bands, then is
	// distributed proportionally back to each category and benefit.
	if err := applySchemeWideCommission(intQuoteID(quoteId), groupQuote, logger); err != nil {
		logger.WithField("error", err.Error()).Error("Failed to apply scheme-wide commission")
	}

	// Stamp the calculation completion time on the quote so the UI can show
	// users when the results were last produced. Done before emitting the
	// completion event so a client that queries the quote on receipt sees the
	// updated timestamp.
	completedAt := time.Now()
	if err := DB.Model(&models.GroupPricingQuote{}).
		Where("id = ?", intQuoteID(quoteId)).
		Update("calculation_completed_at", completedAt).Error; err != nil {
		logger.WithField("error", err.Error()).Error("Failed to stamp calculation_completed_at on quote")
	}

	// Send final 100% completion event
	sendCalculationProgress(user.UserEmail, CalculationProgress{
		QuoteID:             quoteId,
		TotalCategories:     totalCategories,
		CompletedCategories: totalCategories,
		Phase:               "completed",
		Progress:            100,
	})

	elapsed := time.Since(startTime)
	logger.WithField("elapsed_ms", elapsed.Milliseconds()).Info("Group pricing quote calculation completed successfully")
	go func(id int) { ScoreQuote(id) }(intQuoteID(quoteId))
	return nil
}

func calculateForCategory(quoteId string, basis string, credibility float64, user models.AppUser, logger *logrus.Entry, groupQuote models.GroupPricingQuote, dbStartTime time.Time, dbElapsed time.Duration, selectedSchemeCategory string, taxTable []models.TaxTable, taxRetirementBands []models.TaxRetirementTable, tieredIncomeTiers []models.TieredIncomeReplacement, customTieredIncomeTiers []models.TieredIncomeReplacement, schemeSizeLevel int, schemeSizeRow models.SchemeSizeLevel, emitProgress func(category, phase string)) error {
	var memberMps []models.GPricingMemberData
	var indicativeMemberMps []models.MemberIndicativeDataSet
	var experienceMemberMps []models.GPricingMemberData
	var groupParameter models.GroupPricingParameters
	var groupPricingReinsuranceStructure models.GroupPricingReinsuranceStructure

	var memberDataResults []models.MemberRatingResult
	var mdrs models.MemberRatingResultSummary
	var historicalCredibilityData models.HistoricalCredibilityData
	var memberPremiumSchedule []models.MemberPremiumSchedule
	var bordereaux []models.Bordereaux
	var incomeLevels []models.IncomeLevel
	var ageBands []models.GroupPricingAgeBands
	var experienceClaimsData []models.GroupPricingClaimsExperience
	var experienceRateOverrides []models.GroupPricingExperienceRateOverride
	var experienceRateOverrideLookup map[string]map[string]models.GroupPricingExperienceRateOverride
	var educatorBenefitStructure models.EducatorBenefitStructure
	var annualGlaExperienceWeightedRate, credibilityRate float64
	var annualPtdExperienceWeightedRate, annualCiExperienceWeightedRate float64
	//var ptdTheoreticalRate, ciTheoreticalRate, ttdTheoreticalRate, phiTheoreticalRate float64
	var glaTheoreticalRate float64
	var err error
	var theoreticalRiskTotal TheoreticalRiskTotal
	var GroupPricingInsurerDetail models.GroupPricingInsurerDetail

	// populate the scheme category basis with the quote basis for now

	logger.Debug("Retrieving group parameters for basis")
	err = DB.Where("basis=? and risk_rate_code = ?", basis, groupQuote.RiskRateCode).Find(&groupParameter).Error

	if err != nil {
		logger.Errorf("Error retrieving group parameters: %v", err)
		return err
	}

	restriction, err := GetRestrictionByRiskRateCode(groupQuote.RiskRateCode)
	if err != nil {
		appLog.Error("Error retrieving restriction for risk rate code: ", groupQuote.RiskRateCode, " error: ", err.Error())
	}
	reinsCoverCaps := LoadReinsuranceCoverCaps(groupQuote.RiskRateCode, groupQuote.MemberDataCount)

	premiumLoading := GetPremiumLoading(groupParameter, schemeSizeLevel, string(groupQuote.DistributionChannel))
	logger.WithFields(map[string]interface{}{
		"scheme_size_level": schemeSizeLevel,
		"channel":           groupQuote.DistributionChannel,
		"found":             premiumLoading.ID != 0,
	}).Debug("Retrieved PremiumLoading for group parameter")

	if groupQuote.DistributionChannel == models.ChannelDirect {
		groupQuote.Loadings.CommissionLoading = 0
	} else {
		groupQuote.Loadings.CommissionLoading = premiumLoading.CommissionLoading * 100
	}
	groupQuote.Loadings.ExpenseLoading = premiumLoading.ExpenseLoading * 100
	groupQuote.Loadings.ProfitLoading = premiumLoading.ProfitLoading * 100
	groupQuote.Loadings.AdminLoading = premiumLoading.AdminLoading * 100
	groupQuote.Loadings.OtherLoading = premiumLoading.OtherLoading * 100

	logger.Debug("Retrieving insurer detail for basis")
	dbStartTime = time.Now()
	if err := DB.Limit(1).Find(&GroupPricingInsurerDetail).Error; err != nil {
		logger.WithField("error", err.Error()).Warn("Failed to load insurer detail")
	}
	dbElapsed = time.Since(dbStartTime)
	logger.WithField("elapsed_ms", dbElapsed.Milliseconds()).Debug("Retrieved insurer detail")
	insurerYearEndMonth := GroupPricingInsurerDetail.YearEndMonth

	// update the group quote with the basis
	groupQuote.Basis = basis

	// find the matching scheme category by name and update basis
	var category *models.SchemeCategory
	categoryIndex := -1
	for idx := range groupQuote.SchemeCategories {
		if groupQuote.SchemeCategories[idx].SchemeCategory == selectedSchemeCategory {
			category = &groupQuote.SchemeCategories[idx]
			categoryIndex = idx
			break
		}
	}

	groupQuote.OccupationClass = GetOccupationClass(groupParameter, groupQuote, selectedSchemeCategory)
	logger.WithField("occupation_class", groupQuote.OccupationClass).Debug("Updated group quote with basis and occupation class")

	// Pre-load region loadings for this category at the category level (single DB hit, keyed by gender).
	// Skipped when the table is configured as not required — downstream variables resolve to zero
	// via the existing graceful-degradation path.
	regionLoadingByGender := make(map[string]models.RegionLoading)
	if IsTableRequired("regionLoading") && category != nil && category.Region != "" {
		var categoryRegionLoadings []models.RegionLoading
		DB.Where("risk_rate_code = ? AND region = ?",
			groupParameter.RiskRateCode,
			strings.TrimSpace(category.Region),
		).Find(&categoryRegionLoadings)
		for _, r := range categoryRegionLoadings {
			if len(r.Gender) > 0 {
				regionLoadingByGender[strings.ToUpper(r.Gender[:1])] = r
			}
		}
	}

	// Pre-load reinsurance region loadings for this category (single DB hit, keyed by gender).
	// Skipped when the table is configured as not required.
	reinsRegionLoadingByGender := make(map[string]models.ReinsuranceRegionLoading)
	if IsTableRequired("reinsuranceRegionLoading") && category != nil && category.Region != "" {
		var categoryReinsRegionLoadings []models.ReinsuranceRegionLoading
		DB.Where("risk_rate_code = ? AND region = ?",
			groupParameter.RiskRateCode,
			strings.TrimSpace(category.Region),
		).Find(&categoryReinsRegionLoadings)
		for _, r := range categoryReinsRegionLoadings {
			if len(r.Gender) > 0 {
				reinsRegionLoadingByGender[strings.ToUpper(r.Gender[:1])] = r
			}
		}
	}

	// Pre-load industry loadings for this category (single DB hit, keyed by gender).
	// Skipped when the table is configured as not required.
	industryLoadingByGender := make(map[string]models.IndustryLoading)
	if IsTableRequired("industryLoading") {
		var categoryIndustryLoadings []models.IndustryLoading
		DB.Where("risk_rate_code = ? AND occupation_class = ?",
			groupParameter.RiskRateCode,
			groupQuote.OccupationClass,
		).Find(&categoryIndustryLoadings)
		for _, il := range categoryIndustryLoadings {
			if len(il.Gender) > 0 {
				industryLoadingByGender[strings.ToUpper(il.Gender[:1])] = il
			}
		}
	}

	dbStartTime = time.Now()
	DB.Save(&groupQuote)

	dbElapsed = time.Since(dbStartTime)
	logger.WithField("elapsed_ms", dbElapsed.Milliseconds()).Debug("Saved group quote")

	// Run independent DB queries in parallel using errgroup
	logger.Debug("Loading reference data in parallel")
	dbParallelStart := time.Now()
	var memberDataErr error
	eg := new(errgroup.Group)

	// Each parallel load is gated on its TableConfiguration.IsRequired flag.
	// When a table is configured as not required the slice stays empty and
	// downstream variables resolve to zero through the existing graceful
	// degradation path.
	if IsTableRequired("gpReinsuranceStructure") {
		eg.Go(func() error {
			return DB.Where("risk_rate_code=? and basis=?", groupParameter.RiskRateCode, basis).Find(&groupPricingReinsuranceStructure).Error
		})
	}

	if IsTableRequired("incomeLevel") {
		eg.Go(func() error {
			return DB.Where("risk_rate_code=?", groupParameter.RiskRateCode).Find(&incomeLevels).Error
		})
	}

	if IsTableRequired("ageBands") {
		eg.Go(func() error {
			return DB.Find(&ageBands).Error
		})
	}

	if IsTableRequired("gpEducatorBenefitStructure") {
		if code := educatorCodeForCategory(category); code != "" {
			eg.Go(func() error {
				return DB.Where("risk_rate_code=? and educator_benefit_code=?", groupParameter.RiskRateCode, code).Find(&educatorBenefitStructure).Error
			})
		}
	}

	if groupQuote.ExperienceRating == "Yes" {
		eg.Go(func() error {
			return DB.Where("quote_id=?", groupQuote.ID).Find(&experienceClaimsData).Error
		})
	}

	if groupQuote.ExperienceRating == "Override" {
		eg.Go(func() error {
			return DB.Where("quote_id=?", groupQuote.ID).Find(&experienceRateOverrides).Error
		})
	}

	// Member data loading (runs concurrently with the above)
	eg.Go(func() error {
		if groupQuote.MemberIndicativeData {
			var membermp models.GPricingMemberData
			if err := DB.Where("quote_id = ? and scheme_category = ?", groupQuote.ID, selectedSchemeCategory).Find(&indicativeMemberMps).Error; err != nil {
				memberDataErr = err
				return nil
			}
			membermp.AnnualSalary = indicativeMemberMps[0].MemberAverageIncome
			membermp.SchemeName = groupQuote.SchemeName
			membermp.SchemeId = groupQuote.SchemeID
			membermp.QuoteId = groupQuote.ID
			membermp.Gender = "M"
			memberMps = append(memberMps, membermp)
		} else {
			if groupQuote.QuoteType == "Renewal" {
				memberDataErr = DB.Model(&models.GPricingMemberDataInForce{}).Where("scheme_id = ? and scheme_category=?", groupQuote.SchemeID, selectedSchemeCategory).Scan(&memberMps).Error
			} else {
				memberDataErr = DB.Where("quote_id = ? and scheme_category=?", groupQuote.ID, selectedSchemeCategory).Find(&memberMps).Error
			}

			if groupQuote.ExperienceRating == "Yes" && len(memberMps) > 0 {
				if groupQuote.QuoteType == "Renewal" {
					DB.Model(&models.GPricingMemberDataInForce{}).Where("scheme_id = ?", groupQuote.SchemeID).Scan(&experienceMemberMps)
				} else {
					DB.Where("quote_id = ?", groupQuote.ID).Find(&experienceMemberMps)
				}
			}
		}
		return nil
	})

	if egErr := eg.Wait(); egErr != nil {
		logger.WithField("error", egErr.Error()).Error("Failed to load reference data")
		return egErr
	}
	logger.WithField("elapsed_ms", time.Since(dbParallelStart).Milliseconds()).Debug("Loaded all reference data in parallel")
	emitProgress(selectedSchemeCategory, "loading_data")

	if groupQuote.ExperienceRating == "Override" {
		experienceRateOverrideLookup = buildExperienceRateOverrideLookup(experienceRateOverrides)
	}

	if memberDataErr != nil {
		logger.WithField("error", memberDataErr.Error()).Error("Failed to retrieve member data")
	}
	logger.WithField("member_count", len(memberMps)).Debug("Retrieved member data")

	if !groupQuote.MemberIndicativeData && len(memberMps) == 0 {
		return nil
	}

	// Extended family funeral benefit: compute per-age-band straight-averaged
	// loaded rates and attach to the scheme category so the UI can render them.
	// Uses a representative income level (first member) as the loading driver.
	if category != nil && category.ExtendedFamilyBenefit {
		var bands []models.ExtendedFamilyAgeBand
		if category.ExtendedFamilyAgeBandSource == "custom" {
			bands = category.ExtendedFamilyCustomAgeBands
		} else {
			bands = make([]models.ExtendedFamilyAgeBand, 0, len(ageBands))
			for _, b := range ageBands {
				// Filter by the selected age-band type. Empty type matches
				// untyped rows for backward compatibility.
				if category.ExtendedFamilyAgeBandType != "" &&
					b.Type != category.ExtendedFamilyAgeBandType {
					continue
				}
				bands = append(bands, models.ExtendedFamilyAgeBand{MinAge: b.MinAge, MaxAge: b.MaxAge})
			}
		}
		if len(bands) > 0 {
			repIncomeLevel := 0
			if len(memberMps) > 0 {
				repIncomeLevel = GetIncomeLevel(memberMps[0], incomeLevels)
			}
			var sums []models.ExtendedFamilyBandSumAssured
			if category.ExtendedFamilyPricingMethod == "sum_assured" {
				sums = category.ExtendedFamilySumsAssured
			}
			// Combined premium loading from the premium_loading table, mirroring
			// the formula used for every other benefit's office premium (see the
			// per-member rating at memberDataPointResult.TotalLoading below).
			// Direct-channel schemes zero out commission exactly as they do
			// there.
			effectiveCommission := premiumLoading.CommissionLoading
			if groupQuote.DistributionChannel == models.ChannelDirect {
				effectiveCommission = 0
			}
			efTotalLoading := math.Max(
				premiumLoading.ExpenseLoading+
					premiumLoading.AdminLoading+
					effectiveCommission+
					premiumLoading.ProfitLoading+
					premiumLoading.OtherLoading,
				premiumLoading.MinimumPremiumLoading,
			)
			bandRates, efErr := CalculateExtendedFamilyAgeBandRates(
				groupParameter.RiskRateCode,
				repIncomeLevel,
				groupParameter.ExtendedFamilyMaleProp,
				efTotalLoading,
				bands,
				sums,
			)
			if efErr != nil {
				logger.WithField("error", efErr.Error()).Warn("Failed to compute extended family band rates")
			} else {
				// For rate_per_1000 method the CalculateExtendedFamilyAgeBandRates
				// helper does not populate MonthlyPremium (no sum assured), so
				// derive it here: rate * 1000 / 12 per extended-family member
				// on both the risk rate (AverageRate) and the office rate.
				if category.ExtendedFamilyPricingMethod != "sum_assured" {
					for i := range bandRates {
						bandRates[i].MonthlyPremium = bandRates[i].AverageRate * 1000.0 / 12.0
						bandRates[i].OfficeMonthlyPremium = bandRates[i].OfficeRate * 1000.0 / 12.0
					}
				} else {
					// Divide-by-12 already applied; ensure zero premium rows
					// (bands with no sum assured provided) remain zero.
					for i := range bandRates {
						if bandRates[i].SumAssured == 0 {
							bandRates[i].MonthlyPremium = 0
							bandRates[i].OfficeMonthlyPremium = 0
						}
					}
				}
				category.ExtendedFamilyBandRates = bandRates
				if err := DB.Model(&models.SchemeCategory{}).
					Where("id = ?", category.ID).
					Update("extended_family_band_rates", bandRates).Error; err != nil {
					logger.WithField("error", err.Error()).Warn("Failed to persist extended_family_band_rates")
				}
				// Mirror into the preloaded slice so downstream code sees it.
				if categoryIndex >= 0 {
					groupQuote.SchemeCategories[categoryIndex].ExtendedFamilyBandRates = bandRates
				}
				// Mirror the config + computed rates onto the member rating
				// result summary so the Premiums Summary screen can render the
				// per-category extended-family section from a single payload.
				mdrs.ExtendedFamilyBenefit = true
				mdrs.ExtendedFamilyAgeBandSource = category.ExtendedFamilyAgeBandSource
				mdrs.ExtendedFamilyAgeBandType = category.ExtendedFamilyAgeBandType
				mdrs.ExtendedFamilyPricingMethod = category.ExtendedFamilyPricingMethod
				mdrs.ExtendedFamilyBandRates = bandRates
				var totalMonthly float64
				for _, b := range bandRates {
					totalMonthly += b.MonthlyPremium
				}
				mdrs.TotalExtendedFamilyMonthlyPremium = totalMonthly
			}
		}
	}

	// Additional GLA Cover: compute per-age-band rate per 1,000 using the
	// main GLA benefit type + waiting period, the gender split from the
	// uploaded member data for this category, and the same region /
	// industry / contingency / premium loading stack as base GLA. This is
	// a rate-only product — no per-member premium aggregation.
	if category != nil && category.AdditionalGlaCoverBenefit && category.GlaBenefitType != "" {
		var aglaBands []models.AdditionalGlaCoverAgeBand
		if category.AdditionalGlaCoverAgeBandSource == "custom" {
			aglaBands = category.AdditionalGlaCoverCustomAgeBands
		} else {
			for _, b := range ageBands {
				if category.AdditionalGlaCoverAgeBandType != "" &&
					b.Type != category.AdditionalGlaCoverAgeBandType {
					continue
				}
				aglaBands = append(aglaBands, models.AdditionalGlaCoverAgeBand{MinAge: b.MinAge, MaxAge: b.MaxAge})
			}
		}
		if len(aglaBands) > 0 {
			repIncomeLevelAgla := 0
			if len(memberMps) > 0 {
				repIncomeLevelAgla = GetIncomeLevel(memberMps[0], incomeLevels)
			}
			// Blend the gender-weighted loadings using the uploaded-member
			// male proportion. Falls back to
			// group_pricing_parameters.main_member_male_prop when the member
			// list has no gender rows.
			var maleCount, totalGender int
			for _, m := range memberMps {
				g := strings.ToUpper(strings.TrimSpace(m.Gender))
				if g == "" {
					continue
				}
				totalGender++
				if g == "M" || g == "MALE" {
					maleCount++
				}
			}
			maleProp := groupParameter.MainMemberMaleProp
			if totalGender > 0 {
				maleProp = float64(maleCount) / float64(totalGender)
			}

			effectiveCommissionAgla := premiumLoading.CommissionLoading
			if groupQuote.DistributionChannel == models.ChannelDirect {
				effectiveCommissionAgla = 0
			}
			aglaBinderRate, aglaOutsourceRate := binderAndOutsourceRates(&groupQuote)
			aglaOtherLoadings := premiumLoading.ExpenseLoading +
				premiumLoading.AdminLoading +
				premiumLoading.ProfitLoading +
				premiumLoading.OtherLoading
			// Snapshot the pre-existing smoothed values (and their factors)
			// keyed by (min_age, max_age) so they survive a recalc when the
			// same band still exists. New/resized bands start unsmoothed.
			type aglaSmoothCarry struct {
				smoothed, smoothedM, smoothedF *float64
				factor, factorM, factorF       *float64
			}
			priorSmoothed := make(map[[2]int]aglaSmoothCarry, len(category.AdditionalGlaCoverBandRates))
			for _, prev := range category.AdditionalGlaCoverBandRates {
				priorSmoothed[[2]int{prev.MinAge, prev.MaxAge}] = aglaSmoothCarry{
					smoothed:  prev.SmoothedOfficeRatePer1000,
					smoothedM: prev.SmoothedOfficeRatePer1000Male,
					smoothedF: prev.SmoothedOfficeRatePer1000Female,
					factor:    prev.SmoothingFactor,
					factorM:   prev.SmoothingFactorMale,
					factorF:   prev.SmoothingFactorFemale,
				}
			}

			aglaBandRates, aglaErr := CalculateAdditionalGlaCoverBandRates(
				groupParameter.RiskRateCode,
				category.GlaBenefitType,
				category.Region,
				category.GlaWaitingPeriod,
				repIncomeLevelAgla,
				groupQuote.OccupationClass,
				maleProp,
				aglaBands,
				effectiveCommissionAgla,
				aglaBinderRate,
				aglaOutsourceRate,
				aglaOtherLoadings,
				premiumLoading.MinimumPremiumLoading,
				memberMps,
				groupQuote.CommencementDate,
				category.GlaSalaryMultiple,
			)
			if aglaErr != nil {
				logger.WithFields(map[string]interface{}{
					"error":        aglaErr.Error(),
					"category":     selectedSchemeCategory,
					"benefit_type": category.GlaBenefitType,
					"band_count":   len(aglaBands),
					"male_prop":    maleProp,
				}).Warn("Failed to compute additional GLA cover band rates")
			} else {
				for i := range aglaBandRates {
					if carry, ok := priorSmoothed[[2]int{aglaBandRates[i].MinAge, aglaBandRates[i].MaxAge}]; ok {
						aglaBandRates[i].SmoothedOfficeRatePer1000 = carry.smoothed
						aglaBandRates[i].SmoothedOfficeRatePer1000Male = carry.smoothedM
						aglaBandRates[i].SmoothedOfficeRatePer1000Female = carry.smoothedF
						aglaBandRates[i].SmoothingFactor = carry.factor
						aglaBandRates[i].SmoothingFactorMale = carry.factorM
						aglaBandRates[i].SmoothingFactorFemale = carry.factorF
					}
				}
				category.AdditionalGlaCoverBandRates = aglaBandRates
				mpUsed := maleProp
				category.AdditionalGlaCoverMalePropUsed = &mpUsed
				if err := DB.Model(&models.SchemeCategory{}).
					Where("id = ?", category.ID).
					Updates(map[string]interface{}{
						"additional_gla_cover_band_rates":     aglaBandRates,
						"additional_gla_cover_male_prop_used": mpUsed,
					}).Error; err != nil {
					logger.WithField("error", err.Error()).Warn("Failed to persist additional_gla_cover_band_rates")
				}
				if categoryIndex >= 0 {
					groupQuote.SchemeCategories[categoryIndex].AdditionalGlaCoverBandRates = aglaBandRates
					groupQuote.SchemeCategories[categoryIndex].AdditionalGlaCoverMalePropUsed = &mpUsed
				}
				// Mirror onto the MRRS row being built so the Premium
				// Summary can read everything from a single payload
				// (matches the Extended Family Funeral mirror above).
				mdrs.AdditionalGlaCoverBenefit = true
				mdrs.AdditionalGlaCoverAgeBandSource = category.AdditionalGlaCoverAgeBandSource
				mdrs.AdditionalGlaCoverAgeBandType = category.AdditionalGlaCoverAgeBandType
				mdrs.AdditionalGlaCoverBandRates = aglaBandRates
				mdrs.AdditionalGlaCoverMalePropUsed = &mpUsed
				logger.WithFields(map[string]interface{}{
					"category":     selectedSchemeCategory,
					"band_count":   len(aglaBandRates),
					"male_prop":    maleProp,
					"benefit_type": category.GlaBenefitType,
				}).Info("Additional GLA Cover band rates computed and mirrored to MRRS")
			}
		} else {
			// Benefit is enabled but we couldn't build a band list — most
			// often because source=standard with no age_band_type selected,
			// or a custom list that's empty. Log it so the Premium Summary
			// silence makes sense.
			logger.WithFields(map[string]interface{}{
				"category":        selectedSchemeCategory,
				"source":          category.AdditionalGlaCoverAgeBandSource,
				"band_type":       category.AdditionalGlaCoverAgeBandType,
				"custom_count":    len(category.AdditionalGlaCoverCustomAgeBands),
				"available_bands": len(ageBands),
			}).Warn("Additional GLA Cover enabled but band list resolved to empty")
		}
	} else if category != nil && category.AdditionalGlaCoverBenefit {
		// Enabled but no GlaBenefitType — usually means base GLA wasn't
		// fully configured.
		logger.WithField("category", selectedSchemeCategory).
			Warn("Additional GLA Cover enabled but scheme category has no gla_benefit_type; skipping calc")
	}

	var weightedLifeYears float64

	if groupQuote.ExperienceRating == "Yes" {

		var weightedPeriod, experienceRate, annualExperienceRate float64
		var ptdExperienceRate, ciExperienceRate, ptdAnnualExperienceRate, ciAnnualExperienceRate float64
		for _, claimsDataPoint := range experienceClaimsData {

			start, err := time.Parse("2006/01/02", claimsDataPoint.StartDate)
			if err != nil {
				fmt.Println("Error parsing start date:", err)

			}
			end, err := time.Parse("2006/01/02", claimsDataPoint.EndDate)
			if err != nil {
				fmt.Println("Error parsing start date:", err)
			}
			duration := end.Sub(start)
			timePeriodYears := (duration.Hours() / 24) / 365.25
			weightedPeriod += timePeriodYears * claimsDataPoint.Weighting
			weightedLifeYears += float64(claimsDataPoint.NumberOfMembers) * timePeriodYears * claimsDataPoint.Weighting
			if claimsDataPoint.TotalGlaSumAssured > 0 {
				experienceRate = claimsDataPoint.GlaClaimsAmount / (claimsDataPoint.TotalGlaSumAssured)
			}
			if claimsDataPoint.TotalPtdSumAssured > 0 {
				ptdExperienceRate = claimsDataPoint.PtdClaimsAmount / claimsDataPoint.TotalPtdSumAssured
			}

			if claimsDataPoint.TotalCiSumAssured > 0 {
				ciExperienceRate = claimsDataPoint.CiClaimsAmount / claimsDataPoint.TotalCiSumAssured
			}

			if timePeriodYears > 0 && claimsDataPoint.Weighting > 0 {
				annualExperienceRate = experienceRate / timePeriodYears
				ptdAnnualExperienceRate = ptdExperienceRate / timePeriodYears
				ciAnnualExperienceRate = ciExperienceRate / timePeriodYears
			}
			annualGlaExperienceWeightedRate += annualExperienceRate * claimsDataPoint.Weighting
			annualPtdExperienceWeightedRate += ptdAnnualExperienceRate * claimsDataPoint.Weighting
			annualCiExperienceWeightedRate += ciAnnualExperienceRate * claimsDataPoint.Weighting

			historicalCredibilityData.ClaimCount += claimsDataPoint.NumberOfGlaClaims
			historicalCredibilityData.DurationInForce += utils.RoundUp(timePeriodYears)
			historicalCredibilityData.ExperiencePeriod += utils.RoundUp(timePeriodYears)
		}

		if groupParameter.FullCredibilityThreshold > 0 {
			credibilityRate = math.Min(math.Sqrt(weightedLifeYears/groupParameter.FullCredibilityThreshold), 1)
		}
		if groupParameter.GlobalGlaExperienceRate > 0 {
			annualGlaExperienceWeightedRate = groupParameter.GlobalGlaExperienceRate
			annualPtdExperienceWeightedRate = groupParameter.GlobalPtdExperienceRate
			annualCiExperienceWeightedRate = groupParameter.GlobalCiExperienceRate
		}
	}

	mdrs.QuoteId = groupQuote.ID
	mdrs.SchemeId = groupQuote.SchemeID

	//Calculating Distribution Based Free Cover Limit
	var Query string
	if groupQuote.UseGlobalSalaryMultiple {
		Query = fmt.Sprintf("select annual_salary * %v FROM g_pricing_member_data as annual_salary where quote_id = '%d' order by annual_salary", category.GlaSalaryMultiple, groupQuote.ID)
	}
	if !groupQuote.UseGlobalSalaryMultiple {
		if groupQuote.QuoteType == "Renewal" {
			// This query prioritizes g_pricing_member_data.
			// If it has rows for this quote_id, those are returned.
			// Otherwise, it pulls from g_pricing_member_data_in_force.
			Query = fmt.Sprintf(`
    SELECT annual_salary * benefits_gla_multiple
    FROM (
        SELECT annual_salary, benefits_gla_multiple, 1 as priority
        FROM g_pricing_member_data
        WHERE quote_id = '%d'

        UNION ALL

        SELECT annual_salary, benefits_gla_multiple, 2 as priority
        FROM g_pricing_member_data_in_forces
        WHERE scheme_id = '%d'
        AND NOT EXISTS (SELECT 1 FROM g_pricing_member_data WHERE quote_id = '%d')
    ) as combined_data
    ORDER BY annual_salary`,
				groupQuote.ID, groupQuote.SchemeID, groupQuote.ID,
			)
		}
		if groupQuote.QuoteType == "New Business" {
			Query = fmt.Sprintf("select annual_salary * benefits_gla_multiple FROM g_pricing_member_data as annual_salary where quote_id = '%d' order by annual_salary", groupQuote.ID)
		}
	}

	var salaryData []float64
	var nthPercentileSalary float64
	var meanSalary, calculatedFreeCoverLimit float64
	var memberCount int
	var indicativeRatesCount float64
	logger.Debug("Executing raw SQL query for salary data")
	dbStartTime = time.Now()
	err = DB.Raw(Query).Scan(&salaryData).Error
	dbElapsed = time.Since(dbStartTime)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"error":      err.Error(),
			"elapsed_ms": dbElapsed.Milliseconds(),
		}).Error("Failed to execute raw SQL query for salary data")
	} else {
		logger.WithField("elapsed_ms", dbElapsed.Milliseconds()).Debug("Executed raw SQL query for salary data")
	}
	//median, _ = stats.Median(salaryData)
	if groupParameter.FreeCoverLimitPercentile == 0 {
		nthPercentileSalary = 0
	} else {
		nthPercentileSalary, _ = stats.Percentile(salaryData, groupParameter.FreeCoverLimitPercentile*100)
	}
	if !groupQuote.MemberIndicativeData {
		meanSalary, _ = stats.Mean(salaryData)
		//maxSalary, _ = stats.Max(salaryData)
		memberCount = len(salaryData)
		indicativeRatesCount = 1
	}

	if groupQuote.MemberIndicativeData {
		meanSalary = indicativeMemberMps[0].MemberAverageIncome
		//maxSalary, _ = stats.Max(salaryData)
		memberCount = indicativeMemberMps[0].MemberDataCount
		indicativeRatesCount = float64(memberCount)
	}

	if groupQuote.MemberIndicativeData {
		calculatedFreeCoverLimit = groupQuote.FreeCoverLimit
	}

	var scalingTerm, statisticalOutlierThreshold, maximumSumAssured float64
	maxAllowedFCL := restriction.MaximumAllowedFCL
	fclOverrideTolerance := GetFCLOverrideTolerance()

	if !groupQuote.MemberIndicativeData {
		scalingTerm = groupParameter.FreeCoverLimitScalingFactor * math.Sqrt(float64(memberCount)) * meanSalary

		var rawFCL float64
		switch GetFCLMethod() {
		case models.FCLMethodOutlier:
			// Log-normal +3σ upper bound on sum assured, plus a max-SA-scaled cap.
			// Skip non-positive entries before LN.
			lnSalaryData := make([]float64, 0, len(salaryData))
			for _, v := range salaryData {
				if v > 0 {
					lnSalaryData = append(lnSalaryData, math.Log(v))
				}
			}
			lnMean, _ := stats.Mean(lnSalaryData)
			lnStdevP, _ := stats.StandardDeviationPopulation(lnSalaryData)
			statisticalOutlierThreshold = math.Exp(lnMean + 3*lnStdevP)
			maximumSumAssured, _ = stats.Max(salaryData)

			rawFCL = math.Min(scalingTerm, statisticalOutlierThreshold)
			if groupParameter.FCLMaximumCoverScalingFactor > 0 {
				rawFCL = math.Min(rawFCL, maximumSumAssured*groupParameter.FCLMaximumCoverScalingFactor)
			}
		default:
			// Percentile method (existing behaviour, default).
			rawFCL = math.Min(scalingTerm, nthPercentileSalary)
		}

		if groupParameter.FreeCoverLimitNearestMultiple > 0 {
			calculatedFreeCoverLimit = math.Ceil(rawFCL/groupParameter.FreeCoverLimitNearestMultiple) * groupParameter.FreeCoverLimitNearestMultiple
		} else {
			calculatedFreeCoverLimit = rawFCL
		}
	}

	if category != nil {
		category.Basis = basis
		if groupQuote.FreeCoverLimit > 0 {
			switch GetFCLMethod() {
			case models.FCLMethodOutlier:
				maxthreshold1 := math.Max(statisticalOutlierThreshold, maximumSumAssured*groupParameter.FCLMaximumCoverScalingFactor)
				maxthrehold2 := math.Max(scalingTerm, maxthreshold1)

				if groupQuote.FreeCoverLimit > (1+fclOverrideTolerance)*maxthrehold2 {
					category.FreeCoverLimit = math.Min(maxthrehold2, maxAllowedFCL)
				} else {
					category.FreeCoverLimit = math.Min(groupQuote.FreeCoverLimit, maxAllowedFCL)
				}
			default:
				category.FreeCoverLimit = math.Min(groupQuote.FreeCoverLimit, maxAllowedFCL)
			}
		} else {
			category.FreeCoverLimit = math.Min(calculatedFreeCoverLimit, maxAllowedFCL)
		}
		DB.Save(category)
	} else {
		logger.WithField("selected_scheme_category", selectedSchemeCategory).Warn("Selected scheme category not found in SchemeCategories")
	}

	//where experience rating is chosen — collect-then-reduce pattern (no shared-state mutation)

	if groupQuote.ExperienceRating == "Yes" {
		expResults := make([]TheoreticalRiskTotal, len(experienceMemberMps))
		experienceWorkerPool := workerpool.New(min(runtime.NumCPU(), 8))
		for idx, mp := range experienceMemberMps {
			idx, mp := idx, mp
			experienceWorkerPool.Submit(func() {
				expResults[idx] = PopulateRatesPerMemberForExperienceRating(categoryIndex, indicativeRatesCount, mp, groupQuote, groupParameter, incomeLevels, ageBands, calculatedFreeCoverLimit, insurerYearEndMonth, restriction, reinsCoverCaps, taxTable, tieredIncomeTiers, customTieredIncomeTiers, premiumLoading, regionLoadingByGender, industryLoadingByGender)
			})
		}
		experienceWorkerPool.StopWait()

		// Single-threaded reduce — no mutex needed
		for _, r := range expResults {
			theoreticalRiskTotal.GlaSumRiskPremium += r.GlaSumRiskPremium
			theoreticalRiskTotal.PtdSumRiskPremium += r.PtdSumRiskPremium
			theoreticalRiskTotal.TtdSumRiskPremium += r.TtdSumRiskPremium
			theoreticalRiskTotal.PhiSumRiskPremium += r.PhiSumRiskPremium
			theoreticalRiskTotal.CiSumRiskPremium += r.CiSumRiskPremium
			theoreticalRiskTotal.GlaSumAssured += r.GlaSumAssured
			theoreticalRiskTotal.PtdSumAssured += r.PtdSumAssured
			theoreticalRiskTotal.TtdSumAssured += r.TtdSumAssured
			theoreticalRiskTotal.PhiSumAssured += r.PhiSumAssured
			theoreticalRiskTotal.CiSumAssured += r.CiSumAssured
		}
		if theoreticalRiskTotal.GlaSumAssured > 0 {
			glaTheoreticalRate = theoreticalRiskTotal.GlaSumRiskPremium * 1000 / theoreticalRiskTotal.GlaSumAssured
		}
	}

	// Collect-then-reduce: each goroutine writes to its own pre-allocated slot — no mutex needed
	rateResults := make([]MemberRateResult, len(memberMps))
	rateWorkerPool := workerpool.New(min(runtime.NumCPU(), 8))
	categoryExperienceOverrides := experienceRateOverrideLookup[selectedSchemeCategory]
	for idx, mp := range memberMps {
		idx, mp := idx, mp
		rateWorkerPool.Submit(func() {
			rateResults[idx] = PopulateRatesPerMember(categoryIndex, indicativeRatesCount, indicativeMemberMps, selectedSchemeCategory, mp, groupQuote, groupParameter, groupPricingReinsuranceStructure, incomeLevels, ageBands, credibilityRate, annualGlaExperienceWeightedRate, annualPtdExperienceWeightedRate, annualCiExperienceWeightedRate, calculatedFreeCoverLimit, educatorBenefitStructure, credibility, glaTheoreticalRate, insurerYearEndMonth, user, restriction, reinsCoverCaps, taxTable, taxRetirementBands, tieredIncomeTiers, customTieredIncomeTiers, premiumLoading, regionLoadingByGender, industryLoadingByGender, reinsRegionLoadingByGender, categoryExperienceOverrides, schemeSizeRow)
		})
	}
	rateWorkerPool.StopWait()
	emitProgress(selectedSchemeCategory, "rating_members")

	// Single-threaded reduce: aggregate results into slices and summary
	memberDataResults = make([]models.MemberRatingResult, 0, len(rateResults))
	memberPremiumSchedule = make([]models.MemberPremiumSchedule, 0, len(rateResults))
	bordereaux = make([]models.Bordereaux, 0, len(rateResults))

	for _, r := range rateResults {
		memberDataResults = append(memberDataResults, r.Rating)
		memberPremiumSchedule = append(memberPremiumSchedule, r.PremiumSchedule)
		bordereaux = append(bordereaux, r.Bordereaux)

		mr := r.Rating
		mdrs.Category = selectedSchemeCategory
		if category != nil {
			mdrs.TaxSaverBenefit = category.TaxSaverBenefit
		}
		mdrs.FinancialYear = mr.FinancialYear
		if groupQuote.MemberIndicativeData {
			mdrs.MemberCount = float64(indicativeMemberMps[0].MemberDataCount)
		} else {
			mdrs.MemberCount++
		}
		mdrs.ExceedsFreeCoverLimitIndicator += mr.ExceedsFreeCoverLimitIndicator
		mdrs.ExceedsNormalRetirementAgeIndicator += mr.ExceedsNormalRetirementAgeIndicator
		mdrs.TotalGlaRiskRate += mr.LoadedGlaRate
		mdrs.ExpTotalGlaRiskRate += mr.ExpAdjLoadedGlaRate
		mdrs.TotalGlaAnnualRiskPremium += mr.GlaRiskPremium
		mdrs.ExpTotalGlaAnnualRiskPremium += mr.ExpAdjGlaRiskPremium

		// TaxSaver is a slice of the GLA premium (already baked into
		// TotalGlaAnnualOfficePremium via LoadedGlaRate) — these rollups
		// exist only for reporting the slice, not for adding to any total.
		mdrs.TotalTaxSaverAnnualRiskPremium += mr.TaxSaverRiskPremium
		mdrs.ExpTotalTaxSaverAnnualRiskPremium += mr.ExpAdjTaxSaverRiskPremium

		mdrs.TotalAdditionalAccidentalGlaRiskRate += mr.LoadedAdditionalAccidentalGlaRate
		mdrs.ExpTotalAdditionalAccidentalGlaRiskRate += mr.ExpAdjLoadedAdditionalAccidentalGlaRate
		mdrs.TotalAdditionalAccidentalGlaAnnualRiskPremium += mr.AdditionalAccidentalGlaRiskPremium
		mdrs.ExpTotalAdditionalAccidentalGlaAnnualRiskPremium += mr.ExpAdjAdditionalAccidentalGlaRiskPremium

		mdrs.TotalPtdRiskRate += mr.LoadedPtdRate
		mdrs.ExpTotalPtdRiskRate += mr.ExpAdjLoadedPtdRate
		mdrs.TotalPtdAnnualRiskPremium += mr.PtdRiskPremium
		mdrs.ExpTotalPtdAnnualRiskPremium += mr.ExpAdjPtdRiskPremium

		mdrs.TotalTtdRiskRate += mr.LoadedTtdRate
		mdrs.ExpTotalTtdRiskRate += mr.ExpAdjLoadedTtdRate
		mdrs.TotalTtdAnnualRiskPremium += mr.TtdRiskPremium
		mdrs.ExpTotalTtdAnnualRiskPremium += mr.ExpAdjTtdRiskPremium

		mdrs.TotalPhiRiskRate += mr.LoadedPhiRate
		mdrs.ExpTotalPhiRiskRate += mr.ExpAdjLoadedPhiRate
		mdrs.TotalPhiAnnualRiskPremium += mr.PhiRiskPremium
		mdrs.ExpTotalPhiAnnualRiskPremium += mr.ExpAdjPhiRiskPremium

		mdrs.TotalCiRiskRate += mr.LoadedCiRate
		mdrs.ExpTotalCiRiskRate += mr.ExpAdjLoadedCiRate
		mdrs.TotalCiAnnualRiskPremium += mr.CiRiskPremium
		mdrs.ExpTotalCiAnnualRiskPremium += mr.ExpAdjCiRiskPremium

		mdrs.TotalSglaRiskRate += mr.LoadedSpouseGlaRate
		mdrs.ExpTotalSglaRiskRate += mr.ExpAdjLoadedSpouseGlaRate
		mdrs.TotalSglaAnnualRiskPremium += mr.SpouseGlaRiskPremium
		mdrs.ExpTotalSglaAnnualRiskPremium += mr.ExpAdjSpouseGlaRiskPremium

		mdrs.TotalFunAnnualRiskPremium += mr.TotalFuneralRiskPremium
		mdrs.ExpTotalFunAnnualRiskPremium += mr.ExpAdjTotalFuneralRiskPremium

		// Reinsurance premium sums per benefit. Funeral is the roll-up of
		// the five relationship-level reinsurance premiums.
		mdrs.TotalGlaReinsurancePremium += mr.GlaReinsurancePremium
		mdrs.TotalPtdReinsurancePremium += mr.PtdReinsurancePremium
		mdrs.TotalCiReinsurancePremium += mr.CiReinsurancePremium
		mdrs.TotalSglaReinsurancePremium += mr.SpouseGlaReinsurancePremium
		mdrs.TotalPhiReinsurancePremium += mr.PhiReinsurancePremium
		mdrs.TotalTtdReinsurancePremium += mr.TtdReinsurancePremium
		mdrs.TotalFunReinsurancePremium += mr.MainMemberReinsurancePremium + mr.SpouseReinsurancePremium + mr.ChildReinsurancePremium + mr.ParentReinsurancePremium + mr.DependantReinsurancePremium

		// Ceded sum-assured aggregates. The bordereaux datapoint carries the
		// per-member ceded values computed by GroupPricingReinsurance, so we
		// sum across all members in the category here.
		mdrs.TotalGlaCededSumAssured += r.Bordereaux.GlaCededSumAssured
		mdrs.TotalPtdCededSumAssured += r.Bordereaux.PtdCededSumAssured
		mdrs.TotalCiCededSumAssured += r.Bordereaux.CiCededSumAssured
		mdrs.TotalSglaCededSumAssured += r.Bordereaux.SglaCededSumAssured
		mdrs.TotalTtdCededMonthlyBenefit += r.Bordereaux.TtdCededMonthlyBenefit
		mdrs.TotalPhiCededMonthlyBenefit += r.Bordereaux.PhiCededMonthlyBenefit
		mdrs.TotalFunCededSumAssured += r.Bordereaux.MainMemberCededSumAssured + r.Bordereaux.SpouseCededSumAssured + r.Bordereaux.ChildCededSumAssured + r.Bordereaux.ParentCededSumAssured + r.Bordereaux.DependantCededSumAssured

		mdrs.TotalGlaSumAssured += mr.GlaSumAssured
		mdrs.TotalGlaCappedSumAssured += mr.GlaCappedSumAssured
		mdrs.TotalAdditionalAccidentalGlaSumAssured += mr.AdditionalAccidentalGlaSumAssured
		mdrs.TotalAdditionalAccidentalGlaCappedSumAssured += mr.AdditionalAccidentalGlaCappedSumAssured
		mdrs.TotalPtdSumAssured += mr.PtdSumAssured
		mdrs.TotalPtdCappedSumAssured += mr.PtdCappedSumAssured
		mdrs.TotalCiSumAssured += mr.CiSumAssured
		mdrs.TotalCiCappedSumAssured += mr.CiCappedSumAssured
		mdrs.TotalSglaSumAssured += mr.SpouseGlaSumAssured
		mdrs.TotalSglaCappedSumAssured += mr.SpouseGlaCappedSumAssured
		mdrs.TotalTtdIncome += mr.TtdIncome
		mdrs.TotalTtdCappedIncome += mr.TtdCappedIncome
		mdrs.TotalPhiIncome += mr.PhiIncome
		mdrs.TotalPhiCappedIncome += mr.PhiCappedIncome
		mdrs.TotalAnnualSalary += mr.AnnualSalary

		// Per-benefit annual-salary totals exclude members whose age has
		// passed the benefit's max cover age (their covered SA / income is
		// 0, so their salary must not inflate the proportion-of-salary
		// denominator). SGLA + AAGla key off GlaMaxCoverAge (no separate
		// column); MaxCoverAge of 0 means "no limit".
		mdrs.TotalAnnualSalaryGla += applyCoverAgeLimit(mr.AnnualSalary, mr.AgeNextBirthday, restriction.GlaMaxCoverAge)
		mdrs.TotalAnnualSalaryPtd += applyCoverAgeLimit(mr.AnnualSalary, mr.AgeNextBirthday, restriction.PtdMaxCoverAge)
		mdrs.TotalAnnualSalaryCi += applyCoverAgeLimit(mr.AnnualSalary, mr.AgeNextBirthday, restriction.CiMaxCoverAge)
		mdrs.TotalAnnualSalarySgla += applyCoverAgeLimit(mr.AnnualSalary, mr.AgeNextBirthday, restriction.GlaMaxCoverAge)
		mdrs.TotalAnnualSalaryTtd += applyCoverAgeLimit(mr.AnnualSalary, mr.AgeNextBirthday, restriction.TtdMaxCoverAge)
		mdrs.TotalAnnualSalaryPhi += applyCoverAgeLimit(mr.AnnualSalary, mr.AgeNextBirthday, restriction.PhiMaxCoverAge)
		mdrs.TotalAnnualSalaryFun += applyCoverAgeLimit(mr.AnnualSalary, mr.AgeNextBirthday, restriction.FunMaxCoverAge)

		if category != nil && (category.GlaEducatorBenefit == "Yes" || category.PtdEducatorBenefit == "Yes") {
			//educatorRates := GetEducatorRate(groupParameter, mr.AgeNextBirthday)
			mdrs.TotalEducatorSumAssured += mr.EducatorSumAtRisk
			mdrs.TotalGlaEducatorRiskPremium += mr.GlaEducatorRiskPremium
			mdrs.ExpAdjTotalGlaEducatorRiskPremium += mr.ExpAdjGlaEducatorRiskPremium
			mdrs.TotalPtdEducatorRiskPremium += mr.PtdEducatorRiskPremium
			mdrs.ExpAdjTotalPtdEducatorRiskPremium += mr.ExpAdjPtdEducatorRiskPremium
			mdrs.TotalGlaEducatorBinderAmount += mr.GlaEducatorBinderAmount
			mdrs.TotalGlaEducatorOutsourcedAmount += mr.GlaEducatorOutsourcedAmount
			mdrs.ExpAdjTotalGlaEducatorBinderAmount += mr.ExpAdjGlaEducatorBinderAmount
			mdrs.ExpAdjTotalGlaEducatorOutsourcedAmount += mr.ExpAdjGlaEducatorOutsourcedAmount
			mdrs.TotalPtdEducatorBinderAmount += mr.PtdEducatorBinderAmount
			mdrs.TotalPtdEducatorOutsourcedAmount += mr.PtdEducatorOutsourcedAmount
			mdrs.ExpAdjTotalPtdEducatorBinderAmount += mr.ExpAdjPtdEducatorBinderAmount
			mdrs.ExpAdjTotalPtdEducatorOutsourcedAmount += mr.ExpAdjPtdEducatorOutsourcedAmount
		}

		// Funeral sum assured aggregate — used as denominator for
		// FunConversionOnWithdrawal rate-per-1000.
		if category != nil && category.FamilyFuneralBenefit {
			mdrs.TotalFamilyFuneralSumAssured += mr.MemberFuneralSumAssured + mr.SpouseFuneralSumAssured + mr.ChildFuneralSumAssured + mr.ParentFuneralSumAssured + mr.ParentFuneralSumAssured
		}

		// Conversion / continuity slice rollups — risk + office (both ExpAdj
		// and non-ExpAdj). NOT added to any benefit total; these are a
		// reportable slice of the benefit's premium.
		mdrs.TotalGlaConversionOnWithdrawalAnnualRiskPremium += mr.GlaConversionOnWithdrawalRiskPremium
		mdrs.ExpAdjTotalGlaConversionOnWithdrawalAnnualRiskPremium += mr.ExpAdjGlaConversionOnWithdrawalRiskPremium
		mdrs.TotalGlaConversionOnRetirementAnnualRiskPremium += mr.GlaConversionOnRetirementRiskPremium
		mdrs.ExpAdjTotalGlaConversionOnRetirementAnnualRiskPremium += mr.ExpAdjGlaConversionOnRetirementRiskPremium
		mdrs.TotalGlaContinuityDuringDisabilityAnnualRiskPremium += mr.GlaContinuityDuringDisabilityRiskPremium
		mdrs.ExpAdjTotalGlaContinuityDuringDisabilityAnnualRiskPremium += mr.ExpAdjGlaContinuityDuringDisabilityRiskPremium
		mdrs.TotalGlaEducatorConversionOnWithdrawalAnnualRiskPremium += mr.GlaEducatorConversionOnWithdrawalRiskPremium
		mdrs.ExpAdjTotalGlaEducatorConversionOnWithdrawalAnnualRiskPremium += mr.ExpAdjGlaEducatorConversionOnWithdrawalRiskPremium
		mdrs.TotalGlaEducatorConversionOnRetirementAnnualRiskPremium += mr.GlaEducatorConversionOnRetirementRiskPremium
		mdrs.ExpAdjTotalGlaEducatorConversionOnRetirementAnnualRiskPremium += mr.ExpAdjGlaEducatorConversionOnRetirementRiskPremium
		mdrs.TotalGlaEducatorContinuityDuringDisabilityAnnualRiskPremium += mr.GlaEducatorContinuityDuringDisabilityRiskPremium
		mdrs.ExpAdjTotalGlaEducatorContinuityDuringDisabilityAnnualRiskPremium += mr.ExpAdjGlaEducatorContinuityDuringDisabilityRiskPremium
		mdrs.TotalPtdConversionOnWithdrawalAnnualRiskPremium += mr.PtdConversionOnWithdrawalRiskPremium
		mdrs.ExpAdjTotalPtdConversionOnWithdrawalAnnualRiskPremium += mr.ExpAdjPtdConversionOnWithdrawalRiskPremium
		mdrs.TotalPtdEducatorConversionOnWithdrawalAnnualRiskPremium += mr.PtdEducatorConversionOnWithdrawalRiskPremium
		mdrs.ExpAdjTotalPtdEducatorConversionOnWithdrawalAnnualRiskPremium += mr.ExpAdjPtdEducatorConversionOnWithdrawalRiskPremium
		mdrs.TotalPtdEducatorConversionOnRetirementAnnualRiskPremium += mr.PtdEducatorConversionOnRetirementRiskPremium
		mdrs.ExpAdjTotalPtdEducatorConversionOnRetirementAnnualRiskPremium += mr.ExpAdjPtdEducatorConversionOnRetirementRiskPremium
		mdrs.TotalPhiConversionOnWithdrawalAnnualRiskPremium += mr.PhiConversionOnWithdrawalRiskPremium
		mdrs.ExpAdjTotalPhiConversionOnWithdrawalAnnualRiskPremium += mr.ExpAdjPhiConversionOnWithdrawalRiskPremium
		mdrs.TotalTtdConversionOnWithdrawalAnnualRiskPremium += mr.TtdConversionOnWithdrawalRiskPremium
		mdrs.ExpAdjTotalTtdConversionOnWithdrawalAnnualRiskPremium += mr.ExpAdjTtdConversionOnWithdrawalRiskPremium
		mdrs.TotalCiConversionOnWithdrawalAnnualRiskPremium += mr.CiConversionOnWithdrawalRiskPremium
		mdrs.ExpAdjTotalCiConversionOnWithdrawalAnnualRiskPremium += mr.ExpAdjCiConversionOnWithdrawalRiskPremium
		mdrs.TotalSglaConversionOnWithdrawalAnnualRiskPremium += mr.SglaConversionOnWithdrawalRiskPremium
		mdrs.ExpAdjTotalSglaConversionOnWithdrawalAnnualRiskPremium += mr.ExpAdjSglaConversionOnWithdrawalRiskPremium
		mdrs.TotalFunConversionOnWithdrawalAnnualRiskPremium += mr.FunConversionOnWithdrawalRiskPremium
		mdrs.ExpAdjTotalFunConversionOnWithdrawalAnnualRiskPremium += mr.ExpAdjFunConversionOnWithdrawalRiskPremium

		// Binder and outsource fee aggregates per benefit (see applyBinderOutsourceAmounts).
		// Non-zero only when the quote's distribution channel is "binder".
		mdrs.TotalGlaAnnualBinderAmount += mr.GlaBinderAmount
		mdrs.TotalGlaAnnualOutsourcedAmount += mr.GlaOutsourcedAmount
		mdrs.ExpTotalGlaAnnualBinderAmount += mr.ExpAdjGlaBinderAmount
		mdrs.ExpTotalGlaAnnualOutsourcedAmount += mr.ExpAdjGlaOutsourcedAmount
		mdrs.TotalAdditionalAccidentalGlaAnnualBinderAmount += mr.AdditionalAccidentalGlaBinderAmount
		mdrs.TotalAdditionalAccidentalGlaAnnualOutsourcedAmt += mr.AdditionalAccidentalGlaOutsourcedAmount
		mdrs.ExpTotalAdditionalAccidentalGlaAnnualBinderAmount += mr.ExpAdjAdditionalAccidentalGlaBinderAmount
		mdrs.ExpTotalAdditionalAccidentalGlaAnnualOutsourcedAmt += mr.ExpAdjAdditionalAccidentalGlaOutsourcedAmt
		mdrs.TotalPtdAnnualBinderAmount += mr.PtdBinderAmount
		mdrs.TotalPtdAnnualOutsourcedAmount += mr.PtdOutsourcedAmount
		mdrs.ExpTotalPtdAnnualBinderAmount += mr.ExpAdjPtdBinderAmount
		mdrs.ExpTotalPtdAnnualOutsourcedAmount += mr.ExpAdjPtdOutsourcedAmount
		mdrs.TotalCiAnnualBinderAmount += mr.CiBinderAmount
		mdrs.TotalCiAnnualOutsourcedAmount += mr.CiOutsourcedAmount
		mdrs.ExpTotalCiAnnualBinderAmount += mr.ExpAdjCiBinderAmount
		mdrs.ExpTotalCiAnnualOutsourcedAmount += mr.ExpAdjCiOutsourcedAmount
		mdrs.TotalSglaAnnualBinderAmount += mr.SpouseGlaBinderAmount
		mdrs.TotalSglaAnnualOutsourcedAmount += mr.SpouseGlaOutsourcedAmount
		mdrs.ExpTotalSglaAnnualBinderAmount += mr.ExpAdjSpouseGlaBinderAmount
		mdrs.ExpTotalSglaAnnualOutsourcedAmount += mr.ExpAdjSpouseGlaOutsourcedAmount
		mdrs.TotalTtdAnnualBinderAmount += mr.TtdBinderAmount
		mdrs.TotalTtdAnnualOutsourcedAmount += mr.TtdOutsourcedAmount
		mdrs.ExpTotalTtdAnnualBinderAmount += mr.ExpAdjTtdBinderAmount
		mdrs.ExpTotalTtdAnnualOutsourcedAmount += mr.ExpAdjTtdOutsourcedAmount
		mdrs.TotalPhiAnnualBinderAmount += mr.PhiBinderAmount
		mdrs.TotalPhiAnnualOutsourcedAmount += mr.PhiOutsourcedAmount
		mdrs.ExpTotalPhiAnnualBinderAmount += mr.ExpAdjPhiBinderAmount
		mdrs.ExpTotalPhiAnnualOutsourcedAmount += mr.ExpAdjPhiOutsourcedAmount
		mdrs.TotalFunAnnualBinderAmount += mr.TotalFuneralBinderAmount
		mdrs.TotalFunAnnualOutsourcedAmount += mr.TotalFuneralOutsourcedAmount
		mdrs.ExpTotalFunAnnualBinderAmount += mr.ExpAdjTotalFuneralBinderAmount
		mdrs.ExpTotalFunAnnualOutsourcedAmount += mr.ExpAdjTotalFuneralOutsourcedAmount
		mdrs.TotalAnnualBinderAmount += mr.TotalBinderAmount
		mdrs.TotalAnnualOutsourcedAmount += mr.TotalOutsourcedAmount

		if mdrs.MemberCount == 1 {
			mdrs.MinGlaSumAssured = mr.GlaSumAssured
			mdrs.MinAdditionalAccidentalGlaSumAssured = mr.AdditionalAccidentalGlaSumAssured
			mdrs.MinPtdSumAssured = mr.PtdSumAssured
			mdrs.MinCiSumAssured = mr.CiSumAssured
			mdrs.MinSglaSumAssured = mr.SpouseGlaSumAssured
			mdrs.MinPhiIncome = mr.PhiIncome
			mdrs.MinTtdIncome = mr.TtdIncome
		} else {
			mdrs.MinGlaSumAssured = math.Min(mdrs.MinGlaSumAssured, mr.GlaSumAssured)
			mdrs.MinAdditionalAccidentalGlaSumAssured = math.Min(mdrs.MinAdditionalAccidentalGlaSumAssured, mr.AdditionalAccidentalGlaSumAssured)
			mdrs.MinPtdSumAssured = math.Min(mdrs.MinPtdSumAssured, mr.PtdSumAssured)
			mdrs.MinCiSumAssured = math.Min(mdrs.MinCiSumAssured, mr.CiSumAssured)
			mdrs.MinSglaSumAssured = math.Min(mdrs.MinSglaSumAssured, mr.SpouseGlaSumAssured)
			mdrs.MinPhiIncome = math.Min(mdrs.MinPhiIncome, mr.PhiIncome)
			mdrs.MinTtdIncome = math.Min(mdrs.MinTtdIncome, mr.TtdIncome)
		}
		mdrs.MaxGlaSumAssured = math.Max(mdrs.MaxGlaSumAssured, mr.GlaSumAssured)
		mdrs.MaxGlaCappedSumAssured = math.Max(mdrs.MaxGlaCappedSumAssured, mr.GlaCappedSumAssured)
		mdrs.MaxAdditionalAccidentalGlaSumAssured = math.Max(mdrs.MaxAdditionalAccidentalGlaSumAssured, mr.AdditionalAccidentalGlaSumAssured)
		mdrs.MaxAdditionalAccidentalGlaCappedSumAssured = math.Max(mdrs.MaxAdditionalAccidentalGlaCappedSumAssured, mr.AdditionalAccidentalGlaCappedSumAssured)
		mdrs.MaxPtdSumAssured = math.Max(mdrs.MaxPtdSumAssured, mr.PtdSumAssured)
		mdrs.MaxPtdCappedSumAssured = math.Max(mdrs.MaxPtdCappedSumAssured, mr.PtdCappedSumAssured)
		mdrs.MaxCiSumAssured = math.Max(mdrs.MaxCiSumAssured, mr.CiSumAssured)
		mdrs.MaxCiCappedSumAssured = math.Max(mdrs.MaxCiCappedSumAssured, mr.CiCappedSumAssured)
		mdrs.MaxSglaSumAssured = math.Max(mdrs.MaxSglaSumAssured, mr.SpouseGlaSumAssured)
		mdrs.MaxSglaCappedSumAssured = math.Max(mdrs.MaxSglaCappedSumAssured, mr.SpouseGlaCappedSumAssured)
		mdrs.MaxPhiIncome = math.Max(mdrs.MaxPhiIncome, mr.PhiIncome)
		mdrs.MaxPhiCappedIncome = math.Max(mdrs.MaxPhiCappedIncome, mr.PhiCappedIncome)
		mdrs.MaxTtdIncome = math.Max(mdrs.MaxTtdIncome, mr.TtdIncome)
		mdrs.MaxTtdCappedIncome = math.Max(mdrs.MaxTtdCappedIncome, mr.TtdCappedIncome)
	}

	// Compute derived summary fields after the reduce loop
	mdrs.TotalAnnualPremium = models.ComputeOfficePremium(mdrs.TotalGlaAnnualRiskPremium, &mdrs) + models.ComputeOfficePremium(mdrs.TotalPtdAnnualRiskPremium, &mdrs) + models.ComputeOfficePremium(mdrs.TotalTtdAnnualRiskPremium, &mdrs) + models.ComputeOfficePremium(mdrs.TotalPhiAnnualRiskPremium, &mdrs) + models.ComputeOfficePremium(mdrs.TotalCiAnnualRiskPremium, &mdrs) + models.ComputeOfficePremium(mdrs.TotalSglaAnnualRiskPremium, &mdrs)
	// Sub Total (Excl. Funeral) includes the Additional Accidental GLA rider,
	// the GLA TaxSaver rider, GLA Educator and PTD Educator components on top
	// of the six core benefits.
	mdrs.ExpTotalAnnualPremiumExclFuneral = models.ComputeOfficePremium(mdrs.ExpTotalGlaAnnualRiskPremium, &mdrs) +
		models.ComputeOfficePremium(mdrs.ExpTotalAdditionalAccidentalGlaAnnualRiskPremium, &mdrs) +
		models.ComputeOfficePremium(mdrs.ExpTotalPtdAnnualRiskPremium, &mdrs) +
		models.ComputeOfficePremium(mdrs.ExpTotalTtdAnnualRiskPremium, &mdrs) +
		models.ComputeOfficePremium(mdrs.ExpTotalPhiAnnualRiskPremium, &mdrs) +
		models.ComputeOfficePremium(mdrs.ExpTotalCiAnnualRiskPremium, &mdrs) +
		models.ComputeOfficePremium(mdrs.ExpTotalSglaAnnualRiskPremium, &mdrs) +
		models.ComputeOfficePremium(mdrs.ExpTotalTaxSaverAnnualRiskPremium, &mdrs) +
		models.ComputeOfficePremium(mdrs.ExpAdjTotalGlaEducatorRiskPremium, &mdrs) +
		models.ComputeOfficePremium(mdrs.ExpAdjTotalPtdEducatorRiskPremium, &mdrs)
	// TotalCommission is populated by applySchemeWideCommission after all
	// categories are computed: the tiered CommissionStructure bands run on
	// the scheme-wide total premium, and each category is allocated its
	// proportional share.
	mdrs.TotalCommission = 0
	mdrs.TotalExpenses = (mdrs.ExpTotalAnnualPremiumExclFuneral + models.ComputeOfficePremium(mdrs.ExpTotalFunAnnualRiskPremium, &mdrs)) * premiumLoading.ExpenseLoading
	mdrs.TotalExpectedClaims = mdrs.ExpTotalGlaAnnualRiskPremium + mdrs.ExpTotalPtdAnnualRiskPremium + mdrs.ExpTotalCiAnnualRiskPremium + mdrs.ExpTotalSglaAnnualRiskPremium + mdrs.ExpTotalTtdAnnualRiskPremium + mdrs.ExpTotalPhiAnnualRiskPremium + mdrs.ExpTotalFunAnnualRiskPremium

	// Delete results before saving new set of results
	emitProgress(selectedSchemeCategory, "saving_results")
	logger.Debug("Deleting existing results before saving new set")
	dbStartTime = time.Now()
	DB.Where("quote_id = ? and category = ?", quoteId, selectedSchemeCategory).Delete(&models.MemberRatingResult{})
	DB.Where("quote_id = ? and category = ?", quoteId, selectedSchemeCategory).Delete(&models.MemberPremiumSchedule{})
	DB.Where("quote_id = ? and category = ?", quoteId, selectedSchemeCategory).Delete(&models.MemberRatingResultSummary{})
	if groupQuote.ExperienceRating == "Yes" || groupQuote.ExperienceRating == "Override" {
		DB.Where("quote_id = ?", quoteId).Delete(&models.HistoricalCredibilityData{})
	}
	dbElapsed = time.Since(dbStartTime)
	logger.WithField("elapsed_ms", dbElapsed.Milliseconds()).Debug("Deleted existing results")

	// Save the results to the database

	logger.Debug("Saving member data results to database")
	dbStartTime = time.Now()
	err = DB.CreateInBatches(memberDataResults, 100).Error
	dbElapsed = time.Since(dbStartTime)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"error":      err.Error(),
			"elapsed_ms": dbElapsed.Milliseconds(),
		}).Error("Failed to save member data results")
	} else {
		logger.WithFields(map[string]interface{}{
			"record_count": len(memberDataResults),
			"elapsed_ms":   dbElapsed.Milliseconds(),
		}).Info("Successfully saved member data results")
	}

	// Save member premium schedule
	logger.Debug("Saving member premium schedule to database")
	dbStartTime = time.Now()
	err = DB.CreateInBatches(memberPremiumSchedule, 100).Error
	dbElapsed = time.Since(dbStartTime)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"error":      err.Error(),
			"elapsed_ms": dbElapsed.Milliseconds(),
		}).Error("Failed to save member premium schedule")
	} else {
		logger.WithFields(map[string]interface{}{
			"record_count": len(memberPremiumSchedule),
			"elapsed_ms":   dbElapsed.Milliseconds(),
		}).Info("Successfully saved member premium schedule")
	}

	// Bordereaux rows are no longer persisted — they are projected from
	// MemberRatingResult on-the-fly via BuildBordereauxRowsForQuote. The
	// in-memory `bordereaux` slice is still built during member iteration
	// because GroupPricingReinsurance writes the cession back into each
	// MemberRatingResult; the slice itself is now discarded.
	_ = bordereaux

	// Do averages for the member rating result summary
	mdrs.CredibilityRate = credibilityRate
	mdrs.ManuallyAddedCredibility = credibility
	mdrs.AnnualGlaExperienceWeightedRate = annualGlaExperienceWeightedRate
	if groupQuote.FreeCoverLimit > 0 {
		mdrs.FreeCoverLimit = groupQuote.FreeCoverLimit
	}
	if groupQuote.FreeCoverLimit == 0 {
		mdrs.FreeCoverLimit = calculatedFreeCoverLimit
	}

	// Persist cap limits so bordereaux's CoveredSumsAssured can reapply the
	// same caps the pricing flow used. A 0 value means "no limit".
	mdrs.MaximumGlaCover = restriction.MaximumGlaCover
	mdrs.MaximumPtdCover = restriction.MaximumPtdCover
	mdrs.SevereIllnessMaximumBenefit = restriction.SevereIllnessMaximumBenefit
	mdrs.SpouseGlaMaximumBenefit = restriction.SpouseGlaMaximumBenefit
	mdrs.TtdMaximumMonthlyBenefit = restriction.TtdMaximumMonthlyBenefit
	mdrs.PhiMaximumMonthlyBenefit = restriction.PhiMaximumMonthlyBenefit
	mdrs.GlaMaxCoverAge = restriction.GlaMaxCoverAge
	mdrs.PtdMaxCoverAge = restriction.PtdMaxCoverAge
	mdrs.CiMaxCoverAge = restriction.CiMaxCoverAge
	mdrs.TtdMaxCoverAge = restriction.TtdMaxCoverAge
	mdrs.PhiMaxCoverAge = restriction.PhiMaxCoverAge
	mdrs.FunMaxCoverAge = restriction.FunMaxCoverAge
	if category != nil {
		mdrs.ReinsMaxGlaCover = reinsCoverCaps[benefitTypeKey(category.GlaAlias, models.BenefitTypeGla)]
		mdrs.ReinsMaxPtdCover = reinsCoverCaps[benefitTypeKey(category.PtdAlias, models.BenefitTypePtd)]
		mdrs.ReinsMaxCiCover = reinsCoverCaps[benefitTypeKey(category.CiAlias, models.BenefitTypeCi)]
		mdrs.ReinsMaxSglaCover = reinsCoverCaps[benefitTypeKey(category.SglaAlias, models.BenefitTypeSgla)]
		mdrs.ReinsMaxTtdCover = reinsCoverCaps[benefitTypeKey(category.TtdAlias, models.BenefitTypeTtd)]
		mdrs.ReinsMaxPhiCover = reinsCoverCaps[benefitTypeKey(category.PhiAlias, models.BenefitTypePhi)]
		mdrs.ReinsMaxFunCover = reinsCoverCaps[benefitTypeKey(category.FamilyFuneralAlias, models.BenefitTypeFun)]
	}

	if mdrs.MemberCount > 0 {
		mdrs.AverageGlaCappedSumAssured = mdrs.TotalGlaCappedSumAssured / mdrs.MemberCount
		mdrs.AverageAdditionalAccidentalGlaCappedSumAssured = mdrs.TotalAdditionalAccidentalGlaCappedSumAssured / mdrs.MemberCount
		mdrs.AveragePtdCappedSumAssured = mdrs.TotalPtdCappedSumAssured / mdrs.MemberCount
		mdrs.AverageCiCappedSumAssured = mdrs.TotalCiCappedSumAssured / mdrs.MemberCount
		mdrs.AveragePhiCappedIncome = mdrs.TotalPhiCappedIncome / mdrs.MemberCount
		mdrs.AverageTtdCappedIncome = mdrs.TotalTtdCappedIncome / mdrs.MemberCount
		mdrs.AverageSglaCappedSumAssured = mdrs.TotalSglaCappedSumAssured / mdrs.MemberCount
		mdrs.TotalFunAnnualPremiumPerMember = models.ComputeOfficePremium(mdrs.TotalFunAnnualRiskPremium, &mdrs) / mdrs.MemberCount
		mdrs.TotalFunMonthlyPremiumPerMember = mdrs.TotalFunAnnualPremiumPerMember / 12.0
		mdrs.ExpTotalFunAnnualPremiumPerMember = models.ComputeOfficePremium(mdrs.ExpTotalFunAnnualRiskPremium, &mdrs) / mdrs.MemberCount
		mdrs.ExpTotalFunMonthlyPremiumPerMember = mdrs.ExpTotalFunAnnualPremiumPerMember / 12.0

	}
	mdrs.IndicativeRatesCount = indicativeRatesCount
	// Per-benefit denominators for Proportion*Salary fields. Members whose
	// age has passed a benefit's max cover age are excluded from that
	// benefit's denominator (they contribute 0 premium for it). SGLA and
	// AAGla follow the GLA total; conversion / continuity slices and
	// educator slices reuse the parent benefit's total. A 0 denominator
	// means no contributing members — leave the proportion at 0.
	glaDenomSalary := mdrs.TotalAnnualSalaryGla * indicativeRatesCount
	ptdDenomSalary := mdrs.TotalAnnualSalaryPtd * indicativeRatesCount
	ciDenomSalary := mdrs.TotalAnnualSalaryCi * indicativeRatesCount
	sglaDenomSalary := mdrs.TotalAnnualSalarySgla * indicativeRatesCount
	ttdDenomSalary := mdrs.TotalAnnualSalaryTtd * indicativeRatesCount
	phiDenomSalary := mdrs.TotalAnnualSalaryPhi * indicativeRatesCount
	funDenomSalary := mdrs.TotalAnnualSalaryFun * indicativeRatesCount

	if glaDenomSalary > 0 {
		mdrs.ProportionGlaAnnualRiskPremiumSalary = mdrs.TotalGlaAnnualRiskPremium / glaDenomSalary
		mdrs.ExpProportionGlaAnnualRiskPremiumSalary = mdrs.ExpTotalGlaAnnualRiskPremium / glaDenomSalary

		mdrs.ProportionAdditionalAccidentalGlaAnnualRiskPremiumSalary = mdrs.TotalAdditionalAccidentalGlaAnnualRiskPremium / glaDenomSalary
		mdrs.ExpProportionAdditionalAccidentalGlaAnnualRiskPremiumSalary = mdrs.ExpTotalAdditionalAccidentalGlaAnnualRiskPremium / glaDenomSalary

		// GLA conversion / continuity slices and educator slices.
		mdrs.ProportionGlaConversionOnWithdrawalRiskPremiumSalary = mdrs.TotalGlaConversionOnWithdrawalAnnualRiskPremium / glaDenomSalary
		mdrs.ExpAdjProportionGlaConversionOnWithdrawalRiskPremiumSalary = mdrs.ExpAdjTotalGlaConversionOnWithdrawalAnnualRiskPremium / glaDenomSalary
		mdrs.ProportionGlaConversionOnRetirementRiskPremiumSalary = mdrs.TotalGlaConversionOnRetirementAnnualRiskPremium / glaDenomSalary
		mdrs.ExpAdjProportionGlaConversionOnRetirementRiskPremiumSalary = mdrs.ExpAdjTotalGlaConversionOnRetirementAnnualRiskPremium / glaDenomSalary
		mdrs.ProportionGlaContinuityDuringDisabilityRiskPremiumSalary = mdrs.TotalGlaContinuityDuringDisabilityAnnualRiskPremium / glaDenomSalary
		mdrs.ExpAdjProportionGlaContinuityDuringDisabilityRiskPremiumSalary = mdrs.ExpAdjTotalGlaContinuityDuringDisabilityAnnualRiskPremium / glaDenomSalary
		mdrs.ProportionGlaEducatorConversionOnWithdrawalRiskPremiumSalary = mdrs.TotalGlaEducatorConversionOnWithdrawalAnnualRiskPremium / glaDenomSalary
		mdrs.ExpAdjProportionGlaEducatorConversionOnWithdrawalRiskPremiumSalary = mdrs.ExpAdjTotalGlaEducatorConversionOnWithdrawalAnnualRiskPremium / glaDenomSalary
		mdrs.ProportionGlaEducatorConversionOnRetirementRiskPremiumSalary = mdrs.TotalGlaEducatorConversionOnRetirementAnnualRiskPremium / glaDenomSalary
		mdrs.ExpAdjProportionGlaEducatorConversionOnRetirementRiskPremiumSalary = mdrs.ExpAdjTotalGlaEducatorConversionOnRetirementAnnualRiskPremium / glaDenomSalary
		mdrs.ProportionGlaEducatorContinuityDuringDisabilityRiskPremiumSalary = mdrs.TotalGlaEducatorContinuityDuringDisabilityAnnualRiskPremium / glaDenomSalary
		mdrs.ExpAdjProportionGlaEducatorContinuityDuringDisabilityRiskPremiumSalary = mdrs.ExpAdjTotalGlaEducatorContinuityDuringDisabilityAnnualRiskPremium / glaDenomSalary

		// Educator proportions of salary (split only — no combined)
		mdrs.ProportionGlaEducatorRiskPremiumSalary = mdrs.TotalGlaEducatorRiskPremium / glaDenomSalary
		mdrs.ExpAdjProportionGlaEducatorRiskPremiumSalary = mdrs.ExpAdjTotalGlaEducatorRiskPremium / glaDenomSalary
	}

	if ptdDenomSalary > 0 {
		mdrs.ProportionPtdAnnualRiskPremiumSalary = mdrs.TotalPtdAnnualRiskPremium / ptdDenomSalary
		mdrs.ExpProportionPtdAnnualRiskPremiumSalary = mdrs.ExpTotalPtdAnnualRiskPremium / ptdDenomSalary

		mdrs.ProportionPtdConversionOnWithdrawalRiskPremiumSalary = mdrs.TotalPtdConversionOnWithdrawalAnnualRiskPremium / ptdDenomSalary
		mdrs.ExpAdjProportionPtdConversionOnWithdrawalRiskPremiumSalary = mdrs.ExpAdjTotalPtdConversionOnWithdrawalAnnualRiskPremium / ptdDenomSalary
		mdrs.ProportionPtdEducatorConversionOnWithdrawalRiskPremiumSalary = mdrs.TotalPtdEducatorConversionOnWithdrawalAnnualRiskPremium / ptdDenomSalary
		mdrs.ExpAdjProportionPtdEducatorConversionOnWithdrawalRiskPremiumSalary = mdrs.ExpAdjTotalPtdEducatorConversionOnWithdrawalAnnualRiskPremium / ptdDenomSalary
		mdrs.ProportionPtdEducatorConversionOnRetirementRiskPremiumSalary = mdrs.TotalPtdEducatorConversionOnRetirementAnnualRiskPremium / ptdDenomSalary
		mdrs.ExpAdjProportionPtdEducatorConversionOnRetirementRiskPremiumSalary = mdrs.ExpAdjTotalPtdEducatorConversionOnRetirementAnnualRiskPremium / ptdDenomSalary

		mdrs.ProportionPtdEducatorRiskPremiumSalary = mdrs.TotalPtdEducatorRiskPremium / ptdDenomSalary
		mdrs.ExpAdjProportionPtdEducatorRiskPremiumSalary = mdrs.ExpAdjTotalPtdEducatorRiskPremium / ptdDenomSalary
	}

	if ttdDenomSalary > 0 {
		mdrs.ProportionTtdAnnualRiskPremiumSalary = mdrs.TotalTtdAnnualRiskPremium / ttdDenomSalary
		mdrs.ExpProportionTtdAnnualRiskPremiumSalary = mdrs.ExpTotalTtdAnnualRiskPremium / ttdDenomSalary

		mdrs.ProportionTtdConversionOnWithdrawalRiskPremiumSalary = mdrs.TotalTtdConversionOnWithdrawalAnnualRiskPremium / ttdDenomSalary
		mdrs.ExpAdjProportionTtdConversionOnWithdrawalRiskPremiumSalary = mdrs.ExpAdjTotalTtdConversionOnWithdrawalAnnualRiskPremium / ttdDenomSalary
	}

	if phiDenomSalary > 0 {
		mdrs.ProportionPhiAnnualRiskPremiumSalary = mdrs.TotalPhiAnnualRiskPremium / phiDenomSalary
		mdrs.ExpProportionPhiAnnualRiskPremiumSalary = mdrs.ExpTotalPhiAnnualRiskPremium / phiDenomSalary

		mdrs.ProportionPhiConversionOnWithdrawalRiskPremiumSalary = mdrs.TotalPhiConversionOnWithdrawalAnnualRiskPremium / phiDenomSalary
		mdrs.ExpAdjProportionPhiConversionOnWithdrawalRiskPremiumSalary = mdrs.ExpAdjTotalPhiConversionOnWithdrawalAnnualRiskPremium / phiDenomSalary
	}

	if ciDenomSalary > 0 {
		mdrs.ProportionCiAnnualRiskPremiumSalary = mdrs.TotalCiAnnualRiskPremium / ciDenomSalary
		mdrs.ExpProportionCiAnnualRiskPremiumSalary = mdrs.ExpTotalCiAnnualRiskPremium / ciDenomSalary

		mdrs.ProportionCiConversionOnWithdrawalRiskPremiumSalary = mdrs.TotalCiConversionOnWithdrawalAnnualRiskPremium / ciDenomSalary
		mdrs.ExpAdjProportionCiConversionOnWithdrawalRiskPremiumSalary = mdrs.ExpAdjTotalCiConversionOnWithdrawalAnnualRiskPremium / ciDenomSalary
	}

	if funDenomSalary > 0 {
		mdrs.ProportionFunAnnualRiskPremiumSalary = mdrs.TotalFunAnnualRiskPremium / funDenomSalary
		mdrs.ExpProportionFunAnnualRiskPremiumSalary = mdrs.ExpTotalFunAnnualRiskPremium / funDenomSalary

		mdrs.ProportionFunConversionOnWithdrawalRiskPremiumSalary = mdrs.TotalFunConversionOnWithdrawalAnnualRiskPremium / funDenomSalary
		mdrs.ExpAdjProportionFunConversionOnWithdrawalRiskPremiumSalary = mdrs.ExpAdjTotalFunConversionOnWithdrawalAnnualRiskPremium / funDenomSalary
	}

	if sglaDenomSalary > 0 {
		mdrs.ProportionSglaAnnualRiskPremiumSalary = mdrs.TotalSglaAnnualRiskPremium / sglaDenomSalary
		mdrs.ExpProportionSglaAnnualRiskPremiumSalary = mdrs.ExpTotalSglaAnnualRiskPremium / sglaDenomSalary

		mdrs.ProportionSglaConversionOnWithdrawalRiskPremiumSalary = mdrs.TotalSglaConversionOnWithdrawalAnnualRiskPremium / sglaDenomSalary
		mdrs.ExpAdjProportionSglaConversionOnWithdrawalRiskPremiumSalary = mdrs.ExpAdjTotalSglaConversionOnWithdrawalAnnualRiskPremium / sglaDenomSalary
	}

	// Reinsurance premium proportions per benefit:
	//   proportion = sum(<benefit>_reinsurance_premium) / sum(exp_adj_<benefit>_office_premium)
	// Each ratio is guarded against a zero denominator so disabled benefits
	// leave the proportion at 0 instead of NaN.
	if models.ComputeOfficePremium(mdrs.ExpTotalGlaAnnualRiskPremium, &mdrs) > 0 {
		mdrs.GlaReinsurancePremiumProportion = mdrs.TotalGlaReinsurancePremium / models.ComputeOfficePremium(mdrs.ExpTotalGlaAnnualRiskPremium, &mdrs)
	}
	if models.ComputeOfficePremium(mdrs.ExpTotalPtdAnnualRiskPremium, &mdrs) > 0 {
		mdrs.PtdReinsurancePremiumProportion = mdrs.TotalPtdReinsurancePremium / models.ComputeOfficePremium(mdrs.ExpTotalPtdAnnualRiskPremium, &mdrs)
	}
	if models.ComputeOfficePremium(mdrs.ExpTotalCiAnnualRiskPremium, &mdrs) > 0 {
		mdrs.CiReinsurancePremiumProportion = mdrs.TotalCiReinsurancePremium / models.ComputeOfficePremium(mdrs.ExpTotalCiAnnualRiskPremium, &mdrs)
	}
	if models.ComputeOfficePremium(mdrs.ExpTotalSglaAnnualRiskPremium, &mdrs) > 0 {
		mdrs.SglaReinsurancePremiumProportion = mdrs.TotalSglaReinsurancePremium / models.ComputeOfficePremium(mdrs.ExpTotalSglaAnnualRiskPremium, &mdrs)
	}
	if models.ComputeOfficePremium(mdrs.ExpTotalPhiAnnualRiskPremium, &mdrs) > 0 {
		mdrs.PhiReinsurancePremiumProportion = mdrs.TotalPhiReinsurancePremium / models.ComputeOfficePremium(mdrs.ExpTotalPhiAnnualRiskPremium, &mdrs)
	}
	if models.ComputeOfficePremium(mdrs.ExpTotalTtdAnnualRiskPremium, &mdrs) > 0 {
		mdrs.TtdReinsurancePremiumProportion = mdrs.TotalTtdReinsurancePremium / models.ComputeOfficePremium(mdrs.ExpTotalTtdAnnualRiskPremium, &mdrs)
	}
	if models.ComputeOfficePremium(mdrs.ExpTotalFunAnnualRiskPremium, &mdrs) > 0 {
		mdrs.FunReinsurancePremiumProportion = mdrs.TotalFunReinsurancePremium / models.ComputeOfficePremium(mdrs.ExpTotalFunAnnualRiskPremium, &mdrs)
	}

	if mdrs.TotalGlaCappedSumAssured > 0 {
		mdrs.GlaRiskRatePer1000SA = mdrs.TotalGlaAnnualRiskPremium * 1000.0 / mdrs.TotalGlaCappedSumAssured
		mdrs.ExpGlaRiskRatePer1000SA = mdrs.ExpTotalGlaAnnualRiskPremium * 1000.0 / mdrs.TotalGlaCappedSumAssured
	}

	if mdrs.TotalAdditionalAccidentalGlaCappedSumAssured > 0 {
		mdrs.AdditionalAccidentalGlaRiskRatePer1000SA = mdrs.TotalAdditionalAccidentalGlaAnnualRiskPremium * 1000.0 / mdrs.TotalAdditionalAccidentalGlaCappedSumAssured
		mdrs.ExpAdditionalAccidentalGlaRiskRatePer1000SA = mdrs.ExpTotalAdditionalAccidentalGlaAnnualRiskPremium * 1000.0 / mdrs.TotalAdditionalAccidentalGlaCappedSumAssured
	}

	if mdrs.TotalPtdCappedSumAssured > 0 {
		mdrs.PtdRiskRatePer1000SA = mdrs.TotalPtdAnnualRiskPremium * 1000.0 / mdrs.TotalPtdCappedSumAssured
		mdrs.ExpPtdRiskRatePer1000SA = mdrs.ExpTotalPtdAnnualRiskPremium * 1000.0 / mdrs.TotalPtdCappedSumAssured

	}

	if mdrs.TotalCiCappedSumAssured > 0 {
		mdrs.CiRiskRatePer1000SA = mdrs.TotalCiAnnualRiskPremium * 1000.0 / mdrs.TotalCiCappedSumAssured
		mdrs.ExpCiRiskRatePer1000SA = mdrs.ExpTotalCiAnnualRiskPremium * 1000.0 / mdrs.TotalCiCappedSumAssured

	}

	if mdrs.TotalSglaCappedSumAssured > 0 {
		mdrs.SglaRiskRatePer1000SA = mdrs.TotalSglaAnnualRiskPremium * 1000.0 / mdrs.TotalSglaCappedSumAssured
		mdrs.ExpSglaRiskRatePer1000SA = mdrs.ExpTotalSglaAnnualRiskPremium * 1000.0 / mdrs.TotalSglaCappedSumAssured

	}

	// Educator rate per 1000 (split between GLA and PTD components). Uses
	// TotalEducatorSumAssured (Grade0+17+812+Tertiary summed across members)
	// as the denominator so the rates are directly comparable.
	if mdrs.TotalEducatorSumAssured > 0 {
		mdrs.GlaEducatorRiskRatePer1000SA = mdrs.TotalGlaEducatorRiskPremium * 1000.0 / mdrs.TotalEducatorSumAssured
		mdrs.ExpGlaEducatorRiskRatePer1000SA = mdrs.ExpAdjTotalGlaEducatorRiskPremium * 1000.0 / mdrs.TotalEducatorSumAssured
		mdrs.PtdEducatorRiskRatePer1000SA = mdrs.TotalPtdEducatorRiskPremium * 1000.0 / mdrs.TotalEducatorSumAssured
		mdrs.ExpPtdEducatorRiskRatePer1000SA = mdrs.ExpAdjTotalPtdEducatorRiskPremium * 1000.0 / mdrs.TotalEducatorSumAssured
	}

	// Conversion / continuity slice rate-per-1000. Denominator uses the
	// relevant benefit's TotalCappedSumAssured (or TotalEducatorSumAssured
	// for educator slices, TotalPhiCappedIncome for PHI, and
	// TotalFamilyFuneralSumAssured for funeral).
	if mdrs.TotalGlaCappedSumAssured > 0 {
		d := mdrs.TotalGlaCappedSumAssured
		mdrs.GlaConversionOnWithdrawalRiskRatePer1000SA = mdrs.TotalGlaConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
		mdrs.ExpGlaConversionOnWithdrawalRiskRatePer1000SA = mdrs.ExpAdjTotalGlaConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
		mdrs.GlaConversionOnRetirementRiskRatePer1000SA = mdrs.TotalGlaConversionOnRetirementAnnualRiskPremium * 1000.0 / d
		mdrs.ExpGlaConversionOnRetirementRiskRatePer1000SA = mdrs.ExpAdjTotalGlaConversionOnRetirementAnnualRiskPremium * 1000.0 / d
		mdrs.GlaContinuityDuringDisabilityRiskRatePer1000SA = mdrs.TotalGlaContinuityDuringDisabilityAnnualRiskPremium * 1000.0 / d
		mdrs.ExpGlaContinuityDuringDisabilityRiskRatePer1000SA = mdrs.ExpAdjTotalGlaContinuityDuringDisabilityAnnualRiskPremium * 1000.0 / d
	}
	if mdrs.TotalEducatorSumAssured > 0 {
		d := mdrs.TotalEducatorSumAssured
		mdrs.GlaEducatorConversionOnWithdrawalRiskRatePer1000SA = mdrs.TotalGlaEducatorConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
		mdrs.ExpGlaEducatorConversionOnWithdrawalRiskRatePer1000SA = mdrs.ExpAdjTotalGlaEducatorConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
		mdrs.GlaEducatorConversionOnRetirementRiskRatePer1000SA = mdrs.TotalGlaEducatorConversionOnRetirementAnnualRiskPremium * 1000.0 / d
		mdrs.ExpGlaEducatorConversionOnRetirementRiskRatePer1000SA = mdrs.ExpAdjTotalGlaEducatorConversionOnRetirementAnnualRiskPremium * 1000.0 / d
		mdrs.GlaEducatorContinuityDuringDisabilityRiskRatePer1000SA = mdrs.TotalGlaEducatorContinuityDuringDisabilityAnnualRiskPremium * 1000.0 / d
		mdrs.ExpGlaEducatorContinuityDuringDisabilityRiskRatePer1000SA = mdrs.ExpAdjTotalGlaEducatorContinuityDuringDisabilityAnnualRiskPremium * 1000.0 / d
		mdrs.PtdEducatorConversionOnWithdrawalRiskRatePer1000SA = mdrs.TotalPtdEducatorConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
		mdrs.ExpPtdEducatorConversionOnWithdrawalRiskRatePer1000SA = mdrs.ExpAdjTotalPtdEducatorConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
		mdrs.PtdEducatorConversionOnRetirementRiskRatePer1000SA = mdrs.TotalPtdEducatorConversionOnRetirementAnnualRiskPremium * 1000.0 / d
		mdrs.ExpPtdEducatorConversionOnRetirementRiskRatePer1000SA = mdrs.ExpAdjTotalPtdEducatorConversionOnRetirementAnnualRiskPremium * 1000.0 / d
	}
	if mdrs.TotalPtdCappedSumAssured > 0 {
		d := mdrs.TotalPtdCappedSumAssured
		mdrs.PtdConversionOnWithdrawalRiskRatePer1000SA = mdrs.TotalPtdConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
		mdrs.ExpPtdConversionOnWithdrawalRiskRatePer1000SA = mdrs.ExpAdjTotalPtdConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
	}
	if mdrs.TotalCiCappedSumAssured > 0 {
		d := mdrs.TotalCiCappedSumAssured
		mdrs.CiConversionOnWithdrawalRiskRatePer1000SA = mdrs.TotalCiConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
		mdrs.ExpCiConversionOnWithdrawalRiskRatePer1000SA = mdrs.ExpAdjTotalCiConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
	}
	if mdrs.TotalSglaCappedSumAssured > 0 {
		d := mdrs.TotalSglaCappedSumAssured
		mdrs.SglaConversionOnWithdrawalRiskRatePer1000SA = mdrs.TotalSglaConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
		mdrs.ExpSglaConversionOnWithdrawalRiskRatePer1000SA = mdrs.ExpAdjTotalSglaConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
	}
	if mdrs.TotalPhiCappedIncome > 0 {
		d := mdrs.TotalPhiCappedIncome
		mdrs.PhiConversionOnWithdrawalRiskRatePer1000SA = mdrs.TotalPhiConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
		mdrs.ExpPhiConversionOnWithdrawalRiskRatePer1000SA = mdrs.ExpAdjTotalPhiConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
	}
	if mdrs.TotalTtdCappedIncome > 0 {
		d := mdrs.TotalTtdCappedIncome
		mdrs.TtdConversionOnWithdrawalRiskRatePer1000SA = mdrs.TotalTtdConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
		mdrs.ExpTtdConversionOnWithdrawalRiskRatePer1000SA = mdrs.ExpAdjTotalTtdConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
	}
	if mdrs.TotalFamilyFuneralSumAssured > 0 {
		d := mdrs.TotalFamilyFuneralSumAssured
		mdrs.FunConversionOnWithdrawalRiskRatePer1000SA = mdrs.TotalFunConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
		mdrs.ExpFunConversionOnWithdrawalRiskRatePer1000SA = mdrs.ExpAdjTotalFunConversionOnWithdrawalAnnualRiskPremium * 1000.0 / d
	}

	mdrs.IfStatus = groupQuote.Status
	mdrs.QuoteType = groupQuote.QuoteType
	mdrs.TotalAnnualPremium = models.ComputeOfficePremium(mdrs.ExpTotalGlaAnnualRiskPremium, &mdrs) + models.ComputeOfficePremium(mdrs.ExpTotalAdditionalAccidentalGlaAnnualRiskPremium, &mdrs) + models.ComputeOfficePremium(mdrs.ExpTotalPtdAnnualRiskPremium, &mdrs) + models.ComputeOfficePremium(mdrs.ExpTotalTtdAnnualRiskPremium, &mdrs) + models.ComputeOfficePremium(mdrs.ExpTotalPhiAnnualRiskPremium, &mdrs) + models.ComputeOfficePremium(mdrs.ExpTotalCiAnnualRiskPremium, &mdrs) + models.ComputeOfficePremium(mdrs.ExpTotalSglaAnnualRiskPremium, &mdrs) + models.ComputeOfficePremium(mdrs.ExpTotalFunAnnualRiskPremium, &mdrs)
	mdrs.TotalSumAssured = mdrs.TotalGlaSumAssured //+ mdrs.TotalPtdSumAssured + mdrs.TotalCiSumAssured + mdrs.TotalSpouseGlaSumAssured
	mdrs.PremiumRatesGuaranteedPeriodMonths = groupParameter.PremiumRatesGuaranteedPeriodMonths
	mdrs.QuoteValidityPeriodMonths = groupParameter.QuoteValidityPeriodMonths
	mdrs.ExpenseLoading = premiumLoading.ExpenseLoading
	mdrs.CommissionLoading = premiumLoading.CommissionLoading
	mdrs.ProfitLoading = premiumLoading.ProfitLoading
	mdrs.AdminLoading = premiumLoading.AdminLoading
	mdrs.OtherLoading = premiumLoading.OtherLoading
	summaryBinderRate, summaryOutsourceRate := binderAndOutsourceRates(&groupQuote)
	mdrs.BinderFeeRate = summaryBinderRate
	mdrs.OutsourceFeeRate = summaryOutsourceRate
	// AcceleratedBenefitDiscount is now per (age, gender) in GeneralLoading; no single summary value applies.

	//Adding credibility Data
	historicalCredibilityData.Year = mdrs.FinancialYear
	historicalCredibilityData.QuoteID = mdrs.QuoteId
	historicalCredibilityData.QuoteType = groupQuote.QuoteType
	historicalCredibilityData.SchemeName = groupQuote.SchemeName
	historicalCredibilityData.SchemeID = groupQuote.SchemeID
	var totalSchemeMemberCount int64
	if groupQuote.QuoteType == "Renewal" {
		DB.Model(&models.GPricingMemberDataInForce{}).Where("scheme_name = ?", groupQuote.SchemeName).Count(&totalSchemeMemberCount)
	} else {
		DB.Model(&models.GPricingMemberData{}).Where("quote_id = ?", groupQuote.ID).Count(&totalSchemeMemberCount)
	}
	historicalCredibilityData.MemberCount = int(totalSchemeMemberCount)
	historicalCredibilityData.CalculatedCredibility = mdrs.CredibilityRate
	historicalCredibilityData.ManuallyAddedCredibility = mdrs.ManuallyAddedCredibility
	historicalCredibilityData.WeightedLifeYears = weightedLifeYears
	historicalCredibilityData.FullCredibilityThreshold = groupParameter.FullCredibilityThreshold
	historicalCredibilityData.AnnualGlaExperienceRate = annualGlaExperienceWeightedRate
	historicalCredibilityData.AnnualPtdExperienceRate = annualPtdExperienceWeightedRate
	historicalCredibilityData.AnnualCiExperienceRate = annualCiExperienceWeightedRate

	logger.Debug("Saving member rating result summary to database")
	dbStartTime = time.Now()
	err = DB.Create(&mdrs).Error
	dbElapsed = time.Since(dbStartTime)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"error":      err.Error(),
			"elapsed_ms": dbElapsed.Milliseconds(),
		}).Error("Failed to save member rating result summary")
	} else {
		logger.WithFields(map[string]interface{}{
			"quote_id":   mdrs.QuoteId,
			"elapsed_ms": dbElapsed.Milliseconds(),
		}).Info("Successfully saved member rating result summary")
	}

	if groupQuote.ExperienceRating == "Yes" || groupQuote.ExperienceRating == "Override" {
		logger.Debug("Saving credibility data")
		dbStartTime = time.Now()
		historicalCredibilityData.Basis = basis
		historicalCredibilityData.CreationDate = time.Now()

		// Override mode has no claims-based weighted rate to blend; record the
		// actuary's manually-entered credibility so it surfaces in the
		// historical credibility table for future reference. Calculated and
		// claims-derived fields stay zero (meaningful: nothing was claims-
		// weighted).
		if groupQuote.ExperienceRating == "Override" {
			historicalCredibilityData.QuoteID = groupQuote.ID
			historicalCredibilityData.QuoteType = groupQuote.QuoteType
			historicalCredibilityData.SchemeName = groupQuote.SchemeName
			historicalCredibilityData.SchemeID = groupQuote.SchemeID
			historicalCredibilityData.ManuallyAddedCredibility = groupQuote.ExperienceOverrideCredibility
			// Wipe claims-derived fields so the row clearly reflects an
			// override-mode entry, not a stale Yes-mode one.
			historicalCredibilityData.CalculatedCredibility = 0
			historicalCredibilityData.AnnualGlaExperienceRate = 0
			historicalCredibilityData.AnnualPtdExperienceRate = 0
			historicalCredibilityData.AnnualCiExperienceRate = 0
			historicalCredibilityData.ClaimCount = 0
			historicalCredibilityData.ExperiencePeriod = 0
			historicalCredibilityData.WeightedLifeYears = 0

			// Per-benefit credibility = simple average of the per-(category)
			// credibility values for that benefit across every override row.
			avgs := averageOverrideCredibilityByBenefit(experienceRateOverrides)
			historicalCredibilityData.GlaCredibility = avgs[models.ExperienceRateOverrideBenefitGla]
			historicalCredibilityData.AaglaCredibility = avgs[models.ExperienceRateOverrideBenefitAagla]
			historicalCredibilityData.SglaCredibility = avgs[models.ExperienceRateOverrideBenefitSgla]
			historicalCredibilityData.PtdCredibility = avgs[models.ExperienceRateOverrideBenefitPtd]
			historicalCredibilityData.TtdCredibility = avgs[models.ExperienceRateOverrideBenefitTtd]
			historicalCredibilityData.PhiCredibility = avgs[models.ExperienceRateOverrideBenefitPhi]
			historicalCredibilityData.CiCredibility = avgs[models.ExperienceRateOverrideBenefitCi]
			historicalCredibilityData.FunCredibility = avgs[models.ExperienceRateOverrideBenefitFun]

			// Replace any prior rows for this quote so we don't accumulate.
			DB.Where("quote_id = ?", groupQuote.ID).Delete(&models.HistoricalCredibilityData{})
		} else {
			// Yes mode: claims-based credibility is a single quote-wide
			// number. Stamp it onto every per-benefit column so the audit
			// row is internally consistent (per-benefit credibility now
			// always populated, regardless of mode).
			c := historicalCredibilityData.ManuallyAddedCredibility
			if c == 0 {
				c = historicalCredibilityData.CalculatedCredibility
			}
			historicalCredibilityData.GlaCredibility = c
			historicalCredibilityData.AaglaCredibility = c
			historicalCredibilityData.SglaCredibility = c
			historicalCredibilityData.PtdCredibility = c
			historicalCredibilityData.TtdCredibility = c
			historicalCredibilityData.PhiCredibility = c
			historicalCredibilityData.CiCredibility = c
			historicalCredibilityData.FunCredibility = c
		}

		err = DB.Create(&historicalCredibilityData).Error
		dbElapsed = time.Since(dbStartTime)
		if err != nil {
			logger.WithFields(map[string]interface{}{
				"error":      err.Error(),
				"elapsed_ms": dbElapsed.Milliseconds(),
			}).Error("Failed to save credibility data")
		} else {
			logger.WithFields(map[string]interface{}{
				"quote_id":   historicalCredibilityData.QuoteID,
				"elapsed_ms": dbElapsed.Milliseconds(),
			}).Info("Successfully saved credibility data")
		}
	}

	logger.Debug("Clearing group pricing cache")
	return nil
}

func calculateAgeNextBirthday(commencementDate, dob time.Time) int {
	age := (commencementDate.Year() - dob.Year())

	// If the next birthday hasn't occurred yet this year, subtract 1 from the age
	birthdayThisYear := time.Date(commencementDate.Year(), dob.Month(), dob.Day(), 0, 0, 0, 0, commencementDate.Location())
	if commencementDate.Before(birthdayThisYear) {
		// Birthday is still upcoming this year, so next birthday is at current age
		return age
	}
	return age + 1
}

// applyCoverAgeLimit returns 0 when the member has aged past the benefit's
// max cover age. A maxCoverAge of 0 means "no limit" (field unset).
func applyCoverAgeLimit(value float64, ageNextBirthday, maxCoverAge int) float64 {
	if maxCoverAge > 0 && ageNextBirthday > maxCoverAge {
		return 0
	}
	return value
}

// applyMaxCoverCap caps value at maxCap when maxCap > 0. A maxCap of 0 means
// "no limit" (field unset).
func applyMaxCoverCap(value, maxCap float64) float64 {
	if maxCap > 0 && value > maxCap {
		return maxCap
	}
	return value
}

// benefitTypeKey returns the customised benefit code (alias) when set,
// otherwise the standard fallback code. Used so reinsurance cover-restriction
// lookups key by the same code the scheme uses for the benefit. Aliases are
// upper-cased for case-insensitive matching against stored benefit_type rows.
func benefitTypeKey(alias, fallback string) string {
	if alias != "" {
		return strings.ToUpper(alias)
	}
	return fallback
}

// resolveConvContLoadings populates per-member conversion-on-withdrawal,
// conversion-on-retirement, and continuity-during-disability slice loadings
// from the matching GeneralLoading row. Each slice is gated by its
// SchemeCategory flag and the underlying base benefit flag; missing flags
// leave the member loading at zero (no premium impact).
func resolveConvContLoadings(m *models.MemberRatingResult, sc *models.SchemeCategory, gl *models.GeneralLoading) {
	if sc.GlaBenefit {
		if sc.GlaConversionOnWithdrawal {
			m.GlaConversionOnWithdrawalLoading = gl.GlaConversionOnWithdrawalLoadingRate
		}
		if sc.GlaConversionOnRetirement {
			m.GlaConversionOnRetirementLoading = gl.GlaConversionOnRetirementLoadingRate
		}
		if sc.GlaContinuityDuringDisability {
			m.GlaContinuityDuringDisabilityLoading = gl.GlaContinuityDuringDisabilityLoadingRate
		}
		if sc.GlaEducatorBenefit == "Yes" {
			m.GlaEducatorLoading = gl.GlaEducatorLoadingRate
			if sc.GlaEducatorConversionOnWithdrawal {
				m.GlaEducatorConversionOnWithdrawalLoading = gl.GlaEducatorConversionOnWithdrawalLoadingRate
			}
			if sc.GlaEducatorConversionOnRetirement {
				m.GlaEducatorConversionOnRetirementLoading = gl.GlaEducatorConversionOnRetirementLoadingRate
			}
			if sc.GlaEducatorContinuityDuringDisability {
				m.GlaEducatorContinuityDuringDisabilityLoading = gl.GlaEducatorContinuityDuringDisabilityLoadingRate
			}
		}
	}
	if sc.PtdBenefit {
		if sc.PtdConversionOnWithdrawal {
			m.PtdConversionOnWithdrawalLoading = gl.PtdConversionOnWithdrawalLoadingRate
		}
		if sc.PtdEducatorBenefit == "Yes" {
			m.PtdEducatorLoading = gl.PtdEducatorLoadingRate
			if sc.PtdEducatorConversionOnWithdrawal {
				m.PtdEducatorConversionOnWithdrawalLoading = gl.PtdEducatorConversionOnWithdrawalLoadingRate
			}
			if sc.PtdEducatorConversionOnRetirement {
				m.PtdEducatorConversionOnRetirementLoading = gl.PtdEducatorConversionOnRetirementLoadingRate
			}
		}
	}
	if sc.CiBenefit && sc.CiConversionOnWithdrawal {
		m.CiConversionOnWithdrawalLoading = gl.CiConversionOnWithdrawalLoadingRate
	}
	if sc.PhiBenefit && sc.PhiConversionOnWithdrawal {
		m.PhiConversionOnWithdrawalLoading = gl.PhiConversionOnWithdrawalLoadingRate
	}
	if sc.TtdBenefit && sc.TtdConversionOnWithdrawal {
		m.TtdConversionOnWithdrawalLoading = gl.TtdConversionOnWithdrawalLoadingRate
	}
	if sc.SglaBenefit && sc.SglaConversionOnWithdrawal {
		m.SglaConversionOnWithdrawalLoading = gl.SglaConversionOnWithdrawalLoadingRate
	}
	if sc.FamilyFuneralBenefit && sc.FunConversionOnWithdrawal {
		m.FunConversionOnWithdrawalLoading = gl.FunConversionOnWithdrawalLoadingRate
	}
}

// computeEducatorLoadedRates folds the educator-specific conversion and
// continuity loadings into the GLA- and PTD-educator loaded rate columns.
// Must be called after LoadedGlaRate, ExpAdjLoadedGlaRate, LoadedPtdRate
// and ExpAdjLoadedPtdRate are set for the member.
func computeEducatorLoadedRates(m *models.MemberRatingResult) {
	edGlaMul := 1.0 + m.GlaEducatorLoading +
		m.GlaEducatorConversionOnWithdrawalLoading +
		m.GlaEducatorConversionOnRetirementLoading +
		m.GlaEducatorContinuityDuringDisabilityLoading
	m.LoadedGlaEducatorRate = m.LoadedGlaRate * edGlaMul
	m.ExpAdjLoadedGlaEducatorRate = m.ExpAdjLoadedGlaRate * edGlaMul

	edPtdMul := 1.0 + m.PtdEducatorLoading +
		m.PtdEducatorConversionOnWithdrawalLoading +
		m.PtdEducatorConversionOnRetirementLoading
	m.LoadedPtdEducatorRate = m.LoadedPtdRate * edPtdMul
	m.ExpAdjLoadedPtdEducatorRate = m.ExpAdjLoadedPtdRate * edPtdMul
}

// computeConvContSlicePremiums computes all non-educator conversion/continuity
// slice risk and office premiums. Educator slice premiums are computed inline
// in the educator block (they need the per-member riskWeightedSA factor).
// Must be called after every Loaded*Rate and TotalFuneralRiskCost are final.
func computeConvContSlicePremiums(m *models.MemberRatingResult, groupParameter models.GroupPricingParameters) {
	divisor := 1.0 - m.TotalPremiumLoading
	if divisor == 0 {
		divisor = 1.0
	}

	// Slice 1: GLA conversion on withdrawal
	m.GlaConversionOnWithdrawalRiskPremium = m.LoadedGlaRate * m.GlaConversionOnWithdrawalLoading * m.GlaCappedSumAssured
	m.ExpAdjGlaConversionOnWithdrawalRiskPremium = m.ExpAdjLoadedGlaRate * m.GlaConversionOnWithdrawalLoading * m.GlaCappedSumAssured

	// Slice 9: GLA conversion on retirement
	m.GlaConversionOnRetirementRiskPremium = m.LoadedGlaRate * m.GlaConversionOnRetirementLoading * m.GlaCappedSumAssured
	m.ExpAdjGlaConversionOnRetirementRiskPremium = m.ExpAdjLoadedGlaRate * m.GlaConversionOnRetirementLoading * m.GlaCappedSumAssured

	// Slice 12: GLA continuity during disability
	m.GlaContinuityDuringDisabilityRiskPremium = m.LoadedGlaRate * m.GlaContinuityDuringDisabilityLoading * m.GlaCappedSumAssured
	m.ExpAdjGlaContinuityDuringDisabilityRiskPremium = m.ExpAdjLoadedGlaRate * m.GlaContinuityDuringDisabilityLoading * m.GlaCappedSumAssured

	// Slice 4: PTD conversion on withdrawal
	m.PtdConversionOnWithdrawalRiskPremium = m.LoadedPtdRate * m.PtdConversionOnWithdrawalLoading * m.PtdCappedSumAssured
	m.ExpAdjPtdConversionOnWithdrawalRiskPremium = m.ExpAdjLoadedPtdRate * m.PtdConversionOnWithdrawalLoading * m.PtdCappedSumAssured

	// Slice 5: PHI conversion on withdrawal (base = PhiMonthlyBenefit, not SA)
	m.PhiConversionOnWithdrawalRiskPremium = m.LoadedPhiRate * m.PhiConversionOnWithdrawalLoading * m.PhiMonthlyBenefit
	m.ExpAdjPhiConversionOnWithdrawalRiskPremium = m.ExpAdjLoadedPhiRate * m.PhiConversionOnWithdrawalLoading * m.PhiMonthlyBenefit

	// Slice: TTD conversion on withdrawal (base = TtdCappedIncome × TtdNumberMonthlyPayments, parallels parent TTD premium)
	ttdIncomeBase := m.TtdCappedIncome * groupParameter.TtdNumberMonthlyPayments
	m.TtdConversionOnWithdrawalRiskPremium = m.LoadedTtdRate * m.TtdConversionOnWithdrawalLoading * ttdIncomeBase
	m.ExpAdjTtdConversionOnWithdrawalRiskPremium = m.ExpAdjLoadedTtdRate * m.TtdConversionOnWithdrawalLoading * ttdIncomeBase

	// Slice 6: CI conversion on withdrawal
	m.CiConversionOnWithdrawalRiskPremium = m.LoadedCiRate * m.CiConversionOnWithdrawalLoading * m.CiCappedSumAssured
	m.ExpAdjCiConversionOnWithdrawalRiskPremium = m.ExpAdjLoadedCiRate * m.CiConversionOnWithdrawalLoading * m.CiCappedSumAssured

	// Slice 7: SGLA conversion on withdrawal
	m.SglaConversionOnWithdrawalRiskPremium = m.LoadedSpouseGlaRate * m.SglaConversionOnWithdrawalLoading * m.SpouseGlaCappedSumAssured
	m.ExpAdjSglaConversionOnWithdrawalRiskPremium = m.ExpAdjLoadedSpouseGlaRate * m.SglaConversionOnWithdrawalLoading * m.SpouseGlaCappedSumAssured

	// Slice 8: FUN conversion on withdrawal — scales TotalFuneralRiskCost
	m.FunConversionOnWithdrawalRiskPremium = m.TotalFuneralRiskPremium * m.FunConversionOnWithdrawalLoading
	m.ExpAdjFunConversionOnWithdrawalRiskPremium = m.ExpAdjTotalFuneralRiskPremium * m.FunConversionOnWithdrawalLoading
}

// computeEducatorSlicePremiums computes the 5 educator-based slice premiums
// (GLA-ed withdrawal/retirement/continuity, PTD-ed withdrawal/retirement).
// Pass riskWeightedSA computed at the educator block call site.
func computeEducatorSlicePremiums(m *models.MemberRatingResult, riskWeightedSA float64) {
	divisor := 1.0 - m.TotalPremiumLoading
	if divisor == 0 {
		divisor = 1.0
	}

	// Slice 2: GLA Educator conversion on withdrawal
	m.GlaEducatorConversionOnWithdrawalRiskPremium = riskWeightedSA * m.LoadedGlaEducatorRate * m.GlaEducatorConversionOnWithdrawalLoading
	m.ExpAdjGlaEducatorConversionOnWithdrawalRiskPremium = riskWeightedSA * m.ExpAdjLoadedGlaEducatorRate * m.GlaEducatorConversionOnWithdrawalLoading

	// Slice 10: GLA Educator conversion on retirement
	m.GlaEducatorConversionOnRetirementRiskPremium = riskWeightedSA * m.LoadedGlaEducatorRate * m.GlaEducatorConversionOnRetirementLoading
	m.ExpAdjGlaEducatorConversionOnRetirementRiskPremium = riskWeightedSA * m.ExpAdjLoadedGlaEducatorRate * m.GlaEducatorConversionOnRetirementLoading

	// Slice 13: GLA Educator continuity during disability
	m.GlaEducatorContinuityDuringDisabilityRiskPremium = riskWeightedSA * m.LoadedGlaEducatorRate * m.GlaEducatorContinuityDuringDisabilityLoading
	m.ExpAdjGlaEducatorContinuityDuringDisabilityRiskPremium = riskWeightedSA * m.ExpAdjLoadedGlaEducatorRate * m.GlaEducatorContinuityDuringDisabilityLoading

	// Slice 3: PTD Educator conversion on withdrawal
	m.PtdEducatorConversionOnWithdrawalRiskPremium = riskWeightedSA * m.LoadedPtdEducatorRate * m.PtdEducatorConversionOnWithdrawalLoading
	m.ExpAdjPtdEducatorConversionOnWithdrawalRiskPremium = riskWeightedSA * m.ExpAdjLoadedPtdEducatorRate * m.PtdEducatorConversionOnWithdrawalLoading

	// Slice 11: PTD Educator conversion on retirement
	m.PtdEducatorConversionOnRetirementRiskPremium = riskWeightedSA * m.LoadedPtdEducatorRate * m.PtdEducatorConversionOnRetirementLoading
	m.ExpAdjPtdEducatorConversionOnRetirementRiskPremium = riskWeightedSA * m.ExpAdjLoadedPtdEducatorRate * m.PtdEducatorConversionOnRetirementLoading
}

// applyExperienceRateOverridesToMember mutates the per-benefit
// ExpAdjLoaded{Benefit}Rate fields based on the user-supplied per-(category,
// benefit) overrides. Must run after the LoadedRate × ExperienceAdjustment
// assignments and before the educator/conversion cascade reads those values.
//
// Mode == experience_rated → ExpAdjLoaded{Benefit}Rate is replaced by
//
//	OverrideRate (a direct annual-fraction substitution).
//
// Mode == theoretical (or row missing) → ExpAdjLoaded{Benefit}Rate is reset to
//
//	Loaded{Benefit}Rate so the rest of the pipeline behaves as if no
//	experience adjustment applied for that benefit.
//
// FUN is intentionally not handled here — funeral has no ExpAdjLoaded fun
// rate field; see applyFuneralExperienceRateOverride.
func applyExperienceRateOverridesToMember(m *models.MemberRatingResult, overrides map[string]models.GroupPricingExperienceRateOverride) {
	if len(overrides) == 0 {
		return
	}
	if row, ok := overrides[models.ExperienceRateOverrideBenefitGla]; ok {
		if row.Mode == models.ExperienceRateOverrideModeExperienceRated {
			m.ExpAdjLoadedGlaRate = row.OverrideRate
		} else {
			m.ExpAdjLoadedGlaRate = m.LoadedGlaRate
		}
	}
	if row, ok := overrides[models.ExperienceRateOverrideBenefitAagla]; ok {
		if row.Mode == models.ExperienceRateOverrideModeExperienceRated {
			m.ExpAdjLoadedAdditionalAccidentalGlaRate = row.OverrideRate
		} else {
			m.ExpAdjLoadedAdditionalAccidentalGlaRate = m.LoadedAdditionalAccidentalGlaRate
		}
	}
	if row, ok := overrides[models.ExperienceRateOverrideBenefitSgla]; ok {
		if row.Mode == models.ExperienceRateOverrideModeExperienceRated {
			m.ExpAdjLoadedSpouseGlaRate = row.OverrideRate
		} else {
			m.ExpAdjLoadedSpouseGlaRate = m.LoadedSpouseGlaRate
		}
	}
	if row, ok := overrides[models.ExperienceRateOverrideBenefitPtd]; ok {
		if row.Mode == models.ExperienceRateOverrideModeExperienceRated {
			m.ExpAdjLoadedPtdRate = row.OverrideRate
		} else {
			m.ExpAdjLoadedPtdRate = m.LoadedPtdRate
		}
	}
	if row, ok := overrides[models.ExperienceRateOverrideBenefitTtd]; ok {
		if row.Mode == models.ExperienceRateOverrideModeExperienceRated {
			m.ExpAdjLoadedTtdRate = row.OverrideRate
		} else {
			m.ExpAdjLoadedTtdRate = m.LoadedTtdRate
		}
	}
	if row, ok := overrides[models.ExperienceRateOverrideBenefitPhi]; ok {
		if row.Mode == models.ExperienceRateOverrideModeExperienceRated {
			m.ExpAdjLoadedPhiRate = row.OverrideRate
		} else {
			m.ExpAdjLoadedPhiRate = m.LoadedPhiRate
		}
	}
	if row, ok := overrides[models.ExperienceRateOverrideBenefitCi]; ok {
		if row.Mode == models.ExperienceRateOverrideModeExperienceRated {
			m.ExpAdjLoadedCiRate = row.OverrideRate
		} else {
			m.ExpAdjLoadedCiRate = m.LoadedCiRate
		}
	}
}

// applyFuneralExperienceRateOverride overrides the experience-adjusted
// funeral risk premium for a member when an FUN override row is present, and
// re-derives the dependent office and final office premiums. Theoretical
// mode collapses ExpAdjTotalFuneralRiskPremium to the unadjusted
// TotalFuneralRiskPremium; experience-rated mode applies the supplied
// override rate to the total funeral sum-assured for the member-category.
func applyFuneralExperienceRateOverride(
	m *models.MemberRatingResult,
	overrides map[string]models.GroupPricingExperienceRateOverride,
	schemeCategory models.SchemeCategory,
	discountFraction float64,
) {
	row, ok := overrides[models.ExperienceRateOverrideBenefitFun]
	if !ok {
		return
	}
	if row.Mode == models.ExperienceRateOverrideModeTheoretical {
		m.ExpAdjTotalFuneralRiskPremium = m.TotalFuneralRiskPremium
	} else {
		funeralSA := schemeCategory.FamilyFuneralMainMemberFuneralSumAssured +
			schemeCategory.FamilyFuneralSpouseFuneralSumAssured +
			schemeCategory.FamilyFuneralChildrenFuneralSumAssured*math.Min(m.AverageNumberChildren, float64(schemeCategory.FamilyFuneralMaxNumberChildren)) +
			schemeCategory.FamilyFuneralAdultDependantSumAssured*m.AverageNumberDependants
		m.ExpAdjTotalFuneralRiskPremium = row.OverrideRate * funeralSA * (1 + m.FunConversionOnWithdrawalLoading)
	}
	if denom := 1.0 - m.TotalPremiumLoading; denom != 0 {
		m.ExpAdjTotalFuneralOfficePremium = m.ExpAdjTotalFuneralRiskPremium / denom
	}
	if denom := 1.0 - (m.TotalPremiumLoading + discountFraction); denom != 0 {
		m.FinalTotalFuneralOfficePremium = m.ExpAdjTotalFuneralRiskPremium / denom
	}
}

func MovementPopulateRatesPerMember(memberDataPointResult *models.MemberRatingResult, bordereauxDatapoint models.Bordereaux, memberPremiumScheduleDatapoint *models.MemberPremiumSchedule, addedMemberInForce models.GPricingMemberDataInForce, groupQuote models.GroupPricingQuote,
	groupParameter models.GroupPricingParameters,
	groupPricingReinsuranceStructure models.GroupPricingReinsuranceStructure,
	incomeLevels []models.IncomeLevel, ageBands []models.GroupPricingAgeBands, credibilityRate, annualGlaExperienceWeightedRate, annualPtdExperienceWeightedRate, annualCiExperienceWeightedRate, manuallyAddedCredibility float64, finacialYear int, restriction models.Restriction, reinsCoverCaps map[string]float64, premiumLoading models.PremiumLoading, industryLoadingByGender map[string]models.IndustryLoading, regionLoadingByGender map[string]models.RegionLoading, reinsRegionLoadingByGender map[string]models.ReinsuranceRegionLoading) {

	//var memberPremiumScheduleDatapoint models.MemberPremiumSchedule
	var mpIncomeLevel int
	var mpAgeBand string
	var tempMp models.GPricingMemberData
	var originalMemberDataPointResult models.MemberRatingResult
	var groupFuneralParameter models.FuneralParameters
	var binderFeeRate, outsourceFeeRate float64
	//var bordereauxDatapoint models.Bordereaux

	tempMp.MemberName = addedMemberInForce.MemberName
	tempMp.AnnualSalary = addedMemberInForce.AnnualSalary

	schemeCategory := groupQuote.SchemeCategories[0]
	//memberDataPointResult.MovementType = "NewMember"
	memberDataPointResult.QuoteId = addedMemberInForce.QuoteId
	memberDataPointResult.FinancialYear = finacialYear
	memberDataPointResult.SchemeId = addedMemberInForce.SchemeId
	memberDataPointResult.MemberName = addedMemberInForce.MemberName
	memberDataPointResult.Gender = addedMemberInForce.Gender
	memberDataPointResult.DateOfBirth = addedMemberInForce.DateOfBirth
	memberDataPointResult.AnnualSalary = addedMemberInForce.AnnualSalary

	if !groupQuote.UseGlobalSalaryMultiple {
		memberDataPointResult.GlaSalaryMultiple = addedMemberInForce.Benefits.GlaMultiple
		if schemeCategory.PtdBenefit {
			memberDataPointResult.PtdSalaryMultiple = addedMemberInForce.Benefits.PtdMultiple
		}
		if schemeCategory.SglaBenefit {
			memberDataPointResult.SglaSalaryMultiple = addedMemberInForce.Benefits.SglaMultiple
		}
		if schemeCategory.CiBenefit {
			memberDataPointResult.CiSalaryMultiple = addedMemberInForce.Benefits.CiMultiple
		}
	}
	if groupQuote.UseGlobalSalaryMultiple {
		memberDataPointResult.GlaSalaryMultiple = schemeCategory.GlaSalaryMultiple
		memberDataPointResult.PtdSalaryMultiple = schemeCategory.PtdSalaryMultiple
		memberDataPointResult.SglaSalaryMultiple = schemeCategory.SglaSalaryMultiple
		memberDataPointResult.CiSalaryMultiple = schemeCategory.CiCriticalIllnessSalaryMultiple
	}

	//memberDataPointResult.GlaSalaryMultiple = addedMemberInForce.GlaSalaryMultiple
	memberDataPointResult.Occupation = groupQuote.Industry
	memberDataPointResult.OccupationClass = groupQuote.OccupationClass
	memberDataPointResult.Industry = groupQuote.Industry
	memberDataPointResult.ExpenseLoading = premiumLoading.ExpenseLoading
	memberDataPointResult.AdminLoading = premiumLoading.AdminLoading
	effectiveCommission := premiumLoading.CommissionLoading
	if groupQuote.DistributionChannel == models.ChannelDirect {
		effectiveCommission = 0
	}
	memberDataPointResult.CommissionLoading = effectiveCommission
	memberDataPointResult.ProfitLoading = premiumLoading.ProfitLoading
	memberDataPointResult.OtherLoading = premiumLoading.OtherLoading
	memberDataPointResult.Discount = -(groupQuote.Loadings.Discount / 100.0)

	if groupQuote.DistributionChannel == "binder" {
		binderFeeRate, outsourceFeeRate = binderAndOutsourceRates(&groupQuote)
		memberDataPointResult.BinderFeeRate = binderFeeRate
		memberDataPointResult.OutsourceFeeRate = outsourceFeeRate
	}

	// Commission is excluded from TotalLoading — it is computed at the scheme
	// level on total premium via the tiered CommissionStructure bands and
	// distributed per benefit in applySchemeWideCommission.
	memberDataPointResult.TotalPremiumLoading = math.Max(premiumLoading.ExpenseLoading+premiumLoading.AdminLoading+premiumLoading.ProfitLoading+premiumLoading.OtherLoading+binderFeeRate+outsourceFeeRate, premiumLoading.MinimumPremiumLoading)
	memberDataPointResult.ExpCredibility = credibilityRate
	memberDataPointResult.GlaWeightedExperienceCrudeRate = annualGlaExperienceWeightedRate
	memberDataPointResult.PtdExperienceCrudeRate = annualPtdExperienceWeightedRate
	memberDataPointResult.CiExperienceCrudeRate = annualCiExperienceWeightedRate
	memberDataPointResult.ManuallyAddedCredibility = manuallyAddedCredibility

	if manuallyAddedCredibility > 0 {
		credibilityRate = manuallyAddedCredibility
	}

	//memberPremiumScheduleDatapoint.MemberName = memberDataPointResult.MemberName
	//memberPremiumScheduleDatapoint.QuoteId = memberDataPointResult.QuoteId
	//memberPremiumScheduleDatapoint.Gender = memberDataPointResult.Gender
	//memberPremiumScheduleDatapoint.SchemeName = memberDataPointResult.SchemeName

	//dob, err := utils.ParseDateString(addedMemberInForce.DateOfBirth)
	//if err != nil {
	//	fmt.Println("error encountered parsing date of birth: ", err)
	//	//return err
	//}
	memberDataPointResult.AgeNextBirthday = calculateAgeNextBirthday(groupQuote.CommencementDate, addedMemberInForce.DateOfBirth)
	mpAgeBand = GetAgeBand(memberDataPointResult.AgeNextBirthday, ageBands)
	groupFuneralParameter = GetFuneralParameter(groupParameter, memberDataPointResult.AgeNextBirthday)
	memberDataPointResult.AgeBand = mpAgeBand
	if addedMemberInForce.Gender == "Male" {
		memberDataPointResult.SpouseGender = "Female"
		memberDataPointResult.SpouseAgeNextBirthday = int(math.Min(math.Max(float64(memberDataPointResult.AgeNextBirthday)-
			float64(groupParameter.SpouseAgeGap), float64(restriction.MinEntryAge)), float64(restriction.MaxEntryAge)))

	}
	if addedMemberInForce.Gender == "Female" {
		memberDataPointResult.SpouseGender = "Male"
		memberDataPointResult.SpouseAgeNextBirthday = int(math.Min(math.Max(float64(memberDataPointResult.AgeNextBirthday)+
			float64(groupParameter.SpouseAgeGap), float64(restriction.MinEntryAge)), float64(restriction.MaxEntryAge)))
	}

	mpIncomeLevel = GetIncomeLevel(tempMp, incomeLevels)

	memberDataPointResult.GlaExperienceAdjustment = 1
	memberDataPointResult.PtdExperienceAdjustment = 1
	memberDataPointResult.CiExperienceAdjustment = 1
	memberDataPointResult.PhiExperienceAdjustment = 1
	memberDataPointResult.TtdExperienceAdjustment = 1

	memberDataPointResult.AverageDependantAgeNextBirthday = groupFuneralParameter.AverageDependantAge
	memberDataPointResult.AverageChildAgeNextBirthday = groupFuneralParameter.AverageChildAge
	memberDataPointResult.AverageNumberDependants = groupFuneralParameter.NumberDependants
	memberDataPointResult.AverageNumberChildren = groupFuneralParameter.NumberChildren
	memberDataPointResult.GlaSumAssured = addedMemberInForce.AnnualSalary * memberDataPointResult.GlaSalaryMultiple
	memberDataPointResult.GlaCappedSumAssured = applyMaxCoverCap(applyMaxCoverCap(math.Min(memberDataPointResult.GlaSumAssured, groupQuote.FreeCoverLimit), restriction.MaximumGlaCover), reinsCoverCaps[benefitTypeKey(schemeCategory.GlaAlias, models.BenefitTypeGla)])
	memberDataPointResult.PtdSumAssured = addedMemberInForce.AnnualSalary * memberDataPointResult.PtdSalaryMultiple
	memberDataPointResult.PtdCappedSumAssured = applyMaxCoverCap(applyMaxCoverCap(math.Min(memberDataPointResult.PtdSumAssured, memberDataPointResult.GlaCappedSumAssured), restriction.MaximumPtdCover), reinsCoverCaps[benefitTypeKey(schemeCategory.PtdAlias, models.BenefitTypePtd)])
	memberDataPointResult.CiSumAssured = addedMemberInForce.AnnualSalary * memberDataPointResult.CiSalaryMultiple
	memberDataPointResult.CiCappedSumAssured = applyMaxCoverCap(math.Min(math.Min(memberDataPointResult.CiSumAssured, restriction.SevereIllnessMaximumBenefit), memberDataPointResult.GlaCappedSumAssured), reinsCoverCaps[benefitTypeKey(schemeCategory.CiAlias, models.BenefitTypeCi)])
	memberDataPointResult.SpouseGlaSumAssured = addedMemberInForce.AnnualSalary * memberDataPointResult.SglaSalaryMultiple
	memberDataPointResult.SpouseGlaCappedSumAssured = applyMaxCoverCap(math.Min(math.Min(memberDataPointResult.SpouseGlaSumAssured, restriction.SpouseGlaMaximumBenefit), memberDataPointResult.GlaCappedSumAssured), reinsCoverCaps[benefitTypeKey(schemeCategory.SglaAlias, models.BenefitTypeSgla)])

	memberDataPointResult.TtdIncome = addedMemberInForce.AnnualSalary * schemeCategory.TtdIncomeReplacementPercentage / 12.0 / 100.0
	memberDataPointResult.TtdCappedIncome = applyMaxCoverCap(math.Min(memberDataPointResult.TtdIncome, restriction.TtdMaximumMonthlyBenefit), reinsCoverCaps[benefitTypeKey(schemeCategory.TtdAlias, models.BenefitTypeTtd)])

	memberDataPointResult.PhiIncome = addedMemberInForce.AnnualSalary * schemeCategory.PhiIncomeReplacementPercentage / 12.0 / 100.0
	memberDataPointResult.PhiCappedIncome = applyMaxCoverCap(math.Min(memberDataPointResult.PhiIncome, restriction.PhiMaximumMonthlyBenefit), reinsCoverCaps[benefitTypeKey(schemeCategory.PhiAlias, models.BenefitTypePhi)])

	mpAge := memberDataPointResult.AgeNextBirthday
	memberDataPointResult.GlaCappedSumAssured = applyCoverAgeLimit(memberDataPointResult.GlaCappedSumAssured, mpAge, restriction.GlaMaxCoverAge)
	memberDataPointResult.PtdCappedSumAssured = applyCoverAgeLimit(memberDataPointResult.PtdCappedSumAssured, mpAge, restriction.PtdMaxCoverAge)
	memberDataPointResult.CiCappedSumAssured = applyCoverAgeLimit(memberDataPointResult.CiCappedSumAssured, mpAge, restriction.CiMaxCoverAge)
	memberDataPointResult.SpouseGlaCappedSumAssured = applyCoverAgeLimit(memberDataPointResult.SpouseGlaCappedSumAssured, mpAge, restriction.GlaMaxCoverAge)
	memberDataPointResult.TtdCappedIncome = applyCoverAgeLimit(memberDataPointResult.TtdCappedIncome, mpAge, restriction.TtdMaxCoverAge)
	memberDataPointResult.PhiCappedIncome = applyCoverAgeLimit(memberDataPointResult.PhiCappedIncome, mpAge, restriction.PhiMaxCoverAge)

	memberDataPointResult.IncomeReplacementRatio = schemeCategory.PhiIncomeReplacementPercentage / 100.0

	if schemeCategory.PhiPremiumWaiver == "Yes" {
		memberDataPointResult.PhiContributionWaiver = math.Min(addedMemberInForce.AnnualSalary*addedMemberInForce.ContributionWaiverProportion, restriction.PhiMaximumMonthlyContributionWaiver)
	}
	if schemeCategory.PhiMedicalAidPremiumWaiver == "Yes" {
		var waiver float64
		if GetMedicalAidWaiverMethod() == models.MedicalAidWaiverMethodTableLookup {
			waiver = GetMedicalWaiverSumAtRisk(memberDataPointResult, groupParameter, mpIncomeLevel)
		} else {
			waiver = addedMemberInForce.AnnualSalary*groupParameter.MedicalAidWaiverProportion + groupParameter.MedicalAidWaiverAmount
		}
		memberDataPointResult.PhiMedicalAidWaiver = math.Min(waiver, restriction.MaxMedicalAidWaiver)
	}

	memberDataPointResult.PhiMonthlyBenefit = memberDataPointResult.PhiCappedIncome + memberDataPointResult.PhiContributionWaiver + memberDataPointResult.PhiMedicalAidWaiver

	data, _ := json.Marshal(memberDataPointResult)
	json.Unmarshal(data, &originalMemberDataPointResult)
	memberIndustryLoading := industryLoadingByGender[strings.ToUpper(originalMemberDataPointResult.Gender[:1])]
	memberRegionLoading := regionLoadingByGender[strings.ToUpper(originalMemberDataPointResult.Gender[:1])]

	// Persist the per-member region and industry loadings so formulas below
	// can read them directly from MemberRatingResult.
	memberDataPointResult.GlaRegionLoading = memberRegionLoading.GlaRegionLoadingRate
	memberDataPointResult.GlaAidsRegionLoading = memberRegionLoading.GlaAidsRegionLoadingRate
	memberDataPointResult.PtdRegionLoading = memberRegionLoading.PtdRegionLoadingRate
	memberDataPointResult.CiRegionLoading = memberRegionLoading.CiRegionLoadingRate
	memberDataPointResult.TtdRegionLoading = memberRegionLoading.TtdRegionLoadingRate
	memberDataPointResult.PhiRegionLoading = memberRegionLoading.PhiRegionLoadingRate
	memberDataPointResult.FunRegionLoading = memberRegionLoading.FunRegionLoadingRate
	memberDataPointResult.FunAidsRegionLoading = memberRegionLoading.FunAidsRegionLoadingRate

	memberDataPointResult.GlaIndustryLoading = memberIndustryLoading.GlaIndustryLoadingRate
	memberDataPointResult.PtdIndustryLoading = memberIndustryLoading.PtdIndustryLoadingRate
	memberDataPointResult.CiIndustryLoading = memberIndustryLoading.CiIndustryLoadingRate
	memberDataPointResult.TtdIndustryLoading = memberIndustryLoading.TtdIndustryLoadingRate
	memberDataPointResult.PhiIndustryLoading = memberIndustryLoading.PhiIndustryLoadingRate

	if schemeCategory.GlaBenefit {
		memberDataPointResult.GlaQx = applyCoverAgeLimit(GetGlaRate(&originalMemberDataPointResult, groupParameter, mpIncomeLevel, groupQuote, schemeCategory), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.GlaAidsQx = applyCoverAgeLimit(GetGlaAidsRate(memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.BaseGlaRate = memberDataPointResult.GlaQx*(1+memberDataPointResult.GlaIndustryLoading+memberDataPointResult.GlaRegionLoading) + memberDataPointResult.GlaAidsQx*(1+memberDataPointResult.GlaAidsRegionLoading)
	} else if schemeCategory.FamilyFuneralBenefit {
		funQx := applyCoverAgeLimit(GetFuneralRate(memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
		funAidsQx := applyCoverAgeLimit(GetFuneralAidsRate(memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
		memberDataPointResult.MainMemberFuneralBaseRate = funQx*(1+memberDataPointResult.FunRegionLoading) + funAidsQx*(1+memberDataPointResult.FunAidsRegionLoading)
		if len(memberDataPointResult.SpouseGender) > 0 {
			spouseRegionLoading := regionLoadingByGender[strings.ToUpper(memberDataPointResult.SpouseGender[:1])]
			spouseFunQx := applyCoverAgeLimit(GetSpouseFuneralRate(memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
			spouseFunAidsQx := applyCoverAgeLimit(GetSpouseFuneralAidsRate(memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
			memberDataPointResult.SpouseFuneralBaseRate = spouseFunQx*(1+spouseRegionLoading.FunRegionLoadingRate) + spouseFunAidsQx*(1+spouseRegionLoading.FunAidsRegionLoadingRate)
		}
	}

	gl1 := GetGeneralLoading(groupParameter.RiskRateCode, memberDataPointResult.AgeNextBirthday, memberDataPointResult.Gender)
	if schemeCategory.GlaTerminalIllnessBenefit == "Yes" {
		memberDataPointResult.GlaTerminalIllnessLoading = gl1.TerminalIllnessLoadingRate
	}

	// Persist the per-member contingency loadings alongside the already-set
	// region/industry loadings.
	memberDataPointResult.GlaContingencyLoading = gl1.GlaContigencyLoadingRate
	memberDataPointResult.PtdContingencyLoading = gl1.PtdContigencyLoadingRate
	memberDataPointResult.CiContingencyLoading = gl1.CiContigencyLoadingRate
	memberDataPointResult.TtdContingencyLoading = gl1.TtdContigencyLoadingRate
	memberDataPointResult.PhiContingencyLoading = gl1.PhiContigencyLoadingRate
	memberDataPointResult.FunContingencyLoading = gl1.FunContigencyLoadingRate

	// Voluntary loadings only apply when the quote's obligation type is
	// Voluntary. For Compulsory (or any other value) they stay zero so the
	// (1 + ContingencyLoading + VoluntaryLoading + ...) multiplier collapses
	// back to the compulsory-only formula.
	if groupQuote.ObligationType == "Voluntary" {
		memberDataPointResult.GlaVoluntaryLoading = gl1.GlaVoluntaryLoadingRate
		memberDataPointResult.PtdVoluntaryLoading = gl1.PtdVoluntaryLoadingRate
		memberDataPointResult.CiVoluntaryLoading = gl1.CiVoluntaryLoadingRate
		memberDataPointResult.TtdVoluntaryLoading = gl1.TtdVoluntaryLoadingRate
		memberDataPointResult.PhiVoluntaryLoading = gl1.PhiVoluntaryLoadingRate
		memberDataPointResult.FunVoluntaryLoading = gl1.FunVoluntaryLoadingRate
	}

	// TaxSaver: resolved per-member from GeneralLoading (by age/gender/
	// risk_rate_code) when the category opts into the rider. Added to the
	// GLA loading chain so LoadedGlaRate reflects the full premium.
	if schemeCategory.GlaBenefit && schemeCategory.TaxSaverBenefit {
		memberDataPointResult.TaxSaverLoading = gl1.TaxSaverLoadingRate
	}

	// Conversion / continuity slice loadings (folded into each benefit's
	// Loaded*Rate below; slice premiums computed after rates are final).
	resolveConvContLoadings(memberDataPointResult, &schemeCategory, &gl1)

	memberDataPointResult.LoadedGlaRate = memberDataPointResult.BaseGlaRate * (1 + memberDataPointResult.GlaContingencyLoading + memberDataPointResult.GlaVoluntaryLoading + memberDataPointResult.GlaTerminalIllnessLoading + memberDataPointResult.GlaContinuityDuringDisabilityLoading + memberDataPointResult.GlaConversionOnWithdrawalLoading + memberDataPointResult.GlaConversionOnRetirementLoading)

	memberDataPointResult.ExpAdjLoadedGlaRate = memberDataPointResult.LoadedGlaRate * memberDataPointResult.GlaExperienceAdjustment

	// Additional Accidental GLA — optional sub-benefit that re-uses every GLA
	// parameter (sum assured, loadings, experience adjustment) but prices
	// against a different benefit_type row in gla_rates. Skipped when the
	// additional benefit type equals the main GLA benefit type — that case
	// would duplicate the main GLA premium.
	if schemeCategory.GlaBenefit && schemeCategory.AdditionalAccidentalGlaBenefit &&
		len(schemeCategory.AdditionalAccidentalGlaBenefitType) > 0 &&
		schemeCategory.AdditionalAccidentalGlaBenefitType != schemeCategory.GlaBenefitType {
		memberDataPointResult.AdditionalAccidentalGlaSumAssured = memberDataPointResult.GlaSumAssured
		memberDataPointResult.AdditionalAccidentalGlaCappedSumAssured = memberDataPointResult.GlaCappedSumAssured
		memberDataPointResult.AdditionalAccidentalGlaQx = applyCoverAgeLimit(GetAdditionalAccidentalGlaRate(&originalMemberDataPointResult, groupParameter, mpIncomeLevel, groupQuote, schemeCategory), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.BaseAdditionalAccidentalGlaRate = memberDataPointResult.AdditionalAccidentalGlaQx * (1 + memberDataPointResult.GlaIndustryLoading + memberDataPointResult.GlaRegionLoading)
		memberDataPointResult.LoadedAdditionalAccidentalGlaRate = memberDataPointResult.BaseAdditionalAccidentalGlaRate * (1 + memberDataPointResult.GlaContingencyLoading)
		memberDataPointResult.AdditionalAccidentalGlaExperienceAdjustment = memberDataPointResult.GlaExperienceAdjustment
		memberDataPointResult.ExpAdjLoadedAdditionalAccidentalGlaRate = memberDataPointResult.LoadedAdditionalAccidentalGlaRate * memberDataPointResult.AdditionalAccidentalGlaExperienceAdjustment
	}

	if schemeCategory.PtdBenefit {
		ptdRate := applyCoverAgeLimit(GetPtdRate(&originalMemberDataPointResult, groupParameter, groupQuote, schemeCategory, mpIncomeLevel), memberDataPointResult.AgeNextBirthday, restriction.PtdMaxCoverAge)
		if schemeCategory.PtdBenefitType == "Accelerated" {
			memberDataPointResult.BasePtdRate = ptdRate * (1 + memberDataPointResult.PtdIndustryLoading + memberDataPointResult.PtdRegionLoading - gl1.PtdAcceleratedBenefitDiscount)
		} else {
			memberDataPointResult.BasePtdRate = ptdRate * (1 + memberDataPointResult.PtdIndustryLoading + memberDataPointResult.PtdRegionLoading)
		}
		memberDataPointResult.LoadedPtdRate = memberDataPointResult.BasePtdRate * (1 + memberDataPointResult.PtdContingencyLoading + memberDataPointResult.PtdVoluntaryLoading + memberDataPointResult.PtdConversionOnWithdrawalLoading)
		memberDataPointResult.ExpAdjLoadedPtdRate = memberDataPointResult.LoadedPtdRate * memberDataPointResult.PtdExperienceAdjustment
	}

	if schemeCategory.TtdBenefit {
		ttdRate := applyCoverAgeLimit(GetTtdRate(&originalMemberDataPointResult, groupParameter, groupQuote, schemeCategory, mpIncomeLevel), memberDataPointResult.AgeNextBirthday, restriction.TtdMaxCoverAge)
		memberDataPointResult.BaseTtdRate = ttdRate * (1 + memberDataPointResult.TtdIndustryLoading + memberDataPointResult.TtdRegionLoading)
		memberDataPointResult.LoadedTtdRate = memberDataPointResult.BaseTtdRate * (1 + memberDataPointResult.TtdContingencyLoading + memberDataPointResult.TtdVoluntaryLoading + memberDataPointResult.TtdConversionOnWithdrawalLoading)
		memberDataPointResult.ExpAdjLoadedTtdRate = memberDataPointResult.LoadedTtdRate * memberDataPointResult.TtdExperienceAdjustment
	}

	if schemeCategory.PhiBenefit {
		phiRate := applyCoverAgeLimit(GetPhiRate(&originalMemberDataPointResult, groupParameter, groupQuote, schemeCategory, mpIncomeLevel), memberDataPointResult.AgeNextBirthday, restriction.PhiMaxCoverAge)
		memberDataPointResult.PhiSalaryLevel = float64(mpIncomeLevel)
		memberDataPointResult.BasePhiRate = phiRate * (1 + memberDataPointResult.PhiIndustryLoading + memberDataPointResult.PhiRegionLoading)
		memberDataPointResult.LoadedPhiRate = memberDataPointResult.BasePhiRate * (1 + memberDataPointResult.PhiContingencyLoading + memberDataPointResult.PhiVoluntaryLoading + memberDataPointResult.PhiConversionOnWithdrawalLoading)
		memberDataPointResult.ExpAdjLoadedPhiRate = memberDataPointResult.LoadedPhiRate * memberDataPointResult.PhiExperienceAdjustment
	}

	if schemeCategory.CiBenefit {
		ciRate := applyCoverAgeLimit(GetCiRate(&originalMemberDataPointResult, groupParameter, groupQuote, schemeCategory, mpIncomeLevel), memberDataPointResult.AgeNextBirthday, restriction.CiMaxCoverAge)
		if schemeCategory.CiBenefitStructure == "Accelerated" {
			memberDataPointResult.BaseCiRate = ciRate * (1 + memberDataPointResult.CiIndustryLoading + memberDataPointResult.CiRegionLoading - gl1.CiAcceleratedBenefitDiscount)
		} else {
			memberDataPointResult.BaseCiRate = ciRate * (1 + memberDataPointResult.CiIndustryLoading + memberDataPointResult.CiRegionLoading)
		}
		memberDataPointResult.LoadedCiRate = memberDataPointResult.BaseCiRate * (1 + memberDataPointResult.CiContingencyLoading + memberDataPointResult.CiVoluntaryLoading + memberDataPointResult.CiConversionOnWithdrawalLoading)
		memberDataPointResult.ExpAdjLoadedCiRate = memberDataPointResult.LoadedCiRate * memberDataPointResult.CiExperienceAdjustment
	}

	if schemeCategory.SglaBenefit && len(memberDataPointResult.SpouseGender) > 0 {
		spouseIndustryLoading := industryLoadingByGender[strings.ToUpper(memberDataPointResult.SpouseGender[:1])]
		spouseRegionLoading := regionLoadingByGender[strings.ToUpper(memberDataPointResult.SpouseGender[:1])]
		memberDataPointResult.SpouseGlaQx = applyCoverAgeLimit(GetSpouseGlaRate(&originalMemberDataPointResult, groupParameter, mpIncomeLevel, groupQuote, schemeCategory), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.SpouseGlaAidsQx = applyCoverAgeLimit(GetSpouseGlaAidsRate(memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.SpouseGlaLoading = spouseIndustryLoading.GlaIndustryLoadingRate
		memberDataPointResult.BaseSpouseGlaRate = memberDataPointResult.SpouseGlaQx*(1+memberDataPointResult.SpouseGlaLoading+spouseRegionLoading.GlaRegionLoadingRate) + memberDataPointResult.SpouseGlaAidsQx*(1+spouseRegionLoading.GlaAidsRegionLoadingRate)
		memberDataPointResult.LoadedSpouseGlaRate = memberDataPointResult.BaseSpouseGlaRate * (1 + memberDataPointResult.SglaConversionOnWithdrawalLoading)
		memberDataPointResult.ExpAdjLoadedSpouseGlaRate = memberDataPointResult.LoadedSpouseGlaRate * memberDataPointResult.GlaExperienceAdjustment
	}

	// Compute educator loaded rates (includes educator-specific conversion
	// and continuity loadings). Safe to call even when educator isn't enabled —
	// missing loadings are zero, rate equals LoadedGlaRate / LoadedPtdRate.
	computeEducatorLoadedRates(memberDataPointResult)

	memberDataPointResult.GlaRiskPremium = memberDataPointResult.LoadedGlaRate * memberDataPointResult.GlaCappedSumAssured
	memberDataPointResult.ExpAdjGlaRiskPremium = memberDataPointResult.ExpAdjLoadedGlaRate * memberDataPointResult.GlaCappedSumAssured

	// TaxSaver slice of GlaRiskPremium. Per business spec, computed against
	// the already-loaded rate (LoadedGlaRate × TaxSaverLoading × SumAssured),
	// loadedRate × TaxSaverSA × SumAssured
	// over-attributes relative to a pure linear share.
	memberDataPointResult.TaxSaverRiskPremium = memberDataPointResult.LoadedGlaRate * memberDataPointResult.TaxSaverSumAssured
	memberDataPointResult.ExpAdjTaxSaverRiskPremium = memberDataPointResult.ExpAdjLoadedGlaRate * memberDataPointResult.TaxSaverSumAssured

	memberDataPointResult.AdditionalAccidentalGlaRiskPremium = memberDataPointResult.LoadedAdditionalAccidentalGlaRate * memberDataPointResult.AdditionalAccidentalGlaCappedSumAssured
	memberDataPointResult.ExpAdjAdditionalAccidentalGlaRiskPremium = memberDataPointResult.ExpAdjLoadedAdditionalAccidentalGlaRate * memberDataPointResult.AdditionalAccidentalGlaCappedSumAssured

	memberDataPointResult.PtdRiskPremium = memberDataPointResult.LoadedPtdRate * memberDataPointResult.PtdCappedSumAssured
	memberDataPointResult.ExpAdjPtdRiskPremium = memberDataPointResult.ExpAdjLoadedPtdRate * memberDataPointResult.PtdCappedSumAssured

	memberDataPointResult.TtdNumberOfMonthlyPayments = groupParameter.TtdNumberMonthlyPayments
	memberDataPointResult.TtdRiskPremium = memberDataPointResult.LoadedTtdRate * memberDataPointResult.TtdCappedIncome * groupParameter.TtdNumberMonthlyPayments
	memberDataPointResult.ExpAdjTtdRiskPremium = memberDataPointResult.ExpAdjLoadedTtdRate * memberDataPointResult.TtdCappedIncome * groupParameter.TtdNumberMonthlyPayments

	memberDataPointResult.PhiRiskPremium = memberDataPointResult.LoadedPhiRate * memberDataPointResult.PhiMonthlyBenefit
	memberDataPointResult.ExpAdjPhiRiskPremium = memberDataPointResult.ExpAdjLoadedPhiRate * memberDataPointResult.PhiMonthlyBenefit

	memberDataPointResult.CiRiskPremium = memberDataPointResult.LoadedCiRate * memberDataPointResult.CiCappedSumAssured
	memberDataPointResult.ExpAdjCiRiskPremium = memberDataPointResult.ExpAdjLoadedCiRate * memberDataPointResult.CiCappedSumAssured

	memberDataPointResult.SpouseGlaRiskPremium = memberDataPointResult.LoadedSpouseGlaRate * memberDataPointResult.SpouseGlaCappedSumAssured
	memberDataPointResult.ExpAdjSpouseGlaRiskPremium = memberDataPointResult.ExpAdjLoadedSpouseGlaRate * memberDataPointResult.SpouseGlaCappedSumAssured

	// ── Reinsurance rates & loadings ───────────────────────────────────────
	// Compute reinsurance base & loaded rates per benefit using the new
	// reinsurance rate/loading tables (reinsurance_*_rates,
	// reinsurance_general_loadings, reinsurance_industry_loadings,
	// reinsurance_region_loadings). If any reinsurance table is empty for
	// this risk code the getters return zero values, leaving reinsurance
	// rates at 0 so non-reinsured schemes are unaffected.
	reinsIndustryLoadingMain := GetReinsuranceIndustryLoading(groupParameter.RiskRateCode, groupQuote.OccupationClass, memberDataPointResult.Gender[:1])
	memberDataPointResult.ReinsGlaIndustryLoading = reinsIndustryLoadingMain.GlaIndustryLoadingRate
	memberDataPointResult.ReinsPtdIndustryLoading = reinsIndustryLoadingMain.PtdIndustryLoadingRate
	memberDataPointResult.ReinsCiIndustryLoading = reinsIndustryLoadingMain.CiIndustryLoadingRate
	memberDataPointResult.ReinsTtdIndustryLoading = reinsIndustryLoadingMain.TtdIndustryLoadingRate
	memberDataPointResult.ReinsPhiIndustryLoading = reinsIndustryLoadingMain.PhiIndustryLoadingRate

	// Reinsurance region loadings (pre-loaded per-gender map built once per category).
	reinsRegionLoadingMain := reinsRegionLoadingByGender[strings.ToUpper(memberDataPointResult.Gender[:1])]
	memberDataPointResult.ReinsGlaRegionLoading = reinsRegionLoadingMain.GlaRegionLoadingRate
	memberDataPointResult.ReinsGlaAidsRegionLoading = reinsRegionLoadingMain.GlaAidsRegionLoadingRate
	memberDataPointResult.ReinsPtdRegionLoading = reinsRegionLoadingMain.PtdRegionLoadingRate
	memberDataPointResult.ReinsCiRegionLoading = reinsRegionLoadingMain.CiRegionLoadingRate
	memberDataPointResult.ReinsTtdRegionLoading = reinsRegionLoadingMain.TtdRegionLoadingRate
	memberDataPointResult.ReinsPhiRegionLoading = reinsRegionLoadingMain.PhiRegionLoadingRate
	memberDataPointResult.ReinsFunRegionLoading = reinsRegionLoadingMain.FunRegionLoadingRate
	memberDataPointResult.ReinsFunAidsRegionLoading = reinsRegionLoadingMain.FunAidsRegionLoadingRate

	reinsGL := GetReinsuranceGeneralLoading(groupParameter.RiskRateCode, memberDataPointResult.AgeNextBirthday, memberDataPointResult.Gender)
	memberDataPointResult.ReinsGlaContingencyLoading = reinsGL.GlaContigencyLoadingRate
	memberDataPointResult.ReinsPtdContingencyLoading = reinsGL.PtdContigencyLoadingRate
	memberDataPointResult.ReinsCiContingencyLoading = reinsGL.CiContigencyLoadingRate
	memberDataPointResult.ReinsTtdContingencyLoading = reinsGL.TtdContigencyLoadingRate
	memberDataPointResult.ReinsPhiContingencyLoading = reinsGL.PhiContigencyLoadingRate
	memberDataPointResult.ReinsFunContingencyLoading = reinsGL.FunContigencyLoadingRate

	// Voluntary loadings only apply when the quote's obligation type is
	// Voluntary; otherwise they stay zero.
	if groupQuote.ObligationType == "Voluntary" {
		memberDataPointResult.ReinsGlaVoluntaryLoading = reinsGL.GlaVoluntaryLoadingRate
		memberDataPointResult.ReinsPtdVoluntaryLoading = reinsGL.PtdVoluntaryLoadingRate
		memberDataPointResult.ReinsCiVoluntaryLoading = reinsGL.CiVoluntaryLoadingRate
		memberDataPointResult.ReinsTtdVoluntaryLoading = reinsGL.TtdVoluntaryLoadingRate
		memberDataPointResult.ReinsPhiVoluntaryLoading = reinsGL.PhiVoluntaryLoadingRate
		memberDataPointResult.ReinsFunVoluntaryLoading = reinsGL.FunVoluntaryLoadingRate
	}
	if schemeCategory.GlaTerminalIllnessBenefit == "Yes" {
		memberDataPointResult.ReinsGlaTerminalIllnessLoading = reinsGL.TerminalIllnessLoadingRate
	}

	// GLA: BaseReinsGlaRate = ReinsGlaQx × (1 + ReinsGlaIndustryLoading + ReinsGlaRegionLoading)
	//                        + ReinsGlaAidsQx × (1 + ReinsGlaAidsRegionLoading)
	if schemeCategory.GlaBenefit {
		memberDataPointResult.ReinsGlaQx = applyCoverAgeLimit(GetReinsuranceGlaRate(&originalMemberDataPointResult, groupParameter, mpIncomeLevel, schemeCategory), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.ReinsGlaAidsQx = applyCoverAgeLimit(GetReinsuranceGlaAidsRate(memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.BaseReinsGlaRate = memberDataPointResult.ReinsGlaQx*(1+memberDataPointResult.ReinsGlaIndustryLoading+memberDataPointResult.ReinsGlaRegionLoading) +
			memberDataPointResult.ReinsGlaAidsQx*(1+memberDataPointResult.ReinsGlaAidsRegionLoading)
		memberDataPointResult.LoadedReinsGlaRate = memberDataPointResult.BaseReinsGlaRate *
			(1 + memberDataPointResult.ReinsGlaContingencyLoading + memberDataPointResult.ReinsGlaVoluntaryLoading + memberDataPointResult.ReinsGlaTerminalIllnessLoading)
	}

	// PTD
	if schemeCategory.PtdBenefit {
		memberDataPointResult.ReinsPtdRate = applyCoverAgeLimit(GetReinsurancePtdRate(&originalMemberDataPointResult, groupParameter, mpIncomeLevel, schemeCategory), memberDataPointResult.AgeNextBirthday, restriction.PtdMaxCoverAge)
		if schemeCategory.PtdBenefitType == "Accelerated" {
			memberDataPointResult.BaseReinsPtdRate = memberDataPointResult.ReinsPtdRate * (1 + memberDataPointResult.ReinsPtdIndustryLoading + memberDataPointResult.ReinsPtdRegionLoading - reinsGL.PtdAcceleratedBenefitDiscount)
		} else {
			memberDataPointResult.BaseReinsPtdRate = memberDataPointResult.ReinsPtdRate * (1 + memberDataPointResult.ReinsPtdIndustryLoading + memberDataPointResult.ReinsPtdRegionLoading)
		}
		memberDataPointResult.LoadedReinsPtdRate = memberDataPointResult.BaseReinsPtdRate * (1 + memberDataPointResult.ReinsPtdContingencyLoading + memberDataPointResult.ReinsPtdVoluntaryLoading)
	}

	// CI
	if schemeCategory.CiBenefit {
		memberDataPointResult.ReinsCiRate = applyCoverAgeLimit(GetReinsuranceCiRate(&originalMemberDataPointResult, groupParameter, mpIncomeLevel, schemeCategory), memberDataPointResult.AgeNextBirthday, restriction.CiMaxCoverAge)
		if schemeCategory.CiBenefitStructure == "Accelerated" {
			memberDataPointResult.BaseReinsCiRate = memberDataPointResult.ReinsCiRate * (1 + memberDataPointResult.ReinsCiIndustryLoading + memberDataPointResult.ReinsCiRegionLoading - reinsGL.CiAcceleratedBenefitDiscount)
		} else {
			memberDataPointResult.BaseReinsCiRate = memberDataPointResult.ReinsCiRate * (1 + memberDataPointResult.ReinsCiIndustryLoading + memberDataPointResult.ReinsCiRegionLoading)
		}
		memberDataPointResult.LoadedReinsCiRate = memberDataPointResult.BaseReinsCiRate * (1 + memberDataPointResult.ReinsCiContingencyLoading + memberDataPointResult.ReinsCiVoluntaryLoading)
	}

	// TTD — reinsurance TTD rate table is not yet present; base uses direct TTD Qx scaled by reinsurance industry/region loadings.
	if schemeCategory.TtdBenefit {
		ttdReinsQx := memberDataPointResult.BaseTtdRate / math.Max(1+memberDataPointResult.TtdIndustryLoading+memberDataPointResult.TtdRegionLoading, 1e-9)
		memberDataPointResult.BaseReinsTtdRate = ttdReinsQx * (1 + memberDataPointResult.ReinsTtdIndustryLoading + memberDataPointResult.ReinsTtdRegionLoading)
		memberDataPointResult.LoadedReinsTtdRate = memberDataPointResult.BaseReinsTtdRate * (1 + memberDataPointResult.ReinsTtdContingencyLoading + memberDataPointResult.ReinsTtdVoluntaryLoading)
	}

	// PHI
	if schemeCategory.PhiBenefit {
		memberDataPointResult.ReinsPhiRate = applyCoverAgeLimit(GetReinsurancePhiRate(&originalMemberDataPointResult, groupParameter, mpIncomeLevel, schemeCategory), memberDataPointResult.AgeNextBirthday, restriction.PhiMaxCoverAge)
		memberDataPointResult.BaseReinsPhiRate = memberDataPointResult.ReinsPhiRate * (1 + memberDataPointResult.ReinsPhiIndustryLoading + memberDataPointResult.ReinsPhiRegionLoading)
		memberDataPointResult.LoadedReinsPhiRate = memberDataPointResult.BaseReinsPhiRate * (1 + memberDataPointResult.ReinsPhiContingencyLoading + memberDataPointResult.ReinsPhiVoluntaryLoading)
	}

	// Spouse GLA
	if schemeCategory.SglaBenefit && len(memberDataPointResult.SpouseGender) > 0 {
		spouseReinsIndustry := GetReinsuranceIndustryLoading(groupParameter.RiskRateCode, groupQuote.OccupationClass, memberDataPointResult.SpouseGender[:1])
		memberDataPointResult.ReinsSpouseGlaQx = applyCoverAgeLimit(GetReinsuranceSpouseGlaRate(memberDataPointResult, groupParameter, mpIncomeLevel, schemeCategory), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.ReinsSpouseGlaAidsQx = applyCoverAgeLimit(GetReinsuranceSpouseGlaAidsRate(memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.ReinsSpouseGlaLoading = spouseReinsIndustry.GlaIndustryLoadingRate
		spouseReinsRegion := reinsRegionLoadingByGender[strings.ToUpper(memberDataPointResult.SpouseGender[:1])]
		memberDataPointResult.BaseReinsSpouseGlaRate = memberDataPointResult.ReinsSpouseGlaQx*(1+memberDataPointResult.ReinsSpouseGlaLoading+spouseReinsRegion.GlaRegionLoadingRate) +
			memberDataPointResult.ReinsSpouseGlaAidsQx*(1+spouseReinsRegion.GlaAidsRegionLoadingRate)
		memberDataPointResult.LoadedReinsSpouseGlaRate = memberDataPointResult.BaseReinsSpouseGlaRate
	}

	// Family funeral per-relationship reinsurance rates. When the scheme
	// uses the GLA benefit as the main-member funeral rate we reuse
	// LoadedReinsGlaRate; otherwise compute a dedicated reinsurance funeral
	// rate from reinsurance_funeral_rates + reinsurance_funeral_aids_rates.
	if schemeCategory.GlaBenefit {
		memberDataPointResult.MainMemberReinsuranceRate = memberDataPointResult.LoadedReinsGlaRate
		memberDataPointResult.SpouseReinsuranceRate = memberDataPointResult.LoadedReinsSpouseGlaRate
	} else if schemeCategory.FamilyFuneralBenefit {
		reinsFunQx := applyCoverAgeLimit(GetReinsuranceFuneralRate(memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
		reinsFunAidsQx := applyCoverAgeLimit(GetReinsuranceFuneralAidsRate(memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
		memberDataPointResult.MainMemberReinsuranceBaseRate = reinsFunQx*(1+memberDataPointResult.ReinsFunRegionLoading) + reinsFunAidsQx*(1+memberDataPointResult.ReinsFunAidsRegionLoading)
		memberDataPointResult.MainMemberReinsuranceRate = memberDataPointResult.MainMemberReinsuranceBaseRate * (1 + memberDataPointResult.ReinsFunContingencyLoading + memberDataPointResult.ReinsFunVoluntaryLoading)
		if len(memberDataPointResult.SpouseGender) > 0 {
			spouseReinsRegion := reinsRegionLoadingByGender[strings.ToUpper(memberDataPointResult.SpouseGender[:1])]
			spouseReinsFunQx := applyCoverAgeLimit(GetReinsuranceSpouseFuneralRate(memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
			spouseReinsFunAidsQx := applyCoverAgeLimit(GetReinsuranceSpouseFuneralAidsRate(memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
			memberDataPointResult.SpouseReinsuranceBaseRate = spouseReinsFunQx*(1+spouseReinsRegion.FunRegionLoadingRate) + spouseReinsFunAidsQx*(1+spouseReinsRegion.FunAidsRegionLoadingRate)
			memberDataPointResult.SpouseReinsuranceRate = memberDataPointResult.SpouseReinsuranceBaseRate * (1 + memberDataPointResult.ReinsFunContingencyLoading + memberDataPointResult.ReinsFunVoluntaryLoading)
		}
	}
	// Child/parent/dependant — fall back to the direct-pricing rate scaled
	// by the funeral reinsurance contingency loading, since those don't have
	// per-life reinsurance rate tables.
	memberDataPointResult.ChildReinsuranceBaseRate = memberDataPointResult.ChildFuneralBaseRate
	memberDataPointResult.ChildReinsuranceRate = memberDataPointResult.ChildReinsuranceBaseRate * (1 + memberDataPointResult.ReinsFunContingencyLoading + memberDataPointResult.ReinsFunVoluntaryLoading)
	memberDataPointResult.ParentReinsuranceBaseRate = memberDataPointResult.ParentFuneralBaseRate
	memberDataPointResult.ParentReinsuranceRate = memberDataPointResult.ParentReinsuranceBaseRate * (1 + memberDataPointResult.ReinsFunContingencyLoading + memberDataPointResult.ReinsFunVoluntaryLoading)
	memberDataPointResult.DependantReinsuranceBaseRate = memberDataPointResult.ParentFuneralBaseRate
	memberDataPointResult.DependantReinsuranceRate = memberDataPointResult.DependantReinsuranceBaseRate * (1 + memberDataPointResult.ReinsFunContingencyLoading + memberDataPointResult.ReinsFunVoluntaryLoading)

	discountFraction := -(groupQuote.Loadings.Discount / 100.0)

	//memberDataPointResult.MarriageProportion = groupFuneralParameter.ProportionMarried

	memberDataPointResult.ChildFuneralBaseRate = applyCoverAgeLimit(GetChildFuneralRate(&originalMemberDataPointResult, groupParameter, memberDataPointResult.AverageChildAgeNextBirthday), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
	memberDataPointResult.ChildFuneralSumAssured = schemeCategory.FamilyFuneralChildrenFuneralSumAssured
	memberDataPointResult.ParentFuneralBaseRate = applyCoverAgeLimit(GetDependantMortalityRate(&originalMemberDataPointResult, groupParameter, memberDataPointResult.AverageDependantAgeNextBirthday, mpIncomeLevel), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
	memberDataPointResult.ParentFuneralSumAssured = schemeCategory.FamilyFuneralAdultDependantSumAssured

	mainMemberFuneralSA := applyMaxCoverCap(applyCoverAgeLimit(schemeCategory.FamilyFuneralMainMemberFuneralSumAssured, memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge), reinsCoverCaps[benefitTypeKey(schemeCategory.FamilyFuneralAlias, models.BenefitTypeFun)])
	if schemeCategory.GlaBenefit {
		memberDataPointResult.MainMemberFuneralRiskPremium = memberDataPointResult.LoadedGlaRate * mainMemberFuneralSA
		memberDataPointResult.SpouseFuneralRiskPremium = memberDataPointResult.LoadedSpouseGlaRate * schemeCategory.FamilyFuneralSpouseFuneralSumAssured
	} else {
		memberDataPointResult.MainMemberFuneralRiskPremium = memberDataPointResult.MainMemberFuneralBaseRate * mainMemberFuneralSA
		memberDataPointResult.SpouseFuneralRiskPremium = memberDataPointResult.SpouseFuneralBaseRate * schemeCategory.FamilyFuneralSpouseFuneralSumAssured
	}
	memberDataPointResult.ChildFuneralRiskPremium = memberDataPointResult.ChildFuneralBaseRate * schemeCategory.FamilyFuneralChildrenFuneralSumAssured * math.Min(memberDataPointResult.AverageNumberChildren, float64(schemeCategory.FamilyFuneralMaxNumberChildren))
	memberDataPointResult.ParentFuneralRiskPremium = memberDataPointResult.ParentFuneralBaseRate * schemeCategory.FamilyFuneralAdultDependantSumAssured * memberDataPointResult.AverageNumberDependants

	memberDataPointResult.TotalFuneralRiskPremium = (memberDataPointResult.MainMemberFuneralRiskPremium + memberDataPointResult.SpouseFuneralRiskPremium + memberDataPointResult.ChildFuneralRiskPremium + memberDataPointResult.ParentFuneralRiskPremium) * (1 + memberDataPointResult.FunConversionOnWithdrawalLoading)
	memberDataPointResult.ExpAdjTotalFuneralRiskPremium = memberDataPointResult.GlaExperienceAdjustment * (memberDataPointResult.MainMemberFuneralRiskPremium + memberDataPointResult.SpouseFuneralRiskPremium + memberDataPointResult.ChildFuneralRiskPremium + memberDataPointResult.ParentFuneralRiskPremium) * (1 + memberDataPointResult.FunConversionOnWithdrawalLoading)
	memberDataPointResult.TotalFuneralOfficePremium = memberDataPointResult.TotalFuneralRiskPremium / (1.0 - memberDataPointResult.TotalPremiumLoading)
	memberDataPointResult.ExpAdjTotalFuneralOfficePremium = memberDataPointResult.ExpAdjTotalFuneralRiskPremium / (1.0 - memberDataPointResult.TotalPremiumLoading)
	memberDataPointResult.FinalTotalFuneralOfficePremium = memberDataPointResult.ExpAdjTotalFuneralRiskPremium / (1.0 - (memberDataPointResult.TotalPremiumLoading + discountFraction))

	// Compute all non-educator conversion / continuity slice premiums now
	// that Loaded*Rates and TotalFuneralRiskCost are final. Educator slice
	// premiums are computed inline in the educator block (they need
	// riskWeightedSA).
	computeConvContSlicePremiums(memberDataPointResult, groupParameter)

	applyBinderOutsourceAmounts(memberDataPointResult, &groupQuote, binderFeeRate, outsourceFeeRate)

	if memberDataPointResult.AgeNextBirthday > groupQuote.NormalRetirementAge {
		memberDataPointResult.ExceedsNormalRetirementAgeIndicator = 1
	}
	if memberDataPointResult.GlaSumAssured > groupQuote.FreeCoverLimit {
		memberDataPointResult.ExceedsFreeCoverLimitIndicator = 1
	}

	//memberDataPointResult.GlaExperienceAdjustedAnnualPremium = memberDataPointResult.GlaRiskPremium
	//memberDataPointResult.PtdExperienceAdjustedAnnualPremium = memberDataPointResult.PtdRiskPremium
	//memberDataPointResult.TtdExperienceAdjustedAnnualPremium = memberDataPointResult.TtdRiskPremium
	//memberDataPointResult.PhiExperienceAdjustedAnnualPremium = memberDataPointResult.PhiRiskPremium
	//memberDataPointResult.CiExperienceAdjustedAnnualPremium = memberDataPointResult.CiRiskPremium
	//memberDataPointResult.SpouseExperienceAdjustedAnnualPremium = memberDataPointResult.SpouseGlaRiskPremium
	//memberDataPointResult.FuneralExperienceAdjustedAnnualPremium = memberDataPointResult.TotalFuneralRiskCost
	memberPremiumScheduleDatapoint.IsOriginalMember = false
	memberPremiumScheduleDatapoint.GlaCoveredSumAssured = memberDataPointResult.GlaCappedSumAssured
	memberPremiumScheduleDatapoint.GlaAnnualPremium = computeMemberOfficePremium(memberDataPointResult.ExpAdjGlaRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.PtdCoveredSumAssured = memberDataPointResult.PtdCappedSumAssured
	memberPremiumScheduleDatapoint.PtdAnnualPremium = computeMemberOfficePremium(memberDataPointResult.ExpAdjPtdRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.CiCoveredSumAssured = memberDataPointResult.CiCappedSumAssured
	memberPremiumScheduleDatapoint.CiAnnualPremium = computeMemberOfficePremium(memberDataPointResult.ExpAdjCiRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.TtdCoveredIncome = memberDataPointResult.TtdCappedIncome
	memberPremiumScheduleDatapoint.TtdAnnualPremium = computeMemberOfficePremium(memberDataPointResult.ExpAdjTtdRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.PhiCoveredIncome = memberDataPointResult.PhiCappedIncome
	memberPremiumScheduleDatapoint.PhiAnnualPremium = computeMemberOfficePremium(memberDataPointResult.ExpAdjPhiRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.SpouseGlaCoveredSumAssured = memberDataPointResult.SpouseGlaCappedSumAssured
	memberPremiumScheduleDatapoint.SpouseGlaAnnualPremium = computeMemberOfficePremium(memberDataPointResult.ExpAdjSpouseGlaRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.MainMemberFuneralSumAssured = schemeCategory.FamilyFuneralMainMemberFuneralSumAssured
	memberPremiumScheduleDatapoint.MainMemberFuneralAnnualPremium = computeMemberOfficePremium(memberDataPointResult.MainMemberFuneralRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.SpouseFuneralSumAssured = schemeCategory.FamilyFuneralSpouseFuneralSumAssured
	memberPremiumScheduleDatapoint.SpouseFuneralAnnualPremium = computeMemberOfficePremium(memberDataPointResult.SpouseFuneralRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.ChildFuneralSumAssured = memberDataPointResult.ChildFuneralSumAssured
	memberPremiumScheduleDatapoint.ChildrenFuneralAnnualPremium = computeMemberOfficePremium(memberDataPointResult.ChildFuneralRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.DependantsFuneralSumAssured = memberDataPointResult.ParentFuneralSumAssured
	memberPremiumScheduleDatapoint.DependantsFuneralAnnualPremium = computeMemberOfficePremium(memberDataPointResult.ParentFuneralRiskPremium, &groupQuote)

	memberPremiumScheduleDatapoint.TotalAnnualPremiumPayable = memberPremiumScheduleDatapoint.GlaAnnualPremium + memberPremiumScheduleDatapoint.PtdAnnualPremium + memberPremiumScheduleDatapoint.CiAnnualPremium + memberPremiumScheduleDatapoint.TtdAnnualPremium + memberPremiumScheduleDatapoint.PhiAnnualPremium + memberPremiumScheduleDatapoint.SpouseGlaAnnualPremium + memberPremiumScheduleDatapoint.MainMemberFuneralAnnualPremium + memberPremiumScheduleDatapoint.SpouseFuneralAnnualPremium + memberPremiumScheduleDatapoint.ChildrenFuneralAnnualPremium + memberPremiumScheduleDatapoint.DependantsFuneralAnnualPremium

	GroupPricingReinsurance(&originalMemberDataPointResult, &bordereauxDatapoint, groupQuote, schemeCategory, groupPricingReinsuranceStructure, groupParameter)

	bordereauxDatapoint.IsOriginalMember = false
	bordereauxDatapoint.SchemeId = memberDataPointResult.SchemeId
	bordereauxDatapoint.QuoteId = memberDataPointResult.QuoteId
	bordereauxDatapoint.Gender = memberDataPointResult.Gender
	bordereauxDatapoint.AgeNextBirthday = float64(memberDataPointResult.AgeNextBirthday)
	bordereauxDatapoint.AnnualSalary = memberDataPointResult.AnnualSalary
	bordereauxDatapoint.RenewalDate = "" //groupQuote.CommencementDate
	bordereauxDatapoint.Currency = groupQuote.Currency
	bordereauxDatapoint.Industry = groupQuote.Industry
	bordereauxDatapoint.IndustryClass = groupQuote.Industry //
	bordereauxDatapoint.GlaMultiple = schemeCategory.GlaSalaryMultiple
	bordereauxDatapoint.GlaCoveredSumAssured = memberDataPointResult.GlaCappedSumAssured
	bordereauxDatapoint.LoadedGlaRiskRate = memberDataPointResult.LoadedGlaRate
	bordereauxDatapoint.GlaRetainedRiskPremium = bordereauxDatapoint.GlaRetainedSumAssured * memberDataPointResult.LoadedGlaRate
	// Ceded premium uses the LOADED REINSURANCE rate × ceded sum assured.
	bordereauxDatapoint.GlaCededRiskPremium = bordereauxDatapoint.GlaCededSumAssured * memberDataPointResult.LoadedReinsGlaRate
	memberDataPointResult.GlaReinsurancePremium = bordereauxDatapoint.GlaCededRiskPremium

	bordereauxDatapoint.PtdMultiple = schemeCategory.PtdSalaryMultiple
	bordereauxDatapoint.PtdCoveredSumAssured = memberDataPointResult.PtdCappedSumAssured
	bordereauxDatapoint.LoadedPtdRiskRate = memberDataPointResult.LoadedPtdRate
	bordereauxDatapoint.PtdRetainedRiskPremium = bordereauxDatapoint.PtdRetainedSumAssured * memberDataPointResult.LoadedPtdRate
	bordereauxDatapoint.PtdCededRiskPremium = bordereauxDatapoint.PtdCededSumAssured * memberDataPointResult.LoadedReinsPtdRate
	memberDataPointResult.PtdReinsurancePremium = bordereauxDatapoint.PtdCededRiskPremium

	bordereauxDatapoint.CiMultiple = schemeCategory.CiCriticalIllnessSalaryMultiple
	bordereauxDatapoint.CiCoveredSumAssured = memberDataPointResult.CiCappedSumAssured
	bordereauxDatapoint.LoadedCiRiskRate = memberDataPointResult.LoadedCiRate
	bordereauxDatapoint.CiRetainedRiskPremium = bordereauxDatapoint.CiRetainedSumAssured * memberDataPointResult.LoadedCiRate
	bordereauxDatapoint.CiCededRiskPremium = bordereauxDatapoint.CiCededSumAssured * memberDataPointResult.LoadedReinsCiRate
	memberDataPointResult.CiReinsurancePremium = bordereauxDatapoint.CiCededRiskPremium

	bordereauxDatapoint.SglaMultiple = schemeCategory.SglaSalaryMultiple
	bordereauxDatapoint.SglaCoveredSumAssured = memberDataPointResult.SpouseGlaCappedSumAssured
	bordereauxDatapoint.LoadedSglaRiskRate = memberDataPointResult.LoadedSpouseGlaRate
	bordereauxDatapoint.SglaRetainedRiskPremium = bordereauxDatapoint.SglaRetainedSumAssured * memberDataPointResult.LoadedSpouseGlaRate
	bordereauxDatapoint.SglaCededRiskPremium = bordereauxDatapoint.SglaCededSumAssured * memberDataPointResult.LoadedReinsSpouseGlaRate
	memberDataPointResult.SpouseGlaReinsurancePremium = bordereauxDatapoint.SglaCededRiskPremium

	bordereauxDatapoint.TtdReplacementMultiple = schemeCategory.TtdIncomeReplacementPercentage
	bordereauxDatapoint.TtdMonthlyBenefit = memberDataPointResult.TtdCappedIncome
	bordereauxDatapoint.LoadedTtdRiskRate = memberDataPointResult.LoadedTtdRate
	bordereauxDatapoint.TtdRetainedRiskPremium = bordereauxDatapoint.TtdRetainedMonthlyBenefit * memberDataPointResult.LoadedTtdRate * groupParameter.TtdNumberMonthlyPayments
	bordereauxDatapoint.TtdCededRiskPremium = bordereauxDatapoint.TtdCededMonthlyBenefit * memberDataPointResult.LoadedReinsTtdRate * groupParameter.TtdNumberMonthlyPayments
	memberDataPointResult.TtdReinsurancePremium = bordereauxDatapoint.TtdCededRiskPremium

	bordereauxDatapoint.PhiReplacementMultiple = schemeCategory.PhiIncomeReplacementPercentage
	bordereauxDatapoint.PhiMonthlyBenefit = memberDataPointResult.PhiCappedIncome
	bordereauxDatapoint.LoadedPhiRiskRate = memberDataPointResult.LoadedPhiRate
	bordereauxDatapoint.PhiRetainedRiskPremium = bordereauxDatapoint.PhiRetainedMonthlyBenefit * memberDataPointResult.LoadedPhiRate
	bordereauxDatapoint.PhiCededRiskPremium = bordereauxDatapoint.PhiCededMonthlyBenefit * memberDataPointResult.LoadedReinsPhiRate
	memberDataPointResult.PhiReinsurancePremium = bordereauxDatapoint.PhiCededRiskPremium

	bordereauxDatapoint.MainMemberFuneralSumAssured = schemeCategory.FamilyFuneralMainMemberFuneralSumAssured
	bordereauxDatapoint.MainMemberRiskRate = memberDataPointResult.LoadedGlaRate
	bordereauxDatapoint.MainMemberRetainedRiskPremium = bordereauxDatapoint.MainMemberRetainedSumAssured * memberDataPointResult.LoadedGlaRate
	bordereauxDatapoint.MainMemberCededRiskPremium = bordereauxDatapoint.MainMemberCededSumAssured * memberDataPointResult.MainMemberReinsuranceRate
	memberDataPointResult.MainMemberReinsurancePremium = bordereauxDatapoint.MainMemberCededRiskPremium

	bordereauxDatapoint.SpouseFuneralSumAssured = schemeCategory.FamilyFuneralSpouseFuneralSumAssured
	bordereauxDatapoint.SpouseRiskRate = memberDataPointResult.LoadedSpouseGlaRate
	bordereauxDatapoint.SpouseRetainedRiskPremium = bordereauxDatapoint.SpouseRetainedSumAssured * memberDataPointResult.LoadedSpouseGlaRate
	bordereauxDatapoint.SpouseCededRiskPremium = bordereauxDatapoint.SpouseCededSumAssured * memberDataPointResult.SpouseReinsuranceRate
	memberDataPointResult.SpouseReinsurancePremium = bordereauxDatapoint.SpouseCededRiskPremium

	bordereauxDatapoint.ChildFuneralSumAssured = schemeCategory.FamilyFuneralChildrenFuneralSumAssured
	bordereauxDatapoint.ChildRiskRate = memberDataPointResult.ChildFuneralBaseRate
	bordereauxDatapoint.ChildRetainedRiskPremium = bordereauxDatapoint.ChildRetainedSumAssured * memberDataPointResult.ChildFuneralBaseRate
	bordereauxDatapoint.ChildCededRiskPremium = bordereauxDatapoint.ChildCededSumAssured * memberDataPointResult.ChildReinsuranceRate
	memberDataPointResult.ChildReinsurancePremium = bordereauxDatapoint.ChildCededRiskPremium

	bordereauxDatapoint.ParentFuneralSumAssured = schemeCategory.FamilyFuneralParentFuneralSumAssured
	bordereauxDatapoint.ParentRiskRate = memberDataPointResult.ParentFuneralBaseRate
	bordereauxDatapoint.ParentRetainedRiskPremium = bordereauxDatapoint.ParentRetainedSumAssured * memberDataPointResult.ParentFuneralBaseRate
	bordereauxDatapoint.ParentCededRiskPremium = bordereauxDatapoint.ParentCededSumAssured * memberDataPointResult.ParentReinsuranceRate
	memberDataPointResult.ParentReinsurancePremium = bordereauxDatapoint.ParentCededRiskPremium

	bordereauxDatapoint.DependantFuneralSumAssured = schemeCategory.FamilyFuneralAdultDependantSumAssured
	bordereauxDatapoint.DependantRiskRate = memberDataPointResult.ParentFuneralBaseRate
	bordereauxDatapoint.DependantRetainedRiskPremium = bordereauxDatapoint.DependantRetainedSumAssured * memberDataPointResult.ParentFuneralBaseRate
	bordereauxDatapoint.DependantCededRiskPremium = bordereauxDatapoint.DependantCededSumAssured * memberDataPointResult.DependantReinsuranceRate
	memberDataPointResult.DependantReinsurancePremium = bordereauxDatapoint.DependantCededRiskPremium

	//UpdateRatedMemberData(mpIndex, memberDataPointResult, memberDataResults)
}

func PopulateRatesPerMemberForExperienceRating(i int, indicativeRatesCount float64, mp models.GPricingMemberData,
	groupQuote models.GroupPricingQuote,
	groupParameter models.GroupPricingParameters, incomeLevels []models.IncomeLevel, ageBands []models.GroupPricingAgeBands, calculatedFreeCoverLimit float64,
	insurerYearEndMonth int, restriction models.Restriction, reinsCoverCaps map[string]float64,
	taxTable []models.TaxTable, tieredIncomeTiers []models.TieredIncomeReplacement, customTieredIncomeTiers []models.TieredIncomeReplacement, premiumLoading models.PremiumLoading, regionLoadingByGender map[string]models.RegionLoading, industryLoadingByGender map[string]models.IndustryLoading) TheoreticalRiskTotal {

	var memberDataPointResult models.MemberRatingResult
	var memberPremiumScheduleDatapoint models.MemberPremiumSchedule
	var mpIncomeLevel int
	var mpAgeBand string
	var groupFuneralParameter models.FuneralParameters
	var binderFeeRate2, outsourceFeeRate2 float64

	memberDataPointResult.IsOriginalMember = true
	memberDataPointResult.QuoteId = groupQuote.ID
	_, memberDataPointResult.FinancialYear = getGroupRiskQuotingFinancialYear(groupQuote.CommencementDate, insurerYearEndMonth)

	memberDataPointResult.SchemeId = mp.SchemeId
	memberDataPointResult.MemberName = mp.MemberName
	memberDataPointResult.Gender = mp.Gender
	memberDataPointResult.DateOfBirth = mp.DateOfBirth
	memberDataPointResult.AnnualSalary = mp.AnnualSalary
	// Look up region loading for this member by gender using the category-level pre-loaded map
	memberRegionLoadingExp := regionLoadingByGender[strings.ToUpper(memberDataPointResult.Gender[:1])]
	if !groupQuote.UseGlobalSalaryMultiple {
		memberDataPointResult.GlaSalaryMultiple = mp.Benefits.GlaMultiple
		if groupQuote.SchemeCategories[i].PtdBenefit {
			memberDataPointResult.PtdSalaryMultiple = mp.Benefits.PtdMultiple
		}
		if groupQuote.SchemeCategories[i].SglaBenefit {
			memberDataPointResult.SglaSalaryMultiple = mp.Benefits.SglaMultiple
		}
		if groupQuote.SchemeCategories[i].CiBenefit {
			memberDataPointResult.CiSalaryMultiple = mp.Benefits.CiMultiple
		}
	}
	if groupQuote.UseGlobalSalaryMultiple {
		memberDataPointResult.GlaSalaryMultiple = groupQuote.SchemeCategories[i].GlaSalaryMultiple
		memberDataPointResult.PtdSalaryMultiple = groupQuote.SchemeCategories[i].PtdSalaryMultiple
		memberDataPointResult.SglaSalaryMultiple = groupQuote.SchemeCategories[i].SglaSalaryMultiple
		memberDataPointResult.CiSalaryMultiple = groupQuote.SchemeCategories[i].CiCriticalIllnessSalaryMultiple
	}
	memberDataPointResult.Occupation = groupQuote.Industry
	memberDataPointResult.OccupationClass = groupQuote.OccupationClass
	memberDataPointResult.Industry = groupQuote.Industry
	memberDataPointResult.ExpenseLoading = premiumLoading.ExpenseLoading
	memberDataPointResult.AdminLoading = premiumLoading.AdminLoading
	effectiveCommission2 := premiumLoading.CommissionLoading
	if groupQuote.DistributionChannel == models.ChannelDirect {
		effectiveCommission2 = 0
	}
	memberDataPointResult.CommissionLoading = effectiveCommission2
	memberDataPointResult.ProfitLoading = premiumLoading.ProfitLoading
	memberDataPointResult.OtherLoading = premiumLoading.OtherLoading
	memberDataPointResult.Discount = -(groupQuote.Loadings.Discount / 100.0)

	if groupQuote.DistributionChannel == "binder" {
		binderFeeRate2, outsourceFeeRate2 = binderAndOutsourceRates(&groupQuote)
		memberDataPointResult.BinderFeeRate = binderFeeRate2
		memberDataPointResult.OutsourceFeeRate = outsourceFeeRate2
	}

	// Commission is excluded from TotalLoading — see applySchemeWideCommission.
	memberDataPointResult.TotalPremiumLoading = math.Max(premiumLoading.ExpenseLoading+premiumLoading.AdminLoading+premiumLoading.ProfitLoading+premiumLoading.OtherLoading+binderFeeRate2+outsourceFeeRate2, premiumLoading.MinimumPremiumLoading)

	memberDataPointResult.CalculatedFreeCoverLimit = calculatedFreeCoverLimit
	if groupQuote.FreeCoverLimit > 0 {
		memberDataPointResult.AppliedFreeCoverLimit = groupQuote.FreeCoverLimit
	}
	if groupQuote.FreeCoverLimit == 0 {
		memberDataPointResult.AppliedFreeCoverLimit = calculatedFreeCoverLimit
	}

	memberPremiumScheduleDatapoint.MemberName = memberDataPointResult.MemberName
	memberPremiumScheduleDatapoint.QuoteId = memberDataPointResult.QuoteId
	memberPremiumScheduleDatapoint.Gender = memberDataPointResult.Gender
	memberPremiumScheduleDatapoint.SchemeId = memberDataPointResult.SchemeId

	//dob, err := utils.ParseDateString(mp.DateOfBirth)
	//if err != nil {
	//	fmt.Println("error encountered parsing date of birth: ", err)
	//	//return err
	//}
	memberDataPointResult.AgeNextBirthday = calculateAgeNextBirthday(groupQuote.CommencementDate, mp.DateOfBirth)
	groupFuneralParameter = GetFuneralParameter(groupParameter, memberDataPointResult.AgeNextBirthday)
	mpAgeBand = GetAgeBand(memberDataPointResult.AgeNextBirthday, ageBands)
	memberDataPointResult.AgeBand = mpAgeBand
	if mp.Gender == "M" {
		memberDataPointResult.SpouseGender = "F"
		memberDataPointResult.SpouseAgeNextBirthday = int(math.Min(math.Max(float64(memberDataPointResult.AgeNextBirthday)-
			float64(groupParameter.SpouseAgeGap), float64(restriction.MinEntryAge)), float64(restriction.MaxEntryAge)))

	}
	if mp.Gender == "F" {
		memberDataPointResult.SpouseGender = "M"
		memberDataPointResult.SpouseAgeNextBirthday = int(math.Min(math.Max(float64(memberDataPointResult.AgeNextBirthday)+
			float64(groupParameter.SpouseAgeGap), float64(restriction.MinEntryAge)), float64(restriction.MaxEntryAge)))
	}

	mpIncomeLevel = GetIncomeLevel(mp, incomeLevels)

	memberDataPointResult.IncomeLevel = mpIncomeLevel

	memberDataPointResult.GlaExperienceAdjustment = 1
	memberDataPointResult.PtdExperienceAdjustment = 1
	memberDataPointResult.CiExperienceAdjustment = 1
	memberDataPointResult.PhiExperienceAdjustment = 1
	memberDataPointResult.TtdExperienceAdjustment = 1

	memberDataPointResult.AverageDependantAgeNextBirthday = groupFuneralParameter.AverageDependantAge
	memberDataPointResult.AverageChildAgeNextBirthday = groupFuneralParameter.AverageChildAge
	memberDataPointResult.AverageNumberDependants = groupFuneralParameter.NumberDependants
	memberDataPointResult.AverageNumberChildren = groupFuneralParameter.NumberChildren
	memberDataPointResult.GlaSumAssured = mp.AnnualSalary * memberDataPointResult.GlaSalaryMultiple * indicativeRatesCount
	unScaledGlaCappedSumAssured := applyMaxCoverCap(applyMaxCoverCap(math.Min(mp.AnnualSalary*memberDataPointResult.GlaSalaryMultiple, memberDataPointResult.AppliedFreeCoverLimit), restriction.MaximumGlaCover), reinsCoverCaps[benefitTypeKey(groupQuote.SchemeCategories[i].GlaAlias, models.BenefitTypeGla)])
	memberDataPointResult.GlaCappedSumAssured = unScaledGlaCappedSumAssured * indicativeRatesCount
	memberDataPointResult.PtdSumAssured = mp.AnnualSalary * memberDataPointResult.PtdSalaryMultiple * indicativeRatesCount
	memberDataPointResult.PtdCappedSumAssured = applyMaxCoverCap(applyMaxCoverCap(math.Min(mp.AnnualSalary*memberDataPointResult.PtdSalaryMultiple, unScaledGlaCappedSumAssured), restriction.MaximumPtdCover), reinsCoverCaps[benefitTypeKey(groupQuote.SchemeCategories[i].PtdAlias, models.BenefitTypePtd)]) * indicativeRatesCount
	memberDataPointResult.CiSumAssured = mp.AnnualSalary * memberDataPointResult.CiSalaryMultiple * indicativeRatesCount
	memberDataPointResult.CiCappedSumAssured = applyMaxCoverCap(math.Min(math.Min(mp.AnnualSalary*memberDataPointResult.CiSalaryMultiple, restriction.SevereIllnessMaximumBenefit), unScaledGlaCappedSumAssured), reinsCoverCaps[benefitTypeKey(groupQuote.SchemeCategories[i].CiAlias, models.BenefitTypeCi)]) * indicativeRatesCount
	memberDataPointResult.SpouseGlaSumAssured = mp.AnnualSalary * memberDataPointResult.SglaSalaryMultiple * indicativeRatesCount
	memberDataPointResult.SpouseGlaCappedSumAssured = applyMaxCoverCap(math.Min(math.Min(mp.AnnualSalary*memberDataPointResult.SglaSalaryMultiple, restriction.SpouseGlaMaximumBenefit), unScaledGlaCappedSumAssured), reinsCoverCaps[benefitTypeKey(groupQuote.SchemeCategories[i].SglaAlias, models.BenefitTypeSgla)]) * indicativeRatesCount

	takeHomePay, _ := computeTakeHomePayFromBands(mp.AnnualSalary, taxTable)

	memberDataPointResult.TtdIncome = mp.AnnualSalary * groupQuote.SchemeCategories[i].TtdIncomeReplacementPercentage * indicativeRatesCount / 12.0 / 100.0
	var unScaledTtdIncome float64
	if groupQuote.SchemeCategories[i].TtdUseTieredIncomeReplacementRatio {
		ttdTiers := tieredIncomeTiers
		if groupQuote.SchemeCategories[i].TtdTieredIncomeReplacementType == "custom" {
			ttdTiers = customTieredIncomeTiers
		}
		unScaledTtdIncome = math.Min(takeHomePay, computeCoveredIncomeFromTiers(mp.AnnualSalary, ttdTiers)) / 12.0
	} else {
		unScaledTtdIncome = math.Min(takeHomePay, mp.AnnualSalary*groupQuote.SchemeCategories[i].TtdIncomeReplacementPercentage/100.0) / 12.0
	}
	memberDataPointResult.TtdCappedIncome = applyMaxCoverCap(math.Min(unScaledTtdIncome, restriction.TtdMaximumMonthlyBenefit), reinsCoverCaps[benefitTypeKey(groupQuote.SchemeCategories[i].TtdAlias, models.BenefitTypeTtd)]) * indicativeRatesCount

	memberDataPointResult.PhiIncome = mp.AnnualSalary * groupQuote.SchemeCategories[i].PhiIncomeReplacementPercentage * indicativeRatesCount / 12.0 / 100.0
	var unScaledPhiIncome float64
	if groupQuote.SchemeCategories[i].PhiUseTieredIncomeReplacementRatio {
		phiTiers := tieredIncomeTiers
		if groupQuote.SchemeCategories[i].PhiTieredIncomeReplacementType == "custom" {
			phiTiers = customTieredIncomeTiers
		}
		unScaledPhiIncome = math.Min(takeHomePay, computeCoveredIncomeFromTiers(mp.AnnualSalary, phiTiers)) / 12.0
	} else {
		unScaledPhiIncome = math.Min(takeHomePay, mp.AnnualSalary*groupQuote.SchemeCategories[i].PhiIncomeReplacementPercentage/100.0) / 12.0
	}
	memberDataPointResult.PhiCappedIncome = applyMaxCoverCap(math.Min(unScaledPhiIncome, restriction.PhiMaximumMonthlyBenefit), reinsCoverCaps[benefitTypeKey(groupQuote.SchemeCategories[i].PhiAlias, models.BenefitTypePhi)]) * indicativeRatesCount

	mpAge := memberDataPointResult.AgeNextBirthday
	memberDataPointResult.GlaCappedSumAssured = applyCoverAgeLimit(memberDataPointResult.GlaCappedSumAssured, mpAge, restriction.GlaMaxCoverAge)
	memberDataPointResult.PtdCappedSumAssured = applyCoverAgeLimit(memberDataPointResult.PtdCappedSumAssured, mpAge, restriction.PtdMaxCoverAge)
	memberDataPointResult.CiCappedSumAssured = applyCoverAgeLimit(memberDataPointResult.CiCappedSumAssured, mpAge, restriction.CiMaxCoverAge)
	memberDataPointResult.SpouseGlaCappedSumAssured = applyCoverAgeLimit(memberDataPointResult.SpouseGlaCappedSumAssured, mpAge, restriction.GlaMaxCoverAge)
	memberDataPointResult.TtdCappedIncome = applyCoverAgeLimit(memberDataPointResult.TtdCappedIncome, mpAge, restriction.TtdMaxCoverAge)
	memberDataPointResult.PhiCappedIncome = applyCoverAgeLimit(memberDataPointResult.PhiCappedIncome, mpAge, restriction.PhiMaxCoverAge)

	memberDataPointResult.IncomeReplacementRatio = groupQuote.SchemeCategories[i].PhiIncomeReplacementPercentage / 100.0
	memberDataPointResult.PhiContributionWaiver = math.Min(mp.AnnualSalary*mp.ContributionWaiverProportion, restriction.PhiMaximumMonthlyContributionWaiver) * indicativeRatesCount
	memberDataPointResult.PhiMonthlyBenefit = memberDataPointResult.PhiCappedIncome + memberDataPointResult.PhiContributionWaiver

	memberIndustryLoading := industryLoadingByGender[strings.ToUpper(memberDataPointResult.Gender[:1])]

	// Persist the per-member region and industry loadings so formulas below
	// can read them directly from MemberRatingResult.
	memberDataPointResult.GlaRegionLoading = memberRegionLoadingExp.GlaRegionLoadingRate
	memberDataPointResult.GlaAidsRegionLoading = memberRegionLoadingExp.GlaAidsRegionLoadingRate
	memberDataPointResult.PtdRegionLoading = memberRegionLoadingExp.PtdRegionLoadingRate
	memberDataPointResult.CiRegionLoading = memberRegionLoadingExp.CiRegionLoadingRate
	memberDataPointResult.TtdRegionLoading = memberRegionLoadingExp.TtdRegionLoadingRate
	memberDataPointResult.PhiRegionLoading = memberRegionLoadingExp.PhiRegionLoadingRate
	memberDataPointResult.FunRegionLoading = memberRegionLoadingExp.FunRegionLoadingRate
	memberDataPointResult.FunAidsRegionLoading = memberRegionLoadingExp.FunAidsRegionLoadingRate

	memberDataPointResult.GlaIndustryLoading = memberIndustryLoading.GlaIndustryLoadingRate
	memberDataPointResult.PtdIndustryLoading = memberIndustryLoading.PtdIndustryLoadingRate
	memberDataPointResult.CiIndustryLoading = memberIndustryLoading.CiIndustryLoadingRate
	memberDataPointResult.TtdIndustryLoading = memberIndustryLoading.TtdIndustryLoadingRate
	memberDataPointResult.PhiIndustryLoading = memberIndustryLoading.PhiIndustryLoadingRate

	memberDataPointResult.GlaQx = applyCoverAgeLimit(GetGlaRate(&memberDataPointResult, groupParameter, mpIncomeLevel, groupQuote, groupQuote.SchemeCategories[i]), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
	memberDataPointResult.GlaAidsQx = applyCoverAgeLimit(GetGlaAidsRate(&memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
	memberDataPointResult.BaseGlaRate = memberDataPointResult.GlaQx*(1+memberDataPointResult.GlaIndustryLoading+memberDataPointResult.GlaRegionLoading) + memberDataPointResult.GlaAidsQx*(1+memberDataPointResult.GlaAidsRegionLoading)

	gl2 := GetGeneralLoading(groupParameter.RiskRateCode, memberDataPointResult.AgeNextBirthday, memberDataPointResult.Gender)
	if groupQuote.SchemeCategories[i].GlaTerminalIllnessBenefit == "Yes" {
		memberDataPointResult.GlaTerminalIllnessLoading = gl2.TerminalIllnessLoadingRate
	}

	// Persist the per-member contingency loadings alongside the already-set
	// region/industry loadings.
	memberDataPointResult.GlaContingencyLoading = gl2.GlaContigencyLoadingRate
	memberDataPointResult.PtdContingencyLoading = gl2.PtdContigencyLoadingRate
	memberDataPointResult.CiContingencyLoading = gl2.CiContigencyLoadingRate
	memberDataPointResult.TtdContingencyLoading = gl2.TtdContigencyLoadingRate
	memberDataPointResult.PhiContingencyLoading = gl2.PhiContigencyLoadingRate
	memberDataPointResult.FunContingencyLoading = gl2.FunContigencyLoadingRate

	// Voluntary loadings only apply when the quote's obligation type is
	// Voluntary; otherwise they stay zero.
	if groupQuote.ObligationType == "Voluntary" {
		memberDataPointResult.GlaVoluntaryLoading = gl2.GlaVoluntaryLoadingRate
		memberDataPointResult.PtdVoluntaryLoading = gl2.PtdVoluntaryLoadingRate
		memberDataPointResult.CiVoluntaryLoading = gl2.CiVoluntaryLoadingRate
		memberDataPointResult.TtdVoluntaryLoading = gl2.TtdVoluntaryLoadingRate
		memberDataPointResult.PhiVoluntaryLoading = gl2.PhiVoluntaryLoadingRate
		memberDataPointResult.FunVoluntaryLoading = gl2.FunVoluntaryLoadingRate
	}

	if groupQuote.SchemeCategories[i].GlaBenefit && groupQuote.SchemeCategories[i].TaxSaverBenefit {
		memberDataPointResult.TaxSaverLoading = gl2.TaxSaverLoadingRate
	}
	resolveConvContLoadings(&memberDataPointResult, &groupQuote.SchemeCategories[i], &gl2)

	memberDataPointResult.LoadedGlaRate = memberDataPointResult.BaseGlaRate * (1 + memberDataPointResult.GlaContingencyLoading + memberDataPointResult.GlaVoluntaryLoading + memberDataPointResult.GlaTerminalIllnessLoading + memberDataPointResult.GlaContinuityDuringDisabilityLoading + memberDataPointResult.GlaConversionOnWithdrawalLoading + memberDataPointResult.GlaConversionOnRetirementLoading)

	if groupQuote.SchemeCategories[i].PtdBenefit {
		ptdRate := applyCoverAgeLimit(GetPtdRate(&memberDataPointResult, groupParameter, groupQuote, groupQuote.SchemeCategories[i], mpIncomeLevel), memberDataPointResult.AgeNextBirthday, restriction.PtdMaxCoverAge)
		if groupQuote.SchemeCategories[i].PtdBenefitType == "Accelerated" {
			memberDataPointResult.BasePtdRate = ptdRate * (1 + memberDataPointResult.PtdIndustryLoading + memberDataPointResult.PtdRegionLoading - gl2.PtdAcceleratedBenefitDiscount)
		} else {
			memberDataPointResult.BasePtdRate = ptdRate * (1 + memberDataPointResult.PtdIndustryLoading + memberDataPointResult.PtdRegionLoading)
		}
		memberDataPointResult.LoadedPtdRate = memberDataPointResult.BasePtdRate * (1 + memberDataPointResult.PtdContingencyLoading + memberDataPointResult.PtdVoluntaryLoading + memberDataPointResult.PtdConversionOnWithdrawalLoading)
	}

	if groupQuote.SchemeCategories[i].TtdBenefit {
		ttdRate := applyCoverAgeLimit(GetTtdRate(&memberDataPointResult, groupParameter, groupQuote, groupQuote.SchemeCategories[i], mpIncomeLevel), memberDataPointResult.AgeNextBirthday, restriction.TtdMaxCoverAge)
		memberDataPointResult.BaseTtdRate = ttdRate * (1 + memberDataPointResult.TtdIndustryLoading + memberDataPointResult.TtdRegionLoading)
		memberDataPointResult.LoadedTtdRate = memberDataPointResult.BaseTtdRate * (1 + memberDataPointResult.TtdContingencyLoading + memberDataPointResult.TtdVoluntaryLoading + memberDataPointResult.TtdConversionOnWithdrawalLoading)
	}

	if groupQuote.SchemeCategories[i].PhiBenefit {
		phiRate := applyCoverAgeLimit(GetPhiRate(&memberDataPointResult, groupParameter, groupQuote, groupQuote.SchemeCategories[i], mpIncomeLevel), memberDataPointResult.AgeNextBirthday, restriction.PhiMaxCoverAge)
		memberDataPointResult.PhiSalaryLevel = float64(mpIncomeLevel)
		memberDataPointResult.BasePhiRate = phiRate * (1 + memberDataPointResult.PhiIndustryLoading + memberDataPointResult.PhiRegionLoading)
		memberDataPointResult.LoadedPhiRate = memberDataPointResult.BasePhiRate * (1 + memberDataPointResult.PhiContingencyLoading + memberDataPointResult.PhiVoluntaryLoading + memberDataPointResult.PhiConversionOnWithdrawalLoading)
	}

	if groupQuote.SchemeCategories[i].CiBenefit {
		ciRate := applyCoverAgeLimit(GetCiRate(&memberDataPointResult, groupParameter, groupQuote, groupQuote.SchemeCategories[i], mpIncomeLevel), memberDataPointResult.AgeNextBirthday, restriction.CiMaxCoverAge)
		if groupQuote.SchemeCategories[i].CiBenefitStructure == "Accelerated" {
			memberDataPointResult.BaseCiRate = ciRate * (1 + memberDataPointResult.CiIndustryLoading + memberDataPointResult.CiRegionLoading - gl2.CiAcceleratedBenefitDiscount)
		} else {
			memberDataPointResult.BaseCiRate = ciRate * (1 + memberDataPointResult.CiIndustryLoading + memberDataPointResult.CiRegionLoading)
		}
		memberDataPointResult.LoadedCiRate = memberDataPointResult.BaseCiRate * (1 + memberDataPointResult.CiContingencyLoading + memberDataPointResult.CiVoluntaryLoading + memberDataPointResult.CiConversionOnWithdrawalLoading)
	}

	if groupQuote.SchemeCategories[i].SglaBenefit && len(memberDataPointResult.SpouseGender) > 0 {
		spouseIndustryLoading := industryLoadingByGender[strings.ToUpper(memberDataPointResult.SpouseGender[:1])]
		spouseRegionLoading := regionLoadingByGender[strings.ToUpper(memberDataPointResult.SpouseGender[:1])]
		memberDataPointResult.SpouseGlaQx = applyCoverAgeLimit(GetSpouseGlaRate(&memberDataPointResult, groupParameter, mpIncomeLevel, groupQuote, groupQuote.SchemeCategories[i]), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.SpouseGlaAidsQx = applyCoverAgeLimit(GetSpouseGlaAidsRate(&memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.SpouseGlaLoading = spouseIndustryLoading.GlaIndustryLoadingRate
		memberDataPointResult.BaseSpouseGlaRate = memberDataPointResult.SpouseGlaQx*(1+memberDataPointResult.SpouseGlaLoading+spouseRegionLoading.GlaRegionLoadingRate) + memberDataPointResult.SpouseGlaAidsQx*(1+spouseRegionLoading.GlaAidsRegionLoadingRate)
		memberDataPointResult.LoadedSpouseGlaRate = memberDataPointResult.BaseSpouseGlaRate * (1 + memberDataPointResult.SglaConversionOnWithdrawalLoading)
	}

	memberDataPointResult.GlaRiskPremium = memberDataPointResult.LoadedGlaRate * memberDataPointResult.GlaCappedSumAssured
	memberDataPointResult.PtdRiskPremium = memberDataPointResult.LoadedPtdRate * memberDataPointResult.PtdCappedSumAssured
	memberDataPointResult.TtdNumberOfMonthlyPayments = groupParameter.TtdNumberMonthlyPayments
	memberDataPointResult.TtdRiskPremium = memberDataPointResult.LoadedTtdRate * memberDataPointResult.TtdCappedIncome * groupParameter.TtdNumberMonthlyPayments
	memberDataPointResult.PhiRiskPremium = memberDataPointResult.LoadedPhiRate * memberDataPointResult.PhiMonthlyBenefit
	memberDataPointResult.CiRiskPremium = memberDataPointResult.LoadedCiRate * memberDataPointResult.CiCappedSumAssured

	return TheoreticalRiskTotal{
		GlaSumRiskPremium: memberDataPointResult.GlaRiskPremium,
		PtdSumRiskPremium: memberDataPointResult.PtdRiskPremium,
		TtdSumRiskPremium: memberDataPointResult.TtdRiskPremium,
		PhiSumRiskPremium: memberDataPointResult.PhiRiskPremium,
		CiSumRiskPremium:  memberDataPointResult.CiRiskPremium,
		GlaSumAssured:     memberDataPointResult.GlaCappedSumAssured,
		PtdSumAssured:     memberDataPointResult.PtdCappedSumAssured,
		TtdSumAssured:     memberDataPointResult.TtdCappedIncome * groupParameter.TtdNumberMonthlyPayments,
		PhiSumAssured:     memberDataPointResult.PhiCappedIncome,
		CiSumAssured:      memberDataPointResult.CiCappedSumAssured,
	}
}

func PopulateRatesPerMember(i int, indicativeRatesCount float64, indicativeMemberMps []models.MemberIndicativeDataSet, selectedSchemeCategory string, mp models.GPricingMemberData, groupQuote models.GroupPricingQuote, groupParameter models.GroupPricingParameters, groupPricingReinsuranceStructure models.GroupPricingReinsuranceStructure, incomeLevels []models.IncomeLevel, ageBands []models.GroupPricingAgeBands, credibilityRate, annualGlaExperienceWeightedRate, annualPtdExperienceWeightedRate, annualCiExperienceWeightedRate, calculatedFreeCoverLimit float64, educatorBenefitStructure models.EducatorBenefitStructure, credibility, glatheoreticalRate float64, insurerYearEndMonth int, user models.AppUser, restriction models.Restriction, reinsCoverCaps map[string]float64, taxTable []models.TaxTable, taxRetirementBands []models.TaxRetirementTable, tieredIncomeTiers []models.TieredIncomeReplacement, customTieredIncomeTiers []models.TieredIncomeReplacement, premiumLoading models.PremiumLoading, regionLoadingByGender map[string]models.RegionLoading, industryLoadingByGender map[string]models.IndustryLoading, reinsRegionLoadingByGender map[string]models.ReinsuranceRegionLoading, experienceRateOverrides map[string]models.GroupPricingExperienceRateOverride, schemeSizeRow models.SchemeSizeLevel) MemberRateResult {

	var memberDataPointResult models.MemberRatingResult
	var memberPremiumScheduleDatapoint models.MemberPremiumSchedule
	var bordereauxDatapoint models.Bordereaux
	var mpIncomeLevel int
	var mpAgeBand string
	var groupFuneralParameter models.FuneralParameters
	var binderFeeRate3, outsourceFeeRate3 float64

	memberDataPointResult.IsOriginalMember = true
	memberDataPointResult.QuoteId = groupQuote.ID
	memberDataPointResult.CreatedBy = user.UserName
	_, memberDataPointResult.FinancialYear = getGroupRiskQuotingFinancialYear(groupQuote.CommencementDate, insurerYearEndMonth)

	memberDataPointResult.Category = selectedSchemeCategory
	memberDataPointResult.SchemeId = mp.SchemeId
	memberDataPointResult.MemberName = mp.MemberName
	memberDataPointResult.Gender = mp.Gender
	memberDataPointResult.DateOfBirth = mp.DateOfBirth
	memberDataPointResult.AnnualSalary = mp.AnnualSalary
	// Look up region loading for this member by gender using the category-level pre-loaded map
	memberRegionLoading := regionLoadingByGender[strings.ToUpper(memberDataPointResult.Gender[:1])]
	if !groupQuote.UseGlobalSalaryMultiple {
		memberDataPointResult.GlaSalaryMultiple = mp.Benefits.GlaMultiple
		if groupQuote.SchemeCategories[i].PtdBenefit {
			memberDataPointResult.PtdSalaryMultiple = mp.Benefits.PtdMultiple
		}
		if groupQuote.SchemeCategories[i].SglaBenefit {
			memberDataPointResult.SglaSalaryMultiple = mp.Benefits.SglaMultiple
		}
		if groupQuote.SchemeCategories[i].CiBenefit {
			memberDataPointResult.CiSalaryMultiple = mp.Benefits.CiMultiple
		}
	}
	if groupQuote.UseGlobalSalaryMultiple {
		memberDataPointResult.GlaSalaryMultiple = groupQuote.SchemeCategories[i].GlaSalaryMultiple
		memberDataPointResult.PtdSalaryMultiple = groupQuote.SchemeCategories[i].PtdSalaryMultiple
		memberDataPointResult.SglaSalaryMultiple = groupQuote.SchemeCategories[i].SglaSalaryMultiple
		memberDataPointResult.CiSalaryMultiple = groupQuote.SchemeCategories[i].CiCriticalIllnessSalaryMultiple

	}
	memberDataPointResult.Occupation = groupQuote.Industry
	memberDataPointResult.OccupationClass = groupQuote.OccupationClass
	memberDataPointResult.Industry = groupQuote.Industry
	memberDataPointResult.ExpenseLoading = premiumLoading.ExpenseLoading
	memberDataPointResult.AdminLoading = premiumLoading.AdminLoading
	//effectiveCommission3 := premiumLoading.CommissionLoading
	//if groupQuote.DistributionChannel == models.ChannelDirect {
	//	effectiveCommission3 = 0
	//}
	//memberDataPointResult.CommissionLoading = effectiveCommission3
	memberDataPointResult.ProfitLoading = premiumLoading.ProfitLoading
	memberDataPointResult.OtherLoading = premiumLoading.OtherLoading
	memberDataPointResult.Discount = -(groupQuote.Loadings.Discount / 100.0)

	if groupQuote.DistributionChannel == "binder" {
		binderFeeRate3, outsourceFeeRate3 = binderAndOutsourceRates(&groupQuote)
		memberDataPointResult.BinderFeeRate = binderFeeRate3
		memberDataPointResult.OutsourceFeeRate = outsourceFeeRate3
	}

	// Commission is excluded from TotalLoading — see applySchemeWideCommission.
	memberDataPointResult.TotalPremiumLoading = math.Max(premiumLoading.ExpenseLoading+premiumLoading.AdminLoading+premiumLoading.ProfitLoading+premiumLoading.OtherLoading+binderFeeRate3+outsourceFeeRate3, premiumLoading.MinimumPremiumLoading)

	memberDataPointResult.ExpCredibility = credibilityRate
	memberDataPointResult.ManuallyAddedCredibility = credibility

	if credibility > 0 {
		credibilityRate = memberDataPointResult.ManuallyAddedCredibility
	}
	memberDataPointResult.GlaWeightedExperienceCrudeRate = annualGlaExperienceWeightedRate
	memberDataPointResult.PtdExperienceCrudeRate = annualPtdExperienceWeightedRate
	memberDataPointResult.CiExperienceCrudeRate = annualCiExperienceWeightedRate

	memberDataPointResult.CalculatedFreeCoverLimit = calculatedFreeCoverLimit
	if groupQuote.FreeCoverLimit > 0 {
		memberDataPointResult.AppliedFreeCoverLimit = groupQuote.FreeCoverLimit
	}
	if groupQuote.MemberIndicativeData {
		memberDataPointResult.MemberCount = indicativeMemberMps[0].MemberDataCount
		memberDataPointResult.AppliedFreeCoverLimit = memberDataPointResult.AnnualSalary * groupQuote.SchemeCategories[i].GlaSalaryMultiple
	} else {
		memberDataPointResult.MemberCount = 1
	}

	if groupQuote.FreeCoverLimit == 0 {
		memberDataPointResult.AppliedFreeCoverLimit = calculatedFreeCoverLimit
	}

	memberPremiumScheduleDatapoint.MemberName = memberDataPointResult.MemberName
	memberPremiumScheduleDatapoint.QuoteId = memberDataPointResult.QuoteId
	memberPremiumScheduleDatapoint.Gender = memberDataPointResult.Gender
	memberPremiumScheduleDatapoint.SchemeId = memberDataPointResult.SchemeId

	//dob, err := utils.ParseDateString(mp.DateOfBirth)
	//if err != nil {
	//	fmt.Println("error encountered parsing date of birth: ", err)
	//	//return err
	//}

	if groupQuote.MemberIndicativeData {
		memberDataPointResult.AgeNextBirthday = indicativeMemberMps[0].MemberAverageAge
	}

	if !groupQuote.MemberIndicativeData {
		memberDataPointResult.AgeNextBirthday = calculateAgeNextBirthday(groupQuote.CommencementDate, mp.DateOfBirth)
	}
	groupFuneralParameter = GetFuneralParameter(groupParameter, memberDataPointResult.AgeNextBirthday)
	mpAgeBand = GetAgeBand(memberDataPointResult.AgeNextBirthday, ageBands)
	memberDataPointResult.AgeBand = mpAgeBand
	if mp.Gender == "M" {
		memberDataPointResult.SpouseGender = "F"
		memberDataPointResult.SpouseAgeNextBirthday = int(math.Min(math.Max(float64(memberDataPointResult.AgeNextBirthday)-
			float64(groupParameter.SpouseAgeGap), float64(restriction.MinEntryAge)), float64(restriction.MaxEntryAge)))

	}
	if mp.Gender == "F" {
		memberDataPointResult.SpouseGender = "M"
		memberDataPointResult.SpouseAgeNextBirthday = int(math.Min(math.Max(float64(memberDataPointResult.AgeNextBirthday)+
			float64(groupParameter.SpouseAgeGap), float64(restriction.MinEntryAge)), float64(restriction.MaxEntryAge)))
	}

	mpIncomeLevel = GetIncomeLevel(mp, incomeLevels)

	memberDataPointResult.IncomeLevel = mpIncomeLevel

	if glatheoreticalRate > 0 && groupQuote.ExperienceRating == "Yes" {
		memberDataPointResult.GlaTheoreticalRate = glatheoreticalRate
		memberDataPointResult.GlaExperienceAdjustment = (glatheoreticalRate*(1-credibilityRate) + credibilityRate*annualGlaExperienceWeightedRate) / glatheoreticalRate
	} else {
		memberDataPointResult.GlaExperienceAdjustment = 1
	}

	memberDataPointResult.PtdExperienceAdjustment = 1
	memberDataPointResult.CiExperienceAdjustment = 1
	memberDataPointResult.PhiExperienceAdjustment = 1
	memberDataPointResult.TtdExperienceAdjustment = 1

	memberDataPointResult.AverageDependantAgeNextBirthday = groupFuneralParameter.AverageDependantAge
	memberDataPointResult.AverageChildAgeNextBirthday = groupFuneralParameter.AverageChildAge
	memberDataPointResult.AverageNumberDependants = groupFuneralParameter.NumberDependants
	memberDataPointResult.AverageNumberChildren = groupFuneralParameter.NumberChildren
	unScaledGlaSumAssured := mp.AnnualSalary * memberDataPointResult.GlaSalaryMultiple

	memberIndustryLoading := industryLoadingByGender[strings.ToUpper(memberDataPointResult.Gender[:1])]

	// Persist the per-member region and industry loadings so formulas below
	// can read them directly from MemberRatingResult.
	memberDataPointResult.GlaRegionLoading = memberRegionLoading.GlaRegionLoadingRate
	memberDataPointResult.GlaAidsRegionLoading = memberRegionLoading.GlaAidsRegionLoadingRate
	memberDataPointResult.PtdRegionLoading = memberRegionLoading.PtdRegionLoadingRate
	memberDataPointResult.CiRegionLoading = memberRegionLoading.CiRegionLoadingRate
	memberDataPointResult.TtdRegionLoading = memberRegionLoading.TtdRegionLoadingRate
	memberDataPointResult.PhiRegionLoading = memberRegionLoading.PhiRegionLoadingRate
	memberDataPointResult.FunRegionLoading = memberRegionLoading.FunRegionLoadingRate
	memberDataPointResult.FunAidsRegionLoading = memberRegionLoading.FunAidsRegionLoadingRate

	memberDataPointResult.GlaIndustryLoading = memberIndustryLoading.GlaIndustryLoadingRate
	memberDataPointResult.PtdIndustryLoading = memberIndustryLoading.PtdIndustryLoadingRate
	memberDataPointResult.CiIndustryLoading = memberIndustryLoading.CiIndustryLoadingRate
	memberDataPointResult.TtdIndustryLoading = memberIndustryLoading.TtdIndustryLoadingRate
	memberDataPointResult.PhiIndustryLoading = memberIndustryLoading.PhiIndustryLoadingRate

	if groupQuote.SchemeCategories[i].GlaBenefit {
		memberDataPointResult.GlaSumAssured = mp.AnnualSalary * memberDataPointResult.GlaSalaryMultiple * indicativeRatesCount
		memberDataPointResult.GlaCappedSumAssured = applyMaxCoverCap(applyMaxCoverCap(math.Min(unScaledGlaSumAssured, memberDataPointResult.AppliedFreeCoverLimit), restriction.MaximumGlaCover), reinsCoverCaps[benefitTypeKey(groupQuote.SchemeCategories[i].GlaAlias, models.BenefitTypeGla)]) * indicativeRatesCount

		// TaxSaver grosses up the GLA covered sum assured so the post-retirement
		// payout net of the retirement-fund lump-sum tax matches the covered SA.
		// Only runs when the category has TaxSaverBenefit on and the retirement
		// tax bands for the quote's risk rate code have been loaded.
		if groupQuote.SchemeCategories[i].TaxSaverBenefit && len(taxRetirementBands) > 0 {
			memberDataPointResult.TaxSaverSumAssured = computeTaxSaverSumAssured(memberDataPointResult.GlaCappedSumAssured, taxRetirementBands)
		}
	}

	if groupQuote.SchemeCategories[i].PtdBenefit {
		memberDataPointResult.PtdSumAssured = mp.AnnualSalary * memberDataPointResult.PtdSalaryMultiple * indicativeRatesCount
		if groupQuote.SchemeCategories[i].GlaBenefit {
			memberDataPointResult.PtdCappedSumAssured = applyMaxCoverCap(applyMaxCoverCap(math.Min(mp.AnnualSalary*memberDataPointResult.PtdSalaryMultiple, math.Min(unScaledGlaSumAssured, memberDataPointResult.AppliedFreeCoverLimit)), restriction.MaximumPtdCover), reinsCoverCaps[benefitTypeKey(groupQuote.SchemeCategories[i].PtdAlias, models.BenefitTypePtd)]) * indicativeRatesCount
		}
		if !groupQuote.SchemeCategories[i].GlaBenefit {
			memberDataPointResult.PtdCappedSumAssured = applyMaxCoverCap(applyMaxCoverCap(math.Min(mp.AnnualSalary*memberDataPointResult.PtdSalaryMultiple, memberDataPointResult.AppliedFreeCoverLimit), restriction.MaximumPtdCover), reinsCoverCaps[benefitTypeKey(groupQuote.SchemeCategories[i].PtdAlias, models.BenefitTypePtd)]) * indicativeRatesCount
		}
	}
	if groupQuote.SchemeCategories[i].CiBenefit {
		memberDataPointResult.CiSumAssured = mp.AnnualSalary * memberDataPointResult.CiSalaryMultiple * indicativeRatesCount
		memberDataPointResult.CiCappedSumAssured = applyMaxCoverCap(math.Min(mp.AnnualSalary*memberDataPointResult.CiSalaryMultiple, restriction.SevereIllnessMaximumBenefit), reinsCoverCaps[benefitTypeKey(groupQuote.SchemeCategories[i].CiAlias, models.BenefitTypeCi)]) * indicativeRatesCount

	}
	if groupQuote.SchemeCategories[i].SglaBenefit {
		memberDataPointResult.SpouseGlaSumAssured = mp.AnnualSalary * memberDataPointResult.SglaSalaryMultiple * indicativeRatesCount
		memberDataPointResult.SpouseGlaCappedSumAssured = applyMaxCoverCap(math.Min(math.Min(mp.AnnualSalary*memberDataPointResult.SglaSalaryMultiple, restriction.SpouseGlaMaximumBenefit), math.Min(unScaledGlaSumAssured, memberDataPointResult.AppliedFreeCoverLimit)), reinsCoverCaps[benefitTypeKey(groupQuote.SchemeCategories[i].SglaAlias, models.BenefitTypeSgla)]) * indicativeRatesCount

	}

	takeHomePay, _ := computeTakeHomePayFromBands(mp.AnnualSalary, taxTable)

	if groupQuote.SchemeCategories[i].TtdBenefit {
		//memberDataPointResult.TtdIncome = mp.AnnualSalary * groupQuote.SchemeCategories[i].TtdIncomeReplacementPercentage * indicativeRatesCount / 12.0 / 100.0
		var unScaledTtdIncome float64
		if groupQuote.SchemeCategories[i].TtdUseTieredIncomeReplacementRatio {
			ttdTiers := tieredIncomeTiers
			if groupQuote.SchemeCategories[i].TtdTieredIncomeReplacementType == "custom" {
				ttdTiers = customTieredIncomeTiers
			}
			unScaledTtdIncome = math.Min(takeHomePay, computeCoveredIncomeFromTiers(mp.AnnualSalary, ttdTiers)) / 12.0
		} else {
			unScaledTtdIncome = math.Min(takeHomePay, mp.AnnualSalary*groupQuote.SchemeCategories[i].TtdIncomeReplacementPercentage/100.0) / 12.0
		}
		memberDataPointResult.TtdIncome = unScaledTtdIncome * indicativeRatesCount
		memberDataPointResult.TtdCappedIncome = applyMaxCoverCap(math.Min(unScaledTtdIncome, restriction.TtdMaximumMonthlyBenefit), reinsCoverCaps[benefitTypeKey(groupQuote.SchemeCategories[i].TtdAlias, models.BenefitTypeTtd)]) * indicativeRatesCount
	}

	if groupQuote.SchemeCategories[i].PhiBenefit {
		//memberDataPointResult.PhiIncome = mp.AnnualSalary * groupQuote.SchemeCategories[i].PhiIncomeReplacementPercentage * indicativeRatesCount / 12.0 / 100.0
		var unScaledPhiIncome float64
		if groupQuote.SchemeCategories[i].PhiUseTieredIncomeReplacementRatio {
			phiTiers := tieredIncomeTiers
			if groupQuote.SchemeCategories[i].PhiTieredIncomeReplacementType == "custom" {
				phiTiers = customTieredIncomeTiers
			}
			unScaledPhiIncome = math.Min(takeHomePay, computeCoveredIncomeFromTiers(mp.AnnualSalary, phiTiers)) / 12.0
		} else {
			unScaledPhiIncome = math.Min(takeHomePay, mp.AnnualSalary*groupQuote.SchemeCategories[i].PhiIncomeReplacementPercentage/100.0) / 12.0
		}
		memberDataPointResult.PhiIncome = unScaledPhiIncome * indicativeRatesCount
		memberDataPointResult.PhiCappedIncome = applyMaxCoverCap(math.Min(unScaledPhiIncome, restriction.PhiMaximumMonthlyBenefit), reinsCoverCaps[benefitTypeKey(groupQuote.SchemeCategories[i].PhiAlias, models.BenefitTypePhi)]) * indicativeRatesCount
		memberDataPointResult.IncomeReplacementRatio = groupQuote.SchemeCategories[i].PhiIncomeReplacementPercentage / 100.0
		if groupQuote.SchemeCategories[i].PhiPremiumWaiver == "Yes" {
			memberDataPointResult.PhiContributionWaiver = math.Min(mp.AnnualSalary*mp.ContributionWaiverProportion/12.0, restriction.PhiMaximumMonthlyContributionWaiver) * indicativeRatesCount
		}
		if groupQuote.SchemeCategories[i].PhiMedicalAidPremiumWaiver == "Yes" {
			var waiver float64
			if GetMedicalAidWaiverMethod() == models.MedicalAidWaiverMethodTableLookup {
				waiver = GetMedicalWaiverSumAtRisk(&memberDataPointResult, groupParameter, mpIncomeLevel)
			} else {
				waiver = mp.AnnualSalary*groupParameter.MedicalAidWaiverProportion/12.0 + groupParameter.MedicalAidWaiverAmount
			}
			memberDataPointResult.PhiMedicalAidWaiver = math.Min(waiver, restriction.MaxMedicalAidWaiver) * indicativeRatesCount
		}
		memberDataPointResult.PhiMonthlyBenefit = memberDataPointResult.PhiCappedIncome + memberDataPointResult.PhiContributionWaiver + memberDataPointResult.PhiMedicalAidWaiver
	}

	// Cover-age zero-out: a member whose AgeNextBirthday exceeds a benefit's
	// max cover age contributes no covered SA / income for that benefit. A
	// MaxCoverAge of 0 means "no limit" (column unset). SGLA reuses
	// GlaMaxCoverAge since restriction has no separate SGLA column. Funeral
	// is gated downstream where MemberFuneralSumAssured is computed.
	mpAge := memberDataPointResult.AgeNextBirthday
	memberDataPointResult.GlaCappedSumAssured = applyCoverAgeLimit(memberDataPointResult.GlaCappedSumAssured, mpAge, restriction.GlaMaxCoverAge)
	memberDataPointResult.PtdCappedSumAssured = applyCoverAgeLimit(memberDataPointResult.PtdCappedSumAssured, mpAge, restriction.PtdMaxCoverAge)
	memberDataPointResult.CiCappedSumAssured = applyCoverAgeLimit(memberDataPointResult.CiCappedSumAssured, mpAge, restriction.CiMaxCoverAge)
	memberDataPointResult.SpouseGlaCappedSumAssured = applyCoverAgeLimit(memberDataPointResult.SpouseGlaCappedSumAssured, mpAge, restriction.GlaMaxCoverAge)
	memberDataPointResult.TtdCappedIncome = applyCoverAgeLimit(memberDataPointResult.TtdCappedIncome, mpAge, restriction.TtdMaxCoverAge)
	memberDataPointResult.PhiCappedIncome = applyCoverAgeLimit(memberDataPointResult.PhiCappedIncome, mpAge, restriction.PhiMaxCoverAge)

	// PhiMonthlyBenefit was computed inside the PhiBenefit block above; recompute
	// it so the gated PhiCappedIncome propagates into the monthly benefit total.
	if groupQuote.SchemeCategories[i].PhiBenefit {
		memberDataPointResult.PhiMonthlyBenefit = memberDataPointResult.PhiCappedIncome + memberDataPointResult.PhiContributionWaiver + memberDataPointResult.PhiMedicalAidWaiver
	}

	// TaxSaver was computed inline at the GLA calc site with the pre-gate value;
	// recompute now that GlaCappedSumAssured has been age-gated.
	if groupQuote.SchemeCategories[i].TaxSaverBenefit && len(taxRetirementBands) > 0 {
		memberDataPointResult.TaxSaverSumAssured = computeTaxSaverSumAssured(memberDataPointResult.GlaCappedSumAssured, taxRetirementBands)
	}

	if groupQuote.SchemeCategories[i].GlaBenefit {
		memberDataPointResult.GlaQx = applyCoverAgeLimit(GetGlaRate(&memberDataPointResult, groupParameter, mpIncomeLevel, groupQuote, groupQuote.SchemeCategories[i]), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.GlaAidsQx = applyCoverAgeLimit(GetGlaAidsRate(&memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.BaseGlaRate = memberDataPointResult.GlaQx*(1+memberDataPointResult.GlaIndustryLoading+memberDataPointResult.GlaRegionLoading) + memberDataPointResult.GlaAidsQx*(1+memberDataPointResult.GlaAidsRegionLoading)
	} else if groupQuote.SchemeCategories[i].FamilyFuneralBenefit {
		funQx := applyCoverAgeLimit(GetFuneralRate(&memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
		funAidsQx := applyCoverAgeLimit(GetFuneralAidsRate(&memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
		memberDataPointResult.MainMemberFuneralBaseRate = funQx*(1+memberDataPointResult.FunRegionLoading) + funAidsQx*(1+memberDataPointResult.FunAidsRegionLoading)
		if len(memberDataPointResult.SpouseGender) > 0 {
			spouseRegionLoading := regionLoadingByGender[strings.ToUpper(memberDataPointResult.SpouseGender[:1])]
			spouseFunQx := applyCoverAgeLimit(GetSpouseFuneralRate(&memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
			spouseFunAidsQx := applyCoverAgeLimit(GetSpouseFuneralAidsRate(&memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
			memberDataPointResult.SpouseFuneralBaseRate = spouseFunQx*(1+spouseRegionLoading.FunRegionLoadingRate) + spouseFunAidsQx*(1+spouseRegionLoading.FunAidsRegionLoadingRate)
		}
	}

	gl3 := GetGeneralLoading(groupParameter.RiskRateCode, memberDataPointResult.AgeNextBirthday, memberDataPointResult.Gender)
	if groupQuote.SchemeCategories[i].GlaTerminalIllnessBenefit == "Yes" {
		memberDataPointResult.GlaTerminalIllnessLoading = gl3.TerminalIllnessLoadingRate
	}

	// Persist the per-member contingency loadings alongside the already-set
	// region/industry loadings.
	memberDataPointResult.GlaContingencyLoading = gl3.GlaContigencyLoadingRate
	memberDataPointResult.PtdContingencyLoading = gl3.PtdContigencyLoadingRate
	memberDataPointResult.CiContingencyLoading = gl3.CiContigencyLoadingRate
	memberDataPointResult.TtdContingencyLoading = gl3.TtdContigencyLoadingRate
	memberDataPointResult.PhiContingencyLoading = gl3.PhiContigencyLoadingRate
	memberDataPointResult.FunContingencyLoading = gl3.FunContigencyLoadingRate

	// Scheme-size loadings (resolved once per quote from the SchemeSizeLevel
	// row matching the quote's MemberDataCount). Folded into each benefit's
	// LoadedRate multiplier below.
	memberDataPointResult.GlaSchemeSizeLoading = schemeSizeRow.GlaLoading
	memberDataPointResult.PtdSchemeSizeLoading = schemeSizeRow.PtdLoading
	memberDataPointResult.CiSchemeSizeLoading = schemeSizeRow.CiLoading
	memberDataPointResult.TtdSchemeSizeLoading = schemeSizeRow.TtdLoading
	memberDataPointResult.PhiSchemeSizeLoading = schemeSizeRow.PhiLoading
	memberDataPointResult.FunSchemeSizeLoading = schemeSizeRow.FunLoading

	// Voluntary loadings only apply when the quote's obligation type is
	// Voluntary; otherwise they stay zero.
	if groupQuote.ObligationType == "Voluntary" {
		memberDataPointResult.GlaVoluntaryLoading = gl3.GlaVoluntaryLoadingRate
		memberDataPointResult.PtdVoluntaryLoading = gl3.PtdVoluntaryLoadingRate
		memberDataPointResult.CiVoluntaryLoading = gl3.CiVoluntaryLoadingRate
		memberDataPointResult.TtdVoluntaryLoading = gl3.TtdVoluntaryLoadingRate
		memberDataPointResult.PhiVoluntaryLoading = gl3.PhiVoluntaryLoadingRate
		memberDataPointResult.FunVoluntaryLoading = gl3.FunVoluntaryLoadingRate
	}

	if groupQuote.SchemeCategories[i].GlaBenefit && groupQuote.SchemeCategories[i].TaxSaverBenefit {
		memberDataPointResult.TaxSaverLoading = gl3.TaxSaverLoadingRate
	}
	resolveConvContLoadings(&memberDataPointResult, &groupQuote.SchemeCategories[i], &gl3)

	memberDataPointResult.LoadedGlaRate = memberDataPointResult.BaseGlaRate * (1 + memberDataPointResult.GlaContingencyLoading + memberDataPointResult.GlaVoluntaryLoading + memberDataPointResult.GlaTerminalIllnessLoading + memberDataPointResult.GlaContinuityDuringDisabilityLoading + memberDataPointResult.GlaConversionOnWithdrawalLoading + memberDataPointResult.GlaConversionOnRetirementLoading + memberDataPointResult.GlaSchemeSizeLoading)

	if groupQuote.SchemeCategories[i].PtdBenefit {
		ptdRate := applyCoverAgeLimit(GetPtdRate(&memberDataPointResult, groupParameter, groupQuote, groupQuote.SchemeCategories[i], mpIncomeLevel), memberDataPointResult.AgeNextBirthday, restriction.PtdMaxCoverAge)
		if groupQuote.SchemeCategories[i].PtdBenefitType == "Accelerated" {
			memberDataPointResult.BasePtdRate = ptdRate * (1 + memberDataPointResult.PtdIndustryLoading + memberDataPointResult.PtdRegionLoading - gl3.PtdAcceleratedBenefitDiscount)
		} else {
			memberDataPointResult.BasePtdRate = ptdRate * (1 + memberDataPointResult.PtdIndustryLoading + memberDataPointResult.PtdRegionLoading)
		}
		memberDataPointResult.LoadedPtdRate = memberDataPointResult.BasePtdRate * (1 + memberDataPointResult.PtdContingencyLoading + memberDataPointResult.PtdVoluntaryLoading + memberDataPointResult.PtdConversionOnWithdrawalLoading + memberDataPointResult.PtdSchemeSizeLoading)
	}

	if groupQuote.SchemeCategories[i].TtdBenefit {
		ttdRate := applyCoverAgeLimit(GetTtdRate(&memberDataPointResult, groupParameter, groupQuote, groupQuote.SchemeCategories[i], mpIncomeLevel), memberDataPointResult.AgeNextBirthday, restriction.TtdMaxCoverAge)
		memberDataPointResult.BaseTtdRate = ttdRate * (1 + memberDataPointResult.TtdIndustryLoading + memberDataPointResult.TtdRegionLoading)
		memberDataPointResult.LoadedTtdRate = memberDataPointResult.BaseTtdRate * (1 + memberDataPointResult.TtdContingencyLoading + memberDataPointResult.TtdVoluntaryLoading + memberDataPointResult.TtdConversionOnWithdrawalLoading + memberDataPointResult.TtdSchemeSizeLoading)
	}

	if groupQuote.SchemeCategories[i].PhiBenefit {
		phiRate := applyCoverAgeLimit(GetPhiRate(&memberDataPointResult, groupParameter, groupQuote, groupQuote.SchemeCategories[i], mpIncomeLevel), memberDataPointResult.AgeNextBirthday, restriction.PhiMaxCoverAge)
		memberDataPointResult.PhiSalaryLevel = float64(mpIncomeLevel)
		memberDataPointResult.BasePhiRate = phiRate * (1 + memberDataPointResult.PhiIndustryLoading + memberDataPointResult.PhiRegionLoading)
		memberDataPointResult.LoadedPhiRate = memberDataPointResult.BasePhiRate * (1 + memberDataPointResult.PhiContingencyLoading + memberDataPointResult.PhiVoluntaryLoading + memberDataPointResult.PhiConversionOnWithdrawalLoading + memberDataPointResult.PhiSchemeSizeLoading)
	}

	if groupQuote.SchemeCategories[i].CiBenefit {
		ciRate := applyCoverAgeLimit(GetCiRate(&memberDataPointResult, groupParameter, groupQuote, groupQuote.SchemeCategories[i], mpIncomeLevel), memberDataPointResult.AgeNextBirthday, restriction.CiMaxCoverAge)
		if groupQuote.SchemeCategories[i].CiBenefitStructure == "Accelerated" {
			memberDataPointResult.BaseCiRate = ciRate * (1 + memberDataPointResult.CiIndustryLoading + memberDataPointResult.CiRegionLoading - gl3.CiAcceleratedBenefitDiscount)
		} else {
			memberDataPointResult.BaseCiRate = ciRate * (1 + memberDataPointResult.CiIndustryLoading + memberDataPointResult.CiRegionLoading)
		}
		memberDataPointResult.LoadedCiRate = memberDataPointResult.BaseCiRate * (1 + memberDataPointResult.CiContingencyLoading + memberDataPointResult.CiVoluntaryLoading + memberDataPointResult.CiConversionOnWithdrawalLoading + memberDataPointResult.CiSchemeSizeLoading)
	}

	if groupQuote.SchemeCategories[i].SglaBenefit && len(memberDataPointResult.SpouseGender) > 0 {
		spouseIndustryLoading := industryLoadingByGender[strings.ToUpper(memberDataPointResult.SpouseGender[:1])]
		spouseRegionLoading := regionLoadingByGender[strings.ToUpper(memberDataPointResult.SpouseGender[:1])]
		memberDataPointResult.SpouseGlaQx = applyCoverAgeLimit(GetSpouseGlaRate(&memberDataPointResult, groupParameter, mpIncomeLevel, groupQuote, groupQuote.SchemeCategories[i]), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.SpouseGlaAidsQx = applyCoverAgeLimit(GetSpouseGlaAidsRate(&memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.SpouseGlaLoading = spouseIndustryLoading.GlaIndustryLoadingRate
		memberDataPointResult.BaseSpouseGlaRate = memberDataPointResult.SpouseGlaQx*(1+memberDataPointResult.SpouseGlaLoading+spouseRegionLoading.GlaRegionLoadingRate) + memberDataPointResult.SpouseGlaAidsQx*(1+spouseRegionLoading.GlaAidsRegionLoadingRate)
		memberDataPointResult.LoadedSpouseGlaRate = memberDataPointResult.BaseSpouseGlaRate * (1 + gl3.GlaContigencyLoadingRate + memberDataPointResult.SglaConversionOnWithdrawalLoading)
	}

	// ── Reinsurance rates & loadings ───────────────────────────────────────
	// Mirror the direct base/loaded rate pipeline above using the
	// reinsurance-specific rate/loading tables. Region loadings come from
	// the category-level reinsRegionLoadingByGender map; industry and
	// general loadings from the cached reinsurance getters. Each benefit is
	// guarded the same way as the direct block so reinsurance fields stay
	// zero when the relevant benefit is not enabled for this category.
	reinsIndustryLoadingMain := GetReinsuranceIndustryLoading(groupParameter.RiskRateCode, groupQuote.OccupationClass, memberDataPointResult.Gender)
	memberDataPointResult.ReinsGlaIndustryLoading = reinsIndustryLoadingMain.GlaIndustryLoadingRate
	memberDataPointResult.ReinsPtdIndustryLoading = reinsIndustryLoadingMain.PtdIndustryLoadingRate
	memberDataPointResult.ReinsCiIndustryLoading = reinsIndustryLoadingMain.CiIndustryLoadingRate
	memberDataPointResult.ReinsTtdIndustryLoading = reinsIndustryLoadingMain.TtdIndustryLoadingRate
	memberDataPointResult.ReinsPhiIndustryLoading = reinsIndustryLoadingMain.PhiIndustryLoadingRate

	reinsRegionLoadingMain := reinsRegionLoadingByGender[strings.ToUpper(memberDataPointResult.Gender[:1])]
	memberDataPointResult.ReinsGlaRegionLoading = reinsRegionLoadingMain.GlaRegionLoadingRate
	memberDataPointResult.ReinsGlaAidsRegionLoading = reinsRegionLoadingMain.GlaAidsRegionLoadingRate
	memberDataPointResult.ReinsPtdRegionLoading = reinsRegionLoadingMain.PtdRegionLoadingRate
	memberDataPointResult.ReinsCiRegionLoading = reinsRegionLoadingMain.CiRegionLoadingRate
	memberDataPointResult.ReinsTtdRegionLoading = reinsRegionLoadingMain.TtdRegionLoadingRate
	memberDataPointResult.ReinsPhiRegionLoading = reinsRegionLoadingMain.PhiRegionLoadingRate
	memberDataPointResult.ReinsFunRegionLoading = reinsRegionLoadingMain.FunRegionLoadingRate
	memberDataPointResult.ReinsFunAidsRegionLoading = reinsRegionLoadingMain.FunAidsRegionLoadingRate

	reinsGL := GetReinsuranceGeneralLoading(groupParameter.RiskRateCode, memberDataPointResult.AgeNextBirthday, memberDataPointResult.Gender)
	memberDataPointResult.ReinsGlaContingencyLoading = reinsGL.GlaContigencyLoadingRate
	memberDataPointResult.ReinsPtdContingencyLoading = reinsGL.PtdContigencyLoadingRate
	memberDataPointResult.ReinsCiContingencyLoading = reinsGL.CiContigencyLoadingRate
	memberDataPointResult.ReinsTtdContingencyLoading = reinsGL.TtdContigencyLoadingRate
	memberDataPointResult.ReinsPhiContingencyLoading = reinsGL.PhiContigencyLoadingRate
	memberDataPointResult.ReinsFunContingencyLoading = reinsGL.FunContigencyLoadingRate

	// Voluntary loadings only apply when the quote's obligation type is
	// Voluntary; otherwise they stay zero.
	if groupQuote.ObligationType == "Voluntary" {
		memberDataPointResult.ReinsGlaVoluntaryLoading = reinsGL.GlaVoluntaryLoadingRate
		memberDataPointResult.ReinsPtdVoluntaryLoading = reinsGL.PtdVoluntaryLoadingRate
		memberDataPointResult.ReinsCiVoluntaryLoading = reinsGL.CiVoluntaryLoadingRate
		memberDataPointResult.ReinsTtdVoluntaryLoading = reinsGL.TtdVoluntaryLoadingRate
		memberDataPointResult.ReinsPhiVoluntaryLoading = reinsGL.PhiVoluntaryLoadingRate
		memberDataPointResult.ReinsFunVoluntaryLoading = reinsGL.FunVoluntaryLoadingRate
	}
	if groupQuote.SchemeCategories[i].GlaTerminalIllnessBenefit == "Yes" {
		memberDataPointResult.ReinsGlaTerminalIllnessLoading = reinsGL.TerminalIllnessLoadingRate
	}

	if groupQuote.SchemeCategories[i].GlaBenefit {
		memberDataPointResult.ReinsGlaQx = applyCoverAgeLimit(GetReinsuranceGlaRate(&memberDataPointResult, groupParameter, mpIncomeLevel, groupQuote.SchemeCategories[i]), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.ReinsGlaAidsQx = applyCoverAgeLimit(GetReinsuranceGlaAidsRate(&memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.BaseReinsGlaRate = memberDataPointResult.ReinsGlaQx*(1+memberDataPointResult.ReinsGlaIndustryLoading+memberDataPointResult.ReinsGlaRegionLoading) +
			memberDataPointResult.ReinsGlaAidsQx*(1+memberDataPointResult.ReinsGlaAidsRegionLoading)
		memberDataPointResult.LoadedReinsGlaRate = memberDataPointResult.BaseReinsGlaRate *
			(1 + memberDataPointResult.ReinsGlaContingencyLoading + memberDataPointResult.ReinsGlaVoluntaryLoading + memberDataPointResult.ReinsGlaTerminalIllnessLoading)
	}

	if groupQuote.SchemeCategories[i].PtdBenefit {
		memberDataPointResult.ReinsPtdRate = applyCoverAgeLimit(GetReinsurancePtdRate(&memberDataPointResult, groupParameter, mpIncomeLevel, groupQuote.SchemeCategories[i]), memberDataPointResult.AgeNextBirthday, restriction.PtdMaxCoverAge)
		if groupQuote.SchemeCategories[i].PtdBenefitType == "Accelerated" {
			memberDataPointResult.BaseReinsPtdRate = memberDataPointResult.ReinsPtdRate * (1 + memberDataPointResult.ReinsPtdIndustryLoading + memberDataPointResult.ReinsPtdRegionLoading - reinsGL.PtdAcceleratedBenefitDiscount)
		} else {
			memberDataPointResult.BaseReinsPtdRate = memberDataPointResult.ReinsPtdRate * (1 + memberDataPointResult.ReinsPtdIndustryLoading + memberDataPointResult.ReinsPtdRegionLoading)
		}
		memberDataPointResult.LoadedReinsPtdRate = memberDataPointResult.BaseReinsPtdRate * (1 + memberDataPointResult.ReinsPtdContingencyLoading + memberDataPointResult.ReinsPtdVoluntaryLoading)
	}

	if groupQuote.SchemeCategories[i].CiBenefit {
		memberDataPointResult.ReinsCiRate = applyCoverAgeLimit(GetReinsuranceCiRate(&memberDataPointResult, groupParameter, mpIncomeLevel, groupQuote.SchemeCategories[i]), memberDataPointResult.AgeNextBirthday, restriction.CiMaxCoverAge)
		if groupQuote.SchemeCategories[i].CiBenefitStructure == "Accelerated" {
			memberDataPointResult.BaseReinsCiRate = memberDataPointResult.ReinsCiRate * (1 + memberDataPointResult.ReinsCiIndustryLoading + memberDataPointResult.ReinsCiRegionLoading - reinsGL.CiAcceleratedBenefitDiscount)
		} else {
			memberDataPointResult.BaseReinsCiRate = memberDataPointResult.ReinsCiRate * (1 + memberDataPointResult.ReinsCiIndustryLoading + memberDataPointResult.ReinsCiRegionLoading)
		}
		memberDataPointResult.LoadedReinsCiRate = memberDataPointResult.BaseReinsCiRate * (1 + memberDataPointResult.ReinsCiContingencyLoading + memberDataPointResult.ReinsCiVoluntaryLoading)
	}

	if groupQuote.SchemeCategories[i].TtdBenefit {
		// No dedicated reinsurance TTD rate table yet; strip the direct
		// industry/region loading off BaseTtdRate then re-apply the
		// reinsurance-specific industry & region loadings.
		ttdReinsQx := memberDataPointResult.BaseTtdRate / math.Max(1+memberDataPointResult.TtdIndustryLoading+memberDataPointResult.TtdRegionLoading, 1e-9)
		memberDataPointResult.BaseReinsTtdRate = ttdReinsQx * (1 + memberDataPointResult.ReinsTtdIndustryLoading + memberDataPointResult.ReinsTtdRegionLoading)
		memberDataPointResult.LoadedReinsTtdRate = memberDataPointResult.BaseReinsTtdRate * (1 + memberDataPointResult.ReinsTtdContingencyLoading + memberDataPointResult.ReinsTtdVoluntaryLoading)
	}

	if groupQuote.SchemeCategories[i].PhiBenefit {
		memberDataPointResult.ReinsPhiRate = applyCoverAgeLimit(GetReinsurancePhiRate(&memberDataPointResult, groupParameter, mpIncomeLevel, groupQuote.SchemeCategories[i]), memberDataPointResult.AgeNextBirthday, restriction.PhiMaxCoverAge)
		memberDataPointResult.BaseReinsPhiRate = memberDataPointResult.ReinsPhiRate * (1 + memberDataPointResult.ReinsPhiIndustryLoading + memberDataPointResult.ReinsPhiRegionLoading)
		memberDataPointResult.LoadedReinsPhiRate = memberDataPointResult.BaseReinsPhiRate * (1 + memberDataPointResult.ReinsPhiContingencyLoading + memberDataPointResult.ReinsPhiVoluntaryLoading)
	}

	if groupQuote.SchemeCategories[i].SglaBenefit && len(memberDataPointResult.SpouseGender) > 0 {
		spouseReinsIndustry := GetReinsuranceIndustryLoading(groupParameter.RiskRateCode, groupQuote.OccupationClass, memberDataPointResult.SpouseGender)
		memberDataPointResult.ReinsSpouseGlaQx = applyCoverAgeLimit(GetReinsuranceSpouseGlaRate(&memberDataPointResult, groupParameter, mpIncomeLevel, groupQuote.SchemeCategories[i]), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.ReinsSpouseGlaAidsQx = applyCoverAgeLimit(GetReinsuranceSpouseGlaAidsRate(&memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.GlaMaxCoverAge)
		memberDataPointResult.ReinsSpouseGlaLoading = spouseReinsIndustry.GlaIndustryLoadingRate
		spouseReinsRegion := reinsRegionLoadingByGender[strings.ToUpper(memberDataPointResult.SpouseGender[:1])]
		memberDataPointResult.BaseReinsSpouseGlaRate = memberDataPointResult.ReinsSpouseGlaQx*(1+memberDataPointResult.ReinsSpouseGlaLoading+spouseReinsRegion.GlaRegionLoadingRate) +
			memberDataPointResult.ReinsSpouseGlaAidsQx*(1+spouseReinsRegion.GlaAidsRegionLoadingRate)
		memberDataPointResult.LoadedReinsSpouseGlaRate = memberDataPointResult.BaseReinsSpouseGlaRate
	}

	// Main-member / spouse funeral reinsurance rate. When GLA benefit is
	// selected the main-member funeral rides on LoadedReinsGlaRate; otherwise
	// it comes from reinsurance_funeral_rates + reinsurance_funeral_aids_rates.
	if groupQuote.SchemeCategories[i].GlaBenefit {
		memberDataPointResult.MainMemberReinsuranceRate = memberDataPointResult.LoadedReinsGlaRate
		memberDataPointResult.SpouseReinsuranceRate = memberDataPointResult.LoadedReinsSpouseGlaRate
	} else if groupQuote.SchemeCategories[i].FamilyFuneralBenefit {
		reinsFunQx := applyCoverAgeLimit(GetReinsuranceFuneralRate(&memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
		reinsFunAidsQx := applyCoverAgeLimit(GetReinsuranceFuneralAidsRate(&memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
		memberDataPointResult.MainMemberReinsuranceBaseRate = reinsFunQx*(1+memberDataPointResult.ReinsFunRegionLoading) + reinsFunAidsQx*(1+memberDataPointResult.ReinsFunAidsRegionLoading)
		memberDataPointResult.MainMemberReinsuranceRate = memberDataPointResult.MainMemberReinsuranceBaseRate * (1 + memberDataPointResult.ReinsFunContingencyLoading + memberDataPointResult.ReinsFunVoluntaryLoading)
		if len(memberDataPointResult.SpouseGender) > 0 {
			spouseReinsRegion := reinsRegionLoadingByGender[strings.ToUpper(memberDataPointResult.SpouseGender[:1])]
			spouseReinsFunQx := applyCoverAgeLimit(GetReinsuranceSpouseFuneralRate(&memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
			spouseReinsFunAidsQx := applyCoverAgeLimit(GetReinsuranceSpouseFuneralAidsRate(&memberDataPointResult, groupParameter), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
			memberDataPointResult.SpouseReinsuranceBaseRate = spouseReinsFunQx*(1+spouseReinsRegion.FunRegionLoadingRate) + spouseReinsFunAidsQx*(1+spouseReinsRegion.FunAidsRegionLoadingRate)
			memberDataPointResult.SpouseReinsuranceRate = memberDataPointResult.SpouseReinsuranceBaseRate * (1 + memberDataPointResult.ReinsFunContingencyLoading + memberDataPointResult.ReinsFunVoluntaryLoading)
		}
	}

	memberDataPointResult.GlaRiskPremium = memberDataPointResult.LoadedGlaRate * memberDataPointResult.GlaCappedSumAssured
	memberDataPointResult.TaxSaverRiskPremium = memberDataPointResult.LoadedGlaRate * memberDataPointResult.TaxSaverSumAssured
	memberDataPointResult.PtdRiskPremium = memberDataPointResult.LoadedPtdRate * memberDataPointResult.PtdCappedSumAssured
	memberDataPointResult.TtdNumberOfMonthlyPayments = groupParameter.TtdNumberMonthlyPayments
	memberDataPointResult.TtdRiskPremium = memberDataPointResult.LoadedTtdRate * memberDataPointResult.TtdCappedIncome * groupParameter.TtdNumberMonthlyPayments
	memberDataPointResult.PhiRiskPremium = memberDataPointResult.LoadedPhiRate * memberDataPointResult.PhiMonthlyBenefit
	memberDataPointResult.CiRiskPremium = memberDataPointResult.LoadedCiRate * memberDataPointResult.CiCappedSumAssured
	memberDataPointResult.SpouseGlaRiskPremium = memberDataPointResult.LoadedSpouseGlaRate * memberDataPointResult.SpouseGlaCappedSumAssured

	discountFraction := -(groupQuote.Loadings.Discount / 100.0)

	//memberDataPointResult.MarriageProportion = groupFuneralParameter.ProportionMarried

	memberDataPointResult.ChildFuneralBaseRate = applyCoverAgeLimit(GetChildFuneralRate(&memberDataPointResult, groupParameter, memberDataPointResult.AverageChildAgeNextBirthday), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
	memberDataPointResult.ChildFuneralBaseRate *= (1 + memberRegionLoading.FunRegionLoadingRate)
	memberDataPointResult.ChildFuneralSumAssured = groupQuote.SchemeCategories[i].FamilyFuneralChildrenFuneralSumAssured
	memberDataPointResult.ParentFuneralBaseRate = applyCoverAgeLimit(GetDependantMortalityRate(&memberDataPointResult, groupParameter, memberDataPointResult.AverageDependantAgeNextBirthday, mpIncomeLevel), memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge)
	memberDataPointResult.ParentFuneralBaseRate *= (1 + memberRegionLoading.FunRegionLoadingRate)
	memberDataPointResult.ParentFuneralSumAssured = groupQuote.SchemeCategories[i].FamilyFuneralParentFuneralSumAssured
	memberDataPointResult.ParentFuneralSumAssured = groupQuote.SchemeCategories[i].FamilyFuneralAdultDependantSumAssured
	memberDataPointResult.MemberFuneralSumAssured = applyMaxCoverCap(applyCoverAgeLimit(groupQuote.SchemeCategories[i].FamilyFuneralMainMemberFuneralSumAssured, memberDataPointResult.AgeNextBirthday, restriction.FunMaxCoverAge), reinsCoverCaps[benefitTypeKey(groupQuote.SchemeCategories[i].FamilyFuneralAlias, models.BenefitTypeFun)])

	// Child / parent / dependant reinsurance rates. No dedicated per-life
	// reinsurance rate tables exist for these relationships, so we fall back
	// to the direct funeral base rate plus the reinsurance funeral contingency
	// loading. The ceded premium later uses these values × the ceded sum
	// assured computed in GroupPricingReinsurance.
	memberDataPointResult.ChildReinsuranceBaseRate = memberDataPointResult.ChildFuneralBaseRate
	memberDataPointResult.ChildReinsuranceRate = memberDataPointResult.ChildReinsuranceBaseRate * (1 + memberDataPointResult.ReinsFunContingencyLoading + memberDataPointResult.ReinsFunVoluntaryLoading)
	memberDataPointResult.ParentReinsuranceBaseRate = memberDataPointResult.ParentFuneralBaseRate
	memberDataPointResult.ParentReinsuranceRate = memberDataPointResult.ParentReinsuranceBaseRate * (1 + memberDataPointResult.ReinsFunContingencyLoading + memberDataPointResult.ReinsFunVoluntaryLoading)
	memberDataPointResult.DependantReinsuranceBaseRate = memberDataPointResult.ParentFuneralBaseRate
	memberDataPointResult.DependantReinsuranceRate = memberDataPointResult.DependantReinsuranceBaseRate * (1 + memberDataPointResult.ReinsFunContingencyLoading + memberDataPointResult.ReinsFunVoluntaryLoading)

	if groupQuote.SchemeCategories[i].GlaBenefit {
		memberDataPointResult.MainMemberFuneralRiskPremium = memberDataPointResult.LoadedGlaRate * memberDataPointResult.MemberFuneralSumAssured
	} else {
		memberDataPointResult.MainMemberFuneralRiskPremium = memberDataPointResult.MainMemberFuneralBaseRate * memberDataPointResult.MemberFuneralSumAssured
	}
	memberDataPointResult.SpouseFuneralSumAssured = groupQuote.SchemeCategories[i].FamilyFuneralSpouseFuneralSumAssured
	if groupQuote.SchemeCategories[i].GlaBenefit {
		memberDataPointResult.SpouseFuneralRiskPremium = memberDataPointResult.LoadedSpouseGlaRate * groupQuote.SchemeCategories[i].FamilyFuneralSpouseFuneralSumAssured
	} else {
		memberDataPointResult.SpouseFuneralRiskPremium = memberDataPointResult.SpouseFuneralBaseRate * groupQuote.SchemeCategories[i].FamilyFuneralSpouseFuneralSumAssured
	}
	memberDataPointResult.ChildFuneralRiskPremium = memberDataPointResult.ChildFuneralBaseRate * groupQuote.SchemeCategories[i].FamilyFuneralChildrenFuneralSumAssured * math.Min(memberDataPointResult.AverageNumberChildren, float64(groupQuote.SchemeCategories[i].FamilyFuneralMaxNumberChildren))
	memberDataPointResult.ParentFuneralRiskPremium = memberDataPointResult.ParentFuneralBaseRate * groupQuote.SchemeCategories[i].FamilyFuneralAdultDependantSumAssured * memberDataPointResult.AverageNumberDependants

	memberDataPointResult.TotalFuneralRiskPremium = (memberDataPointResult.MainMemberFuneralRiskPremium + memberDataPointResult.SpouseFuneralRiskPremium + memberDataPointResult.ChildFuneralRiskPremium + memberDataPointResult.ParentFuneralRiskPremium) * (1 + memberDataPointResult.FunConversionOnWithdrawalLoading + memberDataPointResult.FunSchemeSizeLoading)
	memberDataPointResult.TotalFuneralOfficePremium = memberDataPointResult.TotalFuneralRiskPremium / (1.0 - memberDataPointResult.TotalPremiumLoading)

	// Derive LoadedGla/PtdEducatorRate now that the parent LoadedRates are
	// set but before the educator block so educator premiums (and their
	// slices) can use the educator-specific loaded rates.
	computeEducatorLoadedRates(&memberDataPointResult)

	if groupQuote.SchemeCategories[i].GlaEducatorBenefit == "Yes" || groupQuote.SchemeCategories[i].PtdEducatorBenefit == "Yes" { //groupQuote.Educator
		educatorRates := GetEducatorRate(groupParameter, educatorCodeForCategory(&groupQuote.SchemeCategories[i]), memberDataPointResult.AgeNextBirthday, memberDataPointResult.IncomeLevel)
		memberDataPointResult.EducatorSumAtRisk = educatorRates.EducatorSumAtRisk

		// Educator premium uses LoadedGla/PtdEducatorRate (which folds in
		// educator-specific conversion / continuity loadings). The split
		// Gla/Ptd educator premiums are computed only when the respective
		// toggle is "Yes". ExpAdj variants come in the Experience Rating block.
		if groupQuote.SchemeCategories[i].GlaEducatorBenefit == "Yes" {
			memberDataPointResult.GlaEducatorRiskPremium = memberDataPointResult.EducatorSumAtRisk * memberDataPointResult.LoadedGlaEducatorRate
		}
		if groupQuote.SchemeCategories[i].PtdEducatorBenefit == "Yes" {
			memberDataPointResult.PtdEducatorRiskPremium = memberDataPointResult.EducatorSumAtRisk * memberDataPointResult.LoadedPtdEducatorRate
		}
		// Educator conversion / continuity slice premiums (risk leg).
		// ExpAdj variants computed in the Experience Rating block below.
		computeEducatorSlicePremiums(&memberDataPointResult, memberDataPointResult.EducatorSumAtRisk)
	}

	//Experience Rating
	memberDataPointResult.ExpAdjLoadedGlaRate = memberDataPointResult.LoadedGlaRate * memberDataPointResult.GlaExperienceAdjustment
	if groupQuote.SchemeCategories[i].PtdBenefit {
		memberDataPointResult.ExpAdjLoadedPtdRate = memberDataPointResult.LoadedPtdRate * memberDataPointResult.PtdExperienceAdjustment
	}
	if groupQuote.SchemeCategories[i].TtdBenefit {
		memberDataPointResult.ExpAdjLoadedTtdRate = memberDataPointResult.LoadedTtdRate * memberDataPointResult.TtdExperienceAdjustment
	}
	if groupQuote.SchemeCategories[i].PhiBenefit {
		memberDataPointResult.ExpAdjLoadedPhiRate = memberDataPointResult.LoadedPhiRate * memberDataPointResult.PhiExperienceAdjustment
	}
	if groupQuote.SchemeCategories[i].CiBenefit {
		memberDataPointResult.ExpAdjLoadedCiRate = memberDataPointResult.LoadedCiRate * memberDataPointResult.CiExperienceAdjustment
	}
	if groupQuote.SchemeCategories[i].SglaBenefit {
		memberDataPointResult.ExpAdjLoadedSpouseGlaRate = memberDataPointResult.LoadedSpouseGlaRate * memberDataPointResult.GlaExperienceAdjustment
	}

	// Override mode: replace ExpAdjLoaded{Benefit}Rate with the user-supplied
	// per-(category, benefit) overrides before downstream risk-premium and
	// educator cascade reads them. Funeral is handled separately below.
	if groupQuote.ExperienceRating == "Override" {
		applyExperienceRateOverridesToMember(&memberDataPointResult, experienceRateOverrides)
	}

	memberDataPointResult.ExpAdjGlaRiskPremium = memberDataPointResult.ExpAdjLoadedGlaRate * memberDataPointResult.GlaCappedSumAssured
	memberDataPointResult.ExpAdjTaxSaverRiskPremium = memberDataPointResult.ExpAdjLoadedGlaRate * memberDataPointResult.TaxSaverSumAssured
	memberDataPointResult.ExpAdjPtdRiskPremium = memberDataPointResult.ExpAdjLoadedPtdRate * memberDataPointResult.PtdCappedSumAssured
	memberDataPointResult.ExpAdjTtdRiskPremium = memberDataPointResult.ExpAdjLoadedTtdRate * memberDataPointResult.TtdCappedIncome * groupParameter.TtdNumberMonthlyPayments
	memberDataPointResult.ExpAdjPhiRiskPremium = memberDataPointResult.ExpAdjLoadedPhiRate * memberDataPointResult.PhiMonthlyBenefit
	memberDataPointResult.ExpAdjCiRiskPremium = memberDataPointResult.ExpAdjLoadedCiRate * memberDataPointResult.CiCappedSumAssured
	memberDataPointResult.ExpAdjSpouseGlaRiskPremium = memberDataPointResult.ExpAdjLoadedSpouseGlaRate * memberDataPointResult.SpouseGlaCappedSumAssured

	memberDataPointResult.ExpAdjTotalFuneralRiskPremium = memberDataPointResult.GlaExperienceAdjustment * (memberDataPointResult.MainMemberFuneralRiskPremium + memberDataPointResult.SpouseFuneralRiskPremium + memberDataPointResult.ChildFuneralRiskPremium + memberDataPointResult.ParentFuneralRiskPremium) * (1 + memberDataPointResult.FunConversionOnWithdrawalLoading + memberDataPointResult.FunSchemeSizeLoading)
	memberDataPointResult.ExpAdjTotalFuneralOfficePremium = memberDataPointResult.ExpAdjTotalFuneralRiskPremium / (1.0 - memberDataPointResult.TotalPremiumLoading)
	memberDataPointResult.FinalTotalFuneralOfficePremium = memberDataPointResult.ExpAdjTotalFuneralRiskPremium / (1.0 - (memberDataPointResult.TotalPremiumLoading + discountFraction))

	// FUN override (replaces ExpAdjTotalFuneralRiskPremium and re-derives
	// the office/final-office premiums). Sits here so it runs after the
	// initial ExpAdj funeral computation but before the educator and
	// conversion cascades read these values.
	if groupQuote.ExperienceRating == "Override" {
		applyFuneralExperienceRateOverride(&memberDataPointResult, experienceRateOverrides, groupQuote.SchemeCategories[i], discountFraction)
	}

	// Re-derive the educator loaded rates now that ExpAdjLoaded*Rate values
	// are set so the ExpAdjLoadedGla/PtdEducatorRate fields are accurate.
	computeEducatorLoadedRates(&memberDataPointResult)

	if groupQuote.SchemeCategories[i].GlaEducatorBenefit == "Yes" || groupQuote.SchemeCategories[i].PtdEducatorBenefit == "Yes" { //groupQuote.Educator
		educatorRates := GetEducatorRate(groupParameter, educatorCodeForCategory(&groupQuote.SchemeCategories[i]), memberDataPointResult.AgeNextBirthday, memberDataPointResult.IncomeLevel)
		// ExpAdj split uses ExpAdjLoadedGla/PtdEducatorRate (the pre-split
		// implementation used ExpAdjLoadedTtdRate by mistake).

		if groupQuote.SchemeCategories[i].GlaEducatorBenefit == "Yes" {
			memberDataPointResult.ExpAdjGlaEducatorRiskPremium = educatorRates.EducatorSumAtRisk * memberDataPointResult.ExpAdjLoadedGlaEducatorRate
		}
		if groupQuote.SchemeCategories[i].PtdEducatorBenefit == "Yes" {
			memberDataPointResult.ExpAdjPtdEducatorRiskPremium = educatorRates.EducatorSumAtRisk * memberDataPointResult.ExpAdjLoadedPtdEducatorRate
		}
		// Recompute educator slice premiums using the finalised ExpAdj
		// educator loaded rates. This overwrites the partial values from the
		// first pass (non-ExpAdj risk + office) so that ExpAdj fields are
		// populated correctly.
		computeEducatorSlicePremiums(&memberDataPointResult, educatorRates.EducatorSumAtRisk)
	}

	// Non-educator conversion / continuity slice premiums (risk + office +
	// ExpAdj variants). Called after all LoadedRates and TotalFuneralRiskCost
	// are final, including ExpAdj* rates set in the Experience Rating block.
	computeConvContSlicePremiums(&memberDataPointResult, groupParameter)

	applyBinderOutsourceAmounts(&memberDataPointResult, &groupQuote, binderFeeRate3, outsourceFeeRate3)

	if memberDataPointResult.AgeNextBirthday > groupQuote.NormalRetirementAge {
		memberDataPointResult.ExceedsNormalRetirementAgeIndicator = 1
	}
	if memberDataPointResult.GlaSumAssured > groupQuote.FreeCoverLimit {
		memberDataPointResult.ExceedsFreeCoverLimitIndicator = 1
	}

	//memberDataPointResult.GlaExperienceAdjustedAnnualPremium = memberDataPointResult.GlaRiskPremium
	//memberDataPointResult.PtdExperienceAdjustedAnnualPremium = memberDataPointResult.PtdRiskPremium
	//memberDataPointResult.TtdExperienceAdjustedAnnualPremium = memberDataPointResult.TtdRiskPremium
	//memberDataPointResult.PhiExperienceAdjustedAnnualPremium = memberDataPointResult.PhiRiskPremium
	//memberDataPointResult.CiExperienceAdjustedAnnualPremium = memberDataPointResult.CiRiskPremium
	//memberDataPointResult.SpouseExperienceAdjustedAnnualPremium = memberDataPointResult.SpouseGlaRiskPremium
	//memberDataPointResult.FuneralExperienceAdjustedAnnualPremium = memberDataPointResult.TotalFuneralRiskCost

	memberPremiumScheduleDatapoint.Category = selectedSchemeCategory
	memberPremiumScheduleDatapoint.IsOriginalMember = true
	memberPremiumScheduleDatapoint.GlaCoveredSumAssured = memberDataPointResult.GlaCappedSumAssured
	memberPremiumScheduleDatapoint.GlaAnnualPremium = computeMemberOfficePremium(memberDataPointResult.ExpAdjGlaRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.PtdCoveredSumAssured = memberDataPointResult.PtdCappedSumAssured
	memberPremiumScheduleDatapoint.PtdAnnualPremium = computeMemberOfficePremium(memberDataPointResult.ExpAdjPtdRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.CiCoveredSumAssured = memberDataPointResult.CiCappedSumAssured
	memberPremiumScheduleDatapoint.CiAnnualPremium = computeMemberOfficePremium(memberDataPointResult.ExpAdjCiRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.TtdCoveredIncome = memberDataPointResult.TtdCappedIncome
	memberPremiumScheduleDatapoint.TtdAnnualPremium = computeMemberOfficePremium(memberDataPointResult.ExpAdjTtdRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.PhiCoveredIncome = memberDataPointResult.PhiCappedIncome
	memberPremiumScheduleDatapoint.PhiAnnualPremium = computeMemberOfficePremium(memberDataPointResult.ExpAdjPhiRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.SpouseGlaCoveredSumAssured = memberDataPointResult.SpouseGlaCappedSumAssured
	memberPremiumScheduleDatapoint.SpouseGlaAnnualPremium = computeMemberOfficePremium(memberDataPointResult.ExpAdjSpouseGlaRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.MainMemberFuneralSumAssured = groupQuote.SchemeCategories[i].FamilyFuneralMainMemberFuneralSumAssured
	memberPremiumScheduleDatapoint.MainMemberFuneralAnnualPremium = computeMemberOfficePremium(memberDataPointResult.MainMemberFuneralRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.SpouseFuneralSumAssured = groupQuote.SchemeCategories[i].FamilyFuneralSpouseFuneralSumAssured
	memberPremiumScheduleDatapoint.SpouseFuneralAnnualPremium = computeMemberOfficePremium(memberDataPointResult.SpouseFuneralRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.ChildFuneralSumAssured = memberDataPointResult.ChildFuneralSumAssured
	memberPremiumScheduleDatapoint.ChildrenFuneralAnnualPremium = computeMemberOfficePremium(memberDataPointResult.ChildFuneralRiskPremium, &groupQuote)
	memberPremiumScheduleDatapoint.DependantsFuneralSumAssured = memberDataPointResult.ParentFuneralSumAssured
	memberPremiumScheduleDatapoint.DependantsFuneralAnnualPremium = computeMemberOfficePremium(memberDataPointResult.ParentFuneralRiskPremium, &groupQuote)

	memberPremiumScheduleDatapoint.TotalAnnualPremiumPayable = memberPremiumScheduleDatapoint.GlaAnnualPremium + memberPremiumScheduleDatapoint.PtdAnnualPremium + memberPremiumScheduleDatapoint.CiAnnualPremium + memberPremiumScheduleDatapoint.TtdAnnualPremium + memberPremiumScheduleDatapoint.PhiAnnualPremium + memberPremiumScheduleDatapoint.SpouseGlaAnnualPremium + memberPremiumScheduleDatapoint.MainMemberFuneralAnnualPremium + memberPremiumScheduleDatapoint.SpouseFuneralAnnualPremium + memberPremiumScheduleDatapoint.ChildrenFuneralAnnualPremium + memberPremiumScheduleDatapoint.DependantsFuneralAnnualPremium

	GroupPricingReinsurance(&memberDataPointResult, &bordereauxDatapoint, groupQuote, groupQuote.SchemeCategories[i], groupPricingReinsuranceStructure, groupParameter)

	bordereauxDatapoint.Category = selectedSchemeCategory
	bordereauxDatapoint.IsOriginalMember = true
	bordereauxDatapoint.MemberName = memberDataPointResult.MemberName
	bordereauxDatapoint.SchemeId = memberDataPointResult.SchemeId
	bordereauxDatapoint.QuoteId = memberDataPointResult.QuoteId
	bordereauxDatapoint.Gender = memberDataPointResult.Gender
	bordereauxDatapoint.AgeNextBirthday = float64(memberDataPointResult.AgeNextBirthday)
	bordereauxDatapoint.AnnualSalary = memberDataPointResult.AnnualSalary
	bordereauxDatapoint.RenewalDate = "" //groupQuote.CommencementDate
	bordereauxDatapoint.Currency = groupQuote.Currency
	bordereauxDatapoint.Industry = groupQuote.Industry
	bordereauxDatapoint.IndustryClass = groupQuote.Industry //
	bordereauxDatapoint.GlaMultiple = groupQuote.SchemeCategories[i].GlaSalaryMultiple
	bordereauxDatapoint.GlaCoveredSumAssured = memberDataPointResult.GlaCappedSumAssured
	bordereauxDatapoint.LoadedGlaRiskRate = memberDataPointResult.LoadedGlaRate
	bordereauxDatapoint.GlaRetainedRiskPremium = bordereauxDatapoint.GlaRetainedSumAssured * memberDataPointResult.LoadedGlaRate
	// Ceded premium uses the LOADED REINSURANCE rate × ceded sum assured.
	bordereauxDatapoint.GlaCededRiskPremium = bordereauxDatapoint.GlaCededSumAssured * memberDataPointResult.LoadedReinsGlaRate
	memberDataPointResult.GlaReinsurancePremium = bordereauxDatapoint.GlaCededRiskPremium

	bordereauxDatapoint.PtdMultiple = groupQuote.SchemeCategories[i].PtdSalaryMultiple
	bordereauxDatapoint.PtdCoveredSumAssured = memberDataPointResult.PtdCappedSumAssured
	bordereauxDatapoint.LoadedPtdRiskRate = memberDataPointResult.LoadedPtdRate
	bordereauxDatapoint.PtdRetainedRiskPremium = bordereauxDatapoint.PtdRetainedSumAssured * memberDataPointResult.LoadedPtdRate
	bordereauxDatapoint.PtdCededRiskPremium = bordereauxDatapoint.PtdCededSumAssured * memberDataPointResult.LoadedReinsPtdRate
	memberDataPointResult.PtdReinsurancePremium = bordereauxDatapoint.PtdCededRiskPremium

	bordereauxDatapoint.CiMultiple = groupQuote.SchemeCategories[i].CiCriticalIllnessSalaryMultiple
	bordereauxDatapoint.CiCoveredSumAssured = memberDataPointResult.CiCappedSumAssured
	bordereauxDatapoint.LoadedCiRiskRate = memberDataPointResult.LoadedCiRate
	bordereauxDatapoint.CiRetainedRiskPremium = bordereauxDatapoint.CiRetainedSumAssured * memberDataPointResult.LoadedCiRate
	bordereauxDatapoint.CiCededRiskPremium = bordereauxDatapoint.CiCededSumAssured * memberDataPointResult.LoadedReinsCiRate
	memberDataPointResult.CiReinsurancePremium = bordereauxDatapoint.CiCededRiskPremium

	bordereauxDatapoint.SglaMultiple = groupQuote.SchemeCategories[i].SglaSalaryMultiple
	bordereauxDatapoint.SglaCoveredSumAssured = memberDataPointResult.SpouseGlaCappedSumAssured
	bordereauxDatapoint.LoadedSglaRiskRate = memberDataPointResult.LoadedSpouseGlaRate
	bordereauxDatapoint.SglaRetainedRiskPremium = bordereauxDatapoint.SglaRetainedSumAssured * memberDataPointResult.LoadedSpouseGlaRate
	bordereauxDatapoint.SglaCededRiskPremium = bordereauxDatapoint.SglaCededSumAssured * memberDataPointResult.LoadedReinsSpouseGlaRate
	memberDataPointResult.SpouseGlaReinsurancePremium = bordereauxDatapoint.SglaCededRiskPremium

	bordereauxDatapoint.TtdReplacementMultiple = groupQuote.SchemeCategories[i].TtdIncomeReplacementPercentage
	bordereauxDatapoint.TtdMonthlyBenefit = memberDataPointResult.TtdCappedIncome
	bordereauxDatapoint.LoadedTtdRiskRate = memberDataPointResult.LoadedTtdRate
	bordereauxDatapoint.TtdRetainedRiskPremium = bordereauxDatapoint.TtdRetainedMonthlyBenefit * memberDataPointResult.LoadedTtdRate * groupParameter.TtdNumberMonthlyPayments
	bordereauxDatapoint.TtdCededRiskPremium = bordereauxDatapoint.TtdCededMonthlyBenefit * memberDataPointResult.LoadedReinsTtdRate * groupParameter.TtdNumberMonthlyPayments
	memberDataPointResult.TtdReinsurancePremium = bordereauxDatapoint.TtdCededRiskPremium

	bordereauxDatapoint.PhiReplacementMultiple = groupQuote.SchemeCategories[i].PhiIncomeReplacementPercentage
	bordereauxDatapoint.PhiMonthlyBenefit = memberDataPointResult.PhiCappedIncome
	bordereauxDatapoint.LoadedPhiRiskRate = memberDataPointResult.LoadedPhiRate
	bordereauxDatapoint.PhiRetainedRiskPremium = bordereauxDatapoint.PhiRetainedMonthlyBenefit * memberDataPointResult.LoadedPhiRate
	bordereauxDatapoint.PhiCededRiskPremium = bordereauxDatapoint.PhiCededMonthlyBenefit * memberDataPointResult.LoadedReinsPhiRate
	memberDataPointResult.PhiReinsurancePremium = bordereauxDatapoint.PhiCededRiskPremium

	bordereauxDatapoint.MainMemberFuneralSumAssured = groupQuote.SchemeCategories[i].FamilyFuneralMainMemberFuneralSumAssured
	bordereauxDatapoint.MainMemberRiskRate = memberDataPointResult.LoadedGlaRate
	bordereauxDatapoint.MainMemberRetainedRiskPremium = bordereauxDatapoint.MainMemberRetainedSumAssured * memberDataPointResult.LoadedGlaRate
	bordereauxDatapoint.MainMemberCededRiskPremium = bordereauxDatapoint.MainMemberCededSumAssured * memberDataPointResult.MainMemberReinsuranceRate
	memberDataPointResult.MainMemberReinsurancePremium = bordereauxDatapoint.MainMemberCededRiskPremium

	bordereauxDatapoint.SpouseFuneralSumAssured = groupQuote.SchemeCategories[i].FamilyFuneralSpouseFuneralSumAssured
	bordereauxDatapoint.SpouseRiskRate = memberDataPointResult.LoadedSpouseGlaRate
	bordereauxDatapoint.SpouseRetainedRiskPremium = bordereauxDatapoint.SpouseRetainedSumAssured * memberDataPointResult.LoadedSpouseGlaRate
	bordereauxDatapoint.SpouseCededRiskPremium = bordereauxDatapoint.SpouseCededSumAssured * memberDataPointResult.SpouseReinsuranceRate
	memberDataPointResult.SpouseReinsurancePremium = bordereauxDatapoint.SpouseCededRiskPremium

	bordereauxDatapoint.ChildFuneralSumAssured = groupQuote.SchemeCategories[i].FamilyFuneralChildrenFuneralSumAssured
	bordereauxDatapoint.ChildRiskRate = memberDataPointResult.ChildFuneralBaseRate
	bordereauxDatapoint.ChildRetainedRiskPremium = bordereauxDatapoint.ChildRetainedSumAssured * memberDataPointResult.ChildFuneralBaseRate
	bordereauxDatapoint.ChildCededRiskPremium = bordereauxDatapoint.ChildCededSumAssured * memberDataPointResult.ChildReinsuranceRate
	memberDataPointResult.ChildReinsurancePremium = bordereauxDatapoint.ChildCededRiskPremium

	bordereauxDatapoint.ParentFuneralSumAssured = groupQuote.SchemeCategories[i].FamilyFuneralParentFuneralSumAssured
	bordereauxDatapoint.ParentRiskRate = memberDataPointResult.ParentFuneralBaseRate
	bordereauxDatapoint.ParentRetainedRiskPremium = bordereauxDatapoint.ParentRetainedSumAssured * memberDataPointResult.ParentFuneralBaseRate
	bordereauxDatapoint.ParentCededRiskPremium = bordereauxDatapoint.ParentCededSumAssured * memberDataPointResult.ParentReinsuranceRate
	memberDataPointResult.ParentReinsurancePremium = bordereauxDatapoint.ParentCededRiskPremium

	bordereauxDatapoint.DependantFuneralSumAssured = groupQuote.SchemeCategories[i].FamilyFuneralAdultDependantSumAssured
	bordereauxDatapoint.DependantRiskRate = memberDataPointResult.ParentFuneralBaseRate
	bordereauxDatapoint.DependantRetainedRiskPremium = bordereauxDatapoint.DependantRetainedSumAssured * memberDataPointResult.ParentFuneralBaseRate
	bordereauxDatapoint.DependantCededRiskPremium = bordereauxDatapoint.DependantCededSumAssured * memberDataPointResult.DependantReinsuranceRate
	memberDataPointResult.DependantReinsurancePremium = bordereauxDatapoint.DependantCededRiskPremium

	return MemberRateResult{
		Rating:          memberDataPointResult,
		PremiumSchedule: memberPremiumScheduleDatapoint,
		Bordereaux:      bordereauxDatapoint,
	}
}

func UpdateRatedMemberData(mpIndex int, memberDataPointResult models.MemberRatingResult, memberDataResults *map[string]models.MemberRatingResult) {
	mutex.Lock()
	defer mutex.Unlock()
	// TODO: This is a hack. The key should be unique.
	key := strconv.Itoa(mpIndex) + "_" + strconv.Itoa(memberDataPointResult.SchemeId)
	agg, exists := (*memberDataResults)[key]
	if exists {
		(*memberDataResults)[key] = agg
	} else {
		(*memberDataResults)[key] = memberDataPointResult
	}
}

func GroupPricingReinsurance(memberDataPointResult *models.MemberRatingResult, bordereaux *models.Bordereaux, groupQuote models.GroupPricingQuote, schemeCategory models.SchemeCategory, groupPricingReinsuranceStructure models.GroupPricingReinsuranceStructure, groupParameters models.GroupPricingParameters) {
	var glalevel1value, glalevel2value, glalevel3value float64
	var ptdlevel1value, ptdlevel2value, ptdlevel3value float64
	var cilevel1value, cilevel2value, cilevel3value float64
	var sglalevel1value, sglalevel2value, sglalevel3value float64
	var ttdlevel1value, ttdlevel2value, ttdlevel3value float64
	var philevel1value, philevel2value, philevel3value float64
	var glaretainedpropotion float64

	glalevel1value = math.Min(memberDataPointResult.GlaCappedSumAssured, groupPricingReinsuranceStructure.Level1Upperbound-groupPricingReinsuranceStructure.Level1Lowerbound)
	glalevel2value = math.Min(memberDataPointResult.GlaCappedSumAssured-glalevel1value, groupPricingReinsuranceStructure.Level2Upperbound-groupPricingReinsuranceStructure.Level2Lowerbound)
	glalevel3value = math.Min(memberDataPointResult.GlaCappedSumAssured-glalevel1value-glalevel2value, groupPricingReinsuranceStructure.Level3Upperbound-groupPricingReinsuranceStructure.Level3Lowerbound)
	bordereaux.GlaCededSumAssured = glalevel1value*groupPricingReinsuranceStructure.Level1CededProportion + glalevel2value*groupPricingReinsuranceStructure.Level2CededProportion + glalevel3value*groupPricingReinsuranceStructure.Level3CededProportion
	bordereaux.GlaRetainedSumAssured = math.Max(memberDataPointResult.GlaCappedSumAssured-bordereaux.GlaCededSumAssured, 0)
	if memberDataPointResult.GlaCappedSumAssured > 0 {
		glaretainedpropotion = bordereaux.GlaRetainedSumAssured / memberDataPointResult.GlaCappedSumAssured
	}

	if groupParameters.IsLumpsumReinsGLADependent {
		bordereaux.PtdRetainedSumAssured = memberDataPointResult.PtdCappedSumAssured * glaretainedpropotion
		bordereaux.PtdCededSumAssured = math.Max(memberDataPointResult.PtdCappedSumAssured-bordereaux.PtdRetainedSumAssured, 0)

		bordereaux.CiRetainedSumAssured = memberDataPointResult.CiCappedSumAssured * glaretainedpropotion
		bordereaux.CiCededSumAssured = math.Max(memberDataPointResult.CiCappedSumAssured-bordereaux.CiRetainedSumAssured, 0)

		bordereaux.SglaRetainedSumAssured = memberDataPointResult.SpouseGlaCappedSumAssured * glaretainedpropotion
		bordereaux.SglaCededSumAssured = math.Max(memberDataPointResult.SpouseGlaCappedSumAssured-bordereaux.SglaRetainedSumAssured, 0)
	}

	if !groupParameters.IsLumpsumReinsGLADependent {
		ptdlevel1value = math.Min(memberDataPointResult.PtdCappedSumAssured, groupPricingReinsuranceStructure.Level1Upperbound-groupPricingReinsuranceStructure.Level1Lowerbound)
		ptdlevel2value = math.Min(memberDataPointResult.PtdCappedSumAssured-ptdlevel1value, groupPricingReinsuranceStructure.Level2Upperbound-groupPricingReinsuranceStructure.Level2Lowerbound)
		ptdlevel3value = math.Min(memberDataPointResult.PtdCappedSumAssured-ptdlevel1value-ptdlevel2value, groupPricingReinsuranceStructure.Level3Upperbound-groupPricingReinsuranceStructure.Level3Lowerbound)
		bordereaux.PtdCededSumAssured = ptdlevel1value*groupPricingReinsuranceStructure.Level1CededProportion + ptdlevel2value*groupPricingReinsuranceStructure.Level2CededProportion + ptdlevel3value*groupPricingReinsuranceStructure.Level3CededProportion
		bordereaux.PtdRetainedSumAssured = math.Max(memberDataPointResult.PtdCappedSumAssured-bordereaux.PtdCededSumAssured, 0)

		cilevel1value = math.Min(memberDataPointResult.CiCappedSumAssured, groupPricingReinsuranceStructure.Level1Upperbound-groupPricingReinsuranceStructure.Level1Lowerbound)
		cilevel2value = math.Min(memberDataPointResult.CiCappedSumAssured-cilevel1value, groupPricingReinsuranceStructure.Level2Upperbound-groupPricingReinsuranceStructure.Level2Lowerbound)
		cilevel3value = math.Min(memberDataPointResult.CiCappedSumAssured-cilevel1value-cilevel2value, groupPricingReinsuranceStructure.Level3Upperbound-groupPricingReinsuranceStructure.Level3Lowerbound)
		bordereaux.CiCededSumAssured = cilevel1value*groupPricingReinsuranceStructure.Level1CededProportion + cilevel2value*groupPricingReinsuranceStructure.Level2CededProportion + cilevel3value*groupPricingReinsuranceStructure.Level3CededProportion
		bordereaux.CiRetainedSumAssured = math.Max(memberDataPointResult.CiCappedSumAssured-bordereaux.CiCededSumAssured, 0)

		sglalevel1value = math.Min(memberDataPointResult.SpouseGlaCappedSumAssured, groupPricingReinsuranceStructure.Level1Upperbound-groupPricingReinsuranceStructure.Level1Lowerbound)
		sglalevel2value = math.Min(memberDataPointResult.SpouseGlaCappedSumAssured-sglalevel1value, groupPricingReinsuranceStructure.Level2Upperbound-groupPricingReinsuranceStructure.Level2Lowerbound)
		sglalevel3value = math.Min(memberDataPointResult.SpouseGlaCappedSumAssured-sglalevel1value-sglalevel2value, groupPricingReinsuranceStructure.Level3Upperbound-groupPricingReinsuranceStructure.Level3Lowerbound)
		bordereaux.SglaCededSumAssured = sglalevel1value*groupPricingReinsuranceStructure.Level1CededProportion + sglalevel2value*groupPricingReinsuranceStructure.Level2CededProportion + sglalevel3value*groupPricingReinsuranceStructure.Level3CededProportion
		bordereaux.SglaRetainedSumAssured = math.Max(memberDataPointResult.SpouseGlaCappedSumAssured-bordereaux.SglaCededSumAssured, 0)
	}

	ttdlevel1value = math.Min(memberDataPointResult.TtdCappedIncome, groupPricingReinsuranceStructure.IncomeLevel1Upperbound-groupPricingReinsuranceStructure.IncomeLevel1Lowerbound)
	ttdlevel2value = math.Min(memberDataPointResult.TtdCappedIncome-ttdlevel1value, groupPricingReinsuranceStructure.IncomeLevel2Upperbound-groupPricingReinsuranceStructure.IncomeLevel2Lowerbound)
	ttdlevel3value = math.Min(memberDataPointResult.TtdCappedIncome-ttdlevel1value-ttdlevel2value, groupPricingReinsuranceStructure.IncomeLevel3Upperbound-groupPricingReinsuranceStructure.IncomeLevel3Lowerbound)
	bordereaux.TtdCededMonthlyBenefit = ttdlevel1value*groupPricingReinsuranceStructure.IncomeLevel1CededProportion + ttdlevel2value*groupPricingReinsuranceStructure.IncomeLevel2CededProportion + ttdlevel3value*groupPricingReinsuranceStructure.IncomeLevel3CededProportion
	bordereaux.TtdRetainedMonthlyBenefit = math.Max(memberDataPointResult.TtdCappedIncome-bordereaux.TtdCededMonthlyBenefit, 0)

	philevel1value = math.Min(memberDataPointResult.PhiCappedIncome, groupPricingReinsuranceStructure.IncomeLevel1Upperbound-groupPricingReinsuranceStructure.IncomeLevel1Lowerbound)
	philevel2value = math.Min(memberDataPointResult.PhiCappedIncome-philevel1value, groupPricingReinsuranceStructure.IncomeLevel2Upperbound-groupPricingReinsuranceStructure.IncomeLevel2Lowerbound)
	philevel3value = math.Min(memberDataPointResult.PhiCappedIncome-philevel1value-philevel2value, groupPricingReinsuranceStructure.IncomeLevel3Upperbound-groupPricingReinsuranceStructure.IncomeLevel3Lowerbound)
	bordereaux.PhiCededMonthlyBenefit = philevel1value*groupPricingReinsuranceStructure.IncomeLevel1CededProportion + philevel2value*groupPricingReinsuranceStructure.IncomeLevel2CededProportion + philevel3value*groupPricingReinsuranceStructure.IncomeLevel3CededProportion
	bordereaux.PhiRetainedMonthlyBenefit = math.Max(memberDataPointResult.PhiCappedIncome-bordereaux.PhiCededMonthlyBenefit, 0)

	if groupPricingReinsuranceStructure.FuneralReinsuranceInclusionIndicator {
		mmlevel1value := math.Min(schemeCategory.FamilyFuneralMainMemberFuneralSumAssured, groupPricingReinsuranceStructure.Level1Upperbound-groupPricingReinsuranceStructure.Level1Lowerbound)
		mmlevel2value := math.Min(schemeCategory.FamilyFuneralMainMemberFuneralSumAssured-mmlevel1value, groupPricingReinsuranceStructure.Level2Upperbound-groupPricingReinsuranceStructure.Level2Lowerbound)
		mmlevel3value := math.Min(schemeCategory.FamilyFuneralMainMemberFuneralSumAssured-mmlevel1value-mmlevel2value, groupPricingReinsuranceStructure.Level3Upperbound-groupPricingReinsuranceStructure.Level3Lowerbound)
		bordereaux.MainMemberCededSumAssured = mmlevel1value*groupPricingReinsuranceStructure.Level1CededProportion + mmlevel2value*groupPricingReinsuranceStructure.Level2CededProportion + mmlevel3value*groupPricingReinsuranceStructure.Level3CededProportion

		splevel1value := math.Min(schemeCategory.FamilyFuneralSpouseFuneralSumAssured, groupPricingReinsuranceStructure.Level1Upperbound-groupPricingReinsuranceStructure.Level1Lowerbound)
		splevel2value := math.Min(schemeCategory.FamilyFuneralSpouseFuneralSumAssured-splevel1value, groupPricingReinsuranceStructure.Level2Upperbound-groupPricingReinsuranceStructure.Level2Lowerbound)
		splevel3value := math.Min(schemeCategory.FamilyFuneralSpouseFuneralSumAssured-splevel1value-splevel2value, groupPricingReinsuranceStructure.Level3Upperbound-groupPricingReinsuranceStructure.Level3Lowerbound)
		bordereaux.SpouseCededSumAssured = splevel1value*groupPricingReinsuranceStructure.Level1CededProportion + splevel2value*groupPricingReinsuranceStructure.Level2CededProportion + splevel3value*groupPricingReinsuranceStructure.Level3CededProportion

		chlevel1value := math.Min(schemeCategory.FamilyFuneralChildrenFuneralSumAssured, groupPricingReinsuranceStructure.Level1Upperbound-groupPricingReinsuranceStructure.Level1Lowerbound)
		chlevel2value := math.Min(schemeCategory.FamilyFuneralChildrenFuneralSumAssured-chlevel1value, groupPricingReinsuranceStructure.Level2Upperbound-groupPricingReinsuranceStructure.Level2Lowerbound)
		chlevel3value := math.Min(schemeCategory.FamilyFuneralChildrenFuneralSumAssured-chlevel1value-chlevel2value, groupPricingReinsuranceStructure.Level3Upperbound-groupPricingReinsuranceStructure.Level3Lowerbound)
		bordereaux.ChildCededSumAssured = chlevel1value*groupPricingReinsuranceStructure.Level1CededProportion + chlevel2value*groupPricingReinsuranceStructure.Level2CededProportion + chlevel3value*groupPricingReinsuranceStructure.Level3CededProportion

		parlevel1value := math.Min(schemeCategory.FamilyFuneralParentFuneralSumAssured, groupPricingReinsuranceStructure.Level1Upperbound-groupPricingReinsuranceStructure.Level1Lowerbound)
		parlevel2value := math.Min(schemeCategory.FamilyFuneralParentFuneralSumAssured-parlevel1value, groupPricingReinsuranceStructure.Level2Upperbound-groupPricingReinsuranceStructure.Level2Lowerbound)
		parlevel3value := math.Min(schemeCategory.FamilyFuneralParentFuneralSumAssured-parlevel1value-parlevel2value, groupPricingReinsuranceStructure.Level3Upperbound-groupPricingReinsuranceStructure.Level3Lowerbound)
		bordereaux.ParentCededSumAssured = parlevel1value*groupPricingReinsuranceStructure.Level1CededProportion + parlevel2value*groupPricingReinsuranceStructure.Level2CededProportion + parlevel3value*groupPricingReinsuranceStructure.Level3CededProportion

		deplevel1value := math.Min(schemeCategory.FamilyFuneralAdultDependantSumAssured, groupPricingReinsuranceStructure.Level1Upperbound-groupPricingReinsuranceStructure.Level1Lowerbound)
		deplevel2value := math.Min(schemeCategory.FamilyFuneralAdultDependantSumAssured-deplevel1value, groupPricingReinsuranceStructure.Level2Upperbound-groupPricingReinsuranceStructure.Level2Lowerbound)
		deplevel3value := math.Min(schemeCategory.FamilyFuneralAdultDependantSumAssured-deplevel1value-deplevel2value, groupPricingReinsuranceStructure.Level3Upperbound-groupPricingReinsuranceStructure.Level3Lowerbound)
		bordereaux.DependantCededSumAssured = deplevel1value*groupPricingReinsuranceStructure.Level1CededProportion + deplevel2value*groupPricingReinsuranceStructure.Level2CededProportion + deplevel3value*groupPricingReinsuranceStructure.Level3CededProportion
	}
	if !groupPricingReinsuranceStructure.FuneralReinsuranceInclusionIndicator {
		bordereaux.MainMemberCededSumAssured = 0
		bordereaux.SpouseCededSumAssured = 0
		bordereaux.ChildCededSumAssured = 0
		bordereaux.ParentCededSumAssured = 0
		bordereaux.DependantCededSumAssured = 0
	}
	bordereaux.MainMemberRetainedSumAssured = math.Max(schemeCategory.FamilyFuneralMainMemberFuneralSumAssured-bordereaux.MainMemberCededSumAssured, 0)
	bordereaux.SpouseRetainedSumAssured = math.Max(schemeCategory.FamilyFuneralSpouseFuneralSumAssured-bordereaux.SpouseCededSumAssured, 0)
	bordereaux.ChildRetainedSumAssured = math.Max(schemeCategory.FamilyFuneralChildrenFuneralSumAssured-bordereaux.ChildCededSumAssured, 0)
	bordereaux.ParentRetainedSumAssured = math.Max(schemeCategory.FamilyFuneralParentFuneralSumAssured-bordereaux.ParentCededSumAssured, 0)
	bordereaux.DependantRetainedSumAssured = math.Max(schemeCategory.FamilyFuneralAdultDependantSumAssured-bordereaux.DependantCededSumAssured, 0)
}

// BuildBordereauxFromRatingResult projects a Bordereaux row from a persisted
// MemberRatingResult plus the quote/category/reinsurance context. It is the
// single source of truth for the field-mapping from rating outputs to
// bordereaux columns and is used by both quote-stage previews and on-the-fly
// downloads — replacing the previously persisted `bordereauxes` table.
func BuildBordereauxFromRatingResult(
	mrr models.MemberRatingResult,
	schemeCategory models.SchemeCategory,
	groupQuote models.GroupPricingQuote,
	reinsStructure models.GroupPricingReinsuranceStructure,
	groupParameter models.GroupPricingParameters,
) models.Bordereaux {
	var b models.Bordereaux

	// Identity / spine fields.
	b.Category = mrr.Category
	b.IsOriginalMember = mrr.IsOriginalMember
	b.MemberName = mrr.MemberName
	b.SchemeId = mrr.SchemeId
	b.QuoteId = mrr.QuoteId
	b.Gender = mrr.Gender
	b.AgeNextBirthday = float64(mrr.AgeNextBirthday)
	b.AnnualSalary = mrr.AnnualSalary
	if !mrr.DateOfBirth.IsZero() {
		dob := mrr.DateOfBirth
		b.DateOfBirth = &dob
	}
	b.RenewalDate = ""
	b.Currency = groupQuote.Currency
	b.Industry = groupQuote.Industry
	b.IndustryClass = groupQuote.Industry

	// Reinsurance splits — populates retained/ceded sums assured on `b`.
	GroupPricingReinsurance(&mrr, &b, groupQuote, schemeCategory, reinsStructure, groupParameter)

	// GLA
	b.GlaMultiple = schemeCategory.GlaSalaryMultiple
	b.GlaCoveredSumAssured = mrr.GlaCappedSumAssured
	b.LoadedGlaRiskRate = mrr.LoadedGlaRate
	b.ExpAdjLoadedGlaRiskRate = mrr.ExpAdjLoadedGlaRate
	b.GlaRetainedRiskPremium = b.GlaRetainedSumAssured * mrr.LoadedGlaRate
	b.GlaCededRiskPremium = b.GlaCededSumAssured * mrr.LoadedReinsGlaRate
	b.GlaAnnualPremium = models.ComputeOfficePremium(b.GlaRetainedRiskPremium, mrrSummaryShim(&mrr))

	// PTD
	b.PtdMultiple = schemeCategory.PtdSalaryMultiple
	b.PtdCoveredSumAssured = mrr.PtdCappedSumAssured
	b.LoadedPtdRiskRate = mrr.LoadedPtdRate
	b.ExpAdjLoadedPtdRiskRate = mrr.ExpAdjLoadedPtdRate
	b.PtdRetainedRiskPremium = b.PtdRetainedSumAssured * mrr.LoadedPtdRate
	b.PtdCededRiskPremium = b.PtdCededSumAssured * mrr.LoadedReinsPtdRate
	b.PtdAnnualPremium = models.ComputeOfficePremium(b.PtdRetainedRiskPremium, mrrSummaryShim(&mrr))

	// CI
	b.CiMultiple = schemeCategory.CiCriticalIllnessSalaryMultiple
	b.CiCoveredSumAssured = mrr.CiCappedSumAssured
	b.LoadedCiRiskRate = mrr.LoadedCiRate
	b.ExpAdjLoadedCiRiskRate = mrr.ExpAdjLoadedCiRate
	b.CiRetainedRiskPremium = b.CiRetainedSumAssured * mrr.LoadedCiRate
	b.CiCededRiskPremium = b.CiCededSumAssured * mrr.LoadedReinsCiRate
	b.CiAnnualPremium = models.ComputeOfficePremium(b.CiRetainedRiskPremium, mrrSummaryShim(&mrr))

	// SGLA (Spouse GLA)
	b.SglaMultiple = schemeCategory.SglaSalaryMultiple
	b.SglaCoveredSumAssured = mrr.SpouseGlaCappedSumAssured
	b.LoadedSglaRiskRate = mrr.LoadedSpouseGlaRate
	b.ExpAdjLoadedSglaRiskRate = mrr.ExpAdjLoadedSpouseGlaRate
	b.SglaRetainedRiskPremium = b.SglaRetainedSumAssured * mrr.LoadedSpouseGlaRate
	b.SglaCededRiskPremium = b.SglaCededSumAssured * mrr.LoadedReinsSpouseGlaRate
	b.SglaAnnualPremium = models.ComputeOfficePremium(b.SglaRetainedRiskPremium, mrrSummaryShim(&mrr))

	// TTD
	b.TtdReplacementMultiple = schemeCategory.TtdIncomeReplacementPercentage
	b.TtdMonthlyBenefit = mrr.TtdCappedIncome
	b.LoadedTtdRiskRate = mrr.LoadedTtdRate
	b.ExpAdjLoadedTtdRiskRate = mrr.ExpAdjLoadedTtdRate
	b.TtdRetainedRiskPremium = b.TtdRetainedMonthlyBenefit * mrr.LoadedTtdRate * groupParameter.TtdNumberMonthlyPayments
	b.TtdCededRiskPremium = b.TtdCededMonthlyBenefit * mrr.LoadedReinsTtdRate * groupParameter.TtdNumberMonthlyPayments

	// PHI
	b.PhiReplacementMultiple = schemeCategory.PhiIncomeReplacementPercentage
	b.PhiMonthlyBenefit = mrr.PhiCappedIncome
	b.LoadedPhiRiskRate = mrr.LoadedPhiRate
	b.ExpAdjLoadedPhiRiskRate = mrr.ExpAdjLoadedPhiRate
	b.PhiRetainedRiskPremium = b.PhiRetainedMonthlyBenefit * mrr.LoadedPhiRate
	b.PhiCededRiskPremium = b.PhiCededMonthlyBenefit * mrr.LoadedReinsPhiRate

	// Family funeral lives — sums assured come from scheme category, rates
	// come from member rating result (per-life rates on the member).
	b.MainMemberFuneralSumAssured = schemeCategory.FamilyFuneralMainMemberFuneralSumAssured
	b.MainMemberRiskRate = mrr.LoadedGlaRate
	b.MainMemberRetainedRiskPremium = b.MainMemberRetainedSumAssured * mrr.LoadedGlaRate
	b.MainMemberCededRiskPremium = b.MainMemberCededSumAssured * mrr.MainMemberReinsuranceRate

	b.SpouseFuneralSumAssured = schemeCategory.FamilyFuneralSpouseFuneralSumAssured
	b.SpouseRiskRate = mrr.LoadedSpouseGlaRate
	b.SpouseRetainedRiskPremium = b.SpouseRetainedSumAssured * mrr.LoadedSpouseGlaRate
	b.SpouseCededRiskPremium = b.SpouseCededSumAssured * mrr.SpouseReinsuranceRate

	b.ChildFuneralSumAssured = schemeCategory.FamilyFuneralChildrenFuneralSumAssured
	b.ChildRiskRate = mrr.ChildFuneralBaseRate
	b.ChildRetainedRiskPremium = b.ChildRetainedSumAssured * mrr.ChildFuneralBaseRate
	b.ChildCededRiskPremium = b.ChildCededSumAssured * mrr.ChildReinsuranceRate

	b.ParentFuneralSumAssured = schemeCategory.FamilyFuneralParentFuneralSumAssured
	b.ParentRiskRate = mrr.ParentFuneralBaseRate
	b.ParentRetainedRiskPremium = b.ParentRetainedSumAssured * mrr.ParentFuneralBaseRate
	b.ParentCededRiskPremium = b.ParentCededSumAssured * mrr.ParentReinsuranceRate

	b.DependantFuneralSumAssured = schemeCategory.FamilyFuneralAdultDependantSumAssured
	b.DependantRiskRate = mrr.ParentFuneralBaseRate
	b.DependantRetainedRiskPremium = b.DependantRetainedSumAssured * mrr.ParentFuneralBaseRate
	b.DependantCededRiskPremium = b.DependantCededSumAssured * mrr.DependantReinsuranceRate

	return b
}

// mrrSummaryShim adapts a MemberRatingResult into the loadings-bearing struct
// expected by ComputeOfficePremium. The member's per-row loadings are an
// authoritative snapshot of what was used to produce the persisted premiums,
// so the gross-up matches the calculation engine without requiring an extra
// summary lookup at projection time.
func mrrSummaryShim(m *models.MemberRatingResult) *models.MemberRatingResultSummary {
	return &models.MemberRatingResultSummary{
		ExpenseLoading:    m.ExpenseLoading,
		AdminLoading:      m.AdminLoading,
		CommissionLoading: m.CommissionLoading,
		ProfitLoading:     m.ProfitLoading,
		OtherLoading:      m.OtherLoading,
		BinderFeeRate:     m.BinderFeeRate,
		OutsourceFeeRate:  m.OutsourceFeeRate,
		Discount:          m.Discount,
	}
}

// BuildBordereauxRowsForQuote projects bordereaux rows for a quote on-the-fly
// from persisted MemberRatingResult rows. It loads the quote's reinsurance
// structure and group parameters once and applies them to every row, so the
// returned bordereaux always reflects the latest loadings/cession structure.
//
// Pagination matches the previous DB-backed read path: pass limit=0 to fetch
// all rows.
func BuildBordereauxRowsForQuote(quoteId int, offset, limit int) ([]models.Bordereaux, error) {
	var quote models.GroupPricingQuote
	if err := DB.Preload("SchemeCategories").Where("id = ?", quoteId).First(&quote).Error; err != nil {
		return nil, fmt.Errorf("load quote %d: %w", quoteId, err)
	}

	var groupParameter models.GroupPricingParameters
	if err := DB.Where("basis = ? and risk_rate_code = ?", quote.Basis, quote.RiskRateCode).
		First(&groupParameter).Error; err != nil {
		// A missing group parameter only zeroes out parameter-dependent fields;
		// it should not block the projection.
		appLog.WithField("quote_id", quoteId).
			WithField("error", err.Error()).
			Warn("BuildBordereauxRowsForQuote: group parameter not found; continuing with zero values")
	}

	var reinsStructure models.GroupPricingReinsuranceStructure
	if IsTableRequired("gpReinsuranceStructure") {
		if err := DB.Where("risk_rate_code = ? and basis = ?", quote.RiskRateCode, quote.Basis).
			First(&reinsStructure).Error; err != nil {
			appLog.WithField("quote_id", quoteId).
				WithField("error", err.Error()).
				Warn("BuildBordereauxRowsForQuote: reinsurance structure not found; cession will be zero")
		}
	}

	categoryByName := make(map[string]models.SchemeCategory, len(quote.SchemeCategories))
	for _, sc := range quote.SchemeCategories {
		categoryByName[sc.SchemeCategory] = sc
	}

	// MemberRatingResult has no auto-increment primary key, so we follow the
	// same pattern as the existing case "member_rating_results" reader and
	// rely on insertion order for stable pagination within a single quote.
	q := DB.Where("quote_id = ?", quoteId)
	if limit > 0 {
		q = q.Offset(offset).Limit(limit)
	}
	var ratings []models.MemberRatingResult
	if err := q.Find(&ratings).Error; err != nil {
		return nil, fmt.Errorf("load member rating results for quote %d: %w", quoteId, err)
	}

	out := make([]models.Bordereaux, 0, len(ratings))
	for _, r := range ratings {
		sc := categoryByName[r.Category]
		out = append(out, BuildBordereauxFromRatingResult(r, sc, quote, reinsStructure, groupParameter))
	}
	return out, nil
}

// quoteBordereauxHiddenColumns lists json tags that are populated on the
// projected Bordereaux row but intentionally suppressed from the quote-stage
// preview and xlsx export. The data still drives downstream calculations
// (e.g. annual office premiums grossed up via ComputeOfficePremium); it just
// adds clutter for the user reviewing per-member sums assured and ceded
// premiums during pricing.
var quoteBordereauxHiddenColumns = map[string]struct{}{
	"gla_annual_premium":                {},
	"ptd_annual_premium":                {},
	"gla_retained_risk_premium":         {},
	"loaded_gla_risk_rate":              {},
	"exp_adj_loaded_gla_risk_rate":      {},
	"loaded_ptd_risk_rate":              {},
	"exp_adj_loaded_ptd_risk_rate":      {},
	"ptd_retained_risk_premium":         {},
	"loaded_ci_risk_rate":               {},
	"exp_adj_loaded_ci_risk_rate":       {},
	"ci_retained_risk_premium":          {},
	"ci_annual_premium":                 {},
	"loaded_sgla_risk_rate":             {},
	"exp_adj_loaded_sgla_risk_rate":     {},
	"sgla_retained_risk_premium":        {},
	"sgla_annual_premium":               {},
	"loaded_ttd_risk_rate":              {},
	"exp_adj_loaded_ttd_risk_rate":      {},
	"ttd_retained_risk_premium":         {},
	"loaded_phi_risk_rate":              {},
	"exp_adj_loaded_phi_risk_rate":      {},
	"phi_retained_risk_premium":         {},
	"main_member_risk_rate":             {},
	"main_member_retained_risk_premium": {},
	"spouse_risk_rate":                  {},
	"spouse_retained_risk_premium":      {},
	"child_risk_rate":                   {},
	"child_retained_risk_premium":       {},
	"parent_risk_rate":                  {},
	"parent_retained_risk_premium":      {},
}

// QuoteBordereauxJSONTags returns the json tag list for a quote-stage
// bordereaux row, with hidden columns filtered out. The shared list is used
// by the JSON preview path (so the grid hides the columns) and the xlsx
// export (so headers and values stay aligned).
func QuoteBordereauxJSONTags() []string {
	all := getJSONTags(models.Bordereaux{})
	out := make([]string, 0, len(all))
	for _, t := range all {
		if _, hide := quoteBordereauxHiddenColumns[t]; hide {
			continue
		}
		out = append(out, t)
	}
	return out
}

// ExportQuoteBordereauxXLSX projects all bordereaux rows for a quote and
// streams them to an xlsx workbook. Column headers come from the Bordereaux
// struct's json tags so the export matches the on-screen preview shape.
func ExportQuoteBordereauxXLSX(quoteId int) ([]byte, error) {
	rows, err := BuildBordereauxRowsForQuote(quoteId, 0, 0)
	if err != nil {
		return nil, err
	}

	headers := QuoteBordereauxJSONTags()

	f := excelize.NewFile()
	sheetName := f.GetSheetName(f.GetActiveSheetIndex())
	sw, err := f.NewStreamWriter(sheetName)
	if err != nil {
		return nil, err
	}

	headerRow := make([]interface{}, len(headers))
	for i, h := range headers {
		headerRow[i] = h
	}
	if err := sw.SetRow("A1", headerRow); err != nil {
		return nil, err
	}

	for idx, r := range rows {
		values := bordereauxValuesForJSONTags(r, headers)
		axis, _ := excelize.CoordinatesToCellName(1, idx+2)
		if err := sw.SetRow(axis, values); err != nil {
			return nil, err
		}
	}
	if err := sw.Flush(); err != nil {
		return nil, err
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// bordereauxValuesForJSONTags returns the values of a Bordereaux row in the
// same order as the supplied json tags, so the row aligns with the header
// row written by ExportQuoteBordereauxXLSX.
func bordereauxValuesForJSONTags(b models.Bordereaux, jsonTags []string) []interface{} {
	v := reflect.ValueOf(b)
	t := v.Type()
	indexByTag := make(map[string]int, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		if comma := strings.IndexByte(tag, ','); comma >= 0 {
			tag = tag[:comma]
		}
		indexByTag[tag] = i
	}
	out := make([]interface{}, len(jsonTags))
	for i, tag := range jsonTags {
		fieldIdx, ok := indexByTag[tag]
		if !ok {
			out[i] = ""
			continue
		}
		fv := v.Field(fieldIdx)
		switch fv.Kind() {
		case reflect.Ptr:
			if fv.IsNil() {
				out[i] = ""
			} else {
				out[i] = fmt.Sprintf("%v", fv.Elem().Interface())
			}
		case reflect.Float64, reflect.Float32:
			out[i] = fv.Float()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			out[i] = fv.Int()
		case reflect.Bool:
			out[i] = fv.Bool()
		case reflect.String:
			out[i] = fv.String()
		default:
			out[i] = fmt.Sprintf("%v", fv.Interface())
		}
	}
	return out
}

func saveFileToDB(file *multipart.FileHeader, table string, groupQuote models.GroupPricingQuote) error {
	// get delimiter
	var delimiter rune
	delimiterFile, err := file.Open()
	if err != nil {
		return err
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)
	//delimiter = ','
	// save file to database
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.TrimLeadingSpace = true
	reader.Comma = delimiter
	dec, _ := csvutil.NewDecoder(reader)
	dec.Header()

	if table == "member_data" {
		DB.Where("year=? and scheme_name=?", 2024, groupQuote.SchemeName).Delete(models.MemberRatingResult{})
		for {
			var memberData models.GPricingMemberData
			if err := dec.Decode(&memberData); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			memberData.CreationDate = time.Now()
			memberData.SchemeName = groupQuote.SchemeName
			memberData.Year = 2024
			memberData.QuoteId = groupQuote.ID
			err = DB.Create(&memberData).Error
			if err != nil {
				return err
			}
		}
	}

	if table == "claims_experience" {
		DB.Where("year=? and scheme_name=?", 2024, groupQuote.SchemeName).Delete(models.GroupPricingClaimsExperience{})
		for {
			var claimsExperience models.GroupPricingClaimsExperience
			if err := dec.Decode(&claimsExperience); err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err)
			}
			claimsExperience.CreationDate = time.Now()
			claimsExperience.SchemeName = groupQuote.SchemeName
			claimsExperience.QuoteId = groupQuote.ID
			err = DB.Create(&claimsExperience).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func convertNumbers(input interface{}) interface{} {
	switch v := input.(type) {
	case map[string]interface{}:
		for key, value := range v {
			v[key] = convertNumbers(value)
		}
		return v
	case []interface{}:
		for i, value := range v {
			v[i] = convertNumbers(value)
		}
		return v
	case string:
		// Attempt to parse the string as a number
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num
		}
		return v
	default:
		return v
	}
}

func IsTableEmpty(model interface{}) (bool, error) {
	var empty bool
	err := DB.Model(&model).
		Select("count(*) > 0").
		Find(&empty).
		Error

	if err != nil {
		return false, err
	}
	if empty {
		return true, nil
	} else {
		return false, nil
	}
}

// gpTableSpec maps a rate table's internal stat key to its display name,
// delete-endpoint slug, UI category, and a zero-value model for COUNT queries.
type gpTableSpec struct {
	statName    string
	displayName string
	deleteKey   string
	category    string
	model       interface{}
}

// gpTableSpecs is the canonical ordered list of all tracked GP rate tables.
var gpTableSpecs = []gpTableSpec{
	{"gpParameters", "Group Pricing Parameters", "grouppricingparameters", "group_pricing", models.GroupPricingParameters{}},
	{"glaRate", "Group Life Assurance", "grouplifeassurance", "group_pricing", models.GlaRate{}},
	{"ptdRate", "Permanent and Total Disability", "permanentandtotaldisability", "group_pricing", models.PtdRate{}},
	{"ciRate", "Critical Illness", "criticalillness", "group_pricing", models.CiRate{}},
	{"accidentalTtdRate", "Accidental Temporary Total Disability", "accidentaltemporarytotaldisability", "group_pricing", models.AccidentalTtdRate{}},
	{"ttdRate", "Temporary Total Disability", "temporarytotaldisability", "group_pricing", models.TtdRate{}},
	{"phiRate", "Permanent Health Insurance", "permanenthealthinsurance", "group_pricing", models.PhiRate{}},
	{"childMortality", "Child Mortality", "childmortality", "group_pricing", models.ChildMortality{}},
	{"industryLoading", "Industry Loading", "industryloading", "group_pricing", models.IndustryLoading{}},
	{"funeralParameters", "Funeral Parameters", "funeralparameters", "group_pricing", models.FuneralParameters{}},
	{"gpReinsuranceStructure", "Group Pricing Reinsurance Structure", "grouppricingreinsurancestructure", "group_pricing", models.GroupPricingReinsuranceStructure{}},
	{"incomeLevel", "Income Level", "incomelevel", "group_pricing", models.IncomeLevel{}},
	{"occupationClass", "Occupation Class", "occupationclass", "group_pricing", models.OccupationClass{}},
	{"gpEducatorBenefitStructure", "Group Pricing Educator Structure", "grouppricingeducatorstructure", "group_pricing", models.EducatorBenefitStructure{}},
	{"educatorRiskRates", "Educator Risk Rates", "educatorriskrates", "group_pricing", models.EducatorRate{}},
	{"escalations", "Escalations", "escalations", "group_pricing", models.Escalations{}},
	{"funeralAidsRate", "Funeral Aids Rates", "funeralaidsrates", "group_pricing", models.FuneralAidsRate{}},
	{"funeralRate", "Funeral Rates", "funeralrates", "group_pricing", models.FuneralRate{}},
	{"generalLoading", "General Loadings", "generalloadings", "group_pricing", models.GeneralLoading{}},
	{"glaAidsRate", "GLA Aids Rates", "glaaidsrates", "group_pricing", models.GlaAidsRate{}},
	{"regionLoading", "Region Loading", "regionloading", "group_pricing", models.RegionLoading{}},
	{"tieredIncomeReplacement", "Tiered Income Replacement", "tieredincomereplacement", "group_pricing", models.TieredIncomeReplacement{}},
	{"customTieredIncomeReplacement", "Custom Tiered Income Replacement", "customtieredincomereplacement", "group_pricing", models.CustomTieredIncomeReplacement{}},
	{"discountAuthority", "Discount Authority", "discountauthority", "group_pricing", models.DiscountAuthority{}},
	{"restriction", "Restrictions", "restrictions", "group_pricing", models.Restriction{}},
	{"premiumLoading", "Premium Loadings", "premiumloadings", "group_pricing", models.PremiumLoading{}},
	{"schemeSizeLevel", "Scheme Size Levels", "schemesizelevels", "group_pricing", models.SchemeSizeLevel{}},
	{"taxTable", "Tax Table", "taxtable", "group_pricing", models.TaxTable{}},
	{"ageBands", "Age Bands", "agebands", "group_pricing", models.GroupPricingAgeBands{}},
	{"medicalWaiver", "Medical Waiver", "medicalwaiver", "group_pricing", models.MedicalWaiver{}},
	{"reinsuranceGlaRate", "Reinsurance GLA Rate", "reinsuranceglarate", "reinsurance", models.ReinsuranceGlaRate{}},
	{"reinsuranceCiRate", "Reinsurance CI Rate", "reinsurancecirate", "reinsurance", models.ReinsuranceCiRate{}},
	{"reinsurancePtdRate", "Reinsurance PTD Rate", "reinsuranceptdrate", "reinsurance", models.ReinsurancePtdRate{}},
	{"reinsurancePhiRate", "Reinsurance PHI Rate", "reinsurancephirate", "reinsurance", models.ReinsurancePhiRate{}},
	{"reinsuranceFuneralAidsRate", "Reinsurance Funeral Aids Rate", "reinsurancefuneralaidsrates", "reinsurance", models.ReinsuranceFuneralAidsRate{}},
	{"reinsuranceFuneralRate", "Reinsurance Funeral Rate", "reinsurancefuneralrates", "reinsurance", models.ReinsuranceFuneralRate{}},
	{"reinsuranceGlaAidsRate", "Reinsurance GLA Aids Rate", "reinsuranceglaaidsrates", "reinsurance", models.ReinsuranceGlaAidsRate{}},
	{"reinsuranceGeneralLoading", "Reinsurance General Loading", "reinsurancegeneralloadings", "reinsurance", models.ReinsuranceGeneralLoading{}},
	{"reinsuranceIndustryLoading", "Reinsurance Industry Loading", "reinsuranceindustryloadings", "reinsurance", models.ReinsuranceIndustryLoading{}},
	{"reinsuranceRegionLoading", "Reinsurance Region Loading", "reinsuranceregionloadings", "reinsurance", models.ReinsuranceRegionLoading{}},
	{"reinsuranceCoverRestriction", "Reinsurance Cover Restrictions", "reinsurancecoverrestrictions", "reinsurance", models.ReinsuranceCoverRestriction{}},
}

// setGPTableStat upserts the row count for a single table into gp_table_stats.
func setGPTableStat(statName string, rowCount int64) {
	var stat models.GPTableStat
	DB.Where("table_name = ?", statName).FirstOrCreate(&stat, models.GPTableStat{TableName: statName})
	stat.RowCount = rowCount
	stat.UpdatedAt = time.Now()
	DB.Save(&stat)
}

// refreshGPTableStatByDisplayName re-counts the live table and updates its stat.
// Called after uploads (which use the display name, e.g. "Group Life Assurance").
func refreshGPTableStatByDisplayName(displayName string) {
	for _, spec := range gpTableSpecs {
		if spec.displayName == displayName {
			var count int64
			DB.Model(&spec.model).Count(&count)
			setGPTableStat(spec.statName, count)
			return
		}
	}
}

// refreshGPTableStatByDeleteKey re-counts the live table and updates its stat.
// Called after deletes (which use the lowercase slug, e.g. "grouplifeassurance").
func refreshGPTableStatByDeleteKey(deleteKey string) {
	for _, spec := range gpTableSpecs {
		if spec.deleteKey == deleteKey {
			var count int64
			DB.Model(&spec.model).Count(&count)
			setGPTableStat(spec.statName, count)
			return
		}
	}
}

// EnsureGPTableStats is called on startup. It creates the gp_table_stats table
// if it does not exist (handles existing databases that predate this feature)
// and seeds any missing stat rows. A full rebuild happens only when the
// table is completely empty; otherwise we self-heal per-spec so newly-added
// entries in gpTableSpecs (e.g. new reinsurance tables) get a correct
// populated flag on first load without a manual rebuild.
func EnsureGPTableStats() {
	if err := DB.AutoMigrate(&models.GPTableStat{}); err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to auto-migrate gp_table_stats")
		return
	}
	var count int64
	DB.Model(&models.GPTableStat{}).Count(&count)
	if count == 0 {
		appLog.Info("gp_table_stats is empty — rebuilding from live counts")
		if err := RebuildGPTableStats(); err != nil {
			appLog.WithField("error", err.Error()).Error("Failed to rebuild gp_table_stats on startup")
		}
		return
	}

	// Fill in any spec rows that are missing (idempotent per-spec seed).
	var existing []models.GPTableStat
	DB.Find(&existing)
	have := make(map[string]struct{}, len(existing))
	for _, s := range existing {
		have[s.TableName] = struct{}{}
	}
	for _, spec := range gpTableSpecs {
		if _, ok := have[spec.statName]; ok {
			continue
		}
		var rowCount int64
		if err := DB.Model(&spec.model).Count(&rowCount).Error; err != nil {
			appLog.WithField("spec", spec.statName).WithField("error", err.Error()).Warn("Failed to seed missing gp_table_stats row")
			continue
		}
		setGPTableStat(spec.statName, rowCount)
	}
}

// RebuildGPTableStats recomputes row counts for all tracked tables from scratch
// and writes them to gp_table_stats. Run once after the initial deploy, or to
// repair drift caused by out-of-band database changes.
func RebuildGPTableStats() error {
	for _, spec := range gpTableSpecs {
		var count int64
		if err := DB.Model(&spec.model).Count(&count).Error; err != nil {
			return fmt.Errorf("rebuild stats: count %s: %w", spec.statName, err)
		}
		setGPTableStat(spec.statName, count)
	}
	return nil
}

// GetGPTableMetaData returns metadata (populated status + per-table-type
// IsRequired configuration) for all tracked GP rate tables. It reads from
// gp_table_stats in a single query instead of issuing 32 sequential COUNT
// queries, and joins in table_configurations so the UI can render the
// "Required" toggle and last-updated info without a second round-trip.
func GetGPTableMetaData() (map[string]interface{}, error) {
	// Fetch all stat rows in one query.
	var stats []models.GPTableStat
	DB.Find(&stats)

	statsMap := make(map[string]int64, len(stats))
	for _, s := range stats {
		statsMap[s.TableName] = s.RowCount
	}

	// Fetch all table_configuration rows in one query so we can decorate
	// each metadata entry with IsRequired / UpdatedBy / UpdatedAt.
	var configs []models.TableConfiguration
	DB.Find(&configs)
	configMap := make(map[string]models.TableConfiguration, len(configs))
	for _, c := range configs {
		configMap[c.TableType] = c
	}

	metadata := make([]models.TableMetaData, 0, len(gpTableSpecs))
	for _, spec := range gpTableSpecs {
		count, exists := statsMap[spec.statName]
		populated := exists && count > 0

		entry := models.TableMetaData{
			TableType:  spec.displayName,
			Category:   spec.category,
			Populated:  populated,
			TableKey:   spec.statName,
			DeleteKey:  spec.deleteKey,
			IsRequired: true, // safe default if no configuration row yet
		}
		if cfg, ok := configMap[spec.statName]; ok {
			entry.IsRequired = cfg.IsRequired
			entry.UpdatedBy = cfg.UpdatedBy
			ua := cfg.UpdatedAt
			entry.UpdatedAt = &ua
		}
		metadata = append(metadata, entry)
	}

	return map[string]interface{}{"associated_tables": metadata}, nil
}

func GetGPTableMetaDataLEGACY() (map[string]interface{}, error) {
	var legacyMetadata []models.TableMetaData
	var results = make(map[string]interface{})
	tables := []string{"gpParameters", "glaRate", "ptdRate", "ciRate", "accidentalTtdRate", "ttdRate", "phiRate",
		"childMortality", "industryLoading", "funeralParameters", "gpReinsuranceStructure", "incomeLevel", "occupationClass",
		"gpEducatorBenefitStructure", "educatorRiskRates", "escalations",
		"funeralAidsRate", "funeralRate", "generalLoading", "glaAidsRate", "regionLoading", "tieredIncomeReplacement", "customTieredIncomeReplacement",
		"discountAuthority", "restriction",
		"premiumLoading", "schemeSizeLevel", "taxTable",
		"reinsuranceGlaRate", "reinsuranceCiRate", "reinsurancePtdRate", "reinsurancePhiRate"}
	//tableMap := make(map[string]interface{})
	//tableMap := map[string]models.TableMetaData{"gpParameters": {}, "glaRate": {}, "ptdRate": {}, "ciRate": {}, "accidentalTtdRate": {}, "ttdRate": {},
	//	"phiRate": {}, "childMortality": {}, "industryLoading": {}, "funeralParameters": {}, "gpReinsuranceStructure": {}, "incomeLevel": {},
	//	"occupationClass": {}, "gpEducatorBenefitStructure": {}, "educatorRiskRates": {}, "escalations": {}}
	for _, table := range tables {
		var tableType string
		category := "group_pricing"
		switch table {
		case "gpParameters":
			tableType = "Group Pricing Parameters"
			var m models.GroupPricingParameters
			populated, err := IsTableEmpty(m)
			if err != nil {
				return nil, err
			}
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "glaRate":
			tableType = "Group Life Assurance"
			var m models.GlaRate
			populated, err := IsTableEmpty(m)
			if err != nil {
				return nil, err
			}
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "ptdRate":
			tableType = "Permanent and Total Disability"
			var m models.PtdRate
			populated, err := IsTableEmpty(m)
			if err != nil {
				return nil, err
			}
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "ciRate":
			tableType = "Critical Illness"
			var m models.CiRate
			populated, err := IsTableEmpty(m)
			if err != nil {
				return nil, err
			}
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "accidentalTtdRate":
			tableType = "Accidental Temporary Total Disability"
			var m models.AccidentalTtdRate
			populated, err := IsTableEmpty(m)
			if err != nil {
				return nil, err
			}
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "ttdRate":
			tableType = "Temporary Total Disability"
			var m models.TtdRate
			populated, err := IsTableEmpty(m)
			if err != nil {
				return nil, err
			}
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "phiRate":
			tableType = "Permanent Health Insurance"
			var m models.PhiRate
			populated, err := IsTableEmpty(m)
			if err != nil {
				return nil, err
			}
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "childMortality":
			tableType = "Child Mortality"
			var m models.ChildMortality
			populated, err := IsTableEmpty(m)
			if err != nil {
				return nil, err
			}
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "industryLoading":
			tableType = "Industry Loading"
			var m models.IndustryLoading
			populated, err := IsTableEmpty(m)
			if err != nil {
				return nil, err
			}
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "funeralParameters":
			tableType = "Funeral Parameters"
			var m models.FuneralParameters
			populated, err := IsTableEmpty(m)
			if err != nil {
				return nil, err
			}
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "gpReinsuranceStructure":
			tableType = "Group Pricing Reinsurance Structure"
			var m models.GroupPricingReinsuranceStructure
			populated, err := IsTableEmpty(m)
			if err != nil {
				return nil, err
			}
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "incomeLevel":
			tableType = "Income Level"
			var m models.IncomeLevel
			populated, err := IsTableEmpty(m)
			if err != nil {
				return nil, err
			}
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "occupationClass":
			tableType = "Occupation Class"
			var m models.OccupationClass
			populated, err := IsTableEmpty(m)
			if err != nil {
				return nil, err
			}
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "gpEducatorBenefitStructure":
			tableType = "Group Pricing Educator Structure"
			var m models.EducatorBenefitStructure
			populated, err := IsTableEmpty(m)
			if err != nil {
				return nil, err
			}
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "educatorRiskRates":
			tableType = "Educator Risk Rates"
			var m models.EducatorRate
			populated, err := IsTableEmpty(m)
			if err != nil {
				return nil, err
			}
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "escalations":
			tableType = "Escalations"
			var m models.Escalations
			populated, err := IsTableEmpty(m)
			if err != nil {
				return nil, err
			}
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "funeralAidsRate":
			tableType = "Funeral Aids Rates"
			var m models.FuneralAidsRate
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "funeralRate":
			tableType = "Funeral Rates"
			var m models.FuneralRate
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "generalLoading":
			tableType = "General Loadings"
			var m models.GeneralLoading
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "glaAidsRate":
			tableType = "GLA Aids Rates"
			var m models.GlaAidsRate
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "regionLoading":
			tableType = "Region Loading"
			var m models.RegionLoading
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "tieredIncomeReplacement":
			tableType = "Tiered Income Replacement"
			var m models.TieredIncomeReplacement
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "customTieredIncomeReplacement":
			tableType = "Custom Tiered Income Replacement"
			var m models.CustomTieredIncomeReplacement
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "discountAuthority":
			tableType = "Discount Authority"
			var m models.DiscountAuthority
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "restriction":
			tableType = "Restrictions"
			var m models.Restriction
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "premiumLoading":
			tableType = "Premium Loadings"
			var m models.PremiumLoading
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "schemeSizeLevel":
			tableType = "Scheme Size Levels"
			var m models.SchemeSizeLevel
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "taxTable":
			tableType = "Tax Table"
			var m models.TaxTable
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "reinsuranceGlaRate":
			tableType = "Reinsurance GLA Rate"
			category = "reinsurance"
			var m models.ReinsuranceGlaRate
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "reinsuranceCiRate":
			tableType = "Reinsurance CI Rate"
			category = "reinsurance"
			var m models.ReinsuranceCiRate
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "reinsurancePtdRate":
			tableType = "Reinsurance PTD Rate"
			category = "reinsurance"
			var m models.ReinsurancePtdRate
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "reinsurancePhiRate":
			tableType = "Reinsurance PHI Rate"
			category = "reinsurance"
			var m models.ReinsurancePhiRate
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "reinsuranceFuneralAidsRate":
			tableType = "Reinsurance Funeral Aids Rate"
			category = "reinsurance"
			var m models.ReinsuranceFuneralAidsRate
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "reinsuranceFuneralRate":
			tableType = "Reinsurance Funeral Rate"
			category = "reinsurance"
			var m models.ReinsuranceFuneralRate
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "reinsuranceGlaAidsRate":
			tableType = "Reinsurance GLA Aids Rate"
			category = "reinsurance"
			var m models.ReinsuranceGlaAidsRate
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "reinsuranceGeneralLoading":
			tableType = "Reinsurance General Loading"
			category = "reinsurance"
			var m models.ReinsuranceGeneralLoading
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "reinsuranceIndustryLoading":
			tableType = "Reinsurance Industry Loading"
			category = "reinsurance"
			var m models.ReinsuranceIndustryLoading
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		case "reinsuranceRegionLoading":
			tableType = "Reinsurance Region Loading"
			category = "reinsurance"
			var m models.ReinsuranceRegionLoading
			populated, _ := IsTableEmpty(m)
			legacyMetadata = append(legacyMetadata, models.TableMetaData{TableType: tableType, Category: category, Data: nil, Populated: populated})
		}
	}
	// Test additions
	//var gpParameters = models.TableMetaData{TableType: "Group Pricing Parameters", Data: nil, Populated: true}
	//var glaRate = models.TableMetaData{TableType: "Group Life Assurance", Data: nil, Populated: true}
	//var ptdRate = models.TableMetaData{TableType: "Permanent and Total Disability", Data: nil, Populated: true}
	//var ciRate = models.TableMetaData{TableType: "Critical Illness", Data: nil, Populated: true}
	//var accidentalTtdRate = models.TableMetaData{TableType: "Accidental Temporary Total Disability", Data: nil, Populated: true}
	//var ttdRate = models.TableMetaData{TableType: "Temporary Total Disability", Data: nil, Populated: true}
	//var phiRate = models.TableMetaData{TableType: "Permanent Health Insurance", Data: nil, Populated: true}
	//var childMortality = models.TableMetaData{TableType: "Child Mortality", Data: nil, Populated: true}
	//var industryLoading = models.TableMetaData{TableType: "Industry Loading", Data: nil, Populated: true}
	//var funeralParameters = models.TableMetaData{TableType: "Funeral Parameters", Data: nil, Populated: true}
	//var gpReinsuranceStructure = models.TableMetaData{TableType: "Group Pricing Reinsurance Structure", Data: nil, Populated: true}
	//var incomeLevel = models.TableMetaData{TableType: "Income Level", Data: nil, Populated: true}
	//var occupationClass = models.TableMetaData{TableType: "Occupation Class", Data: nil, Populated: true}
	//var gpEducatorBenefitStructure = models.TableMetaData{TableType: "Group Pricing Educator Structure", Data: nil, Populated: true}
	//var educatorRiskRates = models.TableMetaData{TableType: "Educator Risk Rates", Data: nil, Populated: true}
	//var escalations = models.TableMetaData{TableType: "Escalations", Data: nil, Populated: true}
	//metadata = append(metadata, gpParameters, glaRate, ptdRate, ciRate, accidentalTtdRate, ttdRate, phiRate,
	//	childMortality, industryLoading, funeralParameters, gpReinsuranceStructure, incomeLevel, occupationClass, gpEducatorBenefitStructure, educatorRiskRates, escalations)

	results["associated_tables"] = legacyMetadata
	return results, nil
}

func SaveQuoteTables(v *multipart.FileHeader, tableType string, quoteId int, user models.AppUser) (error, int) {
	var count int64
	var delimiter rune
	delimiterFile, err := v.Open()
	if err != nil {
		return fmt.Errorf("failed to open file for delimiter detection: %v", err), int(count)
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)
	if err != nil {
		return fmt.Errorf("failed to detect CSV delimiter: %v", err), int(count)
	}

	file, err := v.Open()
	if err != nil {
		return fmt.Errorf("failed to open file for reading: %v", err), int(count)
	}
	defer file.Close()

	var reader io.Reader = file
	magic := make([]byte, 2)
	if n, err := file.Read(magic); err == nil && n == 2 && magic[0] == 0x1f && magic[1] == 0x8b {
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek: %v", err), int(count)
		}
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %v", err), int(count)
		}
		defer gzReader.Close()
		reader = gzReader
	} else {
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek: %v", err), int(count)
		}
	}

	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true
	csvReader.Comma = delimiter
	dec, _ := csvutil.NewDecoder(csvReader)
	headers := dec.Header()

	var quote models.GroupPricingQuote
	DB.First(&quote, quoteId)
	switch tableType {
	case "Member Data":
		// Member Data uses a local shadow struct memberDataCSV for decoding
		type memberDataCSV struct {
			Year                         int            `csv:"year"`
			SchemeName                   string         `csv:"scheme_name"`
			MemberName                   string         `csv:"member_name"`
			MemberIdNumber               string         `csv:"member_id_number"`
			MemberIdType                 string         `csv:"member_id_type"`
			SchemeCategory               string         `csv:"scheme_category"`
			Gender                       string         `csv:"gender"`
			Email                        string         `csv:"email"`
			EmployeeNumber               string         `csv:"employee_number"`
			DateOfBirth                  models.CsvTime `csv:"date_of_birth"`
			AnnualSalary                 float64        `csv:"annual_salary"`
			ContributionWaiverProportion float64        `csv:"contribution_waiver_proportion"`
			EntryDate                    models.CsvTime `csv:"entry_date"`
			ExitDate                     models.CsvTime `csv:"exit_date"`
			EffectiveExitDate            models.CsvTime `csv:"effective_exit_date"`

			// Flattened MemberBenefits columns (option 1)
			BenefitsGlaMultiple  float64 `csv:"benefits_gla_multiple"`
			BenefitsSglaMultiple float64 `csv:"benefits_sgla_multiple"`
			BenefitsPtdMultiple  float64 `csv:"benefits_ptd_multiple"`
			BenefitsCiMultiple   float64 `csv:"benefits_ci_multiple"`
			BenefitsTtdMultiple  float64 `csv:"benefits_ttd_multiple"`
			BenefitsPhiMultiple  float64 `csv:"benefits_phi_multiple"`
		}

		if validationErr := utils.ValidateCSVHeaders(headers, memberDataCSV{}); validationErr != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, validationErr), 0
		}

		DB.Where("quote_id = ?", quoteId).Delete(&models.GPricingMemberData{})
		var membersData []models.GPricingMemberData
		var validationErrors []string
		for i := 1; ; i++ {
			var row memberDataCSV
			if err := dec.Decode(&row); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Member Data at row %d: %v", i, err), 0
			}

			// Sanity check: if a scheme_name is provided in the CSV row, ensure it
			// exists in the quote's SelectedSchemeCategories. If not, skip this row.
			csvSchemeName := strings.TrimSpace(row.SchemeCategory)
			if csvSchemeName != "" {
				valid := false
				for _, cat := range quote.SelectedSchemeCategories {
					if strings.EqualFold(strings.TrimSpace(cat), csvSchemeName) {
						valid = true
						break
					}
				}
				if !valid {
					// Skip adding this member data row as it doesn't match allowed categories
					continue
				}
			}
			// Validate RSA ID
			idType := strings.ToUpper(strings.TrimSpace(row.MemberIdType))
			idNumber := strings.TrimSpace(row.MemberIdNumber)
			if (idType == "RSA_ID" || idType == "ID") && idNumber != "" {
				valid, checkErr := utils.ValidateRSAID(idNumber)
				if checkErr != nil {
					return fmt.Errorf("ID validation service error: %v", checkErr), 0
				}
				if !valid {
					validationErrors = append(validationErrors, fmt.Sprintf("invalid RSA ID '%s' at row %d", idNumber, i))
					continue
				}
			}

			pp := models.GPricingMemberData{
				Year:                         row.Year,
				MemberName:                   row.MemberName,
				MemberIdNumber:               row.MemberIdNumber,
				MemberIdType:                 row.MemberIdType,
				SchemeCategory:               strings.TrimSpace(row.SchemeCategory),
				Gender:                       row.Gender,
				DateOfBirth:                  time.Time(row.DateOfBirth),
				Email:                        strings.TrimSpace(row.Email),
				EmployeeNumber:               strings.TrimSpace(row.EmployeeNumber),
				AnnualSalary:                 row.AnnualSalary,
				ContributionWaiverProportion: row.ContributionWaiverProportion,
				EntryDate:                    time.Time(row.EntryDate),
				ExitDate: func() *time.Time {
					if t := time.Time(row.ExitDate); !t.IsZero() {
						return &t
					}
					return nil
				}(),
				EffectiveExitDate: func() *time.Time {
					if t := time.Time(row.EffectiveExitDate); !t.IsZero() {
						return &t
					}
					return nil
				}(),
				CreatedBy: user.UserName,
				QuoteId:   quoteId,
				SchemeId:  quote.SchemeID,
				// Keep existing behavior of setting SchemeName from quote but we could
				// also persist the CSV-provided scheme_name if needed in future.
				SchemeName: quote.SchemeName,
				Benefits: models.MemberBenefits{
					GlaMultiple:  row.BenefitsGlaMultiple,
					SglaMultiple: row.BenefitsSglaMultiple,
					PtdMultiple:  row.BenefitsPtdMultiple,
					CiMultiple:   row.BenefitsCiMultiple,
					TtdMultiple:  row.BenefitsTtdMultiple,
					PhiMultiple:  row.BenefitsPhiMultiple,
				},
			}
			membersData = append(membersData, pp)
			//err = DB.Create(&pp).Error
			//if err != nil {
			//	appLog.Error("Save Quote Tables error: ", err.Error())
			//}
		}

		if len(validationErrors) > 0 {
			var errMsg string
			if len(validationErrors) > 3 {
				errMsg = strings.Join(validationErrors[:3], "; ") + fmt.Sprintf(" +%d more", len(validationErrors)-3)
			} else {
				errMsg = strings.Join(validationErrors, "; ")
			}
			return fmt.Errorf("validation failed: %s", errMsg), 0
		}

		//err := utils.ValidMemberIdColValues(membersData)
		//if err != nil {
		//	return fmt.Errorf("Member ID validation failed: %v", err), 0
		//}

		err = DB.CreateInBatches(&membersData, 100).Error
		if err != nil {
			appLog.Error("Save Quote Tables error: ", err.Error())
			return fmt.Errorf("failed to save Member Data: %v", err), 0
		}
		err = DB.Model(&models.GPricingMemberData{}).Where("quote_id = ?", quoteId).Count(&count).Error
		if err == nil {
			// we remove all indicative data
			quote.MemberIndicativeData = false
			quote.MemberDataCount = 0
			quote.MemberAverageAge = 0
			quote.MemberAverageIncome = 0
			quote.MemberMaleFemaleDistribution = 0
			DB.Save(&quote)
		}
	case "Claims Experience":
		if err := utils.ValidateCSVHeaders(headers, models.GroupPricingClaimsExperience{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err), 0
		}
		DB.Where("quote_id = ?", quoteId).Delete(&models.GroupPricingClaimsExperience{})
		for i := 1; ; i++ {
			var pp models.GroupPricingClaimsExperience
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Claims Experience at row %d: %v", i, err), 0
			}
			pp.QuoteId = quoteId
			pp.CreatedBy = user.UserName
			err = DB.Create(&pp).Error
			if err != nil {
				appLog.Error("Save Quote Tables error: ", err.Error())
				return fmt.Errorf("failed to save Claims Experience data at row %d: %v", i, err), 0
			}
		}
		DB.Model(&models.GroupPricingClaimsExperience{}).Where("quote_id = ?", quoteId).Count(&count)
	}
	// Drop any cached per-quote lookups so the next calculation reads the
	// freshly-uploaded data rather than the pre-upload snapshot.
	if GroupPricingCache != nil {
		GroupPricingCache.Clear()
	}
	return nil, int(count)
}

func DeleteQuoteTableData(tableType string, quoteId int) error {
	switch tableType {
	case "Member Data":
		DB.Where("quote_id = ?", quoteId).Delete(&models.GPricingMemberData{})
	case "Claims Experience":
		DB.Where("quote_id = ?", quoteId).Delete(&models.GroupPricingClaimsExperience{})
	case "Member Rating Results":
		DB.Where("quote_id = ?", quoteId).Delete(&models.MemberRatingResult{})
		DB.Where("quote_id = ?", quoteId).Delete(&models.MemberRatingResultSummary{})
	case "Member Premium Schedules":
		DB.Where("quote_id = ?", quoteId).Delete(&models.MemberPremiumSchedule{})
	case "Bordereauxes":
		DB.Where("quote_id = ?", quoteId).Delete(&models.Bordereaux{})
	case "Group Scheme Exposures":
		DB.Where("quote_id = ?", quoteId).Delete(&models.GroupSchemeExposure{})
	}
	// Drop any cached per-quote lookups so the next calculation cannot
	// reuse rows that have just been deleted.
	if GroupPricingCache != nil {
		GroupPricingCache.Clear()
	}
	return nil
}

func SaveGPTables(v *multipart.FileHeader, tableType string, riskRateCode string, user models.AppUser, schemeName string) error {
	var delimiter rune
	delimiterFile, err := v.Open()
	if err != nil {
		return fmt.Errorf("failed to open file for delimiter detection: %v", err)
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)
	if err != nil {
		return fmt.Errorf("failed to detect CSV delimiter: %v", err)
	}

	//var HeaderMismatch = errors.New("Columns of uploaded table do not match.")
	file, err := v.Open()
	if err != nil {
		return fmt.Errorf("failed to open file for reading: %v", err)
	}
	defer file.Close()

	var reader io.Reader = file
	magic := make([]byte, 2)
	if n, err := file.Read(magic); err == nil && n == 2 && magic[0] == 0x1f && magic[1] == 0x8b {
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek: %v", err)
		}
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %v", err)
		}
		defer gzReader.Close()
		reader = gzReader
	} else {
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek: %v", err)
		}
	}

	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true
	csvReader.Comma = delimiter
	dec, _ := csvutil.NewDecoder(csvReader)
	headers := dec.Header()

	switch tableType {
	case "Group Life Assurance":
		if err := utils.ValidateCSVHeaders(headers, models.GlaRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.GlaRate{})
		var pps []models.GlaRate
		for i := 1; ; i++ {
			var pp models.GlaRate
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Group Life Assurance at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName

			pps = append(pps, pp)
		}

		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Group Life Assurance error: ", err.Error())
			return fmt.Errorf("failed to save Group Life Assurance data: %v", err)
		}

	case "Permanent and Total Disability":
		if err := utils.ValidateCSVHeaders(headers, models.PtdRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.PtdRate{})
		var pps []models.PtdRate
		for i := 1; ; i++ {
			var pp models.PtdRate
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Permanent and Total Disability at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}

		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Permanent and Total Disability error: ", err.Error())
			return fmt.Errorf("failed to save Permanent and Total Disability data: %v", err)
		}

	case "Critical Illness":
		if err := utils.ValidateCSVHeaders(headers, models.CiRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.CiRate{})
		var pps []models.CiRate
		for i := 1; ; i++ {
			var pp models.CiRate
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Critical Illness at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}

		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Critical Illness error: ", err.Error())
			return fmt.Errorf("failed to save Critical Illness data: %v", err)
		}

	case "Accidental Temporary Total Disability":
		if err := utils.ValidateCSVHeaders(headers, models.AccidentalTtdRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.AccidentalTtdRate{})
		var pps []models.AccidentalTtdRate
		for i := 1; ; i++ {
			var pp models.AccidentalTtdRate

			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Accidental Temporary Total Disability at row %d: %v", i, err)
			}

			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}

		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Accidental Temporary Total Disability error: ", err.Error())
			return fmt.Errorf("failed to save Accidental Temporary Total Disability data: %v", err)
		}

	case "Temporary Total Disability":
		if err := utils.ValidateCSVHeaders(headers, models.TtdRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.TtdRate{})
		var pps []models.TtdRate
		for i := 1; ; i++ {
			var pp models.TtdRate
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Temporary Total Disability at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}

		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Temporary Total Disability error: ", err.Error())
			return fmt.Errorf("failed to save Temporary Total Disability data: %v", err)
		}

	case "Permanent Health Insurance":
		if err := utils.ValidateCSVHeaders(headers, models.PhiRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.PhiRate{})
		var pps []models.PhiRate
		for i := 1; ; i++ {
			var pp models.PhiRate
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Permanent Health Insurance at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName

			pps = append(pps, pp)
		}

		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Permanent Health Insurance error: ", err.Error())
			return fmt.Errorf("failed to save Permanent Health Insurance data: %v", err)
		}

	case "Child Mortality":
		if err := utils.ValidateCSVHeaders(headers, models.ChildMortality{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.ChildMortality{})
		var pps []models.ChildMortality
		for i := 1; ; i++ {
			var pp models.ChildMortality
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Child Mortality at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}

		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save ChildMortality error: ", err.Error())
			return fmt.Errorf("failed to save Child Mortality data: %v", err)
		}

	case "Industry Loading":
		if err := utils.ValidateCSVHeaders(headers, models.IndustryLoading{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.IndustryLoading{})
		var pps []models.IndustryLoading
		for i := 1; ; i++ {
			var pp models.IndustryLoading
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Industry Loading at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}

		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Industry Loading error: ", err.Error())
			return fmt.Errorf("failed to save Industry Loading data: %v", err)
		}

	case "Funeral Parameters":
		if err := utils.ValidateCSVHeaders(headers, models.FuneralParameters{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.FuneralParameters{})
		var pps []models.FuneralParameters
		for i := 1; ; i++ {
			var pp models.FuneralParameters
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Funeral Parameters at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}

		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error(err.Error())
			return fmt.Errorf("failed to save Funeral Parameters data: %v", err)
		}

	case "Group Pricing Reinsurance Structure":
		if err := utils.ValidateCSVHeaders(headers, models.GroupPricingReinsuranceStructure{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.GroupPricingReinsuranceStructure{})
		var pps []models.GroupPricingReinsuranceStructure
		for i := 1; ; i++ {
			var pp models.GroupPricingReinsuranceStructure
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Group Pricing Reinsurance Structure at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Info("Save Group Pricing Reinsurance Structure error: ", err.Error())
			return fmt.Errorf("failed to save Group Pricing Reinsurance Structure data: %v", err)
		}

	case "Group Pricing Parameters":
		if err := utils.ValidateCSVHeaders(headers, models.GroupPricingParameters{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.GroupPricingParameters{})
		var pps []models.GroupPricingParameters
		for i := 1; ; i++ {
			var pp models.GroupPricingParameters
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Group Pricing Parameters at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Group Pricing Parameter Error: ", err.Error())
			return fmt.Errorf("failed to save Group Pricing Parameters data: %v", err)
		}

	case "Income Level":
		if err := utils.ValidateCSVHeaders(headers, models.IncomeLevel{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.IncomeLevel{})
		for i := 1; ; i++ {
			var pp models.IncomeLevel
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Income Level at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			err = DB.Create(&pp).Error
			if err != nil {
				appLog.Info("Save Income Level error: ", err.Error())
				return fmt.Errorf("failed to save Income Level data at row %d: %v", i, err)
			}
		}
	case "Occupation Class":
		if err := utils.ValidateCSVHeaders(headers, models.OccupationClass{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.OccupationClass{})
		var ocs []models.OccupationClass
		for i := 1; ; i++ {
			var pp models.OccupationClass
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Occupation Class at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			ocs = append(ocs, pp)
		}
		err = DB.CreateInBatches(&ocs, 100).Error
		if err != nil {
			appLog.Info("Save Occupation Class error: ", err.Error())
			return fmt.Errorf("failed to save Occupation Class data: %v", err)
		}
	case "Group Pricing Educator Structure":
		if err := utils.ValidateCSVHeaders(headers, models.EducatorBenefitStructure{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.EducatorBenefitStructure{})
		for i := 1; ; i++ {
			var pp models.EducatorBenefitStructure
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Group Pricing Educator Structure at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			err = DB.Create(&pp).Error
			if err != nil {
				appLog.Error(err.Error())
				return fmt.Errorf("failed to save Group Pricing Educator Structure data at row %d: %v", i, err)
			}
		}
	case "Educator Risk Rates":
		if err := utils.ValidateCSVHeaders(headers, models.EducatorRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.EducatorRate{})
		var pps []models.EducatorRate
		for i := 1; ; i++ {
			var pp models.EducatorRate
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Educator Risk Rates at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error(err.Error())
			return fmt.Errorf("failed to save Educator Risk Rates data: %v", err)
		}

	case "Escalations":
		if err := utils.ValidateCSVHeaders(headers, models.Escalations{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.Escalations{})
		var pps []models.Escalations
		for i := 1; ; i++ {
			var pp models.Escalations
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Escalations at row %d: %v", i, err)
			}
			//pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Escalations error: ", err.Error())
			return fmt.Errorf("failed to save Escalations data: %v", err)
		}

	case "Funeral Aids Rates":
		if err := utils.ValidateCSVHeaders(headers, models.FuneralAidsRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.FuneralAidsRate{})
		var pps []models.FuneralAidsRate
		for i := 1; ; i++ {
			var pp models.FuneralAidsRate
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Funeral Aids Rates at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Funeral Aids Rates error: ", err.Error())
			return fmt.Errorf("failed to save Funeral Aids Rates data: %v", err)
		}

	case "Funeral Rates":
		if err := utils.ValidateCSVHeaders(headers, models.FuneralRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.FuneralRate{})
		var pps []models.FuneralRate
		for i := 1; ; i++ {
			var pp models.FuneralRate
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Funeral Rates at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Funeral Rates error: ", err.Error())
			return fmt.Errorf("failed to save Funeral Rates data: %v", err)
		}

	case "General Loadings":
		if err := utils.ValidateCSVHeaders(headers, models.GeneralLoading{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.GeneralLoading{})
		var pps []models.GeneralLoading
		for i := 1; ; i++ {
			var pp models.GeneralLoading
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding General Loadings at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save General Loadings error: ", err.Error())
			return fmt.Errorf("failed to save General Loadings data: %v", err)
		}

	case "GLA Aids Rates":
		if err := utils.ValidateCSVHeaders(headers, models.GlaAidsRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.GlaAidsRate{})
		var pps []models.GlaAidsRate
		for i := 1; ; i++ {
			var pp models.GlaAidsRate
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding GLA Aids Rates at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save GLA Aids Rates error: ", err.Error())
			return fmt.Errorf("failed to save GLA Aids Rates data: %v", err)
		}

	case "Region Loading":
		if err := utils.ValidateCSVHeaders(headers, models.RegionLoading{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.RegionLoading{})
		var pps []models.RegionLoading
		for i := 1; ; i++ {
			var pp models.RegionLoading
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Region Loading at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Region Loading error: ", err.Error())
			return fmt.Errorf("failed to save Region Loading data: %v", err)
		}

	case "Tiered Income Replacement":
		if err := utils.ValidateCSVHeaders(headers, models.TieredIncomeReplacement{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.TieredIncomeReplacement{})
		var pps []models.TieredIncomeReplacement
		for i := 1; ; i++ {
			var pp models.TieredIncomeReplacement
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Tiered Income Replacement at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Tiered Income Replacement error: ", err.Error())
			return fmt.Errorf("failed to save Tiered Income Replacement data: %v", err)
		}

	case "Custom Tiered Income Replacement":
		if schemeName == "" {
			return fmt.Errorf("scheme_name is required for Custom Tiered Income Replacement uploads")
		}
		if err := utils.ValidateCSVHeaders(headers, models.CustomTieredIncomeReplacement{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("scheme_name = ? AND risk_rate_code = ?", schemeName, riskRateCode).Delete(&models.CustomTieredIncomeReplacement{})
		var pps []models.CustomTieredIncomeReplacement
		for i := 1; ; i++ {
			var pp models.CustomTieredIncomeReplacement
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Custom Tiered Income Replacement at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.SchemeName = schemeName
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Custom Tiered Income Replacement error: ", err.Error())
			return fmt.Errorf("failed to save Custom Tiered Income Replacement data: %v", err)
		}
		// Notify the requesting user that the custom table has been uploaded
		go NotifyCustomTieredTableUploaded(schemeName, riskRateCode, user)

	case "Discount Authority":
		if err := utils.ValidateCSVHeaders(headers, models.DiscountAuthority{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.DiscountAuthority{})
		var pps []models.DiscountAuthority
		for i := 1; ; i++ {
			var pp models.DiscountAuthority
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Discount Authority at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Discount Authority error: ", err.Error())
			return fmt.Errorf("failed to save Discount Authority data: %v", err)
		}

	case "Restrictions":
		if err := utils.ValidateCSVHeaders(headers, models.Restriction{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.Restriction{})
		var pps []models.Restriction
		for i := 1; ; i++ {
			var pp models.Restriction
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Restrictions at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Restrictions error: ", err.Error())
			return fmt.Errorf("failed to save Restrictions data: %v", err)
		}

	case "Reinsurance Cover Restrictions":
		if err := utils.ValidateCSVHeaders(headers, models.ReinsuranceCoverRestriction{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.ReinsuranceCoverRestriction{})
		var pps []models.ReinsuranceCoverRestriction
		for i := 1; ; i++ {
			var pp models.ReinsuranceCoverRestriction
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Reinsurance Cover Restrictions at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Reinsurance Cover Restrictions error: ", err.Error())
			return fmt.Errorf("failed to save Reinsurance Cover Restrictions data: %v", err)
		}

	case "Reinsurance GLA Rate":
		if err := utils.ValidateCSVHeaders(headers, models.ReinsuranceGlaRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.ReinsuranceGlaRate{})
		var pps []models.ReinsuranceGlaRate
		for i := 1; ; i++ {
			var pp models.ReinsuranceGlaRate
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Reinsurance GLA Rate at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Reinsurance GLA Rate error: ", err.Error())
			return fmt.Errorf("failed to save Reinsurance GLA Rate data: %v", err)
		}

	case "Reinsurance CI Rate":
		if err := utils.ValidateCSVHeaders(headers, models.ReinsuranceCiRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.ReinsuranceCiRate{})
		var pps []models.ReinsuranceCiRate
		for i := 1; ; i++ {
			var pp models.ReinsuranceCiRate
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Reinsurance CI Rate at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Reinsurance CI Rate error: ", err.Error())
			return fmt.Errorf("failed to save Reinsurance CI Rate data: %v", err)
		}

	case "Reinsurance PTD Rate":
		if err := utils.ValidateCSVHeaders(headers, models.ReinsurancePtdRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.ReinsurancePtdRate{})
		var pps []models.ReinsurancePtdRate
		for i := 1; ; i++ {
			var pp models.ReinsurancePtdRate
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Reinsurance PTD Rate at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Reinsurance PTD Rate error: ", err.Error())
			return fmt.Errorf("failed to save Reinsurance PTD Rate data: %v", err)
		}

	case "Reinsurance PHI Rate":
		if err := utils.ValidateCSVHeaders(headers, models.ReinsurancePhiRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.ReinsurancePhiRate{})
		var pps []models.ReinsurancePhiRate
		for i := 1; ; i++ {
			var pp models.ReinsurancePhiRate
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Reinsurance PHI Rate at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Reinsurance PHI Rate error: ", err.Error())
			return fmt.Errorf("failed to save Reinsurance PHI Rate data: %v", err)
		}

	case "Reinsurance Funeral Aids Rate":
		if err := utils.ValidateCSVHeaders(headers, models.ReinsuranceFuneralAidsRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.ReinsuranceFuneralAidsRate{})
		var pps []models.ReinsuranceFuneralAidsRate
		for i := 1; ; i++ {
			var pp models.ReinsuranceFuneralAidsRate
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Reinsurance Funeral Aids Rate at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Reinsurance Funeral Aids Rate error: ", err.Error())
			return fmt.Errorf("failed to save Reinsurance Funeral Aids Rate data: %v", err)
		}

	case "Reinsurance Funeral Rate":
		if err := utils.ValidateCSVHeaders(headers, models.ReinsuranceFuneralRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.ReinsuranceFuneralRate{})
		var pps []models.ReinsuranceFuneralRate
		for i := 1; ; i++ {
			var pp models.ReinsuranceFuneralRate
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Reinsurance Funeral Rate at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Reinsurance Funeral Rate error: ", err.Error())
			return fmt.Errorf("failed to save Reinsurance Funeral Rate data: %v", err)
		}

	case "Reinsurance GLA Aids Rate":
		if err := utils.ValidateCSVHeaders(headers, models.ReinsuranceGlaAidsRate{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.ReinsuranceGlaAidsRate{})
		var pps []models.ReinsuranceGlaAidsRate
		for i := 1; ; i++ {
			var pp models.ReinsuranceGlaAidsRate
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Reinsurance GLA Aids Rate at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Reinsurance GLA Aids Rate error: ", err.Error())
			return fmt.Errorf("failed to save Reinsurance GLA Aids Rate data: %v", err)
		}

	case "Reinsurance General Loading":
		if err := utils.ValidateCSVHeaders(headers, models.ReinsuranceGeneralLoading{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.ReinsuranceGeneralLoading{})
		var pps []models.ReinsuranceGeneralLoading
		for i := 1; ; i++ {
			var pp models.ReinsuranceGeneralLoading
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Reinsurance General Loading at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Reinsurance General Loading error: ", err.Error())
			return fmt.Errorf("failed to save Reinsurance General Loading data: %v", err)
		}

	case "Reinsurance Industry Loading":
		if err := utils.ValidateCSVHeaders(headers, models.ReinsuranceIndustryLoading{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.ReinsuranceIndustryLoading{})
		var pps []models.ReinsuranceIndustryLoading
		for i := 1; ; i++ {
			var pp models.ReinsuranceIndustryLoading
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Reinsurance Industry Loading at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Reinsurance Industry Loading error: ", err.Error())
			return fmt.Errorf("failed to save Reinsurance Industry Loading data: %v", err)
		}

	case "Reinsurance Region Loading":
		if err := utils.ValidateCSVHeaders(headers, models.ReinsuranceRegionLoading{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.ReinsuranceRegionLoading{})
		var pps []models.ReinsuranceRegionLoading
		for i := 1; ; i++ {
			var pp models.ReinsuranceRegionLoading
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Reinsurance Region Loading at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Reinsurance Region Loading error: ", err.Error())
			return fmt.Errorf("failed to save Reinsurance Region Loading data: %v", err)
		}

	case "Premium Loadings":
		if err := utils.ValidateCSVHeaders(headers, models.PremiumLoading{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.PremiumLoading{})
		var pps []models.PremiumLoading
		for i := 1; ; i++ {
			var pp models.PremiumLoading
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Premium Loadings at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Premium Loadings error: ", err.Error())
			return fmt.Errorf("failed to save Premium Loadings data: %v", err)
		}

	case "Scheme Size Levels":
		if err := utils.ValidateCSVHeaders(headers, models.SchemeSizeLevel{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.SchemeSizeLevel{})
		var pps []models.SchemeSizeLevel
		for i := 1; ; i++ {
			var pp models.SchemeSizeLevel
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Scheme Size Levels at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Scheme Size Levels error: ", err.Error())
			return fmt.Errorf("failed to save Scheme Size Levels data: %v", err)
		}

	case "Tax Table":
		if err := utils.ValidateCSVHeaders(headers, models.TaxTable{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.TaxTable{})
		var pps []models.TaxTable
		for i := 1; ; i++ {
			var pp models.TaxTable
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Tax Table at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Tax Table error: ", err.Error())
			return fmt.Errorf("failed to save Tax Table data: %v", err)
		}

	case "Age Bands":
		// Age bands are keyed by `type` (e.g. GLA vs funeral sets). A fresh
		// upload only replaces rows whose `type` appears in the incoming CSV
		// — other types are left untouched so uploading one band set does
		// not clobber another.
		if err := utils.ValidateCSVHeaders(headers, models.GroupPricingAgeBands{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		var pps []models.GroupPricingAgeBands
		for i := 1; ; i++ {
			var pp models.GroupPricingAgeBands
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Age Bands at row %d: %v", i, err)
			}
			// Stamp the uploading user on each row. CreationDate is filled by
			// GORM via the autoCreateTime tag.
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		// Gather the distinct types present in the upload. An empty string is
		// its own bucket so legacy untyped rows still replace legacy untyped
		// rows on re-upload.
		typeSet := make(map[string]struct{}, len(pps))
		for _, p := range pps {
			typeSet[p.Type] = struct{}{}
		}
		types := make([]string, 0, len(typeSet))
		for t := range typeSet {
			types = append(types, t)
		}
		if len(types) > 0 {
			if err := DB.Where("type IN ?", types).Delete(&models.GroupPricingAgeBands{}).Error; err != nil {
				appLog.Error("Delete Age Bands by type error: ", err.Error())
				return fmt.Errorf("failed to clear existing Age Bands for types %v: %v", types, err)
			}
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Age Bands error: ", err.Error())
			return fmt.Errorf("failed to save Age Bands data: %v", err)
		}

	case "Medical Waiver":
		if err := utils.ValidateCSVHeaders(headers, models.MedicalWaiver{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		DB.Where("risk_rate_code = ?", riskRateCode).Delete(&models.MedicalWaiver{})
		var pps []models.MedicalWaiver
		for i := 1; ; i++ {
			var pp models.MedicalWaiver
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Medical Waiver at row %d: %v", i, err)
			}
			pp.RiskRateCode = riskRateCode
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
		}
		err = DB.CreateInBatches(&pps, 100).Error
		if err != nil {
			appLog.Error("Save Medical Waiver error: ", err.Error())
			return fmt.Errorf("failed to save Medical Waiver data: %v", err)
		}

	case "Commission Structure":
		if err := utils.ValidateCSVHeaders(headers, models.CommissionStructure{}); err != nil {
			return fmt.Errorf("%s validation failed: %v", tableType, err)
		}
		var pps []models.CommissionStructure
		type chKey struct {
			Channel string
			Holder  string
		}
		affectedGroups := map[chKey]struct{}{}
		for i := 1; ; i++ {
			var pp models.CommissionStructure
			if err := dec.Decode(&pp); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("error decoding Commission Structure at row %d: %v", i, err)
			}
			pp.Channel = strings.ToLower(strings.TrimSpace(pp.Channel))
			pp.HolderName = strings.TrimSpace(pp.HolderName)
			if pp.Channel == "" {
				return fmt.Errorf("error decoding Commission Structure at row %d: channel is required", i)
			}
			if pp.Channel == "direct" {
				return fmt.Errorf("error decoding Commission Structure at row %d: direct channel cannot have bands", i)
			}
			pp.CreatedBy = user.UserName
			pps = append(pps, pp)
			affectedGroups[chKey{Channel: pp.Channel, Holder: pp.HolderName}] = struct{}{}
		}
		if err := DB.Transaction(func(tx *gorm.DB) error {
			// Replace rows per (channel, holder_name) group that appears
			// in the file — other groups are left untouched.
			for k := range affectedGroups {
				if err := tx.Where("channel = ? AND holder_name = ?", k.Channel, k.Holder).
					Delete(&models.CommissionStructure{}).Error; err != nil {
					return err
				}
			}
			if err := tx.CreateInBatches(&pps, 100).Error; err != nil {
				return err
			}
			for k := range affectedGroups {
				if err := validateCommissionBandsForChannelHolderTx(tx, k.Channel, k.Holder); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			appLog.Error("Save Commission Structure error: ", err.Error())
			return fmt.Errorf("failed to save Commission Structure data: %v", err)
		}

	}
	// Update the stats table so GetGPTableMetaData reflects the new upload.
	go refreshGPTableStatByDisplayName(tableType)
	// Drop any per-row lookups (region/industry/general/premium loadings,
	// etc.) cached by previous calculation runs — without this an in-flight
	// or immediately-following calc could reuse the pre-upload values.
	if GroupPricingCache != nil {
		GroupPricingCache.Clear()
	}
	return nil
}

func DeleteGPTableData(tableType, riskCode string) error {
	switch tableType {
	case "grouplifeassurance":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.GlaRate{})
	case "permanenthealthinsurance":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.PhiRate{})
	case "childmortality":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.ChildMortality{})
	case "industryloading":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.IndustryLoading{})
	case "grouppricingreinsurancestructure":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.GroupPricingReinsuranceStructure{})
	case "accidentaltemporarytotaldisability":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.AccidentalTtdRate{})
	case "permanentandtotaldisability":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.PtdRate{})
	case "criticalillness":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.CiRate{})
	case "temporarytotaldisability":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.TtdRate{})
	case "funeralparameters":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.FuneralParameters{})
	case "grouppricingparameters":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.GroupPricingParameters{})
	case "incomelevel":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.IncomeLevel{})
	case "occupationclass":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.OccupationClass{})
	case "educatorriskrates":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.EducatorRate{})
	case "grouppricingeducatorstructure":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.EducatorBenefitStructure{})
	case "escalations":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.Escalations{})
	case "funeralaidsrates":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.FuneralAidsRate{})
	case "funeralrates":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.FuneralRate{})
	case "generalloadings":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.GeneralLoading{})
	case "glaaidsrates":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.GlaAidsRate{})
	case "regionloading":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.RegionLoading{})
	case "tieredincomereplacement":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.TieredIncomeReplacement{})
	case "customtieredincomereplacement":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.CustomTieredIncomeReplacement{})
	case "discountauthority":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.DiscountAuthority{})
	case "restrictions":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.Restriction{})
	case "reinsurancecoverrestrictions":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.ReinsuranceCoverRestriction{})
	case "reinsuranceglarate":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.ReinsuranceGlaRate{})
	case "reinsurancecirate":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.ReinsuranceCiRate{})
	case "reinsuranceptdrate":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.ReinsurancePtdRate{})
	case "reinsurancephirate":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.ReinsurancePhiRate{})
	case "reinsurancefuneralaidsrates":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.ReinsuranceFuneralAidsRate{})
	case "reinsurancefuneralrates":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.ReinsuranceFuneralRate{})
	case "reinsuranceglaaidsrates":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.ReinsuranceGlaAidsRate{})
	case "reinsurancegeneralloadings":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.ReinsuranceGeneralLoading{})
	case "reinsuranceindustryloadings":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.ReinsuranceIndustryLoading{})
	case "reinsuranceregionloadings":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.ReinsuranceRegionLoading{})
	case "premiumloadings":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.PremiumLoading{})
	case "schemesizelevels":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.SchemeSizeLevel{})
	case "taxtable":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.TaxTable{})
	case "agebands":
		// Age bands are keyed by type. The delete picker re-uses the
		// risk-code slot to carry the selected type.
		if riskCode != "" {
			DB.Where("type = ?", riskCode).Delete(&models.GroupPricingAgeBands{})
		}
	case "medicalwaiver":
		DB.Where("risk_rate_code = ?", riskCode).Delete(&models.MedicalWaiver{})
	}
	// Update the stats table so GetGPTableMetaData reflects the deletion.
	go refreshGPTableStatByDeleteKey(tableType)
	// Drop any per-row lookups cached by previous calculation runs so the
	// next calc cannot reuse rows that have just been deleted.
	if GroupPricingCache != nil {
		GroupPricingCache.Clear()
	}
	return nil
}

func GetGPTableRiskCodes(tableType string) []string {
	var riskCodes []string
	switch tableType {
	case "grouppricingparameters":
		DB.Model(&models.GroupPricingParameters{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "grouplifeassurance":
		DB.Model(&models.GlaRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "permanenthealthinsurance":
		DB.Model(&models.PhiRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "childmortality":
		DB.Model(&models.ChildMortality{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "industryloading":
		DB.Model(&models.IndustryLoading{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "grouppricingreinsurancestructure":
		DB.Model(&models.GroupPricingReinsuranceStructure{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "accidentaltemporarytotaldisability":
		DB.Model(&models.AccidentalTtdRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "permanentandtotaldisability":
		DB.Model(&models.PtdRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "criticalillness":
		DB.Model(&models.CiRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "temporarytotaldisability":
		DB.Model(&models.TtdRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "funeralparameters":
		DB.Model(&models.FuneralParameters{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "incomelevel":
		DB.Model(&models.IncomeLevel{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "occupationclass":
		DB.Model(&models.OccupationClass{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "educatorriskrates":
		DB.Model(&models.EducatorRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "grouppricingeducatorstructure":
		DB.Model(&models.EducatorBenefitStructure{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "escalations":
		DB.Model(&models.Escalations{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "funeralaidsrates":
		DB.Model(&models.FuneralAidsRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "funeralrates":
		DB.Model(&models.FuneralRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "generalloadings":
		DB.Model(&models.GeneralLoading{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "glaaidsrates":
		DB.Model(&models.GlaAidsRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "regionloading":
		DB.Model(&models.RegionLoading{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "tieredincomereplacement":
		DB.Model(&models.TieredIncomeReplacement{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "customtieredincomereplacement":
		DB.Model(&models.CustomTieredIncomeReplacement{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "discountauthority":
		DB.Model(&models.DiscountAuthority{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "restrictions":
		DB.Model(&models.Restriction{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "reinsurancecoverrestrictions":
		DB.Model(&models.ReinsuranceCoverRestriction{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "reinsuranceglarate":
		DB.Model(&models.ReinsuranceGlaRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "reinsurancecirate":
		DB.Model(&models.ReinsuranceCiRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "reinsuranceptdrate":
		DB.Model(&models.ReinsurancePtdRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "reinsurancephirate":
		DB.Model(&models.ReinsurancePhiRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "reinsurancefuneralaidsrates":
		DB.Model(&models.ReinsuranceFuneralAidsRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "reinsurancefuneralrates":
		DB.Model(&models.ReinsuranceFuneralRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "reinsuranceglaaidsrates":
		DB.Model(&models.ReinsuranceGlaAidsRate{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "reinsurancegeneralloadings":
		DB.Model(&models.ReinsuranceGeneralLoading{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "reinsuranceindustryloadings":
		DB.Model(&models.ReinsuranceIndustryLoading{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "reinsuranceregionloadings":
		DB.Model(&models.ReinsuranceRegionLoading{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "premiumloadings":
		DB.Model(&models.PremiumLoading{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "schemesizelevels":
		DB.Model(&models.SchemeSizeLevel{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "taxtable":
		DB.Model(&models.TaxTable{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	case "agebands":
		// Age bands are keyed by type, not risk_rate_code. The frontend
		// re-uses this endpoint to populate the delete picker, so we return
		// distinct types under the same key.
		DB.Model(&models.GroupPricingAgeBands{}).Select("DISTINCT type").Order("type asc").Find(&riskCodes)
	case "medicalwaiver":
		DB.Model(&models.MedicalWaiver{}).Select("DISTINCT risk_rate_code").Order("risk_rate_code desc").Find(&riskCodes)
	}
	return riskCodes
}

func GetGPTableYears(tableType string) []int {
	var years []int
	switch tableType {
	case "grouplifeassurance":
		DB.Model(&models.GlaRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "permanenthealthinsurance":
		DB.Model(&models.PhiRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "childmortality":
		DB.Model(&models.ChildMortality{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "industryloading":
		DB.Model(&models.IndustryLoading{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "grouppricingreinsurancestructure":
		DB.Model(&models.GroupPricingReinsuranceStructure{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "accidentaltemporarytotaldisability":
		DB.Model(&models.AccidentalTtdRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "permanentandtotaldisability":
		DB.Model(&models.PtdRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "criticalillness":
		DB.Model(&models.CiRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "temporarytotaldisability":
		DB.Model(&models.TtdRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "funeralparameters":
		DB.Model(&models.FuneralParameters{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "grouppricingparameters":
		DB.Model(&models.GroupPricingParameters{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "incomelevel":
		DB.Model(&models.IncomeLevel{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "occupationclass":
		DB.Model(&models.OccupationClass{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "educatorriskrates":
		DB.Model(&models.EducatorRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "grouppricingeducatorstructure":
		DB.Model(&models.EducatorBenefitStructure{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "escalations":
		DB.Model(&models.Escalations{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "funeralaidsrates":
		DB.Model(&models.FuneralAidsRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "funeralrates":
		DB.Model(&models.FuneralRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "generalloadings":
		DB.Model(&models.GeneralLoading{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "glaaidsrates":
		DB.Model(&models.GlaAidsRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "regionloading":
		DB.Model(&models.RegionLoading{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "tieredincomereplacement":
		DB.Model(&models.TieredIncomeReplacement{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "customtieredincomereplacement":
		DB.Model(&models.CustomTieredIncomeReplacement{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "discountauthority":
		DB.Model(&models.DiscountAuthority{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "restrictions":
		DB.Model(&models.Restriction{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "reinsurancecoverrestrictions":
		DB.Model(&models.ReinsuranceCoverRestriction{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "reinsuranceglarate":
		DB.Model(&models.ReinsuranceGlaRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "reinsurancecirate":
		DB.Model(&models.ReinsuranceCiRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "reinsuranceptdrate":
		DB.Model(&models.ReinsurancePtdRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "reinsurancephirate":
		DB.Model(&models.ReinsurancePhiRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "reinsurancefuneralaidsrates":
		DB.Model(&models.ReinsuranceFuneralAidsRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "reinsurancefuneralrates":
		DB.Model(&models.ReinsuranceFuneralRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "reinsuranceglaaidsrates":
		DB.Model(&models.ReinsuranceGlaAidsRate{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "reinsurancegeneralloadings":
		DB.Model(&models.ReinsuranceGeneralLoading{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "reinsuranceindustryloadings":
		DB.Model(&models.ReinsuranceIndustryLoading{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "reinsuranceregionloadings":
		DB.Model(&models.ReinsuranceRegionLoading{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "premiumloadings":
		DB.Model(&models.PremiumLoading{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "schemesizelevels":
		DB.Model(&models.SchemeSizeLevel{}).Select("DISTINCT year").Order("year desc").Find(&years)
	case "taxtable":
		DB.Model(&models.TaxTable{}).Select("DISTINCT year").Order("year desc").Find(&years)
	}
	return years
}

// GetGPTableData returns rows for the requested table_type. The return type is
// `any` so each case can return its native struct slice — Gin marshals it
// directly to JSON, preserving struct field order on the wire (which the UI
// uses to lay out columns). Earlier this function unmarshalled through
// []map[string]interface{}, which sorted keys alphabetically and forced the UI
// to maintain a per-table preferred-column-order list. Returning the typed
// slice removes that workaround and keeps the column order in sync with the
// model definition.
func GetGPTableData(tableType string) any {
	switch tableType {
	case "grouplifeassurance":
		var data []models.GlaRate
		DB.Find(&data)
		return data
	case "permanenthealthinsurance":
		var data []models.PhiRate
		if err := DB.Find(&data).Error; err != nil {
			fmt.Println(err)
		}
		return data
	case "childmortality":
		var data []models.ChildMortality
		DB.Find(&data)
		return data
	case "industryloading":
		var data []models.IndustryLoading
		DB.Find(&data)
		return data
	case "grouppricingreinsurancestructure":
		var data []models.GroupPricingReinsuranceStructure
		DB.Find(&data)
		return data
	case "accidentaltemporarytotaldisability":
		var data []models.AccidentalTtdRate
		DB.Find(&data)
		return data
	case "permanentandtotaldisability":
		var data []models.PtdRate
		DB.Find(&data)
		return data
	case "criticalillness":
		var data []models.CiRate
		DB.Find(&data)
		return data
	case "temporarytotaldisability":
		var data []models.TtdRate
		DB.Find(&data)
		return data
	case "funeralparameters":
		var data []models.FuneralParameters
		DB.Find(&data)
		return data
	case "grouppricingparameters":
		var data []models.GroupPricingParameters
		DB.Find(&data)
		return data
	case "incomelevel":
		var data []models.IncomeLevel
		DB.Find(&data)
		return data
	case "occupationclass":
		var data []models.OccupationClass
		DB.Find(&data)
		return data
	case "bordereauxes":
		var data []models.Bordereaux
		DB.Find(&data)
		return data
	case "educatorriskrates":
		var data []models.EducatorRate
		DB.Find(&data)
		return data
	case "grouppricingeducatorstructure":
		var data []models.EducatorBenefitStructure
		DB.Find(&data)
		return data
	case "escalations":
		var data []models.Escalations
		DB.Find(&data)
		return data
	case "funeralaidsrates":
		var data []models.FuneralAidsRate
		DB.Find(&data)
		return data
	case "funeralrates":
		var data []models.FuneralRate
		DB.Find(&data)
		return data
	case "generalloadings":
		var data []models.GeneralLoading
		DB.Find(&data)
		return data
	case "glaaidsrates":
		var data []models.GlaAidsRate
		DB.Find(&data)
		return data
	case "regionloading":
		var data []models.RegionLoading
		DB.Find(&data)
		return data
	case "tieredincomereplacement":
		var data []models.TieredIncomeReplacement
		DB.Find(&data)
		return data
	case "customtieredincomereplacement":
		var data []models.CustomTieredIncomeReplacement
		DB.Find(&data)
		return data
	case "discountauthority":
		var data []models.DiscountAuthority
		DB.Find(&data)
		return data
	case "restrictions":
		var data []models.Restriction
		DB.Find(&data)
		return data
	case "reinsurancecoverrestrictions":
		var data []models.ReinsuranceCoverRestriction
		DB.Find(&data)
		return data
	case "reinsuranceglarate":
		var data []models.ReinsuranceGlaRate
		DB.Find(&data)
		return data
	case "reinsurancecirate":
		var data []models.ReinsuranceCiRate
		DB.Find(&data)
		return data
	case "reinsuranceptdrate":
		var data []models.ReinsurancePtdRate
		DB.Find(&data)
		return data
	case "reinsurancephirate":
		var data []models.ReinsurancePhiRate
		DB.Find(&data)
		return data
	case "reinsurancefuneralaidsrates":
		var data []models.ReinsuranceFuneralAidsRate
		DB.Find(&data)
		return data
	case "reinsurancefuneralrates":
		var data []models.ReinsuranceFuneralRate
		DB.Find(&data)
		return data
	case "reinsuranceglaaidsrates":
		var data []models.ReinsuranceGlaAidsRate
		DB.Find(&data)
		return data
	case "reinsurancegeneralloadings":
		var data []models.ReinsuranceGeneralLoading
		DB.Find(&data)
		return data
	case "reinsuranceindustryloadings":
		var data []models.ReinsuranceIndustryLoading
		DB.Find(&data)
		return data
	case "reinsuranceregionloadings":
		var data []models.ReinsuranceRegionLoading
		DB.Find(&data)
		return data
	case "premiumloadings":
		var data []models.PremiumLoading
		DB.Find(&data)
		return data
	case "schemesizelevels":
		var data []models.SchemeSizeLevel
		DB.Find(&data)
		return data
	case "taxtable":
		var data []models.TaxTable
		DB.Find(&data)
		return data
	case "agebands":
		var data []models.GroupPricingAgeBands
		DB.Find(&data)
		return data
	case "medicalwaiver":
		var data []models.MedicalWaiver
		DB.Find(&data)
		return data
	}
	return nil
}

func GetIncomeLevel(mp models.GPricingMemberData, incomeLevels []models.IncomeLevel) int {
	for _, level := range incomeLevels {
		if mp.AnnualSalary >= level.MinIncome && mp.AnnualSalary <= level.MaxIncome {
			return level.Level
		}
	}
	return 0
}

func GetRestrictionByRiskRateCode(riskRateCode string) (models.Restriction, error) {
	var restriction models.Restriction
	err := DB.Where("risk_rate_code = ?", riskRateCode).First(&restriction).Error
	return restriction, err
}

// LoadReinsuranceCoverCaps returns a benefit-type → maximum cover map for the
// scheme size band that contains memberCount. A row is in-band when
// MinSchemeSize <= memberCount and (MaxSchemeSize = 0 OR memberCount <= MaxSchemeSize)
// — MaxSchemeSize 0 represents the open-ended top band. Missing benefit types
// and stored zero values both mean "no restriction".
func LoadReinsuranceCoverCaps(riskRateCode string, memberCount int) map[string]float64 {
	caps := map[string]float64{}
	if riskRateCode == "" {
		return caps
	}
	var rows []models.ReinsuranceCoverRestriction
	err := DB.Where(
		"risk_rate_code = ? AND min_scheme_size <= ? AND (max_scheme_size = 0 OR max_scheme_size >= ?)",
		riskRateCode, memberCount, memberCount,
	).Find(&rows).Error
	if err != nil {
		return caps
	}
	for _, r := range rows {
		caps[strings.ToUpper(r.BenefitType)] = float64(r.MaximumCover)
	}
	return caps
}

func GetFuneralParameter(groupPricingParameter models.GroupPricingParameters, anb int) models.FuneralParameters {
	tableName := "funeral_parameters"
	var funeralParameters models.FuneralParameters
	var keyString strings.Builder

	//keyString.WriteString(strconv.Itoa(groupPricingParameter.Year) + "_")
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(anb) + "_")
	key := keyString.String()
	cacheKey := tableName + "_" + key
	cached, found := GroupPricingCache.Get(cacheKey)

	if found {
		result := cached.(models.FuneralParameters)
		//if result > 0 {
		return result
		//}
	} else {
		query := "risk_rate_code=? and age_next_birthday=?"
		err := DB.Table(tableName).Where(query, groupPricingParameter.RiskRateCode, anb).First(&funeralParameters).Error
		if err != nil {
			fmt.Println(err)
		}
		GroupPricingCache.Set(cacheKey, funeralParameters, 1)
	}
	return funeralParameters
}

func GetOccupationClass(groupPricingParameter models.GroupPricingParameters, qroupQuote models.GroupPricingQuote, selectedSchemeCategory string) int {
	tableName := "occupation_classes"
	var occupationClass models.OccupationClass
	var keyString strings.Builder

	//keyString.WriteString(strconv.Itoa(groupPricingParameter.Year) + "_")
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(qroupQuote.Industry + "_")
	keyString.WriteString(selectedSchemeCategory + "_")
	key := keyString.String()
	cacheKey := tableName + "_" + key
	cached, found := GroupPricingCache.Get(cacheKey)

	if found {
		result := cached.(int)
		//if result > 0 {
		return result
		//}
	} else {
		//fmt.Println("cache missed: ", key)
	}
	query := "risk_rate_code=? and lower(industry)=lower(?) and lower(category)= lower(?)"
	err := DB.Table(tableName).Where(query, groupPricingParameter.RiskRateCode, qroupQuote.Industry, selectedSchemeCategory).First(&occupationClass).Error
	if err != nil {
		fmt.Println(err)
		GroupPricingCache.Set(cacheKey, 0, 1)
		return 0
	}

	GroupPricingCache.Set(cacheKey, occupationClass.Class, 1)
	return occupationClass.Class
}

func GetAgeBand(ageNextBirthday int, ageBand []models.GroupPricingAgeBands) string {
	for _, band := range ageBand {
		if ageNextBirthday >= band.MinAge && ageNextBirthday <= band.MaxAge {
			return band.Name
		}
	}
	return "0"
}

func GetGlaRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, incomeLevel int, groupQuote models.GroupPricingQuote, schemeCategory models.SchemeCategory) float64 {
	tableName := "gla_rates"
	var keyString strings.Builder

	keyString.WriteString(schemeCategory.SchemeCategory + "_")
	//keyString.WriteString(strconv.Itoa(groupPricingParameter.Year) + "_")
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.AgeNextBirthday) + "_")
	keyString.WriteString(strconv.Itoa(incomeLevel) + "_") //incomelevel
	keyString.WriteString(memberResultData.Gender[:1] + "_")
	keyString.WriteString(strconv.Itoa(schemeCategory.GlaWaitingPeriod) + "_")
	keyString.WriteString(schemeCategory.GlaBenefitType + "_")
	key := keyString.String()
	cacheKey := tableName + "_" + key
	cached, found := GroupPricingCache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	} else {
		//fmt.Println("cache missed: ", key)
	}
	var qx float64
	query := "risk_rate_code=? and age_next_birthday=? and income_level=? and gender=? and waiting_period=? and benefit_type=?"
	err := DB.Table(tableName).Where(query, groupPricingParameter.RiskRateCode, memberResultData.AgeNextBirthday, incomeLevel, memberResultData.Gender, schemeCategory.GlaWaitingPeriod, schemeCategory.GlaBenefitType).Pluck("qx", &qx).Error //.Select("qx").Row()
	if err != nil {
		fmt.Println(err)
	}
	GroupPricingCache.Set(cacheKey, qx, 1)
	//time.Sleep(5 * time.Millisecond)
	return qx
}

// GetMedicalWaiverSumAtRisk looks up the medical waiver sum-at-risk from the
// medical_waivers table by (risk_rate_code, age_next_birthday, income_level,
// gender). Returns 0 when no row matches, which the caller treats as a zero
// PhiMedicalAidWaiver under the table_lookup methodology.
func GetMedicalWaiverSumAtRisk(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, incomeLevel int) float64 {
	tableName := "medical_waivers"
	var keyString strings.Builder
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.AgeNextBirthday) + "_")
	keyString.WriteString(strconv.Itoa(incomeLevel) + "_")
	if len(memberResultData.Gender) > 0 {
		keyString.WriteString(memberResultData.Gender[:1] + "_")
	}
	cacheKey := tableName + "_" + keyString.String()
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}
	var sumAtRisk float64
	query := "risk_rate_code=? and age_next_birthday=? and income_level=? and gender=?"
	if err := DB.Table(tableName).Where(query,
		groupPricingParameter.RiskRateCode,
		memberResultData.AgeNextBirthday,
		incomeLevel,
		memberResultData.Gender,
	).Pluck("medicalwaiver_sum_at_risk", &sumAtRisk).Error; err != nil {
		fmt.Println(err)
	}
	GroupPricingCache.Set(cacheKey, sumAtRisk, 1)
	return sumAtRisk
}

// GetAdditionalAccidentalGlaRate looks up the qx from gla_rates using the
// benefit type configured for the optional Additional Accidental GLA layer.
// The additional layer re-uses every other GLA parameter (risk rate code,
// waiting period, income level, age, gender); only the benefit_type differs.
func GetAdditionalAccidentalGlaRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, incomeLevel int, groupQuote models.GroupPricingQuote, schemeCategory models.SchemeCategory) float64 {
	tableName := "gla_rates"
	var keyString strings.Builder

	keyString.WriteString(schemeCategory.SchemeCategory + "_")
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.AgeNextBirthday) + "_")
	keyString.WriteString(strconv.Itoa(incomeLevel) + "_")
	keyString.WriteString(memberResultData.Gender[:1] + "_")
	keyString.WriteString(strconv.Itoa(schemeCategory.GlaWaitingPeriod) + "_")
	keyString.WriteString(schemeCategory.AdditionalAccidentalGlaBenefitType + "_")
	key := keyString.String()
	cacheKey := "additional_accidental_" + tableName + "_" + key
	cached, found := GroupPricingCache.Get(cacheKey)

	if found {
		result := cached.(float64)
		return result
	}
	var qx float64
	query := "risk_rate_code=? and age_next_birthday=? and income_level=? and gender=? and waiting_period=? and benefit_type=?"
	err := DB.Table(tableName).Where(query, groupPricingParameter.RiskRateCode, memberResultData.AgeNextBirthday, incomeLevel, memberResultData.Gender, schemeCategory.GlaWaitingPeriod, schemeCategory.AdditionalAccidentalGlaBenefitType).Pluck("qx", &qx).Error
	if err != nil {
		fmt.Println(err)
	}
	GroupPricingCache.Set(cacheKey, qx, 1)
	return qx
}

func GetGlaAidsRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters) float64 {
	tableName := "gla_aids_rates"

	var keyString strings.Builder
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.AgeNextBirthday) + "_")
	keyString.WriteString(memberResultData.Gender[:1] + "_")
	key := keyString.String()

	cacheKey := tableName + "_" + key
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}

	var qx float64
	err := DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND gender = ?",
			groupPricingParameter.RiskRateCode,
			memberResultData.AgeNextBirthday,
			memberResultData.Gender[:1],
		).
		Pluck("gla_aids_qx", &qx).Error
	if err != nil {
		fmt.Println(err)
	}

	GroupPricingCache.Set(cacheKey, qx, 1)
	return qx
}

func GetSpouseGlaAidsRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters) float64 {
	tableName := "gla_aids_rates"

	var keyString strings.Builder
	keyString.WriteString("spouse_")
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.SpouseAgeNextBirthday) + "_")
	keyString.WriteString(memberResultData.SpouseGender[:1] + "_")
	key := keyString.String()

	cacheKey := tableName + "_" + key
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}

	var qx float64
	err := DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND gender = ?",
			groupPricingParameter.RiskRateCode,
			memberResultData.SpouseAgeNextBirthday,
			memberResultData.SpouseGender[:1],
		).
		Pluck("gla_aids_qx", &qx).Error
	if err != nil {
		fmt.Println(err)
	}

	GroupPricingCache.Set(cacheKey, qx, 1)
	return qx
}

func GetRegionLoading(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, region string) models.RegionLoading {
	// Short-circuit when the table is configured as not required: skip both
	// the cache lookup and the DB read, returning a zero-value struct so
	// downstream loadings resolve to 0.
	if !IsTableRequired("regionLoading") {
		return models.RegionLoading{}
	}

	tableName := "region_loadings"

	var keyString strings.Builder
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(memberResultData.Gender[:1] + "_")
	keyString.WriteString(strings.TrimSpace(region) + "_")
	key := keyString.String()

	cacheKey := tableName + "_" + key
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(models.RegionLoading)
	}

	var regionLoading models.RegionLoading
	err := DB.Table(tableName).
		Where("risk_rate_code = ? AND gender = ? AND region = ?",
			groupPricingParameter.RiskRateCode,
			memberResultData.Gender[:1],
			strings.TrimSpace(region),
		).
		First(&regionLoading).Error
	if err != nil {
		fmt.Println(err)
	}

	GroupPricingCache.Set(cacheKey, regionLoading, 1)
	return regionLoading
}

func GetPremiumLoading(groupPricingParameter models.GroupPricingParameters, schemeSizeLevel int, channel string) models.PremiumLoading {
	if !IsTableRequired("premiumLoading") {
		return models.PremiumLoading{}
	}

	tableName := "premium_loadings"

	var keyString strings.Builder
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(schemeSizeLevel) + "_")
	keyString.WriteString(strings.TrimSpace(channel) + "_")
	key := keyString.String()

	cacheKey := tableName + "_" + key
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(models.PremiumLoading)
	}

	var premiumLoading models.PremiumLoading
	err := DB.Table(tableName).
		Where("risk_rate_code = ? AND scheme_size_level = ? AND channel = ?",
			groupPricingParameter.RiskRateCode,
			schemeSizeLevel,
			strings.TrimSpace(channel),
		).
		First(&premiumLoading).Error
	if err != nil {
		fmt.Println(err)
	}

	GroupPricingCache.Set(cacheKey, premiumLoading, 1)
	return premiumLoading
}

func GetGeneralLoading(riskRateCode string, age int, gender string) models.GeneralLoading {
	if !IsTableRequired("generalLoading") {
		return models.GeneralLoading{}
	}

	tableName := "general_loadings"

	cacheKey := tableName + "_" + riskRateCode + "_" + strconv.Itoa(age) + "_" + gender
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(models.GeneralLoading)
	}

	var generalLoading models.GeneralLoading
	err := DB.Table(tableName).
		Where("risk_rate_code = ? AND age = ? AND gender = ?", riskRateCode, age, gender).
		First(&generalLoading).Error
	if err != nil {
		fmt.Println(err)
	}

	GroupPricingCache.Set(cacheKey, generalLoading, 1)
	return generalLoading
}

// benefitAccessor is the canonical per-benefit getter/setter table used by
// both the Exp* commission allocation in applySchemeWideCommission and the
// Final* premium + commission allocation in recomputeFinalPremiumsAndCommission.
// Adding a new benefit (one that participates in TotalAnnualPremium /
// FinalTotalAnnualPremium) is a single-row change here.
type benefitAccessor struct {
	name              string
	bookRiskPremium   func(*models.MemberRatingResultSummary) float64
	expRiskPremium    func(*models.MemberRatingResultSummary) float64
	expBinder         func(*models.MemberRatingResultSummary) float64
	expOutsource      func(*models.MemberRatingResultSummary) float64
	setBookComm       func(*models.MemberRatingResultSummary, float64)
	setExpComm        func(*models.MemberRatingResultSummary, float64)
	setFinalPremium   func(*models.MemberRatingResultSummary, float64)
	setFinalComm      func(*models.MemberRatingResultSummary, float64)
	setFinalBinder    func(*models.MemberRatingResultSummary, float64)
	setFinalOutsource func(*models.MemberRatingResultSummary, float64)
	includeInExclFun  bool
}

// benefitAccessors returns the 11-benefit table mirroring the composition of
// ExpTotalAnnualPremiumExclFuneral (+ funeral). The order is significant:
// the last accessor receives the residual when commission is sliced.
func benefitAccessors() []benefitAccessor {
	return []benefitAccessor{
		{
			name:              "gla",
			bookRiskPremium:   func(s *models.MemberRatingResultSummary) float64 { return s.TotalGlaAnnualRiskPremium },
			expRiskPremium:    func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalGlaAnnualRiskPremium },
			expBinder:         func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalGlaAnnualBinderAmount },
			expOutsource:      func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalGlaAnnualOutsourcedAmount },
			setBookComm:       func(s *models.MemberRatingResultSummary, v float64) { s.TotalGlaAnnualCommissionAmount = v },
			setExpComm:        func(s *models.MemberRatingResultSummary, v float64) { s.ExpTotalGlaAnnualCommissionAmount = v },
			setFinalPremium:   func(s *models.MemberRatingResultSummary, v float64) { s.FinalGlaAnnualOfficePremium = v },
			setFinalComm:      func(s *models.MemberRatingResultSummary, v float64) { s.FinalGlaAnnualCommissionAmount = v },
			setFinalBinder:    func(s *models.MemberRatingResultSummary, v float64) { s.FinalGlaAnnualBinderAmount = v },
			setFinalOutsource: func(s *models.MemberRatingResultSummary, v float64) { s.FinalGlaAnnualOutsourcedAmount = v },
			includeInExclFun:  true,
		},
		{
			name: "add_acc_gla",
			bookRiskPremium: func(s *models.MemberRatingResultSummary) float64 {
				return s.TotalAdditionalAccidentalGlaAnnualRiskPremium
			},
			expRiskPremium: func(s *models.MemberRatingResultSummary) float64 {
				return s.ExpTotalAdditionalAccidentalGlaAnnualRiskPremium
			},
			expBinder: func(s *models.MemberRatingResultSummary) float64 {
				return s.ExpTotalAdditionalAccidentalGlaAnnualBinderAmount
			},
			expOutsource: func(s *models.MemberRatingResultSummary) float64 {
				return s.ExpTotalAdditionalAccidentalGlaAnnualOutsourcedAmt
			},
			setBookComm: func(s *models.MemberRatingResultSummary, v float64) {
				s.TotalAdditionalAccidentalGlaAnnualCommissionAmount = v
			},
			setExpComm: func(s *models.MemberRatingResultSummary, v float64) {
				s.ExpTotalAdditionalAccidentalGlaAnnualCommissionAmount = v
			},
			setFinalPremium: func(s *models.MemberRatingResultSummary, v float64) {
				s.FinalAdditionalAccidentalGlaAnnualOfficePremium = v
			},
			setFinalComm: func(s *models.MemberRatingResultSummary, v float64) {
				s.FinalAdditionalAccidentalGlaAnnualCommissionAmount = v
			},
			setFinalBinder: func(s *models.MemberRatingResultSummary, v float64) {
				s.FinalAdditionalAccidentalGlaAnnualBinderAmount = v
			},
			setFinalOutsource: func(s *models.MemberRatingResultSummary, v float64) {
				s.FinalAdditionalAccidentalGlaAnnualOutsourcedAmt = v
			},
			includeInExclFun: true,
		},
		{
			name:              "ptd",
			bookRiskPremium:   func(s *models.MemberRatingResultSummary) float64 { return s.TotalPtdAnnualRiskPremium },
			expRiskPremium:    func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalPtdAnnualRiskPremium },
			expBinder:         func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalPtdAnnualBinderAmount },
			expOutsource:      func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalPtdAnnualOutsourcedAmount },
			setBookComm:       func(s *models.MemberRatingResultSummary, v float64) { s.TotalPtdAnnualCommissionAmount = v },
			setExpComm:        func(s *models.MemberRatingResultSummary, v float64) { s.ExpTotalPtdAnnualCommissionAmount = v },
			setFinalPremium:   func(s *models.MemberRatingResultSummary, v float64) { s.FinalPtdAnnualOfficePremium = v },
			setFinalComm:      func(s *models.MemberRatingResultSummary, v float64) { s.FinalPtdAnnualCommissionAmount = v },
			setFinalBinder:    func(s *models.MemberRatingResultSummary, v float64) { s.FinalPtdAnnualBinderAmount = v },
			setFinalOutsource: func(s *models.MemberRatingResultSummary, v float64) { s.FinalPtdAnnualOutsourcedAmount = v },
			includeInExclFun:  true,
		},
		{
			name:              "ci",
			bookRiskPremium:   func(s *models.MemberRatingResultSummary) float64 { return s.TotalCiAnnualRiskPremium },
			expRiskPremium:    func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalCiAnnualRiskPremium },
			expBinder:         func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalCiAnnualBinderAmount },
			expOutsource:      func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalCiAnnualOutsourcedAmount },
			setBookComm:       func(s *models.MemberRatingResultSummary, v float64) { s.TotalCiAnnualCommissionAmount = v },
			setExpComm:        func(s *models.MemberRatingResultSummary, v float64) { s.ExpTotalCiAnnualCommissionAmount = v },
			setFinalPremium:   func(s *models.MemberRatingResultSummary, v float64) { s.FinalCiAnnualOfficePremium = v },
			setFinalComm:      func(s *models.MemberRatingResultSummary, v float64) { s.FinalCiAnnualCommissionAmount = v },
			setFinalBinder:    func(s *models.MemberRatingResultSummary, v float64) { s.FinalCiAnnualBinderAmount = v },
			setFinalOutsource: func(s *models.MemberRatingResultSummary, v float64) { s.FinalCiAnnualOutsourcedAmount = v },
			includeInExclFun:  true,
		},
		{
			name:              "sgla",
			bookRiskPremium:   func(s *models.MemberRatingResultSummary) float64 { return s.TotalSglaAnnualRiskPremium },
			expRiskPremium:    func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalSglaAnnualRiskPremium },
			expBinder:         func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalSglaAnnualBinderAmount },
			expOutsource:      func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalSglaAnnualOutsourcedAmount },
			setBookComm:       func(s *models.MemberRatingResultSummary, v float64) { s.TotalSglaAnnualCommissionAmount = v },
			setExpComm:        func(s *models.MemberRatingResultSummary, v float64) { s.ExpTotalSglaAnnualCommissionAmount = v },
			setFinalPremium:   func(s *models.MemberRatingResultSummary, v float64) { s.FinalSglaAnnualOfficePremium = v },
			setFinalComm:      func(s *models.MemberRatingResultSummary, v float64) { s.FinalSglaAnnualCommissionAmount = v },
			setFinalBinder:    func(s *models.MemberRatingResultSummary, v float64) { s.FinalSglaAnnualBinderAmount = v },
			setFinalOutsource: func(s *models.MemberRatingResultSummary, v float64) { s.FinalSglaAnnualOutsourcedAmount = v },
			includeInExclFun:  true,
		},
		{
			// Tax saver has no persisted Exp* binder/outsource aggregates
			// (it's an informational rider on top of GLA), so the Exp
			// getters return 0 — finalBinder/finalOutsource scale to 0
			// for this benefit, matching the calc-time behaviour.
			name:              "tax_saver",
			bookRiskPremium:   func(s *models.MemberRatingResultSummary) float64 { return s.TotalTaxSaverAnnualRiskPremium },
			expRiskPremium:    func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalTaxSaverAnnualRiskPremium },
			expBinder:         func(s *models.MemberRatingResultSummary) float64 { return 0 },
			expOutsource:      func(s *models.MemberRatingResultSummary) float64 { return 0 },
			setBookComm:       func(s *models.MemberRatingResultSummary, v float64) { s.TotalTaxSaverAnnualCommissionAmount = v },
			setExpComm:        func(s *models.MemberRatingResultSummary, v float64) { s.ExpTotalTaxSaverAnnualCommissionAmount = v },
			setFinalPremium:   func(s *models.MemberRatingResultSummary, v float64) { s.FinalTaxSaverAnnualOfficePremium = v },
			setFinalComm:      func(s *models.MemberRatingResultSummary, v float64) { s.FinalTaxSaverAnnualCommissionAmount = v },
			setFinalBinder:    func(s *models.MemberRatingResultSummary, v float64) { s.FinalTaxSaverAnnualBinderAmount = v },
			setFinalOutsource: func(s *models.MemberRatingResultSummary, v float64) { s.FinalTaxSaverAnnualOutsourcedAmount = v },
			includeInExclFun:  true,
		},
		{
			name:              "ttd",
			bookRiskPremium:   func(s *models.MemberRatingResultSummary) float64 { return s.TotalTtdAnnualRiskPremium },
			expRiskPremium:    func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalTtdAnnualRiskPremium },
			expBinder:         func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalTtdAnnualBinderAmount },
			expOutsource:      func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalTtdAnnualOutsourcedAmount },
			setBookComm:       func(s *models.MemberRatingResultSummary, v float64) { s.TotalTtdAnnualCommissionAmount = v },
			setExpComm:        func(s *models.MemberRatingResultSummary, v float64) { s.ExpTotalTtdAnnualCommissionAmount = v },
			setFinalPremium:   func(s *models.MemberRatingResultSummary, v float64) { s.FinalTtdAnnualOfficePremium = v },
			setFinalComm:      func(s *models.MemberRatingResultSummary, v float64) { s.FinalTtdAnnualCommissionAmount = v },
			setFinalBinder:    func(s *models.MemberRatingResultSummary, v float64) { s.FinalTtdAnnualBinderAmount = v },
			setFinalOutsource: func(s *models.MemberRatingResultSummary, v float64) { s.FinalTtdAnnualOutsourcedAmount = v },
			includeInExclFun:  true,
		},
		{
			name:              "phi",
			bookRiskPremium:   func(s *models.MemberRatingResultSummary) float64 { return s.TotalPhiAnnualRiskPremium },
			expRiskPremium:    func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalPhiAnnualRiskPremium },
			expBinder:         func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalPhiAnnualBinderAmount },
			expOutsource:      func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalPhiAnnualOutsourcedAmount },
			setBookComm:       func(s *models.MemberRatingResultSummary, v float64) { s.TotalPhiAnnualCommissionAmount = v },
			setExpComm:        func(s *models.MemberRatingResultSummary, v float64) { s.ExpTotalPhiAnnualCommissionAmount = v },
			setFinalPremium:   func(s *models.MemberRatingResultSummary, v float64) { s.FinalPhiAnnualOfficePremium = v },
			setFinalComm:      func(s *models.MemberRatingResultSummary, v float64) { s.FinalPhiAnnualCommissionAmount = v },
			setFinalBinder:    func(s *models.MemberRatingResultSummary, v float64) { s.FinalPhiAnnualBinderAmount = v },
			setFinalOutsource: func(s *models.MemberRatingResultSummary, v float64) { s.FinalPhiAnnualOutsourcedAmount = v },
			includeInExclFun:  true,
		},
		{
			name:              "fun",
			bookRiskPremium:   func(s *models.MemberRatingResultSummary) float64 { return s.TotalFunAnnualRiskPremium },
			expRiskPremium:    func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalFunAnnualRiskPremium },
			expBinder:         func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalFunAnnualBinderAmount },
			expOutsource:      func(s *models.MemberRatingResultSummary) float64 { return s.ExpTotalFunAnnualOutsourcedAmount },
			setBookComm:       func(s *models.MemberRatingResultSummary, v float64) { s.TotalFunAnnualCommissionAmount = v },
			setExpComm:        func(s *models.MemberRatingResultSummary, v float64) { s.ExpTotalFunAnnualCommissionAmount = v },
			setFinalPremium:   func(s *models.MemberRatingResultSummary, v float64) { s.FinalFunAnnualOfficePremium = v },
			setFinalComm:      func(s *models.MemberRatingResultSummary, v float64) { s.FinalFunAnnualCommissionAmount = v },
			setFinalBinder:    func(s *models.MemberRatingResultSummary, v float64) { s.FinalFunAnnualBinderAmount = v },
			setFinalOutsource: func(s *models.MemberRatingResultSummary, v float64) { s.FinalFunAnnualOutsourcedAmount = v },
			includeInExclFun:  false,
		},
		{
			name:              "gla_educator",
			bookRiskPremium:   func(s *models.MemberRatingResultSummary) float64 { return s.TotalGlaEducatorRiskPremium },
			expRiskPremium:    func(s *models.MemberRatingResultSummary) float64 { return s.ExpAdjTotalGlaEducatorRiskPremium },
			expBinder:         func(s *models.MemberRatingResultSummary) float64 { return s.ExpAdjTotalGlaEducatorBinderAmount },
			expOutsource:      func(s *models.MemberRatingResultSummary) float64 { return s.ExpAdjTotalGlaEducatorOutsourcedAmount },
			setBookComm:       func(s *models.MemberRatingResultSummary, v float64) { s.TotalGlaEducatorCommissionAmount = v },
			setExpComm:        func(s *models.MemberRatingResultSummary, v float64) { s.ExpAdjTotalGlaEducatorCommissionAmount = v },
			setFinalPremium:   func(s *models.MemberRatingResultSummary, v float64) { s.FinalGlaEducatorAnnualOfficePremium = v },
			setFinalComm:      func(s *models.MemberRatingResultSummary, v float64) { s.FinalGlaEducatorAnnualCommissionAmount = v },
			setFinalBinder:    func(s *models.MemberRatingResultSummary, v float64) { s.FinalGlaEducatorAnnualBinderAmount = v },
			setFinalOutsource: func(s *models.MemberRatingResultSummary, v float64) { s.FinalGlaEducatorAnnualOutsourcedAmount = v },
			includeInExclFun:  true,
		},
		{
			name:              "ptd_educator",
			bookRiskPremium:   func(s *models.MemberRatingResultSummary) float64 { return s.TotalPtdEducatorRiskPremium },
			expRiskPremium:    func(s *models.MemberRatingResultSummary) float64 { return s.ExpAdjTotalPtdEducatorRiskPremium },
			expBinder:         func(s *models.MemberRatingResultSummary) float64 { return s.ExpAdjTotalPtdEducatorBinderAmount },
			expOutsource:      func(s *models.MemberRatingResultSummary) float64 { return s.ExpAdjTotalPtdEducatorOutsourcedAmount },
			setBookComm:       func(s *models.MemberRatingResultSummary, v float64) { s.TotalPtdEducatorCommissionAmount = v },
			setExpComm:        func(s *models.MemberRatingResultSummary, v float64) { s.ExpAdjTotalPtdEducatorCommissionAmount = v },
			setFinalPremium:   func(s *models.MemberRatingResultSummary, v float64) { s.FinalPtdEducatorAnnualOfficePremium = v },
			setFinalComm:      func(s *models.MemberRatingResultSummary, v float64) { s.FinalPtdEducatorAnnualCommissionAmount = v },
			setFinalBinder:    func(s *models.MemberRatingResultSummary, v float64) { s.FinalPtdEducatorAnnualBinderAmount = v },
			setFinalOutsource: func(s *models.MemberRatingResultSummary, v float64) { s.FinalPtdEducatorAnnualOutsourcedAmount = v },
			includeInExclFun:  true,
		},
	}
}

// recomputeFinalPremiumsAndCommission writes Final*OfficePremium and
// Final*CommissionAmount onto every summary for the quote using the Discount
// currently stored on each summary. Each persisted Final*OfficePremium is the
// pre-commission office premium (RiskPremium / (1 - FinalSchemeTotalLoading))
// plus that benefit's commission slice — so the per-benefit values include
// commission and reconcile exactly to FinalTotalAnnualPremium{,ExclFuneral}.
// The progressive commission rate is derived against the scheme-wide pre-comm
// total via ComputeProgressiveCommission, then split scheme → category →
// benefit by pre-comm Final premium share with the last slice absorbing the
// residual (mirroring the Exp* allocation pattern).
//
// Note: Exp*OfficePremium remains pre-commission, so post-change the
// "Discount == 0 ⇒ Final == Exp" invariant no longer holds — the gap equals
// SchemeTotalCommission.
func recomputeFinalPremiumsAndCommission(quoteID int, groupQuote models.GroupPricingQuote) error {
	var summaries []models.MemberRatingResultSummary
	if err := DB.Where("quote_id = ?", quoteID).Order("id ASC").Find(&summaries).Error; err != nil {
		return fmt.Errorf("load summaries for final recompute: %w", err)
	}
	if len(summaries) == 0 {
		return nil
	}

	accessors := benefitAccessors()
	discountMethod := GetDiscountMethod()

	// 1. Seed pre-comm Final*OfficePremium for every benefit on every summary
	//    and the FinalTotalAnnualPremium{ExclFuneral,} rollups derived from
	//    them. These are overwritten in pass 3 once commission is allocated;
	//    the seed is kept so pass 1's persisted values are self-consistent if
	//    a later pass bails out. Also derive Final*Binder/Outsourced amounts
	//    by scaling each benefit's persisted Exp*Binder/Outsource by the
	//    pre-comm discount ratio OfficeFinalPreComm / OfficePre. At
	//    Discount == 0 the ratio is 1, so Final == Exp exactly; at Discount
	//    != 0 it shrinks (or grows) the Exp amount by exactly the discount-
	//    driven loading change. This avoids depending on s.BinderFeeRate /
	//    s.OutsourceFeeRate, which are 0 on legacy summaries that pre-date
	//    the rate columns and on non-binder distribution channels.
	categoryTotals := make([]float64, len(summaries))
	schemeTotal := 0.0
	for i := range summaries {
		s := &summaries[i]
		var exclFun, total float64
		for _, acc := range accessors {
			risk := acc.expRiskPremium(s)
			expOffice := models.ComputeOfficePremium(risk, s)
			var finalOfficePreComm float64
			switch discountMethod {
			case models.DiscountMethodProrata:
				finalOfficePreComm = models.ComputeProrataFinalOfficePremium(risk, s)
			default:
				finalOfficePreComm = models.ComputeFinalOfficePremium(risk, s)
			}
			acc.setFinalPremium(s, finalOfficePreComm)

			scale := 0.0
			if expOffice > 0 {
				scale = finalOfficePreComm / expOffice
			}
			acc.setFinalBinder(s, acc.expBinder(s)*scale)
			acc.setFinalOutsource(s, acc.expOutsource(s)*scale)

			total += finalOfficePreComm
			if acc.includeInExclFun {
				exclFun += finalOfficePreComm
			}
		}
		s.FinalTotalAnnualPremiumExclFuneral = exclFun
		s.FinalTotalAnnualPremium = total
		categoryTotals[i] = total
		schemeTotal += total
	}

	// 2. Re-derive the progressive commission rate against the new scheme total.
	channel := string(groupQuote.DistributionChannel)
	holderName := strings.TrimSpace(groupQuote.QuoteBroker.Name)
	rate := 0.0
	if schemeTotal > 0 {
		r, err := ComputeProgressiveCommission(channel, holderName, schemeTotal)
		if err != nil {
			return fmt.Errorf("compute final progressive commission: %w", err)
		}
		rate = r
	}
	overallCommission := schemeTotal * rate

	// 3. Split scheme commission into per-category commissions by share, with
	//    the last category absorbing any rounding residual.
	categoryCommissions := make([]float64, len(summaries))
	lastCatIdx := len(summaries) - 1
	remainingScheme := overallCommission
	for i := range summaries {
		if i == lastCatIdx {
			categoryCommissions[i] = remainingScheme
		} else if schemeTotal > 0 {
			categoryCommissions[i] = overallCommission * (categoryTotals[i] / schemeTotal)
			remainingScheme -= categoryCommissions[i]
		}
	}

	// 4. Within each category, split commission by benefit share (last benefit
	//    absorbs the per-category residual). Then gross up Final*OfficePremium
	//    in place to (pre-comm + commission slice) so the final premium itself
	//    includes commission, and rewrite FinalTotalAnnualPremium{,ExclFuneral}
	//    as the sum of the with-commission per-benefit values. Per-benefit
	//    values reconcile to the rollup because the last slice absorbs all
	//    floating-point residual.
	lastBenIdx := len(accessors) - 1
	for catIdx := range summaries {
		s := &summaries[catIdx]
		categoryCommission := categoryCommissions[catIdx]
		categoryTotal := categoryTotals[catIdx]
		remainingCat := categoryCommission

		var newTotal, newExclFun float64
		for benIdx, acc := range accessors {
			risk := acc.expRiskPremium(s)
			var benefitPremium float64
			switch discountMethod {
			case models.DiscountMethodProrata:
				benefitPremium = models.ComputeProrataFinalOfficePremium(risk, s)
			default:
				benefitPremium = models.ComputeFinalOfficePremium(risk, s)
			}

			var slice float64
			if benIdx == lastBenIdx {
				slice = remainingCat
			} else if categoryTotal > 0 {
				slice = categoryCommission * (benefitPremium / categoryTotal)
				remainingCat -= slice
			}
			acc.setFinalComm(s, slice)

			withComm := benefitPremium + slice
			acc.setFinalPremium(s, withComm)
			newTotal += withComm
			if acc.includeInExclFun {
				newExclFun += withComm
			}
		}
		s.FinalTotalAnnualPremium = newTotal
		s.FinalTotalAnnualPremiumExclFuneral = newExclFun
		s.FinalSchemeTotalCommission = overallCommission
		s.FinalSchemeTotalCommissionRate = rate
	}

	// 5. Back-patch Additional GLA Cover band rates so commission_per1000
	//    reflects the actual progressive commission rate rather than the
	//    configured commission loading used at initial band-rate compute
	//    time. SchemeTotalLoading() covers expense + admin + profit + other +
	//    binder + outsource (no commission); progressive commission is added
	//    on top via `rate`. Office, binder, outsource, and commission are
	//    recomputed together so the columns reconcile to office_rate_per1000.
	for i := range summaries {
		s := &summaries[i]
		if len(s.AdditionalGlaCoverBandRates) == 0 {
			continue
		}
		divisor := 1.0 - (s.SchemeTotalLoading() + rate)
		if divisor <= 0 {
			divisor = 1.0
		}
		// In prorata mode every component (including AGLA per-1000 rates)
		// shrinks by (1 - d), so multiply the post-divisor rate by
		// (1 + s.Discount). At Discount == 0 prorataMul == 1 and behaviour is
		// identical to loading_adjustment.
		prorataMul := 1.0
		if discountMethod == models.DiscountMethodProrata {
			prorataMul = 1.0 + s.Discount
		}
		for j := range s.AdditionalGlaCoverBandRates {
			b := &s.AdditionalGlaCoverBandRates[j]
			if b.RiskRatePer1000 == 0 {
				continue
			}
			risk := b.RiskRatePer1000 / 1000.0
			officePer1000 := (risk / divisor) * 1000.0 * prorataMul
			// Snapshot the freshly computed pre-smoothed rate so the
			// smoothing comparison view always has the unsmoothed
			// reference, regardless of whether smoothing has been
			// applied. Same approach for the gender variants below.
			origC := officePer1000
			b.OriginalOfficeRatePer1000 = &origC

			riskM := b.RiskRatePer1000Male / 1000.0
			officeM := (riskM / divisor) * 1000.0 * prorataMul
			origM := officeM
			b.OriginalOfficeRatePer1000Male = &origM

			riskF := b.RiskRatePer1000Female / 1000.0
			officeF := (riskF / divisor) * 1000.0 * prorataMul
			origF := officeF
			b.OriginalOfficeRatePer1000Female = &origF

			// If smoothing has been applied for this band/gender, the
			// final office rate equals the smoothed value; otherwise it
			// equals the freshly computed rate. Downstream consumers
			// (PDFs, certificates, exports) just read OfficeRatePer1000
			// and get the right thing automatically.
			if b.SmoothedOfficeRatePer1000 != nil {
				b.OfficeRatePer1000 = *b.SmoothedOfficeRatePer1000
			} else {
				b.OfficeRatePer1000 = officePer1000
			}
			if b.SmoothedOfficeRatePer1000Male != nil {
				b.OfficeRatePer1000Male = *b.SmoothedOfficeRatePer1000Male
			} else {
				b.OfficeRatePer1000Male = officeM
			}
			if b.SmoothedOfficeRatePer1000Female != nil {
				b.OfficeRatePer1000Female = *b.SmoothedOfficeRatePer1000Female
			} else {
				b.OfficeRatePer1000Female = officeF
			}
			b.BinderFeePer1000 = b.OfficeRatePer1000 * s.BinderFeeRate
			b.BinderFeePer1000Male = b.OfficeRatePer1000Male * s.BinderFeeRate
			b.BinderFeePer1000Female = b.OfficeRatePer1000Female * s.BinderFeeRate
			b.OutsourceFeePer1000 = b.OfficeRatePer1000 * s.OutsourceFeeRate
			b.OutsourceFeePer1000Male = b.OfficeRatePer1000Male * s.OutsourceFeeRate
			b.OutsourceFeePer1000Female = b.OfficeRatePer1000Female * s.OutsourceFeeRate
			b.CommissionPer1000 = b.OfficeRatePer1000 * rate
			b.CommissionPer1000Male = b.OfficeRatePer1000Male * rate
			b.CommissionPer1000Female = b.OfficeRatePer1000Female * rate
		}
	}

	// 6. Conversion / continuity slice office premiums. Sliced from each
	//    parent benefit's exp-adj risk premium via the same discount method as
	//    the parent benefit. Pre-discount these collapse to the
	//    ComputeOfficePremium values; post-discount they re-derive against
	//    the discounted loading denominator. Doc-template tokens read these
	//    directly so render-time math doesn't have to gross up risk values.
	finalSliceOffice := func(s *models.MemberRatingResultSummary, risk float64) float64 {
		switch discountMethod {
		case models.DiscountMethodProrata:
			return models.ComputeProrataFinalOfficePremium(risk, s)
		default:
			return models.ComputeFinalOfficePremium(risk, s)
		}
	}
	for i := range summaries {
		s := &summaries[i]
		s.FinalGlaConversionOnWithdrawalOfficePremium = finalSliceOffice(s, s.ExpAdjTotalGlaConversionOnWithdrawalAnnualRiskPremium)
		s.FinalGlaConversionOnRetirementOfficePremium = finalSliceOffice(s, s.ExpAdjTotalGlaConversionOnRetirementAnnualRiskPremium)
		s.FinalGlaContinuityDuringDisabilityOfficePremium = finalSliceOffice(s, s.ExpAdjTotalGlaContinuityDuringDisabilityAnnualRiskPremium)
		s.FinalGlaEducatorConversionOnWithdrawalOfficePremium = finalSliceOffice(s, s.ExpAdjTotalGlaEducatorConversionOnWithdrawalAnnualRiskPremium)
		s.FinalGlaEducatorConversionOnRetirementOfficePremium = finalSliceOffice(s, s.ExpAdjTotalGlaEducatorConversionOnRetirementAnnualRiskPremium)
		s.FinalGlaEducatorContinuityDuringDisabilityOfficePremium = finalSliceOffice(s, s.ExpAdjTotalGlaEducatorContinuityDuringDisabilityAnnualRiskPremium)
		s.FinalPtdConversionOnWithdrawalOfficePremium = finalSliceOffice(s, s.ExpAdjTotalPtdConversionOnWithdrawalAnnualRiskPremium)
		s.FinalPtdEducatorConversionOnWithdrawalOfficePremium = finalSliceOffice(s, s.ExpAdjTotalPtdEducatorConversionOnWithdrawalAnnualRiskPremium)
		s.FinalPtdEducatorConversionOnRetirementOfficePremium = finalSliceOffice(s, s.ExpAdjTotalPtdEducatorConversionOnRetirementAnnualRiskPremium)
		s.FinalPhiConversionOnWithdrawalOfficePremium = finalSliceOffice(s, s.ExpAdjTotalPhiConversionOnWithdrawalAnnualRiskPremium)
		s.FinalTtdConversionOnWithdrawalOfficePremium = finalSliceOffice(s, s.ExpAdjTotalTtdConversionOnWithdrawalAnnualRiskPremium)
		s.FinalCiConversionOnWithdrawalOfficePremium = finalSliceOffice(s, s.ExpAdjTotalCiConversionOnWithdrawalAnnualRiskPremium)
		s.FinalSglaConversionOnWithdrawalOfficePremium = finalSliceOffice(s, s.ExpAdjTotalSglaConversionOnWithdrawalAnnualRiskPremium)
		s.FinalFunConversionOnWithdrawalOfficePremium = finalSliceOffice(s, s.ExpAdjTotalFunConversionOnWithdrawalAnnualRiskPremium)
	}

	for i := range summaries {
		if err := DB.Save(&summaries[i]).Error; err != nil {
			return fmt.Errorf("save summary %d: %w", summaries[i].ID, err)
		}
	}

	if err := refreshMemberPremiumSchedules(quoteID); err != nil {
		return fmt.Errorf("refresh member premium schedules: %w", err)
	}
	return nil
}

// refreshMemberPremiumSchedules updates each MemberPremiumSchedule row for the
// quote so its per-benefit annual premiums and TotalAnnualPremiumPayable
// reflect the post-commission, post-discount Final*OfficePremium totals on
// MemberRatingResultSummary. Each member's pre-update premium share within a
// (category, benefit) bucket is used as the apportioning weight; the last
// member in each bucket absorbs rounding residual so per-category sums match
// the summary exactly. Funeral is consolidated on the summary as
// FinalFunAnnualOfficePremium, so we apportion the per-member funeral total
// first and then split each member's new total back across the four funeral
// sub-types in the same ratio they had before the refresh.
//
// Called at the tail of recomputeFinalPremiumsAndCommission, which is the
// single sink for both ApplyDiscountToQuote and applySchemeWideCommission.
func refreshMemberPremiumSchedules(quoteID int) error {
	var schedules []models.MemberPremiumSchedule
	if err := DB.Where("quote_id = ?", quoteID).
		Order("category ASC, member_name ASC").
		Find(&schedules).Error; err != nil {
		return fmt.Errorf("load member premium schedules: %w", err)
	}
	if len(schedules) == 0 {
		return nil
	}

	var summaries []models.MemberRatingResultSummary
	if err := DB.Where("quote_id = ?", quoteID).Find(&summaries).Error; err != nil {
		return fmt.Errorf("load summaries for schedule refresh: %w", err)
	}
	summaryByCategory := make(map[string]*models.MemberRatingResultSummary, len(summaries))
	for i := range summaries {
		summaryByCategory[summaries[i].Category] = &summaries[i]
	}

	categoryRows := make(map[string][]int)
	categoriesInOrder := make([]string, 0)
	for i := range schedules {
		cat := schedules[i].Category
		if _, ok := categoryRows[cat]; !ok {
			categoriesInOrder = append(categoriesInOrder, cat)
		}
		categoryRows[cat] = append(categoryRows[cat], i)
	}

	for _, cat := range categoriesInOrder {
		idxs := categoryRows[cat]
		sum, ok := summaryByCategory[cat]
		if !ok {
			continue
		}

		gla := make([]float64, len(idxs))
		ptd := make([]float64, len(idxs))
		ci := make([]float64, len(idxs))
		ttd := make([]float64, len(idxs))
		phi := make([]float64, len(idxs))
		sgla := make([]float64, len(idxs))
		funMember := make([]float64, len(idxs))
		funSplit := make([][4]float64, len(idxs))
		for k, idx := range idxs {
			r := &schedules[idx]
			gla[k] = r.GlaAnnualPremium
			ptd[k] = r.PtdAnnualPremium
			ci[k] = r.CiAnnualPremium
			ttd[k] = r.TtdAnnualPremium
			phi[k] = r.PhiAnnualPremium
			sgla[k] = r.SpouseGlaAnnualPremium
			funSplit[k] = [4]float64{
				r.MainMemberFuneralAnnualPremium,
				r.SpouseFuneralAnnualPremium,
				r.ChildrenFuneralAnnualPremium,
				r.DependantsFuneralAnnualPremium,
			}
			funMember[k] = funSplit[k][0] + funSplit[k][1] + funSplit[k][2] + funSplit[k][3]
		}

		newGla := apportionByWeight(gla, sum.FinalGlaAnnualOfficePremium)
		newPtd := apportionByWeight(ptd, sum.FinalPtdAnnualOfficePremium)
		newCi := apportionByWeight(ci, sum.FinalCiAnnualOfficePremium)
		newTtd := apportionByWeight(ttd, sum.FinalTtdAnnualOfficePremium)
		newPhi := apportionByWeight(phi, sum.FinalPhiAnnualOfficePremium)
		newSgla := apportionByWeight(sgla, sum.FinalSglaAnnualOfficePremium)
		newFunMember := apportionByWeight(funMember, sum.FinalFunAnnualOfficePremium)

		for k, idx := range idxs {
			r := &schedules[idx]
			r.GlaAnnualPremium = newGla[k]
			r.PtdAnnualPremium = newPtd[k]
			r.CiAnnualPremium = newCi[k]
			r.TtdAnnualPremium = newTtd[k]
			r.PhiAnnualPremium = newPhi[k]
			r.SpouseGlaAnnualPremium = newSgla[k]

			oldFun := funMember[k]
			if oldFun > 0 {
				scale := newFunMember[k] / oldFun
				r.MainMemberFuneralAnnualPremium = funSplit[k][0] * scale
				r.SpouseFuneralAnnualPremium = funSplit[k][1] * scale
				r.ChildrenFuneralAnnualPremium = funSplit[k][2] * scale
				r.DependantsFuneralAnnualPremium = funSplit[k][3] * scale
			} else {
				// Member had no funeral cover before; if the category total is
				// non-zero (defensive — shouldn't happen) drop the apportioned
				// share onto the main-member field so the category sum still
				// matches FinalFunAnnualOfficePremium.
				r.MainMemberFuneralAnnualPremium = newFunMember[k]
				r.SpouseFuneralAnnualPremium = 0
				r.ChildrenFuneralAnnualPremium = 0
				r.DependantsFuneralAnnualPremium = 0
			}

			r.TotalAnnualPremiumPayable = r.GlaAnnualPremium + r.PtdAnnualPremium +
				r.CiAnnualPremium + r.TtdAnnualPremium + r.PhiAnnualPremium +
				r.SpouseGlaAnnualPremium + r.MainMemberFuneralAnnualPremium +
				r.SpouseFuneralAnnualPremium + r.ChildrenFuneralAnnualPremium +
				r.DependantsFuneralAnnualPremium
		}
	}

	// MemberPremiumSchedule has no struct-level primary key. A naive
	// Updates(...) keyed on (quote_id, category, member_name) silently writes
	// the same value to every matching row, which corrupts schedule totals
	// when a member has multiple rows from movements (each with a different
	// entry_date). Delete the whole quote's slice and re-insert from the
	// modified in-memory slice so every row gets its own apportioned value
	// and ancillary columns (member name, dates, sums assured) are preserved
	// verbatim. Pattern mirrors ApplyDiscountToQuote's MemberRatingResult
	// delete-then-CreateInBatches at line ~7826.
	if err := DB.Where("quote_id = ?", quoteID).
		Delete(&models.MemberPremiumSchedule{}).Error; err != nil {
		return fmt.Errorf("delete stale member premium schedules: %w", err)
	}
	if err := DB.CreateInBatches(&schedules, 100).Error; err != nil {
		return fmt.Errorf("recreate member premium schedules: %w", err)
	}

	// Sanity log: for each benefit, the per-quote sum across schedule rows
	// should equal the summary roll-up. A non-trivial gap means apportionment
	// drifted (or some other writer touched the schedule between recompute
	// and now). Warn rather than fail so we don't break ApplyDiscountToQuote.
	const eps = 0.05
	check := func(name string, got, want float64) {
		if math.Abs(got-want) > eps {
			appLog.WithFields(map[string]interface{}{
				"quote_id":      quoteID,
				"benefit":       name,
				"schedule_sum":  got,
				"summary_total": want,
				"diff":          got - want,
			}).Warn("MemberPremiumSchedule apportionment drift")
		}
	}
	var schedGla, schedPtd, schedCi, schedTtd, schedPhi, schedSgla, schedFun float64
	for i := range schedules {
		r := &schedules[i]
		schedGla += r.GlaAnnualPremium
		schedPtd += r.PtdAnnualPremium
		schedCi += r.CiAnnualPremium
		schedTtd += r.TtdAnnualPremium
		schedPhi += r.PhiAnnualPremium
		schedSgla += r.SpouseGlaAnnualPremium
		schedFun += r.MainMemberFuneralAnnualPremium + r.SpouseFuneralAnnualPremium +
			r.ChildrenFuneralAnnualPremium + r.DependantsFuneralAnnualPremium
	}
	var sumGla, sumPtd, sumCi, sumTtd, sumPhi, sumSgla, sumFun float64
	for i := range summaries {
		s := &summaries[i]
		sumGla += s.FinalGlaAnnualOfficePremium
		sumPtd += s.FinalPtdAnnualOfficePremium
		sumCi += s.FinalCiAnnualOfficePremium
		sumTtd += s.FinalTtdAnnualOfficePremium
		sumPhi += s.FinalPhiAnnualOfficePremium
		sumSgla += s.FinalSglaAnnualOfficePremium
		sumFun += s.FinalFunAnnualOfficePremium
	}
	check("gla", schedGla, sumGla)
	check("ptd", schedPtd, sumPtd)
	check("ci", schedCi, sumCi)
	check("ttd", schedTtd, sumTtd)
	check("phi", schedPhi, sumPhi)
	check("sgla", schedSgla, sumSgla)
	check("fun", schedFun, sumFun)
	return nil
}

// apportionByWeight distributes newTotal across len(weights) buckets in
// proportion to weights, with the last bucket absorbing any rounding residual
// so the returned slice sums to newTotal exactly. When all weights are zero
// the total is split evenly (defensive — should not happen in practice since
// the same pre-state weights drive the per-benefit Final* allocation upstream).
func apportionByWeight(weights []float64, newTotal float64) []float64 {
	out := make([]float64, len(weights))
	if len(weights) == 0 {
		return out
	}
	last := len(weights) - 1
	oldTotal := 0.0
	for _, w := range weights {
		oldTotal += w
	}
	if oldTotal == 0 {
		if newTotal == 0 {
			return out
		}
		even := newTotal / float64(len(weights))
		remaining := newTotal
		for i := 0; i < last; i++ {
			out[i] = even
			remaining -= even
		}
		out[last] = remaining
		return out
	}
	remaining := newTotal
	for i := 0; i < last; i++ {
		out[i] = weights[i] * (newTotal / oldTotal)
		remaining -= out[i]
	}
	out[last] = remaining
	return out
}

// applySchemeWideCommission computes the total commission for the quote by
// running the tiered CommissionStructure bands against the scheme-wide total
// premium, then distributes that commission down two levels: scheme → category
// (by total-premium share), and category → benefit (by category-level
// premium share). This guarantees that on every category
// TotalCommission == Σ per-benefit commissions, and across the scheme
// SchemeTotalCommission == Σ TotalCommission == Σ all per-benefit commissions.
//
// Flow:
//  1. Load all MemberRatingResultSummary rows for the quote.
//  2. For each category, sum Exp-Adj office premium across every benefit to
//     get categoryTotals[i]; sum those to get the scheme-wide total.
//  3. Call ComputeProgressiveCommission(channel, broker, total) to get the
//     blended commission rate (also returned so reports can show it).
//  4. overallCommission = total * rate.
//  5. Split overallCommission into per-category commissions in proportion to
//     categoryTotals[i] / schemeTotal. The last category absorbs the residual
//     so Σ category commissions == overallCommission exactly.
//  6. Within each category, split that category's commission across its
//     benefits in proportion to acc.premium(s) / categoryTotals[i]. The last
//     benefit absorbs the per-category residual so Σ benefit commissions in
//     the category == that category's commission exactly.
//
// Called once per quote after all category workers have finished.
func applySchemeWideCommission(quoteID int, groupQuote models.GroupPricingQuote, logger *logrus.Entry) error {
	var summaries []models.MemberRatingResultSummary
	if err := DB.Where("quote_id = ?", quoteID).Order("id ASC").Find(&summaries).Error; err != nil {
		return fmt.Errorf("load summaries for commission: %w", err)
	}
	if len(summaries) == 0 {
		return nil
	}

	accessors := benefitAccessors()

	// Per-benefit Exp office premium for the Exp* commission allocation. This
	// is the pre-discount value: ComputeOfficePremium uses SchemeTotalLoading
	// (no discount) so Exp* commissions remain frozen at quote-calc time.
	expOfficePremium := func(s *models.MemberRatingResultSummary, acc benefitAccessor) float64 {
		return models.ComputeOfficePremium(acc.expRiskPremium(s), s)
	}

	// Step 2: per-category total premium (sum across all benefits in the
	// TotalAnnualPremium formula) and the scheme-wide total. categoryTotals[i]
	// equals s.TotalAnnualPremium for category i once the recompute below has
	// run; computing it here directly avoids depending on stale stored values.
	categoryTotals := make([]float64, len(summaries))
	schemeTotal := 0.0
	for i := range summaries {
		s := &summaries[i]
		for _, acc := range accessors {
			categoryTotals[i] += expOfficePremium(s, acc)
		}
		schemeTotal += categoryTotals[i]
	}

	// Step 3: blended commission rate from the tiered CommissionStructure
	// bands keyed by (channel, broker name).
	channel := string(groupQuote.DistributionChannel)
	holderName := strings.TrimSpace(groupQuote.QuoteBroker.Name)
	rate := 0.0
	if schemeTotal > 0 {
		r, err := ComputeProgressiveCommission(channel, holderName, schemeTotal)
		if err != nil {
			return fmt.Errorf("compute progressive commission: %w", err)
		}
		rate = r
	}
	overallCommission := schemeTotal * rate

	// Step 4: split the scheme commission into per-category commissions by
	// each category's share of total premium. The last category absorbs any
	// floating-point residual so Σ category commissions == overallCommission
	// exactly.
	categoryCommissions := make([]float64, len(summaries))
	lastCatIdx := len(summaries) - 1
	remainingScheme := overallCommission
	for i := range summaries {
		if i == lastCatIdx {
			categoryCommissions[i] = remainingScheme
		} else if schemeTotal > 0 {
			categoryCommissions[i] = overallCommission * (categoryTotals[i] / schemeTotal)
			remainingScheme -= categoryCommissions[i]
		}
	}

	// Step 5: within each category, split the category commission across
	// benefits by each benefit's share of the category-level total premium.
	// The last benefit absorbs the per-category residual so Σ benefit
	// commissions on a category == that category's commission exactly.
	lastBenIdx := len(accessors) - 1
	for catIdx := range summaries {
		s := &summaries[catIdx]
		categoryCommission := categoryCommissions[catIdx]
		categoryTotal := categoryTotals[catIdx]
		remainingCat := categoryCommission
		for benIdx, acc := range accessors {
			var slice float64
			if benIdx == lastBenIdx {
				// Final slice takes whatever is left. When categoryTotal is 0
				// remainingCat is also 0, so slice is 0 too.
				slice = remainingCat
			} else if categoryTotal > 0 {
				slice = categoryCommission * (expOfficePremium(s, acc) / categoryTotal)
				remainingCat -= slice
			}
			acc.setExpComm(s, slice)
		}
		s.TotalCommission = categoryCommission
		s.SchemeTotalCommission = overallCommission
		s.SchemeTotalCommissionRate = rate

		// Post-distribution: refresh the roll-up totals that downstream
		// reports read. Rate-per-1000 and proportion-of-salary fields
		// that were derived during per-category calc are also stale now —
		// recompute from the updated benefit totals.
		if s.TotalGlaCappedSumAssured > 0 {
		}
		if s.TotalPtdCappedSumAssured > 0 {
		}
		if s.TotalCiCappedSumAssured > 0 {
		}
		if s.TotalSglaCappedSumAssured > 0 {
		}
		// TTD / PHI use income bases rather than sum assured.
		s.ExpTotalAnnualPremiumExclFuneral = models.ComputeOfficePremium(s.ExpTotalGlaAnnualRiskPremium, s) +
			models.ComputeOfficePremium(s.ExpTotalAdditionalAccidentalGlaAnnualRiskPremium, s) +
			models.ComputeOfficePremium(s.ExpTotalPtdAnnualRiskPremium, s) +
			models.ComputeOfficePremium(s.ExpTotalTtdAnnualRiskPremium, s) +
			models.ComputeOfficePremium(s.ExpTotalPhiAnnualRiskPremium, s) +
			models.ComputeOfficePremium(s.ExpTotalCiAnnualRiskPremium, s) +
			models.ComputeOfficePremium(s.ExpTotalSglaAnnualRiskPremium, s) +
			models.ComputeOfficePremium(s.ExpTotalTaxSaverAnnualRiskPremium, s) +
			models.ComputeOfficePremium(s.ExpAdjTotalGlaEducatorRiskPremium, s) +
			models.ComputeOfficePremium(s.ExpAdjTotalPtdEducatorRiskPremium, s)
		s.TotalAnnualPremium = s.ExpTotalAnnualPremiumExclFuneral + models.ComputeOfficePremium(s.ExpTotalFunAnnualRiskPremium, s)
	}

	// Theoretical (book-rate) commission allocation. Mirrors the Exp* allocation
	// above but uses bookRiskPremium (TotalAnnualRiskPremium / pre-credibility
	// book rates) so the OutputSummary Theoretical column has a commission slice
	// that reconciles to a book-rate progressive band rather than the experience-
	// blended one. Independent from the Exp* allocation: a scheme whose book
	// premium falls in a different commission band than its experience-blended
	// premium will get a different theoretical rate.
	bookOfficePremium := func(s *models.MemberRatingResultSummary, acc benefitAccessor) float64 {
		return models.ComputeOfficePremium(acc.bookRiskPremium(s), s)
	}
	bookCategoryTotals := make([]float64, len(summaries))
	bookSchemeTotal := 0.0
	for i := range summaries {
		s := &summaries[i]
		for _, acc := range accessors {
			bookCategoryTotals[i] += bookOfficePremium(s, acc)
		}
		bookSchemeTotal += bookCategoryTotals[i]
	}
	bookRate := 0.0
	if bookSchemeTotal > 0 {
		r, err := ComputeProgressiveCommission(channel, holderName, bookSchemeTotal)
		if err != nil {
			return fmt.Errorf("compute book progressive commission: %w", err)
		}
		bookRate = r
	}
	bookOverallCommission := bookSchemeTotal * bookRate
	bookCategoryCommissions := make([]float64, len(summaries))
	remainingBookScheme := bookOverallCommission
	for i := range summaries {
		if i == lastCatIdx {
			bookCategoryCommissions[i] = remainingBookScheme
		} else if bookSchemeTotal > 0 {
			bookCategoryCommissions[i] = bookOverallCommission * (bookCategoryTotals[i] / bookSchemeTotal)
			remainingBookScheme -= bookCategoryCommissions[i]
		}
	}
	for catIdx := range summaries {
		s := &summaries[catIdx]
		bookCategoryCommission := bookCategoryCommissions[catIdx]
		bookCategoryTotal := bookCategoryTotals[catIdx]
		remainingCat := bookCategoryCommission
		for benIdx, acc := range accessors {
			var slice float64
			if benIdx == lastBenIdx {
				slice = remainingCat
			} else if bookCategoryTotal > 0 {
				slice = bookCategoryCommission * (bookOfficePremium(s, acc) / bookCategoryTotal)
				remainingCat -= slice
			}
			acc.setBookComm(s, slice)
		}
	}

	for i := range summaries {
		if err := DB.Save(&summaries[i]).Error; err != nil {
			return fmt.Errorf("save summary %d: %w", summaries[i].ID, err)
		}
	}

	// Mirror the Exp* allocation onto the Final* fields. With Discount == 0
	// (the case at calc time), Final*OfficePremium == Exp*OfficePremium and
	// Final*CommissionAmount == Exp*CommissionAmount. ApplyDiscountToQuote
	// re-runs this same helper after updating the discount on each summary.
	if err := recomputeFinalPremiumsAndCommission(quoteID, groupQuote); err != nil {
		return err
	}

	if logger != nil {
		logger.WithFields(map[string]interface{}{
			"scheme_total_premium": schemeTotal,
			"commission_rate":      rate,
			"overall_commission":   overallCommission,
			"category_count":       len(summaries),
		}).Info("Applied scheme-wide commission")
	}
	return nil
}

// binderAndOutsourceRates returns the binder-fee and outsource-fee fractions
// (e.g. 5% → 0.05) to apply on top of the regular loadings when a quote is
// sold through the binder distribution channel. Returns zeros for any other
// channel so non-binder quotes are unaffected.
func binderAndOutsourceRates(quote *models.GroupPricingQuote) (float64, float64) {
	if quote == nil || quote.DistributionChannel != models.ChannelBinder {
		return 0, 0
	}
	return quote.Loadings.BinderFee / 100.0, quote.Loadings.OutsourceFee / 100.0
}

// applyBinderOutsourceAmounts decomposes each benefit's office premium into
// its binder-fee and outsource-fee slices and writes them into the matching
// fields on the MemberRatingResult. Office premium is derived from the
// per-member risk premium and the scheme-level loading on the quote.
// Totals across benefits are rolled up into TotalFuneral*Amount (family
// funeral) and the all-benefit TotalBinderAmount / TotalOutsourcedAmount.
func applyBinderOutsourceAmounts(r *models.MemberRatingResult, quote *models.GroupPricingQuote, binderRate, outsourceRate float64) {
	op := func(risk float64) float64 { return computeMemberOfficePremium(risk, quote) }

	r.GlaBinderAmount = op(r.GlaRiskPremium) * binderRate
	r.GlaOutsourcedAmount = op(r.GlaRiskPremium) * outsourceRate
	r.ExpAdjGlaBinderAmount = op(r.ExpAdjGlaRiskPremium) * binderRate
	r.ExpAdjGlaOutsourcedAmount = op(r.ExpAdjGlaRiskPremium) * outsourceRate

	r.AdditionalAccidentalGlaBinderAmount = op(r.AdditionalAccidentalGlaRiskPremium) * binderRate
	r.AdditionalAccidentalGlaOutsourcedAmount = op(r.AdditionalAccidentalGlaRiskPremium) * outsourceRate
	r.ExpAdjAdditionalAccidentalGlaBinderAmount = op(r.ExpAdjAdditionalAccidentalGlaRiskPremium) * binderRate
	r.ExpAdjAdditionalAccidentalGlaOutsourcedAmt = op(r.ExpAdjAdditionalAccidentalGlaRiskPremium) * outsourceRate

	r.PtdBinderAmount = op(r.PtdRiskPremium) * binderRate
	r.PtdOutsourcedAmount = op(r.PtdRiskPremium) * outsourceRate
	r.ExpAdjPtdBinderAmount = op(r.ExpAdjPtdRiskPremium) * binderRate
	r.ExpAdjPtdOutsourcedAmount = op(r.ExpAdjPtdRiskPremium) * outsourceRate

	r.CiBinderAmount = op(r.CiRiskPremium) * binderRate
	r.CiOutsourcedAmount = op(r.CiRiskPremium) * outsourceRate
	r.ExpAdjCiBinderAmount = op(r.ExpAdjCiRiskPremium) * binderRate
	r.ExpAdjCiOutsourcedAmount = op(r.ExpAdjCiRiskPremium) * outsourceRate

	r.SpouseGlaBinderAmount = op(r.SpouseGlaRiskPremium) * binderRate
	r.SpouseGlaOutsourcedAmount = op(r.SpouseGlaRiskPremium) * outsourceRate
	r.ExpAdjSpouseGlaBinderAmount = op(r.ExpAdjSpouseGlaRiskPremium) * binderRate
	r.ExpAdjSpouseGlaOutsourcedAmount = op(r.ExpAdjSpouseGlaRiskPremium) * outsourceRate

	r.TtdBinderAmount = op(r.TtdRiskPremium) * binderRate
	r.TtdOutsourcedAmount = op(r.TtdRiskPremium) * outsourceRate
	r.ExpAdjTtdBinderAmount = op(r.ExpAdjTtdRiskPremium) * binderRate
	r.ExpAdjTtdOutsourcedAmount = op(r.ExpAdjTtdRiskPremium) * outsourceRate

	r.PhiBinderAmount = op(r.PhiRiskPremium) * binderRate
	r.PhiOutsourcedAmount = op(r.PhiRiskPremium) * outsourceRate
	r.ExpAdjPhiBinderAmount = op(r.ExpAdjPhiRiskPremium) * binderRate
	r.ExpAdjPhiOutsourcedAmount = op(r.ExpAdjPhiRiskPremium) * outsourceRate

	r.MainMemberFuneralBinderAmount = op(r.MainMemberFuneralRiskPremium) * binderRate
	r.MainMemberFuneralOutsourcedAmount = op(r.MainMemberFuneralRiskPremium) * outsourceRate
	r.SpouseFuneralBinderAmount = op(r.SpouseFuneralRiskPremium) * binderRate
	r.SpouseFuneralOutsourcedAmount = op(r.SpouseFuneralRiskPremium) * outsourceRate
	r.ChildrenFuneralBinderAmount = op(r.ChildFuneralRiskPremium) * binderRate
	r.ChildrenFuneralOutsourcedAmount = op(r.ChildFuneralRiskPremium) * outsourceRate
	r.DependantsFuneralBinderAmount = op(r.ParentFuneralRiskPremium) * binderRate
	r.DependantsFuneralOutsourcedAmount = op(r.ParentFuneralRiskPremium) * outsourceRate
	r.TotalFuneralBinderAmount = r.TotalFuneralOfficePremium * binderRate
	r.TotalFuneralOutsourcedAmount = r.TotalFuneralOfficePremium * outsourceRate
	r.ExpAdjTotalFuneralBinderAmount = r.FinalTotalFuneralOfficePremium * binderRate
	r.ExpAdjTotalFuneralOutsourcedAmount = r.FinalTotalFuneralOfficePremium * outsourceRate

	r.GlaEducatorBinderAmount = op(r.GlaEducatorRiskPremium) * binderRate
	r.GlaEducatorOutsourcedAmount = op(r.GlaEducatorRiskPremium) * outsourceRate
	r.ExpAdjGlaEducatorBinderAmount = op(r.ExpAdjGlaEducatorRiskPremium) * binderRate
	r.ExpAdjGlaEducatorOutsourcedAmount = op(r.ExpAdjGlaEducatorRiskPremium) * outsourceRate
	r.PtdEducatorBinderAmount = op(r.PtdEducatorRiskPremium) * binderRate
	r.PtdEducatorOutsourcedAmount = op(r.PtdEducatorRiskPremium) * outsourceRate
	r.ExpAdjPtdEducatorBinderAmount = op(r.ExpAdjPtdEducatorRiskPremium) * binderRate
	r.ExpAdjPtdEducatorOutsourcedAmount = op(r.ExpAdjPtdEducatorRiskPremium) * outsourceRate

	r.TotalBinderAmount = r.GlaBinderAmount + r.AdditionalAccidentalGlaBinderAmount +
		r.PtdBinderAmount + r.CiBinderAmount + r.SpouseGlaBinderAmount +
		r.TtdBinderAmount + r.PhiBinderAmount + r.TotalFuneralBinderAmount +
		r.GlaEducatorBinderAmount + r.PtdEducatorBinderAmount
	r.TotalOutsourcedAmount = r.GlaOutsourcedAmount + r.AdditionalAccidentalGlaOutsourcedAmount +
		r.PtdOutsourcedAmount + r.CiOutsourcedAmount + r.SpouseGlaOutsourcedAmount +
		r.TtdOutsourcedAmount + r.PhiOutsourcedAmount + r.TotalFuneralOutsourcedAmount +
		r.GlaEducatorOutsourcedAmount + r.PtdEducatorOutsourcedAmount
}

// GetSchemeSizeLoading returns the SchemeSizeLevel row matching the quote's
// member count band for the given risk_rate_code. Caches the full row so the
// per-benefit loadings (gla_loading, ptd_loading, …) and the size_level can be
// read with a single DB hit. Returns a zero-value struct when no band matches —
// callers can read SizeLevel and the *Loading fields safely (they'll be 0).
func GetSchemeSizeLoading(groupPricingParameter models.GroupPricingParameters, memberCount int) models.SchemeSizeLevel {
	tableName := "scheme_size_levels"

	var keyString strings.Builder
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(memberCount) + "_")
	key := keyString.String()

	cacheKey := tableName + "_row_" + key
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(models.SchemeSizeLevel)
	}

	var schemeSizeLevel models.SchemeSizeLevel
	err := DB.Table(tableName).
		Where("risk_rate_code = ? AND min_count <= ? AND max_count >= ?",
			groupPricingParameter.RiskRateCode,
			memberCount,
			memberCount,
		).
		First(&schemeSizeLevel).Error
	if err != nil {
		fmt.Println(err)
		GroupPricingCache.Set(cacheKey, models.SchemeSizeLevel{}, 1)
		return models.SchemeSizeLevel{}
	}

	GroupPricingCache.Set(cacheKey, schemeSizeLevel, 1)
	return schemeSizeLevel
}

func GetSchemeSizeLevel(groupPricingParameter models.GroupPricingParameters, memberCount int) int {
	return GetSchemeSizeLoading(groupPricingParameter, memberCount).SizeLevel
}

// GetDiscountAuthorityForUser returns the DiscountAuthority row for the user's GP role
// and the given risk_rate_code. The MaxDiscount is stored as a percentage (e.g. 5 = 5%).
func GetDiscountAuthorityForUser(userEmail string, riskRateCode string) (models.DiscountAuthority, error) {
	var orgUser models.OrgUser
	if err := DB.Where("email = ?", userEmail).First(&orgUser).Error; err != nil {
		return models.DiscountAuthority{}, err
	}
	var da models.DiscountAuthority
	if err := DB.Where("role = ? AND risk_rate_code = ?", orgUser.GPRole, riskRateCode).First(&da).Error; err != nil {
		return models.DiscountAuthority{}, err
	}
	return da, nil
}

// GetDiscountMethod returns the globally configured discount calculation
// method from the GroupPricingSetting singleton row. Defaults to
// loading_adjustment if the row is missing or the column is empty so that
// quotes computed before the setting was introduced behave identically to
// historical output.
func GetDiscountMethod() string {
	var s models.GroupPricingSetting
	if err := DB.First(&s, 1).Error; err != nil {
		return models.DiscountMethodLoadingAdjustment
	}
	if s.DiscountMethod == "" {
		return models.DiscountMethodLoadingAdjustment
	}
	return s.DiscountMethod
}

// GetFCLMethod returns the globally configured Free Cover Limit calculation
// method from the GroupPricingSetting singleton row. Defaults to the
// percentile method (existing behaviour) when the row is missing or the
// column is empty so quotes computed before the setting was introduced
// behave identically to historical output.
func GetFCLMethod() string {
	var s models.GroupPricingSetting
	if err := DB.First(&s, 1).Error; err != nil {
		return models.FCLMethodPercentile
	}
	if s.FCLMethod == "" {
		return models.FCLMethodPercentile
	}
	return s.FCLMethod
}

// GetMedicalAidWaiverMethod returns the globally configured PHI medical aid
// waiver methodology from the GroupPricingSetting singleton row. Defaults to
// the formula method (salary * proportion + amount) when the row is missing or
// the column is empty so quotes computed before the setting was introduced
// behave identically to historical output.
func GetMedicalAidWaiverMethod() string {
	var s models.GroupPricingSetting
	if err := DB.First(&s, 1).Error; err != nil {
		return models.MedicalAidWaiverMethodFormula
	}
	if s.MedicalAidWaiverMethod == "" {
		return models.MedicalAidWaiverMethodFormula
	}
	return s.MedicalAidWaiverMethod
}

// FCLOverrideToleranceDefault is the headroom (as a fraction) allowed above
// Restriction.MaximumAllowedFCL before a quote-level FCL override is clamped.
// 0.2 means a 20% allowance.
const FCLOverrideToleranceDefault = 0.2

// RiskProfileVariationToleranceDefault is the percentage variation in member
// data profile that an insurer is willing to accept between quotation and
// implementation before reserving the right to revise rates. Used as the
// fallback for the Acceptance Form text in the quote PDF/DOCX.
const RiskProfileVariationToleranceDefault = 7.0

// GetRiskProfileVariationTolerancePct returns the configured profile-variation
// tolerance (as a percentage, e.g. 7.0 for 7%) from the GroupPricingSetting
// singleton row. Falls back to RiskProfileVariationToleranceDefault when the
// row is missing or the stored value is non-positive.
func GetRiskProfileVariationTolerancePct() float64 {
	var s models.GroupPricingSetting
	if err := DB.First(&s, 1).Error; err != nil {
		return RiskProfileVariationToleranceDefault
	}
	if s.RiskProfileVariationTolerancePct <= 0 {
		return RiskProfileVariationToleranceDefault
	}
	return s.RiskProfileVariationTolerancePct
}

// GetFCLOverrideTolerance returns the configured headroom above
// Restriction.MaximumAllowedFCL that a quote-level FCL override is allowed
// to claim before being clamped. Falls back to FCLOverrideToleranceDefault
// when the singleton row is missing or the value is non-positive.
func GetFCLOverrideTolerance() float64 {
	var s models.GroupPricingSetting
	if err := DB.First(&s, 1).Error; err != nil {
		return FCLOverrideToleranceDefault
	}
	if s.FCLOverrideTolerance <= 0 {
		return FCLOverrideToleranceDefault
	}
	return s.FCLOverrideTolerance
}

// ApplyDiscountToQuote applies a discount (in percentage, e.g. 5.0 for 5%) to a
// quote: persists the percentage on the quote loadings, mirrors it onto every
// MemberRatingResultSummary as a negative fraction so SchemeTotalLoading +
// Discount drives Final*Premium and Final*Commission via the shared helper,
// and recomputes the per-member funeral office premium + binder/outsource
// amounts using the same Risk / (1 - (TotalPremiumLoading + Discount)) formula.
func ApplyDiscountToQuote(quoteId string, discountPct float64, user models.AppUser) error {
	var quote models.GroupPricingQuote
	if err := DB.Where("id = ?", quoteId).First(&quote).Error; err != nil {
		return err
	}

	var groupParam models.GroupPricingParameters
	DB.Where("basis = ? AND risk_rate_code = ?", quote.Basis, quote.RiskRateCode).First(&groupParam)

	// Persist the discount percentage on the quote so it is available during (re)calculation
	quote.Loadings.Discount = discountPct
	DB.Save(&quote)

	// Store discount as a negative fraction so it reduces TotalLoading
	discount := -(discountPct / 100.0)

	// Mirror the discount onto every summary so SchemeTotalLoading() picks it
	// up via FinalSchemeTotalLoading(). The recomputeFinalPremiumsAndCommission
	// call below reads s.Discount when deriving every Final*Premium and
	// re-allocating Final*Commission against the new scheme total.
	if err := DB.Model(&models.MemberRatingResultSummary{}).
		Where("quote_id = ?", quoteId).
		Update("discount", discount).Error; err != nil {
		return err
	}

	var results []models.MemberRatingResult
	if err := DB.Where("quote_id = ?", quoteId).Find(&results).Error; err != nil {
		return err
	}

	schemeSizeLevel := GetSchemeSizeLevel(groupParam, quote.MemberDataCount)
	premiumLoading := GetPremiumLoading(groupParam, schemeSizeLevel, string(quote.DistributionChannel))

	binderFeeRate, outsourceFeeRate := binderAndOutsourceRates(&quote)
	discountMethod := GetDiscountMethod()

	for i := range results {
		r := &results[i]
		r.Discount = discount
		// Commission is excluded from TotalLoading — see applySchemeWideCommission.
		r.TotalPremiumLoading = math.Max(
			r.ExpenseLoading+r.AdminLoading+r.ProfitLoading+r.OtherLoading+binderFeeRate+outsourceFeeRate,
			premiumLoading.MinimumPremiumLoading,
		)
		divisor := 1.0 - r.TotalPremiumLoading
		if divisor == 0 {
			divisor = 1.0
		}
		r.TotalFuneralOfficePremium = r.TotalFuneralRiskPremium / divisor
		r.ExpAdjTotalFuneralOfficePremium = r.ExpAdjTotalFuneralRiskPremium / divisor
		// Final office premium under either discount method:
		//  - loading_adjustment: Risk / (1 - (TotalPremiumLoading + discount)).
		//    discount is the negative fraction so the denominator grows and
		//    Final shrinks below ExpAdj.
		//  - prorata: ExpAdj * (1 + discount) — i.e. the un-discounted office
		//    premium multiplied by (1 - d), proportionally reducing risk +
		//    every loading component by the same fraction.
		var finalFunOffice float64
		if discountMethod == models.DiscountMethodProrata {
			finalFunOffice = r.ExpAdjTotalFuneralOfficePremium * (1.0 + discount)
		} else {
			finalDivisor := 1.0 - (r.TotalPremiumLoading + discount)
			if finalDivisor == 0 {
				finalDivisor = 1.0
			}
			finalFunOffice = r.ExpAdjTotalFuneralRiskPremium / finalDivisor
		}
		r.FinalTotalFuneralOfficePremium = finalFunOffice
		applyBinderOutsourceAmounts(r, &quote, binderFeeRate, outsourceFeeRate)
	}

	// Delete existing member rows and recreate with updated values.
	if err := DB.Where("quote_id = ?", quoteId).Delete(&models.MemberRatingResult{}).Error; err != nil {
		return err
	}
	if len(results) > 0 {
		if err := DB.CreateInBatches(&results, 100).Error; err != nil {
			return err
		}
	}

	// Refresh Final*Premium and Final*Commission on every summary using the
	// new discount. Pre-discount this was a no-op (Final == Exp); now the
	// Final values diverge while Exp* commissions stay frozen.
	return recomputeFinalPremiumsAndCommission(intQuoteID(quoteId), quote)
}

func GetFuneralRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters) float64 {
	tableName := "funeral_rates"

	var keyString strings.Builder
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.AgeNextBirthday) + "_")
	keyString.WriteString(memberResultData.Gender[:1] + "_")
	key := keyString.String()

	cacheKey := tableName + "_" + key
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}

	var qx float64
	err := DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND gender = ?",
			groupPricingParameter.RiskRateCode,
			memberResultData.AgeNextBirthday,
			memberResultData.Gender[:1],
		).
		Pluck("fun_qx", &qx).Error
	if err != nil {
		fmt.Println(err)
	}

	GroupPricingCache.Set(cacheKey, qx, 1)
	return qx
}

func GetFuneralAidsRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters) float64 {
	tableName := "funeral_aids_rates"

	var keyString strings.Builder
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.AgeNextBirthday) + "_")
	keyString.WriteString(memberResultData.Gender[:1] + "_")
	key := keyString.String()

	cacheKey := tableName + "_" + key
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}

	var qx float64
	err := DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND gender = ?",
			groupPricingParameter.RiskRateCode,
			memberResultData.AgeNextBirthday,
			memberResultData.Gender[:1],
		).
		Pluck("fun_aids_qx", &qx).Error
	if err != nil {
		fmt.Println(err)
	}

	GroupPricingCache.Set(cacheKey, qx, 1)
	return qx
}

func GetSpouseFuneralRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters) float64 {
	tableName := "funeral_rates"

	var keyString strings.Builder
	keyString.WriteString("spouse_")
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.SpouseAgeNextBirthday) + "_")
	keyString.WriteString(memberResultData.SpouseGender[:1] + "_")
	key := keyString.String()

	cacheKey := tableName + "_" + key
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}

	var qx float64
	err := DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND gender = ?",
			groupPricingParameter.RiskRateCode,
			memberResultData.SpouseAgeNextBirthday,
			memberResultData.SpouseGender[:1],
		).
		Pluck("fun_qx", &qx).Error
	if err != nil {
		fmt.Println(err)
	}

	GroupPricingCache.Set(cacheKey, qx, 1)
	return qx
}

func GetSpouseFuneralAidsRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters) float64 {
	tableName := "funeral_aids_rates"

	var keyString strings.Builder
	keyString.WriteString("spouse_")
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.SpouseAgeNextBirthday) + "_")
	keyString.WriteString(memberResultData.SpouseGender[:1] + "_")
	key := keyString.String()

	cacheKey := tableName + "_" + key
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}

	var qx float64
	err := DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND gender = ?",
			groupPricingParameter.RiskRateCode,
			memberResultData.SpouseAgeNextBirthday,
			memberResultData.SpouseGender[:1],
		).
		Pluck("fun_aids_qx", &qx).Error
	if err != nil {
		fmt.Println(err)
	}

	GroupPricingCache.Set(cacheKey, qx, 1)
	return qx
}

// GetReinsuranceGlaAidsRate looks up the reinsurance-specific GLA aids Qx
// for the main member. Returns 0 if no row is found (so missing tables do
// not break pricing for schemes without reinsurance).
func GetReinsuranceGlaAidsRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters) float64 {
	tableName := "reinsurance_gla_aids_rates"
	cacheKey := tableName + "_" + groupPricingParameter.RiskRateCode + "_" + strconv.Itoa(memberResultData.AgeNextBirthday) + "_" + memberResultData.Gender[:1]
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}
	var qx float64
	DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND gender = ?",
			groupPricingParameter.RiskRateCode, memberResultData.AgeNextBirthday, memberResultData.Gender[:1]).
		Pluck("gla_aids_qx", &qx)
	GroupPricingCache.Set(cacheKey, qx, 1)
	return qx
}

// GetReinsuranceSpouseGlaAidsRate — spouse variant of GetReinsuranceGlaAidsRate.
func GetReinsuranceSpouseGlaAidsRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters) float64 {
	tableName := "reinsurance_gla_aids_rates"
	cacheKey := tableName + "_spouse_" + groupPricingParameter.RiskRateCode + "_" + strconv.Itoa(memberResultData.SpouseAgeNextBirthday) + "_" + memberResultData.SpouseGender[:1]
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}
	var qx float64
	DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND gender = ?",
			groupPricingParameter.RiskRateCode, memberResultData.SpouseAgeNextBirthday, memberResultData.SpouseGender[:1]).
		Pluck("gla_aids_qx", &qx)
	GroupPricingCache.Set(cacheKey, qx, 1)
	return qx
}

// GetReinsuranceFuneralRate looks up the reinsurance-specific funeral Qx for the main member.
func GetReinsuranceFuneralRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters) float64 {
	tableName := "reinsurance_funeral_rates"
	cacheKey := tableName + "_" + groupPricingParameter.RiskRateCode + "_" + strconv.Itoa(memberResultData.AgeNextBirthday) + "_" + memberResultData.Gender[:1]
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}
	var qx float64
	DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND gender = ?",
			groupPricingParameter.RiskRateCode, memberResultData.AgeNextBirthday, memberResultData.Gender[:1]).
		Pluck("fun_qx", &qx)
	GroupPricingCache.Set(cacheKey, qx, 1)
	return qx
}

// GetReinsuranceSpouseFuneralRate — spouse variant.
func GetReinsuranceSpouseFuneralRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters) float64 {
	tableName := "reinsurance_funeral_rates"
	cacheKey := tableName + "_spouse_" + groupPricingParameter.RiskRateCode + "_" + strconv.Itoa(memberResultData.SpouseAgeNextBirthday) + "_" + memberResultData.SpouseGender[:1]
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}
	var qx float64
	DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND gender = ?",
			groupPricingParameter.RiskRateCode, memberResultData.SpouseAgeNextBirthday, memberResultData.SpouseGender[:1]).
		Pluck("fun_qx", &qx)
	GroupPricingCache.Set(cacheKey, qx, 1)
	return qx
}

// GetReinsuranceFuneralAidsRate looks up the reinsurance-specific funeral aids Qx for the main member.
func GetReinsuranceFuneralAidsRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters) float64 {
	tableName := "reinsurance_funeral_aids_rates"
	cacheKey := tableName + "_" + groupPricingParameter.RiskRateCode + "_" + strconv.Itoa(memberResultData.AgeNextBirthday) + "_" + memberResultData.Gender[:1]
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}
	var qx float64
	DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND gender = ?",
			groupPricingParameter.RiskRateCode, memberResultData.AgeNextBirthday, memberResultData.Gender[:1]).
		Pluck("fun_aids_qx", &qx)
	GroupPricingCache.Set(cacheKey, qx, 1)
	return qx
}

// GetReinsuranceSpouseFuneralAidsRate — spouse variant.
func GetReinsuranceSpouseFuneralAidsRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters) float64 {
	tableName := "reinsurance_funeral_aids_rates"
	cacheKey := tableName + "_spouse_" + groupPricingParameter.RiskRateCode + "_" + strconv.Itoa(memberResultData.SpouseAgeNextBirthday) + "_" + memberResultData.SpouseGender[:1]
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}
	var qx float64
	DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND gender = ?",
			groupPricingParameter.RiskRateCode, memberResultData.SpouseAgeNextBirthday, memberResultData.SpouseGender[:1]).
		Pluck("fun_aids_qx", &qx)
	GroupPricingCache.Set(cacheKey, qx, 1)
	return qx
}

// GetReinsuranceGeneralLoading returns the per-age/gender reinsurance contingency &
// continuation loadings row. Returns a zero-value struct (all loadings 0) if no
// row matches, so reinsurance pricing gracefully degrades when the table is empty.
func GetReinsuranceGeneralLoading(riskRateCode string, age int, gender string) models.ReinsuranceGeneralLoading {
	if !IsTableRequired("reinsuranceGeneralLoading") {
		return models.ReinsuranceGeneralLoading{}
	}
	tableName := "reinsurance_general_loadings"
	cacheKey := tableName + "_" + riskRateCode + "_" + strconv.Itoa(age) + "_" + gender
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(models.ReinsuranceGeneralLoading)
	}
	var loading models.ReinsuranceGeneralLoading
	DB.Table(tableName).
		Where("risk_rate_code = ? AND age = ? AND gender = ?", riskRateCode, age, gender).
		First(&loading)
	GroupPricingCache.Set(cacheKey, loading, 1)
	return loading
}

// GetReinsuranceIndustryLoading returns one row from reinsurance_industry_loadings
// keyed by (risk_rate_code, occupation_class, gender). Returns a zero-value
// struct if no row matches, so reinsurance pricing gracefully degrades when
// the table is empty.
//
// The function fetches all rows for (risk, occupation_class) in a single
// cached query and picks the one whose first-letter-upper gender matches,
// mirroring the normalization done by the direct-pricing industry-loading
// map built at api/services/group_pricing.go:1157 — the DB can store
// "Male"/"Female"/"M"/"F"/"m" and the lookup still finds the right row.
func GetReinsuranceIndustryLoading(riskRateCode string, occupationClass int, gender string) models.ReinsuranceIndustryLoading {
	if !IsTableRequired("reinsuranceIndustryLoading") {
		return models.ReinsuranceIndustryLoading{}
	}
	tableName := "reinsurance_industry_loadings"
	mapKey := tableName + "_bygender_" + riskRateCode + "_" + strconv.Itoa(occupationClass)
	var byGender map[string]models.ReinsuranceIndustryLoading
	if cached, found := GroupPricingCache.Get(mapKey); found {
		byGender = cached.(map[string]models.ReinsuranceIndustryLoading)
	} else {
		byGender = make(map[string]models.ReinsuranceIndustryLoading)
		var rows []models.ReinsuranceIndustryLoading
		DB.Table(tableName).
			Where("risk_rate_code = ? AND occupation_class = ?", riskRateCode, occupationClass).
			Find(&rows)
		for _, r := range rows {
			if len(r.Gender) > 0 {
				byGender[strings.ToUpper(r.Gender[:1])] = r
			}
		}
		GroupPricingCache.Set(mapKey, byGender, 1)
	}
	if len(gender) == 0 {
		return models.ReinsuranceIndustryLoading{}
	}
	return byGender[strings.ToUpper(gender[:1])]
}

// GetReinsuranceRegionLoading returns one row from reinsurance_region_loadings
// keyed by (risk_rate_code, gender, region). Returns a zero-value struct when
// no row matches, so reinsurance pricing degrades gracefully when the table
// is empty. Mirrors GetRegionLoading for the reinsurance side.
//
// Also short-circuits when the table is configured as not required: skips
// the cache and DB read entirely, returning zero.
func GetReinsuranceRegionLoading(riskRateCode string, gender string, region string) models.ReinsuranceRegionLoading {
	if !IsTableRequired("reinsuranceRegionLoading") {
		return models.ReinsuranceRegionLoading{}
	}
	tableName := "reinsurance_region_loadings"
	cacheKey := tableName + "_" + riskRateCode + "_" + strings.ToUpper(gender) + "_" + strings.TrimSpace(region)
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(models.ReinsuranceRegionLoading)
	}
	var loading models.ReinsuranceRegionLoading
	DB.Table(tableName).
		Where("risk_rate_code = ? AND gender = ? AND region = ?",
			riskRateCode, strings.ToUpper(gender), strings.TrimSpace(region)).
		First(&loading)
	GroupPricingCache.Set(cacheKey, loading, 1)
	return loading
}

// GetReinsuranceGlaRate looks up the reinsurance GLA Qx (re_qx) from
// reinsurance_gla_rates. The table has richer keys (income_level,
// waiting_period) so we honor them here; returns 0 when no match.
// Matches the direct-pricing GetGlaRate convention: full gender string,
// raw int for income_level / waiting_period (MySQL coerces to the column
// type). See api/services/group_pricing.go:6115 (GetGlaRate) for the
// reference pattern used by the existing gla_rates table.
func GetReinsuranceGlaRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, incomeLevel int, schemeCategory models.SchemeCategory) float64 {
	tableName := "reinsurance_gla_rates"
	cacheKey := tableName + "_" + groupPricingParameter.RiskRateCode + "_" + strconv.Itoa(memberResultData.AgeNextBirthday) + "_" + strconv.Itoa(incomeLevel) + "_" + memberResultData.Gender + "_" + strconv.Itoa(schemeCategory.GlaWaitingPeriod)
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}
	var qx float64
	DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND income_level = ? AND gender = ? AND waiting_period = ?",
			groupPricingParameter.RiskRateCode, memberResultData.AgeNextBirthday, incomeLevel, memberResultData.Gender, schemeCategory.GlaWaitingPeriod).
		Pluck("re_qx", &qx)
	GroupPricingCache.Set(cacheKey, qx, 1)
	return qx
}

// GetReinsuranceSpouseGlaRate — spouse variant.
func GetReinsuranceSpouseGlaRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, incomeLevel int, schemeCategory models.SchemeCategory) float64 {
	tableName := "reinsurance_gla_rates"
	cacheKey := tableName + "_spouse_" + groupPricingParameter.RiskRateCode + "_" + strconv.Itoa(memberResultData.SpouseAgeNextBirthday) + "_" + strconv.Itoa(incomeLevel) + "_" + memberResultData.SpouseGender + "_" + strconv.Itoa(schemeCategory.GlaWaitingPeriod)
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}
	var qx float64
	DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND income_level = ? AND gender = ? AND waiting_period = ?",
			groupPricingParameter.RiskRateCode, memberResultData.SpouseAgeNextBirthday, incomeLevel, memberResultData.SpouseGender, schemeCategory.GlaWaitingPeriod).
		Pluck("re_qx", &qx)
	GroupPricingCache.Set(cacheKey, qx, 1)
	return qx
}

// GetReinsurancePtdRate looks up the reinsurance PTD rate. Matches direct
// GetPtdRate (api/services/group_pricing.go:7040) convention: full gender
// string and raw int occupation_class.
func GetReinsurancePtdRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, incomeLevel int, schemeCategory models.SchemeCategory) float64 {
	tableName := "reinsurance_ptd_rates"
	cacheKey := tableName + "_" + groupPricingParameter.RiskRateCode + "_" + strconv.Itoa(memberResultData.AgeNextBirthday) + "_" + strconv.Itoa(incomeLevel) + "_" + memberResultData.Gender + "_" + strconv.Itoa(memberResultData.OccupationClass)
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}
	var rate float64
	DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND income_level = ? AND gender = ? AND occupation_class = ?",
			groupPricingParameter.RiskRateCode, memberResultData.AgeNextBirthday, incomeLevel, memberResultData.Gender, memberResultData.OccupationClass).
		Pluck("ptd_rate", &rate)
	GroupPricingCache.Set(cacheKey, rate, 1)
	return rate
}

// GetReinsuranceCiRate looks up the reinsurance CI rate.
func GetReinsuranceCiRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, incomeLevel int, schemeCategory models.SchemeCategory) float64 {
	tableName := "reinsurance_ci_rates"
	cacheKey := tableName + "_" + groupPricingParameter.RiskRateCode + "_" + strconv.Itoa(memberResultData.AgeNextBirthday) + "_" + strconv.Itoa(incomeLevel) + "_" + memberResultData.Gender + "_" + strconv.Itoa(memberResultData.OccupationClass)
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}
	var rate float64
	DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND income_level = ? AND gender = ? AND occupation_class = ?",
			groupPricingParameter.RiskRateCode, memberResultData.AgeNextBirthday, incomeLevel, memberResultData.Gender, memberResultData.OccupationClass).
		Pluck("ci_rate", &rate)
	GroupPricingCache.Set(cacheKey, rate, 1)
	return rate
}

// GetReinsurancePhiRate looks up the reinsurance PHI rate.
func GetReinsurancePhiRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, incomeLevel int, schemeCategory models.SchemeCategory) float64 {
	tableName := "reinsurance_phi_rates"
	cacheKey := tableName + "_" + groupPricingParameter.RiskRateCode + "_" + strconv.Itoa(memberResultData.AgeNextBirthday) + "_" + strconv.Itoa(incomeLevel) + "_" + memberResultData.Gender + "_" + strconv.Itoa(memberResultData.OccupationClass)
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.(float64)
	}
	var rate float64
	DB.Table(tableName).
		Where("risk_rate_code = ? AND age_next_birthday = ? AND income_level = ? AND gender = ? AND occupation_class = ?",
			groupPricingParameter.RiskRateCode, memberResultData.AgeNextBirthday, incomeLevel, memberResultData.Gender, memberResultData.OccupationClass).
		Pluck("phi_rate", &rate)
	GroupPricingCache.Set(cacheKey, rate, 1)
	return rate
}

// computeTakeHomePayFromBands calculates annual take-home pay from pre-fetched tax bands.
// Returns (takeHomePay, tax). Falls back to (annualSalary, 0) when bands are empty.
func computeTakeHomePayFromBands(annualSalary float64, taxBands []models.TaxTable) (float64, float64) {
	if len(taxBands) == 0 || annualSalary <= 0 {
		return annualSalary, 0
	}
	var tax float64
	remainingSalary := annualSalary
	for _, band := range taxBands {
		if remainingSalary <= 0 {
			break
		}
		if annualSalary <= band.Min {
			continue
		}
		applicableAmount := math.Min(remainingSalary, band.Max-band.Min)
		if band.Max <= 0 {
			applicableAmount = remainingSalary
		}
		if applicableAmount < 0 {
			applicableAmount = 0
		}
		tax += applicableAmount * band.TaxRate
		remainingSalary -= applicableAmount
	}
	return annualSalary - tax, tax
}

// computeCoveredIncomeFromTiers calculates annual covered income from pre-fetched tiered bands.
func computeCoveredIncomeFromTiers(annualSalary float64, tiers []models.TieredIncomeReplacement) float64 {
	if len(tiers) == 0 || annualSalary <= 0 {
		return 0
	}
	var coveredIncome float64
	remainingSalary := annualSalary
	for _, tier := range tiers {
		if remainingSalary <= 0 {
			break
		}
		if annualSalary <= tier.AnnualLowerBound {
			continue
		}
		applicableAmount := math.Min(remainingSalary, tier.AnnualUpperBound-tier.AnnualLowerBound)
		if applicableAmount < 0 {
			applicableAmount = 0
		}
		coveredIncome += applicableAmount * tier.IncomeReplacementRatio
		remainingSalary -= applicableAmount
	}
	return coveredIncome
}

func GetTieredIncomeReplacementTiers(riskRateCode string) ([]models.TieredIncomeReplacement, error) {
	tableName := "tiered_income_replacements"

	cacheKey := tableName + "_" + strings.TrimSpace(riskRateCode)
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.([]models.TieredIncomeReplacement), nil
	}

	var tiers []models.TieredIncomeReplacement
	err := DB.Table(tableName).
		Where("risk_rate_code = ?", strings.TrimSpace(riskRateCode)).
		Order("annual_lower_bound asc").
		Find(&tiers).Error
	if err != nil {
		return nil, err
	}

	GroupPricingCache.Set(cacheKey, tiers, 1)
	return tiers, nil
}

// GetCustomTieredIncomeReplacementTiers fetches custom tiered income replacement bands
// for a specific scheme name and risk rate code.
func GetCustomTieredIncomeReplacementTiers(schemeName, riskRateCode string) ([]models.TieredIncomeReplacement, error) {
	tableName := "custom_tiered_income_replacements"

	cacheKey := tableName + "_" + strings.TrimSpace(schemeName) + "_" + strings.TrimSpace(riskRateCode)
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.([]models.TieredIncomeReplacement), nil
	}

	var tiers []models.TieredIncomeReplacement
	err := DB.Table(tableName).
		Where("scheme_name = ? AND risk_rate_code = ?", strings.TrimSpace(schemeName), strings.TrimSpace(riskRateCode)).
		Order("annual_lower_bound asc").
		Find(&tiers).Error
	if err != nil {
		return nil, err
	}

	GroupPricingCache.Set(cacheKey, tiers, 1)
	return tiers, nil
}

// HasCustomTieredIncomeReplacementTable checks whether a custom tiered income replacement
// table exists for the given scheme name and risk rate code.
func HasCustomTieredIncomeReplacementTable(schemeName, riskRateCode string) bool {
	var count int64
	DB.Table("custom_tiered_income_replacements").
		Where("scheme_name = ? AND risk_rate_code = ?", strings.TrimSpace(schemeName), strings.TrimSpace(riskRateCode)).
		Count(&count)
	return count > 0
}

// CustomTirStatusResult is the result of a custom TIR preflight check for a quote.
type CustomTirStatusResult struct {
	NeedsCustomTir bool   `json:"needs_custom_tir"`
	HasTable       bool   `json:"has_table"`
	SchemeName     string `json:"scheme_name"`
	RiskRateCode   string `json:"risk_rate_code"`
}

// CheckCustomTirTableStatus inspects a quote's scheme categories to determine whether
// any of them require a custom tiered income replacement table, and if so, whether
// the table has been uploaded for the quote's scheme name and risk rate code.
func CheckCustomTirTableStatus(quoteId string) (CustomTirStatusResult, error) {
	var quote models.GroupPricingQuote
	if err := DB.Where("id = ?", quoteId).Preload("SchemeCategories").First(&quote).Error; err != nil {
		return CustomTirStatusResult{}, err
	}

	result := CustomTirStatusResult{
		SchemeName:   quote.SchemeName,
		RiskRateCode: quote.RiskRateCode,
		HasTable:     true, // default: not blocking
	}

	for _, sc := range quote.SchemeCategories {
		if (sc.PhiUseTieredIncomeReplacementRatio && sc.PhiTieredIncomeReplacementType == "custom") ||
			(sc.TtdUseTieredIncomeReplacementRatio && sc.TtdTieredIncomeReplacementType == "custom") {
			result.NeedsCustomTir = true
			break
		}
	}

	if result.NeedsCustomTir {
		result.HasTable = HasCustomTieredIncomeReplacementTable(quote.SchemeName, quote.RiskRateCode)
	}

	return result, nil
}

func GetCoveredIncomeFromTieredIncomeReplacement(annualSalary float64, riskRateCode string) float64 {
	tiers, err := GetTieredIncomeReplacementTiers(riskRateCode)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	if len(tiers) == 0 || annualSalary <= 0 {
		return 0
	}

	var coveredIncome float64
	remainingSalary := annualSalary

	for _, tier := range tiers {
		if remainingSalary <= 0 {
			break
		}

		// Portion of salary that falls into this tier
		tierStart := tier.AnnualLowerBound
		tierEnd := tier.AnnualUpperBound

		if annualSalary <= tierStart {
			continue
		}

		applicableAmount := math.Min(remainingSalary, tierEnd-tierStart)
		if applicableAmount < 0 {
			applicableAmount = 0
		}

		coveredIncome += applicableAmount * tier.IncomeReplacementRatio
		remainingSalary -= applicableAmount
	}

	return coveredIncome
}

func GetTaxTableByRiskRateCode(riskRateCode string) ([]models.TaxTable, error) {
	tableName := "tax_tables"

	cacheKey := tableName + "_" + strings.TrimSpace(riskRateCode)
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.([]models.TaxTable), nil
	}

	var taxTable []models.TaxTable
	err := DB.Table(tableName).
		Where("risk_rate_code = ?", strings.TrimSpace(riskRateCode)).
		Order("level asc").
		Find(&taxTable).Error
	if err != nil {
		return nil, err
	}

	GroupPricingCache.Set(cacheKey, taxTable, 1)
	return taxTable, nil
}

// computeTaxSaverSumAssured returns the extra cover required so the covered
// sum assured survives the retirement-fund lump-sum tax. Bands must be sorted
// by lower_bound ascending; the last band is treated as the open-ended top
// bracket (Case Is > lowerbound_N in the spec). The first band is paid in
// full with no gross-up, so a member whose covered SA falls in the tax-free
// tier gets TaxSaverSumAssured == 0.
func computeTaxSaverSumAssured(sumAssured float64, bands []models.TaxRetirementTable) float64 {
	if sumAssured <= 0 || len(bands) == 0 {
		return 0
	}
	totalBenefit := sumAssured
	for idx, band := range bands {
		isLast := idx == len(bands)-1
		inBracket := sumAssured >= band.LowerBound && (isLast || sumAssured <= band.UpperBound)
		if !inBracket {
			continue
		}
		if idx == 0 {
			totalBenefit = sumAssured
		} else if band.TaxRate < 1 {
			totalBenefit = (sumAssured - band.CumulativeTaxRelief) / (1 - band.TaxRate)
		}
		break
	}
	return totalBenefit - sumAssured
}

// GetTaxRetirementTableByRiskRateCode loads the retirement-tax bands used by
// the TaxSaver benefit to gross up GLA cover so the post-tax payout on
// retirement matches the member's covered sum assured. Rows come back ordered
// by lower_bound ascending; the last row is the open-ended top bracket.
func GetTaxRetirementTableByRiskRateCode(riskRateCode string) ([]models.TaxRetirementTable, error) {
	tableName := "tax_retirement_tables"

	cacheKey := tableName + "_" + strings.TrimSpace(riskRateCode)
	if cached, found := GroupPricingCache.Get(cacheKey); found {
		return cached.([]models.TaxRetirementTable), nil
	}

	var bands []models.TaxRetirementTable
	err := DB.Table(tableName).
		Where("risk_rate_code = ?", strings.TrimSpace(riskRateCode)).
		Order("lower_bound asc").
		Find(&bands).Error
	if err != nil {
		return nil, err
	}

	GroupPricingCache.Set(cacheKey, bands, 1)
	return bands, nil
}

func GetTakeHomePayFromTaxTable(annualSalary float64, riskRateCode string) (float64, float64) {
	taxBands, err := GetTaxTableByRiskRateCode(riskRateCode)
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}

	if len(taxBands) == 0 || annualSalary <= 0 {
		return annualSalary, 0
	}

	var tax float64
	remainingSalary := annualSalary

	for _, band := range taxBands {
		if remainingSalary <= 0 {
			break
		}

		bandMin := band.Min
		bandMax := band.Max

		if annualSalary <= bandMin {
			continue
		}

		// Amount of salary falling into this band
		applicableAmount := math.Min(remainingSalary, bandMax-bandMin)
		if bandMax <= 0 {
			// Safety for open-ended top band
			applicableAmount = remainingSalary
		}
		if applicableAmount < 0 {
			applicableAmount = 0
		}

		tax += applicableAmount * band.TaxRate
		remainingSalary -= applicableAmount
	}

	takeHomePay := annualSalary - tax
	return takeHomePay, tax
}

func GetPtdRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, groupQuote models.GroupPricingQuote, schemeCategory models.SchemeCategory, incomeLevel int) float64 {
	tableName := "ptd_rates"

	var keyString strings.Builder
	keyString.WriteString(schemeCategory.SchemeCategory + "_")
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(schemeCategory.PtdRiskType + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.AgeNextBirthday) + "_")
	keyString.WriteString(memberResultData.Gender[:1] + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.OccupationClass) + "_")
	keyString.WriteString(strconv.Itoa(incomeLevel) + "_")
	//keyString.WriteString(strconv.Itoa(groupQuote.Ptd.WaitingPeriod) + "_")
	keyString.WriteString(strconv.Itoa(schemeCategory.PtdDeferredPeriod) + "_")
	keyString.WriteString(schemeCategory.PtdDisabilityDefinition + "_")
	key := keyString.String()
	cacheKey := tableName + "_" + key
	cached, found := GroupPricingCache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	} else {
		//fmt.Println("cache missed: ", key)
	}
	var ptdRate float64
	query := "risk_rate_code=? and risk_type=? and age_next_birthday=? and gender=? and occupation_class=? and income_level=? and deferred_period=? and disability_definition=?"
	err := DB.Table(tableName).Where(query, groupPricingParameter.RiskRateCode, schemeCategory.PtdRiskType, memberResultData.AgeNextBirthday, memberResultData.Gender, memberResultData.OccupationClass, incomeLevel, schemeCategory.PtdDeferredPeriod, schemeCategory.PtdDisabilityDefinition).Pluck("ptd_rate", &ptdRate).Error
	if err != nil {
		fmt.Println(err)
	}
	GroupPricingCache.Set(cacheKey, ptdRate, 1)
	//time.Sleep(5 * time.Millisecond)
	return ptdRate
}

func GetTtdRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, groupQuote models.GroupPricingQuote, schemeCategory models.SchemeCategory, incomeLevel int) float64 {
	tableName := "ttd_rates"

	var keyString strings.Builder

	keyString.WriteString(schemeCategory.SchemeCategory + "_")
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.AgeNextBirthday) + "_")
	keyString.WriteString(memberResultData.Gender[:1] + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.OccupationClass) + "_")
	keyString.WriteString(strconv.Itoa(incomeLevel) + "_")
	keyString.WriteString(strconv.Itoa(schemeCategory.TtdWaitingPeriod) + "_")
	keyString.WriteString(schemeCategory.TtdDisabilityDefinition + "_")
	keyString.WriteString(schemeCategory.TtdRiskType + "_")
	key := keyString.String()
	cacheKey := tableName + "_" + key
	cached, found := GroupPricingCache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	} else {
		//fmt.Println("cache missed: ", key)
	}
	var ttdRate float64
	query := "risk_rate_code=? and age_next_birthday=? and gender=? and occupation_class=? and income_level=?  and waiting_period=? and disability_definition=? and risk_type=?"
	err := DB.Table(tableName).Where(query, groupPricingParameter.RiskRateCode, memberResultData.AgeNextBirthday, memberResultData.Gender, memberResultData.OccupationClass, incomeLevel, schemeCategory.TtdWaitingPeriod, schemeCategory.TtdDisabilityDefinition, schemeCategory.TtdRiskType).Pluck("ttd_rate", &ttdRate).Error
	if err != nil {
		fmt.Println(err)
	}
	GroupPricingCache.Set(cacheKey, ttdRate, 1)
	//time.Sleep(5 * time.Millisecond)
	return ttdRate
}

func GetPhiRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, groupQuote models.GroupPricingQuote, schemeCategory models.SchemeCategory, incomeLevel int) float64 {
	tableName := "phi_rates"

	var keyString strings.Builder

	keyString.WriteString(schemeCategory.SchemeCategory + "_")
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(schemeCategory.PhiRiskType + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.AgeNextBirthday) + "_")
	keyString.WriteString(memberResultData.Gender[:1] + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.OccupationClass) + "_")
	keyString.WriteString(strconv.Itoa(incomeLevel) + "_")
	keyString.WriteString(strconv.Itoa(schemeCategory.PhiWaitingPeriod) + "_")
	keyString.WriteString(strconv.Itoa(schemeCategory.PhiDeferredPeriod) + "_")
	keyString.WriteString(strconv.Itoa(schemeCategory.PhiNormalRetirementAge) + "_")
	keyString.WriteString(schemeCategory.PhiBenefitEscalation + "_")
	keyString.WriteString(schemeCategory.PhiDisabilityDefinition + "_")
	key := keyString.String()
	cacheKey := tableName + "_" + key
	cached, found := GroupPricingCache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	} else {
		//fmt.Println("cache missed: ", key)
	}
	var phiRate float64
	query := "risk_rate_code=? and risk_type=? and age_next_birthday=? and gender=? and occupation_class=? and income_level=? and waiting_period=? and deferred_period=? and normal_retirement_age=? and benefit_escalation_option=? and disability_definition=?"
	err := DB.Table(tableName).Where(query, groupPricingParameter.RiskRateCode, schemeCategory.PhiRiskType, memberResultData.AgeNextBirthday, memberResultData.Gender, memberResultData.OccupationClass, incomeLevel, schemeCategory.PhiWaitingPeriod, schemeCategory.PhiDeferredPeriod, schemeCategory.PhiNormalRetirementAge, schemeCategory.PhiBenefitEscalation, schemeCategory.PhiDisabilityDefinition).Pluck("phi_rate", &phiRate).Error
	if err != nil {
		fmt.Println(err)
	}
	GroupPricingCache.Set(cacheKey, phiRate, 1)
	//time.Sleep(5 * time.Millisecond)
	return phiRate
}

func GetCiRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, groupQuote models.GroupPricingQuote, schemeCategory models.SchemeCategory, incomeLevel int) float64 {
	tableName := "ci_rates"

	var keyString strings.Builder
	keyString.WriteString(schemeCategory.SchemeCategory + "_")
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.AgeNextBirthday) + "_")
	keyString.WriteString(memberResultData.Gender[:1] + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.OccupationClass) + "_")
	keyString.WriteString(strconv.Itoa(incomeLevel) + "_")
	keyString.WriteString(strconv.Itoa(schemeCategory.GlaWaitingPeriod) + "_")
	keyString.WriteString(schemeCategory.CiBenefitDefinition + "_")
	key := keyString.String()
	cacheKey := tableName + "_" + key
	cached, found := GroupPricingCache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	} else {
		//fmt.Println("cache missed: ", key)
	}
	var ciRate float64
	query := "risk_rate_code=? and age_next_birthday=? and gender=? and occupation_class=? and income_level=? and waiting_period=? and benefit_definition=?"
	err := DB.Table(tableName).Where(query, groupPricingParameter.RiskRateCode, memberResultData.AgeNextBirthday, memberResultData.Gender, memberResultData.OccupationClass, incomeLevel, schemeCategory.GlaWaitingPeriod, schemeCategory.CiBenefitDefinition).Pluck("ci_rate", &ciRate).Error
	if err != nil {
		fmt.Println(err)
	}
	GroupPricingCache.Set(cacheKey, ciRate, 1)
	//time.Sleep(5 * time.Millisecond)
	return ciRate
}

func GetSpouseGlaRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, incomeLevel int, groupQuote models.GroupPricingQuote, schemeCategory models.SchemeCategory) float64 {
	tableName := "gla_rates"

	var keyString strings.Builder

	keyString.WriteString(schemeCategory.SchemeCategory + "_")
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.SpouseAgeNextBirthday) + "_")
	keyString.WriteString(strconv.Itoa(incomeLevel) + "_") //incomelevel
	keyString.WriteString(memberResultData.SpouseGender[:1] + "_")
	keyString.WriteString(strconv.Itoa(schemeCategory.GlaWaitingPeriod) + "_")
	key := keyString.String()
	cacheKey := tableName + "_" + key
	cached, found := GroupPricingCache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	} else {
		//fmt.Println("cache missed: ", key)
	}
	var qx float64
	query := "risk_rate_code=? and age_next_birthday=? and income_level=? and gender=? and waiting_period=?"
	err := DB.Table(tableName).Where(query, groupPricingParameter.RiskRateCode, memberResultData.SpouseAgeNextBirthday, incomeLevel, memberResultData.SpouseGender, schemeCategory.GlaWaitingPeriod).Pluck("qx", &qx).Error
	if err != nil {
		fmt.Println(err)
	}
	GroupPricingCache.Set(cacheKey, qx, 1)
	//time.Sleep(5 * time.Millisecond)
	return qx
}

func GetEducatorRate(groupPricingParameter models.GroupPricingParameters, educatorBenefitCode string, anb int, incomeLevel int) models.EducatorRate {
	tableName := "educator_rates"
	var educatorRates models.EducatorRate
	var keyString strings.Builder

	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(educatorBenefitCode + "_")
	keyString.WriteString(strconv.Itoa(anb) + "_")
	keyString.WriteString(strconv.Itoa(incomeLevel) + "_")
	key := keyString.String()
	cacheKey := tableName + "_" + key
	cached, found := GroupPricingCache.Get(cacheKey)

	if found {
		result := cached.(models.EducatorRate)
		return result
	} else {
		query := "risk_rate_code=? and educator_benefit_code=? and age_next_birthday=? and income_level=?"
		err := DB.Table(tableName).Where(query, groupPricingParameter.RiskRateCode, educatorBenefitCode, anb, incomeLevel).First(&educatorRates).Error
		if err != nil {
			fmt.Println(err)
		}
		GroupPricingCache.Set(cacheKey, educatorRates, 1)
	}
	return educatorRates
}

func GetFuneralParameters(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, parameterVariable string) float64 {
	tableName := "funeral_parameters"
	var funeralParameter models.FuneralParameters
	var keyString strings.Builder

	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(memberResultData.AgeNextBirthday) + "_")
	key := keyString.String()
	cacheKey := tableName + "_" + key
	cached, found := GroupPricingCache.Get(cacheKey)

	if found {
		funeralParameter = cached.(models.FuneralParameters)
	} else {
		query := "risk_rate_code=? and age_next_birthday=?"
		err := DB.Table(tableName).Where(query, groupPricingParameter.RiskRateCode, memberResultData.SpouseAgeNextBirthday).First(&funeralParameter).Error
		if err == nil {
			GroupPricingCache.Set(cacheKey, funeralParameter, 1)
		}
	}
	switch parameterVariable {
	case "ProportionMarried":
		return funeralParameter.ProportionMarried
	case "NumberChildren":
		return funeralParameter.NumberChildren
	case "NumberDependants":
		return funeralParameter.NumberDependants
	case "AverageChildAge":
		return funeralParameter.AverageChildAge
	case "AverageDependantAge":
		return funeralParameter.AverageDependantAge

	}
	return 0
}

// CalculateExtendedFamilyAgeBandRates produces a loaded funeral rate for each
// requested age band. For every age present in the funeral_rates table it
// blends male and female FunQx using the extended-family male proportion from
// GroupPricingParameters:
//
//	qx[age] = maleProp * maleQx[age] + (1 - maleProp) * femaleQx[age]
//
// It then applies the extended-family loading from funeral_parameters (keyed
// by risk_rate_code + member_income_level + age_next_birthday) as
// `qx * (1 + loading)`, and averages the loaded rates across the integer ages
// in each band (AverageRate — the per-band risk rate). If only one gender has
// a rate for a given age (e.g. unisex data), that rate is used as-is.
//
// totalLoading is the combined premium loading from the premium_loading table
// (expense + admin + commission + profit + other − discount, floored at
// MinimumPremiumLoading). It is grossed up onto AverageRate to produce the
// office rate:
//
//	OfficeRate = AverageRate / (1 - totalLoading)
//
// sumsAssured, when non-empty, is merged into the result so callers can emit
// per-band risk and office monthly premiums without looking them up again.
// Rate-per-1000 callers pass an empty slice and derive the premiums as
// rate * 1000 / 12.
func CalculateExtendedFamilyAgeBandRates(
	riskRateCode string,
	memberIncomeLevel int,
	maleProp float64,
	totalLoading float64,
	bands []models.ExtendedFamilyAgeBand,
	sumsAssured []models.ExtendedFamilyBandSumAssured,
) ([]models.ExtendedFamilyBandRate, error) {
	if len(bands) == 0 {
		return nil, nil
	}

	// Clamp male proportion to [0, 1] so bad data in
	// group_pricing_parameters.extended_family_male_prop can't produce
	// negative or >100% weighted rates.
	if maleProp < 0 {
		maleProp = 0
	} else if maleProp > 1 {
		maleProp = 1
	}
	femaleProp := 1 - maleProp

	// 1) Pull all funeral_rates rows for this risk code, then blend
	//    male/female FunQx per age using maleProp.
	var rateRows []models.FuneralRate
	if err := DB.Table("funeral_rates").
		Where("risk_rate_code = ?", riskRateCode).
		Find(&rateRows).Error; err != nil {
		return nil, err
	}
	if len(rateRows) == 0 {
		return nil, fmt.Errorf("no funeral rates found for risk_rate_code=%s", riskRateCode)
	}

	type genderAgg struct {
		maleSum     float64
		maleCount   int
		femaleSum   float64
		femaleCount int
	}
	ageToQx := make(map[int]float64)
	{
		agg := make(map[int]*genderAgg)
		for _, r := range rateRows {
			a := agg[r.AgeNextBirthday]
			if a == nil {
				a = &genderAgg{}
				agg[r.AgeNextBirthday] = a
			}
			// Gender values are typically "M"/"F" but normalise defensively.
			switch strings.ToUpper(strings.TrimSpace(r.Gender)) {
			case "M", "MALE":
				a.maleSum += r.FunQx
				a.maleCount++
			case "F", "FEMALE":
				a.femaleSum += r.FunQx
				a.femaleCount++
			default:
				// Unisex/unknown gender — count in both buckets so it still
				// contributes when the opposite gender is absent.
				a.maleSum += r.FunQx
				a.maleCount++
				a.femaleSum += r.FunQx
				a.femaleCount++
			}
		}
		for age, a := range agg {
			hasMale := a.maleCount > 0
			hasFemale := a.femaleCount > 0
			switch {
			case hasMale && hasFemale:
				maleQx := a.maleSum / float64(a.maleCount)
				femaleQx := a.femaleSum / float64(a.femaleCount)
				ageToQx[age] = maleProp*maleQx + femaleProp*femaleQx
			case hasMale:
				ageToQx[age] = a.maleSum / float64(a.maleCount)
			case hasFemale:
				ageToQx[age] = a.femaleSum / float64(a.femaleCount)
			}
		}
	}

	// 2) Pull loadings for (risk_rate_code, member_income_level) and index by age.
	var paramRows []models.FuneralParameters
	if err := DB.Table("funeral_parameters").
		Where("risk_rate_code = ? AND member_income_level = ?", riskRateCode, memberIncomeLevel).
		Find(&paramRows).Error; err != nil {
		return nil, err
	}
	ageToLoading := make(map[int]float64, len(paramRows))
	for _, p := range paramRows {
		ageToLoading[p.AgeNextBirthday] = p.ExtendedFamilyLoading
	}

	// 3) Build loadedQx[age] = qx * (1 + loading) for every age we have a rate.
	loadedQx := make(map[int]float64, len(ageToQx))
	maxAge := 0
	for age, qx := range ageToQx {
		loadedQx[age] = qx * (1 + ageToLoading[age])
		if age > maxAge {
			maxAge = age
		}
	}

	// 4) Map sums-assured by (min,max) so we can merge into results.
	sumByKey := make(map[[2]int]float64, len(sumsAssured))
	for _, s := range sumsAssured {
		sumByKey[[2]int{s.MinAge, s.MaxAge}] = s.SumAssured
	}

	// Office-rate gross-up divisor. A totalLoading outside [0, 1) would either
	// yield zero/negative premiums or a negative divisor, so fall back to the
	// risk rate (divisor = 1) in those degenerate cases.
	loadingDivisor := 1.0 - totalLoading
	if loadingDivisor <= 0 {
		loadingDivisor = 1.0
	}

	// 5) Straight-average loadedQx over the integer ages in each band.
	results := make([]models.ExtendedFamilyBandRate, 0, len(bands))
	for _, b := range bands {
		lo := b.MinAge
		hi := b.MaxAge
		if hi < lo {
			continue
		}
		// Clip upper bound to the maximum age available in the rate table.
		if hi > maxAge {
			hi = maxAge
		}
		var sum float64
		var count int
		for age := lo; age <= hi; age++ {
			if q, ok := loadedQx[age]; ok {
				sum += q
				count++
			}
		}
		var avg float64
		if count > 0 {
			avg = sum / float64(count)
		}
		officeRate := avg / loadingDivisor
		res := models.ExtendedFamilyBandRate{
			MinAge:      b.MinAge,
			MaxAge:      b.MaxAge,
			AverageRate: avg,
			OfficeRate:  officeRate,
		}
		if sa, ok := sumByKey[[2]int{b.MinAge, b.MaxAge}]; ok {
			res.SumAssured = sa
			res.MonthlyPremium = avg * sa / 12.0
			res.OfficeMonthlyPremium = officeRate * sa / 12.0
		}
		results = append(results, res)
	}
	return results, nil
}

// CalculateAdditionalGlaCoverBandRates produces per-age-band office-rate-per-1000
// figures for the optional rate-only Additional GLA Cover. It mirrors
// CalculateExtendedFamilyAgeBandRates but draws from gla_rates + gla_aids_rates
// and layers the same region / industry / contingency loadings that base GLA
// uses, then grosses the risk rate up by (1 - totalPremiumLoading).
//
// All rates are quoted per 1,000 sum assured. No per-member premium is
// computed here — this is a rate-only product.
//
// Ages used come from gla_rates itself (whatever the table actually supplies
// for the given risk code / benefit type / waiting period / income level).
// Bands outside that range are clipped and empty bands are skipped.
func CalculateAdditionalGlaCoverBandRates(
	riskRateCode, benefitType, region string,
	waitingPeriod, incomeLevel, occupationClass int,
	maleProp float64,
	bands []models.AdditionalGlaCoverAgeBand,
	commissionLoading, binderFeeRate, outsourceFeeRate, otherLoadingsSum, minimumPremiumLoading float64,
	members []models.GPricingMemberData,
	referenceDate time.Time,
	categoryGlaSalaryMultiple float64,
) ([]models.AdditionalGlaCoverBandRate, error) {
	totalPremiumLoading := math.Max(
		commissionLoading+binderFeeRate+outsourceFeeRate+otherLoadingsSum,
		minimumPremiumLoading,
	)
	if len(bands) == 0 {
		return nil, nil
	}

	// Clamp male proportion to [0, 1] so a bad UI or DB value can't
	// produce negative or > 100% weighted rates.
	if maleProp < 0 {
		maleProp = 0
	} else if maleProp > 1 {
		maleProp = 1
	}
	femaleProp := 1 - maleProp

	// 1) Pull gla_rates rows across whatever ages the table has for this
	// risk code / benefit type / waiting period / income level. The
	// base-GLA path (GetGlaRate) binds the int income_level directly and
	// relies on the DB's implicit conversion; do the same here.
	var glaRows []models.GlaRate
	if err := DB.Table("gla_rates").
		Where("risk_rate_code = ? AND benefit_type = ? AND waiting_period = ? AND income_level = ?",
			riskRateCode, benefitType, waitingPeriod, incomeLevel).
		Find(&glaRows).Error; err != nil {
		return nil, err
	}
	if len(glaRows) == 0 {
		return nil, fmt.Errorf("no gla_rates found for risk_rate_code=%s, benefit_type=%s, waiting_period=%d, income_level=%d",
			riskRateCode, benefitType, waitingPeriod, incomeLevel)
	}

	type genderAgg struct {
		maleSum     float64
		maleCount   int
		femaleSum   float64
		femaleCount int
	}
	// pickMFC reduces a per-age male/female aggregation to (male-only,
	// female-only, blended) rates. When one side is missing the other side
	// stands in for it, so the male-only / female-only columns degrade to
	// the available data instead of falling to zero.
	pickMFC := func(a *genderAgg) (m, f, c float64, ok bool) {
		hasM := a.maleCount > 0
		hasF := a.femaleCount > 0
		switch {
		case hasM && hasF:
			m = a.maleSum / float64(a.maleCount)
			f = a.femaleSum / float64(a.femaleCount)
			c = maleProp*m + femaleProp*f
			return m, f, c, true
		case hasM:
			m = a.maleSum / float64(a.maleCount)
			return m, m, m, true
		case hasF:
			f = a.femaleSum / float64(a.femaleCount)
			return f, f, f, true
		}
		return 0, 0, 0, false
	}
	ageToQx := make(map[int]float64)
	ageToQxMale := make(map[int]float64)
	ageToQxFemale := make(map[int]float64)
	{
		agg := make(map[int]*genderAgg)
		for _, r := range glaRows {
			a := agg[r.AgeNextBirthday]
			if a == nil {
				a = &genderAgg{}
				agg[r.AgeNextBirthday] = a
			}
			switch strings.ToUpper(strings.TrimSpace(r.Gender)) {
			case "M", "MALE":
				a.maleSum += r.Qx
				a.maleCount++
			case "F", "FEMALE":
				a.femaleSum += r.Qx
				a.femaleCount++
			default:
				// Unisex/unknown gender — count in both buckets so it still
				// contributes when the opposite gender is absent.
				a.maleSum += r.Qx
				a.maleCount++
				a.femaleSum += r.Qx
				a.femaleCount++
			}
		}
		for age, a := range agg {
			if m, f, c, ok := pickMFC(a); ok {
				ageToQxMale[age] = m
				ageToQxFemale[age] = f
				ageToQx[age] = c
			}
		}
	}

	// Record the min/max age actually present in gla_rates so bands can
	// be clipped without a hardcoded assumption.
	minTableAge, maxTableAge := 0, 0
	firstAge := true
	for age := range ageToQx {
		if firstAge {
			minTableAge = age
			maxTableAge = age
			firstAge = false
			continue
		}
		if age < minTableAge {
			minTableAge = age
		}
		if age > maxTableAge {
			maxTableAge = age
		}
	}

	// 2) Pull gla_aids_rates across all ages for this risk code and blend
	// male/female. Missing ages are implicitly 0 (matches the per-member
	// path which simply looks up 0 when nothing is stored).
	var aidsRows []models.GlaAidsRate
	if err := DB.Table("gla_aids_rates").
		Where("risk_rate_code = ?", riskRateCode).
		Find(&aidsRows).Error; err != nil {
		return nil, err
	}
	ageToAidsQx := make(map[int]float64)
	ageToAidsQxMale := make(map[int]float64)
	ageToAidsQxFemale := make(map[int]float64)
	{
		agg := make(map[int]*genderAgg)
		for _, r := range aidsRows {
			a := agg[r.AgeNextBirthday]
			if a == nil {
				a = &genderAgg{}
				agg[r.AgeNextBirthday] = a
			}
			switch strings.ToUpper(strings.TrimSpace(r.Gender)) {
			case "M", "MALE":
				a.maleSum += r.GlaAidsQx
				a.maleCount++
			case "F", "FEMALE":
				a.femaleSum += r.GlaAidsQx
				a.femaleCount++
			default:
				a.maleSum += r.GlaAidsQx
				a.maleCount++
				a.femaleSum += r.GlaAidsQx
				a.femaleCount++
			}
		}
		for age, a := range agg {
			if m, f, c, ok := pickMFC(a); ok {
				ageToAidsQxMale[age] = m
				ageToAidsQxFemale[age] = f
				ageToAidsQx[age] = c
			}
		}
	}

	// 3) Resolve region loadings (keyed by gender) and blend by male prop.
	// The loading is applied uniformly across ages — it doesn't vary by age.
	var regionRows []models.RegionLoading
	if region != "" {
		if err := DB.Table("region_loadings").
			Where("risk_rate_code = ? AND region = ?", riskRateCode, strings.TrimSpace(region)).
			Find(&regionRows).Error; err != nil {
			return nil, err
		}
	}
	var regionLoading, aidsRegionLoading float64
	var regionLoadingMale, regionLoadingFemale float64
	var aidsRegionLoadingMale, aidsRegionLoadingFemale float64
	{
		var rm, rf, rmAids, rfAids float64
		var hasM, hasF bool
		for _, r := range regionRows {
			g := strings.ToUpper(strings.TrimSpace(r.Gender))
			if g == "M" || g == "MALE" {
				rm = r.GlaRegionLoadingRate
				rmAids = r.GlaAidsRegionLoadingRate
				hasM = true
			} else if g == "F" || g == "FEMALE" {
				rf = r.GlaRegionLoadingRate
				rfAids = r.GlaAidsRegionLoadingRate
				hasF = true
			}
		}
		switch {
		case hasM && hasF:
			regionLoadingMale, regionLoadingFemale = rm, rf
			aidsRegionLoadingMale, aidsRegionLoadingFemale = rmAids, rfAids
			regionLoading = maleProp*rm + femaleProp*rf
			aidsRegionLoading = maleProp*rmAids + femaleProp*rfAids
		case hasM:
			regionLoadingMale, regionLoadingFemale = rm, rm
			aidsRegionLoadingMale, aidsRegionLoadingFemale = rmAids, rmAids
			regionLoading = rm
			aidsRegionLoading = rmAids
		case hasF:
			regionLoadingMale, regionLoadingFemale = rf, rf
			aidsRegionLoadingMale, aidsRegionLoadingFemale = rfAids, rfAids
			regionLoading = rf
			aidsRegionLoading = rfAids
		}
	}

	// 4) Industry loadings — keyed by (risk_rate_code, occupation_class, gender).
	// Blend by male prop, uniform across ages.
	var industryRows []models.IndustryLoading
	if err := DB.Table("industry_loadings").
		Where("risk_rate_code = ? AND occupation_class = ?", riskRateCode, occupationClass).
		Find(&industryRows).Error; err != nil {
		return nil, err
	}
	var industryLoading, industryLoadingMale, industryLoadingFemale float64
	{
		var im, iFem float64
		var hasM, hasF bool
		for _, il := range industryRows {
			g := strings.ToUpper(strings.TrimSpace(il.Gender))
			if g == "M" || g == "MALE" {
				im = il.GlaIndustryLoadingRate
				hasM = true
			} else if g == "F" || g == "FEMALE" {
				iFem = il.GlaIndustryLoadingRate
				hasF = true
			}
		}
		switch {
		case hasM && hasF:
			industryLoadingMale, industryLoadingFemale = im, iFem
			industryLoading = maleProp*im + femaleProp*iFem
		case hasM:
			industryLoadingMale, industryLoadingFemale = im, im
			industryLoading = im
		case hasF:
			industryLoadingMale, industryLoadingFemale = iFem, iFem
			industryLoading = iFem
		}
	}

	// 5) Per-age contingency loading from general_loadings (varies by age
	// and gender) — keep the male and female rates separately and blend
	// the Combined value with male prop so it matches the qx blend above.
	ageToContingency := make(map[int]float64, len(ageToQx))
	ageToContingencyMale := make(map[int]float64, len(ageToQx))
	ageToContingencyFemale := make(map[int]float64, len(ageToQx))
	for age := range ageToQx {
		m := GetGeneralLoading(riskRateCode, age, "M").GlaContigencyLoadingRate
		f := GetGeneralLoading(riskRateCode, age, "F").GlaContigencyLoadingRate
		ageToContingencyMale[age] = m
		ageToContingencyFemale[age] = f
		ageToContingency[age] = maleProp*m + femaleProp*f
	}

	// 6) Build per-age loaded rates for the male, female and blended cohorts.
	//    For each cohort:
	//        baseRate   = qx * (1 + industry + region) + aidsQx * (1 + aidsRegion)
	//        loadedRate = baseRate * (1 + contingency)
	loadedRate := make(map[int]float64, len(ageToQx))
	loadedRateMale := make(map[int]float64, len(ageToQx))
	loadedRateFemale := make(map[int]float64, len(ageToQx))
	for age, qxC := range ageToQx {
		qxM, qxF := ageToQxMale[age], ageToQxFemale[age]
		aC, aM, aF := ageToAidsQx[age], ageToAidsQxMale[age], ageToAidsQxFemale[age]
		cC, cM, cF := ageToContingency[age], ageToContingencyMale[age], ageToContingencyFemale[age]

		baseC := qxC*(1+industryLoading+regionLoading) + aC*(1+aidsRegionLoading)
		baseM := qxM*(1+industryLoadingMale+regionLoadingMale) + aM*(1+aidsRegionLoadingMale)
		baseF := qxF*(1+industryLoadingFemale+regionLoadingFemale) + aF*(1+aidsRegionLoadingFemale)

		loadedRate[age] = baseC * (1 + cC)
		loadedRateMale[age] = baseM * (1 + cM)
		loadedRateFemale[age] = baseF * (1 + cF)
	}

	// 7) Office-rate gross-up divisor. Out-of-range totalPremiumLoading
	// falls back to the risk rate (divisor = 1) so we never produce a
	// zero or negative denominator.
	loadingDivisor := 1.0 - totalPremiumLoading
	if loadingDivisor <= 0 {
		loadingDivisor = 1.0
	}

	// 8) Straight-average each cohort's loaded rate across the integer ages
	// in each band, clipped to the range the gla_rates table actually
	// covers, then gross up to the office rate and derive the per-1,000
	// fee splits. The Combined column is blended at maleProp; the Male
	// and Female columns are computed as if the population were 100% male
	// or 100% female.
	avgBand := func(m map[int]float64, lo, hi int) float64 {
		var sum float64
		var count int
		for age := lo; age <= hi; age++ {
			if r, ok := m[age]; ok {
				sum += r
				count++
			}
		}
		if count == 0 {
			return 0
		}
		return sum / float64(count)
	}
	// Per-band exposure-weighted office rate, using each member's exact-age
	// office rate weighted by their GLA covered sum assured (= multiple ×
	// salary, uncapped). Members are already filtered to the scheme category
	// upstream, so we don't gate on Benefits.GlaEnabled — that flag is
	// populated only after the calc engine seeds it from the category, and
	// when reading raw rows from DB it's typically false.
	// Returns nil pointers when the band has no members contributing SA so
	// the UI can render an empty cell instead of zero.
	weightedForBand := func(minAge, maxAge int) (m, f, c *float64) {
		var maleSARate, maleSA, femaleSARate, femaleSA float64
		// Combined is the SA-weighted average across both genders, using
		// each member's own gender-specific rate. When only one gender has
		// members in the band, Combined collapses to that gender's value.
		var combinedSARate, combinedSA float64
		for _, mp := range members {
			multiple := mp.Benefits.GlaMultiple
			if multiple <= 0 {
				multiple = categoryGlaSalaryMultiple
			}
			sa := multiple * mp.AnnualSalary
			if sa <= 0 {
				continue
			}
			age := calculateAgeNextBirthday(referenceDate, mp.DateOfBirth)
			if age < minAge || age > maxAge {
				continue
			}
			g := strings.ToUpper(strings.TrimSpace(mp.Gender))
			switch g {
			case "M", "MALE":
				if r, ok := loadedRateMale[age]; ok {
					rate := (r / loadingDivisor) * 1000.0
					maleSARate += rate * sa
					maleSA += sa
					combinedSARate += rate * sa
					combinedSA += sa
				}
			case "F", "FEMALE":
				if r, ok := loadedRateFemale[age]; ok {
					rate := (r / loadingDivisor) * 1000.0
					femaleSARate += rate * sa
					femaleSA += sa
					combinedSARate += rate * sa
					combinedSA += sa
				}
			}
		}
		if maleSA > 0 {
			v := maleSARate / maleSA
			m = &v
		}
		if femaleSA > 0 {
			v := femaleSARate / femaleSA
			f = &v
		}
		if combinedSA > 0 {
			v := combinedSARate / combinedSA
			c = &v
		}
		return
	}

	results := make([]models.AdditionalGlaCoverBandRate, 0, len(bands))
	for _, b := range bands {
		lo := b.MinAge
		hi := b.MaxAge
		if hi < lo {
			continue
		}
		if lo < minTableAge {
			lo = minTableAge
		}
		if hi > maxTableAge {
			hi = maxTableAge
		}
		if hi < lo {
			// Band fell entirely outside the rate table — emit a zero
			// row so the caller can see the clip took the band to empty.
			results = append(results, models.AdditionalGlaCoverBandRate{
				MinAge:       b.MinAge,
				MaxAge:       b.MaxAge,
				MalePropUsed: maleProp,
			})
			continue
		}
		avgC := avgBand(loadedRate, lo, hi)
		avgM := avgBand(loadedRateMale, lo, hi)
		avgF := avgBand(loadedRateFemale, lo, hi)

		officeC := (avgC / loadingDivisor) * 1000.0
		officeM := (avgM / loadingDivisor) * 1000.0
		officeF := (avgF / loadingDivisor) * 1000.0

		weightedM, weightedF, weightedC := weightedForBand(b.MinAge, b.MaxAge)

		results = append(results, models.AdditionalGlaCoverBandRate{
			MinAge:                          b.MinAge,
			MaxAge:                          b.MaxAge,
			RiskRatePer1000:                 avgC * 1000.0,
			RiskRatePer1000Male:             avgM * 1000.0,
			RiskRatePer1000Female:           avgF * 1000.0,
			BinderFeePer1000:                officeC * binderFeeRate,
			BinderFeePer1000Male:            officeM * binderFeeRate,
			BinderFeePer1000Female:          officeF * binderFeeRate,
			OutsourceFeePer1000:             officeC * outsourceFeeRate,
			OutsourceFeePer1000Male:         officeM * outsourceFeeRate,
			OutsourceFeePer1000Female:       officeF * outsourceFeeRate,
			CommissionPer1000:               officeC * commissionLoading,
			CommissionPer1000Male:           officeM * commissionLoading,
			CommissionPer1000Female:         officeF * commissionLoading,
			OfficeRatePer1000:               officeC,
			OfficeRatePer1000Male:           officeM,
			OfficeRatePer1000Female:         officeF,
			MalePropUsed:                    maleProp,
			WeightedOfficeRatePer1000:       weightedC,
			WeightedOfficeRatePer1000Male:   weightedM,
			WeightedOfficeRatePer1000Female: weightedF,
		})
	}
	return results, nil
}

// AdditionalGlaSmoothedRateRow is one ageband update from the UI. min_age and
// max_age identify the band; the smoothed/factor pointers are applied verbatim
// (nil leaves the persisted value unchanged via the explicit `clear` flags).
type AdditionalGlaSmoothedRateRow struct {
	MinAge int `json:"min_age"`
	MaxAge int `json:"max_age"`

	SmoothedOfficeRatePer1000       *float64 `json:"smoothed_office_rate_per1000,omitempty"`
	SmoothedOfficeRatePer1000Male   *float64 `json:"smoothed_office_rate_per1000_male,omitempty"`
	SmoothedOfficeRatePer1000Female *float64 `json:"smoothed_office_rate_per1000_female,omitempty"`

	SmoothingFactor       *float64 `json:"smoothing_factor,omitempty"`
	SmoothingFactorMale   *float64 `json:"smoothing_factor_male,omitempty"`
	SmoothingFactorFemale *float64 `json:"smoothing_factor_female,omitempty"`

	// When true, the corresponding pointer field above is cleared to nil
	// regardless of the value. Lets the UI revert a single cell back to the
	// computed OfficeRate without sending a sentinel.
	ClearSmoothed       bool `json:"clear_smoothed,omitempty"`
	ClearSmoothedMale   bool `json:"clear_smoothed_male,omitempty"`
	ClearSmoothedFemale bool `json:"clear_smoothed_female,omitempty"`
	ClearFactor         bool `json:"clear_factor,omitempty"`
	ClearFactorMale     bool `json:"clear_factor_male,omitempty"`
	ClearFactorFemale   bool `json:"clear_factor_female,omitempty"`
}

// ErrAglaSmoothingLocked is returned when the caller tries to mutate the
// smoothed rates of a quote whose status no longer allows underwriter edits
// (approved / accepted / in-force). The HTTP layer maps this to 409 Conflict
// so the UI can show a friendly message rather than a generic 500.
var ErrAglaSmoothingLocked = errors.New("additional GLA smoothed rates are locked for this quote status")

// AdditionalGlaSmoothedSaveResult is the structured response from a save: the
// updated band rates plus the audit metadata so the UI can refresh the
// "last updated by X at Y" banner without refetching the whole quote.
type AdditionalGlaSmoothedSaveResult struct {
	BandRates []models.AdditionalGlaCoverBandRate `json:"additional_gla_cover_band_rates"`
	UpdatedAt *time.Time                          `json:"additional_gla_smoothed_updated_at"`
	UpdatedBy string                              `json:"additional_gla_smoothed_updated_by"`
}

// SaveAdditionalGlaCoverSmoothedRates merges per-band smoothed rates and
// smoothing factors into the persisted scheme_categories row for the given
// quote+category and mirrors the result onto the matching member_rating_result_summaries
// row so the Premium Summary refreshes without requiring a full quote recalc.
// Only the band rate JSON column is touched; the rest of the snapshot is
// untouched.
//
// updatedBy is recorded against the scheme_categories row alongside the
// timestamp so the UI can show "last updated by X at Y". Pass an empty
// string when there's no user context (e.g. background jobs).
func SaveAdditionalGlaCoverSmoothedRates(quoteID int, category string, rows []AdditionalGlaSmoothedRateRow, updatedBy string) (*AdditionalGlaSmoothedSaveResult, error) {
	// Reject the save once the quote is locked. Smoothed rates feed into the
	// agreed pricing, so we can't let them shift after the quote is approved
	// or accepted.
	var quote models.GroupPricingQuote
	if err := DB.Select("status").Where("id = ?", quoteID).First(&quote).Error; err != nil {
		return nil, fmt.Errorf("quote %d not found: %w", quoteID, err)
	}
	switch quote.Status {
	case models.StatusApproved, models.StatusAccepted, models.StatusInForce:
		return nil, ErrAglaSmoothingLocked
	}

	var sc models.SchemeCategory
	if err := DB.Where("quote_id = ? AND scheme_category = ?", quoteID, category).
		First(&sc).Error; err != nil {
		return nil, fmt.Errorf("scheme category not found for quote %d, category %q: %w", quoteID, category, err)
	}
	if len(sc.AdditionalGlaCoverBandRates) == 0 {
		return nil, fmt.Errorf("no additional GLA cover band rates persisted for quote %d, category %q — recalculate the quote first", quoteID, category)
	}

	updates := make(map[[2]int]AdditionalGlaSmoothedRateRow, len(rows))
	for _, r := range rows {
		updates[[2]int{r.MinAge, r.MaxAge}] = r
	}

	for i := range sc.AdditionalGlaCoverBandRates {
		band := &sc.AdditionalGlaCoverBandRates[i]
		u, ok := updates[[2]int{band.MinAge, band.MaxAge}]
		if !ok {
			continue
		}
		if u.ClearSmoothed {
			band.SmoothedOfficeRatePer1000 = nil
		} else if u.SmoothedOfficeRatePer1000 != nil {
			band.SmoothedOfficeRatePer1000 = u.SmoothedOfficeRatePer1000
		}
		if u.ClearSmoothedMale {
			band.SmoothedOfficeRatePer1000Male = nil
		} else if u.SmoothedOfficeRatePer1000Male != nil {
			band.SmoothedOfficeRatePer1000Male = u.SmoothedOfficeRatePer1000Male
		}
		if u.ClearSmoothedFemale {
			band.SmoothedOfficeRatePer1000Female = nil
		} else if u.SmoothedOfficeRatePer1000Female != nil {
			band.SmoothedOfficeRatePer1000Female = u.SmoothedOfficeRatePer1000Female
		}
		if u.ClearFactor {
			band.SmoothingFactor = nil
		} else if u.SmoothingFactor != nil {
			band.SmoothingFactor = u.SmoothingFactor
		}
		if u.ClearFactorMale {
			band.SmoothingFactorMale = nil
		} else if u.SmoothingFactorMale != nil {
			band.SmoothingFactorMale = u.SmoothingFactorMale
		}
		if u.ClearFactorFemale {
			band.SmoothingFactorFemale = nil
		} else if u.SmoothingFactorFemale != nil {
			band.SmoothingFactorFemale = u.SmoothingFactorFemale
		}

		// Propagate the smoothed-or-original office rate into the
		// regular OfficeRatePer1000* fields and recompute the
		// associated Binder / Outsource / Commission per-1,000 so any
		// downstream consumer reading the band rates gets the final
		// (smoothed-applied) value without an extra recalc. The
		// pre-smoothed reference is held in OriginalOfficeRatePer1000*
		// — capture it on first save when missing so the comparison
		// view still has it.
		if band.OriginalOfficeRatePer1000 == nil {
			v := band.OfficeRatePer1000
			band.OriginalOfficeRatePer1000 = &v
		}
		if band.OriginalOfficeRatePer1000Male == nil {
			v := band.OfficeRatePer1000Male
			band.OriginalOfficeRatePer1000Male = &v
		}
		if band.OriginalOfficeRatePer1000Female == nil {
			v := band.OfficeRatePer1000Female
			band.OriginalOfficeRatePer1000Female = &v
		}

		// Recover per-category Binder / Outsource / Commission rates
		// from the existing band before we overwrite the office rate.
		// Falls back to 0 when the band hasn't been computed yet.
		var binderRate, outsourceRate, commissionRate float64
		if band.OfficeRatePer1000 > 0 {
			binderRate = band.BinderFeePer1000 / band.OfficeRatePer1000
			outsourceRate = band.OutsourceFeePer1000 / band.OfficeRatePer1000
			commissionRate = band.CommissionPer1000 / band.OfficeRatePer1000
		}

		if band.SmoothedOfficeRatePer1000 != nil {
			band.OfficeRatePer1000 = *band.SmoothedOfficeRatePer1000
		} else if band.OriginalOfficeRatePer1000 != nil {
			band.OfficeRatePer1000 = *band.OriginalOfficeRatePer1000
		}
		if band.SmoothedOfficeRatePer1000Male != nil {
			band.OfficeRatePer1000Male = *band.SmoothedOfficeRatePer1000Male
		} else if band.OriginalOfficeRatePer1000Male != nil {
			band.OfficeRatePer1000Male = *band.OriginalOfficeRatePer1000Male
		}
		if band.SmoothedOfficeRatePer1000Female != nil {
			band.OfficeRatePer1000Female = *band.SmoothedOfficeRatePer1000Female
		} else if band.OriginalOfficeRatePer1000Female != nil {
			band.OfficeRatePer1000Female = *band.OriginalOfficeRatePer1000Female
		}

		band.BinderFeePer1000 = band.OfficeRatePer1000 * binderRate
		band.BinderFeePer1000Male = band.OfficeRatePer1000Male * binderRate
		band.BinderFeePer1000Female = band.OfficeRatePer1000Female * binderRate
		band.OutsourceFeePer1000 = band.OfficeRatePer1000 * outsourceRate
		band.OutsourceFeePer1000Male = band.OfficeRatePer1000Male * outsourceRate
		band.OutsourceFeePer1000Female = band.OfficeRatePer1000Female * outsourceRate
		band.CommissionPer1000 = band.OfficeRatePer1000 * commissionRate
		band.CommissionPer1000Male = band.OfficeRatePer1000Male * commissionRate
		band.CommissionPer1000Female = band.OfficeRatePer1000Female * commissionRate
	}

	now := time.Now()
	if err := DB.Model(&models.SchemeCategory{}).
		Where("id = ?", sc.ID).
		Updates(map[string]interface{}{
			"additional_gla_cover_band_rates":    sc.AdditionalGlaCoverBandRates,
			"additional_gla_smoothed_updated_at": now,
			"additional_gla_smoothed_updated_by": updatedBy,
		}).Error; err != nil {
		return nil, fmt.Errorf("persist scheme_categories.additional_gla_cover_band_rates: %w", err)
	}
	sc.AdditionalGlaSmoothedUpdatedAt = &now
	sc.AdditionalGlaSmoothedUpdatedBy = updatedBy

	// Mirror onto the MRRS row so the result-summary GET surfaces the
	// updated band rates (including the smoothed snapshot) on the next
	// reload. Best-effort — log and continue if no MRRS row.
	// Note: MRRS keys the category by `category` (not `scheme_category`).
	if err := DB.Model(&models.MemberRatingResultSummary{}).
		Where("quote_id = ? AND category = ?", quoteID, category).
		Updates(map[string]interface{}{
			"additional_gla_cover_band_rates": sc.AdditionalGlaCoverBandRates,
		}).Error; err != nil {
		logger := appLog.WithContext(context.Background())
		logger.WithFields(map[string]interface{}{
			"error":    err.Error(),
			"quote_id": quoteID,
			"category": category,
		}).Warn("Failed to mirror smoothed AGLA rates onto MemberRatingResultSummary; Premium Summary may show stale values until next recalc")
	}

	return &AdditionalGlaSmoothedSaveResult{
		BandRates: sc.AdditionalGlaCoverBandRates,
		UpdatedAt: sc.AdditionalGlaSmoothedUpdatedAt,
		UpdatedBy: sc.AdditionalGlaSmoothedUpdatedBy,
	}, nil
}

// ResolveAdditionalGlaCoverMaleProp returns the male proportion that should
// be used for this scheme category's Additional GLA Cover calc. The value is
// always derived from the uploaded member data — the benefit is only offered
// for populations whose member data has already been loaded. Falls back to
// group_pricing_parameters.main_member_male_prop (default 0.5) when the
// member list is empty or has no gender rows.
func ResolveAdditionalGlaCoverMaleProp(parameters models.GroupPricingParameters, members []models.GPricingMemberDataInForce) float64 {
	var male, total int
	for _, m := range members {
		if strings.TrimSpace(m.Gender) == "" {
			continue
		}
		total++
		g := strings.ToUpper(strings.TrimSpace(m.Gender))
		if g == "M" || g == "MALE" {
			male++
		}
	}
	if total == 0 {
		return parameters.MainMemberMaleProp
	}
	return float64(male) / float64(total)
}

func GetChildFuneralRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, childAge float64) float64 {
	tableName := "child_mortalities"

	var keyString strings.Builder
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(int(childAge)) + "_")
	key := keyString.String()
	cacheKey := tableName + "_" + key
	cached, found := GroupPricingCache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	} else {
		//fmt.Println("cache missed: ", key)
	}
	var childRate float64
	query := "risk_rate_code=? and age_next_birthday=?"
	err := DB.Table(tableName).Where(query, groupPricingParameter.RiskRateCode, int(childAge)).Pluck("child_rate", &childRate).Error
	if err != nil {
		fmt.Println(err)
	}
	GroupPricingCache.Set(cacheKey, childRate, 1)
	//time.Sleep(5 * time.Millisecond)
	return childRate
}

func GetDependantMortalityRate(memberResultData *models.MemberRatingResult, groupPricingParameter models.GroupPricingParameters, dependantAverageAge float64, incomeLevel int) float64 {
	tableName := "gla_rates"

	var keyString strings.Builder
	keyString.WriteString(groupPricingParameter.RiskRateCode + "_")
	keyString.WriteString(strconv.Itoa(int(dependantAverageAge)) + "_")
	keyString.WriteString(strconv.Itoa(incomeLevel) + "_") //incomelevel
	keyString.WriteString(memberResultData.Gender[:1] + "_")
	key := keyString.String()
	cacheKey := tableName + "_" + key
	cached, found := GroupPricingCache.Get(cacheKey)

	if found {
		result := cached.(float64)
		//if result > 0 {
		return result
		//}
	} else {
		//fmt.Println("cache missed: ", key)
	}
	var qx float64
	query := "risk_rate_code=? and age_next_birthday=? and income_level=? and gender=?"
	err := DB.Table(tableName).Where(query, groupPricingParameter.RiskRateCode, memberResultData.AgeNextBirthday, incomeLevel, memberResultData.Gender).Pluck("qx", &qx).Error
	if err != nil {
		fmt.Println(err)
	}
	GroupPricingCache.Set(cacheKey, qx, 1)
	//time.Sleep(5 * time.Millisecond)
	return qx
}

func CreateBroker(broker models.Broker, appUser models.AppUser) error {
	broker.CreatedBy = appUser.UserName
	result := DB.Where(models.Broker{ContactEmail: broker.ContactEmail}).FirstOrCreate(&broker)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

func GetBrokers() ([]models.Broker, error) {
	var brokers []models.Broker
	err := DB.Find(&brokers).Error
	if err != nil {
		return nil, err
	}
	return brokers, nil
}

func GetBroker(id string) (models.Broker, error) {
	var broker models.Broker
	err := DB.Where("id = ?", id).First(&broker).Error
	if err != nil {
		return broker, err
	}
	return broker, nil
}
func EditBroker(id string, payload models.Broker) (models.Broker, error) {
	// Ensure the claim exists
	var existing models.Broker
	if err := DB.First(&existing, id).Error; err != nil {
		return existing, err
	}

	// Force IDs to match and protect immutable fields
	payload.ID = existing.ID

	// Perform full-field update (including zero values) while omitting immutable columns
	if err := DB.Model(&existing).Select("*").Omit("id", "creation_date", "created_by").Updates(payload).Error; err != nil {
		return existing, err
	}
	// Return the refreshed Broker record
	var updated models.Broker
	if err := DB.First(&updated, id).Error; err != nil {
		return updated, err
	}
	return updated, nil
}

type DeleteError struct {
	Field   string
	Message string
}

func (e *DeleteError) Error() string {
	return fmt.Sprintf("Delete operation failed on %s: %s", e.Field, e.Message)
}
func DeleteBroker(id string) error {
	var activeQuotes []models.GroupPricingQuote
	var broker, err = GetBroker(id)
	if err != nil {
		return err
	}
	result := DB.Where("broker_name = ? AND status != ?", broker.Name, "expired").Find(&activeQuotes)
	//fmt.Printf("Broker affected was: %s.\n", broker.Name)
	//fmt.Printf("Rows affected was:   %d.\n", result.RowsAffected)
	if result.RowsAffected == 0 {
		err := DB.Where("id = ?", id).Delete(&models.Broker{}).Error
		if err != nil {
			return err
		}
	} else {
		//fmt.Println(activeQuotes[0].Status)
		//fmt.Println("Task failed successfully")
		return &DeleteError{Field: "Broker: " + broker.Name, Message: "Cannot delete brokers with active quotes."}
	}

	return nil
}

// BinderFee CRUD operations

func CreateBinderFee(fee models.BinderFee, appUser models.AppUser) error {
	fee.CreatedBy = appUser.UserName
	result := DB.Where(models.BinderFee{
		BinderholderName: fee.BinderholderName,
		RiskRateCode:     fee.RiskRateCode,
	}).FirstOrCreate(&fee)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

func GetBinderFees() ([]models.BinderFee, error) {
	var fees []models.BinderFee
	err := DB.Order("binderholder_name, risk_rate_code").Find(&fees).Error
	if err != nil {
		return nil, err
	}
	return fees, nil
}

func GetBinderFee(id string) (models.BinderFee, error) {
	var fee models.BinderFee
	err := DB.Where("id = ?", id).First(&fee).Error
	if err != nil {
		return fee, err
	}
	return fee, nil
}

func EditBinderFee(id string, payload models.BinderFee) (models.BinderFee, error) {
	var existing models.BinderFee
	if err := DB.First(&existing, id).Error; err != nil {
		return existing, err
	}

	payload.ID = existing.ID

	if err := DB.Model(&existing).Select("*").
		Omit("id", "creation_date", "created_by").
		Updates(payload).Error; err != nil {
		return existing, err
	}

	var updated models.BinderFee
	if err := DB.First(&updated, id).Error; err != nil {
		return updated, err
	}
	return updated, nil
}

func DeleteBinderFee(id string) error {
	return DB.Where("id = ?", id).Delete(&models.BinderFee{}).Error
}

// CommissionStructure CRUD + domain helpers

// sortBandsByLowerBound sorts in-place by lower_bound ascending.
func sortBandsByLowerBound(bands []models.CommissionStructure) {
	sort.Slice(bands, func(i, j int) bool {
		return bands[i].LowerBound < bands[j].LowerBound
	})
}

// holderLabel returns a human-readable label for error messages —
// "default" for empty holders, or the holder_name quoted.
func holderLabel(holder string) string {
	if holder == "" {
		return "default"
	}
	return fmt.Sprintf("%q", holder)
}

// ValidateCommissionBandsForChannelHolder returns nil when the current set of
// persisted bands for a given (channel, holder_name) group is contiguous
// from 0, non-overlapping, and each row has applicable_rate <=
// maximum_commission. Called after any Create/Edit/Delete so the
// invariant holds transactionally on a per-holder basis.
func ValidateCommissionBandsForChannelHolder(channel, holderName string) error {
	channel = strings.ToLower(strings.TrimSpace(channel))
	if channel == "" {
		return fmt.Errorf("channel is required")
	}
	if channel == "direct" {
		return fmt.Errorf("direct channel is always 0%% and cannot have bands")
	}

	var bands []models.CommissionStructure
	if err := DB.Where("channel = ? AND holder_name = ?", channel, holderName).Find(&bands).Error; err != nil {
		return err
	}
	return validateBandsSlice(channel, holderName, bands)
}

// validateBandsSlice is the pure validation used by both the exported
// validator and the transactional validator — keeps the rules in one place.
func validateBandsSlice(channel, holderName string, bands []models.CommissionStructure) error {
	if len(bands) == 0 {
		return nil
	}
	sortBandsByLowerBound(bands)
	for i, b := range bands {
		if b.LowerBound < 0 {
			return fmt.Errorf("band %d: lower_bound must be >= 0", b.ID)
		}
		if b.UpperBound != nil && *b.UpperBound <= b.LowerBound {
			return fmt.Errorf("band %d: upper_bound must be greater than lower_bound", b.ID)
		}
		if b.ApplicableRate < 0 || b.MaximumCommission < 0 {
			return fmt.Errorf("band %d: rates must be >= 0", b.ID)
		}
		if b.ApplicableRate > b.MaximumCommission {
			return fmt.Errorf("band %d: applicable_rate (%.6f) exceeds maximum_commission (%.6f)", b.ID, b.ApplicableRate, b.MaximumCommission)
		}
		if i == 0 && b.LowerBound != 0 {
			return fmt.Errorf("first band must start at lower_bound = 0 (channel %q holder %s starts at %.2f)", channel, holderLabel(holderName), b.LowerBound)
		}
		if i > 0 {
			prev := bands[i-1]
			if prev.UpperBound == nil {
				return fmt.Errorf("channel %q holder %s has a band after an unbounded band", channel, holderLabel(holderName))
			}
			if *prev.UpperBound != b.LowerBound {
				return fmt.Errorf("gap or overlap between bands (channel %q holder %s): previous ends at %.2f, next starts at %.2f", channel, holderLabel(holderName), *prev.UpperBound, b.LowerBound)
			}
		}
	}
	return nil
}

func CreateCommissionBand(band models.CommissionStructure, appUser models.AppUser) error {
	band.Channel = strings.ToLower(strings.TrimSpace(band.Channel))
	band.HolderName = strings.TrimSpace(band.HolderName)
	if band.Channel == "direct" {
		return fmt.Errorf("direct channel is always 0%% and cannot have bands")
	}
	band.CreatedBy = appUser.UserName
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&band).Error; err != nil {
			return err
		}
		if err := validateCommissionBandsForChannelHolderTx(tx, band.Channel, band.HolderName); err != nil {
			return err
		}
		return nil
	})
}

// GetCommissionBands returns commission bands filtered by channel and/or
// holder. holderFilterProvided distinguishes between "no filter" (show
// everything for the channel) and "filter to empty holder" (show only
// defaults). When allHolders is true, channel rows across every holder
// are returned regardless of the holderFilterProvided flag.
func GetCommissionBands(channel, holderName string, holderFilterProvided, allHolders bool) ([]models.CommissionStructure, error) {
	var bands []models.CommissionStructure
	q := DB.Order("channel, holder_name, lower_bound")
	if channel != "" {
		q = q.Where("channel = ?", strings.ToLower(strings.TrimSpace(channel)))
	}
	if !allHolders && holderFilterProvided {
		q = q.Where("holder_name = ?", strings.TrimSpace(holderName))
	}
	if err := q.Find(&bands).Error; err != nil {
		return nil, err
	}
	return bands, nil
}

func GetCommissionBand(id string) (models.CommissionStructure, error) {
	var band models.CommissionStructure
	if err := DB.Where("id = ?", id).First(&band).Error; err != nil {
		return band, err
	}
	return band, nil
}

func EditCommissionBand(id string, payload models.CommissionStructure) (models.CommissionStructure, error) {
	payload.Channel = strings.ToLower(strings.TrimSpace(payload.Channel))
	payload.HolderName = strings.TrimSpace(payload.HolderName)
	if payload.Channel == "direct" {
		return payload, fmt.Errorf("direct channel is always 0%% and cannot have bands")
	}

	var existing models.CommissionStructure
	if err := DB.First(&existing, id).Error; err != nil {
		return existing, err
	}

	payload.ID = existing.ID

	if err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&existing).Select("*").
			Omit("id", "creation_date", "created_by").
			Updates(payload).Error; err != nil {
			return err
		}
		// Validate the new group.
		if err := validateCommissionBandsForChannelHolderTx(tx, payload.Channel, payload.HolderName); err != nil {
			return err
		}
		// If the group changed (channel or holder moved), validate the old
		// one too — the row we just moved out may have left a gap there.
		if existing.Channel != payload.Channel || existing.HolderName != payload.HolderName {
			if err := validateCommissionBandsForChannelHolderTx(tx, existing.Channel, existing.HolderName); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return existing, err
	}

	var updated models.CommissionStructure
	if err := DB.First(&updated, id).Error; err != nil {
		return updated, err
	}
	return updated, nil
}

func DeleteCommissionBand(id string) error {
	var existing models.CommissionStructure
	if err := DB.First(&existing, id).Error; err != nil {
		return err
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&models.CommissionStructure{}).Error; err != nil {
			return err
		}
		// After delete, the remaining bands in the (channel, holder) group
		// may leave a gap. We validate but DO NOT hard-fail the delete —
		// operators may be removing bands as a first step of a reshape.
		// Leave the warning to the frontend contiguity check and the next
		// Create/Edit call.
		_ = validateCommissionBandsForChannelHolderTx(tx, existing.Channel, existing.HolderName)
		return nil
	})
}

// validateCommissionBandsForChannelHolderTx is the transactional variant
// used inside Create/Edit flows — validates one (channel, holder) group.
func validateCommissionBandsForChannelHolderTx(tx *gorm.DB, channel, holderName string) error {
	channel = strings.ToLower(strings.TrimSpace(channel))
	if channel == "" || channel == "direct" {
		return nil
	}
	var bands []models.CommissionStructure
	if err := tx.Where("channel = ? AND holder_name = ?", channel, holderName).Find(&bands).Error; err != nil {
		return err
	}
	return validateBandsSlice(channel, holderName, bands)
}

// ComputeProgressiveCommission returns the blended commission loading for
// the given annualised premium under a marginal (progressive) application
// of a (channel, holder)'s bands. Lookup order:
//  1. rows matching (channel, holderName) — use them if any exist.
//  2. else rows matching (channel, ""), i.e. the channel default — use them.
//  3. else return 0.
//
// Returns 0 for channel="direct" or when no bands are configured.
//
// The return value is a decimal fraction (e.g. 0.053 for 5.3%).
func ComputeProgressiveCommission(channel, holderName string, annualPremium float64) (float64, error) {
	channel = strings.ToLower(strings.TrimSpace(channel))
	holderName = strings.TrimSpace(holderName)
	if channel == "direct" || annualPremium <= 0 {
		return 0, nil
	}
	// 1. Holder-specific.
	if holderName != "" {
		var bands []models.CommissionStructure
		if err := DB.Where("channel = ? AND holder_name = ?", channel, holderName).Find(&bands).Error; err != nil {
			return 0, err
		}
		if len(bands) > 0 {
			return computeProgressiveCommissionFromBands(bands, annualPremium), nil
		}
	}
	// 2. Channel default.
	var bands []models.CommissionStructure
	if err := DB.Where("channel = ? AND holder_name = ?", channel, "").Find(&bands).Error; err != nil {
		return 0, err
	}
	return computeProgressiveCommissionFromBands(bands, annualPremium), nil
}

// computeProgressiveCommissionFromBands is the pure math used by
// ComputeProgressiveCommission — no DB access, safe to unit-test in
// isolation. Returns 0 for empty bands or non-positive premium.
func computeProgressiveCommissionFromBands(bands []models.CommissionStructure, annualPremium float64) float64 {
	if annualPremium <= 0 || len(bands) == 0 {
		return 0
	}
	sortBandsByLowerBound(bands)
	weighted := 0.0
	for _, b := range bands {
		if annualPremium <= b.LowerBound {
			break
		}
		upper := annualPremium
		if b.UpperBound != nil && *b.UpperBound < upper {
			upper = *b.UpperBound
		}
		segment := upper - b.LowerBound
		if segment <= 0 {
			continue
		}
		weighted += segment * b.ApplicableRate
	}
	return weighted / annualPremium
}

// Reinsurer CRUD operations

func CreateReinsurer(reinsurer models.Reinsurer, appUser models.AppUser) error {
	reinsurer.CreatedBy = appUser.UserName
	result := DB.Where(models.Reinsurer{Code: reinsurer.Code}).FirstOrCreate(&reinsurer)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

func GetReinsurers() ([]models.Reinsurer, error) {
	var reinsurers []models.Reinsurer
	// Only return active reinsurers
	err := DB.Where("is_active = ?", true).Order("name ASC").Find(&reinsurers).Error
	if err != nil {
		return nil, err
	}
	return reinsurers, nil
}

func GetReinsurer(id string) (models.Reinsurer, error) {
	var reinsurer models.Reinsurer
	err := DB.Where("id = ?", id).First(&reinsurer).Error
	if err != nil {
		return reinsurer, err
	}
	return reinsurer, nil
}

func EditReinsurer(id string, payload models.Reinsurer) (models.Reinsurer, error) {
	// Ensure the reinsurer exists
	var existing models.Reinsurer
	if err := DB.First(&existing, id).Error; err != nil {
		return existing, err
	}

	// Force IDs to match and protect immutable fields
	payload.ID = existing.ID

	// Perform full-field update (including zero values) while omitting immutable columns
	if err := DB.Model(&existing).Select("*").Omit("id", "creation_date", "created_by").Updates(payload).Error; err != nil {
		return existing, err
	}

	// Return the refreshed Reinsurer record
	var updated models.Reinsurer
	if err := DB.First(&updated, id).Error; err != nil {
		return updated, err
	}
	return updated, nil
}

func DeleteReinsurer(id string) error {
	err := DB.Where("id = ?", id).Delete(&models.Reinsurer{}).Error
	if err != nil {
		return err
	}
	return nil
}

// DeactivateReinsurer deactivates a reinsurer if not involved in any active treaties
func DeactivateReinsurer(id string, reason string, appUser models.AppUser) error {
	// First check if reinsurer exists and is active
	var reinsurer models.Reinsurer
	if err := DB.First(&reinsurer, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("reinsurer not found")
		}
		return err
	}

	if !reinsurer.IsActive {
		return fmt.Errorf("reinsurer is already deactivated")
	}

	// Check if reinsurer is involved in any active treaties
	var activeTreatyCount int64
	err := DB.Model(&models.ReinsuranceTreaty{}).
		Where("reinsurer_name = ? AND status IN (?)", reinsurer.Name, []string{"active", "draft", "under_negotiation"}).
		Count(&activeTreatyCount).Error

	if err != nil {
		return fmt.Errorf("failed to check active treaties: %w", err)
	}

	if activeTreatyCount > 0 {
		return fmt.Errorf("cannot deactivate reinsurer: involved in %d active treaty/treaties", activeTreatyCount)
	}

	// Deactivate the reinsurer
	now := time.Now()
	updates := map[string]interface{}{
		"is_active":           false,
		"deactivated_at":      now,
		"deactivated_by":      appUser.UserName,
		"deactivation_reason": reason,
	}

	if err := DB.Model(&reinsurer).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to deactivate reinsurer: %w", err)
	}

	return nil
}

func CreateGroupScheme(groupScheme models.GroupScheme, appUser models.AppUser) error {
	groupScheme.CreatedBy = appUser.UserName
	err := DB.Create(&groupScheme).Error
	if err != nil {
		return err
	}
	return nil
}

func CheckGroupSchemeName(name string) bool {
	var groupScheme models.GroupScheme
	err := DB.Where("name = ?", name).First(&groupScheme).Error
	if err != nil {
		return false
	}
	return true
}

func GetAllGroupSchemes() ([]models.GroupScheme, error) {
	var groupSchemes []models.GroupScheme
	err := DB.Preload("Quote").Find(&groupSchemes).Error
	if err != nil {
		return nil, err
	}
	hydrateSchemesTreatyLinkFlag(groupSchemes)
	return groupSchemes, nil
}

func GetGroupSchemesInforce() ([]models.GroupScheme, error) {
	var groupSchemes []models.GroupScheme
	err := DB.Where("status = ? ", models.StatusInForce).Find(&groupSchemes).Error
	if err != nil {
		return nil, err
	}
	hydrateSchemesTreatyLinkFlag(groupSchemes)
	return groupSchemes, nil
}

// claimsAggForITD is the per-scheme claims roll-up used to compute ITD ALR.
// Status set is {approved, paid} — both represent claims the insurer has
// committed to on the book. Declined / withdrawn / pending are excluded.
type claimsAggForITD struct {
	SchemeID      int     `gorm:"column:scheme_id"`
	TotalAmount   float64 `gorm:"column:total_amount"`
	ClaimsCount   int     `gorm:"column:claims_count"`
	LastClaimDate string  `gorm:"column:last_claim_date"`
}

// computeSchemePerformanceRows builds per-scheme rows for the in-force book.
// Shared by GetSchemePerformanceRows and GetInForceRiskProfile so that the
// portfolio-level KPIs and the table use exactly the same source of truth.
func computeSchemePerformanceRows() ([]models.SchemePerformanceRow, error) {
	var schemes []models.GroupScheme
	if err := DB.Where("in_force = ? OR status = ?", true, models.StatusInForce).
		Find(&schemes).Error; err != nil {
		return nil, err
	}
	if len(schemes) == 0 {
		return []models.SchemePerformanceRow{}, nil
	}

	schemeIDs := make([]int, len(schemes))
	for i, s := range schemes {
		schemeIDs[i] = s.ID
	}

	var aggs []claimsAggForITD
	if err := DB.Table("group_scheme_claims").
		Select(`scheme_id,
			COALESCE(SUM(claim_amount), 0) AS total_amount,
			COUNT(*) AS claims_count,
			COALESCE(MAX(date_of_event), '') AS last_claim_date`).
		Where("scheme_id IN ? AND status IN ?", schemeIDs, []string{"approved", "paid"}).
		Group("scheme_id").
		Scan(&aggs).Error; err != nil {
		return nil, err
	}
	aggBySchemeID := make(map[int]claimsAggForITD, len(aggs))
	for _, a := range aggs {
		aggBySchemeID[a.SchemeID] = a
	}

	// Rolling-12-month claims per scheme — date_of_event in the last 12
	// months, same status filter as ITD. date_of_event is a string
	// (YYYY-MM-DD) so we compare lexicographically against today minus 12.
	r12mCutoff := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	var r12mAggs []claimsAggForITD
	if err := DB.Table("group_scheme_claims").
		Select(`scheme_id,
			COALESCE(SUM(claim_amount), 0) AS total_amount,
			COUNT(*) AS claims_count`).
		Where("scheme_id IN ? AND status IN ? AND date_of_event >= ?",
			schemeIDs, []string{"approved", "paid"}, r12mCutoff).
		Group("scheme_id").
		Scan(&r12mAggs).Error; err != nil {
		return nil, err
	}
	r12mBySchemeID := make(map[int]claimsAggForITD, len(r12mAggs))
	for _, a := range r12mAggs {
		r12mBySchemeID[a.SchemeID] = a
	}

	now := time.Now()
	rows := make([]models.SchemePerformanceRow, 0, len(schemes))
	for i := range schemes {
		s := schemes[i]

		var coverStart *time.Time
		if !s.CoverStartDate.IsZero() {
			coverStart = &s.CoverStartDate
		}

		startForMonths := s.CoverStartDate
		if startForMonths.IsZero() {
			startForMonths = s.CreationDate
		}
		months := 0
		if !startForMonths.IsZero() && now.After(startForMonths) {
			months = int(now.Sub(startForMonths).Hours()/(24*30)) + 1
			if months < 1 {
				months = 1
			}
		}

		itdEarned := s.EarnedPremium
		if itdEarned <= 0 && s.AnnualPremium > 0 && months > 0 {
			itdEarned = s.AnnualPremium * float64(months) / 12.0
		}

		commissionPct := 0.0
		if s.AnnualPremium > 0 {
			commissionPct = FloatPrecision(s.Commission/s.AnnualPremium*100, 2)
		}

		agg := aggBySchemeID[s.ID]

		var avgSeverity *float64
		if agg.ClaimsCount > 0 {
			v := FloatPrecision(agg.TotalAmount/float64(agg.ClaimsCount), 2)
			avgSeverity = &v
		}

		var freq *float64
		if s.MemberCount > 0 {
			v := FloatPrecision(float64(agg.ClaimsCount)/s.MemberCount*1000, 2)
			freq = &v
		}

		var alr, delta *float64
		if itdEarned > 0 {
			v := FloatPrecision(agg.TotalAmount/itdEarned*100, 1)
			alr = &v
			d := FloatPrecision(v-s.ExpectedLossRatio, 1)
			delta = &d
		}

		// Rolling-12-month ALR. Denominator = annual_premium when scheme
		// has been in force the full window, otherwise pro-rated by
		// months_in_force / 12 — same shape as ITD earned premium logic.
		r12m := r12mBySchemeID[s.ID]
		var r12mAlr *float64
		var r12mDenom float64
		if s.AnnualPremium > 0 {
			if months >= 12 {
				r12mDenom = s.AnnualPremium
			} else if months > 0 {
				r12mDenom = s.AnnualPremium * float64(months) / 12.0
			}
		}
		if r12mDenom > 0 {
			v := FloatPrecision(r12m.TotalAmount/r12mDenom*100, 1)
			r12mAlr = &v
		}

		rows = append(rows, models.SchemePerformanceRow{
			SchemeID:          s.ID,
			SchemeName:        s.Name,
			Status:            string(s.Status),
			CoverStartDate:    coverStart,
			MonthsInForce:     months,
			MemberCount:       s.MemberCount,
			AnnualPremium:     FloatPrecision(s.AnnualPremium, 2),
			EarnedPremium:     FloatPrecision(s.EarnedPremium, 2),
			ITDEarnedPremium:  FloatPrecision(itdEarned, 2),
			CommissionPct:     commissionPct,
			ItdClaimsPaid:     FloatPrecision(agg.TotalAmount, 2),
			ItdClaimsCount:    agg.ClaimsCount,
			AvgClaimSeverity:  avgSeverity,
			ClaimsFrequency:   freq,
			ExpectedLossRatio: FloatPrecision(s.ExpectedLossRatio, 1),
			ActualLossRatio:   alr,
			LossRatioDelta:    delta,
			R12mClaimsPaid:    FloatPrecision(r12m.TotalAmount, 2),
			R12mClaimsCount:   r12m.ClaimsCount,
			R12mAlr:           r12mAlr,
			LastClaimDate:     agg.LastClaimDate,
		})
	}
	return rows, nil
}

// GetSchemePerformanceRows returns the per-scheme performance table along
// with portfolio-level totals (used by the headline KPI cards).
func GetSchemePerformanceRows() (models.SchemePerformanceResponse, error) {
	rows, err := computeSchemePerformanceRows()
	if err != nil {
		return models.SchemePerformanceResponse{}, err
	}
	resp := models.SchemePerformanceResponse{
		Rows:         rows,
		TotalSchemes: len(rows),
	}
	var totalEarned, totalR12mDenom float64
	for _, r := range rows {
		resp.TotalPremium += r.AnnualPremium
		resp.TotalItdClaims += r.ItdClaimsPaid
		totalEarned += r.ITDEarnedPremium
		resp.TotalR12mClaims += r.R12mClaimsPaid
		// Match the per-scheme R12M denominator pro-rating so the portfolio
		// ratio is internally consistent with the per-scheme rows.
		if r.AnnualPremium > 0 {
			if r.MonthsInForce >= 12 {
				totalR12mDenom += r.AnnualPremium
			} else if r.MonthsInForce > 0 {
				totalR12mDenom += r.AnnualPremium * float64(r.MonthsInForce) / 12.0
			}
		}
	}
	resp.TotalPremium = FloatPrecision(resp.TotalPremium, 2)
	resp.TotalItdClaims = FloatPrecision(resp.TotalItdClaims, 2)
	resp.TotalR12mClaims = FloatPrecision(resp.TotalR12mClaims, 2)
	if totalEarned > 0 {
		v := FloatPrecision(resp.TotalItdClaims/totalEarned*100, 1)
		resp.PortfolioALR = &v
	}
	if totalR12mDenom > 0 {
		v := FloatPrecision(resp.TotalR12mClaims/totalR12mDenom*100, 1)
		resp.PortfolioR12mALR = &v
	}
	return resp, nil
}

// GetInForceRiskProfile derives all risk-profile signals from the per-scheme
// rows: concentration KPIs, Pareto curve, loss-ratio distribution, top-10
// worst, frequency-severity scatter, industry/region heatmap, and the
// deteriorating-schemes panel.
func GetInForceRiskProfile() (models.RiskProfileResult, error) {
	result := models.RiskProfileResult{
		Pareto:            []models.ParetoPoint{},
		LossRatioBuckets:  []models.LossRatioBucket{},
		Top10Worst:        []models.SchemePerformanceRow{},
		IndustryRegion:    []models.IndustryRegionCell{},
		FrequencySeverity: []models.FreqSeverityPoint{},
		Deteriorating:     []models.DeterioratingScheme{},
	}

	// Resolve company-level deteriorating-scheme thresholds from the
	// singleton settings row. Falls back to the documented defaults if the
	// row is missing or has zeros (e.g. fresh DB before settings have been
	// touched).
	alrCeiling := 100.0
	alrDelta := 20.0
	var settings models.GroupPricingSetting
	if err := DB.First(&settings, 1).Error; err == nil {
		if settings.RiskAlrCeilingPct > 0 {
			alrCeiling = settings.RiskAlrCeilingPct
		}
		if settings.RiskAlrDeltaPp > 0 {
			alrDelta = settings.RiskAlrDeltaPp
		}
	}
	result.Thresholds = models.RiskWatchlistThresholds{
		AlrCeilingPct: alrCeiling,
		AlrDeltaPp:    alrDelta,
	}

	rows, err := computeSchemePerformanceRows()
	if err != nil {
		return result, err
	}
	result.Concentration.TotalSchemes = len(rows)
	if len(rows) == 0 {
		result.LossRatioBuckets = makeEmptyLossRatioBuckets()
		return result, nil
	}

	// ── Concentration KPIs and Pareto ─────────────────────────────────────
	byPremium := make([]models.SchemePerformanceRow, len(rows))
	copy(byPremium, rows)
	sort.SliceStable(byPremium, func(i, j int) bool {
		return byPremium[i].AnnualPremium > byPremium[j].AnnualPremium
	})
	var totalPremium float64
	for _, r := range byPremium {
		totalPremium += r.AnnualPremium
	}
	if totalPremium > 0 {
		var top5, top10, hhi, cum float64
		paretoLimit := len(byPremium)
		if paretoLimit > 50 {
			paretoLimit = 50
		}
		result.Pareto = make([]models.ParetoPoint, 0, paretoLimit)
		for i, r := range byPremium {
			share := r.AnnualPremium / totalPremium
			if i < 5 {
				top5 += share
			}
			if i < 10 {
				top10 += share
			}
			hhi += share * share
			cum += share
			if i < paretoLimit {
				result.Pareto = append(result.Pareto, models.ParetoPoint{
					Rank:            i + 1,
					SchemeName:      r.SchemeName,
					Premium:         r.AnnualPremium,
					CumulativeShare: FloatPrecision(cum, 4),
				})
			}
		}
		result.Concentration.Top5PremiumShare = FloatPrecision(top5*100, 2)
		result.Concentration.Top10PremiumShare = FloatPrecision(top10*100, 2)
		result.Concentration.HHI = FloatPrecision(hhi*10000, 1)
	}

	// ── Loss-ratio distribution buckets ───────────────────────────────────
	buckets := makeEmptyLossRatioBuckets()
	for _, r := range rows {
		if r.ActualLossRatio == nil {
			continue
		}
		alr := *r.ActualLossRatio
		idx := bucketIndexForALR(alr)
		buckets[idx].SchemeCount++
		buckets[idx].Premium = FloatPrecision(buckets[idx].Premium+r.AnnualPremium, 2)
	}
	result.LossRatioBuckets = buckets

	// ── Top 10 worst ──────────────────────────────────────────────────────
	withALR := make([]models.SchemePerformanceRow, 0, len(rows))
	for _, r := range rows {
		if r.ActualLossRatio != nil {
			withALR = append(withALR, r)
		}
	}
	sort.SliceStable(withALR, func(i, j int) bool {
		if *withALR[i].ActualLossRatio == *withALR[j].ActualLossRatio {
			return withALR[i].ItdClaimsPaid > withALR[j].ItdClaimsPaid
		}
		return *withALR[i].ActualLossRatio > *withALR[j].ActualLossRatio
	})
	if len(withALR) > 10 {
		result.Top10Worst = withALR[:10]
	} else {
		result.Top10Worst = withALR
	}

	// ── Frequency / severity scatter ──────────────────────────────────────
	for _, r := range rows {
		if r.ClaimsFrequency == nil || r.AvgClaimSeverity == nil {
			continue
		}
		result.FrequencySeverity = append(result.FrequencySeverity, models.FreqSeverityPoint{
			SchemeID:      r.SchemeID,
			SchemeName:    r.SchemeName,
			Frequency:     *r.ClaimsFrequency,
			AvgSeverity:   *r.AvgClaimSeverity,
			AnnualPremium: r.AnnualPremium,
		})
	}

	// ── Deteriorating schemes ─────────────────────────────────────────────
	ceilingLabel := strconv.FormatFloat(alrCeiling, 'f', -1, 64)
	deltaLabel := strconv.FormatFloat(alrDelta, 'f', -1, 64)
	for _, r := range rows {
		reasons := []string{}
		if r.ActualLossRatio != nil && *r.ActualLossRatio > alrCeiling {
			reasons = append(reasons, "ALR > "+ceilingLabel+"%")
		}
		if r.LossRatioDelta != nil && *r.LossRatioDelta > alrDelta {
			reasons = append(reasons, "ALR exceeds ELR by >"+deltaLabel+"pp")
		}
		if len(reasons) == 0 {
			continue
		}
		result.Deteriorating = append(result.Deteriorating, models.DeterioratingScheme{
			SchemeID:          r.SchemeID,
			SchemeName:        r.SchemeName,
			TriggerReasons:    reasons,
			ExpectedLossRatio: r.ExpectedLossRatio,
			ActualLossRatio:   r.ActualLossRatio,
			LossRatioDelta:    r.LossRatioDelta,
			ItdClaimsPaid:     r.ItdClaimsPaid,
			LastClaimDate:     r.LastClaimDate,
		})
	}
	sort.SliceStable(result.Deteriorating, func(i, j int) bool {
		ai := safeFloatPtr(result.Deteriorating[i].ActualLossRatio)
		aj := safeFloatPtr(result.Deteriorating[j].ActualLossRatio)
		return ai > aj
	})

	// ── Industry / region heatmap ─────────────────────────────────────────
	heatmap, hmErr := computeIndustryRegionHeatmap(rows)
	if hmErr == nil {
		result.IndustryRegion = heatmap
	}

	return result, nil
}

// trendWindowMonths controls how many trailing months are returned by the
// loss-ratio trend endpoint. 24 months gives the user enough history to
// distinguish a steady scheme from one that's deteriorating without making
// the chart unreadable.
const trendWindowMonths = 24

// GetLossRatioTrend returns rolling-12-month ALR per in-force scheme over the
// last `trendWindowMonths` months, plus a portfolio aggregate line. Computed
// from monthly per-scheme claims sums + each scheme's annual premium so the
// frontend has everything it needs to render line charts.
func GetLossRatioTrend() (models.LossRatioTrendResult, error) {
	result := models.LossRatioTrendResult{
		Months:    []string{},
		Portfolio: []*float64{},
		Schemes:   []models.LossRatioTrendSeries{},
	}

	var schemes []models.GroupScheme
	if err := DB.Where("in_force = ? OR status = ?", true, models.StatusInForce).
		Find(&schemes).Error; err != nil {
		return result, err
	}
	if len(schemes) == 0 {
		return result, nil
	}
	schemeIDs := make([]int, len(schemes))
	for i, s := range schemes {
		schemeIDs[i] = s.ID
	}

	// Build the display month axis (oldest → newest) and a lookback table to
	// collect per-scheme monthly claims back to (display oldest − 11 months)
	// so we can compute rolling-12 sums for the very first display point.
	now := time.Now()
	displayStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).
		AddDate(0, -(trendWindowMonths - 1), 0)
	dataStart := displayStart.AddDate(0, -11, 0)
	cutoff := dataStart.Format("2006-01-02")

	// Pull monthly per-scheme sums in a single query. date_of_event is
	// stored as a YYYY-MM-DD string so we extract YYYY-MM with LEFT() —
	// portable across the project's three supported dialects.
	type monthRow struct {
		SchemeID int     `gorm:"column:scheme_id"`
		YM       string  `gorm:"column:ym"`
		Total    float64 `gorm:"column:total"`
	}
	var monthRows []monthRow
	if err := DB.Table("group_scheme_claims").
		Select("scheme_id, LEFT(date_of_event, 7) AS ym, COALESCE(SUM(claim_amount), 0) AS total").
		Where("scheme_id IN ? AND status IN ? AND date_of_event >= ?",
			schemeIDs, []string{"approved", "paid"}, cutoff).
		Group("scheme_id, LEFT(date_of_event, 7)").
		Scan(&monthRows).Error; err != nil {
		return result, err
	}

	// claimsByScheme[schemeID][ym] = sum
	claimsByScheme := make(map[int]map[string]float64, len(schemes))
	for _, mr := range monthRows {
		if _, ok := claimsByScheme[mr.SchemeID]; !ok {
			claimsByScheme[mr.SchemeID] = map[string]float64{}
		}
		claimsByScheme[mr.SchemeID][mr.YM] = mr.Total
	}

	// Display months: trendWindowMonths entries, formatted YYYY-MM.
	displayMonths := make([]string, trendWindowMonths)
	displayCursor := displayStart
	for i := 0; i < trendWindowMonths; i++ {
		displayMonths[i] = displayCursor.Format("2006-01")
		displayCursor = displayCursor.AddDate(0, 1, 0)
	}
	result.Months = displayMonths

	// Look-up months: 12 months prior up to the current display month, used
	// for the rolling sum at each display point.
	lookbackKey := func(end time.Time, offset int) string {
		return end.AddDate(0, -offset, 0).Format("2006-01")
	}

	// Sort schemes by current ITD claims as a cheap proxy for "biggest
	// signal" — caller can pick top N from this order if it wants a default.
	sort.SliceStable(schemes, func(i, j int) bool {
		return schemes[i].AnnualPremium > schemes[j].AnnualPremium
	})

	// Compute per-scheme rolling-12 ALR series, plus accumulate portfolio
	// numerator and denominator per display month.
	portNum := make([]float64, trendWindowMonths)
	portDen := make([]float64, trendWindowMonths)
	for _, s := range schemes {
		months := claimsByScheme[s.ID]
		series := models.LossRatioTrendSeries{
			SchemeID:      s.ID,
			SchemeName:    s.Name,
			AnnualPremium: FloatPrecision(s.AnnualPremium, 2),
			Values:        make([]*float64, trendWindowMonths),
		}

		startForMonths := s.CoverStartDate
		if startForMonths.IsZero() {
			startForMonths = s.CreationDate
		}

		for i := 0; i < trendWindowMonths; i++ {
			endOfWindow := displayStart.AddDate(0, i, 0).AddDate(0, 1, 0).Add(-time.Hour)
			// Skip points before the scheme existed.
			if !startForMonths.IsZero() && endOfWindow.Before(startForMonths) {
				continue
			}

			// Sum claims in trailing 12 months ending in display-month i.
			var num float64
			pivot := displayStart.AddDate(0, i, 0)
			for k := 0; k < 12; k++ {
				num += months[lookbackKey(pivot, k)]
			}

			// Time-weighted denominator: months in force as of pivot,
			// capped at 12.
			minf := 0
			if !startForMonths.IsZero() {
				diff := endOfWindow.Sub(startForMonths).Hours() / (24 * 30)
				if diff > 0 {
					minf = int(diff) + 1
				}
			}
			if minf > 12 {
				minf = 12
			}
			denom := 0.0
			if s.AnnualPremium > 0 && minf > 0 {
				denom = s.AnnualPremium * float64(minf) / 12.0
			}
			if denom > 0 {
				v := FloatPrecision(num/denom*100, 1)
				series.Values[i] = &v
				portNum[i] += num
				portDen[i] += denom
			}
		}

		// Latest non-nil value is the "current" R12M ALR for this scheme —
		// useful for the UI to default to top N by current value.
		for i := trendWindowMonths - 1; i >= 0; i-- {
			if series.Values[i] != nil {
				v := *series.Values[i]
				series.CurrentR12mALR = &v
				break
			}
		}

		result.Schemes = append(result.Schemes, series)
	}

	result.Portfolio = make([]*float64, trendWindowMonths)
	for i := 0; i < trendWindowMonths; i++ {
		if portDen[i] > 0 {
			v := FloatPrecision(portNum[i]/portDen[i]*100, 1)
			result.Portfolio[i] = &v
		}
	}

	return result, nil
}

func makeEmptyLossRatioBuckets() []models.LossRatioBucket {
	return []models.LossRatioBucket{
		{Label: "<50%", LowerBound: 0, UpperBound: 50},
		{Label: "50-80%", LowerBound: 50, UpperBound: 80},
		{Label: "80-100%", LowerBound: 80, UpperBound: 100},
		{Label: "100-150%", LowerBound: 100, UpperBound: 150},
		{Label: ">150%", LowerBound: 150, UpperBound: -1},
	}
}

func bucketIndexForALR(alr float64) int {
	switch {
	case alr < 50:
		return 0
	case alr < 80:
		return 1
	case alr < 100:
		return 2
	case alr < 150:
		return 3
	default:
		return 4
	}
}

func safeFloatPtr(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

// computeIndustryRegionHeatmap maps each in-force scheme to its (industry,
// region) pair — industry from the parent group_pricing_quotes row, region
// from the scheme's first scheme_category — and aggregates premium and ITD
// claims by that pair.
func computeIndustryRegionHeatmap(rows []models.SchemePerformanceRow) ([]models.IndustryRegionCell, error) {
	if len(rows) == 0 {
		return []models.IndustryRegionCell{}, nil
	}
	rowBySchemeID := make(map[int]models.SchemePerformanceRow, len(rows))
	schemeIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		rowBySchemeID[r.SchemeID] = r
		schemeIDs = append(schemeIDs, r.SchemeID)
	}

	type schemeIndustryRow struct {
		SchemeID int    `gorm:"column:scheme_id"`
		Industry string `gorm:"column:industry"`
	}
	var inds []schemeIndustryRow
	if err := DB.Table("group_schemes gs").
		Select("gs.id AS scheme_id, COALESCE(q.industry, '') AS industry").
		Joins("LEFT JOIN group_pricing_quotes q ON q.id = gs.quote_id").
		Where("gs.id IN ?", schemeIDs).
		Scan(&inds).Error; err != nil {
		return nil, err
	}
	industryBySchemeID := make(map[int]string, len(inds))
	for _, ir := range inds {
		industryBySchemeID[ir.SchemeID] = ir.Industry
	}

	// Region: per scheme, take the lowest-id scheme_category's region (a
	// pragmatic V1 simplification — schemes that span multiple regions get
	// represented by their first category).
	type schemeRegionRow struct {
		SchemeID int    `gorm:"column:scheme_id"`
		Region   string `gorm:"column:region"`
	}
	var regs []schemeRegionRow
	if err := DB.Table("group_schemes gs").
		Select("gs.id AS scheme_id, COALESCE(MIN(sc.region), '') AS region").
		Joins("LEFT JOIN scheme_categories sc ON sc.quote_id = gs.quote_id").
		Where("gs.id IN ?", schemeIDs).
		Group("gs.id").
		Scan(&regs).Error; err != nil {
		return nil, err
	}
	regionBySchemeID := make(map[int]string, len(regs))
	for _, rr := range regs {
		regionBySchemeID[rr.SchemeID] = rr.Region
	}

	// Build (industry, region) pairs per scheme.
	type schemeCatRow struct {
		SchemeID int
		Industry string
		Region   string
	}
	cats := make([]schemeCatRow, 0, len(rows))
	for _, r := range rows {
		cats = append(cats, schemeCatRow{
			SchemeID: r.SchemeID,
			Industry: industryBySchemeID[r.SchemeID],
			Region:   regionBySchemeID[r.SchemeID],
		})
	}

	type cellKey struct {
		Industry string
		Region   string
	}
	cellMap := make(map[cellKey]*models.IndustryRegionCell)
	for _, c := range cats {
		ind := c.Industry
		if ind == "" {
			ind = "Unknown"
		}
		reg := c.Region
		if reg == "" {
			reg = "Unknown"
		}
		k := cellKey{Industry: ind, Region: reg}
		row := rowBySchemeID[c.SchemeID]
		cell, ok := cellMap[k]
		if !ok {
			cell = &models.IndustryRegionCell{Industry: ind, Region: reg}
			cellMap[k] = cell
		}
		cell.Premium += row.AnnualPremium
		cell.ClaimsPaid += row.ItdClaimsPaid
		cell.MemberCount += int64(row.MemberCount)
	}

	cells := make([]models.IndustryRegionCell, 0, len(cellMap))
	for _, c := range cellMap {
		c.Premium = FloatPrecision(c.Premium, 2)
		c.ClaimsPaid = FloatPrecision(c.ClaimsPaid, 2)
		if c.Premium > 0 {
			c.LossRatio = FloatPrecision(c.ClaimsPaid/c.Premium*100, 1)
		}
		cells = append(cells, *c)
	}
	sort.SliceStable(cells, func(i, j int) bool {
		if cells[i].Industry == cells[j].Industry {
			return cells[i].Region < cells[j].Region
		}
		return cells[i].Industry < cells[j].Industry
	})
	return cells, nil
}

// hydrateSchemesTreatyLinkFlag sets HasTreatyLink on each scheme by checking
// treaty_scheme_links for linked scheme IDs in a single query.
func hydrateSchemesTreatyLinkFlag(schemes []models.GroupScheme) {
	if len(schemes) == 0 {
		return
	}
	schemeIDs := make([]int, len(schemes))
	for i, s := range schemes {
		schemeIDs[i] = s.ID
	}
	var linkedIDs []int
	DB.Model(&models.TreatySchemeLink{}).Where("scheme_id IN ?", schemeIDs).
		Distinct("scheme_id").Pluck("scheme_id", &linkedIDs)
	linkedSet := make(map[int]bool, len(linkedIDs))
	for _, id := range linkedIDs {
		linkedSet[id] = true
	}
	for i := range schemes {
		schemes[i].HasTreatyLink = linkedSet[schemes[i].ID]
	}
}

func GetGroupScheme(id string) (models.GroupScheme, error) {
	var groupScheme models.GroupScheme
	err := DB.Where("id = ?", id).First(&groupScheme).Error
	if err != nil {
		return groupScheme, err
	}
	var claims []models.GroupSchemeClaim
	err = DB.Where("scheme_id = ?", id).Find(&claims).Error

	var totalClaims float64 = 0

	for _, claim := range claims {
		if claim.Status == APPROVED {
			totalClaims += claim.ClaimAmount
		}
	}

	groupScheme.ActualClaims = totalClaims

	if groupScheme.BrokerId > 0 {
		var broker models.Broker
		err = DB.Where("id = ?", groupScheme.BrokerId).First(&broker).Error
		if err != nil {
			return groupScheme, err
		}
		groupScheme.Broker = broker
	}

	var quote models.GroupPricingQuote
	err = DB.Where("id = ?", groupScheme.QuoteId).First(&quote).Error
	if err != nil {
		return groupScheme, err
	}
	groupScheme.Quote = quote
	return groupScheme, nil
}

func GetGroupSchemeCategories(id string) ([]models.SchemeCategory, error) {
	var schemeCategories []models.SchemeCategory
	err := DB.Where("quote_id = ?", id).Find(&schemeCategories).Error
	if err != nil {
		return nil, err
	}

	if len(schemeCategories) == 0 {
		insurer, _ := GetInsurerDetails()
		masters, _ := GetSchemeCategoryMasters(insurer.ID)
		if len(masters) > 0 {
			for _, m := range masters {
				// Only include active categories
				if m.Active {
					schemeCategories = append(schemeCategories, models.SchemeCategory{
						SchemeCategory: m.Name,
					})
				}
			}
			if len(schemeCategories) > 0 {
				return schemeCategories, nil
			}
		}

		// Fallback to default categories if no masters configured or no active masters
		categories := []string{"Management", "Administration", "General"}
		for _, cat := range categories {
			schemeCategories = append(schemeCategories, models.SchemeCategory{
				// QuoteId will be set by the caller if needed, or we can use the id param
				SchemeCategory: cat,
			})
		}
	}

	return schemeCategories, nil
}

func AddMemberToScheme(member models.GPricingMemberDataInForce, user models.AppUser) (models.GPricingMemberDataInForce, error) {
	var created models.GPricingMemberDataInForce
	var schemecategory models.SchemeCategory
	// Validate RSA ID before proceeding
	idType := strings.ToUpper(strings.TrimSpace(member.MemberIdType))
	if (idType == "RSA_ID" || idType == "ID" || idType == "RSA_ISD") && strings.TrimSpace(member.MemberIdNumber) != "" {
		valid, checkErr := utils.ValidateRSAID(strings.TrimSpace(member.MemberIdNumber))
		if checkErr != nil {
			return created, fmt.Errorf("ID validation service error: %v", checkErr)
		}
		if !valid {
			return created, fmt.Errorf("invalid RSA ID '%s'", member.MemberIdNumber)
		}
	}

	// Wrap the whole operation in a transaction to ensure atomicity and audit consistency
	if err := DB.Transaction(func(tx *gorm.DB) error {
		// Prevent duplicates: a member with the same ID Number must not already exist for this scheme
		if member.MemberIdNumber != "" && member.SchemeId != 0 {
			var existing models.GPricingMemberDataInForce
			if err := tx.Where("scheme_id = ? AND member_id_number = ?", member.SchemeId, member.MemberIdNumber).First(&existing).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
			} else {
				return fmt.Errorf("member with ID Number %s already exists for this scheme", member.MemberIdNumber)
			}
		}

		var groupQuote models.GroupPricingQuote
		if err := tx.Preload("SchemeCategories").Where("scheme_id = ? and scheme_quote_status = ? ", member.SchemeId, models.StatusInEffect).First(&groupQuote).Error; err != nil {
			return err
		}

		if member.EntryDate.Before(groupQuote.CommencementDate) {
			return fmt.Errorf("Entry data cannot be before sheme commencement date")
		}

		err := DB.Where("quote_id = ? and scheme_category=?", groupQuote.ID, member.SchemeCategory).First(&schemecategory).Error

		member.Status = "Active"
		member.IsOriginalMember = false
		member.SchemeName = groupQuote.SchemeName
		member.CreationDate = time.Now()
		member.Year = time.Now().Year()
		member.CreatedBy = user.UserName
		//member.ExitDate = (1900, time.January, 15, 0, 0, 0, 0, time.UTC) //time.Now() // hack: solve later
		//member.EffectiveExitDate = (1900, time.January, 15, 0, 0, 0, 0, time.UTC)         //time.Now()
		member.Benefits.GlaEnabled = schemecategory.GlaBenefit
		member.Benefits.PtdEnabled = schemecategory.PtdBenefit
		member.Benefits.CiEnabled = schemecategory.CiBenefit
		member.Benefits.SglaEnabled = schemecategory.SglaBenefit
		member.Benefits.PhiEnabled = schemecategory.PtdBenefit
		member.Benefits.TtdEnabled = schemecategory.TtdBenefit
		member.Benefits.GffEnabled = schemecategory.FamilyFuneralBenefit

		if groupQuote.UseGlobalSalaryMultiple {
			member.Benefits.GlaMultiple = schemecategory.GlaSalaryMultiple
			member.Benefits.PtdMultiple = schemecategory.PtdSalaryMultiple
			member.Benefits.CiMultiple = schemecategory.CiCriticalIllnessSalaryMultiple
			member.Benefits.PhiMultiple = schemecategory.PhiIncomeReplacementPercentage / 100
			member.Benefits.TtdMultiple = schemecategory.TtdIncomeReplacementPercentage / 100
			member.Benefits.SglaMultiple = schemecategory.SglaSalaryMultiple
		}

		if !groupQuote.UseGlobalSalaryMultiple {
			member.Benefits.PhiMultiple = schemecategory.PhiIncomeReplacementPercentage / 100
			member.Benefits.TtdMultiple = schemecategory.TtdIncomeReplacementPercentage / 100
		}

		// the groupQuote.CommencementDate is  GMT time, so we need to convert it to local time
		sastLocation, err := time.LoadLocation("Africa/Johannesburg")
		if err != nil {
			panic(err)
		}

		// 2. Use the .In() method to convert the time
		//sastTime := 	groupQuote.CommencementDate.In(sastLocation)

		sastTime := member.DateOfBirth.In(sastLocation)
		fmt.Println(sastTime)
		member.DateOfBirth = sastTime

		// Persist member first to get ID
		if err := tx.Create(&member).Error; err != nil {
			return err
		}
		created = member

		// Generic audit for member creation
		if err := writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "g_pricing_member_data_in_forces",
			EntityID:  strconv.Itoa(created.ID),
			Action:    "CREATE",
			ChangedBy: user.UserName,
		}, models.GPricingMemberDataInForce{}, created); err != nil {
			return err
		}

		// Log structured Enrollment activity
		enrollmentDetails, _ := json.Marshal(map[string]interface{}{
			"scheme":          created.SchemeName,
			"initialSalary":   created.AnnualSalary,
			"entranceMedical": "Completed", // Defaulting for example, should be based on real data if available
		})
		_ = tx.Create(&models.MemberActivity{
			MemberID:       created.ID,
			MemberIDNumber: created.MemberIdNumber,
			Type:           "enrollment",
			Title:          "Member Enrollment",
			Description:    "Member enrolled in the group insurance scheme",
			Details:        enrollmentDetails,
			PerformedBy:    user.UserName,
		})

		GroupPricingCache.Clear()
		//var groupParameter models.GroupPricingParameters
		//var groupPricingReinsuranceStructure models.GroupPricingReinsuranceStructure
		//var memberRatingResultSummary models.MemberRatingResultSummary
		//
		//var incomeLevels []models.IncomeLevel
		//var ageBands []models.GroupPricingAgeBands

		//tx.Where("basis=?", groupQuote.Basis).Find(&groupParameter)
		//tx.Where("risk_rate_code=?", groupParameter.RiskRateCode).Find(&groupPricingReinsuranceStructure)

		//if err := tx.Where("risk_rate_code=?", groupParameter.RiskRateCode).Find(&incomeLevels).Error; err != nil {
		//	appLog.WithFields(map[string]interface{}{
		//		"risk_rate_code": groupParameter.RiskRateCode,
		//		"quote_id":       groupQuote.ID,
		//		"error":          err.Error(),
		//	}).Warn("Failed to find income levels")
		//}

		//if err := tx.Where("quote_id=?", groupQuote.ID).Find(&memberRatingResultSummary).Error; err != nil {
		//	appLog.WithFields(map[string]interface{}{
		//		"quote_id": groupQuote.ID,
		//		"error":    err.Error(),
		//	}).Warn("Failed to find member rating result summary")
		//}

		//if err := tx.Find(&ageBands).Error; err != nil {
		//	appLog.WithFields(map[string]interface{}{
		//		"quote_id": groupQuote.ID,
		//		"error":    err.Error(),
		//	}).Warn("Failed to find age bands")
		//}

		// Generic CREATE audit for the new member
		if err := writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "g_pricing_member_data_in_forces",
			EntityID:  strconv.Itoa(member.ID),
			Action:    "CREATE",
			ChangedBy: user.UserName,
		}, struct{}{}, member); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return created, err
	}
	return created, nil
}

func GetSchemeMembers(schemeId string) ([]models.GPricingMemberDataInForce, error) {
	logger := appLog.WithFields(map[string]interface{}{
		"scheme_id": schemeId,
		"function":  "GetSchemeMembers",
	})

	logger.Debug("Retrieving scheme members")

	var members []models.GPricingMemberDataInForce
	err := DB.Where("scheme_id = ?", schemeId).Find(&members).Error
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to retrieve scheme members")
		return nil, err
	}

	logger.WithField("member_count", len(members)).Info("Successfully retrieved scheme members")
	return members, nil
}

// getSchemeNameByID returns the scheme name for a given scheme ID.
// If the scheme is not found, it returns an empty string and gorm.ErrRecordNotFound.
func getSchemeNameByID(id int) (string, error) {
	var scheme models.GroupScheme
	if err := DB.Select("id, name").Where("id = ?", id).First(&scheme).Error; err != nil {
		return "", err
	}
	return scheme.Name, nil
}

// GetMemberInForceByID retrieves a single GPricingMemberDataInForce by its primary key ID
func GetMemberInForceByID(id string) (models.GPricingMemberDataInForce, error) {
	logger := appLog.WithFields(map[string]interface{}{
		"member_id": id,
		"function":  "GetMemberInForceByID",
	})

	logger.Debug("Retrieving member in-force by ID")

	var member models.GPricingMemberDataInForce
	if err := DB.Where("id = ?", id).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.WithField("error", err.Error()).Warn("Member in-force not found")
			return models.GPricingMemberDataInForce{}, err
		}
		logger.WithField("error", err.Error()).Error("Failed to retrieve member in-force")
		return models.GPricingMemberDataInForce{}, err
	}

	// Enrich SchemeName before returning (do not fail the call if scheme lookup fails)
	if member.SchemeId != 0 {
		if schemeName, err := getSchemeNameByID(member.SchemeId); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.WithFields(map[string]interface{}{
					"scheme_id": member.SchemeId,
				}).Warn("Scheme not found while enriching scheme name")
			} else {
				logger.WithFields(map[string]interface{}{
					"scheme_id": member.SchemeId,
					"error":     err.Error(),
				}).Error("Failed to enrich scheme name")
			}
		} else {
			member.SchemeName = schemeName
		}
	}

	logger.Info("Successfully retrieved member in-force")
	return member, nil
}

// GetMemberInForceByIdNumber retrieves a single GPricingMemberDataInForce by the member's IdNumber
// If multiple records exist for the same IdNumber, it prefers Active status and most recent entry.
func GetMemberInForceByIdNumber(idNumber string) (models.GPricingMemberDataInForce, error) {
	logger := appLog.WithFields(map[string]interface{}{
		"member_id_number": idNumber,
		"function":         "GetMemberInForceByIdNumber",
	})

	logger.Debug("Retrieving member in-force by IdNumber")

	var member models.GPricingMemberDataInForce
	// Prefer Active members, then by most recent entry date
	if err := DB.Where("member_id_number = ?", idNumber).
		Order("CASE WHEN status = 'Active' THEN 0 ELSE 1 END").
		Order("entry_date DESC").
		First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.WithField("error", err.Error()).Warn("Member in-force not found for IdNumber")
			return models.GPricingMemberDataInForce{}, err
		}
		logger.WithField("error", err.Error()).Error("Failed to retrieve member in-force by IdNumber")
		return models.GPricingMemberDataInForce{}, err
	}

	// Enrich SchemeName & SchemeCategoryDetails (do not fail the call if lookups fail)
	var scheme models.GroupScheme
	var haveScheme bool
	if member.SchemeId != 0 {
		// Try by ID first
		if err := DB.Where("id = ?", member.SchemeId).First(&scheme).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.WithField("scheme_id", member.SchemeId).Warn("Scheme not found by ID while enriching member")
			} else {
				logger.WithFields(map[string]interface{}{"scheme_id": member.SchemeId, "error": err.Error()}).Error("Failed to fetch scheme by ID for enrichment")
			}
		} else {
			haveScheme = true
			if member.SchemeName == "" {
				member.SchemeName = scheme.Name
			}
		}
	}
	// If scheme still not resolved, try by name if we have one
	if !haveScheme {
		name := strings.TrimSpace(member.SchemeName)
		if name == "" && member.SchemeId != 0 {
			if schemeName, err := getSchemeNameByID(member.SchemeId); err == nil {
				name = schemeName
				member.SchemeName = schemeName
			}
		}
		if name != "" {
			if err := DB.Where("LOWER(name) = LOWER(?)", name).First(&scheme).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					logger.WithField("scheme_name", name).Warn("Scheme not found by name while enriching member")
				} else {
					logger.WithFields(map[string]interface{}{"scheme_name": name, "error": err.Error()}).Error("Failed to fetch scheme by name for enrichment")
				}
			} else {
				haveScheme = true
			}
		}
	}

	// If we have scheme, try to resolve in-force quote and scheme category details
	if haveScheme {
		// Determine quote in force
		var quote models.GroupPricingQuote
		quoteName := strings.TrimSpace(scheme.QuoteInForce)
		if quoteName != "" {
			if err := DB.Where("quote_name = ?", quoteName).First(&quote).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					logger.WithFields(map[string]interface{}{"quote_name": quoteName, "error": err.Error()}).Error("Failed to fetch in-force quote by name")
				} else {
					logger.WithField("quote_name", quoteName).Warn("In-force quote name on scheme not found")
				}
			}
		}
		// Fallback: latest quote for scheme
		if quote.ID == 0 {
			if err := DB.Where("scheme_id = ?", scheme.ID).Order("id desc").First(&quote).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					logger.WithFields(map[string]interface{}{"scheme_id": scheme.ID, "error": err.Error()}).Error("Failed to fetch latest quote for scheme")
				}
			}
		}

		// With quote resolved and member's SchemeCategory name, fetch details
		if quote.ID != 0 && strings.TrimSpace(member.SchemeCategory) != "" {
			var cat models.SchemeCategory
			if err := DB.Where("quote_id = ? AND LOWER(scheme_category) = LOWER(?)", quote.ID, member.SchemeCategory).First(&cat).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					logger.WithFields(map[string]interface{}{"quote_id": quote.ID, "scheme_category": member.SchemeCategory}).Warn("Scheme category details not found for member")
				} else {
					logger.WithFields(map[string]interface{}{"quote_id": quote.ID, "scheme_category": member.SchemeCategory, "error": err.Error()}).Error("Failed to fetch scheme category details for member")
				}
			} else {
				member.SchemeCategoryDetails = cat
			}
		}
	}

	logger.Info("Successfully retrieved member in-force by IdNumber")
	return member, nil
}

func DeleteGroupScheme(id string) error {
	logger := appLog.WithFields(map[string]interface{}{
		"scheme_id": id,
		"function":  "DeleteGroupScheme",
	})

	logger.Debug("Deleting group scheme")

	return DB.Transaction(func(tx *gorm.DB) error {
		var before models.GroupScheme
		if err := tx.Where("id = ?", id).First(&before).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", id).Delete(&models.GroupScheme{}).Error; err != nil {
			logger.WithField("error", err.Error()).Error("Failed to delete group scheme")
			return err
		}
		// Write generic DELETE audit with snapshot in prev_values
		if err := writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "group_schemes",
			EntityID:  id,
			Action:    "DELETE",
			ChangedBy: "system",
		}, before, struct{}{}); err != nil {
			return err
		}
		logger.Info("Successfully deleted group scheme")
		return nil
	})
}

func GetGroupPricingParameterBases() ([]string, error) {
	var parameterBases []string
	err := DB.Model(models.GroupPricingParameters{}).Distinct("basis").Pluck("basis", &parameterBases).Error
	if err != nil {
		return nil, err
	}
	return parameterBases, nil
}

func GetGroupPricingIndustries() ([]string, error) {
	var industries []string
	err := DB.Model(models.OccupationClass{}).Distinct("industry").Pluck("industry", &industries).Error
	if err != nil {
		return nil, err
	}
	return industries, nil
}

func GetGroupPricingQuoteTableData(quoteId int, tableType string, offset int, limit int) (map[string]interface{}, error) {
	var results []map[string]interface{}
	resultData := make(map[string]interface{})
	var jsonTags []string

	// sanitize
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 0
	}

	switch tableType {
	case "member_data":
		var quote models.GroupPricingQuote
		if err := DB.Where("id = ?", quoteId).First(&quote).Error; err == nil && quote.QuoteType == "Renewal" {
			var memberData []models.GPricingMemberDataInForce
			db := DB.Where("scheme_id = ?", quote.SchemeID)
			if limit > 0 {
				db = db.Offset(offset).Limit(limit)
			}
			db.Find(&memberData)
			b, _ := json.Marshal(&memberData)
			err := json.Unmarshal(b, &results)
			if err != nil {
				return nil, err
			}
			jsonTags = getJSONTags(models.GPricingMemberDataInForce{})
		} else {
			var memberData []models.GPricingMemberData
			db := DB.Where("quote_id = ?", quoteId)
			if limit > 0 {
				db = db.Offset(offset).Limit(limit)
			}
			db.Find(&memberData)
			b, _ := json.Marshal(&memberData)
			err := json.Unmarshal(b, &results)
			if err != nil {
				return nil, err
			}
			jsonTags = getJSONTags(models.GPricingMemberData{})
		}
	case "claims_experience":
		var claimsExperience []models.GroupPricingClaimsExperience
		db := DB.Where("quote_id = ?", quoteId)
		if limit > 0 {
			db = db.Offset(offset).Limit(limit)
		}
		db.Find(&claimsExperience)
		b, _ := json.Marshal(&claimsExperience)
		err := json.Unmarshal(b, &results)
		if err != nil {
			return nil, err
		}
		jsonTags = getJSONTags(models.GroupPricingClaimsExperience{})
	case "member_rating_results":
		var memberRatingResults []models.MemberRatingResult
		db := DB.Where("quote_id = ?", quoteId)
		if limit > 0 {
			db = db.Offset(offset).Limit(limit)
		}
		db.Find(&memberRatingResults)
		b, _ := json.Marshal(&memberRatingResults)
		err := json.Unmarshal(b, &results)
		if err != nil {
			return nil, err
		}
		jsonTags = getJSONTags(models.MemberRatingResult{})
	case "member_premium_schedules":
		var memberPremiumSchedules []models.MemberPremiumSchedule
		db := DB.Where("quote_id = ?", quoteId)
		if limit > 0 {
			db = db.Offset(offset).Limit(limit)
		}
		db.Find(&memberPremiumSchedules)
		b, _ := json.Marshal(&memberPremiumSchedules)
		err := json.Unmarshal(b, &results)
		if err != nil {
			return nil, err
		}
		jsonTags = getJSONTags(models.MemberPremiumSchedule{})
	case "group_pricing_parameters":
		var groupPricingParameters []models.GroupPricingParameters
		db := DB
		if limit > 0 {
			db = db.Offset(offset).Limit(limit)
		}
		db.Find(&groupPricingParameters)
		b, _ := json.Marshal(&groupPricingParameters)
		err := json.Unmarshal(b, &results)
		if err != nil {
			return nil, err
		}
		jsonTags = getJSONTags(models.GroupPricingParameters{})
	case "bordereaux":
		// Bordereaux rows are projected on-the-fly from MemberRatingResult so
		// the latest loadings and reinsurance structure flow through without
		// requiring a recalculation.
		bordereaux, err := BuildBordereauxRowsForQuote(quoteId, offset, limit)
		if err != nil {
			return nil, err
		}
		b, _ := json.Marshal(&bordereaux)
		if err := json.Unmarshal(b, &results); err != nil {
			return nil, err
		}
		// jsonTags drives which columns the frontend renders. The shared
		// QuoteBordereauxJSONTags list omits columns the user asked to hide
		// (loaded/retained rates and annual premiums) so the grid matches
		// the xlsx export shape.
		jsonTags = QuoteBordereauxJSONTags()
	}

	resultData["data"] = results
	resultData["json_tags"] = jsonTags

	return resultData, nil
}

func GetInForceTableData(schemeId, tableType string) (map[string]interface{}, error) {
	var results []map[string]interface{}
	resultData := make(map[string]interface{})
	var jsonTags []string

	switch tableType {
	case "member_data":
		var memberData []models.GPricingMemberDataInForce
		DB.Where("scheme_id = ?", schemeId).Find(&memberData)
		b, _ := json.Marshal(&memberData)
		err := json.Unmarshal(b, &results)
		if err != nil {
			return nil, err
		}
		jsonTags = getJSONTags(models.GPricingMemberData{})
	}

	resultData["data"] = results
	resultData["json_tags"] = jsonTags

	return resultData, nil
}

// GetMembersPaginated returns paginated GPricingMemberDataInForce records with optional filtering
// Filters supported: schemeId, status, search (matches member_name or member_id_number)
// Pagination: page (1-based), pageSize
func GetMembersPaginated(page, pageSize int, search, schemeId, status string) ([]models.GPricingMemberDataInForce, int64, error) {
	logger := appLog.WithFields(map[string]interface{}{
		"function":  "GetMembersPaginated",
		"page":      page,
		"page_size": pageSize,
		"search":    search,
		"scheme_id": schemeId,
		"status":    status,
	})

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	var (
		members []models.GPricingMemberDataInForce
		total   int64
	)

	db := DB.Model(&models.GPricingMemberDataInForce{})

	if schemeId != "" {
		db = db.Where("scheme_id = ?", schemeId)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if search != "" {
		like := "%" + strings.TrimSpace(search) + "%"
		db = db.Where("member_name LIKE ? OR member_id_number LIKE ?", like, like)
	}

	// Count total with filters
	if err := db.Count(&total).Error; err != nil {
		logger.WithField("error", err.Error()).Error("failed to count members")
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("id DESC").Offset(offset).Limit(pageSize).Find(&members).Error; err != nil {
		logger.WithField("error", err.Error()).Error("failed to fetch paginated members")
		return nil, 0, err
	}

	return members, total, nil
}

func GetGroupPricingQuoteResultSummary(quoteId int) ([]models.MemberRatingResultSummary, error) {
	summaries := make([]models.MemberRatingResultSummary, 0)
	err := DB.Where("quote_id = ?", quoteId).Find(&summaries).Error
	if err != nil {
		return summaries, err
	}
	return summaries, nil
}

func SaveInsurerDetails(form *multipart.Form, user models.AppUser) error {
	fmt.Println("form: ", form)
	var insurerDetails models.GroupPricingInsurerDetail

	// New way: extract individual fields from form data
	if ids, ok := form.Value["id"]; ok && len(ids) > 0 {
		id, _ := strconv.Atoi(ids[0])
		insurerDetails.ID = id
	}
	if names, ok := form.Value["name"]; ok && len(names) > 0 {
		insurerDetails.Name = names[0]
	}
	if contactPersons, ok := form.Value["contact_person"]; ok && len(contactPersons) > 0 {
		insurerDetails.ContactPerson = contactPersons[0]
	}
	if address1s, ok := form.Value["address_line_1"]; ok && len(address1s) > 0 {
		insurerDetails.AddressLine1 = address1s[0]
	}
	if address2s, ok := form.Value["address_line_2"]; ok && len(address2s) > 0 {
		insurerDetails.AddressLine2 = address2s[0]
	}
	if address3s, ok := form.Value["address_line_3"]; ok && len(address3s) > 0 {
		insurerDetails.AddressLine3 = address3s[0]
	}
	if postCodes, ok := form.Value["post_code"]; ok && len(postCodes) > 0 {
		insurerDetails.PostCode = postCodes[0]
	}
	if provinces, ok := form.Value["province"]; ok && len(provinces) > 0 {
		insurerDetails.Province = provinces[0]
	}
	if cities, ok := form.Value["city"]; ok && len(cities) > 0 {
		insurerDetails.City = cities[0]
	}
	if countries, ok := form.Value["country"]; ok && len(countries) > 0 {
		insurerDetails.Country = countries[0]
	}
	if telephones, ok := form.Value["telephone"]; ok && len(telephones) > 0 {
		insurerDetails.Telephone = telephones[0]
	}
	if emails, ok := form.Value["email"]; ok && len(emails) > 0 {
		insurerDetails.Email = emails[0]
	}
	if yearEndMonths, ok := form.Value["year_end_month"]; ok && len(yearEndMonths) > 0 {
		yem, _ := strconv.Atoi(yearEndMonths[0])
		insurerDetails.YearEndMonth = yem
	}
	if introductoryTexts, ok := form.Value["introductory_text"]; ok && len(introductoryTexts) > 0 {
		insurerDetails.IntroductoryText = introductoryTexts[0]
	}
	if generalProvisionsTexts, ok := form.Value["general_provisions_text"]; ok && len(generalProvisionsTexts) > 0 {
		insurerDetails.GeneralProvisionsText = generalProvisionsTexts[0]
	}

	// Support old way if "insurer" field is present
	var err error
	if insurerJsonStrings, ok := form.Value["insurer"]; ok && len(insurerJsonStrings) > 0 {
		err = json.Unmarshal([]byte(insurerJsonStrings[0]), &insurerDetails)
		if err != nil {
			return err
		}
	}

	var file multipart.File
	if form.File["logo"] != nil || len(form.File["logo"]) > 0 {
		file, err = form.File["logo"][0].Open()
		if err != nil {
			return err
		}
		defer file.Close()
		header := make([]byte, 512)
		_, err = file.Read(header)
		if err != nil {
			return fmt.Errorf("failed to read file header for type detection: %w", err)
		}

		// Reset file reader back to beginning before full read
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return fmt.Errorf("failed to reset file reader: %w", err)
		}

		// Detect MIME type
		mimeType := http.DetectContentType(header)
		fmt.Println("Detected MIME type:", mimeType)

		// Accept only PNG or JPEG
		if mimeType != "image/png" && mimeType != "image/jpeg" {
			return fmt.Errorf("unsupported file type: %s (only PNG and JPEG allowed)", mimeType)
		}

		fileBytes, err := io.ReadAll(file)

		if err != nil {
			return err
		}

		insurerDetails.Logo = fileBytes
		insurerDetails.LogoMimeType = mimeType
		fmt.Println("fileBytes: ", string(fileBytes))

	}

	insurerDetails.CreationDate = time.Now()

	if insurerDetails.ID > 0 {
		err = DB.Where("id = ?", insurerDetails.ID).Updates(&insurerDetails).Error
		if err != nil {
			return err
		}
		return nil
	} else {
		err = DB.Create(&insurerDetails).Error
	}

	if err != nil {
		return err
	}

	return nil
}

func GetInsurerDetails() (models.GroupPricingInsurerDetail, error) {
	var insurer models.GroupPricingInsurerDetail
	err := DB.First(&insurer).Error
	if err != nil {
		return insurer, err
	}
	return insurer, nil
}

// CategoryEducatorBenefit represents the educator benefit structure for a specific category
type CategoryEducatorBenefit struct {
	SchemeCategory           string                          `json:"scheme_category"`
	EducatorBenefitStructure models.EducatorBenefitStructure `json:"educator_benefit_structure"`
}

// GetGroupPricingQuoteEducatorBenefits retrieves educator benefit structures for each category in a quote
func GetGroupPricingQuoteEducatorBenefits(quoteId int) ([]CategoryEducatorBenefit, error) {
	// Convert quoteId to string for GetGroupPricingQuote function
	quoteIdStr := strconv.Itoa(quoteId)

	// Get the quote
	quote, err := GetGroupPricingQuote(quoteIdStr)
	if err != nil {
		return nil, err
	}

	// Create a slice to hold the results
	results := make([]CategoryEducatorBenefit, 0)

	// For each scheme category in the quote
	for _, category := range quote.SchemeCategories {
		// Get the basis value
		basis := category.Basis

		// Get the group pricing parameter using the basis
		var groupParameter models.GroupPricingParameters
		err := DB.Where("basis = ?", basis).First(&groupParameter).Error
		if err != nil {
			// Skip this category if no matching parameter is found
			continue
		}

		// Determine the educator benefit code from the scheme category
		// (GLA takes priority, then PTD).
		educatorBenefitCode := educatorCodeForCategory(&category)
		if educatorBenefitCode == "" {
			// Skip categories that don't have an educator benefit configured
			continue
		}

		// Get the educator benefit structure using the risk rate code and educator benefit code
		var educatorBenefitStructure models.EducatorBenefitStructure
		err = DB.Where("risk_rate_code = ? AND educator_benefit_code = ?", groupParameter.RiskRateCode, educatorBenefitCode).First(&educatorBenefitStructure).Error
		if err != nil {
			// Skip this category if no matching structure is found
			continue
		}

		// Add the educator benefit structure to the results
		results = append(results, CategoryEducatorBenefit{
			SchemeCategory:           category.SchemeCategory,
			EducatorBenefitStructure: educatorBenefitStructure,
		})
	}

	return results, nil
}

type UpdateError struct {
	Field   string
	Message string
}

func (e *UpdateError) Error() string {
	return fmt.Sprintf("Update operation failed on %s: %s", e.Field, e.Message)
}

func UpdateGroupSchemeCoverEndDate(schemeId int, coverEndDate time.Time, user models.AppUser) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var before models.GroupScheme
		if err := tx.Where("id = ?", schemeId).First(&before).Error; err != nil {
			return err
		}
		if before.CoverStartDate.After(coverEndDate) {
			return &UpdateError{Field: "Cover End Date Update", Message: "Cover End Date cannot be before Cover Start Date."}
		}
		if err := tx.Model(&models.GroupScheme{}).Where("id = ?", schemeId).Update("cover_end_date", coverEndDate).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.GroupScheme{}).Where("id = ?", schemeId).Update("renewal_date", coverEndDate.AddDate(0, 0, 1)).Error; err != nil {
			return err
		}
		var after models.GroupScheme
		if err := tx.Where("id = ?", schemeId).First(&after).Error; err != nil {
			return err
		}
		if err := writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "group_schemes",
			EntityID:  strconv.Itoa(schemeId),
			Action:    "UPDATE",
			ChangedBy: user.UserName,
		}, before, after); err != nil {
			return err
		}
		return nil
	})
}

func GetGroupPricingDashboardData(year int, dataSource string, benefit string) (map[string]interface{}, error) {
	resultData := make(map[string]interface{})

	// get all quotes and quotes that in force
	var returnedQuotes []models.GroupPricingQuote

	var allQuotesCount int64

	//var groupSchemes []models.GroupScheme
	var newBusinessAcceptedQuotes int64
	var renewedAcceptedQuotes int64
	//var newBusinessInforceSchemesPremium float64
	var renewedInforceSchemesPremium float64
	var gs models.GroupScheme
	// get all quotes
	err := DB.Where(" year(creation_date) = ?", year).Find(&returnedQuotes).Count(&allQuotesCount).Error

	//Accepted Quotes

	//New Business Count
	err = DB.Where(" year(creation_date) = ? and status=? and quote_type=?", year, models.StatusAccepted, "New Business").Find(&returnedQuotes).Count(&newBusinessAcceptedQuotes).Error
	if err != nil {
		return nil, err
	}

	//Renewals Count
	err = DB.Where(" year(creation_date) = ? and status=? and quote_type=?", year, models.StatusAccepted, "Renewal").Find(&returnedQuotes).Count(&renewedAcceptedQuotes).Error
	if err != nil {
		return nil, err
	}

	//New Business Premium
	//err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ? and new_business=?", year, true, true).Pluck("sum(annual_premium)", &newBusinessInforceSchemesPremium).Error
	//if err != nil {
	//	return nil, err
	//}

	//err = DB.Table("group_schemes").
	//	Select("SUM(annual_premium) AS annual_premium").
	//	Where(" year(creation_date) = ? and in_force = ? and new_business=?", year, true, true).
	//	Scan(&newBusinessInforceSchemesPremium).Error
	//err = DB.Table("group_schemes").Select("sum(annual_premium) as annual_premium").Where("year(creation_date) = ? and in_force = ? and new_business=?", year, true, true).Scan(&gs).Error
	//newBusinessInforceSchemesPremium = gs.AnnualPremium

	//if err != nil {
	//	return nil, err
	//}

	//Renewals Premium
	//err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ? and new_business=?", year, true, false).Pluck("sum(annual_premium)", &renewedInforceSchemesPremium).Error
	//if err != nil {
	//	return nil, err
	//}

	//err = DB.Table("group_schemes").
	//	Select("SUM(annual_premium) AS annual_premium").
	//	Where(" year(creation_date) = ? and in_force = ? and new_business=?", year, true, false).
	//	Scan(&renewedInforceSchemesPremium).Error
	err = DB.Table("group_schemes").Select("sum(annual_premium) as annual_premium").Where("year(creation_date) = ? and in_force = ? and new_business=?", year, true, false).Scan(&gs).Error
	renewedInforceSchemesPremium = gs.AnnualPremium
	if err != nil {
		return nil, err
	}

	// New Business Quotes
	var newQuotesCountApproved int64
	var newQuotesCountInProgress int64
	var newQuotesCountPendingReview int64
	var newQuotesInForceCount int64
	var newQuotesCountAccepted int64

	//In Progress
	err = DB.Where(" year(creation_date) = ? and status=? and quote_type=?", year, models.StatusInProgress, "New Business").Find(&returnedQuotes).Count(&newQuotesCountInProgress).Error
	if err != nil {
		return nil, err
	}

	//Pending Review
	err = DB.Where(" year(creation_date) = ? and status=? and quote_type=?", year, models.StatusPendingReview, "New Business").Find(&returnedQuotes).Count(&newQuotesCountPendingReview).Error
	if err != nil {
		return nil, err
	}

	//Approved
	err = DB.Where(" year(creation_date) = ? and status=? and quote_type=?", year, models.StatusApproved, "New Business").Find(&returnedQuotes).Count(&newQuotesCountApproved).Error
	if err != nil {
		return nil, err
	}

	//Accepted
	err = DB.Where(" year(creation_date) = ? and status=? and quote_type=?", year, models.StatusAccepted, "New Business").Find(&returnedQuotes).Count(&newQuotesCountAccepted).Error
	if err != nil {
		return nil, err
	}

	// In Force
	//err = DB.Where(" year(creation_date) = ? and status = ? and quote_type=?", year, models.StatusInForce, "New Business").Find(&returnedQuotes).Count(&newQuotesInForceCount).Error
	//if err != nil {
	//	return nil, err
	//}

	// get all quotes by annual premium
	var newQuotesPremiumInProgress float64
	var newQuotesPremiumPendingReview float64
	var newQuotesPremiumApproved float64
	var newQuotesPremiumInForce float64
	var newQuotesPremiumAccepted float64
	//In Progress
	//err = DB.Table("member_rating_result_summaries").Where(" year(creation_date) = ? and if_status = ? and quote_type=?", year, "In Progress", "New Business").Pluck("sum(total_annual_premium) as total_annual_premium", &newQuotesPremiumInProgress).Error
	//if err != nil {
	//	return nil, err
	//}

	var mrrs models.MemberRatingResultSummary
	err = DB.Table("member_rating_result_summaries").
		Select("sum(total_annual_premium) as total_annual_premium").
		Where(" year(creation_date) = ? and if_status = ? and quote_type=?", year, models.StatusInProgress, "New Business").
		Scan(&mrrs).Error

	if err != nil {
		return nil, err
	}

	newQuotesPremiumInProgress = mrrs.TotalAnnualPremium
	fmt.Println("newQuotesPremiumInProgress: ", newQuotesPremiumInProgress)

	//Pending Review
	err = DB.Table("member_rating_result_summaries").
		Select("sum(total_annual_premium) as total_annual_premium").
		Where(" year(creation_date) = ? and if_status = ? and quote_type=?", year, models.StatusPendingReview, "New Business").
		Scan(&mrrs).Error

	if err != nil {
		return nil, err
	}

	newQuotesPremiumPendingReview = mrrs.TotalAnnualPremium
	fmt.Println("newQuotesPremiumInProgress: ", newQuotesPremiumPendingReview)

	//Approved
	err = DB.Table("member_rating_result_summaries").
		Select("sum(total_annual_premium) as total_annual_premium").
		Where(" year(creation_date) = ? and if_status = ? and quote_type=?", year, models.StatusApproved, "New Business").
		Scan(&mrrs).Error

	if err != nil {
		return nil, err
	}

	newQuotesPremiumApproved = mrrs.TotalAnnualPremium

	//Accepted
	err = DB.Table("member_rating_result_summaries").
		Select("sum(total_annual_premium) as total_annual_premium").
		Where(" year(creation_date) = ? and if_status = ? and quote_type=?", year, models.StatusAccepted, "New Business").
		Scan(&mrrs).Error

	if err != nil {
		return nil, err
	}

	newQuotesPremiumAccepted = mrrs.TotalAnnualPremium

	// In Force
	//err = DB.Table("member_rating_result_summaries").
	//	Select("sum(total_annual_premium) as total_annual_premium").
	//	Where(" year(creation_date) = ? and if_status = ? and quote_type=?", year, models.StatusInForce, "New Business").
	//	Scan(&mrrs).Error

	//if err != nil {
	//	return nil, err
	//}
	//
	//newQuotesPremiumInForce = mrrs.TotalAnnualPremium

	//Renewals Quotes

	var renewalQuotesCountApproved int64
	var renewalQuotesCountInProgress int64
	var renewalQuotesCountPendingReview int64
	var renewalQuotesInForceCount int64
	var renewalQuotesCountAccepted int64

	//In Progress
	err = DB.Where(" year(creation_date) = ? and status=? and quote_type=?", year, models.StatusInProgress, "Renewal").Find(&returnedQuotes).Count(&renewalQuotesCountInProgress).Error
	if err != nil {
		return nil, err
	}

	//Pending Review
	err = DB.Where(" year(creation_date) = ? and status=? and quote_type=?", year, models.StatusPendingReview, "Renewal").Find(&returnedQuotes).Count(&renewalQuotesCountPendingReview).Error
	if err != nil {
		return nil, err
	}

	//Approved
	err = DB.Where(" year(creation_date) = ? and status=? and quote_type=?", year, models.StatusApproved, "Renewal").Find(&returnedQuotes).Count(&renewalQuotesCountApproved).Error
	if err != nil {
		return nil, err
	}

	//Accepted
	err = DB.Where(" year(creation_date) = ? and status=? and quote_type=?", year, models.StatusAccepted, "Renewal").Find(&returnedQuotes).Count(&renewalQuotesCountAccepted).Error
	if err != nil {
		return nil, err
	}

	// In Force
	//err = DB.Where(" year(creation_date) = ? and status = ? and quote_type=?", year, models.StatusInForce, "Renewal").Find(&returnedQuotes).Count(&renewalQuotesInForceCount).Error

	//if err != nil {
	//	return nil, err
	//}

	// get all quotes by annual premium
	var renewalQuotesPremiumInProgress float64
	var renewalQuotesPremiumPendingReview float64
	var renewalQuotesPremiumApproved float64
	var renewalQuotesPremiumInForce float64
	var renewalQuotesPremiumAccepted float64
	//In Progress
	//err = DB.Table("member_rating_result_summaries").Where(" year(creation_date) = ? and if_status = ? and quote_type=?", year, "In Progress", "Renewal").Pluck("sum(total_annual_premium)", &renewalQuotesPremiumInProgress).Error
	err = DB.Table("member_rating_result_summaries").
		Select("sum(total_annual_premium) as total_annual_premium").
		Where(" year(creation_date) = ? and if_status = ? and quote_type=?", year, models.StatusInProgress, "Renewal").
		Scan(&mrrs).Error

	if err != nil {
		return nil, err
	}

	renewalQuotesPremiumInProgress = mrrs.TotalAnnualPremium
	fmt.Println("renewalQuotesPremiumInProgress: ", renewalQuotesPremiumInProgress)

	//Pending Review
	//err = DB.Table("member_rating_result_summaries").Where(" year(creation_date) = ? and if_status = ? and quote_type=?", year, "Pending Review", "Renewal").Pluck("sum(total_annual_premium)", &renewalQuotesPremiumPendingReview).Error
	err = DB.Table("member_rating_result_summaries").
		Select("sum(total_annual_premium) as total_annual_premium").
		Where(" year(creation_date) = ? and if_status = ? and quote_type=?", year, models.StatusPendingReview, "Renewal").
		Scan(&mrrs).Error

	if err != nil {
		return nil, err
	}

	renewalQuotesPremiumPendingReview = mrrs.TotalAnnualPremium
	fmt.Println("renewalQuotesPremiumPendingReview: ", renewalQuotesPremiumPendingReview)

	//Approved
	//err = DB.Table("member_rating_result_summaries").Where(" year(creation_date) = ? and if_status = ? and quote_type=?", year, "Approved", "Renewal").Pluck("sum(total_annual_premium)", &renewalQuotesPremiumApproved).Error
	err = DB.Table("member_rating_result_summaries").
		Select("sum(total_annual_premium) as total_annual_premium").
		Where(" year(creation_date) = ? and if_status = ? and quote_type=?", year, models.StatusApproved, "Renewal").
		Scan(&mrrs).Error

	if err != nil {
		return nil, err
	}

	renewalQuotesPremiumApproved = mrrs.TotalAnnualPremium

	//Accepted
	//err = DB.Table("member_rating_result_summaries").Where(" year(creation_date) = ? and if_status = ? and quote_type=?", year, "Approved", "Renewal").Pluck("sum(total_annual_premium)", &renewalQuotesPremiumApproved).Error
	err = DB.Table("member_rating_result_summaries").
		Select("sum(total_annual_premium) as total_annual_premium").
		Where(" year(creation_date) = ? and if_status = ? and quote_type=?", year, models.StatusAccepted, "Renewal").
		Scan(&mrrs).Error

	if err != nil {
		return nil, err
	}

	renewalQuotesPremiumAccepted = mrrs.TotalAnnualPremium

	// In Force
	//err = DB.Table("member_rating_result_summaries").
	//	Select("sum(total_annual_premium) as total_annual_premium").
	//	Where(" year(creation_date) = ? and if_status = ? and quote_type=?", year, models.StatusInForce, "Renewal").
	//	Scan(&mrrs).Error
	//
	//if err != nil {
	//	return nil, err
	//}

	resultData["all_quotes"] = allQuotesCount

	//Accepted Quote count
	resultData["new_business_in_force_count"] = newBusinessAcceptedQuotes
	resultData["renewals_in_force_count"] = renewedAcceptedQuotes
	resultData["total_in_force_count"] = newBusinessAcceptedQuotes + renewedAcceptedQuotes

	//Accepted Quote premium
	resultData["new_business_in_force_premium"] = newQuotesPremiumInForce + newQuotesPremiumAccepted
	resultData["renewals_in_force_premium"] = renewedInforceSchemesPremium + renewalQuotesPremiumAccepted
	resultData["total_in_force_premium"] = newQuotesPremiumInForce + newQuotesPremiumAccepted + renewedInforceSchemesPremium + renewalQuotesPremiumAccepted

	//new business count
	resultData["new_quotes_in_progress_count"] = newQuotesCountInProgress
	resultData["new_quotes_pending_review_count"] = newQuotesCountPendingReview
	resultData["new_quotes_approved_count"] = newQuotesCountApproved
	resultData["new_quotes_in_force_count"] = newQuotesInForceCount
	resultData["new_quotes_accepted_count"] = newQuotesCountAccepted
	resultData["new_quotes_total_count"] = newQuotesCountInProgress + newQuotesCountPendingReview + newQuotesCountAccepted + newQuotesCountApproved + newQuotesInForceCount

	//new business premium
	resultData["new_quotes_in_progress_premium"] = newQuotesPremiumInProgress
	resultData["new_quotes_pending_review_premium"] = newQuotesPremiumPendingReview
	resultData["new_quotes_approved_premium"] = newQuotesPremiumApproved
	resultData["new_quotes_in_force_premium"] = newQuotesPremiumInForce
	resultData["new_quotes_accepted_premium"] = newQuotesPremiumAccepted
	resultData["new_quotes_total_premium"] = newQuotesPremiumInProgress + newQuotesPremiumPendingReview + newQuotesPremiumApproved + newQuotesPremiumInForce + newQuotesPremiumAccepted

	//renewals count
	resultData["renewals_quotes_in_progress_count"] = renewalQuotesCountInProgress
	resultData["renewals_quotes_pending_review_count"] = renewalQuotesCountPendingReview
	resultData["renewals_quotes_approved_count"] = renewalQuotesCountApproved
	resultData["renewals_quotes_accepted_count"] = renewalQuotesCountAccepted
	resultData["renewals_quotes_in_force_count"] = renewalQuotesInForceCount
	resultData["renewals_total_count"] = renewalQuotesCountInProgress + renewalQuotesCountPendingReview + renewalQuotesCountApproved + renewalQuotesCountAccepted + renewalQuotesInForceCount

	//renewals premium
	resultData["renewals_quotes_in_progress_premium"] = renewalQuotesPremiumInProgress
	resultData["renewals_quotes_pending_review_premium"] = renewalQuotesPremiumPendingReview
	resultData["renewals_quotes_approved_premium"] = renewalQuotesPremiumApproved
	resultData["renewals_quotes_accepted_premium"] = renewalQuotesPremiumAccepted
	resultData["renewals_quotes_in_force_premium"] = renewalQuotesPremiumInForce
	resultData["renewals_total_premium"] = renewalQuotesPremiumInProgress + renewalQuotesPremiumPendingReview + renewalQuotesPremiumApproved + renewalQuotesPremiumAccepted + renewalQuotesPremiumInForce

	//Conversion Rates
	resultData["new_quotes_unconverted_count"] = newQuotesCountApproved
	resultData["new_quotes_converted_count_in_force"] = newQuotesInForceCount
	resultData["new_quotes_converted_count_accepted"] = newQuotesCountAccepted
	resultData["new_quotes_total_count"] = newQuotesCountApproved + newQuotesInForceCount + newQuotesCountAccepted

	resultData["new_quotes_unconverted_premium"] = newQuotesPremiumApproved
	resultData["new_quotes_converted_premium_in_force"] = newQuotesPremiumInForce
	resultData["new_quotes_converted_premium_accepted"] = newQuotesPremiumAccepted
	resultData["new_quotes_total_premium"] = newQuotesPremiumApproved + newQuotesPremiumAccepted + newQuotesPremiumInForce

	//Renewal Rates
	resultData["renewals_quotes_unconverted_count"] = renewalQuotesCountApproved
	resultData["renewals_quotes_converted_count_in_force"] = renewalQuotesInForceCount
	resultData["renewals_quotes_converted_count_accepted"] = renewalQuotesCountAccepted
	resultData["renewals_quotes_total_count"] = renewalQuotesCountApproved + renewalQuotesCountAccepted + renewalQuotesInForceCount

	resultData["renewals_quotes_unconverted_premium"] = renewalQuotesPremiumApproved
	resultData["renewals_quotes_converted_premium_in_force"] = renewalQuotesPremiumInForce
	resultData["renewals_quotes_converted_premium_accepted"] = renewalQuotesPremiumAccepted
	resultData["renewals_quotes_total_premium"] = renewalQuotesPremiumApproved + renewalQuotesPremiumInForce + renewalQuotesPremiumAccepted

	resultData["total_quotes_unconverted_count"] = newQuotesCountApproved + renewalQuotesCountApproved
	resultData["total_quotes_converted_count_in_force"] = newQuotesInForceCount + renewalQuotesInForceCount
	resultData["total_quotes_converted_count_accepted"] = newQuotesCountAccepted + renewalQuotesCountAccepted
	resultData["conversion_total_count"] = newQuotesCountApproved + newQuotesInForceCount +
		newQuotesCountAccepted + renewalQuotesCountApproved +
		renewalQuotesInForceCount + renewalQuotesCountAccepted

	resultData["total_quotes_unconverted_premium"] = newQuotesPremiumApproved + renewalQuotesPremiumApproved
	resultData["total_quotes_converted_premium_in_force"] = newQuotesPremiumInForce + renewalQuotesPremiumInForce
	resultData["total_quotes_converted_premium_accepted"] = newQuotesPremiumAccepted + renewalQuotesPremiumAccepted
	resultData["conversion_total_premium"] = newQuotesPremiumApproved + newQuotesPremiumInForce +
		newQuotesPremiumAccepted + renewalQuotesPremiumApproved +
		renewalQuotesPremiumInForce + renewalQuotesPremiumAccepted
	// Data for cards
	var cds []models.CardData

	var newQuotes int64
	err = DB.Table("member_rating_result_summaries").Distinct("quote_id").Where(" year(creation_date) = ?", year).Count(&newQuotes).Error
	var cd1 models.CardData
	cd1.Title = "New Quotes"
	cd1.Value = newQuotes
	cd1.Flex = 3
	cds = append(cds, cd1)

	var nonFunNewQuotePremium float64
	var funNewQuotePremium float64
	err = DB.Table("member_rating_result_summaries").Where(" year(creation_date) = ? ", year).Pluck("sum(exp_total_annual_premium_excl_funeral)", &nonFunNewQuotePremium).Error
	err = DB.Table("member_rating_result_summaries").Where(" year(creation_date) = ? ", year).Pluck("sum(exp_total_fun_annual_office_premium)", &funNewQuotePremium).Error

	var cd2 models.CardData
	cd2.Title = "Quoted Annual Premium"
	cd2.Value = FloatPrecision(nonFunNewQuotePremium+funNewQuotePremium, AccountingPrecision)
	cd2.DataType = "currency"
	cd2.Flex = 3
	cds = append(cds, cd2)

	// schemes count
	var schemesCount int64
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Count(&schemesCount).Error
	var cd3 models.CardData
	cd3.Title = "Schemes In Force"

	cd3.Value = schemesCount
	cd3.DataType = "number"
	cd3.Flex = 3
	cds = append(cds, cd3)

	// schemes count
	var annualPremium float64
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(annual_premium)", &annualPremium).Error
	var cd4 models.CardData
	cd4.Title = "Gross Written Premium"
	cd4.Value = FloatPrecision(annualPremium, AccountingPrecision)
	cd4.DataType = "currency"
	cd4.Flex = 3
	cds = append(cds, cd4)

	resultData["card_data"] = cds

	// Revenue By Benefit
	// Data for cards
	var rbs []models.RevenueBenefit

	var glaRevenue float64
	var glaClaims float64
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(gla_annual_premium)", &glaRevenue).Error
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(expected_gla_claims)", &glaClaims).Error

	var ptdRevenue float64
	var ptdClaims float64
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(ptd_annual_premium)", &ptdRevenue).Error
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(expected_ptd_claims)", &ptdClaims).Error

	var ciRevenue float64
	var ciClaims float64
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(ci_annual_premium)", &ciRevenue).Error
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(expected_ci_claims)", &ciClaims).Error

	var sglaRevenue float64
	var sglaClaims float64
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(sgla_annual_premium)", &sglaRevenue).Error
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(expected_sgla_claims)", &sglaClaims).Error

	var ttdRevenue float64
	var ttdClaims float64
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(ttd_annual_premium)", &ttdRevenue).Error
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(expected_ttd_claims)", &ttdClaims).Error

	var phiRevenue float64
	var phiClaims float64
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(phi_annual_premium)", &phiRevenue).Error
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(expected_phi_claims)", &phiClaims).Error

	var funeralRevenue float64
	var funeralClaims float64
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(funeral_annual_premium)", &funeralRevenue).Error
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(expected_fun_claims)", &funeralClaims).Error

	var rb models.RevenueBenefit
	rb.Benefit = "GLA"
	rb.Revenue = glaRevenue
	rb.Claims = glaClaims

	var rc models.RevenueBenefit
	rc.Benefit = "PTD"
	rc.Revenue = ptdRevenue
	rc.Claims = ptdClaims

	var rd models.RevenueBenefit
	rd.Benefit = "CI"
	rd.Revenue = ciRevenue
	rd.Claims = ciClaims

	var re models.RevenueBenefit
	re.Benefit = "SGLA"
	re.Revenue = sglaRevenue
	re.Claims = sglaClaims

	var rf models.RevenueBenefit
	rf.Benefit = "TTD"
	rf.Revenue = ttdRevenue
	rf.Claims = ttdClaims

	var rg models.RevenueBenefit
	rg.Benefit = "PHI"
	rg.Revenue = phiRevenue
	rg.Claims = phiClaims

	var rh models.RevenueBenefit
	rh.Benefit = "GFF"
	rh.Revenue = funeralRevenue
	rh.Claims = funeralClaims

	rbs = append(rbs, rb, rc, rd, re, rf, rg, rh)

	resultData["revenue_benefits"] = rbs

	// Income Component
	var gic []models.GroupPricingIncomeComponent

	var gic1 models.GroupPricingIncomeComponent
	var expPremium float64
	var actPremium float64
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(annual_premium)", &expPremium).Error
	gic1.Benefit = "Premium"
	gic1.Expected = expPremium
	gic1.Actual = actPremium

	var gic2 models.GroupPricingIncomeComponent
	var expCommission float64
	var actCommission float64
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(commission)", &expCommission).Error
	gic2.Benefit = "Commission"
	gic2.Expected = expCommission
	gic2.Actual = actCommission

	var gic3 models.GroupPricingIncomeComponent
	var expClaims float64
	var actClaims float64
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(expected_claims)", &expClaims).Error

	err = DB.Table("group_scheme_claims").Where(" year(date_of_event) = ? and status = ?", year, models.StatusApproved).Pluck("sum(claim_amount)", &actClaims).Error
	gic3.Benefit = "Claims"
	gic3.Expected = expClaims
	gic3.Actual = actClaims

	var gic4 models.GroupPricingIncomeComponent
	var expExpenses float64
	var actExpenses float64
	err = DB.Table("group_schemes").Where(" year(creation_date) = ? and in_force = ?", year, true).Pluck("sum(expected_expenses)", &expExpenses).Error
	gic4.Benefit = "Expenses"
	gic4.Expected = expExpenses
	gic4.Actual = actExpenses

	var gic5 models.GroupPricingIncomeComponent
	var expNetIncome float64
	var actNetIncome float64
	expNetIncome = expPremium - expCommission - expExpenses - expClaims - expExpenses
	gic5.Benefit = "Net Income"
	gic5.Expected = expNetIncome
	gic5.Actual = actNetIncome

	gic = append(gic, gic1, gic2, gic3, gic4, gic5)

	resultData["income_statement_components"] = gic

	// ── Monthly Quote Trend ─────────────────────────────────────────────────────
	monthNames := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	var monthlyTrend []models.MonthlyQuoteTrend
	for m := 1; m <= 12; m++ {
		var nbCount, rnCount int64
		var mrrsNB, mrrsRN models.MemberRatingResultSummary

		DB.Table("group_pricing_quotes").
			Where("year(creation_date) = ? AND month(creation_date) = ? AND quote_type = ?", year, m, "New Business").
			Count(&nbCount)
		DB.Table("group_pricing_quotes").
			Where("year(creation_date) = ? AND month(creation_date) = ? AND quote_type = ?", year, m, "Renewal").
			Count(&rnCount)
		DB.Table("member_rating_result_summaries").
			Select("sum(total_annual_premium) as total_annual_premium").
			Where("year(creation_date) = ? AND month(creation_date) = ? AND quote_type = ?", year, m, "New Business").
			Scan(&mrrsNB)
		DB.Table("member_rating_result_summaries").
			Select("sum(total_annual_premium) as total_annual_premium").
			Where("year(creation_date) = ? AND month(creation_date) = ? AND quote_type = ?", year, m, "Renewal").
			Scan(&mrrsRN)

		monthlyTrend = append(monthlyTrend, models.MonthlyQuoteTrend{
			Month:              m,
			MonthName:          monthNames[m-1],
			NewBusinessCount:   nbCount,
			RenewalCount:       rnCount,
			TotalCount:         nbCount + rnCount,
			NewBusinessPremium: mrrsNB.TotalAnnualPremium,
			RenewalPremium:     mrrsRN.TotalAnnualPremium,
		})
	}
	resultData["monthly_quote_trend"] = monthlyTrend

	// ── Quote Funnel (all pipeline stages) ─────────────────────────────────────
	funnelStatuses := []string{"in_progress", "pending_review", "approved", "accepted", "rejected", "not_taken_up", "expired", "cancelled"}
	var funnelStages []models.QuoteFunnelStage
	for _, s := range funnelStatuses {
		var cnt int64
		var funnelMrrs models.MemberRatingResultSummary
		DB.Table("group_pricing_quotes").
			Where("year(creation_date) = ? AND status = ?", year, s).Count(&cnt)
		DB.Table("member_rating_result_summaries").
			Select("sum(total_annual_premium) as total_annual_premium").
			Where("year(creation_date) = ? AND if_status = ?", year, s).Scan(&funnelMrrs)
		funnelStages = append(funnelStages, models.QuoteFunnelStage{
			Stage:   s,
			Count:   cnt,
			Premium: funnelMrrs.TotalAnnualPremium,
		})
	}
	resultData["quote_funnel"] = funnelStages

	// ── Broker Metrics (top 10 by quote volume) ─────────────────────────────────
	type brokerRow struct {
		BrokerID       int     `gorm:"column:broker_id"`
		BrokerName     string  `gorm:"column:broker_name"`
		TotalQuotes    int64   `gorm:"column:total_quotes"`
		AcceptedQuotes int64   `gorm:"column:accepted_quotes"`
		TotalPremium   float64 `gorm:"column:total_premium"`
	}
	var brokerRows []brokerRow
	DB.Table("group_pricing_quotes q").
		Select(`q.broker_id, q.broker_name,
			COUNT(DISTINCT q.id) as total_quotes,
			COUNT(DISTINCT CASE WHEN q.status = 'accepted' THEN q.id END) as accepted_quotes,
			COALESCE(SUM(s.total_annual_premium), 0) as total_premium`).
		Joins("LEFT JOIN member_rating_result_summaries s ON s.quote_id = q.id AND YEAR(s.creation_date) = ?", year).
		Where("YEAR(q.creation_date) = ?", year).
		Group("q.broker_id, q.broker_name").
		Order("total_quotes DESC").
		Limit(10).
		Scan(&brokerRows)

	var brokerMetrics []models.BrokerMetric
	for _, b := range brokerRows {
		var cr float64
		if b.TotalQuotes > 0 {
			cr = float64(b.AcceptedQuotes) / float64(b.TotalQuotes) * 100
		}
		brokerMetrics = append(brokerMetrics, models.BrokerMetric{
			BrokerID:       b.BrokerID,
			BrokerName:     b.BrokerName,
			TotalQuotes:    b.TotalQuotes,
			AcceptedQuotes: b.AcceptedQuotes,
			ConversionRate: FloatPrecision(cr, 1),
			TotalPremium:   FloatPrecision(b.TotalPremium, AccountingPrecision),
		})
	}
	resultData["broker_metrics"] = brokerMetrics

	// ── Pricing Metrics ─────────────────────────────────────────────────────────
	// Fix: GORM snake_cases digits without underscore prefix, so Per1000SA → per1000_sa
	var rateRow struct {
		AvgGla  float64 `gorm:"column:avg_gla"`
		AvgPtd  float64 `gorm:"column:avg_ptd"`
		AvgCi   float64 `gorm:"column:avg_ci"`
		AvgSgla float64 `gorm:"column:avg_sgla"`
		AvgComm float64 `gorm:"column:avg_comm"`
	}
	DB.Table("member_rating_result_summaries").
		Select(`AVG(gla_risk_rate_per1000_sa) as avg_gla,
			AVG(ptd_risk_rate_per1000_sa) as avg_ptd,
			AVG(ci_risk_rate_per1000_sa) as avg_ci,
			AVG(sgla_risk_rate_per1000_sa) as avg_sgla,
			AVG(commission_loading) as avg_comm`).
		Where("year(creation_date) = ?", year).Scan(&rateRow)

	var discRow struct {
		AvgDisc float64 `gorm:"column:avg_disc"`
	}
	DB.Table("group_pricing_quotes").
		Select("AVG(loadings_discount) as avg_disc").
		Where("year(creation_date) = ?", year).Scan(&discRow)

	// ELR: expected_claims / annual_premium from in-force schemes
	var elrRow struct {
		ELR float64 `gorm:"column:elr"`
	}
	DB.Table("group_schemes").
		Select("SUM(expected_claims) / NULLIF(SUM(annual_premium), 0) * 100 as elr").
		Where("year(creation_date) = ? AND in_force = ?", year, true).Scan(&elrRow)

	// ALR: actual approved claim amounts / in-force annual premium
	// group_schemes.actual_claims is not auto-updated; query directly from claims table
	var actClaimsTotal float64
	DB.Table("group_scheme_claims").
		Where("year(date_of_event) = ? AND status = ?", year, models.StatusApproved).
		Select("COALESCE(SUM(claim_amount), 0)").
		Scan(&actClaimsTotal)

	var inForcePremiumTotal float64
	DB.Table("group_schemes").
		Where("year(creation_date) = ? AND in_force = ?", year, true).
		Select("COALESCE(SUM(annual_premium), 0)").
		Scan(&inForcePremiumTotal)

	var alr float64
	if inForcePremiumTotal > 0 {
		alr = actClaimsTotal / inForcePremiumTotal * 100
	}

	// commission_loading in member_rating_result_summaries is stored as a decimal fraction
	// (e.g. 0.15 = 15%) — multiply by 100 before displaying as a percentage
	resultData["pricing_metrics"] = models.DashboardPricingMetrics{
		AvgGlaRatePer1000:  FloatPrecision(rateRow.AvgGla, 4),
		AvgPtdRatePer1000:  FloatPrecision(rateRow.AvgPtd, 4),
		AvgCiRatePer1000:   FloatPrecision(rateRow.AvgCi, 4),
		AvgSglaRatePer1000: FloatPrecision(rateRow.AvgSgla, 4),
		AvgDiscount:        FloatPrecision(discRow.AvgDisc, 2),
		AvgCommissionPct:   FloatPrecision(rateRow.AvgComm*100, 2),
		ExpectedLossRatio:  FloatPrecision(elrRow.ELR, 1),
		ActualLossRatio:    FloatPrecision(alr, 1),
	}

	resultData["win_probability_bands"] = QuoteWinProbabilityBandCounts()

	// ── Exposure by Province ──────────────────────────────────────────────────
	// Count members and sum annual salary per province/region.
	// Uses group_scheme_exposures to identify quotes in the selected financial
	// year (consistent with industry_by_age and age-band exposure charts).
	// Falls back to scheme_category.region when member province is not set —
	// region is always populated, province is optional in uploads.
	type provinceRow struct {
		Province    string  `gorm:"column:region" json:"region"`
		MemberCount int64   `gorm:"column:member_count" json:"member_count"`
		TotalSalary float64 `gorm:"column:total_salary" json:"total_salary"`
	}
	var provinceRows []provinceRow
	// Region comes from scheme_categories (sc.region) — the category-level
	// region is reliably populated even when the member-level province is
	// not. Benefit filter also uses scheme_categories per-benefit flags
	// (sc.gla_benefit, sc.ttd_benefit, etc.) which are set at quote-config
	// time and don't depend on whether exposure rebuild has been run.
	benefitCategoryCol := map[string]string{
		"GLA":  "sc.gla_benefit",
		"SGLA": "sc.sgla_benefit",
		"PTD":  "sc.ptd_benefit",
		"CI":   "sc.ci_benefit",
		"TTD":  "sc.ttd_benefit",
		"PHI":  "sc.phi_benefit",
	}
	provinceQ := DB.Table("g_pricing_member_data m").
		Select(`sc.region,
			COUNT(DISTINCT m.id) as member_count,
			COALESCE(SUM(m.annual_salary), 0) as total_salary`).
		Joins("JOIN group_pricing_quotes q ON q.id = m.quote_id").
		Joins("LEFT JOIN scheme_categories sc ON sc.quote_id = m.quote_id AND sc.scheme_category = m.scheme_category").
		Joins(
			"JOIN (SELECT DISTINCT quote_id FROM group_scheme_exposures WHERE financial_year = ?) gse ON gse.quote_id = m.quote_id",
			year,
		).
		Where("sc.region IS NOT NULL AND sc.region != ''")
	if benefit != "" && benefit != "All" {
		if col, ok := benefitCategoryCol[benefit]; ok {
			provinceQ = provinceQ.Where(col+" = ?", true)
		}
	}
	switch dataSource {
	case "inforce":
		provinceQ = provinceQ.Where("q.status IN ?", []string{"accepted", "in_force"})
	case "quotes":
		provinceQ = provinceQ.Where("q.status NOT IN ?", []string{"accepted", "in_force"})
	}
	provinceQ.Group("sc.region").Order("member_count DESC").Scan(&provinceRows)
	if provinceRows == nil {
		provinceRows = []provinceRow{}
	}
	resultData["exposure_by_province"] = provinceRows

	// ── Industry by Age ───────────────────────────────────────────────────────
	// Sum assured cross-tabulated by industry and age band.
	// Uses group_scheme_exposures (filtered by financial_year) for consistency
	// with the age-band and gender exposure charts.
	type industryAgeRow struct {
		Industry        string  `gorm:"column:industry" json:"industry"`
		AgeBand         string  `gorm:"column:age_band" json:"age_band"`
		TotalSumAssured float64 `gorm:"column:total_sum_assured" json:"total_sum_assured"`
		RecordCount     int64   `gorm:"column:record_count" json:"record_count"`
		MinAge          int     `gorm:"column:min_age" json:"min_age"`
	}
	var industryAgeRows []industryAgeRow
	industryAgeQ := DB.Table("group_scheme_exposures").
		Select(`industry, age_band, MIN(min_age) as min_age,
			COUNT(*) as record_count, SUM(total_sum_assured) as total_sum_assured`).
		Where("financial_year = ? AND industry IS NOT NULL AND industry != ''", year)
	if benefit != "" && benefit != "All" {
		industryAgeQ = industryAgeQ.Where("benefit = ?", benefit)
	}
	switch dataSource {
	case "inforce":
		industryAgeQ = industryAgeQ.Where("quote_status IN ?", []string{"accepted", "in_force"})
	case "quotes":
		industryAgeQ = industryAgeQ.Where("quote_status NOT IN ?", []string{"accepted", "in_force"})
	}
	industryAgeQ.Group("industry, age_band").Order("industry, MIN(min_age) ASC").Scan(&industryAgeRows)
	if industryAgeRows == nil {
		industryAgeRows = []industryAgeRow{}
	}
	resultData["industry_by_age"] = industryAgeRows

	// ── Monthly Conversion Rate Trend ─────────────────────────────────────────
	// For each month compute: accepted / (accepted + not_taken_up + rejected +
	// expired + cancelled) — i.e. closed-won / all-closed — for NB and Renewal.
	type monthlyConvRow struct {
		Month       int     `json:"month"`
		MonthName   string  `json:"month_name"`
		NBConvRate  float64 `json:"nb_conv_rate"`
		RenConvRate float64 `json:"ren_conv_rate"`
		NBTotal     int64   `json:"nb_total"`
		RenTotal    int64   `json:"ren_total"`
		NBAccepted  int64   `json:"nb_accepted"`
		RenAccepted int64   `json:"ren_accepted"`
	}
	terminalStatuses := []string{"accepted", "not_taken_up", "rejected", "expired", "cancelled"}
	convMonthNames := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	var monthlyConvRows []monthlyConvRow
	for m := 1; m <= 12; m++ {
		var nbTotal, nbAccepted, renTotal, renAccepted int64
		DB.Table("group_pricing_quotes").
			Where("year(creation_date) = ? AND month(creation_date) = ? AND quote_type = ? AND status IN ?",
				year, m, "New Business", terminalStatuses).
			Count(&nbTotal)
		DB.Table("group_pricing_quotes").
			Where("year(creation_date) = ? AND month(creation_date) = ? AND quote_type = ? AND status = ?",
				year, m, "New Business", "accepted").
			Count(&nbAccepted)
		DB.Table("group_pricing_quotes").
			Where("year(creation_date) = ? AND month(creation_date) = ? AND quote_type = ? AND status IN ?",
				year, m, "Renewal", terminalStatuses).
			Count(&renTotal)
		DB.Table("group_pricing_quotes").
			Where("year(creation_date) = ? AND month(creation_date) = ? AND quote_type = ? AND status = ?",
				year, m, "Renewal", "accepted").
			Count(&renAccepted)

		var nbRate, renRate float64
		if nbTotal > 0 {
			nbRate = FloatPrecision(float64(nbAccepted)/float64(nbTotal)*100, 1)
		}
		if renTotal > 0 {
			renRate = FloatPrecision(float64(renAccepted)/float64(renTotal)*100, 1)
		}
		monthlyConvRows = append(monthlyConvRows, monthlyConvRow{
			Month:       m,
			MonthName:   convMonthNames[m-1],
			NBConvRate:  nbRate,
			RenConvRate: renRate,
			NBTotal:     nbTotal,
			RenTotal:    renTotal,
			NBAccepted:  nbAccepted,
			RenAccepted: renAccepted,
		})
	}
	resultData["monthly_conversion_trend"] = monthlyConvRows

	// ── Industry Mix of Active Pipeline ──────────────────────────────────────
	// Active = quotes in a non-terminal stage (in_progress, pending_review, approved).
	type industryPipelineRow struct {
		Industry     string  `gorm:"column:industry" json:"industry"`
		QuoteCount   int64   `gorm:"column:quote_count" json:"quote_count"`
		TotalPremium float64 `gorm:"column:total_premium" json:"total_premium"`
	}
	activeStatuses := []string{"in_progress", "pending_review", "approved"}
	var industryPipelineRows []industryPipelineRow
	DB.Table("group_pricing_quotes q").
		Select(`COALESCE(NULLIF(q.industry, ''), 'Unspecified') as industry,
			COUNT(*) as quote_count,
			COALESCE(SUM(s.total_annual_premium), 0) as total_premium`).
		Joins("LEFT JOIN member_rating_result_summaries s ON s.quote_id = q.id AND YEAR(s.creation_date) = ?", year).
		Where("YEAR(q.creation_date) = ? AND q.status IN ? AND q.industry IS NOT NULL", year, activeStatuses).
		Group("COALESCE(NULLIF(q.industry, ''), 'Unspecified')").
		Order("total_premium DESC").
		Scan(&industryPipelineRows)
	if industryPipelineRows == nil {
		industryPipelineRows = []industryPipelineRow{}
	}
	resultData["industry_pipeline"] = industryPipelineRows

	// ── Scheme Size Distribution ──────────────────────────────────────────────
	// Bucket active pipeline quotes by number of rated lives into size bands.
	type schemeSizeBand struct {
		SizeBand     string  `json:"size_band"`
		QuoteCount   int64   `json:"quote_count"`
		TotalPremium float64 `json:"total_premium"`
		SortOrder    int     `json:"sort_order"`
	}
	type liveCountRow struct {
		QuoteID      int     `gorm:"column:quote_id"`
		LiveCount    int64   `gorm:"column:live_count"`
		TotalPremium float64 `gorm:"column:total_premium"`
	}
	var liveCountRows []liveCountRow
	DB.Table("group_pricing_quotes q").
		Select(`q.id as quote_id,
			COUNT(m.id) as live_count,
			COALESCE(s.total_annual_premium, 0) as total_premium`).
		Joins("LEFT JOIN g_pricing_member_data m ON m.quote_id = q.id").
		Joins("LEFT JOIN member_rating_result_summaries s ON s.quote_id = q.id AND YEAR(s.creation_date) = ?", year).
		Where("YEAR(q.creation_date) = ? AND q.status IN ?", year, activeStatuses).
		Group("q.id, s.total_annual_premium").
		Scan(&liveCountRows)

	sizeBandMap := map[string]*schemeSizeBand{
		"1–50":    {SizeBand: "1–50", SortOrder: 1},
		"51–100":  {SizeBand: "51–100", SortOrder: 2},
		"101–250": {SizeBand: "101–250", SortOrder: 3},
		"251–500": {SizeBand: "251–500", SortOrder: 4},
		"500+":    {SizeBand: "500+", SortOrder: 5},
	}
	for _, r := range liveCountRows {
		var band string
		switch {
		case r.LiveCount <= 50:
			band = "1–50"
		case r.LiveCount <= 100:
			band = "51–100"
		case r.LiveCount <= 250:
			band = "101–250"
		case r.LiveCount <= 500:
			band = "251–500"
		default:
			band = "500+"
		}
		sizeBandMap[band].QuoteCount++
		sizeBandMap[band].TotalPremium += r.TotalPremium
	}
	sizeBands := []schemeSizeBand{
		*sizeBandMap["1–50"],
		*sizeBandMap["51–100"],
		*sizeBandMap["101–250"],
		*sizeBandMap["251–500"],
		*sizeBandMap["500+"],
	}
	resultData["scheme_size_distribution"] = sizeBands

	return resultData, nil
}

func GetGroupSchemeExposureData(year int, benefit string, dataSource string) ([]models.GroupSchemeExposure, error) {
	var exposures []models.GroupSchemeExposure

	// Build base query
	q := DB.Table("group_scheme_exposures").
		Select("age_band, SUM(total_sum_assured) AS total_sum_assured, SUM(male_sum_assured) AS male_sum_assured, SUM(female_sum_assured) AS female_sum_assured").
		Where("financial_year = ?", year)

	if benefit != "All" {
		q = q.Where("benefit = ?", benefit)
	}

	switch dataSource {
	case "inforce":
		q = q.Where("quote_status IN ?", []string{"accepted", "in_force"})
	case "quotes":
		q = q.Where("quote_status NOT IN ?", []string{"accepted", "in_force"})
	}

	err := q.Group("age_band").Scan(&exposures).Error
	if err != nil {
		return nil, err
	}
	if exposures == nil {
		exposures = []models.GroupSchemeExposure{}
	}
	return exposures, nil
}

// RebuildExposureDataForYear re-computes group_scheme_exposures rows for every
// quote whose member rating summary falls in the given financial year.
// It also refreshes quote_status so the column stays current.
func RebuildExposureDataForYear(year int) (int, error) {
	// Find all quote IDs that have rating data for this financial year
	var quoteIDs []int
	if err := DB.Table("member_rating_result_summaries").
		Select("DISTINCT quote_id").
		Where("financial_year = ?", year).
		Pluck("quote_id", &quoteIDs).Error; err != nil {
		return 0, err
	}
	if len(quoteIDs) == 0 {
		return 0, nil
	}

	groupAgeBands, err := GetGroupPricingAgeBands(context.Background())
	if err != nil {
		return 0, err
	}
	groupBenefits, err := GetGroupPricingBenefits(context.Background())
	if err != nil {
		return 0, err
	}

	type BenefitSum struct {
		AgeBand                   string  `gorm:"column:age_band"`
		Gender                    string  `gorm:"column:gender"`
		GlaCappedSumAssured       float64 `gorm:"column:gla_capped_sum_assured"`
		PtdCappedSumAssured       float64 `gorm:"column:ptd_capped_sum_assured"`
		CiCappedSumAssured        float64 `gorm:"column:ci_capped_sum_assured"`
		SpouseGlaCappedSumAssured float64 `gorm:"column:spouse_gla_capped_sum_assured"`
		TtdCappedIncome           float64 `gorm:"column:ttd_capped_income"`
		PhiCappedIncome           float64 `gorm:"column:phi_capped_income"`
	}

	processed := 0
	for _, quoteID := range quoteIDs {
		var quote models.GroupPricingQuote
		if err := DB.First(&quote, quoteID).Error; err != nil {
			continue
		}
		var scheme models.GroupScheme
		DB.Where("quote_id = ?", quoteID).First(&scheme)

		var memberRatingSummary []models.MemberRatingResultSummary
		DB.Where("quote_id = ?", quoteID).Find(&memberRatingSummary)

		financialYear := year
		if len(memberRatingSummary) > 0 {
			financialYear = memberRatingSummary[0].FinancialYear
		}

		var benefitSums []BenefitSum
		DB.Table("member_rating_results").
			Select("age_band, gender, SUM(gla_capped_sum_assured) as gla_capped_sum_assured, SUM(ptd_capped_sum_assured) as ptd_capped_sum_assured, SUM(ci_capped_sum_assured) as ci_capped_sum_assured, SUM(spouse_gla_capped_sum_assured) as spouse_gla_capped_sum_assured, SUM(ttd_capped_income) as ttd_capped_income, SUM(phi_capped_income) as phi_capped_income").
			Where("quote_id = ?", quoteID).
			Group("age_band, gender").
			Scan(&benefitSums)

		sumMap := make(map[string]map[string]BenefitSum)
		for _, bs := range benefitSums {
			if _, ok := sumMap[bs.AgeBand]; !ok {
				sumMap[bs.AgeBand] = make(map[string]BenefitSum)
			}
			sumMap[bs.AgeBand][bs.Gender] = bs
		}

		var exposureData []models.GroupSchemeExposure
		for _, ageBand := range groupAgeBands {
			for _, groupBenefit := range groupBenefits {
				ep := models.GroupSchemeExposure{
					QuoteId:       quoteID,
					SchemeName:    scheme.Name,
					Industry:      quote.Industry,
					AgeBand:       ageBand.Name,
					MinAge:        ageBand.MinAge,
					MaxAge:        ageBand.MaxAge,
					Benefit:       groupBenefit.Name,
					FinancialYear: financialYear,
					QuoteStatus:   string(quote.Status),
				}
				maleBS := sumMap[ageBand.Name]["M"]
				femaleBS := sumMap[ageBand.Name]["F"]
				var maleVal, femaleVal float64
				switch groupBenefit.Name {
				case "GLA":
					maleVal = maleBS.GlaCappedSumAssured
					femaleVal = femaleBS.GlaCappedSumAssured
				case "PTD":
					maleVal = maleBS.PtdCappedSumAssured
					femaleVal = femaleBS.PtdCappedSumAssured
				case "CI":
					maleVal = maleBS.CiCappedSumAssured
					femaleVal = femaleBS.CiCappedSumAssured
				case "SGLA":
					maleVal = maleBS.SpouseGlaCappedSumAssured
					femaleVal = femaleBS.SpouseGlaCappedSumAssured
				case "TTD":
					maleVal = maleBS.TtdCappedIncome
					femaleVal = femaleBS.TtdCappedIncome
				case "PHI":
					maleVal = maleBS.PhiCappedIncome
					femaleVal = femaleBS.PhiCappedIncome
				}
				ep.MaleSumAssured = maleVal
				ep.FemaleSumAssured = femaleVal
				ep.TotalSumAssured = maleVal + femaleVal
				exposureData = append(exposureData, ep)
			}
		}

		DB.Where("quote_id = ? AND financial_year = ?", quoteID, financialYear).Delete(&models.GroupSchemeExposure{})
		if len(exposureData) > 0 {
			DB.CreateInBatches(&exposureData, 100)
		}
		processed++
	}
	return processed, nil
}

// ExposureTimeSeriesRow represents one data point in the year-over-year exposure trend.
type ExposureTimeSeriesRow struct {
	FinancialYear    int     `gorm:"column:financial_year" json:"financial_year"`
	AgeBand          string  `gorm:"column:age_band" json:"age_band"`
	MinAge           int     `gorm:"column:min_age" json:"min_age"`
	TotalSumAssured  float64 `gorm:"column:total_sum_assured" json:"total_sum_assured"`
	MaleSumAssured   float64 `gorm:"column:male_sum_assured" json:"male_sum_assured"`
	FemaleSumAssured float64 `gorm:"column:female_sum_assured" json:"female_sum_assured"`
}

// GetExposureTimeSeries returns sum assured grouped by financial_year and age_band,
// optionally filtered by benefit and data source (all / inforce / quotes).
func GetExposureTimeSeries(benefit string, dataSource string) ([]ExposureTimeSeriesRow, error) {
	var rows []ExposureTimeSeriesRow

	q := DB.Table("group_scheme_exposures").
		Select("financial_year, age_band, MIN(min_age) as min_age, SUM(total_sum_assured) AS total_sum_assured, SUM(male_sum_assured) AS male_sum_assured, SUM(female_sum_assured) AS female_sum_assured")

	if benefit != "" && benefit != "All" {
		q = q.Where("benefit = ?", benefit)
	}
	switch dataSource {
	case "inforce":
		q = q.Where("quote_status IN ?", []string{"accepted", "in_force"})
	case "quotes":
		q = q.Where("quote_status NOT IN ?", []string{"accepted", "in_force"})
	}

	err := q.Group("financial_year, age_band").
		Order("financial_year ASC, MIN(min_age) ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []ExposureTimeSeriesRow{}
	}
	return rows, nil
}

// FinancialYearInfo describes the insurer's financial year for a given calendar year.
type FinancialYearInfo struct {
	YearEndMonth       int    `json:"year_end_month"`
	FinancialYearLabel string `json:"financial_year_label"`
	FinancialYearStart int    `json:"financial_year_start"`
	FinancialYearEnd   int    `json:"financial_year_end"`
}

// GetFinancialYearInfo computes the insurer's financial year boundaries for the
// given calendar year using the configured YearEndMonth.
func GetFinancialYearInfo(calendarYear int) FinancialYearInfo {
	insurer, err := GetInsurerDetails()
	yearEndMonth := 12
	if err == nil && insurer.YearEndMonth > 0 {
		yearEndMonth = insurer.YearEndMonth
	}

	var startYear, endYear int

	if yearEndMonth == 12 {
		// December year-end: the financial year IS the calendar year (Jan–Dec).
		// getGroupRiskQuotingFinancialYear would incorrectly return startYear-1
		// because every month satisfies month <= 12, so we handle it directly.
		startYear = calendarYear
		endYear = calendarYear
	} else {
		// For all other year-end months use a mid-year reference date so the
		// calc returns the financial year that contains the majority of
		// calendarYear's months.
		refDate := time.Date(calendarYear, time.July, 1, 0, 0, 0, 0, time.UTC)
		startYear, endYear = getGroupRiskQuotingFinancialYear(refDate, yearEndMonth)
	}

	startMonth := time.Month(yearEndMonth%12 + 1)
	endMonth := time.Month(yearEndMonth)

	label := fmt.Sprintf("FY%d (%s %d – %s %d)",
		endYear,
		startMonth.String()[:3], startYear,
		endMonth.String()[:3], endYear,
	)

	return FinancialYearInfo{
		YearEndMonth:       yearEndMonth,
		FinancialYearLabel: label,
		FinancialYearStart: startYear,
		FinancialYearEnd:   endYear,
	}
}

func validateClaim(claim models.GroupSchemeClaim) error {
	if strings.TrimSpace(claim.MemberIDNumber) == "" {
		return errors.New("member ID number is required")
	}
	if strings.TrimSpace(claim.SchemeName) == "" && claim.SchemeId == 0 {
		return errors.New("scheme name or ID is required")
	}
	if strings.TrimSpace(claim.BenefitCode) == "" && strings.TrimSpace(claim.BenefitAlias) == "" {
		return errors.New("benefit code or alias is required")
	}
	if strings.TrimSpace(claim.DateOfEvent) == "" {
		return errors.New("date of event is required")
	}
	return nil
}

func GroupSchemeSubmitClaim(claim models.GroupSchemeClaim, user models.AppUser) error {
	if err := validateClaim(claim); err != nil {
		return err
	}

	claim.CreationDate = time.Now()
	claim.CreatedBy = user.UserName
	if strings.TrimSpace(claim.Status) == "" {
		claim.Status = "Pending"
	}
	if strings.TrimSpace(claim.ClaimNumber) == "" {
		// generate a claim number if missing
		cn, err := generateUniqueClaimNumber()
		if err != nil {
			return fmt.Errorf("failed to generate unique claim number: %w", err)
		}
		claim.ClaimNumber = cn
	}
	if err := DB.Create(&claim).Error; err != nil {
		return fmt.Errorf("failed to create claim record: %w", err)
	}
	// record initial status audit
	if err := DB.Create(&models.GroupSchemeClaimStatusAudit{
		ClaimID:       claim.ID,
		OldStatus:     "",
		NewStatus:     claim.Status,
		StatusMessage: "Claim submitted",
		ChangedBy:     user.UserName,
		ChangedAt:     time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("failed to create status audit: %w", err)
	}

	// Log structured claim activity
	claimDetails, _ := json.Marshal(map[string]interface{}{
		"claimNumber":  claim.ClaimNumber,
		"claimType":    claim.BenefitAlias,
		"amount":       claim.ClaimAmount,
		"status":       claim.Status,
		"incidentDate": claim.DateOfEvent,
	})
	// We need to find the internal member ID if possible.
	// claim.MemberIDNumber is what we have.
	var member models.GPricingMemberDataInForce
	_ = DB.Where("member_id_number = ?", claim.MemberIDNumber).First(&member)

	_ = DB.Create(&models.MemberActivity{
		MemberID:       member.ID,
		MemberIDNumber: claim.MemberIDNumber,
		Type:           "claim",
		Title:          "Claim Submitted",
		Description:    fmt.Sprintf("%s claim submitted", claim.BenefitAlias),
		Details:        claimDetails,
		PerformedBy:    user.UserName,
	})

	return nil
}

// GroupSchemeSubmitClaimWithFiles handles multipart submissions with supporting documents.
// It stores the claim record and saves any uploaded files to disk under tmp/uploads/group_claims.
func GroupSchemeSubmitClaimWithFiles(claim models.GroupSchemeClaim, files map[string][]*multipart.FileHeader, user models.AppUser) error {
	if err := validateClaim(claim); err != nil {
		return err
	}

	// Save the claim first (now we need the ID for attachments)
	claim.CreationDate = time.Now()
	claim.CreatedBy = user.UserName
	if strings.TrimSpace(claim.Status) == "" {
		claim.Status = "Pending"
	}
	if strings.TrimSpace(claim.ClaimNumber) == "" {
		cn, err := generateUniqueClaimNumber()
		if err != nil {
			return fmt.Errorf("failed to generate unique claim number: %w", err)
		}
		claim.ClaimNumber = cn
	}
	if err := DB.Create(&claim).Error; err != nil {
		return fmt.Errorf("failed to create claim record: %w", err)
	}

	if err := DB.Create(&models.GroupSchemeClaimStatusAudit{
		ClaimID:       claim.ID,
		OldStatus:     "",
		NewStatus:     claim.Status,
		StatusMessage: "Claim submitted",
		ChangedBy:     user.UserName,
		ChangedAt:     time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("failed to create status audit: %w", err)
	}

	// Log structured claim activity
	claimDetails, _ := json.Marshal(map[string]interface{}{
		"claimNumber":  claim.ClaimNumber,
		"claimType":    claim.BenefitAlias,
		"amount":       claim.ClaimAmount,
		"status":       claim.Status,
		"incidentDate": claim.DateOfEvent,
	})
	var member models.GPricingMemberDataInForce
	_ = DB.Where("member_id_number = ?", claim.MemberIDNumber).First(&member)

	_ = DB.Create(&models.MemberActivity{
		MemberID:       member.ID,
		MemberIDNumber: claim.MemberIDNumber,
		Type:           "claim",
		Title:          "Claim Submitted",
		Description:    fmt.Sprintf("%s claim submitted", claim.BenefitAlias),
		Details:        claimDetails,
		PerformedBy:    user.UserName,
	})

	if len(claim.SupportingDocuments) == 0 {
		return nil
	}

	// Build a destination directory using the claim ID
	baseDir := filepath.Join("tmp", "uploads", "group_claims", fmt.Sprintf("claim_%d", claim.ID))
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory for attachments: %w", err)
	}

	// Save files and create attachment records based on supporting documents metadata
	for _, doc := range claim.SupportingDocuments {
		key := "files"
		fhList, ok := files[key]
		if !ok || len(fhList) == 0 {
			// Fallback to legacy key format if "files" is not present
			key = fmt.Sprintf("file_%d", doc.FileIndex)
			fhList, ok = files[key]
			if !ok || len(fhList) == 0 {
				continue
			}
		}

		// Ensure the requested index is within bounds of the "files" array
		if key == "files" && (doc.FileIndex < 0 || doc.FileIndex >= len(fhList)) {
			continue
		}

		var fh *multipart.FileHeader
		if key == "files" {
			fh = fhList[doc.FileIndex]
		} else {
			fh = fhList[0]
		}

		// Sanitize filename minimally by removing path components
		name := filepath.Base(fh.Filename)
		destPath := filepath.Join(baseDir, name)

		if err := func() error {
			src, err := fh.Open()
			if err != nil {
				return fmt.Errorf("failed to open uploaded file %s: %w", name, err)
			}
			defer src.Close()

			dst, err := os.Create(destPath)
			if err != nil {
				return fmt.Errorf("failed to create destination file %s: %w", name, err)
			}
			defer dst.Close()

			if _, err = io.Copy(dst, src); err != nil {
				return fmt.Errorf("failed to save file %s: %w", name, err)
			}
			return nil
		}(); err != nil {
			return err
		}

		att := models.GroupSchemeClaimAttachment{
			ClaimID:      claim.ID,
			DocumentType: doc.DocumentType,
			DocumentName: doc.DocumentName,
			FileName:     name,
			ContentType:  fh.Header.Get("Content-Type"),
			SizeBytes:    fh.Size,
			StoragePath:  destPath,
			UploadedAt:   time.Now(),
			UploadedBy:   user.UserName,
		}
		if err := DB.Create(&att).Error; err != nil {
			return fmt.Errorf("failed to save attachment record for %s: %w", name, err)
		}
	}

	return nil
}

// Read Accelerated claim Amount
func GetAcceleratedApprovedClaims(memberIdNumber string, schemeId int, claimType string, benefitCode string) float64 {
	//if !strings.EqualFold(claimType, "member") {
	//	return 0 // The only member_type this affects is the main member
	//}

	var member models.GPricingMemberDataInForce
	var scheme models.GroupScheme
	var quote models.GroupPricingQuote
	var totalAmount float64
	var benefitTypes []string

	DB.Where("scheme_id = ? and scheme_quote_status = ? ", schemeId, models.StatusInEffect).First(&quote)

	DB.Where("id = ? and scheme_quote_status = ? ", schemeId, models.StatusInEffect).First(&scheme)

	if err := DB.Where("member_id_number = ? and scheme_name =?", memberIdNumber, scheme.Name).First(&member).Error; err != nil {
		fmt.Println(err)
	}

	var cat models.SchemeCategory
	if len(memberIdNumber) != 0 && strings.TrimSpace(member.SchemeCategory) != "" {
		if err := DB.Where("quote_id = ? AND LOWER(scheme_category) = LOWER(?)", quote.ID, member.SchemeCategory).First(&cat).Error; err != nil {
			fmt.Println(err)
		}
	}

	restriction, restrictionErr := GetRestrictionByRiskRateCode(quote.RiskRateCode)
	if restrictionErr != nil {
		appLog.Error("Error retrieving restriction for risk rate code: ", quote.RiskRateCode, " error: ", restrictionErr.Error())
	}

	// Calculate Sum Assured for the requested benefitCode
	var sumAssured float64
	annualSalary := member.AnnualSalary

	switch {
	case strings.Contains(benefitCode, "PTD"):
		mult := member.Benefits.PtdMultiple
		if mult == 0 {
			mult = member.Benefits.PtdMultiple
		}
		if member.Benefits.PtdEnabled && mult > 0 {
			sumAssured = math.Min(annualSalary*mult, cat.FreeCoverLimit)
		}
	case strings.Contains(benefitCode, "CI"):
		mult := member.Benefits.CiMultiple
		if mult == 0 {
			mult = member.Benefits.CiMultiple
		}
		if member.Benefits.CiEnabled && mult > 0 {
			sumAssured = math.Min(math.Min(annualSalary*mult, restriction.SevereIllnessMaximumBenefit), cat.FreeCoverLimit)
		}
	case strings.Contains(benefitCode, "GLA"):
		mult := member.Benefits.GlaMultiple
		if mult == 0 {
			mult = member.Benefits.GlaMultiple
		}
		if member.Benefits.GlaEnabled && mult > 0 {
			sumAssured = math.Min(annualSalary*mult, cat.FreeCoverLimit)
		}
	case strings.Contains(benefitCode, "SGLA"):
		mult := member.Benefits.SglaMultiple
		if mult == 0 {
			mult = member.Benefits.SglaMultiple
		}
		if member.Benefits.SglaEnabled && mult > 0 {
			sumAssured = math.Min(math.Min(annualSalary*mult, restriction.SpouseGlaMaximumBenefit), cat.FreeCoverLimit)
		}
	case strings.Contains(benefitCode, "PHI"):
		mult := member.Benefits.PhiMultiple
		if member.Benefits.PhiEnabled && mult > 0 {
			sumAssured = math.Min(annualSalary*mult, cat.FreeCoverLimit)
		}
	case strings.Contains(benefitCode, "TTD"):
		mult := member.Benefits.TtdMultiple
		if member.Benefits.TtdEnabled && mult > 0 {
			sumAssured = math.Min(annualSalary*mult, cat.FreeCoverLimit)
		}
	case strings.Contains(benefitCode, "GFF"):
		if claimType == "member" && member.Benefits.GffEnabled {
			sumAssured = cat.FamilyFuneralMainMemberFuneralSumAssured //member.SchemeCategoryDetails.FamilyFuneralMainMemberFuneralSumAssured
		}
		if claimType == "spouse" && member.Benefits.GffEnabled {
			sumAssured = cat.FamilyFuneralSpouseFuneralSumAssured //member.SchemeCategoryDetails.FamilyFuneralSpouseFuneralSumAssured
		}
		if claimType == "child" && member.Benefits.GffEnabled {
			sumAssured = cat.FamilyFuneralChildrenFuneralSumAssured //member.SchemeCategoryDetails.FamilyFuneralChildrenFuneralSumAssured
		}
		if claimType == "parent" && member.Benefits.GffEnabled {
			sumAssured = cat.FamilyFuneralParentFuneralSumAssured //member.SchemeCategoryDetails.FamilyFuneralParentFuneralSumAssured
		}

		if claimType == "dependant" && member.Benefits.GffEnabled {
			sumAssured = cat.FamilyFuneralAdultDependantSumAssured //member.SchemeCategoryDetails.FamilyFuneralAdultDependantSumAssured
		}
	}

	if cat.PtdBenefitType == "Accelerated" || cat.CiBenefitStructure == "Accelerated" {
		// Check for Accelerated PTD
		if cat.PtdBenefitType == "Accelerated" {
			benefitTypes = append(benefitTypes, "Permanent Total Disability (PTD)")
		}

		// Check for Accelerated CI
		if cat.CiBenefitStructure == "Accelerated" {
			benefitTypes = append(benefitTypes, "Critical Illness (CI)")
		}

		if len(memberIdNumber) != 0 && cat.PtdBenefitType != "Standalone" {
			// We use .Row().Scan() to capture the sum directly from the database
			err := DB.Model(&models.GroupSchemeClaim{}).
				Where("member_id_number = ?", memberIdNumber).
				Where("benefit_type IN ?", benefitTypes).
				Where("status = ?", models.StatusApproved).
				Where("date_of_event BETWEEN ? AND ?", scheme.CoverStartDate, scheme.CoverEndDate).
				Select("COALESCE(sum(claim_amount), 0)").
				Row().
				Scan(&totalAmount)

			if err != nil {
				return sumAssured // If error fetching claims, return full sum assured
			}
		}
	}

	return sumAssured - totalAmount
}

// GroupSchemeSubmitClaimsBatch processes multiple claim submissions in one call.
// It generates a claim number for any item missing one and records initial status audits.
func GroupSchemeSubmitClaimsBatch(claims []models.GroupSchemeClaim, user models.AppUser) ([]models.GroupSchemeClaim, error) {
	if len(claims) == 0 {
		return []models.GroupSchemeClaim{}, nil
	}
	now := time.Now()
	err := DB.Transaction(func(tx *gorm.DB) error {
		for i := range claims {
			c := &claims[i]

			if err := validateClaim(*c); err != nil {
				return fmt.Errorf("validation failed for claim at index %d: %w", i, err)
			}

			c.CreationDate = now
			c.CreatedBy = user.UserName
			if strings.TrimSpace(c.Status) == "" {
				c.Status = "Pending"
			}
			if strings.TrimSpace(c.ClaimNumber) == "" {
				cn, err := generateUniqueClaimNumberWithTX(tx)
				if err != nil {
					return fmt.Errorf("failed to generate unique claim number for claim at index %d: %w", i, err)
				}
				c.ClaimNumber = cn
			}

			// With quote resolved and member's SchemeCategory name, fetch details
			var member models.GPricingMemberDataInForce
			var scheme models.GroupScheme

			if err := tx.Where("member_id_number = ? and scheme_name =?", c.MemberIDNumber, c.SchemeName).First(&member).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("member with ID %s not found in scheme %s for claim at index %d", c.MemberIDNumber, c.SchemeName, i)
				}
				return fmt.Errorf("failed to fetch member details for claim at index %d: %w", i, err)
			}

			if err := tx.Where("name = ? and scheme_quote_status =?", c.SchemeName, models.StatusInEffect).First(&scheme).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("in-force scheme %s not found for claim at index %d", c.SchemeName, i)
				}
				return fmt.Errorf("failed to fetch scheme details for claim at index %d: %w", i, err)
			}

			var cat models.SchemeCategory
			if len(c.MemberIDNumber) != 0 && strings.TrimSpace(member.SchemeCategory) != "" {
				if err := tx.Where("quote_id = ? AND LOWER(scheme_category) = LOWER(?)", scheme.QuoteId, member.SchemeCategory).First(&cat).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return fmt.Errorf("scheme category %s not found for claim at index %d", member.SchemeCategory, i)
					}
					return fmt.Errorf("failed to fetch scheme category details for claim at index %d: %w", i, err)
				}
			}

			var claimQuote models.GroupPricingQuote
			if err := tx.Where("id = ?", scheme.QuoteId).First(&claimQuote).Error; err != nil {
				return fmt.Errorf("failed to fetch quote for scheme %s at index %d: %w", c.SchemeName, i, err)
			}
			claimRestriction, claimRestrictionErr := GetRestrictionByRiskRateCode(claimQuote.RiskRateCode)
			if claimRestrictionErr != nil {
				appLog.Error("Error retrieving restriction for risk rate code: ", claimQuote.RiskRateCode, " error: ", claimRestrictionErr.Error())
			}

			// Source caps from the per-quote rating summary (frozen at quote-rating
			// time). Fall back to the live restriction for legacy quotes that
			// pre-date the cap-limit / max-cover-age migrations.
			var ratingSummary models.MemberRatingResultSummary
			summaryErr := tx.Where("quote_id = ? AND LOWER(category) = LOWER(?)", scheme.QuoteId, member.SchemeCategory).First(&ratingSummary).Error
			if summaryErr != nil && !errors.Is(summaryErr, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to fetch rating summary for claim at index %d: %w", i, summaryErr)
			}
			hasSummary := summaryErr == nil

			maxGlaCover := claimRestriction.MaximumGlaCover
			maxPtdCover := claimRestriction.MaximumPtdCover
			sevIllMax := claimRestriction.SevereIllnessMaximumBenefit
			spouseGlaMax := claimRestriction.SpouseGlaMaximumBenefit
			ttdMonthlyMax := claimRestriction.TtdMaximumMonthlyBenefit
			phiMonthlyMax := claimRestriction.PhiMaximumMonthlyBenefit
			glaAgeLimit := claimRestriction.GlaMaxCoverAge
			ptdAgeLimit := claimRestriction.PtdMaxCoverAge
			ciAgeLimit := claimRestriction.CiMaxCoverAge
			ttdAgeLimit := claimRestriction.TtdMaxCoverAge
			phiAgeLimit := claimRestriction.PhiMaxCoverAge
			funAgeLimit := claimRestriction.FunMaxCoverAge
			if hasSummary {
				maxGlaCover = ratingSummary.MaximumGlaCover
				maxPtdCover = ratingSummary.MaximumPtdCover
				sevIllMax = ratingSummary.SevereIllnessMaximumBenefit
				spouseGlaMax = ratingSummary.SpouseGlaMaximumBenefit
				ttdMonthlyMax = ratingSummary.TtdMaximumMonthlyBenefit
				phiMonthlyMax = ratingSummary.PhiMaximumMonthlyBenefit
				glaAgeLimit = ratingSummary.GlaMaxCoverAge
				ptdAgeLimit = ratingSummary.PtdMaxCoverAge
				ciAgeLimit = ratingSummary.CiMaxCoverAge
				ttdAgeLimit = ratingSummary.TtdMaxCoverAge
				phiAgeLimit = ratingSummary.PhiMaxCoverAge
				funAgeLimit = ratingSummary.FunMaxCoverAge
			}

			ageNextBirthday := calculateAgeNextBirthday(member.DateOfBirth, scheme.CommencementDate)

			acceleratedApprovedClaims := GetAcceleratedApprovedClaims(c.MemberIDNumber, scheme.ID, c.MemberType, c.BenefitAlias)

			switch c.BenefitCode {
			case "GLA":
				glaSA := applyMaxCoverCap(math.Min(member.AnnualSalary*member.Benefits.GlaMultiple, cat.FreeCoverLimit), maxGlaCover)
				c.ClaimAmount = math.Max(applyCoverAgeLimit(glaSA, ageNextBirthday, glaAgeLimit)-acceleratedApprovedClaims, 0)
			case "SGLA":
				sglaSA := math.Min(member.AnnualSalary*member.Benefits.SglaMultiple, spouseGlaMax)
				c.ClaimAmount = applyCoverAgeLimit(sglaSA, ageNextBirthday, glaAgeLimit)
			case "PTD":
				ptdSA := applyMaxCoverCap(math.Min(member.AnnualSalary*member.Benefits.PtdMultiple, cat.FreeCoverLimit), maxPtdCover)
				c.ClaimAmount = applyCoverAgeLimit(ptdSA, ageNextBirthday, ptdAgeLimit)
			case "CI":
				ciSA := math.Min(member.AnnualSalary*member.Benefits.CiMultiple, sevIllMax)
				c.ClaimAmount = applyCoverAgeLimit(ciSA, ageNextBirthday, ciAgeLimit)
			case "TTD":
				ttdSA := math.Min(member.AnnualSalary*member.Benefits.TtdMultiple, ttdMonthlyMax)
				c.ClaimAmount = applyCoverAgeLimit(ttdSA, ageNextBirthday, ttdAgeLimit)
			case "PHI":
				phiSA := math.Min(member.AnnualSalary*member.Benefits.PhiMultiple, phiMonthlyMax)
				c.ClaimAmount = applyCoverAgeLimit(phiSA, ageNextBirthday, phiAgeLimit)
			case "GFF":
				var funSA float64
				switch c.MemberType {
				case "Member":
					funSA = cat.FamilyFuneralMainMemberFuneralSumAssured
				case "Spouse":
					funSA = cat.FamilyFuneralSpouseFuneralSumAssured
				case "Child":
					funSA = cat.FamilyFuneralChildrenFuneralSumAssured
				case "Parent":
					funSA = cat.FamilyFuneralParentFuneralSumAssured
				case "Dependant":
					funSA = cat.FamilyFuneralAdultDependantSumAssured
				}
				c.ClaimAmount = applyCoverAgeLimit(funSA, ageNextBirthday, funAgeLimit)
			default:
				c.ClaimAmount = 0

			}

			if err := tx.Create(c).Error; err != nil {
				return fmt.Errorf("failed to create claim record at index %d: %w", i, err)
			}
			audit := models.GroupSchemeClaimStatusAudit{
				ClaimID:       c.ID,
				OldStatus:     "",
				NewStatus:     c.Status,
				StatusMessage: "Claim submitted",
				ChangedBy:     user.UserName,
				ChangedAt:     now,
			}
			if err := tx.Create(&audit).Error; err != nil {
				return fmt.Errorf("failed to create status audit for claim at index %d: %w", i, err)
			}

			// Log structured claim activity
			claimDetails, _ := json.Marshal(map[string]interface{}{
				"claimNumber":  c.ClaimNumber,
				"claimType":    c.BenefitAlias,
				"amount":       c.ClaimAmount,
				"status":       c.Status,
				"incidentDate": c.DateOfEvent,
			})

			_ = tx.Create(&models.MemberActivity{
				MemberID:       member.ID,
				MemberIDNumber: c.MemberIDNumber,
				Type:           "claim",
				Title:          "Claim Submitted",
				Description:    fmt.Sprintf("%s claim submitted", c.BenefitAlias),
				Details:        claimDetails,
				PerformedBy:    user.UserName,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// generateUniqueClaimNumber creates a unique claim number of the form GSC-YYYYMMDD-XXXXXX
// using a time component and a random suffix. It checks the database to avoid collisions.
func generateUniqueClaimNumber() (string, error) {
	return generateUniqueClaimNumberWithTX(DB)
}

func generateUniqueClaimNumberWithTX(tx *gorm.DB) (string, error) {
	const maxAttempts = 5
	datePart := time.Now().Format("20060102")
	for i := 0; i < maxAttempts; i++ {
		// 6-digit pseudo-random based on time and loop to reduce collisions
		suffix := fmt.Sprintf("%06d", (time.Now().UnixNano()+int64(i))%1000000)
		candidate := fmt.Sprintf("CLM-%s-%s", datePart, suffix)
		var count int64
		if err := tx.Model(&models.GroupSchemeClaim{}).Where("claim_number = ?", candidate).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return candidate, nil
		}
		// small sleep to change nano time in tight loops
		time.Sleep(1 * time.Millisecond)
	}
	return "", errors.New("could not generate unique claim number after several attempts")
}

func GetGroupSchemeClaims() ([]models.GroupSchemeClaim, error) {
	var claims []models.GroupSchemeClaim
	err := DB.Preload("Attachments").Preload("Assessments").Preload("Communications").Preload("Declines").Preload("StatusAudits").Find(&claims).Error
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// GetGroupSchemeClaimByID fetches a single claim by its ID including attachments and assessments
func GetGroupSchemeClaimByID(claimID int) (models.GroupSchemeClaim, error) {
	var claim models.GroupSchemeClaim
	if err := DB.Preload("Attachments").Preload("Assessments").Preload("Communications").Preload("Declines").Preload("StatusAudits").First(&claim, claimID).Error; err != nil {
		return claim, err
	}
	return claim, nil
}

// GetClaimAttachments returns all attachments for a given claim ID
func GetClaimAttachments(claimID int) ([]models.GroupSchemeClaimAttachment, error) {
	var atts []models.GroupSchemeClaimAttachment
	if err := DB.Where("claim_id = ?", claimID).Find(&atts).Error; err != nil {
		return nil, err
	}
	return atts, nil
}

// GetAttachmentByID loads a single attachment by its primary key
func GetAttachmentByID(attachmentID int) (models.GroupSchemeClaimAttachment, error) {
	var att models.GroupSchemeClaimAttachment
	if err := DB.First(&att, attachmentID).Error; err != nil {
		return att, err
	}
	return att, nil
}

// AppendClaimAttachments saves uploaded files to disk and creates attachment records for a claim.
// Files are stored under tmp/uploads/group_claims/claim_<id> and metadata is persisted in DB.
func AppendClaimAttachments(claimID int, files map[string][]*multipart.FileHeader, fileValues map[string][]string, user models.AppUser) ([]models.GroupSchemeClaimAttachment, error) {
	// Ensure claim exists
	var claim models.GroupSchemeClaim
	if err := DB.First(&claim, claimID).Error; err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return []models.GroupSchemeClaimAttachment{}, nil
	}

	baseDir := filepath.Join("tmp", "uploads", "group_claims", fmt.Sprintf("claim_%d", claim.ID))
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, err
	}

	var created []models.GroupSchemeClaimAttachment

	for _, fhList := range files {
		for _, fh := range fhList {
			name := filepath.Base(fh.Filename)
			destPath := filepath.Join(baseDir, name)

			if err := func() error {
				src, err := fh.Open()
				if err != nil {
					return err
				}
				defer src.Close()

				dst, err := os.Create(destPath)
				if err != nil {
					return err
				}
				defer dst.Close()

				if _, err = io.Copy(dst, src); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return nil, err
			}

			att := models.GroupSchemeClaimAttachment{
				ClaimID:      claim.ID,
				FileName:     name,
				DocumentType: fileValues["document_type"][0],
				DocumentName: fileValues["document_name"][0],
				ContentType:  fh.Header.Get("Content-Type"),
				SizeBytes:    fh.Size,
				StoragePath:  destPath,
				UploadedAt:   time.Now(),
				UploadedBy:   user.UserName,
			}
			if err := DB.Create(&att).Error; err != nil {
				return nil, err
			}
			created = append(created, att)
		}
	}

	return created, nil
}

// ClaimsDashboardFilters defines optional filters for claims analytics
type ClaimsDashboardFilters struct {
	SchemeID *int
	Benefit  string
	From     *time.Time
	To       *time.Time
	// Limit controls the number of rows for the top claims table
	Limit int
}

func getDOBFromSAID(id string) (time.Time, error) {
	if !utils.IsValidRSAID(id) {
		return time.Time{}, errors.New("invalid RSA ID")
	}
	yearPart := id[0:2]
	monthPart := id[2:4]
	dayPart := id[4:6]

	year, _ := strconv.Atoi(yearPart)
	month, _ := strconv.Atoi(monthPart)
	day, _ := strconv.Atoi(dayPart)

	if month < 1 || month > 12 || day < 1 || day > 31 {
		return time.Time{}, errors.New("invalid date in ID")
	}

	// Determine century (simple rule: if year <= current year, assume 2000s, else 1900s)
	// For insurance/claims, claimants are likely < 100 years old.
	currentYearFull := time.Now().Year()
	currentYear := currentYearFull % 100
	fullYear := 1900 + year
	if year <= currentYear {
		fullYear = 2000 + year
	}
	// If calculated DOB is in the future, it must be 1900s (e.g. ID says 99 but it's 2024)
	dob := time.Date(fullYear, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if dob.After(time.Now()) {
		dob = time.Date(fullYear-100, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	}

	return dob, nil
}

// GetGroupSchemeClaimsDashboardData aggregates claims to support the front-end dashboard
// Returned structure keys are stable and UI-friendly.
func GetGroupSchemeClaimsDashboardData(f ClaimsDashboardFilters) (map[string]interface{}, error) {
	// temporary holder for processing efficiency card data
	var _processingEfficiency map[string]interface{}
	// helper to apply filters consistently on a clean query builder
	applyClaimFilters := func(q *gorm.DB, f ClaimsDashboardFilters) *gorm.DB {
		if f.SchemeID != nil {
			q = q.Where("group_scheme_claims.scheme_id = ?", *f.SchemeID)
		}
		if f.Benefit != "" {
			q = q.Where("group_scheme_claims.benefit_type = ?", f.Benefit)
		}
		if f.From != nil {
			q = q.Where("group_scheme_claims.creation_date >= ?", *f.From)
		}
		if f.To != nil {
			q = q.Where("group_scheme_claims.creation_date <= ?", *f.To)
		}
		return q
	}

	// Total claims
	var totalClaims int64
	if err := applyClaimFilters(DB.Model(&models.GroupSchemeClaim{}), f).Count(&totalClaims).Error; err != nil {
		return nil, err
	}

	// Breakdown by status
	type kv struct {
		K string `gorm:"column:k"`
		V int64  `gorm:"column:v"`
	}
	var byStatusRaw []kv
	statusBreakdownQ := applyClaimFilters(DB.Model(&models.GroupSchemeClaim{}), f)
	if err := statusBreakdownQ.Select("LOWER(status) as k, COUNT(*) as v").Group("LOWER(status)").Scan(&byStatusRaw).Error; err != nil {
		return nil, err
	}
	byStatus := map[string]int64{}
	for _, r := range byStatusRaw {
		byStatus[r.K] = r.V
	}

	// Breakdown by benefit type
	var byBenefitRaw []kv
	benefitBreakdownQ := applyClaimFilters(DB.Model(&models.GroupSchemeClaim{}), f)
	if err := benefitBreakdownQ.Select("benefit_alias as k, COUNT(*) as v").Group("benefit_alias").Scan(&byBenefitRaw).Error; err != nil {
		return nil, err
	}
	byBenefit := map[string]int64{}
	for _, r := range byBenefitRaw {
		byBenefit[r.K] = r.V
	}

	// Total paid amount (use status in approved or paid)
	var totalPaid float64
	paidQ := applyClaimFilters(DB.Model(&models.GroupSchemeClaim{}), f)
	if err := paidQ.Where("LOWER(status) IN ?", []string{"approved", "paid"}).Select("COALESCE(SUM(claim_amount),0)").Scan(&totalPaid).Error; err != nil {
		return nil, err
	}

	// Approval rate = (approved + paid)/total
	var approvedOrPaid int64
	apprQ := applyClaimFilters(DB.Model(&models.GroupSchemeClaim{}), f)
	if err := apprQ.Where("LOWER(status) IN ?", []string{"approved", "paid"}).Count(&approvedOrPaid).Error; err != nil {
		return nil, err
	}
	var approvalRate float64
	if totalClaims > 0 {
		approvalRate = float64(approvedOrPaid) / float64(totalClaims)
	}

	// Decline rate = declined / total
	var declinedCount int64
	declQ := applyClaimFilters(DB.Model(&models.GroupSchemeClaim{}), f)
	if err := declQ.Where("LOWER(status) = ?", "declined").Count(&declinedCount).Error; err != nil {
		return nil, err
	}
	var declineRate float64
	if totalClaims > 0 {
		declineRate = float64(declinedCount) / float64(totalClaims)
	}

	// Average processing days: from creation to first terminal status (approved/declined/paid)
	// Fetch claims ids and creation dates in scope
	var baseClaims []struct {
		ID           int
		CreationDate time.Time
	}
	// Use a fresh, ungrouped query with the same filters to avoid only_full_group_by issues
	baseQ := applyClaimFilters(DB.Model(&models.GroupSchemeClaim{}), f)
	if err := baseQ.Select("id, creation_date").Scan(&baseClaims).Error; err != nil {
		return nil, err
	}
	avgProcessingDays := 0.0
	if len(baseClaims) > 0 {
		ids := make([]int, 0, len(baseClaims))
		creation := make(map[int]time.Time, len(baseClaims))
		for _, c := range baseClaims {
			ids = append(ids, c.ID)
			creation[c.ID] = c.CreationDate
		}

		var audits []models.GroupSchemeClaimStatusAudit
		if err := DB.Where("claim_id IN ? AND LOWER(new_status) IN ?", ids, []string{"approved", "declined", "paid"}).
			Order("claim_id, changed_at").Find(&audits).Error; err != nil {
			return nil, err
		}
		// Pick first terminal audit per claim
		firstTerminal := make(map[int]time.Time)
		for _, a := range audits {
			if _, ok := firstTerminal[a.ClaimID]; !ok {
				firstTerminal[a.ClaimID] = a.ChangedAt
			}
		}
		// Compute mean of durations for claims that have terminal status
		var durations []float64
		sumDays := 0.0
		n := 0.0
		onTimeCount := 0.0
		slaDays := 30.0
		for id, end := range firstTerminal {
			start := creation[id]
			if end.After(start) {
				diff := end.Sub(start).Hours() / 24.0
				sumDays += diff
				durations = append(durations, diff)
				n += 1.0
				if diff <= slaDays {
					onTimeCount += 1.0
				}
			}
		}

		processingTimeMean := 0.0
		processingTimeMedian := 0.0
		var processingTimeDistribution []map[string]interface{}

		if n > 0 {
			avgProcessingDays = utils.FloatPrecision(sumDays/n, 2)
			processingTimeMean = avgProcessingDays
			median, _ := stats.Median(durations)
			processingTimeMedian = utils.FloatPrecision(median, 2)

			// Distribution buckets (5-day intervals)
			maxD := 0.0
			for _, d := range durations {
				if d > maxD {
					maxD = d
				}
			}
			// Cap at 100 or use maxD
			limit := 90
			if int(maxD) > limit {
				limit = (int(maxD)/5 + 1) * 5
			}

			distMap := make(map[int]int)
			for _, d := range durations {
				bucket := int(d) / 5
				distMap[bucket]++
			}

			for i := 0; i <= limit/5; i++ {
				label := fmt.Sprintf("%d-%d", i*5, (i*5)+4)
				processingTimeDistribution = append(processingTimeDistribution, map[string]interface{}{
					"label": label,
					"count": distMap[i],
				})
			}
		}

		// Compute processing efficiency related metrics
		// On-time rate (SLA compliance)
		onTimeRate := 0.0
		if n > 0 {
			onTimeRate = utils.FloatPrecision(onTimeCount/n, 4)
		}

		// Throughput per week: closed in period divided by number of weeks in period
		// Determine period window
		var fromRef, toRef time.Time
		if f.From != nil {
			fromRef = *f.From
		} else {
			// default last 12 months if not provided (consistent with controller)
			now := time.Now()
			fromRef = now.AddDate(-1, 0, 0)
		}
		if f.To != nil {
			toRef = *f.To
		} else {
			toRef = time.Now()
		}
		totalWeeks := toRef.Sub(fromRef).Hours() / (24.0 * 7.0)
		if totalWeeks < 1.0 {
			totalWeeks = 1.0
		}
		throughputPerWeek := utils.FloatPrecision(n/totalWeeks, 4)

		// Work-in-progress: open claims count (non-terminal statuses) at query snapshot within filters
		var openCount int64
		if err := applyClaimFilters(DB.Model(&models.GroupSchemeClaim{}), f).
			Where("LOWER(status) NOT IN ?", []string{"approved", "declined", "paid"}).
			Count(&openCount).Error; err != nil {
			return nil, err
		}

		// Composite efficiency score (v1): weighted combination
		// Components: onTimeRate (0..1), speed component (1 - avgDays/SLA), throughput normalized to a target
		speedComponent := 1.0
		if avgProcessingDays > 0 {
			ratio := avgProcessingDays / slaDays
			if ratio > 1.0 {
				ratio = 1.0
			}
			speedComponent = 1.0 - ratio
		}
		throughputTarget := 25.0 // adjustable target throughput per week
		throughputNorm := throughputPerWeek / throughputTarget
		if throughputNorm > 1.0 {
			throughputNorm = 1.0
		}
		if throughputNorm < 0.0 {
			throughputNorm = 0.0
		}
		efficiencyScore := 0.6*onTimeRate + 0.3*speedComponent + 0.1*throughputNorm
		efficiencyScore = utils.FloatPrecision(efficiencyScore, 4)

		// Attach processing_efficiency object into result later
		_processingEfficiency = map[string]interface{}{
			"sla_days":                     int(slaDays),
			"on_time_rate":                 onTimeRate,
			"throughput_per_week":          throughputPerWeek,
			"wip_open_claims":              openCount,
			"efficiency_score":             efficiencyScore,
			"closed_in_period":             int(n),
			"processing_time_mean":         processingTimeMean,
			"processing_time_median":       processingTimeMedian,
			"processing_time_distribution": processingTimeDistribution,
		}
	}

	// Build top claims table rows (configurable limit, default 10)
	topLimit := f.Limit
	if topLimit <= 0 {
		topLimit = 10
	}
	if topLimit > 1000 {
		topLimit = 1000
	}
	type topClaimRow struct {
		ID           int       `json:"id"`
		ClaimNumber  string    `json:"claim_number"`
		MemberName   string    `json:"member"`
		SchemeName   string    `json:"scheme"`
		BenefitType  string    `json:"benefit_type"`
		ClaimAmount  float64   `json:"amount"`
		Status       string    `json:"status"`
		CreationDate time.Time `json:"creation_date"`
		DateNotified string    `json:"date_notified"`
	}
	var topClaims []topClaimRow
	if err := applyClaimFilters(DB.Model(&models.GroupSchemeClaim{}), f).
		Select("id, claim_number, member_name, scheme_name, benefit_type, claim_amount, status, creation_date, date_notified").
		Order("claim_amount DESC").
		Limit(topLimit).
		Scan(&topClaims).Error; err != nil {
		return nil, err
	}

	// Top 5 performing assessors for the given period (by assessment count)
	type assessorRow struct {
		Name string `gorm:"column:name" json:"name"`
		Cnt  int64  `gorm:"column:cnt" json:"count"`
	}
	var topAssessors []assessorRow
	assessQ := DB.Model(&models.GroupSchemeClaimAssessment{})
	// Filter by assessment date range if provided (stored as YYYY-MM-DD string)
	if f.From != nil {
		assessQ = assessQ.Where("assessment_date >= ?", f.From.Format("2006-01-02"))
	}
	if f.To != nil {
		assessQ = assessQ.Where("assessment_date <= ?", f.To.Format("2006-01-02"))
	}
	// Restrict to scheme/benefit if applied in filters by joining claims
	if f.SchemeID != nil || f.Benefit != "" {
		assessQ = assessQ.Joins("JOIN group_scheme_claims c ON c.id = group_scheme_claim_assessments.claim_id")
		if f.SchemeID != nil {
			assessQ = assessQ.Where("c.scheme_id = ?", *f.SchemeID)
		}
		if f.Benefit != "" {
			assessQ = assessQ.Where("c.benefit_type = ?", f.Benefit)
		}
	}
	if err := assessQ.Select("assessor_name AS name, COUNT(*) AS cnt").
		Where("TRIM(assessor_name) <> ''").
		Group("assessor_name").
		Order("cnt DESC").
		Limit(5).
		Scan(&topAssessors).Error; err != nil {
		return nil, err
	}

	// Top decline reasons (by primary_reason) within the given period and filters
	type declineReasonRow struct {
		Reason string `gorm:"column:reason" json:"reason"`
		Cnt    int64  `gorm:"column:cnt" json:"count"`
	}
	var topDeclineReasons []declineReasonRow
	drQ := DB.Model(&models.GroupSchemeClaimDecline{})
	// Filter by declined_at timestamp when provided
	if f.From != nil {
		drQ = drQ.Where("declined_at >= ?", *f.From)
	}
	if f.To != nil {
		drQ = drQ.Where("declined_at <= ?", *f.To)
	}
	// Join to claims when scheme/benefit filters are used
	if f.SchemeID != nil || f.Benefit != "" {
		drQ = drQ.Joins("JOIN group_scheme_claims c ON c.id = group_scheme_claim_declines.claim_id")
		if f.SchemeID != nil {
			drQ = drQ.Where("c.scheme_id = ?", *f.SchemeID)
		}
		if f.Benefit != "" {
			drQ = drQ.Where("c.benefit_type = ?", f.Benefit)
		}
	}
	if err := drQ.Select("primary_reason AS reason, COUNT(*) AS cnt").
		Where("TRIM(primary_reason) <> ''").
		Group("primary_reason").
		Order("cnt DESC").
		Limit(5).
		Scan(&topDeclineReasons).Error; err != nil {
		return nil, err
	}

	// include processing efficiency card info (if computed) via temporary var
	var processingEfficiency any
	if _processingEfficiency != nil {
		processingEfficiency = _processingEfficiency
	}

	// Claims trend data (Monthly)
	type trendRow struct {
		Month  string  `gorm:"column:month"`
		Count  int64   `gorm:"column:count"`
		Amount float64 `gorm:"column:amount"`
	}
	var trendRows []trendRow
	// Group by month using DATE_FORMAT for MySQL, filtering for approved claims only
	trendQ := applyClaimFilters(DB.Model(&models.GroupSchemeClaim{}), f)
	if err := trendQ.Where("LOWER(status) IN ?", []string{"approved", "paid"}).
		Select("DATE_FORMAT(creation_date, '%Y-%m') as month, COUNT(*) as count, SUM(claim_amount) as amount").
		Group("month").
		Order("month").
		Scan(&trendRows).Error; err != nil {
		return nil, err
	}

	// Ensure all months in the requested period are represented (even with 0)
	var finalTrend []map[string]interface{}
	var fromTrend, toTrend time.Time
	if f.From != nil {
		fromTrend = *f.From
	} else {
		fromTrend = time.Now().AddDate(-1, 0, 0)
	}
	if f.To != nil {
		toTrend = *f.To
	} else {
		toTrend = time.Now()
	}

	// Start from the 'from' month and go up to the 'to' month
	curr := time.Date(fromTrend.Year(), fromTrend.Month(), 1, 0, 0, 0, 0, time.UTC)
	endTrend := time.Date(toTrend.Year(), toTrend.Month(), 1, 0, 0, 0, 0, time.UTC)

	trendMap := make(map[string]trendRow)
	for _, tr := range trendRows {
		trendMap[tr.Month] = tr
	}

	for !curr.After(endTrend) {
		mStr := curr.Format("2006-01")
		if tr, ok := trendMap[mStr]; ok {
			finalTrend = append(finalTrend, map[string]interface{}{
				"month":  mStr,
				"count":  tr.Count,
				"amount": utils.FloatPrecision(tr.Amount, AccountingPrecision),
			})
		} else {
			finalTrend = append(finalTrend, map[string]interface{}{
				"month":  mStr,
				"count":  0,
				"amount": 0.0,
			})
		}
		curr = curr.AddDate(0, 1, 0)
	}

	// Claimant Age Distribution
	// Only for approved/paid claims.
	// We use the main member's DateOfBirth from GPricingMemberDataInForce, even for beneficiaries.
	var memberDOBs []struct {
		DateOfBirth time.Time
	}
	// Use fresh query and join with GPricingMemberDataInForce to get member's DOB
	ageQ := applyClaimFilters(DB.Model(&models.GroupSchemeClaim{}), f)
	if err := ageQ.Where("LOWER(group_scheme_claims.status) IN ?", []string{"approved", "paid"}).
		Joins("JOIN g_pricing_member_data_in_forces m ON m.member_id_number = group_scheme_claims.member_id_number").
		Select("DISTINCT m.date_of_birth, group_scheme_claims.id").
		Scan(&memberDOBs).Error; err != nil {
		return nil, err
	}

	ageFreq := make(map[int]int)
	for _, m := range memberDOBs {
		dob := m.DateOfBirth
		if dob.IsZero() {
			continue
		}

		age := time.Now().Year() - dob.Year()
		if time.Now().Before(time.Date(time.Now().Year(), dob.Month(), dob.Day(), 0, 0, 0, 0, time.Now().Location())) {
			age--
		}

		if age < 0 || age > 120 {
			continue
		}
		ageFreq[age]++
	}

	// Bucket the ages into 5-year intervals (e.g., 15-19, 20-24, ...)
	// The chart shows 20, 25, 30... as labels
	type ageBucket struct {
		Label string `json:"label"`
		Count int    `json:"count"`
	}
	var ageDistribution []ageBucket
	minAge, maxAge := 15, 85
	for start := minAge; start <= maxAge; start += 5 {
		bucketCount := 0
		for age := start; age < start+5; age++ {
			bucketCount += ageFreq[age]
		}
		// If it's the last bucket, include everything above it?
		// Or just stick to the range.
		label := fmt.Sprintf("%d-%d", start, start+4)
		ageDistribution = append(ageDistribution, ageBucket{Label: label, Count: bucketCount})
	}

	result := map[string]interface{}{
		"total_claims":              totalClaims,
		"total_paid_amount":         utils.FloatPrecision(totalPaid, AccountingPrecision),
		"avg_processing_days":       avgProcessingDays,
		"approval_rate":             utils.FloatPrecision(approvalRate, 4),
		"decline_rate":              utils.FloatPrecision(declineRate, 4),
		"by_status":                 byStatus,
		"by_benefit":                byBenefit,
		"top_claims":                topClaims,
		"top_assessors":             topAssessors,
		"topDeclineReasons":         topDeclineReasons,
		"processing_efficiency":     processingEfficiency,
		"claims_trend":              finalTrend,
		"claimant_age_distribution": ageDistribution,
	}
	return result, nil
}

// UpdateGroupSchemeClaim updates a claim by ID. It preserves immutable fields like ID, CreationDate, and CreatedBy.
func UpdateGroupSchemeClaim(claimID int, payload models.GroupSchemeClaim, user models.AppUser) (models.GroupSchemeClaim, error) {
	// Ensure the claim exists
	var existing models.GroupSchemeClaim
	if err := DB.First(&existing, claimID).Error; err != nil {
		return existing, err
	}

	// Force IDs to match and protect immutable fields
	payload.ID = existing.ID

	// Detect status change
	statusChanged := strings.TrimSpace(payload.Status) != "" && payload.Status != existing.Status

	// Perform full-field update (including zero values) while omitting immutable columns
	if err := DB.Model(&existing).Select("*").Omit("id", "creation_date", "created_by").Updates(payload).Error; err != nil {
		return existing, err
	}

	// Record status audit if changed
	if statusChanged {
		if err := DB.Create(&models.GroupSchemeClaimStatusAudit{
			ClaimID:       existing.ID,
			OldStatus:     existing.Status,
			NewStatus:     payload.Status,
			StatusMessage: "Status updated",
			ChangedBy:     user.UserName,
			ChangedAt:     time.Now(),
		}).Error; err != nil {
			return existing, err
		}

		// Log structured claim status activity
		claimDetails, _ := json.Marshal(map[string]interface{}{
			"claimNumber": existing.ClaimNumber,
			"claimType":   existing.BenefitAlias,
			"oldStatus":   existing.Status,
			"newStatus":   payload.Status,
		})
		var m models.GPricingMemberDataInForce
		_ = DB.Where("member_id_number = ?", existing.MemberIDNumber).First(&m)

		_ = DB.Create(&models.MemberActivity{
			MemberID:       m.ID,
			MemberIDNumber: existing.MemberIDNumber,
			Type:           "claim",
			Title:          "Claim Status Updated",
			Description:    fmt.Sprintf("%s claim status changed from %s to %s", existing.BenefitAlias, existing.Status, payload.Status),
			Details:        claimDetails,
			PerformedBy:    user.UserName,
		})

		// Auto-generate reinsurer recovery when claim is approved or paid
		newStatus := strings.ToLower(payload.Status)
		if newStatus == "approved" || newStatus == "paid" {
			var refreshed models.GroupSchemeClaim
			if err := DB.First(&refreshed, claimID).Error; err != nil {
				appLog.WithFields(map[string]interface{}{
					"claim_id": claimID,
					"error":    err.Error(),
				}).Error("Failed to reload claim for recovery generation")
			} else if _, err := GenerateClaimRecovery(DB, refreshed, user); err != nil {
				appLog.WithFields(map[string]interface{}{
					"claim_id":     claimID,
					"claim_number": refreshed.ClaimNumber,
					"error":        err.Error(),
				}).Error("GenerateClaimRecovery failed on claim status transition")
			}
		}
	}

	// Return the refreshed record with relations
	var updated models.GroupSchemeClaim
	if err := DB.Preload("Attachments").Preload("Assessments").Preload("Communications").First(&updated, claimID).Error; err != nil {
		return updated, err
	}
	return updated, nil
}

// UpdateGroupSchemeClaimWithFiles updates a claim and optionally appends new attachments from multipart upload.
// Immutable fields (id, creation_date, created_by) are preserved. New files are stored under tmp/uploads/group_claims/claim_<id>.
func UpdateGroupSchemeClaimWithFiles(claimID int, payload models.GroupSchemeClaim, files map[string][]*multipart.FileHeader, user models.AppUser) (models.GroupSchemeClaim, error) {
	// Start a transaction to ensure atomicity for DB updates and attachment rows
	tx := DB.Begin()
	if tx.Error != nil {
		return models.GroupSchemeClaim{}, tx.Error
	}

	var existing models.GroupSchemeClaim
	if err := tx.First(&existing, claimID).Error; err != nil {
		tx.Rollback()
		return existing, err
	}

	// Apply updates while protecting immutable fields
	payload.ID = existing.ID
	// Track status change before updating
	statusChanged := strings.TrimSpace(payload.Status) != "" && payload.Status != existing.Status

	if err := tx.Model(&existing).Select("*").Omit("id", "creation_date", "created_by").Updates(payload).Error; err != nil {
		tx.Rollback()
		return existing, err
	}

	if statusChanged {
		if err := tx.Create(&models.GroupSchemeClaimStatusAudit{
			ClaimID:       existing.ID,
			OldStatus:     existing.Status,
			NewStatus:     payload.Status,
			StatusMessage: "Status updated",
			ChangedBy:     user.UserName,
			ChangedAt:     time.Now(),
		}).Error; err != nil {
			tx.Rollback()
			return existing, err
		}

		// Log structured claim status activity
		claimDetails, _ := json.Marshal(map[string]interface{}{
			"claimNumber": existing.ClaimNumber,
			"claimType":   existing.BenefitAlias,
			"oldStatus":   existing.Status,
			"newStatus":   payload.Status,
		})
		var m models.GPricingMemberDataInForce
		_ = tx.Where("member_id_number = ?", existing.MemberIDNumber).First(&m)

		_ = tx.Create(&models.MemberActivity{
			MemberID:       m.ID,
			MemberIDNumber: existing.MemberIDNumber,
			Type:           "claim",
			Title:          "Claim Status Updated",
			Description:    fmt.Sprintf("%s claim status changed from %s to %s", existing.BenefitAlias, existing.Status, payload.Status),
			Details:        claimDetails,
			PerformedBy:    user.UserName,
		})

		// Auto-generate reinsurer recovery when claim is approved or paid. Use the
		// active transaction so the recovery row commits (or rolls back) with the
		// claim update — previously this wrote via the package-level DB and could
		// leave an orphan recovery if the outer tx rolled back.
		newStatus := strings.ToLower(payload.Status)
		if newStatus == "approved" || newStatus == "paid" {
			var refreshed models.GroupSchemeClaim
			if err := tx.First(&refreshed, claimID).Error; err != nil {
				appLog.WithFields(map[string]interface{}{
					"claim_id": claimID,
					"error":    err.Error(),
				}).Error("Failed to reload claim for recovery generation")
			} else if _, err := GenerateClaimRecovery(tx, refreshed, user); err != nil {
				appLog.WithFields(map[string]interface{}{
					"claim_id":     claimID,
					"claim_number": refreshed.ClaimNumber,
					"error":        err.Error(),
				}).Error("GenerateClaimRecovery failed on claim status transition")
			}
		}
	}

	// Handle file uploads (if any)
	if len(payload.SupportingDocuments) > 0 {
		baseDir := filepath.Join("tmp", "uploads", "group_claims", fmt.Sprintf("claim_%d", existing.ID))
		if err := os.MkdirAll(baseDir, 0o755); err != nil {
			tx.Rollback()
			return existing, err
		}

		for i, doc := range payload.SupportingDocuments {
			key := fmt.Sprintf("file_%d", i)
			fhList, ok := files[key]
			if !ok || len(fhList) == 0 {
				continue
			}
			fh := fhList[0]

			name := filepath.Base(fh.Filename)
			destPath := filepath.Join(baseDir, name)

			if err := func() error {
				src, err := fh.Open()
				if err != nil {
					return err
				}
				defer src.Close()

				dst, err := os.Create(destPath)
				if err != nil {
					return err
				}
				defer dst.Close()

				if _, err = io.Copy(dst, src); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				tx.Rollback()
				return existing, err
			}

			att := models.GroupSchemeClaimAttachment{
				ClaimID:      existing.ID,
				DocumentType: doc.DocumentType,
				DocumentName: doc.DocumentName,
				FileName:     name,
				ContentType:  fh.Header.Get("Content-Type"),
				SizeBytes:    fh.Size,
				StoragePath:  destPath,
				UploadedAt:   time.Now(),
				UploadedBy:   user.UserName,
			}
			if err := tx.Create(&att).Error; err != nil {
				tx.Rollback()
				return existing, err
			}
		}
	} else if len(files) > 0 {
		// Fallback for when SupportingDocuments is not provided but files are
		baseDir := filepath.Join("tmp", "uploads", "group_claims", fmt.Sprintf("claim_%d", existing.ID))
		if err := os.MkdirAll(baseDir, 0o755); err != nil {
			tx.Rollback()
			return existing, err
		}

		for _, fhList := range files {
			for _, fh := range fhList {
				name := filepath.Base(fh.Filename)
				destPath := filepath.Join(baseDir, name)

				if err := func() error {
					src, err := fh.Open()
					if err != nil {
						return err
					}
					defer src.Close()

					dst, err := os.Create(destPath)
					if err != nil {
						return err
					}
					defer dst.Close()

					if _, err = io.Copy(dst, src); err != nil {
						return err
					}
					return nil
				}(); err != nil {
					tx.Rollback()
					return existing, err
				}

				att := models.GroupSchemeClaimAttachment{
					ClaimID:     existing.ID,
					FileName:    name,
					ContentType: fh.Header.Get("Content-Type"),
					SizeBytes:   fh.Size,
					StoragePath: destPath,
					UploadedAt:  time.Now(),
					UploadedBy:  user.UserName,
				}
				if err := tx.Create(&att).Error; err != nil {
					tx.Rollback()
					return existing, err
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return existing, err
	}

	// Reload updated record with relations
	var updated models.GroupSchemeClaim
	if err := DB.Preload("Attachments").Preload("Assessments").Preload("Communications").First(&updated, existing.ID).Error; err != nil {
		return updated, err
	}
	return updated, nil
}

// CreateClaimCommunication creates a new communication log linked to a claim
func CreateClaimCommunication(comm models.GroupSchemeClaimCommunication, user models.AppUser) (models.GroupSchemeClaimCommunication, error) {
	// Ensure claim exists
	var claim models.GroupSchemeClaim
	if err := DB.First(&claim, comm.ClaimID).Error; err != nil {
		return comm, err
	}
	// Set audit fields
	comm.CreatedBy = user.UserName
	// If Timestamp not provided, leave nil; CreatedAt auto-set by GORM
	if err := DB.Create(&comm).Error; err != nil {
		return comm, err
	}
	return comm, nil
}

// GetClaimCommunicationsByClaim returns communications linked to the given claim ID
func GetClaimCommunicationsByClaim(claimID int) ([]models.GroupSchemeClaimCommunication, error) {
	var count int64
	if err := DB.Model(&models.GroupSchemeClaim{}).Where("id = ?", claimID).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	var comms []models.GroupSchemeClaimCommunication
	if err := DB.Where("claim_id = ?", claimID).Order("id ASC").Find(&comms).Error; err != nil {
		return nil, err
	}
	if comms == nil {
		comms = []models.GroupSchemeClaimCommunication{}
	}
	return comms, nil
}

func GetSchemeMemberRating(schemeId, quoteId, memberId string) (models.MemberRatingResult, error) {
	var member models.GPricingMemberDataInForce
	var memberRating models.MemberRatingResult

	err := DB.Where("scheme_name = ? AND quote_id = ? AND id = ?", schemeId, quoteId, memberId).First(&member).Error
	if err != nil {
		return memberRating, err
	}

	err = DB.Where("scheme_name = ? AND quote_id = ? AND member_name = ? AND date_of_birth = ?", schemeId, quoteId, member.MemberName, member.DateOfBirth).First(&memberRating).Error
	if err != nil {
		return memberRating, err
	}
	return memberRating, nil
}

// CreateClaimAssessment creates a new assessment linked to a claim
func CreateClaimAssessment(assessment models.GroupSchemeClaimAssessment, user models.AppUser) (models.GroupSchemeClaimAssessment, error) {
	// Ensure claim exists
	var claim models.GroupSchemeClaim
	if err := DB.First(&claim, assessment.ClaimID).Error; err != nil {
		return assessment, err
	}
	// Set audit fields
	assessment.CreatedBy = user.UserName
	now := time.Now()
	if assessment.AssessmentTimestamp == nil {
		assessment.AssessmentTimestamp = &now
	}
	if err := DB.Create(&assessment).Error; err != nil {
		return assessment, err
	}

	// Write audit trail for assessment creation
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "group_scheme_claim_assessments",
		EntityID:  strconv.Itoa(assessment.ID),
		Action:    "CREATE",
		ChangedBy: user.UserName,
	}, struct{}{}, assessment)

	// Log structured claim assessment activity
	assessmentDetails, _ := json.Marshal(map[string]interface{}{
		"claimNumber":  claim.ClaimNumber,
		"claimType":    claim.BenefitAlias,
		"assessor":     assessment.AssessorName,
		"outcome":      assessment.AssessmentOutcome,
		"amount":       assessment.RecommendedAmount,
		"medicalNotes": assessment.MedicalNotes,
	})
	var m models.GPricingMemberDataInForce
	_ = DB.Where("member_id_number = ?", claim.MemberIDNumber).First(&m)

	_ = DB.Create(&models.MemberActivity{
		MemberID:       m.ID,
		MemberIDNumber: claim.MemberIDNumber,
		Type:           "claim",
		Title:          "Claim Assessment Created",
		Description:    fmt.Sprintf("%s claim assessment performed with outcome: %s", claim.BenefitAlias, assessment.AssessmentOutcome),
		Details:        assessmentDetails,
		PerformedBy:    user.UserName,
	})

	return assessment, nil
}

// UpdateClaimAssessment updates an existing assessment by ID
func UpdateClaimAssessment(assessmentID int, payload models.GroupSchemeClaimAssessment, user models.AppUser) (models.GroupSchemeClaimAssessment, error) {
	var existing models.GroupSchemeClaimAssessment
	if err := DB.First(&existing, assessmentID).Error; err != nil {
		return existing, err
	}

	// Snapshot before state for audit
	before := existing

	// Only update fields that can change
	existing.AssessorName = payload.AssessorName
	existing.AssessmentDate = payload.AssessmentDate
	existing.AssessmentOutcome = payload.AssessmentOutcome
	existing.RecommendedAmount = payload.RecommendedAmount
	existing.MedicalOfficer = payload.MedicalOfficer
	existing.MedicalAssessmentDate = payload.MedicalAssessmentDate
	existing.DisabilityPercentage = payload.DisabilityPercentage
	existing.MedicalCondition = payload.MedicalCondition
	existing.MedicalNotes = payload.MedicalNotes
	existing.DocumentsVerified = payload.DocumentsVerified
	existing.FraudRiskLevel = payload.FraudRiskLevel
	existing.RequiresInvestigation = payload.RequiresInvestigation
	existing.RiskNotes = payload.RiskNotes
	existing.AssessmentNotes = payload.AssessmentNotes
	existing.NextActions = payload.NextActions
	existing.Checklist = payload.Checklist
	existing.AssessmentTimestamp = payload.AssessmentTimestamp
	// keep ClaimID unchanged; CreatedBy remains original

	if err := DB.Save(&existing).Error; err != nil {
		return existing, err
	}

	// Write audit trail for assessment update
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "group_scheme_claim_assessments",
		EntityID:  strconv.Itoa(assessmentID),
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, existing)

	// Log structured claim assessment activity
	var claim models.GroupSchemeClaim
	if err := DB.First(&claim, existing.ClaimID).Error; err == nil {
		assessmentDetails, _ := json.Marshal(map[string]interface{}{
			"claimNumber":  claim.ClaimNumber,
			"claimType":    claim.BenefitAlias,
			"assessor":     existing.AssessorName,
			"outcome":      existing.AssessmentOutcome,
			"amount":       existing.RecommendedAmount,
			"medicalNotes": existing.MedicalNotes,
		})
		var m models.GPricingMemberDataInForce
		_ = DB.Where("member_id_number = ?", claim.MemberIDNumber).First(&m)

		_ = DB.Create(&models.MemberActivity{
			MemberID:       m.ID,
			MemberIDNumber: claim.MemberIDNumber,
			Type:           "claim",
			Title:          "Claim Assessment Updated",
			Description:    fmt.Sprintf("%s claim assessment updated with outcome: %s", claim.BenefitAlias, existing.AssessmentOutcome),
			Details:        assessmentDetails,
			PerformedBy:    user.UserName,
		})
	}

	return existing, nil
}

// GetClaimAssessmentsByClaim returns assessments linked to the given claim ID
func GetClaimAssessmentsByClaim(claimID int) ([]models.GroupSchemeClaimAssessment, error) {
	// Optionally ensure claim exists for clearer error messages
	var count int64
	if err := DB.Model(&models.GroupSchemeClaim{}).Where("id = ?", claimID).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	var assessments []models.GroupSchemeClaimAssessment
	if err := DB.Where("claim_id = ?", claimID).Order("id ASC").Find(&assessments).Error; err != nil {
		return nil, err
	}
	if assessments == nil {
		assessments = []models.GroupSchemeClaimAssessment{}
	}
	return assessments, nil
}

// CreateClaimDecline creates a decline record linked to a claim and optionally updates claim status
func CreateClaimDecline(decline models.GroupSchemeClaimDecline, user models.AppUser) (models.GroupSchemeClaimDecline, error) {
	// Ensure claim exists
	var claim models.GroupSchemeClaim
	if err := DB.First(&claim, decline.ClaimID).Error; err != nil {
		return decline, err
	}

	// Populate audit fields
	if strings.TrimSpace(decline.DeclinedBy) == "" {
		decline.DeclinedBy = user.UserName
	}
	if decline.DeclinedAt == nil {
		now := time.Now()
		decline.DeclinedAt = &now
	}

	if err := DB.Create(&decline).Error; err != nil {
		return decline, err
	}

	// Write audit trail for claim decline
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "group_scheme_claim_declines",
		EntityID:  strconv.Itoa(decline.ID),
		Action:    "CREATE",
		ChangedBy: user.UserName,
	}, struct{}{}, decline)

	// Update claim status to Declined if not already and audit it
	if !strings.EqualFold(claim.Status, DECLINED) {
		oldStatus := claim.Status
		if err := DB.Model(&claim).Update("status", DECLINED).Error; err != nil {
			return decline, err
		}
		if err := DB.Create(&models.GroupSchemeClaimStatusAudit{
			ClaimID:       claim.ID,
			OldStatus:     oldStatus,
			NewStatus:     DECLINED,
			StatusMessage: "Claim declined",
			ChangedBy:     decline.DeclinedBy,
			ChangedAt:     time.Now(),
		}).Error; err != nil {
			return decline, err
		}

		// Log structured claim decline activity
		claimDetails, _ := json.Marshal(map[string]interface{}{
			"claimNumber": claim.ClaimNumber,
			"claimType":   claim.BenefitAlias,
			"status":      DECLINED,
			"reason":      decline.PrimaryReason,
		})
		var m models.GPricingMemberDataInForce
		_ = DB.Where("member_id_number = ?", claim.MemberIDNumber).First(&m)

		_ = DB.Create(&models.MemberActivity{
			MemberID:       m.ID,
			MemberIDNumber: claim.MemberIDNumber,
			Type:           "claim",
			Title:          "Claim Declined",
			Description:    fmt.Sprintf("%s claim declined", claim.BenefitAlias),
			Details:        claimDetails,
			PerformedBy:    decline.DeclinedBy,
		})
	}

	return decline, nil
}

// GetClaimDeclinesByClaim returns declines for a given claim ID
func GetClaimDeclinesByClaim(claimID int) ([]models.GroupSchemeClaimDecline, error) {
	var count int64
	if err := DB.Model(&models.GroupSchemeClaim{}).Where("id = ?", claimID).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	var declines []models.GroupSchemeClaimDecline
	if err := DB.Where("claim_id = ?", claimID).Order("id ASC").Find(&declines).Error; err != nil {
		return nil, err
	}
	if declines == nil {
		declines = []models.GroupSchemeClaimDecline{}
	}
	return declines, nil
}

func GetBenefitMaps() ([]models.GroupBenefitMapper, error) {
	var benefitMaps []models.GroupBenefitMapper
	err := DB.Find(&benefitMaps).Error
	if err != nil {
		return nil, err
	}

	// if benefitMaps is empty, generate the base list of benefits
	if len(benefitMaps) == 0 {
		benefitMaps = getBaseBenefitMaps()
	}

	return benefitMaps, nil
}

func GetBenefitMapsByScheme(schemeId string) ([]models.GroupBenefitMapper, error) {
	var scheme models.GroupScheme
	if err := DB.Where("id = ?", schemeId).First(&scheme).Error; err != nil {
		return nil, err
	}

	var categories []models.SchemeCategory
	if err := DB.Where("quote_id = ?", scheme.QuoteId).Find(&categories).Error; err != nil {
		return nil, err
	}

	enabledBenefits := make(map[string]bool)
	benefitAliases := make(map[string]string)

	for _, cat := range categories {
		// Only consider categories that are in ActiveSchemeCategories
		//isActive := false
		//for _, activeCat := range scheme.ActiveSchemeCategories {
		//	if activeCat == cat.SchemeCategory {
		//		isActive = true
		//		break
		//	}
		//}
		//if !isActive {
		//	continue
		//}

		if cat.GlaBenefit {
			enabledBenefits["GLA"] = true
			if cat.GlaAlias != "" {
				benefitAliases["GLA"] = cat.GlaAlias
			}
		}
		if cat.AdditionalAccidentalGlaBenefit {
			enabledBenefits["AAGLA"] = true
		}
		if cat.AdditionalGlaCoverBenefit {
			enabledBenefits["AGLA"] = true
		}
		if cat.SglaBenefit {
			enabledBenefits["SGLA"] = true
			if cat.SglaAlias != "" {
				benefitAliases["SGLA"] = cat.SglaAlias
			}
		}
		if cat.PtdBenefit {
			enabledBenefits["PTD"] = true
			if cat.PtdAlias != "" {
				benefitAliases["PTD"] = cat.PtdAlias
			}
		}
		if cat.TtdBenefit {
			enabledBenefits["TTD"] = true
			if cat.TtdAlias != "" {
				benefitAliases["TTD"] = cat.TtdAlias
			}
		}
		if cat.PhiBenefit {
			enabledBenefits["PHI"] = true
			if cat.PhiAlias != "" {
				benefitAliases["PHI"] = cat.PhiAlias
			}
		}
		if cat.CiBenefit {
			enabledBenefits["CI"] = true
			if cat.CiAlias != "" {
				benefitAliases["CI"] = cat.CiAlias
			}
		}
		if cat.FamilyFuneralBenefit {
			enabledBenefits["GFF"] = true
			if cat.FamilyFuneralAlias != "" {
				benefitAliases["GFF"] = cat.FamilyFuneralAlias
			}
		}
		// Educator benefits are split per attachment: GLA_EDU is enabled
		// whenever the GLA educator is configured and PTD_EDU whenever the
		// PTD educator is configured, so users can customise display names
		// independently.
		if cat.GlaEducatorBenefit != "" {
			enabledBenefits["GLA_EDU"] = true
		}
		if cat.PtdEducatorBenefit != "" {
			enabledBenefits["PTD_EDU"] = true
		}
	}

	allMaps, err := GetBenefitMaps()
	if err != nil {
		return nil, err
	}

	var filteredMaps []models.GroupBenefitMapper
	for _, m := range allMaps {
		if enabledBenefits[m.BenefitCode] {
			if alias, ok := benefitAliases[m.BenefitCode]; ok {
				m.BenefitAlias = alias
			}
			filteredMaps = append(filteredMaps, m)
		}
	}

	return filteredMaps, nil
}

func GetBenefitMapsBySchemeCategory(schemeId string, categoryId string) ([]models.GroupBenefitMapper, error) {
	var category models.SchemeCategory
	if err := DB.Where("id = ?", categoryId).First(&category).Error; err != nil {
		return nil, err
	}

	enabledBenefits := make(map[string]bool)
	benefitAliases := make(map[string]string)

	if category.GlaBenefit {
		enabledBenefits["GLA"] = true
		if category.GlaAlias != "" {
			benefitAliases["GLA"] = category.GlaAlias
		}
	}
	if category.AdditionalAccidentalGlaBenefit {
		enabledBenefits["AAGLA"] = true
	}
	if category.AdditionalGlaCoverBenefit {
		enabledBenefits["AGLA"] = true
	}
	if category.SglaBenefit {
		enabledBenefits["SGLA"] = true
		if category.SglaAlias != "" {
			benefitAliases["SGLA"] = category.SglaAlias
		}
	}
	if category.PtdBenefit {
		enabledBenefits["PTD"] = true
		if category.PtdAlias != "" {
			benefitAliases["PTD"] = category.PtdAlias
		}
	}
	if category.TtdBenefit {
		enabledBenefits["TTD"] = true
		if category.TtdAlias != "" {
			benefitAliases["TTD"] = category.TtdAlias
		}
	}
	if category.PhiBenefit {
		enabledBenefits["PHI"] = true
		if category.PhiAlias != "" {
			benefitAliases["PHI"] = category.PhiAlias
		}
	}
	if category.CiBenefit {
		enabledBenefits["CI"] = true
		if category.CiAlias != "" {
			benefitAliases["CI"] = category.CiAlias
		}
	}
	if category.FamilyFuneralBenefit {
		enabledBenefits["GFF"] = true
		if category.FamilyFuneralAlias != "" {
			benefitAliases["GFF"] = category.FamilyFuneralAlias
		}
	}
	if category.GlaEducatorBenefit != "" {
		enabledBenefits["GLA_EDU"] = true
	}
	if category.PtdEducatorBenefit != "" {
		enabledBenefits["PTD_EDU"] = true
	}

	allMaps, err := GetBenefitMaps()
	if err != nil {
		return nil, err
	}

	var filteredMaps []models.GroupBenefitMapper
	for _, m := range allMaps {
		if enabledBenefits[m.BenefitCode] {
			if alias, ok := benefitAliases[m.BenefitCode]; ok {
				m.BenefitAlias = alias
			}
			filteredMaps = append(filteredMaps, m)
		}
	}

	return filteredMaps, nil
}
func SaveBenefitMaps(benefitMaps []models.GroupBenefitMapper) error {
	// Delete all existing records
	err := DB.Exec("DELETE FROM group_benefit_mappers").Error
	if err != nil {
		return err
	}

	// Insert the new benefit maps
	for _, benefitMap := range benefitMaps {
		if benefitMap.BenefitAlias != "" {
			benefitMap.IsMapped = true
		}
		err = DB.Create(&benefitMap).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func GetBenefitDefinitions() ([]string, error) {
	var benefitDefinitions []string
	err := DB.Table("ci_rates").Distinct().Pluck("benefit_definition", &benefitDefinitions).Error
	if err != nil {
		return nil, err
	}
	return benefitDefinitions, nil
}

func GetGPPermissions() ([]models.GPPermission, error) {
	var permissions []models.GPPermission
	err := DB.Find(&permissions).Error
	if err != nil {
		return permissions, err
	}
	return permissions, nil
}

func GetGPUserRoles() ([]models.GPUserRole, error) {
	var roles []models.GPUserRole
	err := DB.Preload("Permissions").Find(&roles).Error
	if err != nil {
		return roles, err
	}
	return roles, nil
}

func CreateGPUserRole(role models.GPUserRole) (models.GPUserRole, error) {

	permissions := role.Permissions
	role.Permissions = nil // avoid GORM trying to insert join table before Role exists

	if role.ID == 0 {
		if err := DB.Create(&role).Error; err != nil {
			// if there is a duplicate key error, return the error formatted
			if strings.Contains(err.Error(), "Duplicate entry") {
				return role, fmt.Errorf("role with name '%s' already exists", role.RoleName)
			}
			return role, err
		}

	} else {
		if err := DB.Save(&role).Error; err != nil {
			// if there is a duplicate key error, return the error formatted
			if strings.Contains(err.Error(), "Duplicate entry") {
				return role, fmt.Errorf("role with name '%s' already exists", role.RoleName)
			}
			return role, err
		}
	}

	// remove any previous associations
	if err := DB.Model(&role).Association("Permissions").Clear(); err != nil {
		return role, err
	}

	if err := DB.Model(&role).Association("Permissions").Append(permissions); err != nil {
		return role, err
	}

	return role, nil
}

func DeleteGPUserRole(roleId string) error {
	var role models.GPUserRole
	err := DB.Where("id = ?", roleId).First(&role).Error
	if err != nil {
		return err
	}

	// check if role is in use
	var users []models.OrgUser
	err = DB.Where("gp_role_id = ?", roleId).Find(&users).Error

	if err != nil {
		fmt.Println(err)
	}
	if len(users) > 0 {
		return errors.New("role is in use")
	}

	// delete associated permissions
	err = DB.Model(&role).Association("Permissions").Clear()
	if err != nil {
		fmt.Println(err)
	}

	err = DB.Delete(&role).Error
	if err != nil {
		return err
	}
	return nil
}

func GetRolePermissions(roleId string) ([]models.GPPermission, error) {
	var permissions []models.GPPermission

	var role models.GPUserRole
	err := DB.Where("id = ?", roleId).Preload("Permissions").First(&role).Error
	if err != nil {
		return permissions, err
	}

	return filterByBaseline(role.Permissions), nil
}

// filterByBaseline drops any "special" permission whose parent baseline is
// not also present in the slug set. Baselines, system:admin, and any slug with
// no parent_slug pass through. Called wherever a user's effective permission
// set is materialised (role lookup, license lookup, email lookup) so that
// the rule "specials require their baseline" is enforced in one place — the
// frontend permission map and the RequirePermission middleware then both see
// a self-consistent set without needing their own checks.
func filterByBaseline(perms []models.GPPermission) []models.GPPermission {
	have := make(map[string]bool, len(perms))
	for _, p := range perms {
		have[p.Slug] = true
	}
	out := make([]models.GPPermission, 0, len(perms))
	for _, p := range perms {
		if p.Tier == "special" && p.ParentSlug != "" && !have[p.ParentSlug] {
			continue
		}
		out = append(out, p)
	}
	return out
}

func AssignRoleToUser(user models.OrgUser) error {
	fmt.Println(user)
	err := DB.Save(&user).Error
	if err != nil {
		appLog.Error(err.Error())
		return err
	}
	return nil
}

func RemoveRoleFromUser(user models.OrgUser) error {
	user.GPRoleId = 0
	user.GPRole = "None"
	err := DB.Where("id = ?", user.ID).Updates(&user).Error
	if err != nil {
		appLog.Error(err.Error())
		return err
	}
	return nil
}

func GetRoleForUserLicense(licenseId string) (models.GPUserRole, error) {
	var orgUser models.OrgUser
	err := DB.Where("license_id = ?", licenseId).First(&orgUser).Error
	if err != nil {
		appLog.Error("Error getting user role for license: ", err.Error())
		return models.GPUserRole{}, err
	}
	var role models.GPUserRole
	err = DB.Where("id = ?", orgUser.GPRoleId).Preload("Permissions").First(&role).Error
	if err != nil {
		appLog.Error("Error getting user role: ", err.Error())
		return models.GPUserRole{}, err
	}
	return role, nil
}

// GetPermissionsForLicense resolves the permission slugs held by the user
// associated with the given license_id. Mirrors GetPermissionsForEmail and
// is the preferred lookup — it matches the frontend's loadUserPermissions
// behaviour, so frontend and backend can never disagree about which role
// applies to the active user.
func GetPermissionsForLicense(licenseId string) (hasRole bool, slugs []string, err error) {
	var orgUser models.OrgUser
	if err = DB.Where("license_id = ?", licenseId).First(&orgUser).Error; err != nil {
		return false, nil, nil
	}
	if orgUser.GPRoleId == 0 {
		return false, nil, nil
	}

	var role models.GPUserRole
	if err = DB.Where("id = ?", orgUser.GPRoleId).Preload("Permissions").First(&role).Error; err != nil {
		return false, nil, err
	}

	effective := filterByBaseline(role.Permissions)
	slugs = make([]string, 0, len(effective))
	for _, p := range effective {
		slugs = append(slugs, p.Slug)
	}
	return true, slugs, nil
}

// GetPermissionsForEmail resolves the permission slugs held by the user with
// the given email. Returns hasRole=false when the user has no role assigned —
// callers treat that as bootstrap mode (fresh install, open) mirroring the
// frontend usePermissionCheck behaviour.
func GetPermissionsForEmail(email string) (hasRole bool, slugs []string, err error) {
	var orgUser models.OrgUser
	if err = DB.Where("email = ?", email).First(&orgUser).Error; err != nil {
		// No org_user row yet — treat as no-role (bootstrap) rather than a
		// hard failure. A missing user is expected during initial setup.
		return false, nil, nil
	}
	if orgUser.GPRoleId == 0 {
		return false, nil, nil
	}

	var role models.GPUserRole
	if err = DB.Where("id = ?", orgUser.GPRoleId).Preload("Permissions").First(&role).Error; err != nil {
		return false, nil, err
	}

	effective := filterByBaseline(role.Permissions)
	slugs = make([]string, 0, len(effective))
	for _, p := range effective {
		slugs = append(slugs, p.Slug)
	}
	return true, slugs, nil
}

func getBaseBenefitMaps() []models.GroupBenefitMapper {
	// Define the base list of benefits. Educator is split into GLA_EDU and
	// PTD_EDU because the SchemeCategory already tracks separate GLA and PTD
	// educator benefits and users need to customise their display names
	// independently (e.g. "School Fees Cover" under GLA, "Education
	// Disability" under PTD).
	baseBenefits := []models.GroupBenefitMapper{
		{BenefitName: "Group Life Assurance", BenefitCode: "GLA", BenefitAlias: ""},
		{BenefitName: "Additional Accidental Group Life Assurance", BenefitCode: "AAGLA", BenefitAlias: ""},
		{BenefitName: "Additional Group Life Assurance", BenefitCode: "AGLA", BenefitAlias: ""},
		{BenefitName: "Spouse Group Life Assurance", BenefitCode: "SGLA", BenefitAlias: ""},
		{BenefitName: "Permanent Total Disability", BenefitCode: "PTD", BenefitAlias: ""},
		{BenefitName: "Temporary Total Disability", BenefitCode: "TTD", BenefitAlias: ""},
		{BenefitName: "Personal Health Insurance", BenefitCode: "PHI", BenefitAlias: ""},
		{BenefitName: "Critical Illness", BenefitCode: "CI", BenefitAlias: ""},
		{BenefitName: "Group Family Funeral", BenefitCode: "GFF", BenefitAlias: ""},
		{BenefitName: "GLA Educator", BenefitCode: "GLA_EDU", BenefitAlias: ""},
		{BenefitName: "PTD Educator", BenefitCode: "PTD_EDU", BenefitAlias: ""},
	}
	return baseBenefits
}

// EnsureBaseBenefitMapsSeeded inserts any base benefit maps that don't yet
// exist in group_benefit_mappers, and refreshes the benefit_name of existing
// rows when it differs from the base list (benefit_name is not user-editable,
// so renames in the base list should propagate automatically on startup
// without clobbering the user's benefit_alias / benefit_alias_code).
// Called from main.go after migrations so additions and renames in the base
// list (e.g. AAGLA, AGLA) surface on existing installs without manual work.
func EnsureBaseBenefitMapsSeeded() error {
	// Legacy code migrations: if a previous seed used a different
	// benefit_code for what is now a base entry, rename the code on the
	// existing row so the user's custom alias (and any data referencing
	// it) survives. Only renames when the new code is NOT already
	// present.
	codeRenames := map[string]string{
		"AGLC": "AGLA",    // "Additional Group Life Cover" -> "... Assurance"
		"EDU":  "GLA_EDU", // Split "Educator Risk Rates" into "GLA Educator" + "PTD Educator"; carry any user-set alias to the GLA side and insert PTD_EDU fresh.
	}
	for oldCode, newCode := range codeRenames {
		var oldCount, newCount int64
		if err := DB.Model(&models.GroupBenefitMapper{}).
			Where("benefit_code = ?", oldCode).
			Count(&oldCount).Error; err != nil {
			return err
		}
		if oldCount == 0 {
			continue
		}
		if err := DB.Model(&models.GroupBenefitMapper{}).
			Where("benefit_code = ?", newCode).
			Count(&newCount).Error; err != nil {
			return err
		}
		if newCount > 0 {
			// New code already seeded (e.g. from a clean install). Drop the
			// legacy row to avoid a duplicate entry for the same benefit.
			if err := DB.Where("benefit_code = ?", oldCode).
				Delete(&models.GroupBenefitMapper{}).Error; err != nil {
				return err
			}
			continue
		}
		if err := DB.Model(&models.GroupBenefitMapper{}).
			Where("benefit_code = ?", oldCode).
			Update("benefit_code", newCode).Error; err != nil {
			return err
		}
	}

	base := getBaseBenefitMaps()
	for _, b := range base {
		var existing models.GroupBenefitMapper
		err := DB.Where("benefit_code = ?", b.BenefitCode).First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				row := b
				if err := DB.Create(&row).Error; err != nil {
					return err
				}
				continue
			}
			return err
		}
		if existing.BenefitName != b.BenefitName {
			if err := DB.Model(&models.GroupBenefitMapper{}).
				Where("id = ?", existing.ID).
				Update("benefit_name", b.BenefitName).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func GetGroupPricingIndustriesForQuotes() ([]string, error) {
	var industries []string
	err := DB.Model(models.OccupationClass{}).Distinct("industry").Pluck("industry", &industries).Error
	if err != nil {
		return nil, err
	}
	return industries, nil
}

func GetBenefitEscalationsOptions() ([]string, error) {
	var escalations []string
	err := DB.Model(models.PhiRate{}).Distinct("benefit_escalation_option").Pluck("benefit_escalation_option", &escalations).Error
	if err != nil {
		return nil, err
	}
	return escalations, nil
}

func GetTTDDisabilityDefinitions(riskRateCode string) ([]string, error) {
	var definitions []string
	err := DB.Model(models.TtdRate{}).Where("risk_rate_code = ?", riskRateCode).Distinct("disability_definition").Pluck("disability_definition", &definitions).Error
	if err != nil {
		return nil, err
	}
	return definitions, nil
}

func GetPTDDisabilityDefinitions(riskRateCode string) ([]string, error) {
	var definitions []string
	err := DB.Model(models.PtdRate{}).Where("risk_rate_code = ?", riskRateCode).Distinct("disability_definition").Pluck("disability_definition", &definitions).Error
	if err != nil {
		return nil, err
	}
	return definitions, nil
}

func GetPhiDisabilityDefinitions(riskRateCode string) ([]string, error) {
	var definitions []string
	err := DB.Model(models.PhiRate{}).Where("risk_rate_code  = ?", riskRateCode).Distinct("disability_definition").Pluck("disability_definition", &definitions).Error
	if err != nil {
		return nil, err
	}
	return definitions, nil
}

// GetEducatorBenefitTypes returns the distinct educator benefit codes available
// for a given risk rate code, used to populate the Educator Benefit Type
// dropdown under GLA and PTD in the quote form.
func GetEducatorBenefitTypes(riskRateCode string) ([]string, error) {
	var codes []string
	err := DB.Model(&models.EducatorBenefitStructure{}).
		Where("risk_rate_code = ?", riskRateCode).
		Distinct("educator_benefit_code").
		Order("educator_benefit_code").
		Pluck("educator_benefit_code", &codes).Error
	if err != nil {
		return nil, err
	}
	return codes, nil
}

// educatorCodeForCategory returns the educator_benefit_code that should be used
// for a scheme category's educator benefit calculations. GLA takes priority
// over PTD when both are enabled, matching how the current pricing code shares
// a single educator_benefit_structure between the two benefits.
func educatorCodeForCategory(c *models.SchemeCategory) string {
	if c == nil {
		return ""
	}
	if c.GlaEducatorBenefit == "Yes" && c.GlaEducatorBenefitType != "" {
		return c.GlaEducatorBenefitType
	}
	if c.PtdEducatorBenefit == "Yes" && c.PtdEducatorBenefitType != "" {
		return c.PtdEducatorBenefitType
	}
	return ""
}

/// utility functions

func getJSONTags(obj interface{}) []string {
	var float64Fields []string
	val := reflect.ValueOf(obj)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		// Strip any options like ",omitempty" so the tag matches the JSON key.
		if comma := strings.IndexByte(jsonTag, ','); comma >= 0 {
			jsonTag = jsonTag[:comma]
		}
		float64Fields = append(float64Fields, jsonTag)
	}

	return float64Fields
}

// getStructDBColumns returns the DB column names for an object's exported,
// non-skipped struct fields in declaration order. Honours the gorm:"column:..."
// override; otherwise falls back to snake_case of the field name (matching
// GORM's default convention). Fields tagged json:"-" or gorm:"-" are skipped
// so the output matches what would actually be persisted.
func getStructDBColumns(obj interface{}) []string {
	val := reflect.ValueOf(obj)
	typ := val.Type()
	cols := make([]string, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		gormTag := field.Tag.Get("gorm")
		if gormTag == "-" {
			continue
		}
		colName := ""
		for _, part := range strings.Split(gormTag, ";") {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "column:") {
				colName = strings.TrimPrefix(part, "column:")
				break
			}
		}
		if colName == "" {
			colName = DB.NamingStrategy.ColumnName("", field.Name)
		}
		cols = append(cols, colName)
	}
	return cols
}

func FloatPrecision(x float64, precision float64) float64 {
	return math.Round(x*math.Pow(10, precision)) / math.Pow(10, precision)
}

func GetGroupPricingAgeBands(ctx context.Context) ([]models.GroupPricingAgeBands, error) {
	var variables []models.GroupPricingAgeBands
	err := DB.Find(&variables).Error
	if err != nil {
		return variables, err
	}
	return variables, nil
}

// QuoteMemberGenderSplit is the gender breakdown of the uploaded member
// data for a single quote. Counts exclude rows with an empty gender.
type QuoteMemberGenderSplit struct {
	QuoteID     int     `json:"quote_id"`
	MaleCount   int     `json:"male_count"`
	FemaleCount int     `json:"female_count"`
	OtherCount  int     `json:"other_count"`
	TotalCount  int     `json:"total_count"`
	MaleProp    float64 `json:"male_prop"`
	FemaleProp  float64 `json:"female_prop"`
}

// GetQuoteMemberGenderSplit counts the male/female members uploaded for a
// quote. New Business quotes live in g_pricing_member_data (keyed by
// quote_id); renewals/in-force quotes live in g_pricing_member_data_in_forces
// (keyed by scheme_id). We check the new-business table first; if nothing is
// found and the quote has a scheme_id, we fall back to the in-force table.
func GetQuoteMemberGenderSplit(quoteID int) (QuoteMemberGenderSplit, error) {
	split := QuoteMemberGenderSplit{QuoteID: quoteID}

	type genderRow struct {
		Gender string
		N      int
	}

	countByGender := func(rows []genderRow) {
		for _, r := range rows {
			g := strings.ToUpper(strings.TrimSpace(r.Gender))
			switch g {
			case "M", "MALE":
				split.MaleCount += r.N
			case "F", "FEMALE":
				split.FemaleCount += r.N
			default:
				split.OtherCount += r.N
			}
		}
	}

	// New-business member data.
	var nbRows []genderRow
	if err := DB.Table("g_pricing_member_data").
		Select("gender, COUNT(*) AS n").
		Where("quote_id = ?", quoteID).
		Group("gender").
		Scan(&nbRows).Error; err != nil {
		return split, err
	}
	countByGender(nbRows)

	// If nothing for this quote_id, try the in-force members (scheme-scoped).
	if split.MaleCount+split.FemaleCount+split.OtherCount == 0 {
		var quote models.GroupPricingQuote
		if err := DB.Select("scheme_id").
			Where("id = ?", quoteID).
			First(&quote).Error; err == nil && quote.SchemeID != 0 {
			var ifRows []genderRow
			if err := DB.Table("g_pricing_member_data_in_forces").
				Select("gender, COUNT(*) AS n").
				Where("scheme_id = ?", quote.SchemeID).
				Group("gender").
				Scan(&ifRows).Error; err != nil {
				return split, err
			}
			countByGender(ifRows)
		}
	}

	split.TotalCount = split.MaleCount + split.FemaleCount + split.OtherCount
	if split.TotalCount > 0 {
		split.MaleProp = float64(split.MaleCount) / float64(split.TotalCount)
		split.FemaleProp = float64(split.FemaleCount) / float64(split.TotalCount)
	}
	return split, nil
}

func GetGroupPricingBenefits(ctx context.Context) ([]models.GroupBusinessBenefits, error) {
	var variables []models.GroupBusinessBenefits
	err := DB.Find(&variables).Error
	if err != nil {
		return variables, err
	}
	return variables, nil
}

func getGroupRiskQuotingFinancialYear(commencementDate time.Time, yearEndMonth int) (int, int) {
	date := commencementDate

	var startYear, endYear int

	if int(date.Month()) <= yearEndMonth {
		startYear = date.Year() - 1
		endYear = date.Year()
	} else {
		startYear = date.Year()
		endYear = date.Year() + 1
	}

	return startYear, endYear
}

func addMappers(groupQuote *models.GroupPricingQuote, benefitMaps []models.GroupBenefitMapper) {
	//GLA
	for _, benefit := range benefitMaps {
		for i, _ := range groupQuote.SchemeCategories {
			if benefit.BenefitCode == "GLA" {
				//groupQuote.GlaAlias = benefit.BenefitAlias
				groupQuote.SchemeCategories[i].GlaAlias = benefit.BenefitName
				if benefit.BenefitAlias != "" {
					groupQuote.SchemeCategories[i].GlaAlias = benefit.BenefitAlias
				}
			}
			if benefit.BenefitCode == "SGLA" {
				groupQuote.SchemeCategories[i].SglaAlias = benefit.BenefitName
				if benefit.BenefitAlias != "" {
					groupQuote.SchemeCategories[i].SglaAlias = benefit.BenefitAlias
				}
			}
			if benefit.BenefitCode == "PTD" {
				groupQuote.SchemeCategories[i].PtdAlias = benefit.BenefitName
				if benefit.BenefitAlias != "" {
					groupQuote.SchemeCategories[i].PtdAlias = benefit.BenefitAlias
				}
			}
			if benefit.BenefitCode == "TTD" {
				groupQuote.SchemeCategories[i].TtdAlias = benefit.BenefitName
				if benefit.BenefitAlias != "" {
					groupQuote.SchemeCategories[i].TtdAlias = benefit.BenefitAlias
				}
			}
			if benefit.BenefitCode == "PHI" {
				groupQuote.SchemeCategories[i].PhiAlias = benefit.BenefitName
				if benefit.BenefitAlias != "" {
					groupQuote.SchemeCategories[i].PhiAlias = benefit.BenefitAlias
				}
			}
			if benefit.BenefitCode == "CI" {
				groupQuote.SchemeCategories[i].CiAlias = benefit.BenefitName
				if benefit.BenefitAlias != "" {
					groupQuote.SchemeCategories[i].CiAlias = benefit.BenefitAlias
				}
			}
			if benefit.BenefitCode == "GFF" {
				groupQuote.SchemeCategories[i].FamilyFuneralAlias = benefit.BenefitName
				if benefit.BenefitAlias != "" {
					groupQuote.SchemeCategories[i].FamilyFuneralAlias = benefit.BenefitAlias
				}
			}
		}
	}

}

// GetDistinctWaitingPeriods returns a slice of ints that represent the distinct waiting_period values from the specified table
// Valid table names are: "gla_rates", "ptd_rates", "ci_rates", "ttd_rates", "phi_rates"
func GetDistinctWaitingPeriods(tableName, riskRateCode string) ([]int, error) {
	var waitingPeriods []int

	// Validate table name
	validTables := map[string]bool{
		"gla_rates": true,
		"ptd_rates": true,
		"ci_rates":  true,
		"ttd_rates": true,
		"phi_rates": true,
	}

	if !validTables[tableName] {
		return nil, fmt.Errorf("invalid table name: %s", tableName)
	}

	// Query the database for distinct waiting_period values
	err := DB.Table(tableName).Where("risk_rate_code = ?", riskRateCode).Distinct("waiting_period").Pluck("waiting_period", &waitingPeriods).Error
	if err != nil {
		return nil, err
	}

	return waitingPeriods, nil
}

// GetDistinctGlaBenefitTypes returns a slice of strings that represent the distinct benefit_type values from the gla_rates table
func GetDistinctGlaBenefitTypes(riskRateCode string) ([]string, error) {
	var benefitTypes []string

	err := DB.Table("gla_rates").Where("risk_rate_code = ?", riskRateCode).Distinct("benefit_type").Pluck("benefit_type", &benefitTypes).Error
	if err != nil {
		return nil, err
	}

	return benefitTypes, nil
}

// GetDistinctDeferredPeriods returns a slice of ints that represent the distinct deferred_period values from the specified table
// Valid table names are: "gla_rates", "ptd_rates", "ci_rates", "ttd_rates", "phi_rates"
func GetDistinctDeferredPeriods(tableName, riskRateCode string) ([]int, error) {
	var deferredPeriods []int

	// Validate table name
	validTables := map[string]bool{
		"gla_rates": true,
		"ptd_rates": true,
		"ci_rates":  true,
		"ttd_rates": true,
		"phi_rates": true,
	}

	if !validTables[tableName] {
		return nil, fmt.Errorf("invalid table name: %s", tableName)
	}

	// Query the database for distinct deferred_period values
	err := DB.Table(tableName).Where("risk_rate_code = ?", riskRateCode).Distinct("deferred_period").Pluck("deferred_period", &deferredPeriods).Error
	if err != nil {
		return nil, err
	}

	return deferredPeriods, nil
}

// GetDistinctNormalRetirementAges returns a slice of ints that represent the distinct normal_retirement_age values from the phi_rates table
func GetDistinctNormalRetirementAges() ([]int, error) {
	var normalRetirementAges []int

	// Query the database for distinct normal_retirement_age values from phi_rates table
	err := DB.Table("phi_rates").Distinct("normal_retirement_age").Pluck("normal_retirement_age", &normalRetirementAges).Error
	if err != nil {
		return nil, err
	}

	return normalRetirementAges, nil
}

// RiskTypes holds the distinct risk types from different tables
type RiskTypes struct {
	PhiRates []string `json:"phi_rates"`
	TtdRates []string `json:"ttd_rates"`
	PtdRates []string `json:"ptd_rates"`
}

// GetDistinctRiskTypes returns the distinct risk_type values from the phi_rates, ttd_rates, and ptd_rates tables separately
func GetDistinctRiskTypes() (RiskTypes, error) {
	var result RiskTypes
	var phiRiskTypes []string
	var ttdRiskTypes []string
	var ptdRiskTypes []string

	// Query the database for distinct risk_type values from phi_rates table
	err := DB.Table("phi_rates").Distinct("risk_type").Pluck("risk_type", &phiRiskTypes).Error
	if err != nil {
		return RiskTypes{}, err
	}

	// Query the database for distinct risk_type values from ttd_rates table
	err = DB.Table("ttd_rates").Distinct("risk_type").Pluck("risk_type", &ttdRiskTypes).Error
	if err != nil {
		return RiskTypes{}, err
	}

	// Query the database for distinct risk_type values from ptd_rates table
	err = DB.Table("ptd_rates").Distinct("risk_type").Pluck("risk_type", &ptdRiskTypes).Error
	if err != nil {
		return RiskTypes{}, err
	}

	result.PhiRates = phiRiskTypes
	result.TtdRates = ttdRiskTypes
	result.PtdRates = ptdRiskTypes

	return result, nil
}

// GetHistoricalCredibilityData returns all historical credibility data
func GetHistoricalCredibilityData() ([]models.HistoricalCredibilityData, error) {
	var data []models.HistoricalCredibilityData

	// Query the database for all historical credibility data
	err := DB.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

// SearchSchemeMembers searches for members in a scheme that match the query string
func SearchSchemeMembers(schemeId string, quoteId string, query string) ([]models.GPricingMemberDataInForce, error) {
	logger := appLog.WithFields(map[string]interface{}{
		"scheme_id": schemeId,
		"quote_id":  quoteId,
		"query":     query,
		"function":  "SearchSchemeMembers",
	})

	logger.Debug("Searching for scheme members")

	var members []models.GPricingMemberDataInForce

	// Convert quoteId to int
	quoteIdInt, err := strconv.Atoi(quoteId)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to convert quote_id to integer")
		return nil, err
	}

	// Search for members with a partial match on member_name
	err = DB.Where("scheme_id = ? AND quote_id = ? AND member_name LIKE ? AND status = 'Active'",
		schemeId, quoteIdInt, "%"+query+"%").
		Find(&members).Error

	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to search for scheme members")
		return nil, err
	}

	logger.WithField("member_count", len(members)).Info("Successfully searched for scheme members")
	return members, nil
}

func UpdateGroupSchemeStatus(schemeId string, schemeStatus models.SchemeStatusUpdate, user models.AppUser) error {
	// Use a transaction to ensure audit and update are atomic
	return DB.Transaction(func(tx *gorm.DB) error {
		var before models.GroupScheme
		if err := tx.Where("id = ?", schemeId).First(&before).Error; err != nil {
			return err
		}

		// Determine updates
		updates := map[string]interface{}{}
		statusChanged := before.Status != schemeStatus.Status
		messageChanged := before.SchemeStatusMessage != schemeStatus.SchemeStatusMessage && schemeStatus.SchemeStatusMessage != ""

		if statusChanged {
			updates["status"] = schemeStatus.Status
		}
		if statusChanged || messageChanged {
			updates["scheme_status_message"] = schemeStatus.SchemeStatusMessage
		}

		// Specialized audit row for status changes
		if statusChanged || messageChanged {
			audit := models.GroupSchemeStatusAudit{
				SchemeID:      before.ID,
				OldStatus:     before.Status,
				NewStatus:     schemeStatus.Status,
				StatusMessage: schemeStatus.SchemeStatusMessage,
				ChangedBy:     user.UserName,
				ChangedAt:     time.Now(),
			}
			if err := tx.Create(&audit).Error; err != nil {
				return err
			}
		}

		// Apply updates
		if len(updates) > 0 {
			if err := tx.Model(&models.GroupScheme{}).Where("id = ?", schemeId).Updates(updates).Error; err != nil {
				return err
			}
		}

		// Load after state for generic audit
		var after models.GroupScheme
		if err := tx.Where("id = ?", schemeId).First(&after).Error; err != nil {
			return err
		}

		// Generic audit log
		if err := writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "group_schemes",
			EntityID:  schemeId,
			Action:    "UPDATE",
			ChangedBy: user.UserName,
		}, before, after); err != nil {
			return err
		}

		return nil
	})
}

// DeactivateSchemeMember deletes a member from a scheme by scheme id and member id
func DeactivateSchemeMember(schemeId string, member models.GPricingMemberDataInForce, user models.AppUser) error {
	logger := appLog.WithFields(map[string]interface{}{
		"scheme_id": schemeId,
		"member_id": member.ID,
		"function":  "DeactivateSchemeMember",
	})

	logger.Debug("Deleting scheme member")

	return DB.Transaction(func(tx *gorm.DB) error {
		var before models.GPricingMemberDataInForce
		if err := tx.Where("id = ? AND scheme_id = ?", member.ID, schemeId).First(&before).Error; err != nil {
			return err
		}
		beforeStatus := before.Status
		member.Status = "Inactive"
		if err := tx.Save(&member).Error; err != nil {
			logger.WithField("error", err.Error()).Error("Failed to deactivate scheme member")
			return err
		}
		var after models.GPricingMemberDataInForce
		if err := tx.Where("id = ?", member.ID).First(&after).Error; err != nil {
			return err
		}
		// Only write audit if status actually changed
		if beforeStatus != after.Status {
			if err := writeAudit(tx, AuditContext{
				Area:      "group-pricing",
				Entity:    "g_pricing_member_data_in_forces",
				EntityID:  strconv.Itoa(member.ID),
				Action:    "UPDATE",
				ChangedBy: user.UserName,
			}, before, after); err != nil {
				return err
			}

			// Record the exit as a member history event.
			exitDetails, _ := json.Marshal(map[string]interface{}{
				"scheme_id":   before.SchemeId,
				"scheme_name": before.SchemeName,
				"exit_date": func() string {
					if before.ExitDate != nil {
						return before.ExitDate.Format("2006-01-02")
					}
					return ""
				}(),
				"effective_exit_date": func() string {
					if before.EffectiveExitDate != nil {
						return before.EffectiveExitDate.Format("2006-01-02")
					}
					return ""
				}(),
				"previous_status": beforeStatus,
			})
			_ = tx.Create(&models.MemberActivity{
				MemberID:       before.ID,
				MemberIDNumber: before.MemberIdNumber,
				Type:           "exit",
				Title:          "Member Exited Scheme",
				Description:    fmt.Sprintf("Member deactivated from scheme '%s'", before.SchemeName),
				Details:        exitDetails,
				PerformedBy:    user.UserName,
			})
		}

		// Clear the cache
		GroupPricingCache.Clear()

		logger.Info("Successfully deleted scheme member")
		return nil
	})
}

// MemberIndicativeDataSet represents the payload for updating member stats on a GroupPricingQuote

// UpdateGroupPricingQuoteMemberStats updates the member statistics on a GroupPricingQuote by ID
func UpdateGroupPricingQuoteMemberStats(input []models.MemberIndicativeDataSet, user models.AppUser) (models.GroupPricingQuote, error) {
	var quote models.GroupPricingQuote
	if len(input) == 0 {
		return quote, errors.New("no data provided")
	}
	quoteID := input[0].QuoteID

	// Find quote by primary ID using the provided quote_id
	if err := DB.Where("id = ?", quoteID).First(&quote).Error; err != nil {
		return quote, err
	}

	// Transactionally replace rows and audit digest
	type digest struct {
		QuoteID  int `json:"quote_id"`
		RowCount int `json:"row_count"`
	}

	err := DB.Transaction(func(tx *gorm.DB) error {
		var beforeCount int64
		if err := tx.Model(&models.MemberIndicativeDataSet{}).Where("quote_id = ?", quoteID).Count(&beforeCount).Error; err != nil {
			return err
		}

		// Replace rows
		if err := tx.Where("quote_id = ?", quote.ID).Delete(&models.MemberIndicativeDataSet{}).Error; err != nil {
			return err
		}
		for _, data := range input {
			if err := tx.Create(&data).Error; err != nil {
				return err
			}
		}

		// Update quote audit fields
		if err := tx.Model(&models.GroupPricingQuote{}).
			Where("id = ?", quote.ID).
			Updates(map[string]interface{}{
				"member_indicative_data": true,
				"modified_by":            user.UserName,
				"modification_date":      time.Now(),
			}).Error; err != nil {
			return err
		}

		var afterCount int64
		if err := tx.Model(&models.MemberIndicativeDataSet{}).Where("quote_id = ?", quoteID).Count(&afterCount).Error; err != nil {
			return err
		}

		before := digest{QuoteID: quoteID, RowCount: int(beforeCount)}
		after := digest{QuoteID: quoteID, RowCount: int(afterCount)}
		if err := writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "member_indicative_data_sets",
			EntityID:  strconv.Itoa(quoteID),
			Action:    "UPDATE",
			ChangedBy: user.UserName,
		}, before, after); err != nil {
			return err
		}

		return nil
	})
	return quote, err
}

// UpdateGroupPricingQuoteIndicativeFlag updates the MemberIndicativeData boolean on a GroupPricingQuote by ID
func UpdateGroupPricingQuoteIndicativeFlag(quoteID int, enabled bool, user models.AppUser) (models.GroupPricingQuote, error) {
	var quote models.GroupPricingQuote

	// Ensure quote exists first
	if err := DB.Where("id = ?", quoteID).First(&quote).Error; err != nil {
		return quote, err
	}

	// Prepare before/after for audit (only the fields we change)
	type snapshot struct {
		ID                   int  `json:"id"`
		MemberIndicativeData bool `json:"member_indicative_data"`
	}
	before := snapshot{ID: quote.ID, MemberIndicativeData: quote.MemberIndicativeData}

	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.GroupPricingQuote{}).
			Where("id = ?", quote.ID).
			Updates(map[string]interface{}{
				"member_indicative_data": enabled,
				"modified_by":            user.UserName,
				"modification_date":      time.Now(),
			}).Error; err != nil {
			return err
		}

		// Reload the updated value
		if err := tx.Where("id = ?", quote.ID).First(&quote).Error; err != nil {
			return err
		}

		after := snapshot{ID: quote.ID, MemberIndicativeData: quote.MemberIndicativeData}
		if err := writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "group_pricing_quotes",
			EntityID:  strconv.Itoa(quote.ID),
			Action:    "UPDATE",
			ChangedBy: user.UserName,
		}, before, after); err != nil {
			return err
		}
		return nil
	})
	return quote, err
}

// GetGroupSchemeStatusAudit returns the status change history for a scheme
func GetGroupSchemeStatusAudit(schemeId string) ([]models.GroupSchemeStatusAudit, error) {
	var audits []models.GroupSchemeStatusAudit
	if err := DB.Where("scheme_id = ?", schemeId).Order("changed_at desc").Find(&audits).Error; err != nil {
		return nil, err
	}
	return audits, nil
}

func DeleteIndicativeData(quoteId int, user models.AppUser) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("quote_id = ?", quoteId).Delete(&models.MemberIndicativeDataSet{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// SaveExperienceRateOverrides replaces the full set of experience-rate
// override rows for a quote. The frontend submits the in-memory list verbatim
// after each Save, so the simplest and least error-prone strategy is
// delete-then-insert inside one transaction. CreatedBy is stamped from the
// caller; UpdatedBy is set on every save so it always reflects the most
// recent author.
func SaveExperienceRateOverrides(quoteId int, rows []models.GroupPricingExperienceRateOverride, user models.AppUser) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("quote_id = ?", quoteId).
			Delete(&models.GroupPricingExperienceRateOverride{}).Error; err != nil {
			return err
		}
		now := time.Now()
		for i := range rows {
			rows[i].QuoteId = quoteId
			rows[i].ID = 0
			rows[i].CreatedAt = now
			rows[i].UpdatedAt = now
			if rows[i].CreatedBy == "" {
				rows[i].CreatedBy = user.UserEmail
			}
			rows[i].UpdatedBy = user.UserEmail
		}
		if len(rows) == 0 {
			return nil
		}
		return tx.Create(&rows).Error
	})
}

// buildExperienceRateOverrideLookup converts a flat slice of override rows
// into a category -> benefit -> row map for O(1) lookup during member rating.
// The two-level map is keyed exactly as the rate-population code references
// (member.Category for the outer key, the benefit constant for the inner).
func buildExperienceRateOverrideLookup(rows []models.GroupPricingExperienceRateOverride) map[string]map[string]models.GroupPricingExperienceRateOverride {
	lookup := make(map[string]map[string]models.GroupPricingExperienceRateOverride, len(rows))
	for _, r := range rows {
		if _, ok := lookup[r.SchemeCategory]; !ok {
			lookup[r.SchemeCategory] = make(map[string]models.GroupPricingExperienceRateOverride, 7)
		}
		lookup[r.SchemeCategory][r.Benefit] = r
	}
	return lookup
}

// averageOverrideCredibilityByBenefit collapses the per-(category, benefit)
// override rows to one credibility per benefit by simple average across
// categories. Used to populate the per-benefit credibility columns on
// HistoricalCredibilityData (one row per quote) when ExperienceRating ==
// "Override". Benefits with no override rows are absent from the result map.
func averageOverrideCredibilityByBenefit(rows []models.GroupPricingExperienceRateOverride) map[string]float64 {
	sums := map[string]float64{}
	counts := map[string]int{}
	for _, r := range rows {
		sums[r.Benefit] += r.Credibility
		counts[r.Benefit]++
	}
	out := make(map[string]float64, len(sums))
	for b, s := range sums {
		if counts[b] > 0 {
			out[b] = s / float64(counts[b])
		}
	}
	return out
}

// UpdateMemberInForce updates a member's details and writes an audit trail with before/after differences
func UpdateMemberInForce(memberID string, input models.GPricingMemberDataInForce, user models.AppUser) (models.GPricingMemberDataInForce, error) {
	var current models.GPricingMemberDataInForce
	// Try find by member_id_number first; if not found, try by numeric ID
	if err := DB.Where("member_id_number = ?", memberID).First(&current).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err2 := DB.Where("id = ?", memberID).First(&current).Error; err2 != nil {
				return current, err2
			}
		} else {
			return current, err
		}
	}

	// Build an update map from input, excluding immutable fields
	// We go through JSON marshal/unmarshal to a map to avoid zero-value skipping problems
	var update map[string]interface{}
	{
		b, _ := json.Marshal(input)
		_ = json.Unmarshal(b, &update)
		// Remove non-updatable fields
		delete(update, "id")
		delete(update, "year")
		delete(update, "scheme_name")
		delete(update, "creation_date")
		delete(update, "created_by")
		delete(update, "quote_id")
		delete(update, "scheme_id")
		delete(update, "scheme_category_details")
		// Keep status editable via dedicated flows; avoid silent changes here
		// Internal flags
		delete(update, "is_original_member")

		// Flatten embedded benefits struct into column names (gorm uses embeddedPrefix: benefits_)
		if ben, ok := update["benefits"].(map[string]interface{}); ok {
			for k, v := range ben {
				update["benefits_"+k] = v
			}
			delete(update, "benefits")
		}

		update["address_line1"] = input.AddressLine1 // to fix the column mismatch problem
		update["address_line2"] = input.AddressLine2
		delete(update, "address_line_1")
		delete(update, "address_line_2")

		// Normalize any time fields that may have been converted to strings by the JSON roundtrip.
		// MySQL DATETIME won't accept RFC3339 strings with 'T'/'Z'. Ensure we pass time.Time values.
		parseFlexible := func(v any) (time.Time, bool) {
			s, ok := v.(string)
			if !ok {
				return time.Time{}, false
			}
			s = strings.TrimSpace(s)
			if s == "" {
				return time.Time{}, true // treat empty as zero time
			}
			layouts := []string{
				time.RFC3339Nano,
				time.RFC3339,
				"2006-01-02 15:04:05",
				"2006-01-02",
			}
			for _, layout := range layouts {
				if t, err := time.Parse(layout, s); err == nil {
					return t, true
				}
			}
			return time.Time{}, false
		}

		// List of date/time keys we care about
		for _, key := range []string{"effective_exit_date", "exit_date", "entry_date", "creation_date", "date_of_birth"} {
			if v, ok := update[key]; ok {
				if t, ok2 := parseFlexible(v); ok2 {
					if t.IsZero() {
						// If empty string provided, drop the field to avoid pushing invalid zero time
						delete(update, key)
					} else {
						update[key] = t
					}
				}
			}
		}
	}

	// Prepare before snapshot for audit
	before := current

	// Execute transactional update and audit
	err := DB.Transaction(func(tx *gorm.DB) error {
		if len(update) > 0 {
			if err := tx.Model(&models.GPricingMemberDataInForce{}).Where("id = ?", current.ID).Updates(update).Error; err != nil {
				return err
			}
		}
		// Reload after state
		if err := tx.Where("id = ?", current.ID).First(&current).Error; err != nil {
			return err
		}

		// Write generic audit log for this member change
		if err := writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "g_pricing_member_data_in_forces",
			EntityID:  strconv.Itoa(current.ID),
			Action:    "UPDATE",
			ChangedBy: user.UserName,
		}, before, current); err != nil {
			return err
		}

		// Log structured member activity if specific fields changed
		logStructuredActivity(tx, before, current, user.UserName)

		return nil
	})
	return current, err
}

func logStructuredActivity(tx *gorm.DB, before, after models.GPricingMemberDataInForce, performedBy string) {
	// Salary Adjustment
	if before.AnnualSalary != after.AnnualSalary && before.AnnualSalary > 0 {
		details, _ := json.Marshal(map[string]interface{}{
			"previousValue": before.AnnualSalary,
			"newValue":      after.AnnualSalary,
			"reason":        "Salary Adjustment", // Could be more dynamic if we had a reason field
			"effectiveDate": time.Now().Format("2006-01-02"),
		})
		_ = tx.Create(&models.MemberActivity{
			MemberID:       after.ID,
			MemberIDNumber: after.MemberIdNumber,
			Type:           "salary_change",
			Title:          "Salary Adjustment",
			Description:    "Annual salary increase processed",
			Details:        details,
			PerformedBy:    performedBy,
		})
	}

	// Status Change
	if before.Status != after.Status {
		details, _ := json.Marshal(map[string]interface{}{
			"previousStatus": before.Status,
			"newStatus":      after.Status,
			"effectiveDate":  time.Now().Format("2006-01-02"),
		})
		_ = tx.Create(&models.MemberActivity{
			MemberID:       after.ID,
			MemberIDNumber: after.MemberIdNumber,
			Type:           "status_change",
			Title:          "Status Updated",
			Description:    fmt.Sprintf("Member status changed from %s to %s", before.Status, after.Status),
			Details:        details,
			PerformedBy:    performedBy,
		})
	}

	// Contact Update - Address
	if (before.AddressLine1 != after.AddressLine1) || (before.AddressLine2 != after.AddressLine2) {
		details, _ := json.Marshal(map[string]interface{}{
			"previousAddress": strings.TrimSpace(before.AddressLine1 + " " + before.AddressLine2),
			"newAddress":      strings.TrimSpace(after.AddressLine1 + " " + after.AddressLine2),
		})
		_ = tx.Create(&models.MemberActivity{
			MemberID:       after.ID,
			MemberIDNumber: after.MemberIdNumber,
			Type:           "contact_update",
			Title:          "Address Updated",
			Description:    "Member updated residential address",
			Details:        details,
			PerformedBy:    performedBy,
		})
	}

	// Contact Update - Phone
	if before.PhoneNumber != after.PhoneNumber {
		details, _ := json.Marshal(map[string]interface{}{
			"previousPhone": before.PhoneNumber,
			"newPhone":      after.PhoneNumber,
		})
		_ = tx.Create(&models.MemberActivity{
			MemberID:       after.ID,
			MemberIDNumber: after.MemberIdNumber,
			Type:           "contact_update",
			Title:          "Phone Number Updated",
			Description:    "Member updated contact phone number",
			Details:        details,
			PerformedBy:    performedBy,
		})
	}

	// Benefit changes (simple detection for now)
	if !reflect.DeepEqual(before.Benefits, after.Benefits) {
		// This is a bit complex as we have many benefits.
		// For now, let's log a generic benefit change if any benefit field changed.
		details, _ := json.Marshal(map[string]interface{}{
			"action": "Updated",
		})
		_ = tx.Create(&models.MemberActivity{
			MemberID:       after.ID,
			MemberIDNumber: after.MemberIdNumber,
			Type:           "benefit_change",
			Title:          "Benefits Updated",
			Description:    "Member benefits configuration updated",
			Details:        details,
			PerformedBy:    performedBy,
		})
	}
}

func GetQuoteTableDataExcel(quoteId int, tableType string) ([]byte, error) {
	var excelData []byte
	var tableName string
	var columns []string
	switch tableType {
	case "member_data":
		var quote models.GroupPricingQuote
		if err := DB.Where("id = ?", quoteId).First(&quote).Error; err != nil {
			return excelData, err
		}
		if quote.QuoteType == "New Business" {
			tableName = "g_pricing_member_data"
			columns = getStructDBColumns(models.GPricingMemberData{})
		} else {
			tableName = "g_pricing_member_data_in_forces"
			columns = getStructDBColumns(models.GPricingMemberDataInForce{})
		}
	case "member_rating_results":
		tableName = "member_rating_results"
		columns = getStructDBColumns(models.MemberRatingResult{})
	case "bordereaux":
		// Bordereaux is no longer persisted — project rows on-the-fly and
		// stream them to xlsx in struct-field order.
		return ExportQuoteBordereauxXLSX(quoteId)
	case "member_premium_schedules":
		tableName = "member_premium_schedules"
		columns = getStructDBColumns(models.MemberPremiumSchedule{})
	default:
		return nil, fmt.Errorf("invalid table type: %s", tableType)
	}

	// Project columns in struct field order so the Excel output matches the
	// Go model layout regardless of the order columns were created in the DB.
	colList := strings.Join(quoteIdent(columns), ", ")
	dQuery := fmt.Sprintf("SELECT %s FROM %s WHERE quote_id = %d", colList, tableName, quoteId)
	excelData, err := exportTableToExcel(dQuery)
	if err != nil {
		return nil, err
	}

	return excelData, nil
}

// quoteIdent backtick-quotes each identifier so columns whose names happen to
// collide with reserved words still parse. Backticks are accepted by MySQL;
// PostgreSQL and SQL Server use double quotes — for portability we intentionally
// leave bare identifiers unless the underlying driver complains. Currently this
// is a passthrough that just returns its input.
func quoteIdent(cols []string) []string {
	return cols
}

// -------------------------
// Beneficiaries - Services
// -------------------------

// List beneficiaries for a member
func GetBeneficiariesByMemberID(memberID int) ([]models.Beneficiary, error) {
	var list []models.Beneficiary
	if err := DB.Where("member_id = ?", memberID).Order("id asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// Get a single beneficiary ensuring it belongs to the member
func GetBeneficiaryByID(memberID int, id int) (models.Beneficiary, error) {
	var b models.Beneficiary
	if err := DB.Where("id = ? AND member_id = ?", id, memberID).First(&b).Error; err != nil {
		return b, err
	}
	return b, nil
}

// validateAllocation ensures the total allocation for a member does not exceed 100%
func validateAllocation(tx *gorm.DB, memberID int, excludeID *int, newAlloc float64) error {
	var total float64
	q := tx.Model(&models.Beneficiary{}).Where("member_id = ?", memberID)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	if err := q.Select("COALESCE(SUM(allocation_percentage),0)").Scan(&total).Error; err != nil {
		return err
	}
	if total+newAlloc > 100.0+1e-9 { // allow tiny epsilon
		return fmt.Errorf("total allocation percentage for member %d exceeds 100%%", memberID)
	}
	return nil
}

// -------------------------
// Benefit Summary - DTOs and Services
// -------------------------

// MemberBenefitSummaryDTO is a normalized response for a member's benefit summary
// per current requirements: show for non-GFF benefits the benefit type (mapped name),
// salary multiple, and the covered sum assured computed from multiples. For GFF,
// provide a breakdown by lives.
type MemberBenefitSummaryDTO struct {
	MemberID int    `json:"member_id"`
	SchemeID int    `json:"scheme_id"`
	Source   string `json:"source"` // "in_force" | "quote"
	Status   string `json:"status"`
	// Additional member and scheme context requested by client
	MemberIdNumber        string                   `json:"member_id_number"`
	MemberName            string                   `json:"member_name"`
	SchemeCategory        string                   `json:"scheme_category"`
	AnnualSalary          float64                  `json:"annual_salary"`
	AnnualPremium         float64                  `json:"annual_premium"`
	MonthlyPremium        float64                  `json:"monthly_premium"`
	PremiumSalaryProp     float64                  `json:"premium_salary_prop"`
	FuneralAnnualPremium  float64                  `json:"funeral_annual_premium"`
	FuneralMonthlyPremium float64                  `json:"funeral_monthly_premium"`
	SchemeName            string                   `json:"scheme_name"`
	Benefits              []BenefitMultipleItemDTO `json:"benefits"`
	GFF                   *GffBreakdownDTO         `json:"gff,omitempty"`
}

// BenefitMultipleItemDTO represents a single non-GFF benefit with salary multiple and computed cover
type BenefitMultipleItemDTO struct {
	Code              string   `json:"code"`
	Name              string   `json:"name"`
	IsActive          bool     `json:"is_active"`
	SalaryMultiple    *float64 `json:"salary_multiple,omitempty"`
	CoveredSumAssured *float64 `json:"covered_sum_assured,omitempty"` // For TTD/PHI this is monthly benefit
}

// GffBreakdownDTO represents funeral cover amounts by life type
type GffBreakdownDTO struct {
	Currency   string   `json:"currency,omitempty"`
	MainMember *float64 `json:"main_member,omitempty"`
	Spouse     *float64 `json:"spouse,omitempty"`
	Children   *float64 `json:"children,omitempty"`
	Parents    *float64 `json:"parents,omitempty"`
	Dependants *float64 `json:"dependants,omitempty"`
	// Counts of covered lives (when available from scheme/category configuration)
	ChildrenCount   *int `json:"children_count,omitempty"`
	ParentsCount    *int `json:"parents_count,omitempty"`
	DependantsCount *int `json:"dependants_count,omitempty"`
}

// GetMemberBenefitSummaryInForce returns the in-force benefit summary for a member by member primary key.
// It prioritizes financials from Bordereaux using (SchemeId, MemberName, Category).
func GetMemberBenefitSummaryInForce(memberID int) (MemberBenefitSummaryDTO, error) {
	var dto MemberBenefitSummaryDTO

	var m models.GPricingMemberDataInForce
	if err := DB.Where("id = ?", memberID).First(&m).Error; err != nil {
		return dto, err
	}

	var quote models.GroupPricingQuote
	err := DB.Where("scheme_id = ? and scheme_quote_status = ?", m.SchemeId, models.StatusInEffect).First(&quote).Error

	var schemeCategory models.SchemeCategory
	err = DB.Where("quote_id = ? AND scheme_category = ?", quote.ID, m.SchemeCategory).First(&schemeCategory).Error

	if err != nil {
		// Retry with case-insensitive match as a fallback (handles capitalization mismatches)
		err = DB.Where("quote_id = ? AND LOWER(scheme_category) = LOWER(?)", quote.ID, m.SchemeCategory).First(&schemeCategory).Error
	}

	restriction, restrictionErr := GetRestrictionByRiskRateCode(quote.RiskRateCode)
	if restrictionErr != nil {
		appLog.Error("Error retrieving restriction for risk rate code: ", quote.RiskRateCode, " error: ", restrictionErr.Error())
	}

	var memberRatingResultSummary models.MemberRatingResultSummary
	err = DB.Where("quote_id = ? AND category = ?", quote.ID, m.SchemeCategory).First(&memberRatingResultSummary).Error

	if err != nil {
		// Retry with case-insensitive match as a fallback (handles capitalization mismatches)
		err = DB.Where("quote_id = ? AND LOWER(category) = LOWER(?)", quote.ID, m.SchemeCategory).First(&memberRatingResultSummary).Error
	}

	dto.MemberID = m.ID
	dto.SchemeID = m.SchemeId
	dto.Source = "in_force"
	dto.Status = m.Status
	// Populate additional requested fields
	dto.MemberIdNumber = m.MemberIdNumber
	dto.MemberName = m.MemberName
	dto.SchemeCategory = m.SchemeCategory
	dto.AnnualSalary = m.AnnualSalary
	premiumSalaryProp := 0.0
	if memberRatingResultSummary.TotalAnnualSalary > 0 {
		premiumSalaryProp = memberRatingResultSummary.ExpTotalAnnualPremiumExclFuneral / memberRatingResultSummary.TotalAnnualSalary
	}
	dto.AnnualPremium = utils.FloatPrecision(m.AnnualSalary*premiumSalaryProp, AccountingPrecision)
	dto.MonthlyPremium = utils.FloatPrecision(dto.AnnualPremium/12.0, AccountingPrecision)
	dto.PremiumSalaryProp = premiumSalaryProp
	dto.FuneralAnnualPremium = utils.FloatPrecision(memberRatingResultSummary.ExpTotalFunAnnualPremiumPerMember, AccountingPrecision)
	dto.FuneralMonthlyPremium = utils.FloatPrecision(memberRatingResultSummary.ExpTotalFunMonthlyPremiumPerMember, AccountingPrecision)
	dto.SchemeName = m.SchemeName

	// Benefit name mapper (code -> name)
	nameFor := func(code string) string {
		maps, err := GetBenefitMaps()
		if err != nil || len(maps) == 0 {
			return code
		}
		for _, mm := range maps {
			if strings.EqualFold(mm.BenefitCode, code) || strings.EqualFold(mm.BenefitAliasCode, code) {
				if strings.TrimSpace(mm.BenefitAlias) != "" {
					return mm.BenefitAlias
				}
				if strings.TrimSpace(mm.BenefitName) != "" {
					return mm.BenefitName
				}
				return code
			}
		}
		return code
	}

	// helpers
	firstNonZero := func(vals ...float64) *float64 {
		for _, v := range vals {
			if v > 0 {
				vv := v
				return &vv
			}
		}
		return nil
	}
	nz := func(v float64) *float64 {
		if v > 0 {
			vv := v
			return &vv
		}
		return nil
	}
	// ptr returns a pointer to v regardless of zero-ness (useful when we want to include zero values explicitly)
	ptr := func(v float64) *float64 { vv := v; return &vv }

	addItem := func(code string, enabled bool, multiple *float64, covered *float64) {
		item := BenefitMultipleItemDTO{Code: code, Name: nameFor(code), IsActive: enabled}
		if multiple != nil {
			item.SalaryMultiple = multiple
		}
		if covered != nil {
			item.CoveredSumAssured = covered
		}
		dto.Benefits = append(dto.Benefits, item)
	}

	annualSalary := m.AnnualSalary

	// Non-GFF benefits
	// Always include benefit items with IsActive reflecting enabled status
	{
		enabled := m.Benefits.GlaEnabled
		mult := firstNonZero(m.Benefits.GlaMultiple, m.Benefits.GlaMultiple)
		var sa *float64
		if enabled && mult != nil {
			sa = nz(math.Min(annualSalary*(*mult), schemeCategory.FreeCoverLimit))
		}
		addItem("GLA", enabled, mult, sa)
	}
	{
		enabled := m.Benefits.PtdEnabled
		mult := firstNonZero(m.Benefits.PtdMultiple, m.Benefits.PtdMultiple)
		var sa *float64
		if enabled && mult != nil {
			sa = nz(math.Min(annualSalary*(*mult), schemeCategory.FreeCoverLimit))
		}
		addItem("PTD", enabled, mult, sa)
	}
	{
		enabled := m.Benefits.CiEnabled
		mult := firstNonZero(m.Benefits.CiMultiple, m.Benefits.CiMultiple)
		var sa *float64
		if enabled && mult != nil {
			sa = nz(math.Min(math.Min(annualSalary*(*mult), restriction.SevereIllnessMaximumBenefit), schemeCategory.FreeCoverLimit))
		}
		addItem("CI", enabled, mult, sa)
	}
	{
		enabled := m.Benefits.SglaEnabled
		mult := firstNonZero(m.Benefits.SglaMultiple, m.Benefits.SglaMultiple)
		var sa *float64
		if enabled && mult != nil {
			sa = nz(math.Min(math.Min(annualSalary*(*mult), restriction.SpouseGlaMaximumBenefit), schemeCategory.FreeCoverLimit))
		}
		addItem("SGLA", enabled, mult, sa)
	}
	{
		enabled := m.Benefits.TtdEnabled
		mult := firstNonZero(m.Benefits.TtdMultiple)
		var sa *float64
		if enabled && mult != nil {
			sa = nz(math.Min((annualSalary/12.0)*(*mult), restriction.TtdMaximumMonthlyBenefit))
		} // monthly
		addItem("TTD", enabled, mult, sa)
	}
	{
		enabled := m.Benefits.PhiEnabled
		mult := firstNonZero(m.Benefits.PhiMultiple)
		var sa *float64
		if enabled && mult != nil {
			sa = nz(math.Min((annualSalary/12.0)*(*mult), restriction.PhiMaximumMonthlyBenefit))
		} // monthly
		addItem("PHI", enabled, mult, sa)
	}

	// GFF breakdown from Bordereaux if available
	if m.Benefits.GffEnabled {
		// Initialize GFF so the client always sees the dedicated object when GFF is enabled
		gff := &GffBreakdownDTO{}

		if err == nil {
			gff.Currency = quote.Currency
			// Use ptr to include zero values explicitly if present in data
			gff.MainMember = ptr(schemeCategory.FamilyFuneralMainMemberFuneralSumAssured) //b.MainMemberFuneralSumAssured)
			gff.Spouse = ptr(schemeCategory.FamilyFuneralSpouseFuneralSumAssured)         //b.SpouseFuneralSumAssured)
			gff.Children = ptr(schemeCategory.FamilyFuneralChildrenFuneralSumAssured)     //b.ChildFuneralSumAssured)
			gff.Parents = ptr(schemeCategory.FamilyFuneralParentFuneralSumAssured)        //b.ParentFuneralSumAssured)
			gff.Dependants = ptr(schemeCategory.FamilyFuneralAdultDependantSumAssured)    //b.DependantFuneralSumAssured)
			// Populate counts where available from SchemeCategory
			if schemeCategory.FamilyFuneralMaxNumberChildren > 0 {
				cc := schemeCategory.FamilyFuneralMaxNumberChildren
				gff.ChildrenCount = &cc
			}
			if schemeCategory.FamilyFuneralMaxNumberAdultDependants > 0 {
				dc := schemeCategory.FamilyFuneralMaxNumberAdultDependants
				gff.DependantsCount = &dc
			}
			// ParentsCount not defined in SchemeCategory; left nil unless added later
		} else {
			// Log for observability; keep gff present even if amounts are unknown
			appLog.WithFields(map[string]interface{}{
				"member_id":   m.ID,
				"scheme_id":   m.SchemeId,
				"member_name": m.MemberName,
				"category":    m.SchemeCategory,
			}).Info("GFF enabled but no matching Bordereaux row found for breakdown")
		}
		dto.GFF = gff
	}

	// If no non-GFF benefits and no GFF data at all, fallback to quote context
	if len(dto.Benefits) == 0 && dto.GFF == nil {
		qdto, err := GetMemberBenefitSummaryQuote(m.QuoteId, m.ID)
		if err == nil {
			return qdto, nil
		}
	}

	return dto, nil
}

// GetMemberBenefitSummaryQuote returns the benefit summary for a member in the context of a quote.
// Uses MemberRatingResult for financials and GPricingMemberData for flags/multiples.
func GetMemberBenefitSummaryQuote(quoteID, memberID int) (MemberBenefitSummaryDTO, error) {
	var dto MemberBenefitSummaryDTO

	var m models.GPricingMemberData
	if err := DB.Where("id = ? AND quote_id = ?", memberID, quoteID).First(&m).Error; err != nil {
		return dto, err
	}

	// Populate common header fields
	dto.MemberID = m.ID
	dto.SchemeID = m.SchemeId
	dto.Source = "quote"
	dto.MemberIdNumber = m.MemberIdNumber
	dto.MemberName = m.MemberName
	dto.SchemeCategory = m.SchemeCategory
	dto.AnnualSalary = m.AnnualSalary
	dto.SchemeName = m.SchemeName

	// Benefit name mapper
	nameFor := func(code string) string {
		maps, err := GetBenefitMaps()
		if err != nil || len(maps) == 0 {
			return code
		}
		for _, mm := range maps {
			if strings.EqualFold(mm.BenefitCode, code) || strings.EqualFold(mm.BenefitAliasCode, code) {
				if strings.TrimSpace(mm.BenefitAlias) != "" {
					return mm.BenefitAlias
				}
				if strings.TrimSpace(mm.BenefitName) != "" {
					return mm.BenefitName
				}
				return code
			}
		}
		return code
	}

	// helpers
	firstNonZero := func(vals ...float64) *float64 {
		for _, v := range vals {
			if v > 0 {
				vv := v
				return &vv
			}
		}
		return nil
	}
	nz := func(v float64) *float64 {
		if v > 0 {
			vv := v
			return &vv
		}
		return nil
	}

	addItem := func(code string, enabled bool, multiple *float64, covered *float64) {
		item := BenefitMultipleItemDTO{Code: code, Name: nameFor(code), IsActive: enabled}
		if multiple != nil {
			item.SalaryMultiple = multiple
		}
		if covered != nil {
			item.CoveredSumAssured = covered
		}
		dto.Benefits = append(dto.Benefits, item)
	}

	annualSalary := m.AnnualSalary

	// Non-GFF benefits from multiples (always included with is_active flag)
	{
		enabled := m.Benefits.GlaEnabled
		mult := firstNonZero(m.Benefits.GlaMultiple, 0)
		var sa *float64
		if enabled && mult != nil {
			sa = nz(annualSalary * (*mult))
		}
		addItem("GLA", enabled, mult, sa)
	}
	{
		enabled := m.Benefits.PtdEnabled
		mult := firstNonZero(m.Benefits.PtdMultiple, 0)
		var sa *float64
		if enabled && mult != nil {
			sa = nz(annualSalary * (*mult))
		}
		addItem("PTD", enabled, mult, sa)
	}
	{
		enabled := m.Benefits.CiEnabled
		mult := firstNonZero(m.Benefits.CiMultiple, 0)
		var sa *float64
		if enabled && mult != nil {
			sa = nz(annualSalary * (*mult))
		}
		addItem("CI", enabled, mult, sa)
	}
	{
		enabled := m.Benefits.SglaEnabled
		mult := firstNonZero(m.Benefits.SglaMultiple, 0)
		var sa *float64
		if enabled && mult != nil {
			sa = nz(annualSalary * (*mult))
		}
		addItem("SGLA", enabled, mult, sa)
	}
	{
		enabled := m.Benefits.TtdEnabled
		mult := firstNonZero(m.Benefits.TtdMultiple)
		var sa *float64
		if enabled && mult != nil {
			sa = nz((annualSalary / 12.0) * (*mult))
		}
		addItem("TTD", enabled, mult, sa)
	}
	{
		enabled := m.Benefits.PhiEnabled
		mult := firstNonZero(m.Benefits.PhiMultiple)
		var sa *float64
		if enabled && mult != nil {
			sa = nz((annualSalary / 12.0) * (*mult))
		}
		addItem("PHI", enabled, mult, sa)
	}

	// GFF breakdown and counts
	if m.Benefits.GffEnabled {
		gff := &GffBreakdownDTO{}

		// Try populate amounts from MemberPremiumSchedule if available
		var ps models.MemberPremiumSchedule
		if err := DB.Where("quote_id = ? AND member_name = ? AND category = ?", m.QuoteId, m.MemberName, m.SchemeCategory).First(&ps).Error; err == nil {
			gff.MainMember = nz(ps.MainMemberFuneralSumAssured)
			gff.Spouse = nz(ps.SpouseFuneralSumAssured)
			gff.Children = nz(ps.ChildFuneralSumAssured)
			gff.Dependants = nz(ps.DependantsFuneralSumAssured)
		}

		// Fetch SchemeCategory to get counts and default amounts if needed
		var schemeCategory models.SchemeCategory
		if err := DB.Where("quote_id = ? AND scheme_category = ?", m.QuoteId, m.SchemeCategory).First(&schemeCategory).Error; err == nil {
			if schemeCategory.FamilyFuneralMainMemberFuneralSumAssured > 0 && gff.MainMember == nil {
				v := schemeCategory.FamilyFuneralMainMemberFuneralSumAssured
				gff.MainMember = &v
			}
			if schemeCategory.FamilyFuneralSpouseFuneralSumAssured > 0 && gff.Spouse == nil {
				v := schemeCategory.FamilyFuneralSpouseFuneralSumAssured
				gff.Spouse = &v
			}
			if schemeCategory.FamilyFuneralChildrenFuneralSumAssured > 0 && gff.Children == nil {
				v := schemeCategory.FamilyFuneralChildrenFuneralSumAssured
				gff.Children = &v
			}
			if schemeCategory.FamilyFuneralAdultDependantSumAssured > 0 && gff.Dependants == nil {
				v := schemeCategory.FamilyFuneralAdultDependantSumAssured
				gff.Dependants = &v
			}
			if schemeCategory.FamilyFuneralParentFuneralSumAssured > 0 && gff.Parents == nil {
				v := schemeCategory.FamilyFuneralParentFuneralSumAssured
				gff.Parents = &v
			}
			if schemeCategory.FamilyFuneralMaxNumberChildren > 0 {
				cc := schemeCategory.FamilyFuneralMaxNumberChildren
				gff.ChildrenCount = &cc
			}
			if schemeCategory.FamilyFuneralMaxNumberAdultDependants > 0 {
				dc := schemeCategory.FamilyFuneralMaxNumberAdultDependants
				gff.DependantsCount = &dc
			}
			// No explicit parents count in model; left nil by design
		}

		dto.GFF = gff
	}

	return dto, nil
}

// Create a beneficiary
func CreateBeneficiary(b models.Beneficiary, user models.AppUser) (models.Beneficiary, error) {
	// basic sanity
	if strings.TrimSpace(b.FullName) == "" {
		return b, fmt.Errorf("full_name is required")
	}
	if b.MemberID == 0 {
		return b, fmt.Errorf("memberId is required")
	}

	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := validateAllocation(tx, b.MemberID, nil, b.AllocationPercentage); err != nil {
			return err
		}
		if err := tx.Create(&b).Error; err != nil {
			return err
		}
		// Audit CREATE
		if err := writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "member_beneficiaries",
			EntityID:  strconv.Itoa(b.ID),
			Action:    "CREATE",
			ChangedBy: user.UserName,
		}, map[string]interface{}{}, b); err != nil {
			return err
		}

		// Log structured beneficiary activity
		beneficiaryDetails, _ := json.Marshal(map[string]interface{}{
			"beneficiaryName": b.FullName,
			"relationship":    b.Relationship,
			"action":          "Added",
			"allocation":      b.AllocationPercentage,
		})
		_ = tx.Create(&models.MemberActivity{
			MemberID:    b.MemberID,
			Type:        "beneficiary_change",
			Title:       "Beneficiary Added",
			Description: "New beneficiary added to policy",
			Details:     beneficiaryDetails,
			PerformedBy: user.UserName,
		})

		return nil
	})
	return b, err
}

// Update a beneficiary
func UpdateBeneficiary(memberID int, id int, input models.Beneficiary, user models.AppUser) (models.Beneficiary, error) {
	var b models.Beneficiary
	if err := DB.Where("id = ? AND member_id = ?", id, memberID).First(&b).Error; err != nil {
		return b, err
	}

	// keep ids consistent
	input.ID = b.ID
	input.MemberID = b.MemberID

	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := validateAllocation(tx, memberID, &id, input.AllocationPercentage); err != nil {
			return err
		}
		if err := tx.Model(&models.Beneficiary{}).Where("id = ? AND member_id = ?", id, memberID).Updates(map[string]interface{}{
			"full_name":             input.FullName,
			"relationship":          input.Relationship,
			"id_type":               input.IDType,
			"id_number":             input.IDNumber,
			"gender":                input.Gender,
			"date_of_birth":         input.DateOfBirth,
			"contact_number":        input.ContactNumber,
			"email":                 input.Email,
			"address":               input.Address,
			"allocation_percentage": input.AllocationPercentage,
			"benefit_types":         input.BenefitTypes,
			"guardian_name":         input.GuardianName,
			"guardian_relationship": input.GuardianRelationship,
			"guardian_id_number":    input.GuardianIDNumber,
			"guardian_contact":      input.GuardianContact,
			"bank_name":             input.BankName,
			"branch_code":           input.BranchCode,
			"account_number":        input.AccountNumber,
			"account_type":          input.AccountType,
			"status":                input.Status,
		}).Error; err != nil {
			return err
		}
		var after models.Beneficiary
		if err := tx.Where("id = ?", id).First(&after).Error; err != nil {
			return err
		}
		// Audit UPDATE
		if err := writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "member_beneficiaries",
			EntityID:  strconv.Itoa(after.ID),
			Action:    "UPDATE",
			ChangedBy: user.UserName,
		}, b, after); err != nil {
			return err
		}

		// Log structured beneficiary activity
		beneficiaryDetails, _ := json.Marshal(map[string]interface{}{
			"beneficiaryName": after.FullName,
			"relationship":    after.Relationship,
			"action":          "Updated",
			"allocation":      after.AllocationPercentage,
		})
		_ = tx.Create(&models.MemberActivity{
			MemberID:    after.MemberID,
			Type:        "beneficiary_change",
			Title:       "Beneficiary Updated",
			Description: "Beneficiary allocation percentage updated",
			Details:     beneficiaryDetails,
			PerformedBy: user.UserName,
		})

		return nil
	})
	if err != nil {
		return b, err
	}
	// reload
	if err := DB.Where("id = ?", id).First(&b).Error; err != nil {
		return b, err
	}
	return b, nil
}

// Delete a beneficiary
func DeleteBeneficiary(memberID int, id int, user models.AppUser) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND member_id = ?", id, memberID).Delete(&models.Beneficiary{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// CreateBenefitDocumentType creates a new benefit document type
func CreateBenefitDocumentType(docType models.BenefitDocumentType) (models.BenefitDocumentType, error) {
	if err := DB.Create(&docType).Error; err != nil {
		return docType, err
	}
	return docType, nil
}

// GetBenefitDocumentTypes returns all benefit document types
func GetBenefitDocumentTypes() ([]models.BenefitDocumentType, error) {
	var docTypes []models.BenefitDocumentType
	if err := DB.Find(&docTypes).Error; err != nil {
		return nil, err
	}
	return docTypes, nil
}

// GetBenefitDocumentTypesByBenefitCode returns benefit document types for a specific benefit code
func GetBenefitDocumentTypesByBenefitCode(benefitCode string) ([]models.BenefitDocumentType, error) {
	var docTypes []models.BenefitDocumentType
	if err := DB.Where("benefit_code = ?", benefitCode).Find(&docTypes).Error; err != nil {
		return nil, err
	}
	return docTypes, nil
}

// GetBenefitDocumentTypesByClaimID returns benefit document types for a specific claim
func GetBenefitDocumentTypesByClaimID(claimID int) ([]models.BenefitDocumentType, error) {
	var claim models.GroupSchemeClaim
	if err := DB.First(&claim, claimID).Error; err != nil {
		return nil, err
	}
	return GetBenefitDocumentTypesByBenefitCode(claim.BenefitCode)
}

// UpdateBenefitDocumentType updates an existing benefit document type
func UpdateBenefitDocumentType(id int, docType models.BenefitDocumentType) (models.BenefitDocumentType, error) {
	docType.ID = id
	if err := DB.Save(&docType).Error; err != nil {
		return docType, err
	}
	return docType, nil
}

// DeleteBenefitDocumentType deletes a benefit document type
func DeleteBenefitDocumentType(id int) error {
	if err := DB.Delete(&models.BenefitDocumentType{}, id).Error; err != nil {
		return err
	}
	return nil
}

// CreateSchemeCategoryMaster creates a new scheme category master record
func CreateSchemeCategoryMaster(category models.SchemeCategoryMaster, user models.AppUser) (models.SchemeCategoryMaster, error) {
	category.CreatedBy = user.UserName
	if err := DB.Create(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

// GetSchemeCategoryMasters returns all scheme category master records, optionally filtered by insurer_id
func GetSchemeCategoryMasters(insurerId int) ([]models.SchemeCategoryMaster, error) {
	var categories []models.SchemeCategoryMaster
	query := DB
	if insurerId > 0 {
		query = query.Where("insurer_id = ? OR insurer_id = 0", insurerId)
	}
	// Only return active categories for normal usage, or let frontend decide?
	// The request says "These fields are expected in the front end application".
	// Usually, frontend handles display of active/inactive, but maybe we should filter by default?
	// Looking at GetGroupSchemeCategories, it calls this.
	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetSchemeCategoryMasterByID returns a single scheme category master record by ID
func GetSchemeCategoryMasterByID(id int) (models.SchemeCategoryMaster, error) {
	var category models.SchemeCategoryMaster
	if err := DB.First(&category, id).Error; err != nil {
		return category, err
	}
	return category, nil
}

// UpdateSchemeCategoryMaster updates an existing scheme category master record
func UpdateSchemeCategoryMaster(id int, category models.SchemeCategoryMaster, user models.AppUser) (models.SchemeCategoryMaster, error) {
	var existing models.SchemeCategoryMaster
	if err := DB.First(&existing, id).Error; err != nil {
		return existing, err
	}
	existing.Name = category.Name
	existing.Description = category.Description
	existing.Active = category.Active
	existing.InsurerId = category.InsurerId
	// UpdatedAt is handled by GORM autoUpdateTime
	if err := DB.Save(&existing).Error; err != nil {
		return existing, err
	}
	return existing, nil
}

// DeleteSchemeCategoryMaster deletes a scheme category master record
func DeleteSchemeCategoryMaster(id int) error {
	if err := DB.Delete(&models.SchemeCategoryMaster{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetDistinctRegionsForRiskCode returns distinct region names from region_loadings
// for the given risk_rate_code, deduplicated and sorted alphabetically.
func GetDistinctRegionsForRiskCode(riskRateCode string) ([]string, error) {
	var regions []string
	err := DB.Model(&models.RegionLoading{}).
		Where("risk_rate_code = ?", riskRateCode).
		Distinct("region").
		Order("region asc").
		Pluck("region", &regions).Error
	return regions, err
}

// schemeLoadingFromQuote returns the per-member office-premium loading
// fraction (ExpenseLoading + ProfitLoading) used to derive the *pre-commission*
// office premium. CommissionLoading is excluded — commission is added on top
// of the pre-comm office premium via the progressive commission allocation,
// not baked into the gross-up denominator. The quote stores loadings as
// percentages (e.g. 5 for 5%) so we divide by 100.
func schemeLoadingFromQuote(quote *models.GroupPricingQuote) float64 {
	if quote == nil {
		return 0
	}
	binderRate, outsourceRate := binderAndOutsourceRates(quote)
	return (quote.Loadings.ExpenseLoading+quote.Loadings.ProfitLoading+
		quote.Loadings.AdminLoading+quote.Loadings.OtherLoading)/100.0 +
		binderRate + outsourceRate
}

// computeMemberOfficePremium derives the *pre-commission* office premium for a
// per-member risk premium using the scheme-level loading on the quote.
// Commission is added on top via the progressive allocation, not via this
// denominator.
func computeMemberOfficePremium(riskPremium float64, quote *models.GroupPricingQuote) float64 {
	denom := 1.0 - schemeLoadingFromQuote(quote)
	if denom <= 0 {
		return 0
	}
	return riskPremium / denom
}

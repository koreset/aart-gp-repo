package services

import (
	appLog "api/log"
	"api/models"
	"api/utils"
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/jszwec/csvutil"
	"gorm.io/gorm"
)

// GetBrokerByName looks up a broker by name. Returns gorm.ErrRecordNotFound if not found.
func GetBrokerByName(name string) (models.Broker, error) {
	var broker models.Broker
	err := DB.Where("name = ?", name).First(&broker).Error
	return broker, err
}

// openCSVDecoder handles delimiter detection, gzip, and returns a csvutil decoder + headers.
func openCSVDecoder(fh *multipart.FileHeader) (*csvutil.Decoder, []string, error) {
	// Detect delimiter
	delimFile, err := fh.Open()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file for delimiter detection: %v", err)
	}
	defer delimFile.Close()
	delimiter, err := utils.GetDelimiter(delimFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to detect CSV delimiter: %v", err)
	}

	// Open file for reading
	file, err := fh.Open()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file for reading: %v", err)
	}

	var reader io.Reader = file
	magic := make([]byte, 2)
	if n, readErr := file.Read(magic); readErr == nil && n == 2 && magic[0] == 0x1f && magic[1] == 0x8b {
		if _, seekErr := file.Seek(0, io.SeekStart); seekErr != nil {
			return nil, nil, fmt.Errorf("failed to seek: %v", seekErr)
		}
		gzReader, gzErr := gzip.NewReader(file)
		if gzErr != nil {
			return nil, nil, fmt.Errorf("failed to create gzip reader: %v", gzErr)
		}
		reader = gzReader
	} else {
		if _, seekErr := file.Seek(0, io.SeekStart); seekErr != nil {
			return nil, nil, fmt.Errorf("failed to seek: %v", seekErr)
		}
	}

	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true
	csvReader.Comma = delimiter
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CSV decoder: %v", err)
	}
	return dec, dec.Header(), nil
}

// decodeAllRows decodes all CSV rows into a slice of the given type.
func decodeAllRows[T any](dec *csvutil.Decoder) ([]T, error) {
	var rows []T
	for i := 1; ; i++ {
		var row T
		if err := dec.Decode(&row); err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("error decoding row %d: %v", i, err)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// ValidateMigration parses all CSV templates, validates headers, cross-references, and field values.
// It does NOT write to the database.
func ValidateMigration(files map[string]*multipart.FileHeader, user models.AppUser) (*models.MigrationResult, error) {
	result := &models.MigrationResult{Valid: true}

	// --- Parse Template 1: Scheme Setup (required) ---
	schemeFile, ok := files["scheme_setup"]
	if !ok || schemeFile == nil {
		return nil, fmt.Errorf("scheme_setup CSV file is required")
	}
	schemeDec, schemeHeaders, err := openCSVDecoder(schemeFile)
	if err != nil {
		return nil, fmt.Errorf("scheme_setup: %v", err)
	}
	if validationErr := utils.ValidateCSVHeaders(schemeHeaders, models.SchemeSetupRow{}); validationErr != nil {
		return nil, fmt.Errorf("scheme_setup header validation failed: %v", validationErr)
	}
	schemeRows, err := decodeAllRows[models.SchemeSetupRow](schemeDec)
	if err != nil {
		return nil, fmt.Errorf("scheme_setup: %v", err)
	}
	if len(schemeRows) == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, models.MigrationValidationError{
			Template: "scheme_setup", Row: 0, Column: "", Message: "no data rows found",
		})
		return result, nil
	}

	// Build scheme name set and validate required fields
	schemeNames := make(map[string]bool)
	for i, row := range schemeRows {
		rowNum := i + 2 // 1-indexed + header row
		if strings.TrimSpace(row.SchemeName) == "" {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "scheme_setup", Row: rowNum, Column: "scheme_name", Message: "scheme_name is required",
			})
		} else if schemeNames[row.SchemeName] {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "scheme_setup", Row: rowNum, Column: "scheme_name", Message: fmt.Sprintf("duplicate scheme_name '%s'", row.SchemeName),
			})
		} else {
			schemeNames[row.SchemeName] = true
		}
		if strings.TrimSpace(row.DistributionChannel) == "" {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "scheme_setup", Row: rowNum, Column: "distribution_channel", Message: "distribution_channel is required",
			})
		}
		channel := strings.ToLower(strings.TrimSpace(row.DistributionChannel))
		if channel != "" && channel != "broker" && channel != "direct" && channel != "binder" && channel != "tied_agent" {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "scheme_setup", Row: rowNum, Column: "distribution_channel", Message: fmt.Sprintf("invalid distribution_channel '%s'", row.DistributionChannel),
			})
		}
		if channel == "broker" && strings.TrimSpace(row.BrokerName) == "" {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "scheme_setup", Row: rowNum, Column: "broker_name", Message: "broker_name is required when distribution_channel is 'broker'",
			})
		}
		if time.Time(row.CommencementDate).IsZero() {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "scheme_setup", Row: rowNum, Column: "commencement_date", Message: "commencement_date is required",
			})
		}
		if time.Time(row.CoverStartDate).IsZero() {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "scheme_setup", Row: rowNum, Column: "cover_start_date", Message: "cover_start_date is required",
			})
		}
		if time.Time(row.CoverEndDate).IsZero() {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "scheme_setup", Row: rowNum, Column: "cover_end_date", Message: "cover_end_date is required",
			})
		}
		if strings.TrimSpace(row.Currency) == "" {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "scheme_setup", Row: rowNum, Column: "currency", Message: "currency is required",
			})
		}
	}
	result.SchemeCount = len(schemeRows)

	// --- Parse Template 2: Categories (required) ---
	catFile, ok := files["categories"]
	if !ok || catFile == nil {
		return nil, fmt.Errorf("categories CSV file is required")
	}
	catDec, catHeaders, err := openCSVDecoder(catFile)
	if err != nil {
		return nil, fmt.Errorf("categories: %v", err)
	}
	if validationErr := utils.ValidateCSVHeaders(catHeaders, models.SchemeCategoryRow{}); validationErr != nil {
		return nil, fmt.Errorf("categories header validation failed: %v", validationErr)
	}
	catRows, err := decodeAllRows[models.SchemeCategoryRow](catDec)
	if err != nil {
		return nil, fmt.Errorf("categories: %v", err)
	}
	if len(catRows) == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, models.MigrationValidationError{
			Template: "categories", Row: 0, Column: "", Message: "no data rows found",
		})
	}

	// Cross-reference: category scheme_name must exist in scheme setup
	catKey := make(map[string]map[string]bool) // scheme_name -> set of category_names
	for i, row := range catRows {
		rowNum := i + 2
		sn := strings.TrimSpace(row.SchemeName)
		cn := strings.TrimSpace(row.CategoryName)
		if sn == "" {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "categories", Row: rowNum, Column: "scheme_name", Message: "scheme_name is required",
			})
			continue
		}
		if !schemeNames[sn] {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "categories", Row: rowNum, Column: "scheme_name", Message: fmt.Sprintf("scheme_name '%s' not found in scheme_setup", sn),
			})
		}
		if cn == "" {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "categories", Row: rowNum, Column: "category_name", Message: "category_name is required",
			})
			continue
		}
		if catKey[sn] == nil {
			catKey[sn] = make(map[string]bool)
		}
		if catKey[sn][cn] {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "categories", Row: rowNum, Column: "category_name", Message: fmt.Sprintf("duplicate category '%s' for scheme '%s'", cn, sn),
			})
		}
		catKey[sn][cn] = true
	}
	result.CategoryCount = len(catRows)

	// --- Parse Template 3: Member Data (required) ---
	memberFile, ok := files["member_data"]
	if !ok || memberFile == nil {
		return nil, fmt.Errorf("member_data CSV file is required")
	}
	memberDec, memberHeaders, err := openCSVDecoder(memberFile)
	if err != nil {
		return nil, fmt.Errorf("member_data: %v", err)
	}
	if validationErr := utils.ValidateCSVHeaders(memberHeaders, models.MemberDataRow{}); validationErr != nil {
		return nil, fmt.Errorf("member_data header validation failed: %v", validationErr)
	}
	memberRows, err := decodeAllRows[models.MemberDataRow](memberDec)
	if err != nil {
		return nil, fmt.Errorf("member_data: %v", err)
	}

	// Collect RSA IDs for bulk validation via CheckID API
	var rsaIDsForValidation []string
	for _, row := range memberRows {
		idType := strings.ToUpper(strings.TrimSpace(row.MemberIdType))
		idNum := strings.TrimSpace(row.MemberIdNumber)
		if (idType == "RSA_ID" || idType == "ID") && idNum != "" {
			rsaIDsForValidation = append(rsaIDsForValidation, idNum)
		}
	}
	rsaIDResults, rsaIDErr := utils.ValidateRSAIDsBulk(rsaIDsForValidation)
	if rsaIDErr != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.MigrationValidationError{
			Template: "member_data", Row: 0, Column: "member_id_number", Message: fmt.Sprintf("ID validation service error: %v", rsaIDErr),
		})
		return result, nil
	}

	// Validate member data: cross-references and field values
	memberIdSet := make(map[string]map[string]bool) // scheme_name -> set of member_id_numbers
	for i, row := range memberRows {
		rowNum := i + 2
		sn := strings.TrimSpace(row.SchemeName)
		cat := strings.TrimSpace(row.SchemeCategory)
		idNum := strings.TrimSpace(row.MemberIdNumber)

		if sn != "" && !schemeNames[sn] {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "member_data", Row: rowNum, Column: "scheme_name", Message: fmt.Sprintf("scheme_name '%s' not found in scheme_setup", sn),
			})
		}
		if cat != "" && catKey[sn] != nil && !catKey[sn][cat] {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "member_data", Row: rowNum, Column: "scheme_category", Message: fmt.Sprintf("scheme_category '%s' not found in categories for scheme '%s'", cat, sn),
			})
		}
		if strings.TrimSpace(row.MemberName) == "" {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "member_data", Row: rowNum, Column: "member_name", Message: "member_name is required",
			})
		}
		if idNum == "" {
			result.Valid = false
			result.Errors = append(result.Errors, models.MigrationValidationError{
				Template: "member_data", Row: rowNum, Column: "member_id_number", Message: "member_id_number is required",
			})
		} else {
			idType := strings.ToUpper(strings.TrimSpace(row.MemberIdType))
			if (idType == "RSA_ID" || idType == "ID") {
				if valid, ok := rsaIDResults[idNum]; ok && !valid {
					result.Valid = false
					result.Errors = append(result.Errors, models.MigrationValidationError{
						Template: "member_data", Row: rowNum, Column: "member_id_number", Message: fmt.Sprintf("invalid RSA ID '%s'", idNum),
					})
				}
			}
			if memberIdSet[sn] == nil {
				memberIdSet[sn] = make(map[string]bool)
			}
			if memberIdSet[sn][idNum] {
				result.Valid = false
				result.Errors = append(result.Errors, models.MigrationValidationError{
					Template: "member_data", Row: rowNum, Column: "member_id_number", Message: fmt.Sprintf("duplicate member_id_number '%s' in scheme '%s'", idNum, sn),
				})
			}
			memberIdSet[sn][idNum] = true
		}
	}
	result.MemberCount = len(memberRows)

	// --- Parse Template 4: Beneficiaries (optional) ---
	if benFile, hasBen := files["beneficiaries"]; hasBen && benFile != nil {
		benDec, benHeaders, benErr := openCSVDecoder(benFile)
		if benErr != nil {
			return nil, fmt.Errorf("beneficiaries: %v", benErr)
		}
		if validationErr := utils.ValidateCSVHeaders(benHeaders, models.BeneficiaryRow{}); validationErr != nil {
			return nil, fmt.Errorf("beneficiaries header validation failed: %v", validationErr)
		}
		benRows, decErr := decodeAllRows[models.BeneficiaryRow](benDec)
		if decErr != nil {
			return nil, fmt.Errorf("beneficiaries: %v", decErr)
		}

		// Track allocation totals per (scheme_name, member_id_number)
		allocations := make(map[string]float64)
		for i, row := range benRows {
			rowNum := i + 2
			sn := strings.TrimSpace(row.SchemeName)
			memId := strings.TrimSpace(row.MemberIdNumber)

			if sn != "" && !schemeNames[sn] {
				result.Valid = false
				result.Errors = append(result.Errors, models.MigrationValidationError{
					Template: "beneficiaries", Row: rowNum, Column: "scheme_name", Message: fmt.Sprintf("scheme_name '%s' not found in scheme_setup", sn),
				})
			}
			if memId != "" && memberIdSet[sn] != nil && !memberIdSet[sn][memId] {
				result.Valid = false
				result.Errors = append(result.Errors, models.MigrationValidationError{
					Template: "beneficiaries", Row: rowNum, Column: "member_id_number", Message: fmt.Sprintf("member_id_number '%s' not found in member_data for scheme '%s'", memId, sn),
				})
			}
			if strings.TrimSpace(row.BeneficiaryFullName) == "" {
				result.Valid = false
				result.Errors = append(result.Errors, models.MigrationValidationError{
					Template: "beneficiaries", Row: rowNum, Column: "beneficiary_full_name", Message: "beneficiary_full_name is required",
				})
			}
			if row.AllocationPercentage <= 0 || row.AllocationPercentage > 100 {
				result.Valid = false
				result.Errors = append(result.Errors, models.MigrationValidationError{
					Template: "beneficiaries", Row: rowNum, Column: "allocation_percentage", Message: "allocation_percentage must be between 0 and 100",
				})
			}

			// Check minor guardian requirements
			dob := time.Time(row.DateOfBirth)
			if !dob.IsZero() {
				age := time.Now().Year() - dob.Year()
				if age < 18 && strings.TrimSpace(row.GuardianName) == "" {
					result.Valid = false
					result.Errors = append(result.Errors, models.MigrationValidationError{
						Template: "beneficiaries", Row: rowNum, Column: "guardian_name", Message: "guardian_name is required for minor beneficiaries (under 18)",
					})
				}
			}

			key := sn + "|" + memId
			allocations[key] += row.AllocationPercentage
		}

		// Check allocation sums
		for key, total := range allocations {
			if total > 100.01 { // small tolerance for floating point
				parts := strings.SplitN(key, "|", 2)
				result.Valid = false
				result.Errors = append(result.Errors, models.MigrationValidationError{
					Template: "beneficiaries", Row: 0, Column: "allocation_percentage",
					Message: fmt.Sprintf("total allocation for member '%s' in scheme '%s' is %.1f%% (exceeds 100%%)", parts[1], parts[0], total),
				})
			}
		}
		result.BeneficiaryCount = len(benRows)
	}

	// --- Parse Template 5: Claims Experience (optional) ---
	if claimsFile, hasClaims := files["claims_experience"]; hasClaims && claimsFile != nil {
		claimsDec, claimsHeaders, claimsErr := openCSVDecoder(claimsFile)
		if claimsErr != nil {
			return nil, fmt.Errorf("claims_experience: %v", claimsErr)
		}
		if validationErr := utils.ValidateCSVHeaders(claimsHeaders, models.GroupPricingClaimsExperience{}); validationErr != nil {
			return nil, fmt.Errorf("claims_experience header validation failed: %v", validationErr)
		}
		claimsRows, decErr := decodeAllRows[models.GroupPricingClaimsExperience](claimsDec)
		if decErr != nil {
			return nil, fmt.Errorf("claims_experience: %v", decErr)
		}
		result.ClaimsCount = len(claimsRows)
	}

	return result, nil
}

// ExecuteMigration validates then atomically creates all entities in a single DB transaction.
func ExecuteMigration(files map[string]*multipart.FileHeader, user models.AppUser) (*models.MigrationResult, error) {
	logger := appLog.WithField("action", "ExecuteMigration").WithField("user", user.UserName)

	// Validate first
	valResult, err := ValidateMigration(files, user)
	if err != nil {
		return nil, err
	}
	if !valResult.Valid {
		return valResult, nil
	}

	// Re-parse all templates for execution (validation consumed the readers)
	schemeRows, err := parseCSVFile[models.SchemeSetupRow](files["scheme_setup"])
	if err != nil {
		return nil, fmt.Errorf("scheme_setup re-parse: %v", err)
	}
	catRows, err := parseCSVFile[models.SchemeCategoryRow](files["categories"])
	if err != nil {
		return nil, fmt.Errorf("categories re-parse: %v", err)
	}
	memberRows, err := parseCSVFile[models.MemberDataRow](files["member_data"])
	if err != nil {
		return nil, fmt.Errorf("member_data re-parse: %v", err)
	}

	var benRows []models.BeneficiaryRow
	if benFile, ok := files["beneficiaries"]; ok && benFile != nil {
		benRows, err = parseCSVFile[models.BeneficiaryRow](benFile)
		if err != nil {
			return nil, fmt.Errorf("beneficiaries re-parse: %v", err)
		}
	}

	var claimsRows []models.GroupPricingClaimsExperience
	if claimsFile, ok := files["claims_experience"]; ok && claimsFile != nil {
		claimsRows, err = parseCSVFile[models.GroupPricingClaimsExperience](claimsFile)
		if err != nil {
			return nil, fmt.Errorf("claims_experience re-parse: %v", err)
		}
	}

	result := &models.MigrationResult{
		Valid:       true,
		SchemeCount: len(schemeRows),
	}

	txErr := DB.Transaction(func(tx *gorm.DB) error {
		for _, schemeRow := range schemeRows {
			schemeName := strings.TrimSpace(schemeRow.SchemeName)
			logger.Infof("Migrating scheme: %s", schemeName)

			// --- 1. Resolve broker ---
			var brokerId int
			channel := models.DistributionChannel(strings.ToLower(strings.TrimSpace(schemeRow.DistributionChannel)))
			if channel != models.ChannelDirect {
				brokerName := strings.TrimSpace(schemeRow.BrokerName)
				var broker models.Broker
				err := tx.Where("name = ?", brokerName).First(&broker).Error
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						broker = models.Broker{
							Name:         brokerName,
							ContactEmail: strings.TrimSpace(schemeRow.BrokerEmail),
							CreatedBy:    user.UserName,
						}
						if createErr := tx.Create(&broker).Error; createErr != nil {
							return fmt.Errorf("failed to create broker '%s': %v", brokerName, createErr)
						}
						logger.Infof("Created new broker: %s (ID %d)", brokerName, broker.ID)
					} else {
						return fmt.Errorf("broker lookup failed: %v", err)
					}
				}
				brokerId = broker.ID
			}

			// --- 2. Create GroupPricingQuote ---
			commDate := time.Time(schemeRow.CommencementDate)
			coverStart := time.Time(schemeRow.CoverStartDate)
			coverEnd := time.Time(schemeRow.CoverEndDate)

			// Build selected categories list for this scheme
			var selectedCats models.StringArray
			for _, cat := range catRows {
				if strings.TrimSpace(cat.SchemeName) == schemeName {
					selectedCats = append(selectedCats, strings.TrimSpace(cat.CategoryName))
				}
			}

			quoteName := strings.ReplaceAll(schemeName, " ", "_") + "_migration_1"
			quote := models.GroupPricingQuote{
				QuoteName:                quoteName,
				QuoteType:                "New Business",
				SchemeName:               schemeName,
				SchemeContact:            strings.TrimSpace(schemeRow.ContactPerson),
				SchemeEmail:              strings.TrimSpace(schemeRow.ContactEmail),
				DistributionChannel:      channel,
				CommencementDate:         commDate,
				CoverEndDate:             coverEnd,
				Industry:                 strings.TrimSpace(schemeRow.Industry),
				OccupationClass:          schemeRow.OccupationClass,
				FreeCoverLimit:           schemeRow.FreeCoverLimit,
				Currency:                 strings.TrimSpace(schemeRow.Currency),
				NormalRetirementAge:      schemeRow.NormalRetirementAge,
				ExperienceRating:         strings.TrimSpace(schemeRow.ExperienceRating),
				RiskRateCode:             strings.TrimSpace(schemeRow.RiskRateCode),
				UseGlobalSalaryMultiple:  schemeRow.UseGlobalSalaryMultiple,
				CreatedBy:                user.UserName,
				CreationDate:             time.Now(),
				ModificationDate:         time.Now(),
				Status:                   models.StatusApproved,
				SchemeQuoteStatus:        models.StatusNotInEffect,
				SelectedSchemeCategories: selectedCats,
			}
			if channel != models.ChannelDirect {
				quote.QuoteBroker = models.QuoteBroker{ID: brokerId}
			}

			if err := tx.Create(&quote).Error; err != nil {
				return fmt.Errorf("failed to create quote for scheme '%s': %v", schemeName, err)
			}
			result.CreatedQuoteIDs = append(result.CreatedQuoteIDs, quote.ID)

			// --- 3. Create GroupScheme ---
			scheme := models.GroupScheme{
				Name:                schemeName,
				ContactPerson:       strings.TrimSpace(schemeRow.ContactPerson),
				ContactEmail:        strings.TrimSpace(schemeRow.ContactEmail),
				DistributionChannel: channel,
				BrokerId:            brokerId,
				CreationDate:        time.Now(),
				CreatedBy:           user.UserName,
				QuoteId:             quote.ID,
				Status:              models.StatusQuoted,
				QuoteInForce:        quoteName,
				RenewalDate:         time.Time(schemeRow.RenewalDate),
				CoverStartDate:      coverStart,
				CoverEndDate:        coverEnd,
				CommencementDate:    commDate,
			}
			if err := tx.Create(&scheme).Error; err != nil {
				return fmt.Errorf("failed to create scheme '%s': %v", schemeName, err)
			}
			result.CreatedSchemeIDs = append(result.CreatedSchemeIDs, scheme.ID)

			// Link quote to scheme
			quote.SchemeID = scheme.ID
			if err := tx.Save(&quote).Error; err != nil {
				return fmt.Errorf("failed to link quote to scheme '%s': %v", schemeName, err)
			}

			// --- 4. Create SchemeCategories ---
			var categories []models.SchemeCategory
			for _, catRow := range catRows {
				if strings.TrimSpace(catRow.SchemeName) != schemeName {
					continue
				}
				cat := models.SchemeCategory{
					QuoteId:                                  quote.ID,
					SchemeCategory:                           strings.TrimSpace(catRow.CategoryName),
					Basis:                                    strings.TrimSpace(catRow.Basis),
					FreeCoverLimit:                           catRow.FreeCoverLimit,
					Region:                                   strings.TrimSpace(catRow.Region),
					GlaBenefit:                               catRow.GlaBenefit,
					GlaSalaryMultiple:                        catRow.GlaSalaryMultiple,
					GlaBenefitType:                           strings.TrimSpace(catRow.GlaBenefitType),
					GlaTerminalIllnessBenefit:                strings.TrimSpace(catRow.GlaTerminalIllnessBenefit),
					GlaWaitingPeriod:                         catRow.GlaWaitingPeriod,
					GlaEducatorBenefit:                       strings.TrimSpace(catRow.GlaEducatorBenefit),
					SglaBenefit:                              catRow.SglaBenefit,
					SglaSalaryMultiple:                       catRow.SglaSalaryMultiple,
					SglaMaxBenefit:                           catRow.SglaMaxBenefit,
					PtdBenefit:                               catRow.PtdBenefit,
					PtdSalaryMultiple:                        catRow.PtdSalaryMultiple,
					PtdBenefitType:                           strings.TrimSpace(catRow.PtdBenefitType),
					PtdRiskType:                              strings.TrimSpace(catRow.PtdRiskType),
					PtdDeferredPeriod:                        catRow.PtdDeferredPeriod,
					PtdDisabilityDefinition:                  strings.TrimSpace(catRow.PtdDisabilityDefinition),
					CiBenefit:                                catRow.CiBenefit,
					CiCriticalIllnessSalaryMultiple:          catRow.CiSalaryMultiple,
					CiBenefitStructure:                       strings.TrimSpace(catRow.CiBenefitStructure),
					CiBenefitDefinition:                      strings.TrimSpace(catRow.CiBenefitDefinition),
					CiMaxBenefit:                             catRow.CiMaxBenefit,
					TtdBenefit:                               catRow.TtdBenefit,
					TtdRiskType:                              strings.TrimSpace(catRow.TtdRiskType),
					TtdMaximumBenefit:                        catRow.TtdMaximumBenefit,
					TtdIncomeReplacementPercentage:           catRow.TtdIncomeReplacementPercentage,
					TtdWaitingPeriod:                         catRow.TtdWaitingPeriod,
					TtdDeferredPeriod:                        catRow.TtdDeferredPeriod,
					PhiBenefit:                               catRow.PhiBenefit,
					PhiRiskType:                              strings.TrimSpace(catRow.PhiRiskType),
					PhiMaximumBenefit:                        catRow.PhiMaximumBenefit,
					PhiIncomeReplacementPercentage:           catRow.PhiIncomeReplacementPercentage,
					PhiWaitingPeriod:                         catRow.PhiWaitingPeriod,
					PhiDeferredPeriod:                        catRow.PhiDeferredPeriod,
					PhiNormalRetirementAge:                   catRow.PhiNormalRetirementAge,
					PhiPremiumWaiver:                         strings.TrimSpace(catRow.PhiPremiumWaiver),
					FamilyFuneralBenefit:                     catRow.FamilyFuneralBenefit,
					FamilyFuneralMainMemberFuneralSumAssured: catRow.FuneralMainMemberSumAssured,
					FamilyFuneralSpouseFuneralSumAssured:     catRow.FuneralSpouseSumAssured,
					FamilyFuneralChildrenFuneralSumAssured:   catRow.FuneralChildrenSumAssured,
					FamilyFuneralMaxNumberChildren:           catRow.FuneralMaxChildren,
				}
				categories = append(categories, cat)
			}
			if len(categories) > 0 {
				if err := tx.CreateInBatches(&categories, 100).Error; err != nil {
					return fmt.Errorf("failed to create categories for scheme '%s': %v", schemeName, err)
				}
			}
			result.CategoryCount += len(categories)

			// --- 5. Insert member data (Template 3) ---
			var membersData []models.GPricingMemberData
			for _, row := range memberRows {
				if strings.TrimSpace(row.SchemeName) != schemeName {
					continue
				}
				pp := models.GPricingMemberData{
					Year:           row.Year,
					SchemeName:     schemeName,
					MemberName:     strings.TrimSpace(row.MemberName),
					MemberIdNumber: strings.TrimSpace(row.MemberIdNumber),
					MemberIdType:   strings.TrimSpace(row.MemberIdType),
					SchemeCategory: strings.TrimSpace(row.SchemeCategory),
					Gender:         strings.TrimSpace(row.Gender),
					DateOfBirth:    time.Time(row.DateOfBirth),
					Email:          strings.TrimSpace(row.Email),
					EmployeeNumber: strings.TrimSpace(row.EmployeeNumber),
					AnnualSalary:   row.AnnualSalary,
					ContributionWaiverProportion: row.ContributionWaiverProportion,
					EntryDate: time.Time(row.EntryDate),
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
					QuoteId:   quote.ID,
					SchemeId:  scheme.ID,
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
			}
			if len(membersData) > 0 {
				if err := tx.CreateInBatches(&membersData, 100).Error; err != nil {
					return fmt.Errorf("failed to insert member data for scheme '%s': %v", schemeName, err)
				}
			}
			result.MemberCount += len(membersData)

			// --- 6. Insert claims experience (Template 5, optional) ---
			var schemeClaims int
			for _, clRow := range claimsRows {
				clRow.QuoteId = quote.ID
				clRow.CreatedBy = user.UserName
				if err := tx.Create(&clRow).Error; err != nil {
					return fmt.Errorf("failed to insert claims experience for scheme '%s': %v", schemeName, err)
				}
				schemeClaims++
			}
			result.ClaimsCount += schemeClaims

			// --- 7. Transition to in-force (mirrors AcceptGroupPricingQuote lines 390-495) ---
			// Delete any existing in-force data for this quote
			tx.Where("quote_id = ?", quote.ID).Delete(&models.GPricingMemberDataInForce{})

			var gmdif []models.GPricingMemberDataInForce
			err = tx.Model(&models.GPricingMemberData{}).Where("quote_id = ?", quote.ID).Find(&gmdif).Error
			if err != nil {
				return fmt.Errorf("failed to load member data for in-force copy: %v", err)
			}

			for i := range gmdif {
				gmdif[i].ID = 0
				gmdif[i].IsOriginalMember = true
				gmdif[i].Status = "Active"

				// Apply benefit mappings from categories
				for _, cat := range categories {
					if cat.SchemeCategory == gmdif[i].SchemeCategory {
						if quote.UseGlobalSalaryMultiple {
							if cat.GlaBenefit {
								gmdif[i].Benefits.GlaEnabled = true
								gmdif[i].Benefits.GlaMultiple = cat.GlaSalaryMultiple
							}
							if cat.SglaBenefit {
								gmdif[i].Benefits.SglaEnabled = true
								gmdif[i].Benefits.SglaMultiple = cat.SglaSalaryMultiple
							}
							if cat.PtdBenefit {
								gmdif[i].Benefits.PtdEnabled = true
								gmdif[i].Benefits.PtdMultiple = cat.PtdSalaryMultiple
							}
							if cat.CiBenefit {
								gmdif[i].Benefits.CiEnabled = true
								gmdif[i].Benefits.CiMultiple = cat.CiCriticalIllnessSalaryMultiple
							}
							if cat.PhiBenefit {
								gmdif[i].Benefits.PhiEnabled = true
								gmdif[i].Benefits.PhiMultiple = cat.PhiIncomeReplacementPercentage / 100
							}
							if cat.TtdBenefit {
								gmdif[i].Benefits.TtdEnabled = true
								gmdif[i].Benefits.TtdMultiple = cat.TtdIncomeReplacementPercentage / 100
							}
							if cat.FamilyFuneralBenefit {
								gmdif[i].Benefits.GffEnabled = true
							}
						} else {
							if cat.GlaBenefit {
								gmdif[i].Benefits.GlaEnabled = true
							}
							if cat.SglaBenefit {
								gmdif[i].Benefits.SglaEnabled = true
							}
							if cat.PtdBenefit {
								gmdif[i].Benefits.PtdEnabled = true
							}
							if cat.CiBenefit {
								gmdif[i].Benefits.CiEnabled = true
							}
							if cat.PhiBenefit {
								gmdif[i].Benefits.PhiEnabled = true
								gmdif[i].Benefits.PhiMultiple = cat.PhiIncomeReplacementPercentage / 100
							}
							if cat.TtdBenefit {
								gmdif[i].Benefits.TtdEnabled = true
								gmdif[i].Benefits.TtdMultiple = cat.TtdIncomeReplacementPercentage / 100
							}
							if cat.FamilyFuneralBenefit {
								gmdif[i].Benefits.GffEnabled = true
							}
						}
						break
					}
				}
			}

			if len(gmdif) > 0 {
				if err := tx.CreateInBatches(&gmdif, 100).Error; err != nil {
					return fmt.Errorf("failed to insert in-force member data for scheme '%s': %v", schemeName, err)
				}
			}

			// Update scheme to in-force
			scheme.InForce = true
			scheme.Status = models.StatusInForce
			scheme.ActiveSchemeCategories = selectedCats
			if err := tx.Save(&scheme).Error; err != nil {
				return fmt.Errorf("failed to update scheme status for '%s': %v", schemeName, err)
			}

			// Update quote status
			quote.Status = models.StatusAccepted
			quote.SchemeQuoteStatus = models.StatusInEffect
			quote.MemberDataCount = len(membersData)
			quote.ClaimsExperienceCount = schemeClaims
			if err := tx.Save(&quote).Error; err != nil {
				return fmt.Errorf("failed to update quote status for scheme '%s': %v", schemeName, err)
			}

			// --- 8. Insert beneficiaries (Template 4, optional) ---
			if len(benRows) > 0 {
				// Build member ID number -> in-force member ID mapping
				var inForceMembers []models.GPricingMemberDataInForce
				if err := tx.Where("quote_id = ?", quote.ID).Find(&inForceMembers).Error; err != nil {
					return fmt.Errorf("failed to load in-force members for beneficiary linking: %v", err)
				}
				memberIdMap := make(map[string]int) // member_id_number -> DB id
				for _, m := range inForceMembers {
					memberIdMap[strings.TrimSpace(m.MemberIdNumber)] = m.ID
				}

				var beneficiaries []models.Beneficiary
				for _, benRow := range benRows {
					if strings.TrimSpace(benRow.SchemeName) != schemeName {
						continue
					}
					memIdNum := strings.TrimSpace(benRow.MemberIdNumber)
					memberDBId, found := memberIdMap[memIdNum]
					if !found {
						continue // validated already, should not happen
					}

					// Parse benefit types from comma-separated string
					var benefitTypes models.StringArray
					if bt := strings.TrimSpace(benRow.BenefitTypes); bt != "" {
						for _, t := range strings.Split(bt, ",") {
							benefitTypes = append(benefitTypes, strings.TrimSpace(t))
						}
					}

					dob := time.Time(benRow.DateOfBirth)
					var dobPtr *time.Time
					if !dob.IsZero() {
						dobPtr = &dob
					}

					ben := models.Beneficiary{
						FullName:             strings.TrimSpace(benRow.BeneficiaryFullName),
						Relationship:         strings.TrimSpace(benRow.Relationship),
						IDType:               strings.TrimSpace(benRow.IdType),
						IDNumber:             strings.TrimSpace(benRow.IdNumber),
						Gender:               strings.TrimSpace(benRow.Gender),
						DateOfBirth:          dobPtr,
						AllocationPercentage: benRow.AllocationPercentage,
						BenefitTypes:         benefitTypes,
						ContactNumber:        strings.TrimSpace(benRow.ContactNumber),
						Email:                strings.TrimSpace(benRow.Email),
						Address:              strings.TrimSpace(benRow.Address),
						BankName:             strings.TrimSpace(benRow.BankName),
						BranchCode:           strings.TrimSpace(benRow.BranchCode),
						AccountNumber:        strings.TrimSpace(benRow.AccountNumber),
						AccountType:          strings.TrimSpace(benRow.AccountType),
						GuardianName:         strings.TrimSpace(benRow.GuardianName),
						GuardianRelationship: strings.TrimSpace(benRow.GuardianRelationship),
						GuardianIDNumber:     strings.TrimSpace(benRow.GuardianIdNumber),
						GuardianContact:      strings.TrimSpace(benRow.GuardianContact),
						Status:               "active",
						MemberID:             memberDBId,
					}
					beneficiaries = append(beneficiaries, ben)
				}

				if len(beneficiaries) > 0 {
					if err := tx.CreateInBatches(&beneficiaries, 100).Error; err != nil {
						return fmt.Errorf("failed to insert beneficiaries for scheme '%s': %v", schemeName, err)
					}
				}
				result.BeneficiaryCount += len(beneficiaries)
			}

			// Write audit trail
			_ = writeAudit(tx, AuditContext{
				Area:      "Group Pricing",
				Entity:    "Migration",
				EntityID:  strconv.Itoa(scheme.ID),
				Action:    "CREATE",
				ChangedBy: user.UserName,
			}, nil, map[string]interface{}{
				"scheme_name":      schemeName,
				"scheme_id":        scheme.ID,
				"quote_id":         quote.ID,
				"member_count":     len(membersData),
				"category_count":   len(categories),
				"beneficiary_count": result.BeneficiaryCount,
			})

			logger.Infof("Successfully migrated scheme '%s' (ID %d)", schemeName, scheme.ID)
		}

		return nil
	})

	if txErr != nil {
		return nil, fmt.Errorf("migration transaction failed: %v", txErr)
	}

	// Clear cache after successful migration
	if GroupPricingCache != nil {
		GroupPricingCache.Clear()
	}

	return result, nil
}

// parseCSVFile is a convenience wrapper that opens a CSV file header and decodes all rows.
func parseCSVFile[T any](fh *multipart.FileHeader) ([]T, error) {
	dec, _, err := openCSVDecoder(fh)
	if err != nil {
		return nil, err
	}
	return decodeAllRows[T](dec)
}

// GenerateMigrationTemplate builds a CSV template with headers and one example row.
func GenerateMigrationTemplate(templateName string) ([]byte, string, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	switch templateName {
	case "scheme_setup":
		w.Write([]string{
			"scheme_name", "distribution_channel", "broker_name", "broker_email",
			"contact_person", "contact_email", "commencement_date", "cover_start_date",
			"cover_end_date", "renewal_date", "industry", "occupation_class",
			"currency", "normal_retirement_age", "free_cover_limit", "risk_rate_code",
			"experience_rating", "use_global_salary_multiple",
		})
		w.Write([]string{
			"Acme Corp Group Life", "broker", "XYZ Brokers", "broker@xyz.co.za",
			"Jane Smith", "jane@acme.com", "2024-01-01", "2024-01-01",
			"2025-01-01", "2025-01-01", "Manufacturing", "1",
			"ZAR", "65", "5000000", "DEFAULT",
			"yes", "true",
		})
		w.Write([]string{
			"Beta Holdings Staff Cover", "direct", "", "",
			"Tom Brown", "tom@beta.co.za", "2024-06-01", "2024-06-01",
			"2025-06-01", "2025-06-01", "Financial Services", "2",
			"ZAR", "60", "3000000", "DEFAULT",
			"no", "false",
		})

	case "scheme_categories":
		w.Write([]string{
			"scheme_name", "category_name", "basis", "free_cover_limit", "region",
			"gla_benefit", "gla_salary_multiple", "gla_benefit_type", "gla_terminal_illness_benefit", "gla_waiting_period", "gla_educator_benefit",
			"sgla_benefit", "sgla_salary_multiple", "sgla_max_benefit",
			"ptd_benefit", "ptd_salary_multiple", "ptd_benefit_type", "ptd_risk_type", "ptd_deferred_period", "ptd_disability_definition",
			"ci_benefit", "ci_salary_multiple", "ci_benefit_structure", "ci_benefit_definition", "ci_max_benefit",
			"ttd_benefit", "ttd_risk_type", "ttd_maximum_benefit", "ttd_income_replacement_percentage", "ttd_waiting_period", "ttd_deferred_period",
			"phi_benefit", "phi_risk_type", "phi_maximum_benefit", "phi_income_replacement_percentage", "phi_waiting_period", "phi_deferred_period", "phi_normal_retirement_age", "phi_premium_waiver",
			"family_funeral_benefit", "funeral_main_member_sum_assured", "funeral_spouse_sum_assured", "funeral_children_sum_assured", "funeral_max_children",
		})
		w.Write([]string{
			"Acme Corp Group Life", "Management", "actuals", "5000000", "Gauteng",
			"true", "3", "lump_sum", "included", "0", "none",
			"true", "1.5", "3000000",
			"true", "3", "lump_sum", "own_occupation", "6", "standard",
			"true", "1.5", "standalone", "standard", "2000000",
			"true", "own_occupation", "500000", "75", "7", "14",
			"true", "own_occupation", "600000", "75", "7", "30", "65", "included",
			"true", "50000", "50000", "25000", "4",
		})
		w.Write([]string{
			"Acme Corp Group Life", "General Staff", "actuals", "3000000", "Gauteng",
			"true", "2", "lump_sum", "excluded", "0", "none",
			"false", "0", "0",
			"true", "2", "lump_sum", "any_occupation", "3", "standard",
			"false", "0", "standalone", "standard", "0",
			"false", "any_occupation", "0", "0", "0", "0",
			"false", "any_occupation", "0", "0", "0", "0", "65", "excluded",
			"true", "30000", "30000", "15000", "3",
		})

	case "member_data":
		w.Write([]string{
			"year", "scheme_name", "member_name", "member_id_number", "member_id_type",
			"scheme_category", "gender", "email", "employee_number", "date_of_birth",
			"annual_salary", "contribution_waiver_proportion", "entry_date", "exit_date", "effective_exit_date",
			"benefits_gla_multiple", "benefits_sgla_multiple", "benefits_ptd_multiple",
			"benefits_ci_multiple", "benefits_ttd_multiple", "benefits_phi_multiple",
		})
		w.Write([]string{
			"2024", "Acme Corp Group Life", "John Smith", "8001015009087", "RSA_ID",
			"Management", "M", "john@acme.com", "EMP001", "1980-01-01",
			"450000", "0", "2024-01-01", "", "",
			"3", "1.5", "3", "1.5", "0", "0",
		})
		w.Write([]string{
			"2024", "Acme Corp Group Life", "Mary Jones", "8506230009089", "RSA_ID",
			"General Staff", "F", "mary@acme.com", "EMP002", "1985-06-23",
			"320000", "0", "2024-03-01", "", "",
			"2", "0", "2", "0", "0", "0",
		})

	case "beneficiaries":
		w.Write([]string{
			"scheme_name", "member_id_number", "beneficiary_full_name", "relationship",
			"id_type", "id_number", "gender", "date_of_birth", "allocation_percentage",
			"benefit_types", "contact_number", "email", "address",
			"bank_name", "branch_code", "account_number", "account_type",
			"guardian_name", "guardian_relationship", "guardian_id_number", "guardian_contact",
		})
		w.Write([]string{
			"Acme Corp Group Life", "8001015009087", "Sarah Smith", "Spouse",
			"RSA_ID", "8205125009088", "F", "1982-05-12", "60",
			"GLA,PTD", "0821234567", "sarah@email.com", "123 Main St, Sandton",
			"FNB", "250655", "62345678901", "cheque",
			"", "", "", "",
		})
		w.Write([]string{
			"Acme Corp Group Life", "8001015009087", "James Smith", "Child",
			"RSA_ID", "1505015009083", "M", "2015-05-01", "40",
			"GLA", "0821234567", "", "123 Main St, Sandton",
			"FNB", "250655", "62345678901", "savings",
			"Sarah Smith", "Mother", "8205125009088", "0821234567",
		})

	case "claims_experience":
		w.Write([]string{
			"year", "start_date", "end_date",
			"total_gla_sum_assured", "total_ptd_sum_assured", "total_ci_sum_assured",
			"number_of_members", "number_of_gla_claims", "gla_claims_amount",
			"number_of_ptd_claims", "ptd_claims_amount", "number_of_ci_claims", "ci_claims_amount",
			"number_of_phi_claims", "phi_claims_amount", "weighting",
		})
		w.Write([]string{
			"2022", "2022-01-01", "2022-12-31",
			"120000000", "120000000", "60000000",
			"480", "1", "1800000",
			"0", "0", "1", "750000",
			"2", "450000", "0.8",
		})
		w.Write([]string{
			"2023", "2023-01-01", "2023-12-31",
			"150000000", "150000000", "75000000",
			"500", "2", "3500000",
			"1", "1200000", "0", "0",
			"0", "0", "1.0",
		})

	default:
		return nil, "", fmt.Errorf("unknown template: %s", templateName)
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, "", fmt.Errorf("failed to generate template CSV: %v", err)
	}

	return buf.Bytes(), templateName + ".csv", nil
}

package services

import (
	"api/installer"
	"api/models"
	"encoding/json"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
	"time"
)

type PublicKey struct {
	Id  int
	Key string
}

//func LoadDummyData() {
//	var product models.Product
//	DB.Where("product_code = ?", DUMMY_CODE).First(&product)
//
//	if product.ProductCode != strings.ToLower(DUMMY_CODE) {
//		product = CreateDummyProduct()
//		CreateDummyRatingTables(product)
//		AddOtherTableData(product)
//	}
//}

func BaseData(initTables bool) {
	if !AcquireLock("base_data_loading", 5*time.Minute) {
		log.Info().Msg("Base data is already being loaded by another instance, skipping")
		return
	}
	defer ReleaseLock("base_data_loading")

	log.Info().Msg("Loading Base Data")
	//BaseFeatures
	file, err := installer.Files.Open("base_features.json")
	//file, err := installer.Files.Open("/installer/base_features.json")
	if err != nil {
		fmt.Println(err)
	}
	featuresBody, err := io.ReadAll(file)
	baseFeatures := make([]models.BaseFeature, 0)

	_ = json.Unmarshal(featuresBody, &baseFeatures)
	DB.Where("id > 0").Delete(&models.BaseFeature{})
	DB.CreateInBatches(baseFeatures, 100)

	//BaseAssumptionVariables
	file, err = installer.Files.Open("base_assumption_variables.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err := io.ReadAll(file)
	baseAssumptions := make([]models.BaseAssumptionVariable, 0)

	_ = json.Unmarshal(body, &baseAssumptions)
	DB.Where("id > 0").Delete(&models.BaseAssumptionVariable{})
	DB.CreateInBatches(baseAssumptions, 100)

	//BaseMortalityBands
	file, err = installer.Files.Open("base_mortality_bands.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(file)
	baseMortalities := make([]models.BaseMortalityBand, 0)

	_ = json.Unmarshal(body, &baseMortalities)
	DB.Where("id > 0").Delete(&models.BaseMortalityBand{})
	DB.CreateInBatches(baseMortalities, 100)

	//MarkovStates
	file, err = installer.Files.Open("markov_states.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(file)
	markovStates := make([]models.MarkovState, 0)

	_ = json.Unmarshal(body, &markovStates)
	DB.Where("id > 0").Delete(&models.MarkovState{})
	DB.CreateInBatches(markovStates, 100)

	//RatingFactors
	file, err = installer.Files.Open("rating_factors.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(file)
	ratingFactors := make([]models.RatingFactor, 0)

	_ = json.Unmarshal(body, &ratingFactors)
	DB.Where("id > 0").Delete(&models.RatingFactor{})
	DB.CreateInBatches(ratingFactors, 100)

	// Group Pricing User Permissions
	file, err = installer.Files.Open("gp_permissions.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(file)
	gpPermissions := make([]models.GPPermission, 0)

	_ = json.Unmarshal(body, &gpPermissions)

	var gpCount int64
	DB.Model(&models.GPPermission{}).Count(&gpCount)
	if gpCount == 0 {
		DB.Where("id > 0").Delete(&models.GPPermission{})
		err = DB.CreateInBatches(gpPermissions, 100).Error
	} else {
		for _, gp := range gpPermissions {
			var gpPermission models.GPPermission
			DB.Where("slug = ?", gp.Slug).First(&gpPermission)
			if gpPermission.ID == 0 {
				err = DB.Create(&gp).Error
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	// Seed default GP roles (only if no roles exist yet)
	seedDefaultRoles()

	// Valuation User Permissions
	file, err = installer.Files.Open("val_permissions.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(file)
	valPermissions := make([]models.ValPermission, 0)

	_ = json.Unmarshal(body, &valPermissions)

	var valCount int64
	DB.Model(&models.ValPermission{}).Count(&valCount)
	if valCount == 0 {
		DB.Where("id > 0").Delete(&models.ValPermission{})
		err = DB.CreateInBatches(valPermissions, 100).Error
	} else {
		for _, vp := range valPermissions {
			var valPermission models.ValPermission
			DB.Where("slug = ?", vp.Slug).First(&valPermission)
			if valPermission.ID == 0 {
				err = DB.Create(&vp).Error
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	// Seed default valuation roles
	seedValDefaultRoles()

	//ModelPointVariables
	file, err = installer.Files.Open("model_point_variables.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(file)
	mps := make([]models.ModelPointVariable, 0)

	_ = json.Unmarshal(body, &mps)
	DB.Where("id > 0").Delete(&models.ModelPointVariable{})
	err = DB.CreateInBatches(&mps, 100).Error
	if err != nil {
		fmt.Println(err)
	}

	//ConsolidateResults
	file, err = installer.Files.Open("consolidate_results.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(file)
	crs := make([]models.ConsolidateResult, 0)

	_ = json.Unmarshal(body, &crs)
	DB.Where("month > -1").Delete(&models.ConsolidateResult{})
	DB.CreateInBatches(crs, 100)

	//Annual ConsolidateResults
	file, err = installer.Files.Open("annual_consolidated.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(file)
	acrs := make([]models.AnnualConsolidatedResult, 0)

	_ = json.Unmarshal(body, &acrs)
	DB.Where("year > 0").Delete(&models.AnnualConsolidatedResult{})
	DB.CreateInBatches(acrs, 100)

	//CumulativeConsolidateResults
	file, err = installer.Files.Open("cumulative_consolidate_results.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(file)
	ccrs := make([]models.CumulativeConsolidatedResult, 0)

	_ = json.Unmarshal(body, &ccrs)
	DB.Where("month > -1").Delete(&models.CumulativeConsolidatedResult{})
	DB.CreateInBatches(ccrs, 100)

	//AoS Variables
	file, err = installer.Files.Open("aos_variables.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(file)
	aosvars := make([]models.BaseAosVariable, 0)

	_ = json.Unmarshal(body, &aosvars)
	DB.Where("id > 0").Delete(&models.BaseAosVariable{})
	DB.CreateInBatches(aosvars, 100)

	//Lic Variables
	file, err = installer.Files.Open("lic_variables.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(file)
	licvars := make([]models.LicBaseVariable, 0)

	_ = json.Unmarshal(body, &licvars)
	DB.Where("id > 0").Delete(&models.LicBaseVariable{})
	DB.CreateInBatches(licvars, 100)

	//Base Bel Buildup
	file, err = installer.Files.Open("bel_buildup_variables.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(file)
	belvars := make([]models.BelBuildupBaseVariable, 0)

	_ = json.Unmarshal(body, &belvars)
	DB.Where("id > 0").Delete(&models.BelBuildupBaseVariable{})
	DB.CreateInBatches(belvars, 100)

	//LiabilityMovementLines
	file, err = installer.Files.Open("liability_movement_lines.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(file)
	lmls := make([]models.LiabilityMovementLine, 0)

	_ = json.Unmarshal(body, &lmls)
	DB.Where("code > 0").Delete(&models.LiabilityMovementLine{})
	DB.CreateInBatches(lmls, 100)

	//ProductFamilies
	file, err = installer.Files.Open("product_families.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(file)
	pfs := make([]models.ProductFamily, 0)

	_ = json.Unmarshal(body, &pfs)
	var count int64
	DB.Model(&models.ProductFamily{}).Count(&count)
	if count == 0 {
		DB.CreateInBatches(pfs, 100)
	}

	//Group Pricing Age Bands
	file, err = installer.Files.Open("group_business_agebands.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(file)
	groupPricingAgeBands := make([]models.GroupPricingAgeBands, 0)

	_ = json.Unmarshal(body, &groupPricingAgeBands)
	DB.Where("id > 0").Delete(&models.GroupPricingAgeBands{})
	DB.CreateInBatches(groupPricingAgeBands, 100)

	//Group Pricing Benefits
	file, err = installer.Files.Open("group_business_benefits.json")
	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(file)
	groupPricingBenefits := make([]models.GroupBusinessBenefits, 0)

	_ = json.Unmarshal(body, &groupPricingBenefits)
	DB.Where("id > 0").Delete(&models.GroupBusinessBenefits{})
	DB.CreateInBatches(groupPricingBenefits, 100)

	SeedBenefitDocumentTypes()
}

// SeedBenefitDocumentTypes populates the benefit document types if they don't exist
func SeedBenefitDocumentTypes() {
	log.Info().Msg("Seeding Benefit Document Types")

	var count int64
	DB.Model(&models.BenefitDocumentType{}).Count(&count)
	if count > 0 {
		log.Info().Msg("Benefit document types already exist, skipping seeding")
		return
	}

	documentTypesMapping := map[string][]models.BenefitDocumentType{
		"GLA": {
			{Code: "claim_form", Name: "Claim Form (official insurer form)", Required: true},
			{Code: "certified_id_deceased", Name: "Certified ID - Deceased", Required: true},
			{Code: "certified_id_claimant", Name: "Certified ID - Claimant/Beneficiaries", Required: true},
			{Code: "death_certificate", Name: "Death Certificate (BI-5)", Required: true},
			{Code: "dha_notification", Name: "DHA-1663 Notification of Death", Required: true},
			{Code: "beneficiary_form", Name: "Beneficiary Nomination Form / Employer Beneficiary Statement", Required: true},
			{Code: "employment_proof", Name: "Proof of Employment / HR Letter", Required: true},
			{Code: "salary_confirmation", Name: "Salary Confirmation / CTC / Pensionable Salary", Required: true},
			{Code: "banking_details", Name: "Banking Details - beneficiary or member", Required: true},
			{Code: "accident_report", Name: "Accident Report / Police Report (if accidental cause)", Required: false},
			{Code: "post_mortem", Name: "Post-mortem / Final BI-1680/1683", Required: false},
		},
		"SGLA": {
			{Code: "claim_form", Name: "Claim Form (official insurer form)", Required: true},
			{Code: "certified_id_deceased", Name: "Certified ID - Deceased", Required: true},
			{Code: "certified_id_claimant", Name: "Certified ID - Claimant/Beneficiaries", Required: true},
			{Code: "death_certificate", Name: "Death Certificate (BI-5)", Required: true},
			{Code: "dha_notification", Name: "DHA-1663 Notification of Death", Required: true},
			{Code: "beneficiary_form", Name: "Beneficiary Nomination Form / Employer Beneficiary Statement", Required: true},
			{Code: "relationship_proof", Name: "Proof of Relationship (Spouse/Child)", Required: true},
			{Code: "employment_proof", Name: "Proof of Employment / HR Letter", Required: true},
			{Code: "salary_confirmation", Name: "Salary Confirmation / CTC / Pensionable Salary", Required: true},
			{Code: "banking_details", Name: "Banking Details - beneficiary or member", Required: true},
			{Code: "accident_report", Name: "Accident Report / Police Report (if accidental cause)", Required: false},
			{Code: "post_mortem", Name: "Post-mortem / Final BI-1680/1683", Required: false},
		},
		"GFF": {
			{Code: "claim_form", Name: "Claim Form (official insurer form)", Required: true},
			{Code: "certified_id_deceased", Name: "Certified ID - Deceased", Required: true},
			{Code: "certified_id_claimant", Name: "Certified ID - Claimant/Beneficiaries", Required: true},
			{Code: "death_certificate", Name: "Death Certificate (BI-5)", Required: true},
			{Code: "dha_notification", Name: "DHA-1663 Notification of Death", Required: true},
			{Code: "beneficiary_form", Name: "Beneficiary Nomination Form / Employer Beneficiary Statement", Required: true},
			{Code: "relationship_proof", Name: "Proof of Relationship (Spouse/Child)", Required: true},
			{Code: "employment_proof", Name: "Proof of Employment / HR Letter", Required: true},
			{Code: "salary_confirmation", Name: "Salary Confirmation / CTC / Pensionable Salary", Required: true},
			{Code: "banking_details", Name: "Banking Details - beneficiary or member", Required: true},
			{Code: "accident_report", Name: "Accident Report / Police Report (if accidental cause)", Required: false},
			{Code: "post_mortem", Name: "Post-mortem / Final BI-1680/1683", Required: false},
		},
		"PTD": {
			{Code: "claim_form", Name: "Claim Form (official insurer form)", Required: true},
			{Code: "certified_id_member", Name: "Certified ID - Member", Required: true},
			{Code: "beneficiary_form", Name: "Beneficiary Nomination Form / Employer Beneficiary Statement", Required: false},
			{Code: "employment_proof", Name: "Proof of Employment / HR Letter", Required: true},
			{Code: "salary_confirmation", Name: "Salary Confirmation / CTC / Pensionable Salary", Required: true},
			{Code: "banking_details", Name: "Banking Details - beneficiary or member", Required: true},
			{Code: "accident_report", Name: "Accident Report / Police Report (if accidental cause)", Required: false},
			{Code: "medical_reports", Name: "Medical Reports - treating doctor report", Required: true},
			{Code: "attending_doctor_statement", Name: "Attending Doctor's Statement (Disability/CI Report)", Required: true},
			{Code: "specialist_report", Name: "Specialist Medical Report (e.g., Oncologist, Cardiologist, Neurologist)", Required: false},
			{Code: "employer_duties_statement", Name: "Employer Statement of Duties / Job Description", Required: true},
			{Code: "functional_capacity_assessment", Name: "Functional Capacity Assessment (FCE)", Required: false},
			{Code: "occupational_therapist_report", Name: "Occupational Therapist Report", Required: false},
		},
		"CI": {
			{Code: "claim_form", Name: "Claim Form (official insurer form)", Required: true},
			{Code: "certified_id_member", Name: "Certified ID - Member", Required: true},
			{Code: "beneficiary_form", Name: "Beneficiary Nomination Form / Employer Beneficiary Statement", Required: false},
			{Code: "employment_proof", Name: "Proof of Employment / HR Letter", Required: true},
			{Code: "salary_confirmation", Name: "Salary Confirmation / CTC / Pensionable Salary", Required: true},
			{Code: "banking_details", Name: "Banking Details - beneficiary or member", Required: true},
			{Code: "accident_report", Name: "Accident Report / Police Report (if accidental cause)", Required: false},
			{Code: "medical_reports", Name: "Medical Reports - treating doctor report", Required: true},
			{Code: "attending_doctor_statement", Name: "Attending Doctor's Statement (Disability/CI Report)", Required: true},
			{Code: "specialist_report", Name: "Specialist Medical Report (e.g., Oncologist, Cardiologist, Neurologist)", Required: true},
			{Code: "employer_duties_statement", Name: "Employer Statement of Duties / Job Description", Required: false},
			{Code: "functional_capacity_assessment", Name: "Functional Capacity Assessment (FCE)", Required: false},
			{Code: "occupational_therapist_report", Name: "Occupational Therapist Report", Required: false},
		},
		"PHI": {
			{Code: "claim_form", Name: "Claim Form (official insurer form)", Required: true},
			{Code: "certified_id_member", Name: "Certified ID - Member", Required: true},
			{Code: "beneficiary_form", Name: "Beneficiary Nomination Form / Employer Beneficiary Statement", Required: false},
			{Code: "employment_proof", Name: "Proof of Employment / HR Letter", Required: true},
			{Code: "salary_confirmation", Name: "Salary Confirmation / CTC / Pensionable Salary", Required: true},
			{Code: "banking_details", Name: "Banking Details - beneficiary or member", Required: true},
			{Code: "accident_report", Name: "Accident Report / Police Report (if accidental cause)", Required: false},
			{Code: "medical_reports", Name: "Medical Reports - treating doctor report", Required: true},
			{Code: "specialist_report", Name: "Specialist Medical Report (e.g., Oncologist, Cardiologist, Neurologist)", Required: true},
			{Code: "employer_duties_statement", Name: "Employer Statement of Duties / Job Description", Required: true},
			{Code: "functional_capacity_assessment", Name: "Functional Capacity Assessment (FCE)", Required: false},
			{Code: "occupational_therapist_report", Name: "Occupational Therapist Report", Required: false},
			{Code: "psychiatric_report", Name: "Psychiatric Report (if mental illness claim)", Required: false},
			{Code: "income_loss_proof", Name: "Proof of Income Loss / Sick Leave Records", Required: true},
		},
		"TTD": {
			{Code: "claim_form", Name: "Claim Form (official insurer form)", Required: true},
			{Code: "certified_id_member", Name: "Certified ID - Member", Required: true},
			{Code: "medical_reports", Name: "Medical Reports - treating doctor report", Required: true},
			{Code: "banking_details", Name: "Banking Details - beneficiary or member", Required: true},
			{Code: "employment_proof", Name: "Proof of Employment / HR Letter", Required: true},
			{Code: "salary_confirmation", Name: "Salary Confirmation / CTC / Pensionable Salary", Required: true},
		},
	}

	for benefitCode, docTypes := range documentTypesMapping {
		for _, docType := range docTypes {
			docType.BenefitCode = benefitCode
			if err := DB.Create(&docType).Error; err != nil {
				log.Error().Err(err).Msgf("Failed to seed benefit document type: %s for benefit: %s", docType.Code, benefitCode)
			}
		}
	}

	log.Info().Msg("Successfully seeded benefit document types")
}

type DefaultRole struct {
	RoleName    string   `json:"role_name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

func seedDefaultRoles() {
	file, err := installer.Files.Open("gp_default_roles.json")
	if err != nil {
		log.Error().Err(err).Msg("Failed to open gp_default_roles.json")
		return
	}
	body, err := io.ReadAll(file)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read gp_default_roles.json")
		return
	}

	var defaultRoles []DefaultRole
	if err := json.Unmarshal(body, &defaultRoles); err != nil {
		log.Error().Err(err).Msg("Failed to parse gp_default_roles.json")
		return
	}

	for _, dr := range defaultRoles {
		var role models.GPUserRole
		DB.Where("role_name = ?", dr.RoleName).Preload("Permissions").First(&role)

		if role.ID == 0 {
			// Role doesn't exist — create it
			role = models.GPUserRole{
				RoleName:    dr.RoleName,
				Description: dr.Description,
			}
			if err := DB.Create(&role).Error; err != nil {
				log.Error().Err(err).Msgf("Failed to create role: %s", dr.RoleName)
				continue
			}
			log.Info().Msgf("Created default role: %s", dr.RoleName)
		}

		// Associate permissions if the role currently has none
		if len(role.Permissions) == 0 && len(dr.Permissions) > 0 {
			var perms []models.GPPermission
			DB.Where("slug IN ?", dr.Permissions).Find(&perms)
			if len(perms) > 0 {
				if err := DB.Model(&role).Association("Permissions").Replace(&perms); err != nil {
					log.Error().Err(err).Msgf("Failed to associate permissions for role: %s", dr.RoleName)
				} else {
					log.Info().Msgf("Associated %d permissions to role: %s", len(perms), dr.RoleName)
				}
			} else {
				log.Warn().Msgf("No matching permissions found for role: %s", dr.RoleName)
			}
		}
	}

	log.Info().Msg("Default role seeding complete")
}

func seedValDefaultRoles() {
	file, err := installer.Files.Open("val_default_roles.json")
	if err != nil {
		log.Error().Err(err).Msg("Failed to open val_default_roles.json")
		return
	}
	body, err := io.ReadAll(file)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read val_default_roles.json")
		return
	}

	var defaultRoles []DefaultRole
	if err := json.Unmarshal(body, &defaultRoles); err != nil {
		log.Error().Err(err).Msg("Failed to parse val_default_roles.json")
		return
	}

	for _, dr := range defaultRoles {
		var role models.ValUserRole
		DB.Where("role_name = ?", dr.RoleName).Preload("Permissions").First(&role)

		if role.ID == 0 {
			role = models.ValUserRole{
				RoleName:    dr.RoleName,
				Description: dr.Description,
			}
			if err := DB.Create(&role).Error; err != nil {
				log.Error().Err(err).Msgf("Failed to create val role: %s", dr.RoleName)
				continue
			}
			log.Info().Msgf("Created default val role: %s", dr.RoleName)
		}

		if len(role.Permissions) == 0 && len(dr.Permissions) > 0 {
			var perms []models.ValPermission
			DB.Where("slug IN ?", dr.Permissions).Find(&perms)
			if len(perms) > 0 {
				if err := DB.Model(&role).Association("Permissions").Replace(&perms); err != nil {
					log.Error().Err(err).Msgf("Failed to associate permissions for val role: %s", dr.RoleName)
				} else {
					log.Info().Msgf("Associated %d permissions to val role: %s", len(perms), dr.RoleName)
				}
			} else {
				log.Warn().Msgf("No matching permissions found for val role: %s", dr.RoleName)
			}
		}
	}

	log.Info().Msg("Default valuation role seeding complete")
}

//func AddOtherTableData(product models.Product) {
//	// Product Features
//	file, _ := installer.Files.Open("dummy/features.json")
//	body, _ := io.ReadAll(file)
//	var features models.ProductFeatures
//	err := json.Unmarshal(body, &features)
//	if err != nil {
//		fmt.Println(err)
//	}
//	features.ProductCode = product.ProductCode
//	DB.Save(&features)
//
//	// Yield Curve
//	file, _ = installer.Files.Open("dummy/yield_curve.json")
//	body, _ = io.ReadAll(file)
//	var yieldData []models.YieldCurve
//	_ = json.Unmarshal(body, &yieldData)
//
//	for _, dataPoint := range yieldData {
//		err = DB.Save(&dataPoint).Error
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//
//	// Parameters
//	file, _ = installer.Files.Open("dummy/parameters.json")
//	body, _ = io.ReadAll(file)
//	var parameters models.ProductParameters
//	_ = json.Unmarshal(body, &parameters)
//
//	parameters.ProductCode = product.ProductCode
//	DB.Save(&parameters)
//
//	// Margins
//	file, _ = installer.Files.Open("dummy/margins.json")
//	body, _ = io.ReadAll(file)
//	var margins models.ProductMargins
//
//	_ = json.Unmarshal(body, &margins)
//	margins.ProductCode = product.ProductCode
//	DB.Save(&margins)
//
//	// Accident Benefit Multipliers
//	file, _ = installer.Files.Open("dummy/accident_multipliers.json")
//	body, _ = io.ReadAll(file)
//	var multipliers models.ProductAccidentalBenefitMultiplier
//
//	_ = json.Unmarshal(body, &multipliers)
//	multipliers.ProductCode = product.ProductCode
//	DB.Save(&multipliers)
//
//	// Lapse Margins
//	file, _ = installer.Files.Open("dummy/lapse_margins.json")
//	body, _ = io.ReadAll(file)
//	var lapseMargins []models.ProductLapseMargin
//
//	_ = json.Unmarshal(body, &lapseMargins)
//	for _, lapseMargin := range lapseMargins {
//		lapseMargin.ProductCode = product.ProductCode
//		err = DB.Save(&lapseMargin).Error
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//
//	// Clawback
//	file, _ = installer.Files.Open("dummy/clawback.json")
//	body, _ = io.ReadAll(file)
//	var clawbacks []models.ProductClawback
//
//	_ = json.Unmarshal(body, &clawbacks)
//
//	for _, clawback := range clawbacks {
//		clawback.ProductCode = product.ProductCode
//		err = DB.Save(&clawback).Error
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//
//	// Child Sum Assured
//	file, _ = installer.Files.Open("dummy/child_sum_assured.json")
//	body, _ = io.ReadAll(file)
//	var childSumAssureds []models.ProductChildSumAssured
//
//	_ = json.Unmarshal(body, &childSumAssureds)
//
//	for _, csa := range childSumAssureds {
//		csa.ProductCode = product.ProductCode
//		err = DB.Save(&csa).Error
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//
//	// Child Funeral Service
//	file, _ = installer.Files.Open("dummy/child_additional_sum_assured.json")
//	body, _ = io.ReadAll(file)
//	var childFuneralServices []models.ProductChildAdditionalSumAssured
//
//	_ = json.Unmarshal(body, &childFuneralServices)
//
//	for _, cfs := range childFuneralServices {
//		cfs.ProductCode = product.ProductCode
//		err = DB.Save(&cfs).Error
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//
//	// Additional Sum Assured
//	file, _ = installer.Files.Open("dummy/additional_sum_assured.json")
//	body, _ = io.ReadAll(file)
//	var pfs models.ProductAdditionalSumAssured
//
//	_ = json.Unmarshal(body, &pfs)
//	pfs.ProductCode = product.ProductCode
//	DB.Save(&pfs)
//
//	//Product ModelPoint Variables
//	file, err = installer.Files.Open("model_point_variables.json")
//	if err != nil {
//		fmt.Println(err)
//	}
//	body, err = io.ReadAll(file)
//	mps := make([]models.ProductModelpointVariable, 0)
//
//	_ = json.Unmarshal(body, &mps)
//	for _, mp := range mps {
//		mp.ProductID = product.ID
//		_ = DB.Save(&mp)
//	}
//
//	// product riders
//	file, err = installer.Files.Open("dummy/dummy_product_riders.json")
//	if err != nil {
//		fmt.Println(err)
//	}
//	body, err = io.ReadAll(file)
//	prs := make([]models.ProductRider, 0)
//
//	_ = json.Unmarshal(body, &prs)
//	for _, pr := range prs {
//		pr.ProductCode = product.ProductCode
//		_ = DB.Save(&pr)
//	}
//
//}

//func CreateDummyProduct() models.Product {
//	file, err := installer.Files.Open("dummy/dummyconfig.json")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	body, err := io.ReadAll(file)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	fmt.Println(string(body))
//	var product models.Product
//
//	//err = json.Unmarshal(body, &product)
//	//err = CreateProduct(product)
//	return product
//}

//func CreateDummyRatingTables(product models.Product) {
//	// Create transition state tables.
//	err := DB.Exec(fmt.Sprintf("CREATE TABLE `%s_disability` ( `year_anb_sec_occ_class_gender` varchar(100) NOT NULL, `incidence_rate` double DEFAULT NULL, PRIMARY KEY (`year_anb_sec_occ_class_gender`))", strings.ToLower(product.ProductCode))).Error
//	err = DB.Exec(fmt.Sprintf("CREATE TABLE `%s_lapse` (`year_duration_if_m` varchar(100) NOT NULL,`lapse_rate` double DEFAULT NULL,PRIMARY KEY (`year_duration_if_m`))", strings.ToLower(product.ProductCode))).Error
//	err = DB.Exec(fmt.Sprintf("CREATE TABLE `%s_mortality_accidental` (`year_anb_gender` varchar(100) NOT NULL,`acc_qx` double DEFAULT NULL,PRIMARY KEY (`year_anb_gender`))", strings.ToLower(product.ProductCode))).Error
//	err = DB.Exec(fmt.Sprintf("CREATE TABLE `%s_mort_table` (`year_anb_gender` varchar(100) NOT NULL,`qx` double DEFAULT NULL,PRIMARY KEY (`year_anb_gender`))", strings.ToLower(product.ProductCode))).Error
//	err = DB.Exec(fmt.Sprintf("CREATE TABLE `%s_retrenchment` (`year_duration_if_y` varchar(100) NOT NULL,`retr_rate` double DEFAULT NULL,PRIMARY KEY (`year_duration_if_y`))", strings.ToLower(product.ProductCode))).Error
//
//	//Populate the respective tables...
//	file, _ := installer.Files.Open("dummy/dummy_disability.json")
//	body, _ := io.ReadAll(file)
//	var disabilities []map[string]interface{}
//	err = json.Unmarshal(body, &disabilities)
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(disabilities)
//		for _, disability := range disabilities {
//			_ = AddToRatingTable(strings.ToLower(product.ProductCode)+"_disability", disability["year_anb_gender_sec_occ_class"].(string), disability["incidence_rate"].(float64), "year_anb_sec_occ_class_gender")
//		}
//	}
//
//	file, _ = installer.Files.Open("dummy/dummy_lapse.json")
//	body, _ = io.ReadAll(file)
//	var lapses []map[string]interface{}
//	err = json.Unmarshal(body, &lapses)
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		for _, lapse := range lapses {
//			AddToRatingTable(strings.ToLower(product.ProductCode)+"_lapse", lapse["year_duration_if_m"].(string), lapse["lapse_rate"].(float64), "year_duration_if_m")
//		}
//	}
//
//	file, _ = installer.Files.Open("dummy/dummy_retrenchment.json")
//	body, _ = io.ReadAll(file)
//	var retrenchments []map[string]interface{}
//	err = json.Unmarshal(body, &retrenchments)
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		for _, retrenchment := range retrenchments {
//			AddToRatingTable(strings.ToLower(product.ProductCode)+"_retrenchment", retrenchment["year_duration_if_y"].(string), retrenchment["retr_rate"].(float64), "year_duration_if_y")
//		}
//	}
//
//	file, _ = installer.Files.Open("dummy/dummy_mortality_accidental.json")
//	body, _ = io.ReadAll(file)
//	var mortalities []map[string]interface{}
//	err = json.Unmarshal(body, &mortalities)
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		for _, mortality := range mortalities {
//			AddToRatingTable(strings.ToLower(product.ProductCode)+"_mortality_accidental", mortality["year_anb_gender"].(string), mortality["acc_qx"].(float64), "year_anb_gender")
//		}
//	}
//
//	file, _ = installer.Files.Open("dummy/dummy_mort_table.json")
//	body, _ = io.ReadAll(file)
//	var morts []map[string]interface{}
//	err = json.Unmarshal(body, &morts)
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		for _, mortality := range morts {
//			AddToRatingTable(strings.ToLower(product.ProductCode)+"_mort_table", mortality["year_anb_gender"].(string), mortality["qx"].(float64), "year_anb_gender")
//		}
//	}
//
//	file, _ = installer.Files.Open("dummy/dummy_modelpoints.json")
//	body, _ = io.ReadAll(file)
//	var modelpoints []models.ProductModelPoint
//	err = json.Unmarshal(body, &modelpoints)
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		for _, modelpoint := range modelpoints {
//			err := DB.Table(strings.ToLower(product.ProductCode) + "_modelpoints").Save(&modelpoint).Error
//			if err != nil {
//				fmt.Println(err)
//			}
//		}
//	}
//
//}

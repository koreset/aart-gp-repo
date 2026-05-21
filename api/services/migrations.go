package services

import (
	appLog "api/log"
	"api/models"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// MigrateBaseTables migrates the base tables
func MigrateBaseTables() error {
	appLog.Info("Migrating base tables")
	err := DB.AutoMigrate(
		&models.SystemLock{},
		&models.System{},
		&models.AppUser{},
		&models.UserToken{},
		&models.GroupPricingAgeBands{},
		&models.MedicalWaiver{},
		&models.GroupBusinessBenefits{},
		&models.Activity{},
		&models.GPPermission{},
		&models.AuditLog{},
	)

	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to migrate base tables")
		return err
	}

	appLog.Info("Successfully migrated base tables")
	return nil
}

// MigrateEmailTables migrates the email tables
func MigrateEmailTables() error {
	appLog.Info("Migrating email tables")
	err := DB.AutoMigrate(
		&models.EmailSettings{},
		&models.EmailTemplate{},
		&models.EmailOutbox{},
	)

	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to migrate email tables")
		return err
	}

	appLog.Info("Successfully migrated email tables")
	return nil
}

// MigrateGroupPricingTables migrates the group pricing tables
func MigrateGroupPricingTables() error {
	appLog.Info("Migrating group pricing tables")
	err := DB.Migrator().AutoMigrate(
		&models.CalculationJob{},
		&models.GroupPricingQuote{},
		&models.GroupRiskQuoteStats{},
		&models.SchemeCategory{},
		&models.SchemeCategoryMaster{},
		&models.Broker{},
		&models.GPricingMemberData{},
		&models.GPricingMemberDataInForce{},
		&models.BulkEnrollmentBatch{},
		&models.GroupPricingClaimsExperience{},
		&models.GroupPricingExperienceRateOverride{},
		&models.MemberRatingResult{},
		&models.MemberRatingResultSummary{},
		&models.MovementMemberRatingResult{},
		&models.GlaRate{},
		&models.PhiRate{},
		&models.SalaryContinuationRate{},
		&models.PtdRate{},
		&models.CiRate{},
		&models.AccidentalTtdRate{},
		&models.TtdRate{},
		&models.ChildMortality{},
		&models.IndustryLoading{},
		&models.FuneralParameters{},
		&models.GroupPricingReinsuranceStructure{},
		&models.GroupPricingParameters{},
		&models.GroupScheme{},
		&models.GroupPricingInsurerDetail{},
		&models.MemberPremiumSchedule{},
		&models.IncomeLevel{},
		&models.OccupationClass{},
		&models.EducatorBenefitStructure{},
		&models.EducatorRate{},
		&models.BinderFee{},
		&models.CommissionStructure{},
		&models.Escalations{},
		&models.FuneralAidsRate{},
		&models.FuneralRate{},
		&models.GeneralLoading{},
		&models.GlaAidsRate{},
		&models.RegionLoading{},
		&models.UnderwritingTierConfig{},
		&models.UnderwritingCase{},
		&models.UnderwritingDecision{},
		&models.UnderwritingCaseEvent{},
		&models.UnderwritingCaseAttachment{},
		&models.UWRuleSet{},
		&models.UWRule{},
		&models.UWConditionCode{},
		&models.QuoteReRateEvent{},
		&models.MemberDisclosure{},
		&models.ActivelyAtWorkAttestation{},
		&models.ConsentRecord{},
		&models.VendorRequest{},
		&models.VendorWebhook{},
		&models.PriorInsurerSchedule{},
		&models.PriorInsurerMember{},
		&models.PolicyHandoffSnapshot{},
		&models.TieredIncomeReplacement{},
		&models.CustomTieredIncomeReplacement{},
		&models.DiscountAuthority{},
		&models.Restriction{},
		&models.ReinsuranceCoverRestriction{},
		&models.ReinsuranceGlaRate{},
		&models.ReinsuranceCiRate{},
		&models.ReinsurancePtdRate{},
		&models.ReinsurancePhiRate{},
		&models.ReinsuranceSalaryContinuationRate{},
		&models.ReinsuranceFuneralAidsRate{},
		&models.ReinsuranceFuneralRate{},
		&models.ReinsuranceGlaAidsRate{},
		&models.ReinsuranceGeneralLoading{},
		&models.ReinsuranceIndustryLoading{},
		&models.ReinsuranceRegionLoading{},
		&models.PremiumLoading{},
		&models.SchemeSizeLevel{},
		&models.TaxTable{},
		&models.TaxRetirementTable{},
		&models.GroupSchemeExposure{},
		&models.GroupPricingAgeBands{},
		&models.MedicalWaiver{},
		&models.GroupBusinessBenefits{},
		&models.GroupSchemeClaim{},
		&models.GroupSchemeClaimAttachment{},
		&models.GroupSchemeClaimAssessment{},
		&models.GroupSchemeClaimCommunication{},
		&models.GroupSchemeClaimDecline{},
		&models.GroupSchemeClaimStatusAudit{},
		&models.FraudRiskModel{},
		&models.FraudRiskRule{},
		&models.FraudRiskAssessment{},
		&models.ClaimPaymentSchedule{},
		&models.ClaimPaymentScheduleItem{},
		&models.ClaimPaymentProof{},
		&models.ACBBankProfile{},
		&models.ACBFileRecord{},
		&models.ACBReconciliationResult{},
		&models.PremiumBordereauxData{},
		&models.MemberBordereauxData{},
		&models.ClaimBordereauxData{},
		&models.HistoricalCredibilityData{},
		&models.MemberIndicativeDataSet{},
		&models.GroupSchemeStatusAudit{},
		&models.GroupPricingQuoteStatusAudit{},
		&models.QuoteSlaTarget{},
		&models.MemberActivity{},
		&models.BordereauxTemplate{},
		&models.GeneratedBordereaux{},
		&models.BordereauxTimeline{},
		&models.BordereauxConfiguration{},
		&models.BordereauxConfirmation{},
		&models.BordereauxConfirmationNote{},
		&models.BordereauxConfirmationRecord{},
		&models.BordereauxReconciliationResult{},
		&models.SchemeType{},
		&models.BenefitDocumentType{},
		&models.AuditLog{},
		&models.OnRiskLetter{},
		&models.WinProbabilityModel{},
		&models.QuoteWinProbability{},
		&models.MemberBeneficiary{},
		&models.Beneficiary{},
		&models.Bordereaux{},
		&models.Reinsurer{},
		&models.GroupBenefitMapper{},
		&models.GPTableStat{},
		&models.InsurerQuoteTemplate{},
		&models.InsurerOnRiskLetterTemplate{},
		&models.BAVVerificationLog{},
		&models.EmailSettings{},
		&models.EmailTemplate{},
		&models.EmailOutbox{},
		&models.GroupPricingSetting{},
		&models.PaymentLetterSetting{},
		&models.ClaimPaymentLetter{},
		&models.ClaimPaymentLetterDelivery{},
		// Generic table-requirement framework (also auto-migrated lazily by
		// services.EnsureTableConfigurations on every boot). Listed here so a
		// fresh install creates these in the same bootstrap pass as everything
		// else, instead of relying on the post-SetupTables call.
		&models.TableConfiguration{},
		&models.TableConfigurationAuditLog{},
	)

	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to migrate group pricing tables")
		return err
	}

	// Seed the singleton GroupPricingSetting row with the default discount
	// method so existing quotes recompute identically to historical behaviour
	// until the user explicitly switches to prorata via MetaData.
	var settingCount int64
	DB.Model(&models.GroupPricingSetting{}).Count(&settingCount)
	if settingCount == 0 {
		now := time.Now()
		DB.Create(&models.GroupPricingSetting{
			ID:                      1,
			DiscountMethod:          models.DiscountMethodLoadingAdjustment,
			DiscountMethodUpdatedAt: &now,
			DiscountMethodUpdatedBy: "system",
			FCLMethod:               models.FCLMethodPercentile,
			FCLMethodUpdatedAt:      &now,
			FCLMethodUpdatedBy:      "system",
			FCLOverrideTolerance:    FCLOverrideToleranceDefault,
		})
	} else {
		// Backfill audit timestamps on the existing singleton so the
		// "Last updated" caption renders even before the user touches the
		// panel. Use updated_at as a best-effort stand-in for the historical
		// touch time and tag the actor as "system" so it's clear the row
		// was not explicitly set by a person.
		DB.Model(&models.GroupPricingSetting{}).
			Where("id = ? AND discount_method_updated_at IS NULL", 1).
			Updates(map[string]interface{}{
				"discount_method_updated_at": gorm.Expr("updated_at"),
				"discount_method_updated_by": "system",
			})
		DB.Model(&models.GroupPricingSetting{}).
			Where("id = ? AND fcl_method_updated_at IS NULL", 1).
			Updates(map[string]interface{}{
				"fcl_method_updated_at": gorm.Expr("updated_at"),
				"fcl_method_updated_by": "system",
			})
	}

	// Ensure Beneficiary columns exist on the member_beneficiaries table.
	// MemberBeneficiary and Beneficiary share the same table; the richer
	// Beneficiary model may have columns that a prior migration did not add.
	if !DB.Migrator().HasColumn(&models.Beneficiary{}, "FullName") {
		if colErr := DB.Migrator().AddColumn(&models.Beneficiary{}, "FullName"); colErr != nil {
			appLog.WithField("error", colErr.Error()).Warn("Could not add full_name column to member_beneficiaries (may already exist)")
		}
		// Backfill: copy legacy 'name' column into 'full_name' for existing rows
		if DB.Migrator().HasColumn(&models.MemberBeneficiary{}, "Name") {
			DB.Exec("UPDATE member_beneficiaries SET full_name = name WHERE full_name IS NULL OR full_name = ''")
		}
	}

	appLog.Info("Successfully migrated group pricing tables")
	return nil
}

// MigrateGroupPricingUserTables migrates the group pricing user tables
func MigrateGroupPricingUserTables() error {
	appLog.Info("Migrating group pricing user tables")
	err := DB.Migrator().AutoMigrate(
		&models.GPUserRole{},
		&models.GPUser{},
		&models.OrgUser{},
		&models.GPPermission{},
	)

	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to migrate group pricing user tables")
		return err
	}

	appLog.Info("Successfully migrated group pricing user tables")
	return nil
}

func MigratePhiValuationTables() error {
	appLog.Info("Migrating PHI valuation tables")
	err := DB.Migrator().AutoMigrate(
		&models.PhiModelPoint{},
		&models.PhiParameter{},
		&models.PhiYieldCurve{},
		&models.PhiRecoveryRate{},
		&models.PhiMortality{},
		&models.PhiShock{},
		&models.PhiReinsurance{},
		&models.PhiShockSetting{},
		&models.PhiProjections{},
		&models.PhiAggregatedProjections{},
		&models.PhiScopedAggregatedProjections{},
		&models.PhiRunParameters{},
		&models.PhiRunConfig{},
		// Job envelope for PHI projection runs. Holds the user, the AppUser
		// snapshot, and the list of PhiRunParameters that the engine iterates
		// over. Used by services.RunPhiProjection (see services/phi_valuations.go).
		&models.RunPhiJob{},
	)

	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to migrate PHI valuation tables")
		return err
	}

	appLog.Info("Successfully migrated PHI valuation tables")
	return nil
}

// MigrateGroupPremiumTables migrates the group premium lifecycle tables
func MigrateGroupPremiumTables() error {
	appLog.Info("Migrating group premium tables")
	err := DB.AutoMigrate(
		&models.ContributionConfig{},
		&models.PremiumSchedule{},
		&models.ScheduleMemberRow{},
		&models.Invoice{},
		&models.InvoiceAdjustment{},
		&models.Payment{},
		&models.ArrearsHistory{},
		&models.PaymentPlan{},
		&models.PaymentPlanInstalment{},
		&models.EmployerSubmission{},
		&models.EmployerSubmissionRecord{},
		&models.SubmissionDeltaRecord{},
		&models.BordereauxDeadline{},
		&models.ReinsurerAcceptance{},
		&models.ReinsurerRecovery{},
		&models.ClaimNotificationLog{},
		&models.NewJoinerDetail{},
		&models.RegisterDiffSnapshot{},
		&models.ReinsuranceTreaty{},
		&models.TreatySchemeLink{},
		&models.RIBordereauxRun{},
		&models.RIBordereauxMemberRow{},
		&models.RIBordereauxClaimsRow{},
		&models.RIValidationResult{},
		&models.LargeClaimNotice{},
		&models.TechnicalAccount{},
		&models.SettlementPayment{},
		// Premium Reconciliation v2 (ledger-based allocation)
		&models.ReconciliationRun{},
		&models.PaymentAllocation{},
		&models.ReconciliationItem{},
		&models.MatchingRule{},
		// In-app communications
		&models.Notification{},
		&models.Conversation{},
		&models.ConversationParticipant{},
		&models.ConversationMessage{},
	)

	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to migrate group premium tables")
		return err
	}

	appLog.Info("Successfully migrated group premium tables")

	// Backfill reconciliation items for any invoices/payments that existed
	// before the v2 reconciliation system was introduced.
	BackfillReconciliationItems()

	return nil
}

// UpdateModelPointTablesForMigration updates model point tables for all products
func UpdateModelPointTablesForMigration() error {
	appLog.Info("Updating model point tables for all products")

	//Get a list of all products
	var products []models.Product
	err := DB.Find(&products).Error
	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to fetch products")
		return err
	}

	//loop through each product and construct the model point table from the product code and _modelpoint
	for _, product := range products {
		mpTableName := strings.ToLower(product.ProductCode) + "_modelpoints"
		pricingMpTableName := strings.ToLower(product.ProductCode) + "_pricing_modelpoints"

		//check if table exists
		if DB.Migrator().HasTable(mpTableName) {
			err := DB.Table(mpTableName).AutoMigrate(&models.ProductModelPoint{})
			if err != nil {
				appLog.WithFields(map[string]interface{}{
					"error": err.Error(),
					"table": mpTableName,
				}).Error("Failed to migrate product model point table")
				return err
			}
			appLog.WithField("table", mpTableName).Info("Successfully migrated product model point table")
		}

		if DB.Migrator().HasTable(pricingMpTableName) {
			err := DB.Table(pricingMpTableName).AutoMigrate(&models.ProductPricingModelPoint{})
			if err != nil {
				appLog.WithFields(map[string]interface{}{
					"error": err.Error(),
					"table": pricingMpTableName,
				}).Error("Failed to migrate product pricing model point table")
				return err
			}
			appLog.WithField("table", pricingMpTableName).Info("Successfully migrated product pricing model point table")
		}
	}

	appLog.Info("Successfully updated all model point tables")
	return nil
}

// CreateDatabaseIndexesForMigration creates indexes on frequently queried tables
func CreateDatabaseIndexesForMigration() error {
	appLog.Info("Creating database indexes for performance optimization")

	// List of indexes to create
	indexes := []struct {
		table      string
		columns    []string
		indexName  string
		unique     bool
		conditions string
	}{
		// Aggregated projections indexes
		{table: "aggregated_projections", columns: []string{"run_id"}, indexName: "idx_agg_proj_run_id", unique: false},
		{table: "aggregated_projections", columns: []string{"run_id", "product_code"}, indexName: "idx_agg_proj_run_product", unique: false},
		{table: "aggregated_projections", columns: []string{"run_id", "product_code", "sp_code"}, indexName: "idx_agg_proj_run_product_sp", unique: false},
		{table: "aggregated_projections", columns: []string{"projection_month"}, indexName: "idx_agg_proj_month", unique: false},

		// Group pricing indexes
		{table: "g_pricing_member_data", columns: []string{"quote_id"}, indexName: "idx_member_data_quote_id", unique: false},
		{table: "group_pricing_claims_experiences", columns: []string{"quote_id"}, indexName: "idx_claims_exp_quote_id", unique: false},
		{table: "member_rating_results", columns: []string{"quote_id"}, indexName: "idx_rating_results_quote_id", unique: false},
		{table: "member_premium_schedules", columns: []string{"quote_id"}, indexName: "idx_premium_schedules_quote_id", unique: false},
		{table: "bordereaux", columns: []string{"quote_id"}, indexName: "idx_bordereaux_quote_id", unique: false},

		// User management indexes
		{table: "org_users", columns: []string{"email"}, indexName: "idx_org_users_email", unique: true},
		{table: "app_users", columns: []string{"user_email"}, indexName: "idx_app_users_email", unique: true},
		{table: "user_tokens", columns: []string{"subject"}, indexName: "idx_user_tokens_subject", unique: true},

		// Activity tracking indexes
		{table: "activities", columns: []string{"user_email"}, indexName: "idx_activities_user_email", unique: false},
		{table: "activities", columns: []string{"date"}, indexName: "idx_activities_date", unique: false},
		{table: "activities", columns: []string{"object_type", "object_id"}, indexName: "idx_activities_object", unique: false},
	}

	// Create each index
	for _, idx := range indexes {
		var indexSQL string

		if DbBackend == "postgresql" {
			// PostgreSQL index creation
			indexType := ""
			if idx.unique {
				indexType = "UNIQUE "
			}

			indexSQL = fmt.Sprintf("CREATE %sINDEX IF NOT EXISTS %s ON %s (%s)",
				indexType,
				idx.indexName,
				idx.table,
				strings.Join(idx.columns, ", "))

			if idx.conditions != "" {
				indexSQL += " " + idx.conditions
			}
		} else {
			// MySQL index creation
			indexType := "INDEX"
			if idx.unique {
				indexType = "UNIQUE INDEX"
			}

			indexSQL = fmt.Sprintf("CREATE %s %s ON %s (%s)",
				indexType,
				idx.indexName,
				idx.table,
				strings.Join(idx.columns, ", "))

			if idx.conditions != "" {
				indexSQL += " " + idx.conditions
			}
		}

		// Execute the index creation
		err := DB.Exec(indexSQL).Error
		if err != nil {
			appLog.WithFields(map[string]interface{}{
				"error":      err.Error(),
				"table":      idx.table,
				"index_name": idx.indexName,
				"columns":    strings.Join(idx.columns, ", "),
			}).Warn("Failed to create index")
			// Continue with other indexes even if one fails
		} else {
			appLog.WithFields(map[string]interface{}{
				"table":      idx.table,
				"index_name": idx.indexName,
				"columns":    strings.Join(idx.columns, ", "),
			}).Info("Successfully created index")
		}
	}

	appLog.Info("Database index creation completed")
	return nil
}

// SeedGeneralLedger bootstraps the operational General Ledger with a default
// chart of accounts, the MVP posting rules that hook claim payments and
// premium reconciliation into the ledger, and one open accounting period for
// the current month.
//
// Idempotent: every record is keyed on a unique column (Code, EventKey, or
// Name) and inserted only if missing. Safe to run on every boot. Errors are
// logged but never propagate; a broken seed should not crash the API, and
// finance can correct config via the /gl admin UI.
func SeedGeneralLedger() {
	if DB == nil {
		appLog.Warn("SeedGeneralLedger called with nil DB — skipping")
		return
	}

	defaultAccounts := []models.GLAccount{
		{Code: "1000", Name: "Bank", AccountType: "asset", NormalBalance: "debit", IsActive: true, Description: "Operational bank account"},
		{Code: "1100", Name: "Premium Receivable", AccountType: "asset", NormalBalance: "debit", IsActive: true, Description: "Outstanding premium invoices"},
		{Code: "1200", Name: "Suspense", AccountType: "asset", NormalBalance: "debit", IsActive: true, Description: "Unallocated premium receipts"},
		{Code: "2000", Name: "Claims Liability", AccountType: "liability", NormalBalance: "credit", IsActive: true, Description: "Approved claims awaiting payment"},
		{Code: "2100", Name: "Refunds Payable", AccountType: "liability", NormalBalance: "credit", IsActive: true, Description: "Premium overpayments due for refund"},
		{Code: "4000", Name: "Premium Income", AccountType: "income", NormalBalance: "credit", IsActive: true, Description: "Earned premium revenue"},
		{Code: "6000", Name: "Bad Debt Expense", AccountType: "expense", NormalBalance: "debit", IsActive: true, Description: "Written-off premium receivables"},
	}
	for _, a := range defaultAccounts {
		var existing models.GLAccount
		err := DB.Where("code = ?", a.Code).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			appLog.WithField("code", a.Code).WithField("error", err.Error()).Warn("Failed to check existing GL account")
			continue
		}
		if err := DB.Create(&a).Error; err != nil {
			appLog.WithField("code", a.Code).WithField("error", err.Error()).Warn("Failed to seed GL account")
		}
	}

	accountIDByCode := make(map[string]int)
	var allAccounts []models.GLAccount
	if err := DB.Find(&allAccounts).Error; err == nil {
		for _, a := range allAccounts {
			accountIDByCode[a.Code] = a.ID
		}
	}

	defaultRules := []struct {
		EventKey  string
		Debit     string
		Credit    string
		Notes     string
	}{
		{"claim_payment.confirmed", "2000", "1000", "Confirmed claim payment: clear liability, drain bank"},
		{"premium_allocation.created", "1000", "1100", "Premium allocated to invoice: bank up, receivable down"},
		{"premium_allocation.reversed", "1100", "1000", "Allocation reversed: restore receivable, undo bank entry"},
		{"premium.write_off", "6000", "1100", "Write-off uncollectable premium"},
		{"premium.refund", "2100", "1000", "Refund overpayment to payer"},
	}
	for _, r := range defaultRules {
		drID, drOK := accountIDByCode[r.Debit]
		crID, crOK := accountIDByCode[r.Credit]
		if !drOK || !crOK {
			appLog.WithField("event_key", r.EventKey).Warn("Skipping posting rule seed — referenced account not present")
			continue
		}
		var existing models.PostingRule
		err := DB.Where("event_key = ?", r.EventKey).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			appLog.WithField("event_key", r.EventKey).WithField("error", err.Error()).Warn("Failed to check existing posting rule")
			continue
		}
		rule := models.PostingRule{
			EventKey:        r.EventKey,
			DebitAccountID:  drID,
			CreditAccountID: crID,
			IsActive:        true,
			Notes:           r.Notes,
		}
		if err := DB.Create(&rule).Error; err != nil {
			appLog.WithField("event_key", r.EventKey).WithField("error", err.Error()).Warn("Failed to seed posting rule")
		}
	}

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Second)
	periodName := now.Format("2006-01")
	var existing models.AccountingPeriod
	err := DB.Where("name = ?", periodName).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		period := models.AccountingPeriod{
			Name:      periodName,
			StartDate: monthStart,
			EndDate:   monthEnd,
			Status:    "open",
		}
		if err := DB.Create(&period).Error; err != nil {
			appLog.WithField("name", periodName).WithField("error", err.Error()).Warn("Failed to seed accounting period")
		}
	} else if err != nil {
		appLog.WithField("name", periodName).WithField("error", err.Error()).Warn("Failed to check existing accounting period")
	}

	appLog.Info("General Ledger seeding completed")
}

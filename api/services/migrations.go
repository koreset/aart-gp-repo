package services

import (
	appLog "api/log"
	"api/models"
	"fmt"
	"strings"
)

// MigrateBaseTables migrates the base tables
func MigrateBaseTables() error {
	appLog.Info("Migrating base tables")
	err := DB.AutoMigrate(
		&models.SystemLock{},
		&models.BaseAssumptionVariable{},
		&models.BaseFeature{},
		&models.ModelPointVariable{},
		&models.MarkovState{},
		&models.BaseMortalityBand{},
		&models.BaseAosVariable{},
		&models.LicBaseVariable{},
		&models.System{},
		&models.AppUser{},
		&models.UserToken{},
		&models.ConsolidateResult{},
		&models.CumulativeConsolidatedResult{},
		&models.AnnualConsolidatedResult{},
		&models.BelBuildupBaseVariable{},
		&models.ProductFamily{},
		&models.LiabilityMovementLine{},
		&models.GroupPricingAgeBands{},
		&models.GroupBusinessBenefits{},
		&models.Activity{},
		&models.RatingFactor{},
		&models.GPPermission{},
		&models.AuditLog{},
		&models.TransitionState{},
		&models.RetrenchmentRate{},
		&models.DisabilityIncidenceFactors{},
	)

	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to migrate base tables")
		return err
	}

	appLog.Info("Successfully migrated base tables")
	return nil
}

// MigrateProductModelTables migrates the product model tables
func MigrateProductModelTables() error {
	appLog.Info("Migrating product model tables")
	err := DB.AutoMigrate(
		&models.ProductModelPointCount{},
		&models.AggregatedProjection{},
		&models.ScopedAggregatedProjection{},
		&models.Projection{},

		&models.ProductModelpointVariable{},
		&models.ProductTransitionState{},
		&models.ProductTransition{},
		&models.ProductTable{},
		&models.ProductRatingFactor{},
		&models.ProductFeatures{},
		&models.ProductAccidentalBenefitMultiplier{},
		&models.ProductPricingAccidentalBenefitMultiplier{},
		&models.Fds{},
		&models.GlobalTable{},
		&models.ProductClawback{},
		&models.ProductChildSumAssured{},
		&models.ProductChildAdditionalSumAssured{},
		&models.ProductAdditionalSumAssured{},
		&models.ProductReinsurance{},
		&models.ShockSetting{},
		&models.ProductPricingShock{},
		&models.Task{},
		&models.ManualScopedAggregatedProjection{},
		&models.LICAggregatedProjections{},
		&models.ProductRider{},
		&models.JobsTemplate{},
		&models.JobTemplateContent{},
		&models.ProductNonLifeRating{},
		&models.ProductBenefitMultiplier{},
		&models.ProductSpecialDecrementMargin{},
		&models.ProductRenewableProfitAdjustment{},
		&models.ProductUnitFundCharge{},
		&models.ProductInvestmentReturn{},
		&models.ProductFundAssetDistribution{},
		&models.ProductMaturityPattern{},
		&models.ProductSurrenderValueCoefficient{},
		&models.IBNRAverageClaimAmount{},
		&models.IBNRClaimHistorySummary{},
		&models.LICClaimsAnalysisOfChange{},
		&models.IBNRPaidVsOutstandingClaims{},
		&models.IBNRProportionOutstandingClaims{},
		&models.IBNRIncurredClaims{},
		&models.AggregatedModifiedGMMProjection{},
		&models.ModifiedGMMProjection{},
		&models.ModifiedGMMScopedAggregation{},
		&models.Escalations{},
		&models.YieldCurve{},
		&models.ProductParameters{},
		&models.ProductCommissionStructure{},
		&models.ProductMargins{},
		&models.ProjectionJob{},
		&models.JobProduct{},
		&models.JobProductRunError{},
		&models.CachedReserveResults{},
		&models.RunParameters{},
		&models.ModelPointPricing{},
		&models.Profitability{},
		&models.PricingRun{},
		&models.PricingConfig{},
		&models.RiskDriver{},
		&models.RiskAdjustmentFactor{},
		&models.RiskAdjustmentDriver{},
		&models.RAConfiguration{},
		&models.AOSConfiguration{},
		&models.ProductLapseMargin{},
		&models.AggregatedVariableGroup{},
		&models.Product{},
	)

	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to migrate product model tables")
		return err
	}

	appLog.Info("Successfully migrated product model tables")
	return nil
}

// MigratePricingTables migrates the pricing tables
func MigratePricingTables() error {
	appLog.Info("Migrating pricing tables")
	err := DB.AutoMigrate(
		// Currently empty as commented out in original code
	)

	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to migrate pricing tables")
		return err
	}

	appLog.Info("Successfully migrated pricing tables")
	return nil
}

// MigrateEscalationTables migrates the escalation tables
func MigrateEscalationTables() error {
	appLog.Info("Migrating escalation tables")
	err := DB.AutoMigrate(
		&models.ProductParameters{},
		&models.ProductPricingProductLevelEscalation{},
		&models.PricingPoint{},
		&models.PAAFinance{},
		&models.ProductModelPointVariableStats{},
		&models.PAAResult{},
		&models.BalanceSheetRecord{},
		&models.ProductShock{},
	)

	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to migrate escalation tables")
		return err
	}

	appLog.Info("Successfully migrated escalation tables")
	return nil
}

// MigrateModelPointTables migrates the model point tables
func MigrateModelPointTables() error {
	appLog.Info("Migrating model point tables")
	err := DB.Migrator().AutoMigrate(
		&models.BalanceSheetRecord{},
		&models.ProductPricingMargins{},
		&models.FinanceVariables{},
		&models.PaaModelPointVariableStats{},
		&models.ProductPricingModelPoint{},
		&models.CsmRun{},
		&models.PricingParameter{},
		&models.ProductPricingTable{},
		&models.ProductPricingChildAdditionalSumAssured{},
		&models.ProductPricingChildSumAssured{},
		&models.ProductPricingClawback{},
		&models.ProductPricingAdditionalSumAssured{},
		&models.ProductPricingSpecialDecrementMargin{},
		&models.ProductPricingRenewableProfitAdjustment{},
		&models.ProductPricingLapseMargin{},
		&models.ProductPricingRider{},
		&models.ProductPricingParameters{},
		&models.ProductPricingReinsurance{},
		&models.PricingYieldCurve{},
		&models.PricingRetrenchmentRate{},
		&models.AosVariableSet{},
		&models.AosVariable{},
		&models.CsmAosVariable{},
		&models.AOSStepResult{},
		&models.ChartAccountItem{},
		&models.JournalTransactions{},
		&models.CsmProjection{},
		&models.InitialRecognition{},
		&models.InsuranceRevenue{},
		&models.LiabilityMovement{},
		&models.PremiumEarningPattern{},
		&models.ReserveSummary{},
		&models.PAABuildUp{},
		&models.PAALapse{},
		&models.PAAEligibilityTestResult{},
		&models.LicExpectedSimulation{},
		&models.LicStandardisedResiduals{},
		&models.LicGeneratedRandomResiduals{},
		&models.LicRandomResiduals{},
		&models.LicBootStrappedCumulative{},
		&models.LicBootstrappedDevelopmentFactor{},
		&models.LicBootStrappedCumulativeProjection{},
		&models.LicBootStrappedIncremental{},
		&models.LicBootStrappedIncrementalProjection{},
		&models.LicBootStrappedIncrementalInflatedProjection{},
		&models.LicBootStrappedIncrementalInflatedDiscountedProjection{},
		&models.LicBootstrappedResults{},
		&models.LicBootstrappedResultSummary{},
		&models.IbnrFrequency{},
		&models.LicIndividualDevelopmentFactors{},
		&models.LicBiasAdjustmentFactor{},
		&models.LicMackModelCalculatedParameters{},
		&models.LicBiasAdjustedResiduals{},
		&models.LicMeanBiasAdjustedResiduals{},
		&models.LicLogNormalSigmas{},
		&models.LicLogNormalMeans{},
		&models.LicLogNormalStandardDeviations{},
		&models.LicPseudoRatios{},
		&models.LicMackModelSimulatedDevelopmentFactor{},
		&models.LicMackCumulativeProjection{},
		&models.LicMackSimulationResults{},
		&models.LicMackSimulationSummaryStats{},
		&models.MackIbnrFrequency{},
		&models.ProductPricingNewBusinessProfile{},
		&models.ProductPricingProfitMargin{},
		&models.AggregatedPricingPoint{},
		&models.PricingPolicyDemographic{},
		&models.Lic2Parameter{},
		&models.IFRS17AuditLog{},
	)

	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to migrate model point tables")
		return err
	}

	appLog.Info("Successfully migrated model point tables")
	return nil
}

// MigrateGMMTables migrates the GMM tables
func MigrateGMMTables() error {
	appLog.Info("Migrating GMM tables")
	err := DB.Migrator().AutoMigrate(
		&models.PaaPortfolio{},
		&models.ModifiedGMMModelPoint{},
		&models.ModifiedGMMParameter{},
		&models.ReinsuranceParameter{},
		&models.ModifiedGMMShockSetting{},
		&models.ModifiedGMMShock{},
		&models.MgmmRun{},
		&models.GMMRunSetting{},
		&models.PaaYieldCurve{},
		&models.PAAYearVersion{},
		&models.TransitionAdjustment{},
		&models.IFRS17Amendment{},
		&models.SCRRABridgeEntry{},
		&models.SARBCodeMapping{},
		&models.DeferredTaxEntry{},
	)

	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to migrate GMM tables")
		return err
	}

	appLog.Info("Successfully migrated GMM tables")
	return nil
}

// MigrateLICTables migrates the LIC tables
func MigrateLICTables() error {
	appLog.Info("Migrating LIC tables")
	err := DB.Migrator().AutoMigrate(
		&models.LicPortfolio{},
		&models.LICClaimsInput{},
		&models.LicVariableSet{},
		&models.LicVariable{},
		&models.LICParameter{},
		&models.LicCPI{},
		&models.IBNRRunSetting{},
		&models.LicModelPoint{},
		&models.LicTriangulation{},
		&models.LicTriangulationClaimCount{},
		&models.LicCumulativeTriangulation{},
		&models.LicCumulativeTriangulationClaimCount{},
		&models.LicCumulativeTriangulationAverageClaim{},
		&models.LicDevelopmentFactor{},
		&models.LicDevelopmentFactorClaimCount{},
		&models.LicDevelopmentFactorAverageClaim{},
		&models.LicCumulativeProjection{},
		&models.LicCumulativeProjectionClaimCount{},
		&models.LicCumulativeProjectionAverageClaim{},
		&models.LicCumulativeProjectionAveragetoTotalClaim{},
		&models.LicIncrementalProjectionAveragetoTotalClaim{},
		&models.LicIncrementalProjection{},
		&models.LicIncrementalInflated{},
		&models.LicDiscountedIncrementalInflated{},
		&models.LicIncrementalInflatedAveragetoTotalClaim{},
		&models.LicDiscountedIncrementalInflatedAveragetoTotalClaim{},
		&models.LicIbnrReserve{},
		&models.IbnrReserveReport{},
		&models.IbnrYieldCurve{},
		&models.LicRunSetting{},
		&models.IBNRShockSetting{},
		&models.IBNRShock{},
		&models.LicBuildupResult{},
		&models.LICEarnedPremium{},
		&models.ReinsuranceParameter{},
		&models.BelBuildupVariableSet{},
		&models.BelBuildupVariable{},
		&models.LicJournalTransactions{},
		&models.LicClaimsYearVersion{},
		&models.LicEarnedPremiumYearVersion{},
		&models.LicIbnrMethodAssignment{},
	)

	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to migrate LIC tables")
		return err
	}

	appLog.Info("Successfully migrated LIC tables")
	return nil
}

// MigrateExposureAnalysisTables migrates the exposure analysis tables
func MigrateExposureAnalysisTables() error {
	appLog.Info("Migrating exposure analysis tables")
	err := DB.Migrator().AutoMigrate(
		&models.ExpConfiguration{},
		&models.ExpExposureData{},
		&models.ExpActualData{},
		&models.ExpExpDataYearVersion{},
		&models.ExpActualDataYearVersion{},
		&models.ExpAnalysisRunSetting{},
		&models.ExposureModelPoint{},
		&models.ExpCrudeResult{},
		&models.ExpLapseCrudeResult{},
		&models.ExpAgeBand{},
		&models.ExpCurrentMortality{},
		&models.ExpCurrentLapse{},
		&models.ExpCurrentLapseYearVersion{},
		&models.ExpCurrentMortalityYearVersion{},
		&models.ExpRunGroup{},
		&models.TotalMortalityExpAnalysisResult{},
		&models.TotalLapseExpAnalysisResult{},
	)

	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to migrate exposure analysis tables")
		return err
	}

	appLog.Info("Successfully migrated exposure analysis tables")
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
		&models.GroupPricingClaimsExperience{},
		&models.MemberRatingResult{},
		&models.MemberRatingResultSummary{},
		&models.MovementMemberRatingResult{},
		&models.GlaRate{},
		&models.PhiRate{},
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
		&models.Escalations{},
		&models.FuneralAidsRate{},
		&models.FuneralRate{},
		&models.GeneralLoading{},
		&models.GlaAidsRate{},
		&models.RegionLoading{},
		&models.TieredIncomeReplacement{},
		&models.CustomTieredIncomeReplacement{},
		&models.DiscountAuthority{},
		&models.Restriction{},
		&models.ReinsuranceGlaRate{},
		&models.ReinsuranceCiRate{},
		&models.ReinsurancePtdRate{},
		&models.ReinsurancePhiRate{},
		&models.PremiumLoading{},
		&models.SchemeSizeLevel{},
		&models.TaxTable{},
		&models.GroupSchemeExposure{},
		&models.GroupPricingAgeBands{},
		&models.GroupBusinessBenefits{},
		&models.GroupSchemeClaim{},
		&models.GroupSchemeClaimAttachment{},
		&models.GroupSchemeClaimAssessment{},
		&models.GroupSchemeClaimCommunication{},
		&models.GroupSchemeClaimDecline{},
		&models.GroupSchemeClaimStatusAudit{},
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
		&models.MemberActivity{},
		&models.BordereauxTemplate{},
		&models.GeneratedBordereaux{},
		&models.BordereauxTimeline{},
		&models.BordereauxConfiguration{},
		&models.BordereauxConfirmation{},
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
	)

	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to migrate group pricing tables")
		return err
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

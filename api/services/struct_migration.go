package services

import (
	appLog "api/log"
	"api/models"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// GenerateMigrationForStruct generates SQL migration files for a given struct
func GenerateMigrationForStruct(structName, migrationName, dbType string) error {
	// Get the struct type using reflection
	structType, err := getStructType(structName)
	if err != nil {
		return err
	}

	// Generate SQL for the specified database type(s)
	if dbType == "all" {
		// Generate for all supported database types
		for _, db := range []string{"postgresql", "mysql", "mssql"} {
			if err := generateSQLForStruct(structType, migrationName, db); err != nil {
				return err
			}
		}
	} else {
		// Generate for a specific database type
		if err := generateSQLForStruct(structType, migrationName, dbType); err != nil {
			return err
		}
	}

	return nil
}

// getStructType returns the reflect.Type for a struct by name
func getStructType(structName string) (reflect.Type, error) {
	// Get the models package type
	modelsType := reflect.TypeOf(models.ProductMargins{})
	modelsPackage := modelsType.PkgPath()

	// Use a switch statement to match the struct name
	switch structName {

	// Base structs
	case "SystemLock":
		return reflect.TypeOf(models.SystemLock{}), nil
	case "AuditLog":
		return reflect.TypeOf(models.AuditLog{}), nil
	// Product related structs
	case "ProductMargins":
		return reflect.TypeOf(models.ProductMargins{}), nil
	case "Product":
		return reflect.TypeOf(models.Product{}), nil
	case "ProductFamily":
		return reflect.TypeOf(models.ProductFamily{}), nil
	case "GlobalTable":
		return reflect.TypeOf(models.GlobalTable{}), nil
	case "ProductTransitionState":
		return reflect.TypeOf(models.ProductTransitionState{}), nil
	case "ProductTransition":
		return reflect.TypeOf(models.ProductTransition{}), nil
	case "Fds":
		return reflect.TypeOf(models.Fds{}), nil
	case "ProductRatingFactor":
		return reflect.TypeOf(models.ProductRatingFactor{}), nil
	case "ProductModelpointVariable":
		return reflect.TypeOf(models.ProductModelpointVariable{}), nil
	case "ProductTable":
		return reflect.TypeOf(models.ProductTable{}), nil
	case "ProductChildSumAssured":
		return reflect.TypeOf(models.ProductChildSumAssured{}), nil
	case "ProductChildAdditionalSumAssured":
		return reflect.TypeOf(models.ProductChildAdditionalSumAssured{}), nil
	case "ProductAdditionalSumAssured":
		return reflect.TypeOf(models.ProductAdditionalSumAssured{}), nil
	case "ProductRider":
		return reflect.TypeOf(models.ProductRider{}), nil
	case "ProductClawback":
		return reflect.TypeOf(models.ProductClawback{}), nil
	case "ProductLapseMargin":
		return reflect.TypeOf(models.ProductLapseMargin{}), nil
	case "ProductAccidentalBenefitMultiplier":
		return reflect.TypeOf(models.ProductAccidentalBenefitMultiplier{}), nil
	case "ProductReinsurance":
		return reflect.TypeOf(models.ProductReinsurance{}), nil
	case "ProductFeatures":
		return reflect.TypeOf(models.ProductFeatures{}), nil
	case "ProductParameters":
		return reflect.TypeOf(models.ProductParameters{}), nil
	case "ProductNonLifeRating":
		return reflect.TypeOf(models.ProductNonLifeRating{}), nil
	case "ProductBenefitMultiplier":
		return reflect.TypeOf(models.ProductBenefitMultiplier{}), nil
	case "ProjectionJob":
		return reflect.TypeOf(models.ProjectionJob{}), nil
	case "Projection":
		return reflect.TypeOf(models.Projection{}), nil
	case "AggregatedProjection":
		return reflect.TypeOf(models.AggregatedProjection{}), nil
	case "ScopedAggregatedProjection":
		return reflect.TypeOf(models.ScopedAggregatedProjection{}), nil
	case "ProductCommissionStructure":
		return reflect.TypeOf(models.ProductCommissionStructure{}), nil
	case "ProductShock":
		return reflect.TypeOf(models.ProductShock{}), nil

	// IFRS17
	case "IFRS17AuditLog":
		return reflect.TypeOf(models.IFRS17AuditLog{}), nil
	case "IFRS17Amendment":
		return reflect.TypeOf(models.IFRS17Amendment{}), nil
	case "TransitionAdjustment":
		return reflect.TypeOf(models.TransitionAdjustment{}), nil

	// Pricing related structs
	case "ProductPricingTable":
		return reflect.TypeOf(models.ProductPricingTable{}), nil
	case "ProductPricingModelPoint":
		return reflect.TypeOf(models.ProductPricingModelPoint{}), nil
	case "ProductPricingMargins":
		return reflect.TypeOf(models.ProductPricingMargins{}), nil
	case "ProductPricingNewBusinessProfile":
		return reflect.TypeOf(models.ProductPricingNewBusinessProfile{}), nil
	case "ProductPricingProfitMargin":
		return reflect.TypeOf(models.ProductPricingProfitMargin{}), nil
	case "ProductPricingProductLevelEscalation":
		return reflect.TypeOf(models.ProductPricingProductLevelEscalation{}), nil
	case "ProductPricingParameters":
		return reflect.TypeOf(models.ProductPricingParameters{}), nil
	case "ProductPricingChildSumAssured":
		return reflect.TypeOf(models.ProductPricingChildSumAssured{}), nil
	case "ProductPricingChildAdditionalSumAssured":
		return reflect.TypeOf(models.ProductPricingChildAdditionalSumAssured{}), nil
	case "ProductPricingAdditionalSumAssured":
		return reflect.TypeOf(models.ProductPricingAdditionalSumAssured{}), nil
	case "ProductPricingRider":
		return reflect.TypeOf(models.ProductPricingRider{}), nil
	case "ProductPricingClawback":
		return reflect.TypeOf(models.ProductPricingClawback{}), nil
	case "ProductPricingLapseMargin":
		return reflect.TypeOf(models.ProductPricingLapseMargin{}), nil
	case "PricingYieldCurve":
		return reflect.TypeOf(models.PricingYieldCurve{}), nil
	case "PricingRetrenchmentRate":
		return reflect.TypeOf(models.PricingRetrenchmentRate{}), nil
	case "ProductPricingAccidentalBenefitMultiplier":
		return reflect.TypeOf(models.ProductPricingAccidentalBenefitMultiplier{}), nil
	case "ProductPricingReinsurance":
		return reflect.TypeOf(models.ProductPricingReinsurance{}), nil
	case "RiskAdjustmentFactor":
		return reflect.TypeOf(models.RiskAdjustmentFactor{}), nil
	case "ProductModelPointCount":
		return reflect.TypeOf(models.ProductModelPointCount{}), nil
	case "MemberIndicativeDataSet":
		return reflect.TypeOf(models.MemberIndicativeDataSet{}), nil
	case "PricingParameter":
		return reflect.TypeOf(models.PricingParameter{}), nil

	// User related structs
	case "AppUser":
		return reflect.TypeOf(models.AppUser{}), nil
	case "OrgUser":
		return reflect.TypeOf(models.OrgUser{}), nil
	case "Organisation":
		return reflect.TypeOf(models.Organisation{}), nil
	case "Activity":
		return reflect.TypeOf(models.Activity{}), nil

	// Task related structs
	case "Task":
		return reflect.TypeOf(models.Task{}), nil
	case "Comment":
		return reflect.TypeOf(models.Comment{}), nil

	// Other structs
	case "YieldCurve":
		return reflect.TypeOf(models.YieldCurve{}), nil
	case "RatingFactor":
		return reflect.TypeOf(models.RatingFactor{}), nil
	case "System":
		return reflect.TypeOf(models.System{}), nil
	case "RetrenchmentRate":
		return reflect.TypeOf(models.RetrenchmentRate{}), nil
	case "DisabilityIncidenceFactors":
		return reflect.TypeOf(models.DisabilityIncidenceFactors{}), nil
	case "ColumnNames":
		return reflect.TypeOf(models.ColumnNames{}), nil
	case "RiskAdjustmentDriver":
		return reflect.TypeOf(models.RiskAdjustmentDriver{}), nil
	case "CsmRun":
		return reflect.TypeOf(models.CsmRun{}), nil

	// IFRS17 related structs
	case "BaseAosVariable":
		return reflect.TypeOf(models.BaseAosVariable{}), nil
	case "AosVariable":
		return reflect.TypeOf(models.AosVariable{}), nil
	case "AosVariableSet":
		return reflect.TypeOf(models.AosVariableSet{}), nil
	case "CsmAosVariable":
		return reflect.TypeOf(models.CsmAosVariable{}), nil
	case "AOSStepResult":
		return reflect.TypeOf(models.AOSStepResult{}), nil
	case "PAAResult":
		return reflect.TypeOf(models.PAAResult{}), nil
	case "PAAEligibilityTestResult":
		return reflect.TypeOf(models.PAAEligibilityTestResult{}), nil
	case "IncomeStatementEntry":
		return reflect.TypeOf(models.IncomeStatementEntry{}), nil
	case "CSMModelPoint":
		return reflect.TypeOf(models.CSMModelPoint{}), nil
	case "ProductList":
		return reflect.TypeOf(models.ProductList{}), nil
	case "IFRS17List":
		return reflect.TypeOf(models.IFRS17List{}), nil
	case "GroupResults":
		return reflect.TypeOf(models.GroupResults{}), nil
	case "ProductStepResult":
		return reflect.TypeOf(models.ProductStepResult{}), nil
	case "LiabilityMovement":
		return reflect.TypeOf(models.LiabilityMovement{}), nil
	case "LiabilityMovementLine":
		return reflect.TypeOf(models.LiabilityMovementLine{}), nil
	case "InsuranceRevenue":
		return reflect.TypeOf(models.InsuranceRevenue{}), nil
	case "InitialRecognition":
		return reflect.TypeOf(models.InitialRecognition{}), nil
	case "CsmProjection":
		return reflect.TypeOf(models.CsmProjection{}), nil
	case "ChartAccountItem":
		return reflect.TypeOf(models.ChartAccountItem{}), nil
	case "FinanceVariables":
		return reflect.TypeOf(models.FinanceVariables{}), nil

		// PHI Valuation related structs
	case "PhiYieldCurve":
		return reflect.TypeOf(models.PhiYieldCurve{}), nil
	case "PhiParameter":
		return reflect.TypeOf(models.PhiParameter{}), nil
	case "PhiRecoveryRate":
		return reflect.TypeOf(models.PhiRecoveryRate{}), nil
	case "PhiMortality":
		return reflect.TypeOf(models.PhiMortality{}), nil
	case "PhiModelPoint":
		return reflect.TypeOf(models.PhiModelPoint{}), nil
	case "PhiShockSetting":
		return reflect.TypeOf(models.PhiShockSetting{}), nil
	case "PhiShock":
		return reflect.TypeOf(models.PhiShock{}), nil
	case "PhiRunConfig":
		return reflect.TypeOf(models.PhiRunConfig{}), nil

	// Group Pricing related structs
	case "PhiRate":
		return reflect.TypeOf(models.PhiRate{}), nil
	case "PtdRate":
		return reflect.TypeOf(models.PtdRate{}), nil
	case "MemberRatingResultSummary":
		return reflect.TypeOf(models.MemberRatingResultSummary{}), nil
	case "GroupPricingParameters":
		return reflect.TypeOf(models.GroupPricingParameters{}), nil
	case "GPricingMemberData":
		return reflect.TypeOf(models.GPricingMemberData{}), nil
	case "GPricingMemberDataInForce":
		return reflect.TypeOf(models.GPricingMemberDataInForce{}), nil
	case "CiRate":
		return reflect.TypeOf(models.CiRate{}), nil
	case "GlaRate":
		return reflect.TypeOf(models.GlaRate{}), nil
	case "TtdRate":
		return reflect.TypeOf(models.TtdRate{}), nil
	case "CalculationJob":
		return reflect.TypeOf(models.CalculationJob{}), nil
	case "GroupPricingQuote":
		return reflect.TypeOf(models.GroupPricingQuote{}), nil
	case "MovementMemberRatingResult":
		return reflect.TypeOf(models.MovementMemberRatingResult{}), nil
	case "MemberRatingResult":
		return reflect.TypeOf(models.MemberRatingResult{}), nil
	case "GroupScheme":
		return reflect.TypeOf(models.GroupScheme{}), nil
	case "Broker":
		return reflect.TypeOf(models.Broker{}), nil
	case "PremiumLoading":
		return reflect.TypeOf(models.PremiumLoading{}), nil
	case "SchemeSizeLevel":
		return reflect.TypeOf(models.SchemeSizeLevel{}), nil
	case "TaxTable":
		return reflect.TypeOf(models.TaxTable{}), nil
	case "SchemeCategory":
		return reflect.TypeOf(models.SchemeCategory{}), nil
	case "SchemeType":
		return reflect.TypeOf(models.SchemeType{}), nil
	case "Bordereaux":
		return reflect.TypeOf(models.Bordereaux{}), nil
	case "MemberPremiumSchedule":
		return reflect.TypeOf(models.MemberPremiumSchedule{}), nil
	case "GroupSchemeClaim":
		return reflect.TypeOf(models.GroupSchemeClaim{}), nil
	case "GroupPricingInsurerDetail":
		return reflect.TypeOf(models.GroupPricingInsurerDetail{}), nil
	case "Loadings":
		return reflect.TypeOf(models.Loadings{}), nil
	case "PhiRunParameters":
		return reflect.TypeOf(models.PhiRunParameters{}), nil
	case "PhiReinsurance" +
		"":
		return reflect.TypeOf(models.PhiReinsurance{}), nil
	case "PhiProjections":
		return reflect.TypeOf(models.PhiProjections{}), nil
	case "PhiAggregatedProjections":
		return reflect.TypeOf(models.PhiAggregatedProjections{}), nil
	case "PhiScopedAggregatedProjections":
		return reflect.TypeOf(models.PhiScopedAggregatedProjections{}), nil
	case "RunPhiJob":
		return reflect.TypeOf(models.RunPhiJob{}), nil

	case "FuneralAidsRate":
		return reflect.TypeOf(models.FuneralAidsRate{}), nil
	case "GlaAidsRate":
		return reflect.TypeOf(models.GlaAidsRate{}), nil
	case "TieredIncomeReplacement":
		return reflect.TypeOf(models.TieredIncomeReplacement{}), nil
	case "DiscountAuthority":
		return reflect.TypeOf(models.DiscountAuthority{}), nil
	case "Restriction":
		return reflect.TypeOf(models.Restriction{}), nil
	case "ReinsuranceGlaRate":
		return reflect.TypeOf(models.ReinsuranceGlaRate{}), nil
	case "ReinsuranceCiRate":
		return reflect.TypeOf(models.ReinsuranceCiRate{}), nil
	case "ReinsurancePtdRate":
		return reflect.TypeOf(models.ReinsurancePtdRate{}), nil
	case "ReinsurancePhiRate":
		return reflect.TypeOf(models.ReinsurancePhiRate{}), nil
	case "FuneralRate":
		return reflect.TypeOf(models.FuneralRate{}), nil
	case "GeneralLoading":
		return reflect.TypeOf(models.GeneralLoading{}), nil
	case "RegionLoading":
		return reflect.TypeOf(models.RegionLoading{}), nil
	case "PremiumLoadings":
		return reflect.TypeOf(models.PremiumLoading{}), nil

	case "Escalations":
		return reflect.TypeOf(models.Escalations{}), nil
	case "GroupPricingClaimsExperience":
		return reflect.TypeOf(models.GroupPricingClaimsExperience{}), nil
	case "GroupRiskQuoteStats":
		return reflect.TypeOf(models.GroupRiskQuoteStats{}), nil
	case "GroupSchemeClaimAttachment":
		return reflect.TypeOf(models.GroupSchemeClaimAttachment{}), nil
	case "GroupSchemeClaimAssessment":
		return reflect.TypeOf(models.GroupSchemeClaimAssessment{}), nil
	case "BordereauxTemplate":
		return reflect.TypeOf(models.BordereauxTemplate{}), nil
	case "GroupSchemeClaimCommunication":
		return reflect.TypeOf(models.GroupSchemeClaimCommunication{}), nil
	case "GroupSchemeClaimStatusAudit":
		return reflect.TypeOf(models.GroupSchemeClaimStatusAudit{}), nil
	case "GroupSchemeClaimDecline":
		return reflect.TypeOf(models.GroupSchemeClaimDecline{}), nil
	case "MemberActivity":
		return reflect.TypeOf(models.MemberActivity{}), nil
	case "GeneratedBordereaux":
		return reflect.TypeOf(models.GeneratedBordereaux{}), nil
	case "BordereauxConfiguration":
		return reflect.TypeOf(models.BordereauxConfiguration{}), nil
	case "PremiumBordereauxData":
		return reflect.TypeOf(models.PremiumBordereauxData{}), nil
	case "MemberBordereauxData":
		return reflect.TypeOf(models.MemberBordereauxData{}), nil
	case "ClaimBordereauxData":
		return reflect.TypeOf(models.ClaimBordereauxData{}), nil
	case "BordereauxTimeline":
		return reflect.TypeOf(models.BordereauxTimeline{}), nil
	case "BordereauxConfirmation":
		return reflect.TypeOf(models.BordereauxConfirmation{}), nil
	case "BordereauxReconciliationResult":
		return reflect.TypeOf(models.BordereauxReconciliationResult{}), nil
	case "BordereauxConfirmationRecord":
		return reflect.TypeOf(models.BordereauxConfirmationRecord{}), nil
	case "BenefitDocumentType":
		return reflect.TypeOf(models.BenefitDocumentType{}), nil
	case "SchemeCategoryMaster":
		return reflect.TypeOf(models.SchemeCategoryMaster{}), nil
	case "GroupBenefitMapper":
		return reflect.TypeOf(models.GroupBenefitMapper{}), nil
	case "RegisterDiffSnapshot":
		return reflect.TypeOf(models.RegisterDiffSnapshot{}), nil
	case "MemberBeneficiary":
		return reflect.TypeOf(models.MemberBeneficiary{}), nil
	case "OccupationClass":
		return reflect.TypeOf(models.OccupationClass{}), nil
	case "ReinsuranceTreaty":
		return reflect.TypeOf(models.ReinsuranceTreaty{}), nil
	case "TreatySchemeLink":
		return reflect.TypeOf(models.TreatySchemeLink{}), nil
	case "RIBordereauxRun":
		return reflect.TypeOf(models.RIBordereauxRun{}), nil
	case "LargeClaimNotice":
		return reflect.TypeOf(models.LargeClaimNotice{}), nil
	case "TechnicalAccount":
		return reflect.TypeOf(models.TechnicalAccount{}), nil
	case "SettlementPayment":
		return reflect.TypeOf(models.SettlementPayment{}), nil
	case "RIBordereauxMemberRow":
		return reflect.TypeOf(models.RIBordereauxMemberRow{}), nil
	case "RIBordereauxClaimsRow":
		return reflect.TypeOf(models.RIBordereauxClaimsRow{}), nil
	case "RIValidationResult":
		return reflect.TypeOf(models.RIValidationResult{}), nil
	case "Reinsurer":
		return reflect.TypeOf(models.Reinsurer{}), nil
	case "ClaimPaymentSchedule":
		return reflect.TypeOf(models.ClaimPaymentSchedule{}), nil
	case "ClaimPaymentScheduleItem":
		return reflect.TypeOf(models.ClaimPaymentScheduleItem{}), nil
	case "ClaimPaymentProof":
		return reflect.TypeOf(models.ClaimPaymentProof{}), nil
	case "ACBBankProfile":
		return reflect.TypeOf(models.ACBBankProfile{}), nil
	case "ACBFileRecord":
		return reflect.TypeOf(models.ACBFileRecord{}), nil
	case "ACBReconciliationResult":
		return reflect.TypeOf(models.ACBReconciliationResult{}), nil
	case "Notification":
		return reflect.TypeOf(models.Notification{}), nil
	case "Conversation":
		return reflect.TypeOf(models.Conversation{}), nil
	case "ConversationParticipant":
		return reflect.TypeOf(models.ConversationParticipant{}), nil
	case "ConversationMessage":
		return reflect.TypeOf(models.ConversationMessage{}), nil
	case "MonthlyQuoteTrend":
		return reflect.TypeOf(models.MonthlyQuoteTrend{}), nil
	case "QuoteFunnelStage":
		return reflect.TypeOf(models.QuoteFunnelStage{}), nil
	case "BrokerMetric":
		return reflect.TypeOf(models.BrokerMetric{}), nil
	case "DashboardPricingMetrics":
		return reflect.TypeOf(models.DashboardPricingMetrics{}), nil
	case "WinProbabilityModel":
		return reflect.TypeOf(models.WinProbabilityModel{}), nil
	case "QuoteWinProbability":
		return reflect.TypeOf(models.QuoteWinProbability{}), nil
	case "CustomTieredIncomeReplacement":
		return reflect.TypeOf(models.CustomTieredIncomeReplacement{}), nil
	case "GPTableStat":
		return reflect.TypeOf(models.GPTableStat{}), nil
	// Add more cases for other structs as needed

	// IBNR
	case "LicIbnrReserve":
		return reflect.TypeOf(models.LicIbnrReserve{}), nil
	case "IbnrReserveReport":
		return reflect.TypeOf(models.IbnrReserveReport{}), nil
	case "LicIbnrMethodAssignment":
		return reflect.TypeOf(models.LicIbnrMethodAssignment{}), nil
	case "LicBootStrappedIncremental":
		return reflect.TypeOf(models.LicBootStrappedIncremental{}), nil

	// Lic related structs
	case "LicBuildupResult":
		return reflect.TypeOf(models.LicBuildupResult{}), nil
	case "Lic2Parameter":
		return reflect.TypeOf(models.Lic2Parameter{}), nil

	// Group Premium Tables
	case "ContributionConfig":
		return reflect.TypeOf(models.ContributionConfig{}), nil
	case "PremiumSchedule":
		return reflect.TypeOf(models.PremiumSchedule{}), nil
	case "ScheduleMemberRow":
		return reflect.TypeOf(models.ScheduleMemberRow{}), nil
	case "Invoice":
		return reflect.TypeOf(models.Invoice{}), nil
	case "InvoiceAdjustment":
		return reflect.TypeOf(models.InvoiceAdjustment{}), nil
	case "Payment":
		return reflect.TypeOf(models.Payment{}), nil
	case "ArrearsHistory":
		return reflect.TypeOf(models.ArrearsHistory{}), nil
	case "PaymentPlan":
		return reflect.TypeOf(models.PaymentPlan{}), nil
	case "PaymentPlanInstalment":
		return reflect.TypeOf(models.PaymentPlanInstalment{}), nil
	case "EmployerSubmission":
		return reflect.TypeOf(models.EmployerSubmission{}), nil
	case "EmployerSubmissionRecord":
		return reflect.TypeOf(models.EmployerSubmissionRecord{}), nil
	case "SubmissionDeltaRecord":
		return reflect.TypeOf(models.SubmissionDeltaRecord{}), nil
	case "BordereauxDeadline":
		return reflect.TypeOf(models.BordereauxDeadline{}), nil
	case "ReinsurerAcceptance":
		return reflect.TypeOf(models.ReinsurerAcceptance{}), nil
	case "ReinsurerRecovery":
		return reflect.TypeOf(models.ReinsurerRecovery{}), nil
	case "ClaimNotificationLog":
		return reflect.TypeOf(models.ClaimNotificationLog{}), nil
	case "NewJoinerDetail":
		return reflect.TypeOf(models.NewJoinerDetail{}), nil
	case "ReconciliationRun":
		return reflect.TypeOf(models.ReconciliationRun{}), nil
	case "PaymentAllocation":
		return reflect.TypeOf(models.PaymentAllocation{}), nil
	case "ReconciliationItem":
		return reflect.TypeOf(models.ReconciliationItem{}), nil
	case "MatchingRule":
		return reflect.TypeOf(models.MatchingRule{}), nil
	case "OnRiskLetter":
		return reflect.TypeOf(models.OnRiskLetter{}), nil

	// --- Base tables ---
	case "BaseAssumptionVariable":
		return reflect.TypeOf(models.BaseAssumptionVariable{}), nil
	case "BaseFeature":
		return reflect.TypeOf(models.BaseFeature{}), nil
	case "ModelPointVariable":
		return reflect.TypeOf(models.ModelPointVariable{}), nil
	case "MarkovState":
		return reflect.TypeOf(models.MarkovState{}), nil
	case "BaseMortalityBand":
		return reflect.TypeOf(models.BaseMortalityBand{}), nil
	case "LicBaseVariable":
		return reflect.TypeOf(models.LicBaseVariable{}), nil
	case "UserToken":
		return reflect.TypeOf(models.UserToken{}), nil
	case "ConsolidateResult":
		return reflect.TypeOf(models.ConsolidateResult{}), nil
	case "CumulativeConsolidatedResult":
		return reflect.TypeOf(models.CumulativeConsolidatedResult{}), nil
	case "AnnualConsolidatedResult":
		return reflect.TypeOf(models.AnnualConsolidatedResult{}), nil
	case "BelBuildupBaseVariable":
		return reflect.TypeOf(models.BelBuildupBaseVariable{}), nil
	case "GroupPricingAgeBands":
		return reflect.TypeOf(models.GroupPricingAgeBands{}), nil
	case "GroupBusinessBenefits":
		return reflect.TypeOf(models.GroupBusinessBenefits{}), nil
	case "GPPermission":
		return reflect.TypeOf(models.GPPermission{}), nil
	case "TransitionState":
		return reflect.TypeOf(models.TransitionState{}), nil

	// --- Product model tables ---
	case "ShockSetting":
		return reflect.TypeOf(models.ShockSetting{}), nil
	case "ProductPricingShock":
		return reflect.TypeOf(models.ProductPricingShock{}), nil
	case "ManualScopedAggregatedProjection":
		return reflect.TypeOf(models.ManualScopedAggregatedProjection{}), nil
	case "LICAggregatedProjections":
		return reflect.TypeOf(models.LICAggregatedProjections{}), nil
	case "JobsTemplate":
		return reflect.TypeOf(models.JobsTemplate{}), nil
	case "JobTemplateContent":
		return reflect.TypeOf(models.JobTemplateContent{}), nil
	case "ProductSpecialDecrementMargin":
		return reflect.TypeOf(models.ProductSpecialDecrementMargin{}), nil
	case "ProductRenewableProfitAdjustment":
		return reflect.TypeOf(models.ProductRenewableProfitAdjustment{}), nil
	case "ProductUnitFundCharge":
		return reflect.TypeOf(models.ProductUnitFundCharge{}), nil
	case "ProductInvestmentReturn":
		return reflect.TypeOf(models.ProductInvestmentReturn{}), nil
	case "ProductFundAssetDistribution":
		return reflect.TypeOf(models.ProductFundAssetDistribution{}), nil
	case "ProductMaturityPattern":
		return reflect.TypeOf(models.ProductMaturityPattern{}), nil
	case "ProductSurrenderValueCoefficient":
		return reflect.TypeOf(models.ProductSurrenderValueCoefficient{}), nil
	case "IBNRAverageClaimAmount":
		return reflect.TypeOf(models.IBNRAverageClaimAmount{}), nil
	case "IBNRClaimHistorySummary":
		return reflect.TypeOf(models.IBNRClaimHistorySummary{}), nil
	case "LICClaimsAnalysisOfChange":
		return reflect.TypeOf(models.LICClaimsAnalysisOfChange{}), nil
	case "IBNRPaidVsOutstandingClaims":
		return reflect.TypeOf(models.IBNRPaidVsOutstandingClaims{}), nil
	case "IBNRProportionOutstandingClaims":
		return reflect.TypeOf(models.IBNRProportionOutstandingClaims{}), nil
	case "IBNRIncurredClaims":
		return reflect.TypeOf(models.IBNRIncurredClaims{}), nil
	case "AggregatedModifiedGMMProjection":
		return reflect.TypeOf(models.AggregatedModifiedGMMProjection{}), nil
	case "ModifiedGMMProjection":
		return reflect.TypeOf(models.ModifiedGMMProjection{}), nil
	case "ModifiedGMMScopedAggregation":
		return reflect.TypeOf(models.ModifiedGMMScopedAggregation{}), nil
	case "JobProduct":
		return reflect.TypeOf(models.JobProduct{}), nil
	case "JobProductRunError":
		return reflect.TypeOf(models.JobProductRunError{}), nil
	case "CachedReserveResults":
		return reflect.TypeOf(models.CachedReserveResults{}), nil
	case "RunParameters":
		return reflect.TypeOf(models.RunParameters{}), nil
	case "ModelPointPricing":
		return reflect.TypeOf(models.ModelPointPricing{}), nil
	case "Profitability":
		return reflect.TypeOf(models.Profitability{}), nil
	case "PricingRun":
		return reflect.TypeOf(models.PricingRun{}), nil
	case "PricingConfig":
		return reflect.TypeOf(models.PricingConfig{}), nil
	case "RiskDriver":
		return reflect.TypeOf(models.RiskDriver{}), nil
	case "RAConfiguration":
		return reflect.TypeOf(models.RAConfiguration{}), nil
	case "AOSConfiguration":
		return reflect.TypeOf(models.AOSConfiguration{}), nil
	case "AggregatedVariableGroup":
		return reflect.TypeOf(models.AggregatedVariableGroup{}), nil

	// --- Escalation tables ---
	case "BalanceSheetRecord":
		return reflect.TypeOf(models.BalanceSheetRecord{}), nil
	case "PAAFinance":
		return reflect.TypeOf(models.PAAFinance{}), nil
	case "PricingPoint":
		return reflect.TypeOf(models.PricingPoint{}), nil
	case "ProductModelPointVariableStats":
		return reflect.TypeOf(models.ProductModelPointVariableStats{}), nil

	// --- Model point tables ---
	case "PaaModelPointVariableStats":
		return reflect.TypeOf(models.PaaModelPointVariableStats{}), nil
	case "ProductPricingSpecialDecrementMargin":
		return reflect.TypeOf(models.ProductPricingSpecialDecrementMargin{}), nil
	case "ProductPricingRenewableProfitAdjustment":
		return reflect.TypeOf(models.ProductPricingRenewableProfitAdjustment{}), nil
	case "JournalTransactions":
		return reflect.TypeOf(models.JournalTransactions{}), nil
	case "PremiumEarningPattern":
		return reflect.TypeOf(models.PremiumEarningPattern{}), nil
	case "ReserveSummary":
		return reflect.TypeOf(models.ReserveSummary{}), nil
	case "PAABuildUp":
		return reflect.TypeOf(models.PAABuildUp{}), nil
	case "PAALapse":
		return reflect.TypeOf(models.PAALapse{}), nil
	case "LicExpectedSimulation":
		return reflect.TypeOf(models.LicExpectedSimulation{}), nil
	case "LicStandardisedResiduals":
		return reflect.TypeOf(models.LicStandardisedResiduals{}), nil
	case "LicGeneratedRandomResiduals":
		return reflect.TypeOf(models.LicGeneratedRandomResiduals{}), nil
	case "LicRandomResiduals":
		return reflect.TypeOf(models.LicRandomResiduals{}), nil
	case "LicBootStrappedCumulative":
		return reflect.TypeOf(models.LicBootStrappedCumulative{}), nil
	case "LicBootstrappedDevelopmentFactor":
		return reflect.TypeOf(models.LicBootstrappedDevelopmentFactor{}), nil
	case "LicBootStrappedCumulativeProjection":
		return reflect.TypeOf(models.LicBootStrappedCumulativeProjection{}), nil
	case "LicBootStrappedIncrementalProjection":
		return reflect.TypeOf(models.LicBootStrappedIncrementalProjection{}), nil
	case "LicBootStrappedIncrementalInflatedProjection":
		return reflect.TypeOf(models.LicBootStrappedIncrementalInflatedProjection{}), nil
	case "LicBootStrappedIncrementalInflatedDiscountedProjection":
		return reflect.TypeOf(models.LicBootStrappedIncrementalInflatedDiscountedProjection{}), nil
	case "LicBootstrappedResults":
		return reflect.TypeOf(models.LicBootstrappedResults{}), nil
	case "LicBootstrappedResultSummary":
		return reflect.TypeOf(models.LicBootstrappedResultSummary{}), nil
	case "IbnrFrequency":
		return reflect.TypeOf(models.IbnrFrequency{}), nil
	case "LicIndividualDevelopmentFactors":
		return reflect.TypeOf(models.LicIndividualDevelopmentFactors{}), nil
	case "LicBiasAdjustmentFactor":
		return reflect.TypeOf(models.LicBiasAdjustmentFactor{}), nil
	case "LicMackModelCalculatedParameters":
		return reflect.TypeOf(models.LicMackModelCalculatedParameters{}), nil
	case "LicBiasAdjustedResiduals":
		return reflect.TypeOf(models.LicBiasAdjustedResiduals{}), nil
	case "LicMeanBiasAdjustedResiduals":
		return reflect.TypeOf(models.LicMeanBiasAdjustedResiduals{}), nil
	case "LicLogNormalSigmas":
		return reflect.TypeOf(models.LicLogNormalSigmas{}), nil
	case "LicLogNormalMeans":
		return reflect.TypeOf(models.LicLogNormalMeans{}), nil
	case "LicLogNormalStandardDeviations":
		return reflect.TypeOf(models.LicLogNormalStandardDeviations{}), nil
	case "LicPseudoRatios":
		return reflect.TypeOf(models.LicPseudoRatios{}), nil
	case "LicMackModelSimulatedDevelopmentFactor":
		return reflect.TypeOf(models.LicMackModelSimulatedDevelopmentFactor{}), nil
	case "LicMackCumulativeProjection":
		return reflect.TypeOf(models.LicMackCumulativeProjection{}), nil
	case "LicMackSimulationResults":
		return reflect.TypeOf(models.LicMackSimulationResults{}), nil
	case "LicMackSimulationSummaryStats":
		return reflect.TypeOf(models.LicMackSimulationSummaryStats{}), nil
	case "MackIbnrFrequency":
		return reflect.TypeOf(models.MackIbnrFrequency{}), nil
	case "AggregatedPricingPoint":
		return reflect.TypeOf(models.AggregatedPricingPoint{}), nil
	case "PricingPolicyDemographic":
		return reflect.TypeOf(models.PricingPolicyDemographic{}), nil

	// --- GMM tables ---
	case "PaaPortfolio":
		return reflect.TypeOf(models.PaaPortfolio{}), nil
	case "ModifiedGMMModelPoint":
		return reflect.TypeOf(models.ModifiedGMMModelPoint{}), nil
	case "ModifiedGMMParameter":
		return reflect.TypeOf(models.ModifiedGMMParameter{}), nil
	case "ReinsuranceParameter":
		return reflect.TypeOf(models.ReinsuranceParameter{}), nil
	case "ModifiedGMMShockSetting":
		return reflect.TypeOf(models.ModifiedGMMShockSetting{}), nil
	case "ModifiedGMMShock":
		return reflect.TypeOf(models.ModifiedGMMShock{}), nil
	case "MgmmRun":
		return reflect.TypeOf(models.MgmmRun{}), nil
	case "GMMRunSetting":
		return reflect.TypeOf(models.GMMRunSetting{}), nil
	case "PaaYieldCurve":
		return reflect.TypeOf(models.PaaYieldCurve{}), nil
	case "PAAYearVersion":
		return reflect.TypeOf(models.PAAYearVersion{}), nil
	case "SCRRABridgeEntry":
		return reflect.TypeOf(models.SCRRABridgeEntry{}), nil
	case "SARBCodeMapping":
		return reflect.TypeOf(models.SARBCodeMapping{}), nil
	case "DeferredTaxEntry":
		return reflect.TypeOf(models.DeferredTaxEntry{}), nil

	// --- LIC tables ---
	case "LicPortfolio":
		return reflect.TypeOf(models.LicPortfolio{}), nil
	case "LICClaimsInput":
		return reflect.TypeOf(models.LICClaimsInput{}), nil
	case "LicVariableSet":
		return reflect.TypeOf(models.LicVariableSet{}), nil
	case "LicVariable":
		return reflect.TypeOf(models.LicVariable{}), nil
	case "LICParameter":
		return reflect.TypeOf(models.LICParameter{}), nil
	case "LicCPI":
		return reflect.TypeOf(models.LicCPI{}), nil
	case "IBNRRunSetting":
		return reflect.TypeOf(models.IBNRRunSetting{}), nil
	case "LicModelPoint":
		return reflect.TypeOf(models.LicModelPoint{}), nil
	case "LicTriangulation":
		return reflect.TypeOf(models.LicTriangulation{}), nil
	case "LicTriangulationClaimCount":
		return reflect.TypeOf(models.LicTriangulationClaimCount{}), nil
	case "LicCumulativeTriangulation":
		return reflect.TypeOf(models.LicCumulativeTriangulation{}), nil
	case "LicCumulativeTriangulationClaimCount":
		return reflect.TypeOf(models.LicCumulativeTriangulationClaimCount{}), nil
	case "LicCumulativeTriangulationAverageClaim":
		return reflect.TypeOf(models.LicCumulativeTriangulationAverageClaim{}), nil
	case "LicDevelopmentFactor":
		return reflect.TypeOf(models.LicDevelopmentFactor{}), nil
	case "LicDevelopmentFactorClaimCount":
		return reflect.TypeOf(models.LicDevelopmentFactorClaimCount{}), nil
	case "LicDevelopmentFactorAverageClaim":
		return reflect.TypeOf(models.LicDevelopmentFactorAverageClaim{}), nil
	case "LicCumulativeProjection":
		return reflect.TypeOf(models.LicCumulativeProjection{}), nil
	case "LicCumulativeProjectionClaimCount":
		return reflect.TypeOf(models.LicCumulativeProjectionClaimCount{}), nil
	case "LicCumulativeProjectionAverageClaim":
		return reflect.TypeOf(models.LicCumulativeProjectionAverageClaim{}), nil
	case "LicCumulativeProjectionAveragetoTotalClaim":
		return reflect.TypeOf(models.LicCumulativeProjectionAveragetoTotalClaim{}), nil
	case "LicIncrementalProjectionAveragetoTotalClaim":
		return reflect.TypeOf(models.LicIncrementalProjectionAveragetoTotalClaim{}), nil
	case "LicIncrementalProjection":
		return reflect.TypeOf(models.LicIncrementalProjection{}), nil
	case "LicIncrementalInflated":
		return reflect.TypeOf(models.LicIncrementalInflated{}), nil
	case "LicDiscountedIncrementalInflated":
		return reflect.TypeOf(models.LicDiscountedIncrementalInflated{}), nil
	case "LicIncrementalInflatedAveragetoTotalClaim":
		return reflect.TypeOf(models.LicIncrementalInflatedAveragetoTotalClaim{}), nil
	case "LicDiscountedIncrementalInflatedAveragetoTotalClaim":
		return reflect.TypeOf(models.LicDiscountedIncrementalInflatedAveragetoTotalClaim{}), nil
	case "IbnrYieldCurve":
		return reflect.TypeOf(models.IbnrYieldCurve{}), nil
	case "LicRunSetting":
		return reflect.TypeOf(models.LicRunSetting{}), nil
	case "IBNRShockSetting":
		return reflect.TypeOf(models.IBNRShockSetting{}), nil
	case "IBNRShock":
		return reflect.TypeOf(models.IBNRShock{}), nil
	case "BelBuildupVariableSet":
		return reflect.TypeOf(models.BelBuildupVariableSet{}), nil
	case "BelBuildupVariable":
		return reflect.TypeOf(models.BelBuildupVariable{}), nil
	case "LicJournalTransactions":
		return reflect.TypeOf(models.LicJournalTransactions{}), nil
	case "LicClaimsYearVersion":
		return reflect.TypeOf(models.LicClaimsYearVersion{}), nil
	case "LicEarnedPremiumYearVersion":
		return reflect.TypeOf(models.LicEarnedPremiumYearVersion{}), nil
	case "LICEarnedPremium":
		return reflect.TypeOf(models.LICEarnedPremium{}), nil

	// --- Exposure analysis tables ---
	case "ExpConfiguration":
		return reflect.TypeOf(models.ExpConfiguration{}), nil
	case "ExpExposureData":
		return reflect.TypeOf(models.ExpExposureData{}), nil
	case "ExpActualData":
		return reflect.TypeOf(models.ExpActualData{}), nil
	case "ExpExpDataYearVersion":
		return reflect.TypeOf(models.ExpExpDataYearVersion{}), nil
	case "ExpActualDataYearVersion":
		return reflect.TypeOf(models.ExpActualDataYearVersion{}), nil
	case "ExpAnalysisRunSetting":
		return reflect.TypeOf(models.ExpAnalysisRunSetting{}), nil
	case "ExposureModelPoint":
		return reflect.TypeOf(models.ExposureModelPoint{}), nil
	case "ExpCrudeResult":
		return reflect.TypeOf(models.ExpCrudeResult{}), nil
	case "ExpLapseCrudeResult":
		return reflect.TypeOf(models.ExpLapseCrudeResult{}), nil
	case "ExpAgeBand":
		return reflect.TypeOf(models.ExpAgeBand{}), nil
	case "ExpCurrentMortality":
		return reflect.TypeOf(models.ExpCurrentMortality{}), nil
	case "ExpCurrentLapse":
		return reflect.TypeOf(models.ExpCurrentLapse{}), nil
	case "ExpCurrentLapseYearVersion":
		return reflect.TypeOf(models.ExpCurrentLapseYearVersion{}), nil
	case "ExpCurrentMortalityYearVersion":
		return reflect.TypeOf(models.ExpCurrentMortalityYearVersion{}), nil
	case "ExpRunGroup":
		return reflect.TypeOf(models.ExpRunGroup{}), nil
	case "TotalMortalityExpAnalysisResult":
		return reflect.TypeOf(models.TotalMortalityExpAnalysisResult{}), nil
	case "TotalLapseExpAnalysisResult":
		return reflect.TypeOf(models.TotalLapseExpAnalysisResult{}), nil

	// --- Group pricing tables (missing) ---
	case "ChildMortality":
		return reflect.TypeOf(models.ChildMortality{}), nil
	case "IndustryLoading":
		return reflect.TypeOf(models.IndustryLoading{}), nil
	case "FuneralParameters":
		return reflect.TypeOf(models.FuneralParameters{}), nil
	case "GroupPricingReinsuranceStructure":
		return reflect.TypeOf(models.GroupPricingReinsuranceStructure{}), nil
	case "IncomeLevel":
		return reflect.TypeOf(models.IncomeLevel{}), nil
	case "EducatorBenefitStructure":
		return reflect.TypeOf(models.EducatorBenefitStructure{}), nil
	case "EducatorRate":
		return reflect.TypeOf(models.EducatorRate{}), nil
	case "AccidentalTtdRate":
		return reflect.TypeOf(models.AccidentalTtdRate{}), nil
	case "GroupSchemeExposure":
		return reflect.TypeOf(models.GroupSchemeExposure{}), nil
	case "GroupSchemeStatusAudit":
		return reflect.TypeOf(models.GroupSchemeStatusAudit{}), nil
	case "HistoricalCredibilityData":
		return reflect.TypeOf(models.HistoricalCredibilityData{}), nil
	case "Beneficiary":
		return reflect.TypeOf(models.Beneficiary{}), nil

	// --- Group pricing user tables ---
	case "GPUser":
		return reflect.TypeOf(models.GPUser{}), nil
	case "GPUserRole":
		return reflect.TypeOf(models.GPUserRole{}), nil

	// -- Valuation User Management
	case "ValUserRole":
		return reflect.TypeOf(models.ValUserRole{}), nil
	case "ValPermission":
		return reflect.TypeOf(models.ValPermission{}), nil

	default:
		// If we couldn't find the struct, return an error
		return nil, fmt.Errorf("struct %s not found in package %s", structName, modelsPackage)
	}
}

// generateSQLForStruct generates SQL for a struct for a specific database type
func generateSQLForStruct(structType reflect.Type, migrationName, dbType string) error {
	// Create migration directory if it doesn't exist
	migrationDir := filepath.Join("migrations", dbType)
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		return err
	}

	// Generate version based on timestamp
	version := time.Now().Format("20060102150405")

	// Create filename
	filename := fmt.Sprintf("%s_%s.sql", version, migrationName)
	filePath := filepath.Join(migrationDir, filename)

	// Generate SQL based on struct fields and database type
	sql := generateSQL(structType, dbType)

	// Write SQL to file
	if err := os.WriteFile(filePath, []byte(sql), 0644); err != nil {
		return err
	}

	appLog.WithFields(map[string]interface{}{
		"db_type": dbType,
		"file":    filePath,
	}).Info("Generated migration file for struct")

	return nil
}

// generateSQL generates SQL for a struct based on the database type
func generateSQL(structType reflect.Type, dbType string) string {
	tableName := getTableName(structType)
	var sql strings.Builder

	sql.WriteString(fmt.Sprintf("-- Migration for struct: %s\n\n", structType.Name()))
	sql.WriteString(fmt.Sprintf("-- Table: %s\n\n", tableName))

	// First, ensure the table exists
	sql.WriteString("-- Ensure table exists\n")
	switch dbType {
	case "postgresql":
		sql.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", tableName))
		// Add a placeholder ID column if the table is being created
		sql.WriteString("    id SERIAL PRIMARY KEY\n")
		sql.WriteString(");\n\n")
	case "mysql":
		sql.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", tableName))
		// Add a placeholder ID column if the table is being created
		sql.WriteString("    id INT AUTO_INCREMENT PRIMARY KEY\n")
		sql.WriteString(");\n\n")
	case "mssql":
		sql.WriteString(fmt.Sprintf("IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = '%s')\n", tableName))
		sql.WriteString("BEGIN\n")
		sql.WriteString(fmt.Sprintf("    CREATE TABLE %s (\n", tableName))
		// Add a placeholder ID column if the table is being created
		sql.WriteString("        id INT IDENTITY(1,1) PRIMARY KEY\n")
		sql.WriteString("    );\n")
		sql.WriteString("END;\n\n")
	}

	// Collect fields including embedded structs with prefixes
	fields := collectDBFields(structType, "")

	// Generate ALTER TABLE statements for each collected field
	for _, f := range fields {
		// Skip the top-level ID field if it's already included in the CREATE TABLE statement
		if f.Field.Name == "ID" && f.Prefix == "" {
			baseName := getColumnName(f.Field)
			if baseName == "id" {
				continue
			}
		}

		baseColumnName := getColumnName(f.Field)
		columnName := f.Prefix + baseColumnName
		columnType := getColumnType(f.Field, dbType)

		// Generate ALTER TABLE statement
		sql.WriteString(fmt.Sprintf("-- Add or modify column for field: %s\n", f.Field.Name))

		switch dbType {
		case "postgresql":
			// PostgreSQL supports IF NOT EXISTS for ADD COLUMN
			sql.WriteString(fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS %s %s;\n",
				tableName, columnName, columnType))

			// Also add an ALTER COLUMN statement to modify the column type if it exists
			sql.WriteString("-- Update column type if it exists\n")
			sql.WriteString("DO $$\n")
			sql.WriteString("BEGIN\n")
			sql.WriteString(fmt.Sprintf("    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='%s' AND column_name='%s') THEN\n",
				tableName, columnName))
			sql.WriteString(fmt.Sprintf("        ALTER TABLE %s ALTER COLUMN %s TYPE %s;\n",
				tableName, columnName, columnType))
			sql.WriteString("    END IF;\n")
			sql.WriteString("END $$;\n\n")

		case "mysql":
			// MySQL doesn't support IF NOT EXISTS for ADD COLUMN
			// Use ALTER TABLE MODIFY to add or modify the column
			sql.WriteString("-- MySQL: Add or modify column\n")
			sql.WriteString("SET @s = (SELECT IF(\n")
			sql.WriteString(fmt.Sprintf("    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='%s' AND COLUMN_NAME='%s' AND TABLE_SCHEMA = DATABASE()),\n",
				tableName, columnName))
			sql.WriteString(fmt.Sprintf("    'ALTER TABLE %s MODIFY COLUMN %s %s;',\n",
				tableName, columnName, columnType))
			sql.WriteString(fmt.Sprintf("    'ALTER TABLE %s ADD COLUMN %s %s;'\n",
				tableName, columnName, columnType))
			sql.WriteString("));\n")
			sql.WriteString("PREPARE stmt FROM @s;\n")
			sql.WriteString("EXECUTE stmt;\n")
			sql.WriteString("DEALLOCATE PREPARE stmt;\n\n")

		case "mssql":
			// SQL Server requires a different approach
			sql.WriteString("-- SQL Server: Add column if it doesn't exist\n")
			sql.WriteString(fmt.Sprintf("IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '%s' AND COLUMN_NAME = '%s')\n",
				tableName, columnName))
			sql.WriteString(fmt.Sprintf("BEGIN\n    ALTER TABLE %s ADD %s %s;\nEND;\n",
				tableName, columnName, columnType))

			// Also add an ALTER COLUMN statement to modify the column type if it exists
			sql.WriteString("ELSE\n")
			sql.WriteString(fmt.Sprintf("BEGIN\n    ALTER TABLE %s ALTER COLUMN %s %s;\nEND;\n\n",
				tableName, columnName, columnType))
		}
	}

	return sql.String()
}

// fieldInfo holds a struct field with an accumulated column prefix
type fieldInfo struct {
	Field  reflect.StructField
	Prefix string
}

// collectDBFields flattens a struct type into a list of fields to generate columns for,
// handling gorm embedded structs with optional embeddedPrefix.
func collectDBFields(t reflect.Type, prefix string) []fieldInfo {
	// If we have a pointer, dereference
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var result []fieldInfo

	// Only process struct types
	if t.Kind() != reflect.Struct {
		return result
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip ignored fields
		if field.Tag.Get("gorm") == "-" {
			continue
		}

		// Detect embedded handling via gorm tag
		gormSettings := parseGormTag(field.Tag.Get("gorm"))
		if _, isEmbedded := gormSettings["embedded"]; isEmbedded {
			// Determine the embedded prefix for this field
			embeddedPrefix := gormSettings["embeddedPrefix"]
			newPrefix := prefix + embeddedPrefix

			// Recurse into the embedded struct type
			ft := field.Type
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
			if ft.Kind() == reflect.Struct {
				// Special-case: avoid recursing into time.Time
				if ft.PkgPath() == "time" && ft.Name() == "Time" {
					// Treat as regular field
					result = append(result, fieldInfo{Field: field, Prefix: prefix})
				} else {
					// Recurse into its fields, carrying the newPrefix
					nested := collectDBFields(ft, newPrefix)
					result = append(result, nested...)
				}
			} else {
				// Not a struct, just add as regular field
				result = append(result, fieldInfo{Field: field, Prefix: prefix})
			}
			continue
		}

		// Regular field
		result = append(result, fieldInfo{Field: field, Prefix: prefix})
	}

	return result
}

// parseGormTag parses a gorm struct tag into a key-value map.
// Flags like "embedded" are represented with value "true".
func parseGormTag(tag string) map[string]string {
	res := make(map[string]string)
	if tag == "" {
		return res
	}
	parts := strings.Split(tag, ";")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.Contains(p, ":") {
			kv := strings.SplitN(p, ":", 2)
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])
			res[key] = val
		} else {
			// boolean flag
			res[p] = "true"
		}
	}
	return res
}

// getTableName returns the table name for a struct using Gorm's schema API
func getTableName(structType reflect.Type) string {
	// Create a new instance of the struct
	modelValue := reflect.New(structType).Interface()

	// Try to get the table name using Gorm's schema API
	tableName := ""

	// If DB is initialized, use it to get the table name
	if DB != nil {
		stmt := &gorm.Statement{DB: DB}
		if err := stmt.Parse(modelValue); err == nil {
			tableName = stmt.Table
		}
	}

	// Fallback to manual determination if DB is not available or parsing fails
	if tableName == "" {
		// Create a new schema namer
		namer := schema.NamingStrategy{}

		// Check if the model implements Tabler interface
		if tabler, ok := modelValue.(schema.Tabler); ok {
			tableName = tabler.TableName()
		} else {
			// Use Gorm's default naming strategy
			tableName = namer.TableName(structType.Name())
		}
	}

	return tableName
}

// getColumnName returns the column name for a struct field using Gorm's schema API
func getColumnName(field reflect.StructField) string {
	// Try to get the column name from the gorm tag first
	tag, ok := field.Tag.Lookup("gorm")
	if ok {
		for _, str := range strings.Split(tag, ";") {
			if strings.HasPrefix(str, "column:") {
				return strings.TrimPrefix(str, "column:")
			}
		}
	}

	// Create a new schema namer
	namer := schema.NamingStrategy{}

	// Use Gorm's naming strategy to get the column name
	return namer.ColumnName("", field.Name)
}

// getColumnType returns the SQL column type for a struct field based on the database type
func getColumnType(field reflect.StructField, dbType string) string {
	// Get Go type, preserving special handling for time.Time (including pointers)
	var goType string
	timeType := reflect.TypeOf(time.Time{})
	t := field.Type
	if t == timeType || (t.Kind() == reflect.Ptr && t.Elem() == timeType) {
		goType = "time.Time"
	} else {
		goType = t.Kind().String()
	}

	// Check for specific type mappings in gorm tag
	gormTag := field.Tag.Get("gorm")
	if gormTag != "" {
		for _, tag := range strings.Split(gormTag, ";") {
			if strings.HasPrefix(tag, "type:") {
				return strings.TrimPrefix(tag, "type:")
			}
		}
	}

	// Map Go types to SQL types based on database type
	switch dbType {
	case "postgresql":
		return mapGoTypeToPostgreSQL(goType)
	case "mysql":
		return mapGoTypeToMySQL(goType)
	case "mssql":
		return mapGoTypeToMSSQL(goType)
	default:
		return "TEXT" // Default fallback
	}
}

// mapGoTypeToPostgreSQL maps Go types to PostgreSQL types
func mapGoTypeToPostgreSQL(goType string) string {
	switch goType {
	case "string":
		return "VARCHAR(255)"
	case "int", "int32", "int64":
		return "INTEGER"
	case "float32", "float64":
		return "NUMERIC(15,5)"
	case "bool":
		return "BOOLEAN"
	case "time.Time":
		return "TIMESTAMP WITH TIME ZONE"
	default:
		return "TEXT"
	}
}

// mapGoTypeToMySQL maps Go types to MySQL types
func mapGoTypeToMySQL(goType string) string {
	switch goType {
	case "string":
		return "VARCHAR(255)"
	case "int", "int32", "int64":
		return "INT"
	case "float32", "float64":
		return "DOUBLE"
	case "bool":
		return "TINYINT(1)"
	case "time.Time":
		return "DATETIME"
	default:
		return "TEXT"
	}
}

// mapGoTypeToMSSQL maps Go types to SQL Server types
func mapGoTypeToMSSQL(goType string) string {
	switch goType {
	case "string":
		return "NVARCHAR(255)"
	case "int", "int32", "int64":
		return "INT"
	case "float32", "float64":
		return "DECIMAL(15,5)"
	case "bool":
		return "BIT"
	case "time.Time":
		return "DATETIME2"
	default:
		return "NVARCHAR(MAX)"
	}
}

// toSnakeCase converts a string from camelCase to snake_case
func toSnakeCase(s string) string {
	// Special case for "ID" to convert to "id" instead of "i_d"
	if s == "ID" {
		return "id"
	}

	var result strings.Builder
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

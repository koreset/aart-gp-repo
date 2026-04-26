package services

import (
	appLog "api/log"
	"api/models"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// GenerateOptions controls the diff-based migration generator's behaviour.
//
// All flags default to the conservative value: only additive, non-destructive
// changes are emitted. Type changes and drops must be opted into explicitly
// because they can churn indexes and lose data.
type GenerateOptions struct {
	// AllowTypeChanges, if true, emits ALTER COLUMN TYPE statements when the
	// live column's declared type differs from what the struct says. Off by
	// default — comparison across dialects is fuzzy and easily produces
	// spurious diffs.
	AllowTypeChanges bool

	// AllowDestructive, if true, emits DROP COLUMN and DROP INDEX statements
	// for objects present in the database but absent from the struct. Off by
	// default — data loss is irreversible.
	AllowDestructive bool

	// Message is used in the generated migration filename. If empty, a default
	// is derived from the struct name(s) being migrated.
	Message string
}

// GenerateMigrationForStruct generates an incremental SQL migration file for
// a single struct. It compares the struct's expected schema to the live
// database and emits only the deltas (ALTER TABLE ... ADD COLUMN, CREATE
// INDEX, etc.).
//
// The dialect is auto-detected from the connected database (services.DbBackend)
// and output is written to migrations/<dialect>/<timestamp>_<message>.sql,
// which is the layout RunMigrationsOnStartup picks up.
//
// Returns the generated file path, or "" with no error if no schema changes
// were detected.
func GenerateMigrationForStruct(structName string, opts GenerateOptions) (string, error) {
	return GenerateMigrationForStructs([]string{structName}, opts)
}

// GenerateMigrationForStructs generates a single migration file containing
// the deltas for every named struct. Useful when a release introduces
// changes across several models — they end up in one timestamped file
// rather than one per struct.
func GenerateMigrationForStructs(structNames []string, opts GenerateOptions) (string, error) {
	if DB == nil {
		return "", fmt.Errorf("database not initialized; call services.SetupTables first")
	}
	if DbBackend == "" {
		return "", fmt.Errorf("database backend not detected (services.DbBackend is empty)")
	}
	if len(structNames) == 0 {
		return "", fmt.Errorf("at least one struct name is required")
	}

	var combined strings.Builder
	var emittedFor []string

	for _, name := range structNames {
		structType, err := getStructType(name)
		if err != nil {
			return "", err
		}
		model := reflect.New(structType).Interface()

		body, err := generateIncrementalSQLForModel(model, opts)
		if err != nil {
			return "", fmt.Errorf("generate %s: %w", name, err)
		}
		if strings.TrimSpace(body) == "" {
			appLog.WithField("struct", name).Info("No schema changes detected for struct")
			continue
		}
		combined.WriteString(body)
		combined.WriteString("\n")
		emittedFor = append(emittedFor, name)
	}

	sqlBody := combined.String()
	if strings.TrimSpace(sqlBody) == "" {
		appLog.Info("No schema changes detected — nothing to write")
		return "", nil
	}

	msg := opts.Message
	if msg == "" {
		msg = "update_" + strings.Join(emittedFor, "_")
	}

	return writeMigrationFile(msg, sqlBody)
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

	default:
		// If we couldn't find the struct, return an error
		return nil, fmt.Errorf("struct %s not found in package %s", structName, modelsPackage)
	}
}

// generateIncrementalSQLForModel returns the SQL needed to bring the live
// table for the given model up to the struct's expected schema. Returns ""
// when nothing has changed.
//
// Only additive changes are emitted by default; opts.AllowTypeChanges and
// opts.AllowDestructive enable the riskier transformations.
func generateIncrementalSQLForModel(model interface{}, opts GenerateOptions) (string, error) {
	s, err := parseSchema(model)
	if err != nil {
		return "", err
	}

	// Brand-new table — render full CREATE TABLE plus its indexes.
	if !DB.Migrator().HasTable(model) {
		return renderCreateTable(s, DbBackend), nil
	}

	// Existing table — diff columns and indexes.
	liveCols, err := DB.Migrator().ColumnTypes(model)
	if err != nil {
		return "", fmt.Errorf("introspect %s: %w", s.Table, err)
	}
	liveByName := make(map[string]gorm.ColumnType, len(liveCols))
	for _, c := range liveCols {
		liveByName[c.Name()] = c
	}

	var sb strings.Builder
	var changes int

	sb.WriteString(fmt.Sprintf("-- Migration for: %s (table: %s)\n\n", s.Name, s.Table))

	// Pass 1: columns in struct but missing from DB → ADD COLUMN.
	for _, f := range s.Fields {
		if !isMigratableField(f) {
			continue
		}
		if live, exists := liveByName[f.DBName]; exists {
			// Column exists. Type changes are out of scope unless explicitly
			// allowed; comparing types reliably across dialects is fuzzy.
			if opts.AllowTypeChanges {
				if alter := renderAlterColumnType(s.Table, f, live, DbBackend); alter != "" {
					sb.WriteString(alter)
					changes++
				}
			}
			continue
		}
		sb.WriteString(renderAddColumn(s.Table, f, DbBackend))
		changes++
	}

	// Pass 2: columns in DB but absent from struct → DROP COLUMN (opt-in).
	if opts.AllowDestructive {
		expected := make(map[string]struct{}, len(s.Fields))
		for _, f := range s.Fields {
			if isMigratableField(f) {
				expected[f.DBName] = struct{}{}
			}
		}
		for name := range liveByName {
			if _, ok := expected[name]; ok {
				continue
			}
			sb.WriteString(renderDropColumn(s.Table, name))
			changes++
		}
	}

	// Pass 3: indexes — additive only by default, drops opt-in.
	expectedIdx := s.ParseIndexes()
	for _, idx := range expectedIdx {
		if idx.Class == "PRIMARY" {
			continue
		}
		if !DB.Migrator().HasIndex(model, idx.Name) {
			sb.WriteString(renderCreateIndex(s.Table, idx))
			changes++
		}
	}
	if opts.AllowDestructive {
		liveIdx, _ := DB.Migrator().GetIndexes(model)
		expectedNames := make(map[string]struct{}, len(expectedIdx))
		for _, idx := range expectedIdx {
			expectedNames[idx.Name] = struct{}{}
		}
		for _, idx := range liveIdx {
			if _, ok := expectedNames[idx.Name()]; ok {
				continue
			}
			// Skip primary key indexes — those follow the table's lifecycle.
			if strings.HasSuffix(idx.Name(), "_pkey") || idx.Name() == "PRIMARY" {
				continue
			}
			sb.WriteString(renderDropIndex(s.Table, idx.Name(), DbBackend))
			changes++
		}
	}

	if changes == 0 {
		return "", nil
	}
	return sb.String(), nil
}

// parseSchema runs gorm's schema parser on the given model, producing the
// canonical expected schema (table name, fields with DBName, indexes, etc.).
func parseSchema(model interface{}) (*schema.Schema, error) {
	cache := &sync.Map{}
	s, err := schema.Parse(model, cache, DB.NamingStrategy)
	if err != nil {
		return nil, fmt.Errorf("parse schema: %w", err)
	}
	return s, nil
}

// isMigratableField reports whether the field should be considered for
// column-level migration. Filters out relations, ignored fields, and fields
// without a DB column (e.g. embedded relation handles).
func isMigratableField(f *schema.Field) bool {
	if f == nil || f.DBName == "" {
		return false
	}
	if f.IgnoreMigration {
		return false
	}
	return true
}

// writeMigrationFile writes the generated SQL to the dialect-specific
// migrations directory under the working directory, using a timestamp prefix
// so the runner applies files in order.
func writeMigrationFile(message, body string) (string, error) {
	dir := filepath.Join("migrations", DbBackend)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	safe := sanitizeFilenameSegment(message)
	filename := fmt.Sprintf("%s_%s.sql", time.Now().Format("20060102150405"), safe)
	path := filepath.Join(dir, filename)

	header := fmt.Sprintf("-- Generated %s for dialect %s\n\n", time.Now().Format(time.RFC3339), DbBackend)
	if err := os.WriteFile(path, []byte(header+body), 0o644); err != nil {
		return "", err
	}

	appLog.WithFields(map[string]interface{}{
		"file":    path,
		"dialect": DbBackend,
	}).Info("Wrote incremental migration file")

	return path, nil
}

func sanitizeFilenameSegment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "migration"
	}
	var sb strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			sb.WriteRune(r)
		case r == '_' || r == '-':
			sb.WriteRune(r)
		default:
			sb.WriteRune('_')
		}
	}
	return sb.String()
}

// ----- SQL rendering helpers -----

// renderCreateTable emits a full CREATE TABLE statement for a brand-new table,
// followed by any indexes declared in struct tags. The runner gates re-runs
// via the migrations history table, so we don't emit IF NOT EXISTS guards.
func renderCreateTable(s *schema.Schema, dialect string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("-- Create table: %s\n", s.Table))
	sb.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", s.Table))

	var lines []string
	var pkCols []string
	for _, f := range s.Fields {
		if !isMigratableField(f) {
			continue
		}
		spec := fmt.Sprintf("    %s %s", f.DBName, resolveSQLType(f, dialect))
		if !f.PrimaryKey {
			if f.NotNull {
				spec += " NOT NULL"
			}
			if f.Unique {
				spec += " UNIQUE"
			}
			if hasNonEmptyDefault(f) {
				spec += " DEFAULT " + f.DefaultValue
			}
		}
		lines = append(lines, spec)
		if f.PrimaryKey {
			pkCols = append(pkCols, f.DBName)
		}
	}
	if len(pkCols) > 0 {
		lines = append(lines, fmt.Sprintf("    PRIMARY KEY (%s)", strings.Join(pkCols, ", ")))
	}
	sb.WriteString(strings.Join(lines, ",\n"))
	sb.WriteString("\n);\n\n")

	for _, idx := range s.ParseIndexes() {
		if idx.Class == "PRIMARY" {
			continue
		}
		sb.WriteString(renderCreateIndex(s.Table, idx))
	}
	sb.WriteString("\n")
	return sb.String()
}

// renderAddColumn emits an ALTER TABLE ... ADD COLUMN. NOT NULL is honored
// only when a default is also declared, otherwise the column is added
// nullable with a comment so manual backfill can complete the migration.
func renderAddColumn(table string, f *schema.Field, dialect string) string {
	sqlType := resolveSQLType(f, dialect)
	spec := sqlType

	if f.NotNull {
		if hasNonEmptyDefault(f) {
			spec += " NOT NULL DEFAULT " + f.DefaultValue
		} else {
			spec += " /* NOT NULL omitted: declare a default or backfill manually before tightening */"
		}
	} else if hasNonEmptyDefault(f) {
		spec += " DEFAULT " + f.DefaultValue
	}

	keyword := "ADD COLUMN"
	if dialect == "mssql" {
		keyword = "ADD"
	}
	return fmt.Sprintf("ALTER TABLE %s %s %s %s;\n", table, keyword, f.DBName, spec)
}

// renderAlterColumnType emits an ALTER COLUMN TYPE when the live column's
// declared type doesn't match the struct's expected type. Returns "" when
// the types are compatible.
func renderAlterColumnType(table string, f *schema.Field, live gorm.ColumnType, dialect string) string {
	expected := resolveSQLType(f, dialect)
	actual, ok := live.ColumnType()
	if !ok || actual == "" {
		actual = live.DatabaseTypeName()
	}
	if normalizeType(actual) == normalizeType(expected) {
		return ""
	}
	switch dialect {
	case "postgresql":
		return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s;\n", table, f.DBName, expected)
	case "mysql":
		return fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s %s;\n", table, f.DBName, expected)
	case "mssql":
		return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s %s;\n", table, f.DBName, expected)
	}
	return ""
}

// renderDropColumn emits ALTER TABLE ... DROP COLUMN. Destructive — only
// reachable via opts.AllowDestructive.
func renderDropColumn(table, column string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;\n", table, column)
}

// renderCreateIndex emits CREATE INDEX (or CREATE UNIQUE INDEX) for a parsed
// gorm index definition.
func renderCreateIndex(table string, idx *schema.Index) string {
	cols := make([]string, 0, len(idx.Fields))
	for _, opt := range idx.Fields {
		cols = append(cols, opt.DBName)
	}
	unique := ""
	if idx.Class == "UNIQUE" {
		unique = "UNIQUE "
	}
	return fmt.Sprintf("CREATE %sINDEX %s ON %s (%s);\n", unique, idx.Name, table, strings.Join(cols, ", "))
}

// renderDropIndex emits DROP INDEX. MySQL requires the table reference;
// PostgreSQL and SQL Server do not.
func renderDropIndex(table, name, dialect string) string {
	if dialect == "mysql" {
		return fmt.Sprintf("DROP INDEX %s ON %s;\n", name, table)
	}
	return fmt.Sprintf("DROP INDEX %s;\n", name)
}

// hasNonEmptyDefault returns true when the field declares a default value
// suitable for SQL emission. GORM uses "(-)" as a sentinel for "explicitly
// no default", which we treat as no default at all.
func hasNonEmptyDefault(f *schema.Field) bool {
	return f.HasDefaultValue && f.DefaultValue != "" && f.DefaultValue != "(-)"
}

// normalizeType reduces a SQL type string to a comparable form so trivial
// formatting differences (case, spacing, "character varying" vs. "VARCHAR")
// don't trigger spurious ALTER COLUMN TYPE statements.
func normalizeType(t string) string {
	t = strings.ToLower(strings.TrimSpace(t))
	t = strings.ReplaceAll(t, "character varying", "varchar")
	t = strings.ReplaceAll(t, "double precision", "double")
	t = strings.ReplaceAll(t, " ", "")
	return t
}

// resolveSQLType returns the dialect-specific SQL type string for a field,
// honouring an explicit `gorm:"type:..."` override when present.
func resolveSQLType(f *schema.Field, dialect string) string {
	if t, ok := f.TagSettings["TYPE"]; ok && strings.TrimSpace(t) != "" {
		return t
	}
	switch f.GORMDataType {
	case schema.String:
		size := f.Size
		if size <= 0 {
			size = 255
		}
		switch dialect {
		case "postgresql", "mysql":
			return fmt.Sprintf("VARCHAR(%d)", size)
		case "mssql":
			return fmt.Sprintf("NVARCHAR(%d)", size)
		}
	case schema.Int:
		big := f.Size >= 64
		switch dialect {
		case "postgresql":
			if f.AutoIncrement {
				if big {
					return "BIGSERIAL"
				}
				return "SERIAL"
			}
			if big {
				return "BIGINT"
			}
			return "INTEGER"
		case "mysql":
			base := "INT"
			if big {
				base = "BIGINT"
			}
			if f.AutoIncrement {
				base += " AUTO_INCREMENT"
			}
			return base
		case "mssql":
			base := "INT"
			if big {
				base = "BIGINT"
			}
			if f.AutoIncrement {
				base += " IDENTITY(1,1)"
			}
			return base
		}
	case schema.Uint:
		big := f.Size >= 64
		switch dialect {
		case "postgresql":
			if big {
				return "BIGINT"
			}
			return "INTEGER"
		case "mysql":
			base := "INT UNSIGNED"
			if big {
				base = "BIGINT UNSIGNED"
			}
			if f.AutoIncrement {
				base += " AUTO_INCREMENT"
			}
			return base
		case "mssql":
			if big {
				return "BIGINT"
			}
			return "INT"
		}
	case schema.Float:
		switch dialect {
		case "postgresql":
			if f.Precision > 0 {
				return fmt.Sprintf("NUMERIC(%d,%d)", f.Precision, f.Scale)
			}
			return "NUMERIC(15,5)"
		case "mysql":
			if f.Precision > 0 {
				return fmt.Sprintf("DECIMAL(%d,%d)", f.Precision, f.Scale)
			}
			return "DOUBLE"
		case "mssql":
			if f.Precision > 0 {
				return fmt.Sprintf("DECIMAL(%d,%d)", f.Precision, f.Scale)
			}
			return "DECIMAL(15,5)"
		}
	case schema.Bool:
		switch dialect {
		case "postgresql":
			return "BOOLEAN"
		case "mysql":
			return "TINYINT(1)"
		case "mssql":
			return "BIT"
		}
	case schema.Time:
		switch dialect {
		case "postgresql":
			return "TIMESTAMP WITH TIME ZONE"
		case "mysql":
			return "DATETIME"
		case "mssql":
			return "DATETIME2"
		}
	case schema.Bytes:
		switch dialect {
		case "postgresql":
			return "BYTEA"
		case "mysql":
			return "BLOB"
		case "mssql":
			return "VARBINARY(MAX)"
		}
	}
	// Fallback for any other types.
	switch dialect {
	case "postgresql", "mysql":
		return "TEXT"
	case "mssql":
		return "NVARCHAR(MAX)"
	}
	return "TEXT"
}

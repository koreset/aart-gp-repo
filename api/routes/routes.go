package routes

import (
	"api/controllers"
	_ "api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConfigureRouter(router *gin.Engine) {
	// Apply request logger middleware to all routes
	router.Use(RequestLoggerMiddleware())

	router.GET("/health", controllers.CheckHealth)
	router.GET("/version", controllers.GetAppVersion)


	// WebSocket endpoint — registered outside auth middleware because
	// the browser WebSocket API cannot send custom headers. Auth is
	// handled inside the controller via the query-param token.
	router.GET("/ws", controllers.HandleWebSocketUpgrade)

	apiv1 := router.Group("", GetActiveUser())



	{

		apiv1.POST("org-users", controllers.GetOrgUsers)
		apiv1.POST("org-users/refresh", RequirePermission("system:manage_users"), controllers.RefreshOrgUsers)

		valuationJobs := apiv1.Group("/valuations/jobs")
		{
			valuationJobs.GET("", controllers.GetAllValuationJobs)
			valuationJobs.DELETE("", controllers.DeleteValuationJobs)

			valuationJobs.GET(":id", controllers.GetValuationJob)
			valuationJobs.POST(":id/restart", controllers.RestartStalledJob)
			valuationJobs.GET(":id/sp-code/:sp-code", controllers.GetValuationJob)
			valuationJobs.GET(":id/control", controllers.GetValuationJobControl)
			valuationJobs.DELETE(":id", controllers.DeleteValuationJob)
			valuationJobs.GET(":id/excel", controllers.GetValuationJobExcel)
			valuationJobs.GET(":id/excel/control", controllers.GetValuationJobExcel)
			valuationJobs.GET("all-jobs/:id/excel", controllers.GetProjectionJobExcel)
			valuationJobs.GET("all-jobs/:id/excel/scoped", controllers.GetProjectionJobExcelScoped)

		}



		// NoCache() forces every group-pricing response to carry
		// Cache-Control: no-store so dynamic per-quote data (result
		// summaries, output summary, premium summary, reinsurance premium
		// summary, table metadata, calculation status, etc.) is never
		// served from the renderer's HTTP cache. Without this a watcher
		// triggered re-fetch after a recalculation can hit a cached pre-
		// recalc payload and show stale numbers until the user manually
		// refreshes. Mutating verbs ignore the headers so the cost on
		// POST/PUT/DELETE is nil.
		groupPricing := apiv1.Group("group-pricing", NoCache())
		{
			groupPricing.POST("generate-quote", controllers.GenerateGroupPricingQuote)
			groupPricing.POST("calculate-quote/:id/basis/:basis", controllers.CalculateGroupPricingQuote)
			groupPricing.POST("calculate-quote/:id/basis/:basis/credibility/:credibility", controllers.CalculateGroupPricingQuote)
			groupPricing.GET("calculation-job/:jobId", controllers.GetCalculationJobStatus)

			groupPricing.POST("quotes/:id/update-status", controllers.UpdateGroupPricingQuote)
			groupPricing.POST("quotes/:id/approve-quote", controllers.ApproveGroupPricingQuote)
			groupPricing.POST("quotes/:id/accept-quote", controllers.AcceptGroupPricingQuote)
			groupPricing.POST("quotes/:id/on-risk-letter", controllers.CreateOnRiskLetter)
			groupPricing.GET("quotes/:id/on-risk-letter", controllers.GetOnRiskLetterData)
			groupPricing.GET("quotes/:id/on-risk-letter/document.docx", controllers.GenerateOnRiskLetterDocxTemplated)
			groupPricing.POST("quotes/:id/apply-discount/:discount", controllers.ApplyDiscountToQuote)
			groupPricing.GET("discount-authority/risk-code/:risk_rate_code", controllers.GetDiscountAuthority)
			groupPricing.GET("settings", controllers.GetGroupPricingSettings)
			groupPricing.PUT("settings", controllers.UpdateGroupPricingSettings)
			groupPricing.DELETE("quotes/:id", controllers.DeleteGroupPricingQuote)

			groupPricing.GET("get-quotes/filter/:filter", controllers.GetGroupPricingQuotes)
			groupPricing.GET("get-quotes", controllers.GetGroupPricingQuotes)
			// Retrieve a quote by its associated scheme name
			groupPricing.GET("get-quote-by-scheme-name/:schemeName", controllers.GetGroupPricingQuoteBySchemeName)
			groupPricing.GET("get-quote/:id", controllers.GetGroupPricingQuote)

			groupPricing.GET("age-bands", controllers.GetGroupPricingAgeBands)
			groupPricing.GET("quotes/:id/member-gender-split", controllers.GetQuoteMemberGenderSplit)
			groupPricing.GET("rate-tables", controllers.GetGPTableMetadata)
			groupPricing.POST("rate-tables/rebuild-stats", controllers.RebuildGPTableStats)
			groupPricing.GET("table-configurations", controllers.ListTableConfigurations)
			groupPricing.PATCH("table-configurations/:table_type", controllers.UpdateTableConfiguration)
			groupPricing.GET("table-configurations/:table_type/audit", controllers.GetTableConfigurationAuditLog)
			groupPricing.POST("rate-tables", controllers.UploadGPRateTables)
			groupPricing.GET("rate-tables/:table_type/get-years", controllers.GetGPTableYears)
			groupPricing.GET("rate-tables/:table_type/get-risk-codes", controllers.GetGPTableRiskCodes)
			groupPricing.DELETE("rate-tables/:table_type/risk-code/:risk_code", controllers.DeleteGPTableData)
			groupPricing.GET("rate-tables/:table_type", controllers.GetGPTableData)
			groupPricing.GET("rate-tables/:table_type/risk-rate-code/:risk_rate_code/waiting-periods", controllers.GetDistinctWaitingPeriods)
			groupPricing.GET("rate-tables/gla/risk-rate-code/:risk_rate_code/benefit-types", controllers.GetDistinctGlaBenefitTypes)
			groupPricing.GET("rate-tables/:table_type/:risk-rate-code/:risk_rate_code/deferred-periods", controllers.GetDistinctDeferredPeriods)
			groupPricing.GET("rate-tables/phi-rates/normal-retirement-ages", controllers.GetDistinctNormalRetirementAges)
			groupPricing.GET("rate-tables/risk-types", controllers.GetDistinctRiskTypes)
			groupPricing.GET("historical-credibility-data", controllers.GetHistoricalCredibilityData)
			groupPricing.GET("custom-tiered-income-replacement/check", controllers.CheckCustomTieredTableExists)
			groupPricing.POST("custom-tiered-income-replacement/request", controllers.RequestCustomTieredTable)
			groupPricing.POST("quote-tables", controllers.UploadQuoteTables)
			groupPricing.DELETE("quote-tables/:table_type/:quote_id", controllers.DeleteQuoteTableData)
			groupPricing.POST("brokers", controllers.CreateBroker)
			groupPricing.GET("brokers", controllers.GetBrokers)
			groupPricing.GET("brokers/:id", controllers.GetBroker)
			groupPricing.PUT("brokers/:id", controllers.EditBroker)
			groupPricing.DELETE("brokers/:id", controllers.DeleteBroker)

			// Binder fee management
			groupPricing.POST("binder-fees", controllers.CreateBinderFee)
			groupPricing.GET("binder-fees", controllers.GetBinderFees)
			groupPricing.GET("binder-fees/:id", controllers.GetBinderFee)
			groupPricing.PUT("binder-fees/:id", controllers.EditBinderFee)
			groupPricing.DELETE("binder-fees/:id", controllers.DeleteBinderFee)

			// Commission structure (sliding-scale per channel) management
			groupPricing.POST("commission-structures", controllers.CreateCommissionBand)
			groupPricing.GET("commission-structures", controllers.GetCommissionBands)
			groupPricing.GET("commission-structures/:id", controllers.GetCommissionBand)
			groupPricing.PUT("commission-structures/:id", controllers.EditCommissionBand)
			groupPricing.DELETE("commission-structures/:id", controllers.DeleteCommissionBand)

			// Reinsurer management
			groupPricing.POST("reinsurers", controllers.CreateReinsurer)
			groupPricing.GET("reinsurers", controllers.GetReinsurers)
			groupPricing.GET("reinsurers/:id", controllers.GetReinsurer)
			groupPricing.PUT("reinsurers/:id", controllers.EditReinsurer)
			groupPricing.POST("reinsurers/:id/deactivate", controllers.DeactivateReinsurer)
			groupPricing.DELETE("reinsurers/:id", controllers.DeleteReinsurer)
			groupPricing.POST("schemes", controllers.CreateGroupScheme)
			groupPricing.GET("schemes/check-name/:name", controllers.CheckGroupSchemeName)
			groupPricing.GET("schemes/in-force", controllers.GetAllGroupSchemes)
			groupPricing.GET("schemes/:id/categories", controllers.GetSchemeCategories)
			groupPricing.POST("scheme-category-masters", controllers.CreateSchemeCategoryMaster)
			groupPricing.GET("scheme-category-masters", controllers.GetSchemeCategoryMasters)
			groupPricing.GET("scheme-category-masters/:id", controllers.GetSchemeCategoryMasterByID)
			groupPricing.PUT("scheme-category-masters/:id", controllers.UpdateSchemeCategoryMaster)
			groupPricing.DELETE("scheme-category-masters/:id", controllers.DeleteSchemeCategoryMaster)
			groupPricing.GET("region-loadings/regions", controllers.GetRegionsForRiskCode)
			groupPricing.GET("schemes/in-force-v2", controllers.GetGroupSchemesInForce)
			groupPricing.GET("schemes/:id", controllers.GetGroupScheme)
			groupPricing.GET("schemes/:id/quotes", controllers.GetQuotesForScheme)
			groupPricing.POST("schemes/:id/members", controllers.AddMemberToScheme)
			groupPricing.GET("schemes/:id/members", controllers.GetSchemeMembers)
			// Paginated members with optional filters: page, pageSize, search, schemeId, status
			groupPricing.GET("members/paginated", controllers.GetMembersPaginated)
			// Get a single member in-force record by member_id (primary key)
			groupPricing.GET("members/:member_id", controllers.GetMemberInForce)
			// Update a member in-force record
			groupPricing.PUT("members/:member_id", controllers.UpdateMemberInForce)
			groupPricing.PATCH("members/:member_id", controllers.UpdateMemberInForce)
			// Member change history (generic audit)
			groupPricing.GET("members/:member_id/history", controllers.GetMemberHistory)
			// Get a single member in-force record by member IdNumber
			groupPricing.GET("members/id-number/:id_number", controllers.GetMemberInForceByIdNumber)
			// Member benefit summaries
			groupPricing.GET("members/:member_id/benefit-summary", controllers.GetMemberBenefitSummaryInForce)
			groupPricing.GET("quotes/:id/members/:member_id/benefit-summary", controllers.GetMemberBenefitSummaryQuote)
			// Beneficiaries CRUD
			groupPricing.GET("members/:member_id/beneficiaries", controllers.GetMemberBeneficiaries)
			groupPricing.POST("members/:member_id/beneficiaries", controllers.CreateMemberBeneficiary)
			groupPricing.GET("members/:member_id/beneficiaries/:id", controllers.GetMemberBeneficiary)
			groupPricing.PUT("members/:member_id/beneficiaries/:id", controllers.UpdateMemberBeneficiary)
			groupPricing.DELETE("members/:member_id/beneficiaries/:id", controllers.DeleteMemberBeneficiary)
			groupPricing.GET("schemes/:id/quotes/:quote_id/members/search", controllers.SearchSchemeMembers)
			groupPricing.PUT("schemes/:id/status", controllers.UpdateGroupSchemeStatus)
			groupPricing.GET("schemes/:id/status-history", controllers.GetGroupSchemeStatusAudit)
			groupPricing.DELETE("schemes/:id", controllers.DeleteGroupScheme)
			groupPricing.PUT("schemes/:id/cover-end-date", controllers.UpdateGroupSchemeCoverEndDate)
			//groupPricing.GET("schemes/:id/members/:member_id", controllers.GetSchemeMember)
			groupPricing.GET("parameter-bases", controllers.GetGroupPricingParameterBases)
			groupPricing.GET("industries", controllers.GetGroupPricingIndustries)
			groupPricing.GET("get-quote/:id/table-type/:table_type", controllers.GetGroupPricingQuoteTableData)
			groupPricing.GET("inforce-data/:id/table-type/:table_type", controllers.GetInForceTableData)
			groupPricing.GET("get-quote/:id/result-summary", controllers.GetGroupPricingQuoteResultSummary)
			groupPricing.GET("get-quote/:id/custom-tir-status", controllers.GetCustomTirTableStatus)
			groupPricing.GET("get-quote/:id/categories/educator-benefits", controllers.GetGroupPricingQuoteEducatorBenefits)
			groupPricing.GET("get-quote/:id/document.docx", controllers.GenerateGroupPricingQuoteDocx)
			groupPricing.GET("get-quote/:id/document.pdf", controllers.GenerateGroupPricingQuotePdf)
			groupPricing.POST("insurers", controllers.SaveInsurerDetails)
			groupPricing.GET("insurers", controllers.GetInsurerDetails)
			groupPricing.POST("insurers/:id/quote-template", controllers.UploadInsurerQuoteTemplate)
			groupPricing.GET("insurers/:id/quote-template/active", controllers.GetActiveInsurerQuoteTemplate)
			groupPricing.GET("insurers/:id/quote-template/versions", controllers.ListInsurerQuoteTemplates)
			groupPricing.GET("insurers/quote-template/:templateId/download", controllers.DownloadInsurerQuoteTemplate)
			groupPricing.POST("insurers/:id/quote-template/:templateId/activate", controllers.ActivateInsurerQuoteTemplate)
			groupPricing.DELETE("insurers/:id/quote-template/inactive", controllers.DeleteInactiveInsurerQuoteTemplates)
			groupPricing.DELETE("insurers/:id/quote-template/:templateId", controllers.DeleteInsurerQuoteTemplate)
			groupPricing.GET("quote-template/sample", controllers.DownloadSampleQuoteTemplate)
			groupPricing.POST("insurers/:id/on-risk-letter-template", controllers.UploadInsurerOnRiskLetterTemplate)
			groupPricing.GET("insurers/:id/on-risk-letter-template/active", controllers.GetActiveInsurerOnRiskLetterTemplate)
			groupPricing.GET("insurers/:id/on-risk-letter-template/versions", controllers.ListInsurerOnRiskLetterTemplates)
			groupPricing.GET("insurers/on-risk-letter-template/:templateId/download", controllers.DownloadInsurerOnRiskLetterTemplate)
			groupPricing.POST("insurers/:id/on-risk-letter-template/:templateId/activate", controllers.ActivateInsurerOnRiskLetterTemplate)
			groupPricing.DELETE("insurers/:id/on-risk-letter-template/inactive", controllers.DeleteInactiveInsurerOnRiskLetterTemplates)
			groupPricing.DELETE("insurers/:id/on-risk-letter-template/:templateId", controllers.DeleteInsurerOnRiskLetterTemplate)
			groupPricing.GET("on-risk-letter-template/sample", controllers.DownloadSampleOnRiskLetterTemplate)
			groupPricing.GET("dashboard/year/:year", controllers.GetGroupPricingDashboardData)
			groupPricing.GET("dashboard/exposures/year/:year/benefit/:benefit", controllers.GetGroupSchemeExposureData)
			groupPricing.POST("dashboard/exposures/rebuild/year/:year", controllers.RebuildExposureData)
			groupPricing.GET("dashboard/exposures/trend", controllers.GetExposureTrend)
			groupPricing.GET("metadata/financial-year-info", controllers.GetFinancialYearInfo)
			groupPricing.GET("quotes/:id/win-probability", controllers.GetQuoteWinProbability)
			groupPricing.POST("win-probability/train", controllers.TrainWinProbabilityModelHandler)
			groupPricing.GET("win-probability/model-info", controllers.GetWinProbabilityModelInfo)
			// Claims analytics dashboard
			groupPricing.GET("claims/dashboard", controllers.GetGroupSchemeClaimsDashboard)
			groupPricing.POST("claims", controllers.GroupSchemeSubmitClaim)
			groupPricing.GET("claims", controllers.GetGroupSchemeClaims)
			groupPricing.POST("claims/calculate-amount", controllers.GetUpdatedClaimAmount)
			groupPricing.GET("claims/:claim_id", controllers.GetGroupSchemeClaim)
			groupPricing.PUT("claims/:claim_id", controllers.UpdateGroupSchemeClaim)
			// Claim attachments
			groupPricing.GET("claims/:claim_id/attachments", controllers.ListGroupSchemeClaimAttachments)
			groupPricing.POST("claims/:claim_id/attachments", controllers.UploadGroupSchemeClaimAttachments)
			groupPricing.GET("claims/attachments/:attachment_id/download", controllers.DownloadGroupSchemeClaimAttachment)
			// Claim assessments
			groupPricing.POST("claims/:claim_id/assessments", controllers.CreateGroupSchemeClaimAssessment)
			groupPricing.GET("claims/:claim_id/assessments", controllers.GetGroupSchemeClaimAssessments)
			groupPricing.PUT("claims/assessments/:assessment_id", controllers.UpdateGroupSchemeClaimAssessment)
			// Claim communications
			groupPricing.POST("claims/:claim_id/communications", controllers.CreateGroupSchemeClaimCommunication)
			groupPricing.GET("claims/:claim_id/communications", controllers.GetGroupSchemeClaimCommunications)
			// Claim declines
			groupPricing.POST("claims/:claim_id/decline", controllers.CreateGroupSchemeClaimDecline)
			groupPricing.GET("claims/:claim_id/declines", controllers.GetGroupSchemeClaimDeclines)
			groupPricing.GET("claims/scheme/:scheme_id/quote/:quote_id/member/:member_id/rating", controllers.GetSchemeMemberRating)
			// Claim payment schedules
			groupPricing.POST("claims/payment-schedules", controllers.CreatePaymentSchedule)
			groupPricing.GET("claims/payment-schedules", controllers.GetPaymentSchedules)
			groupPricing.GET("claims/payment-schedules/:schedule_id", controllers.GetPaymentSchedule)
			groupPricing.GET("claims/payment-schedules/:schedule_id/export", controllers.ExportPaymentScheduleCSV)
			groupPricing.POST("claims/payment-schedules/:schedule_id/proof", controllers.UploadPaymentProof)
			groupPricing.GET("claims/payment-schedules/:schedule_id/proof", controllers.GetPaymentProofs)
			groupPricing.GET("claims/payment-schedules/proof/:proof_id/download", controllers.DownloadPaymentProof)
			// Bank account verification (provider-agnostic; see services/bav).
			// v1 kept for one release to avoid a lockstep api+app deploy; v2
			// below returns the canonical bav.VerifyResult shape.
			groupPricing.POST("claims/verify-bank-account", controllers.VerifyBankAccount)
			// ACB Bank profiles
			groupPricing.POST("claims/bank-profiles", controllers.CreateBankProfile)
			groupPricing.GET("claims/bank-profiles", controllers.GetBankProfiles)
			groupPricing.GET("claims/bank-profiles/:profile_id", controllers.GetBankProfile)
			groupPricing.PATCH("claims/bank-profiles/:profile_id", controllers.UpdateBankProfile)
			groupPricing.DELETE("claims/bank-profiles/:profile_id", controllers.DeleteBankProfile)
			// ACB file generation & reconciliation
			groupPricing.POST("claims/payment-schedules/:schedule_id/acb", controllers.GenerateACBFile)
			groupPricing.GET("claims/payment-schedules/:schedule_id/acb-files", controllers.GetACBFileRecords)
			groupPricing.GET("claims/acb-files/:acb_file_id/download", controllers.DownloadACBFile)
			groupPricing.POST("claims/acb-files/:acb_file_id/reconcile", controllers.ProcessBankResponse)
			groupPricing.GET("claims/acb-files/:acb_file_id/reconciliation", controllers.GetACBReconciliationResults)
			groupPricing.GET("claims/payment-schedules/:schedule_id/reconciliation-summary", controllers.GetACBReconciliationSummary)
			groupPricing.POST("claims/acb-files/:acb_file_id/retry", controllers.RetryFailedPayments)
			groupPricing.GET("benefit-maps", controllers.GetBenefitMaps)
			groupPricing.GET("benefit-maps/scheme/:scheme_id", controllers.GetBenefitMapsByScheme)
			groupPricing.GET("benefit-maps/scheme/:scheme_id/category/:category_id", controllers.GetBenefitMapsBySchemeCategory)
			groupPricing.POST("benefit-maps", controllers.SaveBenefitMaps)
			groupPricing.GET("benefit-definitions", controllers.GetBenefitDefinitions)
			// Email system (Phase 2): SMTP config, templates, outbox monitor
			groupPricing.GET("email/settings", controllers.GetEmailSettings)
			groupPricing.PUT("email/settings", RequirePermission("email:configure"), controllers.SaveEmailSettings)
			groupPricing.POST("email/settings/test", RequirePermission("email:configure"), controllers.SendTestEmail)
			groupPricing.GET("email/templates", controllers.ListEmailTemplates)
			groupPricing.POST("email/templates", RequirePermission("email:templates:manage"), controllers.CreateEmailTemplate)
			groupPricing.GET("email/templates/:code", controllers.GetEmailTemplate)
			groupPricing.PUT("email/templates/:code", RequirePermission("email:templates:manage"), controllers.UpdateEmailTemplate)
			groupPricing.DELETE("email/templates/:code", RequirePermission("email:templates:manage"), controllers.DeleteEmailTemplate)
			groupPricing.POST("email/templates/:code/preview", controllers.PreviewEmailTemplate)
			groupPricing.GET("email/outbox", RequirePermission("email:outbox:view"), controllers.ListEmailOutbox)
			groupPricing.GET("email/outbox/:id", RequirePermission("email:outbox:view"), controllers.GetEmailOutboxItem)
			groupPricing.POST("email/outbox/:id/retry", RequirePermission("email:outbox:view"), controllers.RetryEmailOutbox)
			// Bordereaux templates CRUD
			groupPricing.POST("bordereaux/templates", controllers.CreateBordereauxTemplate)
			groupPricing.GET("bordereaux/templates", controllers.GetBordereauxTemplates)
			groupPricing.GET("bordereaux/templates/:id", controllers.GetBordereauxTemplate)
			groupPricing.PUT("bordereaux/templates/:id", controllers.UpdateBordereauxTemplate)
			groupPricing.DELETE("bordereaux/templates/:id", controllers.DeleteBordereauxTemplate)
			groupPricing.POST("bordereaux/templates/:id/test", controllers.TestBordereauxTemplate)
			// Bordereaux generation
			groupPricing.GET("bordereaux/fields/:type", controllers.GetBordereauxFields)
			groupPricing.POST("bordereaux/generate", controllers.GenerateBordereaux)
			groupPricing.POST("bordereaux/batch-submit", controllers.SubmitBordereauxBatch)
			groupPricing.POST("bordereaux/confirmations/import", controllers.ImportSchemeConfirmations)
			groupPricing.POST("bordereaux/confirmations/reconcile-pending", controllers.ReconcilePendingConfirmations)
			groupPricing.GET("bordereaux/confirmations", controllers.GetBordereauxConfirmations)
			groupPricing.GET("bordereaux/confirmations/:id", controllers.GetBordereauxConfirmation)
			groupPricing.DELETE("bordereaux/confirmations/:id", controllers.DeleteBordereauxConfirmation)
			groupPricing.GET("bordereaux/confirmations/:id/results", controllers.GetReconciliationResults)
			groupPricing.GET("bordereaux/confirmations/:id/unmatched", controllers.GetUnmatchedReconciliationResults)
			groupPricing.POST("bordereaux/confirmations/:id/confirm", controllers.ConfirmReconciliation)
			groupPricing.POST("bordereaux/confirmations/:id/reprocess", controllers.ReprocessReconciliation)
			groupPricing.POST("bordereaux/confirmations/:id/note", controllers.AddReconciliationNote)
			groupPricing.GET("bordereaux/confirmations/:id/notes", controllers.GetReconciliationNotes)
			groupPricing.GET("bordereaux/reconciliation/stats", controllers.GetReconciliationStats)
			// Reconciliation result actions
			groupPricing.POST("bordereaux/reconciliation/results/:id/resolve", controllers.ResolveDiscrepancy)
			groupPricing.POST("bordereaux/reconciliation/results/:id/escalate", controllers.EscalateDiscrepancy)
			groupPricing.POST("bordereaux/reconciliation/results/:id/comment", controllers.AddDiscrepancyComment)
			groupPricing.GET("bordereaux/reconciliation/escalations", controllers.ListEscalations)
			groupPricing.GET("bordereaux/dashboard/stats", controllers.GetBordereauxDashboardStats)
			groupPricing.GET("bordereaux/analytics", controllers.GetBordereauxAnalytics)
			groupPricing.GET("bordereaux/compliance-report", controllers.GetBordereauxComplianceReport)
			groupPricing.GET("bordereaux/download/:filename", controllers.DownloadBordereaux)
			groupPricing.GET("bordereaux/generated", controllers.GetGeneratedBordereauxList)
			groupPricing.GET("bordereaux/generated/:id", controllers.GetGeneratedBordereauxByID)
			groupPricing.GET("bordereaux/generated/:id/data", controllers.GetGeneratedBordereauxData)
			groupPricing.POST("bordereaux/generated/:id/timeline", controllers.AddBordereauxTimelineEntry)
			groupPricing.POST("bordereaux/generated/:id/review", controllers.ReviewGeneratedBordereaux)
			groupPricing.POST("bordereaux/generated/:id/approve", controllers.ApproveGeneratedBordereaux)
			groupPricing.POST("bordereaux/generated/:id/return-to-draft", controllers.ReturnOutboundToDraft)
			groupPricing.POST("bordereaux/generated/:id/regenerate", controllers.RegenerateGeneratedBordereaux)
			// Bordereaux configuration management
			groupPricing.GET("bordereaux/configurations", controllers.GetBordereauxConfigurations)
			groupPricing.GET("bordereaux/configurations/:id", controllers.GetBordereauxConfiguration)
			groupPricing.POST("bordereaux/configurations", controllers.SaveBordereauxConfiguration)
			groupPricing.PUT("bordereaux/configurations/:id", controllers.UpdateBordereauxConfiguration)
			groupPricing.DELETE("bordereaux/configurations/:id", controllers.DeleteBordereauxConfiguration)
			groupPricing.PATCH("bordereaux/configurations/:id/usage", controllers.UpdateConfigurationUsage)
			// Reinsurer acceptance & recovery tracking
			groupPricing.POST("bordereaux/reinsurer/acceptances", controllers.CreateReinsurerAcceptance)
			groupPricing.GET("bordereaux/reinsurer/acceptances", controllers.GetReinsurerAcceptances)
			groupPricing.GET("bordereaux/reinsurer/acceptances/stats", controllers.GetAcceptanceStats)
			groupPricing.PATCH("bordereaux/reinsurer/acceptances/:id", controllers.UpdateReinsurerAcceptance)
			groupPricing.POST("bordereaux/reinsurer/recoveries", controllers.CreateReinsurerRecovery)
			groupPricing.GET("bordereaux/reinsurer/recoveries", controllers.GetReinsurerRecoveries)
			groupPricing.PATCH("bordereaux/reinsurer/recoveries/:id", controllers.UpdateReinsurerRecovery)
			// Submission deadline calendar
			groupPricing.GET("bordereaux/deadlines", controllers.GetBordereauxDeadlines)
			groupPricing.POST("bordereaux/deadlines", controllers.CreateBordereauxDeadline)
			groupPricing.POST("bordereaux/deadlines/generate", controllers.GenerateBordereauxDeadlines)
			groupPricing.PATCH("bordereaux/deadlines/:id/status", controllers.UpdateDeadlineStatus)
			groupPricing.GET("bordereaux/deadlines/stats", controllers.GetDeadlineStats)
			// Claim notification cadence
			groupPricing.POST("bordereaux/claim-notifications", controllers.CreateClaimNotification)
			groupPricing.GET("bordereaux/claim-notifications", controllers.GetClaimNotifications)
			groupPricing.GET("bordereaux/claim-notifications/stats", controllers.GetNotificationStats)
			groupPricing.POST("bordereaux/claim-notifications/generate-month-end", controllers.GenerateMonthEndNotifications)
			groupPricing.POST("bordereaux/claim-notifications/:id/sent", controllers.MarkNotificationSent)
			groupPricing.POST("bordereaux/claim-notifications/:id/acknowledged", controllers.MarkNotificationAcknowledged)
			groupPricing.DELETE("bordereaux/claim-notifications/:id", controllers.DeleteClaimNotification)
			groupPricing.GET("bordereaux/claim-notifications/claims-by-scheme/:scheme_id", controllers.GetClaimsByScheme)
			groupPricing.GET("bordereaux/claim-notifications/export", controllers.ExportNotificationsCSV)
			// Reinsurance Treaty Management
			treaties := groupPricing.Group("reinsurance/treaties")
			{
				treaties.POST("", controllers.CreateTreaty)
				treaties.GET("", controllers.GetTreaties)
				treaties.GET("/stats", controllers.GetTreatyStats)
				treaties.GET("/active/scheme/:scheme_id", controllers.GetActiveTreatiesForScheme)
				treaties.GET("/:id", controllers.GetTreatyByID)
				treaties.PUT("/:id", controllers.UpdateTreaty)
				treaties.DELETE("/:id", controllers.DeleteTreaty)
				treaties.POST("/:id/schemes", controllers.LinkSchemeToTreaty)
				treaties.POST("/:id/schemes/bulk", controllers.BulkLinkSchemesToTreaty)
				treaties.GET("/:id/schemes", controllers.GetTreatySchemeLinks)
				treaties.DELETE("/:id/schemes/bulk", controllers.BulkRemoveSchemeLinks)
				treaties.DELETE("/scheme-links/:link_id", controllers.RemoveSchemeTreatyLink)
			}
			// RI Bordereaux runs (member census + claims)
			riRuns := groupPricing.Group("reinsurance/bordereaux")
			{
				riRuns.POST("/member", controllers.GenerateRIMemberBordereaux)
				riRuns.POST("/claims", controllers.GenerateRIClaimsBordereaux)
				riRuns.GET("", controllers.GetRIBordereauxRuns)
				riRuns.GET("/stats", controllers.GetRIBordereauxStats)
				riRuns.GET("/kpis", controllers.GetRIBordereauxKPIs)
				riRuns.POST("/submit", controllers.SubmitRIBordereaux)
				riRuns.POST("/large-claims/monitor", controllers.MonitorLargeClaims)
				riRuns.GET("/large-claims", controllers.GetLargeClaimNotices)
				riRuns.GET("/large-claims/stats", controllers.GetLargeClaimStats)
				riRuns.PATCH("/large-claims/:id", controllers.UpdateLargeClaimNotice)
				riRuns.POST("/large-claims/:id/accept", controllers.AcceptLargeClaimNotice)
				riRuns.POST("/large-claims/:id/reject", controllers.RejectLargeClaimNotice)
				riRuns.POST("/large-claims/:id/query", controllers.QueryLargeClaimNotice)
				riRuns.GET("/:run_id", controllers.GetRIBordereauxRunByID)
				riRuns.GET("/:run_id/diff", controllers.DiffRIBordereauxRun)
				riRuns.POST("/:run_id/acknowledge", controllers.AcknowledgeRIBordereaux)
				riRuns.GET("/:run_id/members", controllers.GetRIBordereauxMemberRows)
				riRuns.GET("/:run_id/claims", controllers.GetRIBordereauxClaimsRows)
				riRuns.POST("/:run_id/validate", controllers.ValidateRIBordereaux)
				riRuns.GET("/:run_id/validation-results", controllers.GetRIValidationResults)
				riRuns.POST("/:run_id/acknowledge-receipt", controllers.AcknowledgeRIBordereauxReceipt)
				riRuns.POST("/:run_id/amend", controllers.AmendRIBordereaux)
				riRuns.GET("/cat-events", controllers.GetCatastropheClaimsRows)
				riRuns.GET("/run-off-treaties", controllers.GetRunOffTreaties)
			}
			// RI Technical Accounts & Settlement
			riSettlement := groupPricing.Group("reinsurance/settlement")
			{
				riSettlement.POST("", controllers.GenerateTechnicalAccount)
				riSettlement.GET("", controllers.GetTechnicalAccounts)
				riSettlement.GET("/stats", controllers.GetSettlementStats)
				riSettlement.POST("/payments", controllers.RecordSettlementPayment)
				riSettlement.GET("/payments", controllers.GetSettlementPayments)
				riSettlement.GET("/:id", controllers.GetTechnicalAccountByID)
				riSettlement.PATCH("/:id", controllers.UpdateTechnicalAccount)
				riSettlement.POST("/:id/escalate-dispute", controllers.EscalateSettlementDispute)
				riSettlement.POST("/:id/resolve-dispute", controllers.ResolveSettlementDispute)
			}
			// Inbound employer submissions
			submissions := groupPricing.Group("bordereaux/submissions")
			{
				submissions.POST("", controllers.CreateEmployerSubmission)
				submissions.GET("", controllers.GetEmployerSubmissions)
				submissions.GET("/:id", controllers.GetEmployerSubmission)
				submissions.POST("/:id/upload", controllers.UploadEmployerSubmission)
				submissions.POST("/:id/review", controllers.ReviewEmployerSubmission)
				submissions.POST("/:id/query", controllers.RaiseSubmissionQuery)
				submissions.POST("/:id/accept", controllers.AcceptEmployerSubmission)
				submissions.POST("/:id/reject", controllers.RejectEmployerSubmission)
				submissions.GET("/:id/records", controllers.GetSubmissionRecords)
				submissions.POST("/:id/generate-schedule", controllers.GenerateScheduleFromSubmission)
				submissions.POST("/:id/compute-delta", controllers.ComputeSubmissionDelta)
				submissions.GET("/:id/delta", controllers.GetSubmissionDeltaRecords)
				submissions.POST("/:id/sync-members", controllers.SyncSubmissionToMemberRegister)
				submissions.GET("/:id/register-diff", controllers.ComputeRegisterDiff)
				submissions.POST("/:id/snapshot-diff", controllers.SnapshotRegisterDiff)
				submissions.POST("/:id/apply-exits", controllers.ApplySubmissionExits)
				submissions.POST("/:id/apply-amendments", controllers.ApplySubmissionAmendments)
				submissions.GET("/:id/new-joiner-details", controllers.GetNewJoinerDetails)
				submissions.POST("/:id/upload-new-joiner-details", controllers.UploadNewJoinerDetails)
				submissions.POST("/:id/sync-new-joiners", controllers.SyncNewJoiners)
			}
			groupPricing.GET("user-management/permissions", RequirePermission("system:manage_roles"), controllers.GetGPPermissions)
			groupPricing.GET("user-management/roles", RequirePermission("system:manage_roles"), controllers.GetGPUserRoles)
			groupPricing.POST("user-management/roles", RequirePermission("system:manage_roles"), controllers.CreateGPUserRole)
			groupPricing.DELETE("user-management/roles/:role_id", RequirePermission("system:manage_roles"), controllers.DeleteGPUserRole)
			groupPricing.GET("user-management/roles/:role_id/permissions", RequirePermission("system:manage_roles"), controllers.GetRolePermissions)
			groupPricing.PUT("user-management/users/assign_role", RequirePermission("system:manage_users"), controllers.AssignRoleToUser)
			groupPricing.POST("user-management/users/remove_role", RequirePermission("system:manage_users"), controllers.RemoveRoleFromUser)
			// GetRoleForUser is used by every authenticated client during bootstrap
			// to discover its own permissions — gating it would deadlock the system.
			groupPricing.GET("user-management/users/license/:license_id/role", controllers.GetRoleForUser)
			groupPricing.GET("quotes/get-industries", controllers.GetGroupPricingIndustriesForQuotes)
			groupPricing.GET("quotes/benefit-escalations", controllers.GetBenefitEscalationOptions)
			groupPricing.GET("quotes/ttd-disability-definitions/risk-rate-code/:risk_rate_code", controllers.GetTTDDisabilityDefinitions)
			groupPricing.GET("quotes/ptd-disability-definitions/risk-rate-code/:risk_rate_code", controllers.GetPtdDisabilityDefinitions)
			groupPricing.GET("quotes/phi-disability-definitions/risk-rate-code/:risk_rate_code", controllers.GetPhiDisabilityDefinitions)
			groupPricing.GET("rate-tables/educator-benefits/risk-rate-code/:risk_rate_code", controllers.GetEducatorBenefitTypes)
			groupPricing.POST("quotes/indicative-member-data", controllers.UpdateGroupPricingQuoteMemberStats)
			groupPricing.PATCH("quotes/:id/indicative-member-data", controllers.UpdateGroupPricingQuoteIndicativeFlag)
			groupPricing.GET("export-csv/:quote_id/table-type/:table_type", controllers.GetTableDataCsvExport)
			groupPricing.DELETE("quotes/:id/indicative-member-data", controllers.DeleteIndicativeData)
			// Benefit document types CRUD
			groupPricing.POST("benefit-document-types", controllers.CreateBenefitDocumentType)
			groupPricing.GET("benefit-document-types", controllers.GetBenefitDocumentTypes)
			groupPricing.PUT("benefit-document-types/:id", controllers.UpdateBenefitDocumentType)
			groupPricing.DELETE("benefit-document-types/:id", controllers.DeleteBenefitDocumentType)
			groupPricing.GET("claims/:claim_id/required-documents", controllers.GetRequiredDocumentsForClaim)

			// Premiums Management — legacy endpoints
			groupPricing.POST("premiums/generate-schedule", controllers.GenerateMonthlySchedule)
			groupPricing.POST("premiums/generate-invoice/:schedule_id", controllers.GenerateInvoice)
			groupPricing.POST("premiums/record-payment", controllers.RecordPayment)
			groupPricing.GET("premiums/schemes/:scheme_id/arrears", controllers.GetArrearsStatus)

			// Premiums — Contribution Config
			groupPricing.GET("premiums/schemes/:scheme_id/contribution-config", controllers.GetContributionConfig)
			groupPricing.POST("premiums/schemes/:scheme_id/contribution-config", controllers.SaveContributionConfig)

			// Premiums — Schedules
			groupPricing.GET("premiums/schedules", controllers.GetPremiumSchedules)
			groupPricing.POST("premiums/schedules/generate", controllers.GenerateMonthlySchedule)
			groupPricing.POST("premiums/schedules/generate-all", controllers.GenerateAllSchedules)
			groupPricing.POST("premiums/schedules/:schedule_id/generate-invoice", controllers.GenerateInvoice)

			groupPricing.GET("premiums/schedules/:schedule_id", controllers.GetScheduleDetail)
			groupPricing.POST("premiums/schedules/:schedule_id/finalize", controllers.FinalizeSchedule)
			groupPricing.GET("premiums/schedules/:schedule_id/export", controllers.ExportScheduleCSV)
			groupPricing.POST("premiums/schedules/:schedule_id/review", controllers.ReviewSchedule)
			groupPricing.POST("premiums/schedules/:schedule_id/approve", controllers.ApproveSchedule)
			groupPricing.POST("premiums/schedules/:schedule_id/return-to-draft", controllers.ReturnScheduleToDraft)
			groupPricing.POST("premiums/schedules/:schedule_id/void", controllers.VoidSchedule)
			groupPricing.POST("premiums/schedules/:schedule_id/cancel", controllers.CancelSchedule)
			groupPricing.POST("premiums/schedules/:schedule_id/regenerate", controllers.RegenerateSchedule)
			groupPricing.DELETE("premiums/schedules/:schedule_id/members/:row_id", controllers.RemoveScheduleMember)
			groupPricing.PATCH("premiums/schedules/:schedule_id/members/:row_id", controllers.UpdateScheduleMemberRow)

			// Premiums — Invoices
			groupPricing.GET("premiums/invoices", controllers.GetInvoices)
			groupPricing.GET("premiums/invoices/:invoice_id", controllers.GetInvoiceDetail)
			groupPricing.PATCH("premiums/invoices/:invoice_id/mark-sent", controllers.MarkInvoiceSent)
			groupPricing.POST("premiums/invoices/:invoice_id/adjustments", controllers.CreateAdjustmentNote)

			// Premiums — Payments
			groupPricing.GET("premiums/payments", controllers.GetPayments)
			groupPricing.POST("premiums/payments", controllers.RecordPayment)
			groupPricing.DELETE("premiums/payments/:payment_id", controllers.VoidPayment)
			groupPricing.POST("premiums/payments/bulk-import", controllers.BulkImportPayments)

			// Premiums — Reconciliation (legacy — kept for backward compatibility)
			groupPricing.POST("premiums/reconciliation/auto-match", controllers.AutoMatchPayments)
			groupPricing.POST("premiums/reconciliation/manual-match", controllers.ManualMatchPayment)
			groupPricing.POST("premiums/reconciliation/credit-note", controllers.CreateCreditNote)
			groupPricing.POST("premiums/reconciliation/debit-note", controllers.CreateDebitNote)

			// Premiums — Reconciliation v2 (ledger-based allocation engine)
			groupPricing.GET("premiums/reconciliation/v2/summary", controllers.GetReconciliationSummaryV2)
			groupPricing.GET("premiums/reconciliation/v2/items", controllers.GetReconciliationItems)
			groupPricing.POST("premiums/reconciliation/v2/auto-match", controllers.RunAutoMatchV2)
			groupPricing.POST("premiums/reconciliation/v2/allocate", controllers.AllocatePaymentV2)
			groupPricing.POST("premiums/reconciliation/v2/reverse", controllers.ReverseAllocations)
			groupPricing.POST("premiums/reconciliation/v2/write-off", controllers.WriteOffBalance)
			groupPricing.POST("premiums/reconciliation/v2/refund", controllers.RefundOverpayment)
			groupPricing.GET("premiums/reconciliation/v2/history/:entity_type/:entity_id", controllers.GetAllocationHistory)
			groupPricing.GET("premiums/reconciliation/v2/runs", controllers.GetReconciliationRunsV2)
			groupPricing.GET("premiums/reconciliation/v2/runs/:run_id", controllers.GetReconciliationRunDetailV2)
			groupPricing.POST("premiums/reconciliation/v2/runs/:run_id/rollback", controllers.RollbackRunV2)
			groupPricing.PATCH("premiums/reconciliation/v2/items/:item_id/reassign", controllers.ReassignReconciliationItemV2)
			groupPricing.POST("premiums/reconciliation/v2/items/:item_id/suspend", controllers.SuspendReconciliationItemV2)
			groupPricing.GET("premiums/reconciliation/v2/matching-rules", controllers.GetMatchingRulesV2)
			groupPricing.POST("premiums/reconciliation/v2/matching-rules", controllers.SaveMatchingRuleV2)
			groupPricing.DELETE("premiums/reconciliation/v2/matching-rules/:rule_id", controllers.DeleteMatchingRuleV2)

			// Premiums — Arrears
			groupPricing.GET("premiums/arrears", controllers.GetAllArrearsAging)
			groupPricing.POST("premiums/schemes/:scheme_id/send-reminder", controllers.SendArrearsReminder)
			groupPricing.POST("premiums/schemes/:scheme_id/payment-plan", controllers.RecordPaymentPlan)
			groupPricing.POST("premiums/schemes/:scheme_id/suspend", controllers.SuspendSchemeCover)
			groupPricing.POST("premiums/schemes/:scheme_id/reinstate", controllers.ReinstateSchemeCover)
			groupPricing.GET("premiums/schemes/:scheme_id/arrears-history", controllers.GetArrearsHistory)
			groupPricing.GET("premiums/schemes/:scheme_id/payment-plans", controllers.GetPaymentPlans)

			// Premiums — Statements
			groupPricing.GET("premiums/statements/employer/:scheme_id", controllers.GetEmployerStatement)
			groupPricing.GET("premiums/statements/broker/:broker_id", controllers.GetBrokerCommissionStatement)

			// Premiums — Dashboard
			groupPricing.GET("premiums/dashboard", controllers.GetPremiumDashboard)
			groupPricing.GET("premiums/dashboard/collection-rate", controllers.GetCollectionRate)

			// Premiums — Coverage Matrix
			groupPricing.GET("premiums/schedules/coverage-matrix", controllers.GetScheduleCoverageMatrix)

			// Migration — bulk scheme & member import
			groupPricing.POST("migration/validate", controllers.ValidateMigration)
			groupPricing.POST("migration/execute", controllers.ExecuteMigration)
			groupPricing.GET("migration/templates/:template_name", controllers.DownloadMigrationTemplate)
		}

		// Generic audit logs for group-pricing
		apiv1.GET("audit/group-pricing", controllers.GetGroupPricingAuditLogs)

		phiValuations := apiv1.Group("phi-valuation")
		{
			phiValuations.GET("table-meta-data", controllers.GetPhiValuationTableMetadata)
			phiValuations.GET("tables/:table_type", controllers.GetPhiValuationTableData)
			phiValuations.GET("tables/:table_type/years", controllers.GetPhiValuationTableYears)
			phiValuations.GET("tables/:table_type/years/:year/versions", controllers.GetPhiValuationTableVersions)
			phiValuations.DELETE("tables/:table_type/year/:year/version/:version", controllers.DeletePhiValuationTableData)
			phiValuations.POST("tables/upload", controllers.UploadPhiValuationTables)
			phiValuations.POST("shock-settings", controllers.CreatePhiShockSetting)
			phiValuations.GET("shock-settings", controllers.GetPhiShockSettings)
			phiValuations.DELETE("shock-settings/:setting_name", controllers.DeleteShockSetting)
			phiValuations.GET("shock-bases", controllers.GetPhiShockBases)
			phiValuations.GET("model-point-years", controllers.GetPhiModelPointYears)
			phiValuations.GET("model-point-versions/year/:year", controllers.GetPhiModelPointYearVersions)
			phiValuations.GET("model-point-count", controllers.GetPhiModelPointCount)
			phiValuations.GET("model-points/:year/:version", controllers.GetPhiModelPointsData)
			phiValuations.GET("model-points/:year/:version/excel", controllers.GetPhiModelPointsExcel)
			phiValuations.DELETE("model-points/:year/:version", controllers.DeletePhiModelPoints)

			// Years and versions for other PHI data sets
			phiValuations.GET("parameter-years", controllers.GetPhiParameterYears)
			phiValuations.GET("parameter-versions/year/:year", controllers.GetPhiParameterYearVersions)
			phiValuations.GET("mortality-years", controllers.GetPhiMortalityYears)
			phiValuations.GET("mortality-versions/year/:year", controllers.GetPhiMortalityYearVersions)
			phiValuations.GET("recovery-rate-years", controllers.GetPhiRecoveryRateYears)
			phiValuations.GET("recovery-rate-versions/year/:year", controllers.GetPhiRecoveryRateYearVersions)
			phiValuations.GET("yield-curve-years", controllers.GetPhiYieldCurveYears)
			phiValuations.GET("yield-curve-versions/year/:year", controllers.GetPhiYieldCurveYearVersions)
			phiValuations.POST("run-projections", controllers.RunPhiProjections)
			phiValuations.GET("run-jobs", controllers.GetAllPhiRunJobs)
			phiValuations.GET("run-jobs/:id", controllers.GetPhiRunJob)
			phiValuations.GET("run-jobs/:id/control", controllers.GetPhiRunJobControl)
			phiValuations.DELETE("run-jobs/:id", controllers.DeletePhiRunJob)
			phiValuations.GET("run-configs", controllers.ListPhiRunConfigs)
			phiValuations.POST("run-configs", controllers.SavePhiRunConfig)
			phiValuations.DELETE("run-configs/:id", controllers.DeletePhiRunConfig)

		}

		// ─── In-App Communications ──────────────────────────────────────
		notifications := apiv1.Group("notifications")
		{
			notifications.GET("", controllers.GetNotifications)
			notifications.GET("unread-count", controllers.GetUnreadCount)
			notifications.PATCH(":id/read", controllers.MarkNotificationAsRead)
			notifications.POST("read-all", controllers.MarkAllNotificationsAsRead)
			notifications.DELETE(":id", controllers.DeleteNotification)
		}

		conversations := apiv1.Group("conversations")
		{
			conversations.POST("", controllers.CreateConversation)
			conversations.GET("", controllers.GetUserConversations)
			conversations.GET("by-object", controllers.GetObjectConversations)
			conversations.GET(":id", controllers.GetConversation)
			conversations.POST(":id/messages", controllers.SendMessage)
			conversations.PATCH("messages/:message_id", controllers.EditMessage)
			conversations.DELETE("messages/:message_id", controllers.DeleteMessage)
			conversations.POST(":id/participants", controllers.AddParticipant)
			conversations.POST(":id/read", controllers.MarkConversationRead)
		}

	}

	// v2 endpoints emit canonical provider-agnostic shapes. Added incrementally
	// as endpoints are rewritten; v1 equivalents are retained for one release
	// so api and app can deploy independently.
	apiv2 := router.Group("/v2", GetActiveUser())
	{
		apiv2.POST("group-pricing/claims/verify-bank-account", controllers.VerifyBankAccountV2)
		apiv2.POST("group-pricing/claims/verify-bank-account/status/:job_id", controllers.VerifyBankAccountV2Status)
	}

	// BAV webhook endpoint. Scaffolded behind its own route group because
	// webhooks are unauthenticated HTTP — HMAC-verified inside the handler.
	// Inert (501) until Phase 7 wires a real async provider.
	router.POST("/bav/webhook/:provider", controllers.BAVWebhook)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}

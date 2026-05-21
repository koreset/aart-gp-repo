import ErrorScreen from '@/renderer/screens/ErrorScreen.vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import App from '../App.vue'
import AppLogin from '../AppLogin.vue'
import AppSetup from '../AppSetup.vue'
import { useGroupUserPermissionsStore } from '../store/group_user'

const checkPermissions = async (to: any, _from: any) => {
  const store = useGroupUserPermissionsStore()
  const required = to.meta.required_permission as string | undefined
  const requiredAny = to.meta.required_permissions_any as string[] | undefined
  if (!required && (!requiredAny || requiredAny.length === 0)) return true
  await store.waitUntilLoaded()
  // Bootstrap mode: user has no role assigned (fresh install) — allow
  // access so an initial admin can be configured.
  if (!store.hasRole) return true
  if (store.hasPermission('system:admin')) return true
  if (required && store.hasPermission(required)) return true
  if (requiredAny && requiredAny.some((p) => store.hasPermission(p)))
    return true
  return { name: 'group-pricing-dashboard' }
}

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/app',
      name: 'App',
      component: App
    },
    {
      path: '/login',
      name: 'AppLogin',
      component: AppLogin
    },
    {
      path: '/setup',
      name: 'setup',
      component: AppSetup
    },

    // ── Redirect root to GP dashboard ────────────────────────────────────────
    {
      path: '/',
      redirect: '/group-pricing/dashboard'
    },

    // ── Group Pricing ────────────────────────────────────────────────────────
    {
      path: '/group-pricing/dashboard',
      name: 'group-pricing-dashboard',
      component: () => import('../screens/group_pricing/GPDashBoard.vue'),
      meta: { required_permission: 'navigation:view_gp_dashboard' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/dashboard/brokers/:brokerId/year/:year',
      name: 'group-pricing-broker-performance',
      component: () =>
        import('../screens/group_pricing/dashboard/BrokerPerformanceDetail.vue'),
      props: true,
      meta: { required_permission: 'navigation:view_gp_dashboard' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/quote-performance',
      name: 'group-pricing-quote-performance',
      component: () =>
        import('../screens/group_pricing/dashboard/QuotePerformanceDashboard.vue'),
      meta: { required_permission: 'quote:view_performance_dashboard' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/quote-performance/extract',
      name: 'group-pricing-quote-extract',
      component: () =>
        import('../screens/group_pricing/dashboard/QuoteExtractView.vue'),
      meta: { required_permission: 'quote:view_performance_dashboard' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/quote-performance/sla-targets',
      name: 'group-pricing-sla-targets',
      component: () =>
        import('../screens/group_pricing/dashboard/SlaTargetSettings.vue'),
      meta: { required_permission: 'quote:manage_sla_targets' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/quote-performance/user-flags',
      name: 'group-pricing-user-flags',
      component: () =>
        import('../screens/group_pricing/dashboard/UserFlagsAdmin.vue'),
      meta: { required_permission: 'quote:manage_user_flags' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/quote-generation',
      name: 'group-pricing-quote-generation',
      component: () => import('../screens/group_pricing/QuoteGeneration.vue'),
      meta: { required_permission: 'navigation:view_quotes' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/quote-generation/:id',
      name: 'group-pricing-quote-generation-edit',
      component: () => import('../screens/group_pricing/QuoteGeneration.vue'),
      meta: { required_permission: 'navigation:view_quotes' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/metadata',
      name: 'group-pricing-metadata',
      component: () => import('../screens/group_pricing/MetaData.vue'),
      meta: { required_permission: 'navigation:view_metadata' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/schemes',
      name: 'group-pricing-schemes',
      component: () => import('../screens/group_pricing/GroupSchemeList.vue'),
      meta: { required_permission: 'navigation:view_schemes' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/schemes/:id',
      name: 'group-pricing-schemes-detail',
      props: true,
      component: () => import('../screens/group_pricing/GroupSchemeDetail.vue'),
      meta: { required_permission: 'navigation:view_schemes' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/scheme-details/:id',
      name: 'group-pricing-scheme-details',
      props: true,
      component: () => import('../screens/group_pricing/NewQuoteDetail.vue'),
      meta: { required_permission: 'navigation:view_schemes' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/quotes',
      name: 'group-pricing-quotes',
      component: () => import('../screens/group_pricing/GroupSchemeQuotes.vue'),
      meta: { required_permission: 'navigation:view_quotes' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/quotes/output/:quoteId',
      name: 'group-pricing-quotes-generation',
      component: () => import('../screens/group_pricing/QuoteOutput.vue'),
      props: true,
      meta: { required_permission: 'navigation:view_quotes' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/tables',
      name: 'group-pricing-tables',
      component: () => import('../screens/group_pricing/Tables.vue'),
      meta: { required_permission: 'navigation:group_tables' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/claims-list',
      name: 'group-pricing-claims-list',
      component: () => import('../screens/group_pricing/ClaimsList.vue'),
      meta: { required_permission: 'navigation:view_claims' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/lodge-claim',
      name: 'group-pricing-lodge-claim',
      component: () => import('../screens/group_pricing/LodgeClaim.vue'),
      meta: { required_permission: 'navigation:view_claims' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/administration/member-management',
      name: 'group-pricing-member-management',
      component: () =>
        import('../screens/group_pricing/administration/MemberManagement.vue'),
      meta: { required_permission: 'navigation:manage_members' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/administration/member-management/bulk-enrollment',
      name: 'group-pricing-bulk-enrollment',
      component: () =>
        import('../screens/group_pricing/administration/BulkEnrollment.vue'),
      meta: { required_permission: 'navigation:manage_members' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/administration/member-management/bulk-enrollment/batches/:batchId',
      name: 'group-pricing-bulk-enrollment-batch',
      component: () =>
        import('../screens/group_pricing/administration/BulkEnrollmentBatchReview.vue'),
      props: true,
      meta: { required_permission: 'navigation:manage_members' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/administration/member-management/:id',
      name: 'group-pricing-member-details',
      component: () =>
        import('../screens/group_pricing/administration/MemberDetails.vue'),
      props: true,
      meta: { required_permission: 'navigation:manage_members' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/administration/beneficiaries',
      name: 'group-pricing-beneficiaries-overview',
      component: () =>
        import('../screens/group_pricing/administration/BeneficiaryManagement.vue'),
      meta: { required_permission: 'navigation:manage_beneficiaries' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/administration/beneficiary-management/:memberId',
      name: 'group-pricing-beneficiary-management',
      component: () =>
        import('../screens/group_pricing/administration/BeneficiaryManagement.vue'),
      meta: { required_permission: 'navigation:manage_beneficiaries' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/administration/scheme-migration',
      name: 'group-pricing-scheme-migration',
      component: () =>
        import('../screens/group_pricing/administration/SchemeMigration.vue'),
      meta: { required_permission: 'navigation:manage_members' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/administration/binder-fees',
      name: 'group-pricing-binder-fees',
      component: () =>
        import('../screens/group_pricing/administration/BinderFeeManagement.vue'),
      meta: { required_permission: 'navigation:manage_binder_fees' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/administration/commission-structures',
      name: 'group-pricing-commission-structures',
      component: () =>
        import('../screens/group_pricing/administration/CommissionStructureManagement.vue'),
      meta: { required_permission: 'navigation:manage_commission_structures' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/administration/fraud-risk',
      name: 'group-pricing-fraud-risk',
      component: () =>
        import('../screens/group_pricing/administration/FraudRiskConfig.vue'),
      meta: { required_permission: 'navigation:manage_claims' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/claims-management',
      name: 'group-pricing-claims-management',
      component: () =>
        import('../screens/group_pricing/claims_management/ClaimsHome.vue'),
      meta: { required_permission: 'navigation:manage_claims' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/claims-management/regular-income',
      name: 'group-pricing-regular-income-claims',
      component: () =>
        import('../screens/group_pricing/claims_management/RegularIncomeHome.vue'),
      meta: { required_permission: 'claims:view_regular_income' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/claims-management/cpi-index',
      name: 'group-pricing-cpi-index',
      component: () =>
        import('../screens/group_pricing/claims_management/CpiIndexManagement.vue'),
      meta: { required_permission: 'claims:view_regular_income' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      // Quote-grouped queue — the top-level Underwriting entry point.
      // Drills into a quote's Underwriting tab via ?tab=underwriting.
      path: '/group-pricing/underwriting',
      name: 'group-pricing-underwriting-cases',
      component: () =>
        import('../screens/group_pricing/UnderwritingQuoteList.vue'),
      meta: { required_permission: 'underwriting:view' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      // Flat per-member case list. Kept as a deep-link target for
      // analysts and accessible from the "Flat case list" button at
      // the top of the quote queue.
      path: '/group-pricing/underwriting/flat',
      name: 'group-pricing-underwriting-cases-flat',
      component: () =>
        import('../screens/group_pricing/UnderwritingCaseList.vue'),
      meta: { required_permission: 'underwriting:view' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/underwriting/:caseId',
      name: 'group-pricing-underwriting-case-detail',
      component: () =>
        import('../screens/group_pricing/UnderwritingCaseDetail.vue'),
      props: true,
      meta: { required_permission: 'underwriting:view' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/underwriting-rules',
      name: 'group-pricing-underwriting-rules',
      component: () =>
        import('../screens/group_pricing/UnderwritingRulesAdmin.vue'),
      meta: { required_permission: 'underwriting:admin' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/quotes/:quoteId/takeover',
      name: 'group-pricing-takeover-upload',
      component: () => import('../screens/group_pricing/TakeoverUpload.vue'),
      props: true,
      meta: { required_permission: 'underwriting:decide' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/claims-management/payment-schedules',
      name: 'group-pricing-claim-payment-schedules',
      component: () =>
        import('../screens/group_pricing/claims_management/ClaimPaymentSchedules.vue'),
      meta: { required_permission: 'claims_pay:finance_review' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/claims-management/my-submissions',
      name: 'group-pricing-claim-my-submissions',
      component: () =>
        import('../screens/group_pricing/claims_management/ClaimMySubmissions.vue'),
      meta: { required_permission: 'claims_pay:create_schedule' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/claims-management/my-submissions/:scheduleId/claims',
      name: 'group-pricing-claim-my-submissions-claims',
      component: () =>
        import('../screens/group_pricing/claims_management/ClaimScheduleClaimsList.vue'),
      props: true,
      meta: { required_permission: 'claims_pay:create_schedule' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/claims-management/payment-schedules/:scheduleId',
      name: 'group-pricing-claim-payment-schedule-detail',
      component: () =>
        import('../screens/group_pricing/claims_management/ClaimPaymentScheduleLayout.vue'),
      props: true,
      meta: {
        required_permissions_any: [
          'claims_pay:create_schedule',
          'claims_pay:finance_review'
        ]
      },
      beforeEnter: (to, from) => checkPermissions(to, from),
      redirect: (to) => ({
        name: 'group-pricing-claim-payment-schedule-claims',
        params: { scheduleId: to.params.scheduleId }
      }),
      children: [
        {
          path: 'claims',
          name: 'group-pricing-claim-payment-schedule-claims',
          component: () =>
            import('../screens/group_pricing/claims_management/ClaimPaymentScheduleClaims.vue')
        },
        {
          path: 'acb',
          name: 'group-pricing-claim-payment-schedule-acb',
          component: () =>
            import('../screens/group_pricing/claims_management/ClaimPaymentScheduleACB.vue'),
          meta: { required_permission: 'claims_pay:generate_acb' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        },
        {
          path: 'queries',
          name: 'group-pricing-claim-payment-schedule-queries',
          component: () =>
            import('../screens/group_pricing/claims_management/ClaimPaymentScheduleQueries.vue')
        },
        {
          path: 'reconciliation',
          name: 'group-pricing-claim-payment-schedule-reconciliation',
          component: () =>
            import('../screens/group_pricing/claims_management/ClaimPaymentScheduleReconciliation.vue'),
          meta: { required_permission: 'claims_pay:upload_response' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        },
        {
          path: 'proofs',
          name: 'group-pricing-claim-payment-schedule-proofs',
          component: () =>
            import('../screens/group_pricing/claims_management/ClaimPaymentScheduleProofs.vue'),
          meta: { required_permission: 'claims_pay:upload_response' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        }
      ]
    },
    {
      path: '/group-pricing/claims-management/authority-matrix',
      name: 'group-pricing-claim-authority-matrix',
      component: () =>
        import('../screens/group_pricing/claims_management/AuthorityMatrix.vue'),
      meta: { required_permission: 'claims_pay:admin_authority' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/claims-management/payment-cutoff',
      name: 'group-pricing-payment-cutoff-settings',
      component: () =>
        import('../screens/group_pricing/claims_management/PaymentCutoffSettings.vue'),
      meta: { required_permission: 'claims_pay:admin_cutoff' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/claims-management/payment-letter-settings',
      name: 'group-pricing-payment-letter-settings',
      component: () =>
        import('../screens/group_pricing/claims_management/PaymentLetterSettings.vue'),
      meta: { required_permission: 'claims_pay:configure_letter_settings' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/claims-management/payment-exceptions',
      name: 'group-pricing-claim-payment-exceptions',
      component: () =>
        import('../screens/group_pricing/claims_management/ClaimPaymentExceptions.vue'),
      meta: { required_permission: 'claims_pay:view_exceptions' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/claims-management/:id',
      name: 'group-pricing-claim-details',
      component: () =>
        import('../screens/group_pricing/claims_management/ClaimDetails.vue'),
      props: true,
      meta: { required_permission: 'navigation:manage_claims' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/claims-management/:id/assess',
      name: 'group-pricing-claim-assess',
      component: () =>
        import('../screens/group_pricing/claims_management/ClaimAssessment.vue'),
      props: true,
      meta: { required_permission: 'claims:assess' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/claims-management/:id/edit',
      name: 'group-pricing-claim-edit',
      component: () =>
        import('../screens/group_pricing/claims_management/ClaimEdit.vue'),
      props: true,
      meta: { required_permission: 'claims:lodge' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/claims-analytics',
      name: 'group-pricing-claims-analytics',
      component: () =>
        import('../screens/group_pricing/claims_management/ClaimsAnalytics.vue'),
      meta: { required_permission: 'navigation:manage_claims' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },

    // ── Finance: General Ledger (tabbed layout) ─────────────────────────────
    {
      path: '/group-pricing/finance/general-ledger',
      component: () =>
        import('../screens/group_pricing/finance/GeneralLedgerLayout.vue'),
      redirect: { name: 'group-pricing-gl-journals' },
      children: [
        {
          path: 'journals',
          name: 'group-pricing-gl-journals',
          component: () =>
            import('../screens/group_pricing/finance/JournalEntries.vue'),
          meta: { required_permission: 'gl:view' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        },
        {
          path: 'trial-balance',
          name: 'group-pricing-gl-trial-balance',
          component: () =>
            import('../screens/group_pricing/finance/TrialBalance.vue'),
          meta: { required_permission: 'gl:view' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        },
        {
          path: 'chart-of-accounts',
          name: 'group-pricing-gl-chart-of-accounts',
          component: () =>
            import('../screens/group_pricing/finance/ChartOfAccounts.vue'),
          meta: { required_permission: 'gl:view' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        },
        {
          path: 'periods',
          name: 'group-pricing-gl-periods',
          component: () =>
            import('../screens/group_pricing/finance/AccountingPeriods.vue'),
          meta: { required_permission: 'gl:view' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        },
        {
          path: 'posting-rules',
          name: 'group-pricing-gl-posting-rules',
          component: () =>
            import('../screens/group_pricing/finance/PostingRules.vue'),
          meta: { required_permission: 'gl:view' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        },
        {
          path: 'audit-log',
          name: 'group-pricing-gl-audit-log',
          component: () =>
            import('../screens/group_pricing/finance/AuditLog.vue'),
          meta: { required_permission: 'gl:view_audit' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        }
      ]
    },
    // ── Finance: Cash & Banking (tabbed layout) ─────────────────────────────
    {
      path: '/group-pricing/finance/cash-banking',
      component: () =>
        import('../screens/group_pricing/finance/CashAndBankingLayout.vue'),
      redirect: { name: 'group-pricing-gl-bank-accounts' },
      children: [
        {
          path: 'accounts',
          name: 'group-pricing-gl-bank-accounts',
          component: () =>
            import('../screens/group_pricing/finance/BankAccounts.vue'),
          meta: { required_permission: 'gl:view' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        },
        {
          path: 'reconciliation',
          name: 'group-pricing-gl-bank-reconciliation',
          component: () =>
            import('../screens/group_pricing/finance/BankReconciliation.vue'),
          meta: { required_permission: 'gl:bank_rec' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        }
      ]
    },
    // Detail routes stay full-bleed (no tab bar wrapper).
    {
      path: '/group-pricing/finance/journals/:id',
      name: 'group-pricing-gl-journal-detail',
      component: () =>
        import('../screens/group_pricing/finance/JournalEntryDetail.vue'),
      props: true,
      meta: { required_permission: 'gl:view' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/finance/accounts/:id/ledger',
      name: 'group-pricing-gl-account-ledger',
      component: () =>
        import('../screens/group_pricing/finance/GeneralLedger.vue'),
      props: true,
      meta: { required_permission: 'gl:view' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },

    // ── Bordereaux Management ────────────────────────────────────────────────
    {
      path: '/group-pricing/bordereaux-management',
      name: 'group-pricing-bordereaux-management',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/BordereauxManagement.vue'),
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/bordereaux-management/generation',
      name: 'group-pricing-bordereaux-generation',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/components/BordereauxGenerationForm.vue'),
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/bordereaux-management/tracking',
      name: 'group-pricing-bordereaux-tracking',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/components/BordereauxSubmissionTracking.vue'),
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/bordereaux-management/reconciliation',
      name: 'group-pricing-bordereaux-reconciliation',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/components/BordereauxReconciliation.vue'),
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/bordereaux-management/templates',
      name: 'group-pricing-bordereaux-templates',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/components/BordereauxTemplateManager.vue'),
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/bordereaux-management/analytics',
      name: 'group-pricing-bordereaux-analytics',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/components/BordereauxAnalyticsDashboard.vue'),
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/bordereaux-management/deadline-calendar',
      name: 'group-pricing-bordereaux-calendar',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/components/BordereauxCalendar.vue'),
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/bordereaux-management/inbound-submissions',
      name: 'group-pricing-bordereaux-inbound',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/components/BordereauxInboundSubmissions.vue'),
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/bordereaux-management/inbound-submissions/:submissionId',
      name: 'group-pricing-bordereaux-inbound-detail',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/components/BordereauxInboundSubmissionDetail.vue'),
      props: true,
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/bordereaux-management/reinsurer-tracking',
      name: 'group-pricing-bordereaux-reinsurer',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/components/BordereauxReinsurerTracking.vue'),
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/bordereaux-management/ri-treaties',
      name: 'group-pricing-ri-treaties',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/components/RITreatyManagement.vue'),
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/bordereaux-management/ri-kpi-dashboard',
      name: 'group-pricing-ri-kpi-dashboard',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/components/RIBordereauxKPIDashboard.vue'),
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/bordereaux-management/ri-submission-register',
      name: 'group-pricing-ri-submission-register',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/components/RISubmissionRegister.vue'),
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/bordereaux-management/ri-bordereaux',
      name: 'group-pricing-ri-bordereaux',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/components/RIBordereauxGeneration.vue'),
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/bordereaux-management/ri-claims',
      name: 'group-pricing-ri-claims',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/components/RIClaimsBordereaux.vue'),
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/bordereaux-management/ri-settlement',
      name: 'group-pricing-ri-settlement',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/components/RITechnicalAccounts.vue'),
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/bordereaux-management/claim-notifications',
      name: 'group-pricing-bordereaux-claim-notifications',
      component: () =>
        import('../screens/group_pricing/bordereaux_management/components/BordereauxClaimNotifications.vue'),
      meta: { required_permission: 'navigation:manage_bordereaux' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },

    // ── PHI Valuations ───────────────────────────────────────────────────────
    {
      path: '/group-pricing/phi/tables',
      name: 'group-pricing-phi-tables',
      component: () => import('../screens/group_pricing/phi/Tables.vue'),
      meta: { required_permission: 'navigation:view_phi_tables' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/phi/shock-settings',
      name: 'group-pricing-phi-shock-settings',
      component: () => import('../screens/group_pricing/phi/ShockSettings.vue'),
      meta: { required_permission: 'navigation:view_phi_shock_settings' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/phi/run-settings',
      name: 'group-pricing-phi-run-settings',
      component: () => import('../screens/group_pricing/phi/RunSettings.vue'),
      meta: { required_permission: 'navigation:view_phi_run_settings' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/phi/run-results',
      name: 'group-pricing-phi-run-results',
      component: () => import('../screens/group_pricing/phi/RunResults.vue'),
      meta: { required_permission: 'navigation:view_phi_run_results' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/phi/run-detail/:jobId',
      name: 'group-pricing-phi-run-detail',
      component: () => import('../screens/group_pricing/phi/RunDetail.vue'),
      props: true
    },

    // ── Finance: Premium Receipts (tabbed layout) ────────────────────────────
    {
      path: '/group-pricing/premiums',
      component: () =>
        import('../screens/group_pricing/premiums/PremiumReceiptsLayout.vue'),
      children: [
        {
          path: 'dashboard',
          name: 'group-pricing-premium-dashboard',
          component: () =>
            import('../screens/group_pricing/premiums/PremiumDashboard.vue'),
          meta: { required_permission: 'navigation:view_premium_dashboard' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        },
        {
          path: 'schedules',
          name: 'group-pricing-premium-schedules',
          component: () =>
            import('../screens/group_pricing/premiums/PremiumSchedules.vue'),
          meta: { required_permission: 'navigation:manage_premium_schedules' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        },
        {
          path: 'invoices',
          name: 'group-pricing-invoices',
          component: () =>
            import('../screens/group_pricing/premiums/Invoices.vue'),
          meta: { required_permission: 'navigation:manage_invoices' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        },
        {
          path: 'payments',
          name: 'group-pricing-payments',
          component: () =>
            import('../screens/group_pricing/premiums/Payments.vue'),
          meta: { required_permission: 'navigation:manage_payments' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        },
        {
          path: 'reconciliation',
          name: 'group-pricing-premium-reconciliation',
          component: () =>
            import('../screens/group_pricing/premiums/PremiumReconciliation.vue'),
          meta: {
            required_permission: 'navigation:manage_premium_reconciliation'
          },
          beforeEnter: (to, from) => checkPermissions(to, from)
        },
        {
          path: 'arrears',
          name: 'group-pricing-arrears',
          component: () =>
            import('../screens/group_pricing/premiums/ArrearsManagement.vue'),
          meta: { required_permission: 'navigation:manage_arrears' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        },
        {
          path: 'statements',
          name: 'group-pricing-statements',
          component: () =>
            import('../screens/group_pricing/premiums/Statements.vue'),
          meta: { required_permission: 'navigation:view_statements' },
          beforeEnter: (to, from) => checkPermissions(to, from)
        }
      ]
    },
    // Detail routes stay as top-level siblings so they render full-bleed
    // without the Premium Receipts tab bar.
    {
      path: '/group-pricing/premiums/schedules/:scheduleId',
      name: 'group-pricing-premium-schedule-detail',
      component: () =>
        import('../screens/group_pricing/premiums/PremiumScheduleDetail.vue'),
      props: true,
      meta: { required_permission: 'navigation:manage_premium_schedules' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/premiums/invoices/:invoiceId',
      name: 'group-pricing-invoice-detail',
      component: () =>
        import('../screens/group_pricing/premiums/InvoiceDetail.vue'),
      props: true,
      meta: { required_permission: 'navigation:manage_invoices' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },

    // ── Email System ─────────────────────────────────────────────────────────
    {
      path: '/group-pricing/email/settings',
      name: 'group-pricing-email-settings',
      component: () =>
        import('../screens/group_pricing/email/EmailSettings.vue'),
      meta: { required_permission: 'navigation:manage_email' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/email/templates',
      name: 'group-pricing-email-templates',
      component: () =>
        import('../screens/group_pricing/email/EmailTemplates.vue'),
      meta: { required_permission: 'navigation:manage_email' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/email/outbox',
      name: 'group-pricing-email-outbox',
      component: () => import('../screens/group_pricing/email/EmailOutbox.vue'),
      meta: { required_permission: 'navigation:manage_email' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },

    // ── User Management ──────────────────────────────────────────────────────
    {
      path: '/user-management-list',
      name: 'user-management-list',
      component: () => import('../screens/user_management/UserList.vue'),
      meta: { required_permission: 'navigation:manage_users' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/user-management-roles',
      name: 'user-management-roles',
      component: () => import('../screens/user_management/UserRoles.vue'),
      meta: { required_permission: 'navigation:manage_users' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },

    // ── Documentation ────────────────────────────────────────────────────────
    {
      path: '/documentation',
      name: 'documentation',
      component: () => import('../screens/Documentation.vue')
    },

    // ── Communications ──────────────────────────────────────────────────────
    {
      path: '/notifications',
      name: 'notification-center',
      component: () => import('../screens/NotificationCenter.vue')
    },
    {
      path: '/messages',
      name: 'messages-inbox',
      component: () => import('../screens/MessagesInbox.vue')
    },

    // ── Shell ────────────────────────────────────────────────────────────────
    {
      path: '/app-settings',
      name: 'app-settings',
      component: () => import('../screens/AppSettings.vue')
    },
    {
      path: '/no-entitlements',
      name: 'no-entitlements',
      component: () => import('../screens/NoEntitlements.vue')
    },
    {
      path: '/error',
      component: ErrorScreen,
      meta: { titleKey: 'title.error' }
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/group-pricing/dashboard'
    }
  ]
})

export default router

import ErrorScreen from '@/renderer/screens/ErrorScreen.vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import App from '../App.vue'
import AppLogin from '../AppLogin.vue'
import AppSetup from '../AppSetup.vue'
import { useGroupUserPermissionsStore } from '../store/group_user'

const checkPermissions = async (to: any, _from: any) => {
  const store = useGroupUserPermissionsStore()
  const required = to.meta.required_permission as string | undefined
  if (!required) return true
  // Wait for permissions to be fetched before checking
  await store.waitUntilLoaded()
  // If no role is assigned (no permissions), allow access (graceful degradation)
  if (!store.loaded || Object.keys(store.permissions).length === 0) return true
  if (store.hasPermission('system:admin')) return true
  if (store.hasPermission(required)) return true
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
      path: '/group-pricing/claims-management',
      name: 'group-pricing-claims-management',
      component: () =>
        import('../screens/group_pricing/claims_management/ClaimsManagement.vue'),
      meta: { required_permission: 'navigation:manage_claims' },
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
      meta: { required_permission: 'navigation:group-pricing_phi_tables' }
    },
    {
      path: '/group-pricing/phi/shock-settings',
      name: 'group-pricing-phi-shock-settings',
      component: () => import('../screens/group_pricing/phi/ShockSettings.vue'),
      meta: {
        required_permission: 'navigation:group-pricing_phi_shock_settings'
      }
    },
    {
      path: '/group-pricing/phi/run-settings',
      name: 'group-pricing-phi-run-settings',
      component: () => import('../screens/group_pricing/phi/RunSettings.vue'),
      meta: { required_permission: 'navigation:group-pricing_phi_run_settings' }
    },
    {
      path: '/group-pricing/phi/run-results',
      name: 'group-pricing-phi-run-results',
      component: () => import('../screens/group_pricing/phi/RunResults.vue'),
      meta: { required_permission: 'navigation:group-pricing_phi_run_results' }
    },
    {
      path: '/group-pricing/phi/run-detail/:jobId',
      name: 'group-pricing-phi-run-detail',
      component: () => import('../screens/group_pricing/phi/RunDetail.vue'),
      props: true
    },

    // ── Premium Management ───────────────────────────────────────────────────
    {
      path: '/group-pricing/premiums/dashboard',
      name: 'group-pricing-premium-dashboard',
      component: () =>
        import('../screens/group_pricing/premiums/PremiumDashboard.vue'),
      meta: { required_permission: 'navigation:view_premium_dashboard' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/premiums/schedules',
      name: 'group-pricing-premium-schedules',
      component: () =>
        import('../screens/group_pricing/premiums/PremiumSchedules.vue'),
      meta: { required_permission: 'navigation:manage_premium_schedules' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
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
      path: '/group-pricing/premiums/invoices',
      name: 'group-pricing-invoices',
      component: () => import('../screens/group_pricing/premiums/Invoices.vue'),
      meta: { required_permission: 'navigation:manage_invoices' },
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
    {
      path: '/group-pricing/premiums/payments',
      name: 'group-pricing-payments',
      component: () => import('../screens/group_pricing/premiums/Payments.vue'),
      meta: { required_permission: 'navigation:manage_payments' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/premiums/reconciliation',
      name: 'group-pricing-premium-reconciliation',
      component: () =>
        import('../screens/group_pricing/premiums/PremiumReconciliation.vue'),
      meta: { required_permission: 'navigation:manage_premium_reconciliation' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/premiums/arrears',
      name: 'group-pricing-arrears',
      component: () =>
        import('../screens/group_pricing/premiums/ArrearsManagement.vue'),
      meta: { required_permission: 'navigation:manage_arrears' },
      beforeEnter: (to, from) => checkPermissions(to, from)
    },
    {
      path: '/group-pricing/premiums/statements',
      name: 'group-pricing-statements',
      component: () =>
        import('../screens/group_pricing/premiums/Statements.vue'),
      meta: { required_permission: 'navigation:view_statements' },
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

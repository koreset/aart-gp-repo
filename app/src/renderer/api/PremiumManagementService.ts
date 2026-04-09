import Api from './Api'

const PremiumManagementService = {
  // ─── Contribution Configuration ────────────────────────────────────────────

  getContributionConfig(schemeId: number) {
    return Api.get(
      `/group-pricing/premiums/schemes/${schemeId}/contribution-config`
    )
  },

  saveContributionConfig(schemeId: number, config: object) {
    return Api.post(
      `/group-pricing/premiums/schemes/${schemeId}/contribution-config`,
      config
    )
  },

  // ─── Premium Schedules ──────────────────────────────────────────────────────

  getPremiumSchedules(params: {
    scheme_id?: number
    month?: number
    year?: number
    status?: string
  }) {
    return Api.get('/group-pricing/premiums/schedules', { params })
  },

  generateSchedule(data: { scheme_id: number; month: number; year: number }) {
    return Api.post('/group-pricing/premiums/schedules/generate', data)
  },

  generateAllSchedules(data: { month: number; year: number }) {
    return Api.post('/group-pricing/premiums/schedules/generate-all', data)
  },

  getScheduleDetail(scheduleId: number) {
    return Api.get(`/group-pricing/premiums/schedules/${scheduleId}`)
  },

  finalizeSchedule(scheduleId: number) {
    return Api.post(`/group-pricing/premiums/schedules/${scheduleId}/finalize`)
  },
  generateInvoice(scheduleId: number, dueDate?: string) {
    return Api.post(
      `/group-pricing/premiums/schedules/${scheduleId}/generate-invoice`,
      {
        due_date: dueDate ?? ''
      }
    )
  },
  reviewSchedule(scheduleId: number) {
    return Api.post(`/group-pricing/premiums/schedules/${scheduleId}/review`)
  },
  approveSchedule(scheduleId: number) {
    return Api.post(`/group-pricing/premiums/schedules/${scheduleId}/approve`)
  },
  returnScheduleToDraft(scheduleId: number) {
    return Api.post(
      `/group-pricing/premiums/schedules/${scheduleId}/return-to-draft`
    )
  },
  voidSchedule(scheduleId: number, data: { reason: string }) {
    return Api.post(
      `/group-pricing/premiums/schedules/${scheduleId}/void`,
      data
    )
  },
  cancelSchedule(scheduleId: number, data: { reason: string }) {
    return Api.post(
      `/group-pricing/premiums/schedules/${scheduleId}/cancel`,
      data
    )
  },
  regenerateSchedule(scheduleId: number) {
    return Api.post(
      `/group-pricing/premiums/schedules/${scheduleId}/regenerate`
    )
  },
  removeScheduleMember(scheduleId: number, rowId: number) {
    return Api.delete(
      `/group-pricing/premiums/schedules/${scheduleId}/members/${rowId}`
    )
  },
  updateScheduleMemberRow(scheduleId: number, rowId: number, data: object) {
    return Api.patch(
      `/group-pricing/premiums/schedules/${scheduleId}/members/${rowId}`,
      data
    )
  },

  exportSchedule(scheduleId: number) {
    return Api.get(`/group-pricing/premiums/schedules/${scheduleId}/export`, {
      responseType: 'blob'
    })
  },

  // ─── Invoices ───────────────────────────────────────────────────────────────

  getInvoices(params: {
    scheme_id?: number
    month?: number
    year?: number
    status?: string
    from?: string
    to?: string
  }) {
    return Api.get('/group-pricing/premiums/invoices', { params })
  },

  getInvoiceDetail(invoiceId: number) {
    return Api.get(`/group-pricing/premiums/invoices/${invoiceId}`)
  },

  getInvoicePdf(invoiceId: number) {
    return Api.get(`/group-pricing/premiums/invoices/${invoiceId}/pdf`, {
      responseType: 'blob'
    })
  },

  markInvoiceSent(invoiceId: number) {
    return Api.patch(`/group-pricing/premiums/invoices/${invoiceId}/mark-sent`)
  },

  emailInvoice(invoiceId: number, data: { message?: string }) {
    return Api.post(`/group-pricing/premiums/invoices/${invoiceId}/email`, data)
  },

  // ─── Payments ───────────────────────────────────────────────────────────────

  getPayments(params: {
    scheme_id?: number
    method?: string
    status?: string
    from?: string
    to?: string
  }) {
    return Api.get('/group-pricing/premiums/payments', { params })
  },

  recordPayment(data: {
    scheme_id: number
    invoice_id?: number | null
    payment_date: string
    method: string
    amount: number
    bank_reference: string
    notes?: string
  }) {
    return Api.post('/group-pricing/premiums/payments', data)
  },

  voidPayment(paymentId: number) {
    return Api.delete(`/group-pricing/premiums/payments/${paymentId}`)
  },

  bulkImportPayments(formData: FormData) {
    return Api.post('/group-pricing/premiums/payments/bulk-import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  downloadPaymentTemplate() {
    return Api.get('/group-pricing/premiums/payments/template', {
      responseType: 'blob'
    })
  },

  // ─── Reconciliation ─────────────────────────────────────────────────────────

  autoMatch() {
    return Api.post('/group-pricing/premiums/reconciliation/auto-match')
  },

  manualMatch(data: { payment_id: number; invoice_id: number }) {
    return Api.post('/group-pricing/premiums/reconciliation/manual-match', data)
  },

  createCreditNote(data: {
    scheme_id: number
    amount: number
    reason: string
    reference: string
  }) {
    return Api.post('/group-pricing/premiums/reconciliation/credit-note', data)
  },

  createDebitNote(data: {
    scheme_id: number
    amount: number
    reason: string
    reference: string
  }) {
    return Api.post('/group-pricing/premiums/reconciliation/debit-note', data)
  },

  // ─── Reconciliation v2 (ledger-based allocation engine) ─────────────────────

  getReconciliationSummary() {
    return Api.get('/group-pricing/premiums/reconciliation/v2/summary')
  },

  getReconciliationItems(params: {
    type?: string
    status?: string
    scheme_id?: number
  }) {
    return Api.get('/group-pricing/premiums/reconciliation/v2/items', {
      params
    })
  },

  runAutoMatchV2(data?: {
    scheme_id?: number
    rule_set?: string
    dry_run?: boolean
  }) {
    return Api.post(
      '/group-pricing/premiums/reconciliation/v2/auto-match',
      data || {}
    )
  },

  allocatePayment(data: {
    payment_id: number
    allocations: { invoice_id: number; amount: number }[]
    notes?: string
  }) {
    return Api.post('/group-pricing/premiums/reconciliation/v2/allocate', data)
  },

  reverseAllocations(data: { allocation_ids: number[]; reason: string }) {
    return Api.post('/group-pricing/premiums/reconciliation/v2/reverse', data)
  },

  writeOffBalance(data: {
    reconciliation_item_id: number
    amount: number
    reason: string
    invoice_id?: number
  }) {
    return Api.post('/group-pricing/premiums/reconciliation/v2/write-off', data)
  },

  refundOverpayment(data: {
    reconciliation_item_id: number
    amount: number
    reason: string
    refund_method: string
    bank_details?: string
  }) {
    return Api.post('/group-pricing/premiums/reconciliation/v2/refund', data)
  },

  getAllocationHistory(entityType: 'payment' | 'invoice', entityId: number) {
    return Api.get(
      `/group-pricing/premiums/reconciliation/v2/history/${entityType}/${entityId}`
    )
  },

  getReconciliationRuns(params?: { page?: number; page_size?: number }) {
    return Api.get('/group-pricing/premiums/reconciliation/v2/runs', { params })
  },

  getReconciliationRunDetail(runId: number) {
    return Api.get(`/group-pricing/premiums/reconciliation/v2/runs/${runId}`)
  },

  rollbackRun(runId: number) {
    return Api.post(
      `/group-pricing/premiums/reconciliation/v2/runs/${runId}/rollback`
    )
  },

  reassignReconciliationItem(
    itemId: number,
    data: { assigned_to: string; priority?: string }
  ) {
    return Api.patch(
      `/group-pricing/premiums/reconciliation/v2/items/${itemId}/reassign`,
      data
    )
  },

  suspendReconciliationItem(itemId: number, data: { reason: string }) {
    return Api.post(
      `/group-pricing/premiums/reconciliation/v2/items/${itemId}/suspend`,
      data
    )
  },

  getMatchingRules(ruleSet?: string) {
    return Api.get('/group-pricing/premiums/reconciliation/v2/matching-rules', {
      params: { rule_set: ruleSet || 'default' }
    })
  },

  saveMatchingRule(data: Record<string, unknown>) {
    return Api.post(
      '/group-pricing/premiums/reconciliation/v2/matching-rules',
      data
    )
  },

  deleteMatchingRule(ruleId: number) {
    return Api.delete(
      `/group-pricing/premiums/reconciliation/v2/matching-rules/${ruleId}`
    )
  },

  // ─── Arrears ────────────────────────────────────────────────────────────────

  getArrearsAging(params: { status?: string; min_days_overdue?: number }) {
    return Api.get('/group-pricing/premiums/arrears', { params })
  },

  sendReminder(schemeId: number, data: { message?: string }) {
    return Api.post(
      `/group-pricing/premiums/schemes/${schemeId}/send-reminder`,
      data
    )
  },

  recordPaymentPlan(
    schemeId: number,
    data: {
      instalments: Array<{ date: string; amount: number }>
      notes?: string
    }
  ) {
    return Api.post(
      `/group-pricing/premiums/schemes/${schemeId}/payment-plan`,
      data
    )
  },

  suspendCover(
    schemeId: number,
    data: { effective_date: string; reason: string }
  ) {
    return Api.post(`/group-pricing/premiums/schemes/${schemeId}/suspend`, data)
  },

  reinstateCover(
    schemeId: number,
    data: { reinstatement_date: string; back_premium?: number; notes?: string }
  ) {
    return Api.post(
      `/group-pricing/premiums/schemes/${schemeId}/reinstate`,
      data
    )
  },

  getArrearsHistory(schemeId: number) {
    return Api.get(
      `/group-pricing/premiums/schemes/${schemeId}/arrears-history`
    )
  },

  // ─── Statements ─────────────────────────────────────────────────────────────

  getEmployerStatement(
    schemeId: number,
    params: { from: string; to: string; detail?: boolean }
  ) {
    return Api.get(`/group-pricing/premiums/statements/employer/${schemeId}`, {
      params
    })
  },

  getBrokerStatement(brokerId: number, params: { from: string; to: string }) {
    return Api.get(`/group-pricing/premiums/statements/broker/${brokerId}`, {
      params
    })
  },

  emailStatement(data: {
    type: 'employer' | 'broker'
    recipient_id: number
    from_date: string
    to_date: string
  }) {
    return Api.post('/group-pricing/premiums/statements/email', data)
  },

  // ─── Dashboard ──────────────────────────────────────────────────────────────

  getPremiumDashboard(year: number) {
    return Api.get('/group-pricing/premiums/dashboard', { params: { year } })
  },

  getCollectionRate(year: number) {
    return Api.get('/group-pricing/premiums/dashboard/collection-rate', {
      params: { year }
    })
  },

  // ─── Schemes ────────────────────────────────────────────────────────────────

  getInForceSchemes() {
    return Api.get('/group-pricing/schemes/in-force-v2')
  },

  // ─── Cross-Module Integration ───────────────────────────────────────────────

  getSchemePaymentSummary(schemeId: number) {
    return Api.get(
      `/group-pricing/premiums/schemes/${schemeId}/payment-summary`
    )
  },

  getMemberContribution(memberId: number) {
    return Api.get(`/group-pricing/premiums/members/${memberId}/contribution`)
  },

  getOutstandingInvoicesForScheme(schemeId: number) {
    return Api.get('/group-pricing/premiums/invoices', {
      params: { scheme_id: schemeId, status: 'unpaid' }
    })
  },

  // ─── Coverage Matrix ──────────────────────────────────────────────────────

  getScheduleCoverageMatrix() {
    return Api.get('/group-pricing/premiums/schedules/coverage-matrix')
  }
}

export default PremiumManagementService

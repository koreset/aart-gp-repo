import type { InjectionKey, Ref } from 'vue'

export interface RiskFlags {
  banking_change_30d?: boolean
  contestable?: boolean
  recent_reinstatement?: boolean
  fraud_risk_level?: string
  // Phase 5: cross-claim duplicate signals. Refs arrays carry short
  // pre-formatted strings ready for the tooltip — never raw rows.
  id_paid_before?: boolean
  id_paid_before_refs?: string[]
  account_used_before?: boolean
  account_used_before_refs?: string[]
}

export interface BankVerificationStatus {
  has_result: boolean
  status: 'complete' | 'failed' | 'pending' | ''
  verified: boolean
  verified_at?: string | null
  stale: boolean
  stale_reason?: string
  provider_request_id?: string
  last_attempt: number
}

export interface ScheduleItem {
  id: number
  claim_id: number
  claim_number: string
  member_name: string
  member_id_number: string
  benefit_name: string
  scheme_name: string
  scheme_id: number
  claim_amount: number
  gross_amount?: number
  premium_arrears_deduction?: number
  policy_loan_deduction?: number
  tax_withheld?: number
  net_payable?: number
  beneficiary_name?: string
  beneficiary_id_number?: string
  risk_flags?: RiskFlags | string | null
  approval_reference?: string
  line_status?: 'pending' | 'verified' | 'queried' | 'rejected' | string
  verified_by?: string
  verified_at?: string
  query_reason_code?: string
  query_notes?: string
  queried_by?: string
  queried_at?: string
  // Phase 3
  reinsurance_recovery_required?: boolean
  reinsurance_recovery_amount?: number
  reinsurance_recovery_raised_by?: string
  reinsurance_recovery_raised_at?: string
  duplicate_beneficiary_flag?: boolean
  duplicate_beneficiary_cleared?: boolean
  // Phase 5 — pre-authorisation amount drift
  approved_amount_snapshot?: number
  amount_drift_resolved?: boolean
  amount_drift_resolved_by?: string
  amount_drift_resolved_at?: string
  bank_name?: string
  bank_branch_code?: string
  bank_account_number?: string
  bank_account_type?: string
  account_holder_name?: string
}

export interface ScheduleQuery {
  id: number
  schedule_id: number
  schedule_item_id: number
  claim_id: number
  claim_number: string
  reason_code: string
  notes: string
  outcome: string
  raised_by: string
  raised_at: string
  resolution_notes?: string
  resolved_by?: string
  resolved_at?: string | null
}

export interface ScheduleAuditRow {
  id: number
  schedule_id: number
  from_status: string
  to_status: string
  actor: string
  notes: string
  changed_at: string
}

export interface TaxCertificate {
  id: number
  schedule_id: number
  schedule_item_id: number
  claim_id: number
  claim_number: string
  benefit_name: string
  beneficiary_name: string
  beneficiary_id_number: string
  tax_year: number
  gross_amount: number
  tax_withheld: number
  certificate_ref: string
  file_name: string
  content_type: string
  generated_by?: string
  generated_at?: string
}

export interface SanctionsScreening {
  id: number
  schedule_id: number
  schedule_item_id: number
  claim_id: number
  provider: string
  status: 'pending' | 'clear' | 'hit' | 'manual_clear' | 'skipped' | string
  provider_ref?: string
  hit_summary?: string
  notes?: string
  screened_by?: string
  screened_at?: string
  cleared_by?: string
  cleared_at?: string
}

export interface PaymentProof {
  id: number
  schedule_id: number
  file_name: string
  content_type: string
  size_bytes: number
  notes: string
  uploaded_by: string
  uploaded_at: string
}

export interface PaymentSchedule {
  id: number
  schedule_number: string
  description: string
  status: string
  total_amount: number
  gross_total?: number
  deductions_total?: number
  net_total?: number
  claims_count: number
  locked_at?: string
  head_of_claims_signed_by?: string
  head_of_claims_signed_at?: string
  finance_review_started_by?: string
  finance_review_started_at?: string
  finance_first_auth_by?: string
  finance_first_auth_at?: string
  finance_second_auth_by?: string
  finance_second_auth_at?: string
  submitted_to_bank_at?: string
  archived_by?: string
  archived_at?: string
  exported_at?: string
  exported_by?: string
  acb_file_generated?: boolean
  acb_generated_at?: string
  acb_generated_by?: string
  created_by: string
  created_at: string
  items: ScheduleItem[]
  proof_of_payments: PaymentProof[]
}

export interface BankProfile {
  id: number
  profile_name: string
  bank_name: string
  user_code: string
  user_branch_code: string
  user_account_number: string
  user_account_type: string
  bank_type_code: string
  service_type: string
  generation_number: number
  is_active: boolean
}

export interface ACBFile {
  id: number
  schedule_id: number
  file_name: string
  action_date: string
  transaction_count: number
  total_amount: number
  status: string
  is_retry: boolean
  generated_by: string
  generated_at: string
}

export interface ReconResult {
  id: number
  claim_number: string
  account_number: string
  amount: number
  status: string
  failure_reason: string
  bank_reference: string
}

export interface ReconSummary {
  total_transactions: number
  paid: number
  failed: number
  unmatched: number
  total_paid: number
  total_failed: number
}

export interface PaymentScheduleContext {
  // Shared state
  schedule: Ref<PaymentSchedule | null>

  // Helpers
  formatCurrency: (val: number) => string
  formatDate: (val?: string) => string
  hasPermission: (slug: string) => boolean
  notify: (message: string, color?: string) => void

  // ACB tab
  acbFiles: Ref<ACBFile[]>
  loadingACBFiles: Ref<boolean>
  downloadingACB: Ref<number | null>
  loadACBFiles: () => Promise<void>
  downloadACBFile: (acb: ACBFile) => Promise<void>

  // Reconciliation tab
  reconResults: Ref<ReconResult[]>
  reconSummary: Ref<ReconSummary | null>
  loadingRecon: Ref<boolean>
  retrying: Ref<boolean>
  loadReconData: () => Promise<void>
  retryFailed: () => Promise<void>
  reconStatusColor: (status: string) => string

  // Proofs tab
  downloadingProof: Ref<number | null>
  downloadProof: (proof: PaymentProof) => Promise<void>

  // Lifecycle (Phase 1)
  signOff: () => Promise<void>
  startFinanceReview: () => Promise<void>
  verifyLineItem: (itemId: number) => Promise<void>
  queryLineItem: (
    itemId: number,
    reasonCode: string,
    notes: string
  ) => Promise<void>
  rejectLineItem: (
    itemId: number,
    reasonCode: string,
    notes: string
  ) => Promise<void>
  authoriseFirst: () => Promise<void>
  authoriseSecond: () => Promise<void>
  archive: () => Promise<void>
  refreshSchedule: () => Promise<void>

  // Queries tab
  queries: Ref<ScheduleQuery[]>
  loadingQueries: Ref<boolean>
  loadQueries: () => Promise<void>

  // Sanctions / reinsurance / duplicates (Phase 3)
  sanctions: Ref<SanctionsScreening[]>
  loadSanctions: () => Promise<void>
  screenLineItem: (itemId: number) => Promise<void>
  recordSanctionsOutcome: (
    itemId: number,
    status: string,
    notes: string
  ) => Promise<void>
  setReinsuranceRecovery: (
    itemId: number,
    required: boolean,
    amount: number
  ) => Promise<void>
  confirmReinsuranceRaised: (itemId: number) => Promise<void>
  clearDuplicateBeneficiary: (itemId: number) => Promise<void>

  // Tax certificates (Phase 4)
  taxCertificates: Ref<TaxCertificate[]>
  loadTaxCertificates: () => Promise<void>
  downloadTaxCertificate: (cert: TaxCertificate) => Promise<void>
}

export const PAYMENT_SCHEDULE_CONTEXT: InjectionKey<PaymentScheduleContext> =
  Symbol('PaymentScheduleContext')

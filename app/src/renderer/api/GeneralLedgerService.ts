import Api from '@/renderer/api/Api'

// Operational General Ledger — chart of accounts, accounting periods,
// posting rules, journals (draft/submit/approve/post + reversal),
// trial balance, account ledger drill, bank reconciliation, and audit log.
// Mirrors the /gl and /bank-accounts route groups in api/routes/routes.go.
//
// All mutating endpoints follow a maker/checker pattern: one user requests
// the change, a different user approves. Self-approval returns HTTP 409.
//
// Every API response is wrapped by controllers.OK in
// PremiumResponse { success, data, message }; this service unwraps that
// envelope so callers receive the payload directly.

export type ApprovalStatus =
  | 'active'
  | 'pending_create'
  | 'pending_update'
  | 'pending_deactivate'
  | 'pending_delete'

export interface GLAccount {
  id?: number
  code: string
  name: string
  account_type: 'asset' | 'liability' | 'equity' | 'income' | 'expense'
  normal_balance: 'debit' | 'credit'
  parent_id?: number | null
  is_active?: boolean
  description?: string
  created_at?: string
  created_by?: string
  updated_at?: string
  updated_by?: string
  approval_status?: ApprovalStatus
  pending_change_json?: string
  pending_requested_by?: string
  pending_requested_at?: string | null
}

export interface AccountingPeriod {
  id?: number
  name?: string
  start_date: string
  end_date?: string
  status?: 'open' | 'close_requested' | 'closed'
  close_requested_by?: string
  close_requested_at?: string | null
  closed_at?: string | null
  closed_by?: string
}

export interface PostingRule {
  id?: number
  event_key: string
  debit_account_id: number
  credit_account_id: number
  is_active?: boolean
  notes?: string
  created_at?: string
  created_by?: string
  updated_at?: string
  updated_by?: string
  approval_status?: ApprovalStatus
  pending_change_json?: string
  pending_requested_by?: string
  pending_requested_at?: string | null
}

export interface JournalLine {
  id?: number
  entry_id?: number
  account_id: number
  debit: number
  credit: number
  description?: string
  scheme_id?: number
  cost_centre?: string
  line_order?: number
}

export type JournalStatus =
  | 'draft'
  | 'submitted'
  | 'approved'
  | 'posted'
  | 'reversal_pending'
  | 'reversal_approved'
  | 'reversed'

export interface JournalEntry {
  id?: number
  entry_number?: string
  period_id?: number
  status?: JournalStatus
  posted_at?: string
  posted_by?: string
  source_type?: string
  source_id?: number
  description?: string
  is_reversed?: boolean
  reversed_by_entry_id?: number | null
  total_debit?: number
  total_credit?: number
  lines?: JournalLine[]
  created_at?: string
  created_by?: string
  updated_at?: string
  updated_by?: string
  submitted_by?: string
  submitted_at?: string | null
  approved_by?: string
  approved_at?: string | null
  reversal_reason?: string
  reversal_requested_by?: string
  reversal_requested_at?: string | null
  reversal_approved_by?: string
  reversal_approved_at?: string | null
}

export interface ManualJournalLineInput {
  account_id: number
  debit: number
  credit: number
  description?: string
  scheme_id?: number
  cost_centre?: string
}

export interface ManualJournalRequest {
  description: string
  period_id?: number
  lines: ManualJournalLineInput[]
}

export interface TrialBalanceRow {
  account_id: number
  account_code: string
  account_name: string
  account_type: string
  normal_balance: string
  total_debit: number
  total_credit: number
  net_balance: number
}

export interface LedgerRow {
  entry_id: number
  entry_number: string
  posted_at: string
  source_type: string
  source_id: number
  description: string
  line_description: string
  debit: number
  credit: number
  running_balance: number
}

export interface BankAccount {
  id?: number
  code: string
  name: string
  bank_name?: string
  account_number?: string
  gl_account_id: number
  currency?: string
  is_active?: boolean
  created_at?: string
  created_by?: string
  updated_at?: string
  updated_by?: string
  approval_status?: ApprovalStatus
  pending_change_json?: string
  pending_requested_by?: string
  pending_requested_at?: string | null
}

export interface BankStatementLine {
  id: number
  bank_account_id: number
  statement_date: string
  value_date?: string | null
  description?: string
  amount: number
  reference?: string
  import_batch_id?: string
  imported_by?: string
  matched_journal_line_id?: number | null
  match_status: 'unmatched' | 'matched' | 'ignored'
  matched_at?: string | null
  matched_by?: string
  review_status: 'not_required' | 'pending_review' | 'reviewed' | 'rejected'
  reviewed_by?: string
  reviewed_at?: string | null
  review_notes?: string
}

export interface StatementImportRow {
  statement_date: string
  value_date?: string
  description?: string
  amount: number
  reference?: string
}

export interface ListJournalsParams {
  period_id?: number
  source_type?: string
  status?: JournalStatus | string
  account_id?: number
  from?: string
  to?: string
  limit?: number
}

export interface GLAuditLogEntry {
  id: number
  event_type: string
  object_type:
    | 'journal_entry'
    | 'accounting_period'
    | 'gl_account'
    | 'posting_rule'
    | 'bank_account'
    | 'bank_statement_line'
    | string
  object_id: number
  object_name: string
  changed_by: string
  changed_at: string
  details?: string
}

export interface ListAuditLogParams {
  event_type?: string
  object_type?: string
  object_id?: number
  changed_by?: string
  from?: string
  to?: string
  limit?: number
}

// unwrap pulls `data` out of the controllers.OK envelope
// (`{ success, data, message }`). Defensive against three shapes:
//   * the wrapped envelope — most endpoints
//   * a bare payload (legacy endpoints)
//   * an error envelope with no `data` field — returns undefined so the
//     caller's `|| []` fallback can kick in instead of leaking the envelope
//     into a component that expects an array
function unwrap<T>(r: { data: unknown }): T | undefined {
  const body = r.data as { success?: boolean; data?: T } | T
  if (body && typeof body === 'object' && 'success' in (body as object)) {
    return (body as { data?: T }).data
  }
  return body as T
}

export default {
  // ─── Chart of accounts (maker/checker) ──────────────────────────────
  listAccounts(): Promise<GLAccount[]> {
    return Api.get('/gl/accounts').then((r) => unwrap<GLAccount[]>(r) || [])
  },
  getAccount(id: number): Promise<GLAccount> {
    return Api.get(`/gl/accounts/${id}`).then((r) => unwrap<GLAccount>(r) as GLAccount)
  },
  requestCreateAccount(payload: GLAccount): Promise<GLAccount> {
    return Api.post('/gl/accounts', payload).then((r) => unwrap<GLAccount>(r) as GLAccount)
  },
  requestUpdateAccount(id: number, payload: GLAccount): Promise<GLAccount> {
    return Api.put(`/gl/accounts/${id}`, payload).then((r) => unwrap<GLAccount>(r) as GLAccount)
  },
  requestDeactivateAccount(id: number): Promise<GLAccount> {
    return Api.delete(`/gl/accounts/${id}`).then((r) => unwrap<GLAccount>(r) as GLAccount)
  },
  approveAccountChange(id: number, notes?: string): Promise<GLAccount> {
    return Api.post(`/gl/accounts/${id}/approve-change`, { notes: notes || '' }).then(
      (r) => unwrap<GLAccount>(r) as GLAccount
    )
  },
  getAccountLedger(
    id: number,
    params: { from?: string; to?: string } = {}
  ): Promise<LedgerRow[]> {
    return Api.get(`/gl/accounts/${id}/ledger`, { params }).then(
      (r) => unwrap<LedgerRow[]>(r) || []
    )
  },

  // ─── Periods (two-step close) ──────────────────────────────────────
  listPeriods(): Promise<AccountingPeriod[]> {
    return Api.get('/gl/periods').then((r) => unwrap<AccountingPeriod[]>(r) || [])
  },
  createPeriod(payload: AccountingPeriod): Promise<AccountingPeriod> {
    return Api.post('/gl/periods', payload).then(
      (r) => unwrap<AccountingPeriod>(r) as AccountingPeriod
    )
  },
  requestClosePeriod(id: number): Promise<AccountingPeriod> {
    return Api.post(`/gl/periods/${id}/request-close`).then(
      (r) => unwrap<AccountingPeriod>(r) as AccountingPeriod
    )
  },
  closePeriod(id: number): Promise<AccountingPeriod> {
    return Api.post(`/gl/periods/${id}/close`).then(
      (r) => unwrap<AccountingPeriod>(r) as AccountingPeriod
    )
  },

  // ─── Posting rules (maker/checker) ─────────────────────────────────
  listPostingRules(): Promise<PostingRule[]> {
    return Api.get('/gl/posting-rules').then((r) => unwrap<PostingRule[]>(r) || [])
  },
  requestCreatePostingRule(payload: PostingRule): Promise<PostingRule> {
    return Api.post('/gl/posting-rules', payload).then(
      (r) => unwrap<PostingRule>(r) as PostingRule
    )
  },
  requestUpdatePostingRule(id: number, payload: PostingRule): Promise<PostingRule> {
    return Api.put(`/gl/posting-rules/${id}`, payload).then(
      (r) => unwrap<PostingRule>(r) as PostingRule
    )
  },
  requestDeletePostingRule(id: number): Promise<PostingRule> {
    return Api.delete(`/gl/posting-rules/${id}`).then(
      (r) => unwrap<PostingRule>(r) as PostingRule
    )
  },
  approvePostingRuleChange(id: number, notes?: string): Promise<PostingRule> {
    return Api.post(`/gl/posting-rules/${id}/approve-change`, { notes: notes || '' }).then(
      (r) => unwrap<PostingRule>(r) as PostingRule
    )
  },

  // ─── Journals (draft → submit → approve → post; request → approve reverse) ─
  listJournals(params: ListJournalsParams = {}): Promise<JournalEntry[]> {
    return Api.get('/gl/journals', { params }).then(
      (r) => unwrap<JournalEntry[]>(r) || []
    )
  },
  getJournalEntry(id: number): Promise<JournalEntry> {
    return Api.get(`/gl/journals/${id}`).then((r) => unwrap<JournalEntry>(r) as JournalEntry)
  },
  draftManualJournal(payload: ManualJournalRequest): Promise<JournalEntry> {
    return Api.post('/gl/journals', payload).then((r) => unwrap<JournalEntry>(r) as JournalEntry)
  },
  updateDraftJournal(id: number, payload: ManualJournalRequest): Promise<JournalEntry> {
    return Api.put(`/gl/journals/${id}`, payload).then(
      (r) => unwrap<JournalEntry>(r) as JournalEntry
    )
  },
  discardDraftJournal(id: number): Promise<unknown> {
    return Api.delete(`/gl/journals/${id}`).then((r) => unwrap(r))
  },
  submitManualJournal(id: number): Promise<JournalEntry> {
    return Api.post(`/gl/journals/${id}/submit`).then(
      (r) => unwrap<JournalEntry>(r) as JournalEntry
    )
  },
  approveManualJournal(id: number): Promise<JournalEntry> {
    return Api.post(`/gl/journals/${id}/approve`).then(
      (r) => unwrap<JournalEntry>(r) as JournalEntry
    )
  },
  postApprovedJournal(id: number): Promise<JournalEntry> {
    return Api.post(`/gl/journals/${id}/post`).then(
      (r) => unwrap<JournalEntry>(r) as JournalEntry
    )
  },
  requestReverseJournal(id: number, reason: string): Promise<JournalEntry> {
    return Api.post(`/gl/journals/${id}/request-reverse`, { reason }).then(
      (r) => unwrap<JournalEntry>(r) as JournalEntry
    )
  },
  approveReverseJournal(id: number): Promise<JournalEntry> {
    return Api.post(`/gl/journals/${id}/approve-reverse`).then(
      (r) => unwrap<JournalEntry>(r) as JournalEntry
    )
  },

  // ─── Reports ───────────────────────────────────────────────────────
  getTrialBalance(periodId?: number): Promise<TrialBalanceRow[]> {
    return Api.get('/gl/trial-balance', {
      params: periodId ? { period_id: periodId } : {}
    }).then((r) => unwrap<TrialBalanceRow[]>(r) || [])
  },

  // ─── Bank accounts + reconciliation (maker/checker + reviewer sign-off) ──
  listBankAccounts(): Promise<BankAccount[]> {
    return Api.get('/bank-accounts').then((r) => unwrap<BankAccount[]>(r) || [])
  },
  requestCreateBankAccount(payload: BankAccount): Promise<BankAccount> {
    return Api.post('/bank-accounts', payload).then(
      (r) => unwrap<BankAccount>(r) as BankAccount
    )
  },
  requestUpdateBankAccount(id: number, payload: BankAccount): Promise<BankAccount> {
    return Api.put(`/bank-accounts/${id}`, payload).then(
      (r) => unwrap<BankAccount>(r) as BankAccount
    )
  },
  requestDeactivateBankAccount(id: number): Promise<BankAccount> {
    return Api.delete(`/bank-accounts/${id}`).then(
      (r) => unwrap<BankAccount>(r) as BankAccount
    )
  },
  approveBankAccountChange(id: number, notes?: string): Promise<BankAccount> {
    return Api.post(`/bank-accounts/${id}/approve-change`, { notes: notes || '' }).then(
      (r) => unwrap<BankAccount>(r) as BankAccount
    )
  },
  importStatement(
    bankAccountId: number,
    rows: StatementImportRow[]
  ): Promise<{ batch_id: string; imported: number }> {
    return Api.post(`/bank-accounts/${bankAccountId}/statement`, {
      bank_account_id: bankAccountId,
      rows
    }).then(
      (r) =>
        unwrap<{ batch_id: string; imported: number }>(r) as {
          batch_id: string
          imported: number
        }
    )
  },
  listStatementLines(
    bankAccountId: number,
    status?: string
  ): Promise<BankStatementLine[]> {
    return Api.get(`/bank-accounts/${bankAccountId}/statement-lines`, {
      params: status ? { status } : {}
    }).then((r) => unwrap<BankStatementLine[]>(r) || [])
  },
  matchStatementLine(
    statementLineId: number,
    journalLineId: number
  ): Promise<unknown> {
    return Api.post('/bank-accounts/statement-lines/match', {
      statement_line_id: statementLineId,
      journal_line_id: journalLineId
    }).then((r) => unwrap(r))
  },
  ignoreStatementLine(statementLineId: number): Promise<unknown> {
    return Api.post(
      `/bank-accounts/statement-lines/${statementLineId}/ignore`
    ).then((r) => unwrap(r))
  },
  reviewStatementLine(
    statementLineId: number,
    outcome: 'reviewed' | 'rejected',
    notes?: string
  ): Promise<BankStatementLine> {
    return Api.post(`/bank-accounts/statement-lines/${statementLineId}/review`, {
      outcome,
      notes: notes || ''
    }).then((r) => unwrap<BankStatementLine>(r) as BankStatementLine)
  },

  // ─── Audit log ─────────────────────────────────────────────────────
  listAuditLog(params: ListAuditLogParams = {}): Promise<GLAuditLogEntry[]> {
    return Api.get('/gl/audit-log', { params }).then(
      (r) => unwrap<GLAuditLogEntry[]>(r) || []
    )
  },
  listAuditLogUsers(): Promise<string[]> {
    return Api.get('/gl/audit-log/users').then(
      (r) => unwrap<string[]>(r) || []
    )
  }
}

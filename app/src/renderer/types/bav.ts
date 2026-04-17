// Bank Account Verification — canonical types shared with api/services/bav.
// The shape here mirrors the JSON emitted by the v2 endpoint
// (POST /v2/group-pricing/claims/verify-bank-account). TriState values are
// deliberately 'yes' | 'no' | 'unknown' so that unknowns render distinctly
// from outright 'no' answers.

export type TriState = 'yes' | 'no' | 'unknown'

export type VerifyStatus = 'complete' | 'pending' | 'failed'

export interface VerifyResult {
  status: VerifyStatus
  verified: boolean
  summary: string
  accountFound: TriState
  accountOpen: TriState
  identityMatch: TriState
  accountTypeMatch: TriState
  acceptsCredits: TriState
  acceptsDebits: TriState
  provider: string
  providerRequestId: string
  providerJobId?: string
}

export interface VerifyBankAccountRequest {
  first_name: string
  surname: string
  identity_number: string
  identity_type?: string
  bank_account_number: string
  bank_branch_code: string
  bank_account_type: string
}

export interface TriStateIcon {
  icon: string
  color: string
  label: string
}

// triStateIcon maps a TriState to the Vuetify icon + theme colour + label
// the UI should use when rendering a tick, cross or question mark.
export function triStateIcon(t: TriState): TriStateIcon {
  switch (t) {
    case 'yes':
      return { icon: 'mdi-check-circle', color: 'success', label: 'Yes' }
    case 'no':
      return { icon: 'mdi-close-circle', color: 'error', label: 'No' }
    default:
      return { icon: 'mdi-help-circle', color: 'warning', label: 'Unknown' }
  }
}

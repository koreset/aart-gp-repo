// Status helpers for group-scheme claims. These are kept in sync with the
// backend's claimEditableStatuses / claimSubmittableStatuses maps in
// api/services/group_pricing.go — when adding a new status, update both
// sides together.

export const EDITABLE_CLAIM_STATUSES = [
  'draft',
  'pending_assessment',
  'under_assessment',
  'additional_info_required'
] as const

export const SUBMITTABLE_CLAIM_STATUSES = [
  'draft',
  'additional_info_required'
] as const

export type EditableClaimStatus = (typeof EDITABLE_CLAIM_STATUSES)[number]
export type SubmittableClaimStatus = (typeof SUBMITTABLE_CLAIM_STATUSES)[number]

export const isEditableClaimStatus = (status?: string | null): boolean =>
  !!status &&
  (EDITABLE_CLAIM_STATUSES as readonly string[]).includes(status)

export const isSubmittableClaimStatus = (status?: string | null): boolean =>
  !!status &&
  (SUBMITTABLE_CLAIM_STATUSES as readonly string[]).includes(status)

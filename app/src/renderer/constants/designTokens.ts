export const STATUS_COLORS: Record<string, { hex: string; label: string }> = {
  draft: { hex: '#9E9E9E', label: 'Draft' },
  pending: { hex: '#FF9800', label: 'Pending' },
  reviewed: { hex: '#009688', label: 'Reviewed' },
  approved: { hex: '#8BC34A', label: 'Approved' },
  active: { hex: '#4CAF50', label: 'Active' },
  finalized: { hex: '#2196F3', label: 'Finalized' },
  submitted: { hex: '#3F51B5', label: 'Submitted' },
  overdue: { hex: '#F44336', label: 'Overdue' },
  voided: { hex: '#B71C1C', label: 'Voided' },
  suspended: { hex: '#E91E63', label: 'Suspended' },
  matched: { hex: '#4CAF50', label: 'Matched' },
  unmatched: { hex: '#FF9800', label: 'Unmatched' },
  generated: { hex: '#1976D2', label: 'Generated' },
  invoiced: { hex: '#00BCD4', label: 'Invoiced' },
  paid: { hex: '#4CAF50', label: 'Paid' },
  partial: { hex: '#FF9800', label: 'Partial' },
  void: { hex: '#B71C1C', label: 'Void' },
  cancelled: { hex: '#FF5722', label: 'Cancelled' },
  accepted: { hex: '#4CAF50', label: 'Accepted' },
  rejected: { hex: '#F44336', label: 'Rejected' },
  queried: { hex: '#FF9800', label: 'Queried' },
  failed: { hex: '#D32F2F', label: 'Failed' },
  in_force: { hex: '#4CAF50', label: 'In Force' },
  lapsed: { hex: '#9E9E9E', label: 'Lapsed' },
  pending_receipt: { hex: '#9E9E9E', label: 'Pending Receipt' },
  received: { hex: '#2196F3', label: 'Received' },
  under_review: { hex: '#FF9800', label: 'Under Review' },
  queries_raised: { hex: '#FFC107', label: 'Queries Raised' },
  under_assessment: { hex: '#FF9800', label: 'Under Assessment' },
  additional_info_required: { hex: '#FB8C00', label: 'Info Required' },
  referred_to_committee: { hex: '#7B1FA2', label: 'Referred to Committee' },
  declined: { hex: '#D32F2F', label: 'Declined' },
  payment_failed: { hex: '#D32F2F', label: 'Payment Failed' },
  submitted_for_payment: { hex: '#1976D2', label: 'Submitted for Payment' }
}

export const CHART_COLORS = [
  '#003F58',
  '#2196F3',
  '#4CAF50',
  '#FF9800',
  '#9C27B0',
  '#00BCD4'
]

export const DIALOG_SIZES = {
  xs: '400px',
  sm: '560px',
  md: '800px',
  lg: '1000px',
  xl: '1200px'
} as const

import { computed, type Ref } from 'vue'

export interface PipelineStep {
  status: string
  label: string
  sub: string
}

export const SCHEDULE_PIPELINE_STEPS: PipelineStep[] = [
  { status: 'draft', label: 'Draft', sub: 'Created' },
  {
    status: 'claims_signed_off',
    label: 'Payment Schedule Signed Off',
    sub: 'Head of Claims'
  },
  { status: 'finance_in_review', label: 'Finance Review', sub: 'In review' },
  {
    status: 'finance_first_authorised',
    label: '1st Authorisation',
    sub: 'Finance'
  },
  {
    status: 'finance_second_authorised',
    label: '2nd Authorisation',
    sub: 'Finance'
  },
  { status: 'submitted_to_bank', label: 'Submitted to Bank', sub: 'ACB run' },
  { status: 'confirmed', label: 'Paid / Confirmed', sub: 'Proof uploaded' }
]

const PIPELINE_STATUSES = SCHEDULE_PIPELINE_STEPS.map((s) => s.status)

export const CLAIMS_OWNED_STATUSES: readonly string[] = [
  'draft',
  'claims_signed_off'
]

export function isClaimsStep(status: string): boolean {
  return CLAIMS_OWNED_STATUSES.includes(status)
}

export interface SchedulePipelineSummary {
  pipelineSteps: PipelineStep[]
  currentStepIndex: Ref<number>
  stepTone: (idx: number) => 'success' | 'current' | 'muted'
  statusLabel: (status: string | undefined | null) => string
}

const STATUS_LABELS: Record<string, string> = SCHEDULE_PIPELINE_STEPS.reduce(
  (acc, step) => {
    acc[step.status] = step.label
    return acc
  },
  { archived: 'Archived', submitted: 'Submitted to Bank' } as Record<
    string,
    string
  >
)

export function useScheduleLifecycle(
  status: Ref<string | undefined | null>
): SchedulePipelineSummary {
  const currentStepIndex = computed(() => {
    const s = status.value
    if (!s) return -1
    const idx = PIPELINE_STATUSES.indexOf(s)
    if (idx >= 0) return idx
    if (s === 'submitted') return 5
    if (s === 'archived') return PIPELINE_STATUSES.length
    return -1
  })

  function stepTone(idx: number): 'success' | 'current' | 'muted' {
    if (idx < currentStepIndex.value) return 'success'
    if (idx === currentStepIndex.value) return 'current'
    return 'muted'
  }

  function statusLabel(s: string | undefined | null): string {
    if (!s) return '—'
    return STATUS_LABELS[s] ?? s
  }

  return {
    pipelineSteps: SCHEDULE_PIPELINE_STEPS,
    currentStepIndex,
    stepTone,
    statusLabel
  }
}

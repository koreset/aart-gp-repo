import { ref, computed } from 'vue'
import { useWebSocket } from './useWebSocket'

export interface CalculationProgress {
  quoteId: string
  jobId?: number
  totalCategories: number
  completedCategories: number
  currentCategory: string
  phase: string
  progress: number
  queuePosition?: number
}

const progress = ref<CalculationProgress | null>(null)
const isCalculating = ref(false)

const phaseLabel = computed(() => {
  if (!progress.value) return ''
  switch (progress.value.phase) {
    case 'queued':
      return progress.value.queuePosition
        ? `Queued — position ${progress.value.queuePosition}`
        : 'Queued'
    case 'loading_data':
      return 'Loading data'
    case 'rating_members':
      return 'Rating members'
    case 'saving_results':
      return 'Saving results'
    case 'category_done':
      return 'Category complete'
    case 'completed':
      return 'Completed'
    case 'failed':
      return 'Calculation failed'
    default:
      return progress.value.phase
  }
})

const isQueued = computed(() => progress.value?.phase === 'queued')

const progressPercent = computed(() => {
  return Math.round(progress.value?.progress ?? 0)
})

function handleProgress(payload: CalculationProgress) {
  // Only handle events for the quote we're tracking
  if (progress.value && payload.quoteId !== progress.value.quoteId) return

  progress.value = payload

  if (payload.phase === 'completed') {
    setTimeout(() => {
      if (progress.value?.phase === 'completed') {
        isCalculating.value = false
        progress.value = null
      }
    }, 1500)
  }

  if (payload.phase === 'failed') {
    setTimeout(() => {
      if (progress.value?.phase === 'failed') {
        isCalculating.value = false
        progress.value = null
      }
    }, 3000)
  }
}

let registered = false

function startTracking(quoteId: string) {
  isCalculating.value = true
  progress.value = {
    quoteId,
    totalCategories: 0,
    completedCategories: 0,
    currentCategory: '',
    phase: 'queued',
    progress: 0
  }

  if (!registered) {
    const { on } = useWebSocket()
    on('calculation_progress', handleProgress)
    registered = true
  }
}

function stopTracking() {
  isCalculating.value = false
  progress.value = null
}

export function useCalculationProgress() {
  return {
    progress,
    isCalculating,
    isQueued,
    phaseLabel,
    progressPercent,
    startTracking,
    stopTracking
  }
}

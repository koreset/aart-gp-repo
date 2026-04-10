import { ref, computed } from 'vue'
import { useWebSocket } from './useWebSocket'

export interface CalculationProgress {
  quoteId: string
  totalCategories: number
  completedCategories: number
  currentCategory: string
  phase: string // "loading_data" | "rating_members" | "saving_results" | "category_done" | "completed"
  progress: number // 0-100
}

const progress = ref<CalculationProgress | null>(null)
const isCalculating = ref(false)

const phaseLabel = computed(() => {
  if (!progress.value) return ''
  switch (progress.value.phase) {
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
    default:
      return progress.value.phase
  }
})

const progressPercent = computed(() => {
  return Math.round(progress.value?.progress ?? 0)
})

function handleProgress(payload: CalculationProgress) {
  progress.value = payload
  if (payload.phase === 'completed') {
    // Keep the completed state visible briefly, then clear
    setTimeout(() => {
      if (progress.value?.phase === 'completed') {
        isCalculating.value = false
        progress.value = null
      }
    }, 1500)
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
    phase: 'loading_data',
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
    phaseLabel,
    progressPercent,
    startTracking,
    stopTracking
  }
}

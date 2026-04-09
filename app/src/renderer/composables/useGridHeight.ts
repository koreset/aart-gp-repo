import { ref, computed, onMounted, onUnmounted } from 'vue'
import type { ComputedRef } from 'vue'

export function useGridHeight(subtractPx = 320): ComputedRef<string> {
  const h = ref(window.innerHeight)
  const handler = () => {
    h.value = window.innerHeight
  }
  onMounted(() => window.addEventListener('resize', handler))
  onUnmounted(() => window.removeEventListener('resize', handler))
  return computed(() => `${Math.max(300, h.value - subtractPx)}px`)
}

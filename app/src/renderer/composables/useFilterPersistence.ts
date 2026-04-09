import { ref, watch, onMounted } from 'vue'

export function useFilterPersistence<T extends object>(
  routeName: string,
  defaults: T
) {
  const filters = ref<T>({ ...defaults })

  onMounted(() => {
    const saved = sessionStorage.getItem(routeName)
    if (saved) {
      try {
        filters.value = { ...defaults, ...JSON.parse(saved) }
      } catch {
        // ignore malformed JSON
      }
    }
  })

  watch(filters, (v) => sessionStorage.setItem(routeName, JSON.stringify(v)), {
    deep: true
  })

  const resetFilters = () => {
    sessionStorage.removeItem(routeName)
    filters.value = { ...defaults }
  }

  return { filters, resetFilters }
}

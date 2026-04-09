import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

// 1. Define a clear type for your permission map
export type UserPermissionsMap = Record<string, boolean>

export const useGroupUserPermissionsStore = defineStore(
  'group_user_permissions',
  () => {
    // --- State ---
    const permissions = ref<UserPermissionsMap>({})
    const loaded = ref(false)

    // Internal resolve function for the loading promise
    let _resolve: (() => void) | null = null
    const _promise = new Promise<void>((resolve) => {
      _resolve = resolve
    })

    // --- Getters ---
    const getPermissions = computed(() => permissions.value)
    const isLoaded = computed(() => loaded.value)

    /**
     * Check if a specific permission exists.
     * In Setup Stores, getters with arguments are just regular functions.
     */
    const hasPermission = (permission: string): boolean => {
      return permissions.value[permission] ?? false
    }

    // --- Actions ---
    const setPermissions = (newPermissions: UserPermissionsMap) => {
      permissions.value = newPermissions
      loaded.value = true
      if (_resolve) {
        _resolve()
        _resolve = null
      }
    }

    const markLoaded = () => {
      loaded.value = true
      if (_resolve) {
        _resolve()
        _resolve = null
      }
    }

    const clearPermissions = () => {
      permissions.value = {}
      loaded.value = false
    }

    /**
     * Returns a promise that resolves when permissions have been loaded.
     * Used by the router guard to wait for async permission fetch.
     */
    const waitUntilLoaded = (): Promise<void> => {
      if (loaded.value) return Promise.resolve()
      return _promise
    }

    return {
      permissions,
      loaded,
      setPermissions,
      markLoaded,
      hasPermission,
      getPermissions,
      isLoaded,
      clearPermissions,
      waitUntilLoaded
    }
  }
)

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
    // True once we know the user has a role assigned. When loaded && !hasRole,
    // the app is in bootstrap mode (fresh install, no role yet) and gates
    // should open so an initial admin can be configured.
    const hasRole = ref(false)

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
      hasRole.value = true
      loaded.value = true
      if (_resolve) {
        _resolve()
        _resolve = null
      }
    }

    const markLoaded = () => {
      hasRole.value = false
      loaded.value = true
      if (_resolve) {
        _resolve()
        _resolve = null
      }
    }

    const clearPermissions = () => {
      permissions.value = {}
      hasRole.value = false
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
      hasRole,
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

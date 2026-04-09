import { useGroupUserPermissionsStore } from '@/renderer/store/group_user'

export function usePermissionCheck() {
  const store = useGroupUserPermissionsStore()

  const hasPermission = (permission: string): boolean => {
    // If permissions haven't been loaded yet, allow access (graceful degradation)
    if (!store.loaded) return true
    // If no permissions exist (no role assigned), allow access
    if (Object.keys(store.permissions).length === 0) return true
    // system:admin bypasses all permission checks
    if (store.hasPermission('system:admin')) return true
    return store.hasPermission(permission)
  }

  return { hasPermission }
}

import { useGroupUserPermissionsStore } from '@/renderer/store/group_user'

export function usePermissionCheck() {
  const store = useGroupUserPermissionsStore()

  const hasPermission = (permission: string): boolean => {
    // Fail closed while permissions are still being fetched — we don't yet
    // know whether this user has a role, so we must not flash gated UI open.
    if (!store.loaded) return false
    // Bootstrap mode: user has no role assigned. This is expected on a fresh
    // installation where the first system administrator still needs to be
    // set up, so allow access until a role is assigned.
    if (!store.hasRole) return true
    // system:admin bypasses all permission checks.
    if (store.hasPermission('system:admin')) return true
    return store.hasPermission(permission)
  }

  return { hasPermission }
}

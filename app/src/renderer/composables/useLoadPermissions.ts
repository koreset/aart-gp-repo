import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useGroupUserPermissionsStore } from '@/renderer/store/group_user'

// Fetches the active user's role + permissions and pushes them into the
// store. Safe to call during app bootstrap — the store's loaded/hasRole
// flags drive the fail-closed-during-load / open-when-no-role behavior
// used by usePermissionCheck.
export async function loadUserPermissions(licenseId?: string | null) {
  const store = useGroupUserPermissionsStore()

  if (!licenseId) {
    console.warn('[RBAC] No licenseId provided — entering bootstrap mode')
    store.markLoaded()
    return
  }

  try {
    console.log('[RBAC] Fetching role for license:', licenseId)
    const response: any = await GroupPricingService.getRoleForUser(licenseId)
    const role = response?.data?.data ?? response?.data
    console.log('[RBAC] Role response:', role)

    if (
      role &&
      Array.isArray(role.permissions) &&
      role.permissions.length > 0
    ) {
      const permMap: Record<string, boolean> = {}
      for (const perm of role.permissions) {
        if (perm?.slug) permMap[perm.slug] = true
      }
      console.log(
        '[RBAC] Loaded',
        Object.keys(permMap).length,
        'permissions for role:',
        role.role_name
      )
      store.setPermissions(permMap)
      return
    }

    console.warn(
      '[RBAC] Role response contained no permissions — entering bootstrap mode'
    )
    store.markLoaded()
  } catch (error) {
    console.error('[RBAC] Error fetching user role:', error)
    store.markLoaded()
  }
}

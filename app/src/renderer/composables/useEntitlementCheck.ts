import { useAppStore } from '@/renderer/store/app'

/**
 * Composable function to check if a user has a specific entitlement
 * @param entitlement The entitlement to check
 * @returns Boolean indicating if the user has the entitlement
 */
export function useEntitlementCheck() {
  const appStore = useAppStore()

  /**
   * Check if the user has a specific entitlement
   * @param entitlement The entitlement to check
   * @returns true if user has the entitlement, false otherwise
   */
  const hasEntitlement = (entitlement: string): boolean => {
    const entitlements: string[] = appStore.getEntitlements() || []

    // Allow access if user has all-features permission
    if (entitlements.includes('all-features')) {
      return true
    }

    return entitlements.includes(entitlement)
  }

  return {
    hasEntitlement
  }
}

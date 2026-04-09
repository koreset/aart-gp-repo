import { RouteLocationNormalized, NavigationGuardNext } from 'vue-router'
import { useEntitlementCheck } from './useEntitlementCheck'
import { useFlashStore } from '@/renderer/store/flash'

/**
 * Navigation guard to check entitlements for routes
 * To use this guard, add requiredEntitlement to the route's meta field
 */
export const entitlementGuard = async (
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
) => {
  // Skip checking if no required entitlement is specified
  if (!to.meta.requiredEntitlement) {
    return next()
  }

  const { hasEntitlement } = useEntitlementCheck()
  const flash = useFlashStore()

  // Get the required entitlement from the route meta
  const requiredEntitlement = to.meta.requiredEntitlement as string

  // Check if the user has the required entitlement
  if (hasEntitlement(requiredEntitlement)) {
    return next()
  }

  // User doesn't have the required entitlement
  // Show an error message
  flash.show(`You don't have permission to access this page`, 'error')

  // Redirect to the dashboard or another fallback route
  return next({ name: 'dashboard' })
}

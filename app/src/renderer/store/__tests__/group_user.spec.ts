import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useGroupUserPermissionsStore } from '../group_user'

describe('useGroupUserPermissionsStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('initialises as not loaded with empty permissions', () => {
    const store = useGroupUserPermissionsStore()
    expect(store.isLoaded).toBe(false)
    expect(store.getPermissions).toEqual({})
  })

  it('setPermissions sets permissions and marks loaded', () => {
    const store = useGroupUserPermissionsStore()
    store.setPermissions({ 'view:quotes': true, 'edit:quotes': false })
    expect(store.isLoaded).toBe(true)
    expect(store.getPermissions).toEqual({
      'view:quotes': true,
      'edit:quotes': false
    })
  })

  it('hasPermission returns true for granted permission', () => {
    const store = useGroupUserPermissionsStore()
    store.setPermissions({ 'view:quotes': true })
    expect(store.hasPermission('view:quotes')).toBe(true)
  })

  it('hasPermission returns false for denied permission', () => {
    const store = useGroupUserPermissionsStore()
    store.setPermissions({ 'view:quotes': false })
    expect(store.hasPermission('view:quotes')).toBe(false)
  })

  it('hasPermission returns false for unknown permission', () => {
    const store = useGroupUserPermissionsStore()
    store.setPermissions({})
    expect(store.hasPermission('nonexistent')).toBe(false)
  })

  it('markLoaded sets loaded without changing permissions', () => {
    const store = useGroupUserPermissionsStore()
    store.markLoaded()
    expect(store.isLoaded).toBe(true)
    expect(store.getPermissions).toEqual({})
  })

  it('clearPermissions resets state', () => {
    const store = useGroupUserPermissionsStore()
    store.setPermissions({ 'view:quotes': true })
    store.clearPermissions()
    expect(store.isLoaded).toBe(false)
    expect(store.getPermissions).toEqual({})
  })

  it('waitUntilLoaded resolves immediately when already loaded', async () => {
    const store = useGroupUserPermissionsStore()
    store.setPermissions({})
    await expect(store.waitUntilLoaded()).resolves.toBeUndefined()
  })

  it('waitUntilLoaded resolves when setPermissions is called', async () => {
    const store = useGroupUserPermissionsStore()
    const promise = store.waitUntilLoaded()
    // Should not resolve yet
    let resolved = false
    promise.then(() => {
      resolved = true
    })
    // Trigger resolution
    store.setPermissions({ a: true })
    await promise
    expect(resolved).toBe(true)
  })
})

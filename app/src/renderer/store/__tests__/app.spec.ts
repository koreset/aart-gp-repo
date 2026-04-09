import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAppStore } from '../app'

describe('useAppStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('initialises with null/empty state', () => {
    const store = useAppStore()
    expect(store.user).toBeNull()
    expect(store.licenseData).toBeNull()
    expect(store.organization).toBeNull()
    expect(store.entitlements).toEqual([])
    expect(store.allProducts).toEqual([])
  })

  it('setUser / getUser', () => {
    const store = useAppStore()
    const user = { name: 'Alice', email: 'alice@test.com' }
    store.setUser(user)
    expect(store.getUser).toEqual(user)
  })

  it('setLicense / getLicenseData', () => {
    const store = useAppStore()
    const lic = { key: 'abc', valid: true }
    store.setLicense(lic)
    expect(store.getLicenseData).toEqual(lic)
  })

  it('setEntitlements / getEntitlements', () => {
    const store = useAppStore()
    const ent = ['feature_a', 'feature_b']
    store.setEntitlements(ent)
    expect(store.getEntitlements()).toEqual(ent)
  })

  it('clearEntitlement resets entitlements', () => {
    const store = useAppStore()
    store.setEntitlements(['x'])
    store.clearEntitlement()
    expect(store.entitlements).toEqual([])
  })

  it('clearUser resets user', () => {
    const store = useAppStore()
    store.setUser({ name: 'Bob' })
    store.clearUser()
    expect(store.user).toBeNull()
  })

  it('clearLicenseData resets licenseData', () => {
    const store = useAppStore()
    store.setLicense({ key: 'x' })
    store.clearLicenseData()
    expect(store.licenseData).toBeNull()
  })

  it('setProducts / getAllProducts', () => {
    const store = useAppStore()
    const prods = [{ id: 1, name: 'P1' }]
    store.setProducts(prods)
    expect(store.getAllProducts).toEqual(prods)
  })

  it('clearAll resets everything', () => {
    const store = useAppStore()
    store.setUser({ name: 'X' })
    store.setLicense({ key: 'Y' })
    store.setEntitlements(['e'])
    store.setProducts([{ id: 1 }])
    store.clearAll()
    expect(store.user).toBeNull()
    expect(store.licenseData).toBeNull()
    expect(store.organization).toBeNull()
    expect(store.entitlements).toEqual([])
    expect(store.allProducts).toEqual([])
  })
})

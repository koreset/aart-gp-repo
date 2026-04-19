import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    user: null,
    licenseData: null,
    entitlements: [],
    allProducts: []
  }),
  actions: {
    setUser(user: any) {
      this.user = user
    },
    setLicense(licenseData: any) {
      this.licenseData = licenseData
    },
    setEntitlements(entitlements: any) {
      this.entitlements = entitlements
    },
    getEntitlements() {
      return this.entitlements
    },
    clearEntitlement() {
      this.entitlements = []
    },
    clearUser() {
      this.user = null
    },
    clearLicenseData() {
      this.licenseData = null
    },
    clearAll() {
      this.user = null
      this.licenseData = null
      this.entitlements = []
      this.allProducts = []
    },
    setProducts(products: any) {
      this.allProducts = products
    }
  },
  getters: {
    getUser: (state): any => state.user,
    getLicenseData: (state): any => state.licenseData,
    getAllProducts: (state): any => state.allProducts,
    // Resolved organisation name from license metadata, set during
    // activation. Supports both the nested Keygen response shape and
    // the legacy flat shape; falls back to userName for single-user
    // licenses without an explicit organization field.
    getOrganisationName: (state): string => {
      const license: any = state.licenseData
      const meta =
        license?.data?.attributes?.metadata || license?.attributes?.metadata
      return meta?.organization || meta?.userName || ''
    }
  }
})

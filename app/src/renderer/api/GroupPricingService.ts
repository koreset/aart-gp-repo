import Api from '@/renderer/api/Api'

export default {
  generateQuote(formData) {
    return Api.post('/group-pricing/generate-quote', formData, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },
  runQuoteCalculations(quoteId, basis) {
    return Api.post(
      '/group-pricing/calculate-quote/' + quoteId + '/basis/' + basis
    )
  },
  runQuoteCalculationsWithCredibility(quoteId, basis, credibility) {
    return Api.post(
      '/group-pricing/calculate-quote/' +
        quoteId +
        '/basis/' +
        basis +
        '/credibility/' +
        credibility
    )
  },
  getDiscountAuthority(riskRateCode: string) {
    return Api.get(
      '/group-pricing/discount-authority/risk-code/' + riskRateCode
    )
  },
  applyDiscount(quoteId: number | string, discountPct: number) {
    return Api.post(
      `/group-pricing/quotes/${quoteId}/apply-discount/${discountPct}`
    )
  },
  getGroupPricingSettings() {
    return Api.get('/group-pricing/settings')
  },
  updateGroupPricingSettings(payload: {
    discount_method?: 'loading_adjustment' | 'prorata'
    fcl_method?: 'percentile' | 'outlier'
    fcl_override_tolerance?: number
    risk_alr_ceiling_pct?: number
    risk_alr_delta_pp?: number
  }) {
    return Api.put('/group-pricing/settings', payload, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },
  getTableMetaData() {
    return Api.get('/group-pricing/rate-tables')
  },
  getTableConfigurations() {
    return Api.get('/group-pricing/table-configurations')
  },
  updateTableConfiguration(
    tableType: string,
    isRequired: boolean,
    note: string = ''
  ) {
    return Api.patch(`/group-pricing/table-configurations/${tableType}`, {
      is_required: isRequired,
      note
    })
  },
  getTableConfigurationAudit(tableType: string) {
    return Api.get(`/group-pricing/table-configurations/${tableType}/audit`)
  },
  getAgeBands() {
    return Api.get('/group-pricing/age-bands')
  },
  getQuoteMemberGenderSplit(quoteId: number) {
    return Api.get(`/group-pricing/quotes/${quoteId}/member-gender-split`)
  },
  getQuoteDocx(id: number) {
    return Api.get(`/group-pricing/get-quote/${id}/document.docx`, {
      responseType: 'blob'
    })
  },
  getQuotePdf(id: number) {
    return Api.get(`/group-pricing/get-quote/${id}/document.pdf`, {
      responseType: 'blob'
    })
  },
  // ----- Per-insurer quote templates -----
  uploadInsurerQuoteTemplate(insurerId: number, file: File) {
    const formData = new FormData()
    formData.append('file', file)
    return Api.post(
      `/group-pricing/insurers/${insurerId}/quote-template`,
      formData,
      { headers: { 'Content-Type': 'multipart/form-data' } }
    )
  },
  getActiveInsurerQuoteTemplate(insurerId: number) {
    return Api.get(`/group-pricing/insurers/${insurerId}/quote-template/active`)
  },
  listInsurerQuoteTemplateVersions(insurerId: number) {
    return Api.get(
      `/group-pricing/insurers/${insurerId}/quote-template/versions`
    )
  },
  downloadInsurerQuoteTemplate(templateId: number) {
    return Api.get(
      `/group-pricing/insurers/quote-template/${templateId}/download`,
      {
        responseType: 'blob'
      }
    )
  },
  activateInsurerQuoteTemplate(insurerId: number, templateId: number) {
    return Api.post(
      `/group-pricing/insurers/${insurerId}/quote-template/${templateId}/activate`
    )
  },
  deleteInsurerQuoteTemplate(insurerId: number, templateId: number) {
    return Api.delete(
      `/group-pricing/insurers/${insurerId}/quote-template/${templateId}`
    )
  },
  deleteInactiveInsurerQuoteTemplates(insurerId: number) {
    return Api.delete(
      `/group-pricing/insurers/${insurerId}/quote-template/inactive`
    )
  },
  downloadSampleQuoteTemplate() {
    return Api.get(`/group-pricing/quote-template/sample`, {
      responseType: 'blob'
    })
  },
  // ----- Per-insurer on-risk letter templates -----
  uploadInsurerOnRiskLetterTemplate(insurerId: number, file: File) {
    const formData = new FormData()
    formData.append('file', file)
    return Api.post(
      `/group-pricing/insurers/${insurerId}/on-risk-letter-template`,
      formData,
      { headers: { 'Content-Type': 'multipart/form-data' } }
    )
  },
  getActiveInsurerOnRiskLetterTemplate(insurerId: number) {
    return Api.get(
      `/group-pricing/insurers/${insurerId}/on-risk-letter-template/active`
    )
  },
  listInsurerOnRiskLetterTemplateVersions(insurerId: number) {
    return Api.get(
      `/group-pricing/insurers/${insurerId}/on-risk-letter-template/versions`
    )
  },
  downloadInsurerOnRiskLetterTemplate(templateId: number) {
    return Api.get(
      `/group-pricing/insurers/on-risk-letter-template/${templateId}/download`,
      { responseType: 'blob' }
    )
  },
  activateInsurerOnRiskLetterTemplate(insurerId: number, templateId: number) {
    return Api.post(
      `/group-pricing/insurers/${insurerId}/on-risk-letter-template/${templateId}/activate`
    )
  },
  deleteInsurerOnRiskLetterTemplate(insurerId: number, templateId: number) {
    return Api.delete(
      `/group-pricing/insurers/${insurerId}/on-risk-letter-template/${templateId}`
    )
  },
  deleteInactiveInsurerOnRiskLetterTemplates(insurerId: number) {
    return Api.delete(
      `/group-pricing/insurers/${insurerId}/on-risk-letter-template/inactive`
    )
  },
  downloadSampleOnRiskLetterTemplate() {
    return Api.get(`/group-pricing/on-risk-letter-template/sample`, {
      responseType: 'blob'
    })
  },
  getOnRiskLetterDocx(quoteId: number) {
    return Api.get(
      `/group-pricing/quotes/${quoteId}/on-risk-letter/document.docx`,
      { responseType: 'blob' }
    )
  },
  uploadTables(formdata) {
    return Api.post('group-pricing/rate-tables', formdata, {
      headers: {
        'Content-Type': 'multipart/form-data',
        Accept: 'multipart/form-data'
      }
    })
  },

  uploadQuoteTable(formdata) {
    return Api.post('group-pricing/quote-tables', formdata, {
      headers: {
        'Content-Type': 'multipart/form-data',
        Accept: 'multipart/form-data'
      }
    })
  },

  deleteQuoteTableData(quoteId, tableType) {
    return Api.delete('group-pricing/quote-tables/' + tableType + '/' + quoteId)
  },

  deleteTable(tableType, riskCode) {
    return Api.delete(
      'group-pricing/rate-tables/' + tableType + '/risk-code/' + riskCode
    )
  },
  getDataForTable(tableType) {
    return Api.get('/group-pricing/rate-tables/' + tableType.toLowerCase())
  },
  createBroker(broker) {
    return Api.post('/group-pricing/brokers', broker)
  },
  getBrokers() {
    return Api.get('/group-pricing/brokers')
  },
  getIndustries() {
    return Api.get('/group-pricing/industries')
  },
  getBroker(id) {
    return Api.get('/group-pricing/brokers/' + id)
  },
  updateBroker(id, brokerDetails) {
    return Api.put(`/group-pricing/brokers/${id}`, brokerDetails, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },
  deleteBroker(id) {
    return Api.delete('/group-pricing/brokers/' + id, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },

  // Binder Fee CRUD operations
  createBinderFee(payload) {
    return Api.post('/group-pricing/binder-fees', payload)
  },
  getBinderFees() {
    return Api.get('/group-pricing/binder-fees')
  },
  getBinderFee(id) {
    return Api.get('/group-pricing/binder-fees/' + id)
  },
  updateBinderFee(id, payload) {
    return Api.put(`/group-pricing/binder-fees/${id}`, payload, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },
  deleteBinderFee(id) {
    return Api.delete('/group-pricing/binder-fees/' + id, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },

  // Commission structure (sliding-scale per channel) CRUD
  createCommissionBand(payload) {
    return Api.post('/group-pricing/commission-structures', payload)
  },
  getCommissionBands(
    channel?: string,
    holderName?: string | null,
    allHolders?: boolean
  ) {
    const params = new URLSearchParams()
    if (channel) params.set('channel', channel)
    // `holderName === undefined` means "don't filter by holder at all".
    // `holderName === ''` means "show default (empty-holder) rows only".
    // A non-empty string narrows to that holder.
    if (holderName !== undefined && holderName !== null) {
      params.set('holder_name', holderName)
    }
    if (allHolders) params.set('all', '1')
    const qs = params.toString()
    return Api.get(
      '/group-pricing/commission-structures' + (qs ? '?' + qs : '')
    )
  },
  getCommissionBand(id) {
    return Api.get('/group-pricing/commission-structures/' + id)
  },
  updateCommissionBand(id, payload) {
    return Api.put(`/group-pricing/commission-structures/${id}`, payload, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },
  deleteCommissionBand(id) {
    return Api.delete('/group-pricing/commission-structures/' + id, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },

  // Reinsurer CRUD operations
  createReinsurer(reinsurer) {
    return Api.post('/group-pricing/reinsurers', reinsurer)
  },
  getReinsurers() {
    return Api.get('/group-pricing/reinsurers')
  },
  getReinsurer(id) {
    return Api.get('/group-pricing/reinsurers/' + id)
  },
  updateReinsurer(id, reinsurerDetails) {
    return Api.put(`/group-pricing/reinsurers/${id}`, reinsurerDetails, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },
  deactivateReinsurer(id, reason) {
    return Api.post(
      `/group-pricing/reinsurers/${id}/deactivate`,
      { reason },
      {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json'
        }
      }
    )
  },
  deleteReinsurer(id) {
    return Api.delete('/group-pricing/reinsurers/' + id, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },

  // Scheme Category Masters API
  createSchemeCategoryMaster(categoryData) {
    return Api.post('/group-pricing/scheme-category-masters', categoryData, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },
  getSchemeCategoryMasters() {
    return Api.get('/group-pricing/scheme-category-masters')
  },
  getSchemeCategoryMaster(id) {
    return Api.get('/group-pricing/scheme-category-masters/' + id)
  },
  updateSchemeCategoryMaster(id, categoryData) {
    return Api.put(
      `/group-pricing/scheme-category-masters/${id}`,
      categoryData,
      {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json'
        }
      }
    )
  },
  deleteSchemeCategoryMaster(id) {
    return Api.delete('/group-pricing/scheme-category-masters/' + id, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },

  createScheme(scheme) {
    return Api.post('/group-pricing/schemes', scheme)
  },
  getSchemesInforce() {
    return Api.get('/group-pricing/schemes/in-force')
  },
  getSchemesInforcev2() {
    return Api.get('/group-pricing/schemes/in-force-v2')
  },
  getScheme(id) {
    return Api.get('/group-pricing/schemes/' + id)
  },
  deleteScheme(id) {
    return Api.delete('/group-pricing/schemes/' + id)
  },
  getSchemeCategories(id) {
    return Api.get('/group-pricing/schemes/' + id + '/categories')
  },
  getParameterBases() {
    return Api.get('/group-pricing/parameter-bases')
  },
  getQuotes(filter) {
    return Api.get('/group-pricing/get-quotes/filter/' + filter)
  },
  getAllQuotes() {
    return Api.get('/group-pricing/get-quotes')
  },
  changeQuoteStatus(quote) {
    return Api.post(
      '/group-pricing/quotes/' + quote.id + '/update-status',
      quote
    )
  },
  getQuote(quoteId) {
    return Api.get('/group-pricing/get-quote/' + quoteId)
  },
  getCustomTirStatus(quoteId) {
    return Api.get('/group-pricing/get-quote/' + quoteId + '/custom-tir-status')
  },
  getQuoteBySchemeName(schemeName) {
    return Api.get('/group-pricing/get-quote-by-scheme-name/' + schemeName)
  },
  deleteQuote(quoteId) {
    return Api.delete('/group-pricing/quotes/' + quoteId)
  },
  getQuoteTable(
    quoteId: any,
    tableType: any,
    params: { offset?: number; limit?: number | null } = {}
  ) {
    const { offset = 0, limit = null } = params
    let url = '/group-pricing/get-quote/' + quoteId + '/table-type/' + tableType
    if (limit !== null) {
      url += `?offset=${offset}&limit=${limit}`
    }
    return Api.get(url)
  },
  exportQuoteTableCsv(quoteId, tableType) {
    return Api.get(
      '/group-pricing/export-csv/' + quoteId + '/table-type/' + tableType,
      {
        responseType: 'blob'
      }
    )
  },
  getResultSummary(quoteId) {
    return Api.get('/group-pricing/get-quote/' + quoteId + '/result-summary')
  },
  getCategoryEducatorBenefits(quoteId) {
    return Api.get(
      '/group-pricing/get-quote/' + quoteId + '/categories/educator-benefits'
    )
  },
  createInsurer(formData) {
    return Api.post('/group-pricing/insurers', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        Accept: 'multipart/form-data'
      }
    })
  },
  getInsurer() {
    return Api.get('/group-pricing/insurers')
  },
  acceptQuote(quoteId, commencementDate, term) {
    return Api.post('/group-pricing/quotes/' + quoteId + '/accept-quote', {
      commencement_date: commencementDate,
      term
    })
  },
  approveQuote(quoteId) {
    return Api.post('/group-pricing/quotes/' + quoteId + '/approve-quote')
  },
  createOnRiskLetter(quoteId: number | string) {
    return Api.post(`/group-pricing/quotes/${quoteId}/on-risk-letter`)
  },
  getOnRiskLetterData(quoteId: number | string) {
    return Api.get(`/group-pricing/quotes/${quoteId}/on-risk-letter`)
  },
  getDashboardData(year, dataSource = 'inforce', benefit = 'All') {
    return Api.get(
      `/group-pricing/dashboard/year/${year}?data_source=${dataSource}&benefit=${encodeURIComponent(benefit)}`
    )
  },
  getSchemePerformance() {
    return Api.get('/group-pricing/dashboard/scheme-performance')
  },
  getRiskProfile() {
    return Api.get('/group-pricing/dashboard/risk-profile')
  },
  getLossRatioTrend() {
    return Api.get('/group-pricing/dashboard/loss-ratio-trend')
  },
  getExposureData(year, benefit, dataSource = 'all') {
    return Api.get(
      `/group-pricing/dashboard/exposures/year/${year}/benefit/${benefit}?data_source=${dataSource}`
    )
  },
  rebuildExposureData(year) {
    return Api.post(`/group-pricing/dashboard/exposures/rebuild/year/${year}`)
  },
  getExposureTrend(benefit = 'All', dataSource = 'inforce') {
    return Api.get(
      `/group-pricing/dashboard/exposures/trend?benefit=${benefit}&data_source=${dataSource}`
    )
  },
  getFinancialYearInfo(year) {
    return Api.get(`/group-pricing/metadata/financial-year-info?year=${year}`)
  },
  checkDuplicateSchemeName(scheme) {
    return Api.get('/group-pricing/schemes/check-name/' + scheme)
  },
  getInforceDataTable(schemeId, tableType) {
    return Api.get(
      '/group-pricing/inforce-data/' + schemeId + '/table-type/' + tableType
    )
  },
  addMember(member) {
    return Api.post(
      '/group-pricing/schemes/' + member.scheme_id + '/members',
      member
    )
  },
  getMembersInForce(schemeId) {
    return Api.get('/group-pricing/schemes/' + schemeId + '/members')
  },
  getMembersPaginated(params) {
    // Build query parameters for server-side pagination and filtering
    const queryParams = new URLSearchParams()

    if (params.page) queryParams.append('page', params.page.toString())
    if (params.pageSize)
      queryParams.append('pageSize', params.pageSize.toString())
    if (params.search) queryParams.append('search', params.search)
    if (params.schemeId)
      queryParams.append('schemeId', params.schemeId.toString())
    if (params.status) queryParams.append('status', params.status)

    const url = `/group-pricing/members/paginated?${queryParams.toString()}`

    return Api.get(url, {
      signal: params.signal // For request cancellation
    })
  },
  submitClaim(formData) {
    return Api.post('/group-pricing/claims', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        Accept: 'multipart/form-data'
      }
    })

    // If there are supporting documents, use FormData for file upload
    // if (claim.supporting_documents && claim.supporting_documents.length > 0) {
    //   const formData = new FormData()

    //   // Add all claim fields except supporting_documents
    //   const claimData = { ...claim }
    //   delete claimData.supporting_documents

    //   // Add claim data as JSON blob
    //   formData.append('claim_data', JSON.stringify(claimData))

    //   // Add each file
    //   claim.supporting_documents.forEach((file, index) => {
    //     formData.append(`supporting_documents[${index}]`, file)
    //   })

    //   return Api.post('/group-pricing/claims', formData, {
    //     headers: {
    //       'Content-Type': 'multipart/form-data',
    //       Accept: 'multipart/form-data'
    //     }
    //   })
    // } else {
    //   // No files, use regular JSON payload
    //   const payload = JSON.stringify(claim)
    //   return Api.post('/group-pricing/claims', payload, {
    //     headers: {
    //       'Content-Type': 'application/json',
    //       Accept: 'application/json'
    //     }
    //   })
    // }
  },
  getClaims() {
    return Api.get('/group-pricing/claims')
  },
  getClaimsAnalytics(
    filters: {
      scheme_id?: any
      period?: any
      benefit_type?: any
      from?: any
      to?: any
    } = {}
  ) {
    // Build query parameters from filters
    const queryParams = new URLSearchParams()
    if (filters.scheme_id)
      queryParams.append('scheme_id', filters.scheme_id.toString())
    if (filters.period) queryParams.append('period', filters.period)
    if (filters.benefit_type)
      queryParams.append('benefit_type', filters.benefit_type)
    if (filters.from) queryParams.append('from', filters.from)
    if (filters.to) queryParams.append('to', filters.to)
    const queryString = queryParams.toString()
    const url = queryString
      ? `/group-pricing/claims/dashboard?${queryString}`
      : '/group-pricing/claims/dashboard'

    return Api.get(url)
  },
  getSchemes() {
    return Api.get('/group-pricing/schemes')
  },
  getMemberRating(schemeId, quoteId, memberId) {
    return Api.get(
      '/group-pricing/claims/scheme/' +
        schemeId +
        '/quote/' +
        quoteId +
        '/member/' +
        memberId +
        '/rating'
    )
  },
  getTableYears(tableType) {
    return Api.get('/group-pricing/rate-tables/' + tableType + '/get-years')
  },
  getRiskRateCodes(tableType) {
    return Api.get(
      '/group-pricing/rate-tables/' + tableType + '/get-risk-codes'
    )
  },
  getBenefitMaps() {
    return Api.get('/group-pricing/benefit-maps')
  },
  getBenefitMapsByScheme(schemeId) {
    return Api.get('/group-pricing/benefit-maps/scheme/' + schemeId)
  },
  getBenefitMapsBySchemeAndCategory(schemeId, categoryId) {
    return Api.get(
      '/group-pricing/benefit-maps/scheme/' +
        schemeId +
        '/category/' +
        categoryId
    )
  },
  saveBenefitMap(benefitMaps) {
    return Api.post('/group-pricing/benefit-maps', benefitMaps)
  },
  getBenefitDefinitions() {
    return Api.get('/group-pricing/benefit-definitions')
  },
  getUserRoles() {
    return Api.get('/group-pricing/user-management/roles')
  },
  getPermissions() {
    return Api.get('/group-pricing/user-management/permissions')
  },
  createUserRole(role) {
    return Api.post('/group-pricing/user-management/roles', role)
  },
  getRolePermissions(roleId) {
    return Api.get(
      '/group-pricing/user-management/roles/' + roleId + '/permissions'
    )
  },
  updateUserRole(userRole) {
    return Api.put('/group-pricing/user-management/users/assign_role', userRole)
  },
  getRoleForUser(licenseId) {
    return Api.get(
      '/group-pricing/user-management/users/license/' + licenseId + '/role'
    )
  },
  deleteUserRole(roleId) {
    return Api.delete('/group-pricing/user-management/roles/' + roleId)
  },
  removeUserRole(userRole) {
    return Api.post(
      '/group-pricing/user-management/users/remove_role',
      userRole
    )
  },
  getBenefitEscalations(_riskRateCode?: any) {
    return Api.get('/group-pricing/quotes/benefit-escalations')
  },
  getTtdDisabilityDefinitions(riskRateCode) {
    return Api.get(
      '/group-pricing/quotes/ttd-disability-definitions/risk-rate-code/' +
        riskRateCode
    )
  },
  getPtdDisabilityDefinitions(riskRateCode) {
    return Api.get(
      '/group-pricing/quotes/ptd-disability-definitions/risk-rate-code/' +
        riskRateCode
    )
  },
  getPhiDisabilityDefinitions(riskRateCode) {
    return Api.get(
      '/group-pricing/quotes/phi-disability-definitions/risk-rate-code/' +
        riskRateCode
    )
  },
  getEducatorBenefitTypes(riskRateCode) {
    return Api.get(
      '/group-pricing/rate-tables/educator-benefits/risk-rate-code/' +
        riskRateCode
    )
  },

  getWaitingPeriods(tableType, riskRateCode) {
    return Api.get(
      '/group-pricing/rate-tables/' +
        tableType +
        '/risk-rate-code/' +
        riskRateCode +
        '/waiting-periods'
    )
  },
  getGlaBenefitTypes(riskRateCode) {
    return Api.get(
      '/group-pricing/rate-tables/gla/risk-rate-code/' +
        riskRateCode +
        '/benefit-types'
    )
  },
  getDeferredPeriods(tableType, riskRateCode) {
    return Api.get(
      '/group-pricing/rate-tables/' +
        tableType +
        '/risk-rate-code/' +
        riskRateCode +
        '/deferred-periods'
    )
  },
  getRiskTypes(_tableType?: any) {
    return Api.get('/group-pricing/rate-tables/risk-types')
  },
  getNormalRetirementAges() {
    return Api.get(
      '/group-pricing/rate-tables/phi-rates/normal-retirement-ages'
    )
  },
  getHistoricalCredibilityData() {
    return Api.get('/group-pricing/historical-credibility-data')
  },
  updateSchemeCoverEndDate(schemeId, coverEndDate) {
    const payload = JSON.stringify({
      cover_end_date: coverEndDate,
      scheme_id: schemeId
    })
    return Api.put(
      '/group-pricing/schemes/' + schemeId + '/cover-end-date',
      payload
    )
  },
  updateSchemeStatus(schemeId, statusData) {
    const payload = JSON.stringify({
      ...statusData,
      scheme_id: schemeId
    })
    return Api.put('/group-pricing/schemes/' + schemeId + '/status', payload, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },
  getSchemeStatusHistory(schemeId) {
    return Api.get('/group-pricing/schemes/' + schemeId + '/status-history')
  },
  getSchemeQuotes(schemeId) {
    return Api.get('/group-pricing/schemes/' + schemeId + '/quotes')
  },
  // This function searches for members in a scheme based on a search term
  searchMembers(schemeId, quoteId, searchTerm) {
    return Api.get(
      '/group-pricing/schemes/' +
        schemeId +
        '/quotes/' +
        quoteId +
        '/members/search?query=' +
        searchTerm
    )
  },
  removeMemberFromScheme(schemeId, member) {
    const payload = JSON.stringify(member)
    return Api.put('/group-pricing/members/' + member.member_id_number, payload)
  },
  deleteQuoteTable(quoteId, tableType) {
    return Api.delete(
      '/group-pricing/quote-tables/' + tableType + '/' + quoteId
    )
  },
  sendIndicativeMemberData(indicativeData) {
    const payload = JSON.stringify(indicativeData)
    return Api.post('/group-pricing/quotes/indicative-member-data', payload)
  },
  deleteIndicativeMemberData(quoteId) {
    return Api.delete(
      '/group-pricing/quotes/' + quoteId + '/indicative-member-data'
    )
  },
  getExperienceRateOverrides(quoteId) {
    return Api.get('/group-pricing/experience-rate-overrides/' + quoteId)
  },
  saveExperienceRateOverrides(rows) {
    return Api.post(
      '/group-pricing/experience-rate-overrides',
      JSON.stringify(rows),
      {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json'
        }
      }
    )
  },
  deleteExperienceRateOverrides(quoteId) {
    return Api.delete('/group-pricing/experience-rate-overrides/' + quoteId)
  },
  updateExperienceOverrideCredibility(quoteId, credibility) {
    return Api.put(
      '/group-pricing/experience-rate-overrides/' + quoteId + '/credibility',
      JSON.stringify({ credibility }),
      {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json'
        }
      }
    )
  },
  updateIndicativeDataFlag(quoteId, indicativeDataEnabled) {
    const payload = JSON.stringify({
      indicative_data_enabled: indicativeDataEnabled
    })
    return Api.patch(
      '/group-pricing/quotes/' + quoteId + '/indicative-member-data',
      payload,
      {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json'
        }
      }
    )
  },
  getMemberInfo(memberId) {
    return Api.get('/group-pricing/members/' + memberId)
  },
  updateMember(memberId, memberData) {
    return Api.put('/group-pricing/members/' + memberId, memberData)
  },
  getMemberBeneficiaries(memberId) {
    return Api.get('/group-pricing/members/' + memberId + '/beneficiaries')
  },
  getMemberBenefitSummary(memberId) {
    return Api.get('/group-pricing/members/' + memberId + '/benefit-summary')
  },
  addBeneficiaryToMember(memberId, beneficiaryData) {
    return Api.post(
      '/group-pricing/members/' + memberId + '/beneficiaries',
      beneficiaryData
    )
  },
  updateBeneficiary(memberId, beneficiaryData) {
    return Api.put(
      '/group-pricing/members/' +
        memberId +
        '/beneficiaries/' +
        beneficiaryData.id,
      beneficiaryData
    )
  },
  deleteBeneficiary(memberId, beneficiaryId) {
    return Api.delete(
      '/group-pricing/members/' + memberId + '/beneficiaries/' + beneficiaryId
    )
  },
  getMemberByIdNumber(idNumber) {
    return Api.get('/group-pricing/members/id-number/' + idNumber)
  },
  getUpdatedClaimAmount(claimRequest) {
    return Api.post('/group-pricing/claims/calculate-amount', claimRequest, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },

  // Claim Assessment methods
  createClaimAssessment(claimId, assessmentData) {
    return Api.post(
      '/group-pricing/claims/' + claimId + '/assessments',
      assessmentData,
      {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json'
        }
      }
    )
  },

  getClaimAssessments(claimId) {
    return Api.get(`/group-pricing/claims/${claimId}/assessments`)
  },

  updateClaim(claimId, claimData) {
    return Api.put(`/group-pricing/claims/${claimId}`, claimData, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },

  createClaimDeclineRecord(claimId, declineData) {
    return Api.post(`/group-pricing/claims/${claimId}/decline`, declineData, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },

  bulkUploadClaims(claimsData) {
    return Api.post('/group-pricing/claims', claimsData, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },
  getClaimsAttachement(claimId, attachmentId) {
    return Api.get(`/group-pricing/claims/attachments/${attachmentId}`, {
      responseType: 'blob'
    })
  },

  // ──────────────────────────────────────────────
  // Claim Payment Schedule methods
  // ──────────────────────────────────────────────

  createPaymentSchedule(payload) {
    return Api.post('/group-pricing/claims/payment-schedules', payload, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },

  getPaymentSchedules() {
    return Api.get('/group-pricing/claims/payment-schedules')
  },

  getPaymentSchedule(scheduleId) {
    return Api.get(`/group-pricing/claims/payment-schedules/${scheduleId}`)
  },

  exportPaymentScheduleCSV(scheduleId) {
    return Api.get(
      `/group-pricing/claims/payment-schedules/${scheduleId}/export`,
      {
        responseType: 'blob'
      }
    )
  },

  uploadPaymentProof(scheduleId, formData) {
    return Api.post(
      `/group-pricing/claims/payment-schedules/${scheduleId}/proof`,
      formData,
      {
        headers: { 'Content-Type': 'multipart/form-data' }
      }
    )
  },

  getPaymentProofs(scheduleId) {
    return Api.get(
      `/group-pricing/claims/payment-schedules/${scheduleId}/proof`
    )
  },

  downloadPaymentProof(proofId) {
    return Api.get(
      `/group-pricing/claims/payment-schedules/proof/${proofId}/download`,
      {
        responseType: 'blob'
      }
    )
  },

  // ──────────────────────────────────────────────
  // ACB Bank Profile methods
  // ──────────────────────────────────────────────

  createBankProfile(payload) {
    return Api.post('/group-pricing/claims/bank-profiles', payload)
  },

  getBankProfiles() {
    return Api.get('/group-pricing/claims/bank-profiles')
  },

  getBankProfile(profileId) {
    return Api.get(`/group-pricing/claims/bank-profiles/${profileId}`)
  },

  updateBankProfile(profileId, payload) {
    return Api.patch(
      `/group-pricing/claims/bank-profiles/${profileId}`,
      payload
    )
  },

  deleteBankProfile(profileId) {
    return Api.delete(`/group-pricing/claims/bank-profiles/${profileId}`)
  },

  // ──────────────────────────────────────────────
  // ACB File Generation & Reconciliation methods
  // ──────────────────────────────────────────────

  generateACBFile(scheduleId, payload) {
    return Api.post(
      `/group-pricing/claims/payment-schedules/${scheduleId}/acb`,
      payload
    )
  },

  getACBFileRecords(scheduleId) {
    return Api.get(
      `/group-pricing/claims/payment-schedules/${scheduleId}/acb-files`
    )
  },

  downloadACBFile(acbFileId) {
    return Api.get(`/group-pricing/claims/acb-files/${acbFileId}/download`, {
      responseType: 'blob'
    })
  },

  processBankResponse(acbFileId, formData) {
    return Api.post(
      `/group-pricing/claims/acb-files/${acbFileId}/reconcile`,
      formData,
      { headers: { 'Content-Type': 'multipart/form-data' } }
    )
  },

  getReconciliationResults(acbFileId) {
    return Api.get(
      `/group-pricing/claims/acb-files/${acbFileId}/reconciliation`
    )
  },

  getReconciliationSummary(scheduleId) {
    return Api.get(
      `/group-pricing/claims/payment-schedules/${scheduleId}/reconciliation-summary`
    )
  },

  retryFailedPayments(acbFileId, payload) {
    return Api.post(
      `/group-pricing/claims/acb-files/${acbFileId}/retry`,
      payload
    )
  },

  // Bordereaux generation methods
  generateBordereaux(formData) {
    return Api.post('/group-pricing/bordereaux/generate', formData, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },

  getBordereauxJobStatus(jobId) {
    return Api.get(`/group-pricing/bordereaux/jobs/${jobId}/status`)
  },

  downloadBordereauxFile(downloadUrl) {
    return Api.get(downloadUrl, {
      responseType: 'blob'
    })
  },

  // Bordereaux Template management methods
  getBordereauxTemplates() {
    return Api.get('/group-pricing/bordereaux/templates')
  },

  getBordereauxTemplate(templateId) {
    return Api.get(`/group-pricing/bordereaux/templates/${templateId}`)
  },

  createBordereauxTemplate(templateData) {
    return Api.post('/group-pricing/bordereaux/templates', templateData, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },

  updateBordereauxTemplate(templateId, templateData) {
    return Api.put(
      `/group-pricing/bordereaux/templates/${templateId}`,
      templateData,
      {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json'
        }
      }
    )
  },

  deleteBordereauxTemplate(templateId) {
    return Api.delete(`/group-pricing/bordereaux/templates/${templateId}`)
  },
  testBordereauxTemplate(
    templateId,
    data: { sample_size?: number; scheme_id?: number } = {}
  ) {
    return Api.post(
      `/group-pricing/bordereaux/templates/${templateId}/test`,
      data
    )
  },

  activateBordereauxTemplate(templateId) {
    return Api.patch(
      `/group-pricing/bordereaux/templates/${templateId}/activate`
    )
  },

  deactivateBordereauxTemplate(templateId) {
    return Api.patch(
      `/group-pricing/bordereaux/templates/${templateId}/deactivate`
    )
  },

  // Bordereaux reconciliation methods
  importSchemeConfirmations(formData) {
    return Api.post(
      '/group-pricing/bordereaux/confirmations/import',
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
          Accept: 'application/json'
        }
      }
    )
  },

  getBordereauxFields(bordereauType) {
    return Api.get(`/group-pricing/bordereaux/fields/${bordereauType}`)
  },

  // Bordereaux configuration management methods
  getBordereauxConfigurations() {
    return Api.get('/group-pricing/bordereaux/configurations')
  },

  getBordereauxConfiguration(configId) {
    return Api.get(`/group-pricing/bordereaux/configurations/${configId}`)
  },

  saveBordereauxConfiguration(configurationData) {
    return Api.post(
      '/group-pricing/bordereaux/configurations',
      configurationData,
      {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json'
        }
      }
    )
  },

  updateBordereauxConfiguration(configId, configurationData) {
    return Api.put(
      `/group-pricing/bordereaux/configurations/${configId}`,
      configurationData,
      {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json'
        }
      }
    )
  },

  deleteBordereauxConfiguration(configId) {
    return Api.delete(`/group-pricing/bordereaux/configurations/${configId}`)
  },

  updateConfigurationUsage(configId) {
    return Api.patch(
      `/group-pricing/bordereaux/configurations/${configId}/usage`
    )
  },

  // Bordereaux dashboard and activity methods
  getBordereauxActivity(params = {}) {
    const queryString = new URLSearchParams(params).toString()
    const url = queryString
      ? `/group-pricing/bordereaux/generated?${queryString}`
      : '/group-pricing/bordereaux/generated'
    return Api.get(url)
  },

  getBordereauxDashboardStats() {
    return Api.get('/group-pricing/bordereaux/dashboard/stats')
  },

  getBordereauxById(generatedId) {
    return Api.get(`/group-pricing/bordereaux/generated/${generatedId}`)
  },

  submitBordereauxBatch(submissionData) {
    return Api.post('/group-pricing/bordereaux/batch-submit', submissionData, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },

  getClaimCommunications(claimId) {
    return Api.get(`/group-pricing/claims/${claimId}/communications`)
  },
  createClaimCommunication(claimId, communicationData) {
    return Api.post(
      `/group-pricing/claims/${claimId}/communications`,
      communicationData,
      {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json'
        }
      }
    )
  },

  uploadClaimDocument(claimId, formData) {
    return Api.post(`/group-pricing/claims/${claimId}/attachments`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        Accept: 'application/json'
      }
    })
  },
  getMemberHistory(memberId) {
    return Api.get(`/group-pricing/members/${memberId}/history?format=activity`)
  },
  getReconciliationItems() {
    // return window.mainApi.invoke('get-reconciliation-items')
    return Api.get('/group-pricing/bordereaux/confirmations')
  },
  runAutoReconciliation() {
    // return window.mainApi.invoke('run-auto-reconciliation')
    return Api.post('/group-pricing/bordereaux/confirmations/reconcile-pending')
  },
  // In GroupPricingService.js
  getDiscrepancyDetails(reconciliationItemId) {
    // return window.mainApi.invoke('get-discrepancy-details', { id: reconciliationItemId })
    return Api.get(
      `/group-pricing/bordereaux/confirmations/${reconciliationItemId}/unmatched`
    )
  },
  deleteReconciliationItem(reconciliationItemId) {
    return Api.delete(
      `/group-pricing/bordereaux/confirmations/${reconciliationItemId}`
    )
  },

  // Reconciliation resolution methods
  resolveDiscrepancy(resultId, data) {
    return Api.post(
      `/group-pricing/bordereaux/reconciliation/results/${resultId}/resolve`,
      data
    )
  },
  escalateDiscrepancy(resultId, data) {
    return Api.post(
      `/group-pricing/bordereaux/reconciliation/results/${resultId}/escalate`,
      data
    )
  },
  listBordereauxEscalations(
    params: {
      assigned_to?: string
      priority?: string
      overdue_only?: boolean
    } = {}
  ) {
    return Api.get('/group-pricing/bordereaux/reconciliation/escalations', {
      params
    })
  },
  getBordereauxAnalytics(
    params: {
      period?: string
      from?: string
      to?: string
      scheme_id?: number
    } = {}
  ) {
    return Api.get('/group-pricing/bordereaux/analytics', { params })
  },
  addDiscrepancyComment(resultId, data) {
    return Api.post(
      `/group-pricing/bordereaux/reconciliation/results/${resultId}/comment`,
      data
    )
  },
  confirmReconciliation(confirmationId) {
    return Api.post(
      `/group-pricing/bordereaux/confirmations/${confirmationId}/confirm`
    )
  },
  reprocessReconciliation(confirmationId) {
    return Api.post(
      `/group-pricing/bordereaux/confirmations/${confirmationId}/reprocess`
    )
  },
  addReconciliationNote(confirmationId, data) {
    return Api.post(
      `/group-pricing/bordereaux/confirmations/${confirmationId}/note`,
      data
    )
  },
  getReconciliationNotes(confirmationId) {
    return Api.get(
      `/group-pricing/bordereaux/confirmations/${confirmationId}/notes`
    )
  },

  downloadBordereaux(filename) {
    return Api.get(`/group-pricing/bordereaux/download/${filename}`, {
      responseType: 'blob'
    })
  },
  downloadBordereauxComplianceReport(
    params: { from?: string; to?: string } = {}
  ) {
    return Api.get('/group-pricing/bordereaux/compliance-report', {
      params,
      responseType: 'blob'
    })
  },
  getBordereauxData(bordereauxId) {
    return Api.get(`/group-pricing/bordereaux/generated/${bordereauxId}/data`)
  },

  // Inbound Employer Submissions
  createEmployerSubmission(data) {
    return Api.post('/group-pricing/bordereaux/submissions', data)
  },
  getEmployerSubmissions(params) {
    return Api.get('/group-pricing/bordereaux/submissions', { params })
  },
  getEmployerSubmission(id) {
    return Api.get(`/group-pricing/bordereaux/submissions/${id}`)
  },
  uploadEmployerSubmission(id, formData) {
    return Api.post(
      `/group-pricing/bordereaux/submissions/${id}/upload`,
      formData,
      {
        headers: { 'Content-Type': 'multipart/form-data' }
      }
    )
  },
  reviewEmployerSubmission(id, data) {
    return Api.post(`/group-pricing/bordereaux/submissions/${id}/review`, data)
  },
  raiseSubmissionQuery(id, data) {
    return Api.post(`/group-pricing/bordereaux/submissions/${id}/query`, data)
  },
  acceptEmployerSubmission(id, data) {
    return Api.post(`/group-pricing/bordereaux/submissions/${id}/accept`, data)
  },
  rejectEmployerSubmission(id, data) {
    return Api.post(`/group-pricing/bordereaux/submissions/${id}/reject`, data)
  },
  getSubmissionRecords(id) {
    return Api.get(`/group-pricing/bordereaux/submissions/${id}/records`)
  },
  generateScheduleFromSubmission(id) {
    return Api.post(
      `/group-pricing/bordereaux/submissions/${id}/generate-schedule`
    )
  },
  computeSubmissionDelta(id) {
    return Api.post(`/group-pricing/bordereaux/submissions/${id}/compute-delta`)
  },
  getSubmissionDeltaRecords(id) {
    return Api.get(`/group-pricing/bordereaux/submissions/${id}/delta`)
  },
  syncSubmissionToMemberRegister(id) {
    return Api.post(`/group-pricing/bordereaux/submissions/${id}/sync-members`)
  },
  computeRegisterDiff(id) {
    return Api.get(`/group-pricing/bordereaux/submissions/${id}/register-diff`)
  },
  snapshotRegisterDiff(id) {
    return Api.post(`/group-pricing/bordereaux/submissions/${id}/snapshot-diff`)
  },
  applySubmissionExits(id) {
    return Api.post(`/group-pricing/bordereaux/submissions/${id}/apply-exits`)
  },
  applySubmissionAmendments(id) {
    return Api.post(
      `/group-pricing/bordereaux/submissions/${id}/apply-amendments`
    )
  },
  getNewJoinerDetails(id) {
    return Api.get(
      `/group-pricing/bordereaux/submissions/${id}/new-joiner-details`
    )
  },
  uploadNewJoinerDetails(id, formData) {
    return Api.post(
      `/group-pricing/bordereaux/submissions/${id}/upload-new-joiner-details`,
      formData,
      {
        headers: { 'Content-Type': 'multipart/form-data' }
      }
    )
  },
  syncNewJoiners(id) {
    return Api.post(
      `/group-pricing/bordereaux/submissions/${id}/sync-new-joiners`
    )
  },

  // Submission Deadline Calendar
  getBordereauxDeadlines(params) {
    return Api.get('/group-pricing/bordereaux/deadlines', { params })
  },
  createBordereauxDeadline(data) {
    return Api.post('/group-pricing/bordereaux/deadlines', data)
  },
  generateBordereauxDeadlines(data) {
    return Api.post('/group-pricing/bordereaux/deadlines/generate', data)
  },
  updateDeadlineStatus(id, data) {
    return Api.patch(`/group-pricing/bordereaux/deadlines/${id}/status`, data)
  },
  getDeadlineStats() {
    return Api.get('/group-pricing/bordereaux/deadlines/stats')
  },

  // Outbound Bordereaux Approval Workflow
  reviewGeneratedBordereaux(generatedId, data) {
    return Api.post(
      `/group-pricing/bordereaux/generated/${generatedId}/review`,
      data
    )
  },
  approveGeneratedBordereaux(generatedId, data) {
    return Api.post(
      `/group-pricing/bordereaux/generated/${generatedId}/approve`,
      data
    )
  },
  returnBordereauxToDraft(generatedId, data) {
    return Api.post(
      `/group-pricing/bordereaux/generated/${generatedId}/return-to-draft`,
      data
    )
  },
  regenerateGeneratedBordereaux(generatedId) {
    return Api.post(
      `/group-pricing/bordereaux/generated/${generatedId}/regenerate`
    )
  },

  // Reinsurer Acceptance & Recovery Tracking
  createReinsurerAcceptance(data) {
    return Api.post('/group-pricing/bordereaux/reinsurer/acceptances', data)
  },
  getReinsurerAcceptances(params) {
    return Api.get('/group-pricing/bordereaux/reinsurer/acceptances', {
      params
    })
  },
  getAcceptanceStats(generatedId) {
    return Api.get('/group-pricing/bordereaux/reinsurer/acceptances/stats', {
      params: { generated_id: generatedId }
    })
  },
  updateReinsurerAcceptance(id, data) {
    return Api.patch(
      `/group-pricing/bordereaux/reinsurer/acceptances/${id}`,
      data
    )
  },
  createReinsurerRecovery(data) {
    return Api.post('/group-pricing/bordereaux/reinsurer/recoveries', data)
  },
  getReinsurerRecoveries(params) {
    return Api.get('/group-pricing/bordereaux/reinsurer/recoveries', { params })
  },
  updateReinsurerRecovery(id, data) {
    return Api.patch(
      `/group-pricing/bordereaux/reinsurer/recoveries/${id}`,
      data
    )
  },

  // Claim Notification Cadence
  createClaimNotification(data) {
    return Api.post('/group-pricing/bordereaux/claim-notifications', data)
  },
  getClaimNotifications(params) {
    return Api.get('/group-pricing/bordereaux/claim-notifications', { params })
  },
  getNotificationStats(params) {
    return Api.get('/group-pricing/bordereaux/claim-notifications/stats', {
      params
    })
  },
  generateMonthEndNotifications(data) {
    return Api.post(
      '/group-pricing/bordereaux/claim-notifications/generate-month-end',
      data
    )
  },
  markNotificationSent(id, data) {
    return Api.post(
      `/group-pricing/bordereaux/claim-notifications/${id}/sent`,
      data
    )
  },
  markNotificationAcknowledged(id, data) {
    return Api.post(
      `/group-pricing/bordereaux/claim-notifications/${id}/acknowledged`,
      data
    )
  },
  deleteClaimNotification(id) {
    return Api.delete(`/group-pricing/bordereaux/claim-notifications/${id}`)
  },
  getClaimsByScheme(schemeId) {
    return Api.get(
      `/group-pricing/bordereaux/claim-notifications/claims-by-scheme/${schemeId}`
    )
  },
  exportNotificationsCSV(params) {
    return Api.get('/group-pricing/bordereaux/claim-notifications/export', {
      params,
      responseType: 'blob'
    })
  },

  // RI Treaty Management
  getTreaties(params = {}) {
    return Api.get('/group-pricing/reinsurance/treaties', { params })
  },
  getTreatyStats() {
    return Api.get('/group-pricing/reinsurance/treaties/stats')
  },
  createTreaty(data) {
    return Api.post('/group-pricing/reinsurance/treaties', data)
  },
  updateTreaty(id, data) {
    return Api.put(`/group-pricing/reinsurance/treaties/${id}`, data)
  },
  deleteTreaty(id) {
    return Api.delete(`/group-pricing/reinsurance/treaties/${id}`)
  },
  getTreatySchemeLinks(treatyId) {
    return Api.get(`/group-pricing/reinsurance/treaties/${treatyId}/schemes`)
  },
  linkSchemeToTreaty(treatyId, data) {
    return Api.post(
      `/group-pricing/reinsurance/treaties/${treatyId}/schemes`,
      data
    )
  },
  bulkLinkSchemesToTreaty(treatyId, data) {
    return Api.post(
      `/group-pricing/reinsurance/treaties/${treatyId}/schemes/bulk`,
      data
    )
  },
  removeSchemeTreatyLink(linkId) {
    return Api.delete(
      `/group-pricing/reinsurance/treaties/scheme-links/${linkId}`
    )
  },
  bulkRemoveSchemeLinks(treatyId, data) {
    return Api.delete(
      `/group-pricing/reinsurance/treaties/${treatyId}/schemes/bulk`,
      { data }
    )
  },
  getActiveTreatiesForScheme(schemeId) {
    return Api.get(
      `/group-pricing/reinsurance/treaties/active/scheme/${schemeId}`
    )
  },

  // RI Bordereaux Generation
  generateRIMemberBordereaux(data) {
    return Api.post('/group-pricing/reinsurance/bordereaux/member', data)
  },
  generateRIClaimsBordereaux(data) {
    return Api.post('/group-pricing/reinsurance/bordereaux/claims', data)
  },
  getRIBordereauxRuns(params = {}) {
    return Api.get('/group-pricing/reinsurance/bordereaux', { params })
  },
  getRIBordereauxStats(params = {}) {
    return Api.get('/group-pricing/reinsurance/bordereaux/stats', { params })
  },
  submitRIBordereaux(data) {
    return Api.post('/group-pricing/reinsurance/bordereaux/submit', data)
  },
  acknowledgeRIBordereaux(runId) {
    return Api.post(
      `/group-pricing/reinsurance/bordereaux/${runId}/acknowledge`
    )
  },
  getRIBordereauxMemberRows(runId) {
    return Api.get(`/group-pricing/reinsurance/bordereaux/${runId}/members`)
  },
  getRIBordereauxClaimsRows(runId) {
    return Api.get(`/group-pricing/reinsurance/bordereaux/${runId}/claims`)
  },
  diffRIBordereauxRun(runId, against?: string) {
    return Api.get(`/group-pricing/reinsurance/bordereaux/${runId}/diff`, {
      params: against ? { against } : {}
    })
  },
  validateRIBordereaux(runId) {
    return Api.post(`/group-pricing/reinsurance/bordereaux/${runId}/validate`)
  },
  getRIValidationResults(runId) {
    return Api.get(
      `/group-pricing/reinsurance/bordereaux/${runId}/validation-results`
    )
  },
  getRIBordereauxKPIs(params = {}) {
    return Api.get('/group-pricing/reinsurance/bordereaux/kpis', { params })
  },
  acknowledgeRIBordereauxReceipt(runId, data) {
    return Api.post(
      `/group-pricing/reinsurance/bordereaux/${runId}/acknowledge-receipt`,
      data
    )
  },
  amendRIBordereaux(runId, data) {
    return Api.post(
      `/group-pricing/reinsurance/bordereaux/${runId}/amend`,
      data
    )
  },

  // RI Large Claim Notices
  monitorLargeClaims(data) {
    return Api.post(
      '/group-pricing/reinsurance/bordereaux/large-claims/monitor',
      data
    )
  },
  getLargeClaimNotices(params = {}) {
    return Api.get('/group-pricing/reinsurance/bordereaux/large-claims', {
      params
    })
  },
  getLargeClaimStats(params = {}) {
    return Api.get('/group-pricing/reinsurance/bordereaux/large-claims/stats', {
      params
    })
  },
  updateLargeClaimNotice(id, data) {
    return Api.patch(
      `/group-pricing/reinsurance/bordereaux/large-claims/${id}`,
      data
    )
  },
  acceptLargeClaimNotice(
    id,
    data: { notes?: string; accepted_amount?: number } = {}
  ) {
    return Api.post(
      `/group-pricing/reinsurance/bordereaux/large-claims/${id}/accept`,
      data
    )
  },
  rejectLargeClaimNotice(id, data: { reason: string; notes?: string }) {
    return Api.post(
      `/group-pricing/reinsurance/bordereaux/large-claims/${id}/reject`,
      data
    )
  },
  queryLargeClaimNotice(id, data: { query_details: string }) {
    return Api.post(
      `/group-pricing/reinsurance/bordereaux/large-claims/${id}/query`,
      data
    )
  },

  // RI Technical Accounts & Settlement
  generateTechnicalAccount(data) {
    return Api.post('/group-pricing/reinsurance/settlement', data)
  },
  getTechnicalAccounts(params = {}) {
    return Api.get('/group-pricing/reinsurance/settlement', { params })
  },
  getTechnicalAccountByID(id) {
    return Api.get(`/group-pricing/reinsurance/settlement/${id}`)
  },
  getSettlementStats(params = {}) {
    return Api.get('/group-pricing/reinsurance/settlement/stats', { params })
  },
  updateTechnicalAccount(id, data) {
    return Api.patch(`/group-pricing/reinsurance/settlement/${id}`, data)
  },
  recordSettlementPayment(data) {
    return Api.post('/group-pricing/reinsurance/settlement/payments', data)
  },
  getSettlementPayments(params = {}) {
    return Api.get('/group-pricing/reinsurance/settlement/payments', { params })
  },
  escalateSettlementDispute(id, data) {
    return Api.post(
      `/group-pricing/reinsurance/settlement/${id}/escalate-dispute`,
      data
    )
  },
  resolveSettlementDispute(id, data) {
    return Api.post(
      `/group-pricing/reinsurance/settlement/${id}/resolve-dispute`,
      data
    )
  },

  // RI Cat Event Register
  getCatastropheClaimsRows(params = {}) {
    return Api.get('/group-pricing/reinsurance/bordereaux/cat-events', {
      params
    })
  },

  // RI Run-Off Treaty Register
  getRunOffTreaties(params = {}) {
    return Api.get('/group-pricing/reinsurance/bordereaux/run-off-treaties', {
      params
    })
  },

  // Win Probability
  getQuoteWinProbability(quoteId: number) {
    return Api.get(`/group-pricing/quotes/${quoteId}/win-probability`)
  },
  getWinProbabilityModelInfo() {
    return Api.get('/group-pricing/win-probability/model-info')
  },
  trainWinProbabilityModel() {
    return Api.post('/group-pricing/win-probability/train')
  },

  // Custom Tiered Income Replacement
  checkCustomTieredTableExists(schemeName: string, riskRateCode: string) {
    return Api.get('/group-pricing/custom-tiered-income-replacement/check', {
      params: { scheme_name: schemeName, risk_rate_code: riskRateCode }
    })
  },
  requestCustomTieredTable(data: {
    scheme_name: string
    scheme_id: number
    risk_rate_code: string
  }) {
    return Api.post(
      '/group-pricing/custom-tiered-income-replacement/request',

      data
    )
  },

  getRegionsForRiskCode(riskRateCode: string) {
    return Api.get('/group-pricing/region-loadings/regions', {
      params: { risk_rate_code: riskRateCode }
    })
  },

  // Migration — bulk scheme & member import
  validateMigration(formData: FormData) {
    return Api.post('/group-pricing/migration/validate', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        Accept: 'application/json'
      }
    })
  },
  executeMigration(formData: FormData) {
    return Api.post('/group-pricing/migration/execute', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        Accept: 'application/json'
      }
    })
  },
  downloadMigrationTemplate(templateName: string) {
    return Api.get(`/group-pricing/migration/templates/${templateName}`, {
      responseType: 'blob'
    })
  },
  verifyBankAccount(
    data: import('@/renderer/types/bav').VerifyBankAccountRequest
  ) {
    return Api.post<{
      success: boolean
      data: import('@/renderer/types/bav').VerifyResult
    }>('/v2/group-pricing/claims/verify-bank-account', data)
  },
  getBankVerificationStatus(jobId: string) {
    return Api.post<{
      success: boolean
      data: import('@/renderer/types/bav').VerifyResult
    }>(
      `/v2/group-pricing/claims/verify-bank-account/status/${encodeURIComponent(jobId)}`
    )
  },
  getOrgUsers(organization) {
    const json = JSON.stringify(organization)
    return Api.post('org-users', json, {
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json'
      }
    })
  },
  refreshOrgUsers(organization) {
    const json = JSON.stringify(organization)
    return Api.post('org-users/refresh', json, {
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json'
      }
    })
  },

  // ─── Email system ─────────────────────────────────────────────────────────
  getEmailSettings() {
    return Api.get('/group-pricing/email/settings')
  },
  saveEmailSettings(settings) {
    return Api.put('/group-pricing/email/settings', settings, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },
  sendTestEmail() {
    return Api.post('/group-pricing/email/settings/test')
  },

  listEmailTemplates() {
    return Api.get('/group-pricing/email/templates')
  },
  getEmailTemplate(code) {
    return Api.get(`/group-pricing/email/templates/${encodeURIComponent(code)}`)
  },
  createEmailTemplate(template) {
    return Api.post('/group-pricing/email/templates', template, {
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    })
  },
  updateEmailTemplate(code, template) {
    return Api.put(
      `/group-pricing/email/templates/${encodeURIComponent(code)}`,
      template,
      {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json'
        }
      }
    )
  },
  deleteEmailTemplate(code) {
    return Api.delete(
      `/group-pricing/email/templates/${encodeURIComponent(code)}`
    )
  },
  previewEmailTemplate(code, payload) {
    return Api.post(
      `/group-pricing/email/templates/${encodeURIComponent(code)}/preview`,
      payload,
      {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json'
        }
      }
    )
  },

  listEmailOutbox(
    params: { status?: string; page?: number; page_size?: number } = {}
  ) {
    return Api.get('/group-pricing/email/outbox', { params })
  },
  getEmailOutboxItem(id) {
    return Api.get(`/group-pricing/email/outbox/${id}`)
  },
  retryEmailOutbox(id) {
    return Api.post(`/group-pricing/email/outbox/${id}/retry`)
  }
}

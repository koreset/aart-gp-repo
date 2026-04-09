import Api from '@/renderer/api/Api'

export default {
  getTableMetaData() {
    return Api.get('/phi-valuation/table-meta-data')
  },
  uploadTables(formdata) {
    return Api.post('/phi-valuation/tables/upload', formdata, {
      headers: {
        'Content-Type': 'multipart/form-data',
        Accept: 'multipart/form-data'
      }
    })
  },
  getDataForTable(tableType) {
    return Api.get(`/phi-valuation/tables/${tableType}`)
  },
  deleteTable(tableType, year, version) {
    return Api.delete(
      `/phi-valuation/tables/${tableType}/year/${year}/version/${version}`
    )
  },
  saveShockSetting(shockSetting) {
    return Api.post('phi-valuation/shock-settings', shockSetting)
  },
  deleteShockSetting(shockSetting) {
    return Api.delete('phi-valuation/shock-settings/' + shockSetting)
  },
  getShockSettings(shockSetting) {
    return Api.get('phi-valuation/shock-settings', shockSetting)
  },
  getAvailableShockBases() {
    return Api.get('phi-valuation/shock-bases')
    // return Api.get("product-tables/" + productCode + "/get-basis")
  },
  getAvailableModelPointYears() {
    return Api.get('phi-valuation/model-point-years')
  },
  getModelPointVersionsForYear(year) {
    return Api.get(`phi-valuation/model-point-versions/year/${year}`)
  },
  getAvailableParameterYears() {
    return Api.get('phi-valuation/parameter-years')
  },
  getAvailableParameterVersions(year) {
    return Api.get(`phi-valuation/parameter-versions/year/${year}`)
  },
  getAvailableMortalityYears() {
    return Api.get('phi-valuation/mortality-years')
  },
  getAvailableMortalityVersions(year) {
    return Api.get(`phi-valuation/mortality-versions/year/${year}`)
  },
  getAvailableRecoveryYears() {
    return Api.get('phi-valuation/recovery-rate-years')
  },
  getAvailableRecoveryVersions(year) {
    return Api.get(`phi-valuation/recovery-rate-versions/year/${year}`)
  },
  getAvailableYieldCurveYears() {
    return Api.get('phi-valuation/yield-curve-years')
  },
  getAvailableYieldCurveVersions(year) {
    return Api.get(`phi-valuation/yield-curve-versions/year/${year}`)
  },
  runProjections(data) {
    const jsonPayload = JSON.stringify(data)

    return Api.post('phi-valuation/run-projections', jsonPayload, {
      headers: {
        'Content-Type': 'application/json'
      }
    })
  },
  getAllRunResults() {
    return Api.get(`phi-valuation/run-jobs`)
  },
  getRunResult(runId) {
    return Api.get(`phi-valuation/run-jobs/${runId}`)
  },
  getControlRunResult(runId) {
    return Api.get(`phi-valuation/run-jobs/${runId}/control`)
  },

  getAvailableYearsForTable(tableType) {
    return Api.get(`phi-valuation/tables/${tableType}/years`)
  },
  getAvailableVersionsForTableYear(tableType, year) {
    return Api.get(`phi-valuation/tables/${tableType}/years/${year}/versions`)
  },
  deleteProjectionJob(runId) {
    return Api.delete(`phi-valuation/run-jobs/${runId}`)
  },
  deleteValuationJobs(runIds) {
    return Api.post(`phi-valuation/run-jobs/delete`, runIds)
  },
  getPhiModelPointCount() {
    return Api.get('phi-valuation/model-point-count')
  },
  getPhiModelPointsForYear(year, version) {
    return Api.get(`phi-valuation/model-points/${year}/${version}`)
  },
  getPhiModelPointsExcel(year, version) {
    return Api.get(`phi-valuation/model-points/${year}/${version}/excel`, {
      responseType: 'blob'
    })
  },
  deletePhiModelPoints(year, version) {
    return Api.delete(`phi-valuation/model-points/${year}/${version}`)
  },
  getPhiRunConfigs() {
    return Api.get('phi-valuation/run-configs')
  },
  savePhiRunConfig(data) {
    return Api.post('phi-valuation/run-configs', data)
  },
  deletePhiRunConfig(id) {
    return Api.delete('phi-valuation/run-configs/' + id)
  },
    getValuationJobWithSpCode(id, spCode) {
    return Api.get('valuations/jobs/' + id + '/sp-code/' + spCode)
  },
    getValuationJobs() {
    return Api.get('valuations/jobs')
  },
  getJobExcelResults(jobId) {
    return Api.get('valuations/jobs/all-jobs/' + jobId + '/excel', {
      responseType: 'blob'
    })
  },
  getJobExcelScopedResults(jobId) {
    return Api.get('valuations/jobs/all-jobs/' + jobId + '/excel/scoped', {
      responseType: 'blob'
    })
  },
    getExcelResults(jobId, control) {
    if (control) {
      return Api.get('valuations/jobs/' + jobId + '/excel/control', {
        responseType: 'blob'
      })
    } else {
      return Api.get('valuations/jobs/' + jobId + '/excel', {
        responseType: 'blob'
      })
    }
  },

}

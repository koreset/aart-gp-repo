Redis caching candidates for DB queries

Overview
- Goal: Identify read-heavy, repetitive, and slow DB queries to benefit from Redis caching. This document lists concrete call sites, proposed cache keys, TTLs, and invalidation hints.
- Scope: Based on code inspection of services and engines using GORM (Find/First/Preload/Raw). Redis client exists at services/redis_client.go but is not yet wired for query-level caching.

General guidelines
- Good candidates: read-mostly reference data; lists with limited cardinality; expensive aggregates; metadata rarely changing; lookups repeated across requests or long-running workers.
- Avoid caching: highly dynamic rows (frequently updated), user-specific rapidly changing data without proper invalidation, or queries with large payloads if memory is a concern.
- Keys: use normalized names like ads:v1:<domain>:<selector>. Include parameters that change results. Prefer shorter keys; avoid spaces. Hash long SQL if needed.
- TTLs: start conservative (5–30 minutes) for lists, 1–24 hours for reference tables, 30–120 seconds for dashboards with near-real-time expectations.
- Invalidation: on writes to the underlying tables, delete the relevant keys. Where write paths aren’t centralized, rely on TTL.

Top priority candidates (high impact)
1) Product with deep preloads
   File: services/products.go
   Functions:
   - GetProductById(id int)
   - GetProductByCode(code string)
   - GetProductsAndFamilies(), GetAllAvailableProducts()
   Query shape: DB.Preload(...) chain for ProductFeatures, ProductTransitionStates, ProductTransitions, ProductModelpointVariables, ProductTables, ProductPricingTables, GlobalTables, ProductParameters.
   Why: Heavy multi-join/tree loads; used in many flows; product metadata changes infrequently.
   Suggested key(s):
   - ads:v1:product:id:{id}
   - ads:v1:product:code:{code}
   - ads:v1:products:families
   TTL: 30–120 minutes (metadata rarely changes).
   Invalidation: on Create/Update/Delete of products or related associations; otherwise rely on TTL.

2) Projection job and run metadata lists
   File: services/projections.go
   Functions:
   - GetJobs(runType int), GetValidJobs(), GetMostRecentJob(), GetProductsForJob(jobId int), GetJobsForProduct(prodCode string)
   - GetJobRunSettings(projectionJobId int), GetJobRunErrors(jobProductId int)
   Query shape: Preload("RunParameters"/"Products") and filters; ordering by id desc.
   Why: Queried frequently from UI and workers; low churn.
   Keys: 
   - ads:v1:jobs:runType:{runType}
   - ads:v1:jobs:mostRecent
   - ads:v1:jobs:forProduct:{prodCode}
   - ads:v1:job:{jobId}:products
   - ads:v1:job:{jobId}:settings
   - ads:v1:jobProduct:{jobProductId}:errors
   TTL: 5–30 minutes (dashboard-like usage).
   Invalidation: on job status updates/products association changes; otherwise rely on TTL.

3) Available years, months, codes (distinct lists)
   File: services/projections.go
   Functions:
   - GetAvailableYieldYears(), GetYieldCurveCodes(year), GetIbnrYieldCurveCodes(year), GetYieldCurveMonths(year, code)
   - GetYieldCurveMonthsv2(productCode, yieldYear, parameterYear, basis)
   - GetAvailableParameterYears(productCode)
   - GetAvailableMarginYears()
   - GetAvailableModelPointYears(productCode)
   - GetAvailableLapseYears/MortalityYears/... for various risk tables
   - GetAvailableBases(productCode), GetGMMShockBases(), GetIBNRBases(portfolio, year), GetIBNRShockBases()
   Query shape: DISTINCT projections over relatively small tables.
   Why: Classic dropdown source data; perfect for caching; recomputed often across sessions.
   Keys: 
   - ads:v1:avail:yieldYears
   - ads:v1:avail:yieldCodes:{year}
   - ads:v1:avail:yieldMonths:{year}:{code}
   - ads:v1:avail:paramYears:{productCode}
   - ads:v1:avail:mpYears:{productCode}
   - ads:v1:avail:bases:{scope}:{arg1}:{arg2}
   TTL: 60 minutes (or longer).
   Invalidation: when ingestion jobs load new years/codes; rely on TTL unless explicit invalidation is wired post-ingest.

4) Aggregated projection results for read APIs
   File: services/projections.go
   Functions:
   - GetAllAggregatedReserves(jobProductIds []interface{}, resultRange int, runVariable string)
   - GetAggregatedReserves(jobProductId int, resultRange int, runVariable string)
   - GetAggregatedProjections(jobProductId int, spCode string)
   - GetExcelAggregatedProjections(...), GetAggregatedProjectionsForDownload(...)
   Query shape: Summaries and joins on aggregated_projections and scoped tables; expensive for large jobs.
   Why: Post-run reporting repeatedly fetches same aggregates; results immutable after job completion.
   Keys:
   - ads:v1:aggRes:jp:{jobProductId}:range:{resultRange}:var:{runVariable}
   - ads:v1:aggProj:jp:{jobProductId}:sp:{spCode}
   - ads:v1:aggExcel:jp:{jobProductId}:sp:{spCode}:vg:{variableGroupId}
   TTL: 24 hours (or longer). Consider no TTL with explicit invalidation when a job is deleted.
   Invalidation: on deletion of a job/job product or re-run overwriting results.

5) Group pricing lists and roles with permissions
   File: services/group_pricing.go
   Functions:
   - Queries with Preload("SchemeCategories"): listing quotes by status/ordering
   - Role and permissions loads: roles, role by id with Permissions, user role with Permissions
   Why: Reused across UI; low churn.
   Keys:
   - ads:v1:gq:list:status:{filter}
   - ads:v1:gq:list:all
   - ads:v1:role:list
   - ads:v1:role:{roleId}
   - ads:v1:userRole:{orgUserId}
   TTL: 10–60 minutes.
   Invalidation: on CRUD in group pricing or permission changes.

Strong secondary candidates
6) IFRS17/LIC/CSM variable sets and configurations
   Files: services/ifrs17.go, services/lic.go, services/lic_engine.go, services/csm_engine.go
   Examples:
   - DB.Preload("LicVariables").Order("id desc").Find(&variableSets)
   - DB.Preload("AosVariables").Where("configuration_name = ?", ...).Find(&aos)
   Why: Configuration sets change infrequently; loaded repeatedly across runs.
   Keys:
   - ads:v1:lic:varsets:latest
   - ads:v1:csm:aos:config:{name}
   TTL: 60–240 minutes.

7) Valuation projection engine reference data
   File: services/valuation_projection_engine.go
   Query hotspots inside PopulateProjectionsForProduct and per-model-point function:
   - RunParameters by job, ShockSettings by id
   - Product states/transitions, features, parameters, multipliers, margins
   - Product shocks/time series references
   Why: These are reference inputs reused across many model points and multiple workers; ideal for cross-process caching so parallel workers don’t all hit DB.
   Keys:
   - ads:v1:run:{jobId}
   - ads:v1:shockSettings:{id}
   - ads:v1:product:{prodId}:states
   - ads:v1:product:{prodId}:features
   - ads:v1:product:{prodId}:parameters:{year}:{basis}
   - ads:v1:product:{prodId}:multipliers:{year}
   - ads:v1:product:{prodId}:margins:{year}
   - ads:v1:product:{prodId}:shocks:{year}:{basis}
   TTL: 120 minutes (or per run window). For long runs, consider no TTL and explicit invalidate on product/config change.

8) SAP/Manual finance and codes
   File: services/products.go (ProcessSAPFile, GetSapFileList, GetSapResultsForRunName), finance-related getters
   Why: UI lists/filters often re-queried; underlying data updated in batch.
   Keys:
   - ads:v1:sap:fileList
   - ads:v1:sap:results:{runName}
   - ads:v1:finance:versions:{year}:{measure}
   TTL: 60–240 minutes.

Implementation notes
- Use services/redis_client.go (InitRedis, RedisAvailable) for connection lifecycle. Add small helpers for JSON marshalling/unmarshalling to store structs/lists. Compress large payloads if needed.
- Pattern to apply around existing queries:
  1. Build key from function parameters.
  2. If RedisAvailable and GET hit, unmarshal and return.
  3. Otherwise query DB, then SET with TTL.
- For aggregates that depend on job completion, skip cache for in-progress jobs (status != Complete) to avoid staleness.

Example code snippet (pseudo-Go around a list query)
  key := fmt.Sprintf("ads:v1:avail:paramYears:%s", productCode)
  if RedisAvailable() {
    if bs, err := redisClient.Get(ctx, key).Bytes(); err == nil {
      var out []models.AvailableYieldResult
      if json.Unmarshal(bs, &out) == nil { return out }
    }
  }
  out := realDBQuery(...)
  if RedisAvailable() { b, _ := json.Marshal(out); redisClient.Set(ctx, key, b, time.Hour) }
  return out

Invalidation checklist
- When products or their dependent tables are updated (upload, process tables, process features), delete product-related keys for that product.
- After projection jobs complete and aggregated results are written, you may populate cache eagerly; when deleting a job, remove its job-related keys.
- For group pricing CRUD, clear relevant list/detail keys.

Next steps
- Prioritize implementing caching in products.go (GetProductById/GetProductByCode), projections.go (GetAvailable* functions), and projections aggregate getters. These give the largest benefit with minimal risk.
- Add a small json helper in services/redis_client.go for generic Get/Set of structs if desired.

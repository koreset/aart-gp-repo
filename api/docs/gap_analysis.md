# AART v5.5.0 — Gap Analysis

**Date:** 2026-02-28
**Scope:** Full application, with detailed focus on the Group Risk Insurance module
**Analyst:** Claude Code (automated codebase analysis)

---

## Executive Summary

The application has a solid actuarial computation foundation but carries significant gaps in
operationalization, error resilience, and several major workflows that are scaffolded but not yet
functional. The group risk insurance module is the most complete but has critical gaps in premium
lifecycle management, audit trail, and renewal pricing. Six other modules have stubbed or missing
core logic.

The core group risk rating engine, experience analysis engine, and IBNR engine are production-quality.
Everything else has material gaps — ranging from a disabled audit trail and missing transaction safety
(blocking all production use) to entire subsystems that exist as UI scaffolding with no backend
implementation (premium lifecycle, PAA, MGMM).

---

## Table of Contents

1. [Group Risk Insurance — Detailed Gap Analysis](#1-group-risk-insurance--detailed-gap-analysis)
   - 1.1 [What Is Fully Working](#11-what-is-fully-working)
   - 1.2 [Critical Gaps — Blockers for Production](#12-critical-gaps--blockers-for-production)
   - 1.3 [High-Priority Gaps](#13-high-priority-gaps)
   - 1.4 [Medium-Priority Gaps](#14-medium-priority-gaps)
2. [Application-Wide Gap Analysis](#2-application-wide-gap-analysis)
   - 2.1 [Module Status Summary](#21-module-status-summary)
   - 2.2 [Cross-Cutting Architectural Gaps](#22-cross-cutting-architectural-gaps)
3. [Prioritized Remediation Roadmap](#3-prioritized-remediation-roadmap)
4. [Appendix: Group Risk Engine Reference](#4-appendix-group-risk-engine-reference)

---

## 1. Group Risk Insurance — Detailed Gap Analysis

### 1.1 What Is Fully Working

The core **quote generation and rating engine** is genuinely complete and production-quality:

- Parallel member-level rating via a worker pool (`gammazero/workerpool`)
- All major benefit types: GLA, PTD, CI, SGLA, TTD, PHI, Family Funeral
- Experience rating with actuarial credibility blending (square-root credibility formula)
- Free cover limit calculation (percentile / scaling factor method)
- 3-level quota share reinsurance for both lump-sum and income benefits
- Full quote lifecycle: draft → in_progress → approved → accepted → in_force
- Member data ingestion via CSV with RSA ID validation and batch inserts
- Claims workflow: registration, assessment, document attachments, timeline tracking
- Broker management, scheme management, scheme exposure tracking (by age band, gender, benefit)
- Quote output PDF generation (`QuoteOutput.vue`)

---

### 1.2 Critical Gaps — Blockers for Production

#### GAP-01: Audit Trail Disabled

**File:** `services/group_pricing.go`, `ApproveGroupPricingQuote()`

The audit write inside `ApproveGroupPricingQuote()` is commented out with the note:

> `// ENSURE AUDITS ARE FIXED BEFORE PRODUCTION!!!`

Claim assessments have no audit trail at all. This is a **compliance and regulatory blocker** — any
IFRS 17 or insurance regulatory framework requires an immutable record of underwriting decisions.

**Impact:** All quote approvals and claim decisions are unaudited.
**Action required:** Re-enable the audit writer, extend it to claim assessments, and add audit entries
for all status transitions across the quote lifecycle.

---

#### GAP-02: Premium Lifecycle — Scaffolded, Not Functional

**Files:** `frontend/src/renderer/screens/group_pricing/premiums/` (all files)

The following frontend screens are marked "backlog" with empty TODO function bodies:

| Screen | Purpose |
|--------|---------|
| `PremiumSchedules.vue` | Monthly schedule list |
| `PremiumScheduleDetail.vue` | Per-schedule detail and line items |
| `PremiumDashboard.vue` | Portfolio premium overview |
| `Invoices.vue` | Invoice list |
| `InvoiceDetail.vue` | Invoice detail |
| `Payments.vue` | Record and view payments |
| `ArrearsManagement.vue` | Arrears aging and follow-up |
| `Statements.vue` | Employer and broker statements |
| `PremiumReconciliation.vue` | Collected vs ceded reconciliation |

The backend `GenerateMonthlySchedule()` and `GenerateInvoice()` functions exist in
`services/group_premiums.go` but have material deficiencies:

- Employer/employee contribution split is **hardcoded to 100% employer** — no configurable split
- Payment matching is one-to-one (one payment per invoice); partial payments and over-payments are
  not handled
- No premium adjustment/amendment workflow for mid-period changes
- No escalation or indexation of premiums between renewal periods
- No reconciliation between collected premiums and amounts ceded to reinsurers

**Impact:** In-force policy administration is not functional. Schemes cannot be operationally managed
after acceptance.
**Action required:** Implement the full premium schedule → invoice → payment → arrears → reconciliation
chain. Externalize the contribution split as a scheme-level configuration field.

---

#### GAP-03: Renewal / Repricing Engine Missing

**File:** `services/group_pricing.go`

Quote type "Renewal" is recognized and prior in-force member data is loaded correctly, but:

- There is no repricing calculation specific to renewals
- Historical claims are not considered when setting renewal rates (only the general experience rating
  mechanism applies, which requires manual credibility input)
- No automatic repricing trigger at policy anniversary
- No premium increase/decrease notification or letter generation
- No non-renewal or cancellation workflow
- No reinstatement workflow

**Impact:** The system cannot be used to manage the full policy lifecycle past the initial in-force
period.
**Action required:** Build a renewal orchestration service that retrieves the in-force quote, loads
updated claims experience, reprices at the new anniversary, and routes through the standard
approval flow.

---

#### GAP-04: Silent Calculation Failures

**File:** `services/group_pricing.go`, `CalculateGroupPricingQuote()` (~line 834)

A `// TODO: Implement the logic` comment exists inside the top-level calculation orchestrator.
Additionally, worker pool errors are logged to stdout but **not propagated** — a silent worker
failure means the UI may report "calculation complete" when only a subset of members were rated.
There is no transaction rollback if a batch of member ratings partially fails.

**Impact:** Calculations can silently produce incomplete results with no user-visible indication.
**Action required:** Propagate worker errors through the result channel; fail the overall calculation
job if any worker errors; wrap composite DB writes in a `db.Transaction()`.

---

### 1.3 High-Priority Gaps

#### GAP-05: Reinsurance — Quota Share Only

The 3-level proportional quota share is well implemented for both lump-sum and income benefits.
Missing capabilities:

- **Facultative reinsurance:** No case-by-case cession of large individual lives above the automatic
  acceptance limit
- **Facultative pricing overrides:** No per-member rate override for facultatively accepted risks
- **Reinsurance commission and profit commission:** Commission receivable from reinsurers is not
  calculated or tracked
- **Reinsurance claim recovery:** Claims are tracked in `GroupSchemeClaim` but there is no logic to
  calculate the ceded portion of an approved claim and generate a recovery receivable
- **Bordereaux reconciliation backend:** `BordereauxReconciliation.vue` (~1800 lines) has a
  sophisticated UI for discrepancy resolution and escalation, but the backend APIs for saving
  reconciliation outcomes appear absent

**Action required:** Implement facultative cession tracking per member, reinsurance commission
calculations, claim recovery logic, and the missing bordereaux reconciliation API endpoints.

---

#### GAP-06: Commission Model — Oversimplified

Commission is calculated as a flat `Total Premium × CommissionLoading` percentage. Missing:

- Tiered commission by benefit type, member count, or premium volume band
- Clawback/adjustment on early scheme exits
- Broker/intermediary commission split
- Commission statement generation (marked TODO in `premiums/Statements.vue`)
- Reinsurance (override) commission tracking
- Profit commission calculations

**Action required:** Extend the `Loadings` model to support tiered commission structures. Build
commission statement generation as part of the premium lifecycle (GAP-02).

---

#### GAP-07: Member Movements — Disconnected

**File:** `services/group_pricing.go`, `MovementPopulateRatesPerMember()` (~lines 1504–1870)

This function exists but is **never called** in the main calculation flow. Consequences:

- No pro-rata premium calculation for mid-period joiners or early exits
- Member exit dates default to `1900-01-01` (hardcoded sentinel value)
- No benefit deferment or suspension logic
- Member movements do not trigger premium adjustments on the in-force schedule

**Action required:** Wire `MovementPopulateRatesPerMember()` into the in-force premium calculation
path. Implement pro-rata logic. Replace the sentinel date with proper nullable handling.

---

#### GAP-08: Medical Underwriting — Not Present

All members are rated on occupation/industry class only. The system has no:

- Individual evidence of insurability workflow for members above the free cover limit
- Medical loading application to base rates
- Non-standard terms tracking (exclusions, loadings, deferred periods per member)
- Exclusion clause management or endorsement records

**Action required:** Design and implement a medical underwriting workflow that intercepts members
above the FCL, collects medical evidence, applies individual loadings, and stores the outcome
against the member's in-force record.

---

#### GAP-09: Dashboard and Analytics — Skeletal

The following analytical components have TODO blocks for their core content:

| Component | Status |
|-----------|--------|
| `GPDashBoard.vue` | Exposure analysis partially implemented |
| `ClaimsAnalyticsDashboard.vue` | TODO blocks for all charts |
| `BordereauxAnalyticsDashboard.vue` | TODO blocks for all charts |
| `PremiumDashboard.vue` | Entire component is backlog |
| `GetGroupPricingDashboardData()` API | Endpoint exists; calculation content unclear |
| `GetGroupSchemeClaimsDashboard()` API | Endpoint exists; calculation content unclear |

**Action required:** Implement the dashboard calculation queries and wire them to the frontend
components.

---

### 1.4 Medium-Priority Gaps

#### GAP-10: Rating Factors Not Utilised

The following rating dimensions exist in the data model but have no corresponding rating logic:

| Factor | Model Field | Rating Logic |
|--------|------------|-------------|
| Smoking status | Present | Absent |
| Health / medical history | Present | Absent |
| Occupation sub-class | Present (beyond industry loading) | Absent |
| Region / territory | Present | Absent |
| Per-member claims history loading | Present | Absent |

---

#### GAP-11: Hardcoded Values

The following values are embedded in code and should be configurable:

| Location | Hardcoded Value | Should Be |
|----------|----------------|-----------|
| `services/group_premiums.go` | 100% employer contribution | Per-scheme config field |
| `services/group_pricing.go` | Year literals 2024 / 2025 | Derived from system date |
| `services/group_pricing.go` | CSV batch size = 100 | App config |
| `services/group_pricing.go` | Worker pool = `min(CPU, 8)` | App config |
| Multiple locations | Default date `1900-01-01` | Nullable with explicit null checks |

---

#### GAP-12: Date Handling — Fragile

A comment in the codebase reads:

> *"v-date-input component has a bug that sends the date as the day before"*

The stated workaround is `AddDate(0, 0, 0)` — which is a no-op. Experience rating date parsing has
a silent failure path that falls back to a zero value with no user notification. Financial year
calculations assume a January year-end, which will produce incorrect results for schemes with
non-January financial year-ends (configurable via `GroupPricingInsurerDetail.YearEndMonth` but not
plumbed through the calculations).

---

#### GAP-13: Cache Invalidation — Naive

The calculation cache (`GroupPricingCache`) is cleared entirely at the start and end of every
calculation run. This means concurrent quote calculations for different schemes invalidate each
other's caches. There is no key-scoped or version-aware cache invalidation strategy.

---

#### GAP-14: Educator Benefit — Skeletal

`PopulateRatesPerMember()` includes educator benefit logic but:

- There is no mapping from a member's occupation field to educator grade (Grade 0, 1–7, 8–12,
  Tertiary)
- Sum assured amounts are read from generic parameters rather than a grade-specific structure
- No validation that the employer is an educational institution

---

#### GAP-15: Quote PDF — No Template Management

`QuoteOutput.vue` (~1700 lines) is monolithic with:

- A hardcoded logo placeholder comment
- No dynamic header/footer based on insurer configuration
- No template versioning or management
- No multi-language support for quote documents

---

## 2. Application-Wide Gap Analysis

### 2.1 Module Status Summary

| Module | Core Logic | Persistence | UI | Error Handling | Production Ready |
|--------|:----------:|:-----------:|:--:|:--------------:|:----------------:|
| Group Risk Pricing Engine | ✅ | ✅ | ✅ | ⚠️ | ❌ Audit / premium gaps |
| Experience Analysis | ✅ | ✅ | ✅ | ⚠️ | ⚠️ Near-ready |
| IBNR Engine | ✅ | ✅ | ✅ | ⚠️ | ⚠️ Near-ready |
| CSM Engine (IFRS 17) | ⚠️ | ✅ | ✅ | ⚠️ | ❌ Stub function (GAP-16) |
| LIC Engine | ⚠️ | ⚠️ | ✅ | ⚠️ | ❌ Incomplete buildup (GAP-19) |
| PHI Valuations | ⚠️ | ✅ | ✅ | ⚠️ | ❌ Calculations absent |
| PAA | ⚠️ | ✅ | ✅ | ⚠️ | ❌ Core logic absent (GAP-17) |
| GMM / VFA (MGMM) | ⚠️ | ⚠️ | ✅ | ⚠️ | ❌ Core logic absent (GAP-18) |
| Non-Group Pricing | ⚠️ | ✅ | ✅ | ⚠️ | ❌ Retrieval only, no rate calc |
| Tasks | ⚠️ | ✅ | ✅ | ❌ | ❌ No priority / dependencies |

---

### 2.2 Cross-Cutting Architectural Gaps

#### GAP-16: CSM Initial Recognition — Empty Stub

**File:** `services/csm_engine.go`, `ProcessInitialRecognitionAnalysis()` (~line 416)

```go
// TODO: build processing function
```

This is a **required IFRS 17 measurement at policy inception**. Its absence means the CSM engine
cannot produce compliant initial recognition measurements. The function is called from the engine
orchestration but does nothing.

**Action required:** Implement initial recognition analysis including CSM at inception, coverage
units, and risk adjustment at first recognition date.

---

#### GAP-17: PAA Core Logic — Absent

**File:** `services/paa.go`

Contains only 2 utility functions (`GetPAAPortfolioNames`, `GetPAARuns`). The actual PAA
calculation logic — eligibility testing, premium allocation, loss component computation, onerous
contract assessment — is entirely absent. The UI screens exist and data persists, but no PAA
measurements are computed.

---

#### GAP-18: GMM/VFA Core Logic — Absent

**File:** `services/mgmm.go`

Contains only 4 validation helper functions. The GMM/VFA expected cash flow projections, actuarial
assumption application, risk adjustment calculation, and financial measure computations are absent.

---

#### GAP-19: LIC Engine — Incomplete Buildup

**File:** `services/lic_engine.go`

The buildup calculation framework is partially written but the switch statement for variance types
is incomplete. Interest accretion placeholders contain hardcoded zero values. Treaty recovery
allocation logic references variables that appear to be undefined in context. The full file contains
two TODOs:

- `// TODO: Delete related tables`
- `// TODO: Add data summaries for the new tables @Motlatsi`

---

#### GAP-20: Database Resilience Pattern — Inconsistently Applied

`services/db.go` defines `DBReadWithResilience()` (exponential backoff with jitter, 5-second
context timeout, 40-connection concurrency gate). This pattern is used in some services but
**not** in:

- `services/csm_engine.go` — all direct GORM calls
- `services/ibnr_engine.go` — all direct GORM calls
- `services/ifrs17.go` — all direct GORM calls
- `services/lic_engine.go` — all direct GORM calls

Under database load or a momentary connection saturation event, these modules will fail hard
rather than retry.

---

#### GAP-21: No Transaction Safety for Composite Writes

The following operations write to multiple tables without a wrapping database transaction:

- `AcceptGroupPricingQuote()` — copies member data, updates scheme status, creates exposure records,
  calculates stats
- CSM engine run persistence — multiple result tables written sequentially
- IBNR run — model points and development factors written separately

A failure mid-sequence leaves the database in an inconsistent state with no rollback.

---

#### GAP-22: Authentication and Authorisation — Effectively Disabled

**Files:** `routes/routes.go`, `frontend/src/renderer/router/index.ts`

The backend `GetActiveUser` middleware is defined but its implementation passes all requests
regardless of token validity. The frontend router's `checkPermissions` function explicitly notes
in a comment that it "always returns true" as a temporary measure (line ~17, `router/index.ts`).

**Impact:** Any authenticated session can access any screen and any API endpoint regardless of
entitlements or role. The entitlement system (Pinia store, route guards) is wired but bypassed.

---

#### GAP-23: No Graceful Shutdown

**File:** `main.go`

The application starts projection workers, Redis connections, and a database connection pool but
registers no `os.Signal` handler for `SIGTERM` or `SIGINT`. In-flight calculations will be
aborted without checkpointing on any controlled restart or container stop event.

---

#### GAP-24: Structured Logging — Not Used

Multiple services use `fmt.Println()` for error output. The application has a `globals.Logger`
(structured logger) but its usage is inconsistent. In production, `fmt.Println` output is not
captured by log aggregators, making incident diagnosis very difficult.

---

#### GAP-25: Frontend — 21 Near-Identical Valuation Components

The `frontend/src/renderer/screens/valuations/` directory has four sub-modules (GMM, PAA, IBNR,
LIC), each with 4–6 near-identical components (RunSettings, RunResults, RunDetail, Tables,
ShockSettings). There is no shared base component. A UX change to the common run management
pattern requires editing 15–21 separate files.

---

#### GAP-26: Frontend — Route Permission Guards Incomplete

The router (`router/index.ts`) defines `entitlementGuard` and per-route permission metadata, but
the underlying `checkPermissions` function always returns `true`. Entitlement-based UI visibility
(Pinia store `v-if` bindings) is working, but **route-level access control is not enforced**.

---

## 3. Prioritized Remediation Roadmap

### P0 — Must Fix Before Any Production Use

These items represent compliance, data integrity, or security blockers.

| # | Item | File(s) |
|---|------|---------|
| P0-1 | Re-enable audit trail in `ApproveGroupPricingQuote()`; add audit to claim assessments and all status transitions | `services/group_pricing.go` |
| P0-2 | Wrap composite DB writes (AcceptQuote, CSM runs, IBNR runs) in `db.Transaction()` | Multiple services |
| P0-3 | Propagate worker pool errors in `CalculateGroupPricingQuote()`; fail the job on any worker error | `services/group_pricing.go` |
| P0-4 | Implement `ProcessInitialRecognitionAnalysis()` (IFRS 17 initial recognition) | `services/csm_engine.go` |
| P0-5 | Enable backend authentication middleware and frontend route permission guards | `routes/routes.go`, `router/index.ts` |
| P0-6 | Add graceful shutdown signal handler | `main.go` |

### P1 — Required for Group Risk Go-Live

| # | Item | File(s) |
|---|------|---------|
| P1-1 | Implement full premium schedule → invoice → payment → arrears → reconciliation chain | `services/group_premiums.go`, `premiums/*.vue` |
| P1-2 | Externalize employer contribution split as a scheme-level config field | `services/group_premiums.go`, `models/group_pricing.go` |
| P1-3 | Build renewal/repricing engine (orchestration, claims experience input, approval routing) | `services/group_pricing.go` |
| P1-4 | Wire `MovementPopulateRatesPerMember()` into in-force calculation; implement pro-rata premiums | `services/group_pricing.go` |
| P1-5 | Implement bordereaux reconciliation backend APIs | `controllers/group_pricing.go` |
| P1-6 | Apply `DBReadWithResilience` to all services not currently using it | All services |
| P1-7 | Fix date handling: date-picker bug workaround, silent parse failure, financial year config | `services/group_pricing.go`, frontend components |
| P1-8 | Externalize all hardcoded configuration values (batch size, worker count, year literals) | `services/group_pricing.go` |

### P2 — Required for Full IFRS 17 Compliance

| # | Item | File(s) |
|---|------|---------|
| P2-1 | Build PAA measurement logic (eligibility, premium allocation, loss component, onerous testing) | `services/paa.go` |
| P2-2 | Build GMM/VFA full calculation engine (cash flow projection, risk adjustment, CSM) | `services/mgmm.go` |
| P2-3 | Complete LIC engine buildup calculations and treaty recovery allocations | `services/lic_engine.go` |
| P2-4 | Complete PHI valuation calculation engine (shock scenarios, recovery rate application) | `services/phi_valuations.go` |

### P3 — Quality, Completeness, and Maintainability

| # | Item | File(s) |
|---|------|---------|
| P3-1 | Implement commission tiering, clawback, and statement generation | `services/group_pricing.go`, `premiums/Statements.vue` |
| P3-2 | Add facultative reinsurance cession and reinsurance claim recovery logic | `services/group_pricing.go` |
| P3-3 | Medical underwriting workflow (FCL excess, individual loading, endorsements) | New service + model |
| P3-4 | Complete dashboard and analytics (claims, bordereaux, premium dashboards) | Multiple backend + frontend |
| P3-5 | Replace all `fmt.Println` error output with `globals.Logger` structured logging | All services |
| P3-6 | Refactor 21 near-identical valuation frontend components into shared base component | `screens/valuations/**` |
| P3-7 | Implement educator benefit grade-to-occupation mapping | `services/group_pricing.go` |
| P3-8 | Quote PDF template management (configurable header/footer, multi-language) | `screens/group_pricing/QuoteOutput.vue` |
| P3-9 | Implement scope-aware cache invalidation (replace clear-all strategy) | `services/group_pricing.go` |

---

## 4. Appendix: Group Risk Engine Reference

> For full field-level configuration and calculation reference — including `risk_rate_code` seeding
> requirements, tax table and tiered income replacement mechanics, loading formula, FCL parameters,
> and output field semantics — see **[group_risk_pricing_reference.md](group_risk_pricing_reference.md)**.

### Supported Benefit Types

| Benefit | Description | Rating Logic | Reinsurance |
|---------|-------------|-------------|-------------|
| GLA | Group Life Assurance | ✅ Full | ✅ |
| SGLA | Spouse GLA | ✅ Full | ✅ |
| PTD | Permanent Total Disability | ✅ Full (accelerated/non-accelerated) | ✅ |
| CI | Critical Illness | ✅ Full (accelerated option) | ✅ |
| TTD | Temporary Total Disability | ✅ Full | ✅ |
| PHI | Permanent Health Insurance | ✅ Full (premium waiver, medical aid waiver) | ✅ |
| Family Funeral | Main member, spouse, children, parents | ✅ Full (5 variants) | ✅ Optional |
| Educator | Grade-linked tuition + allowances | ⚠️ Partial (no occupation mapping) | ❌ |
| Accidental TTD | Separate rate table exists | ❌ Not integrated | ❌ |

### Rating Factors

| Factor | Implemented |
|--------|-------------|
| Age at next birthday | ✅ |
| Gender | ✅ |
| Annual salary / income level | ✅ |
| Industry / occupation loading | ✅ |
| Salary multiples (GLA, PTD, CI) | ✅ |
| Sum assured caps | ✅ |
| Waiting / deferred periods | ✅ |
| Experience / credibility rating | ✅ |
| Commission, expense, profit loadings | ✅ |
| Free cover limit | ✅ |
| Smoking status | ❌ |
| Medical history | ❌ |
| Occupation sub-class | ❌ |
| Region / territory | ❌ |
| Per-member claims loading | ❌ |

### Quote Calculation Flow

```
1. CreateQuote
   Input:  scheme name, industry, benefit configuration, member CSV, basis
   Output: GroupPricingQuote, SchemeCategory records

2. CalculateQuote(quoteId, basis, [credibilityOverride])
   ├── For each SchemeCategory (parallel):
   │   ├── Calculate Free Cover Limit
   │   │   └── percentile(salary distribution) × ScalingFactor × sqrt(n)
   │   ├── For each Member (parallel worker pool):
   │   │   ├── Age at next birthday
   │   │   ├── Income level lookup
   │   │   ├── Per-benefit: BaseRate × IndustryLoading × ExperienceAdj
   │   │   ├── Cap sum assured / income
   │   │   ├── RiskPremium = Rate × CappedSA
   │   │   ├── OfficePremium = RiskPremium / (1 − TotalLoading)
   │   │   └── Reinsurance cession (3-level quota share)
   │   └── Aggregate to MemberRatingResultSummary
   └── Save: MemberRatingResult, MemberPremiumSchedule, Bordereaux, Summary

3. ApproveQuote        → status: approved
4. AcceptQuote         → copy members to in-force, update scheme status, create exposure records
```

### Experience Rating Formula

```
Credibility     = min( sqrt(WeightedLifeYears / FullCredibilityThreshold), 1.0 )
AdjustedRate    = TheorRate × (1 − Credibility) + Credibility × ExperienceRate
ExperienceRatio = AdjustedRate / TheorRate
```

### Reinsurance Structure

```
3 levels for lump-sum benefits (GLA, PTD, CI, SGLA):
  Level i: [Lowerbound_i, Upperbound_i] → CededProportion_i

3 levels for income benefits (TTD, PHI):
  Same structure applied to monthly benefit amount

If IsLumpsumReinsGLADependent:
  Dependent benefits retain proportional to GLA retained/capped SA

Funeral reinsurance (optional): 5 variants, same 3-level structure
```

# Bordereaux Module — Functional Gap Analysis

**Date:** 2026-04-17
**Scope:** End-to-end audit of the bordereaux module — Go API (controllers, services, models, routes) and Vue/Electron frontend (screens, components, services, router).
**Lens:** Functional gaps only — missing features, incomplete workflows, unhandled edge cases, stubs, and wiring disconnects.
**Analyst note:** All claims in this document were verified by reading the referenced files at the stated line numbers. Where earlier automated scans produced false positives (e.g., functions falsely reported as stubs), those are not in this report.

---

## Executive Summary

The bordereaux module is one of the **more feature-complete areas** of the platform: eight distinct sub-capabilities (templates, outbound generation, inbound submissions, reconciliation, reinsurer tracking, RI runs, deadlines, claim notifications) are each wired from backend through UI. It is considerably further along than the premium lifecycle, PAA, or MGMM modules called out in `api/docs/gap_analysis.md`.

However, the module has a cluster of **functional gaps that cluster around three themes**:

1. **"Happy-path UI, silent failure path."** A recurring pattern across many components: the primary action works, but catch blocks contain `// TODO: Show error notification` and log to console — so failures are invisible to users. Template testing, template export, version history, compliance report generation, submission-to-scheme, and analytics export are all **console.log stubs**.

2. **Analytics and dashboards are cosmetic.** `BordereauxAnalyticsDashboard.vue` renders entirely from three hardcoded arrays — no API call, no charts, period filter is a no-op. `BordereauxManagement.vue`'s compliance report button is a stub. These are UI scaffolding, not features.

3. **Downstream wiring is brittle or missing.** Discrepancy escalation sets a status but triggers no notification, no queue, no SLA. Claim recovery auto-generation exists as a function but is never invoked by the claim lifecycle. Overdue deadlines/notifications are swept in-memory during read requests, not by a scheduler. The download endpoint has no per-record authorization. The file-retention story is unbounded disk growth.

The five most-production-blocking items are listed in Section 8. A prioritized roadmap is in Section 9.

---

## 0. Implementation Status (session of 2026-04-17)

Items in the roadmap below have been worked through. Status marks reflect code in this branch, not yet deployed or covered by tests unless noted.

### P0 — Production Blockers

| Item | Status | Notes |
|---|---|---|
| P0-1 Download authorization | ✅ Done | `services.AuthorizeBordereauxDownload` gates `DownloadBordereaux` by creator / reviewer / approver on the `GeneratedBordereaux` row; 404 on missing row prevents filename enumeration; audit entry written on each download. |
| P0-2 Transaction wrap on confirmation import | ✅ Done | `ReconcileConfirmation` now takes `*gorm.DB`; import loop and reprocessing paths wrap it in `DB.Transaction` so partial writes roll back cleanly. |
| P0-3 Claim recovery hook | ✅ Done (hardened) | The hook was already wired at `UpdateGroupSchemeClaim` / `UpdateGroupSchemeClaimWithFiles`; `_, _ =` replaced with structured `appLog` error logging; `GenerateClaimRecovery` accepts `*gorm.DB` so the WithFiles path commits/rolls back the recovery with the claim. |
| P0-4 submitToScheme / checkStatus | ✅ Done | `BordereauxSubmissionTracking.vue` handlers wired to existing `submitBordereauxBatch` / `getBordereauxById`; fixed a latent wrong URL in `getBordereauxById`. |
| P0-5 Notification goroutine errors | ✅ Done | `notifyIfNotSenderE` variant propagates `CreateNotification` error; `NotifySubmission*` return error; inbound call sites wrap in a goroutine that calls new `logNotifyFailure(event, submission_id, err)`. |
| §8 #2 Template delete (local-only) | ✅ Done | Confirmation dialog + real `deleteBordereauxTemplate` call in `BordereauxTemplateManager.vue`; local array only mutates on 2xx. |

### P1 — Required for Go-Live

| Item | Status | Notes |
|---|---|---|
| P1-1 Analytics backend | ✅ Done | New `services/bordereaux_analytics.go` + `GET /bordereaux/analytics`; `BordereauxAnalyticsDashboard.vue` rewritten to consume it with real KPIs, benefit / insurer / monthly tables and an ag-charts monthly trend line. SA-specific compliance / error panels removed (no data source). |
| P1-2 Escalation workflow | ✅ Done | `bordereaux_reconciliation_results` migration adds `assigned_to`, `priority`, `escalated_by`, `escalated_at`, `due_date`; `EscalateDiscrepancy` fills them and computes an SLA due date by priority; new `NotifyDiscrepancyEscalated` hook; new `ResolveEscalationTarget` resolves email / name / role; new `GET /bordereaux/reconciliation/escalations` queue endpoint. |
| P1-3 File retention | ✅ Done | `BORDEREAUX_FILE_RETENTION_DAYS` config (default 90); new `services/bordereaux_file_retention.go` runs daily; three sweeps for `GeneratedBordereaux` (any age), `BordereauxConfirmation` (terminal only), `EmployerSubmission` (terminal only); removes file, clears `FilePath`/`FileName`, logs byte-freed summary. |
| P1-4 Scheduled overdue sweeps | ✅ Done | New `SweepOverdueDeadlines` + `StartDeadlineOverdueSweeper` (15-min ticker) registered in `main.go` alongside the existing notification sweeper. Cross-DB safe (per-row grace computed in Go). |
| P1-5 Member-bordereaux variance | ✅ Done | `ReconcileConfirmation` now populates a `memberExpectedMap` from `MemberBordereauxData`, matches incoming confirmation rows by employee / id / name key, produces matched / missing / extra results instead of phantom counts. |
| P1-6 Silent error catches (flash) | ✅ Done | `useFlashStore` threaded through `BordereauxSubmissionTracking`, `BordereauxTemplateManager`, `BordereauxGenerationForm`, `BordereauxReconciliation`; all `alert()` calls removed, ~30 TODO-commented catches now surface errors; success-path toasts added to main user actions. |
| P1-7 Child component permission checks | ⏳ Blocked | Can't harden frontend enforcement while backend `GetActiveUser` passes every authenticated request (see GAP-22 in `api/docs/gap_analysis.md`). |

### P2 — Completeness

| Item | Status | Notes |
|---|---|---|
| P2-1 Template preview / test / export | ✅ Done (version history deferred) | Preview dialog, JSON export of one / all templates, backend `POST /bordereaux/templates/:id/test` returns a live-data preview of field mappings with unknown-field / missing-in-data warnings. Version history left as an honest "not available yet" flash; needs a history table. |
| P2-2 Compliance report | ✅ Done | New `services/bordereaux_compliance_report.go` builds a 5-sheet xlsx (Summary / Open Discrepancies / Overdue Deadlines / Escalations / Large Claim Notices); `GET /bordereaux/compliance-report?from=&to=` streams it; `BordereauxManagement.vue` "Generate Report" wired to download blob. |
| P2-3 Per-scheme reconciliation tolerance | ✅ Done | `GroupScheme.ReconciliationTolerance` with migration; `reconciliationToleranceForScheme` helper replaces hardcoded `0.001` in `ReconcileConfirmation`; zero falls back to default for backward compatibility. |
| P2-4 Reconciliation notes column | ✅ Done | New `BordereauxConfirmationNote` table + migration; `AddReconciliationNote` writes there instead of fabricating a `_note` synthetic reconciliation row; new `GET /bordereaux/confirmations/:id/notes`; legacy filter kept as a safety net. |
| P2-5 Reinsurer response on large-claim notices | ✅ Done | `LargeClaimNotice` gains `response_status` / `responded_at` / `responded_by` (with migration); three new endpoints `POST large-claims/:id/{accept,reject,query}` with audit entries and note history via `appendResponseNote`. |
| P2-6 RI run amendment diff | ✅ Done | New `DiffRIBordereauxRuns` produces header / member-row / claim-row diffs with add/remove/change grouping; `GET /reinsurance/bordereaux/:run_id/diff?against=<run_id>`; defaults to `ParentRunID` when `against` omitted. |
| P2-7 Status strings → enums + transition guards | ⏳ Pending | Intentionally deferred — wide-reaching refactor across ~15 models and their services; deserves its own branch. |

### Outside the bordereaux module

| Item | Status | Notes |
|---|---|---|
| GAP-22 (auth & authorization disabled) | ⏳ Pending | Explored and planned in this session but deferred; see `api/docs/gap_analysis.md` for the original framing. JWT signature validation is the load-bearing first step. |

### Migrations added this session

| File | Purpose |
|---|---|
| `20260417120000_escalation_fields_bordereaux_reconciliation.sql` | P1-2 escalation fields |
| `20260417130000_group_schemes_reconciliation_tolerance.sql` | P2-3 tolerance column |
| `20260417140000_create_bordereaux_confirmation_notes.sql` | P2-4 notes table |
| `20260417150000_large_claim_notice_response_fields.sql` | P2-5 reinsurer response fields |

All present in `api/migrations/{postgresql,mysql,mssql}/`. Run in timestamp order.

### New background workers

Registered in `main.go` alongside the existing sweepers:
- `StartDeadlineOverdueSweeper` (15-min ticker) — P1-4.
- `StartBordereauxFileRetentionSweeper` (24h ticker, runs at boot) — P1-3.

### Known follow-ups

- Template **version history** (P2-1 residual) needs a history table + audit wiring — currently returns an honest "not available yet" flash.
- **P2-7** status enums + transition guards — deferred refactor.
- **P1-7 / GAP-22** frontend permission enforcement can't be meaningfully tightened until JWT signature verification and route-level enforcement land.
- Export Report on `BordereauxAnalyticsDashboard` was removed (was a no-op); an xlsx export can be added as a follow-up if analysts need one beyond the compliance report.

---

## 1. Module Inventory

### 1.1 Backend (Go)

**Route groups** (all under `api/routes/routes.go`, group `apiv1.Group("", GetActiveUser()) → /group-pricing`):

| Group | Routes | Handler File |
|---|---|---|
| `bordereaux/templates` | 5 | `controllers/bordereaux_templates.go` |
| `bordereaux/generate`, `bordereaux/fields`, `bordereaux/configurations`, `bordereaux/generated`, `bordereaux/download`, `bordereaux/dashboard` | ~20 | `controllers/bordereaux.go` |
| `bordereaux/confirmations`, `bordereaux/reconciliation` | ~14 | `controllers/bordereaux.go`, `controllers/bordereaux_reconciliation.go` |
| `bordereaux/reinsurer/*` | 7 | `controllers/bordereaux_reinsurer.go` |
| `bordereaux/deadlines` | 5 | `controllers/bordereaux_deadlines.go` |
| `bordereaux/claim-notifications` | 9 | `controllers/bordereaux_claim_notifications.go` |
| `reinsurance/bordereaux/*` (RI runs) | ~20 | `controllers/reinsurance_bordereaux.go` |
| `bordereaux/submissions` (inbound) | ~14 | `controllers/bordereaux_inbound.go` |

**Total: ~95 bordereaux-related routes**, all protected by the `GetActiveUser()` middleware (`routes/routes.go:25`).

**Services:** `bordereaux.go`, `bordereaux_inbound.go`, `bordereaux_templates.go`, `bordereaux_reconciliation.go`, `bordereaux_reinsurer.go`, `bordereaux_deadlines.go`, `bordereaux_claim_notifications.go`, plus `reinsurance_validation.go` for the 3-level RI validator.

**Models:** 20+ models spanning `bordereaux_template.go`, `bordereaux_inbound.go`, `bordereaux_deadlines.go`, `bordereaux_reinsurer.go`, `bordereaux_claim_notifications.go`, `reinsurance_bordereaux.go`, and embedded types in `group_pricing.go`.

### 1.2 Frontend (Vue/Electron)

Parent: `app/src/renderer/screens/group_pricing/bordereaux_management/BordereauxManagement.vue`

Children (all under `.../bordereaux_management/components/`):

| Component | Router path | Status |
|---|---|---|
| `BordereauxGenerationForm.vue` | `/generation` | Functional; missing UX feedback |
| `BordereauxSubmissionTracking.vue` | `/tracking` | Functional; submit-to-scheme and status-check are stubs |
| `BordereauxReconciliation.vue` | `/reconciliation` | Functional; no snackbar notifications |
| `BordereauxTemplateManager.vue` | `/templates` | CRUD works; preview/test/export/history are stubs |
| `BordereauxAnalyticsDashboard.vue` | `/analytics` | **Mock data only** |
| `BordereauxCalendar.vue` | `/deadline-calendar` | Functional |
| `BordereauxInboundSubmissions.vue` | `/inbound-submissions` | Functional |
| `BordereauxInboundSubmissionDetail.vue` | `/inbound-submissions/:id` | Functional |
| `BordereauxReinsurerTracking.vue` | `/reinsurer-tracking` | Functional |
| `RITreatyManagement.vue` | `/ri-treaties` | Out of scope (treaties, not bordereaux) |
| `RIBordereauxKPIDashboard.vue` | `/ri-kpi-dashboard` | Functional |
| `RISubmissionRegister.vue` | `/ri-submission-register` | Out of scope |
| `RIBordereauxGeneration.vue` | `/ri-bordereaux` | Functional |
| `RIClaimsBordereaux.vue` | `/ri-claims` | Functional |
| `RITechnicalAccounts.vue` | `/ri-settlement` | Out of scope |
| `BordereauxClaimNotifications.vue` | `/claim-notifications` | Functional |

No dedicated Pinia store for bordereaux. No composable. API calls go through `app/src/renderer/api/GroupPricingService.ts`.

---

## 2. Frontend Functional Gaps

### 2.1 Stub action handlers (buttons that do nothing)

These handlers are wired to UI buttons but their bodies are a `console.log` and a `// TODO` comment:

| # | File | Line | Handler | Effect |
|---|---|---|---|---|
| F-01 | `BordereauxManagement.vue` | 896 | `generateComplianceReport()` | Closes dialog, logs to console. No API call. "Generate Report" button is cosmetic. |
| F-02 | `BordereauxAnalyticsDashboard.vue` | 947 | `loadAnalytics()` | Period-filter `@change` handler does nothing. Mock data stays static. |
| F-03 | `BordereauxAnalyticsDashboard.vue` | 955 | `exportReport()` | "Export Report" button logs to console. No file produced. |
| F-04 | `BordereauxTemplateManager.vue` | 904 | `viewTemplate()` | Template "View" row action is a no-op. No preview dialog implemented. |
| F-05 | `BordereauxTemplateManager.vue` | 927 | `testTemplate()` | "Test Template" action (row and editor) logs to console. No test harness. |
| F-06 | `BordereauxTemplateManager.vue` | 932 | `testCurrentTemplate()` | As above, for the template currently being edited. |
| F-07 | `BordereauxTemplateManager.vue` | 1027 | `exportTemplate()` | "Export" action does not produce a file. |
| F-08 | `BordereauxTemplateManager.vue` | 1032 | `exportTemplates()` | "Export all" does not produce a zip/file. |
| F-09 | `BordereauxTemplateManager.vue` | 1037 | `versionHistory()` | No version history UI or backend exists; button is cosmetic. |
| F-10 | `BordereauxSubmissionTracking.vue` | 937 | `submitToScheme()` | Core action — submitting an outbound bordereaux to the insurer scheme — is a stub. |
| F-11 | `BordereauxSubmissionTracking.vue` | 942 | `checkStatus()` | Status polling for submitted bordereaux is a stub. |
| F-12 | `BordereauxTemplateManager.vue` | 1147 | `deleteTemplate()` | **Delete is local-only**: the handler splices the row from the local `templates` array and logs. No confirmation dialog; **no API call** — deleted templates reappear on refresh. |

**Quote (F-01, `BordereauxManagement.vue:896`):**
```ts
const generateComplianceReport = () => {
  // TODO: Implement compliance report generation
  console.log('Generating compliance report...')
  complianceDialog.value = false
}
```

### 2.2 Silent error paths (UX gap)

Throughout the module, `catch` blocks log errors and leave a TODO for a snackbar/dialog that was never wired:

| File | TODO count | Lines |
|---|---|---|
| `BordereauxGenerationForm.vue` | 8 | 1127, 1203, 1299, 1303, 1331, 1335, 1363, 1367 |
| `BordereauxReconciliation.vue` | 9 | 1459, 1463, 1642, 1676, 1680, 1809, 2030, 2037, 2041 |
| `BordereauxTemplateManager.vue` | 4 | 1081, 1084, 1134, 1137 |
| `BordereauxSubmissionTracking.vue` | 1 | 925 |
| `BordereauxAnalyticsDashboard.vue` | 2 | 948, 956 |
| `BordereauxManagement.vue` | 1 | 897 |

Aggregate: **~30 paths** where a failure (or a success worth confirming) produces no user-visible feedback. A few of them look like this (`BordereauxGenerationForm.vue:1123–1128`):

```ts
} catch (error) {
  console.error('Job submission failed:', error)
  progressDialog.value = false
  loading.value = false
  // TODO: Show error dialog with proper error message
}
```

Note: some components do use the `flash` service (e.g. `BordereauxSubmissionTracking.vue:951`), so there is a shared toast primitive available — it just hasn't been threaded through consistently.

### 2.3 `BordereauxAnalyticsDashboard.vue` — cosmetic, not functional

Three refs hold hardcoded demo data:
- `benefitPerformance` (line 669) — five fixed benefits.
- `insurerMetrics` (line 713) — five fixed insurer names (Old Mutual, Liberty Life, Momentum, Discovery, Sanlam).
- `monthlyKPIs` (line 814) — six hardcoded months.

There is **no `ag-charts` or `AgCharts` import** — chart slots in the template are placeholder divs. The period selector's `@change="loadAnalytics"` is a no-op (F-02 above). The "Export Report" button is a no-op (F-03).

Backend has `GET /group-pricing/bordereaux/dashboard/stats` (`GetBordereauxDashboardStats`) but that endpoint serves the simple stat cards at the top of `BordereauxManagement.vue`, **not** the analytics charts. The analytics dashboard has no backend.

### 2.4 No shared state / no cache / no real-time progress

- **No Pinia store** for the bordereaux module. Every component fetches independently on mount. A user navigating from `BordereauxManagement.vue` → `BordereauxGenerationForm.vue` and back will refetch scheme lists, template lists, configuration lists.
- **Job progress polling is commented out** in `BordereauxGenerationForm.vue:820–825`:
  ```ts
  // const currentJobId = ref<string | null>(null)
  // const jobStatus = ref<string>('')
  // const jobProgress = ref(0)
  ```
  Long-running generation jobs show a modal spinner but no real progress, no ETA, no partial-completion info. The WebSocket endpoint `/ws` exists on the backend but is not subscribed to for bordereaux events.

### 2.5 Permission checks only at parent

`BordereauxManagement.vue` has three granular checks (`bordereaux:generate_outbound` at line 91, `bordereaux:submit_inbound` at 226, `reinsurance:view` at 278). The permission slugs in `api/installer/gp_permissions.json:394–470` define a richer set (`bordereaux:view`, `bordereaux:process_inbound`, `bordereaux:approve_inbound`, `bordereaux:manage_templates`, `reinsurance:manage_cessions`, `reinsurance:manage_recoveries`, etc.).

**Child components do not check these finer permissions.** Combined with the coarse router-level `navigation:manage_bordereaux` guard applied to all 15 routes (`router/index.ts:189, 197, 205, 213, 221, 229, 237, 245, 254, 262, 270, 278, 286, 294, 302, 310, 318`), a user with only "view" permission can open the Template Manager and attempt to click Create/Update/Delete. The backend is the only guard, and — per existing `gap_analysis.md` GAP-22 — backend `GetActiveUser` currently passes all authenticated requests regardless of role.

### 2.6 Form validation gaps

| Component | Gap |
|---|---|
| `BordereauxGenerationForm.vue` | No cross-field validation — no rule enforcing `start_date < end_date` for custom periods; no rule that template type matches selected bordereaux type; no warning when estimated records = 0. |
| `BordereauxTemplateManager.vue` | Field-mapping dialog allows duplicate `target_field` values; `source_field` can be left blank; no validation that `source_field` belongs to the selected bordereaux type. |
| `BordereauxClaimNotifications.vue` (large-claims monitor) | "Run Monitor" checks only that `monitorTreatyId` is truthy — does not verify the treaty exists, is active, or has linked schemes. |

### 2.7 i18n not applied

The module is hardcoded English throughout — headings, button labels, toast text. The repository has `app/src/renderer/locales/` with 13 languages, but grepping the bordereaux components surfaces no `$t(` usage for the primary labels. The repeated `months` arrays in `RIBordereauxGeneration.vue:405`, `BordereauxGenerationForm.vue:925`, and `BordereauxClaimNotifications.vue:510` also duplicate strings that should come from i18n or a shared constants module.

### 2.8 Missing empty / loading states

Several data grids render blank when empty (no "No data found" row) and show no skeleton/spinner while the initial fetch is in flight:

- `BordereauxReinsurerTracking.vue` (acceptances and recoveries grids)
- `RIBordereauxGeneration.vue` (runs grid)
- `RIClaimsBordereaux.vue` (large-claims and cat-events grids)
- `BordereauxClaimNotifications.vue` (notifications grid)

### 2.9 Destructive actions without confirmation

- `BordereauxTemplateManager.vue`: Activate/deactivate and delete actions do not prompt for confirmation (lines 1044, 1097, 1151 — the delete flow calls the API then logs).
- `BordereauxGenerationForm.vue:1276–1305` (saveConfiguration) creates a new config with no check for duplicate names; updateConfiguration silently overwrites.

### 2.10 Context-menu DOM lifecycle risk

`RIBordereauxGeneration.vue:530–608` — the run-context menu is created/destroyed manually via DOM APIs. If the user navigates away mid-menu, cleanup is not guaranteed (no `onBeforeUnmount` removing the element). Risk: orphan nodes and leaked listeners.

---

## 3. Backend Functional Gaps

### 3.1 Discrepancy escalation is a status flag, not a workflow

**File:** `api/services/bordereaux_reconciliation.go:76–105`

`EscalateDiscrepancy` does this and only this:
```go
result.Status = "escalated"
note := fmt.Sprintf("[%s] Escalated to %s by %s (priority: %s): %s", ...)
result.Comments += note
DB.Save(&result)
_ = writeAudit(...)
return result, nil
```

There is **no notification sent, no email, no task queued, no SLA clock started, no assignee**. "Escalated" becomes just another status string. Frontend displays the status but nothing else happens.

### 3.2 Member bordereaux reconciliation has broken variance logic

**File:** `api/services/bordereaux.go:336–342`

```go
} else if bordereaux.Type == "member" {
    totalExpected = bordereaux.Records
    submittedAmount = 0 // Or some logic if member bordereaux has financial impact
}
```

For member-type bordereaux, `submittedAmount` is hardcoded to zero, which means variance (`totalExpected - submittedAmount`) is meaningless. Any discrepancy count for a member bordereaux is therefore computed against a zero baseline. The comment itself acknowledges the gap. Reconciliation for a member census cannot quantify what changed.

### 3.3 Download endpoint has no authorization

**File:** `api/controllers/bordereaux.go:57–82`

```go
func DownloadBordereaux(c *gin.Context) {
    fileName := c.Param("filename")
    ...
    absReportDir, _ := filepath.Abs(reportDir)
    absFilePath, _ := filepath.Abs(filePath)
    if !strings.HasPrefix(absFilePath, absReportDir) {
        c.JSON(http.StatusForbidden, ...)
        return
    }
    ...
    c.File(absFilePath)
}
```

The only check is path traversal (`HasPrefix`). There is no check that the requesting user:
- owns/has view permission on the scheme(s) embedded in the file,
- is listed as the generating user,
- or has `bordereaux:view` permission.

`api/data/reports/` currently contains **50+ member and premium bordereaux files** (visible in the Glob). Any authenticated user can download any of them by guessing or enumerating filenames (which include predictable Unix timestamps, e.g., `bordereaux_premium_1768008162.zip`).

### 3.4 No retention or cleanup for generated files

**Files:** `api/services/bordereaux.go:1042–1048, 1063–1081` and upload handlers throughout.

Generated bordereaux files are written to `data/reports/` and uploaded confirmations are written under `tmp/uploads/` and `data/bordereaux/inbound/`. There is **no scheduled cleanup, no retention policy, no max-age sweep**. The only `os.Remove` calls found in the bordereaux services are `defer os.Remove(tmpPath)` for parse-time temporaries (`bordereaux_inbound.go:1271, 1275`) and a deletion of the saved confirmation file on DeleteBordereauxConfirmation (`bordereaux.go:696`). Disk growth is unbounded.

### 3.5 Claim recovery auto-generation is orphaned

**File:** `api/services/bordereaux_reinsurer.go` — `GenerateClaimRecovery` is defined but grep shows it is **not called from any claim-lifecycle endpoint**. It can only be invoked indirectly through `CreateReinsurerRecovery` (manual entry).

Consequence: the business flow "claim approved → recovery receivable auto-created against the ceding treaty" requires a human to post the recovery by hand. The `ReinsurerRecovery` model supports auto-generation; the wire-up into the claims assessment approval path is missing. This is also flagged in `gap_analysis.md` GAP-05 ("reinsurance claim recovery … absent").

### 3.6 Async notification goroutines with no error handling

**File:** `api/services/bordereaux_inbound.go:329, 348, 383, 406`

```go
go NotifySubmissionReviewed(sub, user)
go NotifySubmissionQueryRaised(sub, user)
go NotifySubmissionAccepted(sub, user)
go NotifySubmissionRejected(sub, user, reason)
```

Each of these fires a notification in a goroutine whose return value is discarded. If the notification fails (SMTP down, template error, panic), there is no retry, no log, no dead-letter queue. The inbound-submission review/query/accept/reject UX reports success to the user even if nobody was actually notified.

### 3.7 Overdue status computed lazily, not on schedule

**File:** `api/services/bordereaux_deadlines.go` — overdue-status computation runs inside `GetBordereauxDeadlines` rather than a scheduled job. If no user opens the calendar, deadlines remain "pending" in the database past their due date. Any downstream query filtering by `status = 'overdue'` will miss them.

The same lazy pattern applies to `bordereaux_claim_notifications.go:StartNotificationOverdueSweeper` — the sweeper exists but per the backend audit it only **marks** notifications as overdue and does not escalate or alert. Status flips to `overdue` but no email/Slack/task is raised.

### 3.8 Confirmation-to-bordereaux matching by fragile heuristic

**File:** `api/services/bordereaux.go:147–170` (confirmation import path).

When ingesting a scheme confirmation file, the service looks up the target bordereaux by `scheme_id + period`. If multiple generated bordereaux exist for the same scheme and period (which is legitimate — e.g., a re-run after a correction), **the first match wins**. No explicit selection by `generated_id`, no most-recent ordering, no warning to the user about ambiguity. If no bordereaux is found, the confirmation is created with `GeneratedBordereauxID = ""` (line 171 area), producing an orphaned record.

### 3.9 Reconciliation notes stored as synthetic result rows

**File:** `api/services/bordereaux_reconciliation.go:216–241` (note creation).

Free-text notes added via `POST /bordereaux/confirmations/:id/note` are stored by inserting a synthetic `BordereauxReconciliationResult` row with `Field = "_note"`. This forces every consumer of the results table to filter out `_note` rows, and bloats result counts. A proper notes table (or a nullable `notes` JSONB column on the confirmation) would be cleaner and is absent.

### 3.10 Hardcoded business values

| Location | Value | Should be |
|---|---|---|
| `api/services/bordereaux.go` (reconciliation tolerance, in the variance check) | `0.001` | Per-scheme configurable tolerance |
| `api/services/bordereaux_claim_notifications.go:236–241` | Fallback notification SLA = 30 days | Configurable default; currently masks missing treaty data |
| `api/services/bordereaux_inbound.go:60–71` | Filename built from `time.Now().UnixNano()` | Nanosecond uniqueness is a race under concurrent uploads; needs UUID |

### 3.11 State strings are not enums

Bordereaux lifecycle statuses (`draft`, `generated`, `reviewed`, `approved`, `submitted`, `acknowledged`, `escalated`, `validated`, `validation_failed`, `settled`, etc.) are plain strings on the models. There is no enum type, no database constraint, no central list, and no state-transition guard. Any handler can advance a record to any status. The RI validator (`reinsurance_validation.go:20`) does guard one transition — it refuses to validate runs already `submitted|acknowledged|settled` — but that is the exception, not the rule.

### 3.12 Audit logging is inconsistent across bordereaux services

`writeAudit` calls are present in:

| Service | writeAudit calls |
|---|---|
| `bordereaux_reconciliation.go` | 5 |
| `bordereaux_claim_notifications.go` | 2 |
| `bordereaux_reinsurer.go` | 1 |
| `bordereaux.go` | 0 |
| `bordereaux_templates.go` | 0 |
| `bordereaux_inbound.go` | 0 |
| `bordereaux_deadlines.go` | 0 |

Template CRUD, outbound bordereaux lifecycle (review/approve/submit), inbound-submission accept/reject, and deadline status changes are all unaudited. Per regulatory expectations (IFRS 17, FSCA), at least the outbound review/approve/submit decisions and the inbound accept/reject decisions should be logged.

### 3.13 Large-claim notices have no reinsurer response mechanism

**File:** `api/models/reinsurance_bordereaux.go` — `LargeClaimNotice` carries `Status` and `QueryDetails` fields, but there is no endpoint for a reinsurer to respond (accept the notice, query it, reject cession). The route set exposes `MonitorLargeClaims`, `GetLargeClaimNotices`, `UpdateLargeClaimNotice`, `GetLargeClaimStats` — all ceding-side. No inbound response. For schemes where the reinsurer communicates acceptance back, the workflow terminates at "notice sent".

### 3.14 No transaction wrapping in confirmation batch import

**File:** `api/services/bordereaux.go:239–262` (confirmation + delta record batch insertion).

The import uses two sequential `CreateInBatches` calls (confirmation records, then delta/result records). These are not wrapped in a `db.Transaction()`. If the second insert fails (unique constraint, disk full, connection drop mid-batch), the first batch is already committed, leaving orphan confirmation records with no reconciliation result rows. This is also consistent with `gap_analysis.md` GAP-21 ("no transaction safety for composite writes") — the bordereaux module inherits the same class of bug.

---

## 4. End-to-End Flow Gaps

### 4.1 Outbound generation → submission → reconciliation chain

| Stage | Status |
|---|---|
| Generate (backend) | Complete |
| Generate (UI form) | Functional — lacks pre-flight validation and progress feedback |
| Review / Approve (backend) | Complete |
| Review / Approve (UI) | Functional via `BordereauxSubmissionTracking.vue` |
| Submit to scheme (backend) | Endpoint exists: `POST bordereaux/batch-submit` |
| Submit to scheme (UI) | **Stub** (`BordereauxSubmissionTracking.vue:937`) |
| Check submission status (UI) | **Stub** (`BordereauxSubmissionTracking.vue:942`) |
| Confirmation import | Complete but fragile matching heuristic (§3.8) |
| Reconcile pending | Complete |
| Resolve / escalate discrepancy | Resolve: complete. Escalate: status-only (§3.1) |

The end-to-end outbound flow can run if the operator submits manually via the backend API or via backend shell — but the primary UI cannot drive it.

### 4.2 Inbound submission → member register sync chain

This is one of the **more complete** flows. The full pipeline exists and is wired:

1. `CreateEmployerSubmission` → upload file → parse → review → accept/reject.
2. `ComputeSubmissionDelta` → classify rows as `new | amendment | ceased | continuing`.
3. `ApplySubmissionExits` (real impl, `bordereaux_inbound.go:1146`) → deactivates members in the live register.
4. `ApplySubmissionAmendments` (real impl, `bordereaux_inbound.go:1197`) → field-level updates.
5. `SyncNewJoiners` → adds members from staged `NewJoinerDetail` rows.
6. `GenerateScheduleFromSubmission` (real impl, `bordereaux_inbound.go:470`) → produces a premium schedule and cross-links.

Gaps in this flow:
- Inbound notification goroutines fire-and-forget (§3.6).
- No idempotency key on `ApplySubmissionExits`/`Amendments` — running twice will flip the register twice. (Guards on `sub.Status != "accepted"` help, but do not prevent repeat runs of a still-accepted submission.)
- Retro submissions take a separate path (`generateRetroSupplementarySchedule`) — worth verifying in its own audit.
- Inbound audit trail (§3.12) is absent.

### 4.3 RI bordereaux run lifecycle

Member and claims runs: create → validate → submit → acknowledge → amend. The 3-level validator is implemented (`reinsurance_validation.go:15` — L1 structural, L2 integrity, L3 business rules). Status-transition guards exist (`reinsurance_validation.go:20`).

Gaps:
- Settlement is out of scope of this audit but the absence of a finalization step (§3.13) for large-claim notices means the RI workflow has a cession-notification edge that terminates in-air.
- Amendment versioning uses `ParentRunID` and `RunVersion` but there is no endpoint to diff a version against its parent. The amendment audit question "what changed between v2 and v3?" cannot be answered from the API.

### 4.4 Bordereaux ↔ Premium Schedule integration

The link exists (`EmployerSubmission.LinkedPremiumScheduleID` and `PremiumSchedule.LinkedSubmissionID`), and `GenerateScheduleFromSubmission` populates both directions. However:

- Premium lifecycle screens (`app/src/renderer/screens/group_pricing/premiums/*`) are per `gap_analysis.md` GAP-02 still "scaffolded, not functional". A submission that successfully generates a schedule will write the link, but the destination UI is largely backlog.
- `BordereauxReconciliation.vue` does not cross-reference the linked premium schedule. A user reconciling a scheme confirmation cannot click through to the schedule that was generated from that month's submission.

### 4.5 Deadlines ↔ Submissions ↔ Notifications

- `BordereauxDeadline` links to a submission via `LinkedSubmissionID` (model complete).
- Overdue transitions are lazy (§3.7).
- Notifications on overdue are absent — per backend audit, `SweepOverdueNotifications` marks status but does not alert.
- No cross-reference in the calendar UI to jump to the submission/notification when a deadline is clicked. (This is UX/linkage, not a data gap.)

### 4.6 Claims ↔ Reinsurer Recoveries

Broken: the auto-generation hook (§3.5) is never invoked from the claim lifecycle. The `BordereauxReinsurerTracking.vue` UI displays recoveries but they can only appear via manual entry.

---

## 5. Data Model Gaps

| Model | Gap |
|---|---|
| `BordereauxTemplate` | No soft-delete; no version field; no audit-trail field. Template changes over time cannot be reconstructed. |
| `BordereauxConfirmation` | No `last_updated_at`; variance total not persisted on confirmation until `updateConfirmationStats` runs. |
| `BordereauxReconciliationResult` | Notes stored as synthetic `_note` rows (§3.9). No dedicated notes column. |
| `BordereauxDeadline` | No escalation fields (last-reminder-sent, reminder-count, suppressed-until). |
| `ClaimNotificationLog` | No reminder schedule; SLA escalation missing. |
| `LargeClaimNotice` | No response sub-record (reinsurer query, reinsurer acceptance, reinsurer rejection). |
| `RIBordereauxRun` | No link from run to its `RIValidationResult` set other than `RunID` join; no `last_validated_at` stamp beyond status. |
| `EmployerSubmission` | `IsRetro`, `LinkedPremiumScheduleID`, `ExitsSyncedAt`, `AmendmentsSyncedAt`, `NewJoinersSyncedAt` are in place — good granularity. No `idempotency_key` for apply-exits/amendments. |
| All models | Statuses are strings, not enums (§3.11). |

---

## 6. Cross-Cutting / Security / Resilience

| # | Concern | Location |
|---|---|---|
| X-01 | Download endpoint has no per-record authorization | `controllers/bordereaux.go:57` (§3.3) |
| X-02 | `GetActiveUser` middleware does not enforce permissions | `routes/routes.go:25`, inherited from `gap_analysis.md` GAP-22 |
| X-03 | Child components skip fine-grained permission gates (§2.5) | multiple frontend components |
| X-04 | No rate limiting on bordereaux generation or RI validation | `routes/routes.go` |
| X-05 | No transaction wrapping on confirmation batch writes (§3.14) | `services/bordereaux.go:239–262` |
| X-06 | Goroutine fire-and-forget notifications (§3.6) | `services/bordereaux_inbound.go:329,348,383,406` |
| X-07 | Overdue status sweeps are lazy, not scheduled (§3.7) | `services/bordereaux_deadlines.go`, `services/bordereaux_claim_notifications.go` |
| X-08 | File retention / cleanup is absent (§3.4) | `api/data/reports/`, upload dirs |
| X-09 | Audit logging is inconsistent (§3.12) | 4 of 7 bordereaux services have no audit writes |
| X-10 | Analytics dashboard has no backend (§2.3) | `BordereauxAnalyticsDashboard.vue` |

---

## 7. Validated Corrections to Prior Automated Scans

For transparency and to avoid downstream error, these items that were flagged by earlier automated scans are **not** real gaps — the functions exist and are implemented. They were verified by reading the source at the stated locations:

| Alleged gap | Reality | Evidence |
|---|---|---|
| `services.GenerateBordereaux` is missing / stub | Fully implemented | `api/services/bordereaux.go:1002` — dispatches on `type` to member/premium/claim sub-generators, handles per-scheme zipping. |
| `GetBordereauxFieldsByType` has no body | Fully implemented | `api/services/bordereaux.go:841` — returns hardcoded field lists per bordereaux type. |
| `ValidateRIBordereaux` is a stub | Fully implemented | `api/services/reinsurance_validation.go:15` — 3-level pipeline, persists results, updates run status. |
| `GenerateScheduleFromSubmission` is a stub | Fully implemented | `api/services/bordereaux_inbound.go:470` — handles retro path, generates schedule, cross-links. |
| `ApplySubmissionExits` is a stub | Fully implemented | `api/services/bordereaux_inbound.go:1146` — deactivates matching members, records sync summary. |
| `ApplySubmissionAmendments` is a stub | Fully implemented | `api/services/bordereaux_inbound.go:1197` — computes changed fields, updates members. |
| `BordereauxCalendar.vue` / `BordereauxInboundSubmissions.vue` / `BordereauxInboundSubmissionDetail.vue` are missing | All three exist and are functional | referenced in `router/index.ts:236, 244, 252` and file-verified. |

---

## 8. Top-Priority Gaps (Production Blockers)

Ordered by the severity of the functional hole they create:

1. **Download endpoint lacks per-record authorization** (§3.3). Any authenticated user can read any scheme's bordereaux file by URL. Given PII and premium data in these files, this is the first item to fix.
2. **Template delete is local-only — no API call** (F-12). Users see templates "deleted" but they return on refresh. Data-integrity bug masquerading as a working button.
3. **Submit-to-scheme is a UI stub** (F-10). The outbound bordereaux flow cannot be completed through the intended UI path.
4. **Analytics Dashboard is mock data** (§2.3). Users seeing this screen will believe they are looking at live portfolio analytics.
5. **Discrepancy escalation is a status flag with no downstream action** (§3.1). Finance team will not be alerted when discrepancies are escalated.
6. **Claim-recovery auto-generation is orphaned** (§3.5). Ceded-claim receivables are not auto-posted, creating an under-recording risk.
7. **Async notification goroutines swallow failures** (§3.6). Inbound submission review/accept/reject reports success even if nobody was notified.
8. **No file retention policy** (§3.4). Unbounded disk growth; sensitive files never purged.

---

## 9. Prioritized Remediation Roadmap

### P0 — Production Blockers (security, correctness)

| # | Item | Files |
|---|---|---|
| P0-1 | Add per-record authorization to `DownloadBordereaux` (check user has `bordereaux:view` and access to the scheme embedded in the file) | `controllers/bordereaux.go:57` |
| P0-2 | Wrap confirmation batch import in `db.Transaction()` | `services/bordereaux.go:239–262` |
| P0-3 | Wire `GenerateClaimRecovery` into the claim-approval service path | `services/bordereaux_reinsurer.go`, `services/group_pricing.go` claim assessment |
| P0-4 | Implement `submitToScheme()` and `checkStatus()` in `BordereauxSubmissionTracking.vue` | frontend + backend endpoint(s) |
| P0-5 | Replace fire-and-forget notification goroutines with a job queue (or at minimum, error logging + retry) | `services/bordereaux_inbound.go` |

### P1 — Required for Go-Live

| # | Item | Files |
|---|---|---|
| P1-1 | Build real analytics backend for `BordereauxAnalyticsDashboard.vue` (queries over generated bordereaux, insurer performance, compliance metrics); replace mock refs | frontend + new controller/service |
| P1-2 | Expand escalation: add notification hook, SLA clock, assignee, and escalation queue | `services/bordereaux_reconciliation.go:76` |
| P1-3 | Add file-retention policy and scheduled cleanup for `data/reports/` and inbound upload dirs | new background job + config |
| P1-4 | Schedule overdue-deadline and overdue-notification sweeps (cron / periodic worker) | `services/bordereaux_deadlines.go`, `services/bordereaux_claim_notifications.go` |
| P1-5 | Fix member-bordereaux reconciliation variance logic (real comparison, not hardcoded zero) | `services/bordereaux.go:336–342` |
| P1-6 | Thread a toast/snackbar primitive through all ~30 silent error catches (§2.2) | frontend components |
| P1-7 | Enforce permission checks in child components (template-manager, reconciliation, generation form) | frontend components |

### P2 — Completeness

| # | Item | Files |
|---|---|---|
| P2-1 | Implement template preview, test, export, and version history | `BordereauxTemplateManager.vue` (F-04 to F-09); needs backend for preview/test |
| P2-2 | Implement compliance-report generation | `BordereauxManagement.vue:896` |
| P2-3 | Add per-scheme reconciliation tolerance config (replace hardcoded 0.001) | `models/group_pricing.go`, `services/bordereaux.go` |
| P2-4 | Add a proper notes column / table for reconciliation instead of `_note` synthetic rows | `models/bordereaux_inbound.go`/reconciliation, `services/bordereaux_reconciliation.go` |
| P2-5 | Add reinsurer response endpoint for large-claim notices | `controllers/reinsurance_bordereaux.go`, model field additions |
| P2-6 | Add amendment diff endpoint for RI runs | `controllers/reinsurance_bordereaux.go` |
| P2-7 | Convert status strings to enums + transition guards | models + services |

### P3 — Quality

| # | Item | Files |
|---|---|---|
| P3-1 | Add Pinia store for bordereaux state to eliminate redundant fetches | new `app/src/renderer/stores/bordereaux.ts` |
| P3-2 | Wire WebSocket-based job progress for generation and RI validation | frontend + backend WebSocket handler |
| P3-3 | Add confirmation dialogs to destructive template/config actions | frontend |
| P3-4 | Extract months array to shared constant or i18n and apply i18n across the module | frontend |
| P3-5 | Add audit writes to `bordereaux.go`, `bordereaux_templates.go`, `bordereaux_inbound.go`, `bordereaux_deadlines.go` for state transitions | services |
| P3-6 | Cross-field validation on generation form (date range, template/type compatibility) | `BordereauxGenerationForm.vue` |
| P3-7 | Empty-state and loading-state UI on data grids | multiple frontend components |
| P3-8 | Unify filename generation with UUID, not `UnixNano` | `services/bordereaux_inbound.go:60–71` |

---

## 10. Summary Table

| Area | Gaps Identified | Severity mix |
|---|---|---|
| Frontend stubs (no-op handlers) | 11 | 2 production-blocking, 9 feature-completeness |
| Frontend silent errors (missing toasts) | ~30 occurrences | UX polish |
| Frontend architecture (no store, no WS) | 3 | Performance / UX |
| Backend — escalation/notification | 4 | Regulatory/ops risk |
| Backend — authorization/security | 2 | Production blocker |
| Backend — data integrity (variance, orphan records, missing transactions) | 4 | Correctness |
| Backend — audit logging inconsistency | 4 service files unaudited | Regulatory |
| Backend — scheduling/retention | 3 | Ops risk |
| End-to-end flow disconnects | 6 | Mixed |
| Data model gaps | 8 | Feature-completeness |

**Overall assessment:** the bordereaux module is well-structured and has most primitives in place, but the polish-layer work (toast system, state management, scheduled jobs, per-record authorization) and the dashboard/compliance/template-testing features have not been finished. The two items to move on immediately are (a) adding proper authorization to the download endpoint, and (b) wiring or removing the seven UI stubs that look like working features but aren't.

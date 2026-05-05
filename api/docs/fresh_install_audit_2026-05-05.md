# Fresh-Install Migration Audit — 2026-05-05

## Goal

Ensure that on an empty database, `services.SetupTables` creates every table that the active code paths need. Symptoms before the audit: one or two tables missing after a clean install, requiring a manual hand-fix.

## Scope

Every struct in `api/models/*.go` (65 files) was scanned. A struct was treated as a database-backed model if it carried at least one `gorm:"…"` field tag or defined a `TableName() string` method. That heuristic flagged **404 candidate structs**.

Bootstrap reference set: every `models.X{}` literal inside the bodies of the AutoMigrate functions actually invoked from `services.SetupTables` on an empty DB:

- `MigrateBaseTables`
- `MigrateGroupPricingTables`
- `MigrateGroupPricingUserTables`
- `MigrateGroupPremiumTables`
- `MigratePhiValuationTables`
- the standalone `DB.AutoMigrate(&models.SystemLock{})` call in `db.go::SetupTables`

Bootstrap set size: **163 unique structs**.

## Findings

The diff produced **242 candidate structs not in the bootstrap**. The vast majority of these — **233** — are intentional: they belong to the `MigrateProductModelTables`, `MigratePricingTables`, `MigrateEscalationTables`, `MigrateModelPointTables`, `MigrateGMMTables`, `MigrateLICTables`, and `MigrateExposureAnalysisTables` functions, all of which are commented out in `services/migrations.go`. Those modules (legacy IFRS17 / LIC / GMM / exposure analysis) are not part of the AART Group Risk product on a fresh install.

After excluding the legacy bundle and unreferenced orphan structs, **four real gaps** remained.

### Tables missing from the fresh-install bootstrap

| Struct | File | Used by | Added to |
|---|---|---|---|
| `BordereauxConfirmationNote` | `models/bordereaux_template.go` | `services/bordereaux_reconciliation.go` (add/get reconciliation notes against a confirmation) | `MigrateGroupPricingTables` |
| `TableConfiguration` | `models/table_configuration.go` | `services/table_configuration.go`, `services/group_pricing.go` (per-table required/optional flags consumed by group pricing screens) | `MigrateGroupPricingTables` |
| `TableConfigurationAuditLog` | `models/table_configuration.go` | `services/table_configuration.go` (audit history for the above) | `MigrateGroupPricingTables` |
| `RunPhiJob` | `models/utils.go` | `controllers/phi_valuations.go`, `services/phi_valuations.go::RunPhiProjection`; also registered in `struct_migration.go` switch | `MigratePhiValuationTables` |

`TableConfiguration` and `TableConfigurationAuditLog` were not strictly broken on fresh install — `main.go:295` calls `services.EnsureTableConfigurations()`, which itself invokes `DB.AutoMigrate(&models.TableConfiguration{}, &models.TableConfigurationAuditLog{})`. They are now also listed in `MigrateGroupPricingTables` so that all tables go through one bootstrap pass and are recorded by `MarkAllMigrationsAsApplied` in the same step. AutoMigrate is idempotent, so the second call from `EnsureTableConfigurations` is a no-op.

`BordereauxConfirmationNote` and `RunPhiJob` were the actual fresh-install gaps — neither had any other migration path, so a clean DB would never have those tables until the first time a piece of code tried to query them.

### Structs deliberately excluded from the patch

| Struct | File | Reason |
|---|---|---|
| `FamilyFuneral` | `models/group_pricing.go` | Defined as `type FamilyFuneral struct { … gorm:… }` but no other file references it. Orphan model, presumably superseded by the inline `FamilyFuneralBenefit`/`FamilyFuneralMainMemberFuneralSumAssured`/etc. fields on `GroupScheme`. Adding it would create an empty unused table. |
| `MemberAddress` | `models/group_pricing.go` | Same pattern: defined once, referenced nowhere. |
| `RunJob`, `RunParameters`, `ShockSetting` | `models/utils.go` | Used only by the legacy projection engine (`services/projections.go`, `services/valuation_projection_engine.go`, `services/products.go`, `services/assumptions_vars.go`). All callers are part of the commented-out `MigrateProductModelTables` / `MigrateModelPointTables` paths and are not exercised on a Group Risk fresh install. |

### Legacy bundle (intentionally not migrated)

233 structs from the following files belong to commented-out Migrate* functions and are therefore not part of fresh install:

`assumption_tables.go`, `basevariables.go`, `basemortalityband.go`, `bel_buildup.go`, `csm_engine.go`, `deferred_tax.go`, `exp_analysis.go`, `global_tables.go`, `ibnr_engine.go`, `ifrs17.go`, `ifrs17_amendment.go`, `ifrs17_audit_log.go`, `lapserate.go`, `lic.go`, `markov_state.go`, `modelpoint.go`, `modelpointvariable.go`, `modified_gmm_engine.go`, `paa.go`, `pricing.go`, `pricingtables.go`, `product.go`, `productfeatures.go`, `productmargins.go`, `productparameters.go`, `projection.go`, `ratingfactor.go`, `reports.go`, `sarb_code_mapping.go`, `scr_ra_bridge.go`, `task.go`, `transition_adjustment.go`, `transitionstate.go`, `winners.go`, `yieldcurve.go`.

If any of these modules ever needs to be re-enabled for the Group Risk product, the corresponding commented-out `Migrate*` function in `services/migrations.go` should be uncommented and its call re-added to `services/db.go::SetupTables`. The audit script used here can then be re-run to confirm coverage.

## Patch summary

`api/services/migrations.go`:

```go
// MigrateGroupPricingTables — added:
&models.BordereauxConfirmationNote{}    // next to BordereauxConfirmation
&models.TableConfiguration{}            // at the end, with comment
&models.TableConfigurationAuditLog{}    // at the end, with comment

// MigratePhiValuationTables — added:
&models.RunPhiJob{}                     // after PhiRunConfig, with comment
```

No other files were touched. The change is purely additive: existing databases that already hold these tables will see no-op AutoMigrate calls.

## Verification

- All four added symbols resolve to declared exported types under `api/models/`:
  - `BordereauxConfirmationNote` — `models/bordereaux_template.go:162`
  - `TableConfiguration` — `models/table_configuration.go:16`
  - `TableConfigurationAuditLog` — `models/table_configuration.go:29`
  - `RunPhiJob` — `models/utils.go:113`
- A `go build ./...` was attempted in this session but Go is not installed in the analysis sandbox; please run it once before merging:
  ```bash
  cd api && go build ./...
  ```
- After merging, the next fresh install should land all 167 tables (163 prior + 4 new) on the first boot and `MarkAllMigrationsAsApplied` will record every file in `migrations/<dialect>/` as baselined.

## Re-running the audit

The diff was produced by extracting all gorm-tagged structs in `api/models/` and intersecting with `models.X{}` references inside the live `Migrate*` function bodies. To re-run after future model additions:

```bash
cd api
python3 - <<'PY'
import os, re
struct_re = re.compile(r'^type\s+([A-Z]\w+)\s+struct\s*\{', re.MULTILINE)
tablename_re = re.compile(
    r'^func\s+\(\s*(?:\w+\s+)?\*?([A-Z]\w+)\s*\)\s*TableName\s*\(\s*\)\s*string',
    re.MULTILINE,
)
# ... see audit transcript for full script
PY
```

The full script is preserved in the conversation that produced this report. Anything new under `api/models/` that has a `gorm:` tag and isn't referenced inside an active `Migrate*` body will surface in the next run.

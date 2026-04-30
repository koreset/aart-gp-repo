# Prompt: Diff-based Migration System Refactor

Paste the following prompt into a Claude Code session running in the target project (the valuations codebase, or any sibling project with the same migration shape).

The prompt assumes a Go API using GORM with a `services/` layer and a `tools/` directory containing CLI utilities, a `migrations/` directory with dialect subfolders (`postgresql/`, `mysql/`, `mssql/`), and a `services.DbBackend` string variable derived from a `db_type` config field. If anything in the project doesn't match those assumptions, ask before changing it.

---

## Prompt to paste

I want you to refactor the database migration system in this project so it stops emitting full-struct SQL dumps and instead generates diffs against the live database. The same approach worked cleanly in a sibling project (AART Group Risk). I'll describe the goal, the architecture, the specific changes, the edge cases you'll hit, and the doc updates. Read this whole prompt first, then walk through the steps in order. Don't skip the pre-flight check.

### Goal

The current `tools/generate_migration.go` (and its `services.GenerateMigrationForStruct` helper) produces SQL that re-declares every column on a struct on every run. This causes:

- Foreign key incompatibilities (placeholder `id INT` vs. real `BIGINT` columns)
- Index churn from redundant `ALTER COLUMN TYPE` statements
- Bloated, hard-to-review migration files
- No real "incremental" semantics — files describe desired state, not changes

We want to replace that with a diff-based generator: read the struct via `gorm.io/gorm/schema.Parse`, introspect the live database via `DB.Migrator()`, and emit only the deltas. The runner that applies migration files (typically `services.RunMigrationsOnStartup`) doesn't need to change architecturally — it already reads files from `migrations/<dialect>/` and tracks applied versions in a `migrations` table.

### Pre-flight check

Before touching anything, confirm the project has the following pieces. If any are missing or different, stop and ask.

1. `services/struct_migration.go` exists with a `GenerateMigrationForStruct` function and a giant `getStructType` switch mapping struct names to `reflect.Type`.
2. `services/migrations.go` exists with `MigrateXTables` functions calling `DB.AutoMigrate(...)` for bootstrap.
3. `services/migration.go` exists with a `MigrationManager`, a `Migration` model, a `RunMigrationsOnStartup` function, a `splitSQLStatements` helper, and a `migrations` directory layout.
4. `services/db.go` has a `SetupTables` function with a "fresh install" branch that calls the AutoMigrate functions when the DB is empty.
5. `services.DbBackend` is set from config and takes values `"postgresql"`, `"mysql"`, or `"mssql"`.
6. There's a `tools/generate_migration.go` and likely a `tools/run_migrations.go` using the `//go:build ignore` convention.
7. The API's `main.go` calls `services.RunMigrationsOnStartup()` during boot.

If those pieces match, proceed. If the project has different file/function names, map them to equivalents and proceed. If something is fundamentally different, ask.

### Architecture target

After the refactor, the system has three roles:

- **Generation**: `services.GenerateMigrationForStructs(names []string, opts GenerateOptions) (string, error)` reads struct definitions, parses them with `schema.Parse`, introspects the live DB via `DB.Migrator()`, and writes a single SQL file with only the deltas to `migrations/<dialect>/<timestamp>_<message>.sql`. Auto-detects dialect from `services.DbBackend`. Never applies SQL. Defaults to additive-only (no type changes, no drops).

- **Application**: `RunMigrationsOnStartup` (already exists) reads files from `migrations/<dialect>/`, compares to the `migrations` table, applies pending files in version order in transactions, records each on success.

- **Baseline**: a one-shot `tools/baseline_migrations.go` that walks every file in `migrations/<dialect>/` and inserts versions not yet in the `migrations` table without executing the SQL. Used for environments whose schema is already at HEAD (typically because earlier AutoMigrate runs put it there). Also called automatically from the fresh-install branch in `SetupTables` so a brand-new DB doesn't try to re-apply baseline files.

### Implementation steps

#### Step 1 — Refactor `services/struct_migration.go`

Replace the existing `GenerateMigrationForStruct` and the entire SQL-generation half (everything below `getStructType`) with a diff-based implementation. **Keep `getStructType` exactly as it is** — it's a giant switch that maps struct names to `reflect.Type` and works fine.

Add a `GenerateOptions` struct:

```go
type GenerateOptions struct {
    AllowTypeChanges bool   // off by default — comparison is fuzzy
    AllowDestructive bool   // off by default — data loss
    Message          string // used in the generated filename
}
```

Replace `GenerateMigrationForStruct` with a thin wrapper plus `GenerateMigrationForStructs([]string, GenerateOptions)` that:

1. Validates DB connection and `DbBackend` are set.
2. For each struct name, resolves the type via `getStructType`, calls `generateIncrementalSQLForModel`, accumulates non-empty bodies.
3. If anything was generated, writes a single file `migrations/<DbBackend>/<timestamp>_<message>.sql` and returns its path.
4. If no changes, returns `("", nil)`.

The core diff function `generateIncrementalSQLForModel(model interface{}, opts GenerateOptions)`:

1. Parses expected schema with `schema.Parse(model, &sync.Map{}, DB.NamingStrategy)`.
2. If `!DB.Migrator().HasTable(model)`, render a full `CREATE TABLE` plus indexes and return.
3. Otherwise, fetch live columns via `DB.Migrator().ColumnTypes(model)`, build a name → ColumnType map.
4. Pass 1: for each `*schema.Field` in expected, if missing from live → emit `ADD COLUMN`. If present and `AllowTypeChanges` is true and types differ → emit `ALTER COLUMN TYPE`.
5. Pass 2 (gated on `AllowDestructive`): live columns missing from expected → `DROP COLUMN`.
6. Pass 3: walk `s.ParseIndexes()` (which returns `[]*schema.Index` in modern GORM versions), emit `CREATE INDEX` for any whose name isn't yet in `DB.Migrator().HasIndex`. With `AllowDestructive`, also emit `DROP INDEX` for live indexes not in the expected set, skipping primary-key indexes (`*_pkey`, `PRIMARY`).
7. Return the buffered SQL, or `""` if no changes.

Helper functions you'll need:

- `parseSchema(model interface{}) (*schema.Schema, error)` — wraps `schema.Parse`.
- `isMigratableField(*schema.Field) bool` — filters out fields with empty `DBName` or `IgnoreMigration`.
- `writeMigrationFile(message, body string) (string, error)` — writes to `migrations/<DbBackend>/`.
- `sanitizeFilenameSegment(string) string` — turns the message into a filesystem-safe slug.
- `renderCreateTable(*schema.Schema, dialect string) string` — full CREATE TABLE with PK + indexes.
- `renderAddColumn(table string, *schema.Field, dialect)` — `ALTER TABLE ... ADD COLUMN`.
- `renderAlterColumnType(table, *schema.Field, gorm.ColumnType, dialect)` — `ALTER COLUMN TYPE`. Use a `normalizeType` helper that lowercases and collapses whitespace, and maps `character varying` → `varchar`, `double precision` → `double`.
- `renderDropColumn`, `renderDropIndex` — for the destructive pass. Note MySQL's DROP INDEX requires the table reference; PG and MSSQL don't.
- `renderCreateIndex(table string, idx *schema.Index)` — IMPORTANT: `idx` must be `*schema.Index`, not `schema.Index`. Modern GORM returns `[]*schema.Index` from `ParseIndexes()`. Iterate with `for _, idx := range expectedIdx { if idx.Class == "PRIMARY" continue; ... }`.
- `hasNonEmptyDefault(*schema.Field) bool` — checks `HasDefaultValue && DefaultValue != "" && DefaultValue != "(-)"`. GORM uses `(-)` as a sentinel for "explicitly no default."
- `resolveSQLType(*schema.Field, dialect string) string` — honors `f.TagSettings["TYPE"]` for explicit overrides, otherwise maps `f.GORMDataType` to dialect SQL types. Cover `schema.String`, `Int`, `Uint`, `Float`, `Bool`, `Time`, `Bytes` for `"postgresql"`, `"mysql"`, `"mssql"`. Use `f.Size`, `f.Precision`, `f.Scale`, `f.AutoIncrement` to render correct types (e.g. `BIGSERIAL` vs `SERIAL` vs `BIGINT`).

NOT NULL handling for new columns: only emit `NOT NULL` when `hasNonEmptyDefault` is also true; otherwise emit a comment `/* NOT NULL omitted: declare a default or backfill manually before tightening */`. Adding NOT NULL to an existing populated table without a default fails on most dialects.

Use `appLog "api/log"` for logging (or whatever the project's logger is).

Imports needed: `appLog "api/log"`, `"api/models"`, `"fmt"`, `"os"`, `"path/filepath"`, `"reflect"`, `"strings"`, `"sync"`, `"time"`, `"gorm.io/gorm"`, `"gorm.io/gorm/schema"`.

#### Step 2 — Refactor `tools/generate_migration.go`

CLI matching the new generator. Flags: `-struct=Name1,Name2`, `-message=text`, `-allow-type-changes`, `-allow-destructive`, `-config=path`. Drop the old `-db=` flag.

Important: this CLI needs to load `config.json` before calling `services.SetupTables(false, false)`. The existing tools don't do this and will panic with nil pointer. Mirror what `main.go::startApplication` does: read `./config.json` (or the path from `-config`), unmarshal into `globals.AppConfig`, then call `SetupTables`.

Inline a `loadConfig(override string) error` helper at the bottom of the file. On non-Windows it reads `./config.json`; on Windows it reads alongside the executable.

#### Step 3 — Create `tools/baseline_migrations.go`

New `//go:build ignore` file. Same config-loading pattern as Step 2. Calls `services.MarkAllMigrationsAsApplied`. Accepts `-config=path` for pointing at remote DBs (with SSH tunneling, etc.).

#### Step 4 — Add `MarkAllMigrationsAsApplied` to `services/migration.go`

```go
func MarkAllMigrationsAsApplied() error {
    if DB == nil { return fmt.Errorf("database not initialized") }
    mgr := NewMigrationManager(DB)
    if err := mgr.Initialize(); err != nil { return err }
    files, err := mgr.getMigrationFiles()
    if err != nil { return err }
    now := time.Now()
    for _, file := range files {
        version := strings.Split(filepath.Base(file), "_")[0]
        var existing Migration
        if err := DB.Where("version = ?", version).First(&existing).Error; err == nil {
            continue
        }
        m := Migration{Version: version, Name: strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)), AppliedAt: now}
        if err := DB.Create(&m).Error; err != nil {
            return fmt.Errorf("record baseline migration %s: %w", version, err)
        }
    }
    return nil
}
```

Idempotent — versions already present are skipped.

#### Step 5 — Wire baseline into the fresh-install bootstrap in `services/db.go`

In `SetupTables`, after the AutoMigrate functions run on an empty database (the `if !dbHasTables` branch), call `MarkAllMigrationsAsApplied()`. This ensures a brand-new install records every existing `.sql` file as already applied, so the runner doesn't try to re-execute them.

```go
if err := MarkAllMigrationsAsApplied(); err != nil {
    appLog.WithField("error", err.Error()).Error("Failed to record baseline migrations")
}
```

#### Step 6 — Make `RunMigrationsOnStartup` fatal on failure

The existing function logs migration errors and returns, letting the API start against a stale schema. Change it to `appLog.Fatal(...)` so any error or nil DB causes the boot to abort with `os.Exit(1)`. Failure modes are: stale schema served to traffic > controlled boot crash.

```go
func RunMigrationsOnStartup() {
    appLog.Info("Running migrations on startup")
    if DB == nil {
        appLog.Fatal("Database not initialized — cannot run migrations; aborting startup")
    }
    mgr := NewMigrationManager(DB)
    if err := mgr.RunMigrations(); err != nil {
        appLog.WithField("error", err.Error()).Fatal("Failed to run migrations — aborting startup")
    }
    appLog.Info("Migrations completed successfully")
}
```

The `appLog.Fatal` from logrus calls `os.Exit(1)` after logging.

#### Step 7 — Teach `splitSQLStatements` about MySQL `DELIMITER`

Hand-written MySQL migrations defining stored procedures, triggers, or functions use `DELIMITER $$ ... DELIMITER ;` to wrap bodies that contain semicolons. `DELIMITER` is a `mysql` CLI client convention, not a SQL keyword — sending it to the server returns error 1064.

Update `splitSQLStatements` to:

1. Track an `atLineStart bool` and a `delim string` (defaulting to `";"`).
2. At the start of each line (skipping leading whitespace), check for the keyword `DELIMITER` (case-insensitive) followed by whitespace. If found, parse the new delimiter token from the rest of the line, flush any in-progress statement, change `delim`, skip the entire directive line. Don't emit it to the output.
3. Use `delim` (multi-character safe) instead of literal `;` for statement splitting.
4. Track line boundaries so DELIMITER detection only fires at the start of a line.
5. Preserve existing behavior: line comments (`--` and `#`), single-quote and double-quote string literal awareness.

Files that don't use DELIMITER continue to behave identically.

#### Step 8 — Add `isIdempotentConflict` and use it in the apply loop

Schema conflicts that mean "already in the desired state" should be soft warnings, not fatal errors. This handles the cross-environment case where a column exists on one DB (from AutoMigrate) but not another (got it via a migration file).

Add to `services/migration.go`:

```go
func isIdempotentConflict(err error) bool {
    if err == nil { return false }
    msg := strings.ToLower(err.Error())
    patterns := []string{
        "duplicate column",         // MySQL 1060
        "duplicate key name",       // MySQL 1061 (index)
        "duplicate index",          // misc MySQL
        "table already exists",     // MySQL 1050
        "already exists",           // PostgreSQL
        "already an object named",  // MSSQL 2714
    }
    for _, p := range patterns {
        if strings.Contains(msg, p) { return true }
    }
    return false
}
```

Modify the inner loop of `MigrationManager.RunMigrations` so when `tx.Exec(stmt).Error` returns an idempotent conflict, log a warning and `continue` instead of returning the error. Other errors stay fatal.

```go
if execErr := tx.Exec(stmt).Error; execErr != nil {
    if isIdempotentConflict(execErr) {
        appLog.WithFields(map[string]interface{}{
            "warn": execErr.Error(),
            "stmt": stmt,
        }).Warn("Migration statement skipped — object already exists")
        continue
    }
    appLog.WithFields(map[string]interface{}{
        "error": execErr.Error(),
        "stmt":  stmt,
    }).Error("Failed to execute migration statement")
    return execErr
}
```

#### Step 9 — Documentation

Rewrite `MIGRATIONS.md` (or create one if it doesn't exist) covering:

- What changed (diff-based vs. full-struct dump)
- Architecture diagram (struct → schema.Parse diff → dialect SQL → runner → DB)
- Day-to-day dev workflow (pull → restart → edit struct → generate → restart → commit)
- Tools reference (generate_migration, baseline_migrations, run_migrations)
- Migration file layout
- Bootstrap vs. incremental paths
- The `migrations` tracking table
- Remote/production usage with `-config` and SSH tunneling
- Cutover guidance (when to baseline, how to handle legacy files)
- Safety notes (additive-only default, NOT NULL caveat, type-comparison fuzziness, fatal-on-failure runner)
- Troubleshooting section covering: nil pointer panic from missing config, MySQL 1064 near DELIMITER, "Migration statement skipped" warnings, "no schema changes detected" trap (local DB ahead of generator-target DB), "API exits on startup" with how to recover.
- Files of interest map.

Trim the in-directory `migrations/README.md` to a short pointer with the dialect layout and a link back to MIGRATIONS.md.

Update `CLAUDE.md` (if present) to replace any one-line migrations bullet with a paragraph describing the diff-based approach and pointing at MIGRATIONS.md.

### Cutover for existing environments

Each environment whose database already has the schema (typically because earlier AutoMigrate runs created it) needs its `migrations` table populated with every existing `.sql` file. Run once per environment:

```bash
cd api
go run tools/baseline_migrations.go -config=config.<env>.json
```

For remote environments, open an SSH tunnel first and point the config at the tunneled address.

If a database is *not* at HEAD, don't baseline wholesale — inspect each file, decide whether the changes are needed, apply the ones that are, and only baseline the ones that the DB already reflects.

### Errors you'll likely see during cutover

These are all real things that came up during the AART cutover. Be ready for them.

1. **Nil pointer panic from `tools/baseline_migrations.go` or `tools/generate_migration.go`**: config not loaded. The CLI must load `config.json` before calling `SetupTables`. See Step 2.

2. **MySQL error 3780 "incompatible columns" during legacy migration apply**: pre-refactor migration tries to add a FK with a hardcoded INT placeholder, but the actual column is BIGINT. The schema is usually already correct (FK in place) and the migration is redundant. Manual fix: verify with `SHOW CREATE TABLE` that both columns are bigint and the FK exists, then `INSERT INTO migrations (version, name, applied_at) VALUES ('<version>', '<name>_legacy_skipped', NOW())` to mark the broken migration as applied.

3. **MySQL error 1170 "BLOB/TEXT column used in key spec without key length"**: the live column is TEXT but the struct says `size:N`. Convert with `ALTER TABLE x MODIFY COLUMN col VARCHAR(N)` after verifying `MAX(LENGTH(col)) <= N`. Then add the index manually and skip the migration.

4. **MySQL error 1064 near `DELIMITER $$`**: splitter doesn't understand DELIMITER. Step 7 fixes this. If the splitter's already updated, the file probably has a missing `DELIMITER ;` reset.

5. **MySQL error 1062 "Duplicate entry" on the `migrations` table INSERT**: two files share a version prefix. The runner applied the first, then tried to apply the second and the INSERT collided. Find the duplicate (`ls migrations/<dialect>/<prefix>*`), rename one to a unique timestamp, or delete if redundant.

6. **MySQL error 1060 "Duplicate column"**: column was added by AutoMigrate on this DB, but a migration file also tries to add it. The `isIdempotentConflict` tolerance from Step 8 handles this — the runner logs a warning and continues.

If the legacy queue contains many FK-mismatch errors (#2 above) and you don't want to skip them one at a time, you can temporarily add `"are incompatible"` to the `isIdempotentConflict` patterns list with a clearly marked TEMPORARY comment indicating it should be removed once every environment's `migrations` table catches up past the legacy date range. This auto-skips the entire legacy FK queue. The trade-off: a real type-mismatch in a future hand-written migration would also get logged-and-skipped instead of erroring. Keep the change clearly marked so it's easy to remove later.

### Verification

After each step:

- `go build ./...` from the project root must succeed. The schema package APIs in modern GORM versions (v1.20+) are stable; if you get errors about `*schema.Index` vs `schema.Index`, you've hit the slice-vs-map ParseIndexes return-type change — use pointer.
- After Step 1, manually run the generator against a struct you know has no diff vs. live DB. It should print "No schema changes detected" and exit 0 without writing a file.
- After Step 2, the same thing via the CLI: `go run tools/generate_migration.go -struct=<KnownStable> -config=<path>`.
- After Step 3, run baseline against a copy of dev DB. Verify `SELECT count(*) FROM migrations` matches the file count under `migrations/<dialect>/`.

### Things NOT to change

- The `getStructType` switch in `services/struct_migration.go`. Hundreds of lines, works fine. Leave it.
- The `MigrateXTables` AutoMigrate functions in `services/migrations.go`. They become bootstrap-only. Don't edit them as part of this refactor.
- The `Migration` model. Keep its existing schema.
- The `MigrationManager.RunMigrations` outer flow (file listing, version comparison, transaction wrapping). Only the inner statement-execution loop needs the `isIdempotentConflict` integration.

### Order to apply

Steps 1, 2, 3 can ship in one PR (the new generator + baseline tool). Steps 4 (`MarkAllMigrationsAsApplied`) is needed for Step 3 to work. Step 5 (wiring into bootstrap) is independent. Steps 6, 7, 8 can be in a second PR (runner robustness). Step 9 (docs) goes with whichever PR introduces the feature being documented.

For verification, do Steps 1-5 first against a dev DB, run baseline, run a real generated migration through to confirm the loop works end-to-end, then add the runner robustness.

If you hit anything not covered here, ask before improvising. The refactor's goal is to make migration files small, dialect-aware, and forward-only — anything that violates that should be flagged.

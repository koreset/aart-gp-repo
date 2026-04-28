# Database Migrations

This document describes the database migration system used by the API: how it works, the day-to-day workflow for adding schema changes, and the operational details for deploying changes to remote environments.

## What changed (April 2026)

The migration system is now **diff-based**. Earlier, the generator emitted SQL containing every column from a struct on every run, regardless of what was already in the database. That produced noisy migration files, churned indexes on type re-declaration, and produced FK type mismatches when placeholder primary keys disagreed with what GORM created elsewhere.

The new generator compares the struct's expected schema (via `gorm.io/gorm/schema.Parse`) against the live database (via `DB.Migrator()`) and emits only the deltas: ADD COLUMN for new fields, CREATE INDEX for new indexes, CREATE TABLE for entirely new tables. Type changes and drops are **off by default** and have to be opted into with explicit flags.

Each environment's database tracks which migration files it has applied in a `migrations` table, so dev, staging, and prod can be at different points in time without conflict — they all converge by applying their pending files on the next deploy.

## Architecture

```
┌──────────────┐   schema.Parse   ┌──────────────────┐
│   Go struct  │ ───────────────▶ │ expected schema  │
│  (api/models)│                  └────────┬─────────┘
└──────────────┘                           │
                                           │  diff
                                           ▼
┌──────────────┐  Migrator API   ┌──────────────────┐
│  Live DB     │ ───────────────▶│ live schema      │ ───▶ delta SQL
└──────────────┘                 └──────────────────┘
                                           │
                                           ▼
                              api/migrations/<dialect>/<timestamp>_<msg>.sql
                                           │
                                           ▼
                              RunMigrationsOnStartup applies pending files
                                           │
                                           ▼
                              records version in migrations table
```

Two halves of the system:

- **Generation** (`services.GenerateMigrationForStructs`, invoked via `tools/generate_migration.go`) produces SQL files. Never applies them.
- **Application** (`services.MigrationManager.RunMigrations`, invoked from `services.RunMigrationsOnStartup` on every API boot) reads files, compares to the `migrations` table, applies what's pending in version order, each in its own transaction, recording success.

The two halves never run in the same process. You always see the SQL before the runner executes it.

## Day-to-day workflow

1. `git pull` — picks up new migration files committed by teammates.
2. Restart the API (or let `air` do it). `RunMigrationsOnStartup` applies any pending files to your local DB.
3. Edit a struct in `api/models/`. Add a field, an index, whatever.
4. Generate the migration:
   ```bash
   cd api
   go run tools/generate_migration.go -struct=Beneficiary -message=add_nickname
   ```
   This writes `migrations/<dialect>/<timestamp>_add_nickname.sql` containing only the diff.
5. Review the generated SQL. Edit by hand if anything looks off (e.g. NOT NULL omitted because no default was declared — see "NOT NULL caveat" below).
6. Restart the API. The runner applies the new file and records the version.
7. Commit the struct change and the SQL file in the same PR.

**Important:** the generator compares the struct to whatever DB it's connected to. If your local DB has been auto-migrated past the previous schema (e.g. from prior boots with the old AutoMigrate path), the diff will be empty — there's nothing for the generator to find. Step 2 is what keeps the local DB at "everyone's HEAD" before step 4.

## Tools

All tools live under `api/tools/` and use the `//go:build ignore` convention so they don't get compiled into the API binary. Run them with `go run` from the `api/` directory, where `config.json` lives.

### `generate_migration.go` — produce a delta migration

```bash
go run tools/generate_migration.go -struct=Beneficiary -message=add_nickname
go run tools/generate_migration.go -struct=Beneficiary,GroupScheme -message=q2_release
go run tools/generate_migration.go -struct=Beneficiary -allow-type-changes
go run tools/generate_migration.go -struct=Beneficiary -config=config.staging.json
```

Flags:

- `-struct=Name1[,Name2,...]` — required. Comma-separated list of model struct names. The list is resolved via the switch in `services/struct_migration.go::getStructType`; if a model isn't in that switch, add a case there.
- `-message=text` — used in the filename. Defaults to `update_<struct>` if omitted.
- `-allow-type-changes` — emits ALTER COLUMN TYPE when the struct's type differs from the live column. Off by default; comparison across dialects is fuzzy and easily produces spurious diffs.
- `-allow-destructive` — emits DROP COLUMN and DROP INDEX for objects in the DB but not in the struct. Off by default; data loss is irreversible.
- `-config=path` — load a different config file (e.g. point at a remote DB). Defaults to `./config.json`.

Output goes to `migrations/<dialect>/<timestamp>_<message>.sql`. Dialect is auto-detected from the connected database via `services.DbBackend`.

### `baseline_migrations.go` — record on-disk files as already applied

Use this once per environment when adopting the new system, or whenever a database already contains the schema described by some `.sql` files (typically because earlier `AutoMigrate` runs created the tables).

```bash
go run tools/baseline_migrations.go
go run tools/baseline_migrations.go -config=config.staging.json
```

Walks every file in `migrations/<dialect>/`, inserts any version not yet recorded into the `migrations` table with `applied_at = now()`. Does **not** execute any SQL. Safe to run repeatedly — files already recorded are skipped.

### `run_migrations.go` — manual apply, mostly legacy

```bash
go run tools/run_migrations.go -type=all
```

Most of this tool's flags (`-type=base`, `-type=pricing`, etc.) drive the **bootstrap** AutoMigrate functions in `services/migrations.go`, which only run on a brand-new database. For day-to-day schema changes, prefer letting `RunMigrationsOnStartup` handle file application on boot.

## Migration file layout

```
api/migrations/
├── README.md
├── postgresql/
│   ├── 20260425_153012_add_nickname.sql
│   └── 20260426_092200_add_broker_index.sql
├── mysql/
│   └── 20260425_153012_add_nickname.sql
└── mssql/
    └── 20260425_153012_add_nickname.sql
```

Files are named `<timestamp>_<message>.sql`. The timestamp prefix determines apply order. Files in a dialect subdirectory only apply when the connected DB matches that dialect (driven by `services.DbBackend`, which is read from `config.json`'s `db_type`).

The generator only writes to one dialect per run — whichever the connected DB is. If you need the same change for multiple dialects, run the generator against each dialect's shadow DB.

## Bootstrap vs. incremental

`services.SetupTables` (called from `main.go` during startup) decides between two paths based on whether the database has any tables:

**Empty database (first install):**
1. Runs every `AutoMigrate`-based function in `services/migrations.go` (`MigrateBaseTables`, `MigrateGroupPricingTables`, etc.) to create all tables from current struct definitions.
2. Calls `MarkAllMigrationsAsApplied` to insert every file in `migrations/<dialect>/` into the `migrations` table as already applied — preventing the runner from trying to re-apply baseline SQL to a fresh DB.

**Existing database:**
1. Skips AutoMigrate.
2. `RunMigrationsOnStartup` runs and applies any files whose versions aren't yet in the `migrations` table.

Treat the AutoMigrate functions in `services/migrations.go` as **bootstrap-only**. Never edit them as part of a normal feature change — the change won't reach existing databases.

## Migration tracking table

The `migrations` table is created automatically on first boot by `MigrationManager.Initialize`. Schema:

| Column     | Type         | Notes                                   |
|------------|--------------|-----------------------------------------|
| id         | INT/SERIAL   | Primary key                             |
| version    | VARCHAR(255) | Unique. The timestamp prefix of a file. |
| name       | VARCHAR      | Filename without extension              |
| applied_at | TIMESTAMP    | When the file was applied or baselined  |

Manually insert a row with `INSERT INTO migrations (version, name, applied_at) VALUES (...)` if you ever need to skip a specific file (e.g. one that was applied out of band).

## Remote and production environments

The tools accept `-config=<path>` so you can point them at a database other than your local one. Concrete recipe for staging or prod:

1. **Create a config file for the environment:**
   ```bash
   cp config.json config.staging.json
   # edit config.staging.json — set db_host, db_user, db_pwd, db_name, db_type
   ```
   Add `config.*.json` to `.gitignore` so credentials never get committed.

2. **Make sure you can reach the DB.** For private prod databases, open an SSH tunnel:
   ```bash
   ssh -L 13306:internal-mysql.example.com:3306 bastion.example.com
   ```
   In the config, set `db_host=127.0.0.1` and `db_port=13306` so the tool connects through the tunnel.

3. **Run the tool:**
   ```bash
   go run tools/baseline_migrations.go -config=config.staging.json
   go run tools/generate_migration.go -struct=Beneficiary -config=config.staging.json
   ```

For automated production deploys, `RunMigrationsOnStartup` runs on every API boot — there's nothing extra to wire up, as long as the new SQL files are in the deployed image.

## Cutover from the old system

If you have legacy `.sql` files generated by the pre-refactor full-struct generator, the runner will try to apply them on first boot of the new code and may fail (placeholder INT primary keys collide with BIGINT FKs, type re-declarations clash with existing indexes, etc.).

To handle the cutover for an environment whose database is already at HEAD:

1. Deploy the refactored code with the new tools available.
2. Run `go run tools/baseline_migrations.go -config=<env>.json` against that environment **before** the next API boot — or accept that the runner will log errors on the first boot and won't apply legacy files until you baseline.
3. Verify with `SELECT count(*) FROM migrations` — should match the number of files in `migrations/<dialect>/`.
4. Restart the API. `RunMigrationsOnStartup` should now log "All migrations applied successfully" with no errors.

If a database is **not** at HEAD (e.g. an older prod that's missing changes from some legacy migration files), don't baseline wholesale. Inspect each file, decide whether the changes are needed, apply the ones that are, and only baseline the ones that the DB already reflects.

## Safety considerations

**Generator (read-only):** safe to run against any database. It only reads. Empty diffs cost nothing.

**Baseline tool:** writes to the `migrations` table only. Never executes migration SQL. Safe to run repeatedly. The risk is *logical*: marking a file as applied when the DB doesn't actually have those changes will cause that schema delta to silently never reach the DB. Verify schema state before baselining unfamiliar environments.

**Runner:** executes SQL. Each file runs in a transaction, so a single failing file rolls back cleanly. **Failures are fatal:** `RunMigrationsOnStartup` calls `os.Exit(1)` if the DB isn't initialized or any migration fails to apply, rather than letting the API continue against a stale or partially-migrated schema. The same applies to a missing/unreadable `migrations` table on a remote DB or a permissions error. Boot-time crashes pre-traffic are preferable to half-migrated apps silently serving requests.

**NOT NULL caveat:** when the generator emits an ADD COLUMN for a field marked `gorm:"not null"` without a default, it omits the NOT NULL clause and adds a comment in the SQL. Existing rows would otherwise prevent the constraint from being added. Either declare a `default:` in the struct tag, or backfill manually and tighten the constraint in a follow-up migration.

**Type comparison fuzziness:** `-allow-type-changes` works but is approximate. PostgreSQL's `character varying(255)` and `VARCHAR(255)` look like the same thing to humans and to the normalizer in `struct_migration.go::normalizeType`, but edge cases exist (numeric precision, timezone handling, JSON/JSONB). Review carefully before applying.

## Troubleshooting

**"unsupported db type: """ + nil pointer panic when running a tool**
Config wasn't loaded. The tool needs to be run from a directory where it can find `config.json`, or you need to pass `-config=<path>` explicitly. The tool reads `globals.AppConfig` via the same path resolution as `main.startApplication`.

**FK type mismatch errors during migration apply (Error 3780 on MySQL, similar on others)**
A legacy file from the old generator with a placeholder INT primary key. Either fix the file's SQL by hand, or baseline past it if the DB already has the right schema.

**Generator says "no schema changes detected" but I just added a field**
Your local DB has already been migrated to the new schema (probably via an earlier `AutoMigrate` boot). Reset the local DB to the previous schema, or run the generator against a shadow/staging DB that's still on the old schema.

**API exits on startup with "Failed to run migrations — aborting startup"**
Intentional. The runner stops on the first failing file and aborts the whole boot. Read the error, fix the offending SQL (or baseline past the file if its changes are already in the DB), and restart. The runner picks up where it left off — files applied before the failure stay applied because each ran in its own transaction.

**MySQL error 1064 near `DELIMITER $$` in a hand-written migration**
The runner's SQL splitter recognises the MySQL `DELIMITER` directive, so stored procedures, triggers, and functions can be defined in migration files. The directive line is stripped before the SQL is sent to the server (since `DELIMITER` is a `mysql` CLI client convention, not actual SQL). If you still see this error, the file is probably missing the matching `DELIMITER ;` reset at the end of the procedure body, or the splitter version is older than April 2026 — update `services/migration.go::splitSQLStatements`.

**Warnings about "Migration statement skipped — object already exists"**
The runner tolerates a known set of idempotent-conflict errors — duplicate column, duplicate index name, table already exists, FK type-incompatibility — and logs them as warnings instead of aborting the boot. This handles the cross-environment case where a column was created by AutoMigrate on one database and by a migration file on another, so the file applies cleanly to one but conflicts on the other. See `services/migration.go::isIdempotentConflict` for the full list of patterns. If you see these warnings unexpectedly, it usually means a migration is redundant relative to the target database's current state, and you can leave the warnings or drop the file.

**Temporary tolerance for MySQL 3780 (FK incompatibility)**
The list of tolerated patterns currently includes `"are incompatible"`, which catches MySQL error 3780 — "Referencing column ... and referenced column ... in foreign key constraint ... are incompatible." This was added in April 2026 specifically to ride out the queue of pre-refactor legacy migration files that try to recreate FK constraints with hardcoded `INT` placeholder types against existing `BIGINT` columns. The error is idempotent-shaped *for the legacy files only*, since the FKs are already in place in the actual schema. **Remove this entry from `isIdempotentConflict` once every environment's `migrations` table has caught up past the legacy date range.** Leaving it in permanently would mask real type-mismatch bugs in future hand-written migrations.

**`migrations` table doesn't exist on a remote DB**
It's created on first call to `MigrationManager.Initialize`, which happens during both `RunMigrationsOnStartup` and `MarkAllMigrationsAsApplied`. Run either against the remote and the table will appear.

## Files of interest

- `api/services/struct_migration.go` — the diff-based generator. `GenerateMigrationForStructs` is the entry point.
- `api/services/migration.go` — `MigrationManager` (file applier), `RunMigrationsOnStartup`, `MarkAllMigrationsAsApplied`.
- `api/services/migrations.go` — bootstrap-only AutoMigrate functions for greenfield databases.
- `api/services/db.go` — `SetupTables` decides between bootstrap and incremental paths.
- `api/tools/generate_migration.go` — generator CLI.
- `api/tools/baseline_migrations.go` — one-shot baseline tool.
- `api/tools/run_migrations.go` — manual bootstrap runner (legacy).
- `api/migrations/<dialect>/` — generated SQL files.

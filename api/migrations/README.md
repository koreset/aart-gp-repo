# `migrations/`

This directory holds generated SQL migration files, organized by database dialect.

```
migrations/
├── postgresql/   # files applied when db_type=postgresql
├── mysql/        # files applied when db_type=mysql
└── mssql/        # files applied when db_type=mssql
```

Files are named `<timestamp>_<message>.sql`. The timestamp prefix determines apply order. The dialect subdirectory is selected automatically at runtime based on `services.DbBackend` (read from `config.json`'s `db_type`).

## Don't edit by hand unless you know what you're doing

These files are produced by `tools/generate_migration.go` and applied by `services.RunMigrationsOnStartup` on every API boot. The runner records each applied version in the `migrations` table inside the database; once a file has been applied, editing it has no effect on databases that already ran it.

If you need to change the schema, generate a new migration file rather than editing an existing one.

## See also

For the full system — how migrations are generated, the dev workflow, the cutover from the old generator, remote/production usage, and troubleshooting — see [`../MIGRATIONS.md`](../MIGRATIONS.md).

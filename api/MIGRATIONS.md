# Database Migrations

This document describes how to use the database migration system in this project.

## Overview

The project uses a combination of SQL-based migrations and Go-based migrations to manage the database schema. The SQL-based migrations are stored in the `migrations` directory and are executed automatically when the application starts. The Go-based migrations are defined in `services/migrations.go` and can be executed using the `run_migrations.go` script.

## Migration Files

Migration files are stored in the `migrations` directory and follow the naming convention:

```
YYYYMMDDHHMMSS_description.sql
```

For example:
```
202506240600001_add_sample_index.sql
```

The timestamp prefix ensures that migrations are executed in the correct order.

## Creating a New Migration

### Using the Utility Script

The easiest way to create a new migration is to use the provided utility script:

```bash
go run create_migration.go add_new_table
```

This will create a new migration file with the current timestamp and the provided description.

## Running Migrations

### Automatic Migrations

Migrations are automatically run when the application starts. The application checks for any new migration files in the `migrations` directory and executes them in order.

### Manual Migrations

You can also run migrations manually using the `run_migrations.go` script:

```bash
go run run_migrations.go -type=all
```

This will run all migrations in the correct order. You can also run specific types of migrations:

```bash
go run run_migrations.go -type=base       # Run base table migrations
go run run_migrations.go -type=pricing    # Run pricing table migrations
go run run_migrations.go -type=escalation # Run escalation table migrations
go run run_migrations.go -type=product    # Run product model table migrations
go run run_migrations.go -type=modelpoint # Run model point table migrations
go run run_migrations.go -type=gmm        # Run GMM table migrations
go run run_migrations.go -type=lic        # Run LIC table migrations
go run run_migrations.go -type=exposure   # Run exposure analysis table migrations
go run run_migrations.go -type=groupprice # Run group pricing table migrations
go run run_migrations.go -type=groupuser  # Run group pricing user table migrations
go run run_migrations.go -type=modelpoints # Run model point table updates
go run run_migrations.go -type=indexes    # Run database index creation
```

## Migration Implementation

The actual migration implementation is defined in `services/migrations.go`. This file contains functions for each type of migration:

- `MigrateBaseTables()`: Migrates base tables
- `MigratePricingTables()`: Migrates pricing tables
- `MigrateEscalationTables()`: Migrates escalation tables
- `MigrateProductModelTables()`: Migrates product model tables
- `MigrateModelPointTables()`: Migrates model point tables
- `MigrateGMMTables()`: Migrates GMM tables
- `MigrateLICTables()`: Migrates LIC tables
- `MigrateExposureAnalysisTables()`: Migrates exposure analysis tables
- `MigrateGroupPricingTables()`: Migrates group pricing tables
- `MigrateGroupPricingUserTables()`: Migrates group pricing user tables
- `UpdateModelPointTablesForMigration()`: Updates model point tables for all products
- `CreateDatabaseIndexesForMigration()`: Creates indexes on frequently queried tables

These functions use GORM's automigration functionality to create and update tables based on the model definitions.

## Migration History

Migrations are tracked in the `migrations` table in the database. This table stores information about which migrations have been applied, when they were applied, and who applied them.

## Best Practices

1. Always create a new migration file for schema changes rather than modifying existing migration files.
2. Make migrations idempotent (can be run multiple times without causing errors) whenever possible.
3. Test migrations thoroughly before deploying to production.
4. Include comments in migration files explaining what the migration does and why it's needed.
5. For complex migrations, consider breaking them down into smaller, more manageable migrations.
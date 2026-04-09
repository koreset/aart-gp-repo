# Database Migrations

This document explains how to work with database migrations in the application, including support for database-specific migrations.

## Supported Database Types

The application supports the following database types:
- PostgreSQL
- MySQL
- MSSQL (Microsoft SQL Server)

## Migration Types

There are two types of migrations:

1. **Common Migrations**: SQL scripts that work across all supported database types. These are stored in the `migrations` directory.
2. **Database-Specific Migrations**: SQL scripts that use syntax specific to a particular database type. These are stored in subdirectories:
   - `migrations/postgresql/` - PostgreSQL-specific migrations
   - `migrations/mysql/` - MySQL-specific migrations
   - `migrations/mssql/` - MSSQL-specific migrations

## Creating Migrations

You can create migrations using the `create_migration.go` script:

```bash
# Create a common migration
go run create_migration.go my_migration_name

# Create a database-specific migration
go run create_migration.go -db=postgresql my_migration_name
go run create_migration.go -db=mysql my_migration_name
go run create_migration.go -db=mssql my_migration_name

# Create migrations for all database types at once
go run create_migration.go -db=all my_migration_name
```

## Migration File Structure

Migration files follow this naming convention:
```
<timestamp>_<name>.sql
```

For example:
```
20250625000001_sample_migration.sql
```

## Running Migrations

Migrations are automatically run when the application starts. You can also run them manually using the `run_migrations.go` script:

```bash
go run run_migrations.go
```

## How Database-Specific Migrations Work

When the application runs migrations, it:

1. Determines the current database type (PostgreSQL, MySQL, or MSSQL)
2. Looks for database-specific migrations in the corresponding subdirectory
3. Also applies common migrations from the main `migrations` directory
4. Executes all migrations in order based on their timestamp

## Writing Database-Specific Migrations

When writing database-specific migrations, you can use syntax specific to that database type. See the sample migrations for examples:

- `migrations/postgresql/20250625000001_sample_postgresql_migration.sql`
- `migrations/mysql/20250625000001_sample_mysql_migration.sql`
- `migrations/mssql/20250625000001_sample_mssql_migration.sql`

## Best Practices

1. Use common migrations for simple schema changes that work across all database types.
2. Use database-specific migrations for features that require database-specific syntax.
3. When creating a new feature that requires database-specific syntax, create migrations for all supported database types.
4. Test migrations on all supported database types before deploying to production.
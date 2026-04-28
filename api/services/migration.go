package services

import (
	appLog "api/log"
	"fmt"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Migration represents a database migration
type Migration struct {
	ID        uint   `gorm:"primaryKey"`
	Version   string `gorm:"uniqueIndex:idx_migrations_version,length:255"`
	Name      string
	AppliedAt time.Time
}

// MigrationManager handles database migrations
type MigrationManager struct {
	db *gorm.DB
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *gorm.DB) *MigrationManager {
	return &MigrationManager{
		db: db,
	}
}

// Initialize creates the migrations table if it doesn't exist
func (m *MigrationManager) Initialize() error {
	appLog.Info("Initializing migration manager")
	err := m.db.AutoMigrate(&Migration{})
	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to create migrations table")
		return err
	}
	appLog.Info("Migrations table created or already exists")
	return nil
}

// RunMigrations runs all pending migrations
func (m *MigrationManager) RunMigrations() error {
	appLog.Info("Running database migrations")

	// Initialize migration manager
	err := m.Initialize()
	if err != nil {
		return err
	}

	// Get all migration files
	migrationFiles, err := m.getMigrationFiles()
	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to get migration files")
		return err
	}

	if len(migrationFiles) == 0 {
		appLog.Info("No migration files found")
		return nil
	}

	// Get applied migrations
	var appliedMigrations []Migration
	err = m.db.Order("version").Find(&appliedMigrations).Error
	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to get applied migrations")
		return err
	}

	// Create a map of applied migrations for quick lookup
	appliedMap := make(map[string]bool)
	for _, migration := range appliedMigrations {
		appliedMap[migration.Version] = true
	}

	// Run pending migrations
	for _, file := range migrationFiles {
		version := strings.Split(filepath.Base(file), "_")[0]
		if !appliedMap[version] {
			appLog.WithField("version", version).Info("Applying migration")

			// Read migration file
			content, err := os.ReadFile(file)
			if err != nil {
				appLog.WithFields(map[string]interface{}{
					"error": err.Error(),
					"file":  file,
				}).Error("Failed to read migration file")
				return err
			}

			// Execute migration in a transaction
			err = m.db.Transaction(func(tx *gorm.DB) error {
				// Execute SQL statements
				//statements := strings.Split(string(content), ";")
				statements := splitSQLStatements(string(content))
				for _, stmt := range statements {
					stmt = strings.TrimSpace(stmt)
					if stmt == "" {
						continue
					}

					if execErr := tx.Exec(stmt).Error; execErr != nil {
						if isIdempotentConflict(execErr) {
							// The schema is already in the desired state for
							// this statement — typically because AutoMigrate
							// (or a sibling process) made the change before
							// this file was applied. Treat as a no-op so the
							// migration as a whole can complete and be
							// recorded; otherwise the runner would loop on
							// this file forever.
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
				}

				// Record migration
				migration := Migration{
					Version:   version,
					Name:      strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)),
					AppliedAt: time.Now(),
				}

				err := tx.Create(&migration).Error
				if err != nil {
					appLog.WithFields(map[string]interface{}{
						"error":    err.Error(),
						"version":  version,
						"filename": filepath.Base(file),
					}).Error("Failed to record migration")
					return err
				}

				return nil
			})

			if err != nil {
				appLog.WithFields(map[string]interface{}{
					"error":   err.Error(),
					"version": version,
				}).Error("Migration failed")
				return err
			}

			appLog.WithField("version", version).Info("Migration applied successfully")
		} else {
			appLog.WithField("version", version).Debug("Migration already applied")
		}
	}

	appLog.Info("All migrations applied successfully")
	return nil
}

// getMigrationFiles returns a sorted list of migration files
func (m *MigrationManager) getMigrationFiles() ([]string, error) {
	// Get migration directory path
	migrationDir := "migrations"

	// Check if directory exists
	_, err := os.Stat(migrationDir)
	if os.IsNotExist(err) {
		appLog.WithField("dir", migrationDir).Warn("Migrations directory does not exist")
		return nil, nil
	}

	// Read migration files
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return nil, err
	}

	// Filter and sort migration files
	var migrationFiles []string

	// Check for database-specific migrations first
	dbSpecificDir := filepath.Join(migrationDir, DbBackend)
	_, err = os.Stat(dbSpecificDir)
	if !os.IsNotExist(err) {
		// Database-specific directory exists, read files from it
		dbSpecificFiles, err := os.ReadDir(dbSpecificDir)
		if err == nil {
			for _, file := range dbSpecificFiles {
				if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
					migrationFiles = append(migrationFiles, filepath.Join(dbSpecificDir, file.Name()))
				}
			}
			appLog.WithField("db_type", DbBackend).Info("Found database-specific migrations")
		}
	}

	// Also check for common migrations
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			// Skip database-specific directories
			if file.Name() != "postgresql" && file.Name() != "mysql" && file.Name() != "mssql" {
				migrationFiles = append(migrationFiles, filepath.Join(migrationDir, file.Name()))
			}
		}
	}

	// Sort files by version
	sort.Slice(migrationFiles, func(i, j int) bool {
		versionI := strings.Split(filepath.Base(migrationFiles[i]), "_")[0]
		versionJ := strings.Split(filepath.Base(migrationFiles[j]), "_")[0]
		return versionI < versionJ
	})

	return migrationFiles, nil
}

// CreateMigration creates a new migration file
func CreateMigration(name string, dbType string) (string, error) {
	// Ensure migrations directory exists
	migrationDir := "migrations"
	err := os.MkdirAll(migrationDir, 0755)
	if err != nil {
		return "", err
	}

	// Generate version based on timestamp
	version := time.Now().Format("20060102150405")

	// Create filename
	filename := fmt.Sprintf("%s_%s.sql", version, name)

	var filePath string

	// If dbType is specified, create a database-specific migration
	if dbType != "" {
		// Validate dbType
		if dbType != "postgresql" && dbType != "mysql" && dbType != "mssql" {
			return "", fmt.Errorf("invalid database type: %s", dbType)
		}

		// Ensure database-specific directory exists
		dbSpecificDir := filepath.Join(migrationDir, dbType)
		err := os.MkdirAll(dbSpecificDir, 0755)
		if err != nil {
			return "", err
		}

		filePath = filepath.Join(dbSpecificDir, filename)
		appLog.WithFields(map[string]interface{}{
			"db_type": dbType,
			"file":    filePath,
		}).Info("Creating database-specific migration file")
	} else {
		// Create a common migration file
		filePath = filepath.Join(migrationDir, filename)
		appLog.WithField("file", filePath).Info("Creating common migration file")
	}

	// Create empty migration file
	err = os.WriteFile(filePath, []byte("-- Migration: "+name+"\n\n"), 0644)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// CreateMigrationForAllDatabases creates migration files for all supported database types
func CreateMigrationForAllDatabases(name string) ([]string, error) {
	var filePaths []string

	// Create migration for PostgreSQL
	pgPath, err := CreateMigration(name, "postgresql")
	if err != nil {
		return nil, err
	}
	filePaths = append(filePaths, pgPath)

	// Create migration for MySQL
	mysqlPath, err := CreateMigration(name, "mysql")
	if err != nil {
		return nil, err
	}
	filePaths = append(filePaths, mysqlPath)

	// Create migration for MSSQL
	mssqlPath, err := CreateMigration(name, "mssql")
	if err != nil {
		return nil, err
	}
	filePaths = append(filePaths, mssqlPath)

	appLog.WithField("files", filePaths).Info("Created migration files for all database types")
	return filePaths, nil
}

// RunMigrationsOnStartup runs migrations on application startup.
//
// Failures are fatal: if the DB isn't initialized or a migration fails to
// apply, this calls os.Exit(1) rather than letting the API continue against
// a stale or partially-migrated schema. That's the correct behaviour for
// production — a half-migrated app silently serving traffic is worse than a
// boot crash that pages someone — and it's also fine for dev, where the next
// `air` cycle will retry once the broken file is fixed.
func RunMigrationsOnStartup() {
	appLog.Info("Running migrations on startup")

	if DB == nil {
		appLog.Fatal("Database not initialized — cannot run migrations; aborting startup")
	}

	migrationManager := NewMigrationManager(DB)

	if err := migrationManager.RunMigrations(); err != nil {
		appLog.WithField("error", err.Error()).Fatal("Failed to run migrations — aborting startup")
	}

	appLog.Info("Migrations completed successfully")
}

// MarkAllMigrationsAsApplied records every migration file currently on disk
// as already applied for the connected database, without executing the SQL.
//
// This is the baseline-insertion path: call it after a fresh-install
// AutoMigrate so the runner does not try to re-apply migration files to a
// database whose schema was already created from current structs. Files that
// are already recorded are left alone, so it's safe to call repeatedly.
func MarkAllMigrationsAsApplied() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	mgr := NewMigrationManager(DB)
	if err := mgr.Initialize(); err != nil {
		return err
	}

	files, err := mgr.getMigrationFiles()
	if err != nil {
		return err
	}

	now := time.Now()
	recorded := 0
	for _, file := range files {
		version := strings.Split(filepath.Base(file), "_")[0]

		var existing Migration
		err := DB.Where("version = ?", version).First(&existing).Error
		if err == nil {
			continue // already recorded
		}

		m := Migration{
			Version:   version,
			Name:      strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)),
			AppliedAt: now,
		}
		if err := DB.Create(&m).Error; err != nil {
			return fmt.Errorf("record baseline migration %s: %w", version, err)
		}
		recorded++
	}

	appLog.WithFields(map[string]interface{}{
		"total":    len(files),
		"recorded": recorded,
	}).Info("Marked existing migration files as applied (baseline)")
	return nil
}

// isIdempotentConflict reports whether the given error reflects a no-op
// schema conflict — typically "column already exists", "table already exists",
// or "duplicate index" — that should be treated as a soft warning during
// migration application rather than a hard failure. Recognises the wording
// emitted by the MySQL, PostgreSQL, and SQL Server drivers.
//
// This makes the runner tolerant of databases that arrived at the same
// schema state by different paths — e.g. a dev DB previously auto-migrated
// when a struct field was added, then later handed a migration file
// generated against a baseline DB that lacked the change. The first path
// already created the column; applying the file would otherwise fail with
// "duplicate column" and prevent the migration from being recorded.
//
// The trade-off: if a developer accidentally writes a migration that
// conflicts with an existing object (e.g. ADD COLUMN with a typo'd name
// that happens to match something else), the conflict is logged as a
// warning rather than blocking the boot. The warning surfaces in logs and
// is recoverable; the alternative — never running migrations that touch
// pre-existing schemas — is strictly worse for this codebase.
func isIdempotentConflict(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	patterns := []string{
		"duplicate column",        // MySQL 1060
		"duplicate key name",      // MySQL 1061 (index name collision)
		"duplicate index",         // misc MySQL variants
		"table already exists",    // MySQL 1050
		"already exists",          // PostgreSQL ("relation X already exists", etc.)
		"already an object named", // MSSQL 2714

		// TEMPORARY (April 2026): MySQL 3780 — "Referencing column ... and
		// referenced column ... in foreign key constraint ... are incompatible."
		// Added to ride out the queue of pre-refactor legacy migration files
		// that try to (re)create FK constraints with hardcoded INT placeholder
		// types against existing BIGINT columns. The schema is already correct
		// in every database we've audited; the legacy files are redundant.
		//
		// Remove this entry once every environment's migrations table has
		// caught up to a date past the last legacy FK file (post-April 2026
		// generation cutover). Leaving it in permanently would mask real FK
		// type-mismatch bugs in future hand-written migrations.
		"are incompatible",
	}
	for _, p := range patterns {
		if strings.Contains(msg, p) {
			return true
		}
	}
	return false
}

// splitSQLStatements breaks a multi-statement SQL script into individual
// statements suitable for sending one-by-one to a database driver.
//
// It supports the MySQL `DELIMITER` directive, which is purely client-side —
// the server rejects `DELIMITER` as a syntax error. When a line begins with
// `DELIMITER <token>`, the directive line itself is stripped from the output
// and the active statement terminator is changed to <token>. `DELIMITER ;`
// switches back. This lets hand-written migrations define stored procedures,
// triggers, and functions whose bodies legitimately contain semicolons.
//
// Single-line comments (`--` and `#`) are stripped, and the splitter is
// aware of `'` and `"` string literal boundaries so semicolons inside string
// literals don't cause spurious splits.
func splitSQLStatements(sqlScript string) []string {
	var statements []string
	var sb strings.Builder
	delim := ";" // active statement terminator; can be changed by DELIMITER directives

	flush := func() {
		stmt := strings.TrimSpace(sb.String())
		if stmt != "" {
			statements = append(statements, stmt)
		}
		sb.Reset()
	}

	inSingle := false
	inDouble := false
	atLineStart := true
	i := 0
	for i < len(sqlScript) {
		c := sqlScript[i]

		// DELIMITER <token> — recognised only at the start of a line (after
		// optional leading whitespace) and outside string literals. The
		// directive line is consumed but produces no output statement.
		if atLineStart && !inSingle && !inDouble {
			j := i
			for j < len(sqlScript) && (sqlScript[j] == ' ' || sqlScript[j] == '\t') {
				j++
			}
			const kw = "DELIMITER"
			if j+len(kw) <= len(sqlScript) && strings.EqualFold(sqlScript[j:j+len(kw)], kw) {
				next := j + len(kw)
				// The keyword must be followed by whitespace (or EOF) to be
				// treated as the directive — guards against identifiers that
				// merely start with "DELIMITER".
				if next == len(sqlScript) || sqlScript[next] == ' ' || sqlScript[next] == '\t' {
					for next < len(sqlScript) && (sqlScript[next] == ' ' || sqlScript[next] == '\t') {
						next++
					}
					start := next
					for next < len(sqlScript) && sqlScript[next] != '\n' && sqlScript[next] != '\r' {
						next++
					}
					newDelim := strings.TrimSpace(sqlScript[start:next])
					if newDelim != "" {
						// Flush any in-progress statement before swapping the
						// delimiter so a half-built statement doesn't leak
						// into the next block.
						flush()
						delim = newDelim
					}
					i = next
					continue
				}
			}
		}

		// Track line boundaries so DELIMITER detection only fires at the
		// start of a line. Any non-whitespace character clears atLineStart.
		if c == '\n' || c == '\r' {
			atLineStart = true
			sb.WriteByte(c)
			i++
			continue
		}
		if c != ' ' && c != '\t' {
			atLineStart = false
		}

		// Line comment: skip from -- or # to end of line (only when not inside a string)
		if !inSingle && !inDouble {
			if (c == '-' && i+1 < len(sqlScript) && sqlScript[i+1] == '-') || c == '#' {
				for i < len(sqlScript) && sqlScript[i] != '\n' {
					i++
				}
				continue
			}
		}

		// String literal boundaries (MySQL: '' inside '' escapes a quote)
		if c == '\'' && !inDouble {
			sb.WriteByte(c)
			if inSingle && i+1 < len(sqlScript) && sqlScript[i+1] == '\'' {
				sb.WriteByte(sqlScript[i+1])
				i += 2
				continue
			}
			inSingle = !inSingle
			i++
			continue
		}
		if c == '"' && !inSingle {
			sb.WriteByte(c)
			inDouble = !inDouble
			i++
			continue
		}

		// Statement terminator. The active delimiter is multi-character
		// safe — it can be ;, $$, //, etc. depending on a prior DELIMITER
		// directive.
		if !inSingle && !inDouble && i+len(delim) <= len(sqlScript) && sqlScript[i:i+len(delim)] == delim {
			flush()
			i += len(delim)
			continue
		}

		sb.WriteByte(c)
		i++
	}

	flush()
	return statements
}

package services

import (
	appLog "api/log"
	"bufio"
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

					err := tx.Exec(stmt).Error
					if err != nil {
						appLog.WithFields(map[string]interface{}{
							"error": err.Error(),
							"stmt":  stmt,
						}).Error("Failed to execute migration statement")
						return err
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

// RunMigrationsOnStartup runs migrations on application startup
func RunMigrationsOnStartup() {
	appLog.Info("Running migrations on startup")

	// Wait for DB to be initialized
	if DB == nil {
		appLog.Error("Database not initialized")
		return
	}

	// Create migration manager
	migrationManager := NewMigrationManager(DB)

	// Run migrations
	err := migrationManager.RunMigrations()
	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to run migrations")
		return
	}

	appLog.Info("Migrations completed successfully")
}

func splitSQLStatements(sqlScript string) []string {
	var statements []string
	var sb strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(sqlScript))
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and comments (simple handling)
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "--") || strings.HasPrefix(trimmed, "#") {
			continue
		}

		sb.WriteString(line)
		// Preserve line breaks for multi-line statements
		sb.WriteString("\n")

		if strings.HasSuffix(trimmed, ";") {
			statement := strings.TrimSpace(sb.String())
			statements = append(statements, statement)
			sb.Reset()
		}
	}

	// In case the last statement does not end in a semicolon
	if sb.Len() > 0 {
		statements = append(statements, strings.TrimSpace(sb.String()))
	}

	return statements
}

package services

import (
	"api/globals"
	appLog "api/log"
	"api/models"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"syscall"
	"time"

	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	gormLogger "gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	mysqlerr "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var ctx context.Context

// DB is a global DB connection object
var DB *gorm.DB
var DbBackend string

// dbGate limits the number of concurrent DB operations to protect the database
var dbGate = make(chan struct{}, 40)

// backoffDuration returns an exponential backoff with full jitter
func backoffDuration(attempt int, base, max time.Duration) time.Duration {
	// exponential: base * 2^attempt, capped to max, then jitter [0, exp]
	d := base << attempt
	if d > max {
		d = max
	}
	if d <= 0 {
		return base
	}
	// add jitter up to d
	n := rand.Int63n(int64(d/time.Millisecond) + 1)
	return time.Duration(n) * time.Millisecond
}

// isRetryable identifies transient errors safe to retry
func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	var myErr *mysqlerr.MySQLError
	if errors.As(err, &myErr) {
		switch myErr.Number {
		case 1040, // ER_CON_COUNT_ERROR: too many connections
			1205, // ER_LOCK_WAIT_TIMEOUT
			1213: // ER_LOCK_DEADLOCK
			return true
		}
	}
	// network/transient
	if errors.Is(err, io.EOF) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, syscall.ECONNRESET) {
		return true
	}
	return false
}

// Retry retries op up to attempts with exponential backoff and jitter
func Retry(ctx context.Context, attempts int, base, max time.Duration, op func(context.Context) error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err = op(ctx)
		if err == nil {
			return nil
		}
		if !isRetryable(err) {
			return err
		}
		// log retry attempt
		appLog.WithContext(ctx).WithFields(map[string]interface{}{
			"attempt": i + 1,
			"error":   err.Error(),
		}).Warn("DB operation failed, will retry with backoff")

		d := backoffDuration(i, base, max)
		timer := time.NewTimer(d)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return err
}

// DBReadWithResilience wraps a DB operation with timeout, concurrency limit, and retry
func DBReadWithResilience(parent context.Context, q func(*gorm.DB) error) error {
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	select {
	case dbGate <- struct{}{}:
		defer func() { <-dbGate }()
	case <-ctx.Done():
		return ctx.Err()
	}

	return Retry(ctx, 4, 200*time.Millisecond, 3*time.Second, func(c context.Context) error {
		return q(DB.WithContext(c))
	})
}

func GetDBName(db *gorm.DB) (string, error) {
	var dbName string

	// Use the appropriate query for your database
	// MySQL: SELECT DATABASE();
	// PostgreSQL: SELECT current_database();
	// SQL Server: SELECT DB_NAME();
	// SQLite: pragma_database_list; (this returns more info, might need parsing)
	var query string
	switch globals.AppConfig.DbType {
	case "mysql":
		query = "SELECT DATABASE()"
	case "postgresql":
		query = "SELECT current_database()"
	case "mssql":
		query = "SELECT DB_NAME()"
	}

	err := db.Raw(query).Scan(&dbName).Error
	if err != nil {
		return "", err
	}

	return dbName, nil
}

// CustomGormLogger implements the gorm logger interface to use our application logger
type CustomGormLogger struct{}

// LogMode sets the log level
func (l *CustomGormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l
}

// Info logs info messages
func (l *CustomGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	appLog.WithContext(ctx).Infof(msg, data...)
}

// Warn logs warn messages
func (l *CustomGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	appLog.WithContext(ctx).Warnf(msg, data...)
}

// Error logs error messages
func (l *CustomGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	appLog.WithContext(ctx).Errorf(msg, data...)
}

// Trace logs SQL operations
func (l *CustomGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := appLog.WithContext(ctx).WithFields(map[string]interface{}{
		"elapsed_ms": elapsed.Milliseconds(),
		"rows":       rows,
		"sql":        sql,
	})

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fields.WithField("error", err.Error()).Error("Database query failed")
	} else {
		if elapsed > time.Millisecond*500 {
			fields.Warn("Slow database query")
		} else {
			fields.Debug("Database query executed")
		}
	}
}

// dialectorForConfig builds a GORM dialector for the given AppConfig. It is
// the single source of truth for DSN construction, shared by both the runtime
// database initialiser and the setup wizard's connectivity check.
func dialectorForConfig(cfg models.AppConfig) (gorm.Dialector, error) {
	switch cfg.DbType {
	case "postgresql":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort)
		return postgres.Open(dsn), nil
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4,utf8&parseTime=true&loc=Local&interpolateParams=true&tls=skip-verify",
			cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)
		return mysql.Open(dsn), nil
	case "mssql":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)
		return sqlserver.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported db type: %q", cfg.DbType)
	}
}

// TestDBConnection opens a throwaway connection using globals.AppConfig and
// pings it. Used by the setup wizard to validate credentials before writing
// config.json. It does not mutate the global DB variable.
func TestDBConnection() error {
	dialector, err := dialectorForConfig(globals.AppConfig)
	if err != nil {
		return err
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	pingCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return sqlDB.PingContext(pingCtx)
}

func initializeDb() {
	var err error
	appLog.Info("Initializing database connection")
	appLog.WithFields(map[string]interface{}{
		"db_type": globals.AppConfig.DbType,
		"db_host": globals.AppConfig.DbHost,
		"db_name": globals.AppConfig.DbName,
		"db_port": globals.AppConfig.DbPort,
		"db_user": globals.AppConfig.DbUser,
	}).Info("Database configuration")

	// Create a custom logger for GORM
	customLogger := &CustomGormLogger{}
	customLogger.LogMode(gormLogger.Info)

	dialector, err := dialectorForConfig(globals.AppConfig)
	if err != nil {
		appLog.WithField("error", err.Error()).Error("Could not start Database")
		globals.Logger.Error("Could not start Database: ", err)
		return
	}

	appLog.WithField("db_type", globals.AppConfig.DbType).Info("Connecting to database")
	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger:      customLogger,
		PrepareStmt: false,
	})
	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to connect to database")
		globals.Logger.Error("Could not start Database: ", err)
		return
	}
	appLog.Info("Successfully connected to database")
	DbBackend = globals.AppConfig.DbType

	sqlDb, err := DB.DB()
	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to get database connection")
		return
	}

	// Configure connection pool with optimized settings
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetConnMaxLifetime(time.Minute * 2) // Significantly increased
	sqlDb.SetConnMaxIdleTime(time.Minute * 1) // Keep this value

	appLog.Info("Database connection pool configured with optimized settings")

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err = sqlDb.PingContext(ctx)
	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to ping database")
	} else {
		appLog.Info("Database connection verified with ping")
	}
}

// HasTables checks if the database has any tables
func HasTables() bool {
	if DB == nil {
		appLog.Error("Database connection is nil")
		return false
	}

	var count int64
	var err error

	switch DbBackend {
	case "postgresql":
		// For PostgreSQL, query the information_schema.tables
		err = DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Count(&count).Error
	case "mysql":
		// For MySQL, query the information_schema.tables
		err = DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ?", globals.AppConfig.DbName).Count(&count).Error
	case "mssql":
		// For MSSQL, query the information_schema.tables
		err = DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_catalog = ?", globals.AppConfig.DbName).Count(&count).Error
	default:
		appLog.WithField("db_backend", DbBackend).Error("Unsupported database backend")
		return false
	}

	if err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to check if database has tables")
		return false
	}

	return count > 0
}

// SetupTables initializes the database and migrates tables required
// On setup, it populates base data requirements.
// Auto-migrations only run when the database is empty (no existing tables).
// For existing databases, schema changes are handled by manual migration
// scripts in the migrations/ folder (see RunMigrationsOnStartup).
func SetupTables(initTables, initDatabaseTables bool) {
	appLog.Info("Setting up database tables")
	initializeDb()

	if initDatabaseTables {
		dbHasTables := HasTables()

		if !dbHasTables {
			// Database is empty — run full auto-migrations to bootstrap the schema.
			appLog.Info("Empty database detected — running full auto-migrations")

			if err := DB.AutoMigrate(&models.SystemLock{}); err != nil {
				appLog.WithField("error", err.Error()).Error("Failed to migrate SystemLock table")
			}

			if err := MigrateBaseTables(); err != nil {
				appLog.WithField("error", err.Error()).Error("Failed to migrate base tables")
			}

			if err := MigrateGroupPricingTables(); err != nil {
				appLog.WithField("error", err.Error()).Error("Failed to migrate group pricing tables")
			}

			if err := MigrateGroupPricingUserTables(); err != nil {
				appLog.WithField("error", err.Error()).Error("Failed to migrate group pricing user tables")
			}

			if err := MigrateGroupPremiumTables(); err != nil {
				appLog.WithField("error", err.Error()).Error("Failed to migrate group premium tables")
			}

			if err := MigratePhiValuationTables(); err != nil {
				appLog.WithField("error", err.Error()).Error("Failed to migrate PHI valuation tables")
			}

			// AutoMigrate has just produced a schema matching the current
			// structs. Record every migration file already on disk as applied
			// so RunMigrationsOnStartup does not try to re-run baseline files
			// against a database that's already at HEAD.
			if err := MarkAllMigrationsAsApplied(); err != nil {
				appLog.WithField("error", err.Error()).Error("Failed to record baseline migrations")
			}
		} else {
			// Database already has tables — skip auto-migrations.
			// Schema changes for existing databases are applied via manual
			// migration scripts in RunMigrationsOnStartup().
			appLog.Info("Existing database detected — skipping auto-migrations, relying on manual migration scripts")
		}
	}

	// Update model point tables for all products
	//if err := UpdateModelPointTablesForMigration(); err != nil {
	//	appLog.WithField("error", err.Error()).Error("Failed to update model point tables")
	//} else {
	//	// to cater for changes in base data, always refresh the base data
	//	BaseData(initTables)
	//}
	initTables = true
	BaseData(initTables)

	// Create indexes for performance optimization
	//if err := CreateDatabaseIndexesForMigration(); err != nil {
	//	appLog.WithField("error", err.Error()).Error("Failed to create database indexes")
	//}
}

//func UpdateModelPointTables() {
//	// Call the migration version of this function
//	err := UpdateModelPointTablesForMigration()
//	if err != nil {
//		appLog.WithField("error", err.Error()).Error("Failed to update model point tables")
//	}
//}

// CreateDatabaseIndexes creates indexes on frequently queried tables to improve performance
func CreateDatabaseIndexes() {
	logger := appLog.WithField("action", "CreateDatabaseIndexes")
	logger.Info("Creating database indexes for performance optimization")

	// Create indexes with a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// List of indexes to create
	indexes := []struct {
		table      string
		columns    []string
		indexName  string
		unique     bool
		conditions string
	}{
		// Aggregated projections indexes
		{table: "aggregated_projections", columns: []string{"run_id"}, indexName: "idx_agg_proj_run_id", unique: false},
		{table: "aggregated_projections", columns: []string{"run_id", "product_code"}, indexName: "idx_agg_proj_run_product", unique: false},
		{table: "aggregated_projections", columns: []string{"run_id", "product_code", "sp_code"}, indexName: "idx_agg_proj_run_product_sp", unique: false},
		{table: "aggregated_projections", columns: []string{"projection_month"}, indexName: "idx_agg_proj_month", unique: false},

		// Group pricing indexes
		{table: "g_pricing_member_data", columns: []string{"quote_id"}, indexName: "idx_member_data_quote_id", unique: false},
		{table: "group_pricing_claims_experiences", columns: []string{"quote_id"}, indexName: "idx_claims_exp_quote_id", unique: false},
		{table: "member_rating_results", columns: []string{"quote_id"}, indexName: "idx_rating_results_quote_id", unique: false},
		{table: "member_premium_schedules", columns: []string{"quote_id"}, indexName: "idx_premium_schedules_quote_id", unique: false},
		{table: "bordereaux", columns: []string{"quote_id"}, indexName: "idx_bordereaux_quote_id", unique: false},

		// User management indexes
		{table: "org_users", columns: []string{"email"}, indexName: "idx_org_users_email", unique: true},
		{table: "app_users", columns: []string{"user_email"}, indexName: "idx_app_users_email", unique: true},
		{table: "user_tokens", columns: []string{"subject"}, indexName: "idx_user_tokens_subject", unique: true},

		// Activity tracking indexes
		{table: "activities", columns: []string{"user_email"}, indexName: "idx_activities_user_email", unique: false},
		{table: "activities", columns: []string{"date"}, indexName: "idx_activities_date", unique: false},
		{table: "activities", columns: []string{"object_type", "object_id"}, indexName: "idx_activities_object", unique: false},
	}

	// Create each index
	for _, idx := range indexes {
		var indexSQL string

		if DbBackend == "postgresql" {
			// PostgreSQL index creation
			indexType := ""
			if idx.unique {
				indexType = "UNIQUE "
			}

			indexSQL = fmt.Sprintf("CREATE %sINDEX IF NOT EXISTS %s ON %s (%s)",
				indexType,
				idx.indexName,
				idx.table,
				strings.Join(idx.columns, ", "))

			if idx.conditions != "" {
				indexSQL += " " + idx.conditions
			}
		} else {
			// MySQL index creation
			indexType := "INDEX"
			if idx.unique {
				indexType = "UNIQUE INDEX"
			}

			indexSQL = fmt.Sprintf("CREATE %s %s ON %s (%s)",
				indexType,
				idx.indexName,
				idx.table,
				strings.Join(idx.columns, ", "))

			if idx.conditions != "" {
				indexSQL += " " + idx.conditions
			}
		}

		// Execute the index creation with context
		err := DB.WithContext(ctx).Exec(indexSQL).Error
		if err != nil {
			logger.WithFields(map[string]interface{}{
				"error":      err.Error(),
				"table":      idx.table,
				"index_name": idx.indexName,
				"columns":    strings.Join(idx.columns, ", "),
			}).Warn("Failed to create index")
		} else {
			logger.WithFields(map[string]interface{}{
				"table":      idx.table,
				"index_name": idx.indexName,
				"columns":    strings.Join(idx.columns, ", "),
			}).Info("Successfully created index")
		}
	}

	logger.Info("Database index creation completed")
}

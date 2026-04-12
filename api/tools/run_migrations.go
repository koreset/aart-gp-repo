//go:build ignore

package main

import (
	"api/log"
	"api/services"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Parse command line arguments
	var migrationType string
	flag.StringVar(&migrationType, "type", "all", "Type of migration to run (all, base, pricing, escalation, product, modelpoint, gmm, lic, exposure, groupprice, groupuser, modelpoints, indexes)")
	flag.Parse()

	// Initialize database connection
	services.SetupTables(false, false)

	log.Info("Running migrations of type: " + migrationType)

	var err error

	switch migrationType {
	case "all":
		err = runAllMigrations()
	case "base":
		err = services.MigrateBaseTables()
	case "pricing":
		err = services.MigratePricingTables()
	case "escalation":
		err = services.MigrateEscalationTables()
	case "product":
		err = services.MigrateProductModelTables()
	case "modelpoint":
		err = services.MigrateModelPointTables()
	case "gmm":
		err = services.MigrateGMMTables()
	case "lic":
		err = services.MigrateLICTables()
	case "exposure":
		err = services.MigrateExposureAnalysisTables()
	case "groupprice":
		err = services.MigrateGroupPricingTables()
	case "groupuser":
		err = services.MigrateGroupPricingUserTables()
	//case "modelpoints":
	//	err = services.UpdateModelPointTablesForMigration()
	case "indexes":
		err = services.CreateDatabaseIndexesForMigration()
	default:
		fmt.Printf("Unknown migration type: %s\n", migrationType)
		os.Exit(1)
	}

	if err != nil {
		log.WithField("error", err.Error()).Error("Migration failed")
		os.Exit(1)
	}

	log.Info("Migration completed successfully")
}

func runAllMigrations() error {
	log.Info("Running all migrations")

	// Run migrations in the correct order
	if err := services.MigrateBaseTables(); err != nil {
		return err
	}

	if err := services.MigratePricingTables(); err != nil {
		return err
	}

	if err := services.MigrateEscalationTables(); err != nil {
		return err
	}

	if err := services.MigrateProductModelTables(); err != nil {
		return err
	}

	if err := services.MigrateModelPointTables(); err != nil {
		return err
	}

	if err := services.MigrateGMMTables(); err != nil {
		return err
	}

	if err := services.MigrateLICTables(); err != nil {
		return err
	}

	if err := services.MigrateExposureAnalysisTables(); err != nil {
		return err
	}

	if err := services.MigrateGroupPricingTables(); err != nil {
		return err
	}

	if err := services.MigrateGroupPricingUserTables(); err != nil {
		return err
	}

	//if err := services.UpdateModelPointTablesForMigration(); err != nil {
	//	return err
	//}

	if err := services.CreateDatabaseIndexesForMigration(); err != nil {
		return err
	}

	log.Info("All migrations completed successfully")
	return nil
}

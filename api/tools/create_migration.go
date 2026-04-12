//go:build ignore

package main

import (
	"api/services"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Define flags
	dbTypePtr := flag.String("db", "", "Database type (postgresql, mysql, mssql, or all)")

	// Parse command line arguments
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Usage: go run create_migration.go [-db=<database_type>] <migration_name>")
		fmt.Println("  database_type: postgresql, mysql, mssql, or all (default: creates a common migration)")
		os.Exit(1)
	}

	// Join all arguments to create the migration name
	migrationName := strings.Join(args, "_")

	// Determine which type of migration to create
	dbType := *dbTypePtr

	if dbType == "all" {
		// Create migrations for all database types
		filepaths, err := services.CreateMigrationForAllDatabases(migrationName)
		if err != nil {
			fmt.Printf("Error creating migrations: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Created migration files:")
		for _, filepath := range filepaths {
			fmt.Printf("- %s\n", filepath)
		}
		fmt.Println("Edit these files to add your database-specific SQL statements.")
	} else {
		// Create a single migration file (either database-specific or common)
		filepath, err := services.CreateMigration(migrationName, dbType)
		if err != nil {
			fmt.Printf("Error creating migration: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Created migration file: %s\n", filepath)
		if dbType != "" {
			fmt.Printf("This is a %s-specific migration. Edit this file to add your SQL statements.\n", dbType)
		} else {
			fmt.Println("This is a common migration. Edit this file to add your SQL statements.")
			fmt.Println("For database-specific migrations, use the -db flag with postgresql, mysql, or mssql.")
		}
	}
}

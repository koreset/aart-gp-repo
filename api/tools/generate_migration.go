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
	dbTypePtr := flag.String("db", "all", "Database type (postgresql, mysql, mssql, or all)")
	structNamePtr := flag.String("struct", "", "Name of the struct to generate migration for")

	// Parse command line arguments
	flag.Parse()

	if *structNamePtr == "" {
		fmt.Println("Usage: go run generate_migration.go -struct=<struct_name> [-db=<database_type>]")
		fmt.Println("  struct_name: Name of the struct to generate migration for (e.g., ProductMargins)")
		fmt.Println("  database_type: postgresql, mysql, mssql, or all (default: all)")
		os.Exit(1)
	}

	// Generate migration name based on struct name
	migrationName := fmt.Sprintf("update_%s", strings.ToLower(*structNamePtr))

	// Generate SQL for the specified struct
	err := services.GenerateMigrationForStruct(*structNamePtr, migrationName, *dbTypePtr)
	if err != nil {
		fmt.Printf("Error generating migration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated migration for struct %s\n", *structNamePtr)
}

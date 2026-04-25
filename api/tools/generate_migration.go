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
	structNamesPtr := flag.String("struct", "", "Comma-separated struct name(s) to migrate (e.g., Beneficiary or Beneficiary,GroupScheme)")
	messagePtr := flag.String("message", "", "Short message used in the generated filename (e.g., add_nickname)")
	allowTypeChangesPtr := flag.Bool("allow-type-changes", false, "Emit ALTER COLUMN TYPE when types differ. Off by default — comparison is fuzzy.")
	allowDestructivePtr := flag.Bool("allow-destructive", false, "Emit DROP COLUMN/DROP INDEX for objects absent from the struct. Off by default.")
	flag.Parse()

	if *structNamesPtr == "" {
		fmt.Println("Usage: go run tools/generate_migration.go -struct=<Name>[,<Name>...] [-message=<text>] [-allow-type-changes] [-allow-destructive]")
		fmt.Println()
		fmt.Println("Notes:")
		fmt.Println("  • Dialect is auto-detected from the connected database.")
		fmt.Println("  • Output goes to migrations/<dialect>/<timestamp>_<message>.sql.")
		fmt.Println("  • Run against a database that is still on the previous schema —")
		fmt.Println("    typically a shadow/staging copy — so the diff captures the new fields.")
		os.Exit(1)
	}

	// Initialize database connection (does not auto-migrate tables).
	services.SetupTables(false, false)

	// Parse comma-separated struct list.
	rawNames := strings.Split(*structNamesPtr, ",")
	names := make([]string, 0, len(rawNames))
	for _, n := range rawNames {
		n = strings.TrimSpace(n)
		if n != "" {
			names = append(names, n)
		}
	}
	if len(names) == 0 {
		fmt.Println("No valid struct names provided.")
		os.Exit(1)
	}

	opts := services.GenerateOptions{
		AllowTypeChanges: *allowTypeChangesPtr,
		AllowDestructive: *allowDestructivePtr,
		Message:          *messagePtr,
	}

	path, err := services.GenerateMigrationForStructs(names, opts)
	if err != nil {
		fmt.Printf("Error generating migration: %v\n", err)
		os.Exit(1)
	}

	if path == "" {
		fmt.Println("No schema changes detected — no migration file was written.")
		return
	}

	fmt.Printf("Generated migration: %s\n", path)
}

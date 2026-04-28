//go:build ignore

package main

import (
	"api/globals"
	"api/services"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	structNamesPtr := flag.String("struct", "", "Comma-separated struct name(s) to migrate (e.g., Beneficiary or Beneficiary,GroupScheme)")
	messagePtr := flag.String("message", "", "Short message used in the generated filename (e.g., add_nickname)")
	allowTypeChangesPtr := flag.Bool("allow-type-changes", false, "Emit ALTER COLUMN TYPE when types differ. Off by default — comparison is fuzzy.")
	allowDestructivePtr := flag.Bool("allow-destructive", false, "Emit DROP COLUMN/DROP INDEX for objects absent from the struct. Off by default.")
	configPathPtr := flag.String("config", "", "Path to config.json. Defaults to ./config.json on Unix, alongside the executable on Windows.")
	flag.Parse()

	if *structNamesPtr == "" {
		fmt.Println("Usage: go run tools/generate_migration.go -struct=<Name>[,<Name>...] [-message=<text>] [-allow-type-changes] [-allow-destructive] [-config=<path>]")
		fmt.Println()
		fmt.Println("Notes:")
		fmt.Println("  • Dialect is auto-detected from the connected database.")
		fmt.Println("  • Output goes to migrations/<dialect>/<timestamp>_<message>.sql.")
		fmt.Println("  • Run against a database that is still on the previous schema —")
		fmt.Println("    typically a shadow/staging copy — so the diff captures the new fields.")
		os.Exit(1)
	}

	if err := loadConfig(*configPathPtr); err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
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

// loadConfig reads a config file into globals.AppConfig. If override is empty,
// it falls back to the same path resolution main.startApplication uses:
// ./config.json on non-Windows, or alongside the executable on Windows. Pass
// an explicit path to point at a different environment (e.g. prod).
func loadConfig(override string) error {
	path := override
	if path == "" {
		path = "config.json"
		if runtime.GOOS == "windows" {
			if exec, err := os.Executable(); err == nil {
				dir, _ := filepath.Split(exec)
				path = filepath.Join(dir, "config.json")
			}
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &globals.AppConfig); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

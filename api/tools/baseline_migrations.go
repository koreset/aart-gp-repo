//go:build ignore

// baseline_migrations is a one-shot tool for the cutover from the old
// full-struct-dump generator to the new diff-based generator.
//
// Run it once against each environment whose database already has the schema
// (typically because earlier AutoMigrate runs created the tables) but whose
// `migrations` table does not yet record the legacy `.sql` files in
// migrations/<dialect>/. After this runs, every file currently on disk is
// recorded as already applied for that database, so the runner won't try to
// re-apply them on the next boot.
//
// Usage (from the api/ directory, with config.json present):
//   go run tools/baseline_migrations.go
//
// Safe to run repeatedly — files already recorded are skipped.
package main

import (
	appLog "api/log"
	"api/globals"
	"api/services"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	if err := loadConfig(); err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize the database connection without running AutoMigrate.
	services.SetupTables(false, false)

	appLog.Info("Marking all migration files as applied (baseline)")

	if err := services.MarkAllMigrationsAsApplied(); err != nil {
		fmt.Printf("Error marking migrations as applied: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done. The migrations table now reflects every file under migrations/<dialect>/")
	fmt.Println("as already applied. RunMigrationsOnStartup will only act on files generated")
	fmt.Println("after this point.")
}

// loadConfig reads config.json into globals.AppConfig, mirroring the path
// resolution in main.startApplication so tools see the same config the API
// does. On non-Windows it reads ./config.json from the working directory; on
// Windows it sits next to the executable.
func loadConfig() error {
	path := "config.json"
	if runtime.GOOS == "windows" {
		if exec, err := os.Executable(); err == nil {
			dir, _ := filepath.Split(exec)
			path = filepath.Join(dir, "config.json")
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

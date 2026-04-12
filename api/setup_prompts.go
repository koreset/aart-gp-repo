package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"api/globals"
	"api/services"

	"github.com/manifoldco/promptui"
)

// runConfigWizard drives the interactive prompts that populate globals.AppConfig.
// When reconfigure is true, current values are used as defaults and the DB
// password may be kept unchanged by entering a blank value.
func runConfigWizard(reconfigure bool) error {
	cfg := &globals.AppConfig

	// --- Database ---
	dbLabels := []string{"MySQL", "PostgreSQL", "SQL Server"}
	labelByType := map[string]string{"mysql": "MySQL", "postgresql": "PostgreSQL", "mssql": "SQL Server"}
	typeByLabel := map[string]string{"MySQL": "mysql", "PostgreSQL": "postgresql", "SQL Server": "mssql"}

	dbLabel, err := promptSelect("Select the database type", dbLabels, labelByType[cfg.DbType])
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}
	cfg.DbType = typeByLabel[dbLabel]

	for {
		if cfg.DbName, err = promptRequired("Database name", cfg.DbName); err != nil {
			return err
		}
		if cfg.DbUser, err = promptRequired("Database user", cfg.DbUser); err != nil {
			return err
		}
		if reconfigure {
			if cfg.DbPassword, err = promptSecretOrKeep("Database password", cfg.DbPassword); err != nil {
				return err
			}
		} else {
			if cfg.DbPassword, err = promptSecret("Database password"); err != nil {
				return err
			}
		}
		if cfg.DbHost, err = promptHost("Database host (IP or hostname)", cfg.DbHost); err != nil {
			return err
		}
		if cfg.DbPort, err = promptPort("Database port", cfg.DbPort); err != nil {
			return err
		}

		fmt.Println("Testing database connection...")
		if pingErr := services.TestDBConnection(); pingErr == nil {
			fmt.Println("✓ Database connection succeeded")
			break
		} else {
			fmt.Printf("✗ Database connection failed: %v\n", pingErr)
			retry, _ := promptYesNo("Re-enter database credentials?", true)
			if !retry {
				return fmt.Errorf("aborted: database connection could not be verified")
			}
		}
	}

	// --- Application ---
	if cfg.AppHost, err = promptHost("Application host", cfg.AppHost); err != nil {
		return err
	}
	if cfg.AppPort, err = promptPort("Application port", cfg.AppPort); err != nil {
		return err
	}

	// --- Redis (optional) ---
	useRedis, err := promptYesNo("Enable Redis? (optional — used for caching and pub/sub)", cfg.RedisEnabled)
	if err != nil {
		return err
	}
	cfg.RedisEnabled = useRedis
	if useRedis {
		if cfg.RedisHost, err = promptHost("Redis host", cfg.RedisHost); err != nil {
			return err
		}
		if cfg.RedisPort, err = promptPort("Redis port", cfg.RedisPort); err != nil {
			return err
		}
		if cfg.RedisPassword, err = promptOptional("Redis password (leave blank if none)", cfg.RedisPassword); err != nil {
			return err
		}
		if cfg.RedisDB, err = promptIntRange("Redis DB number", 0, 15, cfg.RedisDB); err != nil {
			return err
		}
	}

	return nil
}

// writeConfig persists globals.AppConfig to config.json with 0600 permissions.
func writeConfig() error {
	bytes, err := json.MarshalIndent(globals.AppConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := os.WriteFile("config.json", bytes, 0600); err != nil {
		return fmt.Errorf("failed to write config.json: %w", err)
	}
	fmt.Println("Configuration saved to config.json (mode 0600)")
	return nil
}

func validateNonEmpty(label string) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("%s cannot be empty", label)
		}
		return nil
	}
}

func validateHost(label string) func(string) error {
	return func(s string) error {
		s = strings.TrimSpace(s)
		if s == "" {
			return fmt.Errorf("%s cannot be empty", label)
		}
		if strings.ContainsAny(s, " \t") {
			return fmt.Errorf("%s cannot contain whitespace", label)
		}
		if strings.Contains(s, "://") {
			return fmt.Errorf("%s should not include a scheme (e.g. http://)", label)
		}
		return nil
	}
}

func validatePort(label string) func(string) error {
	return func(s string) error {
		n, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil || n < 1 || n > 65535 {
			return fmt.Errorf("%s must be an integer between 1 and 65535", label)
		}
		return nil
	}
}

func validateIntRange(label string, min, max int) func(string) error {
	return func(s string) error {
		n, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil || n < min || n > max {
			return fmt.Errorf("%s must be an integer between %d and %d", label, min, max)
		}
		return nil
	}
}

func promptRequired(label, def string) (string, error) {
	p := promptui.Prompt{Label: label, Default: def, Validate: validateNonEmpty(label)}
	return p.Run()
}

func promptHost(label, def string) (string, error) {
	p := promptui.Prompt{Label: label, Default: def, Validate: validateHost(label)}
	return p.Run()
}

func promptSecret(label string) (string, error) {
	p := promptui.Prompt{
		Label:    label,
		Mask:     '*',
		Validate: validateNonEmpty(label),
	}
	return p.Run()
}

// promptSecretOrKeep returns current unchanged if the user enters a blank value.
// Used by reconfigure so users can keep an existing password without retyping it.
func promptSecretOrKeep(label, current string) (string, error) {
	p := promptui.Prompt{
		Label: label + " (leave blank to keep current)",
		Mask:  '*',
	}
	v, err := p.Run()
	if err != nil {
		return "", err
	}
	if v == "" {
		return current, nil
	}
	return v, nil
}

func promptOptional(label, def string) (string, error) {
	p := promptui.Prompt{Label: label, Default: def}
	return p.Run()
}

func promptPort(label, def string) (string, error) {
	p := promptui.Prompt{Label: label, Default: def, Validate: validatePort(label)}
	return p.Run()
}

func promptIntRange(label string, min, max, def int) (int, error) {
	p := promptui.Prompt{
		Label:    label,
		Default:  strconv.Itoa(def),
		Validate: validateIntRange(label, min, max),
	}
	s, err := p.Run()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(s))
}

func promptYesNo(label string, defYes bool) (bool, error) {
	items := []string{"No", "Yes"}
	s := promptui.Select{Label: label, Items: items}
	if defYes {
		s.CursorPos = 1
	}
	_, res, err := s.Run()
	if err != nil {
		return false, err
	}
	return res == "Yes", nil
}

func promptSelect(label string, items []string, def string) (string, error) {
	s := promptui.Select{Label: label, Items: items}
	for i, it := range items {
		if it == def {
			s.CursorPos = i
			break
		}
	}
	_, res, err := s.Run()
	return res, err
}

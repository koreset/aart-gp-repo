package main

import (
	"api/config"
	"api/docs"
	"api/globals"
	"api/log"
	"api/routes"
	"api/services"
	"api/services/bav"
	"api/services/bav/audit"
	"api/services/bav/providers"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	"github.com/manifoldco/promptui"
)

type application struct{}

func (app *application) Start(s service.Service) error {
	go app.run(s)
	return nil
}

func (app *application) run(s service.Service) {
	startApplication(false, s)
}

func (app *application) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	globals.Logger.Info("We are stopping the service... question is how do we install?")
	return nil
}

func (app *application) Install(s service.Service) error {
	s.Install()
	return nil
}

// @title ADSolutions API
// @version 1.0
// @description ADS API.

// @host localhost:9090
// @BasePath /api/v1/
// Entry point for the Application
func main() {
	svcFlag := flag.String("service", "", "Control the system service.")

	flag.Parse()

	svcConfig := &service.Config{
		Name:        "AartApiServer",
		DisplayName: "AART API Server",
		Description: "Backend Server for AART",
	}

	app := &application{}
	s, err := service.New(app, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	globals.Logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	_ = globals.Logger.Info(os.Getwd())

	if *svcFlag == "uninstall" {
		//Service Name
		promptServiceName := promptui.Prompt{
			Label: "Enter a service name to uninstall",
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("service name is required")
				}
				return nil
			},
		}

		serviceName, err := promptServiceName.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		svcConfig.Name = serviceName

		_ = globals.Logger.Info("Preparing to uninstall service")
		err = s.Stop()
		if err != nil {
			_ = globals.Logger.Error(err)
		}
		err = s.Uninstall()
		_ = globals.Logger.Error(err)
	}

	if *svcFlag == "install" {
		if runtime.GOOS != "windows" {
			fmt.Println("-service install is only supported on Windows. " +
				"On macOS/Linux, run the binary directly or wrap it in a systemd/launchd unit.")
			return
		}

		serviceName, err := promptRequired("Enter a service name for installation", "")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		svcConfig.Name = serviceName

		displayName, err := promptRequired("Enter a display name for the service", "")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		svcConfig.DisplayName = displayName

		if err := s.Install(); err != nil {
			log.Fatal(err)
		}
		if err := s.Start(); err != nil {
			fmt.Printf("Service installed but failed to start: %v\n", err)
		}
	}

	if *svcFlag == "setup" {
		if _, statErr := os.Stat("config.json"); statErr == nil {
			fmt.Println("config.json already exists. Use -service reconfigure to change it.")
			return
		}

		confirm, err := promptYesNo(
			"Installing a new instance of AART. Have your database credentials ready. Continue?", true)
		if err != nil || !confirm {
			fmt.Println("Setup aborted.")
			return
		}

		if err := runConfigWizard(false); err != nil {
			fmt.Println(err)
			return
		}

		if err := writeConfig(); err != nil {
			fmt.Println(err)
			return
		}

		services.SetupTables(true, true)
		fmt.Println("Running Application")
		if err := s.Run(); err != nil {
			fmt.Println(err)
		}
	}

	if *svcFlag == "reconfigure" {
		configBytes, err := os.ReadFile("config.json")
		if err != nil {
			fmt.Printf("Could not read config.json: %v\n", err)
			return
		}
		if err := json.Unmarshal(configBytes, &globals.AppConfig); err != nil {
			fmt.Printf("Could not parse config.json: %v\n", err)
			return
		}

		fmt.Println("Reconfiguring AART. Current values are shown as defaults; press enter to keep.")
		if err := runConfigWizard(true); err != nil {
			fmt.Println(err)
			return
		}

		if err := writeConfig(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Configuration updated. Restart the service to apply changes.")
	}

	if *svcFlag == "" {
		fmt.Println("Running Application")
		err := s.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func startApplication(initTables bool, s service.Service) {
	log.Info("Starting application initialization")

	// Get configuration path
	configPath, err := getConfigPath()
	if runtime.GOOS != "windows" {
		configPath = "config.json"
		log.Info("Using default config.json path for non-Windows platform")
	} else {
		log.WithField("config_path", configPath).Info("Using Windows-specific config path")
	}

	// Read configuration file
	log.WithField("config_path", configPath).Info("Reading configuration file")
	fileByte, err := os.ReadFile(configPath)
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed to read configuration file")
		fmt.Println(err)
		return
	}

	// Parse configuration
	log.Debug("Parsing configuration file")
	err = json.Unmarshal(fileByte, &globals.AppConfig)
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed to parse configuration file")
		fmt.Println(err)
		return
	}

	// Log configuration (excluding sensitive data)
	log.WithFields(map[string]interface{}{
		"app_host":       globals.AppConfig.AppHost,
		"app_port":       globals.AppConfig.AppPort,
		"db_type":        globals.AppConfig.DbType,
		"db_host":        globals.AppConfig.DbHost,
		"db_name":        globals.AppConfig.DbName,
		"redis_host":     globals.AppConfig.RedisHost,
		"redis_port":     globals.AppConfig.RedisPort,
		"redis_password": globals.AppConfig.RedisPassword,
		"redis_db":       globals.AppConfig.RedisDB,
		"redis_enabled":  globals.AppConfig.RedisEnabled,
	}).Info("Application configuration loaded")

	// Initialize Redis (optional)
	services.InitRedis()

	// Wire the Bank Account Verification provider registry from env-var config.
	bavRegistry, err := providers.NewRegistry(providers.Config{
		Provider:          config.BAVProvider,
		APIKey:            config.BAVAPIKey,
		BaseURL:           config.BAVBaseURL,
		Mode:              config.VerifyNowMode,
		OAuthClientID:     config.BAVOAuthClientID,
		OAuthClientSecret: config.BAVOAuthClientSecret,
		OAuthTokenURL:     config.BAVOAuthTokenURL,
		Timeout:           time.Duration(config.BAVTimeoutSeconds) * time.Second,
		MockAsync:         config.MockBAVAsync,
	})
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed to build BAV provider registry")
		fmt.Println(err)
		return
	}
	bavRegistry.WithLogger(audit.NewGormLogger(services.DB))
	bav.SetDefault(bavRegistry)
	log.WithField("provider", config.BAVProvider).Info("Bank Account Verification provider wired")

	// Initialize database and tables
	log.Info("Setting up database tables")
	services.SetupTables(initTables, true)

	// Run database migrations
	log.Info("Running database migrations")
	services.RunMigrationsOnStartup()

	// Ensure gp_table_stats exists and is seeded (idempotent, fast).
	// (Table creation is handled by SetupTables above; this seeds the row counts.)
	services.EnsureGPTableStats()

	// Top up group_benefit_mappers with any newly-added base rows (e.g.
	// AAGLA, AGLC) so Benefits Customization shows them on existing
	// installs without requiring a manual seed.
	if err := services.EnsureBaseBenefitMapsSeeded(); err != nil {
		log.WithField("error", err.Error()).Warn("Failed to seed base benefit maps")
	}

	services.StartGroupSchemeStatusUpdater()
	services.StartNotificationOverdueSweeper()
	services.StartDeadlineOverdueSweeper()
	services.StartBordereauxFileRetentionSweeper()
	services.StartWinProbabilityRetrainingJob()

	// Initialize WebSocket hub for real-time communications
	services.InitWSHub()
	services.StartRedisWSSubscriber()

	// Start the calculation job queue worker (bounded concurrency)
	services.StartCalculationJobWorker()

	// Initialize service logger
	globals.Logger.Info("We are starting up...")
	globals.Logger.Info(globals.AppConfig)

	// Initialize application logger
	log.Info("Initializing application logger")
	log.InitLogger()

	// Create and configure Gin router
	log.Info("Creating and configuring HTTP router")
	router := gin.Default()

	// Configure CORS
	log.Info("Configuring CORS")
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowWildcard = true
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"*"}
	router.Use(cors.New(corsConfig))

	// Configure routes
	log.Info("Configuring API routes")
	routes.ConfigureRouter(router)

	// Get port from config or use default
	port := os.Getenv("APP_PORT")

	if len(port) == 0 {
		// Fallback to config value
		port = globals.AppConfig.AppPort
	}

	if len(strings.Trim(port, "")) == 0 {
		port = "9090"
		log.WithField("port", port).Info("Using default port")
	} else {
		log.WithField("port", port).Info("Using configured port")
	}

	// Configure Swagger
	docs.SwaggerInfo.Host = globals.AppConfig.AppHost + ":" + globals.AppConfig.AppPort
	log.WithField("swagger_host", docs.SwaggerInfo.Host).Info("Configured Swagger host")

	// Start HTTP server
	environment := os.Getenv("ENVIRONMENT")
	log.WithField("environment", environment).Info("Starting HTTP server")

	if environment == "production" {
		log.WithField("host", globals.AppConfig.AppHost).Info("Starting production server with AutoTLS")
		log.Fatal(autotls.Run(router, globals.AppConfig.AppHost))
	} else {
		serverAddress := fmt.Sprintf("0.0.0.0:%s", port)
		log.WithField("address", serverAddress).Info("Starting development server")
		log.Fatal(router.Run(serverAddress))
	}
}

func getConfigPath() (string, error) {
	fullexecpath, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir, _ := filepath.Split(fullexecpath)
	return filepath.Join(dir, "config.json"), nil
}

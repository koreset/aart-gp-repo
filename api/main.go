package main

import (
	"api/docs"
	"api/globals"
	"api/log"
	"api/routes"
	"api/services"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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
		if runtime.GOOS == "windows" {
			//We need to configure and install the service
			//Service Name
			promptServiceName := promptui.Prompt{
				Label: "Enter a service name for installation",
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

			//Display Name
			promptDisplayName := promptui.Prompt{
				Label: "Enter a service name for installation",
				Validate: func(input string) error {
					if input == "" {
						return fmt.Errorf("service name is required")
					}
					return nil
				},
			}

			displayName, err := promptDisplayName.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			svcConfig.DisplayName = displayName

			err = s.Install()
			if err != nil {
				log.Fatal(err)
			}
			err = s.Start()

		}
	}

	if *svcFlag == "setup" {
		_, err := os.Open("config.json")
		if err != nil { //config.json does not exist, clean setup
			//				// Start off installation process and ask for comfirmation from user
			install := promptui.Select{
				Label: "Installing a new instance of AART. " +
					"Please have the database credentials handy as this will be required. " +
					"Continue?",
				Items: []string{"Yes", "No"},
			}

			_, res, err := install.Run()

			if err != nil {
				_ = globals.Logger.Error(err)

			}

			if res == "Yes" {
				//DB Type
				promptDb := promptui.Select{
					Label: "Select the database type to install",
					Items: []string{"MySQL", "PostgreSQL"},
				}

				_, result, err := promptDb.Run()
				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					return
				}
				globals.AppConfig.DbType = strings.ToLower(result)

				//DB Name
				promptDbName := promptui.Prompt{
					Label: "Enter the database name",
					Validate: func(input string) error {
						if input == "" {
							return fmt.Errorf("database name can not be empty")
						}
						return nil
					},
				}

				dbName, err := promptDbName.Run()
				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					return
				}
				globals.AppConfig.DbName = dbName

				//DB UserEmail
				promptDbUser := promptui.Prompt{
					Label: "Enter the database user",
					Validate: func(input string) error {
						if input == "" {
							return fmt.Errorf("database user can not be empty")
						}
						return nil
					},
				}

				dbUser, err := promptDbUser.Run()
				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					return
				}
				globals.AppConfig.DbUser = dbUser

				//DB Password
				promptDbPassword := promptui.Prompt{
					Label: "Enter the database password",
					Validate: func(input string) error {
						if input == "" {
							return fmt.Errorf("database password can not be empty")
						}
						return nil
					},
					Mask: '*',
				}

				dbPassword, err := promptDbPassword.Run()
				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					return
				}
				globals.AppConfig.DbPassword = dbPassword

				//DB Host
				promptDbHost := promptui.Prompt{
					Label: "Enter the database host. This should be a valid IP address or accessible hostname",
					Validate: func(input string) error {
						if input == "" {
							return fmt.Errorf("database host can not be empty")
						}
						return nil
					},
				}

				dbHost, err := promptDbHost.Run()
				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					return
				}
				globals.AppConfig.DbHost = dbHost

				//DB Port
				promptDbPort := promptui.Prompt{
					Label: "Enter the database port",
					Validate: func(input string) error {
						if input == "" {
							return fmt.Errorf("database port can not be empty")
						}
						return nil
					},
				}

				dbPort, err := promptDbPort.Run()
				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					return
				}
				globals.AppConfig.DbPort = dbPort

				//App Port
				promptAppPort := promptui.Prompt{
					Label: "Enter the port number the application will be running on",
					Validate: func(input string) error {
						if input == "" {
							return fmt.Errorf("Application port can not be empty")
						}
						return nil
					},
				}

				appPort, err := promptAppPort.Run()
				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					return
				}
				globals.AppConfig.AppPort = appPort

				//App Host
				promptAppHost := promptui.Prompt{
					Label: "Enter the hostname the application will be running on",
					Validate: func(input string) error {
						if input == "" {
							return fmt.Errorf("Application host can not be empty")
						}
						return nil
					},
				}

				appHost, err := promptAppHost.Run()
				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					return
				}
				globals.AppConfig.AppHost = appHost

				bytes, _ := json.MarshalIndent(globals.AppConfig, "", "  ")
				f, err := os.Create("config.json")
				if err != nil {
					fmt.Println(err)
				}
				_, err = f.Write(bytes)

				if err != nil {
					fmt.Println(err)
				}
				_ = f.Close()
				fmt.Println("Configuration saved to config.json")
				services.SetupTables(true, true)
				fmt.Println("Running Application")
				err = s.Run()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
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

	// Initialize database and tables
	log.Info("Setting up database tables")
	services.SetupTables(initTables, true)

	// Run database migrations
	log.Info("Running database migrations")
	services.RunMigrationsOnStartup()

	// Ensure gp_table_stats exists and is seeded (idempotent, fast).
	// (Table creation is handled by SetupTables above; this seeds the row counts.)
	services.EnsureGPTableStats()

	services.StartGroupSchemeStatusUpdater()
	services.StartNotificationOverdueSweeper()
	services.StartWinProbabilityRetrainingJob()

	// Initialize WebSocket hub for real-time communications
	services.InitWSHub()
	services.StartRedisWSSubscriber()

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

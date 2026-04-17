package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	DbHost     string
	DbUser     string
	DbPort     string
	DbPassword string
	DbName     string

	CheckIDApiKey   string
	VerifyNowApiKey string
	VerifyNowMode   string

	// BAV (Bank Account Verification) — provider-agnostic configuration.
	// BAVProvider selects the active adapter ("verifynow", "lexisnexis", ...).
	// All other BAV* vars are consumed by adapters that need them.
	BAVProvider          string
	BAVAPIKey            string
	BAVBaseURL           string
	BAVOAuthClientID     string
	BAVOAuthClientSecret string
	BAVOAuthTokenURL     string
	BAVTimeoutSeconds    int
	// MockBAVAsync flips the mock provider into async mode. Only consulted
	// when BAVProvider == "mock"; wired via MOCK_BAV_ASYNC.
	MockBAVAsync bool

	// BordereauxFileRetentionDays is the age (in days) after which generated
	// bordereaux output files and terminal-status inbound confirmation/submission
	// uploads are deleted from disk by the retention sweeper. DB rows are kept;
	// only the on-disk artefact is removed. Configured via
	// BORDEREAUX_FILE_RETENTION_DAYS, default 90.
	BordereauxFileRetentionDays int
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	DbHost = os.Getenv("DB_HOST")
	DbUser = os.Getenv("DB_USER")
	DbPort = os.Getenv("DB_PORT")
	DbPassword = os.Getenv("DB_PWD")
	DbName = os.Getenv("DB_NAME")

	CheckIDApiKey = os.Getenv("CHECKID_API_KEY")
	VerifyNowApiKey = os.Getenv("VERIFYNOW_API_KEY")
	VerifyNowMode = os.Getenv("VERIFYNOW_MODE")
	if VerifyNowMode == "" {
		VerifyNowMode = "production"
	}

	BAVProvider = os.Getenv("BAV_PROVIDER")
	BAVAPIKey = os.Getenv("BAV_API_KEY")
	BAVBaseURL = os.Getenv("BAV_BASE_URL")
	BAVOAuthClientID = os.Getenv("BAV_OAUTH_CLIENT_ID")
	BAVOAuthClientSecret = os.Getenv("BAV_OAUTH_CLIENT_SECRET")
	BAVOAuthTokenURL = os.Getenv("BAV_OAUTH_TOKEN_URL")
	if n, err := strconv.Atoi(os.Getenv("BAV_TIMEOUT_SECONDS")); err == nil && n > 0 {
		BAVTimeoutSeconds = n
	} else {
		BAVTimeoutSeconds = 45
	}

	if BAVAPIKey == "" {
		BAVAPIKey = VerifyNowApiKey
	}
	if BAVProvider == "" {
		BAVProvider = "verifynow"
	}

	switch strings.ToLower(os.Getenv("MOCK_BAV_ASYNC")) {
	case "1", "true", "yes", "on":
		MockBAVAsync = true
	}

	if n, err := strconv.Atoi(os.Getenv("BORDEREAUX_FILE_RETENTION_DAYS")); err == nil && n > 0 {
		BordereauxFileRetentionDays = n
	} else {
		BordereauxFileRetentionDays = 90
	}
}

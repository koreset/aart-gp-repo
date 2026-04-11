package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DbHost     string
	DbUser     string
	DbPort     string
	DbPassword string
	DbName     string

	CheckIDApiKey    string
	VerifyNowApiKey  string
	VerifyNowMode    string
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
}




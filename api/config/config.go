package config

import "os"

var (
	DbHost     = os.Getenv("DB_HOST")
	DbUser     = os.Getenv("DB_USER")
	DbPort     = os.Getenv("DB_PORT")
	DbPassword = os.Getenv("DB_PWD")
	DbName     = os.Getenv("DB_NAME")
)




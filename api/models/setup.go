package models

type AppConfig struct {
	DbType       string `json:"db_type"`
	DbName       string `json:"db_name"`
	DbHost       string `json:"db_host"`
	DbUser       string `json:"db_user"`
	DbPassword   string `json:"db_password"`
	DbPort       string `json:"db_port"`
	AppPort      string `json:"app_port"`
	AppHost      string `json:"app_host"`

	// Redis configuration
	RedisEnabled bool   `json:"redis_enabled"`
	RedisHost    string `json:"redis_host"`
	RedisPort    string `json:"redis_port"`
	RedisPassword string `json:"redis_password"`
	RedisDB      int    `json:"redis_db"`
}

package configs

import (
	"os"

	"github.com/joho/godotenv"
)

// Configuration is the project configuration definition
type Configuration struct {
	Application AppConfig
	Database    DBConfig
	Redis       RedisConfig
}

var cfg *Configuration

// Load func loads configs/envvars
func Load() (config Configuration, err error) {
	if os.Getenv("APP_ENV") == "development" {
		dotenvErr := godotenv.Load()
		if dotenvErr != nil {
			panic(dotenvErr)
		}
	}
	// Load Application Configs
	appConfig, _ := loadAppConfig()
	config.Application = appConfig

	// Load Database Configs
	dbConfig, _ := loadDBConfig()
	config.Database = dbConfig

	redisConfig, _ := loadRedisConfig()
	config.Redis = redisConfig

	cfg = &config

	return config, nil
}

// Get func gets config values
func Get() (c *Configuration) {
	return cfg
}

package configs

import (
	"time"

	"github.com/spf13/viper"
)

// RedisConfig holds the redis configuration definition
type RedisConfig struct {
	Host        string
	UserDataTTL time.Duration
}

func loadRedisConfig() (config RedisConfig, err error) {
	provider := viper.New()
	provider.SetEnvPrefix("REDIS")
	provider.AutomaticEnv()
	setDefaultRedisConfig(provider)

	err = provider.Unmarshal(&config)

	return config, err
}

func setDefaultRedisConfig(provider *viper.Viper) {
	provider.SetDefault("Host", "localhost:6379")
	provider.SetDefault("UserDataTTL", 24*time.Hour) // 24h
}

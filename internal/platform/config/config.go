package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv        string
	Port          string
	DatabaseURL   string
	RedisURL      string
	JWTPrivateKey string
	JWTPublicKey  string
}

func Load() *Config {
	viper.AutomaticEnv()

	cfg := &Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getRequired("DATABASE_URL"),
		RedisURL:      getRequired("REDIS_URL"),
		JWTPrivateKey: getRequired("JWT_PRIVATE_KEY"),
		JWTPublicKey:  getRequired("JWT_PUBLIC_KEY"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if viper.GetString(key) == "" {
		return fallback
	}
	return viper.GetString(key)
}

func getRequired(key string) string {
	val := viper.GetString(key)
	if val == "" {
		log.Fatalf("Missing required env variable: %s", key)
	}
	return val
}

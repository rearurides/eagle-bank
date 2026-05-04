package config

import (
	"os"
	"time"
)

type Config struct {
	Port          string
	DBPath        string
	JWTSecret     string
	JWTExpiration time.Duration
}

func LoadConfig() *Config {
	return &Config{
		Port:          getEnv("PORT", "3000"),
		DBPath:        getEnv("DB_PATH", "eagle_bank.db"),
		JWTSecret:     getEnv("JWT_SECRET", "default_jwt_secret_change_me"),
		JWTExpiration: 24 * time.Hour,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}

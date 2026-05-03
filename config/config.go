package config

import "os"

type Config struct {
	Port   string
	DBPath string
}

func LoadConfig() *Config {
	return &Config{
		Port:   getEnv("PORT", "3000"),
		DBPath: getEnv("DB_PATH", "eagle_bank.db"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}

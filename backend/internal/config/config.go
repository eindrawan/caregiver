package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Environment string
	Port        string
	DatabaseURL string
	LogLevel    string
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "caregiver_shift_tracker.db"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

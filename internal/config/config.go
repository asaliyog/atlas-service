package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Environment   string
	DatabaseURL   string
	AzureTenantID string
	AzureClientID string
	JWTSecret     string
	Port          string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Environment:   getEnv("ENVIRONMENT", "development"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/golang_service?sslmode=disable"),
		AzureTenantID: getEnv("AZURE_TENANT_ID", ""),
		AzureClientID: getEnv("AZURE_CLIENT_ID", ""),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		Port:          getEnv("PORT", "8080"),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
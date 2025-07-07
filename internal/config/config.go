package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Environment        string
	DatabaseURL        string
	AzureTenantID      string
	AzureClientID      string
	AzureClientSecret  string
	JWTSecret          string
	Port               string
	BypassAuth         bool
	AzureAuthScope     string
	AzureTokenEndpoint string
}

// Load loads configuration from environment variables
func Load() *Config {
	env := getEnv("ENVIRONMENT", "development")
	
	return &Config{
		Environment:        env,
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/golang_service?sslmode=disable"),
		AzureTenantID:      getEnv("AZURE_TENANT_ID", ""),
		AzureClientID:      getEnv("AZURE_CLIENT_ID", ""),
		AzureClientSecret:  getEnv("AZURE_CLIENT_SECRET", ""),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		Port:               getEnv("PORT", "8080"),
		BypassAuth:         env == "development" || env == "local",
		AzureAuthScope:     getEnv("AZURE_AUTH_SCOPE", "https://graph.microsoft.com/.default"),
		AzureTokenEndpoint: getEnv("AZURE_TOKEN_ENDPOINT", ""),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
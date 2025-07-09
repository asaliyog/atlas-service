package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	DatabaseURL string
	Environment string
	Port        string
	// Auth configuration
	AzureTenantID      string
	AzureClientID      string
	AzureClientSecret  string
	JWTSecret          string
	BypassAuth         bool
	AzureAuthScope     string
	AzureTokenEndpoint string
	// Environment resolution configuration
	EnableEnvironmentResolution bool
	EnvironmentResolutionConfig map[string]bool // API endpoint -> enable/disable
}

// Load loads configuration from environment variables
func Load() *Config {
	env := getEnv("ENVIRONMENT", "development")
	
	return &Config{
		DatabaseURL:                 getEnv("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=atlas_service port=5432 sslmode=disable"),
		Environment:                 env,
		Port:                        getEnv("PORT", "8080"),
		AzureTenantID:               getEnv("AZURE_TENANT_ID", ""),
		AzureClientID:               getEnv("AZURE_CLIENT_ID", ""),
		AzureClientSecret:           getEnv("AZURE_CLIENT_SECRET", ""),
		JWTSecret:                   getEnv("JWT_SECRET", "your-secret-key"),
		BypassAuth:                  env == "development" || env == "local",
		AzureAuthScope:              getEnv("AZURE_AUTH_SCOPE", "https://graph.microsoft.com/.default"),
		AzureTokenEndpoint:          getEnv("AZURE_TOKEN_ENDPOINT", ""),
		EnableEnvironmentResolution: getEnvBool("ENABLE_ENVIRONMENT_RESOLUTION", true),
		EnvironmentResolutionConfig: map[string]bool{
			"/api/v1/vms":          getEnvBool("ENV_RESOLUTION_VMS", true),
			"/api/v1/environments": getEnvBool("ENV_RESOLUTION_ENVIRONMENTS", false),
			"/api/v1/users":        getEnvBool("ENV_RESOLUTION_USERS", false),
		},
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool gets an environment variable and returns a boolean
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if value == "true" {
			return true
		}
		if value == "false" {
			return false
		}
	}
	return defaultValue
}
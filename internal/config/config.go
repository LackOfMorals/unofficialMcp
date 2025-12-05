package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	URI          string
	ClientId     string
	ClientSecret string
	ReadOnly     string // If true, disables tools that make changes to Aura infrastructure
}

// Validate validates the configuration and returns an error if invalid
func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("configuration is required but was nil")
	}

	validations := []struct {
		value string
		name  string
	}{
		{c.URI, "Neo4j URI"},
		{c.ClientId, "Aura API Client Id "},
		{c.ClientSecret, "Aura API Client Secret"},
	}

	for _, v := range validations {
		if v.value == "" {
			return fmt.Errorf("%s is required but was empty", v.name)
		}
	}

	return nil
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() (*Config, error) {
	cfg := &Config{
		URI:          GetEnvWithDefault("AURA_API_URI", "https://api.neo4j.io/v1"),
		ClientId:     GetEnvWithDefault("AURA_API_ID", ""),
		ClientSecret: GetEnvWithDefault("AURA_API_SECRET", ""),
		ReadOnly:     GetEnvWithDefault("AURA_API_READ_ONLY", "true"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

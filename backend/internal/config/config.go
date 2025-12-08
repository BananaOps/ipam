package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

// ServerConfig contains server-related configuration
type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

// DatabaseConfig contains database-related configuration
type DatabaseConfig struct {
	Type             string `yaml:"type"`              // "sqlite" or "mongodb"
	Path             string `yaml:"path"`              // For SQLite
	ConnectionString string `yaml:"connection_string"` // For MongoDB
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	// Read config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() *Config {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Type:             getEnv("DATABASE_TYPE", "sqlite"),
			Path:             getEnv("DATABASE_PATH", "./data/ipam.db"),
			ConnectionString: getEnv("DATABASE_CONNECTION_STRING", ""),
		},
	}

	return config
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate database type
	if c.Database.Type != "sqlite" && c.Database.Type != "mongodb" {
		return fmt.Errorf("invalid database type: %s (must be 'sqlite' or 'mongodb')", c.Database.Type)
	}

	// Validate database-specific configuration
	if c.Database.Type == "sqlite" && c.Database.Path == "" {
		return fmt.Errorf("database path is required for SQLite")
	}

	if c.Database.Type == "mongodb" && c.Database.ConnectionString == "" {
		return fmt.Errorf("connection string is required for MongoDB")
	}

	return nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

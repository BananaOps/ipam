package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bananaops/ipam-bananaops/internal/config"
	"github.com/bananaops/ipam-bananaops/internal/repository"
)

func main() {
	fmt.Println("IPAM by BananaOps - Server")
	log.Println("Server starting...")

	// Load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded: database type=%s", cfg.Database.Type)

	// Initialize database
	repo, err := repository.NewRepository(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer repo.Close()

	log.Printf("Database initialized successfully (%s)", cfg.Database.Type)

	// TODO: Initialize go-ipam
	// TODO: Initialize service layer
	// TODO: Initialize REST gateway
	// TODO: Start HTTP server

	log.Println("Server ready (database layer implemented)")
}

// loadConfiguration loads configuration from file or environment
func loadConfiguration() (*config.Config, error) {
	// Try to load from config file first
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); err == nil {
		log.Printf("Loading configuration from file: %s", configPath)
		return config.LoadConfig(configPath)
	}

	// Fall back to environment variables
	log.Println("Loading configuration from environment variables")
	cfg := config.LoadConfigFromEnv()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

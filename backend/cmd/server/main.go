package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bananaops/ipam-bananaops/internal/cloudprovider"
	"github.com/bananaops/ipam-bananaops/internal/config"
	"github.com/bananaops/ipam-bananaops/internal/gateway"
	"github.com/bananaops/ipam-bananaops/internal/repository"
	"github.com/bananaops/ipam-bananaops/internal/service"
)

func main() {
	fmt.Println("IPAM by BananaOps - Server")
	log.Println("Server starting...")

	ctx := context.Background()

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

	// Initialize IP service
	ipService := service.NewGoIPAMService()
	log.Println("IP service initialized")

	// Initialize cloud provider manager
	cloudManager := cloudprovider.NewManager(cfg, repo)
	log.Println("Cloud provider manager initialized")

	// Start cloud provider manager
	if err := cloudManager.Start(ctx); err != nil {
		log.Printf("Failed to start cloud provider manager: %v", err)
		// Continue without cloud providers if they fail to start
	}
	defer cloudManager.Stop()

	// Initialize service layer
	serviceLayer := service.NewServiceLayer(repo, ipService, cloudManager)
	log.Println("Service layer initialized")

	// Initialize REST gateway with cloud manager
	gatewayHandler := gateway.NewGateway(serviceLayer, cloudManager)
	log.Println("REST gateway initialized")

	// Start HTTP server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting HTTP server on %s", serverAddr)

	if err := http.ListenAndServe(serverAddr, gatewayHandler.Handler()); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
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

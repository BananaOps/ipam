package repository

import (
	"fmt"

	"github.com/bananaops/ipam-bananaops/internal/config"
)

// NewRepository creates a new repository based on the configuration
func NewRepository(cfg *config.DatabaseConfig) (SubnetRepository, error) {
	switch cfg.Type {
	case "sqlite":
		return NewSQLiteRepository(cfg.Path)
	case "mongodb":
		return NewMongoDBRepository(cfg.ConnectionString)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}
}

package cloudprovider

import (
	"context"
	"errors"
)

// Common errors for cloud provider operations
var (
	ErrProviderNotFound     = errors.New("cloud provider not found")
	ErrProviderUnavailable  = errors.New("cloud provider unavailable")
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrRateLimited          = errors.New("rate limited by provider")
	ErrInvalidCredentials   = errors.New("invalid credentials")
)

// CloudProviderType represents the type of cloud provider
type CloudProviderType string

const (
	ProviderAWS      CloudProviderType = "aws"
	ProviderAzure    CloudProviderType = "azure"
	ProviderGCP      CloudProviderType = "gcp"
	ProviderScaleway CloudProviderType = "scaleway"
	ProviderOVH      CloudProviderType = "ovh"
)

// CloudCredentials contains authentication information for cloud providers
type CloudCredentials struct {
	Provider  CloudProviderType
	AccessKey string
	SecretKey string
	Token     string
	Region    string
	// Additional provider-specific fields can be added as needed
	Extra map[string]string
}

// CloudSubnet represents a subnet fetched from a cloud provider
type CloudSubnet struct {
	CIDR      string
	Name      string
	Region    string
	AccountID string
	VPCId     string
	Tags      map[string]string
}

// CloudProvider defines the interface that all cloud provider implementations must satisfy
type CloudProvider interface {
	// GetName returns the name of the cloud provider
	GetName() string

	// GetType returns the type of the cloud provider
	GetType() CloudProviderType

	// FetchSubnets retrieves all subnets from the cloud provider
	FetchSubnets(ctx context.Context, credentials CloudCredentials) ([]*CloudSubnet, error)

	// GetRegions returns the list of available regions for this provider
	GetRegions() []string

	// ValidateCredentials checks if the provided credentials are valid
	ValidateCredentials(ctx context.Context, credentials CloudCredentials) error
}

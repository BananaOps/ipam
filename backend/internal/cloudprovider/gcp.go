package cloudprovider

import (
	"context"
	"fmt"
)

// GCPProvider implements the CloudProvider interface for Google Cloud Platform
type GCPProvider struct {
	name string
}

// NewGCPProvider creates a new GCP cloud provider instance
func NewGCPProvider() *GCPProvider {
	return &GCPProvider{
		name: "Google Cloud Platform",
	}
}

// GetName returns the name of the cloud provider
func (p *GCPProvider) GetName() string {
	return p.name
}

// GetType returns the type of the cloud provider
func (p *GCPProvider) GetType() CloudProviderType {
	return ProviderGCP
}

// FetchSubnets retrieves all subnets from GCP
// This is a stub implementation - actual GCP SDK integration will be added in the future
func (p *GCPProvider) FetchSubnets(ctx context.Context, credentials CloudCredentials) ([]*CloudSubnet, error) {
	// Validate credentials
	if err := p.ValidateCredentials(ctx, credentials); err != nil {
		return nil, err
	}

	// TODO: Implement actual GCP SDK integration
	// For now, return an error indicating the feature is not yet implemented
	return nil, fmt.Errorf("%w: GCP subnet fetching not yet implemented", ErrProviderUnavailable)
}

// GetRegions returns the list of available GCP regions
func (p *GCPProvider) GetRegions() []string {
	return []string{
		"us-central1",
		"us-east1",
		"us-east4",
		"us-west1",
		"us-west2",
		"us-west3",
		"us-west4",
		"europe-west1",
		"europe-west2",
		"europe-west3",
		"europe-west4",
		"europe-west6",
		"europe-north1",
		"asia-east1",
		"asia-east2",
		"asia-northeast1",
		"asia-northeast2",
		"asia-northeast3",
		"asia-south1",
		"asia-southeast1",
		"asia-southeast2",
		"australia-southeast1",
		"southamerica-east1",
	}
}

// ValidateCredentials checks if the provided GCP credentials are valid
func (p *GCPProvider) ValidateCredentials(ctx context.Context, credentials CloudCredentials) error {
	if credentials.Provider != ProviderGCP {
		return fmt.Errorf("invalid provider type: expected %s, got %s", ProviderGCP, credentials.Provider)
	}

	if credentials.Token == "" {
		return ErrInvalidCredentials
	}

	// TODO: Implement actual GCP credential validation
	return nil
}

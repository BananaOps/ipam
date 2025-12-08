package cloudprovider

import (
	"context"
	"fmt"
)

// OVHProvider implements the CloudProvider interface for OVH
type OVHProvider struct {
	name string
}

// NewOVHProvider creates a new OVH cloud provider instance
func NewOVHProvider() *OVHProvider {
	return &OVHProvider{
		name: "OVH",
	}
}

// GetName returns the name of the cloud provider
func (p *OVHProvider) GetName() string {
	return p.name
}

// GetType returns the type of the cloud provider
func (p *OVHProvider) GetType() CloudProviderType {
	return ProviderOVH
}

// FetchSubnets retrieves all subnets from OVH
// This is a stub implementation - actual OVH API integration will be added in the future
func (p *OVHProvider) FetchSubnets(ctx context.Context, credentials CloudCredentials) ([]*CloudSubnet, error) {
	// Validate credentials
	if err := p.ValidateCredentials(ctx, credentials); err != nil {
		return nil, err
	}

	// TODO: Implement actual OVH API integration
	// For now, return an error indicating the feature is not yet implemented
	return nil, fmt.Errorf("%w: OVH subnet fetching not yet implemented", ErrProviderUnavailable)
}

// GetRegions returns the list of available OVH regions
func (p *OVHProvider) GetRegions() []string {
	return []string{
		"GRA1",
		"GRA3",
		"GRA5",
		"GRA7",
		"SBG1",
		"SBG3",
		"SBG5",
		"BHS1",
		"BHS3",
		"BHS5",
		"DE1",
		"UK1",
		"WAW1",
		"SGP1",
		"SYD1",
	}
}

// ValidateCredentials checks if the provided OVH credentials are valid
func (p *OVHProvider) ValidateCredentials(ctx context.Context, credentials CloudCredentials) error {
	if credentials.Provider != ProviderOVH {
		return fmt.Errorf("invalid provider type: expected %s, got %s", ProviderOVH, credentials.Provider)
	}

	if credentials.AccessKey == "" || credentials.SecretKey == "" {
		return ErrInvalidCredentials
	}

	// TODO: Implement actual OVH credential validation
	return nil
}

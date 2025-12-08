package cloudprovider

import (
	"context"
	"fmt"
)

// ScalewayProvider implements the CloudProvider interface for Scaleway
type ScalewayProvider struct {
	name string
}

// NewScalewayProvider creates a new Scaleway cloud provider instance
func NewScalewayProvider() *ScalewayProvider {
	return &ScalewayProvider{
		name: "Scaleway",
	}
}

// GetName returns the name of the cloud provider
func (p *ScalewayProvider) GetName() string {
	return p.name
}

// GetType returns the type of the cloud provider
func (p *ScalewayProvider) GetType() CloudProviderType {
	return ProviderScaleway
}

// FetchSubnets retrieves all subnets from Scaleway
// This is a stub implementation - actual Scaleway SDK integration will be added in the future
func (p *ScalewayProvider) FetchSubnets(ctx context.Context, credentials CloudCredentials) ([]*CloudSubnet, error) {
	// Validate credentials
	if err := p.ValidateCredentials(ctx, credentials); err != nil {
		return nil, err
	}

	// TODO: Implement actual Scaleway SDK integration
	// For now, return an error indicating the feature is not yet implemented
	return nil, fmt.Errorf("%w: Scaleway subnet fetching not yet implemented", ErrProviderUnavailable)
}

// GetRegions returns the list of available Scaleway regions
func (p *ScalewayProvider) GetRegions() []string {
	return []string{
		"fr-par-1",
		"fr-par-2",
		"fr-par-3",
		"nl-ams-1",
		"nl-ams-2",
		"pl-waw-1",
		"pl-waw-2",
	}
}

// ValidateCredentials checks if the provided Scaleway credentials are valid
func (p *ScalewayProvider) ValidateCredentials(ctx context.Context, credentials CloudCredentials) error {
	if credentials.Provider != ProviderScaleway {
		return fmt.Errorf("invalid provider type: expected %s, got %s", ProviderScaleway, credentials.Provider)
	}

	if credentials.AccessKey == "" || credentials.SecretKey == "" {
		return ErrInvalidCredentials
	}

	// TODO: Implement actual Scaleway credential validation
	return nil
}

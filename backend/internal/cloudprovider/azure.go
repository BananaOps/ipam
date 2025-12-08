package cloudprovider

import (
	"context"
	"fmt"
)

// AzureProvider implements the CloudProvider interface for Microsoft Azure
type AzureProvider struct {
	name string
}

// NewAzureProvider creates a new Azure cloud provider instance
func NewAzureProvider() *AzureProvider {
	return &AzureProvider{
		name: "Microsoft Azure",
	}
}

// GetName returns the name of the cloud provider
func (p *AzureProvider) GetName() string {
	return p.name
}

// GetType returns the type of the cloud provider
func (p *AzureProvider) GetType() CloudProviderType {
	return ProviderAzure
}

// FetchSubnets retrieves all subnets from Azure
// This is a stub implementation - actual Azure SDK integration will be added in the future
func (p *AzureProvider) FetchSubnets(ctx context.Context, credentials CloudCredentials) ([]*CloudSubnet, error) {
	// Validate credentials
	if err := p.ValidateCredentials(ctx, credentials); err != nil {
		return nil, err
	}

	// TODO: Implement actual Azure SDK integration
	// For now, return an error indicating the feature is not yet implemented
	return nil, fmt.Errorf("%w: Azure subnet fetching not yet implemented", ErrProviderUnavailable)
}

// GetRegions returns the list of available Azure regions
func (p *AzureProvider) GetRegions() []string {
	return []string{
		"eastus",
		"eastus2",
		"westus",
		"westus2",
		"westus3",
		"centralus",
		"northeurope",
		"westeurope",
		"francecentral",
		"uksouth",
		"ukwest",
		"germanywestcentral",
		"southeastasia",
		"eastasia",
		"australiaeast",
		"australiasoutheast",
		"japaneast",
		"japanwest",
		"koreacentral",
		"canadacentral",
		"brazilsouth",
	}
}

// ValidateCredentials checks if the provided Azure credentials are valid
func (p *AzureProvider) ValidateCredentials(ctx context.Context, credentials CloudCredentials) error {
	if credentials.Provider != ProviderAzure {
		return fmt.Errorf("invalid provider type: expected %s, got %s", ProviderAzure, credentials.Provider)
	}

	if credentials.Token == "" {
		return ErrInvalidCredentials
	}

	// TODO: Implement actual Azure credential validation
	return nil
}

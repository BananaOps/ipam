package cloudprovider

import (
	"context"
	"fmt"
)

// AWSProvider implements the CloudProvider interface for Amazon Web Services
type AWSProvider struct {
	name string
}

// NewAWSProvider creates a new AWS cloud provider instance
func NewAWSProvider() *AWSProvider {
	return &AWSProvider{
		name: "Amazon Web Services",
	}
}

// GetName returns the name of the cloud provider
func (p *AWSProvider) GetName() string {
	return p.name
}

// GetType returns the type of the cloud provider
func (p *AWSProvider) GetType() CloudProviderType {
	return ProviderAWS
}

// FetchSubnets retrieves all subnets from AWS
// This is a stub implementation - actual AWS SDK integration will be added in the future
func (p *AWSProvider) FetchSubnets(ctx context.Context, credentials CloudCredentials) ([]*CloudSubnet, error) {
	// Validate credentials
	if err := p.ValidateCredentials(ctx, credentials); err != nil {
		return nil, err
	}

	// TODO: Implement actual AWS SDK integration
	// For now, return an error indicating the feature is not yet implemented
	return nil, fmt.Errorf("%w: AWS subnet fetching not yet implemented", ErrProviderUnavailable)
}

// GetRegions returns the list of available AWS regions
func (p *AWSProvider) GetRegions() []string {
	return []string{
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"eu-central-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-northeast-1",
		"ap-northeast-2",
		"sa-east-1",
		"ca-central-1",
	}
}

// ValidateCredentials checks if the provided AWS credentials are valid
func (p *AWSProvider) ValidateCredentials(ctx context.Context, credentials CloudCredentials) error {
	if credentials.Provider != ProviderAWS {
		return fmt.Errorf("invalid provider type: expected %s, got %s", ProviderAWS, credentials.Provider)
	}

	if credentials.AccessKey == "" || credentials.SecretKey == "" {
		return ErrInvalidCredentials
	}

	// TODO: Implement actual AWS credential validation
	return nil
}

package cloudprovider

// InitializeDefaultProviders creates and registers all supported cloud providers
func InitializeDefaultProviders() *CloudProviderManager {
	manager := NewCloudProviderManager()

	// Register all supported cloud providers
	providers := []CloudProvider{
		NewAWSProvider(),
		NewAzureProvider(),
		NewGCPProvider(),
		NewScalewayProvider(),
		NewOVHProvider(),
	}

	for _, provider := range providers {
		// Ignore registration errors for now - in production, these should be logged
		_ = manager.Register(provider)
	}

	return manager
}

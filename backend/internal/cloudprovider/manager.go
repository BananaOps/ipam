package cloudprovider

import (
	"context"
	"fmt"
	"sync"
)

// CloudProviderManager manages the registry of cloud providers
type CloudProviderManager struct {
	providers map[CloudProviderType]CloudProvider
	mu        sync.RWMutex
}

// NewCloudProviderManager creates a new CloudProviderManager instance
func NewCloudProviderManager() *CloudProviderManager {
	return &CloudProviderManager{
		providers: make(map[CloudProviderType]CloudProvider),
	}
}

// Register adds a cloud provider to the registry
func (m *CloudProviderManager) Register(provider CloudProvider) error {
	if provider == nil {
		return fmt.Errorf("cannot register nil provider")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	providerType := provider.GetType()
	if _, exists := m.providers[providerType]; exists {
		return fmt.Errorf("provider %s is already registered", providerType)
	}

	m.providers[providerType] = provider
	return nil
}

// Unregister removes a cloud provider from the registry
func (m *CloudProviderManager) Unregister(providerType CloudProviderType) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.providers[providerType]; !exists {
		return fmt.Errorf("%w: %s", ErrProviderNotFound, providerType)
	}

	delete(m.providers, providerType)
	return nil
}

// GetProvider retrieves a cloud provider by type
func (m *CloudProviderManager) GetProvider(providerType CloudProviderType) (CloudProvider, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	provider, exists := m.providers[providerType]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrProviderNotFound, providerType)
	}

	return provider, nil
}

// ListProviders returns all registered provider types
func (m *CloudProviderManager) ListProviders() []CloudProviderType {
	m.mu.RLock()
	defer m.mu.RUnlock()

	types := make([]CloudProviderType, 0, len(m.providers))
	for providerType := range m.providers {
		types = append(types, providerType)
	}

	return types
}

// IsProviderRegistered checks if a provider type is registered
func (m *CloudProviderManager) IsProviderRegistered(providerType CloudProviderType) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.providers[providerType]
	return exists
}

// FetchSubnetsFromProvider fetches subnets from a specific provider with error handling
func (m *CloudProviderManager) FetchSubnetsFromProvider(
	ctx context.Context,
	providerType CloudProviderType,
	credentials CloudCredentials,
) ([]*CloudSubnet, error) {
	provider, err := m.GetProvider(providerType)
	if err != nil {
		return nil, err
	}

	// Attempt to fetch subnets with error handling
	subnets, err := provider.FetchSubnets(ctx, credentials)
	if err != nil {
		// Wrap the error to provide more context
		return nil, fmt.Errorf("%w: failed to fetch subnets from %s: %v",
			ErrProviderUnavailable, providerType, err)
	}

	return subnets, nil
}

// FetchSubnetsFromAllProviders fetches subnets from all registered providers
// Returns a map of provider type to subnets, and a map of provider type to errors
func (m *CloudProviderManager) FetchSubnetsFromAllProviders(
	ctx context.Context,
	credentialsMap map[CloudProviderType]CloudCredentials,
) (map[CloudProviderType][]*CloudSubnet, map[CloudProviderType]error) {
	m.mu.RLock()
	providerTypes := make([]CloudProviderType, 0, len(m.providers))
	for providerType := range m.providers {
		providerTypes = append(providerTypes, providerType)
	}
	m.mu.RUnlock()

	results := make(map[CloudProviderType][]*CloudSubnet)
	errors := make(map[CloudProviderType]error)
	var wg sync.WaitGroup
	var resultsMu sync.Mutex

	for _, providerType := range providerTypes {
		credentials, hasCredentials := credentialsMap[providerType]
		if !hasCredentials {
			// Skip providers without credentials
			continue
		}

		wg.Add(1)
		go func(pt CloudProviderType, creds CloudCredentials) {
			defer wg.Done()

			subnets, err := m.FetchSubnetsFromProvider(ctx, pt, creds)

			resultsMu.Lock()
			defer resultsMu.Unlock()

			if err != nil {
				errors[pt] = err
			} else {
				results[pt] = subnets
			}
		}(providerType, credentials)
	}

	wg.Wait()
	return results, errors
}

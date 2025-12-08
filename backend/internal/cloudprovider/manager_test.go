package cloudprovider

import (
	"context"
	"errors"
	"testing"
)

// mockProvider is a mock implementation of CloudProvider for testing
type mockProvider struct {
	name         string
	providerType CloudProviderType
	regions      []string
	fetchError   error
}

func (m *mockProvider) GetName() string {
	return m.name
}

func (m *mockProvider) GetType() CloudProviderType {
	return m.providerType
}

func (m *mockProvider) FetchSubnets(ctx context.Context, credentials CloudCredentials) ([]*CloudSubnet, error) {
	if m.fetchError != nil {
		return nil, m.fetchError
	}
	return []*CloudSubnet{
		{
			CIDR:      "10.0.0.0/24",
			Name:      "test-subnet",
			Region:    "us-east-1",
			AccountID: "123456",
		},
	}, nil
}

func (m *mockProvider) GetRegions() []string {
	return m.regions
}

func (m *mockProvider) ValidateCredentials(ctx context.Context, credentials CloudCredentials) error {
	return nil
}

func TestNewCloudProviderManager(t *testing.T) {
	manager := NewCloudProviderManager()
	if manager == nil {
		t.Fatal("NewCloudProviderManager returned nil")
	}
	if manager.providers == nil {
		t.Fatal("providers map is nil")
	}
}

func TestRegisterProvider(t *testing.T) {
	manager := NewCloudProviderManager()

	tests := []struct {
		name        string
		provider    CloudProvider
		expectError bool
	}{
		{
			name: "register valid provider",
			provider: &mockProvider{
				name:         "Test Provider",
				providerType: "test",
			},
			expectError: false,
		},
		{
			name:        "register nil provider",
			provider:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.Register(tt.provider)
			if (err != nil) != tt.expectError {
				t.Errorf("Register() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestRegisterDuplicateProvider(t *testing.T) {
	manager := NewCloudProviderManager()
	provider := &mockProvider{
		name:         "Test Provider",
		providerType: "test",
	}

	// First registration should succeed
	err := manager.Register(provider)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	// Second registration should fail
	err = manager.Register(provider)
	if err == nil {
		t.Fatal("Expected error when registering duplicate provider, got nil")
	}
}

func TestGetProvider(t *testing.T) {
	manager := NewCloudProviderManager()
	provider := &mockProvider{
		name:         "Test Provider",
		providerType: "test",
	}

	// Register provider
	err := manager.Register(provider)
	if err != nil {
		t.Fatalf("Failed to register provider: %v", err)
	}

	// Get existing provider
	retrieved, err := manager.GetProvider("test")
	if err != nil {
		t.Fatalf("GetProvider() error = %v", err)
	}
	if retrieved.GetType() != "test" {
		t.Errorf("GetProvider() returned wrong provider type: got %v, want test", retrieved.GetType())
	}

	// Get non-existing provider
	_, err = manager.GetProvider("nonexistent")
	if err == nil {
		t.Fatal("Expected error when getting non-existent provider, got nil")
	}
	if !errors.Is(err, ErrProviderNotFound) {
		t.Errorf("Expected ErrProviderNotFound, got %v", err)
	}
}

func TestUnregisterProvider(t *testing.T) {
	manager := NewCloudProviderManager()
	provider := &mockProvider{
		name:         "Test Provider",
		providerType: "test",
	}

	// Register provider
	err := manager.Register(provider)
	if err != nil {
		t.Fatalf("Failed to register provider: %v", err)
	}

	// Unregister existing provider
	err = manager.Unregister("test")
	if err != nil {
		t.Fatalf("Unregister() error = %v", err)
	}

	// Verify provider is removed
	_, err = manager.GetProvider("test")
	if err == nil {
		t.Fatal("Provider still exists after unregistration")
	}

	// Unregister non-existing provider
	err = manager.Unregister("nonexistent")
	if err == nil {
		t.Fatal("Expected error when unregistering non-existent provider, got nil")
	}
}

func TestListProviders(t *testing.T) {
	manager := NewCloudProviderManager()

	// Empty list
	providers := manager.ListProviders()
	if len(providers) != 0 {
		t.Errorf("Expected empty list, got %d providers", len(providers))
	}

	// Register multiple providers
	provider1 := &mockProvider{name: "Provider 1", providerType: "test1"}
	provider2 := &mockProvider{name: "Provider 2", providerType: "test2"}

	manager.Register(provider1)
	manager.Register(provider2)

	providers = manager.ListProviders()
	if len(providers) != 2 {
		t.Errorf("Expected 2 providers, got %d", len(providers))
	}
}

func TestIsProviderRegistered(t *testing.T) {
	manager := NewCloudProviderManager()
	provider := &mockProvider{
		name:         "Test Provider",
		providerType: "test",
	}

	// Check non-registered provider
	if manager.IsProviderRegistered("test") {
		t.Error("IsProviderRegistered() returned true for non-registered provider")
	}

	// Register provider
	manager.Register(provider)

	// Check registered provider
	if !manager.IsProviderRegistered("test") {
		t.Error("IsProviderRegistered() returned false for registered provider")
	}
}

func TestFetchSubnetsFromProvider(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		provider    *mockProvider
		expectError bool
	}{
		{
			name: "successful fetch",
			provider: &mockProvider{
				name:         "Test Provider",
				providerType: "test-success",
				fetchError:   nil,
			},
			expectError: false,
		},
		{
			name: "provider unavailable",
			provider: &mockProvider{
				name:         "Test Provider",
				providerType: "test-fail",
				fetchError:   errors.New("connection failed"),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new manager for each test to avoid conflicts
			manager := NewCloudProviderManager()
			manager.Register(tt.provider)

			credentials := CloudCredentials{
				Provider:  tt.provider.providerType,
				AccessKey: "test-key",
				SecretKey: "test-secret",
			}

			subnets, err := manager.FetchSubnetsFromProvider(ctx, tt.provider.providerType, credentials)

			if (err != nil) != tt.expectError {
				t.Errorf("FetchSubnetsFromProvider() error = %v, expectError %v", err, tt.expectError)
			}

			if !tt.expectError && len(subnets) == 0 {
				t.Error("Expected subnets, got empty list")
			}
		})
	}
}

func TestFetchSubnetsFromNonExistentProvider(t *testing.T) {
	manager := NewCloudProviderManager()
	ctx := context.Background()

	credentials := CloudCredentials{
		Provider:  "nonexistent",
		AccessKey: "test-key",
		SecretKey: "test-secret",
	}

	_, err := manager.FetchSubnetsFromProvider(ctx, "nonexistent", credentials)
	if err == nil {
		t.Fatal("Expected error when fetching from non-existent provider, got nil")
	}
	if !errors.Is(err, ErrProviderNotFound) {
		t.Errorf("Expected ErrProviderNotFound, got %v", err)
	}
}

func TestFetchSubnetsFromAllProviders(t *testing.T) {
	manager := NewCloudProviderManager()
	ctx := context.Background()

	// Register multiple providers
	provider1 := &mockProvider{
		name:         "Provider 1",
		providerType: "test1",
		fetchError:   nil,
	}
	provider2 := &mockProvider{
		name:         "Provider 2",
		providerType: "test2",
		fetchError:   errors.New("provider unavailable"),
	}

	manager.Register(provider1)
	manager.Register(provider2)

	credentialsMap := map[CloudProviderType]CloudCredentials{
		"test1": {Provider: "test1", AccessKey: "key1", SecretKey: "secret1"},
		"test2": {Provider: "test2", AccessKey: "key2", SecretKey: "secret2"},
	}

	results, errs := manager.FetchSubnetsFromAllProviders(ctx, credentialsMap)

	// Check that we got results from provider1
	if _, ok := results["test1"]; !ok {
		t.Error("Expected results from test1 provider")
	}

	// Check that we got an error from provider2
	if _, ok := errs["test2"]; !ok {
		t.Error("Expected error from test2 provider")
	}
}

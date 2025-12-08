package cloudprovider

import (
	"context"
	"errors"
	"testing"
)

func TestAWSProvider(t *testing.T) {
	provider := NewAWSProvider()

	t.Run("GetName", func(t *testing.T) {
		name := provider.GetName()
		if name == "" {
			t.Error("GetName() returned empty string")
		}
	})

	t.Run("GetType", func(t *testing.T) {
		providerType := provider.GetType()
		if providerType != ProviderAWS {
			t.Errorf("GetType() = %v, want %v", providerType, ProviderAWS)
		}
	})

	t.Run("GetRegions", func(t *testing.T) {
		regions := provider.GetRegions()
		if len(regions) == 0 {
			t.Error("GetRegions() returned empty list")
		}
		// Check for some known AWS regions
		hasUSEast1 := false
		for _, region := range regions {
			if region == "us-east-1" {
				hasUSEast1 = true
				break
			}
		}
		if !hasUSEast1 {
			t.Error("GetRegions() does not include us-east-1")
		}
	})

	t.Run("ValidateCredentials - valid", func(t *testing.T) {
		ctx := context.Background()
		credentials := CloudCredentials{
			Provider:  ProviderAWS,
			AccessKey: "test-access-key",
			SecretKey: "test-secret-key",
		}
		err := provider.ValidateCredentials(ctx, credentials)
		if err != nil {
			t.Errorf("ValidateCredentials() error = %v, want nil", err)
		}
	})

	t.Run("ValidateCredentials - invalid provider", func(t *testing.T) {
		ctx := context.Background()
		credentials := CloudCredentials{
			Provider:  ProviderAzure,
			AccessKey: "test-access-key",
			SecretKey: "test-secret-key",
		}
		err := provider.ValidateCredentials(ctx, credentials)
		if err == nil {
			t.Error("ValidateCredentials() expected error for wrong provider type, got nil")
		}
	})

	t.Run("ValidateCredentials - missing credentials", func(t *testing.T) {
		ctx := context.Background()
		credentials := CloudCredentials{
			Provider: ProviderAWS,
		}
		err := provider.ValidateCredentials(ctx, credentials)
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Errorf("ValidateCredentials() error = %v, want %v", err, ErrInvalidCredentials)
		}
	})

	t.Run("FetchSubnets - not implemented", func(t *testing.T) {
		ctx := context.Background()
		credentials := CloudCredentials{
			Provider:  ProviderAWS,
			AccessKey: "test-access-key",
			SecretKey: "test-secret-key",
		}
		_, err := provider.FetchSubnets(ctx, credentials)
		if err == nil {
			t.Error("FetchSubnets() expected error for unimplemented feature, got nil")
		}
	})
}

func TestAzureProvider(t *testing.T) {
	provider := NewAzureProvider()

	t.Run("GetName", func(t *testing.T) {
		name := provider.GetName()
		if name == "" {
			t.Error("GetName() returned empty string")
		}
	})

	t.Run("GetType", func(t *testing.T) {
		providerType := provider.GetType()
		if providerType != ProviderAzure {
			t.Errorf("GetType() = %v, want %v", providerType, ProviderAzure)
		}
	})

	t.Run("GetRegions", func(t *testing.T) {
		regions := provider.GetRegions()
		if len(regions) == 0 {
			t.Error("GetRegions() returned empty list")
		}
	})

	t.Run("ValidateCredentials - valid", func(t *testing.T) {
		ctx := context.Background()
		credentials := CloudCredentials{
			Provider: ProviderAzure,
			Token:    "test-token",
		}
		err := provider.ValidateCredentials(ctx, credentials)
		if err != nil {
			t.Errorf("ValidateCredentials() error = %v, want nil", err)
		}
	})

	t.Run("ValidateCredentials - missing token", func(t *testing.T) {
		ctx := context.Background()
		credentials := CloudCredentials{
			Provider: ProviderAzure,
		}
		err := provider.ValidateCredentials(ctx, credentials)
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Errorf("ValidateCredentials() error = %v, want %v", err, ErrInvalidCredentials)
		}
	})
}

func TestGCPProvider(t *testing.T) {
	provider := NewGCPProvider()

	t.Run("GetName", func(t *testing.T) {
		name := provider.GetName()
		if name == "" {
			t.Error("GetName() returned empty string")
		}
	})

	t.Run("GetType", func(t *testing.T) {
		providerType := provider.GetType()
		if providerType != ProviderGCP {
			t.Errorf("GetType() = %v, want %v", providerType, ProviderGCP)
		}
	})

	t.Run("GetRegions", func(t *testing.T) {
		regions := provider.GetRegions()
		if len(regions) == 0 {
			t.Error("GetRegions() returned empty list")
		}
	})

	t.Run("ValidateCredentials - valid", func(t *testing.T) {
		ctx := context.Background()
		credentials := CloudCredentials{
			Provider: ProviderGCP,
			Token:    "test-token",
		}
		err := provider.ValidateCredentials(ctx, credentials)
		if err != nil {
			t.Errorf("ValidateCredentials() error = %v, want nil", err)
		}
	})
}

func TestScalewayProvider(t *testing.T) {
	provider := NewScalewayProvider()

	t.Run("GetName", func(t *testing.T) {
		name := provider.GetName()
		if name == "" {
			t.Error("GetName() returned empty string")
		}
	})

	t.Run("GetType", func(t *testing.T) {
		providerType := provider.GetType()
		if providerType != ProviderScaleway {
			t.Errorf("GetType() = %v, want %v", providerType, ProviderScaleway)
		}
	})

	t.Run("GetRegions", func(t *testing.T) {
		regions := provider.GetRegions()
		if len(regions) == 0 {
			t.Error("GetRegions() returned empty list")
		}
	})

	t.Run("ValidateCredentials - valid", func(t *testing.T) {
		ctx := context.Background()
		credentials := CloudCredentials{
			Provider:  ProviderScaleway,
			AccessKey: "test-access-key",
			SecretKey: "test-secret-key",
		}
		err := provider.ValidateCredentials(ctx, credentials)
		if err != nil {
			t.Errorf("ValidateCredentials() error = %v, want nil", err)
		}
	})
}

func TestOVHProvider(t *testing.T) {
	provider := NewOVHProvider()

	t.Run("GetName", func(t *testing.T) {
		name := provider.GetName()
		if name == "" {
			t.Error("GetName() returned empty string")
		}
	})

	t.Run("GetType", func(t *testing.T) {
		providerType := provider.GetType()
		if providerType != ProviderOVH {
			t.Errorf("GetType() = %v, want %v", providerType, ProviderOVH)
		}
	})

	t.Run("GetRegions", func(t *testing.T) {
		regions := provider.GetRegions()
		if len(regions) == 0 {
			t.Error("GetRegions() returned empty list")
		}
	})

	t.Run("ValidateCredentials - valid", func(t *testing.T) {
		ctx := context.Background()
		credentials := CloudCredentials{
			Provider:  ProviderOVH,
			AccessKey: "test-access-key",
			SecretKey: "test-secret-key",
		}
		err := provider.ValidateCredentials(ctx, credentials)
		if err != nil {
			t.Errorf("ValidateCredentials() error = %v, want nil", err)
		}
	})
}

func TestInitializeDefaultProviders(t *testing.T) {
	manager := InitializeDefaultProviders()

	if manager == nil {
		t.Fatal("InitializeDefaultProviders() returned nil")
	}

	// Check that all providers are registered
	expectedProviders := []CloudProviderType{
		ProviderAWS,
		ProviderAzure,
		ProviderGCP,
		ProviderScaleway,
		ProviderOVH,
	}

	for _, providerType := range expectedProviders {
		if !manager.IsProviderRegistered(providerType) {
			t.Errorf("Provider %s is not registered", providerType)
		}
	}

	// Verify we can get each provider
	for _, providerType := range expectedProviders {
		provider, err := manager.GetProvider(providerType)
		if err != nil {
			t.Errorf("Failed to get provider %s: %v", providerType, err)
		}
		if provider.GetType() != providerType {
			t.Errorf("Provider type mismatch: got %v, want %v", provider.GetType(), providerType)
		}
	}
}

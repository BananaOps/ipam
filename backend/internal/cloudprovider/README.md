# Cloud Provider Module

This package provides an extensible interface for integrating with various cloud providers to fetch subnet information.

## Overview

The cloud provider module is designed to support future integration with cloud provider APIs (AWS, Azure, GCP, Scaleway, OVH) for automatic IP subnet discovery. Currently, the module provides the interface and stub implementations that can be extended with actual SDK integrations.

## Architecture

### Core Components

1. **CloudProvider Interface**: Defines the contract that all cloud provider implementations must satisfy
2. **CloudProviderManager**: Manages the registry of cloud providers and provides methods to interact with them
3. **Provider Implementations**: Concrete implementations for each supported cloud provider

## Usage

### Initialize the Manager

```go
import "github.com/bananaops/ipam-bananaops/internal/cloudprovider"

// Initialize with all default providers
manager := cloudprovider.InitializeDefaultProviders()
```

### Register a Custom Provider

```go
// Create a custom provider
customProvider := &MyCustomProvider{}

// Register it
err := manager.Register(customProvider)
if err != nil {
    log.Fatal(err)
}
```

### Fetch Subnets from a Provider

```go
ctx := context.Background()

credentials := cloudprovider.CloudCredentials{
    Provider:  cloudprovider.ProviderAWS,
    AccessKey: "your-access-key",
    SecretKey: "your-secret-key",
    Region:    "us-east-1",
}

subnets, err := manager.FetchSubnetsFromProvider(ctx, cloudprovider.ProviderAWS, credentials)
if err != nil {
    // Handle error - provider might be unavailable
    log.Printf("Failed to fetch subnets: %v", err)
    return
}

for _, subnet := range subnets {
    fmt.Printf("Subnet: %s in region %s\n", subnet.CIDR, subnet.Region)
}
```

### Fetch from All Providers

```go
credentialsMap := map[cloudprovider.CloudProviderType]cloudprovider.CloudCredentials{
    cloudprovider.ProviderAWS: {
        Provider:  cloudprovider.ProviderAWS,
        AccessKey: "aws-key",
        SecretKey: "aws-secret",
    },
    cloudprovider.ProviderAzure: {
        Provider: cloudprovider.ProviderAzure,
        Token:    "azure-token",
    },
}

results, errors := manager.FetchSubnetsFromAllProviders(ctx, credentialsMap)

// Process results
for providerType, subnets := range results {
    fmt.Printf("Got %d subnets from %s\n", len(subnets), providerType)
}

// Handle errors
for providerType, err := range errors {
    log.Printf("Error from %s: %v\n", providerType, err)
}
```

## Supported Providers

- **AWS** (Amazon Web Services)
- **Azure** (Microsoft Azure)
- **GCP** (Google Cloud Platform)
- **Scaleway**
- **OVH**

## Error Handling

The module provides graceful error handling for common scenarios:

- `ErrProviderNotFound`: The requested provider is not registered
- `ErrProviderUnavailable`: The provider is unavailable or returned an error
- `ErrAuthenticationFailed`: Authentication with the provider failed
- `ErrRateLimited`: The provider rate limited the request
- `ErrInvalidCredentials`: The provided credentials are invalid

Example:

```go
subnets, err := manager.FetchSubnetsFromProvider(ctx, providerType, credentials)
if err != nil {
    if errors.Is(err, cloudprovider.ErrProviderUnavailable) {
        // Use cached data or display a warning
        log.Println("Provider unavailable, using cached data")
    } else {
        // Handle other errors
        return err
    }
}
```

## Implementing a New Provider

To add support for a new cloud provider:

1. Create a new file (e.g., `newprovider.go`)
2. Implement the `CloudProvider` interface:

```go
type NewProvider struct {
    name string
}

func NewNewProvider() *NewProvider {
    return &NewProvider{name: "New Provider"}
}

func (p *NewProvider) GetName() string {
    return p.name
}

func (p *NewProvider) GetType() CloudProviderType {
    return "newprovider"
}

func (p *NewProvider) FetchSubnets(ctx context.Context, credentials CloudCredentials) ([]*CloudSubnet, error) {
    // Implement actual API integration here
    return nil, nil
}

func (p *NewProvider) GetRegions() []string {
    return []string{"region1", "region2"}
}

func (p *NewProvider) ValidateCredentials(ctx context.Context, credentials CloudCredentials) error {
    // Implement credential validation
    return nil
}
```

3. Register it with the manager:

```go
manager.Register(NewNewProvider())
```

## Future Enhancements

The current implementation provides stub methods that return `ErrProviderUnavailable`. Future work includes:

1. Integrate AWS SDK for EC2 VPC subnet discovery
2. Integrate Azure SDK for Virtual Network subnet discovery
3. Integrate GCP SDK for VPC subnet discovery
4. Integrate Scaleway SDK for VPC subnet discovery
5. Integrate OVH API for network discovery
6. Add caching layer for fetched subnets
7. Add webhook support for real-time updates
8. Add support for filtering subnets by tags/labels

## Testing

Run the tests:

```bash
go test ./internal/cloudprovider/...
```

Run with coverage:

```bash
go test -cover ./internal/cloudprovider/...
```

## Thread Safety

The `CloudProviderManager` is thread-safe and can be used concurrently. All operations on the provider registry are protected by a read-write mutex.

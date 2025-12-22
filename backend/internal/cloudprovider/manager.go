package cloudprovider

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bananaops/ipam-bananaops/internal/cloudprovider/aws"
	"github.com/bananaops/ipam-bananaops/internal/config"
	"github.com/bananaops/ipam-bananaops/internal/repository"
)

// Manager manages cloud provider integrations
type Manager struct {
	config     *config.Config
	repository repository.SubnetRepository
	awsClients map[string]*aws.Client
	awsSyncs   map[string]*aws.SyncService
	mu         sync.RWMutex
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

// NewManager creates a new cloud provider manager
func NewManager(cfg *config.Config, repo repository.SubnetRepository) *Manager {
	return &Manager{
		config:     cfg,
		repository: repo,
		awsClients: make(map[string]*aws.Client),
		awsSyncs:   make(map[string]*aws.SyncService),
		stopCh:     make(chan struct{}),
	}
}

// Start initializes and starts cloud provider integrations
func (m *Manager) Start(ctx context.Context) error {
	if !m.config.CloudProviders.Enabled {
		log.Println("Cloud providers are disabled in configuration")
		return nil
	}

	log.Println("Starting cloud provider manager...")

	// Initialize AWS clients
	if err := m.initializeAWS(ctx); err != nil {
		return fmt.Errorf("failed to initialize AWS: %w", err)
	}

	// Start periodic sync
	if err := m.startPeriodicSync(ctx); err != nil {
		return fmt.Errorf("failed to start periodic sync: %w", err)
	}

	log.Println("Cloud provider manager started successfully")
	return nil
}

// Stop gracefully stops the cloud provider manager
func (m *Manager) Stop() {
	log.Println("Stopping cloud provider manager...")
	close(m.stopCh)
	m.wg.Wait()
	log.Println("Cloud provider manager stopped")
}

// initializeAWS initializes AWS clients for all configured regions
func (m *Manager) initializeAWS(ctx context.Context) error {
	if !m.config.CloudProviders.AWS.Enabled {
		log.Println("AWS integration is disabled")
		return nil
	}

	log.Printf("Initializing AWS integration for %d regions", len(m.config.CloudProviders.AWS.Regions))

	for _, regionConfig := range m.config.CloudProviders.AWS.Regions {
		awsConfig := aws.AWSConfig{
			Region:          regionConfig.Region,
			AccessKeyID:     regionConfig.AccessKeyID,
			SecretAccessKey: regionConfig.SecretAccessKey,
		}

		client, err := aws.NewClient(ctx, awsConfig)
		if err != nil {
			log.Printf("Failed to create AWS client for region %s: %v", regionConfig.Region, err)
			continue
		}

		// Validate credentials
		if err := client.ValidateCredentials(ctx); err != nil {
			log.Printf("Failed to validate AWS credentials for region %s: %v", regionConfig.Region, err)
			continue
		}

		m.mu.Lock()
		m.awsClients[regionConfig.Region] = client
		m.awsSyncs[regionConfig.Region] = aws.NewSyncService(client, m.repository)
		m.mu.Unlock()

		log.Printf("Successfully initialized AWS client for region: %s", regionConfig.Region)
	}

	if len(m.awsClients) == 0 {
		return fmt.Errorf("no AWS clients were successfully initialized")
	}

	return nil
}

// startPeriodicSync starts the periodic synchronization process
func (m *Manager) startPeriodicSync(ctx context.Context) error {
	syncInterval, err := m.config.CloudProviders.GetSyncInterval()
	if err != nil {
		return fmt.Errorf("invalid sync interval: %w", err)
	}

	log.Printf("Starting periodic sync with interval: %v", syncInterval)

	// Perform initial sync
	if err := m.SyncAll(ctx); err != nil {
		log.Printf("Initial sync failed: %v", err)
	}

	// Start periodic sync goroutine
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		ticker := time.NewTicker(syncInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := m.SyncAll(ctx); err != nil {
					log.Printf("Periodic sync failed: %v", err)
				}
			case <-m.stopCh:
				return
			}
		}
	}()

	return nil
}

// SyncAll synchronizes all cloud providers
func (m *Manager) SyncAll(ctx context.Context) error {
	log.Println("Starting full cloud provider synchronization...")

	var errors []error

	// Sync AWS
	if err := m.syncAWS(ctx); err != nil {
		errors = append(errors, fmt.Errorf("AWS sync failed: %w", err))
	}

	if len(errors) > 0 {
		log.Printf("Synchronization completed with %d errors", len(errors))
		return fmt.Errorf("sync errors: %v", errors)
	}

	log.Println("Full cloud provider synchronization completed successfully")
	return nil
}

// syncAWS synchronizes all AWS regions
func (m *Manager) syncAWS(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.awsSyncs) == 0 {
		return nil
	}

	log.Printf("Synchronizing %d AWS regions", len(m.awsSyncs))

	var errors []error
	for region, syncService := range m.awsSyncs {
		log.Printf("Synchronizing AWS region: %s", region)
		if err := syncService.SyncAll(ctx); err != nil {
			errors = append(errors, fmt.Errorf("region %s: %w", region, err))
			continue
		}
		log.Printf("Successfully synchronized AWS region: %s", region)
	}

	if len(errors) > 0 {
		return fmt.Errorf("AWS sync errors: %v", errors)
	}

	return nil
}

// SyncAWSRegion synchronizes a specific AWS region
func (m *Manager) SyncAWSRegion(ctx context.Context, region string) error {
	m.mu.RLock()
	syncService, exists := m.awsSyncs[region]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("AWS region %s is not configured", region)
	}

	log.Printf("Synchronizing AWS region: %s", region)
	return syncService.SyncAll(ctx)
}

// UpdateUtilization updates utilization data for all cloud providers
func (m *Manager) UpdateUtilization(ctx context.Context) error {
	log.Println("Updating utilization data for all cloud providers...")

	var errors []error

	// Update AWS utilization
	m.mu.RLock()
	for region, syncService := range m.awsSyncs {
		if err := syncService.UpdateUtilization(ctx); err != nil {
			errors = append(errors, fmt.Errorf("AWS region %s: %w", region, err))
		}
	}
	m.mu.RUnlock()

	if len(errors) > 0 {
		return fmt.Errorf("utilization update errors: %v", errors)
	}

	log.Println("Utilization data updated successfully")
	return nil
}

// GetAWSClient returns the AWS client for a specific region
func (m *Manager) GetAWSClient(region string) (*aws.Client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, exists := m.awsClients[region]
	if !exists {
		return nil, fmt.Errorf("AWS client for region %s not found", region)
	}

	return client, nil
}

// ListAWSRegions returns all configured AWS regions
func (m *Manager) ListAWSRegions() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	regions := make([]string, 0, len(m.awsClients))
	for region := range m.awsClients {
		regions = append(regions, region)
	}

	return regions
}

// IsEnabled returns whether cloud providers are enabled
func (m *Manager) IsEnabled() bool {
	return m.config.CloudProviders.Enabled
}

// IsAWSEnabled returns whether AWS integration is enabled
func (m *Manager) IsAWSEnabled() bool {
	return m.config.CloudProviders.AWS.Enabled
}

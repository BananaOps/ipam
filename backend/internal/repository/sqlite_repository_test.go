package repository

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	pb "github.com/bananaops/ipam-bananaops/proto"
)

func TestSQLiteRepository_CreateAndFind(t *testing.T) {
	// Create temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	repo, err := NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// Create a test subnet
	subnet := &pb.Subnet{
		Id:           "test-id-1",
		Cidr:         "192.168.1.0/24",
		Name:         "Test Subnet",
		Description:  "A test subnet",
		Location:     "datacenter-1",
		LocationType: pb.LocationType_DATACENTER,
		Details: &pb.SubnetDetails{
			Address:     "192.168.1.0",
			Netmask:     "255.255.255.0",
			Wildcard:    "0.0.0.255",
			Network:     "192.168.1.0",
			Type:        "IPv4",
			Broadcast:   "192.168.1.255",
			HostMin:     "192.168.1.1",
			HostMax:     "192.168.1.254",
			HostsPerNet: 254,
			IsPublic:    false,
		},
		Utilization: &pb.UtilizationInfo{
			TotalIps:           254,
			AllocatedIps:       0,
			UtilizationPercent: 0.0,
		},
		CreatedAt: 1234567890,
		UpdatedAt: 1234567890,
	}

	ctx := context.Background()

	// Test Create
	err = repo.Create(ctx, subnet)
	if err != nil {
		t.Fatalf("Failed to create subnet: %v", err)
	}

	// Test FindByID
	found, err := repo.FindByID(ctx, "test-id-1")
	if err != nil {
		t.Fatalf("Failed to find subnet: %v", err)
	}

	if found.Id != subnet.Id {
		t.Errorf("Expected ID %s, got %s", subnet.Id, found.Id)
	}
	if found.Cidr != subnet.Cidr {
		t.Errorf("Expected CIDR %s, got %s", subnet.Cidr, found.Cidr)
	}
	if found.Name != subnet.Name {
		t.Errorf("Expected Name %s, got %s", subnet.Name, found.Name)
	}
}

func TestSQLiteRepository_FindAll(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	repo, err := NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()

	// Create multiple subnets
	subnets := []*pb.Subnet{
		{
			Id:           "test-id-1",
			Cidr:         "192.168.1.0/24",
			Name:         "Subnet 1",
			Location:     "datacenter-1",
			LocationType: pb.LocationType_DATACENTER,
			Details:      &pb.SubnetDetails{HostsPerNet: 254},
			Utilization:  &pb.UtilizationInfo{},
			CreatedAt:    1234567890,
			UpdatedAt:    1234567890,
		},
		{
			Id:           "test-id-2",
			Cidr:         "10.0.0.0/16",
			Name:         "Subnet 2",
			Location:     "datacenter-2",
			LocationType: pb.LocationType_DATACENTER,
			Details:      &pb.SubnetDetails{HostsPerNet: 65534},
			Utilization:  &pb.UtilizationInfo{},
			CreatedAt:    1234567891,
			UpdatedAt:    1234567891,
		},
	}

	for _, subnet := range subnets {
		if err := repo.Create(ctx, subnet); err != nil {
			t.Fatalf("Failed to create subnet: %v", err)
		}
	}

	// Test FindAll without filters
	found, err := repo.FindAll(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to find all subnets: %v", err)
	}

	if len(found) != 2 {
		t.Errorf("Expected 2 subnets, got %d", len(found))
	}

	// Test FindAll with location filter
	filters := &SubnetFilters{
		LocationFilter: "datacenter-1",
	}
	found, err = repo.FindAll(ctx, filters)
	if err != nil {
		t.Fatalf("Failed to find filtered subnets: %v", err)
	}

	if len(found) != 1 {
		t.Errorf("Expected 1 subnet, got %d", len(found))
	}
	if len(found) > 0 && found[0].Location != "datacenter-1" {
		t.Errorf("Expected location datacenter-1, got %s", found[0].Location)
	}
}

func TestSQLiteRepository_Update(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	repo, err := NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()

	// Create a subnet
	subnet := &pb.Subnet{
		Id:           "test-id-1",
		Cidr:         "192.168.1.0/24",
		Name:         "Original Name",
		Location:     "datacenter-1",
		LocationType: pb.LocationType_DATACENTER,
		Details:      &pb.SubnetDetails{HostsPerNet: 254},
		Utilization:  &pb.UtilizationInfo{},
		CreatedAt:    1234567890,
		UpdatedAt:    1234567890,
	}

	if err := repo.Create(ctx, subnet); err != nil {
		t.Fatalf("Failed to create subnet: %v", err)
	}

	// Update the subnet
	subnet.Name = "Updated Name"
	subnet.UpdatedAt = 1234567900

	if err := repo.Update(ctx, subnet); err != nil {
		t.Fatalf("Failed to update subnet: %v", err)
	}

	// Verify update
	found, err := repo.FindByID(ctx, "test-id-1")
	if err != nil {
		t.Fatalf("Failed to find subnet: %v", err)
	}

	if found.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %s", found.Name)
	}
}

func TestSQLiteRepository_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	repo, err := NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()

	// Create a subnet
	subnet := &pb.Subnet{
		Id:           "test-id-1",
		Cidr:         "192.168.1.0/24",
		Name:         "Test Subnet",
		Location:     "datacenter-1",
		LocationType: pb.LocationType_DATACENTER,
		Details:      &pb.SubnetDetails{HostsPerNet: 254},
		Utilization:  &pb.UtilizationInfo{},
		CreatedAt:    1234567890,
		UpdatedAt:    1234567890,
	}

	if err := repo.Create(ctx, subnet); err != nil {
		t.Fatalf("Failed to create subnet: %v", err)
	}

	// Delete the subnet
	if err := repo.Delete(ctx, "test-id-1"); err != nil {
		t.Fatalf("Failed to delete subnet: %v", err)
	}

	// Verify deletion
	_, err = repo.FindByID(ctx, "test-id-1")
	if err == nil {
		t.Error("Expected error when finding deleted subnet, got nil")
	}
}

func TestSQLiteRepository_DatabasePath(t *testing.T) {
	// Test that database directory is created if it doesn't exist
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "nested", "dir", "test.db")

	repo, err := NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// Verify database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}
}

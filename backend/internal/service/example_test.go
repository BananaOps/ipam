package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/bananaops/ipam-bananaops/internal/repository"
	pb "github.com/bananaops/ipam-bananaops/proto"
)

// ExampleIPService demonstrates how to use the IPService
func ExampleIPService() {
	// Create a new IPService
	ipService := NewGoIPAMService()

	// Validate a CIDR
	cidr := "192.168.1.0/24"
	err := ipService.ValidateCIDR(cidr)
	if err != nil {
		fmt.Printf("Invalid CIDR: %v\n", err)
		return
	}

	// Calculate subnet details
	details, err := ipService.CalculateSubnetDetails(cidr)
	if err != nil {
		fmt.Printf("Failed to calculate details: %v\n", err)
		return
	}

	fmt.Printf("Network: %s\n", details.Network)
	fmt.Printf("Address: %s\n", details.Address)
	fmt.Printf("Netmask: %s\n", details.Netmask)
	fmt.Printf("Wildcard: %s\n", details.Wildcard)
	fmt.Printf("Broadcast: %s\n", details.Broadcast)
	fmt.Printf("Host Range: %s - %s\n", details.HostMin, details.HostMax)
	fmt.Printf("Hosts per Net: %d\n", details.HostsPerNet)
	fmt.Printf("Type: %s\n", details.Type)
	fmt.Printf("Is Public: %v\n", details.IsPublic)

	// Calculate utilization
	utilization := ipService.CalculateUtilization(details.HostsPerNet, 100)
	fmt.Printf("Utilization: %.2f%%\n", utilization)

	// Output:
	// Network: 192.168.1.0/24
	// Address: 192.168.1.0
	// Netmask: 255.255.255.0
	// Wildcard: 0.0.0.255
	// Broadcast: 192.168.1.255
	// Host Range: 192.168.1.1 - 192.168.1.254
	// Hosts per Net: 254
	// Type: IPv4
	// Is Public: false
	// Utilization: 39.37%
}

// TestServiceLayerIntegration demonstrates the full integration
func TestServiceLayerIntegration(t *testing.T) {
	// This test demonstrates how the IPService integrates with the ServiceLayer
	// In a real scenario, you would use a real repository

	ipService := NewGoIPAMService()

	// Test 1: Validate and calculate for a private IPv4 subnet
	cidr := "10.0.0.0/24"
	err := ipService.ValidateCIDR(cidr)
	if err != nil {
		t.Fatalf("Failed to validate CIDR: %v", err)
	}

	details, err := ipService.CalculateSubnetDetails(cidr)
	if err != nil {
		t.Fatalf("Failed to calculate details: %v", err)
	}

	if details.IsPublic {
		t.Error("Expected private subnet, got public")
	}

	if details.HostsPerNet != 254 {
		t.Errorf("Expected 254 hosts, got %d", details.HostsPerNet)
	}

	// Test 2: Validate and calculate for a public IPv4 subnet
	publicCIDR := "8.8.8.0/24"
	err = ipService.ValidateCIDR(publicCIDR)
	if err != nil {
		t.Fatalf("Failed to validate public CIDR: %v", err)
	}

	publicDetails, err := ipService.CalculateSubnetDetails(publicCIDR)
	if err != nil {
		t.Fatalf("Failed to calculate public details: %v", err)
	}

	if !publicDetails.IsPublic {
		t.Error("Expected public subnet, got private")
	}

	// Test 3: Calculate utilization
	utilization := ipService.CalculateUtilization(254, 127)
	expectedUtilization := float32(50.0)
	if utilization != expectedUtilization {
		t.Errorf("Expected utilization %.2f%%, got %.2f%%", expectedUtilization, utilization)
	}

	t.Log("IPService integration test passed successfully")
}

// TestServiceLayerWithIPService demonstrates how ServiceLayer uses IPService
func TestServiceLayerWithIPService(t *testing.T) {
	// Create mock repository (in real scenario, use actual repository)
	mockRepo := &mockSubnetRepository{
		subnets: make(map[string]*pb.Subnet),
	}

	// Create IPService
	ipService := NewGoIPAMService()

	// Create ServiceLayer
	serviceLayer := NewServiceLayer(mockRepo, ipService, nil)

	// Test CreateSubnet
	ctx := context.Background()
	req := &pb.CreateSubnetRequest{
		Cidr:         "192.168.1.0/24",
		Name:         "Test Subnet",
		Description:  "A test subnet for demonstration",
		Location:     "datacenter-1",
		LocationType: pb.LocationType_DATACENTER,
	}

	resp, err := serviceLayer.CreateSubnet(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create subnet: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("CreateSubnet returned error: %s", resp.Error.Message)
	}

	if resp.Subnet == nil {
		t.Fatal("Expected subnet in response, got nil")
	}

	// Verify subnet details were calculated
	if resp.Subnet.Details == nil {
		t.Fatal("Expected subnet details, got nil")
	}

	if resp.Subnet.Details.HostsPerNet != 254 {
		t.Errorf("Expected 254 hosts, got %d", resp.Subnet.Details.HostsPerNet)
	}

	if resp.Subnet.Details.IsPublic {
		t.Error("Expected private subnet, got public")
	}

	// Verify utilization was initialized
	if resp.Subnet.Utilization == nil {
		t.Fatal("Expected utilization info, got nil")
	}

	if resp.Subnet.Utilization.TotalIps != 254 {
		t.Errorf("Expected 254 total IPs, got %d", resp.Subnet.Utilization.TotalIps)
	}

	if resp.Subnet.Utilization.UtilizationPercent != 0.0 {
		t.Errorf("Expected 0%% utilization, got %.2f%%", resp.Subnet.Utilization.UtilizationPercent)
	}

	t.Log("ServiceLayer integration with IPService test passed successfully")
}

// mockSubnetRepository is a simple in-memory repository for testing
type mockSubnetRepository struct {
	subnets map[string]*pb.Subnet
}

func (m *mockSubnetRepository) Create(ctx context.Context, subnet *pb.Subnet) error {
	m.subnets[subnet.Id] = subnet
	return nil
}

func (m *mockSubnetRepository) FindByID(ctx context.Context, id string) (*pb.Subnet, error) {
	subnet, ok := m.subnets[id]
	if !ok {
		return nil, fmt.Errorf("subnet not found")
	}
	return subnet, nil
}

func (m *mockSubnetRepository) FindAll(ctx context.Context, filters *repository.SubnetFilters) ([]*pb.Subnet, error) {
	result := make([]*pb.Subnet, 0, len(m.subnets))
	for _, subnet := range m.subnets {
		result = append(result, subnet)
	}
	return result, nil
}

func (m *mockSubnetRepository) Update(ctx context.Context, subnet *pb.Subnet) error {
	m.subnets[subnet.Id] = subnet
	return nil
}

func (m *mockSubnetRepository) Delete(ctx context.Context, id string) error {
	delete(m.subnets, id)
	return nil
}

func (m *mockSubnetRepository) Close() error {
	return nil
}

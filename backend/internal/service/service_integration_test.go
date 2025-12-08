package service

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/bananaops/ipam-bananaops/internal/repository"
	pb "github.com/bananaops/ipam-bananaops/proto"
)

// TestServiceLayerFullIntegration tests the complete integration of ServiceLayer with real repository and IPService
func TestServiceLayerFullIntegration(t *testing.T) {
	// Create temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create real SQLite repository
	repo, err := repository.NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// Create real IPService
	ipService := NewGoIPAMService()

	// Create ServiceLayer with real dependencies
	serviceLayer := NewServiceLayer(repo, ipService, nil)

	ctx := context.Background()

	// Test 1: Create a subnet with validation and calculation
	t.Run("CreateSubnet", func(t *testing.T) {
		req := &pb.CreateSubnetRequest{
			Cidr:         "10.0.0.0/24",
			Name:         "Test Network",
			Description:  "Integration test subnet",
			Location:     "datacenter-1",
			LocationType: pb.LocationType_DATACENTER,
		}

		resp, err := serviceLayer.CreateSubnet(ctx, req)
		if err != nil {
			t.Fatalf("CreateSubnet failed: %v", err)
		}

		if resp.Error != nil {
			t.Fatalf("CreateSubnet returned error: %s", resp.Error.Message)
		}

		if resp.Subnet == nil {
			t.Fatal("Expected subnet in response")
		}

		// Verify subnet details were calculated by IPService
		if resp.Subnet.Details == nil {
			t.Fatal("Expected subnet details")
		}

		if resp.Subnet.Details.Address != "10.0.0.0" {
			t.Errorf("Expected address 10.0.0.0, got %s", resp.Subnet.Details.Address)
		}

		if resp.Subnet.Details.Netmask != "255.255.255.0" {
			t.Errorf("Expected netmask 255.255.255.0, got %s", resp.Subnet.Details.Netmask)
		}

		if resp.Subnet.Details.HostsPerNet != 254 {
			t.Errorf("Expected 254 hosts, got %d", resp.Subnet.Details.HostsPerNet)
		}

		if resp.Subnet.Details.IsPublic {
			t.Error("Expected private subnet")
		}

		// Verify utilization was initialized
		if resp.Subnet.Utilization.TotalIps != 254 {
			t.Errorf("Expected 254 total IPs, got %d", resp.Subnet.Utilization.TotalIps)
		}

		if resp.Subnet.Utilization.UtilizationPercent != 0.0 {
			t.Errorf("Expected 0%% utilization, got %.2f%%", resp.Subnet.Utilization.UtilizationPercent)
		}
	})

	// Test 2: Create subnet with invalid CIDR (validation)
	t.Run("CreateSubnetInvalidCIDR", func(t *testing.T) {
		req := &pb.CreateSubnetRequest{
			Cidr:         "invalid-cidr",
			Name:         "Invalid Subnet",
			Location:     "datacenter-1",
			LocationType: pb.LocationType_DATACENTER,
		}

		resp, err := serviceLayer.CreateSubnet(ctx, req)
		if err != nil {
			t.Fatalf("CreateSubnet failed: %v", err)
		}

		if resp.Error == nil {
			t.Fatal("Expected error for invalid CIDR")
		}

		if resp.Error.Code != "INVALID_CIDR" {
			t.Errorf("Expected error code INVALID_CIDR, got %s", resp.Error.Code)
		}
	})

	// Test 3: List subnets
	t.Run("ListSubnets", func(t *testing.T) {
		// Create multiple subnets
		subnets := []struct {
			cidr     string
			name     string
			location string
		}{
			{"192.168.1.0/24", "Subnet 1", "datacenter-1"},
			{"192.168.2.0/24", "Subnet 2", "datacenter-2"},
			{"172.16.0.0/16", "Subnet 3", "datacenter-1"},
		}

		for _, s := range subnets {
			req := &pb.CreateSubnetRequest{
				Cidr:         s.cidr,
				Name:         s.name,
				Location:     s.location,
				LocationType: pb.LocationType_DATACENTER,
			}

			resp, err := serviceLayer.CreateSubnet(ctx, req)
			if err != nil || resp.Error != nil {
				t.Fatalf("Failed to create subnet %s", s.name)
			}
		}

		// List all subnets
		listReq := &pb.ListSubnetsRequest{}
		listResp, err := serviceLayer.ListSubnets(ctx, listReq)
		if err != nil {
			t.Fatalf("ListSubnets failed: %v", err)
		}

		if listResp.Error != nil {
			t.Fatalf("ListSubnets returned error: %s", listResp.Error.Message)
		}

		// Should have at least 4 subnets (1 from Test 1 + 3 from this test)
		if len(listResp.Subnets) < 4 {
			t.Errorf("Expected at least 4 subnets, got %d", len(listResp.Subnets))
		}

		// Test filtering by location
		filterReq := &pb.ListSubnetsRequest{
			LocationFilter: "datacenter-1",
		}
		filterResp, err := serviceLayer.ListSubnets(ctx, filterReq)
		if err != nil {
			t.Fatalf("ListSubnets with filter failed: %v", err)
		}

		if filterResp.Error != nil {
			t.Fatalf("ListSubnets with filter returned error: %s", filterResp.Error.Message)
		}

		// Should have at least 3 subnets in datacenter-1
		if len(filterResp.Subnets) < 3 {
			t.Errorf("Expected at least 3 subnets in datacenter-1, got %d", len(filterResp.Subnets))
		}

		// Verify all returned subnets are from datacenter-1
		for _, subnet := range filterResp.Subnets {
			if subnet.Location != "datacenter-1" {
				t.Errorf("Expected location datacenter-1, got %s", subnet.Location)
			}
		}
	})

	// Test 4: Get subnet by ID
	t.Run("GetSubnet", func(t *testing.T) {
		// Create a subnet
		createReq := &pb.CreateSubnetRequest{
			Cidr:         "10.1.0.0/16",
			Name:         "Get Test Subnet",
			Location:     "datacenter-3",
			LocationType: pb.LocationType_DATACENTER,
		}

		createResp, err := serviceLayer.CreateSubnet(ctx, createReq)
		if err != nil || createResp.Error != nil {
			t.Fatal("Failed to create subnet for get test")
		}

		subnetID := createResp.Subnet.Id

		// Get the subnet
		getReq := &pb.GetSubnetRequest{Id: subnetID}
		getResp, err := serviceLayer.GetSubnet(ctx, getReq)
		if err != nil {
			t.Fatalf("GetSubnet failed: %v", err)
		}

		if getResp.Error != nil {
			t.Fatalf("GetSubnet returned error: %s", getResp.Error.Message)
		}

		if getResp.Subnet.Id != subnetID {
			t.Errorf("Expected subnet ID %s, got %s", subnetID, getResp.Subnet.Id)
		}

		if getResp.Subnet.Name != "Get Test Subnet" {
			t.Errorf("Expected name 'Get Test Subnet', got %s", getResp.Subnet.Name)
		}
	})

	// Test 5: Get non-existent subnet
	t.Run("GetSubnetNotFound", func(t *testing.T) {
		getReq := &pb.GetSubnetRequest{Id: "non-existent-id"}
		getResp, err := serviceLayer.GetSubnet(ctx, getReq)
		if err != nil {
			t.Fatalf("GetSubnet failed: %v", err)
		}

		if getResp.Error == nil {
			t.Fatal("Expected error for non-existent subnet")
		}

		if getResp.Error.Code != "SUBNET_NOT_FOUND" {
			t.Errorf("Expected error code SUBNET_NOT_FOUND, got %s", getResp.Error.Code)
		}
	})

	// Test 6: Update subnet with CIDR change (recalculation)
	t.Run("UpdateSubnetWithCIDRChange", func(t *testing.T) {
		// Create a subnet
		createReq := &pb.CreateSubnetRequest{
			Cidr:         "10.2.0.0/24",
			Name:         "Update Test Subnet",
			Location:     "datacenter-4",
			LocationType: pb.LocationType_DATACENTER,
		}

		createResp, err := serviceLayer.CreateSubnet(ctx, createReq)
		if err != nil || createResp.Error != nil {
			t.Fatal("Failed to create subnet for update test")
		}

		subnetID := createResp.Subnet.Id

		// Update with new CIDR
		updateReq := &pb.UpdateSubnetRequest{
			Id:   subnetID,
			Cidr: "10.2.0.0/16",
			Name: "Updated Subnet Name",
		}

		updateResp, err := serviceLayer.UpdateSubnet(ctx, updateReq)
		if err != nil {
			t.Fatalf("UpdateSubnet failed: %v", err)
		}

		if updateResp.Error != nil {
			t.Fatalf("UpdateSubnet returned error: %s", updateResp.Error.Message)
		}

		// Verify CIDR was updated
		if updateResp.Subnet.Cidr != "10.2.0.0/16" {
			t.Errorf("Expected CIDR 10.2.0.0/16, got %s", updateResp.Subnet.Cidr)
		}

		// Verify name was updated
		if updateResp.Subnet.Name != "Updated Subnet Name" {
			t.Errorf("Expected name 'Updated Subnet Name', got %s", updateResp.Subnet.Name)
		}

		// Verify subnet details were recalculated
		if updateResp.Subnet.Details.HostsPerNet != 65534 {
			t.Errorf("Expected 65534 hosts for /16, got %d", updateResp.Subnet.Details.HostsPerNet)
		}

		// Verify utilization was updated
		if updateResp.Subnet.Utilization.TotalIps != 65534 {
			t.Errorf("Expected 65534 total IPs, got %d", updateResp.Subnet.Utilization.TotalIps)
		}
	})

	// Test 7: Update subnet with invalid CIDR
	t.Run("UpdateSubnetInvalidCIDR", func(t *testing.T) {
		// Create a subnet
		createReq := &pb.CreateSubnetRequest{
			Cidr:         "10.3.0.0/24",
			Name:         "Invalid Update Test",
			Location:     "datacenter-5",
			LocationType: pb.LocationType_DATACENTER,
		}

		createResp, err := serviceLayer.CreateSubnet(ctx, createReq)
		if err != nil || createResp.Error != nil {
			t.Fatal("Failed to create subnet for invalid update test")
		}

		subnetID := createResp.Subnet.Id

		// Try to update with invalid CIDR
		updateReq := &pb.UpdateSubnetRequest{
			Id:   subnetID,
			Cidr: "invalid-cidr",
		}

		updateResp, err := serviceLayer.UpdateSubnet(ctx, updateReq)
		if err != nil {
			t.Fatalf("UpdateSubnet failed: %v", err)
		}

		if updateResp.Error == nil {
			t.Fatal("Expected error for invalid CIDR")
		}

		if updateResp.Error.Code != "INVALID_CIDR" {
			t.Errorf("Expected error code INVALID_CIDR, got %s", updateResp.Error.Code)
		}
	})

	// Test 8: Delete subnet
	t.Run("DeleteSubnet", func(t *testing.T) {
		// Create a subnet
		createReq := &pb.CreateSubnetRequest{
			Cidr:         "10.4.0.0/24",
			Name:         "Delete Test Subnet",
			Location:     "datacenter-6",
			LocationType: pb.LocationType_DATACENTER,
		}

		createResp, err := serviceLayer.CreateSubnet(ctx, createReq)
		if err != nil || createResp.Error != nil {
			t.Fatal("Failed to create subnet for delete test")
		}

		subnetID := createResp.Subnet.Id

		// Delete the subnet
		deleteReq := &pb.DeleteSubnetRequest{Id: subnetID}
		deleteResp, err := serviceLayer.DeleteSubnet(ctx, deleteReq)
		if err != nil {
			t.Fatalf("DeleteSubnet failed: %v", err)
		}

		if deleteResp.Error != nil {
			t.Fatalf("DeleteSubnet returned error: %s", deleteResp.Error.Message)
		}

		if !deleteResp.Success {
			t.Error("Expected success=true")
		}

		// Verify subnet was deleted
		getReq := &pb.GetSubnetRequest{Id: subnetID}
		getResp, err := serviceLayer.GetSubnet(ctx, getReq)
		if err != nil {
			t.Fatalf("GetSubnet failed: %v", err)
		}

		if getResp.Error == nil {
			t.Fatal("Expected error when getting deleted subnet")
		}

		if getResp.Error.Code != "SUBNET_NOT_FOUND" {
			t.Errorf("Expected error code SUBNET_NOT_FOUND, got %s", getResp.Error.Code)
		}
	})

	// Test 9: Delete non-existent subnet
	t.Run("DeleteSubnetNotFound", func(t *testing.T) {
		deleteReq := &pb.DeleteSubnetRequest{Id: "non-existent-id"}
		deleteResp, err := serviceLayer.DeleteSubnet(ctx, deleteReq)
		if err != nil {
			t.Fatalf("DeleteSubnet failed: %v", err)
		}

		if deleteResp.Error == nil {
			t.Fatal("Expected error for non-existent subnet")
		}

		if deleteResp.Error.Code != "SUBNET_NOT_FOUND" {
			t.Errorf("Expected error code SUBNET_NOT_FOUND, got %s", deleteResp.Error.Code)
		}

		if deleteResp.Success {
			t.Error("Expected success=false")
		}
	})

	// Test 10: Create subnet with cloud info
	t.Run("CreateSubnetWithCloudInfo", func(t *testing.T) {
		req := &pb.CreateSubnetRequest{
			Cidr:         "172.31.0.0/16",
			Name:         "AWS VPC Subnet",
			Description:  "Test cloud subnet",
			Location:     "us-east-1",
			LocationType: pb.LocationType_CLOUD,
			CloudInfo: &pb.CloudInfo{
				Provider:  "aws",
				Region:    "us-east-1",
				AccountId: "123456789012",
			},
		}

		resp, err := serviceLayer.CreateSubnet(ctx, req)
		if err != nil {
			t.Fatalf("CreateSubnet failed: %v", err)
		}

		if resp.Error != nil {
			t.Fatalf("CreateSubnet returned error: %s", resp.Error.Message)
		}

		// Verify cloud info was stored
		if resp.Subnet.CloudInfo == nil {
			t.Fatal("Expected cloud info")
		}

		if resp.Subnet.CloudInfo.Provider != "aws" {
			t.Errorf("Expected provider aws, got %s", resp.Subnet.CloudInfo.Provider)
		}

		if resp.Subnet.CloudInfo.Region != "us-east-1" {
			t.Errorf("Expected region us-east-1, got %s", resp.Subnet.CloudInfo.Region)
		}

		// Verify it's classified as private (172.31.0.0/16 is private)
		if resp.Subnet.Details.IsPublic {
			t.Error("Expected private subnet")
		}
	})
}

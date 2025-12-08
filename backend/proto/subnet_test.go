package proto

import (
	"testing"
)

// TestProtobufGeneration verifies that protobuf code was generated correctly
func TestProtobufGeneration(t *testing.T) {
	// Test that we can create a Subnet message
	subnet := &Subnet{
		Id:           "test-id",
		Cidr:         "10.0.0.0/24",
		Name:         "Test Subnet",
		Description:  "Test description",
		Location:     "datacenter-1",
		LocationType: LocationType_DATACENTER,
	}

	if subnet.Id != "test-id" {
		t.Errorf("Expected Id to be 'test-id', got '%s'", subnet.Id)
	}

	if subnet.Cidr != "10.0.0.0/24" {
		t.Errorf("Expected Cidr to be '10.0.0.0/24', got '%s'", subnet.Cidr)
	}

	// Test CloudInfo
	cloudInfo := &CloudInfo{
		Provider:  "aws",
		Region:    "us-east-1",
		AccountId: "123456789",
	}

	if cloudInfo.Provider != "aws" {
		t.Errorf("Expected Provider to be 'aws', got '%s'", cloudInfo.Provider)
	}

	// Test SubnetDetails
	details := &SubnetDetails{
		Address:     "10.0.0.0",
		Netmask:     "255.255.255.0",
		Wildcard:    "0.0.0.255",
		Network:     "10.0.0.0",
		Type:        "private",
		Broadcast:   "10.0.0.255",
		HostMin:     "10.0.0.1",
		HostMax:     "10.0.0.254",
		HostsPerNet: 254,
		IsPublic:    false,
	}

	if details.HostsPerNet != 254 {
		t.Errorf("Expected HostsPerNet to be 254, got %d", details.HostsPerNet)
	}
}

// TestRequestResponseMessages verifies API request/response messages
func TestRequestResponseMessages(t *testing.T) {
	// Test CreateSubnetRequest
	createReq := &CreateSubnetRequest{
		Cidr:         "10.0.0.0/24",
		Name:         "Test Subnet",
		Description:  "Test description",
		Location:     "datacenter-1",
		LocationType: LocationType_DATACENTER,
	}

	if createReq.Cidr != "10.0.0.0/24" {
		t.Errorf("Expected Cidr to be '10.0.0.0/24', got '%s'", createReq.Cidr)
	}

	// Test ListSubnetsRequest
	listReq := &ListSubnetsRequest{
		LocationFilter:      "datacenter-1",
		CloudProviderFilter: "aws",
		SearchQuery:         "test",
		Page:                1,
		PageSize:            10,
	}

	if listReq.Page != 1 {
		t.Errorf("Expected Page to be 1, got %d", listReq.Page)
	}

	// Test GetSubnetRequest
	getReq := &GetSubnetRequest{
		Id: "test-id",
	}

	if getReq.Id != "test-id" {
		t.Errorf("Expected Id to be 'test-id', got '%s'", getReq.Id)
	}
}

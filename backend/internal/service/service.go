package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bananaops/ipam-bananaops/internal/repository"
	pb "github.com/bananaops/ipam-bananaops/proto"
	"github.com/google/uuid"
)

// IPService defines the interface for IP calculations
type IPService interface {
	CalculateSubnetDetails(cidr string) (*pb.SubnetDetails, error)
	ValidateCIDR(cidr string) error
}

// CloudProviderManager defines the interface for cloud provider operations
type CloudProviderManager interface {
	// Future implementation for cloud provider integration
}

// ServiceLayer implements the business logic using Protobuf messages
type ServiceLayer struct {
	subnetRepo   repository.SubnetRepository
	ipService    IPService
	cloudManager CloudProviderManager
}

// NewServiceLayer creates a new service layer instance
func NewServiceLayer(repo repository.SubnetRepository, ipService IPService, cloudManager CloudProviderManager) *ServiceLayer {
	return &ServiceLayer{
		subnetRepo:   repo,
		ipService:    ipService,
		cloudManager: cloudManager,
	}
}

// CreateSubnet creates a new subnet with calculated properties
func (s *ServiceLayer) CreateSubnet(ctx context.Context, req *pb.CreateSubnetRequest) (*pb.CreateSubnetResponse, error) {
	// Validate CIDR
	if err := s.ipService.ValidateCIDR(req.Cidr); err != nil {
		return &pb.CreateSubnetResponse{
			Error: &pb.Error{
				Code:      "INVALID_CIDR",
				Message:   fmt.Sprintf("Invalid CIDR notation: %v", err),
				Timestamp: time.Now().Unix(),
			},
		}, nil
	}

	// Calculate subnet details
	details, err := s.ipService.CalculateSubnetDetails(req.Cidr)
	if err != nil {
		return &pb.CreateSubnetResponse{
			Error: &pb.Error{
				Code:      "CALCULATION_ERROR",
				Message:   fmt.Sprintf("Failed to calculate subnet details: %v", err),
				Timestamp: time.Now().Unix(),
			},
		}, nil
	}

	// Create subnet object
	subnet := &pb.Subnet{
		Id:           uuid.New().String(),
		Cidr:         req.Cidr,
		Name:         req.Name,
		Description:  req.Description,
		Location:     req.Location,
		LocationType: req.LocationType,
		CloudInfo:    req.CloudInfo,
		Details:      details,
		Utilization: &pb.UtilizationInfo{
			TotalIps:           details.HostsPerNet,
			AllocatedIps:       0,
			UtilizationPercent: 0.0,
		},
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Persist to repository
	if err := s.subnetRepo.Create(ctx, subnet); err != nil {
		return &pb.CreateSubnetResponse{
			Error: &pb.Error{
				Code:      "DB_ERROR",
				Message:   fmt.Sprintf("Failed to create subnet: %v", err),
				Timestamp: time.Now().Unix(),
			},
		}, nil
	}

	return &pb.CreateSubnetResponse{
		Subnet: subnet,
	}, nil
}

// ListSubnets retrieves subnets with optional filtering
func (s *ServiceLayer) ListSubnets(ctx context.Context, req *pb.ListSubnetsRequest) (*pb.ListSubnetsResponse, error) {
	// Build filters from request
	filters := &repository.SubnetFilters{
		LocationFilter:      req.LocationFilter,
		CloudProviderFilter: req.CloudProviderFilter,
		SearchQuery:         req.SearchQuery,
		Page:                req.Page,
		PageSize:            req.PageSize,
	}

	// Set default page size if not specified
	if filters.PageSize == 0 {
		filters.PageSize = 50
	}

	// Retrieve subnets from repository
	subnets, err := s.subnetRepo.FindAll(ctx, filters)
	if err != nil {
		return &pb.ListSubnetsResponse{
			Error: &pb.Error{
				Code:      "DB_ERROR",
				Message:   fmt.Sprintf("Failed to retrieve subnets: %v", err),
				Timestamp: time.Now().Unix(),
			},
		}, nil
	}

	return &pb.ListSubnetsResponse{
		Subnets:    subnets,
		TotalCount: int32(len(subnets)),
	}, nil
}

// GetSubnet retrieves a specific subnet by ID
func (s *ServiceLayer) GetSubnet(ctx context.Context, req *pb.GetSubnetRequest) (*pb.GetSubnetResponse, error) {
	if req.Id == "" {
		return &pb.GetSubnetResponse{
			Error: &pb.Error{
				Code:      "INVALID_REQUEST",
				Message:   "Subnet ID is required",
				Timestamp: time.Now().Unix(),
			},
		}, nil
	}

	subnet, err := s.subnetRepo.FindByID(ctx, req.Id)
	if err != nil {
		return &pb.GetSubnetResponse{
			Error: &pb.Error{
				Code:      "SUBNET_NOT_FOUND",
				Message:   fmt.Sprintf("Subnet not found: %v", err),
				Timestamp: time.Now().Unix(),
			},
		}, nil
	}

	return &pb.GetSubnetResponse{
		Subnet: subnet,
	}, nil
}

// UpdateSubnet updates an existing subnet and recalculates properties if CIDR changed
func (s *ServiceLayer) UpdateSubnet(ctx context.Context, req *pb.UpdateSubnetRequest) (*pb.UpdateSubnetResponse, error) {
	if req.Id == "" {
		return &pb.UpdateSubnetResponse{
			Error: &pb.Error{
				Code:      "INVALID_REQUEST",
				Message:   "Subnet ID is required",
				Timestamp: time.Now().Unix(),
			},
		}, nil
	}

	// Retrieve existing subnet
	existing, err := s.subnetRepo.FindByID(ctx, req.Id)
	if err != nil {
		return &pb.UpdateSubnetResponse{
			Error: &pb.Error{
				Code:      "SUBNET_NOT_FOUND",
				Message:   fmt.Sprintf("Subnet not found: %v", err),
				Timestamp: time.Now().Unix(),
			},
		}, nil
	}

	// Check if CIDR changed and recalculate if needed
	var details *pb.SubnetDetails
	if req.Cidr != "" && req.Cidr != existing.Cidr {
		// Validate new CIDR
		if err := s.ipService.ValidateCIDR(req.Cidr); err != nil {
			return &pb.UpdateSubnetResponse{
				Error: &pb.Error{
					Code:      "INVALID_CIDR",
					Message:   fmt.Sprintf("Invalid CIDR notation: %v", err),
					Timestamp: time.Now().Unix(),
				},
			}, nil
		}

		// Recalculate subnet details
		details, err = s.ipService.CalculateSubnetDetails(req.Cidr)
		if err != nil {
			return &pb.UpdateSubnetResponse{
				Error: &pb.Error{
					Code:      "CALCULATION_ERROR",
					Message:   fmt.Sprintf("Failed to calculate subnet details: %v", err),
					Timestamp: time.Now().Unix(),
				},
			}, nil
		}
		existing.Cidr = req.Cidr
		existing.Details = details

		// Update utilization with new total IPs
		existing.Utilization.TotalIps = details.HostsPerNet
		if existing.Utilization.AllocatedIps > 0 {
			existing.Utilization.UtilizationPercent = float32(existing.Utilization.AllocatedIps) / float32(details.HostsPerNet) * 100
		}
	}

	// Update other fields
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Location != "" {
		existing.Location = req.Location
	}
	if req.LocationType != 0 {
		existing.LocationType = req.LocationType
	}
	if req.CloudInfo != nil {
		existing.CloudInfo = req.CloudInfo
	}

	existing.UpdatedAt = time.Now().Unix()

	// Persist changes
	if err := s.subnetRepo.Update(ctx, existing); err != nil {
		return &pb.UpdateSubnetResponse{
			Error: &pb.Error{
				Code:      "DB_ERROR",
				Message:   fmt.Sprintf("Failed to update subnet: %v", err),
				Timestamp: time.Now().Unix(),
			},
		}, nil
	}

	return &pb.UpdateSubnetResponse{
		Subnet: existing,
	}, nil
}

// DeleteSubnet removes a subnet from the system
func (s *ServiceLayer) DeleteSubnet(ctx context.Context, req *pb.DeleteSubnetRequest) (*pb.DeleteSubnetResponse, error) {
	if req.Id == "" {
		return &pb.DeleteSubnetResponse{
			Success: false,
			Error: &pb.Error{
				Code:      "INVALID_REQUEST",
				Message:   "Subnet ID is required",
				Timestamp: time.Now().Unix(),
			},
		}, nil
	}

	// Check if subnet exists
	_, err := s.subnetRepo.FindByID(ctx, req.Id)
	if err != nil {
		return &pb.DeleteSubnetResponse{
			Success: false,
			Error: &pb.Error{
				Code:      "SUBNET_NOT_FOUND",
				Message:   fmt.Sprintf("Subnet not found: %v", err),
				Timestamp: time.Now().Unix(),
			},
		}, nil
	}

	// Delete subnet
	if err := s.subnetRepo.Delete(ctx, req.Id); err != nil {
		return &pb.DeleteSubnetResponse{
			Success: false,
			Error: &pb.Error{
				Code:      "DB_ERROR",
				Message:   fmt.Sprintf("Failed to delete subnet: %v", err),
				Timestamp: time.Now().Unix(),
			},
		}, nil
	}

	return &pb.DeleteSubnetResponse{
		Success: true,
	}, nil
}

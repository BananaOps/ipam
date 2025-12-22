package aws

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bananaops/ipam-bananaops/internal/repository"
	"github.com/google/uuid"
)

// SyncService handles synchronization of AWS resources with IPAM
type SyncService struct {
	client     *Client
	repository repository.SubnetRepository
}

// NewSyncService creates a new AWS sync service
func NewSyncService(client *Client, repo repository.SubnetRepository) *SyncService {
	return &SyncService{
		client:     client,
		repository: repo,
	}
}

// SyncVPCs synchronizes VPCs from AWS to IPAM
func (s *SyncService) SyncVPCs(ctx context.Context) error {
	log.Printf("Starting VPC synchronization for region: %s", s.client.GetRegion())

	vpcs, err := s.client.ListVPCs(ctx)
	if err != nil {
		return fmt.Errorf("failed to list VPCs: %w", err)
	}

	log.Printf("Found %d VPCs in AWS", len(vpcs))

	for _, vpc := range vpcs {
		// Check if VPC already exists in IPAM
		existingSubnet, err := s.repository.GetSubnetByCIDR(ctx, vpc.CIDR)
		if err == nil && existingSubnet != nil {
			log.Printf("VPC %s (%s) already exists in IPAM, skipping", vpc.ID, vpc.CIDR)
			continue
		}

		// Create subnet entry for VPC
		subnet := &repository.Subnet{
			ID:           uuid.New().String(),
			Name:         fmt.Sprintf("VPC-%s", vpc.Name),
			CIDR:         vpc.CIDR,
			Location:     vpc.Region,
			LocationType: "cloud",
			CloudInfo: &repository.CloudInfo{
				Provider:     "aws",
				Region:       vpc.Region,
				AccountID:    "", // Will be populated if available
				ResourceType: "vpc",
				VPCId:        vpc.ID,
				SubnetId:     "", // Empty for VPC entries
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Add tags as metadata
		if len(vpc.Tags) > 0 {
			subnet.Tags = vpc.Tags
		}

		err = s.repository.CreateSubnet(ctx, subnet)
		if err != nil {
			log.Printf("Failed to create VPC %s in IPAM: %v", vpc.ID, err)
			continue
		}

		log.Printf("Successfully synchronized VPC %s (%s) to IPAM", vpc.ID, vpc.CIDR)
	}

	return nil
}

// SyncSubnets synchronizes subnets from AWS to IPAM
func (s *SyncService) SyncSubnets(ctx context.Context) error {
	log.Printf("Starting subnet synchronization for region: %s", s.client.GetRegion())

	subnets, err := s.client.ListSubnets(ctx)
	if err != nil {
		return fmt.Errorf("failed to list subnets: %w", err)
	}

	log.Printf("Found %d subnets in AWS", len(subnets))

	for _, awsSubnet := range subnets {
		// Check if subnet already exists in IPAM
		existingSubnet, err := s.repository.GetSubnetByCIDR(ctx, awsSubnet.CIDR)
		if err == nil && existingSubnet != nil {
			// Update existing subnet with AWS information
			existingSubnet.CloudInfo = &repository.CloudInfo{
				Provider:     "aws",
				Region:       awsSubnet.Region,
				AccountID:    "", // Will be populated if available
				ResourceType: "subnet",
				VPCId:        awsSubnet.VPCId,
				SubnetId:     awsSubnet.ID,
			}
			existingSubnet.Location = awsSubnet.Region
			existingSubnet.LocationType = "cloud"
			existingSubnet.UpdatedAt = time.Now()

			// Find parent VPC
			if parentVPC, err := s.findParentVPC(ctx, awsSubnet.VPCId); err == nil && parentVPC != nil {
				existingSubnet.ParentID = parentVPC.ID
			}

			if len(awsSubnet.Tags) > 0 {
				existingSubnet.Tags = awsSubnet.Tags
			}

			err = s.repository.UpdateSubnet(ctx, existingSubnet.ID, existingSubnet)
			if err != nil {
				log.Printf("Failed to update subnet %s in IPAM: %v", awsSubnet.ID, err)
				continue
			}

			log.Printf("Updated existing subnet %s (%s) with AWS information", awsSubnet.ID, awsSubnet.CIDR)
			continue
		}

		// Get utilization
		utilization, err := s.client.GetSubnetUtilization(ctx, awsSubnet.ID)
		if err != nil {
			log.Printf("Failed to get utilization for subnet %s: %v", awsSubnet.ID, err)
			utilization = 0 // Default to 0 if we can't get utilization
		}

		// Create new subnet entry
		subnet := &repository.Subnet{
			ID:           uuid.New().String(),
			Name:         awsSubnet.Name,
			CIDR:         awsSubnet.CIDR,
			Location:     awsSubnet.Region,
			LocationType: "cloud",
			CloudInfo: &repository.CloudInfo{
				Provider:     "aws",
				Region:       awsSubnet.Region,
				AccountID:    "", // Will be populated if available
				ResourceType: "subnet",
				VPCId:        awsSubnet.VPCId,
				SubnetId:     awsSubnet.ID,
			},
			Utilization: &repository.Utilization{
				UtilizationPercent: utilization,
				LastUpdated:        time.Now(),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Find parent VPC
		if parentVPC, err := s.findParentVPC(ctx, awsSubnet.VPCId); err == nil && parentVPC != nil {
			subnet.ParentID = parentVPC.ID
		}

		// Add tags as metadata
		if len(awsSubnet.Tags) > 0 {
			subnet.Tags = awsSubnet.Tags
		}

		err = s.repository.CreateSubnet(ctx, subnet)
		if err != nil {
			log.Printf("Failed to create subnet %s in IPAM: %v", awsSubnet.ID, err)
			continue
		}

		log.Printf("Successfully synchronized subnet %s (%s) to IPAM", awsSubnet.ID, awsSubnet.CIDR)
	}

	return nil
}

// SyncAll synchronizes both VPCs and subnets
func (s *SyncService) SyncAll(ctx context.Context) error {
	log.Printf("Starting full AWS synchronization for region: %s", s.client.GetRegion())

	// First sync VPCs
	if err := s.SyncVPCs(ctx); err != nil {
		return fmt.Errorf("failed to sync VPCs: %w", err)
	}

	// Then sync subnets
	if err := s.SyncSubnets(ctx); err != nil {
		return fmt.Errorf("failed to sync subnets: %w", err)
	}

	log.Printf("Successfully completed AWS synchronization for region: %s", s.client.GetRegion())
	return nil
}

// UpdateUtilization updates utilization data for all AWS subnets
func (s *SyncService) UpdateUtilization(ctx context.Context) error {
	log.Printf("Updating utilization for AWS subnets in region: %s", s.client.GetRegion())

	// Get all subnets with AWS cloud info
	subnets, err := s.repository.ListSubnets(ctx, repository.SubnetFilters{
		CloudProvider: "aws",
	})
	if err != nil {
		return fmt.Errorf("failed to list AWS subnets: %w", err)
	}

	for _, subnet := range subnets.Subnets {
		if subnet.CloudInfo == nil || subnet.CloudInfo.SubnetId == "" {
			continue // Skip VPC entries or subnets without AWS subnet ID
		}

		utilization, err := s.client.GetSubnetUtilization(ctx, subnet.CloudInfo.SubnetId)
		if err != nil {
			log.Printf("Failed to get utilization for subnet %s: %v", subnet.CloudInfo.SubnetId, err)
			continue
		}

		// Update utilization
		subnet.Utilization = &repository.Utilization{
			UtilizationPercent: utilization,
			LastUpdated:        time.Now(),
		}
		subnet.UpdatedAt = time.Now()

		err = s.repository.UpdateSubnet(ctx, subnet.ID, subnet)
		if err != nil {
			log.Printf("Failed to update utilization for subnet %s: %v", subnet.ID, err)
			continue
		}

		log.Printf("Updated utilization for subnet %s: %.2f%%", subnet.CloudInfo.SubnetId, utilization)
	}

	return nil
}

// findParentVPC finds the parent VPC for a given VPC ID
func (s *SyncService) findParentVPC(ctx context.Context, vpcID string) (*repository.Subnet, error) {
	// List all subnets with AWS cloud info
	subnets, err := s.repository.ListSubnets(ctx, repository.SubnetFilters{
		CloudProvider: "aws",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list subnets: %w", err)
	}

	// Find the VPC entry
	for _, subnet := range subnets.Subnets {
		if subnet.CloudInfo != nil &&
			subnet.CloudInfo.ResourceType == "vpc" &&
			subnet.CloudInfo.VPCId == vpcID {
			return subnet, nil
		}
	}

	return nil, fmt.Errorf("VPC %s not found", vpcID)
}

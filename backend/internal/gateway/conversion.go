package gateway

import (
	"encoding/json"
	"fmt"

	"github.com/bananaops/ipam-bananaops/internal/repository"
	pb "github.com/bananaops/ipam-bananaops/proto"
)

// JSON request/response types for REST API

// CreateSubnetJSON represents the JSON request for creating a subnet
type CreateSubnetJSON struct {
	CIDR         string         `json:"cidr"`
	Name         string         `json:"name"`
	Description  string         `json:"description,omitempty"`
	Location     string         `json:"location,omitempty"`
	LocationType string         `json:"location_type,omitempty"`
	CloudInfo    *CloudInfoJSON `json:"cloud_info,omitempty"`
}

// UpdateSubnetJSON represents the JSON request for updating a subnet
type UpdateSubnetJSON struct {
	CIDR         string         `json:"cidr,omitempty"`
	Name         string         `json:"name,omitempty"`
	Description  string         `json:"description,omitempty"`
	Location     string         `json:"location,omitempty"`
	LocationType string         `json:"location_type,omitempty"`
	CloudInfo    *CloudInfoJSON `json:"cloud_info,omitempty"`
}

// CloudInfoJSON represents cloud provider information in JSON
type CloudInfoJSON struct {
	Provider     string `json:"provider"`
	Region       string `json:"region"`
	AccountID    string `json:"account_id"`
	ResourceType string `json:"resource_type,omitempty"`
	VPCId        string `json:"vpc_id,omitempty"`
	SubnetId     string `json:"subnet_id,omitempty"`
}

// SubnetJSON represents a subnet in JSON format
type SubnetJSON struct {
	ID           string             `json:"id"`
	CIDR         string             `json:"cidr"`
	Name         string             `json:"name"`
	Description  string             `json:"description,omitempty"`
	Location     string             `json:"location,omitempty"`
	LocationType string             `json:"location_type"`
	CloudInfo    *CloudInfoJSON     `json:"cloud_info,omitempty"`
	Details      *SubnetDetailsJSON `json:"details,omitempty"`
	Utilization  *UtilizationJSON   `json:"utilization,omitempty"`
	ParentID     string             `json:"parent_id,omitempty"`
	CreatedAt    int64              `json:"created_at"`
	UpdatedAt    int64              `json:"updated_at"`
}

// SubnetDetailsJSON represents subnet details in JSON format
type SubnetDetailsJSON struct {
	Address     string `json:"address"`
	Netmask     string `json:"netmask"`
	Wildcard    string `json:"wildcard"`
	Network     string `json:"network"`
	Type        string `json:"type"`
	Broadcast   string `json:"broadcast"`
	HostMin     string `json:"host_min"`
	HostMax     string `json:"host_max"`
	HostsPerNet int32  `json:"hosts_per_net"`
	IsPublic    bool   `json:"is_public"`
}

// UtilizationJSON represents utilization info in JSON format
type UtilizationJSON struct {
	TotalIPs           int32   `json:"total_ips"`
	AllocatedIPs       int32   `json:"allocated_ips"`
	UtilizationPercent float32 `json:"utilization_percent"`
}

// ListSubnetsResponseJSON represents the list subnets response in JSON
type ListSubnetsResponseJSON struct {
	Subnets    []*SubnetJSON `json:"subnets"`
	TotalCount int32         `json:"total_count"`
}

// ErrorResponse represents an error response in JSON
type ErrorResponse struct {
	Error *ErrorDetail `json:"error"`
}

// ErrorDetail represents error details in JSON
type ErrorDetail struct {
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Details   map[string]string `json:"details,omitempty"`
	Timestamp int64             `json:"timestamp"`
}

// DeleteResponseJSON represents the delete response in JSON
type DeleteResponseJSON struct {
	Success bool `json:"success"`
}

// JSONToCreateSubnetRequest converts JSON to Protobuf CreateSubnetRequest
func JSONToCreateSubnetRequest(data []byte) (*pb.CreateSubnetRequest, error) {
	var jsonReq CreateSubnetJSON
	if err := json.Unmarshal(data, &jsonReq); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	req := &pb.CreateSubnetRequest{
		Cidr:         jsonReq.CIDR,
		Name:         jsonReq.Name,
		Description:  jsonReq.Description,
		Location:     jsonReq.Location,
		LocationType: stringToLocationType(jsonReq.LocationType),
	}

	if jsonReq.CloudInfo != nil {
		req.CloudInfo = &pb.CloudInfo{
			Provider:  jsonReq.CloudInfo.Provider,
			Region:    jsonReq.CloudInfo.Region,
			AccountId: jsonReq.CloudInfo.AccountID,
		}
	}

	return req, nil
}

// JSONToUpdateSubnetRequest converts JSON to Protobuf UpdateSubnetRequest
func JSONToUpdateSubnetRequest(id string, data []byte) (*pb.UpdateSubnetRequest, error) {
	var jsonReq UpdateSubnetJSON
	if err := json.Unmarshal(data, &jsonReq); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	req := &pb.UpdateSubnetRequest{
		Id:           id,
		Cidr:         jsonReq.CIDR,
		Name:         jsonReq.Name,
		Description:  jsonReq.Description,
		Location:     jsonReq.Location,
		LocationType: stringToLocationType(jsonReq.LocationType),
	}

	if jsonReq.CloudInfo != nil {
		req.CloudInfo = &pb.CloudInfo{
			Provider:  jsonReq.CloudInfo.Provider,
			Region:    jsonReq.CloudInfo.Region,
			AccountId: jsonReq.CloudInfo.AccountID,
		}
	}

	return req, nil
}

// SubnetToJSON converts a Protobuf Subnet to JSON format
func SubnetToJSON(subnet *pb.Subnet) *SubnetJSON {
	if subnet == nil {
		return nil
	}

	result := &SubnetJSON{
		ID:           subnet.Id,
		CIDR:         subnet.Cidr,
		Name:         subnet.Name,
		Description:  subnet.Description,
		Location:     subnet.Location,
		LocationType: locationTypeToString(subnet.LocationType),
		CreatedAt:    subnet.CreatedAt,
		UpdatedAt:    subnet.UpdatedAt,
	}

	if subnet.CloudInfo != nil && subnet.CloudInfo.Provider != "" {
		result.CloudInfo = &CloudInfoJSON{
			Provider:     subnet.CloudInfo.Provider,
			Region:       subnet.CloudInfo.Region,
			AccountID:    subnet.CloudInfo.AccountId,
			ResourceType: "", // Will be populated from repository model
			VPCId:        "", // Will be populated from repository model
			SubnetId:     "", // Will be populated from repository model
		}
	}

	if subnet.Details != nil {
		result.Details = &SubnetDetailsJSON{
			Address:     subnet.Details.Address,
			Netmask:     subnet.Details.Netmask,
			Wildcard:    subnet.Details.Wildcard,
			Network:     subnet.Details.Network,
			Type:        subnet.Details.Type,
			Broadcast:   subnet.Details.Broadcast,
			HostMin:     subnet.Details.HostMin,
			HostMax:     subnet.Details.HostMax,
			HostsPerNet: subnet.Details.HostsPerNet,
			IsPublic:    subnet.Details.IsPublic,
		}
	}

	if subnet.Utilization != nil {
		result.Utilization = &UtilizationJSON{
			TotalIPs:           subnet.Utilization.TotalIps,
			AllocatedIPs:       subnet.Utilization.AllocatedIps,
			UtilizationPercent: subnet.Utilization.UtilizationPercent,
		}
	}

	return result
}

// SubnetsToJSON converts a slice of Protobuf Subnets to JSON format
func SubnetsToJSON(subnets []*pb.Subnet) []*SubnetJSON {
	result := make([]*SubnetJSON, len(subnets))
	for i, subnet := range subnets {
		result[i] = SubnetToJSON(subnet)
	}
	return result
}

// stringToLocationType converts a string to LocationType enum
func stringToLocationType(s string) pb.LocationType {
	switch s {
	case "DATACENTER", "datacenter":
		return pb.LocationType_DATACENTER
	case "SITE", "site":
		return pb.LocationType_SITE
	case "CLOUD", "cloud":
		return pb.LocationType_CLOUD
	default:
		return pb.LocationType_DATACENTER
	}
}

// locationTypeToString converts LocationType enum to string
func locationTypeToString(lt pb.LocationType) string {
	switch lt {
	case pb.LocationType_DATACENTER:
		return "DATACENTER"
	case pb.LocationType_SITE:
		return "SITE"
	case pb.LocationType_CLOUD:
		return "CLOUD"
	default:
		return "DATACENTER"
	}
}

// Repository model conversion functions

// RepositorySubnetToJSON converts a repository Subnet to JSON format
func RepositorySubnetToJSON(subnet *repository.Subnet) *SubnetJSON {
	if subnet == nil {
		return nil
	}

	result := &SubnetJSON{
		ID:           subnet.ID,
		CIDR:         subnet.CIDR,
		Name:         subnet.Name,
		Location:     subnet.Location,
		LocationType: subnet.LocationType,
		ParentID:     subnet.ParentID,
		CreatedAt:    subnet.CreatedAt.Unix(),
		UpdatedAt:    subnet.UpdatedAt.Unix(),
	}

	if subnet.CloudInfo != nil && subnet.CloudInfo.Provider != "" {
		result.CloudInfo = &CloudInfoJSON{
			Provider:     subnet.CloudInfo.Provider,
			Region:       subnet.CloudInfo.Region,
			AccountID:    subnet.CloudInfo.AccountID,
			ResourceType: subnet.CloudInfo.ResourceType,
			VPCId:        subnet.CloudInfo.VPCId,
			SubnetId:     subnet.CloudInfo.SubnetId,
		}
	}

	if subnet.Details != nil {
		result.Details = &SubnetDetailsJSON{
			Address:     subnet.Details.Address,
			Netmask:     subnet.Details.Netmask,
			Wildcard:    subnet.Details.Wildcard,
			Network:     subnet.Details.Network,
			Type:        subnet.Details.Type,
			Broadcast:   subnet.Details.Broadcast,
			HostMin:     subnet.Details.HostMin,
			HostMax:     subnet.Details.HostMax,
			HostsPerNet: subnet.Details.HostsPerNet,
			IsPublic:    subnet.Details.IsPublic,
		}
	}

	if subnet.Utilization != nil {
		result.Utilization = &UtilizationJSON{
			TotalIPs:           subnet.Utilization.TotalIPs,
			AllocatedIPs:       subnet.Utilization.AllocatedIPs,
			UtilizationPercent: float32(subnet.Utilization.UtilizationPercent),
		}
	}

	return result
}

// RepositorySubnetsToJSON converts a slice of repository Subnets to JSON format
func RepositorySubnetsToJSON(subnets []*repository.Subnet) []*SubnetJSON {
	result := make([]*SubnetJSON, len(subnets))
	for i, subnet := range subnets {
		result[i] = RepositorySubnetToJSON(subnet)
	}
	return result
}

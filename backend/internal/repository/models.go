package repository

import (
	"time"
)

// Subnet represents a subnet in the repository layer
type Subnet struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	CIDR         string            `json:"cidr"`
	Location     string            `json:"location"`
	LocationType string            `json:"location_type"`
	CloudInfo    *CloudInfo        `json:"cloud_info,omitempty"`
	Details      *SubnetDetails    `json:"details,omitempty"`
	Utilization  *Utilization      `json:"utilization,omitempty"`
	Tags         map[string]string `json:"tags,omitempty"`
	ParentID     string            `json:"parent_id,omitempty"` // ID du réseau parent
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// SubnetDetails represents calculated subnet information
type SubnetDetails struct {
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

// CloudInfo represents cloud provider information
type CloudInfo struct {
	Provider     string `json:"provider"`
	Region       string `json:"region"`
	AccountID    string `json:"account_id"`
	ResourceType string `json:"resource_type,omitempty"` // "vpc" ou "subnet"
	VPCId        string `json:"vpc_id,omitempty"`
	SubnetId     string `json:"subnet_id,omitempty"`
}

// Utilization represents subnet utilization information
type Utilization struct {
	TotalIPs           int32     `json:"total_ips"`
	AllocatedIPs       int32     `json:"allocated_ips"`
	UtilizationPercent float64   `json:"utilization_percent"`
	LastUpdated        time.Time `json:"last_updated"`
}

// SubnetFilters contains filtering criteria for subnet queries
type SubnetFilters struct {
	LocationFilter      string
	CloudProviderFilter string
	SearchQuery         string
	Page                int32
	PageSize            int32
	CloudProvider       string // For cloud provider specific filtering
}

// SubnetList represents a list of subnets with pagination
type SubnetList struct {
	Subnets    []*Subnet `json:"subnets"`
	TotalCount int32     `json:"total_count"`
}

// Connection represents a connection between subnets
type Connection struct {
	ID             string                 `json:"id"`
	SourceSubnetID string                 `json:"source_subnet_id"`
	TargetSubnetID string                 `json:"target_subnet_id"`
	ConnectionType string                 `json:"connection_type"`
	Status         string                 `json:"status"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description,omitempty"`
	Bandwidth      string                 `json:"bandwidth,omitempty"`
	Latency        int32                  `json:"latency,omitempty"`
	Cost           float64                `json:"cost,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// ConnectionFilters contains filtering criteria for connection queries
type ConnectionFilters struct {
	SourceSubnetID string
	TargetSubnetID string
	ConnectionType string
	Status         string
	Page           int32
	PageSize       int32
}

// ConnectionList represents a list of connections with pagination
type ConnectionList struct {
	Connections []*Connection `json:"connections"`
	TotalCount  int32         `json:"total_count"`
}

package repository

import (
	"context"

	pb "github.com/bananaops/ipam-bananaops/proto"
)

// SubnetRepository defines the interface for subnet data access
type SubnetRepository interface {
	Create(ctx context.Context, subnet *pb.Subnet) error
	FindByID(ctx context.Context, id string) (*pb.Subnet, error)
	FindAll(ctx context.Context, filters *SubnetFilters) ([]*pb.Subnet, error)
	Update(ctx context.Context, subnet *pb.Subnet) error
	Delete(ctx context.Context, id string) error
	Close() error
}

// SubnetFilters contains filtering criteria for subnet queries
type SubnetFilters struct {
	LocationFilter      string
	CloudProviderFilter string
	SearchQuery         string
	Page                int32
	PageSize            int32
}

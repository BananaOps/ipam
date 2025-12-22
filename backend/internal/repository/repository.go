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

	// Extended methods for cloud provider integration
	CreateSubnet(ctx context.Context, subnet *Subnet) error
	GetSubnetByCIDR(ctx context.Context, cidr string) (*Subnet, error)
	GetSubnetByID(ctx context.Context, id string) (*Subnet, error)
	UpdateSubnet(ctx context.Context, id string, subnet *Subnet) error
	ListSubnets(ctx context.Context, filters SubnetFilters) (*SubnetList, error)
	GetSubnetChildren(ctx context.Context, parentID string) ([]*Subnet, error)

	// Connection methods
	CreateConnection(ctx context.Context, connection *Connection) error
	GetConnectionByID(ctx context.Context, id string) (*Connection, error)
	UpdateConnection(ctx context.Context, id string, connection *Connection) error
	DeleteConnection(ctx context.Context, id string) error
	ListConnections(ctx context.Context, filters ConnectionFilters) (*ConnectionList, error)
}

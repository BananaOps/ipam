package repository

import (
	"context"
	"fmt"
	"time"

	pb "github.com/bananaops/ipam-bananaops/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBRepository implements SubnetRepository using MongoDB
type MongoDBRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// subnetDocument represents the MongoDB document structure
type subnetDocument struct {
	ID           string                 `bson:"_id"`
	CIDR         string                 `bson:"cidr"`
	Name         string                 `bson:"name"`
	Description  string                 `bson:"description"`
	Location     string                 `bson:"location"`
	LocationType string                 `bson:"locationType"`
	CloudInfo    *cloudInfoDocument     `bson:"cloudInfo,omitempty"`
	Details      *subnetDetailsDocument `bson:"details"`
	Utilization  *utilizationDocument   `bson:"utilization"`
	CreatedAt    int64                  `bson:"createdAt"`
	UpdatedAt    int64                  `bson:"updatedAt"`
}

type cloudInfoDocument struct {
	Provider  string `bson:"provider"`
	Region    string `bson:"region"`
	AccountID string `bson:"accountId"`
}

type subnetDetailsDocument struct {
	Address     string `bson:"address"`
	Netmask     string `bson:"netmask"`
	Wildcard    string `bson:"wildcard"`
	Network     string `bson:"network"`
	Type        string `bson:"type"`
	Broadcast   string `bson:"broadcast"`
	HostMin     string `bson:"hostMin"`
	HostMax     string `bson:"hostMax"`
	HostsPerNet int32  `bson:"hostsPerNet"`
	IsPublic    bool   `bson:"isPublic"`
}

type utilizationDocument struct {
	TotalIPs           int32   `bson:"totalIps"`
	AllocatedIPs       int32   `bson:"allocatedIps"`
	UtilizationPercent float32 `bson:"utilizationPercent"`
}

// NewMongoDBRepository creates a new MongoDB repository
func NewMongoDBRepository(connectionString string) (*MongoDBRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Get collection
	collection := client.Database("ipam").Collection("subnets")

	repo := &MongoDBRepository{
		client:     client,
		collection: collection,
	}

	// Create indexes
	if err := repo.createIndexes(ctx); err != nil {
		client.Disconnect(ctx)
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return repo, nil
}

// createIndexes creates necessary indexes for the collection
func (r *MongoDBRepository) createIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "location", Value: 1}},
			Options: options.Index().SetName("idx_location"),
		},
		{
			Keys:    bson.D{{Key: "cloudInfo.provider", Value: 1}},
			Options: options.Index().SetName("idx_cloud_provider"),
		},
		{
			Keys:    bson.D{{Key: "cidr", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_cidr_unique"),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// Create inserts a new subnet into the database
func (r *MongoDBRepository) Create(ctx context.Context, subnet *pb.Subnet) error {
	doc := r.toDocument(subnet)

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("subnet with CIDR %s already exists", subnet.Cidr)
		}
		return fmt.Errorf("failed to create subnet: %w", err)
	}

	return nil
}

// FindByID retrieves a subnet by its ID
func (r *MongoDBRepository) FindByID(ctx context.Context, id string) (*pb.Subnet, error) {
	filter := bson.M{"_id": id}

	var doc subnetDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("subnet not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find subnet: %w", err)
	}

	return r.toProto(&doc), nil
}

// FindAll retrieves all subnets with optional filtering
func (r *MongoDBRepository) FindAll(ctx context.Context, filters *SubnetFilters) ([]*pb.Subnet, error) {
	filter := bson.M{}

	// Apply filters
	if filters != nil {
		if filters.LocationFilter != "" {
			filter["location"] = bson.M{"$regex": filters.LocationFilter, "$options": "i"}
		}
		if filters.CloudProviderFilter != "" {
			filter["cloudInfo.provider"] = filters.CloudProviderFilter
		}
		if filters.SearchQuery != "" {
			filter["$or"] = []bson.M{
				{"name": bson.M{"$regex": filters.SearchQuery, "$options": "i"}},
				{"cidr": bson.M{"$regex": filters.SearchQuery, "$options": "i"}},
				{"description": bson.M{"$regex": filters.SearchQuery, "$options": "i"}},
				{"location": bson.M{"$regex": filters.SearchQuery, "$options": "i"}},
			}
		}
	}

	// Set up options - sort by CIDR for logical IP address ordering
	opts := options.Find().SetSort(bson.D{{Key: "cidr", Value: 1}})

	// Apply pagination
	if filters != nil && filters.PageSize > 0 {
		opts.SetLimit(int64(filters.PageSize))
		opts.SetSkip(int64(filters.Page * filters.PageSize))
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query subnets: %w", err)
	}
	defer cursor.Close(ctx)

	var subnets []*pb.Subnet
	for cursor.Next(ctx) {
		var doc subnetDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode subnet: %w", err)
		}
		subnets = append(subnets, r.toProto(&doc))
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return subnets, nil
}

// Update modifies an existing subnet
func (r *MongoDBRepository) Update(ctx context.Context, subnet *pb.Subnet) error {
	filter := bson.M{"_id": subnet.Id}
	doc := r.toDocument(subnet)

	// Remove _id from update document
	update := bson.M{"$set": doc}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update subnet: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("subnet not found")
	}

	return nil
}

// Delete removes a subnet from the database
func (r *MongoDBRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete subnet: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("subnet not found")
	}

	return nil
}

// Close closes the database connection
func (r *MongoDBRepository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.client.Disconnect(ctx)
}

// Connection methods - Not implemented for MongoDB yet
func (r *MongoDBRepository) CreateConnection(ctx context.Context, connection *Connection) error {
	return fmt.Errorf("connection methods not implemented for MongoDB repository")
}

func (r *MongoDBRepository) GetConnectionByID(ctx context.Context, id string) (*Connection, error) {
	return nil, fmt.Errorf("connection methods not implemented for MongoDB repository")
}

func (r *MongoDBRepository) UpdateConnection(ctx context.Context, id string, connection *Connection) error {
	return fmt.Errorf("connection methods not implemented for MongoDB repository")
}

func (r *MongoDBRepository) DeleteConnection(ctx context.Context, id string) error {
	return fmt.Errorf("connection methods not implemented for MongoDB repository")
}

func (r *MongoDBRepository) ListConnections(ctx context.Context, filters ConnectionFilters) (*ConnectionList, error) {
	return nil, fmt.Errorf("connection methods not implemented for MongoDB repository")
}

// toDocument converts a Protobuf Subnet to a MongoDB document
func (r *MongoDBRepository) toDocument(subnet *pb.Subnet) *subnetDocument {
	doc := &subnetDocument{
		ID:           subnet.Id,
		CIDR:         subnet.Cidr,
		Name:         subnet.Name,
		Description:  subnet.Description,
		Location:     subnet.Location,
		LocationType: subnet.LocationType.String(),
		CreatedAt:    subnet.CreatedAt,
		UpdatedAt:    subnet.UpdatedAt,
	}

	if subnet.CloudInfo != nil {
		doc.CloudInfo = &cloudInfoDocument{
			Provider:  subnet.CloudInfo.Provider,
			Region:    subnet.CloudInfo.Region,
			AccountID: subnet.CloudInfo.AccountId,
		}
	}

	if subnet.Details != nil {
		doc.Details = &subnetDetailsDocument{
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
		doc.Utilization = &utilizationDocument{
			TotalIPs:           subnet.Utilization.TotalIps,
			AllocatedIPs:       subnet.Utilization.AllocatedIps,
			UtilizationPercent: subnet.Utilization.UtilizationPercent,
		}
	}

	return doc
}

// toProto converts a MongoDB document to a Protobuf Subnet
func (r *MongoDBRepository) toProto(doc *subnetDocument) *pb.Subnet {
	subnet := &pb.Subnet{
		Id:           doc.ID,
		Cidr:         doc.CIDR,
		Name:         doc.Name,
		Description:  doc.Description,
		Location:     doc.Location,
		LocationType: parseLocationType(doc.LocationType),
		CreatedAt:    doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
	}

	if doc.CloudInfo != nil {
		subnet.CloudInfo = &pb.CloudInfo{
			Provider:  doc.CloudInfo.Provider,
			Region:    doc.CloudInfo.Region,
			AccountId: doc.CloudInfo.AccountID,
		}
	}

	if doc.Details != nil {
		subnet.Details = &pb.SubnetDetails{
			Address:     doc.Details.Address,
			Netmask:     doc.Details.Netmask,
			Wildcard:    doc.Details.Wildcard,
			Network:     doc.Details.Network,
			Type:        doc.Details.Type,
			Broadcast:   doc.Details.Broadcast,
			HostMin:     doc.Details.HostMin,
			HostMax:     doc.Details.HostMax,
			HostsPerNet: doc.Details.HostsPerNet,
			IsPublic:    doc.Details.IsPublic,
		}
	}

	if doc.Utilization != nil {
		subnet.Utilization = &pb.UtilizationInfo{
			TotalIps:           doc.Utilization.TotalIPs,
			AllocatedIps:       doc.Utilization.AllocatedIPs,
			UtilizationPercent: doc.Utilization.UtilizationPercent,
		}
	}

	return subnet
}

// Extended methods for cloud provider integration

// CreateSubnet creates a new subnet using the repository model
func (r *MongoDBRepository) CreateSubnet(ctx context.Context, subnet *Subnet) error {
	doc := r.toRepositoryDocument(subnet)

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to create subnet: %w", err)
	}

	return nil
}

// GetSubnetByCIDR retrieves a subnet by its CIDR
func (r *MongoDBRepository) GetSubnetByCIDR(ctx context.Context, cidr string) (*Subnet, error) {
	filter := bson.M{"cidr": cidr}

	var doc subnetRepositoryDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("subnet not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find subnet: %w", err)
	}

	return r.fromRepositoryDocument(&doc), nil
}

// UpdateSubnet updates an existing subnet using the repository model
func (r *MongoDBRepository) UpdateSubnet(ctx context.Context, id string, subnet *Subnet) error {
	filter := bson.M{"_id": id}
	doc := r.toRepositoryDocument(subnet)

	// Remove _id from update document
	update := bson.M{"$set": doc}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update subnet: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("subnet not found")
	}

	return nil
}

// ListSubnets retrieves subnets with filtering using the repository model
func (r *MongoDBRepository) ListSubnets(ctx context.Context, filters SubnetFilters) (*SubnetList, error) {
	filter := bson.M{}

	// Apply filters
	if filters.LocationFilter != "" {
		filter["location"] = bson.M{"$regex": filters.LocationFilter, "$options": "i"}
	}
	if filters.CloudProviderFilter != "" {
		filter["cloudInfo.provider"] = filters.CloudProviderFilter
	}
	if filters.CloudProvider != "" {
		filter["cloudInfo.provider"] = filters.CloudProvider
	}
	if filters.SearchQuery != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": filters.SearchQuery, "$options": "i"}},
			{"cidr": bson.M{"$regex": filters.SearchQuery, "$options": "i"}},
			{"description": bson.M{"$regex": filters.SearchQuery, "$options": "i"}},
			{"location": bson.M{"$regex": filters.SearchQuery, "$options": "i"}},
		}
	}

	// Count total records
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count subnets: %w", err)
	}

	// Build find options
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"createdAt", -1}})

	// Apply pagination
	if filters.PageSize > 0 {
		findOptions.SetLimit(int64(filters.PageSize))
		findOptions.SetSkip(int64(filters.Page * filters.PageSize))
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to query subnets: %w", err)
	}
	defer cursor.Close(ctx)

	var subnets []*Subnet
	for cursor.Next(ctx) {
		var doc subnetRepositoryDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode subnet: %w", err)
		}
		subnets = append(subnets, r.fromRepositoryDocument(&doc))
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return &SubnetList{
		Subnets:    subnets,
		TotalCount: int32(totalCount),
	}, nil
}

// subnetRepositoryDocument represents the MongoDB document structure for repository model
type subnetRepositoryDocument struct {
	ID           string                           `bson:"_id"`
	CIDR         string                           `bson:"cidr"`
	Name         string                           `bson:"name"`
	Location     string                           `bson:"location"`
	LocationType string                           `bson:"locationType"`
	CloudInfo    *cloudInfoRepositoryDocument     `bson:"cloudInfo,omitempty"`
	Details      *subnetDetailsRepositoryDocument `bson:"details,omitempty"`
	Utilization  *utilizationRepositoryDocument   `bson:"utilization,omitempty"`
	Tags         map[string]string                `bson:"tags,omitempty"`
	ParentID     string                           `bson:"parentId,omitempty"`
	CreatedAt    int64                            `bson:"createdAt"`
	UpdatedAt    int64                            `bson:"updatedAt"`
}

type cloudInfoRepositoryDocument struct {
	Provider     string `bson:"provider"`
	Region       string `bson:"region"`
	AccountID    string `bson:"accountId"`
	ResourceType string `bson:"resourceType,omitempty"`
	VPCId        string `bson:"vpcId,omitempty"`
	SubnetId     string `bson:"subnetId,omitempty"`
}

type subnetDetailsRepositoryDocument struct {
	Address     string `bson:"address"`
	Netmask     string `bson:"netmask"`
	Wildcard    string `bson:"wildcard"`
	Network     string `bson:"network"`
	Type        string `bson:"type"`
	Broadcast   string `bson:"broadcast"`
	HostMin     string `bson:"hostMin"`
	HostMax     string `bson:"hostMax"`
	HostsPerNet int32  `bson:"hostsPerNet"`
	IsPublic    bool   `bson:"isPublic"`
}

type utilizationRepositoryDocument struct {
	TotalIPs           int32   `bson:"totalIps"`
	AllocatedIPs       int32   `bson:"allocatedIps"`
	UtilizationPercent float64 `bson:"utilizationPercent"`
	LastUpdated        int64   `bson:"lastUpdated"`
}

// toRepositoryDocument converts a repository Subnet to a MongoDB document
func (r *MongoDBRepository) toRepositoryDocument(subnet *Subnet) *subnetRepositoryDocument {
	doc := &subnetRepositoryDocument{
		ID:           subnet.ID,
		CIDR:         subnet.CIDR,
		Name:         subnet.Name,
		Location:     subnet.Location,
		LocationType: subnet.LocationType,
		Tags:         subnet.Tags,
		ParentID:     subnet.ParentID,
		CreatedAt:    subnet.CreatedAt.Unix(),
		UpdatedAt:    subnet.UpdatedAt.Unix(),
	}

	if subnet.CloudInfo != nil {
		doc.CloudInfo = &cloudInfoRepositoryDocument{
			Provider:     subnet.CloudInfo.Provider,
			Region:       subnet.CloudInfo.Region,
			AccountID:    subnet.CloudInfo.AccountID,
			ResourceType: subnet.CloudInfo.ResourceType,
			VPCId:        subnet.CloudInfo.VPCId,
			SubnetId:     subnet.CloudInfo.SubnetId,
		}
	}

	if subnet.Details != nil {
		doc.Details = &subnetDetailsRepositoryDocument{
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
		doc.Utilization = &utilizationRepositoryDocument{
			TotalIPs:           subnet.Utilization.TotalIPs,
			AllocatedIPs:       subnet.Utilization.AllocatedIPs,
			UtilizationPercent: subnet.Utilization.UtilizationPercent,
			LastUpdated:        subnet.Utilization.LastUpdated.Unix(),
		}
	}

	return doc
}

// fromRepositoryDocument converts a MongoDB document to a repository Subnet
func (r *MongoDBRepository) fromRepositoryDocument(doc *subnetRepositoryDocument) *Subnet {
	subnet := &Subnet{
		ID:           doc.ID,
		CIDR:         doc.CIDR,
		Name:         doc.Name,
		Location:     doc.Location,
		LocationType: doc.LocationType,
		Tags:         doc.Tags,
		ParentID:     doc.ParentID,
		CreatedAt:    time.Unix(doc.CreatedAt, 0),
		UpdatedAt:    time.Unix(doc.UpdatedAt, 0),
	}

	if doc.CloudInfo != nil {
		subnet.CloudInfo = &CloudInfo{
			Provider:     doc.CloudInfo.Provider,
			Region:       doc.CloudInfo.Region,
			AccountID:    doc.CloudInfo.AccountID,
			ResourceType: doc.CloudInfo.ResourceType,
			VPCId:        doc.CloudInfo.VPCId,
			SubnetId:     doc.CloudInfo.SubnetId,
		}
	}

	if doc.Details != nil {
		subnet.Details = &SubnetDetails{
			Address:     doc.Details.Address,
			Netmask:     doc.Details.Netmask,
			Wildcard:    doc.Details.Wildcard,
			Network:     doc.Details.Network,
			Type:        doc.Details.Type,
			Broadcast:   doc.Details.Broadcast,
			HostMin:     doc.Details.HostMin,
			HostMax:     doc.Details.HostMax,
			HostsPerNet: doc.Details.HostsPerNet,
			IsPublic:    doc.Details.IsPublic,
		}
	}

	if doc.Utilization != nil {
		subnet.Utilization = &Utilization{
			TotalIPs:           doc.Utilization.TotalIPs,
			AllocatedIPs:       doc.Utilization.AllocatedIPs,
			UtilizationPercent: doc.Utilization.UtilizationPercent,
			LastUpdated:        time.Unix(doc.Utilization.LastUpdated, 0),
		}
	}

	return subnet
}

// GetSubnetChildren retrieves child subnets for a given parent subnet ID
func (r *MongoDBRepository) GetSubnetChildren(ctx context.Context, parentID string) ([]*Subnet, error) {
	filter := bson.M{"parentId": parentID}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"cidr", 1}})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to query child subnets: %w", err)
	}
	defer cursor.Close(ctx)

	var subnets []*Subnet
	for cursor.Next(ctx) {
		var doc subnetRepositoryDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode child subnet: %w", err)
		}
		subnets = append(subnets, r.fromRepositoryDocument(&doc))
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return subnets, nil
}

// GetSubnetByID retrieves a subnet by its ID using repository models
func (r *MongoDBRepository) GetSubnetByID(ctx context.Context, id string) (*Subnet, error) {
	filter := bson.M{"_id": id}

	var doc subnetRepositoryDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("subnet not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find subnet: %w", err)
	}

	return r.fromRepositoryDocument(&doc), nil
}

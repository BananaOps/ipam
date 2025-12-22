package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	pb "github.com/bananaops/ipam-bananaops/proto"
	_ "modernc.org/sqlite"
)

// SQLiteRepository implements SubnetRepository using SQLite
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite repository
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	repo := &SQLiteRepository{db: db}

	// Initialize schema
	if err := repo.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

// initSchema creates the database schema
func (r *SQLiteRepository) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS subnets (
		id TEXT PRIMARY KEY,
		cidr TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		location TEXT,
		location_type TEXT,
		cloud_provider TEXT,
		cloud_region TEXT,
		cloud_account_id TEXT,
		cloud_resource_type TEXT,
		cloud_vpc_id TEXT,
		cloud_subnet_id TEXT,
		parent_id TEXT,
		address TEXT,
		netmask TEXT,
		wildcard TEXT,
		network TEXT,
		type TEXT,
		broadcast TEXT,
		host_min TEXT,
		host_max TEXT,
		hosts_per_net INTEGER,
		is_public INTEGER,
		total_ips INTEGER,
		allocated_ips INTEGER,
		utilization_percent REAL,
		created_at INTEGER,
		updated_at INTEGER,
		FOREIGN KEY (parent_id) REFERENCES subnets(id)
	);

	CREATE TABLE IF NOT EXISTS connections (
		id TEXT PRIMARY KEY,
		source_subnet_id TEXT NOT NULL,
		target_subnet_id TEXT NOT NULL,
		connection_type TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'active',
		name TEXT NOT NULL,
		description TEXT,
		bandwidth TEXT,
		latency INTEGER,
		cost REAL,
		metadata TEXT, -- JSON string for additional metadata
		created_at INTEGER,
		updated_at INTEGER,
		FOREIGN KEY (source_subnet_id) REFERENCES subnets(id) ON DELETE CASCADE,
		FOREIGN KEY (target_subnet_id) REFERENCES subnets(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_subnets_location ON subnets(location);
	CREATE INDEX IF NOT EXISTS idx_subnets_cloud_provider ON subnets(cloud_provider);
	CREATE INDEX IF NOT EXISTS idx_subnets_cidr ON subnets(cidr);
	CREATE INDEX IF NOT EXISTS idx_subnets_parent_id ON subnets(parent_id);
	CREATE INDEX IF NOT EXISTS idx_subnets_cloud_resource_type ON subnets(cloud_resource_type);
	
	CREATE INDEX IF NOT EXISTS idx_connections_source ON connections(source_subnet_id);
	CREATE INDEX IF NOT EXISTS idx_connections_target ON connections(target_subnet_id);
	CREATE INDEX IF NOT EXISTS idx_connections_type ON connections(connection_type);
	CREATE INDEX IF NOT EXISTS idx_connections_status ON connections(status);
	`

	_, err := r.db.Exec(schema)
	return err
}

// Create inserts a new subnet into the database
func (r *SQLiteRepository) Create(ctx context.Context, subnet *pb.Subnet) error {
	query := `
		INSERT INTO subnets (
			id, cidr, name, description, location, location_type,
			cloud_provider, cloud_region, cloud_account_id,
			address, netmask, wildcard, network, type, broadcast,
			host_min, host_max, hosts_per_net, is_public,
			total_ips, allocated_ips, utilization_percent,
			created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?,
			?, ?
		)
	`

	cloudProvider := ""
	cloudRegion := ""
	cloudAccountID := ""
	if subnet.CloudInfo != nil {
		cloudProvider = subnet.CloudInfo.Provider
		cloudRegion = subnet.CloudInfo.Region
		cloudAccountID = subnet.CloudInfo.AccountId
	}

	isPublic := 0
	if subnet.Details != nil && subnet.Details.IsPublic {
		isPublic = 1
	}

	_, err := r.db.ExecContext(ctx, query,
		subnet.Id, subnet.Cidr, subnet.Name, subnet.Description,
		subnet.Location, subnet.LocationType.String(),
		cloudProvider, cloudRegion, cloudAccountID,
		subnet.Details.Address, subnet.Details.Netmask, subnet.Details.Wildcard,
		subnet.Details.Network, subnet.Details.Type, subnet.Details.Broadcast,
		subnet.Details.HostMin, subnet.Details.HostMax, subnet.Details.HostsPerNet,
		isPublic,
		subnet.Utilization.TotalIps, subnet.Utilization.AllocatedIps,
		subnet.Utilization.UtilizationPercent,
		subnet.CreatedAt, subnet.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create subnet: %w", err)
	}

	return nil
}

// FindByID retrieves a subnet by its ID
func (r *SQLiteRepository) FindByID(ctx context.Context, id string) (*pb.Subnet, error) {
	query := `
		SELECT 
			id, cidr, name, description, location, location_type,
			cloud_provider, cloud_region, cloud_account_id,
			address, netmask, wildcard, network, type, broadcast,
			host_min, host_max, hosts_per_net, is_public,
			total_ips, allocated_ips, utilization_percent,
			created_at, updated_at
		FROM subnets
		WHERE id = ?
	`

	var subnet pb.Subnet
	var locationType string
	var cloudProvider, cloudRegion, cloudAccountID sql.NullString
	var isPublic int

	subnet.Details = &pb.SubnetDetails{}
	subnet.Utilization = &pb.UtilizationInfo{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subnet.Id, &subnet.Cidr, &subnet.Name, &subnet.Description,
		&subnet.Location, &locationType,
		&cloudProvider, &cloudRegion, &cloudAccountID,
		&subnet.Details.Address, &subnet.Details.Netmask, &subnet.Details.Wildcard,
		&subnet.Details.Network, &subnet.Details.Type, &subnet.Details.Broadcast,
		&subnet.Details.HostMin, &subnet.Details.HostMax, &subnet.Details.HostsPerNet,
		&isPublic,
		&subnet.Utilization.TotalIps, &subnet.Utilization.AllocatedIps,
		&subnet.Utilization.UtilizationPercent,
		&subnet.CreatedAt, &subnet.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subnet not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find subnet: %w", err)
	}

	// Parse location type
	subnet.LocationType = parseLocationType(locationType)

	// Parse cloud info
	if cloudProvider.Valid {
		subnet.CloudInfo = &pb.CloudInfo{
			Provider:  cloudProvider.String,
			Region:    cloudRegion.String,
			AccountId: cloudAccountID.String,
		}
	}

	subnet.Details.IsPublic = isPublic == 1

	return &subnet, nil
}

// FindAll retrieves all subnets with optional filtering
func (r *SQLiteRepository) FindAll(ctx context.Context, filters *SubnetFilters) ([]*pb.Subnet, error) {
	query := `
		SELECT 
			id, cidr, name, description, location, location_type,
			cloud_provider, cloud_region, cloud_account_id,
			address, netmask, wildcard, network, type, broadcast,
			host_min, host_max, hosts_per_net, is_public,
			total_ips, allocated_ips, utilization_percent,
			created_at, updated_at
		FROM subnets
		WHERE 1=1
	`

	args := []interface{}{}

	// Apply filters
	if filters != nil {
		if filters.LocationFilter != "" {
			query += " AND location LIKE ?"
			locationPattern := "%" + filters.LocationFilter + "%"
			args = append(args, locationPattern)
		}
		if filters.CloudProviderFilter != "" {
			query += " AND cloud_provider = ?"
			args = append(args, filters.CloudProviderFilter)
		}
		if filters.SearchQuery != "" {
			query += " AND (name LIKE ? OR cidr LIKE ? OR description LIKE ?)"
			searchPattern := "%" + filters.SearchQuery + "%"
			args = append(args, searchPattern, searchPattern, searchPattern)
		}
	}

	query += " ORDER BY created_at DESC"

	// Apply pagination
	if filters != nil && filters.PageSize > 0 {
		query += " LIMIT ? OFFSET ?"
		offset := filters.Page * filters.PageSize
		args = append(args, filters.PageSize, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query subnets: %w", err)
	}
	defer rows.Close()

	var subnets []*pb.Subnet

	for rows.Next() {
		var subnet pb.Subnet
		var locationType string
		var cloudProvider, cloudRegion, cloudAccountID sql.NullString
		var isPublic int

		subnet.Details = &pb.SubnetDetails{}
		subnet.Utilization = &pb.UtilizationInfo{}

		err := rows.Scan(
			&subnet.Id, &subnet.Cidr, &subnet.Name, &subnet.Description,
			&subnet.Location, &locationType,
			&cloudProvider, &cloudRegion, &cloudAccountID,
			&subnet.Details.Address, &subnet.Details.Netmask, &subnet.Details.Wildcard,
			&subnet.Details.Network, &subnet.Details.Type, &subnet.Details.Broadcast,
			&subnet.Details.HostMin, &subnet.Details.HostMax, &subnet.Details.HostsPerNet,
			&isPublic,
			&subnet.Utilization.TotalIps, &subnet.Utilization.AllocatedIps,
			&subnet.Utilization.UtilizationPercent,
			&subnet.CreatedAt, &subnet.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subnet: %w", err)
		}

		// Parse location type
		subnet.LocationType = parseLocationType(locationType)

		// Parse cloud info
		if cloudProvider.Valid {
			subnet.CloudInfo = &pb.CloudInfo{
				Provider:  cloudProvider.String,
				Region:    cloudRegion.String,
				AccountId: cloudAccountID.String,
			}
		}

		subnet.Details.IsPublic = isPublic == 1

		subnets = append(subnets, &subnet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return subnets, nil
}

// Update modifies an existing subnet
func (r *SQLiteRepository) Update(ctx context.Context, subnet *pb.Subnet) error {
	query := `
		UPDATE subnets SET
			cidr = ?, name = ?, description = ?, location = ?, location_type = ?,
			cloud_provider = ?, cloud_region = ?, cloud_account_id = ?,
			address = ?, netmask = ?, wildcard = ?, network = ?, type = ?, broadcast = ?,
			host_min = ?, host_max = ?, hosts_per_net = ?, is_public = ?,
			total_ips = ?, allocated_ips = ?, utilization_percent = ?,
			updated_at = ?
		WHERE id = ?
	`

	cloudProvider := ""
	cloudRegion := ""
	cloudAccountID := ""
	if subnet.CloudInfo != nil {
		cloudProvider = subnet.CloudInfo.Provider
		cloudRegion = subnet.CloudInfo.Region
		cloudAccountID = subnet.CloudInfo.AccountId
	}

	isPublic := 0
	if subnet.Details != nil && subnet.Details.IsPublic {
		isPublic = 1
	}

	result, err := r.db.ExecContext(ctx, query,
		subnet.Cidr, subnet.Name, subnet.Description,
		subnet.Location, subnet.LocationType.String(),
		cloudProvider, cloudRegion, cloudAccountID,
		subnet.Details.Address, subnet.Details.Netmask, subnet.Details.Wildcard,
		subnet.Details.Network, subnet.Details.Type, subnet.Details.Broadcast,
		subnet.Details.HostMin, subnet.Details.HostMax, subnet.Details.HostsPerNet,
		isPublic,
		subnet.Utilization.TotalIps, subnet.Utilization.AllocatedIps,
		subnet.Utilization.UtilizationPercent,
		subnet.UpdatedAt,
		subnet.Id,
	)

	if err != nil {
		return fmt.Errorf("failed to update subnet: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subnet not found")
	}

	return nil
}

// Delete removes a subnet from the database
func (r *SQLiteRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM subnets WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subnet: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subnet not found")
	}

	return nil
}

// Close closes the database connection
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

// Connection methods

// CreateConnection inserts a new connection into the database
func (r *SQLiteRepository) CreateConnection(ctx context.Context, connection *Connection) error {
	query := `
		INSERT INTO connections (
			id, source_subnet_id, target_subnet_id, connection_type, status,
			name, description, bandwidth, latency, cost, metadata,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	metadataJSON := ""
	if connection.Metadata != nil {
		// Convert metadata to JSON string
		// For simplicity, we'll skip JSON marshaling for now
		// In a real implementation, you'd use json.Marshal
	}

	_, err := r.db.ExecContext(ctx, query,
		connection.ID,
		connection.SourceSubnetID,
		connection.TargetSubnetID,
		connection.ConnectionType,
		connection.Status,
		connection.Name,
		connection.Description,
		connection.Bandwidth,
		connection.Latency,
		connection.Cost,
		metadataJSON,
		connection.CreatedAt.Unix(),
		connection.UpdatedAt.Unix(),
	)

	return err
}

// GetConnectionByID retrieves a connection by its ID
func (r *SQLiteRepository) GetConnectionByID(ctx context.Context, id string) (*Connection, error) {
	query := `
		SELECT id, source_subnet_id, target_subnet_id, connection_type, status,
			   name, description, bandwidth, latency, cost, metadata,
			   created_at, updated_at
		FROM connections
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)

	connection := &Connection{}
	var metadataJSON string
	var createdAt, updatedAt int64

	err := row.Scan(
		&connection.ID,
		&connection.SourceSubnetID,
		&connection.TargetSubnetID,
		&connection.ConnectionType,
		&connection.Status,
		&connection.Name,
		&connection.Description,
		&connection.Bandwidth,
		&connection.Latency,
		&connection.Cost,
		&metadataJSON,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("connection not found")
		}
		return nil, err
	}

	connection.CreatedAt = time.Unix(createdAt, 0)
	connection.UpdatedAt = time.Unix(updatedAt, 0)

	// Parse metadata JSON if needed
	if metadataJSON != "" {
		// In a real implementation, you'd use json.Unmarshal
		connection.Metadata = make(map[string]interface{})
	}

	return connection, nil
}

// UpdateConnection updates an existing connection
func (r *SQLiteRepository) UpdateConnection(ctx context.Context, id string, connection *Connection) error {
	query := `
		UPDATE connections SET
			source_subnet_id = ?, target_subnet_id = ?, connection_type = ?, status = ?,
			name = ?, description = ?, bandwidth = ?, latency = ?, cost = ?,
			metadata = ?, updated_at = ?
		WHERE id = ?
	`

	metadataJSON := ""
	if connection.Metadata != nil {
		// Convert metadata to JSON string
		// For simplicity, we'll skip JSON marshaling for now
	}

	result, err := r.db.ExecContext(ctx, query,
		connection.SourceSubnetID,
		connection.TargetSubnetID,
		connection.ConnectionType,
		connection.Status,
		connection.Name,
		connection.Description,
		connection.Bandwidth,
		connection.Latency,
		connection.Cost,
		metadataJSON,
		time.Now().Unix(),
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("connection not found")
	}

	return nil
}

// DeleteConnection removes a connection from the database
func (r *SQLiteRepository) DeleteConnection(ctx context.Context, id string) error {
	query := `DELETE FROM connections WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("connection not found")
	}

	return nil
}

// ListConnections retrieves connections with optional filtering
func (r *SQLiteRepository) ListConnections(ctx context.Context, filters ConnectionFilters) (*ConnectionList, error) {
	// Build WHERE clause
	var conditions []string
	var args []interface{}

	if filters.SourceSubnetID != "" {
		conditions = append(conditions, "source_subnet_id = ?")
		args = append(args, filters.SourceSubnetID)
	}

	if filters.TargetSubnetID != "" {
		conditions = append(conditions, "target_subnet_id = ?")
		args = append(args, filters.TargetSubnetID)
	}

	if filters.ConnectionType != "" {
		conditions = append(conditions, "connection_type = ?")
		args = append(args, filters.ConnectionType)
	}

	if filters.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filters.Status)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM connections %s", whereClause)
	var totalCount int32
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	// Build main query with pagination
	query := fmt.Sprintf(`
		SELECT id, source_subnet_id, target_subnet_id, connection_type, status,
			   name, description, bandwidth, latency, cost, metadata,
			   created_at, updated_at
		FROM connections
		%s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	// Add pagination parameters
	limit := filters.PageSize
	if limit <= 0 {
		limit = 50 // Default page size
	}
	offset := filters.Page * limit

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []*Connection
	for rows.Next() {
		connection := &Connection{}
		var metadataJSON string
		var createdAt, updatedAt int64

		err := rows.Scan(
			&connection.ID,
			&connection.SourceSubnetID,
			&connection.TargetSubnetID,
			&connection.ConnectionType,
			&connection.Status,
			&connection.Name,
			&connection.Description,
			&connection.Bandwidth,
			&connection.Latency,
			&connection.Cost,
			&metadataJSON,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		connection.CreatedAt = time.Unix(createdAt, 0)
		connection.UpdatedAt = time.Unix(updatedAt, 0)

		// Parse metadata JSON if needed
		if metadataJSON != "" {
			connection.Metadata = make(map[string]interface{})
		}

		connections = append(connections, connection)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &ConnectionList{
		Connections: connections,
		TotalCount:  totalCount,
	}, nil
}

// parseLocationType converts a string to LocationType enum
func parseLocationType(s string) pb.LocationType {
	s = strings.ToUpper(s)
	switch s {
	case "DATACENTER":
		return pb.LocationType_DATACENTER
	case "SITE":
		return pb.LocationType_SITE
	case "CLOUD":
		return pb.LocationType_CLOUD
	default:
		return pb.LocationType_DATACENTER
	}
}

// Extended methods for cloud provider integration

// CreateSubnet creates a new subnet using the repository model
func (r *SQLiteRepository) CreateSubnet(ctx context.Context, subnet *Subnet) error {
	query := `
		INSERT INTO subnets (
			id, cidr, name, description, location, location_type,
			cloud_provider, cloud_region, cloud_account_id, cloud_resource_type, cloud_vpc_id, cloud_subnet_id,
			parent_id, address, netmask, wildcard, network, type, broadcast,
			host_min, host_max, hosts_per_net, is_public,
			total_ips, allocated_ips, utilization_percent, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	cloudProvider := ""
	cloudRegion := ""
	cloudAccountID := ""
	cloudResourceType := ""
	cloudVPCId := ""
	cloudSubnetId := ""
	if subnet.CloudInfo != nil {
		cloudProvider = subnet.CloudInfo.Provider
		cloudRegion = subnet.CloudInfo.Region
		cloudAccountID = subnet.CloudInfo.AccountID
		cloudResourceType = subnet.CloudInfo.ResourceType
		cloudVPCId = subnet.CloudInfo.VPCId
		cloudSubnetId = subnet.CloudInfo.SubnetId
	}

	// Subnet details
	address := ""
	netmask := ""
	wildcard := ""
	network := ""
	subnetType := ""
	broadcast := ""
	hostMin := ""
	hostMax := ""
	var hostsPerNet int32 = 0
	isPublic := 0
	if subnet.Details != nil {
		address = subnet.Details.Address
		netmask = subnet.Details.Netmask
		wildcard = subnet.Details.Wildcard
		network = subnet.Details.Network
		subnetType = subnet.Details.Type
		broadcast = subnet.Details.Broadcast
		hostMin = subnet.Details.HostMin
		hostMax = subnet.Details.HostMax
		hostsPerNet = subnet.Details.HostsPerNet
		if subnet.Details.IsPublic {
			isPublic = 1
		}
	}

	// Utilization
	var totalIPs int32 = 0
	var allocatedIPs int32 = 0
	utilizationPercent := 0.0
	if subnet.Utilization != nil {
		totalIPs = subnet.Utilization.TotalIPs
		allocatedIPs = subnet.Utilization.AllocatedIPs
		utilizationPercent = subnet.Utilization.UtilizationPercent
	}

	_, err := r.db.ExecContext(ctx, query,
		subnet.ID, subnet.CIDR, subnet.Name, "",
		subnet.Location, subnet.LocationType,
		cloudProvider, cloudRegion, cloudAccountID, cloudResourceType, cloudVPCId, cloudSubnetId,
		subnet.ParentID, address, netmask, wildcard, network, subnetType, broadcast,
		hostMin, hostMax, hostsPerNet, isPublic,
		totalIPs, allocatedIPs, utilizationPercent,
		subnet.CreatedAt.Unix(), subnet.UpdatedAt.Unix(),
	)

	if err != nil {
		return fmt.Errorf("failed to create subnet: %w", err)
	}

	return nil
}

// GetSubnetByCIDR retrieves a subnet by its CIDR
func (r *SQLiteRepository) GetSubnetByCIDR(ctx context.Context, cidr string) (*Subnet, error) {
	query := `
		SELECT 
			id, cidr, name, description, location, location_type,
			cloud_provider, cloud_region, cloud_account_id, cloud_resource_type, cloud_vpc_id, cloud_subnet_id,
			parent_id, utilization_percent, created_at, updated_at
		FROM subnets
		WHERE cidr = ?
	`

	var subnet Subnet
	var description sql.NullString
	var cloudProvider, cloudRegion, cloudAccountID, cloudResourceType, cloudVPCId, cloudSubnetId sql.NullString
	var parentID sql.NullString
	var utilizationPercent sql.NullFloat64
	var createdAt, updatedAt int64

	err := r.db.QueryRowContext(ctx, query, cidr).Scan(
		&subnet.ID, &subnet.CIDR, &subnet.Name, &description,
		&subnet.Location, &subnet.LocationType,
		&cloudProvider, &cloudRegion, &cloudAccountID, &cloudResourceType, &cloudVPCId, &cloudSubnetId,
		&parentID, &utilizationPercent, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subnet not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find subnet: %w", err)
	}

	// Parse cloud info
	if cloudProvider.Valid {
		subnet.CloudInfo = &CloudInfo{
			Provider:     cloudProvider.String,
			Region:       cloudRegion.String,
			AccountID:    cloudAccountID.String,
			ResourceType: cloudResourceType.String,
			VPCId:        cloudVPCId.String,
			SubnetId:     cloudSubnetId.String,
		}
	}

	// Parse utilization
	if utilizationPercent.Valid {
		subnet.Utilization = &Utilization{
			UtilizationPercent: utilizationPercent.Float64,
			LastUpdated:        time.Unix(updatedAt, 0),
		}
	}

	if parentID.Valid {
		subnet.ParentID = parentID.String
	}

	subnet.CreatedAt = time.Unix(createdAt, 0)
	subnet.UpdatedAt = time.Unix(updatedAt, 0)

	return &subnet, nil
}

// UpdateSubnet updates an existing subnet using the repository model
func (r *SQLiteRepository) UpdateSubnet(ctx context.Context, id string, subnet *Subnet) error {
	query := `
		UPDATE subnets SET
			cidr = ?, name = ?, location = ?, location_type = ?,
			cloud_provider = ?, cloud_region = ?, cloud_account_id = ?,
			utilization_percent = ?, updated_at = ?
		WHERE id = ?
	`

	cloudProvider := ""
	cloudRegion := ""
	cloudAccountID := ""
	if subnet.CloudInfo != nil {
		cloudProvider = subnet.CloudInfo.Provider
		cloudRegion = subnet.CloudInfo.Region
		cloudAccountID = subnet.CloudInfo.AccountID
	}

	utilizationPercent := 0.0
	if subnet.Utilization != nil {
		utilizationPercent = subnet.Utilization.UtilizationPercent
	}

	result, err := r.db.ExecContext(ctx, query,
		subnet.CIDR, subnet.Name, subnet.Location, subnet.LocationType,
		cloudProvider, cloudRegion, cloudAccountID,
		utilizationPercent, subnet.UpdatedAt.Unix(),
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update subnet: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subnet not found")
	}

	return nil
}

// ListSubnets retrieves subnets with filtering using the repository model
func (r *SQLiteRepository) ListSubnets(ctx context.Context, filters SubnetFilters) (*SubnetList, error) {
	baseQuery := `
		SELECT 
			id, cidr, name, description, location, location_type,
			cloud_provider, cloud_region, cloud_account_id, cloud_resource_type, cloud_vpc_id, cloud_subnet_id,
			parent_id, utilization_percent, created_at, updated_at
		FROM subnets
		WHERE 1=1
	`

	whereClause := ""
	args := []interface{}{}

	// Apply filters
	if filters.LocationFilter != "" {
		whereClause += " AND location LIKE ?"
		args = append(args, "%"+filters.LocationFilter+"%")
	}
	if filters.CloudProviderFilter != "" {
		whereClause += " AND cloud_provider = ?"
		args = append(args, filters.CloudProviderFilter)
	}
	if filters.CloudProvider != "" {
		whereClause += " AND cloud_provider = ?"
		args = append(args, filters.CloudProvider)
	}
	if filters.SearchQuery != "" {
		whereClause += " AND (name LIKE ? OR cidr LIKE ?)"
		searchPattern := "%" + filters.SearchQuery + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// Count total records
	countQuery := "SELECT COUNT(*) FROM subnets WHERE 1=1" + whereClause
	var totalCount int32
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count subnets: %w", err)
	}

	// Build final query
	finalQuery := baseQuery + whereClause + " ORDER BY created_at DESC"

	// Apply pagination
	if filters.PageSize > 0 {
		finalQuery += " LIMIT ? OFFSET ?"
		offset := filters.Page * filters.PageSize
		args = append(args, filters.PageSize, offset)
	}

	rows, err := r.db.QueryContext(ctx, finalQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query subnets: %w", err)
	}
	defer rows.Close()

	var subnets []*Subnet

	for rows.Next() {
		var subnet Subnet
		var description sql.NullString
		var cloudProvider, cloudRegion, cloudAccountID, cloudResourceType, cloudVPCId, cloudSubnetId sql.NullString
		var parentID sql.NullString
		var utilizationPercent sql.NullFloat64
		var createdAt, updatedAt int64

		err := rows.Scan(
			&subnet.ID, &subnet.CIDR, &subnet.Name, &description,
			&subnet.Location, &subnet.LocationType,
			&cloudProvider, &cloudRegion, &cloudAccountID, &cloudResourceType, &cloudVPCId, &cloudSubnetId,
			&parentID, &utilizationPercent, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subnet: %w", err)
		}

		// Parse cloud info
		if cloudProvider.Valid {
			subnet.CloudInfo = &CloudInfo{
				Provider:     cloudProvider.String,
				Region:       cloudRegion.String,
				AccountID:    cloudAccountID.String,
				ResourceType: cloudResourceType.String,
				VPCId:        cloudVPCId.String,
				SubnetId:     cloudSubnetId.String,
			}
		}

		// Parse utilization
		if utilizationPercent.Valid {
			subnet.Utilization = &Utilization{
				UtilizationPercent: utilizationPercent.Float64,
				LastUpdated:        time.Unix(updatedAt, 0),
			}
		}

		if parentID.Valid {
			subnet.ParentID = parentID.String
		}

		subnet.CreatedAt = time.Unix(createdAt, 0)
		subnet.UpdatedAt = time.Unix(updatedAt, 0)

		subnets = append(subnets, &subnet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return &SubnetList{
		Subnets:    subnets,
		TotalCount: totalCount,
	}, nil
}

// GetSubnetChildren retrieves child subnets for a given parent subnet ID
func (r *SQLiteRepository) GetSubnetChildren(ctx context.Context, parentID string) ([]*Subnet, error) {
	query := `
		SELECT 
			id, cidr, name, description, location, location_type,
			cloud_provider, cloud_region, cloud_account_id, cloud_resource_type, cloud_vpc_id, cloud_subnet_id,
			parent_id, utilization_percent, created_at, updated_at
		FROM subnets
		WHERE parent_id = ?
		ORDER BY cidr
	`

	rows, err := r.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query child subnets: %w", err)
	}
	defer rows.Close()

	var subnets []*Subnet

	for rows.Next() {
		var subnet Subnet
		var description sql.NullString
		var cloudProvider, cloudRegion, cloudAccountID, cloudResourceType, cloudVPCId, cloudSubnetId sql.NullString
		var parentID sql.NullString
		var utilizationPercent sql.NullFloat64
		var createdAt, updatedAt int64

		err := rows.Scan(
			&subnet.ID, &subnet.CIDR, &subnet.Name, &description,
			&subnet.Location, &subnet.LocationType,
			&cloudProvider, &cloudRegion, &cloudAccountID, &cloudResourceType, &cloudVPCId, &cloudSubnetId,
			&parentID, &utilizationPercent, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan child subnet: %w", err)
		}

		// Parse cloud info
		if cloudProvider.Valid {
			subnet.CloudInfo = &CloudInfo{
				Provider:     cloudProvider.String,
				Region:       cloudRegion.String,
				AccountID:    cloudAccountID.String,
				ResourceType: cloudResourceType.String,
				VPCId:        cloudVPCId.String,
				SubnetId:     cloudSubnetId.String,
			}
		}

		// Parse utilization
		if utilizationPercent.Valid {
			subnet.Utilization = &Utilization{
				UtilizationPercent: utilizationPercent.Float64,
				LastUpdated:        time.Unix(updatedAt, 0),
			}
		}

		if parentID.Valid {
			subnet.ParentID = parentID.String
		}

		subnet.CreatedAt = time.Unix(createdAt, 0)
		subnet.UpdatedAt = time.Unix(updatedAt, 0)

		subnets = append(subnets, &subnet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating child subnet rows: %w", err)
	}

	return subnets, nil
}

// GetSubnetByID retrieves a subnet by its ID using repository models
func (r *SQLiteRepository) GetSubnetByID(ctx context.Context, id string) (*Subnet, error) {
	query := `
		SELECT 
			id, cidr, name, description, location, location_type,
			cloud_provider, cloud_region, cloud_account_id, cloud_resource_type, cloud_vpc_id, cloud_subnet_id,
			parent_id, address, netmask, wildcard, network, type, broadcast,
			host_min, host_max, hosts_per_net, is_public,
			total_ips, allocated_ips, utilization_percent, created_at, updated_at
		FROM subnets
		WHERE id = ?
	`

	var subnet Subnet
	var description sql.NullString
	var cloudProvider, cloudRegion, cloudAccountID, cloudResourceType, cloudVPCId, cloudSubnetId sql.NullString
	var parentID sql.NullString
	var address, netmask, wildcard, network, subnetType, broadcast sql.NullString
	var hostMin, hostMax sql.NullString
	var hostsPerNet sql.NullInt32
	var isPublic sql.NullInt32
	var totalIPs, allocatedIPs sql.NullInt32
	var utilizationPercent sql.NullFloat64
	var createdAt, updatedAt int64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subnet.ID, &subnet.CIDR, &subnet.Name, &description,
		&subnet.Location, &subnet.LocationType,
		&cloudProvider, &cloudRegion, &cloudAccountID, &cloudResourceType, &cloudVPCId, &cloudSubnetId,
		&parentID, &address, &netmask, &wildcard, &network, &subnetType, &broadcast,
		&hostMin, &hostMax, &hostsPerNet, &isPublic,
		&totalIPs, &allocatedIPs, &utilizationPercent, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subnet not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find subnet: %w", err)
	}

	// Parse cloud info
	if cloudProvider.Valid {
		subnet.CloudInfo = &CloudInfo{
			Provider:     cloudProvider.String,
			Region:       cloudRegion.String,
			AccountID:    cloudAccountID.String,
			ResourceType: cloudResourceType.String,
			VPCId:        cloudVPCId.String,
			SubnetId:     cloudSubnetId.String,
		}
	}

	// Parse subnet details
	if address.Valid {
		subnet.Details = &SubnetDetails{
			Address:     address.String,
			Netmask:     netmask.String,
			Wildcard:    wildcard.String,
			Network:     network.String,
			Type:        subnetType.String,
			Broadcast:   broadcast.String,
			HostMin:     hostMin.String,
			HostMax:     hostMax.String,
			HostsPerNet: hostsPerNet.Int32,
			IsPublic:    isPublic.Int32 == 1,
		}
	}

	// Parse utilization
	if utilizationPercent.Valid {
		subnet.Utilization = &Utilization{
			TotalIPs:           totalIPs.Int32,
			AllocatedIPs:       allocatedIPs.Int32,
			UtilizationPercent: utilizationPercent.Float64,
			LastUpdated:        time.Unix(updatedAt, 0),
		}
	}

	if parentID.Valid {
		subnet.ParentID = parentID.String
	}

	subnet.CreatedAt = time.Unix(createdAt, 0)
	subnet.UpdatedAt = time.Unix(updatedAt, 0)

	return &subnet, nil
}

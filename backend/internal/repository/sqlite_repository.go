package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		updated_at INTEGER
	);

	CREATE INDEX IF NOT EXISTS idx_subnets_location ON subnets(location);
	CREATE INDEX IF NOT EXISTS idx_subnets_cloud_provider ON subnets(cloud_provider);
	CREATE INDEX IF NOT EXISTS idx_subnets_cidr ON subnets(cidr);
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

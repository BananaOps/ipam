// SubnetList component - displays subnets with filtering and search
import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSearch, faFilter, faNetworkWired } from '@fortawesome/free-solid-svg-icons';
import { apiClient } from '../services/api';
import { Subnet, SubnetFilters, LocationType, CloudProviderType, APIError } from '../types';
import CloudProviderIcon from './CloudProviderIcon';
import ErrorMessage from './ErrorMessage';
import { useToast } from '../contexts/ToastContext';
import './SubnetList.css';

interface SubnetListProps {
  filters?: SubnetFilters;
  onFilterChange?: (filters: SubnetFilters) => void;
}

/**
 * Generate a consistent color for a given text string
 * Same text will always produce the same color
 * Uses subtle colors similar to location-type badges
 * Optimized for both light and dark modes
 */
function getColorForText(text: string): { backgroundColor: string; color: string } {
  // Simple hash function to convert string to number
  let hash = 0;
  for (let i = 0; i < text.length; i++) {
    hash = text.charCodeAt(i) + ((hash << 5) - hash);
  }
  
  // Convert hash to HSL color with subtle, pastel tones
  const hue = Math.abs(hash % 360);
  const saturation = 75; // Good saturation for visibility
  
  // Create a subtle background with higher opacity for better visibility
  const backgroundColor = `hsla(${hue}, ${saturation}%, 50%, 0.15)`;
  // Create a brighter, more saturated text color for better contrast in dark mode
  const textColor = `hsl(${hue}, ${saturation}%, 65%)`;
  
  return {
    backgroundColor,
    color: textColor,
  };
}

/**
 * Convert CIDR to a comparable number for sorting
 */
function cidrToNumber(cidr: string): number {
  const [ip, prefix] = cidr.split('/');
  const parts = ip.split('.').map(Number);
  const ipNumber = (parts[0] << 24) + (parts[1] << 16) + (parts[2] << 8) + parts[3];
  // Combine IP and prefix for sorting (IP is primary, prefix is secondary)
  return ipNumber * 1000 + parseInt(prefix);
}

/**
 * Sort subnets by CIDR
 */
function sortSubnetsByCIDR(subnets: Subnet[], order: SortOrder): Subnet[] {
  return [...subnets].sort((a, b) => {
    const aNum = cidrToNumber(a.cidr);
    const bNum = cidrToNumber(b.cidr);
    return order === 'asc' ? aNum - bNum : bNum - aNum;
  });
}

/**
 * Group subnets by the specified criteria
 */
function groupSubnets(subnets: Subnet[], groupBy: GroupBy): { [key: string]: Subnet[] } {
  if (groupBy === 'none') {
    return { 'All Subnets': subnets };
  }

  return subnets.reduce((groups, subnet) => {
    let key: string;
    
    switch (groupBy) {
      case 'provider':
        key = subnet.cloudInfo?.provider ? subnet.cloudInfo.provider.toUpperCase() : 'On-Premise';
        break;
      case 'location':
        key = subnet.location || 'Unknown Location';
        break;
      case 'type':
        key = subnet.locationType || 'Unknown Type';
        break;
      default:
        key = 'All Subnets';
    }

    if (!groups[key]) {
      groups[key] = [];
    }
    groups[key].push(subnet);
    return groups;
  }, {} as { [key: string]: Subnet[] });
}

type SortOrder = 'asc' | 'desc';
type GroupBy = 'none' | 'provider' | 'location' | 'type';

function SubnetList({ filters: externalFilters, onFilterChange }: SubnetListProps) {
  const navigate = useNavigate();
  const [subnets, setSubnets] = useState<Subnet[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<APIError | Error | null>(null);
  const [filters, setFilters] = useState<SubnetFilters>(externalFilters || {});
  const [searchQuery, setSearchQuery] = useState('');
  const [locationQuery, setLocationQuery] = useState('');
  const [sortOrder, setSortOrder] = useState<SortOrder>('asc');
  const [groupBy, setGroupBy] = useState<GroupBy>('none');
  const { showError: showToastError, showSuccess } = useToast();

  // Load subnets from API
  useEffect(() => {
    loadSubnets();
  }, [filters]);

  // Debounce search query
  useEffect(() => {
    const timeoutId = setTimeout(() => {
      if (searchQuery !== (filters.searchQuery || '')) {
        const newFilters = { ...filters, searchQuery: searchQuery || undefined };
        setFilters(newFilters);
        onFilterChange?.(newFilters);
      }
    }, 300); // Wait 300ms after user stops typing

    return () => clearTimeout(timeoutId);
  }, [searchQuery]);

  // Debounce location filter
  useEffect(() => {
    const timeoutId = setTimeout(() => {
      if (locationQuery !== (filters.location || '')) {
        const newFilters = { ...filters, location: locationQuery || undefined };
        setFilters(newFilters);
        onFilterChange?.(newFilters);
      }
    }, 300); // Wait 300ms after user stops typing

    return () => clearTimeout(timeoutId);
  }, [locationQuery]);

  const loadSubnets = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await apiClient.listSubnets(filters);
      setSubnets(response.subnets);
    } catch (err: any) {
      setError(err);
      showToastError(err.message || 'Failed to load subnets');
    } finally {
      setLoading(false);
    }
  };

  const handleLocationFilterChange = (location: string) => {
    setLocationQuery(location);
  };

  const handleCloudProviderFilterChange = (provider: string) => {
    const newFilters = {
      ...filters,
      cloudProvider: provider ? (provider as CloudProviderType) : undefined,
    };
    setFilters(newFilters);
    onFilterChange?.(newFilters);
  };

  const handleSearchChange = (query: string) => {
    setSearchQuery(query);
  };

  const clearFilters = () => {
    setFilters({});
    setSearchQuery('');
    setLocationQuery('');
    onFilterChange?.({});
  };

  if (loading) {
    return (
      <div className="subnet-list-loading">
        <div className="loading-spinner"></div>
        <p>Loading subnets...</p>
      </div>
    );
  }

  if (error) {
    return (
      <ErrorMessage
        error={error}
        onRetry={loadSubnets}
        onDismiss={() => setError(null)}
        showDetails={true}
      />
    );
  }

  return (
    <div className="subnet-list">
      {/* Filters Section */}
      <div className="subnet-list-filters">
        <div className="filter-group">
          <label htmlFor="search-input">
            <FontAwesomeIcon icon={faSearch} />
            <span>Search</span>
          </label>
          <input
            id="search-input"
            type="text"
            placeholder="Search by name or CIDR..."
            value={searchQuery}
            onChange={(e) => handleSearchChange(e.target.value)}
            className="search-input"
          />
        </div>

        <div className="filter-group">
          <label htmlFor="location-filter">
            <FontAwesomeIcon icon={faFilter} />
            <span>Location</span>
          </label>
          <input
            id="location-filter"
            type="text"
            placeholder="Filter by location..."
            value={locationQuery}
            onChange={(e) => handleLocationFilterChange(e.target.value)}
            className="filter-input"
          />
        </div>

        <div className="filter-group filter-group-cloud">
          <label>
            <FontAwesomeIcon icon={faFilter} />
            <span>Cloud Provider</span>
          </label>
          <div className="cloud-provider-filters">
            <button
              className={`cloud-filter-btn ${!filters.cloudProvider ? 'active' : ''}`}
              onClick={() => handleCloudProviderFilterChange('')}
              title="All Providers"
            >
              All
            </button>
            <button
              className={`cloud-filter-btn ${filters.cloudProvider === CloudProviderType.AWS ? 'active' : ''}`}
              onClick={() => handleCloudProviderFilterChange(CloudProviderType.AWS)}
              title="AWS"
            >
              <CloudProviderIcon provider={CloudProviderType.AWS} size="lg" />
            </button>
            <button
              className={`cloud-filter-btn ${filters.cloudProvider === CloudProviderType.AZURE ? 'active' : ''}`}
              onClick={() => handleCloudProviderFilterChange(CloudProviderType.AZURE)}
              title="Azure"
            >
              <CloudProviderIcon provider={CloudProviderType.AZURE} size="lg" />
            </button>
            <button
              className={`cloud-filter-btn ${filters.cloudProvider === CloudProviderType.GCP ? 'active' : ''}`}
              onClick={() => handleCloudProviderFilterChange(CloudProviderType.GCP)}
              title="GCP"
            >
              <CloudProviderIcon provider={CloudProviderType.GCP} size="lg" />
            </button>
            <button
              className={`cloud-filter-btn ${filters.cloudProvider === CloudProviderType.SCALEWAY ? 'active' : ''}`}
              onClick={() => handleCloudProviderFilterChange(CloudProviderType.SCALEWAY)}
              title="Scaleway"
            >
              <CloudProviderIcon provider={CloudProviderType.SCALEWAY} size="lg" />
            </button>
            <button
              className={`cloud-filter-btn ${filters.cloudProvider === CloudProviderType.OVH ? 'active' : ''}`}
              onClick={() => handleCloudProviderFilterChange(CloudProviderType.OVH)}
              title="OVH"
            >
              <CloudProviderIcon provider={CloudProviderType.OVH} size="lg" />
            </button>
          </div>
        </div>

        <div className="filter-group">
          <label htmlFor="sort-order">
            <FontAwesomeIcon icon={faFilter} />
            <span>Sort CIDR</span>
          </label>
          <select
            id="sort-order"
            value={sortOrder}
            onChange={(e) => setSortOrder(e.target.value as SortOrder)}
            className="filter-select"
          >
            <option value="asc">Ascending</option>
            <option value="desc">Descending</option>
          </select>
        </div>

        <div className="filter-group">
          <label htmlFor="group-by">
            <FontAwesomeIcon icon={faFilter} />
            <span>Group By</span>
          </label>
          <select
            id="group-by"
            value={groupBy}
            onChange={(e) => setGroupBy(e.target.value as GroupBy)}
            className="filter-select"
          >
            <option value="none">No Grouping</option>
            <option value="provider">Cloud Provider</option>
            <option value="location">Location</option>
            <option value="type">Location Type</option>
          </select>
        </div>

        {(filters.location || filters.cloudProvider || filters.searchQuery) && (
          <button onClick={clearFilters} className="clear-filters-button">
            Clear Filters
          </button>
        )}
      </div>

      {/* Subnets Display */}
      {subnets.length === 0 ? (
        <div className="subnet-list-empty">
          <FontAwesomeIcon icon={faNetworkWired} className="empty-icon" />
          <h3>No subnets found</h3>
          <p>
            {filters.location || filters.cloudProvider || filters.searchQuery
              ? 'Try adjusting your filters or search query.'
              : 'Get started by creating your first subnet.'}
          </p>
          <Link to="/subnets/create" className="create-subnet-link">
            Create Subnet
          </Link>
        </div>
      ) : (
        (() => {
          // Sort and group subnets
          const sortedSubnets = sortSubnetsByCIDR(subnets, sortOrder);
          const groupedSubnets = groupSubnets(sortedSubnets, groupBy);
          
          return (
            <div className="subnets-display">
              {Object.entries(groupedSubnets).map(([groupName, groupSubnets]) => (
                <div key={groupName} className="subnet-group">
                  {groupBy !== 'none' && (
                    <div className="group-header">
                      <h3 className="group-title">{groupName}</h3>
                      <span className="group-count">({groupSubnets.length} subnet{groupSubnets.length !== 1 ? 's' : ''})</span>
                    </div>
                  )}
                  
                  <div className="subnet-table-container">
                    <table className="subnet-table">
                      {groupBy === 'none' && (
                        <thead>
                          <tr>
                            <th>CIDR</th>
                            <th>Name</th>
                            <th>Location</th>
                            <th>Type</th>
                            <th>Cloud Info</th>
                            <th>Utilization</th>
                          </tr>
                        </thead>
                      )}
                      <tbody>
                        {groupSubnets.map((subnet) => (
                <tr 
                  key={subnet.id} 
                  className="subnet-row clickable-row"
                  onClick={() => navigate(`/subnets/${subnet.id}`)}
                >
                  <td className="subnet-cidr">
                    {subnet.cidr}
                  </td>
                  <td className="subnet-name">{subnet.name}</td>
                  <td className="subnet-location">{subnet.location}</td>
                  <td className="subnet-type">
                    <span className={`location-type-badge ${subnet.locationType ? subnet.locationType.toLowerCase() : 'unknown'}`}>
                      {subnet.locationType || 'Unknown'}
                    </span>
                  </td>
                  <td className="subnet-cloud-info">
                    {subnet.cloudInfo && subnet.cloudInfo.provider ? (
                      <div className="cloud-info-container">
                        <div className="cloud-info-row">
                          <CloudProviderIcon
                            provider={subnet.cloudInfo.provider}
                            size="lg"
                            title={subnet.cloudInfo.provider.toUpperCase()}
                          />
                        </div>
                        <div className="cloud-info-labels">
                          <span 
                            className="cloud-label cloud-region-label"
                            style={getColorForText(subnet.cloudInfo.region)}
                          >
                            {subnet.cloudInfo.region}
                          </span>
                          <span 
                            className="cloud-label cloud-account-label"
                            style={getColorForText(subnet.cloudInfo.accountId)}
                          >
                            {subnet.cloudInfo.accountId}
                          </span>
                        </div>
                      </div>
                    ) : (
                      <span className="no-cloud-info">—</span>
                    )}
                  </td>
                  <td className="subnet-utilization">
                    {subnet.utilization && subnet.utilization.utilizationPercent !== undefined ? (
                      <div className="utilization-display">
                        <div className="utilization-bar">
                          <div
                            className={`utilization-fill ${
                              subnet.utilization.utilizationPercent >= 80 ? 'high' : ''
                            }`}
                            style={{ width: `${subnet.utilization.utilizationPercent}%` }}
                          ></div>
                        </div>
                        <span className="utilization-percent">
                          {subnet.utilization.utilizationPercent.toFixed(1)}%
                        </span>
                      </div>
                    ) : (
                      <span className="no-utilization">N/A</span>
                    )}
                  </td>
                        </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              ))}
            </div>
          );
        })()
      )}
    </div>
  );
}

export default SubnetList;

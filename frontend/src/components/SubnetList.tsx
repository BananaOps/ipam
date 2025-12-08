// SubnetList component - displays subnets with filtering and search
import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
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

function SubnetList({ filters: externalFilters, onFilterChange }: SubnetListProps) {
  const [subnets, setSubnets] = useState<Subnet[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<APIError | Error | null>(null);
  const [filters, setFilters] = useState<SubnetFilters>(externalFilters || {});
  const [searchQuery, setSearchQuery] = useState('');
  const { showError: showToastError, showSuccess } = useToast();

  // Load subnets from API
  useEffect(() => {
    loadSubnets();
  }, [filters]);

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
    const newFilters = { ...filters, location: location || undefined };
    setFilters(newFilters);
    onFilterChange?.(newFilters);
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
    const newFilters = { ...filters, searchQuery: query || undefined };
    setFilters(newFilters);
    onFilterChange?.(newFilters);
  };

  const clearFilters = () => {
    setFilters({});
    setSearchQuery('');
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
            value={filters.location || ''}
            onChange={(e) => handleLocationFilterChange(e.target.value)}
            className="filter-input"
          />
        </div>

        <div className="filter-group">
          <label htmlFor="cloud-provider-filter">
            <FontAwesomeIcon icon={faFilter} />
            <span>Cloud Provider</span>
          </label>
          <select
            id="cloud-provider-filter"
            value={filters.cloudProvider || ''}
            onChange={(e) => handleCloudProviderFilterChange(e.target.value)}
            className="filter-select"
          >
            <option value="">All Providers</option>
            <option value={CloudProviderType.AWS}>AWS</option>
            <option value={CloudProviderType.AZURE}>Azure</option>
            <option value={CloudProviderType.GCP}>Google Cloud</option>
            <option value={CloudProviderType.SCALEWAY}>Scaleway</option>
            <option value={CloudProviderType.OVH}>OVH</option>
          </select>
        </div>

        {(filters.location || filters.cloudProvider || filters.searchQuery) && (
          <button onClick={clearFilters} className="clear-filters-button">
            Clear Filters
          </button>
        )}
      </div>

      {/* Subnets Table */}
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
        <div className="subnet-table-container">
          <table className="subnet-table">
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
            <tbody>
              {subnets.map((subnet) => (
                <tr key={subnet.id} className="subnet-row">
                  <td className="subnet-cidr">
                    <Link to={`/subnets/${subnet.id}`} className="subnet-link">
                      {subnet.cidr}
                    </Link>
                  </td>
                  <td className="subnet-name">{subnet.name}</td>
                  <td className="subnet-location">{subnet.location}</td>
                  <td className="subnet-type">
                    <span className={`location-type-badge ${subnet.locationType.toLowerCase()}`}>
                      {subnet.locationType}
                    </span>
                  </td>
                  <td className="subnet-cloud-info">
                    {subnet.cloudInfo ? (
                      <div className="cloud-info">
                        <CloudProviderIcon
                          provider={subnet.cloudInfo.provider}
                          size="lg"
                        />
                        <span className="cloud-region">{subnet.cloudInfo.region}</span>
                        <span className="cloud-account">{subnet.cloudInfo.accountId}</span>
                      </div>
                    ) : (
                      <span className="no-cloud-info">—</span>
                    )}
                  </td>
                  <td className="subnet-utilization">
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
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

export default SubnetList;

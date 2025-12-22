// SubnetMappingPage - Visual diagram of subnet mapping and relationships
import { useState, useEffect } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { 
  faNetworkWired, 
  faProjectDiagram, 
  faFilter,
  faDownload,
  faExpand,
  faCompress
} from '@fortawesome/free-solid-svg-icons';
import { apiClient } from '../services/api';
import { 
  Subnet, 
  SubnetFilters, 
  CloudProviderType, 
  APIError,
  SubnetConnection 
} from '../types';
import CloudProviderIcon from '../components/CloudProviderIcon';
import ErrorMessage from '../components/ErrorMessage';
import { useToast } from '../contexts/ToastContext';
import SubnetDiagram from '../components/SubnetDiagram';
import './SubnetMappingPage.css';

function SubnetMappingPage() {
  const [subnets, setSubnets] = useState<Subnet[]>([]);
  const [connections, setConnections] = useState<SubnetConnection[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<APIError | Error | null>(null);
  const [filters, setFilters] = useState<SubnetFilters>({});
  const [viewMode, setViewMode] = useState<'hierarchy' | 'network' | 'cloud'>('hierarchy');
  const [isFullscreen, setIsFullscreen] = useState(false);
  const { showError: showToastError } = useToast();

  useEffect(() => {
    loadSubnets();
  }, [filters]);

  const loadSubnets = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const [subnetsResponse, connectionsResponse] = await Promise.all([
        apiClient.listSubnets(filters),
        apiClient.listConnections()
      ]);
      
      setSubnets(subnetsResponse.subnets);
      setConnections(connectionsResponse.connections);
    } catch (err: any) {
      setError(err);
      showToastError(err.message || 'Failed to load subnets');
    } finally {
      setLoading(false);
    }
  };

  const handleCloudProviderFilterChange = (provider: string) => {
    const newFilters = {
      ...filters,
      cloudProvider: provider ? (provider as CloudProviderType) : undefined,
    };
    setFilters(newFilters);
  };

  const handleLocationFilterChange = (location: string) => {
    const newFilters = {
      ...filters,
      location: location || undefined,
    };
    setFilters(newFilters);
  };

  const clearFilters = () => {
    setFilters({});
  };

  const exportDiagram = () => {
    // TODO: Implement diagram export functionality
    console.log('Export diagram functionality to be implemented');
  };

  const toggleFullscreen = () => {
    setIsFullscreen(!isFullscreen);
  };

  if (loading) {
    return (
      <div className="subnet-mapping-loading">
        <div className="loading-spinner"></div>
        <p>Loading subnet mapping...</p>
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
    <div className={`subnet-mapping-page ${isFullscreen ? 'fullscreen' : ''}`}>
      <div className="subnet-mapping-header">
        <div className="header-title">
          <FontAwesomeIcon icon={faProjectDiagram} />
          <h1>Subnet Network Mapping</h1>
          <span className="subnet-count">({subnets.length} subnets)</span>
        </div>
        
        <div className="header-actions">
          <button
            onClick={toggleFullscreen}
            className="action-button"
            title={isFullscreen ? 'Exit Fullscreen' : 'Enter Fullscreen'}
          >
            <FontAwesomeIcon icon={isFullscreen ? faCompress : faExpand} />
          </button>
          <button
            onClick={exportDiagram}
            className="action-button"
            title="Export Diagram"
          >
            <FontAwesomeIcon icon={faDownload} />
            Export
          </button>
        </div>
      </div>

      {/* Filters and View Controls */}
      <div className="subnet-mapping-controls">
        <div className="view-mode-selector">
          <label>View Mode:</label>
          <div className="view-mode-buttons">
            <button
              className={`view-mode-btn ${viewMode === 'hierarchy' ? 'active' : ''}`}
              onClick={() => setViewMode('hierarchy')}
            >
              Hierarchy
            </button>
            <button
              className={`view-mode-btn ${viewMode === 'network' ? 'active' : ''}`}
              onClick={() => setViewMode('network')}
            >
              Network
            </button>
            <button
              className={`view-mode-btn ${viewMode === 'cloud' ? 'active' : ''}`}
              onClick={() => setViewMode('cloud')}
            >
              Cloud
            </button>
          </div>
        </div>

        <div className="mapping-filters">
          <div className="filter-group">
            <label>
              <FontAwesomeIcon icon={faFilter} />
              Location
            </label>
            <input
              type="text"
              placeholder="Filter by location..."
              value={filters.location || ''}
              onChange={(e) => handleLocationFilterChange(e.target.value)}
              className="filter-input"
            />
          </div>

          <div className="filter-group filter-group-cloud">
            <label>
              <FontAwesomeIcon icon={faFilter} />
              Cloud Provider
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

          {(filters.location || filters.cloudProvider) && (
            <button onClick={clearFilters} className="clear-filters-button">
              Clear Filters
            </button>
          )}
        </div>
      </div>

      {/* Diagram Container */}
      <div className="subnet-diagram-container">
        {subnets.length === 0 ? (
          <div className="subnet-mapping-empty">
            <FontAwesomeIcon icon={faNetworkWired} className="empty-icon" />
            <h3>No subnets to display</h3>
            <p>
              {filters.location || filters.cloudProvider
                ? 'Try adjusting your filters to see more subnets.'
                : 'Create some subnets to see the network mapping.'}
            </p>
          </div>
        ) : (
          <SubnetDiagram
            subnets={subnets}
            connections={connections}
            viewMode={viewMode}
            isFullscreen={isFullscreen}
          />
        )}
      </div>
    </div>
  );
}

export default SubnetMappingPage;

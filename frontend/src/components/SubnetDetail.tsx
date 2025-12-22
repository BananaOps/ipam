import { Subnet } from '../types';
import { HIGH_UTILIZATION_THRESHOLD } from '../config/constants';
import CloudProviderIcon from './CloudProviderIcon';
import SubnetChildren from './SubnetChildren';
import './SubnetDetail.css';

interface SubnetDetailProps {
  subnet: Subnet;
}

/**
 * SubnetDetail component displays comprehensive information about a subnet
 * including all calculated properties, utilization, and visual indicators
 */
function SubnetDetail({ subnet }: SubnetDetailProps) {
  const { details, utilization } = subnet;
  
  // Handle missing data gracefully
  if (!details || !utilization) {
    return (
      <div className="subnet-detail">
        <div className="subnet-detail-header">
          <h3 className="subnet-name">{subnet.name}</h3>
          {subnet.description && (
            <p className="subnet-description">{subnet.description}</p>
          )}
        </div>
        <div className="error-message">
          <p>Subnet details are not available. Please try refreshing the page.</p>
        </div>
      </div>
    );
  }
  
  const isHighUtilization = utilization.utilizationPercent !== undefined && 
    utilization.utilizationPercent >= HIGH_UTILIZATION_THRESHOLD;

  return (
    <div className="subnet-detail">
      {/* Header Section */}
      <div className="subnet-detail-header">
        <h3 className="subnet-name">{subnet.name}</h3>
        {subnet.description && (
          <p className="subnet-description">{subnet.description}</p>
        )}
      </div>

      {/* CIDR and Location */}
      <div className="subnet-info-section">
        <div className="info-row">
          <span className="info-label">CIDR:</span>
          <span className="info-value cidr-value">{subnet.cidr}</span>
        </div>
        <div className="info-row">
          <span className="info-label">Location:</span>
          <span className="info-value">{subnet.location} ({subnet.locationType})</span>
        </div>
      </div>

      {/* Cloud Provider Section - Only shown for CLOUD location type */}
      {subnet.locationType === 'CLOUD' && (
        <div className="cloud-provider-section">
          <h4 className="section-title">Cloud Provider Information</h4>
          {subnet.cloudInfo && subnet.cloudInfo.provider ? (
            <div className="cloud-info-grid">
              <div className="cloud-info-item">
                <span className="info-label">Provider:</span>
                <span className="info-value cloud-provider-display">
                  <CloudProviderIcon 
                    provider={subnet.cloudInfo.provider} 
                    size="lg"
                  />
                  <span className="provider-name">{subnet.cloudInfo.provider.toUpperCase()}</span>
                </span>
              </div>
              <div className="cloud-info-item">
                <span className="info-label">Region:</span>
                <span className="info-value">{subnet.cloudInfo.region}</span>
              </div>
              <div className="cloud-info-item">
                <span className="info-label">Account ID:</span>
                <span className="info-value">{subnet.cloudInfo.accountId}</span>
              </div>
              {subnet.cloudInfo.resourceType && (
                <div className="cloud-info-item">
                  <span className="info-label">Resource Type:</span>
                  <span className={`info-value resource-type ${subnet.cloudInfo.resourceType}`}>
                    {subnet.cloudInfo.resourceType.toUpperCase()}
                  </span>
                </div>
              )}
            </div>
          ) : (
            <div className="cloud-info-empty">
              <span style={{ fontStyle: 'italic', color: 'var(--text-secondary)' }}>
                No cloud provider information available
              </span>
            </div>
          )}
        </div>
      )}

      {/* Subnet Properties Section */}
      <div className="subnet-properties-section">
        <h4 className="section-title">Subnet Properties</h4>
        <div className="properties-grid">
          <div className="property-item">
            <span className="property-label">Address:</span>
            <span className="property-value">{details.address}</span>
          </div>
          <div className="property-item">
            <span className="property-label">Netmask:</span>
            <span className="property-value">{details.netmask}</span>
          </div>
          <div className="property-item">
            <span className="property-label">Wildcard:</span>
            <span className="property-value">{details.wildcard}</span>
          </div>
          <div className="property-item">
            <span className="property-label">Network:</span>
            <span className="property-value">{details.network}</span>
          </div>
          <div className="property-item">
            <span className="property-label">Type:</span>
            <span className="property-value">{details.type}</span>
          </div>
          <div className="property-item">
            <span className="property-label">Broadcast:</span>
            <span className="property-value">{details.broadcast}</span>
          </div>
          <div className="property-item">
            <span className="property-label">Host Min:</span>
            <span className="property-value">{details.hostMin}</span>
          </div>
          <div className="property-item">
            <span className="property-label">Host Max:</span>
            <span className="property-value">{details.hostMax}</span>
          </div>
          <div className="property-item">
            <span className="property-label">Hosts/Net:</span>
            <span className="property-value">
              {details.hostsPerNet !== undefined ? details.hostsPerNet.toLocaleString() : 'N/A'}
            </span>
          </div>
          <div className="property-item">
            <span className="property-label">Classification:</span>
            <span className={`property-value classification ${details.isPublic ? 'public' : 'private'}`}>
              {details.isPublic ? 'Public' : 'Private'}
            </span>
          </div>
        </div>
      </div>

      {/* Utilization Section */}
      <div className="subnet-utilization-section">
        <h4 className="section-title">IP Address Utilization</h4>
        <div className="utilization-info">
          <div className="utilization-stats">
            <div className="stat-item">
              <span className="stat-label">Total IPs:</span>
              <span className="stat-value">
                {utilization.totalIps !== undefined ? utilization.totalIps.toLocaleString() : 'N/A'}
              </span>
            </div>
            <div className="stat-item">
              <span className="stat-label">Allocated:</span>
              <span className="stat-value">
                {utilization.allocatedIps !== undefined ? utilization.allocatedIps.toLocaleString() : 'N/A'}
              </span>
            </div>
            <div className="stat-item">
              <span className="stat-label">Available:</span>
              <span className="stat-value">
                {(utilization.totalIps !== undefined && utilization.allocatedIps !== undefined) 
                  ? (utilization.totalIps - utilization.allocatedIps).toLocaleString() 
                  : 'N/A'}
              </span>
            </div>
          </div>
          
          {/* Utilization Percentage Display */}
          {utilization.utilizationPercent !== undefined ? (
            <>
              <div className={`utilization-percentage ${isHighUtilization ? 'high-utilization' : ''}`}>
                <span className="percentage-value">
                  {utilization.utilizationPercent.toFixed(1)}%
                </span>
                {isHighUtilization && (
                  <span className="high-utilization-indicator" title="High utilization warning">
                    ⚠
                  </span>
                )}
              </div>

              {/* Progress Bar */}
              <div className="utilization-progress-container">
                <div 
                  className={`utilization-progress-bar ${isHighUtilization ? 'high' : ''}`}
                  style={{ width: `${Math.min(utilization.utilizationPercent, 100)}%` }}
                  role="progressbar"
                  aria-valuenow={utilization.utilizationPercent}
                  aria-valuemin={0}
                  aria-valuemax={100}
                  aria-label={`IP utilization: ${utilization.utilizationPercent.toFixed(1)}%`}
                />
              </div>
            </>
          ) : (
            <div className="utilization-percentage">
              <span className="percentage-value">N/A</span>
            </div>
          )}
          
          {isHighUtilization && (
            <div className="high-utilization-warning">
              <span className="warning-icon">⚠</span>
              <span className="warning-text">
                High utilization detected. Consider expanding capacity or reviewing allocations.
              </span>
            </div>
          )}
        </div>
      </div>

      {/* Child Subnets Section */}
      <SubnetChildren parentSubnet={subnet} expanded={true} />
    </div>
  );
}

export default SubnetDetail;

import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSync, faCloud, faCheck, faExclamationTriangle } from '@fortawesome/free-solid-svg-icons';
import { apiClient } from '../services/api';
import { useToast } from '../contexts/ToastContext';
import './CloudSync.css';

interface CloudSyncProps {
  onSyncComplete?: () => void;
}

interface CloudStatus {
  enabled: boolean;
  providers: {
    [key: string]: {
      enabled: boolean;
      regions: string[];
    };
  };
}

const CloudSync: React.FC<CloudSyncProps> = ({ onSyncComplete }) => {
  const [syncing, setSyncing] = useState(false);
  const [status, setStatus] = useState<CloudStatus | null>(null);
  const [lastSync, setLastSync] = useState<Date | null>(null);
  const { showSuccess, showError } = useToast();

  // Load cloud status on component mount
  React.useEffect(() => {
    loadCloudStatus();
  }, []);

  const loadCloudStatus = async () => {
    try {
      const response = await fetch('/api/v1/cloud/status');
      if (response.ok) {
        const data = await response.json();
        setStatus(data);
      }
    } catch (error) {
      console.error('Failed to load cloud status:', error);
    }
  };

  const handleSync = async (provider?: string, region?: string) => {
    setSyncing(true);
    try {
      const payload: any = {};
      if (provider) payload.provider = provider;
      if (region) payload.region = region;

      const response = await fetch('/api/v1/cloud/sync', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (response.ok) {
        const data = await response.json();
        showSuccess(data.message || 'Cloud synchronization completed successfully');
        setLastSync(new Date());
        onSyncComplete?.();
      } else {
        const error = await response.json();
        showError(error.message || 'Cloud synchronization failed');
      }
    } catch (error) {
      showError('Failed to sync cloud providers');
      console.error('Sync error:', error);
    } finally {
      setSyncing(false);
    }
  };

  const handleUpdateUtilization = async () => {
    setSyncing(true);
    try {
      const response = await fetch('/api/v1/cloud/utilization/update', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (response.ok) {
        const data = await response.json();
        showSuccess(data.message || 'Utilization data updated successfully');
        onSyncComplete?.();
      } else {
        const error = await response.json();
        showError(error.message || 'Failed to update utilization data');
      }
    } catch (error) {
      showError('Failed to update utilization data');
      console.error('Utilization update error:', error);
    } finally {
      setSyncing(false);
    }
  };

  if (!status) {
    return (
      <div className="cloud-sync-loading">
        <FontAwesomeIcon icon={faSync} spin />
        <span>Loading cloud status...</span>
      </div>
    );
  }

  if (!status.enabled) {
    return (
      <div className="cloud-sync-disabled">
        <FontAwesomeIcon icon={faExclamationTriangle} />
        <span>Cloud providers are disabled</span>
      </div>
    );
  }

  return (
    <div className="cloud-sync">
      <div className="cloud-sync-header">
        <h3>
          <FontAwesomeIcon icon={faCloud} />
          Cloud Provider Synchronization
        </h3>
        {lastSync && (
          <span className="last-sync">
            Last sync: {lastSync.toLocaleString()}
          </span>
        )}
      </div>

      <div className="cloud-sync-actions">
        <button
          className="sync-button sync-all"
          onClick={() => handleSync()}
          disabled={syncing}
        >
          <FontAwesomeIcon icon={syncing ? faSync : faSync} spin={syncing} />
          {syncing ? 'Syncing...' : 'Sync All Providers'}
        </button>

        <button
          className="sync-button update-utilization"
          onClick={handleUpdateUtilization}
          disabled={syncing}
        >
          <FontAwesomeIcon icon={syncing ? faSync : faCheck} spin={syncing} />
          {syncing ? 'Updating...' : 'Update Utilization'}
        </button>
      </div>

      <div className="cloud-providers">
        {Object.entries(status.providers).map(([provider, info]) => (
          <div key={provider} className={`provider-card ${info.enabled ? 'enabled' : 'disabled'}`}>
            <div className="provider-header">
              <h4>{provider.toUpperCase()}</h4>
              <span className={`status ${info.enabled ? 'enabled' : 'disabled'}`}>
                {info.enabled ? 'Enabled' : 'Disabled'}
              </span>
            </div>
            
            {info.enabled && (
              <>
                <div className="provider-regions">
                  <strong>Regions:</strong>
                  <div className="regions-list">
                    {info.regions.map((region) => (
                      <span key={region} className="region-tag">
                        {region}
                      </span>
                    ))}
                  </div>
                </div>

                <div className="provider-actions">
                  <button
                    className="sync-button sync-provider"
                    onClick={() => handleSync(provider)}
                    disabled={syncing}
                  >
                    <FontAwesomeIcon icon={syncing ? faSync : faSync} spin={syncing} />
                    Sync {provider.toUpperCase()}
                  </button>

                  {info.regions.map((region) => (
                    <button
                      key={region}
                      className="sync-button sync-region"
                      onClick={() => handleSync(provider, region)}
                      disabled={syncing}
                    >
                      <FontAwesomeIcon icon={syncing ? faSync : faSync} spin={syncing} />
                      Sync {region}
                    </button>
                  ))}
                </div>
              </>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

export default CloudSync;

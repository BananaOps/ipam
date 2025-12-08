import { useState, useEffect, FormEvent } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faNetworkWired, faArrowLeft, faCheckCircle, faExclamationCircle, faSpinner } from '@fortawesome/free-solid-svg-icons';
import { apiClient } from '../services/api';
import { LocationType, CloudProviderType, UpdateSubnetRequest, Subnet } from '../types';
import './CreateSubnetPage.css';

function EditSubnetPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  
  // Form state
  const [formData, setFormData] = useState<UpdateSubnetRequest>({
    cidr: '',
    name: '',
    description: '',
    location: '',
    locationType: LocationType.DATACENTER,
  });
  
  // UI state
  const [loading, setLoading] = useState(false);
  const [loadingSubnet, setLoadingSubnet] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});
  const [originalSubnet, setOriginalSubnet] = useState<Subnet | null>(null);

  // Load existing subnet data
  useEffect(() => {
    const loadSubnet = async () => {
      if (!id) {
        setError('Subnet ID is missing');
        setLoadingSubnet(false);
        return;
      }

      try {
        setLoadingSubnet(true);
        const subnet = await apiClient.getSubnet(id);
        setOriginalSubnet(subnet);
        
        // Pre-populate form with existing data
        setFormData({
          cidr: subnet.cidr,
          name: subnet.name,
          description: subnet.description || '',
          location: subnet.location,
          locationType: subnet.locationType,
          cloudInfo: subnet.cloudInfo,
        });
        
        setLoadingSubnet(false);
      } catch (err: any) {
        setError(err.message || 'Failed to load subnet. Please try again.');
        setLoadingSubnet(false);
      }
    };

    loadSubnet();
  }, [id]);

  // Handle input changes
  const handleInputChange = (field: keyof UpdateSubnetRequest, value: any) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    // Clear validation error for this field
    if (validationErrors[field]) {
      setValidationErrors(prev => {
        const newErrors = { ...prev };
        delete newErrors[field];
        return newErrors;
      });
    }
  };

  // Handle location type change
  const handleLocationTypeChange = (locationType: LocationType) => {
    setFormData(prev => {
      const newData = { ...prev, locationType };
      // Clear cloud info if not cloud type
      if (locationType !== LocationType.CLOUD) {
        delete newData.cloudInfo;
      }
      return newData;
    });
  };

  // Handle cloud info changes
  const handleCloudInfoChange = (field: 'provider' | 'region' | 'accountId', value: string) => {
    setFormData(prev => ({
      ...prev,
      cloudInfo: {
        provider: prev.cloudInfo?.provider || CloudProviderType.AWS,
        region: prev.cloudInfo?.region || '',
        accountId: prev.cloudInfo?.accountId || '',
        [field]: value,
      },
    }));
  };

  // Validate form
  const validateForm = (): boolean => {
    const errors: Record<string, string> = {};

    // CIDR validation
    if (!formData.cidr?.trim()) {
      errors.cidr = 'CIDR is required';
    } else {
      // Basic CIDR format validation
      const cidrRegex = /^(\d{1,3}\.){3}\d{1,3}\/\d{1,2}$/;
      if (!cidrRegex.test(formData.cidr)) {
        errors.cidr = 'Invalid CIDR format (e.g., 192.168.1.0/24)';
      }
    }

    // Name validation
    if (!formData.name?.trim()) {
      errors.name = 'Name is required';
    }

    // Location validation
    if (!formData.location?.trim()) {
      errors.location = 'Location is required';
    }

    // Cloud info validation
    if (formData.locationType === LocationType.CLOUD) {
      if (!formData.cloudInfo?.provider) {
        errors.cloudProvider = 'Cloud provider is required';
      }
      if (!formData.cloudInfo?.region?.trim()) {
        errors.cloudRegion = 'Region is required for cloud subnets';
      }
      if (!formData.cloudInfo?.accountId?.trim()) {
        errors.cloudAccountId = 'Account ID is required for cloud subnets';
      }
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  // Check if CIDR has changed (for recalculation warning)
  const hasCidrChanged = (): boolean => {
    return originalSubnet !== null && formData.cidr !== originalSubnet.cidr;
  };

  // Handle form submission
  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    
    if (!id) {
      setError('Subnet ID is missing');
      return;
    }
    
    // Clear previous errors
    setError(null);
    setSuccess(false);

    // Validate form
    if (!validateForm()) {
      return;
    }

    try {
      setLoading(true);
      
      // Update subnet via API
      await apiClient.updateSubnet(id, formData);
      
      // Show success message
      setSuccess(true);
      
      // Navigate to the subnet detail page after a short delay
      setTimeout(() => {
        navigate(`/subnets/${id}`);
      }, 1500);
      
    } catch (err: any) {
      setError(err.message || 'Failed to update subnet. Please try again.');
      setLoading(false);
    }
  };

  // Handle cancel
  const handleCancel = () => {
    if (id) {
      navigate(`/subnets/${id}`);
    } else {
      navigate('/subnets');
    }
  };

  // Show loading state while fetching subnet
  if (loadingSubnet) {
    return (
      <div className="create-subnet-page">
        <div className="page-header">
          <h2>
            <FontAwesomeIcon icon={faSpinner} spin />
            Loading Subnet...
          </h2>
        </div>
      </div>
    );
  }

  // Show error if subnet couldn't be loaded
  if (!originalSubnet && error) {
    return (
      <div className="create-subnet-page">
        <div className="page-header">
          <button onClick={() => navigate('/subnets')} className="btn-back">
            <FontAwesomeIcon icon={faArrowLeft} />
            <span>Back to Subnets</span>
          </button>
          <h2>
            <FontAwesomeIcon icon={faExclamationCircle} />
            Error Loading Subnet
          </h2>
        </div>
        <div className="alert alert-error">
          <FontAwesomeIcon icon={faExclamationCircle} />
          <span>{error}</span>
        </div>
      </div>
    );
  }

  return (
    <div className="create-subnet-page">
      <div className="page-header">
        <button onClick={handleCancel} className="btn-back">
          <FontAwesomeIcon icon={faArrowLeft} />
          <span>Back to Subnet</span>
        </button>
        <h2>
          <FontAwesomeIcon icon={faNetworkWired} />
          Edit Subnet
        </h2>
        <p className="page-description">
          Update subnet information. If you change the CIDR, all subnet properties will be recalculated automatically.
        </p>
      </div>

      {/* Success Message */}
      {success && (
        <div className="alert alert-success">
          <FontAwesomeIcon icon={faCheckCircle} />
          <span>Subnet updated successfully! Redirecting...</span>
        </div>
      )}

      {/* Error Message */}
      {error && !loadingSubnet && (
        <div className="alert alert-error">
          <FontAwesomeIcon icon={faExclamationCircle} />
          <span>{error}</span>
        </div>
      )}

      {/* CIDR Change Warning */}
      {hasCidrChanged() && !success && (
        <div className="alert" style={{ 
          backgroundColor: 'rgba(251, 191, 36, 0.1)', 
          color: '#f59e0b', 
          borderColor: '#f59e0b' 
        }}>
          <FontAwesomeIcon icon={faExclamationCircle} />
          <span>
            <strong>Warning:</strong> Changing the CIDR will recalculate all subnet properties 
            (address, netmask, broadcast, host ranges, etc.).
          </span>
        </div>
      )}

      <form onSubmit={handleSubmit} className="subnet-form">
        {/* CIDR Input */}
        <div className="form-group">
          <label htmlFor="cidr" className="form-label required">
            CIDR Notation
          </label>
          <input
            id="cidr"
            type="text"
            className={`form-input ${validationErrors.cidr ? 'error' : ''}`}
            placeholder="e.g., 192.168.1.0/24"
            value={formData.cidr}
            onChange={(e) => handleInputChange('cidr', e.target.value)}
            disabled={loading || success}
          />
          {validationErrors.cidr && (
            <span className="form-error">{validationErrors.cidr}</span>
          )}
          <span className="form-hint">
            Enter the subnet in CIDR notation (IP address/prefix length)
          </span>
        </div>

        {/* Name Input */}
        <div className="form-group">
          <label htmlFor="name" className="form-label required">
            Subnet Name
          </label>
          <input
            id="name"
            type="text"
            className={`form-input ${validationErrors.name ? 'error' : ''}`}
            placeholder="e.g., Production Network"
            value={formData.name}
            onChange={(e) => handleInputChange('name', e.target.value)}
            disabled={loading || success}
          />
          {validationErrors.name && (
            <span className="form-error">{validationErrors.name}</span>
          )}
        </div>

        {/* Description Input */}
        <div className="form-group">
          <label htmlFor="description" className="form-label">
            Description
          </label>
          <textarea
            id="description"
            className="form-textarea"
            placeholder="Optional description of the subnet..."
            rows={3}
            value={formData.description}
            onChange={(e) => handleInputChange('description', e.target.value)}
            disabled={loading || success}
          />
        </div>

        {/* Location Input */}
        <div className="form-group">
          <label htmlFor="location" className="form-label required">
            Location
          </label>
          <input
            id="location"
            type="text"
            className={`form-input ${validationErrors.location ? 'error' : ''}`}
            placeholder="e.g., Paris DC1, New York Office"
            value={formData.location}
            onChange={(e) => handleInputChange('location', e.target.value)}
            disabled={loading || success}
          />
          {validationErrors.location && (
            <span className="form-error">{validationErrors.location}</span>
          )}
        </div>

        {/* Location Type Selection */}
        <div className="form-group">
          <label className="form-label required">Location Type</label>
          <div className="radio-group">
            <label className="radio-label">
              <input
                type="radio"
                name="locationType"
                value={LocationType.DATACENTER}
                checked={formData.locationType === LocationType.DATACENTER}
                onChange={() => handleLocationTypeChange(LocationType.DATACENTER)}
                disabled={loading || success}
              />
              <span>Datacenter</span>
            </label>
            <label className="radio-label">
              <input
                type="radio"
                name="locationType"
                value={LocationType.SITE}
                checked={formData.locationType === LocationType.SITE}
                onChange={() => handleLocationTypeChange(LocationType.SITE)}
                disabled={loading || success}
              />
              <span>Site</span>
            </label>
            <label className="radio-label">
              <input
                type="radio"
                name="locationType"
                value={LocationType.CLOUD}
                checked={formData.locationType === LocationType.CLOUD}
                onChange={() => handleLocationTypeChange(LocationType.CLOUD)}
                disabled={loading || success}
              />
              <span>Cloud</span>
            </label>
          </div>
        </div>

        {/* Cloud Provider Fields (shown only when location type is CLOUD) */}
        {formData.locationType === LocationType.CLOUD && (
          <div className="cloud-info-section">
            <h3 className="section-title">Cloud Provider Information</h3>
            
            {/* Cloud Provider Selection */}
            <div className="form-group">
              <label htmlFor="cloudProvider" className="form-label required">
                Cloud Provider
              </label>
              <select
                id="cloudProvider"
                className={`form-select ${validationErrors.cloudProvider ? 'error' : ''}`}
                value={formData.cloudInfo?.provider || ''}
                onChange={(e) => handleCloudInfoChange('provider', e.target.value)}
                disabled={loading || success}
              >
                <option value="">Select a provider...</option>
                <option value={CloudProviderType.AWS}>AWS</option>
                <option value={CloudProviderType.AZURE}>Azure</option>
                <option value={CloudProviderType.GCP}>Google Cloud</option>
                <option value={CloudProviderType.SCALEWAY}>Scaleway</option>
                <option value={CloudProviderType.OVH}>OVH</option>
              </select>
              {validationErrors.cloudProvider && (
                <span className="form-error">{validationErrors.cloudProvider}</span>
              )}
            </div>

            {/* Region Input */}
            <div className="form-group">
              <label htmlFor="cloudRegion" className="form-label required">
                Region
              </label>
              <input
                id="cloudRegion"
                type="text"
                className={`form-input ${validationErrors.cloudRegion ? 'error' : ''}`}
                placeholder="e.g., us-east-1, eu-west-1"
                value={formData.cloudInfo?.region || ''}
                onChange={(e) => handleCloudInfoChange('region', e.target.value)}
                disabled={loading || success}
              />
              {validationErrors.cloudRegion && (
                <span className="form-error">{validationErrors.cloudRegion}</span>
              )}
            </div>

            {/* Account ID Input */}
            <div className="form-group">
              <label htmlFor="cloudAccountId" className="form-label required">
                Account ID
              </label>
              <input
                id="cloudAccountId"
                type="text"
                className={`form-input ${validationErrors.cloudAccountId ? 'error' : ''}`}
                placeholder="e.g., 123456789012"
                value={formData.cloudInfo?.accountId || ''}
                onChange={(e) => handleCloudInfoChange('accountId', e.target.value)}
                disabled={loading || success}
              />
              {validationErrors.cloudAccountId && (
                <span className="form-error">{validationErrors.cloudAccountId}</span>
              )}
            </div>
          </div>
        )}

        {/* Form Actions */}
        <div className="form-actions">
          <button
            type="button"
            onClick={handleCancel}
            className="btn-secondary"
            disabled={loading || success}
          >
            Cancel
          </button>
          <button
            type="submit"
            className="btn-primary"
            disabled={loading || success}
          >
            {loading ? (
              <>
                <span className="loading-spinner-small"></span>
                <span>Updating...</span>
              </>
            ) : (
              <>
                <FontAwesomeIcon icon={faCheckCircle} />
                <span>Update Subnet</span>
              </>
            )}
          </button>
        </div>
      </form>
    </div>
  );
}

export default EditSubnetPage;

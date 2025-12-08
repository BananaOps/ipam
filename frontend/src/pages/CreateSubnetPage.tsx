import { useState, FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faNetworkWired, faArrowLeft, faCheckCircle, faExclamationCircle } from '@fortawesome/free-solid-svg-icons';
import { apiClient } from '../services/api';
import { LocationType, CloudProviderType, CreateSubnetRequest } from '../types';
import { useToast } from '../contexts/ToastContext';
import './CreateSubnetPage.css';

function CreateSubnetPage() {
  const navigate = useNavigate();
  const { showSuccess, showError } = useToast();
  
  // Form state
  const [formData, setFormData] = useState<CreateSubnetRequest>({
    cidr: '',
    name: '',
    description: '',
    location: '',
    locationType: LocationType.DATACENTER,
  });
  
  // UI state
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});

  // Handle input changes
  const handleInputChange = (field: keyof CreateSubnetRequest, value: any) => {
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
    if (!formData.cidr.trim()) {
      errors.cidr = 'CIDR is required';
    } else {
      // Basic CIDR format validation
      const cidrRegex = /^(\d{1,3}\.){3}\d{1,3}\/\d{1,2}$/;
      if (!cidrRegex.test(formData.cidr)) {
        errors.cidr = 'Invalid CIDR format (e.g., 192.168.1.0/24)';
      }
    }

    // Name validation
    if (!formData.name.trim()) {
      errors.name = 'Name is required';
    }

    // Location validation
    if (!formData.location.trim()) {
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

  // Handle form submission
  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    
    // Clear previous errors
    setError(null);
    setSuccess(false);

    // Validate form
    if (!validateForm()) {
      return;
    }

    try {
      setLoading(true);
      
      // Create subnet via API
      const createdSubnet = await apiClient.createSubnet(formData);
      
      // Show success message
      setSuccess(true);
      showSuccess(`Subnet ${createdSubnet.name} created successfully!`);
      
      // Navigate to the created subnet detail page after a short delay
      setTimeout(() => {
        navigate(`/subnets/${createdSubnet.id}`);
      }, 1500);
      
    } catch (err: any) {
      const errorMessage = err.message || 'Failed to create subnet. Please try again.';
      setError(errorMessage);
      showError(errorMessage);
      setLoading(false);
    }
  };

  // Handle cancel
  const handleCancel = () => {
    navigate('/subnets');
  };

  return (
    <div className="create-subnet-page">
      <div className="page-header">
        <button onClick={handleCancel} className="btn-back">
          <FontAwesomeIcon icon={faArrowLeft} />
          <span>Back to Subnets</span>
        </button>
        <h2>
          <FontAwesomeIcon icon={faNetworkWired} />
          Create New Subnet
        </h2>
        <p className="page-description">
          Add a new subnet to your IPAM inventory. All subnet properties will be calculated automatically.
        </p>
      </div>

      {/* Success Message */}
      {success && (
        <div className="alert alert-success">
          <FontAwesomeIcon icon={faCheckCircle} />
          <span>Subnet created successfully! Redirecting...</span>
        </div>
      )}

      {/* Error Message */}
      {error && (
        <div className="alert alert-error">
          <FontAwesomeIcon icon={faExclamationCircle} />
          <span>{error}</span>
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
                <span>Creating...</span>
              </>
            ) : (
              <>
                <FontAwesomeIcon icon={faNetworkWired} />
                <span>Create Subnet</span>
              </>
            )}
          </button>
        </div>
      </form>
    </div>
  );
}

export default CreateSubnetPage;

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { CloudProviderType } from '../types';
import { getCloudProviderIcon } from '../utils/cloudProviderIcons';
import './CloudProviderIcon.css';

interface CloudProviderIconProps {
  provider: CloudProviderType;
  className?: string;
  size?: 'xs' | 'sm' | 'lg' | '1x' | '2x' | '3x';
  title?: string;
}

/**
 * CloudProviderIcon component renders the appropriate FontAwesome icon
 * for a given cloud provider with consistent styling
 */
function CloudProviderIcon({ 
  provider, 
  className = '', 
  size = '1x',
  title 
}: CloudProviderIconProps) {
  const icon = getCloudProviderIcon(provider);
  const providerTitle = title || provider.toUpperCase();
  
  // Combine classes: base class, provider-specific class, and custom classes
  const combinedClassName = `cloud-provider-icon ${provider.toLowerCase()} ${className}`.trim();
  
  return (
    <FontAwesomeIcon
      icon={icon}
      className={combinedClassName}
      size={size}
      title={providerTitle}
      aria-label={`${providerTitle} cloud provider`}
    />
  );
}

export default CloudProviderIcon;

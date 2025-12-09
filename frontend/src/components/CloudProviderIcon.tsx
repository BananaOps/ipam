import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { CloudProviderType } from '../types';
import { getCloudProviderIcon } from '../utils/cloudProviderIcons';
import ScalewayIcon from './ScalewayIcon';
import AzureIcon from './AzureIcon';
import GCPIcon from './GCPIcon';
import './CloudProviderIcon.css';

interface CloudProviderIconProps {
  provider: CloudProviderType;
  className?: string;
  size?: 'xs' | 'sm' | 'lg' | '1x' | '2x' | '3x';
  title?: string;
}

/**
 * CloudProviderIcon component renders the appropriate icon
 * for a given cloud provider with consistent styling
 */
function CloudProviderIcon({ 
  provider, 
  className = '', 
  size = '1x',
  title 
}: CloudProviderIconProps) {
  const providerTitle = title || provider.toUpperCase();
  const combinedClassName = `cloud-provider-icon ${provider.toLowerCase()} ${className}`.trim();
  
  // Use custom logos for specific providers
  if (provider === CloudProviderType.SCALEWAY) {
    return (
      <ScalewayIcon
        size={size}
        className={combinedClassName}
      />
    );
  }
  
  if (provider === CloudProviderType.AZURE) {
    return (
      <AzureIcon
        size={size}
        className={combinedClassName}
      />
    );
  }
  
  if (provider === CloudProviderType.GCP) {
    return (
      <GCPIcon
        size={size}
        className={combinedClassName}
      />
    );
  }
  
  // Use FontAwesome icons for other providers (AWS, OVH)
  const icon = getCloudProviderIcon(provider);
  
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

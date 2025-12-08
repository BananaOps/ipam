// Cloud provider icon mapping
// Maps cloud provider types to their corresponding FontAwesome icons

import { IconDefinition } from '@fortawesome/fontawesome-svg-core';
import {
  faAws,
  faMicrosoft,
  faGoogle,
} from '@fortawesome/free-brands-svg-icons';
import { faCloud, faServer } from '@fortawesome/free-solid-svg-icons';
import { CloudProviderType } from '../types';

export const CLOUD_PROVIDER_ICONS: Record<CloudProviderType, IconDefinition> = {
  [CloudProviderType.AWS]: faAws,
  [CloudProviderType.AZURE]: faMicrosoft,
  [CloudProviderType.GCP]: faGoogle,
  [CloudProviderType.SCALEWAY]: faCloud,
  [CloudProviderType.OVH]: faServer,
};

export function getCloudProviderIcon(provider: CloudProviderType): IconDefinition {
  return CLOUD_PROVIDER_ICONS[provider] || faCloud;
}

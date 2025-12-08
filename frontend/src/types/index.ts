// Core data types for IPAM application

export enum LocationType {
  DATACENTER = 'DATACENTER',
  SITE = 'SITE',
  CLOUD = 'CLOUD',
}

export enum CloudProviderType {
  AWS = 'aws',
  AZURE = 'azure',
  GCP = 'gcp',
  SCALEWAY = 'scaleway',
  OVH = 'ovh',
}

export interface CloudInfo {
  provider: CloudProviderType;
  region: string;
  accountId: string;
}

export interface SubnetDetails {
  address: string;
  netmask: string;
  wildcard: string;
  network: string;
  type: string;
  broadcast: string;
  hostMin: string;
  hostMax: string;
  hostsPerNet: number;
  isPublic: boolean;
}

export interface UtilizationInfo {
  totalIps: number;
  allocatedIps: number;
  utilizationPercent: number;
}

export interface Subnet {
  id: string;
  cidr: string;
  name: string;
  description?: string;
  location: string;
  locationType: LocationType;
  cloudInfo?: CloudInfo;
  details: SubnetDetails;
  utilization: UtilizationInfo;
  createdAt: number;
  updatedAt: number;
}

export interface SubnetFilters {
  location?: string;
  cloudProvider?: CloudProviderType;
  searchQuery?: string;
}

export interface CreateSubnetRequest {
  cidr: string;
  name: string;
  description?: string;
  location: string;
  locationType: LocationType;
  cloudInfo?: CloudInfo;
}

export interface UpdateSubnetRequest {
  cidr?: string;
  name?: string;
  description?: string;
  location?: string;
  locationType?: LocationType;
  cloudInfo?: CloudInfo;
}

export interface SubnetListResponse {
  subnets: Subnet[];
  totalCount: number;
}

export interface APIError {
  code: string;
  message: string;
  details?: Record<string, string>;
  timestamp?: number;
}

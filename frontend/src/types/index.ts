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

export enum CloudResourceType {
  VPC = 'vpc',
  SUBNET = 'subnet',
}

export enum ConnectionType {
  VPN_SITE_TO_SITE = 'vpn_site_to_site',
  OPENVPN_CLIENT = 'openvpn_client',
  NAT_GATEWAY = 'nat_gateway',
  INTERNET_GATEWAY = 'internet_gateway',
  PEERING = 'peering',
  TRANSIT_GATEWAY = 'transit_gateway',
  DIRECT_CONNECT = 'direct_connect',
  EXPRESSROUTE = 'expressroute',
  CLOUD_INTERCONNECT = 'cloud_interconnect',
  LOAD_BALANCER = 'load_balancer',
  FIREWALL = 'firewall',
  CUSTOM = 'custom'
}

export enum ConnectionStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  PENDING = 'pending',
  ERROR = 'error'
}

export interface CloudInfo {
  provider: CloudProviderType;
  region: string;
  accountId: string;
  resourceType?: CloudResourceType;
  vpcId?: string;
  subnetId?: string;
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

export interface SubnetConnection {
  id: string;
  sourceSubnetId: string;
  targetSubnetId: string; // Peut être 'internet' pour une connexion vers Internet
  connectionType: ConnectionType;
  status: ConnectionStatus;
  name: string;
  description?: string;
  bandwidth?: string; // e.g., "1Gbps", "100Mbps"
  latency?: number; // in ms
  cost?: number; // monthly cost
  metadata?: Record<string, any>; // Additional connection-specific data
  createdAt: number;
  updatedAt: number;
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
  parentId?: string; // ID du réseau parent
  children?: Subnet[]; // Sous-réseaux enfants
  connections?: SubnetConnection[]; // Connexions vers d'autres sous-réseaux
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

export interface CreateConnectionRequest {
  sourceSubnetId: string;
  targetSubnetId: string;
  connectionType: ConnectionType;
  name: string;
  description?: string;
  bandwidth?: string;
  latency?: number;
  cost?: number;
  metadata?: Record<string, any>;
}

export interface UpdateConnectionRequest {
  name?: string;
  description?: string;
  connectionType?: ConnectionType;
  status?: ConnectionStatus;
  bandwidth?: string;
  latency?: number;
  cost?: number;
  metadata?: Record<string, any>;
}

export interface ConnectionListResponse {
  connections: SubnetConnection[];
  totalCount: number;
}

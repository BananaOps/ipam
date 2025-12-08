// Example usage of generated Protobuf types
// This file demonstrates how to use the generated TypeScript types

import { 
  Subnet, 
  LocationType, 
  CloudInfo, 
  SubnetDetails, 
  UtilizationInfo,
  CreateSubnetRequest,
  ListSubnetsRequest,
  GetSubnetRequest,
  UpdateSubnetRequest,
  DeleteSubnetRequest,
  Error as ProtoError
} from './subnet';

// Example: Creating a Subnet object
export function createExampleSubnet(): Subnet {
  const subnet: Subnet = {
    id: 'subnet-123',
    cidr: '10.0.0.0/24',
    name: 'Production Subnet',
    description: 'Main production subnet',
    location: 'datacenter-1',
    locationType: LocationType.DATACENTER,
    cloudInfo: {
      provider: 'aws',
      region: 'us-east-1',
      accountId: '123456789'
    },
    details: {
      address: '10.0.0.0',
      netmask: '255.255.255.0',
      wildcard: '0.0.0.255',
      network: '10.0.0.0',
      type: 'private',
      broadcast: '10.0.0.255',
      hostMin: '10.0.0.1',
      hostMax: '10.0.0.254',
      hostsPerNet: 254,
      isPublic: false
    },
    utilization: {
      totalIps: 254,
      allocatedIps: 100,
      utilizationPercent: 39.37
    },
    createdAt: BigInt(Date.now()),
    updatedAt: BigInt(Date.now())
  };

  return subnet;
}

// Example: Creating a CreateSubnetRequest
export function createSubnetRequest(cidr: string, name: string): CreateSubnetRequest {
  return {
    cidr,
    name,
    description: '',
    location: 'datacenter-1',
    locationType: LocationType.DATACENTER,
    cloudInfo: undefined
  };
}

// Example: Creating a ListSubnetsRequest with filters
export function createListRequest(
  locationFilter?: string,
  cloudProvider?: string
): ListSubnetsRequest {
  return {
    locationFilter: locationFilter || '',
    cloudProviderFilter: cloudProvider || '',
    searchQuery: '',
    page: 1,
    pageSize: 20
  };
}

// Example: Creating an Error object
export function createError(code: string, message: string): ProtoError {
  return {
    code,
    message,
    details: {},
    timestamp: BigInt(Date.now())
  };
}

// Type guards for checking message types
export function isSubnet(obj: any): obj is Subnet {
  return obj && typeof obj.id === 'string' && typeof obj.cidr === 'string';
}

export function isError(obj: any): obj is ProtoError {
  return obj && typeof obj.code === 'string' && typeof obj.message === 'string';
}

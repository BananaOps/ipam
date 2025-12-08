# IP Calculation Service

## Overview

The IP Calculation Service (`ipservice.go`) provides IP address management functionality using Go's standard `net/netip` library. This service implements all the required calculations for subnet management as specified in the IPAM by BananaOps design document.

## Implementation

### GoIPAMService

The `GoIPAMService` struct implements the `IPService` interface with the following methods:

#### ValidateCIDR(cidr string) error

Validates a CIDR notation string to ensure:
- The CIDR is not empty
- The CIDR format is valid (e.g., "192.168.1.0/24")
- The address is the network address (not a host address)

**Example:**
```go
service := NewGoIPAMService()
err := service.ValidateCIDR("192.168.1.0/24") // Valid
err := service.ValidateCIDR("192.168.1.5/24") // Invalid - host bits set
```

#### CalculateSubnetDetails(cidr string) (*pb.SubnetDetails, error)

Calculates all subnet properties from a CIDR notation:

**Calculated Properties:**
- **Address**: The network address
- **Netmask**: The subnet mask (e.g., "255.255.255.0")
- **Wildcard**: The wildcard mask (inverse of netmask)
- **Network**: The CIDR notation
- **Type**: "IPv4" or "IPv6"
- **Broadcast**: The broadcast address (IPv4 only)
- **HostMin**: The first usable host address
- **HostMax**: The last usable host address
- **HostsPerNet**: Total number of usable host addresses
- **IsPublic**: Whether the subnet is in public IP space

**Special Cases:**
- **/32 (IPv4)**: Single host, HostMin = HostMax = network address
- **/31 (IPv4)**: Point-to-point link (RFC 3021), 2 usable addresses
- **IPv6**: No broadcast address, all addresses in range are usable

**Example:**
```go
service := NewGoIPAMService()
details, err := service.CalculateSubnetDetails("192.168.1.0/24")
// details.Address = "192.168.1.0"
// details.Netmask = "255.255.255.0"
// details.Wildcard = "0.0.0.255"
// details.Broadcast = "192.168.1.255"
// details.HostMin = "192.168.1.1"
// details.HostMax = "192.168.1.254"
// details.HostsPerNet = 254
// details.IsPublic = false
```

#### CalculateUtilization(totalIPs, allocatedIPs int32) float32

Calculates the utilization percentage for a subnet:

**Formula:** `(allocatedIPs / totalIPs) * 100`

**Example:**
```go
service := NewGoIPAMService()
utilization := service.CalculateUtilization(254, 127) // Returns 50.0
```

### Public/Private Classification

The service automatically classifies IP addresses as public or private based on:

**Private IPv4 Ranges:**
- 10.0.0.0/8
- 172.16.0.0/12
- 192.168.0.0/16
- 127.0.0.0/8 (loopback)
- 169.254.0.0/16 (link-local)

**Private IPv6 Ranges:**
- fc00::/7 (Unique Local Addresses)
- fe80::/10 (Link-Local)
- ::1/128 (loopback)

All other addresses are considered public.

## Testing

The service includes comprehensive unit tests covering:

### ValidateCIDR Tests
- Valid IPv4 CIDR notations
- Valid IPv6 CIDR notations
- Invalid formats
- Host bits set in CIDR
- Empty CIDR strings

### CalculateSubnetDetails Tests
- IPv4 /24, /16, /32, /31 networks
- IPv6 /64 networks
- Public and private subnets
- All calculated properties
- Edge cases

### CalculateUtilization Tests
- Various utilization percentages
- Zero utilization
- Full utilization
- Edge cases

### Integration Tests
- Full ServiceLayer integration
- Mock repository usage
- End-to-end subnet creation flow

## Usage in ServiceLayer

The IPService is used by the ServiceLayer for:

1. **Subnet Creation**: Validates CIDR and calculates all properties
2. **Subnet Updates**: Recalculates properties when CIDR changes
3. **Validation**: Ensures all IP addresses are valid before persistence

**Example Integration:**
```go
ipService := NewGoIPAMService()
serviceLayer := NewServiceLayer(repo, ipService, cloudManager)

req := &pb.CreateSubnetRequest{
    Cidr: "192.168.1.0/24",
    Name: "Office Network",
    Location: "datacenter-1",
    LocationType: pb.LocationType_DATACENTER,
}

resp, err := serviceLayer.CreateSubnet(ctx, req)
// Subnet is created with all calculated properties
```

## Requirements Coverage

This implementation satisfies the following requirements from the design document:

- **Requirement 2.3**: Uses native Go IP calculation (net/netip)
- **Requirement 4.1**: Provides IP management functionality
- **Requirement 4.3**: Validates and computes IP addresses
- **Requirement 8.1**: Validates CIDR notation
- **Requirement 8.2**: Calculates all subnet properties automatically
- **Requirement 9.1**: Calculates utilization percentage

## Dependencies

- `net/netip`: Go standard library for IP address parsing and manipulation
- `go4.org/netipx`: Extended IP utilities for range operations
- `github.com/bananaops/ipam-bananaops/proto`: Protobuf definitions

## Performance

The service uses efficient algorithms:
- O(1) for CIDR validation
- O(1) for subnet calculations
- O(1) for utilization calculations
- No external API calls or database queries
- All calculations are performed in-memory

## Future Enhancements

Potential improvements:
- Support for subnet overlap detection
- IP allocation tracking within subnets
- Subnet hierarchy calculations
- Advanced IPv6 features
- Custom private IP range definitions

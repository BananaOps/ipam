package service

import (
	"fmt"
	"net"
	"net/netip"

	pb "github.com/bananaops/ipam-bananaops/proto"
	"go4.org/netipx"
)

// GoIPAMService implements IPService using go-ipam for IP calculations
type GoIPAMService struct{}

// NewGoIPAMService creates a new IPService instance
func NewGoIPAMService() *GoIPAMService {
	return &GoIPAMService{}
}

// ValidateCIDR validates a CIDR notation string
func (s *GoIPAMService) ValidateCIDR(cidr string) error {
	if cidr == "" {
		return fmt.Errorf("CIDR cannot be empty")
	}

	// Parse the CIDR using netip
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return fmt.Errorf("invalid CIDR notation: %w", err)
	}

	// Additional validation: ensure the address is the network address
	if prefix.Addr() != prefix.Masked().Addr() {
		return fmt.Errorf("CIDR address must be the network address (got %s, expected %s)",
			prefix.Addr(), prefix.Masked().Addr())
	}

	return nil
}

// CalculateSubnetDetails calculates all subnet properties from a CIDR
func (s *GoIPAMService) CalculateSubnetDetails(cidr string) (*pb.SubnetDetails, error) {
	// Parse the CIDR
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR notation: %w", err)
	}

	// Ensure we're working with the network address
	prefix = prefix.Masked()

	// Get the IP range
	ipRange := netipx.RangeOfPrefix(prefix)

	// Calculate network properties
	networkAddr := prefix.Addr()
	bits := prefix.Bits()

	// Calculate netmask
	netmask := netmaskFromPrefix(prefix)

	// Calculate wildcard mask (inverse of netmask)
	wildcard := wildcardFromNetmask(netmask)

	// Determine subnet type (IPv4 or IPv6)
	subnetType := "IPv4"
	if networkAddr.Is6() {
		subnetType = "IPv6"
	}

	// Calculate broadcast address (IPv4 only)
	var broadcast string
	if networkAddr.Is4() {
		broadcast = ipRange.To().String()
	} else {
		broadcast = "N/A (IPv6)"
	}

	// Calculate host range
	var hostMin, hostMax string
	var hostsPerNet int32

	if networkAddr.Is4() {
		// For IPv4
		if bits == 32 {
			// /32 - single host
			hostMin = networkAddr.String()
			hostMax = networkAddr.String()
			hostsPerNet = 1
		} else if bits == 31 {
			// /31 - point-to-point link (RFC 3021)
			hostMin = ipRange.From().String()
			hostMax = ipRange.To().String()
			hostsPerNet = 2
		} else {
			// Normal subnet - exclude network and broadcast
			hostMin = ipRange.From().Next().String()
			hostMax = ipRange.To().Prev().String()

			// Calculate total hosts (2^(32-bits) - 2)
			totalAddresses := uint64(1) << (32 - bits)
			hostsPerNet = int32(totalAddresses - 2)
		}
	} else {
		// For IPv6
		hostMin = ipRange.From().String()
		hostMax = ipRange.To().String()

		// For IPv6, calculate available addresses
		// Note: For large subnets, this might overflow, so we cap it
		if bits >= 64 {
			totalAddresses := uint64(1) << (128 - bits)
			if totalAddresses > uint64(2147483647) {
				hostsPerNet = 2147483647 // Max int32
			} else {
				hostsPerNet = int32(totalAddresses)
			}
		} else {
			// For very large IPv6 subnets, just use max int32
			hostsPerNet = 2147483647
		}
	}

	// Determine if the subnet is public or private
	isPublic := isPublicIP(networkAddr)

	return &pb.SubnetDetails{
		Address:     networkAddr.String(),
		Netmask:     netmask,
		Wildcard:    wildcard,
		Network:     prefix.String(),
		Type:        subnetType,
		Broadcast:   broadcast,
		HostMin:     hostMin,
		HostMax:     hostMax,
		HostsPerNet: hostsPerNet,
		IsPublic:    isPublic,
	}, nil
}

// CalculateUtilization calculates the utilization percentage for a subnet
func (s *GoIPAMService) CalculateUtilization(totalIPs, allocatedIPs int32) float32 {
	if totalIPs == 0 {
		return 0.0
	}
	return (float32(allocatedIPs) / float32(totalIPs)) * 100.0
}

// netmaskFromPrefix converts a prefix to a netmask string
func netmaskFromPrefix(prefix netip.Prefix) string {
	bits := prefix.Bits()
	addr := prefix.Addr()

	if addr.Is4() {
		// IPv4 netmask
		mask := net.CIDRMask(bits, 32)
		return net.IP(mask).String()
	}

	// IPv6 netmask
	mask := net.CIDRMask(bits, 128)
	return net.IP(mask).String()
}

// wildcardFromNetmask calculates the wildcard mask from a netmask
func wildcardFromNetmask(netmask string) string {
	ip := net.ParseIP(netmask)
	if ip == nil {
		return ""
	}

	// Convert to 4-byte or 16-byte representation
	if ip4 := ip.To4(); ip4 != nil {
		// IPv4
		wildcard := make(net.IP, 4)
		for i := 0; i < 4; i++ {
			wildcard[i] = ^ip4[i]
		}
		return wildcard.String()
	}

	// IPv6
	ip16 := ip.To16()
	wildcard := make(net.IP, 16)
	for i := 0; i < 16; i++ {
		wildcard[i] = ^ip16[i]
	}
	return wildcard.String()
}

// isPublicIP determines if an IP address is public or private
func isPublicIP(addr netip.Addr) bool {
	// Check for private IPv4 ranges
	if addr.Is4() {
		// 10.0.0.0/8
		if addr.As4()[0] == 10 {
			return false
		}
		// 172.16.0.0/12
		if addr.As4()[0] == 172 && addr.As4()[1] >= 16 && addr.As4()[1] <= 31 {
			return false
		}
		// 192.168.0.0/16
		if addr.As4()[0] == 192 && addr.As4()[1] == 168 {
			return false
		}
		// 127.0.0.0/8 (loopback)
		if addr.As4()[0] == 127 {
			return false
		}
		// 169.254.0.0/16 (link-local)
		if addr.As4()[0] == 169 && addr.As4()[1] == 254 {
			return false
		}
	}

	// Check for private IPv6 ranges
	if addr.Is6() {
		// fc00::/7 (Unique Local Addresses)
		if addr.As16()[0] == 0xfc || addr.As16()[0] == 0xfd {
			return false
		}
		// fe80::/10 (Link-Local)
		if addr.As16()[0] == 0xfe && (addr.As16()[1]&0xc0) == 0x80 {
			return false
		}
		// ::1/128 (loopback)
		if addr.IsLoopback() {
			return false
		}
	}

	// Check for loopback and other special addresses
	if addr.IsLoopback() || addr.IsLinkLocalUnicast() || addr.IsLinkLocalMulticast() {
		return false
	}

	return true
}

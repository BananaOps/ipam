package service

import (
	"net/netip"
	"testing"
)

func TestValidateCIDR(t *testing.T) {
	service := NewGoIPAMService()

	tests := []struct {
		name    string
		cidr    string
		wantErr bool
	}{
		{
			name:    "valid IPv4 CIDR",
			cidr:    "192.168.1.0/24",
			wantErr: false,
		},
		{
			name:    "valid IPv4 CIDR /32",
			cidr:    "10.0.0.1/32",
			wantErr: false,
		},
		{
			name:    "valid IPv6 CIDR",
			cidr:    "2001:db8::/32",
			wantErr: false,
		},
		{
			name:    "empty CIDR",
			cidr:    "",
			wantErr: true,
		},
		{
			name:    "invalid CIDR format",
			cidr:    "192.168.1.0",
			wantErr: true,
		},
		{
			name:    "invalid IP in CIDR",
			cidr:    "999.999.999.999/24",
			wantErr: true,
		},
		{
			name:    "CIDR with host bits set",
			cidr:    "192.168.1.5/24",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateCIDR(tt.cidr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCIDR() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCalculateSubnetDetails(t *testing.T) {
	service := NewGoIPAMService()

	tests := []struct {
		name            string
		cidr            string
		wantAddress     string
		wantNetmask     string
		wantType        string
		wantHostsPerNet int32
		wantIsPublic    bool
		wantErr         bool
	}{
		{
			name:            "IPv4 /24 network",
			cidr:            "192.168.1.0/24",
			wantAddress:     "192.168.1.0",
			wantNetmask:     "255.255.255.0",
			wantType:        "IPv4",
			wantHostsPerNet: 254,
			wantIsPublic:    false,
		},
		{
			name:            "IPv4 /32 single host",
			cidr:            "10.0.0.1/32",
			wantAddress:     "10.0.0.1",
			wantNetmask:     "255.255.255.255",
			wantType:        "IPv4",
			wantHostsPerNet: 1,
			wantIsPublic:    false,
		},
		{
			name:            "IPv4 /31 point-to-point",
			cidr:            "10.0.0.0/31",
			wantAddress:     "10.0.0.0",
			wantNetmask:     "255.255.255.254",
			wantType:        "IPv4",
			wantHostsPerNet: 2,
			wantIsPublic:    false,
		},
		{
			name:            "IPv4 /16 network",
			cidr:            "172.16.0.0/16",
			wantAddress:     "172.16.0.0",
			wantNetmask:     "255.255.0.0",
			wantType:        "IPv4",
			wantHostsPerNet: 65534,
			wantIsPublic:    false,
		},
		{
			name:            "public IPv4 /24",
			cidr:            "8.8.8.0/24",
			wantAddress:     "8.8.8.0",
			wantNetmask:     "255.255.255.0",
			wantType:        "IPv4",
			wantHostsPerNet: 254,
			wantIsPublic:    true,
		},
		{
			name:         "IPv6 /64 network",
			cidr:         "2001:db8::/64",
			wantAddress:  "2001:db8::",
			wantType:     "IPv6",
			wantIsPublic: true,
		},
		{
			name:    "invalid CIDR",
			cidr:    "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			details, err := service.CalculateSubnetDetails(tt.cidr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CalculateSubnetDetails() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("CalculateSubnetDetails() unexpected error = %v", err)
			}

			if details.Address != tt.wantAddress {
				t.Errorf("Address = %v, want %v", details.Address, tt.wantAddress)
			}

			if tt.wantNetmask != "" && details.Netmask != tt.wantNetmask {
				t.Errorf("Netmask = %v, want %v", details.Netmask, tt.wantNetmask)
			}

			if details.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", details.Type, tt.wantType)
			}

			if tt.wantHostsPerNet > 0 && details.HostsPerNet != tt.wantHostsPerNet {
				t.Errorf("HostsPerNet = %v, want %v", details.HostsPerNet, tt.wantHostsPerNet)
			}

			if details.IsPublic != tt.wantIsPublic {
				t.Errorf("IsPublic = %v, want %v", details.IsPublic, tt.wantIsPublic)
			}

			// Verify all required fields are populated
			if details.Wildcard == "" {
				t.Error("Wildcard should not be empty")
			}
			if details.Network == "" {
				t.Error("Network should not be empty")
			}
			if details.HostMin == "" {
				t.Error("HostMin should not be empty")
			}
			if details.HostMax == "" {
				t.Error("HostMax should not be empty")
			}
		})
	}
}

func TestCalculateUtilization(t *testing.T) {
	service := NewGoIPAMService()

	tests := []struct {
		name         string
		totalIPs     int32
		allocatedIPs int32
		want         float32
	}{
		{
			name:         "50% utilization",
			totalIPs:     100,
			allocatedIPs: 50,
			want:         50.0,
		},
		{
			name:         "0% utilization",
			totalIPs:     100,
			allocatedIPs: 0,
			want:         0.0,
		},
		{
			name:         "100% utilization",
			totalIPs:     100,
			allocatedIPs: 100,
			want:         100.0,
		},
		{
			name:         "zero total IPs",
			totalIPs:     0,
			allocatedIPs: 0,
			want:         0.0,
		},
		{
			name:         "partial utilization",
			totalIPs:     254,
			allocatedIPs: 127,
			want:         50.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.CalculateUtilization(tt.totalIPs, tt.allocatedIPs)
			if got != tt.want {
				t.Errorf("CalculateUtilization() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPublicIP(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		// Private IPv4 ranges
		{name: "10.0.0.0/8", ip: "10.0.0.1", want: false},
		{name: "172.16.0.0/12", ip: "172.16.0.1", want: false},
		{name: "172.31.255.255", ip: "172.31.255.255", want: false},
		{name: "192.168.0.0/16", ip: "192.168.1.1", want: false},
		{name: "127.0.0.0/8 loopback", ip: "127.0.0.1", want: false},
		{name: "169.254.0.0/16 link-local", ip: "169.254.1.1", want: false},

		// Public IPv4
		{name: "8.8.8.8 public", ip: "8.8.8.8", want: true},
		{name: "1.1.1.1 public", ip: "1.1.1.1", want: true},
		{name: "172.15.0.1 public", ip: "172.15.0.1", want: true},
		{name: "172.32.0.1 public", ip: "172.32.0.1", want: true},

		// Private IPv6 ranges
		{name: "fc00::/7 ULA", ip: "fc00::1", want: false},
		{name: "fd00::/7 ULA", ip: "fd00::1", want: false},
		{name: "fe80::/10 link-local", ip: "fe80::1", want: false},
		{name: "::1 loopback", ip: "::1", want: false},

		// Public IPv6
		{name: "2001:db8:: public", ip: "2001:db8::1", want: true},
		{name: "2606:4700:: public", ip: "2606:4700::1", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := parseIP(tt.ip)
			if err != nil {
				t.Fatalf("Failed to parse IP %s: %v", tt.ip, err)
			}

			got := isPublicIP(addr)
			if got != tt.want {
				t.Errorf("isPublicIP(%s) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

// Helper function to parse IP for testing
func parseIP(ip string) (addr netip.Addr, err error) {
	return netip.ParseAddr(ip)
}

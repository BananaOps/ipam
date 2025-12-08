import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import SubnetDetail from './SubnetDetail';
import type { Subnet } from '../types';
import { LocationType, CloudProviderType } from '../types';

describe('SubnetDetail', () => {
  const mockSubnet: Subnet = {
    id: '1',
    cidr: '192.168.1.0/24',
    name: 'Test Subnet',
    description: 'A test subnet for unit testing',
    location: 'datacenter-1',
    locationType: LocationType.DATACENTER,
    details: {
      address: '192.168.1.0',
      netmask: '255.255.255.0',
      wildcard: '0.0.0.255',
      network: '192.168.1.0',
      type: 'Private',
      broadcast: '192.168.1.255',
      hostMin: '192.168.1.1',
      hostMax: '192.168.1.254',
      hostsPerNet: 254,
      isPublic: false,
    },
    utilization: {
      totalIps: 254,
      allocatedIps: 100,
      utilizationPercent: 39.4,
    },
    createdAt: Date.now(),
    updatedAt: Date.now(),
  };

  it('renders subnet name and description', () => {
    render(<SubnetDetail subnet={mockSubnet} />);
    expect(screen.getByText('Test Subnet')).toBeInTheDocument();
    expect(screen.getByText('A test subnet for unit testing')).toBeInTheDocument();
  });

  it('displays CIDR notation', () => {
    render(<SubnetDetail subnet={mockSubnet} />);
    expect(screen.getByText('192.168.1.0/24')).toBeInTheDocument();
  });

  it('displays all subnet properties', () => {
    render(<SubnetDetail subnet={mockSubnet} />);
    
    // Check for key properties - use getAllByText for duplicates
    expect(screen.getAllByText('192.168.1.0').length).toBeGreaterThan(0); // address and network
    expect(screen.getByText('255.255.255.0')).toBeInTheDocument(); // netmask
    expect(screen.getByText('0.0.0.255')).toBeInTheDocument(); // wildcard
    expect(screen.getByText('192.168.1.255')).toBeInTheDocument(); // broadcast
    expect(screen.getByText('192.168.1.1')).toBeInTheDocument(); // hostMin
    expect(screen.getByText('192.168.1.254')).toBeInTheDocument(); // hostMax
    expect(screen.getAllByText('254').length).toBeGreaterThan(0); // hostsPerNet and totalIps
  });

  it('displays public/private classification', () => {
    render(<SubnetDetail subnet={mockSubnet} />);
    // Check for classification badge specifically
    const classificationElements = screen.getAllByText('Private');
    expect(classificationElements.length).toBeGreaterThan(0);
    // Verify at least one has the classification class
    const hasClassificationClass = classificationElements.some(el => 
      el.classList.contains('classification')
    );
    expect(hasClassificationClass).toBe(true);
  });

  it('displays utilization percentage', () => {
    render(<SubnetDetail subnet={mockSubnet} />);
    expect(screen.getByText('39.4%')).toBeInTheDocument();
  });

  it('displays utilization statistics', () => {
    render(<SubnetDetail subnet={mockSubnet} />);
    expect(screen.getByText('100')).toBeInTheDocument(); // allocated
    expect(screen.getByText('154')).toBeInTheDocument(); // available (254 - 100)
  });

  it('renders progress bar with correct aria attributes', () => {
    render(<SubnetDetail subnet={mockSubnet} />);
    const progressBar = screen.getByRole('progressbar');
    expect(progressBar).toHaveAttribute('aria-valuenow', '39.4');
    expect(progressBar).toHaveAttribute('aria-valuemin', '0');
    expect(progressBar).toHaveAttribute('aria-valuemax', '100');
  });

  it('shows high utilization warning when utilization is above threshold', () => {
    const highUtilSubnet: Subnet = {
      ...mockSubnet,
      utilization: {
        totalIps: 254,
        allocatedIps: 220,
        utilizationPercent: 86.6,
      },
    };

    render(<SubnetDetail subnet={highUtilSubnet} />);
    expect(screen.getByText(/High utilization detected/i)).toBeInTheDocument();
  });

  it('does not show high utilization warning when utilization is below threshold', () => {
    render(<SubnetDetail subnet={mockSubnet} />);
    expect(screen.queryByText(/High utilization detected/i)).not.toBeInTheDocument();
  });

  it('displays cloud provider information when present', () => {
    const cloudSubnet: Subnet = {
      ...mockSubnet,
      locationType: LocationType.CLOUD,
      cloudInfo: {
        provider: CloudProviderType.AWS,
        region: 'us-east-1',
        accountId: '123456789012',
      },
    };

    render(<SubnetDetail subnet={cloudSubnet} />);
    // Use getAllByText since the provider name appears in both the icon title and the text
    const awsElements = screen.getAllByText('AWS');
    expect(awsElements.length).toBeGreaterThan(0);
    expect(screen.getByText('us-east-1')).toBeInTheDocument();
    expect(screen.getByText('123456789012')).toBeInTheDocument();
  });

  it('does not display cloud provider information when not present', () => {
    render(<SubnetDetail subnet={mockSubnet} />);
    expect(screen.queryByText(/Cloud Provider:/i)).not.toBeInTheDocument();
  });
});

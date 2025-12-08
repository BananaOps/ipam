import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import SubnetList from './SubnetList';
import { apiClient } from '../services/api';
import { LocationType, CloudProviderType } from '../types';
import { ToastProvider } from '../contexts/ToastContext';

// Mock the API client
vi.mock('../services/api', () => ({
  apiClient: {
    listSubnets: vi.fn(),
  },
}));

// Helper to render with required providers
const renderWithProviders = (component: React.ReactElement) => {
  return render(
    <BrowserRouter>
      <ToastProvider>
        {component}
      </ToastProvider>
    </BrowserRouter>
  );
};

const mockSubnets = [
  {
    id: '1',
    cidr: '10.0.0.0/24',
    name: 'Test Subnet 1',
    description: 'Test description',
    location: 'datacenter-1',
    locationType: LocationType.DATACENTER,
    details: {
      address: '10.0.0.0',
      netmask: '255.255.255.0',
      wildcard: '0.0.0.255',
      network: '10.0.0.0',
      type: 'Private',
      broadcast: '10.0.0.255',
      hostMin: '10.0.0.1',
      hostMax: '10.0.0.254',
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
  },
  {
    id: '2',
    cidr: '192.168.1.0/24',
    name: 'AWS Subnet',
    description: 'AWS subnet',
    location: 'us-east-1',
    locationType: LocationType.CLOUD,
    cloudInfo: {
      provider: CloudProviderType.AWS,
      region: 'us-east-1',
      accountId: '123456789012',
    },
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
      allocatedIps: 220,
      utilizationPercent: 86.6,
    },
    createdAt: Date.now(),
    updatedAt: Date.now(),
  },
];

describe('SubnetList', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders loading state initially', () => {
    vi.mocked(apiClient.listSubnets).mockImplementation(
      () => new Promise(() => {}) // Never resolves
    );

    renderWithProviders(<SubnetList />);

    expect(screen.getByText('Loading subnets...')).toBeInTheDocument();
  });

  it('renders subnet list after loading', async () => {
    vi.mocked(apiClient.listSubnets).mockResolvedValue({
      subnets: mockSubnets,
      totalCount: 2,
    });

    renderWithProviders(<SubnetList />);

    await waitFor(() => {
      expect(screen.getByText('10.0.0.0/24')).toBeInTheDocument();
      expect(screen.getByText('192.168.1.0/24')).toBeInTheDocument();
      expect(screen.getByText('Test Subnet 1')).toBeInTheDocument();
      expect(screen.getByText('AWS Subnet')).toBeInTheDocument();
    });
  });

  it('renders empty state when no subnets', async () => {
    vi.mocked(apiClient.listSubnets).mockResolvedValue({
      subnets: [],
      totalCount: 0,
    });

    renderWithProviders(<SubnetList />);

    await waitFor(() => {
      expect(screen.getByText('No subnets found')).toBeInTheDocument();
      expect(screen.getByText('Get started by creating your first subnet.')).toBeInTheDocument();
    });
  });

  it('displays cloud provider information for cloud subnets', async () => {
    vi.mocked(apiClient.listSubnets).mockResolvedValue({
      subnets: [mockSubnets[1]], // AWS subnet
      totalCount: 1,
    });

    renderWithProviders(<SubnetList />);

    await waitFor(() => {
      // Check for cloud account ID which is unique to cloud info
      expect(screen.getByText('123456789012')).toBeInTheDocument();
      // Check that cloud provider icon is present
      const cloudInfo = screen.getByText('123456789012').closest('.cloud-info');
      expect(cloudInfo).toBeInTheDocument();
    });
  });

  it('displays high utilization warning', async () => {
    vi.mocked(apiClient.listSubnets).mockResolvedValue({
      subnets: [mockSubnets[1]], // 86.6% utilization
      totalCount: 1,
    });

    renderWithProviders(<SubnetList />);

    await waitFor(() => {
      const utilizationBar = screen.getByText('86.6%').previousElementSibling;
      const fill = utilizationBar?.querySelector('.utilization-fill');
      expect(fill).toHaveClass('high');
    });
  });

  it('renders error state on API failure', async () => {
    vi.mocked(apiClient.listSubnets).mockRejectedValue(
      new Error('Network error')
    );

    renderWithProviders(<SubnetList />);

    await waitFor(() => {
      // Error message appears in both ErrorMessage component and Toast
      const errorMessages = screen.getAllByText('Network error');
      expect(errorMessages.length).toBeGreaterThan(0);
      expect(screen.getByText('Retry')).toBeInTheDocument();
    });
  });
});

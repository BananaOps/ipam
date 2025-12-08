import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import CloudProviderIcon from './CloudProviderIcon';
import { CloudProviderType } from '../types';

describe('CloudProviderIcon', () => {
  it('renders icon for AWS provider', () => {
    render(<CloudProviderIcon provider={CloudProviderType.AWS} />);
    const icon = screen.getByLabelText(/AWS cloud provider/i);
    expect(icon).toBeInTheDocument();
  });

  it('renders icon for Azure provider', () => {
    render(<CloudProviderIcon provider={CloudProviderType.AZURE} />);
    const icon = screen.getByLabelText(/AZURE cloud provider/i);
    expect(icon).toBeInTheDocument();
  });

  it('renders icon for GCP provider', () => {
    render(<CloudProviderIcon provider={CloudProviderType.GCP} />);
    const icon = screen.getByLabelText(/GCP cloud provider/i);
    expect(icon).toBeInTheDocument();
  });

  it('renders icon for Scaleway provider', () => {
    render(<CloudProviderIcon provider={CloudProviderType.SCALEWAY} />);
    const icon = screen.getByLabelText(/SCALEWAY cloud provider/i);
    expect(icon).toBeInTheDocument();
  });

  it('renders icon for OVH provider', () => {
    render(<CloudProviderIcon provider={CloudProviderType.OVH} />);
    const icon = screen.getByLabelText(/OVH cloud provider/i);
    expect(icon).toBeInTheDocument();
  });

  it('applies custom className', () => {
    render(<CloudProviderIcon provider={CloudProviderType.AWS} className="custom-class" />);
    const icon = screen.getByLabelText(/AWS cloud provider/i);
    expect(icon).toHaveClass('custom-class');
  });

  it('applies provider-specific class', () => {
    render(<CloudProviderIcon provider={CloudProviderType.AWS} />);
    const icon = screen.getByLabelText(/AWS cloud provider/i);
    expect(icon).toHaveClass('aws');
  });

  it('uses custom title when provided', () => {
    render(<CloudProviderIcon provider={CloudProviderType.AWS} title="Amazon Web Services" />);
    const icon = screen.getByLabelText(/Amazon Web Services cloud provider/i);
    expect(icon).toBeInTheDocument();
  });

  it('applies size prop correctly', () => {
    render(<CloudProviderIcon provider={CloudProviderType.AWS} size="2x" />);
    const icon = screen.getByLabelText(/AWS cloud provider/i);
    expect(icon).toHaveClass('fa-2x');
  });
});

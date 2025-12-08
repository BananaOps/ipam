import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import ErrorMessage from './ErrorMessage';
import type { APIError } from '../types';

describe('ErrorMessage', () => {
  it('renders string error message', () => {
    render(<ErrorMessage error="Something went wrong" />);
    expect(screen.getByText('Something went wrong')).toBeInTheDocument();
  });

  it('renders Error object', () => {
    const error = new Error('Test error message');
    render(<ErrorMessage error={error} />);
    expect(screen.getByText('Test error message')).toBeInTheDocument();
  });

  it('renders APIError with code', () => {
    const apiError: APIError = {
      code: 'NETWORK_ERROR',
      message: 'Unable to connect',
      timestamp: Date.now(),
    };
    render(<ErrorMessage error={apiError} />);
    expect(screen.getByText('Unable to connect')).toBeInTheDocument();
    expect(screen.getByText('Error Code: NETWORK_ERROR')).toBeInTheDocument();
  });

  it('shows details when showDetails is true', () => {
    const apiError: APIError = {
      code: 'VALIDATION_ERROR',
      message: 'Invalid input',
      details: { field: 'cidr', reason: 'Invalid format' },
      timestamp: Date.now(),
    };
    render(<ErrorMessage error={apiError} showDetails={true} />);
    expect(screen.getByText(/field:/)).toBeInTheDocument();
    expect(screen.getByText(/Invalid format/)).toBeInTheDocument();
  });

  it('calls onRetry when retry button is clicked', () => {
    const onRetry = vi.fn();
    render(<ErrorMessage error="Test error" onRetry={onRetry} />);
    
    const retryButton = screen.getByText('Retry');
    fireEvent.click(retryButton);
    
    expect(onRetry).toHaveBeenCalledTimes(1);
  });

  it('calls onDismiss when dismiss button is clicked', () => {
    const onDismiss = vi.fn();
    render(<ErrorMessage error="Test error" onDismiss={onDismiss} />);
    
    const dismissButton = screen.getByText('Dismiss');
    fireEvent.click(dismissButton);
    
    expect(onDismiss).toHaveBeenCalledTimes(1);
  });

  it('does not render action buttons when callbacks are not provided', () => {
    render(<ErrorMessage error="Test error" />);
    expect(screen.queryByText('Retry')).not.toBeInTheDocument();
    expect(screen.queryByText('Dismiss')).not.toBeInTheDocument();
  });
});

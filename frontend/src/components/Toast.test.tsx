import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import Toast from './Toast';

describe('Toast', () => {
  it('renders nothing when isVisible is false', () => {
    const { container } = render(
      <Toast
        message="Test message"
        type="success"
        isVisible={false}
        onClose={vi.fn()}
      />
    );
    expect(container.firstChild).toBeNull();
  });

  it('renders toast when isVisible is true', () => {
    render(
      <Toast
        message="Operation successful"
        type="success"
        isVisible={true}
        onClose={vi.fn()}
      />
    );
    expect(screen.getByText('Operation successful')).toBeInTheDocument();
  });

  it('calls onClose when close button is clicked', () => {
    const onClose = vi.fn();
    render(
      <Toast
        message="Test message"
        type="info"
        isVisible={true}
        onClose={onClose}
      />
    );
    const closeButton = screen.getByLabelText('Close notification');
    fireEvent.click(closeButton);
    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('auto-dismisses after duration', async () => {
    const onClose = vi.fn();
    render(
      <Toast
        message="Test message"
        type="success"
        isVisible={true}
        onClose={onClose}
        duration={100}
      />
    );
    
    await waitFor(() => expect(onClose).toHaveBeenCalledTimes(1), { timeout: 200 });
  });

  it('applies correct CSS class for each type', () => {
    const { rerender } = render(
      <Toast
        message="Test"
        type="success"
        isVisible={true}
        onClose={vi.fn()}
      />
    );
    expect(screen.getByText('Test').closest('.toast')).toHaveClass('toast-success');

    rerender(
      <Toast
        message="Test"
        type="error"
        isVisible={true}
        onClose={vi.fn()}
      />
    );
    expect(screen.getByText('Test').closest('.toast')).toHaveClass('toast-error');
  });
});

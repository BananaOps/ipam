import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import Logo from './Logo';

describe('Logo', () => {
  it('renders logo image', () => {
    render(<Logo />);
    const logo = screen.getByAltText('IPAM by BananaOps Logo');
    expect(logo).toBeInTheDocument();
  });

  it('renders with text by default', () => {
    render(<Logo />);
    expect(screen.getByText('IPAM')).toBeInTheDocument();
    expect(screen.getByText('by BananaOps')).toBeInTheDocument();
  });

  it('hides text when showText is false', () => {
    render(<Logo showText={false} />);
    expect(screen.queryByText('IPAM')).not.toBeInTheDocument();
    expect(screen.queryByText('by BananaOps')).not.toBeInTheDocument();
  });

  it('uses compact variant by default', () => {
    render(<Logo />);
    const logo = screen.getByAltText('IPAM by BananaOps Logo');
    expect(logo).toHaveAttribute('src', '/logo-horizontal.svg');
  });

  it('uses full variant when specified', () => {
    render(<Logo variant="full" />);
    const logo = screen.getByAltText('IPAM by BananaOps Logo');
    expect(logo).toHaveAttribute('src', '/logo.svg');
  });

  it('applies size class correctly', () => {
    const { container } = render(<Logo size="large" />);
    const logoContainer = container.querySelector('.logo-container');
    expect(logoContainer).toHaveClass('large');
  });

  it('applies custom className', () => {
    const { container } = render(<Logo className="custom-class" />);
    const logoContainer = container.querySelector('.logo-container');
    expect(logoContainer).toHaveClass('custom-class');
  });

  it('renders with medium size by default', () => {
    const { container } = render(<Logo />);
    const logoContainer = container.querySelector('.logo-container');
    expect(logoContainer).toHaveClass('medium');
  });
});

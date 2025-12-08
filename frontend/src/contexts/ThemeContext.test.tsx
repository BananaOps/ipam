import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { ThemeProvider, useTheme } from './ThemeContext';
import { THEME_STORAGE_KEY } from '../config/theme';

// Test component that uses the theme context
function TestComponent() {
  const { theme, effectiveTheme, setTheme } = useTheme();
  
  return (
    <div>
      <div data-testid="theme">{theme}</div>
      <div data-testid="effective-theme">{effectiveTheme}</div>
      <button onClick={() => setTheme('dark')}>Set Dark</button>
      <button onClick={() => setTheme('light')}>Set Light</button>
      <button onClick={() => setTheme('auto')}>Set Auto</button>
    </div>
  );
}

describe('ThemeContext', () => {
  beforeEach(() => {
    // Clear localStorage before each test
    localStorage.clear();
    
    // Reset document classes
    document.documentElement.className = '';
    document.documentElement.removeAttribute('data-theme');
  });

  it('should provide default theme as auto when no stored preference exists', () => {
    render(
      <ThemeProvider>
        <TestComponent />
      </ThemeProvider>
    );

    expect(screen.getByTestId('theme').textContent).toBe('auto');
  });

  it('should persist theme to localStorage when changed', async () => {
    render(
      <ThemeProvider>
        <TestComponent />
      </ThemeProvider>
    );

    const darkButton = screen.getByText('Set Dark');
    darkButton.click();

    await waitFor(() => {
      expect(localStorage.getItem(THEME_STORAGE_KEY)).toBe('dark');
      expect(screen.getByTestId('theme').textContent).toBe('dark');
    });
  });

  it('should apply theme class to document root', async () => {
    render(
      <ThemeProvider>
        <TestComponent />
      </ThemeProvider>
    );

    const lightButton = screen.getByText('Set Light');
    lightButton.click();

    await waitFor(() => {
      expect(document.documentElement.classList.contains('theme-light')).toBe(true);
      expect(document.documentElement.getAttribute('data-theme')).toBe('light');
    });
  });

  it('should throw error when useTheme is used outside ThemeProvider', () => {
    // Suppress console.error for this test
    const consoleError = vi.spyOn(console, 'error').mockImplementation(() => {});

    expect(() => {
      render(<TestComponent />);
    }).toThrow('useTheme must be used within a ThemeProvider');

    consoleError.mockRestore();
  });
});

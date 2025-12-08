import { useState, useEffect } from 'react';
import type { EffectiveTheme } from '../config/theme';

/**
 * Hook to detect and track system theme preference
 * Listens to system theme changes and returns current system preference
 */
export function useSystemTheme(): EffectiveTheme {
  const getSystemTheme = (): EffectiveTheme => {
    if (typeof window === 'undefined') return 'light';
    
    return window.matchMedia('(prefers-color-scheme: dark)').matches
      ? 'dark'
      : 'light';
  };

  const [systemTheme, setSystemTheme] = useState<EffectiveTheme>(getSystemTheme);

  useEffect(() => {
    // Create media query for dark mode
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    
    // Handler for theme changes
    const handleChange = (e: MediaQueryListEvent) => {
      setSystemTheme(e.matches ? 'dark' : 'light');
    };

    // Add listener for system theme changes
    mediaQuery.addEventListener('change', handleChange);

    // Cleanup listener on unmount
    return () => {
      mediaQuery.removeEventListener('change', handleChange);
    };
  }, []);

  return systemTheme;
}

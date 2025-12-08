import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useSystemTheme } from '../hooks/useSystemTheme';
import type { Theme, EffectiveTheme } from '../config/theme';
import { THEME_STORAGE_KEY } from '../config/theme';

interface ThemeContextValue {
  theme: Theme;
  effectiveTheme: EffectiveTheme;
  setTheme: (theme: Theme) => void;
}

const ThemeContext = createContext<ThemeContextValue | undefined>(undefined);

interface ThemeProviderProps {
  children: ReactNode;
}

/**
 * ThemeProvider component that manages theme state and persistence
 * Supports dark, light, and auto (system) themes
 */
export function ThemeProvider({ children }: ThemeProviderProps) {
  const systemTheme = useSystemTheme();
  
  // Load theme from localStorage or default to 'auto'
  const loadTheme = (): Theme => {
    try {
      const stored = localStorage.getItem(THEME_STORAGE_KEY);
      if (stored === 'dark' || stored === 'light' || stored === 'auto') {
        return stored;
      }
    } catch (error) {
      console.warn('Failed to load theme from localStorage:', error);
    }
    return 'auto';
  };

  const [theme, setThemeState] = useState<Theme>(loadTheme);

  // Calculate effective theme based on current theme setting
  const effectiveTheme: EffectiveTheme = theme === 'auto' ? systemTheme : theme;

  // Persist theme to localStorage
  const setTheme = (newTheme: Theme) => {
    try {
      localStorage.setItem(THEME_STORAGE_KEY, newTheme);
      setThemeState(newTheme);
    } catch (error) {
      console.warn('Failed to save theme to localStorage:', error);
      setThemeState(newTheme);
    }
  };

  // Apply theme to document root
  useEffect(() => {
    const root = document.documentElement;
    
    // Remove existing theme classes
    root.classList.remove('theme-dark', 'theme-light');
    
    // Add current theme class
    root.classList.add(`theme-${effectiveTheme}`);
    
    // Set data attribute for CSS
    root.setAttribute('data-theme', effectiveTheme);
  }, [effectiveTheme]);

  const value: ThemeContextValue = {
    theme,
    effectiveTheme,
    setTheme,
  };

  return (
    <ThemeContext.Provider value={value}>
      {children}
    </ThemeContext.Provider>
  );
}

/**
 * Hook to access theme context
 * Must be used within ThemeProvider
 */
export function useTheme(): ThemeContextValue {
  const context = useContext(ThemeContext);
  
  if (context === undefined) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  
  return context;
}

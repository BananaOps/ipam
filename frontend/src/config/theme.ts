// Color palette constants for Cyber Minimal theme
export const colors = {
  // Dark theme colors
  darkPrimary: '#0A1A2F',
  cyanAccent: '#0EA5E9',
  
  // Light theme colors
  lightGray: '#F3F4F6',
  white: '#FFFFFF',
  
  // Shared colors
  darkGray: '#1F2937',
  mediumGray: '#6B7280',
  lightBorder: '#E5E7EB',
  darkBorder: '#374151',
  
  // Status colors
  success: '#10B981',
  warning: '#F59E0B',
  error: '#EF4444',
  info: '#3B82F6',
} as const;

export type Theme = 'dark' | 'light' | 'auto';
export type EffectiveTheme = 'dark' | 'light';

export const THEME_STORAGE_KEY = 'ipam-theme-preference';

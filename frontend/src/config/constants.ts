// Application configuration constants

// API Configuration
export const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1';

// Application Name
export const APP_NAME = import.meta.env.VITE_APP_NAME || 'IPAM by BananaOps';

// Color Palette
export const COLORS = {
  darkPrimary: '#0A1A2F',
  cyanAccent: '#0EA5E9',
  lightGray: '#F3F4F6',
  white: '#FFFFFF',
} as const;

// Theme Configuration
export const THEME_STORAGE_KEY = 'ipam-theme-preference';

// Utilization Threshold
export const HIGH_UTILIZATION_THRESHOLD = parseInt(
  import.meta.env.VITE_HIGH_UTILIZATION_THRESHOLD || '80',
  10
);

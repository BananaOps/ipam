// Application configuration constants

export const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1';

export const COLORS = {
  darkPrimary: '#0A1A2F',
  cyanAccent: '#0EA5E9',
  lightGray: '#F3F4F6',
  white: '#FFFFFF',
} as const;

export const THEME_STORAGE_KEY = 'ipam-theme-preference';

export const HIGH_UTILIZATION_THRESHOLD = 80; // percentage

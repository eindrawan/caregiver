export const colors = {
  // Primary colors
  primary: '#0F766E', // Dark teal for status cards
  primaryLight: '#14B8A6', // Lighter teal
  
  // Accent colors
  accent: '#F97316', // Orange for numbers and highlights
  accentLight: '#FB923C',
  
  // Status colors
  success: '#10B981', // Green for completed
  warning: '#F59E0B', // Yellow for warnings
  error: '#D32F2F', // Red for missed/cancelled
  // cancelled: '#D32F2F', // Red for cancelled
  
  // Neutral colors
  white: '#FFFFFF',
  gray50: '#F9FAFB',
  gray100: '#F3F4F6',
  gray200: '#E5E7EB',
  gray300: '#D1D5DB',
  gray400: '#9CA3AF',
  gray500: '#6B7280',
  gray600: '#4B5563',
  gray700: '#374151',
  gray800: '#1F2937',
  gray900: '#111827',
  
  // Background colors
  background: '#F9FAFB',
  cardBackground: '#FFFFFF',
  accentBackground: '#e5f4ff',
  accentBackgroundDark: '#2DA6FF80',
  accentBackgroundLight: '#2DA6FF0A',

  // Text colors
  textPrimary: '#212121',
  textSecondary: '#4B5563',
  textLight: '#9CA3AF',
  textOnPrimary: '#FFFFFF',
} as const;

export type ColorKey = keyof typeof colors;

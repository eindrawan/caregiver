import React from 'react';
import { View, StyleSheet, ViewStyle } from 'react-native';
import { colors, spacing, borderRadius } from '../../constants';
import Text from './Text';

interface BadgeProps {
  variant?: 'scheduled' | 'in_progress' | 'completed' | 'missed' | 'cancelled';
  children: React.ReactNode;
  style?: ViewStyle;
}

const Badge: React.FC<BadgeProps> = ({
  variant = 'scheduled',
  children,
  style
}) => {
  const badgeStyle = [
    styles.base,
    styles[variant],
    style,
  ];

  const textColor = getTextColor(variant);

  return (
    <View style={badgeStyle}>
      <Text variant="caption" color={textColor}>
        {children}
      </Text>
    </View>
  );
};

const getTextColor = (variant: string): keyof typeof colors => {
  switch (variant) {
    case 'scheduled':
      return 'textOnPrimary';
    case 'in_progress':
      return 'textOnPrimary';
    case 'completed':
      return 'textOnPrimary';
    case 'cancelled':
      return 'textOnPrimary';
    case 'missed':
      return 'textOnPrimary';
    default:
      return 'textPrimary';
  }
};

const styles = StyleSheet.create({
  base: {
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    borderRadius: borderRadius.pill,
    alignSelf: 'flex-start',
  },

  // Status variants
  scheduled: {
    backgroundColor: colors.gray600,
  },
  in_progress: {
    backgroundColor: colors.accent,
  },
  completed: {
    backgroundColor: colors.success,
  },
  missed: {
    backgroundColor: colors.error,
  },
  cancelled: {
    backgroundColor: colors.gray600,
  },
});

export default Badge;

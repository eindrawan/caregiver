import React from 'react';
import { View, StyleSheet } from 'react-native';
import { colors, spacing, borderRadius, shadows } from '../../constants';
import { Text } from '../atoms';

interface SummaryCardProps {
  value: number;
  label: string;
  valueColor?: keyof typeof colors;
}

const SummaryCard: React.FC<SummaryCardProps> = ({
  value,
  label,
  valueColor = 'textPrimary'
}) => {
  return (
    <View style={styles.container}>
      <Text
        variant="h1"
        color={valueColor}
        style={styles.value}
      >
        {value}
      </Text>
      <Text
        variant="bodySmall"
        color="textSecondary"
        style={styles.label}
      >
        {label}
      </Text>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.cardBackground,
    borderRadius: borderRadius.lg,
    padding: spacing.lg,
    alignItems: 'center',
    justifyContent: 'center',
    flex: 1,
    minHeight: 100,
    ...shadows.card,
  },
  value: {
    fontSize: 34,
    fontWeight: '500',
    marginBottom: 10,
    textAlign: 'center',
  },
  label: {
    textAlign: 'center',
    lineHeight: 18,
    flexWrap: 'wrap',
  },
});

export default SummaryCard;

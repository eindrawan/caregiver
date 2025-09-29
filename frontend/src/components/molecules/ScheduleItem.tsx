import React from 'react';
import { View, StyleSheet } from 'react-native';
import { colors, spacing, borderRadius } from '../../constants';
import { Text, Icon } from '../atoms';

interface ScheduleItemProps {
  location: string;
  dateTime: string;
  timeRange: string;
  textColor?: keyof typeof colors;
  iconColor?: keyof typeof colors;
  iconAccentColor?: keyof typeof colors;
}

const ScheduleItem: React.FC<ScheduleItemProps> = ({
  location,
  dateTime,
  timeRange,
  textColor = 'textSecondary',
  iconColor = '#666',
  iconAccentColor = '#02cad1'
}) => {
  return (
    <View style={styles.container}>
      <View style={styles.row}>
        <Icon name="location" size={16} color={iconColor} />
        <Text variant="bodySmall" color={textColor} style={styles.text}>
          {location}
        </Text>
      </View>

      <View style={styles.dateTimeContainer}>
        <View style={styles.row}>
          <Icon name="calendar-outline" size={16} color={iconAccentColor} />
          <Text variant="bodySmall" color="textSecondary" style={styles.text}>
            {dateTime}
          </Text>
        </View>

        <View style={styles.divider} />

        <View style={styles.row}>
          <Icon name="time-outline" size={16} color={iconAccentColor} />
          <Text variant="bodySmall" color="textSecondary" style={styles.text}>
            {timeRange}
          </Text>
        </View>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    gap: spacing.xs,
  },
  row: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  text: {
    marginLeft: spacing.xs,
    flex: 1,
  },
  dateTimeContainer: {
    backgroundColor: '#e5f4ff',
    borderRadius: borderRadius.md,
    padding: spacing.md,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  divider: {
    width: 1,
    height: 20,
    backgroundColor: colors.textSecondary,
    marginHorizontal: spacing.md,
  },
});

export default ScheduleItem;

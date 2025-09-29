import React from 'react';
import { View, StyleSheet, FlatList } from 'react-native';
import { spacing, colors, borderRadius } from '../../constants';
import { Text } from '../atoms';
import ScheduleCard from './ScheduleCard';
import { Schedule } from '../../services/types';

interface ScheduleListProps {
  schedules: Schedule[];
  onClockIn?: (schedule: Schedule) => void;
  onClockOut?: (schedule: Schedule) => void;
  onViewProgress?: (schedule: Schedule) => void;
  onMoreOptions?: (schedule: Schedule) => void;
  onSchedulePress?: (schedule: Schedule) => void;
  showSeeAll?: boolean;
  onSeeAll?: () => void;
}

const ScheduleList: React.FC<ScheduleListProps> = ({
  schedules,
  onClockIn,
  onClockOut,
  onViewProgress,
  onMoreOptions,
  onSchedulePress,
  showSeeAll = true,
  onSeeAll
}) => {
  const renderScheduleCard = ({ item }: { item: Schedule }) => (
    <ScheduleCard
      schedule={item}
      onClockIn={() => onClockIn?.(item)}
      onClockOut={() => onClockOut?.(item)}
      onViewProgress={() => onViewProgress?.(item)}
      onMoreOptions={() => onMoreOptions?.(item)}
      onPress={() => onSchedulePress?.(item)}
    />
  );

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <View style={styles.titleRow}>
          <Text variant="h3" color="textPrimary">
            Schedule
          </Text>
          <View style={styles.badge}>
            <Text variant="caption" color="textOnPrimary">
              {schedules.length}
            </Text>
          </View>
        </View>

        {showSeeAll && (
          <Text
            variant="bodySmall"
            color="primary"
            onPress={onSeeAll}
            style={styles.seeAll}
          >
            See All
          </Text>
        )}
      </View>

      <FlatList
        data={schedules}
        renderItem={renderScheduleCard}
        keyExtractor={(item) => item.id.toString()}
        contentContainerStyle={styles.list}
        showsVerticalScrollIndicator={false}
        scrollEnabled={false} // Disable scroll since this is inside a ScrollView
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    gap: spacing.md,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  titleRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
  },
  badge: {
    backgroundColor: colors.primary,
    borderRadius: borderRadius.pill,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    minWidth: 24,
    alignItems: 'center',
    justifyContent: 'center',
  },
  seeAll: {

  },
  list: {
    gap: spacing.md,
  },
});

export default ScheduleList;

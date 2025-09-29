import React, { useState, useEffect } from 'react';
import { View, StyleSheet, Dimensions } from 'react-native';
import { spacing } from '../../constants';
import { SummaryCard } from '../molecules';
import { ScheduleStats } from '../../services/types';

interface StatsOverviewProps {
  stats: ScheduleStats;
}

const StatsOverview: React.FC<StatsOverviewProps> = ({ stats }) => {
  const [screenWidth, setScreenWidth] = useState(Dimensions.get('window').width);

  useEffect(() => {
    const subscription = Dimensions.addEventListener('change', ({ window }) => {
      setScreenWidth(window.width);
    });

    return () => subscription?.remove();
  }, []);

  const isSmallScreen = screenWidth < 768;

  if (isSmallScreen) {
    return (
      <View style={styles.container}>
        <View style={styles.singleCardRow}>
          <SummaryCard
            value={stats.missed}
            label="Missed Scheduled"
            valueColor="error"
          />
        </View>
        <View style={styles.bottomRow}>
          <SummaryCard
            value={stats.upcoming}
            label="Upcoming Today's Schedule"
            valueColor="accent"
          />
          <SummaryCard
            value={stats.completed}
            label="Today's Completed Schedule"
            valueColor="success"
          />
        </View>
      </View>
    );
  } else {
    return (
      <View style={styles.container}>
        <View style={styles.bottomRow}>
          <SummaryCard
            value={stats.missed}
            label="Missed Scheduled"
            valueColor="error"
          />
          <SummaryCard
            value={stats.upcoming}
            label="Upcoming Today's Schedule"
            valueColor="accent"
          />
          <SummaryCard
            value={stats.completed}
            label="Today's Completed Schedule"
            valueColor="success"
          />
        </View>
      </View>
    );
  }
};

const styles = StyleSheet.create({
  container: {
    gap: spacing.md,
  },
  singleCardRow: {
    // Single card takes full width
  },
  bottomRow: {
    flexDirection: 'row',
    gap: spacing.md,
  },
});

export default StatsOverview;

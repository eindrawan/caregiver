import React from 'react';
import { View, StyleSheet } from 'react-native';
import { colors, spacing, borderRadius, shadows } from '../../constants';
import { Text, Button, Icon } from '../atoms';
import { UserInfo } from '../molecules';

interface StatusCardProps {
  user: {
    name: string;
    avatarUrl?: string;
  };
  location: string;
  timeRange: string;
  timer: string; // Format: "HH:MM:SS"
  onClockOut: () => void;
  isInProgress?: boolean;
}

const StatusCard: React.FC<StatusCardProps> = ({
  user,
  location,
  timeRange,
  timer,
  onClockOut,
  isInProgress = true
}) => {
  return (
    <View style={styles.container}>
      {/* Timer Display */}
      <View style={styles.timerSection}>
        <Text variant="h1" color="textOnPrimary" style={styles.timerText}>
          {timer}
        </Text>
      </View>

      {/* User Info */}
      <View style={styles.userSection}>
        <UserInfo
          name={user.name}
          avatarSource={user.avatarUrl ? { uri: user.avatarUrl } : undefined}
          size="medium"
          textColor="textOnPrimary"
          secondaryTextColor="textOnPrimary"
        />
      </View>

      {/* Location */}
      <View style={styles.locationSection}>
        <View style={styles.row}>
          <Icon name="location" size={16} color="white" />
          <Text variant="bodySmall" color="textOnPrimary" style={styles.locationText}>
            {location}
          </Text>
        </View>
      </View>

      {/* Time Range */}
      <View style={styles.timeSection}>
        <View style={styles.row}>
          <Icon name="time" size={16} color="white" />
          <Text variant="bodySmall" color="textOnPrimary" style={styles.timeText}>
            {timeRange}
          </Text>
        </View>
      </View>

      {/* Clock Out Button */}
      <Button
        variant="secondary"
        onPress={onClockOut}
        fullWidth
      >
        <View style={styles.buttonContent}>
          <Icon name="stopwatch" size={24} color="primary" />
          <Text variant="button" color="primary" style={styles.buttonText}>
            Clock-Out
          </Text>
        </View>
      </Button>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.primary,
    borderRadius: borderRadius.xl,
    padding: spacing.xl,
    ...shadows.card,
  },
  timerSection: {
    alignItems: 'center',
    marginBottom: spacing.lg,
  },
  timerText: {
    fontSize: 24,
    padding: 5,
    fontWeight: '500',
    letterSpacing: 2,
  },
  userSection: {
    marginBottom: spacing.lg,
  },
  locationSection: {
    marginBottom: spacing.md,
  },
  timeSection: {
    marginBottom: spacing.xl,
  },
  row: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  locationText: {
    marginLeft: spacing.sm,
    flex: 1,
  },
  timeText: {
    marginLeft: spacing.sm,
    flex: 1,
  },
  buttonContent: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
  },
  buttonText: {
    marginLeft: spacing.sm,
  },
});

export default StatusCard;

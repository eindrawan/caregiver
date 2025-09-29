import React from 'react';
import { View, StyleSheet, TouchableOpacity } from 'react-native';
import { colors, spacing, borderRadius, shadows } from '../../constants';
import { Text, Button, Badge, Icon } from '../atoms';
import { UserInfo, ScheduleItem } from '../molecules';
import { Schedule } from '../../services/types';

interface ScheduleCardProps {
  schedule: Schedule;
  onClockIn?: () => void;
  onClockOut?: () => void;
  onViewProgress?: () => void;
  onMoreOptions?: () => void;
  onPress?: () => void;
}

const ScheduleCard: React.FC<ScheduleCardProps> = ({
  schedule,
  onClockIn,
  onClockOut,
  onViewProgress,
  onMoreOptions,
  onPress
}) => {
  const formatTime = (timeString: string) => {
    return new Date(timeString).toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    });
  };

  const formatDate = (timeString: string) => {
    return new Date(timeString).toLocaleDateString('en-US', {
      weekday: 'short',
      day: '2-digit',
      month: 'short',
      year: 'numeric'
    });
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'scheduled': return 'Scheduled';
      case 'in_progress': return 'In progress';
      case 'completed': return 'Completed';
      case 'missed': return 'Missed';
      case 'cancelled': return 'Cancelled';
      default: return status;
    }
  };

  const renderActionButton = () => {
    switch (schedule.status) {
      case 'scheduled':
      case 'missed':
        return (
          <Button
            variant="primary"
            onPress={onClockIn}
            fullWidth
            rounded
          >
            Clock-In Now
          </Button>
        );
      case 'in_progress':
        return (
          <View style={styles.buttonRow}>
            <Button
              variant="primary"
              outlined={true}
              onPress={onViewProgress}
              style={styles.halfButton}
              rounded
            >
              View Progress
            </Button>
            <Button
              variant="primary"
              onPress={onClockOut}
              style={styles.halfButton}
              rounded
            >
              Clock-Out Now
            </Button>
          </View>
        );
      case 'completed':
        return (
          <Button
            variant="primary"
            outlined={true}
            onPress={onViewProgress}
            fullWidth
            rounded
          >
            View Report
          </Button>
        );
      default:
        return null;
    }
  };

  return (
    <TouchableOpacity style={styles.container} onPress={onPress} activeOpacity={0.7}>
      <View style={styles.header}>
        <Badge variant={schedule.status as any}>
          {getStatusText(schedule.status)}
        </Badge>

        <TouchableOpacity onPress={onMoreOptions}>
          <Icon name="ellipsis-horizontal" size={20} color="textSecondary" />
        </TouchableOpacity>
      </View>

      <UserInfo
        name={schedule.client?.name || 'Unknown Client'}
        serviceName={schedule.service_name || "Service Name A"}
        size="medium"
      />

      <ScheduleItem
        location={schedule.client?.address || 'No address available'}
        dateTime={formatDate(schedule.start_time)}
        timeRange={`${formatTime(schedule.start_time)} - ${formatTime(schedule.end_time)}`}
      />

      {renderActionButton()}
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.cardBackground,
    borderRadius: borderRadius.lg,
    padding: spacing.lg,
    gap: spacing.md,
    marginBottom: spacing.md,
    ...shadows.card,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: spacing.xs,
  },
  buttonRow: {
    flexDirection: 'row',
    gap: spacing.md,
    marginTop: spacing.sm,
  },
  halfButton: {
    flex: 1,
  },
});

export default ScheduleCard;

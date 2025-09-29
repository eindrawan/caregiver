import React, { useState, useEffect } from 'react';
import {
  View,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
} from 'react-native';
import { StackScreenProps } from '@react-navigation/stack';
import { colors, spacing, borderRadius } from '../constants';
import { ContainerView, Text, Button, Icon } from '../components/atoms';
import { UserInfo, TaskItem, HeaderWithBackButton } from '../components/molecules';
import { HomeStackParamList } from '../navigation/HomeStackNavigator';
import { useScheduleById, useEndVisit, useUpdateTaskStatus, useCancelVisit } from '../hooks';
import { showAlert } from '../utils/alert';

type Props = StackScreenProps<HomeStackParamList, 'ClockOut'>;

const ClockOutScreen: React.FC<Props> = ({ route, navigation }) => {
  const { scheduleId } = route.params;
  const { data: schedule, isLoading, error } = useScheduleById(scheduleId);
  const endVisitMutation = useEndVisit();
  const updateTaskMutation = useUpdateTaskStatus();
  const cancelVisitMutation = useCancelVisit();

  const [elapsedTime, setElapsedTime] = useState(0);
  const [pendingTaskUpdates, setPendingTaskUpdates] = useState<{ [key: number]: { status: 'completed' | 'not_completed', reason?: string } }>({});

  // Calculate elapsed time
  useEffect(() => {
    if (schedule?.visit?.start_time) {
      const startTime = new Date(schedule.visit.start_time).getTime();
      const updateTimer = () => {
        const now = Date.now();
        const elapsed = Math.floor((now - startTime) / 1000);
        setElapsedTime(elapsed);
      };

      updateTimer();
      const interval = setInterval(updateTimer, 1000);
      return () => clearInterval(interval);
    }
  }, [schedule?.visit?.start_time]);

  const formatElapsedTime = (seconds: number) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;
    return `${hours.toString().padStart(2, '0')} : ${minutes.toString().padStart(2, '0')} : ${secs.toString().padStart(2, '0')}`;
  };

  const handleTaskStatusChange = (taskId: number, status: 'completed' | 'not_completed', reason?: string) => {
    setPendingTaskUpdates(prev => ({
      ...prev,
      [taskId]: { status, reason }
    }));
  };

  const handleAddNewTask = () => {
    // TODO: Implement add new task functionality
    showAlert('Add New Task', 'This feature will be implemented soon.');
  };

  const handleCancelClockIn = () => {
    showAlert(
      'Cancel Clock-In',
      'Are you sure you want to cancel this clock-in? This will end your current visit.',
      [
        { text: 'No', style: 'cancel' },
        {
          text: 'Yes, Cancel',
          style: 'destructive',
          onPress: async () => {
            if (!schedule) return;

            try {
              await cancelVisitMutation.mutateAsync(schedule.id);
              navigation.navigate('HomeMain')
            } catch (error) {
              showAlert(
                'Cancel Failed',
                'Failed to cancel clock-in. Please try again.',
                [{ text: 'OK' }]
              );
            }
          }
        },
      ]
    );
  };

  const validateTasks = () => {
    if (!schedule?.tasks || schedule.tasks.length === 0) {
      return { isValid: true, message: '' };
    }

    const allTasks = schedule.tasks;
    const unselectedTasks = [];
    const noTasksWithoutReason = [];

    for (const task of allTasks) {
      const pendingUpdate = pendingTaskUpdates[task.id];
      const currentStatus = pendingUpdate?.status || task.status;
      const currentReason = pendingUpdate?.reason || task.reason;

      // Check if task has been selected (Yes or No)
      if (!pendingUpdate && task.status === 'pending') {
        unselectedTasks.push(task.name);
      }

      // Check if "No" was selected but no reason provided
      if (currentStatus === 'not_completed' && (!currentReason || currentReason.trim() === '')) {
        noTasksWithoutReason.push(task.name);
      }
    }

    if (unselectedTasks.length > 0) {
      return {
        isValid: false,
        message: `Please select Yes or No for the following tasks:\n• ${unselectedTasks.join('\n• ')}`
      };
    }

    if (noTasksWithoutReason.length > 0) {
      return {
        isValid: false,
        message: `Please provide a reason for the following tasks marked as "No":\n• ${noTasksWithoutReason.join('\n• ')}`
      };
    }

    return { isValid: true, message: '' };
  };

  const handleClockOut = async () => {
    if (!schedule) return;

    // Validate all tasks before proceeding
    const validation = validateTasks();
    if (!validation.isValid) {
      showAlert(
        'Incomplete Tasks',
        validation.message,
        [{ text: 'OK' }]
      );
      return;
    }

    try {
      // Update all pending task statuses sequentially to avoid database locking
      for (const [taskId, update] of Object.entries(pendingTaskUpdates)) {
        try {
          await updateTaskMutation.mutateAsync({
            taskId: parseInt(taskId),
            data: { status: update.status, reason: update.reason }
          });
        } catch (error) {
          console.error(`Failed to update task ${taskId}:`, error);
          throw error; // Re-throw to be caught by the outer catch block
        }
      }

      // End the visit
      const mockLocation = {
        end_latitude: 40.7128,
        end_longitude: -74.0060,
        notes: 'Visit completed successfully',
      };

      await endVisitMutation.mutateAsync({
        scheduleId: schedule.id,
        data: mockLocation,
      });

      // Calculate duration for display
      const durationInSeconds = elapsedTime;
      const hours = Math.floor(durationInSeconds / 3600);
      const minutes = Math.floor((durationInSeconds % 3600) / 60);

      let durationText = '';
      if (hours > 0) {
        durationText = `${hours} hour${hours > 1 ? 's' : ''}`;
        if (minutes > 0) {
          durationText += ` ${minutes} minute${minutes > 1 ? 's' : ''}`;
        }
      } else {
        durationText = `${minutes} minute${minutes > 1 ? 's' : ''}`;
      }

      // Navigate to success screen
      navigation.navigate('ScheduleCompleted', {
        schedule,
        duration: durationText,
      });
    } catch (error) {
      showAlert(
        'Clock Out Failed',
        'Failed to clock out. Please try again.',
        [{ text: 'OK' }]
      );
    }
  };

  if (isLoading) {
    return (
      <ContainerView style={styles.container}>
        <View style={styles.loadingContainer}>
          <Text variant="body" color="textSecondary">Loading...</Text>
        </View>
      </ContainerView>
    );
  }

  if (error || !schedule) {
    return (
      <ContainerView style={styles.container}>
        <View style={styles.loadingContainer}>
          <Text variant="body" color="textSecondary">
            {error ? 'Error loading schedule' : 'Schedule not found'}
          </Text>
        </View>
      </ContainerView>
    );
  }

  const currentUser = {
    name: schedule.client?.name || 'Unknown Client',
    avatar: undefined, // Will use default avatar
  };

  return (
    <ContainerView style={styles.container}>
      {/* Header  */}
      <HeaderWithBackButton
        title="Clock-Out"
        onBackPress={() => navigation.goBack()}
      />
      <ScrollView style={styles.scrollView} contentContainerStyle={styles.content} keyboardShouldPersistTaps="handled">
        {/* Timer */}
        <View style={styles.timerContainer}>
          <Text variant="h1" color="textPrimary" style={styles.timer}>
            {formatElapsedTime(elapsedTime)}
          </Text>
        </View>

        {/* Service Info */}
        <View style={styles.serviceSection}>
          <Text variant="h3" color="textPrimary" style={styles.serviceName}>
            {schedule.service_name || '- no service name -'}
          </Text>
          <UserInfo
            name={schedule.client?.name || 'Unknown Client'}
            size="large"
          />
        </View>

        {/* Tasks Section */}
        <View style={styles.tasksSection}>
          <Text variant="title" style={styles.sectionTitle}>
            Tasks:
          </Text>
          <Text variant="body" color="textSecondary" style={styles.sectionSubtitle}>
            Please tick the tasks that you have done
          </Text>

          {schedule.tasks && schedule.tasks.length > 0 ? (
            schedule.tasks.map((task) => (
              <TaskItem
                key={task.id}
                task={{
                  ...task,
                  // Apply pending updates if any
                  status: pendingTaskUpdates[task.id]?.status || task.status,
                  reason: pendingTaskUpdates[task.id]?.reason || task.reason,
                }}
                onStatusChange={handleTaskStatusChange}
              />
            ))
          ) : (
            <Text variant="body" style={styles.noTasksText}>
              No tasks assigned for this visit.
            </Text>
          )}

          {/* Add New Task Button */}
          <TouchableOpacity style={styles.addTaskButton} onPress={handleAddNewTask}>
            <Icon name="add" size={20} color="primary" />
            <Text variant="button" color="primary" style={styles.addTaskText}>
              Add new task
            </Text>
          </TouchableOpacity>
        </View>

        {/* Clock-In Location */}
        <View>
          <Text variant="title" style={styles.sectionTitle}>
            Clock-In Location
          </Text>
          <View style={styles.locationCard}>
            <View style={styles.mapPlaceholder}>
              <Icon name="location" size={40} color="primary" />
            </View>
            <View style={styles.locationInfo}>
              <Text variant="body" color="textPrimary" style={styles.locationAddress}>
                {schedule.client?.address || 'No address available'}
              </Text>
              <Text variant="body" color="textSecondary" style={styles.locationDetails}>
                {schedule.client?.city}, {schedule.client?.state}, {schedule.client?.zip_code}
              </Text>
            </View>
          </View>
        </View>

        {/* Service Notes */}
        {schedule.notes && (
          <View>
            <Text variant="title" color="textPrimary" style={styles.sectionTitle}>
              Service Notes:
            </Text>
            <Text variant="body" color="textSecondary">
              {schedule.notes}
            </Text>
          </View>
        )}

        {/* Action Buttons */}
        <View style={styles.actionButtonsContainer}>
          <Button
            variant="outline"
            onPress={handleCancelClockIn}
            style={styles.cancelButton}
            disabled={endVisitMutation.isPending || updateTaskMutation.isPending || cancelVisitMutation.isPending}
          >
            {cancelVisitMutation.isPending ? 'Cancelling...' : 'Cancel Clock-in'}
          </Button>

          <Button
            variant="primary"
            onPress={handleClockOut}
            style={styles.clockOutButton}
            disabled={endVisitMutation.isPending || updateTaskMutation.isPending || cancelVisitMutation.isPending}
          >
            {endVisitMutation.isPending || updateTaskMutation.isPending ? 'Processing...' : 'Clock Out'}
          </Button>
        </View>
      </ScrollView>
    </ContainerView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  scrollView: {
    flex: 1,
  },
  content: {
    padding: spacing.screenPadding,
    gap: spacing.lg,
    flexGrow: 1,
  },
  headerSection: {
    marginBottom: spacing.md,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  timerContainer: {
    alignItems: 'center',
    paddingVertical: spacing.lg,
  },
  timer: {
    fontSize: 24,
    fontWeight: 'bold',
    letterSpacing: 2,
  },
  serviceSection: {
    alignItems: 'center',
    gap: spacing.md,
    padding: spacing.lg,
    borderRadius: borderRadius.lg,
    backgroundColor: colors.accentBackgroundLight,
  },
  serviceName: {
    fontWeight: '600',
  },
  tasksSection: {
    gap: spacing.sm,
  },
  sectionTitle: {
    marginBottom: spacing.xs,
  },
  sectionSubtitle: {
  },
  noTasksText: {
    textAlign: 'center',
    fontStyle: 'italic',
    paddingVertical: spacing.lg,
  },
  addTaskButton: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: spacing.md,
    gap: spacing.sm,
    marginTop: spacing.sm,
  },
  addTaskText: {
    fontWeight: '600',
  },
  locationCard: {
    flexDirection: 'row',
    backgroundColor: colors.white,
    borderRadius: borderRadius.md,
    padding: spacing.md,
    borderWidth: 1,
    borderColor: colors.gray200,
    gap: spacing.md,
  },
  mapPlaceholder: {
    width: 80,
    height: 80,
    backgroundColor: colors.gray100,
    borderRadius: borderRadius.sm,
    alignItems: 'center',
    justifyContent: 'center',
  },
  locationInfo: {
    flex: 1,
    justifyContent: 'center',
  },
  locationAddress: {
    fontWeight: '600',
    marginBottom: spacing.xs,
  },
  locationDetails: {
    lineHeight: 20,
  },
  actionButtonsContainer: {
    flexDirection: 'row',
    gap: spacing.md,
    marginTop: spacing.lg,
  },
  cancelButton: {
    flex: 1,
    borderColor: colors.gray300,
  },
  clockOutButton: {
    flex: 1,
  },
});

export default ClockOutScreen;

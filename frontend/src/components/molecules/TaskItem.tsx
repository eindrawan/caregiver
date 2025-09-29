import React, { useState } from 'react';
import {
  View,
  StyleSheet,
  TouchableOpacity,
  TextInput,
} from 'react-native';
import { colors, spacing, borderRadius } from '../../constants';
import { Text, Icon } from '../atoms';
import { Task } from '../../services/types';

interface TaskItemProps {
  task: Task;
  onStatusChange: (taskId: number, status: 'completed' | 'not_completed', reason?: string) => void;
  disabled?: boolean;
}

const TaskItem: React.FC<TaskItemProps> = ({ task, onStatusChange, disabled = false }) => {
  const [reason, setReason] = useState(task.reason || '');
  const [selectedButton, setSelectedButton] = useState<'yes' | 'no' | null>(
    task.status === 'completed' ? 'yes' :
      task.status === 'not_completed' ? 'no' :
        null
  );

  const handleYesPress = () => {
    setSelectedButton('yes');
    onStatusChange(task.id, 'completed');
  };

  const handleNoPress = () => {
    setSelectedButton('no');
    // Call onStatusChange immediately with current reason
    onStatusChange(task.id, 'not_completed', reason.trim() || undefined);
  };

  const handleReasonChange = (newReason: string) => {
    setReason(newReason);
    // If "No" is selected and reason changes, update the status
    if (selectedButton === 'no') {
      onStatusChange(task.id, 'not_completed', newReason.trim() || undefined);
    }
  };

  return (
    <View style={styles.container}>
      {/* Task Title */}
      <Text variant="title" color="primary" style={styles.taskTitle}>
        {task.name}
      </Text>

      {/* Task Description */}
      <Text variant="body" color="textSecondary" style={styles.taskDescription}>
        {task.description}
      </Text>

      {/* Action Buttons - Always show for pending tasks */}
      {!disabled && (
        <View style={styles.actionButtons}>
          <TouchableOpacity
            style={[
              styles.actionButton,
              styles.yesButton,
              selectedButton === 'yes' && { backgroundColor: colors.accentBackground }
            ]}
            onPress={handleYesPress}
          >
            <Icon name="checkmark" size={20} color={colors.success} />
            <Text variant="button" style={[styles.buttonText]}>
              Yes
            </Text>
          </TouchableOpacity>

          <View style={styles.separator} />

          <TouchableOpacity
            style={[
              styles.actionButton,
              styles.noButton,
              selectedButton === 'no' && { backgroundColor: colors.accentBackground }
            ]}
            onPress={handleNoPress}
          >
            <Icon name="close" size={20} color={colors.error} />
            <Text variant="button" style={[styles.buttonText]}>
              No
            </Text>
          </TouchableOpacity>
        </View>
      )}

      {/* Reason Input - Always show for pending tasks */}
      {!disabled && selectedButton === 'no' && (
        <TextInput
          style={styles.reasonInput}
          value={reason}
          onChangeText={handleReasonChange}
          placeholder="Add reason..."
          placeholderTextColor={colors.textSecondary}
          multiline
          numberOfLines={2}
          textAlignVertical="top"
        />
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.white,
    borderRadius: borderRadius.md,
    padding: spacing.md,
    borderColor: colors.gray200,
  },
  taskTitle: {
    marginBottom: spacing.xs,
    color: '#0D5D59',
  },
  taskDescription: {
    fontSize: 14,
    lineHeight: 24,
    fontWeight: '400',
    marginBottom: spacing.sm,
    color: '#1D1D1BDE',
  },
  actionButtons: {
    flexDirection: 'row',
    gap: spacing.sm,
    alignItems: 'center',
    justifyContent: 'flex-start',
  },
  actionButton: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: spacing.xs,
    paddingHorizontal: spacing.md,
    borderRadius: borderRadius.sm,
    gap: spacing.xs,
    borderWidth: 0,
    backgroundColor: colors.gray50,
  },
  yesButton: {
    backgroundColor: 'transparent',
  },
  noButton: {
    backgroundColor: 'transparent',
  },
  buttonIcon: {
    fontSize: 18,
  },
  buttonText: {
    fontWeight: '400',
    fontSize: 14,
    marginTop: 1
  },
  reasonInput: {
    backgroundColor: colors.gray50,
    borderWidth: 1,
    borderColor: colors.gray300,
    borderRadius: borderRadius.lg,
    marginTop: spacing.sm,
    padding: spacing.md,
    fontSize: 16,
    color: colors.textSecondary,
    minHeight: 60,
    textAlignVertical: 'top',
  },
  statusBadge: {
    alignSelf: 'flex-start',
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.xs,
    borderRadius: borderRadius.sm,
    marginTop: spacing.sm,
  },
  statusText: {
    fontWeight: '600',
    fontSize: 12,
  },
  reasonContainer: {
    marginTop: spacing.sm,
    padding: spacing.sm,
    backgroundColor: colors.gray50,
    borderRadius: borderRadius.sm,
  },
  reasonLabel: {
    fontWeight: '600',
    marginBottom: spacing.xs,
    fontSize: 12,
  },
  reasonText: {
    lineHeight: 18,
    fontSize: 14,
  },
  separator: {
    width: 1,
    height: 16,
    backgroundColor: colors.gray300,
    marginHorizontal: spacing.xs,
  },
});

export default TaskItem;

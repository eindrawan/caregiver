
import React, { useState } from 'react';
import { View, StyleSheet, TouchableOpacity, Modal, TextInput, Linking } from 'react-native';
import { colors, spacing, borderRadius, shadows } from '../../constants';
import { Text, Button, Badge, Icon } from '../atoms';
import { UserInfo, ScheduleItem } from '../molecules';
import { Schedule, LocationData } from '../../services/types';
import { getCurrentLocation, geocodeAddress } from '../../services/locationService';
import { showAlert } from '../../utils/alert';

interface ScheduleCardProps {
  schedule: Schedule;
  onClockIn?: (location?: LocationData) => void;
  onClockOut?: (location?: LocationData) => void;
  onViewProgress?: () => void;
  onMoreOptions?: () => void;
  onPress?: () => void;
  disabled?: boolean;
}

const ScheduleCard: React.FC<ScheduleCardProps> = ({
  schedule,
  onClockIn,
  onClockOut,
  onViewProgress,
  onMoreOptions,
  onPress,
  disabled = false
}) => {
  const [isLoadingLocation, setIsLoadingLocation] = useState(false);
  const [showManualInput, setShowManualInput] = useState(false);
  const [manualAddress, setManualAddress] = useState(schedule.client?.address || '');

  const isPendingLocation = schedule.visit?.location_status === 'pending';

  const effectiveDisabled = disabled || isPendingLocation || isLoadingLocation;

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

  const handleLocationAction = async (action: (location?: LocationData) => void) => {
    if (effectiveDisabled) return;

    setIsLoadingLocation(true);

    try {
      const location = await getCurrentLocation();
      action(location);
    } catch (error: any) {
      setIsLoadingLocation(false);
      const errorCode = error.code || 'UNKNOWN_ERROR';

      switch (errorCode) {
        case 'PERMISSION_DENIED':
          showAlert(
            'Location Permission Required',
            'Location access is required for visit verification. Please enable permissions in your device settings and try again.',
            [
              { text: 'Cancel', style: 'cancel' },
              { text: 'Retry', onPress: () => handleLocationAction(action) },
              { text: 'Settings', onPress: () => Linking.openSettings() }
            ]
          );
          break;
        case 'POSITION_UNAVAILABLE':
          setShowManualInput(true);
          break;
        case 'TIMEOUT':
          // Auto-retry once
          try {
            const location = await getCurrentLocation();
            action(location);
          } catch (retryError) {
            setIsLoadingLocation(false);
            setShowManualInput(true);
          }
          break;
        default:
          showAlert(
            'Location Error',
            error.message || 'Failed to get location. You can proceed without location, but it will be flagged as pending.',
            [
              { text: 'Cancel', style: 'cancel' },
              { text: 'Proceed Without Location', onPress: () => action(undefined) },
              { text: 'Retry', onPress: () => handleLocationAction(action) }
            ]
          );
          break;
      }
    }
  };

  const handleManualInput = async () => {
    if (!manualAddress.trim()) {
      showAlert('Invalid Address', 'Please enter an address.');
      return;
    }

    setIsLoadingLocation(true);
    setShowManualInput(false);

    try {
      const location = await geocodeAddress(manualAddress);
      // Call the action with the geocoded location
      if (onClockIn) onClockIn(location);
      else if (onClockOut) onClockOut(location);
    } catch (error) {
      setIsLoadingLocation(false);
      showAlert(
        'Geocoding Error',
        'Could not convert address to location. Proceed without location?',
        [
          { text: 'Cancel', style: 'cancel' },
          {
            text: 'Proceed Without Location', onPress: () => {
              if (onClockIn) onClockIn(undefined);
              else if (onClockOut) onClockOut(undefined);
            }
          }
        ]
      );
    }
  };

  const renderActionButton = () => {
    switch (schedule.status) {
      case 'scheduled':
      case 'missed':
        return (
          <Button
            variant="primary"
            onPress={() => handleLocationAction(onClockIn || (() => { }))}
            fullWidth
            rounded
            disabled={effectiveDisabled}
          >
            {isLoadingLocation ? 'Getting Location...' : 'Clock-In Now'}
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
              disabled={disabled}
            >
              View Progress
            </Button>
            <Button
              variant="primary"
              onPress={() => handleLocationAction(onClockOut || (() => { }))}
              style={styles.halfButton}
              rounded
              disabled={effectiveDisabled}
            >
              {isLoadingLocation ? 'Getting Location...' : 'Clock-Out Now'}
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
            disabled={disabled}
          >
            View Report
          </Button>
        );
      default:
        return null;
    }
  };

  return (
    <>
      <TouchableOpacity style={styles.container} onPress={onPress} activeOpacity={0.7}>
        <View style={styles.header}>
          <View style={styles.statusContainer}>
            <Badge variant={schedule.status as any}>
              {getStatusText(schedule.status)}
            </Badge>
            {isPendingLocation && (
              <Badge variant="in_progress" style={styles.pendingBadge}>
                Pending Location
              </Badge>
            )}
          </View>

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

      {/* Manual Input Modal */}
      <Modal visible={showManualInput} animationType="slide" transparent={true}>
        <View style={styles.modalOverlay}>
          <View style={styles.modalContent}>
            <Text style={styles.modalTitle}>Enter Address Manually</Text>
            <Text style={styles.modalSubtitle}>Location unavailable. Enter address to approximate position.</Text>
            <TextInput
              style={styles.textInput}
              value={manualAddress}
              onChangeText={setManualAddress}
              placeholder="Enter full address"
              multiline
            />
            <View style={styles.modalButtons}>
              <Button
                variant="primary"
                outlined={true}
                onPress={() => setShowManualInput(false)}
                style={styles.modalButton}
              >
                Cancel
              </Button>
              <Button
                variant="primary"
                onPress={handleManualInput}
                style={styles.modalButton}
                disabled={!manualAddress.trim()}
              >
                Geocode & Proceed
              </Button>
              <Button
                variant="primary"
                outlined={true}
                onPress={() => {
                  setShowManualInput(false);
                  if (onClockIn) onClockIn(undefined);
                  else if (onClockOut) onClockOut(undefined);
                }}
                style={styles.modalButton}
              >
                Proceed Without Location
              </Button>
            </View>
          </View>
        </View>
      </Modal>
    </>
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
  statusContainer: {
    flexDirection: 'row',
    gap: spacing.xs,
    alignItems: 'center',
  },
  pendingBadge: {
    marginLeft: spacing.xs,
  },
  buttonRow: {
    flexDirection: 'row',
    gap: spacing.md,
    marginTop: spacing.sm,
  },
  halfButton: {
    flex: 1,
  },
  modalOverlay: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: 'rgba(0,0,0,0.5)',
  },
  modalContent: {
    backgroundColor: colors.cardBackground,
    padding: spacing.lg,
    borderRadius: borderRadius.lg,
    width: '90%',
    maxWidth: 400,
    gap: spacing.md,
  },
  modalTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: colors.textPrimary,
    textAlign: 'center',
  },
  modalSubtitle: {
    fontSize: 14,
    color: colors.textSecondary,
    textAlign: 'center',
  },
  textInput: {
    borderWidth: 1,
    borderColor: '#E5E7EB',
    borderRadius: borderRadius.sm,
    padding: spacing.md,
    backgroundColor: colors.background,
    minHeight: 80,
    textAlignVertical: 'top',
  },
  modalButtons: {
    flexDirection: 'row',
    gap: spacing.sm,
    justifyContent: 'space-around',
    flexWrap: 'wrap',
  },
  modalButton: {
    flex: 1,
    minWidth: 100,

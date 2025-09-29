import React, { useState, useEffect } from 'react';
import {
  View,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Linking,
  useWindowDimensions,
} from 'react-native';
import { StackScreenProps } from '@react-navigation/stack';
import { colors, spacing, borderRadius, shadows } from '../constants';
import { Text, Button, Icon } from '../components/atoms';
import { HeaderWithBackButton, ScheduleItem, UserInfo } from '../components/molecules';
import { HomeStackParamList } from '../navigation/HomeStackNavigator';
import { useScheduleById, useStartVisit } from '../hooks/useSchedules';
import { getCurrentLocation, showLocationErrorAlert, LocationError } from '../services/locationService';
import { showAlert } from '../utils/alert';
import { ContainerView, Footer } from '../components/organisms';

type Props = StackScreenProps<HomeStackParamList, 'ScheduleDetails'>;

const ScheduleDetailsScreen: React.FC<Props> = ({ route, navigation }) => {
  const { scheduleId } = route.params;
  const { data: schedule, isLoading, error } = useScheduleById(scheduleId);
  const startVisitMutation = useStartVisit();
  const [isClockingIn, setIsClockingIn] = useState(false);

  // Detect large screen (e.g., tablets)
  const { width } = useWindowDimensions();
  const isLargeScreen = width >= 768;

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

  const handleCall = (phoneNumber: string | undefined) => {
    if (phoneNumber) {
      Linking.openURL(`tel:${phoneNumber}`);
    }
  };

  const handleEmail = (email: string | undefined) => {
    if (email) {
      Linking.openURL(`mailto:${email}`);
    }
  };

  const handleClockIn = async () => {
    if (isClockingIn) return;

    showAlert(
      'Clock In',
      `Clock in for ${schedule?.client?.name}?`,
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Clock In',
          onPress: async () => {
            setIsClockingIn(true);
            try {
              // Get current location
              const location = await getCurrentLocation();

              // Start the visit with location data
              await startVisitMutation.mutateAsync({
                scheduleId,
                data: {
                  start_latitude: location.latitude,
                  start_longitude: location.longitude,
                },
              });

              navigation.goBack();
            } catch (error) {
              console.error('Clock in error:', error);

              if (error instanceof LocationError) {
                showLocationErrorAlert(
                  error,
                  () => handleClockIn(), // Retry function
                  () => setIsClockingIn(false) // Cancel function
                );
              } else {
                // Handle API errors
                showAlert(
                  'Clock In Failed',
                  'Failed to clock in. Please try again.',
                  [
                    {
                      text: 'OK',
                      onPress: () => setIsClockingIn(false),
                    },
                  ]
                );
              }
            } finally {
              setIsClockingIn(false);
            }
          },
        },
      ]
    );
  };

  if (isLoading) {
    return (
      <ContainerView style={styles.container} title="Schedule Details" onBackPress={() => navigation.goBack()}>
        <View style={styles.loadingContainer}>
          <Text variant="body" color="textSecondary">Loading...</Text>
        </View>
      </ContainerView>
    );
  }

  if (error || !schedule) {
    return (
      <ContainerView style={styles.container} title="Schedule Details" onBackPress={() => navigation.goBack()}>
        <View style={styles.loadingContainer}>
          <Text variant="body" color="textSecondary">
            {error ? 'Error loading schedule' : 'Schedule not found'}
          </Text>
        </View>
      </ContainerView>
    );
  }


  return (
    <ContainerView style={styles.container} title="Schedule Details" onBackPress={() => navigation.goBack()}>
      <ScrollView style={styles.scrollView} contentContainerStyle={styles.content}>
        <View style={styles.headerSection}>
          {/* Service Name */}
          <Text variant="h1" color="textSecondary" style={styles.serviceName}>
            {schedule.service_name || 'Service Name A'}
          </Text>

          {/* Client Info */}
          <View style={styles.clientSection}>
            <UserInfo
              name={schedule.client?.name || 'Unknown Client'}
              size="large"
            />
          </View>

          {/* Date and Time */}
          <ScheduleItem
            dateTime={formatDate(schedule.start_time)}
            timeRange={`${formatTime(schedule.start_time)} - ${formatTime(schedule.end_time)}`}
          />
        </View>

        {/* Client Contact */}
        <View style={styles.section}>
          <Text variant="h3" color="textPrimary" style={styles.sectionTitle}>
            Client Contact:
          </Text>

          {schedule.client?.email && (
            <TouchableOpacity
              style={styles.contactItem}
              onPress={() => handleEmail(schedule.client?.email)}
            >
              <Icon name="mail" size={20} color="textSecondary" />
              <Text variant="body" color="textPrimary" style={styles.contactText}>
                {schedule.client.email}
              </Text>
            </TouchableOpacity>
          )}

          {schedule.client?.phone && (
            <TouchableOpacity
              style={styles.contactItem}
              onPress={() => handleCall(schedule.client?.phone)}
            >
              <Icon name="call" size={20} color="textSecondary" />
              <Text variant="body" color="textPrimary" style={styles.contactText}>
                {schedule.client.phone}
              </Text>
            </TouchableOpacity>
          )}
        </View>

        {/* Address */}
        <View style={styles.section}>
          <Text variant="h3" color="textPrimary" style={styles.sectionTitle}>
            Address:
          </Text>
          <Text variant="body" color="textPrimary" style={styles.addressText}>
            {schedule.client?.address || 'No address available'}
          </Text>
          <Text variant="body" color="textPrimary" style={styles.addressText}>
            {schedule.client?.city}, {schedule.client?.state} {schedule.client?.zip_code}
          </Text>
        </View>

        {/* Tasks */}
        {schedule.tasks && schedule.tasks.length > 0 && (
          <View style={styles.section}>
            <Text variant="h3" color="textPrimary" style={styles.sectionTitle}>
              Tasks:
            </Text>
            {schedule.tasks.map((task) => (
              <View key={task.id} style={styles.taskItem}>
                <Text variant="title" color="primary" style={styles.taskTitle}>
                  {task.name}
                </Text>
                <Text variant="body" color="textSecondary" style={styles.taskDescription}>
                  {task.description}
                </Text>
              </View>
            ))}
          </View>
        )}

        {/* Service Notes */}
        {schedule.notes && (
          <View style={styles.section}>
            <Text variant="h3" color="textPrimary" style={styles.sectionTitle}>
              Service Notes
            </Text>
            <Text variant="body" color="textSecondary" style={styles.notesText}>
              {schedule.notes}
            </Text>
          </View>
        )}

        {/* Clock In Button */}
        {schedule.status === 'scheduled' && (
          <Button
            variant="primary"
            onPress={handleClockIn}
            fullWidth
            disabled={isClockingIn || startVisitMutation.isPending}
            style={styles.clockInButton}
          >
            {isClockingIn || startVisitMutation.isPending ? 'Clocking In...' : 'Clock In Now'}
          </Button>
        )}
        {/* Footer */}
        {isLargeScreen && (<Footer />)}
      </ScrollView>
    </ContainerView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  headerSection: {
    marginBottom: spacing.lg,
    backgroundColor: colors.accentBackgroundLight,
    padding: spacing.lg,
    borderRadius: spacing.lg,
  },
  scrollView: {
    flex: 1,
  },
  content: {
    padding: spacing.screenPadding,
    gap: spacing.lg,
    paddingBottom: spacing.xxxl,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  serviceName: {
    textAlign: 'center',
    marginBottom: spacing.md,
  },
  clientSection: {
    alignItems: 'center',
    marginBottom: spacing.lg,
  },
  dateTimeSection: {
    marginBottom: spacing.lg,
  },
  dateTimeRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    gap: spacing.md,
  },
  dateTimeItem: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: colors.gray50,
    padding: spacing.md,
    borderRadius: borderRadius.md,
  },
  dateTimeText: {
    marginLeft: spacing.sm,
    flex: 1,
  },
  section: {
    marginBottom: spacing.lg,
  },
  sectionTitle: {
    marginBottom: spacing.md,
    fontWeight: '600',
  },
  contactItem: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: spacing.sm,
  },
  contactText: {
    marginLeft: spacing.md,
    flex: 1,
  },
  addressText: {
    lineHeight: 24,
  },
  taskItem: {
    backgroundColor: colors.gray50,
    padding: spacing.md,
    borderRadius: borderRadius.md,
    marginBottom: spacing.sm,
  },
  taskTitle: {
    marginBottom: spacing.xs,
    fontWeight: '600',
  },
  taskDescription: {
    lineHeight: 20,
  },
  notesText: {
    lineHeight: 20,
  },
  clockInButton: {
    marginTop: spacing.lg,
  },
});

export default ScheduleDetailsScreen;

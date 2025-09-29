import React, { useState, useEffect } from 'react';
import {
  View,
  StyleSheet,
  ScrollView,
  RefreshControl,
  Dimensions,
} from 'react-native';
import { StackScreenProps } from '@react-navigation/stack';
import { HomeStackParamList } from '../navigation/HomeStackNavigator';
import { colors, spacing } from '../constants';
import { ContainerView, Text } from '../components/atoms';
import {
  StatusCard,
  StatsOverview,
  ScheduleList,
  CustomHeader
} from '../components/organisms';
import { Schedule } from '../services/types';
import { useTodaySchedules, useScheduleStats, useStartVisit, useEndVisit } from '../hooks/useSchedules';
import { showAlert } from '../utils/alert';

type Props = StackScreenProps<HomeStackParamList, 'HomeMain'>;

const HomeScreen: React.FC<Props> = ({ navigation }) => {
  // React Query hooks
  const {
    data: schedules = [],
    isLoading,
    refetch: refetchSchedules
  } = useTodaySchedules();

  const {
    data: stats,
    refetch: refetchStats
  } = useScheduleStats();

  const startVisitMutation = useStartVisit();
  const endVisitMutation = useEndVisit();

  // Mock user data - in real app this would come from auth context
  const currentUser = {
    name: 'Louis',
    avatarUrl: undefined, // Will show initials
  };

  // Current schedule for status card - only show if user has active clock-in
  const currentSchedule = schedules.find(s => s.status === 'in_progress');

  // Timer state for active visits
  const [elapsedTime, setElapsedTime] = useState(0);

  // Screen width state for responsive design
  const [screenWidth, setScreenWidth] = useState(Dimensions.get('window').width);

  // Check if screen is large (tablet/desktop)
  const isLargeScreen = screenWidth >= 768;

  // Effect for listening to screen dimension changes
  useEffect(() => {
    const subscription = Dimensions.addEventListener('change', ({ window }) => {
      setScreenWidth(window.width);
    });

    return () => subscription?.remove();
  }, []);

  // Timer effect for in-progress schedules
  useEffect(() => {
    let interval: NodeJS.Timeout;

    if (currentSchedule?.status === 'in_progress') {
      // Calculate elapsed time from start_time to now
      const startTime = new Date(currentSchedule.start_time).getTime();
      const updateElapsed = () => {
        const now = Date.now();
        const elapsed = Math.floor((now - startTime) / 1000); // in seconds
        setElapsedTime(elapsed);
      };

      updateElapsed(); // Initial calculation
      interval = setInterval(updateElapsed, 1000); // Update every second
    } else {
      setElapsedTime(0);
    }

    return () => {
      if (interval) {
        clearInterval(interval);
      }
    };
  }, [currentSchedule]);

  // Format elapsed time as HH:MM:SS
  const formatElapsedTime = (seconds: number): string => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;

    return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  const onRefresh = async () => {
    await Promise.all([
      refetchSchedules(),
      refetchStats(),
    ]);
  };

  const handleClockIn = (schedule: Schedule) => {
    showAlert(
      'Clock In',
      `Clock in for ${schedule.client?.name}?`,
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Clock In',
          onPress: () => {
            // Mock location data - in real app get from device GPS
            const mockLocation = {
              start_latitude: 40.7128,
              start_longitude: -74.0060,
            };

            startVisitMutation.mutate({
              scheduleId: schedule.id,
              data: mockLocation,
            });
          }
        },
      ]
    );
  };

  const handleClockOut = (schedule?: Schedule) => {
    if (!schedule) return;

    // Navigate to the ClockOut screen
    navigation.navigate('ClockOut', { scheduleId: schedule.id });
  };

  const handleViewProgress = (schedule: Schedule) => {
    navigation.navigate('ScheduleDetails', { scheduleId: schedule.id });
  };

  const handleMoreOptions = (schedule: Schedule) => {
    // TODO: Show action sheet with more options
    console.log('More options for schedule:', schedule.id);
  };

  const handleSchedulePress = (schedule: Schedule) => {
    navigation.navigate('ScheduleDetails', { scheduleId: schedule.id });
  };

  const handleSeeAll = () => {
    // TODO: Navigate to all schedules screen
    console.log('See all schedules');
  };

  if (isLoading) {
    return (
      <ContainerView style={styles.container}>
        <View style={styles.loadingContainer}>
          <Text variant="body" color="textSecondary">
            Loading...
          </Text>
        </View>
      </ContainerView>
    );
  }

  return (
    <ContainerView style={styles.container}>
      {/* Custom Header - only show on large screens */}
      {isLargeScreen && (
        <CustomHeader
          userName={currentUser.name}
          userEmail="Admin@healthcare.io"
        />
      )}

      <ScrollView
        style={styles.scrollView}
        contentContainerStyle={styles.content}
        refreshControl={
          <RefreshControl
            refreshing={startVisitMutation.isPending || endVisitMutation.isPending}
            onRefresh={onRefresh}
          />
        }
        showsVerticalScrollIndicator={false}
      >

        {/* Welcome Header */}
        {!isLargeScreen && (
          <Text variant="h2" color="textPrimary" style={styles.welcome}>
            Welcome {currentUser.name}!
          </Text>
        )}

        {isLargeScreen && (
          <Text variant="h2" color="textPrimary">
            Dashboard
          </Text>
        )}

        {/* Current Status Card */}
        {currentSchedule && (
          <StatusCard
            user={currentUser}
            location={currentSchedule.client?.address || 'No address available'}
            timeRange={`${new Date(currentSchedule.start_time).toLocaleTimeString('en-US', {
              hour: '2-digit',
              minute: '2-digit',
              hour12: false
            })} - ${new Date(currentSchedule.end_time).toLocaleTimeString('en-US', {
              hour: '2-digit',
              minute: '2-digit',
              hour12: false
            })}`}
            timer={formatElapsedTime(elapsedTime)}
            onClockOut={() => handleClockOut(currentSchedule)}
            isInProgress={currentSchedule.status === 'in_progress'}
          />
        )}

        {/* Statistics Overview */}
        {stats && <StatsOverview stats={stats} />}

        {/* Schedule List */}
        <ScheduleList
          schedules={schedules}
          onClockIn={handleClockIn}
          onClockOut={handleClockOut}
          onViewProgress={handleViewProgress}
          onMoreOptions={handleMoreOptions}
          onSchedulePress={handleSchedulePress}
          onSeeAll={handleSeeAll}
        />
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
    gap: spacing.xl,
    paddingBottom: spacing.xxxl,
  },
  welcome: {
    marginBottom: spacing.lg,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
});

export default HomeScreen;

import React, { useEffect, useRef } from 'react';
import {
  View,
  StyleSheet,
  Animated,
  TouchableOpacity,
  Dimensions,
} from 'react-native';
import { StackScreenProps } from '@react-navigation/stack';
import { colors, spacing, borderRadius } from '../constants';
import { Text, Button, Icon } from '../components/atoms';
import { HomeStackParamList } from '../navigation/HomeStackNavigator';
import { ContainerView } from '../components/organisms';

type Props = StackScreenProps<HomeStackParamList, 'ScheduleCompleted'>;

const { width, height } = Dimensions.get('window');

const ScheduleCompletedScreen: React.FC<Props> = ({ route, navigation }) => {
  const { schedule, duration } = route.params;

  // Animation refs
  const fadeAnim = useRef(new Animated.Value(0)).current;
  const scaleAnim = useRef(new Animated.Value(0.5)).current;
  const checkmarkAnim = useRef(new Animated.Value(0)).current;
  const floatingElementsAnim = useRef(new Animated.Value(0)).current;

  useEffect(() => {
    // Start animations sequence
    Animated.sequence([
      // Fade in background
      Animated.timing(fadeAnim, {
        toValue: 1,
        duration: 300,
        useNativeDriver: true,
      }),
      // Scale in main content
      Animated.spring(scaleAnim, {
        toValue: 1,
        tension: 50,
        friction: 7,
        useNativeDriver: true,
      }),
      // Animate checkmark
      Animated.timing(checkmarkAnim, {
        toValue: 1,
        duration: 400,
        useNativeDriver: true,
      }),
      // Animate floating elements
      Animated.timing(floatingElementsAnim, {
        toValue: 1,
        duration: 600,
        useNativeDriver: true,
      }),
    ]).start();
  }, []);

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      weekday: 'short',
      day: '2-digit',
      month: 'long',
      year: 'numeric'
    });
  };

  const formatTime = (startTime: string, endTime: string) => {
    const start = new Date(startTime);
    const end = new Date(endTime);

    const startTimeStr = start.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    });

    const endTimeStr = end.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    });

    return `${startTimeStr} - ${endTimeStr} SGT`;
  };

  const handleGoHome = () => {
    navigation.navigate('HomeMain');
  };

  const renderFloatingElement = (style: any, iconName: string, size: number = 8) => (
    <Animated.View
      style={[
        styles.floatingElement,
        style,
        {
          opacity: floatingElementsAnim,
          transform: [
            {
              translateY: floatingElementsAnim.interpolate({
                inputRange: [0, 1],
                outputRange: [20, 0],
              }),
            },
          ],
        },
      ]}
    >
      <View style={[styles.floatingDot, { width: size, height: size }]} />
    </Animated.View>
  );

  return (
    <ContainerView style={styles.container}>
      <Animated.View style={[styles.content, { opacity: fadeAnim }]}>
        {/* Close Button */}
        <TouchableOpacity style={styles.closeButton} onPress={handleGoHome}>
          <Icon name="close" size={24} color={colors.textOnPrimary} />
        </TouchableOpacity>

        {/* Main Content */}
        <Animated.View
          style={[
            styles.mainContent,
            {
              transform: [{ scale: scaleAnim }],
            },
          ]}
        >
          {/* Success Icon with floating elements */}
          <View style={styles.iconContainer}>
            {/* Floating decorative elements */}
            {renderFloatingElement(styles.floatingElement1, 'ellipse', 6)}
            {renderFloatingElement(styles.floatingElement2, 'ellipse', 4)}
            {renderFloatingElement(styles.floatingElement3, 'ellipse', 8)}
            {renderFloatingElement(styles.floatingElement4, 'ellipse', 5)}
            {renderFloatingElement(styles.floatingElement5, 'ellipse', 7)}
            {renderFloatingElement(styles.floatingElement6, 'ellipse', 4)}

            {/* Curved lines */}
            <Animated.View
              style={[
                styles.curvedLine1,
                {
                  opacity: floatingElementsAnim,
                },
              ]}
            />
            <Animated.View
              style={[
                styles.curvedLine2,
                {
                  opacity: floatingElementsAnim,
                },
              ]}
            />

            {/* Main success icon */}
            <Animated.View
              style={[
                styles.successIcon,
                {
                  opacity: checkmarkAnim,
                  transform: [
                    {
                      scale: checkmarkAnim.interpolate({
                        inputRange: [0, 0.5, 1],
                        outputRange: [0, 1.2, 1],
                      }),
                    },
                  ],
                },
              ]}
            >
              <Icon name="checkmark" size={32} color={colors.textOnPrimary} />
            </Animated.View>
          </View>

          {/* Title */}
          <Text variant="h2" color="textOnPrimary" style={styles.title}>
            Schedule Completed
          </Text>

          {/* Schedule Details Card */}
          <View style={styles.detailsCard}>
            <View style={styles.detailRow}>
              <Icon name="calendar-outline" size={20} color={colors.textOnPrimary} />
              <Text variant="body" color="textOnPrimary" style={styles.detailText}>
                {formatDate(schedule.start_time)}
              </Text>
            </View>

            <View style={styles.detailRow}>
              <Icon name="time" size={20} color="white" />
              <Text variant="body" color="textOnPrimary" style={styles.detailText}>
                {formatTime(schedule.start_time, schedule.end_time)}
              </Text>
            </View>
            <View style={styles.lastRow}>
              <Text variant="bodySmall" color="textOnPrimary" style={styles.durationText}>
                ({duration})
              </Text>
            </View>
          </View>
        </Animated.View>

        {/* Go to Home Button */}
        <View style={styles.buttonContainer}>
          <Button
            variant="outline"
            onPress={handleGoHome}
            fullWidth
            style={styles.homeButton}
          >
            <Text variant="button" color="textOnPrimary">
              Go to Home
            </Text>
          </Button>
        </View>
      </Animated.View>
    </ContainerView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.primary,
  },
  content: {
    flex: 1,
    paddingHorizontal: spacing.screenPadding,
  },
  closeButton: {
    alignSelf: 'flex-end',
    marginTop: spacing.lg,
    padding: spacing.sm,
  },
  mainContent: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingVertical: spacing.xxxl,
  },
  iconContainer: {
    position: 'relative',
    width: 120,
    height: 120,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: spacing.xxxl,
  },
  successIcon: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: colors.accent,
    justifyContent: 'center',
    alignItems: 'center',
    zIndex: 10,
  },
  floatingElement: {
    position: 'absolute',
  },
  floatingDot: {
    backgroundColor: colors.accent,
    borderRadius: 50,
  },
  floatingElement1: {
    top: 10,
    left: 20,
  },
  floatingElement2: {
    top: 30,
    right: 15,
  },
  floatingElement3: {
    bottom: 20,
    left: 10,
  },
  floatingElement4: {
    bottom: 35,
    right: 25,
  },
  floatingElement5: {
    top: 50,
    left: 5,
  },
  floatingElement6: {
    bottom: 10,
    right: 5,
  },
  curvedLine1: {
    position: 'absolute',
    top: 15,
    right: 30,
    width: 20,
    height: 2,
    backgroundColor: colors.accent,
    borderRadius: 1,
    transform: [{ rotate: '45deg' }],
  },
  curvedLine2: {
    position: 'absolute',
    bottom: 25,
    left: 15,
    width: 15,
    height: 2,
    backgroundColor: colors.accent,
    borderRadius: 1,
    transform: [{ rotate: '-30deg' }],
  },
  title: {
    textAlign: 'center',
    marginBottom: spacing.xxxl,
    fontWeight: '600',
  },
  detailsCard: {
    backgroundColor: 'rgba(255, 255, 255, 0.1)',
    borderRadius: borderRadius.lg,
    padding: spacing.lg,
    width: '100%',
    maxWidth: 300,
    gap: spacing.md,
  },
  detailRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
  },
  lastRow: {
    flexDirection: 'row',
    paddingLeft: 30,
    marginTop: -10
  },
  detailText: {
    flex: 1,
  },
  durationText: {
    opacity: 0.8,
  },
  buttonContainer: {
    paddingBottom: spacing.xl,
    paddingTop: spacing.xxxl,
  },
  homeButton: {
    borderColor: '#fff',
    borderWidth: 1,
    backgroundColor: 'transparent',
  },
});

export default ScheduleCompletedScreen;

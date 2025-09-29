import React from 'react';
import { createStackNavigator } from '@react-navigation/stack';
import { HomeScreen, ClockOutScreen, ScheduleDetailsScreen } from '../screens';
import ScheduleCompletedScreen from '../screens/ScheduleCompletedScreen';
import { Schedule } from '../services/types';

export type HomeStackParamList = {
  HomeMain: undefined;
  ScheduleDetails: { scheduleId: number };
  ClockOut: { scheduleId: number };
  ScheduleCompleted: { schedule: Schedule; duration: string };
};

const Stack = createStackNavigator<HomeStackParamList>();

const HomeStackNavigator: React.FC = () => {
  return (
    <Stack.Navigator
      screenOptions={{
        headerShown: false,
        cardStyle: { flex: 1 },
      }}
    >
      <Stack.Screen
        name="HomeMain"
        component={HomeScreen}
      />
      <Stack.Screen
        name="ScheduleDetails"
        component={ScheduleDetailsScreen}
        options={{
          headerShown: false,
          gestureEnabled: false, // Prevent swipe back
        }}
      />
      <Stack.Screen
        name="ClockOut"
        component={ClockOutScreen}
        options={{
          headerShown: false,
        }}
      />
      <Stack.Screen
        name="ScheduleCompleted"
        component={ScheduleCompletedScreen}
        options={{
          headerShown: false,
          gestureEnabled: false, // Prevent swipe back
        }}
      />
    </Stack.Navigator>
  );
};

export default HomeStackNavigator;

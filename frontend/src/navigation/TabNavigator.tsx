import React from 'react';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { getFocusedRouteNameFromRoute } from '@react-navigation/native';
import { Ionicons } from '@expo/vector-icons';
import { colors } from '../constants';
import { ProfileScreen } from '../screens';
import HomeStackNavigator from './HomeStackNavigator';
import useScreenSize from '../hooks/useScreenSize';

export type TabParamList = {
  Home: undefined;
  Profile: undefined;
};

const Tab = createBottomTabNavigator<TabParamList>();

const TabNavigator: React.FC = () => {
  const { isLargeScreen } = useScreenSize();

  return (
    <Tab.Navigator
      screenOptions={({ route }) => {
        // Check if the current route in the Home stack is ScheduleDetails, ClockOut, or ScheduleCompleted
        const routeName = getFocusedRouteNameFromRoute(route);
        const tabBarVisible = routeName !== 'ScheduleDetails' && routeName !== 'ClockOut' && routeName !== 'ScheduleCompleted';

        return {
          tabBarIcon: ({ focused, color, size }) => {
            let iconName: keyof typeof Ionicons.glyphMap;

            if (route.name === 'Home') {
              iconName = focused ? 'home' : 'home-outline';
            } else if (route.name === 'Profile') {
              iconName = focused ? 'person' : 'person-outline';
            } else {
              iconName = 'help-outline';
            }

            return <Ionicons name={iconName} size={size} color={color} />;
          },
          tabBarActiveTintColor: colors.primary,
          tabBarInactiveTintColor: colors.gray400,
          tabBarStyle: {
            ...(!tabBarVisible && { display: 'none' }),
            ...(isLargeScreen && { display: 'none' }), // Hide tab bar on large screens
            backgroundColor: colors.white,
            borderTopWidth: 1,
            borderTopColor: colors.gray200,
            paddingBottom: 8,
            paddingTop: 8,
            height: 60,
          },
          tabBarLabelStyle: {
            fontSize: 12,
            fontWeight: '500',
          },
          headerShown: false,
        };
      }}
    >
      <Tab.Screen
        name="Home"
        component={HomeStackNavigator}
        options={{
          tabBarLabel: 'Home',
        }}
      />
      <Tab.Screen
        name="Profile"
        component={ProfileScreen}
        options={{
          tabBarLabel: 'Profile',
        }}
      />
    </Tab.Navigator>
  );
};

export default TabNavigator;

import React from 'react';
import { View } from 'react-native';
import { StatusBar } from 'expo-status-bar';
import { NavigationContainer } from '@react-navigation/native';
import { QueryProvider } from './src/providers/QueryProvider';
import { TabNavigator } from './src/navigation';
import { MutationErrorHandler } from './src/components/error/MutationErrorHandler';

export default function App() {

  return (
    <QueryProvider>
      <View style={{ flex: 1 }}>
        <NavigationContainer>
          <TabNavigator />
          <StatusBar style="auto" />
        </NavigationContainer>
        <MutationErrorHandler />
      </View>
    </QueryProvider>
  );
}

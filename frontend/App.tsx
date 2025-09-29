import React, { useEffect } from 'react';
import { StatusBar } from 'expo-status-bar';
import { NavigationContainer } from '@react-navigation/native';
import { QueryProvider } from './src/providers/QueryProvider';
import { TabNavigator } from './src/navigation';
import * as Font from 'expo-font';
import { useState } from 'react';

export default function App() {

  return (
    <QueryProvider>
      <NavigationContainer>
        <TabNavigator />
        <StatusBar style="auto" />
      </NavigationContainer>
    </QueryProvider>
  );
}

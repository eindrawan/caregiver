import React from 'react';
import { View, StyleSheet, TouchableOpacity } from 'react-native';
import { colors, spacing } from '../constants';
import { Text, Button } from '../components/atoms';
import { ContainerView } from '../components/organisms';

const ProfileScreen: React.FC = () => {
  const handleLogout = () => {
    // Add logout functionality here
    console.log('Log Out pressed');
  };

  return (
    <ContainerView style={styles.container}>
      <View style={styles.content}>
        <Text variant="h2" color="textPrimary" style={styles.welcomeText}>
          Welcome Louis!
        </Text>
        <Button variant="error" outlined={true} fullWidth>
          Log Out
        </Button>
      </View>
    </ContainerView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  content: {
    flex: 1,
    padding: spacing.screenPadding,
  },
  welcomeText: {
    marginBottom: spacing.md,
  },
  logoutButton: {
    borderWidth: 1,
    borderColor: colors.error,
    borderRadius: spacing.sm,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.lg,
    alignItems: 'center',
  },
});

export default ProfileScreen;

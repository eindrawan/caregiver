import React from 'react';
import { View, StyleSheet } from 'react-native';
import { spacing, colors } from '../../constants';
import { Text, Avatar } from '../atoms';

interface UserInfoProps {
  name: string;
  serviceName?: string;
  avatarSource?: { uri: string };
  size?: 'small' | 'medium' | 'large';
  textColor?: keyof typeof colors;
  secondaryTextColor?: keyof typeof colors;
}

const UserInfo: React.FC<UserInfoProps> = ({
  name,
  serviceName,
  avatarSource,
  size = 'medium',
  textColor = 'textPrimary',
  secondaryTextColor = 'textSecondary'
}) => {

  const nameOnlyStyle = serviceName ? {} : { fontSize: 20 };

  return (
    <View style={styles.container}>
      <Avatar
        source={avatarSource}
        name={name}
        size={size}
      />
      <View style={styles.textContainer}>
        <Text variant="title" color={textColor} style={[styles.nameText, nameOnlyStyle]}>
          {name}
        </Text>
        {serviceName && (
          <Text variant="body" color={secondaryTextColor} style={styles.serviceText}>
            {serviceName}
          </Text>
        )}
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  textContainer: {
    marginLeft: spacing.md,
    flex: 1,
    justifyContent: 'center',
    minHeight: 40, // Ensure minimum height to prevent overlap
  },
  nameText: {
    lineHeight: 20,
    marginBottom: 2,
  },
  serviceText: {
    lineHeight: 16,
    marginTop: 2,
  },
});

export default UserInfo;

import React, { useState, useRef } from 'react';
import { View, StyleSheet, TouchableOpacity, Dimensions, Modal, Image } from 'react-native';
import { colors, spacing, borderRadius } from '../../constants';
import { Text, Icon } from '../atoms';
import { UserInfo } from '../molecules';

interface CustomHeaderProps {
  userName?: string;
  userEmail?: string;
  userAvatar?: string;
}

const CustomHeader: React.FC<CustomHeaderProps> = ({
  userName = 'Admin A',
  userEmail = 'Admin@healthcare.io',
  userAvatar
}) => {
  const boxRef = useRef<View>(null);
  const [showDropdown, setShowDropdown] = useState(false);
  const [dropdownPosition, setDropdownPosition] = useState({ x: 0, y: 0 });

  const handleProfilePress = (event: any) => {
    // Get the position of the profile section to position the dropdown
    if (boxRef.current != null) {
      boxRef.current.measureInWindow((x: number, y: number, width: number, height: number) => {
        setDropdownPosition({ x: x + width - 120, y: y + height + 10 }); // Position menu below box
      });
      setShowDropdown(!showDropdown);
    }
  };

  const handleSignOut = () => {
    // Add logout functionality here
    console.log('Sign out pressed');
    setShowDropdown(false);
    // In a real app, this would clear authentication tokens and navigate to login screen
  };

  const handleBackdropPress = () => {
    setShowDropdown(false);
  };

  return (
    <>
      <View style={styles.headerView}>
        <View style={styles.container}>
          {/* Logo Section */}
          <View style={styles.logoSection}>
            <Image
              source={require('../../assets/careviah.png')}
              style={styles.logoImage}
            />
          </View>

          {/* User Profile Section */}
          <TouchableOpacity
            style={styles.profileSection}
            activeOpacity={1}
            onPress={handleProfilePress}
          >
            <View ref={boxRef}>
              <UserInfo
                name={userName}
                serviceName={userEmail}
                size="small"
                textColor="textPrimary"
                secondaryTextColor="textSecondary"
              />
            </View>
            <View style={styles.dropdownIcon}>
              <Icon name="chevron-down" size={16} color="textSecondary" />
            </View>
          </TouchableOpacity>
        </View>
      </View>

      {/* Dropdown Menu */}
      <Modal
        visible={showDropdown}
        transparent={true}
        animationType="none"
        onRequestClose={() => setShowDropdown(false)}
      >
        <TouchableOpacity
          style={styles.backdrop}
          activeOpacity={1}
          onPress={handleBackdropPress}
        >
          <TouchableOpacity
            style={[
              styles.dropdownMenu,
              {
                top: dropdownPosition.y,
                left: dropdownPosition.x,
              }
            ]}
            activeOpacity={1}
          >
            <TouchableOpacity
              style={styles.dropdownItem}
              onPress={handleSignOut}
            >
              <View style={styles.dropdownItemIcon}>
                <Icon name="log-out" size={20} color={colors.error} />
              </View>
              <Text variant="body" color="error">Sign out</Text>
            </TouchableOpacity>
          </TouchableOpacity>
        </TouchableOpacity>
      </Modal>
    </>
  );
};

const styles = StyleSheet.create({
  headerView: {
    padding: spacing.lg,
    zIndex: 1,
  },
  container: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    backgroundColor: colors.accentBackground,
    borderRadius: borderRadius.xl,
  },
  logoSection: {
    flex: 1,
  },
  logoImage: {
    width: 180,
    height: 50,
    resizeMode: 'contain',
  },

  profileSection: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
    paddingVertical: spacing.xs,
    paddingHorizontal: spacing.sm,
    borderRadius: borderRadius.lg,
  },
  userInfo: {
    alignItems: 'flex-end',
  },
  userName: {
    fontWeight: '600',
    lineHeight: 18,
  },
  userEmail: {
    lineHeight: 16,
    marginTop: 1,
  },
  avatar: {
    width: 32,
    height: 32,
    borderRadius: borderRadius.full,
  },
  avatarFallback: {
    width: 32,
    height: 32,
    borderRadius: borderRadius.full,
    backgroundColor: colors.primary,
    justifyContent: 'center',
    alignItems: 'center',
  },
  initialsText: {
    fontSize: 12,
    fontWeight: '600',
    lineHeight: 14,
  },
  dropdownIcon: {
    marginLeft: spacing.xs,
  },
  backdrop: {
    flex: 1,
  },
  dropdownMenu: {
    position: 'absolute',
    backgroundColor: colors.white,
    borderRadius: borderRadius.md,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 5,
    minWidth: 150,
    zIndex: 10,
  },
  dropdownItem: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.md,
    gap: spacing.sm,
  },
  dropdownItemIcon: {
    marginRight: spacing.sm,
  },
});

export default CustomHeader;

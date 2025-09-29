import React from 'react';
import { View, StyleSheet, TouchableOpacity } from 'react-native';
import { StackScreenProps } from '@react-navigation/stack';
import { spacing, colors } from '../../constants';
import { Text, Icon } from '../atoms';

interface HeaderWithBackButtonProps {
    title: string;
    onBackPress: () => void;
}

const HeaderWithBackButton: React.FC<HeaderWithBackButtonProps> = ({
    title,
    onBackPress
}) => {
    return (
        <View style={styles.header}>
            <TouchableOpacity
                style={styles.backButton}
                onPress={onBackPress}
            >
                <Icon name="arrow-back" size={24} color="textPrimary" />
            </TouchableOpacity>
            <Text variant="h2" color="textPrimary" style={styles.headerTitle}>
                {title}
            </Text>
            <View style={styles.headerSpacer} />
        </View>
    );
};

const styles = StyleSheet.create({
    header: {
        flexDirection: 'row',
        alignItems: 'center',
        paddingHorizontal: spacing.screenPadding,
        paddingVertical: spacing.md,
        borderBottomWidth: 1,
        borderBottomColor: colors.gray200,
    },
    backButton: {
        padding: spacing.xs,
    },
    headerTitle: {
        flex: 1,
        textAlign: 'center',
        fontWeight: '600',
    },
    headerSpacer: {
        width: 40, // Same as back button to center the title
    },
});

export default HeaderWithBackButton;
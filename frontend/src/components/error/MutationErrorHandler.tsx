import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import { useErrorStore } from '../../stores/errorStore';
import { colors } from '../../constants/colors';

export const MutationErrorHandler = () => {
    const { error, clearError } = useErrorStore();

    if (!error) return null;

    return (
        <View style={styles.banner}>
            <Text style={styles.text}>{error}</Text>
            <TouchableOpacity style={styles.dismissButton} onPress={clearError}>
                <Text style={styles.dismissText}>Dismiss</Text>
            </TouchableOpacity>
        </View>
    );
};

const styles = StyleSheet.create({
    banner: {
        position: 'absolute',
        bottom: 20,
        left: 16,
        right: 16,
        backgroundColor: colors.error || '#ef4444',
        padding: 16,
        borderRadius: 8,
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 2 },
        shadowOpacity: 0.1,
        shadowRadius: 4,
        elevation: 3,
        zIndex: 9999,
    },
    text: {
        color: '#fff',
        flex: 1,
        marginRight: 8,
    },
    dismissButton: {
        backgroundColor: 'rgba(255, 255, 255, 0.2)',
        paddingHorizontal: 12,
        paddingVertical: 6,
        borderRadius: 4,
    },
    dismissText: {
        color: '#fff',
        fontWeight: '600',
    },
});
import React from 'react';
import { View, ViewStyle, StyleSheet, TouchableOpacity, Linking, Image } from 'react-native';
import { Schedule } from '../../services/types';
import Text from '../atoms/Text';
import { colors, spacing, borderRadius } from '../../constants';

interface MapCardProps {
    schedule: Schedule;
    style?: ViewStyle;
}

export const MapCard: React.FC<MapCardProps> = ({ schedule, style }) => {
    const openMap = () => {
        const latitude = schedule.visit?.start_latitude
        const longitude = schedule.visit?.start_longitude
        if (latitude && longitude) {
            // Use coordinates if available
            const url = `https://www.google.com/maps/search/?api=1&query=${latitude},${longitude}`;
            Linking.openURL(url).catch(() => {
                // Fallback if Google Maps app is not available
                const fallbackUrl = `https://maps.google.com/maps?q=${latitude},${longitude}`;
                Linking.openURL(fallbackUrl);
            });
        } else if (schedule.client?.address) {
            // Use address as fallback
            const address = `${schedule.client.address}, ${schedule.client.city}, ${schedule.client.state} ${schedule.client.zip_code}`;
            const encodedAddress = encodeURIComponent(address);
            const url = `https://www.google.com/maps/search/?api=1&query=${encodedAddress}`;

            Linking.openURL(url).catch(() => {
                // Fallback if Google Maps app is not available
                const fallbackUrl = `https://maps.google.com/maps?q=${encodedAddress}`;
                Linking.openURL(fallbackUrl);
            });
        }
    };

    return (
        <TouchableOpacity style={[styles.locationCard, style]} onPress={openMap}>
            <View style={styles.mapPlaceholder}>
                <Image source={require('../../assets/map-placeholder.png')} style={styles.mapImage} />
            </View>
            <View style={styles.locationInfo}>
                <Text variant="body" color="textPrimary" style={styles.locationAddress}>
                    {schedule.client?.address || 'No address available'}
                </Text>
                <Text variant="body" color="textSecondary" style={styles.locationDetails}>
                    {schedule.client?.city}, {schedule.client?.state}, {schedule.client?.zip_code}
                </Text>
            </View>
        </TouchableOpacity>
    );
};

const styles = StyleSheet.create({
    locationCard: {
        flexDirection: 'row',
        borderRadius: borderRadius.md,
        gap: spacing.md,
    },
    mapPlaceholder: {
        width: 243,
        height: 178,
        backgroundColor: colors.gray100,
        borderRadius: borderRadius.lg,
        alignItems: 'center',
        justifyContent: 'center',
        overflow: 'hidden',
    },
    mapImage: {
        width: '100%',
        height: '100%',
        resizeMode: 'cover',
    },
    locationInfo: {
        flex: 1,
        justifyContent: 'center',
    },
    locationAddress: {
        fontWeight: '600',
        marginBottom: spacing.xs,
    },
    locationDetails: {
        lineHeight: 20,
    },
});
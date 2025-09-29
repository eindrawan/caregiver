import React from 'react';
import { View, Image, StyleSheet, Text } from 'react-native';
import { colors } from '../../constants';

interface Props {
    size?: 'small' | 'medium' | 'large';
    name?: string;
    imageUrl?: string;
    source?: { uri: string };
}

const Avatar: React.FC<Props> = ({ size = 'medium', name, imageUrl, source }) => {
    const sizes = {
        small: 40,
        medium: 48,
        large: 64,
    };

    const avatarSize = sizes[size];
    const initial = name ? name.charAt(0).toUpperCase() : '?';

    const containerStyle = [styles.container, { width: avatarSize, height: avatarSize }];

    // Use source prop first, then imageUrl for backward compatibility
    const imageSource = source || (imageUrl ? { uri: imageUrl } : null);

    if (imageSource) {
        return (
            <View style={containerStyle}>
                <Image source={imageSource} style={styles.image} />
            </View>
        );
    }

    return (
        <View style={[containerStyle, styles.initialContainer]}>
            <Text style={styles.initial}>{initial}</Text>
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        borderRadius: 9999,
        overflow: 'hidden',
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: colors.primary,
    },
    image: {
        width: '100%',
        height: '100%',
    },
    initialContainer: {
        backgroundColor: colors.gray400,
    },
    initial: {
        color: colors.white,
        fontSize: 18,
        fontWeight: 'bold',
    },
});

export default Avatar;
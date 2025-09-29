import React from 'react';
import { View, StyleSheet } from 'react-native';
import { Text } from '../atoms';

const Footer: React.FC = () => {
    return (
        <View style={styles.footerContainer}>
            <Text style={styles.footerText}>@2025 Careviah, Inc. All rights reserved</Text>
        </View>
    );
};

const styles = StyleSheet.create({
    footerContainer: {
        padding: 16,
    },
    footerText: {
        textAlign: 'center',
        fontSize: 12,
    },
});

export default Footer;
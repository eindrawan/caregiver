import React, { useState } from 'react';
import { View, ViewProps, StyleSheet, Dimensions } from 'react-native';
import { SafeAreaView, SafeAreaViewProps } from 'react-native-safe-area-context';
import CustomHeader from './CustomHeader';
import { HeaderWithBackButton } from '../molecules';

interface ContainerViewProps extends SafeAreaViewProps {
    children: React.ReactNode;
    style?: ViewProps['style'];
    maxWidth?: number;
    userName?: string;
    userEmail?: string;
    title?: string;
    onBackPress?: () => void;
}

const ContainerView: React.FC<ContainerViewProps> = ({
    children,
    style,
    maxWidth = 1200,
    userName = 'Admin',
    userEmail = 'Admin@healthcare.io',
    title,
    onBackPress,
    ...props
}) => {

    // Screen width state for responsive design
    const [screenWidth, setScreenWidth] = useState(Dimensions.get('window').width);

    // Check if screen is large (tablet/desktop)
    const isLargeScreen = screenWidth >= 768;

    return (
        <SafeAreaView style={[styles.container, style]} {...props}>
            {/* Custom Header - only show on large screens */}
            {isLargeScreen && (
                <CustomHeader
                    userName={userName}
                    userEmail={userEmail}
                />
            )}
            {/* Header */}
            {title && (
                <HeaderWithBackButton
                    title={title}
                    align={isLargeScreen ? 'left' : 'center'}
                    onBackPress={onBackPress ?? (() => { })}
                />
            )}
            <View style={[styles.innerContainer, { maxWidth }]}>{children}</View>
        </SafeAreaView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        width: '100%',
    },
    innerContainer: {
        flex: 1,
        width: '100%',
        maxWidth: 1200,
        alignSelf: 'center',
    },
});

export default ContainerView;
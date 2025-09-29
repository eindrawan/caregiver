import React from 'react';
import { View, ViewProps, StyleSheet } from 'react-native';
import { SafeAreaView, SafeAreaViewProps } from 'react-native-safe-area-context';

interface ContainerViewProps extends SafeAreaViewProps {
    children: React.ReactNode;
    style?: ViewProps['style'];
    maxWidth?: number;
}

const ContainerView: React.FC<ContainerViewProps> = ({
    children,
    style,
    maxWidth = 1200,
    ...props
}) => {
    return (
        <SafeAreaView style={[styles.container, style]} {...props}>
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
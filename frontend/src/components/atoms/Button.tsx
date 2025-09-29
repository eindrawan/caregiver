import React from 'react';
import { TouchableOpacity, Text, StyleSheet } from 'react-native';
import { colors, borderRadius, spacing } from '../../constants';
import TextComponent from './Text';

interface Props {
    title?: string;
    onPress?: () => void;
    variant?: 'primary' | 'secondary' | 'error';
    disabled?: boolean;
    rounded?: boolean;
    outlined?: boolean;
    fullWidth?: boolean;
    style?: any;
    children?: React.ReactNode;
}

const Button: React.FC<Props> = ({
    title,
    onPress,
    variant = 'primary',
    disabled = false,
    fullWidth = false,
    rounded = false,
    outlined = false,
    style,
    children
}) => {
    const baseStyle = [
        styles.button,
        fullWidth && styles.fullWidth,
        rounded && styles.rounded,
        style
    ];

    const primaryStyle = {
        backgroundColor: colors.primary,
        borderColor: colors.primary,
        color: colors.textOnPrimary,
    };

    const primaryOutlinedStyle = {
        backgroundColor: 'transparent',
        borderColor: colors.accentBackgroundDark,
    };

    const secondaryStyle = {
        backgroundColor: colors.white,
        borderColor: colors.white,
    };

    const secondaryOutlinedStyle = {
        backgroundColor: colors.white,
        borderColor: colors.white,
    };

    const errorStyle = {
        backgroundColor: colors.error,
        borderColor: colors.error,
        color: colors.textOnPrimary,
    };

    const errorOutlinedStyle = {
        backgroundColor: 'transparent',
        borderColor: colors.error,
        color: colors.error,
    };

    const getButtonStyle = () => {
        switch (variant) {
            case 'primary':
                return outlined ? primaryOutlinedStyle : primaryStyle;
            case 'secondary':
                return outlined ? secondaryOutlinedStyle : secondaryStyle;
            case 'error':
                return outlined ? errorOutlinedStyle : errorStyle;
            default:
                return primaryStyle;
        }
    };

    const getTextColor = () => {
        switch (variant) {
            case 'primary':
                return outlined ? 'primary' : 'textOnPrimary';
            case 'secondary':
                return outlined ? 'textOnPrimary' : 'primary';
            case 'error':
                return outlined ? 'error' : 'textOnPrimary';
            default:
                return 'textOnPrimary';
        }
    };

    return (
        <TouchableOpacity
            style={[baseStyle, getButtonStyle(), disabled && styles.disabled, !onPress && styles.disabled]}
            onPress={onPress}
            disabled={disabled || !onPress}
            activeOpacity={0.7}
        >
            {children ? (
                typeof children === 'string' ? (
                    <TextComponent variant="button" color={getTextColor()}>
                        {children}
                    </TextComponent>
                ) : (
                    children
                )
            ) : (
                <TextComponent variant="button" color={getTextColor()}>
                    {title}
                </TextComponent>
            )}
        </TouchableOpacity>
    );
};

const styles = StyleSheet.create({
    button: {
        paddingVertical: spacing.md,
        paddingHorizontal: spacing.lg,
        borderRadius: borderRadius.xl,
        fontWeight: 500,
        borderWidth: 1,
        alignItems: 'center',
        justifyContent: 'center',
        fontFamily: 'Roboto',
    },
    fullWidth: {
        width: '100%',
    },
    rounded: {
        borderRadius: borderRadius.full,
    },
    disabled: {
        opacity: 0.5,
    },
});

export default Button;
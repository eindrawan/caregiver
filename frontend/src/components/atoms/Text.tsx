import React from 'react';
import { Text as RNText, TextProps } from 'react-native';
import { colors, typography } from '../../constants';

interface Props extends TextProps {
    variant?: 'h1' | 'h2' | 'h3' | 'body' | 'bodySmall' | 'button' | 'caption';
    color?: keyof typeof colors;
}

const Text: React.FC<Props> = ({ variant = 'body', color = 'textPrimary', style, children, ...props }) => {
    const variantStyles = {
        h1: {
            fontSize: 32,
            fontWeight: '700' as '700',
            fontFamily: 'Roboto'
        },
        h2: {
            fontSize: 20,
            fontWeight: '600' as '600',
            fontFamily: 'Roboto'
        },
        h3: {
            fontSize: 18,
            fontWeight: '600' as '600',
            fontFamily: 'Roboto'
        },
        title: {
            fontSize: 16,
            fontWeight: '600' as '600',
            fontFamily: 'Roboto'
        },
        subtitle: {
            fontSize: 14,
            fontWeight: '600' as '600',
            fontFamily: 'Roboto'
        },
        body: {
            fontSize: 14,
            fontFamily: 'Roboto'
        },
        bodySmall: {
            fontSize: 14,
            fontWeight: '400' as '400',
            fontFamily: 'Roboto'
        },
        button: {
            fontSize: 14,
            fontWeight: '600' as '600',
            fontFamily: 'Roboto'
        },
        caption: {
            fontSize: 12,
            fontWeight: '400' as '400',
            fontFamily: 'Roboto'
        },
    };

    const textColor = colors[color] || colors.textPrimary;

    return (
        <RNText
            style={[
                { color: textColor },
                variantStyles[variant],
                style,
            ]}
            {...props}
        >
            {children}
        </RNText>
    );
};

export default Text;
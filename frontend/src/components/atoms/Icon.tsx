import React from 'react';
import { Ionicons } from '@expo/vector-icons';
import { colors } from '../../constants';

interface Props {
    name: keyof typeof Ionicons.glyphMap;
    size?: number;
    color?: keyof typeof colors | string;
}

const Icon: React.FC<Props> = ({ name, size = 24, color = 'textPrimary' }) => {
    const colorValue = typeof color === 'string' && color in colors
        ? colors[color as keyof typeof colors]
        : color;

    return (
        <Ionicons
            name={name}
            size={size}
            color={colorValue}
        />
    );
};

export default Icon;
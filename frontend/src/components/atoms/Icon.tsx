import React from 'react';
import { Ionicons } from '@expo/vector-icons';
import { colors } from '../../constants';

interface Props {
    name: keyof typeof Ionicons.glyphMap;
    size?: number;
    color?: keyof typeof colors;
}

const Icon: React.FC<Props> = ({ name, size = 24, color = colors.textPrimary }) => {
    return (
        <Ionicons
            name={name}
            size={size}
            color={color}
        />
    );
};

export default Icon;
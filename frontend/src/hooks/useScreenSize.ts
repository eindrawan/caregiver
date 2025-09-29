import { useState, useEffect } from 'react';
import { Dimensions, ScaledSize } from 'react-native';

const useScreenSize = () => {
  const [screenSize, setScreenSize] = useState({
    width: Dimensions.get('window').width,
    height: Dimensions.get('window').height,
  });

  useEffect(() => {
    const onChange = (result: { window: ScaledSize; screen: ScaledSize }) => {
      setScreenSize({
        width: result.window.width,
        height: result.window.height,
      });
    };

    const subscription = Dimensions.addEventListener('change', onChange);
    return () => subscription?.remove();
  }, []);

  const isLargeScreen = screenSize.width >= 768; // Common breakpoint for tablets and larger screens

  return {
    ...screenSize,
    isLargeScreen,
  };
};

export default useScreenSize;
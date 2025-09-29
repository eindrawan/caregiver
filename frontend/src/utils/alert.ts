import { Alert, Platform } from 'react-native';

export const showAlert = (
  title: string,
  message?: string,
  buttons?: Array<{
    text: string;
    onPress?: () => void;
    style?: 'default' | 'cancel' | 'destructive';
  }>,
) => {
  if (Platform.OS === 'web') {
    const msg = message ? `${title}\n\n${message}` : title;
    if (buttons && buttons.length === 2 && buttons[0].style === 'cancel' && buttons[1].onPress) {
      // Simulate confirm dialog for cancel/action pairs
      const confirmed = window.confirm(msg);
      if (confirmed) {
        buttons[1].onPress();
      } else if (buttons[0].onPress) {
        buttons[0].onPress();
      }
    } else {
      // Fallback to alert for single button or complex cases
      window.alert(msg);
      buttons?.[0]?.onPress?.();
    }
    return;
  }

  Alert.alert(title, message, buttons);
};
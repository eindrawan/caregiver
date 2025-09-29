import * as Location from 'expo-location';
import { Platform } from 'react-native';
import { showAlert } from '../utils/alert';

export interface LocationData {
  latitude: number;
  longitude: number;
}

/**
 * Custom error class for location-related errors
 */
export class LocationError extends Error {
  code: string;

  constructor(code: string, message: string) {
    super(message);
    this.name = 'LocationError';
    this.code = code;
  }
}

/**
 * Request location permissions from the user
 */
export const requestLocationPermission = async (): Promise<boolean> => {
  try {
    const { status } = await Location.requestForegroundPermissionsAsync();
    return status === 'granted';
  } catch (error) {
    console.error('Error requesting location permission:', error);
    return false;
  }
};

/**
 * Check if location permissions are granted
 */
export const checkLocationPermission = async (): Promise<boolean> => {
  try {
    const { status } = await Location.getForegroundPermissionsAsync();
    return status === 'granted';
  } catch (error) {
    console.error('Error checking location permission:', error);
    return false;
  }
};

/**
 * Get current location with high accuracy
 */
export const getCurrentLocation = async (): Promise<LocationData> => {
  try {
    // For web platform, use browser geolocation API as fallback
    if (Platform.OS === 'web') {
      return await getWebLocation();
    }

    // Check if location services are enabled (mobile only)
    const isEnabled = await Location.hasServicesEnabledAsync();
    if (!isEnabled) {
      throw new LocationError('LOCATION_DISABLED', 'Location services are disabled. Please enable location services in your device settings.');
    }

    // Check permissions
    const hasPermission = await checkLocationPermission();
    if (!hasPermission) {
      const granted = await requestLocationPermission();
      if (!granted) {
        throw new LocationError('PERMISSION_DENIED', 'Location permission is required to clock in. Please grant location access in your device settings.');
      }
    }

    // Get current position with high accuracy
    const location = await Location.getCurrentPositionAsync({
      accuracy: Location.Accuracy.High,
      timeInterval: 5000,
      distanceInterval: 1,
    });

    return {
      latitude: location.coords.latitude,
      longitude: location.coords.longitude,
    };
  } catch (error: any) {
    console.error('Error getting current location:', error);

    // Handle specific error types
    if (error instanceof LocationError) {
      throw error;
    }

    // Handle expo-location specific errors
    if (error.code === 'E_LOCATION_SERVICES_DISABLED') {
      throw new LocationError('LOCATION_DISABLED', 'Location services are disabled. Please enable location services in your device settings.');
    }

    if (error.code === 'E_LOCATION_UNAVAILABLE') {
      throw new LocationError('LOCATION_UNAVAILABLE', 'Unable to determine your location. Please try again.');
    }

    // Generic error
    throw new LocationError('UNKNOWN_ERROR', 'Failed to get your current location. Please try again.');
  }
};

/**
 * Get location using browser's geolocation API (web fallback)
 */
const getWebLocation = async (): Promise<LocationData> => {
  return new Promise((resolve, reject) => {
    if (!navigator.geolocation) {
      reject(new LocationError('NOT_SUPPORTED', 'Geolocation is not supported by this browser.'));
      return;
    }

    navigator.geolocation.getCurrentPosition(
      (position) => {
        resolve({
          latitude: position.coords.latitude,
          longitude: position.coords.longitude,
        });
      },
      (error) => {
        console.error('Web geolocation error:', error);

        switch (error.code) {
          case error.PERMISSION_DENIED:
            reject(new LocationError('PERMISSION_DENIED', 'Location access denied. Please allow location access in your browser.'));
            break;
          case error.POSITION_UNAVAILABLE:
            reject(new LocationError('LOCATION_UNAVAILABLE', 'Location information is unavailable.'));
            break;
          case error.TIMEOUT:
            reject(new LocationError('TIMEOUT', 'Location request timed out. Please try again.'));
            break;
          default:
            reject(new LocationError('UNKNOWN_ERROR', 'An unknown error occurred while retrieving location.'));
            break;
        }
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 60000,
      }
    );
  });
};

/**
 * Show location error alert to user
 */
export const showLocationErrorAlert = (error: LocationError, onRetry?: () => void, onCancel?: () => void) => {
  const buttons: any[] = [
    {
      text: 'Cancel',
      style: 'cancel',
      onPress: onCancel,
    },
  ];

  if (onRetry) {
    buttons.push({
      text: 'Retry',
      onPress: onRetry,
    });
  }

  if (error.code === 'PERMISSION_DENIED' || error.code === 'LOCATION_DISABLED') {
    buttons.push({
      text: 'Settings',
      onPress: () => {
        // On mobile, this would open device settings
        // For now, just show an alert
        showAlert(
          'Open Settings',
          'Please go to your device settings and enable location permissions for this app.',
          [{ text: 'OK' }]
        );
      },
    });
  }

  showAlert(
    'Location Error',
    error.message,
    buttons
  );
};

/**
 * Geocode an address using OpenStreetMap Nominatim (no API key required)
 */
export const geocodeAddress = async (address: string): Promise<LocationData> => {
  const encodedAddress = encodeURIComponent(address);
  const url = `https://nominatim.openstreetmap.org/search?format=json&q=${encodedAddress}&limit=1&addressdetails=1`;

  try {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();

    if (data.length > 0) {
      return {
        latitude: parseFloat(data[0].lat),
        longitude: parseFloat(data[0].lon),
      };
    } else {
      throw new Error('Address not found');
    }
  } catch (error) {
    console.error('Geocoding error:', error);
    throw new LocationError('GEOCODING_FAILED', 'Could not find location for the address provided. Please check the address and try again.');
  }
};


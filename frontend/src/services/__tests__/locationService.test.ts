import { getCurrentLocation, requestLocationPermission, checkLocationPermission, LocationError } from '../locationService';

// Mock expo-location
jest.mock('expo-location', () => ({
  requestForegroundPermissionsAsync: jest.fn(),
  getForegroundPermissionsAsync: jest.fn(),
  hasServicesEnabledAsync: jest.fn(),
  getCurrentPositionAsync: jest.fn(),
  Accuracy: {
    High: 'high',
  },
}));

// Mock react-native Platform
jest.mock('react-native', () => ({
  Platform: {
    OS: 'ios',
  },
  Alert: {
    alert: jest.fn(),
  },
}));

import * as Location from 'expo-location';

const mockLocation = Location as jest.Mocked<typeof Location>;

describe('LocationService', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('requestLocationPermission', () => {
    it('should return true when permission is granted', async () => {
      mockLocation.requestForegroundPermissionsAsync.mockResolvedValue({
        status: 'granted',
        granted: true,
        canAskAgain: true,
        expires: 'never',
      });

      const result = await requestLocationPermission();
      expect(result).toBe(true);
    });

    it('should return false when permission is denied', async () => {
      mockLocation.requestForegroundPermissionsAsync.mockResolvedValue({
        status: 'denied',
        granted: false,
        canAskAgain: true,
        expires: 'never',
      });

      const result = await requestLocationPermission();
      expect(result).toBe(false);
    });

    it('should handle errors gracefully', async () => {
      mockLocation.requestForegroundPermissionsAsync.mockRejectedValue(new Error('Permission error'));

      const result = await requestLocationPermission();
      expect(result).toBe(false);
    });
  });

  describe('checkLocationPermission', () => {
    it('should return true when permission is already granted', async () => {
      mockLocation.getForegroundPermissionsAsync.mockResolvedValue({
        status: 'granted',
        granted: true,
        canAskAgain: true,
        expires: 'never',
      });

      const result = await checkLocationPermission();
      expect(result).toBe(true);
    });

    it('should return false when permission is not granted', async () => {
      mockLocation.getForegroundPermissionsAsync.mockResolvedValue({
        status: 'denied',
        granted: false,
        canAskAgain: true,
        expires: 'never',
      });

      const result = await checkLocationPermission();
      expect(result).toBe(false);
    });
  });

  describe('getCurrentLocation', () => {
    it('should return location data when successful', async () => {
      mockLocation.hasServicesEnabledAsync.mockResolvedValue(true);
      mockLocation.getForegroundPermissionsAsync.mockResolvedValue({
        status: 'granted',
        granted: true,
        canAskAgain: true,
        expires: 'never',
      });
      mockLocation.getCurrentPositionAsync.mockResolvedValue({
        coords: {
          latitude: 40.7128,
          longitude: -74.0060,
          altitude: null,
          accuracy: 5,
          altitudeAccuracy: null,
          heading: null,
          speed: null,
        },
        timestamp: Date.now(),
      });

      const result = await getCurrentLocation();
      expect(result).toEqual({
        latitude: 40.7128,
        longitude: -74.0060,
      });
    });

    it('should throw LocationError when location services are disabled', async () => {
      mockLocation.hasServicesEnabledAsync.mockResolvedValue(false);

      await expect(getCurrentLocation()).rejects.toThrow(LocationError);
      await expect(getCurrentLocation()).rejects.toThrow('Location services are disabled');
    });

    it('should request permission when not granted and succeed', async () => {
      mockLocation.hasServicesEnabledAsync.mockResolvedValue(true);
      mockLocation.getForegroundPermissionsAsync.mockResolvedValue({
        status: 'denied',
        granted: false,
        canAskAgain: true,
        expires: 'never',
      });
      mockLocation.requestForegroundPermissionsAsync.mockResolvedValue({
        status: 'granted',
        granted: true,
        canAskAgain: true,
        expires: 'never',
      });
      mockLocation.getCurrentPositionAsync.mockResolvedValue({
        coords: {
          latitude: 40.7128,
          longitude: -74.0060,
          altitude: null,
          accuracy: 5,
          altitudeAccuracy: null,
          heading: null,
          speed: null,
        },
        timestamp: Date.now(),
      });

      const result = await getCurrentLocation();
      expect(result).toEqual({
        latitude: 40.7128,
        longitude: -74.0060,
      });
      expect(mockLocation.requestForegroundPermissionsAsync).toHaveBeenCalled();
    });

    it('should throw LocationError when permission is denied', async () => {
      mockLocation.hasServicesEnabledAsync.mockResolvedValue(true);
      mockLocation.getForegroundPermissionsAsync.mockResolvedValue({
        status: 'denied',
        granted: false,
        canAskAgain: true,
        expires: 'never',
      });
      mockLocation.requestForegroundPermissionsAsync.mockResolvedValue({
        status: 'denied',
        granted: false,
        canAskAgain: true,
        expires: 'never',
      });

      await expect(getCurrentLocation()).rejects.toThrow(LocationError);
      await expect(getCurrentLocation()).rejects.toThrow('Location permission is required');
    });
  });
});

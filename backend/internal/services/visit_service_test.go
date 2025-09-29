package services

import (
	"caregiver-shift-tracker/internal/models"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVisitService_GetVisitByScheduleID(t *testing.T) {
	// Setup
	mockVisitRepo := new(MockVisitRepository)
	logger := logrus.New()
	service := NewVisitService(mockVisitRepo, logger)

	// Test data
	now := time.Now()
	expectedVisit := &models.Visit{
		ID:         1,
		ScheduleID: 1,
		StartTime:  &now,
		Status:     "in_progress",
	}

	// Mock expectations
	mockVisitRepo.On("GetByScheduleID", 1).Return(expectedVisit, nil)

	// Execute
	result, err := service.GetVisitByScheduleID(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ScheduleID)
	assert.Equal(t, "in_progress", result.Status)

	// Verify mock expectations
	mockVisitRepo.AssertExpectations(t)
}

func TestVisitService_GetVisitByScheduleID_NotFound(t *testing.T) {
	// Setup
	mockVisitRepo := new(MockVisitRepository)
	logger := logrus.New()
	service := NewVisitService(mockVisitRepo, logger)

	// Mock expectations
	mockVisitRepo.On("GetByScheduleID", 999).Return(nil, nil)

	// Execute
	result, err := service.GetVisitByScheduleID(999)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)

	// Verify mock expectations
	mockVisitRepo.AssertExpectations(t)
}

func TestVisitService_CreateVisit(t *testing.T) {
	// Setup
	mockVisitRepo := new(MockVisitRepository)
	logger := logrus.New()
	service := NewVisitService(mockVisitRepo, logger)

	// Test data
	visit := &models.Visit{
		ScheduleID: 1,
		Status:     "not_started",
	}

	// Mock expectations
	mockVisitRepo.On("Create", visit).Return(nil)

	// Execute
	err := service.CreateVisit(visit)

	// Assert
	assert.NoError(t, err)

	// Verify mock expectations
	mockVisitRepo.AssertExpectations(t)
}

func TestVisitService_CreateVisit_ValidationError_MissingScheduleID(t *testing.T) {
	// Setup
	mockVisitRepo := new(MockVisitRepository)
	logger := logrus.New()
	service := NewVisitService(mockVisitRepo, logger)

	// Test data with missing schedule ID
	visit := &models.Visit{
		Status: "not_started",
	}

	// Execute
	err := service.CreateVisit(visit)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "schedule_id is required")
}

func TestVisitService_CreateVisit_ValidationError_InvalidStatus(t *testing.T) {
	// Setup
	mockVisitRepo := new(MockVisitRepository)
	logger := logrus.New()
	service := NewVisitService(mockVisitRepo, logger)

	// Test data with invalid status
	visit := &models.Visit{
		ScheduleID: 1,
		Status:     "invalid_status",
	}

	// Execute
	err := service.CreateVisit(visit)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status")
}

func TestVisitService_CreateVisit_ValidationError_InProgressMissingStartTime(t *testing.T) {
	// Setup
	mockVisitRepo := new(MockVisitRepository)
	logger := logrus.New()
	service := NewVisitService(mockVisitRepo, logger)

	// Test data with in_progress status but missing start time
	visit := &models.Visit{
		ScheduleID: 1,
		Status:     "in_progress",
	}

	// Execute
	err := service.CreateVisit(visit)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "start_time is required when status is in_progress")
}

func TestVisitService_CreateVisit_ValidationError_CompletedMissingEndTime(t *testing.T) {
	// Setup
	mockVisitRepo := new(MockVisitRepository)
	logger := logrus.New()
	service := NewVisitService(mockVisitRepo, logger)

	// Test data with completed status but missing end time
	now := time.Now()
	lat := 40.7128
	lng := -74.0060
	visit := &models.Visit{
		ScheduleID:     1,
		Status:         "completed",
		StartTime:      &now,
		StartLatitude:  &lat,
		StartLongitude: &lng,
		// Missing end time
	}

	// Execute
	err := service.CreateVisit(visit)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "end_time is required when status is completed")
}

func TestVisitService_CreateVisit_ValidationError_InvalidLatitude(t *testing.T) {
	// Setup
	mockVisitRepo := new(MockVisitRepository)
	logger := logrus.New()
	service := NewVisitService(mockVisitRepo, logger)

	// Test data with invalid latitude
	now := time.Now()
	invalidLat := 100.0 // Invalid latitude (> 90)
	lng := -74.0060
	visit := &models.Visit{
		ScheduleID:     1,
		Status:         "in_progress",
		StartTime:      &now,
		StartLatitude:  &invalidLat,
		StartLongitude: &lng,
	}

	// Execute
	err := service.CreateVisit(visit)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "start_latitude must be between -90 and 90")
}

func TestVisitService_CreateVisit_ValidationError_InvalidLongitude(t *testing.T) {
	// Setup
	mockVisitRepo := new(MockVisitRepository)
	logger := logrus.New()
	service := NewVisitService(mockVisitRepo, logger)

	// Test data with invalid longitude
	now := time.Now()
	lat := 40.7128
	invalidLng := 200.0 // Invalid longitude (> 180)
	visit := &models.Visit{
		ScheduleID:     1,
		Status:         "in_progress",
		StartTime:      &now,
		StartLatitude:  &lat,
		StartLongitude: &invalidLng,
	}

	// Execute
	err := service.CreateVisit(visit)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "start_longitude must be between -180 and 180")
}

func TestVisitService_CreateVisit_ValidationError_EndTimeBeforeStartTime(t *testing.T) {
	// Setup
	mockVisitRepo := new(MockVisitRepository)
	logger := logrus.New()
	service := NewVisitService(mockVisitRepo, logger)

	// Test data with end time before start time
	startTime := time.Now()
	endTime := startTime.Add(-time.Hour) // End time is before start time
	lat := 40.7128
	lng := -74.0060
	visit := &models.Visit{
		ScheduleID:     1,
		Status:         "completed",
		StartTime:      &startTime,
		EndTime:        &endTime,
		StartLatitude:  &lat,
		StartLongitude: &lng,
		EndLatitude:    &lat,
		EndLongitude:   &lng,
	}

	// Execute
	err := service.CreateVisit(visit)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "end_time cannot be before start_time")
}

func TestVisitService_UpdateVisit(t *testing.T) {
	// Setup
	mockVisitRepo := new(MockVisitRepository)
	logger := logrus.New()
	service := NewVisitService(mockVisitRepo, logger)

	// Test data
	now := time.Now()
	lat := 40.7128
	lng := -74.0060
	visit := &models.Visit{
		ID:             1,
		ScheduleID:     1,
		Status:         "completed",
		StartTime:      &now,
		EndTime:        &now,
		StartLatitude:  &lat,
		StartLongitude: &lng,
		EndLatitude:    &lat,
		EndLongitude:   &lng,
	}

	// Mock expectations
	mockVisitRepo.On("Update", visit).Return(nil)

	// Execute
	err := service.UpdateVisit(visit)

	// Assert
	assert.NoError(t, err)

	// Verify mock expectations
	mockVisitRepo.AssertExpectations(t)
}

func TestVisitService_EndVisit(t *testing.T) {
	// Setup
	mockVisitRepo := new(MockVisitRepository)
	logger := logrus.New()
	service := NewVisitService(mockVisitRepo, logger)

	// Test data
	now := time.Now()
	lat := 40.7128
	lng := -74.0060
	initialVisit := &models.Visit{
		ID:             1,
		ScheduleID:     1,
		StartTime:      &now,
		StartLatitude:  &lat,
		StartLongitude: &lng,
		Status:         "in_progress",
		LocationStatus: "confirmed",
	}

	// Mock expectations
	mockVisitRepo.On("GetByScheduleID", 1).Return(initialVisit, nil)
	mockVisitRepo.On("Update", mock.AnythingOfType("*models.Visit")).Return(nil).Run(func(args mock.Arguments) {
		updatedVisit := args.Get(0).(*models.Visit)
		updatedVisit.EndTime = &now
		updatedVisit.EndLatitude = &lat
		updatedVisit.EndLongitude = &lng
		updatedVisit.Status = "completed"
		updatedVisit.LocationStatus = "confirmed"
		updatedVisit.Notes = "Test notes"
	})

	// Execute
	err := service.EndVisit(1, lat, lng, "Test notes")

	// Assert
	assert.NoError(t, err)

	// Verify
	updatedCall := mockVisitRepo.Calls[1]
	assert.Len(t, updatedCall.Arguments, 1)
	updatedVisit := updatedCall.Arguments[0].(*models.Visit)
	assert.Equal(t, "completed", updatedVisit.Status)
	assert.Equal(t, "confirmed", updatedVisit.LocationStatus)
	assert.NotNil(t, updatedVisit.EndTime)
	assert.Equal(t, "Test notes", updatedVisit.Notes)

	mockVisitRepo.AssertExpectations(t)
}

func TestVisitService_CreateVisit_LocationStatusPending(t *testing.T) {
	// Setup
	mockVisitRepo := new(MockVisitRepository)
	logger := logrus.New()
	service := NewVisitService(mockVisitRepo, logger)

	// Test data - no location
	visit := &models.Visit{
		ScheduleID: 1,
		Status:     "not_started",
	}

	// Mock expectations
	mockVisitRepo.On("Create", mock.MatchedBy(func(v *models.Visit) bool {
		return v.LocationStatus == "pending"
	})).Return(nil)

	// Execute
	err := service.CreateVisit(visit)

	// Assert
	assert.NoError(t, err)
	mockVisitRepo.AssertExpectations(t)
}

func TestVisitService_CreateVisit_LocationStatusConfirmed(t *testing.T) {
	// Setup
	mockVisitRepo := new(MockVisitRepository)
	logger := logrus.New()
	service := NewVisitService(mockVisitRepo, logger)

	// Test data - with location
	now := time.Now()
	lat := 40.7128
	lng := -74.0060
	visit := &models.Visit{
		ScheduleID:     1,
		Status:         "in_progress",
		StartTime:      &now,
		StartLatitude:  &lat,
		StartLongitude: &lng,
	}

	// Mock expectations
	mockVisitRepo.On("Create", mock.MatchedBy(func(v *models.Visit) bool {
		return v.LocationStatus == "confirmed"
	})).Return(nil)

	// Execute
	err := service.CreateVisit(visit)

	// Assert
	assert.NoError(t, err)
	mockVisitRepo.AssertExpectations(t)
}

func TestVisitService_ValidateVisit_InvalidLocationStatus(t *testing.T) {
	// Setup
	mockVisitRepo := new(MockVisitRepository)
	logger := logrus.New()
	service := NewVisitService(mockVisitRepo, logger)

	// Test data
	visit := &models.Visit{
		ScheduleID:     1,
		Status:         "not_started",
		LocationStatus: "invalid",
	}

	// Execute
	err := service.CreateVisit(visit)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid location_status")

	mockVisitRepo.AssertExpectations(t)
}

func TestVisitService_EndVisit_NotStarted(t *testing.T) {
	// Setup
	mockVisitRepo := new(MockVisitRepository)
	logger := logrus.New()
	service := NewVisitService(mockVisitRepo, logger)

	// Test data - no start time
	visit := &models.Visit{
		ID:         1,
		ScheduleID: 1,
		Status:     "not_started",
	}

	// Mock expectations
	mockVisitRepo.On("GetByScheduleID", 1).Return(visit, nil)

	// Execute
	err := service.EndVisit(1, 40.7128, -74.0060, "")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot end visit that hasn't been started")

	mockVisitRepo.AssertExpectations(t)
}

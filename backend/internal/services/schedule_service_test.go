package services

import (
	"caregiver-shift-tracker/internal/models"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockScheduleRepository is a mock implementation of ScheduleRepository
type MockScheduleRepository struct {
	mock.Mock
}

func (m *MockScheduleRepository) GetAll(filter *models.ScheduleFilter) ([]models.Schedule, error) {
	args := m.Called(filter)
	return args.Get(0).([]models.Schedule), args.Error(1)
}

func (m *MockScheduleRepository) GetByID(id int) (*models.Schedule, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Schedule), args.Error(1)
}

func (m *MockScheduleRepository) GetToday(caregiverID int) ([]models.Schedule, error) {
	args := m.Called(caregiverID)
	return args.Get(0).([]models.Schedule), args.Error(1)
}

func (m *MockScheduleRepository) GetStats(caregiverID int) (*models.ScheduleStats, error) {
	args := m.Called(caregiverID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ScheduleStats), args.Error(1)
}

func (m *MockScheduleRepository) Create(schedule *models.Schedule) error {
	args := m.Called(schedule)
	return args.Error(0)
}

func (m *MockScheduleRepository) Update(schedule *models.Schedule) error {
	args := m.Called(schedule)
	return args.Error(0)
}

func (m *MockScheduleRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockVisitRepository is a mock implementation of VisitRepository
type MockVisitRepository struct {
	mock.Mock
}

func (m *MockVisitRepository) GetByScheduleID(scheduleID int) (*models.Visit, error) {
	args := m.Called(scheduleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Visit), args.Error(1)
}

func (m *MockVisitRepository) Create(visit *models.Visit) error {
	args := m.Called(visit)
	return args.Error(0)
}

func (m *MockVisitRepository) Update(visit *models.Visit) error {
	args := m.Called(visit)
	return args.Error(0)
}

func (m *MockVisitRepository) StartVisit(scheduleID int, latitude, longitude float64) error {
	args := m.Called(scheduleID, latitude, longitude)
	return args.Error(0)
}

func (m *MockVisitRepository) EndVisit(scheduleID int, latitude, longitude float64, notes string) error {
	args := m.Called(scheduleID, latitude, longitude, notes)
	return args.Error(0)
}

func (m *MockVisitRepository) CancelVisit(scheduleID int) error {
	args := m.Called(scheduleID)
	return args.Error(0)
}

// MockTaskRepository is a mock implementation of TaskRepository
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) GetByScheduleID(scheduleID int) ([]models.Task, error) {
	args := m.Called(scheduleID)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskRepository) GetByID(id int) (*models.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) Create(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) Update(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) UpdateStatus(id int, status, reason string) error {
	args := m.Called(id, status, reason)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestScheduleService_GetAllSchedules(t *testing.T) {
	// Setup
	mockScheduleRepo := new(MockScheduleRepository)
	mockVisitRepo := new(MockVisitRepository)
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewScheduleService(mockScheduleRepo, mockVisitRepo, mockTaskRepo, logger)

	// Test data
	expectedSchedules := []models.Schedule{
		{
			ID:       1,
			ClientID: 1,
			Status:   "scheduled",
			Client: &models.Client{
				ID:   1,
				Name: "John Doe",
			},
		},
		{
			ID:       2,
			ClientID: 2,
			Status:   "completed",
			Client: &models.Client{
				ID:   2,
				Name: "Jane Smith",
			},
		},
	}

	filter := &models.ScheduleFilter{CaregiverID: &[]int{1}[0]}

	// Mock expectations
	mockScheduleRepo.On("GetAll", filter).Return(expectedSchedules, nil)
	mockVisitRepo.On("GetByScheduleID", 1).Return(nil, nil)
	mockTaskRepo.On("GetByScheduleID", 1).Return([]models.Task{}, nil)
	mockVisitRepo.On("GetByScheduleID", 2).Return(nil, nil)
	mockTaskRepo.On("GetByScheduleID", 2).Return([]models.Task{}, nil)

	// Execute
	result, err := service.GetAllSchedules(filter)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "John Doe", result[0].Client.Name)
	assert.Equal(t, "Jane Smith", result[1].Client.Name)

	// Verify mock expectations
	mockScheduleRepo.AssertExpectations(t)
	mockVisitRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestScheduleService_GetScheduleByID(t *testing.T) {
	// Setup
	mockScheduleRepo := new(MockScheduleRepository)
	mockVisitRepo := new(MockVisitRepository)
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewScheduleService(mockScheduleRepo, mockVisitRepo, mockTaskRepo, logger)

	// Test data
	expectedSchedule := &models.Schedule{
		ID:       1,
		ClientID: 1,
		Status:   "scheduled",
		Client: &models.Client{
			ID:   1,
			Name: "John Doe",
		},
	}

	// Mock expectations
	mockScheduleRepo.On("GetByID", 1).Return(expectedSchedule, nil)
	mockVisitRepo.On("GetByScheduleID", 1).Return(nil, nil)
	mockTaskRepo.On("GetByScheduleID", 1).Return([]models.Task{}, nil)

	// Execute
	result, err := service.GetScheduleByID(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "John Doe", result.Client.Name)

	// Verify mock expectations
	mockScheduleRepo.AssertExpectations(t)
	mockVisitRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestScheduleService_GetScheduleByID_NotFound(t *testing.T) {
	// Setup
	mockScheduleRepo := new(MockScheduleRepository)
	mockVisitRepo := new(MockVisitRepository)
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewScheduleService(mockScheduleRepo, mockVisitRepo, mockTaskRepo, logger)

	// Mock expectations
	mockScheduleRepo.On("GetByID", 999).Return(nil, nil)

	// Execute
	result, err := service.GetScheduleByID(999)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)

	// Verify mock expectations
	mockScheduleRepo.AssertExpectations(t)
}

func TestScheduleService_StartVisit(t *testing.T) {
	// Setup
	mockScheduleRepo := new(MockScheduleRepository)
	mockVisitRepo := new(MockVisitRepository)
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewScheduleService(mockScheduleRepo, mockVisitRepo, mockTaskRepo, logger)

	// Test data
	schedule := &models.Schedule{
		ID:        1,
		StartTime: time.Now().Add(15 * time.Minute), // Within 30 minutes
		Status:    "scheduled",
	}
	req := &models.VisitStartRequest{
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	// Mock expectations
	mockScheduleRepo.On("GetByID", 1).Return(schedule, nil)
	mockVisitRepo.On("StartVisit", 1, req.Latitude, req.Longitude).Return(nil)
	mockScheduleRepo.On("Update", mock.AnythingOfType("*models.Schedule")).Return(nil)

	// Execute
	err := service.StartVisit(1, req)

	// Assert
	assert.NoError(t, err)

	// Verify mock expectations
	mockScheduleRepo.AssertExpectations(t)
	mockVisitRepo.AssertExpectations(t)
}

func TestScheduleService_StartVisit_ScheduleNotFound(t *testing.T) {
	// Setup
	mockScheduleRepo := new(MockScheduleRepository)
	mockVisitRepo := new(MockVisitRepository)
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewScheduleService(mockScheduleRepo, mockVisitRepo, mockTaskRepo, logger)

	req := &models.VisitStartRequest{
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	// Mock expectations
	mockScheduleRepo.On("GetByID", 999).Return(nil, nil)

	// Execute
	err := service.StartVisit(999, req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "schedule not found")

	// Verify mock expectations
	mockScheduleRepo.AssertExpectations(t)
}

func TestScheduleService_StartVisit_TooEarly(t *testing.T) {
	// Setup
	mockScheduleRepo := new(MockScheduleRepository)
	mockVisitRepo := new(MockVisitRepository)
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewScheduleService(mockScheduleRepo, mockVisitRepo, mockTaskRepo, logger)

	// Test data - schedule starts in 2 hours (more than 30 minutes early)
	schedule := &models.Schedule{
		ID:        1,
		StartTime: time.Now().Add(2 * time.Hour),
		Status:    "scheduled",
	}
	req := &models.VisitStartRequest{
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	// Mock expectations
	mockScheduleRepo.On("GetByID", 1).Return(schedule, nil)

	// Execute
	err := service.StartVisit(1, req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot start visit more than 30 minutes before")

	// Verify mock expectations
	mockScheduleRepo.AssertExpectations(t)
}

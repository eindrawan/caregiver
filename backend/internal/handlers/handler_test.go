package handlers

import (
	"bytes"
	"caregiver-shift-tracker/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockScheduleService is a mock implementation of ScheduleService
type MockScheduleService struct {
	mock.Mock
}

func (m *MockScheduleService) GetAllSchedules(filter *models.ScheduleFilter) ([]models.Schedule, error) {
	args := m.Called(filter)
	return args.Get(0).([]models.Schedule), args.Error(1)
}

func (m *MockScheduleService) GetScheduleByID(id int) (*models.Schedule, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Schedule), args.Error(1)
}

func (m *MockScheduleService) GetTodaySchedules(caregiverID int) ([]models.Schedule, error) {
	args := m.Called(caregiverID)
	return args.Get(0).([]models.Schedule), args.Error(1)
}

func (m *MockScheduleService) GetScheduleStats(caregiverID int) (*models.ScheduleStats, error) {
	args := m.Called(caregiverID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ScheduleStats), args.Error(1)
}

func (m *MockScheduleService) StartVisit(scheduleID int, req *models.VisitStartRequest) error {
	args := m.Called(scheduleID, req)
	return args.Error(0)
}

func (m *MockScheduleService) EndVisit(scheduleID int, req *models.VisitEndRequest) error {
	args := m.Called(scheduleID, req)
	return args.Error(0)
}

func (m *MockScheduleService) CancelVisit(scheduleID int) error {
	args := m.Called(scheduleID)
	return args.Error(0)
}

// MockVisitService is a mock implementation of VisitService
type MockVisitService struct {
	mock.Mock
}

func (m *MockVisitService) GetVisitByScheduleID(scheduleID int) (*models.Visit, error) {
	args := m.Called(scheduleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Visit), args.Error(1)
}

// MockTaskService is a mock implementation of TaskService
type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) GetTaskByID(id int) (*models.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskService) UpdateTaskStatus(id int, req *models.TaskUpdateRequest) (*models.Task, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

// MockClientService is a mock implementation of ClientService
type MockClientService struct {
	mock.Mock
}

func (m *MockClientService) GetAllClients(filter *models.ClientFilter) ([]models.Client, error) {
	args := m.Called(filter)
	return args.Get(0).([]models.Client), args.Error(1)
}

func (m *MockClientService) GetClientByID(id int) (*models.Client, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockClientService) CreateClient(req *models.ClientCreateRequest) (*models.Client, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockClientService) UpdateClient(id int, req *models.ClientUpdateRequest) (*models.Client, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockClientService) DeleteClient(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockClientService) SearchClients(query string) ([]models.Client, error) {
	args := m.Called(query)
	return args.Get(0).([]models.Client), args.Error(1)
}

func setupTestHandler() (*Handler, *MockScheduleService, *MockVisitService, *MockTaskService, *MockClientService) {
	gin.SetMode(gin.TestMode)

	mockScheduleService := new(MockScheduleService)
	mockVisitService := new(MockVisitService)
	mockTaskService := new(MockTaskService)
	mockClientService := new(MockClientService)
	logger := logrus.New()

	handler := NewHandler(mockScheduleService, mockVisitService, mockTaskService, mockClientService, logger)

	return handler, mockScheduleService, mockVisitService, mockTaskService, mockClientService
}

func TestHandler_HealthCheck(t *testing.T) {
	// Setup
	handler, _, _, _, _ := setupTestHandler()
	router := handler.SetupRoutes()

	// Create request
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "caregiver-shift-tracker", response["service"])
}

func TestHandler_GetScheduleByID(t *testing.T) {
	// Setup
	handler, mockScheduleService, _, _, _ := setupTestHandler()
	router := handler.SetupRoutes()

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
	mockScheduleService.On("GetScheduleByID", 1).Return(expectedSchedule, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/schedules/1", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["id"])

	client := data["client"].(map[string]interface{})
	assert.Equal(t, "John Doe", client["name"])

	// Verify mock expectations
	mockScheduleService.AssertExpectations(t)
}

func TestHandler_GetScheduleByID_NotFound(t *testing.T) {
	// Setup
	handler, mockScheduleService, _, _, _ := setupTestHandler()
	router := handler.SetupRoutes()

	// Mock expectations
	mockScheduleService.On("GetScheduleByID", 999).Return(nil, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/schedules/999", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Schedule not found", response["error"])

	// Verify mock expectations
	mockScheduleService.AssertExpectations(t)
}

func TestHandler_StartVisit(t *testing.T) {
	// Setup
	handler, mockScheduleService, _, _, _ := setupTestHandler()
	router := handler.SetupRoutes()

	// Test data
	requestBody := models.VisitStartRequest{
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	// Mock expectations
	mockScheduleService.On("StartVisit", 1, &requestBody).Return(nil)

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/schedules/1/start", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Visit started successfully", response["message"])

	// Verify mock expectations
	mockScheduleService.AssertExpectations(t)
}

func TestHandler_UpdateTaskStatus(t *testing.T) {
	// Setup
	handler, _, _, mockTaskService, _ := setupTestHandler()
	router := handler.SetupRoutes()

	// Test data
	requestBody := models.TaskUpdateRequest{
		Status: "completed",
	}

	// Mock expectations
	expectedTask := &models.Task{
		ID:          1,
		ScheduleID:  1,
		Title:       "Give medication",
		Status:      "completed",
		Description: "Administer morning medications",
	}
	mockTaskService.On("UpdateTaskStatus", 1, &requestBody).Return(expectedTask, nil)

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/api/v1/tasks/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, expectedTask.ID, int(response["data"].(map[string]interface{})["id"].(float64)))

	// Verify mock expectations
	mockTaskService.AssertExpectations(t)
}

func TestHandler_GetTodaySchedules(t *testing.T) {
	// Setup
	handler, mockScheduleService, _, _, _ := setupTestHandler()
	router := handler.SetupRoutes()

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

	// Mock expectations
	mockScheduleService.On("GetTodaySchedules", 1).Return(expectedSchedules, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/schedules/today?caregiver_id=1", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Len(t, data, 2)

	// Verify mock expectations
	mockScheduleService.AssertExpectations(t)
}

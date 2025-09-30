package services

import (
	"caregiver-shift-tracker/internal/models"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestTaskService_GetTasksByScheduleID(t *testing.T) {
	// Setup
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewTaskService(mockTaskRepo, logger)

	// Test data
	expectedTasks := []models.Task{
		{ID: 1, ScheduleID: 1, Title: "Give medication", Status: "pending"},
		{ID: 2, ScheduleID: 1, Title: "Check vitals", Status: "completed"},
	}

	// Mock expectations
	mockTaskRepo.On("GetByScheduleID", 1).Return(expectedTasks, nil)

	// Execute
	result, err := service.GetTasksByScheduleID(1)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Give medication", result[0].Title)
	assert.Equal(t, "Check vitals", result[1].Title)

	// Verify mock expectations
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_GetTaskByID(t *testing.T) {
	// Setup
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewTaskService(mockTaskRepo, logger)

	// Test data
	expectedTask := &models.Task{
		ID:          1,
		ScheduleID:  1,
		Title:       "Give medication",
		Description: "Administer morning medications",
		Status:      "pending",
	}

	// Mock expectations
	mockTaskRepo.On("GetByID", 1).Return(expectedTask, nil)

	// Execute
	result, err := service.GetTaskByID(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Give medication", result.Title)
	assert.Equal(t, "pending", result.Status)

	// Verify mock expectations
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_GetTaskByID_NotFound(t *testing.T) {
	// Setup
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewTaskService(mockTaskRepo, logger)

	// Mock expectations
	mockTaskRepo.On("GetByID", 999).Return(nil, nil)

	// Execute
	result, err := service.GetTaskByID(999)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)

	// Verify mock expectations
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_UpdateTaskStatus_Completed(t *testing.T) {
	// Setup
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewTaskService(mockTaskRepo, logger)

	// Test data
	task := &models.Task{
		ID:         1,
		ScheduleID: 1,
		Title:      "Give medication",
		Status:     "pending",
	}
	req := &models.TaskUpdateRequest{
		Status: "completed",
	}

	// Mock expectations
	updatedTask := &models.Task{
		ID:          1,
		ScheduleID:  1,
		Title:       "Give medication",
		Status:      "completed",
		Description: "Administer morning medications",
	}
	mockTaskRepo.On("GetByID", 1).Return(task, nil).Once() // First call returns original task
	mockTaskRepo.On("UpdateStatus", 1, "completed", "").Return(nil)
	mockTaskRepo.On("GetByID", 1).Return(updatedTask, nil).Once() // Second call returns updated task

	// Execute
	updatedTask, err := service.UpdateTaskStatus(1, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedTask)
	assert.Equal(t, "completed", updatedTask.Status)

	// Verify mock expectations
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_UpdateTaskStatus_NotCompleted(t *testing.T) {
	// Setup
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewTaskService(mockTaskRepo, logger)

	// Test data
	task := &models.Task{
		ID:         1,
		ScheduleID: 1,
		Title:      "Give medication",
		Status:     "pending",
	}
	req := &models.TaskUpdateRequest{
		Status: "not_completed",
		Reason: "Client refused medication",
	}

	// Mock expectations
	updatedTask2 := &models.Task{
		ID:          1,
		ScheduleID:  1,
		Title:       "Give medication",
		Status:      "not_completed",
		Description: "Administer morning medications",
		Reason:      "Client refused medication",
	}
	mockTaskRepo.On("GetByID", 1).Return(task, nil).Once() // First call returns original task
	mockTaskRepo.On("UpdateStatus", 1, "not_completed", "Client refused medication").Return(nil)
	mockTaskRepo.On("GetByID", 1).Return(updatedTask2, nil).Once() // Second call returns updated task

	// Execute
	updatedTask, err := service.UpdateTaskStatus(1, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedTask)
	assert.Equal(t, "not_completed", updatedTask.Status)
	assert.Equal(t, "Client refused medication", updatedTask.Reason)

	// Verify mock expectations
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_UpdateTaskStatus_TaskNotFound(t *testing.T) {
	// Setup
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewTaskService(mockTaskRepo, logger)

	req := &models.TaskUpdateRequest{
		Status: "completed",
	}

	// Mock expectations
	mockTaskRepo.On("GetByID", 999).Return(nil, nil)

	// Execute
	updatedTask, err := service.UpdateTaskStatus(999, req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task not found")
	assert.Nil(t, updatedTask)

	// Verify mock expectations
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_UpdateTaskStatus_ValidationError_MissingReason(t *testing.T) {
	// Setup
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewTaskService(mockTaskRepo, logger)

	req := &models.TaskUpdateRequest{
		Status: "not_completed",
		// Missing reason
	}

	// Execute
	updatedTask, err := service.UpdateTaskStatus(1, req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reason is required when status is not_completed")
	assert.Nil(t, updatedTask)
}

func TestTaskService_UpdateTaskStatus_ValidationError_InvalidStatus(t *testing.T) {
	// Setup
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewTaskService(mockTaskRepo, logger)

	req := &models.TaskUpdateRequest{
		Status: "invalid_status",
	}

	// Execute
	updatedTask, err := service.UpdateTaskStatus(1, req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status")
	assert.Nil(t, updatedTask)
}

func TestTaskService_CreateTask(t *testing.T) {
	// Setup
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewTaskService(mockTaskRepo, logger)

	// Test data
	task := &models.Task{
		ScheduleID:  1,
		Title:       "Give medication",
		Description: "Administer morning medications",
		Status:      "pending",
	}

	// Mock expectations
	mockTaskRepo.On("Create", task).Return(nil)

	// Execute
	err := service.CreateTask(task)

	// Assert
	assert.NoError(t, err)

	// Verify mock expectations
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_CreateTask_ValidationError(t *testing.T) {
	// Setup
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewTaskService(mockTaskRepo, logger)

	// Test data with missing title
	task := &models.Task{
		ScheduleID: 1,
		// Missing title
		Description: "Administer morning medications",
		Status:      "pending",
	}

	// Execute
	err := service.CreateTask(task)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")
}

func TestTaskService_DeleteTask(t *testing.T) {
	// Setup
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewTaskService(mockTaskRepo, logger)

	// Test data
	task := &models.Task{
		ID:         1,
		ScheduleID: 1,
		Title:      "Give medication",
		Status:     "pending",
	}

	// Mock expectations
	mockTaskRepo.On("GetByID", 1).Return(task, nil)
	mockTaskRepo.On("Delete", 1).Return(nil)

	// Execute
	err := service.DeleteTask(1)

	// Assert
	assert.NoError(t, err)

	// Verify mock expectations
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_DeleteTask_NotFound(t *testing.T) {
	// Setup
	mockTaskRepo := new(MockTaskRepository)
	logger := logrus.New()
	service := NewTaskService(mockTaskRepo, logger)

	// Mock expectations
	mockTaskRepo.On("GetByID", 999).Return(nil, nil)

	// Execute
	err := service.DeleteTask(999)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task not found")

	// Verify mock expectations
	mockTaskRepo.AssertExpectations(t)
}

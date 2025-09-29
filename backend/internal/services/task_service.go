package services

import (
	"caregiver-shift-tracker/internal/models"
	"caregiver-shift-tracker/internal/repositories"
	"fmt"

	"github.com/sirupsen/logrus"
)

// TaskService handles business logic for tasks
type TaskService struct {
	taskRepo repositories.TaskRepository
	logger   *logrus.Logger
}

// NewTaskService creates a new task service
func NewTaskService(taskRepo repositories.TaskRepository, logger *logrus.Logger) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
		logger:   logger,
	}
}

// GetTasksByScheduleID retrieves all tasks for a schedule
func (s *TaskService) GetTasksByScheduleID(scheduleID int) ([]models.Task, error) {
	s.logger.WithField("schedule_id", scheduleID).Debug("Getting tasks by schedule ID")

	tasks, err := s.taskRepo.GetByScheduleID(scheduleID)
	if err != nil {
		s.logger.WithError(err).WithField("schedule_id", scheduleID).Error("Failed to get tasks")
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"schedule_id": scheduleID,
		"count":       len(tasks),
	}).Debug("Successfully retrieved tasks")

	return tasks, nil
}

// GetTaskByID retrieves a task by ID
func (s *TaskService) GetTaskByID(id int) (*models.Task, error) {
	s.logger.WithField("task_id", id).Debug("Getting task by ID")

	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).WithField("task_id", id).Error("Failed to get task")
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	if task == nil {
		s.logger.WithField("task_id", id).Debug("Task not found")
		return nil, nil
	}

	s.logger.WithField("task_id", id).Debug("Successfully retrieved task")
	return task, nil
}

// UpdateTaskStatus updates the status of a task
func (s *TaskService) UpdateTaskStatus(id int, req *models.TaskUpdateRequest) error {
	s.logger.WithFields(logrus.Fields{
		"task_id": id,
		"status":  req.Status,
		"reason":  req.Reason,
	}).Info("Updating task status")

	// Validate the request
	if err := s.validateTaskUpdateRequest(req); err != nil {
		s.logger.WithError(err).WithField("task_id", id).Error("Task update validation failed")
		return fmt.Errorf("task update validation failed: %w", err)
	}

	// Check if task exists
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).WithField("task_id", id).Error("Failed to get task")
		return fmt.Errorf("failed to get task: %w", err)
	}

	if task == nil {
		s.logger.WithField("task_id", id).Warn("Task not found")
		return fmt.Errorf("task not found")
	}

	// Update the task status
	if err := s.taskRepo.UpdateStatus(id, req.Status, req.Reason); err != nil {
		s.logger.WithError(err).WithField("task_id", id).Error("Failed to update task status")
		return fmt.Errorf("failed to update task status: %w", err)
	}

	s.logger.WithField("task_id", id).Info("Successfully updated task status")
	return nil
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(task *models.Task) error {
	s.logger.WithField("schedule_id", task.ScheduleID).Debug("Creating task")

	if err := s.validateTask(task); err != nil {
		s.logger.WithError(err).WithField("schedule_id", task.ScheduleID).Error("Task validation failed")
		return fmt.Errorf("task validation failed: %w", err)
	}

	if err := s.taskRepo.Create(task); err != nil {
		s.logger.WithError(err).WithField("schedule_id", task.ScheduleID).Error("Failed to create task")
		return fmt.Errorf("failed to create task: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"task_id":     task.ID,
		"schedule_id": task.ScheduleID,
	}).Info("Successfully created task")

	return nil
}

// UpdateTask updates an existing task
func (s *TaskService) UpdateTask(task *models.Task) error {
	s.logger.WithField("task_id", task.ID).Debug("Updating task")

	if err := s.validateTask(task); err != nil {
		s.logger.WithError(err).WithField("task_id", task.ID).Error("Task validation failed")
		return fmt.Errorf("task validation failed: %w", err)
	}

	if err := s.taskRepo.Update(task); err != nil {
		s.logger.WithError(err).WithField("task_id", task.ID).Error("Failed to update task")
		return fmt.Errorf("failed to update task: %w", err)
	}

	s.logger.WithField("task_id", task.ID).Info("Successfully updated task")
	return nil
}

// DeleteTask deletes a task
func (s *TaskService) DeleteTask(id int) error {
	s.logger.WithField("task_id", id).Info("Deleting task")

	// Check if task exists
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).WithField("task_id", id).Error("Failed to get task")
		return fmt.Errorf("failed to get task: %w", err)
	}

	if task == nil {
		s.logger.WithField("task_id", id).Warn("Task not found")
		return fmt.Errorf("task not found")
	}

	if err := s.taskRepo.Delete(id); err != nil {
		s.logger.WithError(err).WithField("task_id", id).Error("Failed to delete task")
		return fmt.Errorf("failed to delete task: %w", err)
	}

	s.logger.WithField("task_id", id).Info("Successfully deleted task")
	return nil
}

// validateTask validates task data
func (s *TaskService) validateTask(task *models.Task) error {
	if task.ScheduleID <= 0 {
		return fmt.Errorf("schedule_id is required")
	}

	if task.Title == "" {
		return fmt.Errorf("title is required")
	}

	if task.Status == "" {
		return fmt.Errorf("status is required")
	}

	validStatuses := map[string]bool{
		"pending":       true,
		"completed":     true,
		"not_completed": true,
	}

	if !validStatuses[task.Status] {
		return fmt.Errorf("invalid status: %s", task.Status)
	}

	// If status is not_completed, reason is required
	if task.Status == "not_completed" && task.Reason == "" {
		return fmt.Errorf("reason is required when status is not_completed")
	}

	return nil
}

// validateTaskUpdateRequest validates task update request
func (s *TaskService) validateTaskUpdateRequest(req *models.TaskUpdateRequest) error {
	if req.Status == "" {
		return fmt.Errorf("status is required")
	}

	validStatuses := map[string]bool{
		"completed":     true,
		"not_completed": true,
	}

	if !validStatuses[req.Status] {
		return fmt.Errorf("invalid status: %s", req.Status)
	}

	// If status is not_completed, reason is required
	if req.Status == "not_completed" && req.Reason == "" {
		return fmt.Errorf("reason is required when status is not_completed")
	}

	return nil
}

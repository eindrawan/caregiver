package services

import (
	"caregiver-shift-tracker/internal/models"
	"caregiver-shift-tracker/internal/repositories"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// VisitService handles business logic for visits
type VisitService struct {
	visitRepo repositories.VisitRepository
	logger    *logrus.Logger
}

// NewVisitService creates a new visit service
func NewVisitService(visitRepo repositories.VisitRepository, logger *logrus.Logger) *VisitService {
	return &VisitService{
		visitRepo: visitRepo,
		logger:    logger,
	}
}

// GetVisitByScheduleID retrieves a visit by schedule ID
func (s *VisitService) GetVisitByScheduleID(scheduleID int) (*models.Visit, error) {
	s.logger.WithField("schedule_id", scheduleID).Debug("Getting visit by schedule ID")

	visit, err := s.visitRepo.GetByScheduleID(scheduleID)
	if err != nil {
		s.logger.WithError(err).WithField("schedule_id", scheduleID).Error("Failed to get visit")
		return nil, fmt.Errorf("failed to get visit: %w", err)
	}

	if visit == nil {
		s.logger.WithField("schedule_id", scheduleID).Debug("Visit not found")
		return nil, nil
	}

	s.logger.WithField("schedule_id", scheduleID).Debug("Successfully retrieved visit")
	return visit, nil
}

// CreateVisit creates a new visit
func (s *VisitService) CreateVisit(visit *models.Visit) error {
	s.logger.WithField("schedule_id", visit.ScheduleID).Debug("Creating visit")

	if err := s.validateVisit(visit); err != nil {
		s.logger.WithError(err).WithField("schedule_id", visit.ScheduleID).Error("Visit validation failed")
		return fmt.Errorf("visit validation failed: %w", err)
	}

	if err := s.visitRepo.Create(visit); err != nil {
		s.logger.WithError(err).WithField("schedule_id", visit.ScheduleID).Error("Failed to create visit")
		return fmt.Errorf("failed to create visit: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"visit_id":    visit.ID,
		"schedule_id": visit.ScheduleID,
	}).Info("Successfully created visit")

	return nil
}

// UpdateVisit updates an existing visit
func (s *VisitService) UpdateVisit(visit *models.Visit) error {
	s.logger.WithField("visit_id", visit.ID).Debug("Updating visit")

	if err := s.validateVisit(visit); err != nil {
		s.logger.WithError(err).WithField("visit_id", visit.ID).Error("Visit validation failed")
		return fmt.Errorf("visit validation failed: %w", err)
	}

	if err := s.visitRepo.Update(visit); err != nil {
		s.logger.WithError(err).WithField("visit_id", visit.ID).Error("Failed to update visit")
		return fmt.Errorf("failed to update visit: %w", err)
	}

	s.logger.WithField("visit_id", visit.ID).Info("Successfully updated visit")
	return nil
}

// EndVisit ends a visit with timestamp and geolocation
func (s *VisitService) EndVisit(scheduleID int, latitude, longitude float64, notes string) error {
	s.logger.WithField("schedule_id", scheduleID).Debug("Ending visit")

	// Get the existing visit to check if it's started
	visit, err := s.visitRepo.GetByScheduleID(scheduleID)
	if err != nil {
		s.logger.WithError(err).WithField("schedule_id", scheduleID).Error("Failed to get visit")
		return fmt.Errorf("failed to get visit: %w", err)
	}

	if visit == nil {
		s.logger.WithField("schedule_id", scheduleID).Error("Visit not found")
		return fmt.Errorf("visit not found for schedule %d", scheduleID)
	}

	if visit.StartTime == nil {
		s.logger.WithField("schedule_id", scheduleID).Error("Cannot end visit that hasn't been started")
		return fmt.Errorf("cannot end visit that hasn't been started")
	}

	// Update the visit to completed status
	now := time.Now()
	visit.EndTime = &now
	visit.EndLatitude = &latitude
	visit.EndLongitude = &longitude
	visit.Status = "completed"
	if notes != "" {
		visit.Notes = notes
	}

	if err := s.visitRepo.Update(visit); err != nil {
		s.logger.WithError(err).WithField("visit_id", visit.ID).Error("Failed to update visit")
		return fmt.Errorf("failed to update visit: %w", err)
	}

	s.logger.WithField("visit_id", visit.ID).Info("Successfully ended visit")
	return nil
}

// validateVisit validates visit data
func (s *VisitService) validateVisit(visit *models.Visit) error {
	if visit.ScheduleID <= 0 {
		return fmt.Errorf("schedule_id is required")
	}

	if visit.Status == "" {
		return fmt.Errorf("status is required")
	}

	validStatuses := map[string]bool{
		"not_started": true,
		"in_progress": true,
		"completed":   true,
	}

	if !validStatuses[visit.Status] {
		return fmt.Errorf("invalid status: %s", visit.Status)
	}

	// If status is in_progress, start time and location should be set
	if visit.Status == "in_progress" {
		if visit.StartTime == nil {
			return fmt.Errorf("start_time is required when status is in_progress")
		}
		if visit.StartLatitude == nil || visit.StartLongitude == nil {
			return fmt.Errorf("start location is required when status is in_progress")
		}
	}

	// If status is completed, end time and location should be set
	if visit.Status == "completed" {
		if visit.StartTime == nil {
			return fmt.Errorf("start_time is required when status is completed")
		}
		if visit.EndTime == nil {
			return fmt.Errorf("end_time is required when status is completed")
		}
		if visit.StartLatitude == nil || visit.StartLongitude == nil {
			return fmt.Errorf("start location is required when status is completed")
		}
		if visit.EndLatitude == nil || visit.EndLongitude == nil {
			return fmt.Errorf("end location is required when status is completed")
		}
		if visit.EndTime.Before(*visit.StartTime) {
			return fmt.Errorf("end_time cannot be before start_time")
		}
	}

	// Validate latitude and longitude ranges
	if visit.StartLatitude != nil && (*visit.StartLatitude < -90 || *visit.StartLatitude > 90) {
		return fmt.Errorf("start_latitude must be between -90 and 90")
	}
	if visit.StartLongitude != nil && (*visit.StartLongitude < -180 || *visit.StartLongitude > 180) {
		return fmt.Errorf("start_longitude must be between -180 and 180")
	}
	if visit.EndLatitude != nil && (*visit.EndLatitude < -90 || *visit.EndLatitude > 90) {
		return fmt.Errorf("end_latitude must be between -90 and 90")
	}
	if visit.EndLongitude != nil && (*visit.EndLongitude < -180 || *visit.EndLongitude > 180) {
		return fmt.Errorf("end_longitude must be between -180 and 180")
	}

	return nil
}

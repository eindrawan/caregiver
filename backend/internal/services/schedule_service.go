package services

import (
	"caregiver-shift-tracker/internal/models"
	"caregiver-shift-tracker/internal/repositories"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// ScheduleService handles business logic for schedules
type ScheduleService struct {
	scheduleRepo repositories.ScheduleRepository
	visitRepo    repositories.VisitRepository
	taskRepo     repositories.TaskRepository
	logger       *logrus.Logger
}

// NewScheduleService creates a new schedule service
func NewScheduleService(
	scheduleRepo repositories.ScheduleRepository,
	visitRepo repositories.VisitRepository,
	taskRepo repositories.TaskRepository,
	logger *logrus.Logger,
) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
		visitRepo:    visitRepo,
		taskRepo:     taskRepo,
		logger:       logger,
	}
}

// ensureLocalTime ensures time values are in the local timezone
func ensureLocalTime(t time.Time) time.Time {
	// Use the system's local timezone
	return t.In(time.Local)
}

// GetAllSchedules retrieves all schedules with optional filtering
func (s *ScheduleService) GetAllSchedules(filter *models.ScheduleFilter) ([]models.Schedule, error) {
	s.logger.WithFields(logrus.Fields{
		"filter": filter,
	}).Debug("Getting all schedules")

	schedules, err := s.scheduleRepo.GetAll(filter)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get schedules")
		return nil, fmt.Errorf("failed to get schedules: %w", err)
	}

	// Enrich schedules with visit and task data
	for i := range schedules {
		if err := s.enrichSchedule(&schedules[i]); err != nil {
			s.logger.WithError(err).WithField("schedule_id", schedules[i].ID).Warn("Failed to enrich schedule")
		}
	}

	s.logger.WithField("count", len(schedules)).Debug("Successfully retrieved schedules")
	return schedules, nil
}

// GetScheduleByID retrieves a schedule by ID with full details
func (s *ScheduleService) GetScheduleByID(id int) (*models.Schedule, error) {
	s.logger.WithField("schedule_id", id).Debug("Getting schedule by ID")

	schedule, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).WithField("schedule_id", id).Error("Failed to get schedule")
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	if schedule == nil {
		s.logger.WithField("schedule_id", id).Warn("Schedule not found")
		return nil, nil
	}

	// Enrich with visit and task data
	if err := s.enrichSchedule(schedule); err != nil {
		s.logger.WithError(err).WithField("schedule_id", id).Warn("Failed to enrich schedule")
	}

	s.logger.WithField("schedule_id", id).Debug("Successfully retrieved schedule")
	return schedule, nil
}

// GetTodaySchedules retrieves today's schedules for a caregiver
func (s *ScheduleService) GetTodaySchedules(caregiverID int) ([]models.Schedule, error) {
	s.logger.WithField("caregiver_id", caregiverID).Debug("Getting today's schedules")

	schedules, err := s.scheduleRepo.GetToday(caregiverID)
	if err != nil {
		s.logger.WithError(err).WithField("caregiver_id", caregiverID).Error("Failed to get today's schedules")
		return nil, fmt.Errorf("failed to get today's schedules: %w", err)
	}

	// Enrich schedules with visit and task data
	for i := range schedules {
		if err := s.enrichSchedule(&schedules[i]); err != nil {
			s.logger.WithError(err).WithField("schedule_id", schedules[i].ID).Warn("Failed to enrich schedule")
		}
	}

	// Update schedule status based on visit status and time
	for i := range schedules {
		s.updateScheduleStatus(&schedules[i])
	}

	s.logger.WithFields(logrus.Fields{
		"caregiver_id": caregiverID,
		"count":        len(schedules),
	}).Debug("Successfully retrieved today's schedules")

	return schedules, nil
}

// GetScheduleStats retrieves schedule statistics for a caregiver
func (s *ScheduleService) GetScheduleStats(caregiverID int) (*models.ScheduleStats, error) {
	s.logger.WithField("caregiver_id", caregiverID).Debug("Getting schedule stats")

	stats, err := s.scheduleRepo.GetStats(caregiverID)
	if err != nil {
		s.logger.WithError(err).WithField("caregiver_id", caregiverID).Error("Failed to get schedule stats")
		return nil, fmt.Errorf("failed to get schedule stats: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"caregiver_id": caregiverID,
		"stats":        stats,
	}).Debug("Successfully retrieved schedule stats")

	return stats, nil
}

// StartVisit starts a visit for a schedule
func (s *ScheduleService) StartVisit(scheduleID int, req *models.VisitStartRequest) error {
	s.logger.WithFields(logrus.Fields{
		"schedule_id": scheduleID,
		"latitude":    req.Latitude,
		"longitude":   req.Longitude,
	}).Info("Starting visit")

	// Validate schedule exists
	schedule, err := s.scheduleRepo.GetByID(scheduleID)
	if err != nil {
		s.logger.WithError(err).WithField("schedule_id", scheduleID).Error("Failed to get schedule")
		return fmt.Errorf("failed to get schedule: %w", err)
	}

	if schedule == nil {
		s.logger.WithField("schedule_id", scheduleID).Warn("Schedule not found")
		return fmt.Errorf("schedule not found")
	}

	// Check if visit can be started (not too early, not already completed)
	now := time.Now()
	if now.Before(schedule.StartTime.Add(-30 * time.Minute)) {
		return fmt.Errorf("cannot start visit more than 30 minutes before scheduled time")
	}

	// Start the visit
	if err := s.visitRepo.StartVisit(scheduleID, req.Latitude, req.Longitude); err != nil {
		s.logger.WithError(err).WithField("schedule_id", scheduleID).Error("Failed to start visit")
		return fmt.Errorf("failed to start visit: %w", err)
	}

	// Update schedule status
	schedule.Status = "in_progress"
	if err := s.scheduleRepo.Update(schedule); err != nil {
		s.logger.WithError(err).WithField("schedule_id", scheduleID).Error("Failed to update schedule status")
		return fmt.Errorf("failed to update schedule status: %w", err)
	}

	s.logger.WithField("schedule_id", scheduleID).Info("Successfully started visit")
	return nil
}

// EndVisit ends a visit for a schedule
func (s *ScheduleService) EndVisit(scheduleID int, req *models.VisitEndRequest) error {
	s.logger.WithFields(logrus.Fields{
		"schedule_id": scheduleID,
		"latitude":    req.Latitude,
		"longitude":   req.Longitude,
	}).Info("Ending visit")

	// Validate schedule exists
	schedule, err := s.scheduleRepo.GetByID(scheduleID)
	if err != nil {
		s.logger.WithError(err).WithField("schedule_id", scheduleID).Error("Failed to get schedule")
		return fmt.Errorf("failed to get schedule: %w", err)
	}

	if schedule == nil {
		s.logger.WithField("schedule_id", scheduleID).Warn("Schedule not found")
		return fmt.Errorf("schedule not found")
	}

	// End the visit
	if err := s.visitRepo.EndVisit(scheduleID, req.Latitude, req.Longitude, req.Notes); err != nil {
		s.logger.WithError(err).WithField("schedule_id", scheduleID).Error("Failed to end visit")
		return fmt.Errorf("failed to end visit: %w", err)
	}

	// Update schedule status
	schedule.Status = "completed"
	if err := s.scheduleRepo.Update(schedule); err != nil {
		s.logger.WithError(err).WithField("schedule_id", scheduleID).Error("Failed to update schedule status")
		return fmt.Errorf("failed to update schedule status: %w", err)
	}

	s.logger.WithField("schedule_id", scheduleID).Info("Successfully ended visit")
	return nil
}

// CancelVisit cancels an in-progress visit for a schedule
func (s *ScheduleService) CancelVisit(scheduleID int) error {
	s.logger.WithField("schedule_id", scheduleID).Info("Cancelling visit")

	// Validate schedule exists
	schedule, err := s.scheduleRepo.GetByID(scheduleID)
	if err != nil {
		s.logger.WithError(err).WithField("schedule_id", scheduleID).Error("Failed to get schedule")
		return fmt.Errorf("failed to get schedule: %w", err)
	}

	if schedule == nil {
		s.logger.WithField("schedule_id", scheduleID).Warn("Schedule not found")
		return fmt.Errorf("schedule not found")
	}

	// Check if visit can be cancelled
	// Allow cancellation for scheduled visits (not started yet) and in_progress visits
	if schedule.Status != "scheduled" && schedule.Status != "in_progress" {
		s.logger.WithField("schedule_id", scheduleID).Warn("Visit cannot be cancelled in this status")
		return fmt.Errorf("visit cannot be cancelled in status: %s", schedule.Status)
	}

	// For "scheduled" status, we don't need to update the visit table since it hasn't started
	// For "in_progress" status, we need to cancel the visit and update its status
	if schedule.Status == "in_progress" {
		// Cancel the visit (reset to not_started status)
		if err := s.visitRepo.CancelVisit(scheduleID); err != nil {
			s.logger.WithError(err).WithField("schedule_id", scheduleID).Error("Failed to cancel visit")
			return fmt.Errorf("failed to cancel visit: %w", err)
		}
	}

	// Update schedule status back to scheduled (or keep as scheduled if it was scheduled)
	schedule.Status = "scheduled"
	if err := s.scheduleRepo.Update(schedule); err != nil {
		s.logger.WithError(err).WithField("schedule_id", scheduleID).Error("Failed to update schedule status")
		return fmt.Errorf("failed to update schedule status: %w", err)
	}

	s.logger.WithField("schedule_id", scheduleID).Info("Successfully cancelled visit")
	return nil
}

// enrichSchedule adds visit and task data to a schedule
func (s *ScheduleService) enrichSchedule(schedule *models.Schedule) error {
	// Get visit data
	visit, err := s.visitRepo.GetByScheduleID(schedule.ID)
	if err != nil {
		return fmt.Errorf("failed to get visit: %w", err)
	}
	schedule.Visit = visit

	// Get task data
	tasks, err := s.taskRepo.GetByScheduleID(schedule.ID)
	if err != nil {
		return fmt.Errorf("failed to get tasks: %w", err)
	}
	schedule.Tasks = tasks

	return nil
}

// updateScheduleStatus updates schedule status based on current time and visit status
func (s *ScheduleService) updateScheduleStatus(schedule *models.Schedule) {
	now := time.Now()

	// If visit is completed, schedule is completed
	if schedule.Visit != nil && schedule.Visit.Status == "completed" {
		schedule.Status = "completed"
		return
	}

	// If visit is in progress, schedule is in progress
	if schedule.Visit != nil && schedule.Visit.Status == "in_progress" {
		schedule.Status = "in_progress"
		return
	}

	// If current time is past end time and no visit started, mark as missed
	if now.After(schedule.EndTime) && (schedule.Visit == nil || schedule.Visit.StartTime == nil) {
		schedule.Status = "missed"
		return
	}

	// Otherwise, keep as scheduled
	schedule.Status = "scheduled"
}

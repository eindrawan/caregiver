package repositories

import (
	"caregiver-shift-tracker/internal/models"
	"database/sql"
	"fmt"
	"time"
)

type visitRepository struct {
	db *sql.DB
}

// NewVisitRepository creates a new visit repository
func NewVisitRepository(db *sql.DB) VisitRepository {
	return &visitRepository{db: db}
}

// GetByScheduleID retrieves a visit by schedule ID
func (r *visitRepository) GetByScheduleID(scheduleID int) (*models.Visit, error) {
	query := `
		SELECT id, schedule_id, start_time, end_time, start_latitude, start_longitude,
		       end_latitude, end_longitude, status, notes, created_at, updated_at
		FROM visits 
		WHERE schedule_id = ?`

	var v models.Visit
	var notes sql.NullString
	err := r.db.QueryRow(query, scheduleID).Scan(
		&v.ID, &v.ScheduleID, &v.StartTime, &v.EndTime, &v.StartLatitude, &v.StartLongitude,
		&v.EndLatitude, &v.EndLongitude, &v.Status, &notes, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get visit: %w", err)
	}

	if notes.Valid {
		v.Notes = notes.String
	}

	return &v, nil
}

// Create creates a new visit
func (r *visitRepository) Create(visit *models.Visit) error {
	query := `
	INSERT INTO visits (schedule_id, start_time, end_time, start_latitude, start_longitude,
		                   end_latitude, end_longitude, status, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, visit.ScheduleID, visit.StartTime, visit.EndTime,
		visit.StartLatitude, visit.StartLongitude, visit.EndLatitude, visit.EndLongitude,
		visit.Status, visit.Notes)
	if err != nil {
		return fmt.Errorf("failed to create visit: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	visit.ID = int(id)
	return nil
}

// Update updates an existing visit
func (r *visitRepository) Update(visit *models.Visit) error {
	query := `
	UPDATE visits
	SET start_time = ?, end_time = ?, start_latitude = ?, start_longitude = ?,
		    end_latitude = ?, end_longitude = ?, status = ?, notes = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := r.db.Exec(query, visit.StartTime, visit.EndTime, visit.StartLatitude, visit.StartLongitude,
		visit.EndLatitude, visit.EndLongitude, visit.Status, visit.Notes, visit.ID)
	if err != nil {
		return fmt.Errorf("failed to update visit: %w", err)
	}

	return nil
}

// StartVisit starts a visit with timestamp and geolocation
func (r *visitRepository) StartVisit(scheduleID int, latitude, longitude float64) error {
	now := time.Now()

	// First, check if visit exists
	visit, err := r.GetByScheduleID(scheduleID)
	if err != nil {
		return fmt.Errorf("failed to check existing visit: %w", err)
	}

	if visit == nil {
		// Create new visit
		visit = &models.Visit{
			ScheduleID:     scheduleID,
			StartTime:      &now,
			StartLatitude:  &latitude,
			StartLongitude: &longitude,
			Status:         "in_progress",
		}
		return r.Create(visit)
	} else {
		// Update existing visit
		visit.StartTime = &now
		visit.StartLatitude = &latitude
		visit.StartLongitude = &longitude
		visit.Status = "in_progress"
		return r.Update(visit)
	}
}

// EndVisit ends a visit with timestamp and geolocation
func (r *visitRepository) EndVisit(scheduleID int, latitude, longitude float64, notes string) error {
	now := time.Now()

	visit, err := r.GetByScheduleID(scheduleID)
	if err != nil {
		return fmt.Errorf("failed to get visit: %w", err)
	}

	if visit == nil {
		return fmt.Errorf("visit not found for schedule %d", scheduleID)
	}

	if visit.StartTime == nil {
		return fmt.Errorf("cannot end visit that hasn't been started")
	}

	visit.EndTime = &now
	visit.EndLatitude = &latitude
	visit.EndLongitude = &longitude
	visit.Status = "completed"
	if notes != "" {
		visit.Notes = notes
	}

	return r.Update(visit)
}

// CancelVisit cancels an in-progress visit by resetting it to not_started status
func (r *visitRepository) CancelVisit(scheduleID int) error {
	visit, err := r.GetByScheduleID(scheduleID)
	if err != nil {
		return fmt.Errorf("failed to get visit: %w", err)
	}

	if visit == nil {
		return fmt.Errorf("visit not found for schedule %d", scheduleID)
	}

	if visit.Status != "in_progress" {
		return fmt.Errorf("visit is not in progress")
	}

	// Reset visit to not_started state
	visit.StartTime = nil
	visit.EndTime = nil
	visit.StartLatitude = nil
	visit.StartLongitude = nil
	visit.EndLatitude = nil
	visit.EndLongitude = nil
	visit.Status = "not_started"
	visit.Notes = ""

	return r.Update(visit)
}

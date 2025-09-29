package repositories

import (
	"caregiver-shift-tracker/internal/models"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type scheduleRepository struct {
	db *sql.DB
}

// NewScheduleRepository creates a new schedule repository
func NewScheduleRepository(db *sql.DB) ScheduleRepository {
	return &scheduleRepository{db: db}
}

// GetAll retrieves all schedules with optional filtering
func (r *scheduleRepository) GetAll(filter *models.ScheduleFilter) ([]models.Schedule, error) {
	query := `
		SELECT s.id, s.client_id, s.service_name, s.caregiver_id, s.start_time, s.end_time, s.status, s.notes, s.created_at, s.updated_at,
		       c.id, c.name, c.email, c.phone, c.address, c.city, c.state, c.zip_code, c.latitude, c.longitude, c.notes, c.is_active, c.created_at, c.updated_at
		FROM schedules s
		LEFT JOIN clients c ON s.client_id = c.id
		WHERE 1=1`

	args := []interface{}{}
	argIndex := 1

	if filter != nil {
		if filter.CaregiverID != nil {
			query += fmt.Sprintf(" AND s.caregiver_id = ?%d", argIndex)
			args = append(args, *filter.CaregiverID)
			argIndex++
		}
		if filter.Date != nil {
			query += fmt.Sprintf(" AND DATE(s.start_time) = DATE(?%d)", argIndex)
			args = append(args, filter.Date.Format("2006-01-02"))
			argIndex++
		}
		if filter.Status != nil {
			query += fmt.Sprintf(" AND s.status = ?%d", argIndex)
			args = append(args, *filter.Status)
			argIndex++
		}
	}

	query += " ORDER BY s.start_time ASC"

	if filter != nil {
		if filter.Limit != nil {
			query += fmt.Sprintf(" LIMIT ?%d", argIndex)
			args = append(args, *filter.Limit)
			argIndex++
		}
		if filter.Offset != nil {
			query += fmt.Sprintf(" OFFSET ?%d", argIndex)
			args = append(args, *filter.Offset)
		}
	}

	fmt.Println(query)

	// Replace numbered placeholders with actual ? placeholders for SQLite
	for i := len(args); i >= 1; i-- {
		query = strings.ReplaceAll(query, fmt.Sprintf("?%d", i), "?")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query schedules: %w", err)
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var s models.Schedule
		var c models.Client
		var clientNotes, clientEmail, clientPhone sql.NullString

		err := rows.Scan(
			&s.ID, &s.ClientID, &s.ServiceName, &s.CaregiverID, &s.StartTime, &s.EndTime, &s.Status, &s.Notes, &s.CreatedAt, &s.UpdatedAt,
			&c.ID, &c.Name, &clientEmail, &clientPhone, &c.Address, &c.City, &c.State, &c.ZipCode, &c.Latitude, &c.Longitude, &clientNotes, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}

		// Handle nullable fields
		if clientEmail.Valid {
			c.Email = clientEmail.String
		}
		if clientPhone.Valid {
			c.Phone = clientPhone.String
		}
		if clientNotes.Valid {
			c.Notes = clientNotes.String
		}

		// Set client data
		s.Client = &c
		schedules = append(schedules, s)
	}

	return schedules, nil
}

// GetByID retrieves a schedule by ID
func (r *scheduleRepository) GetByID(id int) (*models.Schedule, error) {
	query := `
		SELECT s.id, s.client_id, s.service_name, s.caregiver_id, s.start_time, s.end_time, s.status, s.notes, s.created_at, s.updated_at,
		       c.id, c.name, c.email, c.phone, c.address, c.city, c.state, c.zip_code, c.latitude, c.longitude, c.notes, c.is_active, c.created_at, c.updated_at
		FROM schedules s
		LEFT JOIN clients c ON s.client_id = c.id
		WHERE s.id = ?`

	var s models.Schedule
	var c models.Client
	var clientNotes, clientEmail, clientPhone sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&s.ID, &s.ClientID, &s.ServiceName, &s.CaregiverID, &s.StartTime, &s.EndTime, &s.Status, &s.Notes, &s.CreatedAt, &s.UpdatedAt,
		&c.ID, &c.Name, &clientEmail, &clientPhone, &c.Address, &c.City, &c.State, &c.ZipCode, &c.Latitude, &c.Longitude, &clientNotes, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	// Handle nullable fields
	if clientEmail.Valid {
		c.Email = clientEmail.String
	}
	if clientPhone.Valid {
		c.Phone = clientPhone.String
	}
	if clientNotes.Valid {
		c.Notes = clientNotes.String
	}

	// Set client data
	s.Client = &c
	return &s, nil
}

// GetToday retrieves today's schedules for a caregiver
func (r *scheduleRepository) GetToday(caregiverID int) ([]models.Schedule, error) {
	today := time.Now()
	filter := &models.ScheduleFilter{
		CaregiverID: &caregiverID,
		Date:        &today,
	}
	return r.GetAll(filter)
}

// GetStats retrieves schedule statistics for a caregiver
func (r *scheduleRepository) GetStats(caregiverID int) (*models.ScheduleStats, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'missed' THEN 1 ELSE 0 END) as missed,
			SUM(CASE WHEN status = 'scheduled' AND DATE(start_time) = DATE('now') THEN 1 ELSE 0 END) as upcoming,
			SUM(CASE WHEN status = 'completed' AND DATE(start_time) = DATE('now') THEN 1 ELSE 0 END) as completed
		FROM schedules 
		WHERE caregiver_id = ?`

	var stats models.ScheduleStats
	err := r.db.QueryRow(query, caregiverID).Scan(&stats.Total, &stats.Missed, &stats.Upcoming, &stats.Completed)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule stats: %w", err)
	}

	return &stats, nil
}

// Create creates a new schedule
func (r *scheduleRepository) Create(schedule *models.Schedule) error {
	query := `
		INSERT INTO schedules (client_id, service_name, caregiver_id, start_time, end_time, status, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, schedule.ClientID, schedule.ServiceName, schedule.CaregiverID,
		schedule.StartTime, schedule.EndTime, schedule.Status, schedule.Notes)
	if err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	schedule.ID = int(id)
	return nil
}

// Update updates an existing schedule
func (r *scheduleRepository) Update(schedule *models.Schedule) error {
	query := `
		UPDATE schedules
		SET client_id = ?, service_name = ?, caregiver_id = ?, start_time = ?, end_time = ?,
		    status = ?, notes = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := r.db.Exec(query, schedule.ClientID, schedule.ServiceName, schedule.CaregiverID,
		schedule.StartTime, schedule.EndTime, schedule.Status, schedule.Notes, schedule.ID)
	if err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}

	return nil
}

// Delete deletes a schedule
func (r *scheduleRepository) Delete(id int) error {
	query := "DELETE FROM schedules WHERE id = ?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}
	return nil
}

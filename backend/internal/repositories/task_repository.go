package repositories

import (
	"caregiver-shift-tracker/internal/models"
	"database/sql"
	"fmt"
	"time"
)

type taskRepository struct {
	db *sql.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

// GetByScheduleID retrieves all tasks for a schedule
func (r *taskRepository) GetByScheduleID(scheduleID int) ([]models.Task, error) {
	query := `
		SELECT id, schedule_id, title, description, status, reason, completed_at, created_at, updated_at
		FROM tasks 
		WHERE schedule_id = ?
		ORDER BY created_at ASC`

	rows, err := r.db.Query(query, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var reason sql.NullString
		err := rows.Scan(&t.ID, &t.ScheduleID, &t.Title, &t.Description, &t.Status,
			&reason, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		if reason.Valid {
			t.Reason = reason.String
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// GetByID retrieves a task by ID
func (r *taskRepository) GetByID(id int) (*models.Task, error) {
	query := `
		SELECT id, schedule_id, title, description, status, reason, completed_at, created_at, updated_at
		FROM tasks 
		WHERE id = ?`

	var t models.Task
	var reason sql.NullString
	err := r.db.QueryRow(query, id).Scan(&t.ID, &t.ScheduleID, &t.Title, &t.Description,
		&t.Status, &reason, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	if reason.Valid {
		t.Reason = reason.String
	}

	return &t, nil
}

// Create creates a new task
func (r *taskRepository) Create(task *models.Task) error {
	query := `
		INSERT INTO tasks (schedule_id, title, description, status, reason, completed_at)
		VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, task.ScheduleID, task.Title, task.Description,
		task.Status, task.Reason, task.CompletedAt)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	task.ID = int(id)
	return nil
}

// Update updates an existing task
func (r *taskRepository) Update(task *models.Task) error {
	query := `
	UPDATE tasks
	SET title = ?, description = ?, status = ?, reason = ?, completed_at = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := r.db.Exec(query, task.Title, task.Description, task.Status,
		task.Reason, task.CompletedAt, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// UpdateStatus updates the status of a task
func (r *taskRepository) UpdateStatus(id int, status, reason string) error {
	var completedAt *time.Time
	if status == "completed" {
		now := time.Now()
		completedAt = &now
	}

	query := `
		UPDATE tasks
	SET status = ?, reason = ?, completed_at = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	// Retry mechanism for database locking issues
	maxRetries := 3
	var err error
	for i := 0; i < maxRetries; i++ {
		_, err = r.db.Exec(query, status, reason, completedAt, id)
		if err == nil {
			// Success, exit the retry loop
			break
		}

		// Check if it's a database locked error (SQLITE_BUSY)
		if isDatabaseLockedError(err) {
			// Wait before retrying (with exponential backoff)
			waitTime := time.Millisecond * time.Duration(100*(i+1))
			time.Sleep(waitTime)
			continue
		} else {
			// Not a database locked error, don't retry
			break
		}
	}

	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	return nil
}

// Delete deletes a task
func (r *taskRepository) Delete(id int) error {
	query := "DELETE FROM tasks WHERE id = ?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

// isDatabaseLockedError checks if an error is a SQLite database locked error
func isDatabaseLockedError(err error) bool {
	if err == nil {
		return false
	}

	// Check for SQLite error code 5 (SQLITE_BUSY)
	errStr := err.Error()
	return containsErrorSubstring(errStr, "database is locked") ||
		containsErrorSubstring(errStr, "SQLITE_BUSY") ||
		containsErrorSubstring(errStr, "busy")
}

// containsErrorSubstring checks if the error string contains a specific substring
func containsErrorSubstring(errStr, substr string) bool {
	return len(errStr) >= len(substr) &&
		(errStr == substr ||
			len(errStr) > len(substr) &&
				(errStr[:len(substr)] == substr ||
					errStr[len(errStr)-len(substr):] == substr ||
					containsSubstring(errStr, substr)))
}

// containsSubstring checks if a string contains a substring
func containsSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

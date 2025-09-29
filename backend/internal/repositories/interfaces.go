package repositories

import (
	"caregiver-shift-tracker/internal/models"
)

// ScheduleRepository defines the interface for schedule data access
type ScheduleRepository interface {
	GetAll(filter *models.ScheduleFilter) ([]models.Schedule, error)
	GetByID(id int) (*models.Schedule, error)
	GetToday(caregiverID int) ([]models.Schedule, error)
	GetStats(caregiverID int) (*models.ScheduleStats, error)
	Create(schedule *models.Schedule) error
	Update(schedule *models.Schedule) error
	Delete(id int) error
}

// VisitRepository defines the interface for visit data access
type VisitRepository interface {
	GetByScheduleID(scheduleID int) (*models.Visit, error)
	Create(visit *models.Visit) error
	Update(visit *models.Visit) error
	StartVisit(scheduleID int, latitude, longitude float64) error
	EndVisit(scheduleID int, latitude, longitude float64, notes string) error
	CancelVisit(scheduleID int) error
}

// TaskRepository defines the interface for task data access
type TaskRepository interface {
	GetByScheduleID(scheduleID int) ([]models.Task, error)
	GetByID(id int) (*models.Task, error)
	Create(task *models.Task) error
	Update(task *models.Task) error
	UpdateStatus(id int, status, reason string) error
	Delete(id int) error
}

// ClientRepository defines the interface for client data access
type ClientRepository interface {
	GetAll(filter *models.ClientFilter) ([]models.Client, error)
	GetByID(id int) (*models.Client, error)
	Create(client *models.Client) error
	Update(client *models.Client) error
	Delete(id int) error
	Search(query string) ([]models.Client, error)
}

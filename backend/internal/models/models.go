package models

import (
	"time"
)

// Schedule represents a caregiver's scheduled visit to a client
type Schedule struct {
	ID          int       `json:"id" db:"id"`
	ClientID    int       `json:"client_id" db:"client_id" validate:"required"`
	ServiceName string    `json:"service_name" db:"service_name"`
	CaregiverID int       `json:"caregiver_id" db:"caregiver_id" validate:"required"`
	StartTime   time.Time `json:"start_time" db:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" db:"end_time" validate:"required"`
	Status      string    `json:"status" db:"status" validate:"required,oneof=scheduled in_progress completed missed"`
	Notes       string    `json:"notes" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Related data
	Client *Client `json:"client,omitempty" db:"-"`
	Visit  *Visit  `json:"visit,omitempty" db:"-"`
	Tasks  []Task  `json:"tasks,omitempty" db:"-"`
}

// Client represents a client with their information and location
type Client struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required"`
	Email     string    `json:"email" db:"email" validate:"email"`
	Phone     string    `json:"phone" db:"phone"`
	Address   string    `json:"address" db:"address" validate:"required"`
	City      string    `json:"city" db:"city" validate:"required"`
	State     string    `json:"state" db:"state" validate:"required"`
	ZipCode   string    `json:"zip_code" db:"zip_code" validate:"required"`
	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
	Notes     string    `json:"notes" db:"notes"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Visit represents the actual visit log with timestamps and geolocation
type Visit struct {
	ID             int        `json:"id" db:"id"`
	ScheduleID     int        `json:"schedule_id" db:"schedule_id" validate:"required"`
	StartTime      *time.Time `json:"start_time" db:"start_time"`
	EndTime        *time.Time `json:"end_time" db:"end_time"`
	StartLatitude  *float64   `json:"start_latitude" db:"start_latitude"`
	StartLongitude *float64   `json:"start_longitude" db:"start_longitude"`
	EndLatitude    *float64   `json:"end_latitude" db:"end_latitude"`
	EndLongitude   *float64   `json:"end_longitude" db:"end_longitude"`
	Status         string     `json:"status" db:"status" validate:"required,oneof=not_started in_progress completed"`
	Notes          string     `json:"notes" db:"notes"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// Task represents a care activity that needs to be completed during a visit
type Task struct {
	ID          int        `json:"id" db:"id"`
	ScheduleID  int        `json:"schedule_id" db:"schedule_id" validate:"required"`
	Title       string     `json:"name" db:"title" validate:"required"` // Map title to name for frontend compatibility
	Description string     `json:"description" db:"description"`
	Status      string     `json:"status" db:"status" validate:"required,oneof=pending completed not_completed"`
	Reason      string     `json:"reason" db:"reason"` // Required when status is "not_completed"
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// VisitStartRequest represents the request to start a visit
type VisitStartRequest struct {
	Latitude  float64 `json:"start_latitude" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"start_longitude" validate:"required,min=-180,max=180"`
}

// VisitEndRequest represents the request to end a visit
type VisitEndRequest struct {
	Latitude  float64 `json:"end_latitude" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"end_longitude" validate:"required,min=-180,max=180"`
	Notes     string  `json:"notes"`
}

// TaskUpdateRequest represents the request to update a task
type TaskUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=completed not_completed"`
	Reason string `json:"reason"` // Required when status is "not_completed"
}

// ScheduleStats represents statistics for the dashboard
type ScheduleStats struct {
	Total     int `json:"total"`
	Missed    int `json:"missed"`
	Upcoming  int `json:"upcoming"`
	Completed int `json:"completed"`
}

// ScheduleFilter represents filters for schedule queries
type ScheduleFilter struct {
	CaregiverID *int       `json:"caregiver_id"`
	Date        *time.Time `json:"date"`
	Status      *string    `json:"status"`
	Limit       *int       `json:"limit"`
	Offset      *int       `json:"offset"`
}

// ClientCreateRequest represents the request to create a new client
type ClientCreateRequest struct {
	Name      string  `json:"name" validate:"required"`
	Email     string  `json:"email" validate:"email"`
	Phone     string  `json:"phone"`
	Address   string  `json:"address" validate:"required"`
	City      string  `json:"city" validate:"required"`
	State     string  `json:"state" validate:"required"`
	ZipCode   string  `json:"zip_code" validate:"required"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Notes     string  `json:"notes"`
}

// ClientUpdateRequest represents the request to update a client
type ClientUpdateRequest struct {
	Name      *string  `json:"name"`
	Email     *string  `json:"email" validate:"omitempty,email"`
	Phone     *string  `json:"phone"`
	Address   *string  `json:"address"`
	City      *string  `json:"city"`
	State     *string  `json:"state"`
	ZipCode   *string  `json:"zip_code"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Notes     *string  `json:"notes"`
	IsActive  *bool    `json:"is_active"`
}

// ClientFilter represents filters for client queries
type ClientFilter struct {
	IsActive *bool   `json:"is_active"`
	City     *string `json:"city"`
	State    *string `json:"state"`
	Search   *string `json:"search"` // Search by name, email, or phone
	Limit    *int    `json:"limit"`
	Offset   *int    `json:"offset"`
}

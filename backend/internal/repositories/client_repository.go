package repositories

import (
	"caregiver-shift-tracker/internal/models"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type clientRepository struct {
	db *sql.DB
}

// NewClientRepository creates a new client repository
func NewClientRepository(db *sql.DB) ClientRepository {
	return &clientRepository{db: db}
}

// GetAll retrieves all clients with optional filtering
func (r *clientRepository) GetAll(filter *models.ClientFilter) ([]models.Client, error) {
	query := `
		SELECT id, name, email, phone, address, city, state, zip_code, latitude, longitude, notes, is_active, created_at, updated_at
		FROM clients 
		WHERE 1=1`

	var args []interface{}
	argIndex := 1

	if filter != nil {
		if filter.IsActive != nil {
			query += fmt.Sprintf(" AND is_active = ?%d", argIndex)
			args = append(args, *filter.IsActive)
			argIndex++
		}

		if filter.City != nil {
			query += fmt.Sprintf(" AND LOWER(city) = LOWER(?%d)", argIndex)
			args = append(args, *filter.City)
			argIndex++
		}

		if filter.State != nil {
			query += fmt.Sprintf(" AND LOWER(state) = LOWER(?%d)", argIndex)
			args = append(args, *filter.State)
			argIndex++
		}

		if filter.Search != nil {
			query += fmt.Sprintf(" AND (LOWER(name) LIKE LOWER(?%d) OR LOWER(email) LIKE LOWER(?%d) OR phone LIKE ?%d)", argIndex, argIndex+1, argIndex+2)
			searchTerm := "%" + *filter.Search + "%"
			args = append(args, searchTerm, searchTerm, searchTerm)
			argIndex += 3
		}
	}

	query += " ORDER BY name ASC"

	if filter != nil {
		if filter.Limit != nil {
			query += fmt.Sprintf(" LIMIT ?%d", argIndex)
			args = append(args, *filter.Limit)
			argIndex++
		}

		if filter.Offset != nil {
			query += fmt.Sprintf(" OFFSET ?%d", argIndex)
			args = append(args, *filter.Offset)
			argIndex++
		}
	}

	// Replace numbered placeholders with actual ? placeholders for SQLite
	for i := len(args); i >= 1; i-- {
		query = strings.ReplaceAll(query, fmt.Sprintf("?%d", i), "?")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query clients: %w", err)
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var c models.Client
		var email, phone, notes sql.NullString

		err := rows.Scan(
			&c.ID, &c.Name, &email, &phone, &c.Address, &c.City, &c.State, &c.ZipCode,
			&c.Latitude, &c.Longitude, &notes, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan client: %w", err)
		}

		if email.Valid {
			c.Email = email.String
		}
		if phone.Valid {
			c.Phone = phone.String
		}
		if notes.Valid {
			c.Notes = notes.String
		}

		clients = append(clients, c)
	}

	return clients, nil
}

// GetByID retrieves a client by ID
func (r *clientRepository) GetByID(id int) (*models.Client, error) {
	query := `
		SELECT id, name, email, phone, address, city, state, zip_code, latitude, longitude, notes, is_active, created_at, updated_at
		FROM clients 
		WHERE id = ?`

	var c models.Client
	var email, phone, notes sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.Name, &email, &phone, &c.Address, &c.City, &c.State, &c.ZipCode,
		&c.Latitude, &c.Longitude, &notes, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	if email.Valid {
		c.Email = email.String
	}
	if phone.Valid {
		c.Phone = phone.String
	}
	if notes.Valid {
		c.Notes = notes.String
	}

	return &c, nil
}

// Create creates a new client
func (r *clientRepository) Create(client *models.Client) error {
	query := `
		INSERT INTO clients (name, email, phone, address, city, state, zip_code, latitude, longitude, notes, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, client.Name, client.Email, client.Phone, client.Address,
		client.City, client.State, client.ZipCode, client.Latitude, client.Longitude,
		client.Notes, client.IsActive)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	client.ID = int(id)
	client.CreatedAt = time.Now()
	client.UpdatedAt = time.Now()
	return nil
}

// Update updates an existing client
func (r *clientRepository) Update(client *models.Client) error {
	query := `
		UPDATE clients 
		SET name = ?, email = ?, phone = ?, address = ?, city = ?, state = ?, zip_code = ?, 
		    latitude = ?, longitude = ?, notes = ?, is_active = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := r.db.Exec(query, client.Name, client.Email, client.Phone, client.Address,
		client.City, client.State, client.ZipCode, client.Latitude, client.Longitude,
		client.Notes, client.IsActive, client.ID)
	if err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}

	client.UpdatedAt = time.Now()
	return nil
}

// Delete deletes a client
func (r *clientRepository) Delete(id int) error {
	// Check if client has any schedules
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM schedules WHERE client_id = ?", id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check client schedules: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete client: client has %d associated schedules", count)
	}

	query := "DELETE FROM clients WHERE id = ?"
	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}
	return nil
}

// Search searches for clients by name, email, or phone
func (r *clientRepository) Search(query string) ([]models.Client, error) {
	filter := &models.ClientFilter{
		Search: &query,
	}
	return r.GetAll(filter)
}

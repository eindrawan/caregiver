package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Initialize creates and returns a database connection
func Initialize(databaseURL string) (*sql.DB, error) {
	// Add query parameters for better concurrency handling
	db, err := sql.Open("sqlite", databaseURL+"?_busy_timeout=30000&_txlock=immediate")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings for better concurrency
	db.SetMaxOpenConns(25)   // Allow up to 25 open connections
	db.SetMaxIdleConns(25)   // Keep up to 25 idle connections
	db.SetConnMaxLifetime(0) // Connections can be reused indefinitely

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable foreign keys for SQLite
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Set additional SQLite pragmas for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, fmt.Errorf("failed to set WAL journal mode: %w", err)
	}

	if _, err := db.Exec("PRAGMA synchronous = NORMAL"); err != nil {
		return nil, fmt.Errorf("failed to set synchronous mode: %w", err)
	}

	if _, err := db.Exec("PRAGMA cache_size = 100"); err != nil {
		return nil, fmt.Errorf("failed to set cache size: %w", err)
	}

	if _, err := db.Exec("PRAGMA temp_store = MEMORY"); err != nil {
		return nil, fmt.Errorf("failed to set temp store to memory: %w", err)
	}

	// Configure SQLite to handle datetime values properly
	if _, err := db.Exec("PRAGMA datetime_precision = 'seconds'"); err != nil {
		fmt.Printf("Warning: Could not set datetime precision: %v\n", err)
	}

	return db, nil
}

// Migrate runs database migrations
func Migrate(db *sql.DB) error {
	migrations := []string{
		createClientsTable,
		createSchedulesTable,
		createVisitsTable,
		createTasksTable,
		insertSampleData,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to run migration %d: %w", i+1, err)
		}
	}

	return nil
}

const createClientsTable = `
CREATE TABLE IF NOT EXISTS clients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT,
    phone TEXT,
    address TEXT NOT NULL,
    city TEXT NOT NULL,
    state TEXT NOT NULL,
    zip_code TEXT NOT NULL,
    latitude REAL DEFAULT 0,
    longitude REAL DEFAULT 0,
    notes TEXT,
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

const createSchedulesTable = `
CREATE TABLE IF NOT EXISTS schedules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id INTEGER NOT NULL,
    service_name TEXT,
    caregiver_id INTEGER NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    status TEXT NOT NULL DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'in_progress', 'completed', 'missed')),
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (client_id) REFERENCES clients(id)
);`

const createVisitsTable = `
CREATE TABLE IF NOT EXISTS visits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    schedule_id INTEGER NOT NULL UNIQUE,
    start_time DATETIME,
    end_time DATETIME,
    start_latitude REAL,
    start_longitude REAL,
    end_latitude REAL,
    end_longitude REAL,
    status TEXT NOT NULL DEFAULT 'not_started' CHECK (status IN ('not_started', 'in_progress', 'completed')),
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE CASCADE
);`

const createTasksTable = `
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    schedule_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'not_completed')),
    reason TEXT,
    completed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE CASCADE
);`

const insertSampleData = `
-- Insert sample clients
INSERT OR IGNORE INTO clients (id, name, email, phone, address, city, state, zip_code, latitude, longitude, notes, is_active) VALUES
(101, 'John Smith', 'john.smith@email.com', '+44 1232 212 3233', '123 Main St', 'Springfield', 'IL', '62701', 39.7817, -89.6501, 'Regular client, prefers morning visits', 1),
(102, 'Mary Johnson', 'mary.johnson@email.com', '+44 1232 212 3234', '456 Oak Ave', 'Springfield', 'IL', '62702', 39.7990, -89.6440, 'Needs medication assistance', 1),
(103, 'Robert Brown', 'robert.brown@email.com', '+44 1232 212 3235', '789 Pine Rd', 'Springfield', 'IL', '62703', 39.7665, -89.6808, 'Companionship service client', 1),
(104, 'Sarah Davis', 'sarah.davis@email.com', '+44 1232 212 3236', '123 Main St', 'Springfield', 'IL', '62701', 39.7817, -89.6501, 'Physical therapy client', 1);

-- Insert sample schedules
INSERT OR IGNORE INTO schedules (id, client_id, service_name, caregiver_id, start_time, end_time, status, notes) VALUES
(1, 101, 'Personal Care Service', 1, datetime('now', '+1 hour'), datetime('now', '+3 hours'), 'scheduled', 'Regular morning visit'),
(2, 102, 'Medication Management', 1, datetime('now', '+4 hours'), datetime('now', '+6 hours'), 'scheduled', 'Afternoon medication assistance'),
(3, 103, 'Companionship Service', 1, datetime('now', '-2 hours'), datetime('now'), 'missed', 'Client was not home'),
(4, 104, 'Personal Care Service', 1, datetime('now', '+1 day', '+2 hours'), datetime('now', '+1 day', '+4 hours'), 'scheduled', 'Tomorrow morning visit');

-- Insert sample visits
INSERT OR IGNORE INTO visits (id, schedule_id, status) VALUES
(1, 1, 'not_started'),
(2, 2, 'not_started'),
(3, 3, 'not_started'),
(4, 4, 'not_started');

-- Insert sample tasks
INSERT OR IGNORE INTO tasks (id, schedule_id, title, description, status) VALUES
(1, 1, 'Give medication', 'Administer morning medications as prescribed', 'pending'),
(2, 1, 'Check vital signs', 'Take blood pressure and temperature', 'pending'),
(3, 1, 'Assist with bathing', 'Help client with personal hygiene', 'pending'),
(4, 2, 'Prepare lunch', 'Prepare and serve nutritious lunch', 'pending'),
(5, 2, 'Light housekeeping', 'Tidy up living areas', 'pending'),
(6, 4, 'Physical therapy exercises', 'Guide client through prescribed exercises', 'pending'),
(7, 4, 'Medication review', 'Review medication schedule with client', 'pending');`

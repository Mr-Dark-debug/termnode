package db

import (
	"database/sql"
	"fmt"
	"time"
)

// Repository provides CRUD operations for IoT events.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new repository backed by the given database.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Insert adds a new IoT event and returns its ID.
func (r *Repository) Insert(event IoTEvent) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO events (topic, payload, source, timestamp) VALUES (?, ?, ?, ?)",
		event.Topic, event.Payload, event.Source, event.Timestamp,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert event: %w", err)
	}
	return result.LastInsertId()
}

// Recent returns the most recent events up to the given limit.
func (r *Repository) Recent(limit int) ([]IoTEvent, error) {
	rows, err := r.db.Query(
		"SELECT id, topic, payload, source, timestamp FROM events ORDER BY timestamp DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanEvents(rows)
}

// ByTopic returns recent events for a specific topic.
func (r *Repository) ByTopic(topic string, limit int) ([]IoTEvent, error) {
	rows, err := r.db.Query(
		"SELECT id, topic, payload, source, timestamp FROM events WHERE topic = ? ORDER BY timestamp DESC LIMIT ?",
		topic, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanEvents(rows)
}

// Count returns the total number of events.
func (r *Repository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM events").Scan(&count)
	return count, err
}

// Purge deletes events older than the given duration.
func (r *Repository) Purge(olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)
	result, err := r.db.Exec("DELETE FROM events WHERE timestamp < ?", cutoff)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func scanEvents(rows *sql.Rows) ([]IoTEvent, error) {
	var events []IoTEvent
	for rows.Next() {
		var e IoTEvent
		if err := rows.Scan(&e.ID, &e.Topic, &e.Payload, &e.Source, &e.Timestamp); err != nil {
			return events, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

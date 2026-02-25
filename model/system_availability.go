package model

import (
	"context"
	"time"
)

// SystemAvailability mirrors the system_availability table.
type SystemAvailability struct {
	TS                  time.Time
	WindowInterval      string
	TotalRequests       int64
	SuccessRequests     int64
	ErrorRequests       int64
	AvailabilityPercent float64
	P95LatencyMs        *int64
	P99LatencyMs        *int64
}

// InsertSystemAvailability inserts one row into transactions.system_availability.
func (c *PostgreDb) InsertSystemAvailability(row SystemAvailability) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	query := `INSERT INTO system_availability (ts, window_interval, total_requests, success_requests, error_requests, availability_percent, p95_latency_ms, p99_latency_ms)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := c.db.ExecContext(ctx, query,
		row.TS,
		row.WindowInterval,
		row.TotalRequests,
		row.SuccessRequests,
		row.ErrorRequests,
		row.AvailabilityPercent,
		row.P95LatencyMs,
		row.P99LatencyMs,
	)
	return err
}

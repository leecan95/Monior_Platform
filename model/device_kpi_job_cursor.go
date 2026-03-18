package model

import (
	"context"
	"database/sql"
)

// EnsureDeviceKPIJobCursorTable creates cursor table for incremental jobs.
func (c *PostgreDb) EnsureDeviceKPIJobCursorTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	query := `
CREATE TABLE IF NOT EXISTS device_kpi_job_cursor (
    job_name TEXT PRIMARY KEY,
    last_processed_ts BIGINT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

	_, err := c.db.ExecContext(ctx, query)
	return err
}

// GetDeviceKPIJobCursor returns the cursor value for a job.
func (c *PostgreDb) GetDeviceKPIJobCursor(jobName string) (sql.NullInt64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	var value sql.NullInt64
	err := c.db.QueryRowContext(ctx,
		`SELECT last_processed_ts FROM device_kpi_job_cursor WHERE job_name = $1`,
		jobName,
	).Scan(&value)

	if err == sql.ErrNoRows {
		return sql.NullInt64{Valid: false}, nil
	}
	if err != nil {
		return sql.NullInt64{Valid: false}, err
	}
	return value, nil
}

// UpsertDeviceKPIJobCursor stores cursor value for a job.
func (c *PostgreDb) UpsertDeviceKPIJobCursor(jobName string, lastProcessedTs int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	_, err := c.db.ExecContext(ctx,
		`INSERT INTO device_kpi_job_cursor (job_name, last_processed_ts, updated_at)
VALUES ($1,$2,CURRENT_TIMESTAMP)
ON CONFLICT (job_name) DO UPDATE
SET last_processed_ts = EXCLUDED.last_processed_ts, updated_at = EXCLUDED.updated_at`,
		jobName,
		lastProcessedTs,
	)

	return err
}

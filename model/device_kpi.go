package model

import (
	"context"
	"time"
)

// DeviceKPI stores one KPI snapshot for device conditions.
type DeviceKPI struct {
	Name               string
	MatchedDeviceCount int64
	TotalDeviceCount   int64
	Timestamp          time.Time
}

// EnsureDeviceKPITable creates the device_kpi table if not exists.
func (c *PostgreDb) EnsureDeviceKPITable() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	query := `
CREATE TABLE IF NOT EXISTS device_kpi (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    matched_device_count BIGINT NOT NULL,
    total_device_count BIGINT NOT NULL,
    timestamps TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_device_kpi_name_timestamps
    ON device_kpi (name, timestamps DESC);`

	_, err := c.db.ExecContext(ctx, query)
	return err
}

// InsertDeviceKPI inserts one device KPI row into PostgreSQL.
func (c *PostgreDb) InsertDeviceKPI(row DeviceKPI) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	_, err := c.db.ExecContext(ctx,
		`INSERT INTO device_kpi (name, matched_device_count, total_device_count, timestamps) VALUES ($1,$2,$3,$4)`,
		row.Name,
		row.MatchedDeviceCount,
		row.TotalDeviceCount,
		row.Timestamp,
	)

	return err
}

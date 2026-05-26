package model

import (
	"context"
)

type DeviceKPIImageStoreDailySummary struct {
	SummaryDate          string
	OnlineDeviceCount    int64
	ViolationDeviceCount int64
	ViolationPercent     float64
}

func (c *PostgreDb) EnsureDeviceKPIImageStoreDailySummaryTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	query := `
CREATE TABLE IF NOT EXISTS device_kpi_image_store_daily_summary (
    summary_date DATE PRIMARY KEY,
    online_device_count BIGINT NOT NULL DEFAULT 0,
    violation_device_count BIGINT NOT NULL DEFAULT 0,
    violation_percent NUMERIC(7,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_device_kpi_image_store_daily_summary_date
    ON device_kpi_image_store_daily_summary (summary_date DESC);`

	_, err := c.db.ExecContext(ctx, query)
	return err
}

func (c *PostgreDb) CountDistinctOnlineCameraVehiclesByDate(kpiDate string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	query := `SELECT COUNT(DISTINCT plate_no)
FROM device_kpi_online_camera_daily
WHERE kpi_date = $1::date`

	var total int64
	if err := c.db.QueryRowContext(ctx, query, kpiDate).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (c *PostgreDb) CountDistinctImageStoreViolationVehiclesByDate(kpiDate string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	query := `SELECT COUNT(DISTINCT plate_no)
FROM device_kpi_image_store_violation_daily
WHERE kpi_date = $1::date`

	var total int64
	if err := c.db.QueryRowContext(ctx, query, kpiDate).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (c *PostgreDb) UpsertDeviceKPIImageStoreDailySummary(row DeviceKPIImageStoreDailySummary) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	query := `INSERT INTO device_kpi_image_store_daily_summary
    (summary_date, online_device_count, violation_device_count, violation_percent, created_at, updated_at)
VALUES ($1::date, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (summary_date) DO UPDATE
SET online_device_count = EXCLUDED.online_device_count,
    violation_device_count = EXCLUDED.violation_device_count,
    violation_percent = EXCLUDED.violation_percent,
    updated_at = CURRENT_TIMESTAMP`

	_, err := c.db.ExecContext(ctx, query, row.SummaryDate, row.OnlineDeviceCount, row.ViolationDeviceCount, row.ViolationPercent)
	return err
}

package model

import (
	"context"
	"strings"
	"time"

	"github.com/lib/pq"
)

type DeviceKPIImageStoreViolationDailyRow struct {
	KPIDate           string
	IMEI              string
	PlateNo           string
	MaxGap            int64
	GapViolationCount int64
}

func (c *PostgreDb) EnsureDeviceKPIImageStoreViolationDailyTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	query := `
CREATE TABLE IF NOT EXISTS device_kpi_image_store_violation_daily (
    kpi_date DATE NOT NULL,
    imei TEXT NOT NULL DEFAULT '',
    plate_no TEXT NOT NULL,
    max_gap BIGINT NOT NULL DEFAULT 0,
    gap_violation_count BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (kpi_date, plate_no)
);
CREATE INDEX IF NOT EXISTS idx_device_kpi_image_store_violation_daily_kpi_date
    ON device_kpi_image_store_violation_daily (kpi_date DESC);`

	_, err := c.db.ExecContext(ctx, query)
	return err
}

func (c *PostgreDb) ListMVVehicleViolationRowsByDate(kpiDate string) ([]DeviceKPIImageStoreViolationDailyRow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	query := `SELECT
    COALESCE(imei, ''),
    COALESCE(plate_no, ''),
    COALESCE(max_gap, 0),
    COALESCE(gap_violation_count, 0)
FROM mv_vehicle_violation_full
WHERE kpi_date = $1::date
  AND COALESCE(gap_violation_count, 0) > 0`

	rows, err := c.db.QueryContext(ctx, query, kpiDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]DeviceKPIImageStoreViolationDailyRow, 0)
	for rows.Next() {
		var row DeviceKPIImageStoreViolationDailyRow
		row.KPIDate = kpiDate
		if err := rows.Scan(&row.IMEI, &row.PlateNo, &row.MaxGap, &row.GapViolationCount); err != nil {
			return nil, err
		}
		row.IMEI = strings.TrimSpace(row.IMEI)
		row.PlateNo = strings.TrimSpace(row.PlateNo)
		if row.PlateNo == "" {
			continue
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *PostgreDb) ReplaceDeviceKPIImageStoreViolationDailyByDate(kpiDate string, rows []DeviceKPIImageStoreViolationDailyRow) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM device_kpi_image_store_violation_daily WHERE kpi_date = $1::date`, kpiDate); err != nil {
		_ = tx.Rollback()
		return err
	}

	if len(rows) > 0 {
		kpiDates := make([]string, 0, len(rows))
		imeis := make([]string, 0, len(rows))
		plates := make([]string, 0, len(rows))
		maxGaps := make([]int64, 0, len(rows))
		violCounts := make([]int64, 0, len(rows))
		for _, row := range rows {
			plate := strings.TrimSpace(row.PlateNo)
			if plate == "" {
				continue
			}
			kpiDates = append(kpiDates, kpiDate)
			imeis = append(imeis, strings.TrimSpace(row.IMEI))
			plates = append(plates, plate)
			maxGaps = append(maxGaps, row.MaxGap)
			violCounts = append(violCounts, row.GapViolationCount)
		}

		if len(kpiDates) > 0 {
			insertQuery := `INSERT INTO device_kpi_image_store_violation_daily
    (kpi_date, imei, plate_no, max_gap, gap_violation_count)
SELECT d, i, p, g, c
FROM unnest($1::date[], $2::text[], $3::text[], $4::bigint[], $5::bigint[]) AS t(d, i, p, g, c)`

			if _, err := tx.ExecContext(ctx, insertQuery, pq.Array(kpiDates), pq.Array(imeis), pq.Array(plates), pq.Array(maxGaps), pq.Array(violCounts)); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

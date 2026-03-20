package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

const (
	deviceKpiStatusFlagOffline = "offline"
	deviceKpiStatusFlagBadGPS  = "badgps"
	deviceKpiStatusFlagBadMem  = "bad_memory"
	deviceKpiStatusFlagOnline  = "online"
)

type DeviceKPIDailyStatusVehicle struct {
	PlateNo string
	IMEI    string
}

// EnsureDeviceKPIDailyStatusTable creates the daily status snapshot table if not exists.
func (c *PostgreDb) EnsureDeviceKPIDailyStatusTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	query := `
CREATE TABLE IF NOT EXISTS device_kpi_daily_vehicle_status (
    kpi_date DATE NOT NULL,
    plate_no TEXT NOT NULL,
    imei TEXT NOT NULL DEFAULT '',
    seen_offline BOOLEAN NOT NULL DEFAULT FALSE,
    seen_badgps BOOLEAN NOT NULL DEFAULT FALSE,
    seen_bad_memory BOOLEAN NOT NULL DEFAULT FALSE,
    seen_online BOOLEAN NOT NULL DEFAULT FALSE,
    first_seen_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (kpi_date, plate_no)
);
ALTER TABLE device_kpi_daily_vehicle_status
    ADD COLUMN IF NOT EXISTS imei TEXT NOT NULL DEFAULT '';
ALTER TABLE device_kpi_daily_vehicle_status
    ADD COLUMN IF NOT EXISTS seen_bad_memory BOOLEAN NOT NULL DEFAULT FALSE;
`

	_, err := c.db.ExecContext(ctx, query)
	return err
}

// EnsureDeviceKPIDailyStatusArchiveTable creates archive table for monthly statistics.
func (c *PostgreDb) EnsureDeviceKPIDailyStatusArchiveTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	query := `
CREATE TABLE IF NOT EXISTS device_kpi_daily_vehicle_status_archive (
    kpi_date DATE NOT NULL,
    plate_no TEXT NOT NULL,
    imei TEXT NOT NULL DEFAULT '',
    seen_offline BOOLEAN NOT NULL DEFAULT FALSE,
    seen_badgps BOOLEAN NOT NULL DEFAULT FALSE,
    seen_bad_memory BOOLEAN NOT NULL DEFAULT FALSE,
    seen_online BOOLEAN NOT NULL DEFAULT FALSE,
    first_seen_at TIMESTAMPTZ NULL,
    last_seen_at TIMESTAMPTZ NULL,
    snapshot_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (kpi_date, plate_no)
);
ALTER TABLE device_kpi_daily_vehicle_status_archive
    ADD COLUMN IF NOT EXISTS imei TEXT NOT NULL DEFAULT '';
ALTER TABLE device_kpi_daily_vehicle_status_archive
    ADD COLUMN IF NOT EXISTS seen_bad_memory BOOLEAN NOT NULL DEFAULT FALSE;
CREATE INDEX IF NOT EXISTS idx_device_kpi_daily_vehicle_status_archive_kpi_date
    ON device_kpi_daily_vehicle_status_archive (kpi_date DESC);
`

	_, err := c.db.ExecContext(ctx, query)
	return err
}

// UpsertDeviceKPIDailyStatusFlags upserts the daily seen-flag for a list of vehicles.
func (c *PostgreDb) UpsertDeviceKPIDailyStatusFlags(kpiDate string, vehicles []DeviceKPIDailyStatusVehicle, flag string) error {
	if len(vehicles) == 0 {
		return nil
	}

	column, err := statusFlagToColumn(flag)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	now := time.Now().UTC()
	seen := make(map[string]int, len(vehicles))
	cleanedPlates := make([]string, 0, len(vehicles))
	cleanedIMEIs := make([]string, 0, len(vehicles))
	for _, item := range vehicles {
		plate := strings.TrimSpace(item.PlateNo)
		if plate == "" {
			continue
		}
		imei := strings.TrimSpace(item.IMEI)

		if idx, exists := seen[plate]; exists {
			if imei != "" {
				cleanedIMEIs[idx] = imei
			}
			continue
		}
		seen[plate] = len(cleanedPlates)
		cleanedPlates = append(cleanedPlates, plate)
		cleanedIMEIs = append(cleanedIMEIs, imei)
	}
	if len(cleanedPlates) == 0 {
		return nil
	}

	query := fmt.Sprintf(`INSERT INTO device_kpi_daily_vehicle_status (kpi_date, plate_no, imei, %s, first_seen_at, last_seen_at)
SELECT $1::date, p.plate_no, p.imei, TRUE, $2, $2
FROM unnest($3::text[], $4::text[]) AS p(plate_no, imei)
ON CONFLICT (kpi_date, plate_no) DO UPDATE
SET %s = TRUE,
    last_seen_at = EXCLUDED.last_seen_at,
    imei = CASE
        WHEN EXCLUDED.imei <> '' THEN EXCLUDED.imei
        ELSE device_kpi_daily_vehicle_status.imei
    END
WHERE NOT device_kpi_daily_vehicle_status.%s
   OR (EXCLUDED.imei <> '' AND device_kpi_daily_vehicle_status.imei IS DISTINCT FROM EXCLUDED.imei)`, column, column, column)

	_, err = c.db.ExecContext(ctx, query, kpiDate, now, pq.Array(cleanedPlates), pq.Array(cleanedIMEIs))
	if err != nil {
		return err
	}

	return nil
}

// CountDeviceKPIDailyStatusFlags returns number of unique vehicles seen by each status group on one date.
func (c *PostgreDb) CountDeviceKPIDailyStatusFlags(kpiDate string) (int64, int64, int64, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	query := `SELECT
COALESCE(SUM(CASE WHEN seen_offline THEN 1 ELSE 0 END), 0),
COALESCE(SUM(CASE WHEN seen_badgps THEN 1 ELSE 0 END), 0),
COALESCE(SUM(CASE WHEN seen_bad_memory THEN 1 ELSE 0 END), 0),
COALESCE(SUM(CASE WHEN seen_online THEN 1 ELSE 0 END), 0)
FROM device_kpi_daily_vehicle_status
WHERE kpi_date = $1`

	var offlineCount, badGpsCount, badMemCount, onlineCount int64
	if err := c.db.QueryRowContext(ctx, query, kpiDate).Scan(&offlineCount, &badGpsCount, &badMemCount, &onlineCount); err != nil {
		return 0, 0, 0, 0, err
	}

	return offlineCount, badGpsCount, badMemCount, onlineCount, nil
}

// ResetDeviceKPIDailyStatusFlags resets seen_* flags to false for one KPI date.
func (c *PostgreDb) ResetDeviceKPIDailyStatusFlags(kpiDate string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	query := `UPDATE device_kpi_daily_vehicle_status
SET seen_offline = FALSE,
    seen_badgps = FALSE,
    seen_bad_memory = FALSE,
    seen_online = FALSE
WHERE kpi_date = $1
  AND (seen_offline OR seen_badgps OR seen_bad_memory OR seen_online)`

	_, err := c.db.ExecContext(ctx, query, kpiDate)
	return err
}

// SnapshotDeviceKPIDailyStatus snapshots one KPI date from runtime table to archive table.
func (c *PostgreDb) SnapshotDeviceKPIDailyStatus(kpiDate string, snapshotAt time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := `
INSERT INTO device_kpi_daily_vehicle_status_archive
    (kpi_date, plate_no, imei, seen_offline, seen_badgps, seen_bad_memory, seen_online, first_seen_at, last_seen_at, snapshot_at, created_at, updated_at)
SELECT
    kpi_date, plate_no, imei, seen_offline, seen_badgps, seen_bad_memory, seen_online, first_seen_at, last_seen_at, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
FROM device_kpi_daily_vehicle_status
WHERE kpi_date = $1::date
ON CONFLICT (kpi_date, plate_no) DO UPDATE
SET imei = CASE
        WHEN EXCLUDED.imei <> '' THEN EXCLUDED.imei
        ELSE device_kpi_daily_vehicle_status_archive.imei
    END,
    seen_offline = EXCLUDED.seen_offline,
    seen_badgps = EXCLUDED.seen_badgps,
    seen_bad_memory = EXCLUDED.seen_bad_memory,
    seen_online = EXCLUDED.seen_online,
    first_seen_at = EXCLUDED.first_seen_at,
    last_seen_at = EXCLUDED.last_seen_at,
    snapshot_at = EXCLUDED.snapshot_at,
    updated_at = CURRENT_TIMESTAMP`

	_, err := c.db.ExecContext(ctx, query, kpiDate, snapshotAt.UTC())
	return err
}

// DeleteDeviceKPIDailyStatusByDate deletes runtime rows of one KPI date.
func (c *PostgreDb) DeleteDeviceKPIDailyStatusByDate(kpiDate string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := `DELETE FROM device_kpi_daily_vehicle_status WHERE kpi_date = $1::date`
	res, err := c.db.ExecContext(ctx, query, kpiDate)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func statusFlagToColumn(flag string) (string, error) {
	switch flag {
	case deviceKpiStatusFlagOffline:
		return "seen_offline", nil
	case deviceKpiStatusFlagBadGPS:
		return "seen_badgps", nil
	case deviceKpiStatusFlagBadMem:
		return "seen_bad_memory", nil
	case deviceKpiStatusFlagOnline:
		return "seen_online", nil
	default:
		return "", fmt.Errorf("invalid status flag: %s", flag)
	}
}

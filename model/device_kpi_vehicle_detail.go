package model

import (
	"Monitor_Platform/config"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

// DeviceKPIVehicleDetail stores detailed vehicle snapshot rows per KPI sampling window.
type DeviceKPIVehicleDetail struct {
	KPIName      string
	PlateNo      string
	IMEI         string
	LastUpdateTS int64
	LastUpdateAt time.Time
	SampledAt    time.Time
	SampledHour  time.Time
}

type DeviceVehicleACCState struct {
	PlateNo      string
	IMEI         string
	ACCValue     string
	LastUpdateTS int64
	LastUpdateAt time.Time
}

const (
	deviceKPIVehicleDetailKPIAccOn  = "ACC_ON_OVER_10H"
	deviceKPIVehicleDetailKPIAccOff = "ACC_OFF_OVER_72H"
	deviceKPIVehicleDetailKPIMkn    = "VEHICLE_MKN_OVER_72H"
)

// EnsureDeviceKPIVehicleDetailTable creates/migrates the detail table to one row per plate_no.
func (c *PostgreDb) EnsureDeviceKPIVehicleDetailTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	query := `
CREATE TABLE IF NOT EXISTS device_kpi_vehicle_detail (
    id BIGSERIAL PRIMARY KEY,
    kpi_name TEXT NULL,
    plate_no TEXT NOT NULL,
    imei TEXT NOT NULL DEFAULT '',
    last_update_ts BIGINT NULL,
    last_update_at TIMESTAMPTZ NULL,
    sampled_at TIMESTAMPTZ NOT NULL,
    sampled_hour TIMESTAMPTZ NOT NULL,
    in_acc_on_over_10h BOOLEAN NOT NULL DEFAULT FALSE,
    in_acc_off_over_72h BOOLEAN NOT NULL DEFAULT FALSE,
    in_vehicle_mkn_over_72h BOOLEAN NOT NULL DEFAULT FALSE,
    acc_on_last_update_ts BIGINT NULL,
    acc_on_last_update_at TIMESTAMPTZ NULL,
    acc_on_last_violation_at TIMESTAMPTZ NULL,
    acc_off_last_update_ts BIGINT NULL,
    acc_off_last_update_at TIMESTAMPTZ NULL,
    acc_off_last_violation_at TIMESTAMPTZ NULL,
    mkn_last_update_ts BIGINT NULL,
    mkn_last_update_at TIMESTAMPTZ NULL,
    mkn_last_violation_at TIMESTAMPTZ NULL,
    acc_state_value TEXT NOT NULL DEFAULT '',
    acc_state_since_ts BIGINT NULL,
    acc_state_since_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS kpi_name TEXT NULL;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS last_update_ts BIGINT NULL;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS last_update_at TIMESTAMPTZ NULL;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS in_acc_on_over_10h BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS in_acc_off_over_72h BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS in_vehicle_mkn_over_72h BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS acc_on_last_update_ts BIGINT NULL;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS acc_on_last_update_at TIMESTAMPTZ NULL;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS acc_on_last_violation_at TIMESTAMPTZ NULL;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS acc_off_last_update_ts BIGINT NULL;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS acc_off_last_update_at TIMESTAMPTZ NULL;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS acc_off_last_violation_at TIMESTAMPTZ NULL;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS mkn_last_update_ts BIGINT NULL;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS mkn_last_update_at TIMESTAMPTZ NULL;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS mkn_last_violation_at TIMESTAMPTZ NULL;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS acc_state_value TEXT NOT NULL DEFAULT '';
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS acc_state_since_ts BIGINT NULL;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS acc_state_since_at TIMESTAMPTZ NULL;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS sampled_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE device_kpi_vehicle_detail
    ADD COLUMN IF NOT EXISTS sampled_hour TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE device_kpi_vehicle_detail
    ALTER COLUMN kpi_name DROP NOT NULL;
ALTER TABLE device_kpi_vehicle_detail
    ALTER COLUMN last_update_ts DROP NOT NULL;
ALTER TABLE device_kpi_vehicle_detail
    ALTER COLUMN last_update_at DROP NOT NULL;
ALTER TABLE device_kpi_vehicle_detail
    DROP CONSTRAINT IF EXISTS device_kpi_vehicle_detail_kpi_name_plate_no_sampled_hour_key;
CREATE INDEX IF NOT EXISTS idx_device_kpi_vehicle_detail_sampled_hour
    ON device_kpi_vehicle_detail (sampled_hour DESC);
CREATE INDEX IF NOT EXISTS idx_device_kpi_vehicle_detail_flags
    ON device_kpi_vehicle_detail (in_acc_on_over_10h, in_acc_off_over_72h, in_vehicle_mkn_over_72h);
CREATE INDEX IF NOT EXISTS idx_device_kpi_vehicle_detail_plate
    ON device_kpi_vehicle_detail (plate_no);
`

	_, err := c.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	exists, err := c.hasPlateUniqueIndex(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	if _, err := c.db.ExecContext(ctx, `
WITH ranked AS (
    SELECT ctid,
           ROW_NUMBER() OVER (
               PARTITION BY plate_no
               ORDER BY sampled_hour DESC, sampled_at DESC, id DESC
           ) AS rn
    FROM device_kpi_vehicle_detail
)
DELETE FROM device_kpi_vehicle_detail d
USING ranked r
WHERE d.ctid = r.ctid
  AND r.rn > 1`); err != nil {
		return err
	}

	_, err = c.db.ExecContext(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS uq_device_kpi_vehicle_detail_plate_no ON device_kpi_vehicle_detail (plate_no)`)
	return err
}

func (c *PostgreDb) hasPlateUniqueIndex(ctx context.Context) (bool, error) {
	var exists bool
	err := c.db.QueryRowContext(ctx, `
SELECT EXISTS (
    SELECT 1
    FROM pg_indexes
    WHERE schemaname = current_schema()
      AND tablename = 'device_kpi_vehicle_detail'
      AND indexname = 'uq_device_kpi_vehicle_detail_plate_no'
)`).Scan(&exists)
	return exists, err
}

type deviceACCStatePayload struct {
	plateNos  []string
	imeis     []string
	accValues []string
	ts        []int64
	at        []time.Time
}

func buildDeviceACCStatePayload(rows []DeviceVehicleACCState) deviceACCStatePayload {
	p := deviceACCStatePayload{
		plateNos:  make([]string, 0, len(rows)),
		imeis:     make([]string, 0, len(rows)),
		accValues: make([]string, 0, len(rows)),
		ts:        make([]int64, 0, len(rows)),
		at:        make([]time.Time, 0, len(rows)),
	}
	seen := make(map[string]int, len(rows))

	for _, row := range rows {
		plate := strings.TrimSpace(row.PlateNo)
		if plate == "" {
			continue
		}

		imei := strings.TrimSpace(row.IMEI)
		accValue := strings.ToLower(strings.TrimSpace(row.ACCValue))
		if accValue != "active" && accValue != "inactive" {
			accValue = ""
		}

		ts := row.LastUpdateTS
		at := row.LastUpdateAt.UTC()
		if ts <= 0 {
			continue
		}

		if idx, ok := seen[plate]; ok {
			if ts > p.ts[idx] {
				p.imeis[idx] = imei
				p.accValues[idx] = accValue
				p.ts[idx] = ts
				p.at[idx] = at
			}
			continue
		}
		seen[plate] = len(p.plateNos)
		p.plateNos = append(p.plateNos, plate)
		p.imeis = append(p.imeis, imei)
		p.accValues = append(p.accValues, accValue)
		p.ts = append(p.ts, ts)
		p.at = append(p.at, at)
	}

	return p
}

// SyncDeviceVehicleACCState persists ACC state continuity and refreshes ACC_ON/ACC_OFF flags from state-since timestamps.
func (c *PostgreDb) SyncDeviceVehicleACCState(rows []DeviceVehicleACCState, sampledAtUTC, sampledHourUTC time.Time) error {
	sampledAtUTC = sampledAtUTC.UTC()
	sampledHourUTC = sampledHourUTC.UTC()

	p := buildDeviceACCStatePayload(rows)
	if len(p.plateNos) == 0 {
		return nil
	}

	accOnThresholdTs := sampledAtUTC.Add(-10 * time.Hour).UnixMilli()
	accOffThresholdTs := sampledAtUTC.Add(-72 * time.Hour).UnixMilli()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	upsertStateQuery := `
INSERT INTO device_kpi_vehicle_detail
    (plate_no, imei, sampled_at, sampled_hour, acc_state_value, acc_state_since_ts, acc_state_since_at, created_at, updated_at)
SELECT
    u.plate_no,
    u.imei,
    $1,
    $2,
    u.acc_value,
    u.acc_last_update_ts,
    u.acc_last_update_at,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
FROM unnest($3::text[], $4::text[], $5::text[], $6::bigint[], $7::timestamptz[]) AS u(plate_no, imei, acc_value, acc_last_update_ts, acc_last_update_at)
ON CONFLICT (plate_no) DO UPDATE
SET imei = CASE
        WHEN EXCLUDED.imei <> '' THEN EXCLUDED.imei
        ELSE device_kpi_vehicle_detail.imei
    END,
    sampled_at = EXCLUDED.sampled_at,
    sampled_hour = EXCLUDED.sampled_hour,
    acc_state_value = EXCLUDED.acc_state_value,
    acc_state_since_ts = CASE
        WHEN device_kpi_vehicle_detail.acc_state_value IS DISTINCT FROM EXCLUDED.acc_state_value
            THEN EXCLUDED.acc_state_since_ts
        WHEN device_kpi_vehicle_detail.acc_state_since_ts IS NULL
            THEN EXCLUDED.acc_state_since_ts
        ELSE device_kpi_vehicle_detail.acc_state_since_ts
    END,
    acc_state_since_at = CASE
        WHEN device_kpi_vehicle_detail.acc_state_value IS DISTINCT FROM EXCLUDED.acc_state_value
            THEN EXCLUDED.acc_state_since_at
        WHEN device_kpi_vehicle_detail.acc_state_since_at IS NULL
            THEN EXCLUDED.acc_state_since_at
        ELSE device_kpi_vehicle_detail.acc_state_since_at
    END,
    updated_at = CURRENT_TIMESTAMP`

	if _, err := tx.ExecContext(
		ctx,
		upsertStateQuery,
		sampledAtUTC,
		sampledHourUTC,
		pq.Array(p.plateNos),
		pq.Array(p.imeis),
		pq.Array(p.accValues),
		pq.Array(p.ts),
		pq.Array(p.at),
	); err != nil {
		_ = tx.Rollback()
		return err
	}

	refreshFlagsQuery := `
UPDATE device_kpi_vehicle_detail
SET
    sampled_at = $1,
    sampled_hour = $2,
    in_acc_on_over_10h = (
        acc_state_value = 'active'
        AND acc_state_since_ts IS NOT NULL
        AND acc_state_since_ts <= $3
    ),
    acc_on_last_update_ts = CASE
        WHEN acc_state_value = 'active' THEN acc_state_since_ts
        ELSE NULL
    END,
    acc_on_last_update_at = CASE
        WHEN acc_state_value = 'active' THEN acc_state_since_at
        ELSE NULL
    END,
    acc_on_last_violation_at = CASE
        WHEN (
            acc_state_value = 'active'
            AND acc_state_since_ts IS NOT NULL
            AND acc_state_since_ts <= $3
            AND NOT in_acc_on_over_10h
        ) THEN $1
        ELSE acc_on_last_violation_at
    END,
    in_acc_off_over_72h = (
        acc_state_value = 'inactive'
        AND acc_state_since_ts IS NOT NULL
        AND acc_state_since_ts <= $4
    ),
    acc_off_last_update_ts = CASE
        WHEN acc_state_value = 'inactive' THEN acc_state_since_ts
        ELSE NULL
    END,
    acc_off_last_update_at = CASE
        WHEN acc_state_value = 'inactive' THEN acc_state_since_at
        ELSE NULL
    END,
    acc_off_last_violation_at = CASE
        WHEN (
            acc_state_value = 'inactive'
            AND acc_state_since_ts IS NOT NULL
            AND acc_state_since_ts <= $4
            AND NOT in_acc_off_over_72h
        ) THEN $1
        ELSE acc_off_last_violation_at
    END,
    updated_at = CURRENT_TIMESTAMP
WHERE sampled_hour = $2`

	if _, err := tx.ExecContext(ctx, refreshFlagsQuery, sampledAtUTC, sampledHourUTC, accOnThresholdTs, accOffThresholdTs); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (c *PostgreDb) CountDeviceVehicleACCViolations(sampledHourUTC time.Time) (int64, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	query := `
SELECT
    COALESCE(SUM(CASE WHEN in_acc_on_over_10h THEN 1 ELSE 0 END), 0),
    COALESCE(SUM(CASE WHEN in_acc_off_over_72h THEN 1 ELSE 0 END), 0)
FROM device_kpi_vehicle_detail
WHERE sampled_hour = $1`

	var accOnCount, accOffCount int64
	if err := c.db.QueryRowContext(ctx, query, sampledHourUTC.UTC()).Scan(&accOnCount, &accOffCount); err != nil {
		return 0, 0, err
	}
	return accOnCount, accOffCount, nil
}

func (c *PostgreDb) ListDeviceVehicleACCOnViolationRows(sampledHourUTC time.Time) ([]DeviceKPIVehicleDetail, error) {
	return c.listDeviceVehicleAccRows(sampledHourUTC.UTC(), true)
}

func (c *PostgreDb) ListDeviceVehicleACCOffViolationRows(sampledHourUTC time.Time) ([]DeviceKPIVehicleDetail, error) {
	return c.listDeviceVehicleAccRows(sampledHourUTC.UTC(), false)
}

func (c *PostgreDb) listDeviceVehicleAccRows(sampledHourUTC time.Time, isAccOn bool) ([]DeviceKPIVehicleDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	var query, kpiName string
	if isAccOn {
		kpiName = deviceKPIVehicleDetailKPIAccOn
		query = `
SELECT plate_no, imei, acc_on_last_update_ts, acc_on_last_update_at
FROM device_kpi_vehicle_detail
WHERE sampled_hour = $1
  AND in_acc_on_over_10h = TRUE`
	} else {
		kpiName = deviceKPIVehicleDetailKPIAccOff
		query = `
SELECT plate_no, imei, acc_off_last_update_ts, acc_off_last_update_at
FROM device_kpi_vehicle_detail
WHERE sampled_hour = $1
  AND in_acc_off_over_72h = TRUE`
	}

	rows, err := c.db.QueryContext(ctx, query, sampledHourUTC)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]DeviceKPIVehicleDetail, 0)
	for rows.Next() {
		var plateNo, imei string
		var ts sql.NullInt64
		var at sql.NullTime
		if err := rows.Scan(&plateNo, &imei, &ts, &at); err != nil {
			return nil, err
		}
		if strings.TrimSpace(plateNo) == "" || !ts.Valid || !at.Valid {
			continue
		}
		out = append(out, DeviceKPIVehicleDetail{
			KPIName:      kpiName,
			PlateNo:      strings.TrimSpace(plateNo),
			IMEI:         strings.TrimSpace(imei),
			LastUpdateTS: ts.Int64,
			LastUpdateAt: at.Time.UTC(),
			SampledHour:  sampledHourUTC,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

// EnsureDeviceKPIVehicleDetailSnapshotTable creates archive/snapshot table for daily/monthly statistics.
func (c *PostgreDb) EnsureDeviceKPIVehicleDetailSnapshotTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	query := `
CREATE TABLE IF NOT EXISTS device_kpi_vehicle_detail_snapshot (
    snapshot_date DATE NOT NULL,
    plate_no TEXT NOT NULL,
    imei TEXT NOT NULL DEFAULT '',
    kpi_name TEXT NULL,
    last_update_ts BIGINT NULL,
    last_update_at TIMESTAMPTZ NULL,
    sampled_at TIMESTAMPTZ NOT NULL,
    sampled_hour TIMESTAMPTZ NOT NULL,
    in_acc_on_over_10h BOOLEAN NOT NULL DEFAULT FALSE,
    in_acc_off_over_72h BOOLEAN NOT NULL DEFAULT FALSE,
    in_vehicle_mkn_over_72h BOOLEAN NOT NULL DEFAULT FALSE,
    acc_on_last_update_ts BIGINT NULL,
    acc_on_last_update_at TIMESTAMPTZ NULL,
    acc_on_last_violation_at TIMESTAMPTZ NULL,
    acc_off_last_update_ts BIGINT NULL,
    acc_off_last_update_at TIMESTAMPTZ NULL,
    acc_off_last_violation_at TIMESTAMPTZ NULL,
    mkn_last_update_ts BIGINT NULL,
    mkn_last_update_at TIMESTAMPTZ NULL,
    mkn_last_violation_at TIMESTAMPTZ NULL,
    count_acc_on_over_10h BIGINT NOT NULL DEFAULT 0,
    count_acc_off_over_72h BIGINT NOT NULL DEFAULT 0,
    count_vehicle_mkn_over_72h BIGINT NOT NULL DEFAULT 0,
    snapshot_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (snapshot_date, plate_no)
);
CREATE INDEX IF NOT EXISTS idx_device_kpi_vehicle_detail_snapshot_plate_date
    ON device_kpi_vehicle_detail_snapshot (plate_no, snapshot_date DESC);
CREATE INDEX IF NOT EXISTS idx_device_kpi_vehicle_detail_snapshot_date
    ON device_kpi_vehicle_detail_snapshot (snapshot_date DESC);
`

	_, err := c.db.ExecContext(ctx, query)
	return err
}

// ResetAllDeviceKPIVehicleDetailFlags resets all in_* flags to false.
func (c *PostgreDb) ResetAllDeviceKPIVehicleDetailFlags(sampledAtUTC, sampledHourUTC time.Time) error {
	sampledAtUTC = sampledAtUTC.UTC()
	sampledHourUTC = sampledHourUTC.UTC()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := c.db.ExecContext(ctx, `
UPDATE device_kpi_vehicle_detail
SET in_acc_on_over_10h = FALSE,
    in_acc_off_over_72h = FALSE,
    in_vehicle_mkn_over_72h = FALSE,
    sampled_at = $1,
    sampled_hour = $2,
    updated_at = CURRENT_TIMESTAMP`, sampledAtUTC, sampledHourUTC)

	return err
}

// SnapshotDeviceKPIVehicleDetail snapshots current detail table before day reset and updates cumulative counts.
func (c *PostgreDb) SnapshotDeviceKPIVehicleDetail(snapshotDate string, snapshotAtUTC time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	query := `
WITH src AS (
    SELECT
        plate_no,
        imei,
        kpi_name,
        last_update_ts,
        last_update_at,
        sampled_at,
        sampled_hour,
        in_acc_on_over_10h,
        in_acc_off_over_72h,
        in_vehicle_mkn_over_72h,
        acc_on_last_update_ts,
        acc_on_last_update_at,
        acc_on_last_violation_at,
        acc_off_last_update_ts,
        acc_off_last_update_at,
        acc_off_last_violation_at,
        mkn_last_update_ts,
        mkn_last_update_at,
        mkn_last_violation_at
    FROM device_kpi_vehicle_detail
),
	base AS (
	    SELECT
	        s.*,
	        p.count_acc_on_over_10h AS prev_count_acc_on_over_10h,
	        p.count_acc_off_over_72h AS prev_count_acc_off_over_72h,
	        p.count_vehicle_mkn_over_72h AS prev_count_vehicle_mkn_over_72h,
	        p.acc_on_last_update_at AS prev_acc_on_last_update_at,
	        p.acc_on_last_violation_at AS prev_acc_on_last_violation_at,
	        p.acc_off_last_update_at AS prev_acc_off_last_update_at,
	        p.acc_off_last_violation_at AS prev_acc_off_last_violation_at,
	        p.mkn_last_update_at AS prev_mkn_last_update_at,
	        p.mkn_last_violation_at AS prev_mkn_last_violation_at
	    FROM src s
	    LEFT JOIN LATERAL (
	        SELECT
	            t.count_acc_on_over_10h,
	            t.count_acc_off_over_72h,
	            t.count_vehicle_mkn_over_72h,
	            t.acc_on_last_update_at,
	            t.acc_on_last_violation_at,
	            t.acc_off_last_update_at,
	            t.acc_off_last_violation_at,
	            t.mkn_last_update_at,
	            t.mkn_last_violation_at
	        FROM device_kpi_vehicle_detail_snapshot t
	        WHERE t.plate_no = s.plate_no
	          AND t.snapshot_date < $1::date
        ORDER BY t.snapshot_date DESC, t.snapshot_at DESC
        LIMIT 1
    ) p ON TRUE
)
INSERT INTO device_kpi_vehicle_detail_snapshot
(
    snapshot_date,
    plate_no,
    imei,
    kpi_name,
    last_update_ts,
    last_update_at,
    sampled_at,
    sampled_hour,
    in_acc_on_over_10h,
    in_acc_off_over_72h,
    in_vehicle_mkn_over_72h,
    acc_on_last_update_ts,
    acc_on_last_update_at,
    acc_on_last_violation_at,
    acc_off_last_update_ts,
    acc_off_last_update_at,
    acc_off_last_violation_at,
    mkn_last_update_ts,
    mkn_last_update_at,
    mkn_last_violation_at,
    count_acc_on_over_10h,
    count_acc_off_over_72h,
    count_vehicle_mkn_over_72h,
    snapshot_at,
    created_at,
    updated_at
)
SELECT
    $1::date,
    b.plate_no,
    b.imei,
    b.kpi_name,
    b.last_update_ts,
    b.last_update_at,
    b.sampled_at,
    b.sampled_hour,
    b.in_acc_on_over_10h,
    b.in_acc_off_over_72h,
    b.in_vehicle_mkn_over_72h,
    b.acc_on_last_update_ts,
    b.acc_on_last_update_at,
    b.acc_on_last_violation_at,
    b.acc_off_last_update_ts,
    b.acc_off_last_update_at,
    b.acc_off_last_violation_at,
    b.mkn_last_update_ts,
    b.mkn_last_update_at,
    b.mkn_last_violation_at,
	    COALESCE(b.prev_count_acc_on_over_10h, 0) +
	        CASE
	            WHEN (
	                    b.in_acc_on_over_10h
	                    AND b.acc_on_last_update_at IS NOT NULL
	                    AND (b.prev_acc_on_last_update_at IS NULL OR b.acc_on_last_update_at <> b.prev_acc_on_last_update_at)
	                 ) OR (
	                    b.acc_on_last_violation_at IS NOT NULL
	                    AND (b.prev_acc_on_last_violation_at IS NULL OR b.acc_on_last_violation_at <> b.prev_acc_on_last_violation_at)
	                 )
	                THEN 1
	            ELSE 0
	        END,
	    COALESCE(b.prev_count_acc_off_over_72h, 0) +
	        CASE
	            WHEN (
	                    b.in_acc_off_over_72h
	                    AND b.acc_off_last_update_at IS NOT NULL
	                    AND (b.prev_acc_off_last_update_at IS NULL OR b.acc_off_last_update_at <> b.prev_acc_off_last_update_at)
	                 ) OR (
	                    b.acc_off_last_violation_at IS NOT NULL
	                    AND (b.prev_acc_off_last_violation_at IS NULL OR b.acc_off_last_violation_at <> b.prev_acc_off_last_violation_at)
	                 )
	                THEN 1
	            ELSE 0
	        END,
	    COALESCE(b.prev_count_vehicle_mkn_over_72h, 0) +
	        CASE
	            WHEN (
	                    b.in_vehicle_mkn_over_72h
	                    AND b.mkn_last_update_at IS NOT NULL
	                    AND (b.prev_mkn_last_update_at IS NULL OR b.mkn_last_update_at <> b.prev_mkn_last_update_at)
	                 ) OR (
	                    b.mkn_last_violation_at IS NOT NULL
	                    AND (b.prev_mkn_last_violation_at IS NULL OR b.mkn_last_violation_at <> b.prev_mkn_last_violation_at)
	                 )
	                THEN 1
	            ELSE 0
	        END,
    $2::timestamptz,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
FROM base b
ON CONFLICT (snapshot_date, plate_no) DO UPDATE
SET imei = EXCLUDED.imei,
    kpi_name = EXCLUDED.kpi_name,
    last_update_ts = EXCLUDED.last_update_ts,
    last_update_at = EXCLUDED.last_update_at,
    sampled_at = EXCLUDED.sampled_at,
    sampled_hour = EXCLUDED.sampled_hour,
    in_acc_on_over_10h = EXCLUDED.in_acc_on_over_10h,
    in_acc_off_over_72h = EXCLUDED.in_acc_off_over_72h,
    in_vehicle_mkn_over_72h = EXCLUDED.in_vehicle_mkn_over_72h,
    acc_on_last_update_ts = EXCLUDED.acc_on_last_update_ts,
    acc_on_last_update_at = EXCLUDED.acc_on_last_update_at,
    acc_on_last_violation_at = EXCLUDED.acc_on_last_violation_at,
    acc_off_last_update_ts = EXCLUDED.acc_off_last_update_ts,
    acc_off_last_update_at = EXCLUDED.acc_off_last_update_at,
    acc_off_last_violation_at = EXCLUDED.acc_off_last_violation_at,
    mkn_last_update_ts = EXCLUDED.mkn_last_update_ts,
    mkn_last_update_at = EXCLUDED.mkn_last_update_at,
    mkn_last_violation_at = EXCLUDED.mkn_last_violation_at,
    count_acc_on_over_10h = EXCLUDED.count_acc_on_over_10h,
    count_acc_off_over_72h = EXCLUDED.count_acc_off_over_72h,
    count_vehicle_mkn_over_72h = EXCLUDED.count_vehicle_mkn_over_72h,
    snapshot_at = EXCLUDED.snapshot_at,
    updated_at = CURRENT_TIMESTAMP`

	_, err := c.db.ExecContext(ctx, query, snapshotDate, snapshotAtUTC.UTC())
	return err
}

// UpsertDeviceKPIVehicleDetails updates one-row-per-plate snapshot for current sampled hour.
func (c *PostgreDb) UpsertDeviceKPIVehicleDetails(rows []DeviceKPIVehicleDetail, sampledAtUTC, sampledHourUTC time.Time) error {
	sampledAtUTC = sampledAtUTC.UTC()
	sampledHourUTC = sampledHourUTC.UTC()

	timeout := resolveDetailUpsertTimeout()
	batchSize := resolveDetailUpsertBatchSize()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	accOnRows, accOffRows, mknRows := splitDetailRowsByKPI(rows)

	if err := syncRowsByKPI(ctx, tx, deviceKPIVehicleDetailKPIAccOn, accOnRows, sampledAtUTC, sampledHourUTC, batchSize); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("sync %s rows=%d: %w", deviceKPIVehicleDetailKPIAccOn, len(accOnRows), err)
	}
	if err := syncRowsByKPI(ctx, tx, deviceKPIVehicleDetailKPIAccOff, accOffRows, sampledAtUTC, sampledHourUTC, batchSize); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("sync %s rows=%d: %w", deviceKPIVehicleDetailKPIAccOff, len(accOffRows), err)
	}
	if err := syncRowsByKPI(ctx, tx, deviceKPIVehicleDetailKPIMkn, mknRows, sampledAtUTC, sampledHourUTC, batchSize); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("sync %s rows=%d: %w", deviceKPIVehicleDetailKPIMkn, len(mknRows), err)
	}

	return tx.Commit()
}

func splitDetailRowsByKPI(rows []DeviceKPIVehicleDetail) ([]DeviceKPIVehicleDetail, []DeviceKPIVehicleDetail, []DeviceKPIVehicleDetail) {
	accOn := make([]DeviceKPIVehicleDetail, 0)
	accOff := make([]DeviceKPIVehicleDetail, 0)
	mkn := make([]DeviceKPIVehicleDetail, 0)

	for _, row := range rows {
		switch strings.TrimSpace(row.KPIName) {
		case deviceKPIVehicleDetailKPIAccOn:
			accOn = append(accOn, row)
		case deviceKPIVehicleDetailKPIAccOff:
			accOff = append(accOff, row)
		case deviceKPIVehicleDetailKPIMkn:
			mkn = append(mkn, row)
		}
	}

	return accOn, accOff, mkn
}

type deviceDetailPayload struct {
	plateNos []string
	imeis    []string
	ts       []int64
	at       []time.Time
}

func buildDeviceDetailPayload(rows []DeviceKPIVehicleDetail) deviceDetailPayload {
	p := deviceDetailPayload{
		plateNos: make([]string, 0, len(rows)),
		imeis:    make([]string, 0, len(rows)),
		ts:       make([]int64, 0, len(rows)),
		at:       make([]time.Time, 0, len(rows)),
	}
	seen := make(map[string]int, len(rows))

	for _, row := range rows {
		plate := strings.TrimSpace(row.PlateNo)
		if plate == "" {
			continue
		}
		imei := strings.TrimSpace(row.IMEI)
		at := row.LastUpdateAt.UTC()
		if idx, ok := seen[plate]; ok {
			if row.LastUpdateTS > p.ts[idx] {
				p.imeis[idx] = imei
				p.ts[idx] = row.LastUpdateTS
				p.at[idx] = at
			}
			continue
		}
		seen[plate] = len(p.plateNos)
		p.plateNos = append(p.plateNos, plate)
		p.imeis = append(p.imeis, imei)
		p.ts = append(p.ts, row.LastUpdateTS)
		p.at = append(p.at, at)
	}

	return p
}

func syncRowsByKPI(
	ctx context.Context,
	tx *sql.Tx,
	kpiName string,
	rows []DeviceKPIVehicleDetail,
	sampledAtUTC, sampledHourUTC time.Time,
	batchSize int,
) error {
	flagCol, tsCol, atCol, violationAtCol, err := kpiColumns(kpiName)
	if err != nil {
		return err
	}

	p := buildDeviceDetailPayload(rows)
	if len(p.plateNos) == 0 {
		resetQuery := fmt.Sprintf(`
UPDATE device_kpi_vehicle_detail
SET %s = FALSE,
    sampled_at = $1,
    sampled_hour = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE %s = TRUE`, flagCol, flagCol)
		_, err := tx.ExecContext(ctx, resetQuery, sampledAtUTC, sampledHourUTC)
		return err
	}

	resetEndedQuery := fmt.Sprintf(`
UPDATE device_kpi_vehicle_detail
SET %s = FALSE,
    sampled_at = $1,
    sampled_hour = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE %s = TRUE
  AND NOT (plate_no = ANY($3::text[]))`, flagCol, flagCol)
	if _, err := tx.ExecContext(ctx, resetEndedQuery, sampledAtUTC, sampledHourUTC, pq.Array(p.plateNos)); err != nil {
		return err
	}

	query := fmt.Sprintf(`
INSERT INTO device_kpi_vehicle_detail
    (plate_no, imei, sampled_at, sampled_hour, %s, %s, %s, %s, created_at, updated_at)
SELECT u.plate_no, u.imei, $1, $2, TRUE, u.last_update_ts, u.last_update_at, $1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
FROM unnest($3::text[], $4::text[], $5::bigint[], $6::timestamptz[]) AS u(plate_no, imei, last_update_ts, last_update_at)
ON CONFLICT (plate_no) DO UPDATE
SET imei = CASE WHEN EXCLUDED.imei <> '' THEN EXCLUDED.imei ELSE device_kpi_vehicle_detail.imei END,
    sampled_at = EXCLUDED.sampled_at,
    sampled_hour = EXCLUDED.sampled_hour,
    %s = TRUE,
    %s = CASE
            WHEN device_kpi_vehicle_detail.%s = FALSE OR device_kpi_vehicle_detail.%s IS NULL
                THEN EXCLUDED.%s
            ELSE device_kpi_vehicle_detail.%s
         END,
    %s = CASE
            WHEN device_kpi_vehicle_detail.%s = FALSE OR device_kpi_vehicle_detail.%s IS NULL
                THEN EXCLUDED.%s
            ELSE device_kpi_vehicle_detail.%s
         END,
    %s = CASE
            WHEN device_kpi_vehicle_detail.%s = FALSE OR device_kpi_vehicle_detail.%s IS NULL
                THEN EXCLUDED.%s
            ELSE device_kpi_vehicle_detail.%s
         END,
    updated_at = CURRENT_TIMESTAMP`,
		flagCol, tsCol, atCol, violationAtCol, flagCol,
		tsCol, flagCol, tsCol, tsCol, tsCol,
		atCol, flagCol, atCol, atCol, atCol,
		violationAtCol, flagCol, violationAtCol, violationAtCol, violationAtCol)

	if batchSize <= 0 {
		batchSize = 5000
	}

	for start := 0; start < len(p.plateNos); start += batchSize {
		end := start + batchSize
		if end > len(p.plateNos) {
			end = len(p.plateNos)
		}
		if _, err := tx.ExecContext(
			ctx,
			query,
			sampledAtUTC,
			sampledHourUTC,
			pq.Array(p.plateNos[start:end]),
			pq.Array(p.imeis[start:end]),
			pq.Array(p.ts[start:end]),
			pq.Array(p.at[start:end]),
		); err != nil {
			return fmt.Errorf("upsert batch start=%d end=%d: %w", start, end, err)
		}
	}

	return nil
}

func resolveDetailUpsertTimeout() time.Duration {
	raw := strings.TrimSpace(config.Env(config.EnvKpiDetailUpsertTimeoutSecond, config.DefKpiDetailUpsertTimeoutSecond))
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		fallback, fbErr := strconv.Atoi(config.DefKpiDetailUpsertTimeoutSecond)
		if fbErr != nil || fallback <= 0 {
			return 3 * time.Minute
		}
		return time.Duration(fallback) * time.Second
	}
	return time.Duration(seconds) * time.Second
}

func resolveDetailUpsertBatchSize() int {
	raw := strings.TrimSpace(config.Env(config.EnvKpiDetailUpsertBatchSize, config.DefKpiDetailUpsertBatchSize))
	size, err := strconv.Atoi(raw)
	if err != nil || size <= 0 {
		fallback, fbErr := strconv.Atoi(config.DefKpiDetailUpsertBatchSize)
		if fbErr != nil || fallback <= 0 {
			return 5000
		}
		return fallback
	}
	return size
}

func kpiColumns(kpiName string) (flagCol, tsCol, atCol, violationAtCol string, err error) {
	switch kpiName {
	case deviceKPIVehicleDetailKPIAccOn:
		return "in_acc_on_over_10h", "acc_on_last_update_ts", "acc_on_last_update_at", "acc_on_last_violation_at", nil
	case deviceKPIVehicleDetailKPIAccOff:
		return "in_acc_off_over_72h", "acc_off_last_update_ts", "acc_off_last_update_at", "acc_off_last_violation_at", nil
	case deviceKPIVehicleDetailKPIMkn:
		return "in_vehicle_mkn_over_72h", "mkn_last_update_ts", "mkn_last_update_at", "mkn_last_violation_at", nil
	default:
		return "", "", "", "", fmt.Errorf("unsupported kpi name: %s", kpiName)
	}
}

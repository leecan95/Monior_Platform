package model

import (
	"context"
	"strings"
	"time"

	"github.com/lib/pq"
)

type DeviceKPIOnlineCameraDailyRow struct {
	KPIDate         string
	PlateNo         string
	LicensePlate    string
	IMEI            string
	DeviceModelName string
}

type VehicleCameraRef struct {
	LicensePlate    string
	DeviceModelName string
}

func (c *PostgreDb) EnsureDeviceKPIOnlineCameraDailyTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	query := `
CREATE TABLE IF NOT EXISTS device_kpi_online_camera_daily (
    kpi_date DATE NOT NULL,
    plate_no TEXT NOT NULL,
    license_plate TEXT NOT NULL,
    imei TEXT NOT NULL DEFAULT '',
    device_model_name TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (kpi_date, plate_no)
);
CREATE INDEX IF NOT EXISTS idx_device_kpi_online_camera_daily_kpi_date
    ON device_kpi_online_camera_daily (kpi_date DESC);`

	_, err := c.db.ExecContext(ctx, query)
	return err
}

func (c *PostgreDb) ListDeviceKPIDailySeenOnlineVehicles(kpiDate string) ([]DeviceKPIDailyStatusVehicle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := `SELECT plate_no, imei
FROM device_kpi_daily_vehicle_status
WHERE kpi_date = $1::date
  AND seen_online = TRUE`

	rows, err := c.db.QueryContext(ctx, query, kpiDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]DeviceKPIDailyStatusVehicle, 0)
	for rows.Next() {
		var row DeviceKPIDailyStatusVehicle
		if err := rows.Scan(&row.PlateNo, &row.IMEI); err != nil {
			return nil, err
		}
		row.PlateNo = strings.TrimSpace(row.PlateNo)
		row.IMEI = strings.TrimSpace(row.IMEI)
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

func (c *PostgreDb) ListVehicleCameraRefsByLicensePlates(licensePlates []string) (map[string]VehicleCameraRef, error) {
	if len(licensePlates) == 0 {
		return map[string]VehicleCameraRef{}, nil
	}

	clean := make([]string, 0, len(licensePlates))
	seen := make(map[string]struct{}, len(licensePlates))
	for _, plate := range licensePlates {
		plate = strings.TrimSpace(plate)
		if plate == "" {
			continue
		}
		if _, ok := seen[plate]; ok {
			continue
		}
		seen[plate] = struct{}{}
		clean = append(clean, plate)
	}
	if len(clean) == 0 {
		return map[string]VehicleCameraRef{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	query := `SELECT DISTINCT ON (license_plate)
    license_plate,
    COALESCE(device_model_name, '')
FROM mv_vehicle_camera
WHERE camera = TRUE
  AND license_plate = ANY($1::text[])
ORDER BY license_plate`

	rows, err := c.db.QueryContext(ctx, query, pq.Array(clean))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]VehicleCameraRef, len(clean))
	for rows.Next() {
		var row VehicleCameraRef
		if err := rows.Scan(&row.LicensePlate, &row.DeviceModelName); err != nil {
			return nil, err
		}
		row.LicensePlate = strings.TrimSpace(row.LicensePlate)
		row.DeviceModelName = strings.TrimSpace(row.DeviceModelName)
		if row.LicensePlate == "" {
			continue
		}
		result[row.LicensePlate] = row
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *PostgreDb) ReplaceDeviceKPIOnlineCameraDailyByDate(kpiDate string, rows []DeviceKPIOnlineCameraDailyRow) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM device_kpi_online_camera_daily WHERE kpi_date = $1::date`, kpiDate); err != nil {
		_ = tx.Rollback()
		return err
	}

	if len(rows) > 0 {
		kpiDates := make([]string, 0, len(rows))
		plates := make([]string, 0, len(rows))
		licensePlates := make([]string, 0, len(rows))
		imeis := make([]string, 0, len(rows))
		modelNames := make([]string, 0, len(rows))
		for _, row := range rows {
			plate := strings.TrimSpace(row.PlateNo)
			licensePlate := strings.TrimSpace(row.LicensePlate)
			if plate == "" || licensePlate == "" {
				continue
			}
			kpiDates = append(kpiDates, kpiDate)
			plates = append(plates, plate)
			licensePlates = append(licensePlates, licensePlate)
			imeis = append(imeis, strings.TrimSpace(row.IMEI))
			modelNames = append(modelNames, strings.TrimSpace(row.DeviceModelName))
		}

		if len(kpiDates) > 0 {
			insertQuery := `INSERT INTO device_kpi_online_camera_daily
    (kpi_date, plate_no, license_plate, imei, device_model_name)
SELECT d, p, l, i, m
FROM unnest($1::date[], $2::text[], $3::text[], $4::text[], $5::text[]) AS t(d, p, l, i, m)`
			if _, err := tx.ExecContext(ctx, insertQuery, pq.Array(kpiDates), pq.Array(plates), pq.Array(licensePlates), pq.Array(imeis), pq.Array(modelNames)); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

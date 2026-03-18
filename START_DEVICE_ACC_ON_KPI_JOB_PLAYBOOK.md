# StartDeviceAccOnKpiJob Playbook

This document is the source of truth to re-implement `StartDeviceAccOnKpiJob` on a new branch.

## 1) Entry Point And Schedule

- Entry point: `main.go` calls `services.StartDeviceAccOnKpiJob()`.
- Runtime behavior:
1. Run once immediately.
2. Then run every `5 minutes`.
- Important: Mongo is sampled every 5 minutes, but `device_kpi` write is throttled to every `1 hour` by cursor `device_kpi_hourly_insert`.

## 2) KPI Names Written To `device_kpi`

- `ACC_ON_OVER_10H`
- `ACC_OFF_OVER_72H`
- `VEHICLE_MKN_OVER_72H`
- `VEHICLE_MKN_CURRENT`
- `VEHICLE_BADGPS_OVER_ONLINE_IN_DAY`
- `ACC_ON_OVER_10H_CYCLE_21`
- `ACC_OFF_OVER_72H_CYCLE_21`
- `VEHICLE_MKN_OVER_72H_CYCLE_21`

## 3) Mongo Configuration

Config keys (`config/mongo.go`):
- `ATTR_MONGO_DB_NAME` (default `attributes`)
- `ATTR_MONGO_HOSTS` (default `172.21.5.146:27017,172.21.5.147:27017,172.21.5.148:27017`)
- `ATTR_MONGO_PORT` (default `27017`)
- `ATTR_MONGO_USER` (default `viot`)
- `ATTR_MONGO_PASS` (default `newpassword`)
- `ATTR_MONGO_AUTH_SOURCE` (default empty)
- `ATTR_MONGO_AUTH_MECHANISM` (default empty)

Connection behavior:
- Build URI from hosts + port.
- Try authSource candidates: configured -> dbName -> `admin`.
- Try auth mechanism candidates: configured -> auto(empty) -> `SCRAM-SHA-256` -> `SCRAM-SHA-1`.

## 4) KPI Time Config

Config keys (`config/kpi.go`):
- `KPI_TIMEZONE` (default `Asia/Bangkok`)
- `KPI_INCREMENTAL_LOOKBACK_SECONDS` (default `300`)
- `KPI_CYCLE_RESET_DAY` (default `21`)

## 5) Mongo Query Logic

## 5.1) Basic duration KPIs (5-minute sampling)

- `ACC_ON_OVER_10H`
  - `entity_type = VEHICLE`
  - `acc.value = active`
  - `acc.last_update_ts <= now - 10h`
- `ACC_OFF_OVER_72H`
  - `entity_type = VEHICLE`
  - `acc.value = inactive`
  - `acc.last_update_ts <= now - 72h`
- `VEHICLE_MKN_OVER_72H`
  - `entity_type = VEHICLE`
  - `status.value = offline`
  - `status.last_update_ts <= now - 72h`

Total for the above KPIs: `countDocuments({ entity_type: "VEHICLE" })`.

## 5.2) Daily status KPIs (incremental window)

- Source collection: `attributes`.
- Incremental range on `status.last_update_ts`:
  - Start at day start (`00:00` in KPI timezone), or cursor-based start with lookback.
  - End at current run time.
- Status buckets:
  - Offline set: `status.value = offline`
  - Bad GPS set: `status.value = badgps`
  - Online set: `status.value IN [badgps, overspeed, park, run, stop]`
- Use distinct `plateNo.value`.
- Upsert flags by `kpi_date + plate_no` into `device_kpi_daily_vehicle_status`:
  - `seen_offline`, `seen_badgps`, `seen_online`.
- Then count daily totals from that table:
  - `VEHICLE_MKN_CURRENT` = count `seen_offline = true`
  - `VEHICLE_BADGPS_OVER_ONLINE_IN_DAY` = count `seen_badgps = true`
  - denominator for both = count `seen_online = true`

## 5.3) Cycle-21 KPIs

- Window start = day `KPI_CYCLE_RESET_DAY` at `00:00` in KPI timezone.
- If today day-of-month <= reset day, start from previous month reset day.
- Distinct by `plateNo.value`.
- Conditions:
  - `ACC_ON_OVER_10H_CYCLE_21`: `acc.value=active`, `acc.last_update_ts` in `[cycleStartTs, now-10h]`
  - `ACC_OFF_OVER_72H_CYCLE_21`: `acc.value=inactive`, `acc.last_update_ts` in `[cycleStartTs, now-72h]`
  - `VEHICLE_MKN_OVER_72H_CYCLE_21`: `status.value=offline`, `status.last_update_ts` in `[cycleStartTs, now-72h]`
- Total cycle vehicles: distinct `plateNo.value` with `entity_type=VEHICLE`.

## 6) PostgreSQL Tables Used

Auto-created by code if missing:
- `device_kpi`
- `device_kpi_job_cursor`
- `device_kpi_daily_vehicle_status`
- `device_kpi_daily_vehicle_status_archive`
- `device_kpi_vehicle_detail`
- `device_kpi_vehicle_detail_snapshot`

### 6.1) `device_kpi`
- Snapshot KPI table (`name`, `matched_device_count`, `total_device_count`, `timestamps`).

### 6.2) `device_kpi_job_cursor`
- Generic cursor table:
  - `job_name` (PK), `last_processed_ts`, `updated_at`.

Cursor keys in this job:
- `device_daily_status_plate_snapshot`
- `device_daily_status_day_reset`
- `device_kpi_hourly_insert`
- `device_kpi_vehicle_detail_day_reset`

### 6.3) `device_kpi_daily_vehicle_status`
- PK: `(kpi_date, plate_no)`.
- Flags: `seen_offline`, `seen_badgps`, `seen_online`.

### 6.4) `device_kpi_daily_vehicle_status_archive`
- Daily archive snapshot table for monthly reporting.
- PK: `(kpi_date, plate_no)`.
- Archived on first run after day changes, before resetting daily flags.

### 6.5) `device_kpi_vehicle_detail`
- One row per `plate_no` (unique index on `plate_no`).
- Maintains current `in_*` flags and KPI-specific last update / violation timestamps:
  - `in_acc_on_over_10h`, `in_acc_off_over_72h`, `in_vehicle_mkn_over_72h`
  - `acc_on_*`, `acc_off_*`, `mkn_*`.

### 6.6) `device_kpi_vehicle_detail_snapshot`
- Daily snapshot table before new-day reset of detail flags.
- PK: `(snapshot_date, plate_no)`.
- Stores full detail state + cumulative counters:
  - `count_acc_on_over_10h`
  - `count_acc_off_over_72h`
  - `count_vehicle_mkn_over_72h`
- Counter increment rule per KPI:
1. Current `in_* = true`
2. Current `*_last_update_at` is not null
3. Current `*_last_update_at` differs from nearest previous snapshot of same `plate_no`
4. Then `count_* += 1`

## 7) Reset/Snapshot Behavior At New Day

Timezone basis: `KPI_TIMEZONE` (default Asia/Bangkok).

### 7.1) Daily status table
On first run of new day:
1. Snapshot previous `kpi_date` from `device_kpi_daily_vehicle_status` to archive.
2. Reset current day `seen_*` flags to false.
3. Update cursor `device_daily_status_day_reset`.

### 7.2) Vehicle detail table
On first hourly write path of new day:
1. Snapshot `device_kpi_vehicle_detail` into `device_kpi_vehicle_detail_snapshot` with `snapshot_date = previous day`.
2. Reset all `in_*` flags in `device_kpi_vehicle_detail` to false.
3. Update cursor `device_kpi_vehicle_detail_day_reset`.

## 8) Hourly Insert Gate

- `shouldInsertDeviceKpiSnapshot` checks cursor `device_kpi_hourly_insert`.
- If less than 1 hour from last insert:
  - skip all writes to `device_kpi` and `device_kpi_vehicle_detail`.
- Daily status sync still runs every 5 minutes (before hourly gate).

## 9) Re-Implementation Checklist On New Branch

1. Wire `services.StartDeviceAccOnKpiJob()` in `main.go`.
2. Port full file `services/device_acc_on_kpi_service.go`.
3. Port/verify model files:
   - `model/device_kpi.go`
   - `model/device_kpi_job_cursor.go`
   - `model/device_kpi_daily_status.go`
   - `model/device_kpi_vehicle_detail.go`
4. Port config files:
   - `config/mongo.go`
   - `config/kpi.go`
5. Confirm trans DB connection uses `transactions` DB (`ConnectTransDB`).
6. Ensure all `Ensure*Table()` calls remain before writes.
7. Verify cursor names are identical (to preserve job state).
8. Verify timezone and reset-day env values.

## 10) Smoke Test (After Deploy)

1. Confirm worker logs immediate run and 5-minute runs.
2. Confirm hourly gating message appears between hourly writes.
3. Confirm rows are inserted to `device_kpi` once/hour.
4. Confirm `device_kpi_daily_vehicle_status` updates every 5 minutes.
5. At day rollover, confirm:
   - daily status archive snapshot done
   - detail snapshot done
   - flags reset done
6. Confirm `device_kpi_vehicle_detail_snapshot.count_*` increases only when rule matches.


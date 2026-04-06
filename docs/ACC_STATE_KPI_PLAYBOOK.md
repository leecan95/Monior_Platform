# ACC State KPI Playbook

Mục tiêu: triển khai lại luồng tính `ACC_ON_OVER_10H` và `ACC_OFF_OVER_72H` theo **state continuity** (không phụ thuộc trực tiếp `acc.last_update_ts` bị refresh liên tục trong Mongo).

## 1) Khác biệt so với luồng cũ

- Luồng cũ:
  - Count ACC_ON/ACC_OFF trực tiếp từ Mongo bằng điều kiện `acc.last_update_ts <= threshold`.
- Luồng mới:
  - Đồng bộ trạng thái ACC theo xe vào PostgreSQL:
    - `acc_state_value`
    - `acc_state_since_ts`
    - `acc_state_since_at`
  - Chỉ đổi `acc_state_since_*` khi trạng thái ACC đổi (`active` <-> `inactive` / trạng thái khác).
  - Vi phạm:
    - `ACC_ON_OVER_10H`: `acc_state_value='active'` và `now - acc_state_since >= 10h`
    - `ACC_OFF_OVER_72H`: `acc_state_value='inactive'` và `now - acc_state_since >= 72h`

## 2) File thay đổi

- `model/device_kpi_vehicle_detail.go`
- `services/device_acc_on_kpi_service.go`

## 3) Thay đổi trong model

### 3.1 Cột mới trong `device_kpi_vehicle_detail`

- `acc_state_value TEXT NOT NULL DEFAULT ''`
- `acc_state_since_ts BIGINT NULL`
- `acc_state_since_at TIMESTAMPTZ NULL`

Code đã tự `ALTER TABLE ... ADD COLUMN IF NOT EXISTS` trong `EnsureDeviceKPIVehicleDetailTable`.

### 3.2 Hàm mới

- `SyncDeviceVehicleACCState(rows, sampledAtUTC, sampledHourUTC)`
  - Upsert trạng thái ACC theo `plate_no`.
  - Chỉ reset `acc_state_since_*` khi state đổi.
  - Refresh cờ:
    - `in_acc_on_over_10h`
    - `in_acc_off_over_72h`
  - Cập nhật:
    - `acc_on_last_update_ts/at`, `acc_on_last_violation_at`
    - `acc_off_last_update_ts/at`, `acc_off_last_violation_at`

- `CountDeviceVehicleACCViolations(sampledHourUTC)`
  - Trả về số xe vi phạm ACC_ON/ACC_OFF theo `sampled_hour`.

- `ListDeviceVehicleACCOnViolationRows(sampledHourUTC)`
- `ListDeviceVehicleACCOffViolationRows(sampledHourUTC)`
  - Trả danh sách chi tiết xe vi phạm để ghi vào `device_kpi_vehicle_detail` snapshot logic.

## 4) Thay đổi trong service

### 4.1 Luồng `runDeviceAccOnKpiSnapshot`

1. Lấy `timestamp` + `sampledHourUTC`.
2. Ensure bảng:
   - `device_kpi`
   - `device_kpi_vehicle_detail`
   - `device_kpi_vehicle_detail_snapshot`
3. Chạy reset đầu ngày detail (snapshot + reset flag) như cũ.
4. Đọc trạng thái ACC toàn bộ xe từ Mongo (`fetchVehicleACCStateSnapshot`).
5. Đồng bộ ACC state vào PostgreSQL (`SyncDeviceVehicleACCState`).
6. Count ACC_ON/ACC_OFF từ PostgreSQL (`CountDeviceVehicleACCViolations`).
7. Hourly gate qua thì ghi:
   - `device_kpi` rows
   - `device_kpi_vehicle_detail` rows:
     - ACC_ON/ACC_OFF lấy từ PostgreSQL state list
     - MKN vẫn lấy Mongo (`fetchVehicleKpiDetailSnapshot`)

### 4.2 Hàm mới trong service

- `fetchVehicleACCStateSnapshot(sampledAtUTC)`
  - Mongo aggregate trả:
    - `plateNo.value`
    - `imei.value`
    - `acc.value`
    - `acc.last_update_ts`

### 4.3 Điều chỉnh hàm cũ

- `fetchVehicleKpiDetailSnapshot` chỉ còn lấy MKN (`VEHICLE_MKN_OVER_72H`) từ Mongo.

## 5) SQL migrate thủ công (nếu môi trường không cho app tự ALTER)

```sql
ALTER TABLE device_kpi_vehicle_detail
  ADD COLUMN IF NOT EXISTS acc_state_value TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS acc_state_since_ts BIGINT NULL,
  ADD COLUMN IF NOT EXISTS acc_state_since_at TIMESTAMPTZ NULL;
```

## 6) Checklist port sang nhánh khác

1. Copy thay đổi từ 2 file:
   - `model/device_kpi_vehicle_detail.go`
   - `services/device_acc_on_kpi_service.go`
2. Chạy format:
```bash
gofmt -w model/device_kpi_vehicle_detail.go services/device_acc_on_kpi_service.go
```
3. Build/test:
```bash
GOCACHE=/tmp/go-build go test -mod=mod -vet=off ./model ./services
```
4. Nếu thiếu quyền DDL, chạy SQL migrate thủ công trước khi deploy.

## 7) Kỳ vọng sau deploy

- Xe `acc=active` liên tục sẽ vẫn lên `ACC_ON_OVER_10H` đúng sau 10h, dù Mongo refresh `acc.last_update_ts` mỗi vài giây.
- Xe `acc=inactive` liên tục sẽ lên `ACC_OFF_OVER_72H` đúng sau 72h.
- Không còn phụ thuộc trực tiếp vào việc `acc.last_update_ts` có bị cập nhật lại hay không.

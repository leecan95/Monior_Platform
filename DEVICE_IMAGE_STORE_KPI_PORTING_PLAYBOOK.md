# Device Image-Store KPI Porting Playbook

Mục tiêu: tổng hợp toàn bộ thay đổi để triển khai tương tự trên nhánh khác.

## 1) Phạm vi thay đổi

Luồng `StartDeviceAccOnKpiJob` được mở rộng thêm 3 phần:

1. Đồng bộ `device_kpi_online_camera_daily` theo chu kỳ **1 giờ/lần**.
2. Đồng bộ `device_kpi_image_store_violation_daily` theo chu kỳ **10 phút/lần**.
3. Snapshot tổng hợp theo ngày vào `device_kpi_image_store_daily_summary` khi qua ngày mới.

## 2) File cần mang sang nhánh mới

1. `model/database.go`
2. `services/device_acc_on_kpi_service.go`
3. `model/device_kpi_online_camera_daily.go`
4. `model/device_kpi_image_store_violation_daily.go`
5. `model/device_kpi_image_store_daily_summary.go`

## 3) Các thay đổi chính trong code

### 3.1 Kết nối DB mới

Trong `model/database.go` thêm:

1. `ConnectDevicesDB(cfg)` kết nối DB `devices`
2. `ConnectVtrackingDB(cfg)` kết nối DB `vtracking`

### 3.2 Cursor mới trong job

Trong `services/device_acc_on_kpi_service.go` thêm các job cursor name:

1. `device_online_camera_daily_sync_1h`
2. `device_image_store_violation_sync_10m`
3. `device_image_store_daily_summary_day_snapshot`

### 3.3 Luồng online + camera (1h/lần)

Hàm: `syncOnlineCameraDaily10m(...)` (tên hàm giữ nguyên, nhưng interval là 1h)

1. Ensure bảng `device_kpi_online_camera_daily`
2. Gate theo `shouldRunIntervalSnapshot(..., time.Hour)`
3. Lấy xe online trong ngày từ `device_kpi_daily_vehicle_status` (`seen_online = true`)
4. Join sang `devices.mv_vehicle_camera` theo `plate_no = license_plate` và `camera = true`
5. Replace theo ngày: xóa dữ liệu cũ `kpi_date` rồi insert snapshot mới
6. Cập nhật cursor `device_online_camera_daily_sync_1h`

### 3.4 Luồng vi phạm lưu ảnh (10m/lần)

Hàm: `syncImageStoreViolationDaily10m(...)`

1. Ensure bảng `device_kpi_image_store_violation_daily`
2. Gate theo `shouldRunIntervalSnapshot(..., 10*time.Minute)`
3. Lấy dữ liệu từ `vtracking.mv_vehicle_violation_full` theo `kpi_date` ngày hiện tại, `gap_violation_count > 0`
4. Replace theo ngày: xóa dữ liệu cũ `kpi_date` rồi insert snapshot mới
5. Cập nhật cursor `device_image_store_violation_sync_10m`

### 3.5 Snapshot summary qua ngày mới

Hàm: `snapshotImageStoreDailySummaryOnNewDay(...)`

1. Ensure bảng `device_kpi_image_store_daily_summary`
2. Xác định ngày bắt đầu mới theo timezone KPI
3. Mỗi ngày chỉ chạy 1 lần bằng cursor `device_image_store_daily_summary_day_snapshot`
4. Snapshot cho **ngày hôm trước**:
   - `online_device_count` = distinct `plate_no` từ `device_kpi_online_camera_daily`
   - `violation_device_count` = distinct `plate_no` từ `device_kpi_image_store_violation_daily`
   - `violation_percent` = `violation / online * 100`
5. Upsert theo `summary_date`

## 4) Schema bảng

### 4.1 `device_kpi_online_camera_daily`

```sql
CREATE TABLE IF NOT EXISTS device_kpi_online_camera_daily (
  kpi_date DATE NOT NULL,
  plate_no TEXT NOT NULL,
  license_plate TEXT NOT NULL,
  imei TEXT NOT NULL DEFAULT '',
  device_model_name TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (kpi_date, plate_no)
);
CREATE INDEX IF NOT EXISTS idx_device_kpi_online_camera_daily_kpi_date
ON device_kpi_online_camera_daily (kpi_date DESC);
```

### 4.2 `device_kpi_image_store_violation_daily`

```sql
CREATE TABLE IF NOT EXISTS device_kpi_image_store_violation_daily (
  kpi_date DATE NOT NULL,
  imei TEXT NOT NULL DEFAULT '',
  plate_no TEXT NOT NULL,
  max_gap BIGINT NOT NULL DEFAULT 0,
  gap_violation_count BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (kpi_date, plate_no)
);
CREATE INDEX IF NOT EXISTS idx_device_kpi_image_store_violation_daily_kpi_date
ON device_kpi_image_store_violation_daily (kpi_date DESC);
```

### 4.3 `device_kpi_image_store_daily_summary`

```sql
CREATE TABLE IF NOT EXISTS device_kpi_image_store_daily_summary (
  summary_date DATE PRIMARY KEY,
  online_device_count BIGINT NOT NULL DEFAULT 0,
  violation_device_count BIGINT NOT NULL DEFAULT 0,
  violation_percent NUMERIC(7,2) NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_device_kpi_image_store_daily_summary_date
ON device_kpi_image_store_daily_summary (summary_date DESC);
```

## 5) Cách hoạt động ghi dữ liệu

### 5.1 `device_kpi_online_camera_daily`

1. Chạy 1 giờ/lần
2. Replace theo `kpi_date` (xóa cũ + insert mới)

### 5.2 `device_kpi_image_store_violation_daily`

1. Chạy 10 phút/lần
2. Replace theo `kpi_date` (xóa cũ + insert mới)

### 5.3 `device_kpi_image_store_daily_summary`

1. Chạy 1 lần khi qua ngày mới
2. Upsert theo `summary_date`
3. Cùng `summary_date` thì ghi đè, ngày mới thì thêm dòng mới

## 6) Checklist triển khai sang nhánh mới

1. Copy 5 file trong mục 2.
2. Chạy format:
```bash
gofmt -w model/database.go \
  model/device_kpi_online_camera_daily.go \
  model/device_kpi_image_store_violation_daily.go \
  model/device_kpi_image_store_daily_summary.go \
  services/device_acc_on_kpi_service.go
```
3. Build/test:
```bash
GOCACHE=/tmp/go-build go test -mod=mod ./...
```
4. Nếu cần tạo bảng tay, chạy DDL mục 4.
5. Deploy và kiểm tra log:
   - online-camera sync mỗi 1h
   - image-store sync mỗi 10m
   - snapshot summary chạy khi qua ngày mới

## 7) Câu lệnh kiểm tra sau deploy

```sql
SELECT * FROM device_kpi_job_cursor
WHERE job_name IN (
  'device_online_camera_daily_sync_1h',
  'device_image_store_violation_sync_10m',
  'device_image_store_daily_summary_day_snapshot'
);
```

```sql
SELECT kpi_date, count(*) FROM device_kpi_online_camera_daily
GROUP BY kpi_date ORDER BY kpi_date DESC LIMIT 7;
```

```sql
SELECT kpi_date, count(*) FROM device_kpi_image_store_violation_daily
GROUP BY kpi_date ORDER BY kpi_date DESC LIMIT 7;
```

```sql
SELECT * FROM device_kpi_image_store_daily_summary
ORDER BY summary_date DESC LIMIT 15;
```

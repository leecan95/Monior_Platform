# VHT -> VTNET Full Sync (from current git status)

Thời điểm kiểm tra: nhánh hiện tại `vht`.

## 1) Tóm tắt thay đổi chức năng chính

### `services/device_acc_on_kpi_service.go`
- Thêm KPI mới: `BAD_WRITE_MEMORY`.
- Luồng `runDeviceAccOnKpiSnapshot`:
1. Không `defer db.Close()` nữa (tránh đóng singleton pool giữa các vòng job).
2. Đảm bảo luồng snapshot/reset của `device_kpi_vehicle_detail` chạy trước gate insert theo giờ.
- `syncDailyVehicleStatusAndGetCounts`:
1. Trả thêm `badMemCount`.
2. Upsert theo danh sách vehicle (`plate_no + imei`) thay vì chỉ plate.
- `fetchDailyVehicleStatusPlateNoSnapshot`:
1. Trả 4 nhóm: offline, badgps, bad_memory, online.
2. Filter bad memory dựa trên:
   - `memoryCardStatus.last_update_ts` trong incremental window.
   - `memoryCardStatus.value.nok_latest_ts` thuộc ngày hiện tại KPI.
- `resetDeviceKPIDailyStatusFlagsOnNewDay`:
1. Snapshot `kpi_date` ngày hôm trước sang archive.
2. Xóa runtime rows của đúng ngày vừa snapshot.
3. Reset cờ ngày mới.

### `model/device_kpi_daily_status.go`
- Bổ sung cột và logic:
1. `imei` cho runtime + archive.
2. `seen_bad_memory` cho runtime + archive.
- Upsert flags nhận struct vehicle (`plate_no`, `imei`), cập nhật `imei` nếu giá trị mới non-empty.
- Count flags trả 4 nhóm: offline, badgps, bad_memory, online.
- Snapshot archive copy đủ `imei` + `seen_bad_memory`.
- Thêm hàm xóa runtime theo ngày:
  - `DeleteDeviceKPIDailyStatusByDate(kpiDate string)`.

### `model/device_kpi_vehicle_detail.go`
- Snapshot count logic bổ sung case tăng counter theo `*_last_violation_at` (ngoài case `in_* + *_last_update_at`).
- Mục tiêu: không bỏ sót count khi vi phạm được nhận diện qua mốc violation thay vì update timestamp.

### `go.mod`, `go.sum`
- Bổ sung/cập nhật dependency để build ổn định cho luồng Mongo mới (bao gồm transitive deps của mongo driver).

### `main.go`
- Chỉ đổi chuỗi debug trong `GetCpuUsage()`:
  - `Monitor 10032026` -> `Monitor 18032026`

## 2) Khuyến nghị khi sync sang `vtnet`

### Nên mang sang
- `go.mod`
- `go.sum`
- `model/device_kpi_daily_status.go`
- `model/device_kpi_vehicle_detail.go`
- `services/device_acc_on_kpi_service.go`

### Cân nhắc
- `main.go` (chỉ là debug string, có thể bỏ)




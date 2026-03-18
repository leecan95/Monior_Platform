# VTNET Update: device_kpi_daily_vehicle_status

Mục tiêu: đồng bộ các thay đổi mới nhất của luồng daily status sang nhánh `vtnet`.

## Phạm vi thay đổi

### 1) Snapshot chỉ 1 ngày trước đó (không snapshot toàn bộ <= ngày)
- File: `model/device_kpi_daily_status.go`
- Hàm: `SnapshotDeviceKPIDailyStatus`
- Điều kiện snapshot:
```sql
WHERE kpi_date = $1::date
```

### 2) Sau snapshot, xóa runtime data của đúng ngày vừa snapshot
- File: `model/device_kpi_daily_status.go`
- Hàm mới: `DeleteDeviceKPIDailyStatusByDate(kpiDate string) (int64, error)`
- SQL:
```sql
DELETE FROM device_kpi_daily_vehicle_status
WHERE kpi_date = $1::date
```

### 3) Luồng reset ngày mới gọi snapshot -> delete -> reset flags
- File: `services/device_acc_on_kpi_service.go`
- Hàm: `resetDeviceKPIDailyStatusFlagsOnNewDay`
- Trình tự:
1. `SnapshotDeviceKPIDailyStatus(previousKpiDate, nowLocal.UTC())`
2. `DeleteDeviceKPIDailyStatusByDate(previousKpiDate)`
3. `ResetDeviceKPIDailyStatusFlags(kpiDate)`
4. `UpsertDeviceKPIJobCursor(...)`

## Các file cần mang sang nhánh vtnet
- `model/device_kpi_daily_status.go`
- `services/device_acc_on_kpi_service.go`

## Cách apply nhanh bằng patch

Từ nhánh hiện tại (đã có code đúng):
```bash
git diff -- model/device_kpi_daily_status.go services/device_acc_on_kpi_service.go > /tmp/vtnet_daily_status.patch
```

Sang nhánh `vtnet`:
```bash
git checkout vtnet
git apply /tmp/vtnet_daily_status.patch
gofmt -w model/device_kpi_daily_status.go services/device_acc_on_kpi_service.go
GOCACHE=/tmp/go-build go test -mod=mod -vet=off ./model ./services
```

## Kỳ vọng sau khi deploy
- Qua ngày mới:
  - snapshot `kpi_date` ngày hôm trước sang archive
  - runtime table `device_kpi_daily_vehicle_status` xóa các dòng của ngày hôm trước
  - không còn giữ backlog ngày cũ trong runtime table

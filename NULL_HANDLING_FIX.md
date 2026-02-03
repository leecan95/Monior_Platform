# Fix cho lỗi NULL handling trong Database Scan

## Vấn đề
Lỗi `"sql: Scan error on column index 1, name "latency_95th_percentile": converting NULL to string is unsupported"` xảy ra khi:

1. Materialized views trả về giá trị NULL cho cột `latency_95th_percentile`
2. Code cố gắng scan giá trị NULL vào biến `string` trong Go
3. Go không hỗ trợ scan giá trị NULL trực tiếp vào kiểu `string`

## Nguyên nhân
- Khi không có dữ liệu trong khoảng thời gian được query
- Khi tất cả giá trị latency đều NULL
- Khi materialized view chưa được refresh hoặc có lỗi trong quá trình tính toán

## Giải pháp đã áp dụng

### 1. Thay đổi kiểu dữ liệu scan
**Trước:**
```go
var count string
var latency string
var total string
var result string
if err := rows.Scan(&count, &latency, &total, &result); err != nil {
    fmt.Printf("Loi scan db ", err)
    return data, err
}
```

**Sau:**
```go
var count sql.NullString
var latency sql.NullString
var total sql.NullString
var result sql.NullString
if err := rows.Scan(&count, &latency, &total, &result); err != nil {
    fmt.Printf("Loi scan db %v", err)
    return data, err
}
```

### 2. Thêm helper function
```go
// Helper function để xử lý giá trị NULL từ database
func getStringValue(nullStr sql.NullString) string {
    if nullStr.Valid {
        return nullStr.String
    }
    return "0" // Trả về giá trị mặc định khi NULL
}
```

### 3. Cập nhật cách sử dụng
```go
data = config.LatencyKpi{
    Api:        "API Name",
    Total:      getStringValue(total),
    Count:      getStringValue(count),
    Percentile: getStringValue(latency),
    Result:     getStringValue(result),
}
```

### 4. Cải thiện parseFloat function
```go
func parseFloat(value string) float64 {
    // Xử lý trường hợp giá trị rỗng hoặc NULL
    if value == "" || value == "0" {
        return 0.0
    }
    
    result, err := strconv.ParseFloat(value, 64)
    if err != nil {
        fmt.Printf("❌ Lỗi chuyển đổi giá trị '%s': %v\n", value, err)
        return 0.0
    }
    return result
}
```

## Các function đã được cập nhật

1. `QueryImageLatencyOneDay()`
2. `QueryLoginLatencyOneDay()`
3. `QueryTrackingLatencyOneDay()`
4. `QueryDashboardLatencyOneDay()`
5. `QueryTotalLatencyOneDay()`
6. `QueryGetimageLatency()`
7. `QueryTrackingLatency()`
8. `QueryLoginLatency()`
9. `QueryReportLatency()`
10. `QueryLatency()`
11. `QueryLatencyPool()`

## Lợi ích

1. **Tránh crash**: Ứng dụng không bị crash khi gặp giá trị NULL
2. **Graceful handling**: Xử lý mượt mà các trường hợp edge case
3. **Consistent data**: Đảm bảo dữ liệu luôn có giá trị hợp lệ
4. **Better logging**: Cải thiện thông báo lỗi để debug dễ hơn

## Testing
Chạy file test để kiểm tra:
```bash
go run test_null_handling.go
```

## Khuyến nghị

1. **Monitor materialized views**: Đảm bảo các materialized view được refresh định kỳ
2. **Add data validation**: Thêm validation cho dữ liệu đầu vào
3. **Implement retry logic**: Thêm logic retry khi query thất bại
4. **Add alerting**: Thiết lập cảnh báo khi có quá nhiều giá trị NULL

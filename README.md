# Monior_Platform
Chương trình này sử dụng SNMP (Simple Network Management Protocol) để thu thập dữ liệu hệ thống, sau đó thực hiện tổng hợp và xử lý thông tin.

## Cài Đặt Hệ Thống

Trước khi chạy hệ thống, bạn cần cài đặt SNMP manager và SNMP agent trên hệ thống. Nếu bạn đang sử dụng Ubuntu, sử dụng các lệnh sau:

```bash
sudo apt update
sudo apt install snmp
sudo apt-get install snmpd
```

## Thông Tin MIB

Thông tin MIB (Management Information Base) có thể được tìm thấy tại:

```
/usr/local/share/snmp/mibs
```

## Cấu Hình

Hãy đảm bảo cấu hình chi tiết về máy chủ giám sát và các tuyến đường trong tệp `snmpd.conf`.

---

**Ghi chú**: Thay đổi đường dẫn thực tế nơi thông tin MIB được lưu trữ trên hệ thống của bạn. Điều chỉnh cấu hình theo cần thiết cho môi trường cụ thể của bạn.

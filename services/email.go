package services

import (
	"Monitor_Platform/config"
	"Monitor_Platform/model"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

func SendMail(subject, body string) error {
	from := config.Username
	to := config.ToEmail

	// Chuẩn bị nội dung email
	toEmails := strings.Join(config.ToEmail, ", ")
	msg := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"Subject: " + subject + "\r\n" +
		"From: " + from + "\r\n" +
		"To: " + toEmails + "\r\n\r\n" +
		body

	// Cấu hình TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Nếu server có chứng chỉ hợp lệ, nên để false
		ServerName:         config.SMTPServer,
	}

	// Kết nối đến SMTP server qua TLS
	conn, err := tls.Dial("tcp", config.SMTPServer+":"+config.SMTPPort, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS Dial Error: %v", err)
	}
	defer conn.Close()

	// Tạo client SMTP
	client, err := smtp.NewClient(conn, config.SMTPServer)
	if err != nil {
		return fmt.Errorf("SMTP Client Error: %v", err)
	}

	// Xác thực SMTP
	auth := smtp.PlainAuth("", config.Username, config.Password, config.SMTPServer)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("Auth Error: %v", err)
	}

	// Đặt thông tin người gửi và người nhận
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("Mail From Error: %v", err)
	}
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("Recipient Error: %v", err)
		}
	}

	// Gửi nội dung email
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("Data Error: %v", err)
	}
	_, err = wc.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("Write Error: %v", err)
	}
	err = wc.Close()
	if err != nil {
		return fmt.Errorf("Close Error: %v", err)
	}

	// Đóng kết nối SMTP
	client.Quit()

	fmt.Println("✅ Email đã được gửi thành công!")
	return nil
}

func ContentEmail() config.Mail {
	// Lấy ngày hôm nay
	today := time.Now()

	// Lấy ngày hôm qua bằng cách trừ đi 1 ngày
	yesterday := today.AddDate(0, 0, -1)

	// Định dạng ngày theo dd-MM-yyyy (hoặc yyyy-MM-dd tùy theo yêu cầu)
	formattedDate := yesterday.Format("02-01-2006")

	// Tiêu đề email
	subject := fmt.Sprintf("[Mail tự động]Báo cáo VTRACKING - Ngày %s", formattedDate)
	// Lấy kết quả
	//Giá trị Latency thực tế
	tracking, _ := model.GetTrackingLatencyOneDay()
	image, _ := model.GetImageLatencyOneDay()
	login, _ := model.GetLoginLatencyOneDay()
	dashboard, _ := model.GetDashboardLatencyOneDay()
	total, _ := model.GetTotalLatencyOneDay()

	// Lấy dữ liệu license
	licenseExpired, _ := model.GetLicenseExpiredTodayTotal()
	licenseNew, _ := model.GetLicenseNewTodayTotal()
	licenseNewMonth, _ := model.GetLicenseNewThisMonthTotal()
	licenseNewTotal, _ := model.GetLicenseNewTotal()
	licenseValid, _ := model.GetLicenseValidThisMonthTotal()

	// Ngưỡng đánh giá >=95%
	threshold := 95.0
	// Danh sách các thông số
	kpis := []config.LatencyKpi{tracking, image, login, dashboard, total}

	// Tạo bảng HTML động cho Latency
	rows := ""
	for _, kpi := range kpis {
		percert := parseFloat(kpi.Count) / parseFloat(kpi.Total) * 100
		evaluation := GetEvaluation(percert, threshold)
		rows += fmt.Sprintf(`
			<tr>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%.2f%%</td>
				<td>&gt;=%.2f%%</td>
				<td>%s</td>
			</tr>`, kpi.Api, kpi.Total, kpi.Count, kpi.Percentile, percert, threshold, evaluation)
	}

	// Tạo bảng HTML động cho License
	licenseRows := fmt.Sprintf(`
		<tr>
			<td>Số thuê bao hết hạn trong ngày</td>
			<td>%.0f</td>
		</tr>
		<tr>
			<td>Số thuê bao đăng ký mới trong ngày</td>
			<td>%.0f</td>
		</tr>
		<tr>
			<td>Số thuê bao đăng ký mới trong tháng</td>
			<td>%.0f</td>
		</tr>
		<tr>
			<td>Tổng số thuê bao đăng ký mới</td>
			<td>%.0f</td>
		</tr>
		<tr>
			<td>Số thuê bao còn hạn trong tháng</td>
			<td>%.0f</td>
		</tr>`,
		licenseExpired.Value,
		licenseNew.Value,
		licenseNewMonth.Value,
		licenseNewTotal.Value,
		licenseValid.Value,
	)

	// Format toàn bộ HTML
	htmlBody := fmt.Sprintf(`
		<html>
		<head>
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<style>
				body {
					font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
					font-size: 14px;
					line-height: 1.6;
					color: #333;
					margin: 0;
					padding: 10px;
				}
				table {
					width: 100%%;
					border-collapse: collapse;
					margin-bottom: 20px;
					table-layout: fixed;
				}
				th, td {
					border: 1px solid #ddd;
					padding: 6px;
					text-align: center;
					font-size: 12px;
					overflow: hidden;
					text-overflow: ellipsis;
				}
				th {
					background-color: #f2f2f2;
					font-weight: bold;
				}
				.business-stats {
					width: 100%%;
				}
				.business-stats td:first-child {
					width: 70%%;
					text-align: left;
				}
				.business-stats td:last-child {
					width: 30%%;
					text-align: center;
				}
				h3 {
					color: #2c3e50;
					margin-top: 15px;
					margin-bottom: 10px;
					font-size: 16px;
				}
				.table-container {
					width: 100%%;
				}
				@media screen and (max-width: 480px) {
					body {
						padding: 5px;
					}
					th, td {
						padding: 4px;
						font-size: 11px;
						white-space: nowrap;
						min-width: 60px;
					}
					.table-container {
						overflow-x: auto;
						-webkit-overflow-scrolling: touch;
					}
					.kpi-table {
						width: 600px;
					}
					.kpi-table th:first-child,
					.kpi-table td:first-child {
						position: sticky;
						left: 0;
						background-color: #f9f9f9;
						z-index: 1;
						border-right: 2px solid #ccc;
						white-space: normal; /* Allow text wrapping in first column */
						width: 80px; /* Set fixed width for first column */
						min-width: 80px;
						word-break: break-word; /* Break words if needed */
					}
				}
			</style>
		</head>
		<body>
			<h3>Thống kê KPI trải nghiệm khách hàng</h3>
			<div class="table-container">
				<table class="kpi-table">
					<tr>
						<th>Chỉ số</th>
						<th>Tổng request</th>
						<th>Request < 5s</th>
						<th>Độ trễ trung bình 95%%</th>
						<th>Tỉ lệ</th>
						<th>Ngưỡng</th>
						<th>Đánh giá</th>
					</tr>
					%s
				</table>
			</div>

			<h3>Thống kê chỉ số kinh doanh</h3>
			<table class="business-stats">
				<tr>
					<th>Chỉ số</th>
					<th>Giá trị</th>
				</tr>
				%s
			</table>
		</body>
		</html>`, rows, licenseRows)

	return config.Mail{
		Subject: subject,
		Body:    htmlBody,
	}
}

// Hàm kiểm tra đánh giá dựa trên giá trị
func GetEvaluation(value float64, threshold float64) string {
	if value >= threshold {
		return "Đạt"
	}
	return "<span style='color: red; font-weight: bold;'>Không đạt</span>"
}
func parseFloat(value string) float64 {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Println("❌ Lỗi chuyển đổi:", value)
		return 0.0
	}
	return result
}

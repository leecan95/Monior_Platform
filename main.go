package main

import (
	"Monitor_Platform/config"
	"Monitor_Platform/model"
	"Monitor_Platform/routes"
	"Monitor_Platform/services"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata"
)

func main() {
	// Kiểm tra cấu hình DB bắt buộc qua biến môi trường
	if err := config.ValidateDBConfig(model.LoadDBConfig()); err != nil {
		fmt.Printf("Database config error: %v\n", err)
		return
	}

	// Worker ghi log API (POST /ems/kpi/logs)
	services.StartApiLogWorker()
	// Worker ghi daily api kpi (POST /ems/kpi/daily)
	services.StartDailyApiKpiWorker()
	// Worker tính availability hệ thống hằng ngày lúc 09:59 UTC
	services.StartSystemAvailabilityJob()
	// Worker lấy trạng thái replicas deployment mỗi 10s
	services.StartDeploymentReplicaStatusJob()
	// Worker lấy trạng thái ready nginx-ingress mỗi 30s
	services.StartNginxIngressReadyStatusJob()
	// Worker query mẫu KPI thiết bị mỗi 5 phút, nhưng ghi vào device_kpi theo chu kỳ 1 giờ
	services.StartDeviceAccOnKpiJob()

	r := routes.SetupRouter()
	srv := &http.Server{
		Addr:    ":8933",
		Handler: r,
	}

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Demo loop cũ
	go func() {
		for {
			GetCpuUsage()
			time.Sleep(60 * time.Second)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctx)
}

func GetCpuUsage() {
	fmt.Print("Monitor 10032026\n")
}

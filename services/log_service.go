package services

import (
	"Monitor_Platform/config"
	"Monitor_Platform/model"
	"context"
	"log"
	"time"
)

// Buffer for incoming API logs.
var (
	apiLogQueue = make(chan model.ApiLog, 20000)
	batchSize   = 1000
	flushEvery  = 5 * time.Second

	// queue for daily_api_kpi inserts
	dailyApiKpiQueue   = make(chan model.DailyApiKpi, 5000)
	dailyBatchSize     = 500
	dailyFlushInterval = 10 * time.Second
)

// EnqueueApiLog pushes a log item into queue; drops if full.
func EnqueueApiLog(logItem model.ApiLog) {
	select {
	case apiLogQueue <- logItem:
	default:
		log.Println("[WARN] api log queue full, dropping log")
	}
}

// StartApiLogWorker launches a goroutine that flushes logs in batches.
func StartApiLogWorker() {
	go func() {
		buffer := make([]model.ApiLog, 0, batchSize)
		ticker := time.NewTicker(flushEvery)
		for {
			select {
			case logItem := <-apiLogQueue:
				buffer = append(buffer, logItem)
				if len(buffer) >= batchSize {
					flushApiLogs(buffer)
					buffer = buffer[:0]
				}
			case <-ticker.C:
				if len(buffer) > 0 {
					flushApiLogs(buffer)
					buffer = buffer[:0]
				}
			}
		}
	}()
}

// EnqueueDailyApiKpi pushes a KPI row into queue; drops if full.
func EnqueueDailyApiKpi(row model.DailyApiKpi) {
	select {
	case dailyApiKpiQueue <- row:
	default:
		log.Println("[WARN] daily api kpi queue full, dropping row")
	}
}

// StartDailyApiKpiWorker flushes KPI rows in batches to daily_api_kpi table.
func StartDailyApiKpiWorker() {
	go func() {
		buffer := make([]model.DailyApiKpi, 0, dailyBatchSize)
		ticker := time.NewTicker(dailyFlushInterval)
		for {
			select {
			case item := <-dailyApiKpiQueue:
				buffer = append(buffer, item)
				if len(buffer) >= dailyBatchSize {
					flushDailyApiKpi(buffer)
					buffer = buffer[:0]
				}
			case <-ticker.C:
				if len(buffer) > 0 {
					flushDailyApiKpi(buffer)
					buffer = buffer[:0]
				}
			}
		}
	}()
}

func flushDailyApiKpi(rows []model.DailyApiKpi) {
	cfg := model.LoadDBConfig()
	if err := config.ValidateDBConfig(cfg); err != nil {
		log.Printf("[ERROR] daily api kpi flush skipped: %v", err)
		return
	}
	db, err := model.ConnectTransDB(cfg)
	if err != nil {
		log.Printf("[ERROR] connect trans db: %v", err)
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = ctx

	if err := db.BatchInsertDailyApiKpi(rows); err != nil {
		log.Printf("[ERROR] batch insert daily api kpi failed: %v", err)
	}
}

func flushApiLogs(logs []model.ApiLog) {
	cfg := model.LoadDBConfig()
	if err := config.ValidateDBConfig(cfg); err != nil {
		log.Printf("[ERROR] api log flush skipped: %v", err)
		return
	}
	db, err := model.ConnectTransDB(cfg)
	if err != nil {
		log.Printf("[ERROR] connect trans db: %v", err)
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = ctx // context reserved if we extend ConnectTransDB; BatchInsert uses its own ctx.

	if err := db.BatchInsertApiLogs(logs); err != nil {
		log.Printf("[ERROR] batch insert api logs failed: %v", err)
	}
}

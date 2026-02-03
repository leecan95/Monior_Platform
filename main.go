package main

import (
	"Monitor_Platform/config"
	"Monitor_Platform/routes"
	"Monitor_Platform/services"
	"fmt"
	"time"
)

func main() {
	// Initialize email configuration
	if err := config.LoadEmailConfig(); err != nil {
		fmt.Printf("Failed to load email config: %v\n", err)
		return
	}

	//var cfg config.DBConfig
	//cfg = model.LoadDBConfig()
	//db := model.ConnectToDb()
	//db, _ := model.ConnectPoolDB(cfg)
	//defer db.Close()
	//db2, _ := model.ConnectTransDB(cfg)
	go func() {
		for {
			GetCpuUsage()
			//services.MonitorKpiApi(db)
			//db.QueryData()
			//db2.QueryLatency()
			time.Sleep(60 * time.Second)
		}
	}()
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, time.Local)
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}
			sleepDuration := next.Sub(now)
			fmt.Printf("Next email will be sent at %s (sleeping for %s)\n", next.Format("2006-01-02 15:04:05"), sleepDuration)
			time.Sleep(sleepDuration)
			fmt.Printf("Send mail\n")
			mail := services.ContentEmail()
			err := services.SendMail(mail.Subject, mail.Body)
			if err != nil {
				fmt.Printf("Send mail error %s \n", err)
			}			
		}
	}()
	services.StartApiLogWorker()
	r := routes.SetupRouter()
	err := r.Run(":8933")
	if err != nil {
		panic(err)
	}
	select {}
}

func GetCpuUsage() {
	fmt.Print("Monitor 03022026\n")
}

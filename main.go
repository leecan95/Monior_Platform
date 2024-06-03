package main

import (
	"Monitor_Platform/config"
	"Monitor_Platform/model"
	"Monitor_Platform/routes"
	"fmt"
	"time"
)

func main() {
	var cfg config.DBConfig
	cfg = model.LoadDBConfig()
	//db := model.ConnectToDb()
	db, _ := model.ConnectPoolDB(cfg)
	defer db.Close()
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
			fmt.Print("08042024\n")
			time.Sleep(5 * time.Second)
		}

	}()
	r := routes.SetupRouter(db)
	err := r.Run(":8933")
	if err != nil {
		panic(err)
	}
	select {}

}
func GetCpuUsage() {
	fmt.Print("Monitor 08042024\n")
}

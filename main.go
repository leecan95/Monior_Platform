package main

import (
	"Monitor_Platform/routes"
	"fmt"
	"time"
)

func main() {
	//model.LoadDBConfig()
	//db := model.ConnectToDb()
	//defer db.Close()
	go func() {
		for {
			GetCpuUsage()
			//	services.MonitorKpiApi(db)
			time.Sleep(30 * time.Second)
		}
	}()
	go func() {
		for {
			fmt.Print("14012024\n")
			time.Sleep(5 * time.Second)
		}

	}()
	r := routes.SetupRouter()

	err := r.Run(":8933")
	if err != nil {
		panic(err)
	}
	select {}

}
func GetCpuUsage() {
	fmt.Print("Monitor 14012024\n")
}

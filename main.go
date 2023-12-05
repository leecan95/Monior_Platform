package main

import (
	"Monitor_Platform/model"
	"Monitor_Platform/routes"
	"Monitor_Platform/services"
	"fmt"
	"time"
)

func main() {
	model.LoadDBConfig()
	db := model.ConnectToDb()
	defer db.Close()
	go func() {
		for {
			GetCpuUsage()
			services.MonitorKpiApi(db)
			time.Sleep(15 * time.Second)
		}
	}()
	go func() {
		for {
			fmt.Print("go routine 2\n")
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
	fmt.Print("Test go routine\n")
}

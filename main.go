package main

import (
	"Monitor_Platform/routes"
	"Monitor_Platform/snmp"
	"log"
)

func main() {

	// Kết nối tới SNMP agent
	err := snmp.Manager.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer snmp.Manager.Conn.Close()

	r := routes.SetupRouter()

	err = r.Run("localhost:8080")
	if err != nil {
		panic(err)
	}

}

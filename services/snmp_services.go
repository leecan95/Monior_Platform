package services

import (
	"Monitor_Platform/config"
	"Monitor_Platform/snmp"
	"github.com/gin-gonic/gin"
	"log"
)

func GetValueFromOID(c *gin.Context, oid []string) interface{} {
	// Kết nối tới SNMP agent
	err := snmp.Manager.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer snmp.Manager.Conn.Close()

	result, err2 := snmp.Manager.Get(oid) // Lấy dữ liệu từ SNMP agent
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}
	return result
}

func GetAllCpuUsage(c *gin.Context) interface{} {
	var sum []interface{}
	result := GetDataSNMP(config.Server244, config.MemTotalReal)
	sum = append(sum, result)
	result = GetDataSNMP(config.Server245, config.MemTotalReal)
	sum = append(sum, result)
	result = GetDataSNMP(config.Server246, config.MemTotalReal)
	sum = append(sum, result)
	result = GetDataSNMP(config.Server250, config.MemTotalReal)
	sum = append(sum, result)
	result = GetDataSNMP(config.Server251, config.MemTotalReal)
	sum = append(sum, result)
	result = GetDataSNMP(config.Server252, config.MemTotalReal)
	sum = append(sum, result)
	return sum
}

func GetDataSNMP(target string, oid string) interface{} {

	var oids []string
	oids = append(oids, oid)
	snmp.Manager.Target = target

	err := snmp.Manager.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer snmp.Manager.Conn.Close()

	result, err2 := snmp.Manager.Get(oids) // Lấy dữ liệu từ SNMP agent
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}

	return result
}

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
		log.Printf("Connect() err: %v", err)
		c.JSON(500, "connect to snmp failed")
	}
	defer snmp.Manager.Conn.Close()

	result, err2 := snmp.Manager.Get(oid) // Lấy dữ liệu từ SNMP agent
	if err2 != nil {
		log.Printf("Get() err: %v", err2)
		c.JSON(500, "Get snmp failed")
	}
	return result
}

func GetAllCpuUsage(c *gin.Context) interface{} {
	var sum []interface{}
	result := GetDataSNMP(c, config.Server244, config.MemTotalReal)
	sum = append(sum, result)
	result = GetDataSNMP(c, config.Server245, config.MemTotalReal)
	sum = append(sum, result)
	result = GetDataSNMP(c, config.Server246, config.MemTotalReal)
	sum = append(sum, result)
	result = GetDataSNMP(c, config.Server250, config.MemTotalReal)
	sum = append(sum, result)
	result = GetDataSNMP(c, config.Server251, config.MemTotalReal)
	sum = append(sum, result)
	result = GetDataSNMP(c, config.Server252, config.MemTotalReal)
	sum = append(sum, result)
	return sum
}

func GetDataSNMP(c *gin.Context, target string, oid string) interface{} {

	var oids []string
	oids = append(oids, oid)
	snmp.Manager.Target = target

	err := snmp.Manager.Connect()
	if err != nil {
		log.Printf("Connect() err: %v", err)
		c.JSON(500, "connect SNMP failed")
	}
	defer snmp.Manager.Conn.Close()

	result, err2 := snmp.Manager.Get(oids) // Lấy dữ liệu từ SNMP agent
	if err2 != nil {
		log.Printf("Get() err: %v", err2)
		c.JSON(500, "get SNMP failed")
	}

	return result
}

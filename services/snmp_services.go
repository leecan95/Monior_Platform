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

func GetAllMemInfo(c *gin.Context) []int32 {
	var sum []int32
	val := GetDataSNMP(c, config.Server244, config.MemTotalReal)
	result, _ := val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.MemTotalReal)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.MemTotalReal)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.MemTotalReal)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.MemTotalReal)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.MemTotalReal)
	result, _ = val.(int32)
	sum = append(sum, result)
	return sum
}

func GetAvailMemInfo(c *gin.Context) []int32 {
	var sum []int32
	val := GetDataSNMP(c, config.Server244, config.MemAvailReal)
	result, _ := val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.MemAvailReal)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.MemAvailReal)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.MemAvailReal)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.MemAvailReal)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.MemAvailReal)
	result, _ = val.(int32)
	sum = append(sum, result)
	return sum
}

func GetCpuRawUser(c *gin.Context) []int32 {
	var sum []int32
	val := GetDataSNMP(c, config.Server244, config.SsCpuRawUser)
	result, _ := val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.SsCpuRawUser)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.SsCpuRawUser)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.SsCpuRawUser)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.SsCpuRawUser)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.SsCpuRawUser)
	result, _ = val.(int32)
	sum = append(sum, result)
	return sum
}

func GetCpuRawNice(c *gin.Context) []int32 {
	var sum []int32
	val := GetDataSNMP(c, config.Server244, config.SsCpuRawNice)
	result, _ := val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.SsCpuRawNice)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.SsCpuRawNice)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.SsCpuRawNice)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.SsCpuRawNice)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.SsCpuRawNice)
	result, _ = val.(int32)
	sum = append(sum, result)
	return sum
}

func GetCpuRawSystem(c *gin.Context) []int32 {
	var sum []int32
	val := GetDataSNMP(c, config.Server244, config.SsCpuRawSystem)
	result, _ := val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.SsCpuRawSystem)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.SsCpuRawSystem)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.SsCpuRawSystem)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.SsCpuRawSystem)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.SsCpuRawSystem)
	result, _ = val.(int32)
	sum = append(sum, result)
	return sum
}

func GetCpuRawIdle(c *gin.Context) []int32 {
	var sum []int32
	val := GetDataSNMP(c, config.Server244, config.SsCpuRawIdle)
	result, _ := val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.SsCpuRawIdle)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.SsCpuRawIdle)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.SsCpuRawIdle)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.SsCpuRawIdle)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.SsCpuRawIdle)
	result, _ = val.(int32)
	sum = append(sum, result)
	return sum
}

func GetCpuRawWait(c *gin.Context) []int32 {
	var sum []int32
	val := GetDataSNMP(c, config.Server244, config.SsCpuRawWait)
	result, _ := val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.SsCpuRawWait)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.SsCpuRawWait)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.SsCpuRawWait)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.SsCpuRawWait)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.SsCpuRawWait)
	result, _ = val.(int32)
	sum = append(sum, result)
	return sum
}

func GetCpuRawKernel(c *gin.Context) []int32 {
	var sum []int32
	val := GetDataSNMP(c, config.Server244, config.SsCpuRawKernel)
	result, _ := val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.SsCpuRawKernel)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.SsCpuRawKernel)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.SsCpuRawKernel)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.SsCpuRawKernel)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.SsCpuRawKernel)
	result, _ = val.(int32)
	sum = append(sum, result)
	return sum
}

func GetCpuRawInterrupt(c *gin.Context) []int32 {
	var sum []int32
	val := GetDataSNMP(c, config.Server244, config.SsCpuRawInterrupt)
	result, _ := val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.SsCpuRawInterrupt)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.SsCpuRawInterrupt)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.SsCpuRawInterrupt)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.SsCpuRawInterrupt)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.SsCpuRawInterrupt)
	result, _ = val.(int32)
	sum = append(sum, result)
	return sum
}

func GetCpuRawSoftIrq(c *gin.Context) []int32 {
	var sum []int32
	val := GetDataSNMP(c, config.Server244, config.SsCpuRawSoftIRQ)
	result, _ := val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.SsCpuRawSoftIRQ)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.SsCpuRawSoftIRQ)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.SsCpuRawSoftIRQ)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.SsCpuRawSoftIRQ)
	result, _ = val.(int32)
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.SsCpuRawSoftIRQ)
	result, _ = val.(int32)
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

	return result.Variables[0].Value
}

package services

import (
	"Monitor_Platform/config"
	"Monitor_Platform/snmp"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
)

var mu sync.Mutex

func GetValueFromOID(c *gin.Context, oid []string) interface{} {
	// Kết nối tới SNMP agent
	err := snmp.Manager.Connect()
	if err != nil {
		log.Printf("Connect() err: %v", err)
		c.Error(err)
	}
	defer snmp.Manager.Conn.Close()

	result, err2 := snmp.Manager.Get(oid) // Lấy dữ liệu từ SNMP agent
	if err2 != nil {
		log.Printf("Get() err: %v", err2)
		c.Error(err)
	}
	return result
}

func GetAllMemInfo(c *gin.Context) []int64 {
	var sum []int64
	val := GetDataSNMP(c, config.Server244, config.MemTotalReal)
	result := int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.MemTotalReal)
	result = int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.MemTotalReal)
	result = int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.MemTotalReal)
	result = int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.MemTotalReal)
	result = int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.MemTotalReal)
	result = int64(val.(int))
	sum = append(sum, result)
	return sum
}

func GetAvailMemInfo(c *gin.Context) []int64 {
	var sum []int64
	val := GetDataSNMP(c, config.Server244, config.MemAvailReal)
	result := int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.MemAvailReal)
	result = int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.MemAvailReal)
	result = int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.MemAvailReal)
	result = int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.MemAvailReal)
	result = int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.MemAvailReal)
	result = int64(val.(int))
	sum = append(sum, result)
	return sum
}

func GetCpuRawUser(c *gin.Context) []int64 {
	var sum []int64
	val := GetDataSNMP(c, config.Server244, config.SsCpuRawUser)
	result := int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.SsCpuRawUser)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.SsCpuRawUser)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.SsCpuRawUser)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.SsCpuRawUser)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.SsCpuRawUser)
	result = int64(val.(uint))
	sum = append(sum, result)
	return sum
}

func GetCpuRawNice(c *gin.Context) []int64 {
	var sum []int64
	val := GetDataSNMP(c, config.Server244, config.SsCpuRawNice)
	fmt.Printf("102 gia tri %v\n", val)
	result := int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.SsCpuRawNice)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.SsCpuRawNice)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.SsCpuRawNice)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.SsCpuRawNice)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.SsCpuRawNice)
	result = int64(val.(uint))
	sum = append(sum, result)
	return sum
}

func GetCpuRawSystem(c *gin.Context) []int64 {
	var sum []int64
	val := GetDataSNMP(c, config.Server244, config.SsCpuRawSystem)
	fmt.Printf("126 gia tri %v\n", val)
	result := int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.SsCpuRawSystem)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.SsCpuRawSystem)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.SsCpuRawSystem)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.SsCpuRawSystem)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.SsCpuRawSystem)
	result = int64(val.(uint))
	sum = append(sum, result)
	return sum
}

func GetCpuRawIdle(c *gin.Context) []int64 {
	var sum []int64
	val := GetDataSNMP(c, config.Server244, config.SsCpuRawIdle)
	fmt.Printf("150 gia tri %v\n", val)
	result := int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.SsCpuRawIdle)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.SsCpuRawIdle)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.SsCpuRawIdle)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.SsCpuRawIdle)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.SsCpuRawIdle)
	result = int64(val.(uint))
	sum = append(sum, result)
	return sum
}

func GetCpuRawWait(c *gin.Context) []int64 {
	var sum []int64
	val := GetDataSNMP(c, config.Server244, config.SsCpuRawWait)
	fmt.Printf("174 gia tri %v\n", val)
	result := int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.SsCpuRawWait)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.SsCpuRawWait)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.SsCpuRawWait)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.SsCpuRawWait)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.SsCpuRawWait)
	result = int64(val.(uint))
	sum = append(sum, result)
	return sum
}

func GetCpuRawKernel(c *gin.Context) []int64 {
	var sum []int64
	val := GetDataSNMP(c, config.Server244, config.SsCpuRawKernel)
	fmt.Printf("198 gia tri %v\n", val)
	result := int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.SsCpuRawKernel)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.SsCpuRawKernel)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.SsCpuRawKernel)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.SsCpuRawKernel)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.SsCpuRawKernel)
	result = int64(val.(uint))
	sum = append(sum, result)
	return sum
}

func GetCpuRawInterrupt(c *gin.Context) []int64 {
	var sum []int64
	val := GetDataSNMP(c, config.Server244, config.SsCpuRawInterrupt)
	fmt.Printf("222 gia tri %v\n", val)
	result := int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.SsCpuRawInterrupt)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.SsCpuRawInterrupt)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.SsCpuRawInterrupt)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.SsCpuRawInterrupt)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.SsCpuRawInterrupt)
	result = int64(val.(uint))
	sum = append(sum, result)
	return sum
}

func GetCpuRawSoftIrq(c *gin.Context) []int64 {
	var sum []int64
	val := GetDataSNMP(c, config.Server244, config.SsCpuRawSoftIRQ)
	fmt.Printf("246 gia tri %v\n", val)
	result := int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.SsCpuRawSoftIRQ)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.SsCpuRawSoftIRQ)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.SsCpuRawSoftIRQ)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.SsCpuRawSoftIRQ)
	result = int64(val.(uint))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.SsCpuRawSoftIRQ)
	result = int64(val.(uint))
	sum = append(sum, result)
	return sum
}

func GetDiskTotal(c *gin.Context) []int64 {
	var sum []int64
	val := GetDataSNMP(c, config.Server244, config.DskTotal)
	result := int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.DskTotal)
	result = int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.DskTotal)
	result = int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.DskTotal)
	result = int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.DskTotal)
	result = int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.DskTotal)
	result = int64(val.(int))
	sum = append(sum, result)
	return sum
}

func GetDiskAvail(c *gin.Context) []int64 {
	var sum []int64
	val := GetDataSNMP(c, config.Server244, config.DskAvail)
	result := int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server245, config.DskAvail)
	result = int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server246, config.DskAvail)
	result = int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server250, config.DskAvail)
	result = int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server251, config.DskAvail)
	result = int64(val.(int))
	sum = append(sum, result)
	val = GetDataSNMP(c, config.Server252, config.DskAvail)
	result = int64(val.(int))
	sum = append(sum, result)
	return sum
}

func GetDataSNMP(c *gin.Context, target string, oid string) interface{} {
	mu.Lock()
	defer mu.Unlock()
	var oids []string
	oids = append(oids, oid)
	snmp.Manager.Target = target

	err := snmp.Manager.Connect()
	if err != nil {
		log.Printf("Connect() err: %v", err)
		c.Error(err)
	}
	defer snmp.Manager.Conn.Close()

	result, err2 := snmp.Manager.Get(oids) // Lấy dữ liệu từ SNMP agent
	if err2 != nil {
		log.Printf("Get() err: %v", err2)
		c.Error(err)
	}

	return result.Variables[0].Value
}

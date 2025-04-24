package controllers

import (
	"Monitor_Platform/config"
	"Monitor_Platform/services"
	"Monitor_Platform/validations"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

var previousIdle = make([]int64, 6)
var previousUser = make([]int64, 6)

func GetOIDController(c *gin.Context) {
	var request validations.GetOidSnmp
	if requestErr := c.ShouldBind(&request); requestErr != nil {
		errRes := strings.Split(requestErr.Error(), ": ")
		services.ValidationError(c, "These fields are required!", errRes)
		return
	}
	// Lấy danh sách các OID từ query parameters
	oids := c.QueryArray("oid")
	// Thực hiện lấy giá trị tương ứng với từng OID
	values := services.GetValueFromOID(c, oids)

	c.JSON(200, values)

}

func GetAllMemController(c *gin.Context) {
	values, svs := services.GetAllMemInfo(c)
	var data []MemSnmpAll
	for i, server := range svs {
		if i < len(values) {
			data = append(data, MemSnmpAll{Server: server, Value: values[i]})
		}
	}
	c.JSON(200, data)
}

func GetAvailMemController(c *gin.Context) {
	values, svs := services.GetAvailMemInfo(c)
	var data []MemSnmpAll
	for i, server := range svs {
		if i < len(values) {
			data = append(data, MemSnmpAll{Server: server, Value: values[i]})
		}
	}
	c.JSON(200, data)
}

func GetUsedMemController(c *gin.Context) {
	total, svs := services.GetAllMemInfo(c)
	avail, _ := services.GetAvailMemInfo(c)
	cache, _ := services.GetMemCacheInfo(c)
	buffer, _ := services.GetMemBufferInfo(c)
	var data []MemSnmpAll
	for i, server := range svs {
		if i < len(total) {
			data = append(data, MemSnmpAll{Server: server, Value: total[i] - (avail[i] + cache[i] + buffer[i])})
		}
	}
	//data := []MemSnmpAll{
	//	{Server: config.Server100, Value: total[0] - (avail[0] + cache[0] + buffer[0])},
	//	{Server: config.Server101, Value: total[1] - (avail[1] + cache[1] + buffer[1])},
	//	{Server: config.Server102, Value: total[2] - (avail[2] + cache[2] + buffer[2])},
	//	{Server: config.Server104, Value: total[3] - (avail[3] + cache[3] + buffer[3])},
	//	{Server: config.Server105, Value: total[4] - (avail[4] + cache[4] + buffer[4])},
	//	{Server: config.Server106, Value: total[5] - (avail[5] + cache[5] + buffer[5])},
	//	{Server: config.Server108, Value: total[6] - (avail[6] + cache[6] + buffer[6])},
	//	{Server: config.Server109, Value: total[7] - (avail[7] + cache[7] + buffer[7])},
	//	{Server: config.Server110, Value: total[8] - (avail[8] + cache[8] + buffer[8])},
	//	{Server: config.Server112, Value: total[9] - (avail[9] + cache[9] + buffer[9])},
	//	{Server: config.Server113, Value: total[10] - (avail[10] + cache[10] + buffer[10])},
	//	{Server: config.Server114, Value: total[11] - (avail[11] + cache[11] + buffer[11])},
	//	{Server: config.Server116, Value: total[12] - (avail[12] + cache[12] + buffer[12])},
	//	{Server: config.Server117, Value: total[13] - (avail[13] + cache[13] + buffer[13])},
	//	{Server: config.Server118, Value: total[14] - (avail[14] + cache[14] + buffer[14])},
	//	{Server: config.Server119, Value: total[15] - (avail[15] + cache[15] + buffer[15])},
	//	{Server: config.Server120, Value: total[16] - (avail[16] + cache[16] + buffer[16])},
	//	{Server: config.Server121, Value: total[17] - (avail[17] + cache[17] + buffer[17])},
	//	{Server: config.Server123, Value: total[18] - (avail[18] + cache[18] + buffer[18])},
	//	{Server: config.Server124, Value: total[19] - (avail[19] + cache[19] + buffer[19])},
	//	{Server: config.Server125, Value: total[20] - (avail[20] + cache[20] + buffer[20])},
	//	{Server: config.Server127, Value: total[21] - (avail[21] + cache[21] + buffer[21])},
	//	{Server: config.Server128, Value: total[22] - (avail[22] + cache[22] + buffer[22])},
	//	{Server: config.Server129, Value: total[23] - (avail[23] + cache[23] + buffer[23])},
	//}
	c.JSON(200, data)
}

func GetMemController(c *gin.Context) {
	total, svs := services.GetAllMemInfo(c)
	avail, svs1 := services.GetAvailMemInfo(c)
	cache, svs2 := services.GetMemCacheInfo(c)
	buffer, svs3 := services.GetMemBufferInfo(c)

	// Map để đồng bộ dữ liệu theo server
	memMap := make(map[string]*MemAll)

	// Đồng bộ Total RAM
	for i, server := range svs {
		if i < len(total) {
			memMap[server] = &MemAll{Server: server, Total: toInt64(total[i])}
		}
	}

	// Đồng bộ Available RAM
	for i, server := range svs1 {
		if i < len(avail) {
			if _, exists := memMap[server]; !exists {
				memMap[server] = &MemAll{Server: server}
			}
			memMap[server].Available += toInt64(avail[i])
		}
	}

	// Đồng bộ Cache RAM
	for i, server := range svs2 {
		if i < len(cache) {
			if _, exists := memMap[server]; !exists {
				memMap[server] = &MemAll{Server: server}
			}
			memMap[server].Available += toInt64(cache[i])
		}
	}

	// Đồng bộ Buffer RAM
	for i, server := range svs3 {
		if i < len(buffer) {
			if _, exists := memMap[server]; !exists {
				memMap[server] = &MemAll{Server: server}
			}
			memMap[server].Available += toInt64(buffer[i])
		}
	}

	// Chuyển map thành slice
	var data []MemAll
	for _, mem := range memMap {
		data = append(data, *mem)
	}

	//data := []MemAll{
	//	{Server: config.Server100, Total: total[0], Available: avail[0] + cache[0] + buffer[0]},
	//	{Server: config.Server101, Total: total[1], Available: avail[1] + cache[1] + buffer[1]},
	//	{Server: config.Server102, Total: total[2], Available: avail[2] + cache[2] + buffer[2]},
	//	{Server: config.Server104, Total: total[3], Available: avail[3] + cache[3] + buffer[3]},
	//	{Server: config.Server105, Total: total[4], Available: avail[4] + cache[4] + buffer[4]},
	//	{Server: config.Server106, Total: total[5], Available: avail[5] + cache[5] + buffer[5]},
	//	{Server: config.Server108, Total: total[6], Available: avail[6] + cache[6] + buffer[6]},
	//	{Server: config.Server109, Total: total[7], Available: avail[7] + cache[7] + buffer[7]},
	//	{Server: config.Server110, Total: total[8], Available: avail[8] + cache[8] + buffer[8]},
	//	{Server: config.Server112, Total: total[9], Available: avail[9] + cache[9] + buffer[9]},
	//	{Server: config.Server113, Total: total[10], Available: avail[10] + cache[10] + buffer[10]},
	//	{Server: config.Server114, Total: total[11], Available: avail[11] + cache[11] + buffer[11]},
	//	{Server: config.Server116, Total: total[12], Available: avail[12] + cache[12] + buffer[12]},
	//	{Server: config.Server117, Total: total[13], Available: avail[13] + cache[13] + buffer[13]},
	//	{Server: config.Server118, Total: total[14], Available: avail[14] + cache[14] + buffer[14]},
	//	{Server: config.Server119, Total: total[15], Available: avail[15] + cache[15] + buffer[15]},
	//	{Server: config.Server120, Total: total[16], Available: avail[16] + cache[16] + buffer[16]},
	//	{Server: config.Server121, Total: total[17], Available: avail[17] + cache[17] + buffer[17]},
	//	{Server: config.Server123, Total: total[18], Available: avail[18] + cache[18] + buffer[18]},
	//	{Server: config.Server124, Total: total[19], Available: avail[19] + cache[19] + buffer[19]},
	//	{Server: config.Server125, Total: total[20], Available: avail[20] + cache[20] + buffer[20]},
	//	{Server: config.Server127, Total: total[21], Available: avail[21] + cache[21] + buffer[21]},
	//	{Server: config.Server128, Total: total[22], Available: avail[22] + cache[22] + buffer[22]},
	//	{Server: config.Server129, Total: total[23], Available: avail[23] + cache[23] + buffer[23]},
	//}
	c.JSON(200, data)
}

func GetMemCacheController(c *gin.Context) {
	cache, _ := services.GetMemCacheInfo(c)
	var data []MemSnmpAll
	for i, server := range config.Servers {
		if i < len(cache) {
			data = append(data, MemSnmpAll{Server: server, Value: cache[i]})
		}
	}
	c.JSON(200, data)
}

func GetMemBufferController(c *gin.Context) {
	buffer, _ := services.GetMemBufferInfo(c)
	var data []MemSnmpAll
	for i, server := range config.Servers {
		if i < len(buffer) {
			data = append(data, MemSnmpAll{Server: server, Value: buffer[i]})
		}
	}
	c.JSON(200, data)
}

func GetCpuUsageVerController(c *gin.Context) {
	user, svs := services.GetCpuPercentUser(c)
	sys, svs1 := services.GetCpuPercentSystem(c)
	var data []CpuSnmpAll
	snmpMAp := make(map[string]*CpuSnmpAll)
	for i, server := range svs {
		if i < len(user) {
			snmpMAp[server] = &CpuSnmpAll{Server: server, Value: float64(user[i])}
		}
	}

	for i, server := range svs1 {
		if i < len(sys) {
			if _, exist := snmpMAp[server]; !exist {
				snmpMAp[server].Value += float64(sys[i])
			}
		}
	}

	for _, snmp := range snmpMAp {
		data = append(data, *snmp)
	}
	//data := []CpuSnmpAll{
	//	{Server: config.Server100, Value: float64(user[0] + sys[0])},
	//	{Server: config.Server101, Value: float64(user[1] + sys[1])},
	//	{Server: config.Server102, Value: float64(user[2] + sys[2])},
	//	{Server: config.Server104, Value: float64(user[3] + sys[3])},
	//	{Server: config.Server105, Value: float64(user[4] + sys[4])},
	//	{Server: config.Server106, Value: float64(user[5] + sys[5])},
	//	{Server: config.Server108, Value: float64(user[6] + sys[6])},
	//	{Server: config.Server109, Value: float64(user[7] + sys[7])},
	//	{Server: config.Server110, Value: float64(user[8] + sys[8])},
	//	{Server: config.Server112, Value: float64(user[9] + sys[9])},
	//	{Server: config.Server113, Value: float64(user[10] + sys[10])},
	//	{Server: config.Server114, Value: float64(user[11] + sys[11])},
	//	{Server: config.Server116, Value: float64(user[12] + sys[12])},
	//	{Server: config.Server117, Value: float64(user[13] + sys[13])},
	//	{Server: config.Server118, Value: float64(user[14] + sys[14])},
	//	{Server: config.Server119, Value: float64(user[15] + sys[15])},
	//	{Server: config.Server120, Value: float64(user[16] + sys[16])},
	//	{Server: config.Server121, Value: float64(user[17] + sys[17])},
	//	{Server: config.Server123, Value: float64(user[18] + sys[18])},
	//	{Server: config.Server124, Value: float64(user[19] + sys[19])},
	//	{Server: config.Server125, Value: float64(user[20] + sys[20])},
	//	{Server: config.Server127, Value: float64(user[21] + sys[21])},
	//	{Server: config.Server128, Value: float64(user[22] + sys[22])},
	//	{Server: config.Server129, Value: float64(user[23] + sys[23])},
	//}
	c.JSON(200, data)
}

func GetCpuUsageController(c *gin.Context) {
	user, svs := services.GetCpuRawUser(c)
	idle, _ := services.GetCpuRawIdle(c)
	var data []CpuSnmpAll
	for i, server := range svs {
		if i < len(user) {
			data = append(data, CpuSnmpAll{Server: server, Value: float64((user[i] * 100) / (user[i] + idle[i]))})
		}
	}
	//data = []CpuSnmpAll{
	//	{Server: config.Server100, Value: float64((user[0] * 100) / (user[0] + idle[0]))},
	//	{Server: config.Server101, Value: float64((user[1] * 100) / (user[1] + idle[1]))},
	//	{Server: config.Server102, Value: float64((user[2] * 100) / (user[2] + idle[2]))},
	//	{Server: config.Server104, Value: float64((user[3] * 100) / (user[3] + idle[3]))},
	//	{Server: config.Server105, Value: float64((user[4] * 100) / (user[4] + idle[4]))},
	//	{Server: config.Server106, Value: float64((user[5] * 100) / (user[5] + idle[5]))},
	//	{Server: config.Server108, Value: float64((user[6] * 100) / (user[6] + idle[6]))},
	//	{Server: config.Server109, Value: float64((user[7] * 100) / (user[7] + idle[7]))},
	//	{Server: config.Server110, Value: float64((user[8] * 100) / (user[8] + idle[8]))},
	//	{Server: config.Server112, Value: float64((user[9] * 100) / (user[9] + idle[9]))},
	//	{Server: config.Server113, Value: float64((user[10] * 100) / (user[10] + idle[10]))},
	//	{Server: config.Server114, Value: float64((user[11] * 100) / (user[11] + idle[11]))},
	//	{Server: config.Server116, Value: float64((user[12] * 100) / (user[12] + idle[12]))},
	//	{Server: config.Server117, Value: float64((user[13] * 100) / (user[13] + idle[13]))},
	//	{Server: config.Server118, Value: float64((user[14] * 100) / (user[14] + idle[14]))},
	//	{Server: config.Server119, Value: float64((user[15] * 100) / (user[15] + idle[15]))},
	//	{Server: config.Server120, Value: float64((user[16] * 100) / (user[16] + idle[16]))},
	//	{Server: config.Server121, Value: float64((user[17] * 100) / (user[17] + idle[17]))},
	//	{Server: config.Server123, Value: float64((user[18] * 100) / (user[18] + idle[18]))},
	//	{Server: config.Server124, Value: float64((user[19] * 100) / (user[19] + idle[19]))},
	//	{Server: config.Server125, Value: float64((user[20] * 100) / (user[20] + idle[20]))},
	//	{Server: config.Server127, Value: float64((user[21] * 100) / (user[21] + idle[21]))},
	//	{Server: config.Server128, Value: float64((user[22] * 100) / (user[22] + idle[22]))},
	//	{Server: config.Server129, Value: float64((user[23] * 100) / (user[23] + idle[23]))},
	//}
	//} else {
	//	data = []CpuSnmpAll{
	//		{Server: config.Server244, Value: float64(((user[0] - previousUser[0]) * 100) / ((user[0] - previousUser[0]) + (idle[0] - previousIdle[0])))},
	//		{Server: config.Server245, Value: float64(((user[1] - previousUser[1]) * 100) / ((user[1] - previousUser[1]) + (idle[1] - previousIdle[1])))},
	//		{Server: config.Server246, Value: float64(((user[2] - previousUser[2]) * 100) / ((user[2] - previousUser[2]) + (idle[2] - previousIdle[2])))},
	//		{Server: config.Server250, Value: float64(((user[3] - previousUser[3]) * 100) / ((user[3] - previousUser[3]) + (idle[3] - previousIdle[3])))},
	//		{Server: config.Server251, Value: float64(((user[4] - previousUser[4]) * 100) / ((user[4] - previousUser[4]) + (idle[4] - previousIdle[4])))},
	//		{Server: config.Server252, Value: float64(((user[5] - previousUser[5]) * 100) / ((user[5] - previousUser[5]) + (idle[5] - previousIdle[5])))},
	//	}
	//}

	//previousUser[0] = user[0]
	//previousUser[1] = user[1]
	//previousUser[2] = user[2]
	//previousUser[3] = user[3]
	//previousUser[4] = user[4]
	//previousUser[5] = user[5]
	//previousIdle[0] = idle[0]
	//previousIdle[1] = idle[1]
	//previousIdle[2] = idle[2]
	//previousIdle[3] = idle[3]
	//previousIdle[4] = idle[4]
	//previousIdle[5] = idle[5]
	c.JSON(200, data)
}

func GetDiskController(c *gin.Context) {
	total, svs := services.GetDiskTotal(c)
	avail, svs1 := services.GetDiskAvail(c)
	var data []MemAll

	memMap := make(map[string]*MemAll)

	for i, server := range svs {
		if i < len(total) {
			memMap[server] = &MemAll{Server: server, Total: toInt64(total[i])}
		}
	}

	for i, server := range svs1 {
		if i < len(avail) {
			if _, exists := memMap[server]; !exists {
				memMap[server] = &MemAll{Server: server}
			}
			memMap[server].Available += toInt64(avail[i])
		}
	}

	for _, mem := range memMap {
		data = append(data, *mem)
	}
	//data := []MemAll{
	//	{Server: config.Server100, Total: total[0], Available: avail[0]},
	//	{Server: config.Server101, Total: total[1], Available: avail[1]},
	//	{Server: config.Server102, Total: total[2], Available: avail[2]},
	//	{Server: config.Server104, Total: total[3], Available: avail[3]},
	//	{Server: config.Server105, Total: total[4], Available: avail[4]},
	//	{Server: config.Server106, Total: total[5], Available: avail[5]},
	//	{Server: config.Server108, Total: total[6], Available: avail[6]},
	//	{Server: config.Server109, Total: total[7], Available: avail[7]},
	//	{Server: config.Server110, Total: total[8], Available: avail[8]},
	//	{Server: config.Server112, Total: total[9], Available: avail[9]},
	//	{Server: config.Server113, Total: total[10], Available: avail[10]},
	//	{Server: config.Server114, Total: total[11], Available: avail[11]},
	//	{Server: config.Server116, Total: total[12], Available: avail[12]},
	//	{Server: config.Server117, Total: total[13], Available: avail[13]},
	//	{Server: config.Server118, Total: total[14], Available: avail[14]},
	//	{Server: config.Server119, Total: total[15], Available: avail[15]},
	//	{Server: config.Server120, Total: total[16], Available: avail[16]},
	//	{Server: config.Server121, Total: total[17], Available: avail[17]},
	//	{Server: config.Server123, Total: total[18], Available: avail[19]},
	//	{Server: config.Server127, Total: total[19], Available: avail[19]},
	//	{Server: config.Server128, Total: total[20], Available: avail[20]},
	//	{Server: config.Server129, Total: total[21], Available: avail[21]},
	//	//{Server: config.Server128, Total: total[22], Available: avail[22]},
	//	//{Server: config.Server129, Total: total[23], Available: avail[23]},
	//}
	c.JSON(200, data)
}

func GetIODataController(c *gin.Context) {
	receive, svs := services.GetIoReceiveData(c)
	sent, svs2 := services.GetIoSentData(c)
	var data []IoData

	ioMap := make(map[string]*IoData)
	for i, server := range svs {
		if i < len(receive) {
			ioMap[server] = &IoData{Server: server, Receive: receive[i]}
		}
	}

	for i, server := range svs2 {
		if _, exist := ioMap[server]; !exist {
			ioMap[server].Sent = sent[i]
		}
	}

	for _, io := range ioMap {
		data = append(data, *io)
	}
	//data := []IoData{
	//	{Server: config.Server100, Receive: receive[0], Sent: sent[0]},
	//	{Server: config.Server101, Receive: receive[1], Sent: sent[1]},
	//	{Server: config.Server102, Receive: receive[2], Sent: sent[2]},
	//	{Server: config.Server104, Receive: receive[3], Sent: sent[3]},
	//	{Server: config.Server105, Receive: receive[4], Sent: sent[4]},
	//	{Server: config.Server106, Receive: receive[5], Sent: sent[5]},
	//	{Server: config.Server108, Receive: receive[6], Sent: sent[6]},
	//	{Server: config.Server109, Receive: receive[7], Sent: sent[7]},
	//	{Server: config.Server110, Receive: receive[8], Sent: sent[8]},
	//	{Server: config.Server112, Receive: receive[9], Sent: sent[9]},
	//	{Server: config.Server113, Receive: receive[10], Sent: sent[10]},
	//	{Server: config.Server114, Receive: receive[11], Sent: sent[11]},
	//	{Server: config.Server116, Receive: receive[12], Sent: sent[12]},
	//	{Server: config.Server117, Receive: receive[13], Sent: sent[13]},
	//	{Server: config.Server118, Receive: receive[14], Sent: sent[14]},
	//	{Server: config.Server119, Receive: receive[15], Sent: sent[15]},
	//	{Server: config.Server120, Receive: receive[16], Sent: sent[16]},
	//	{Server: config.Server121, Receive: receive[17], Sent: sent[17]},
	//	{Server: config.Server123, Receive: receive[18], Sent: sent[18]},
	//	{Server: config.Server124, Receive: receive[19], Sent: sent[19]},
	//	{Server: config.Server125, Receive: receive[20], Sent: sent[20]},
	//	{Server: config.Server127, Receive: receive[21], Sent: sent[21]},
	//	{Server: config.Server128, Receive: receive[22], Sent: sent[22]},
	//	{Server: config.Server129, Receive: receive[23], Sent: sent[23]},
	//}
	c.JSON(200, data)
}
func GetMemSwapController(c *gin.Context) {
	avail, svs := services.GetMemSwapAvail(c)
	total, svs1 := services.GetMemSwapTotal(c)
	var data []MemAll

	memMap := make(map[string]*MemAll)

	for i, server := range svs {
		if i < len(total) {
			memMap[server] = &MemAll{Server: server, Total: toInt64(total[i])}
		}
	}

	for i, server := range svs1 {
		if i < len(avail) {
			if _, exists := memMap[server]; !exists {
				memMap[server] = &MemAll{Server: server}
			}
			memMap[server].Available += toInt64(avail[i])
		}
	}

	for _, mem := range memMap {
		data = append(data, *mem)
	}
	//data := []MemAll{
	//	{Server: config.Server100, Total: total[0], Available: avail[0]},
	//	{Server: config.Server101, Total: total[1], Available: avail[1]},
	//	{Server: config.Server102, Total: total[2], Available: avail[2]},
	//	{Server: config.Server104, Total: total[3], Available: avail[3]},
	//	{Server: config.Server105, Total: total[4], Available: avail[4]},
	//	{Server: config.Server106, Total: total[5], Available: avail[5]},
	//	{Server: config.Server108, Total: total[6], Available: avail[6]},
	//	{Server: config.Server109, Total: total[7], Available: avail[7]},
	//	{Server: config.Server110, Total: total[8], Available: avail[8]},
	//	{Server: config.Server112, Total: total[9], Available: avail[9]},
	//	{Server: config.Server113, Total: total[10], Available: avail[10]},
	//	{Server: config.Server114, Total: total[11], Available: avail[11]},
	//	{Server: config.Server116, Total: total[12], Available: avail[12]},
	//	{Server: config.Server117, Total: total[13], Available: avail[13]},
	//	{Server: config.Server118, Total: total[14], Available: avail[14]},
	//	{Server: config.Server119, Total: total[15], Available: avail[15]},
	//	{Server: config.Server120, Total: total[16], Available: avail[16]},
	//	{Server: config.Server121, Total: total[17], Available: avail[17]},
	//	{Server: config.Server123, Total: total[18], Available: avail[18]},
	//	{Server: config.Server124, Total: total[19], Available: avail[19]},
	//	{Server: config.Server125, Total: total[20], Available: avail[20]},
	//	{Server: config.Server127, Total: total[21], Available: avail[21]},
	//	{Server: config.Server128, Total: total[22], Available: avail[22]},
	//	{Server: config.Server129, Total: total[23], Available: avail[23]},
	//}
	c.JSON(200, data)
}

func GetOidWithServerController(c *gin.Context) {
	var request validations.GetEachOidWithServer
	if requestErr := c.ShouldBind(&request); requestErr != nil {
		errRes := strings.Split(requestErr.Error(), ": ")
		services.ValidationError(c, "These fields are required!", errRes)
		return
	}
	oid := c.Query("oid")
	server := c.Query("server")
	values := services.GetDataSNMP(c, server, oid)
	c.JSON(200, values)
}

func toInt64(val interface{}) int64 {
	switch v := val.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case nil:
		return 0
	default:
		log.Printf("[ERROR] Invalid SNMP data type: %T\n", v)
		return 0
	}
}

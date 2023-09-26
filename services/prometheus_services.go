package services

import (
	"Monitor_Platform/config"
	"Monitor_Platform/snmp"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"math/big"
)

var (
	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	})
)

var (
	cpuAvail = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_memory_memory",
		Help: "Current memory of the CPU.",
	})
)
var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total_local",
			Help: "Number of HTTP requests.",
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(cpuAvail)
}

func SetCPUTemperature(temperature float64) {
	cpuTemp.Set(temperature)
}

func IncHttpRequest(c *gin.Context) {
	httpRequests.WithLabelValues(c.Request.Method, c.FullPath()).Inc()
}

func MemToTalAvail() {
	CpuUsage := snmp.GetOIDInt(config.MemTotalReal)
	num := new(big.Float).SetInt(CpuUsage)
	result, _ := num.Float64()
	convert2Gb := result * 0.000001
	cpuAvail.Set(convert2Gb)
}

package controllers

type Response struct {
	Value interface{} `json:"Value"`
}

type MemSnmpAll struct {
	Server string      `json:"server"`
	Value  interface{} `json:"Value"`
}

type MemAll struct {
	Server    string      `json:"server"`
	Available interface{} `json:"available"`
	Total     interface{} `json:"total"`
}

type CpuSnmpAll struct {
	Server string  `json:"server"`
	Value  float64 `json:"Value"`
}

type Kafka struct {
	Server  string  `json:"server"`
	Bytein  float64 `json:"bytein"`
	Byteout float64 `json:"byteout"`
}

type KpiVtrack struct {
	Pod  string  `json:"pod"`
	Rate float64 `json:"rate"`
}

package controllers

type Response struct {
	Value interface{} `json:"Value"`
}

type MemSnmpAll struct {
	Server string      `json:"server"`
	Value  interface{} `json:"Value"`
}

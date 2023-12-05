package snmp

import (
	"Monitor_Platform/config"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"log"
	"math/big"
	"time"
)

// Tạo một đối tượng SNMP manager
var Manager = &gosnmp.GoSNMP{
	Target:    config.Target,                  // Địa chỉ IP của SNMP agent
	Port:      config.SNMP_Port,               // Port mặc định của SNMP
	Version:   gosnmp.Version2c,               // Sử dụng SNMP version 2c
	Community: "public",                       // Community string
	Timeout:   time.Duration(5) * time.Second, // Thời gian chờ tối đa
}

func GetOIDInt(oid string) *big.Int {
	err := Manager.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer Manager.Conn.Close()
	oids := []string{oid}
	result, err2 := Manager.Get(oids) // Lấy dữ liệu từ SNMP agent
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}

	// In dữ liệu lấy được từ agent
	for i, variable := range result.Variables {
		fmt.Printf("%d: oid: %s ", i, variable.Name)
		switch variable.Type {
		case gosnmp.OctetString:
			return nil
		default:
			return gosnmp.ToBigInt(variable.Value)
		}
	}
	return nil
}

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
	Timeout:   time.Duration(2) * time.Second, // Thời gian chờ tối đa
}

func GetValueOID() {
	// OID của các biến mà bạn muốn lấy
	oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0", config.MemTotalReal}
	result, err2 := Manager.Get(oids) // Lấy dữ liệu từ SNMP agent
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}

	// In dữ liệu lấy được từ agent
	for i, variable := range result.Variables {
		fmt.Printf("%d: oid: %s ", i, variable.Name)

		switch variable.Type {
		case gosnmp.OctetString:
			fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
		case gosnmp.Gauge32:
			fmt.Printf("Gauge: %d\n", gosnmp.ToBigInt(variable.Value))
		default:
			fmt.Printf("number: %d\n", gosnmp.ToBigInt(variable.Value))
		}
	}
}

func GetOIDInt(oid string) *big.Int {
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

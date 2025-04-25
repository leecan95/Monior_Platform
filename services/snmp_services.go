package services

import (
	"Monitor_Platform/config"
	"Monitor_Platform/snmp"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gosnmp/gosnmp"
	"log"
	"sync"
	"time"
)

var mu sync.Mutex

const maxConcurrentSNMP = 20

var snmpClients = make(map[string]*gosnmp.GoSNMP)

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

func getSNMPClient(target string) (*gosnmp.GoSNMP, error) {
    if target == "" {
        return nil, fmt.Errorf("empty target address")
    }

    mu.Lock()
    defer mu.Unlock()

    // Nếu đã có kết nối, dùng lại
    if client, exists := snmpClients[target]; exists {
        // Check if the connection is still valid
        if err := client.Connect(); err == nil {
            // Connection is still valid, use it
            return client, nil
        }
        // Connection is invalid, delete it and create a new one
        delete(snmpClients, target)
    }

    // Tạo kết nối SNMP mới
    client := &gosnmp.GoSNMP{
        Target:             target,
        Port:               161,
        Version:            gosnmp.Version2c,
        Community:          "public",
        Timeout:            time.Second * 2,
        Retries:            2,
        ExponentialTimeout: true, // Add exponential backoff for timeouts
    }

    err := client.Connect()
    if err != nil {
        log.Printf("SNMP Connect() err for %s: %v", target, err)
        return nil, err
    }

    snmpClients[target] = client
    return client, nil
}
func GetAllMemInfo(c *gin.Context) ([]int64, []string) {
    var sum []int64
    var svs []string
    var wg sync.WaitGroup
    var svsMutex sync.Mutex // Add mutex for thread-safe server list updates
    results := make(chan int64, len(config.CMP_Servers))    // Channel để thu thập kết quả
    semaphore := make(chan struct{}, maxConcurrentSNMP) // Giới hạn 20 request SNMP chạy đồng thời

    for _, server := range config.CMP_Servers {
        // Skip empty server addresses
        if server == "" {
            log.Printf("[WARNING] Empty server address found in config.Servers")
            continue
        }
        
        wg.Add(1)
        semaphore <- struct{}{} // Giữ chỗ trong semaphore

        go func(server string) {
            defer wg.Done()
            defer func() { <-semaphore }() // Giải phóng chỗ trong semaphore

            client, err := getSNMPClient(server)
            if err != nil {
                log.Printf("Failed to get SNMP client for %s: %v", server, err)
                return
            }

            var resultVal interface{}
            var fetchErr error
            
            for i := 0; i < 3; i++ { // Thử lại tối đa 3 lần
                result, err := client.Get([]string{config.MemTotalReal})
                if err == nil && len(result.Variables) > 0 && result.Variables[0].Value != nil {
                    // Use mutex to safely update the server list
                    svsMutex.Lock()
                    svs = append(svs, server)
                    svsMutex.Unlock()
                    
                    resultVal = result.Variables[0].Value
                    fetchErr = nil
                    break
                }
                
                fetchErr = err
                log.Printf("[WARNING] Retry %d: Failed SNMP from %s: %v", i+1, server, err)
                
                // Instead of client.Close(), reconnect the client
                time.Sleep(500 * time.Millisecond)
                // Reconnect by getting a fresh client
                client, err = getSNMPClient(server)
                if err != nil {
                    log.Printf("Failed to reconnect SNMP client for %s: %v", server, err)
                }
            }

            if fetchErr != nil {
                log.Printf("[ERROR] SNMP Get failed from %s after 3 retries: %v", server, fetchErr)
                return
            }

			if resultVal == nil {
				log.Printf("[WARNING] SNMP response is nil from %s", server)
				return
			}

			// Kiểm tra kiểu dữ liệu trả về
			var result int64
			switch v := resultVal.(type) {
			case int64:
				result = v
			case int:
				result = int64(v)
			case uint:
				result = int64(v)
			case uint64:
				result = int64(v)
			case float64:
				result = int64(v)
			case string:
				log.Printf("[INFO] SNMP data from %s is string: %s", server, v)
				return
			case []byte:
				log.Printf("[INFO] SNMP data from %s is []byte: %s", server, string(v))
				return
			default:
				log.Printf("[ERROR] Invalid SNMP data type from %s: %T", server, v)
				return
			}

			results <- result // Đưa dữ liệu vào channel
		}(server)
	}

	// Đợi tất cả goroutines hoàn thành
	wg.Wait()
	close(results)

	// Lấy dữ liệu từ channel
	for res := range results {
		sum = append(sum, res)
	}

	return sum, svs
}

func GetAvailMemInfo(c *gin.Context) ([]int64, []string) {
	var sum []int64
	var svs []string
	var wg sync.WaitGroup
	results := make(chan int64, len(config.CMP_Servers))    // Channel để thu thập kết quả
	semaphore := make(chan struct{}, maxConcurrentSNMP) // Giới hạn 20 request SNMP chạy đồng thời

	for _, server := range config.CMP_Servers {
		wg.Add(1)
		semaphore <- struct{}{} // Giữ chỗ trong semaphore

		go func(server string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Giải phóng chỗ trong semaphore

			client, err := getSNMPClient(server)
			if err != nil {
				log.Printf("Failed to get SNMP client for %s: %v", server, err)
				return
			}

			var resultVal interface{}
			for i := 0; i < 3; i++ { // Thử lại tối đa 3 lần
				result, err := client.Get([]string{config.MemAvailReal})
				if err == nil && len(result.Variables) > 0 {
					svs = append(svs, server)
					resultVal = result.Variables[0].Value
					break
				}
				log.Printf("[WARNING] Retry %d: Failed SNMP from %s\n", i+1, server)
				time.Sleep(500 * time.Millisecond)
			}

			if err != nil {
				log.Printf("[ERROR] SNMP Get failed from %s: %v\n", server, err)
				return
			}

			if resultVal == nil {
				log.Printf("[WARNING] SNMP response is nil from %s\n", server)
				return
			}

			// Kiểm tra kiểu dữ liệu trả về
			var result int64
			switch v := resultVal.(type) {
			case int64:
				result = v
			case int:
				result = int64(v)
			case float64:
				result = int64(v)
			case string:
				log.Printf("[INFO] SNMP data from %s is string: %s\n", server, v)
				return
			case []byte:
				log.Printf("[INFO] SNMP data from %s is []byte: %s\n", server, string(v))
				return
			default:
				log.Printf("[ERROR] Invalid SNMP data type from %s: %T\n", server, v)
				return
			}

			results <- result // Đưa dữ liệu vào channel
		}(server)
	}

	// Đợi tất cả goroutines hoàn thành
	wg.Wait()
	close(results)

	// Lấy dữ liệu từ channel
	for res := range results {
		sum = append(sum, res)
	}

	return sum, svs
}

func GetMemCacheInfo(c *gin.Context) ([]int64, []string) {
	var sum []int64
	var svs []string
	var wg sync.WaitGroup
	results := make(chan int64, len(config.CMP_Servers))    // Channel để thu thập kết quả
	semaphore := make(chan struct{}, maxConcurrentSNMP) // Giới hạn 20 request SNMP chạy đồng thời

	for _, server := range config.CMP_Servers {
		wg.Add(1)
		semaphore <- struct{}{} // Giữ chỗ trong semaphore

		go func(server string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Giải phóng chỗ trong semaphore

			client, err := getSNMPClient(server)
			if err != nil {
				log.Printf("Failed to get SNMP client for %s: %v", server, err)
				return
			}

			var resultVal interface{}
			for i := 0; i < 3; i++ { // Thử lại tối đa 3 lần
				result, err := client.Get([]string{config.MemCacheReal})
				if err == nil && len(result.Variables) > 0 {
					svs = append(svs, server)
					resultVal = result.Variables[0].Value
					break
				}
				log.Printf("[WARNING] Retry %d: Failed SNMP from %s\n", i+1, server)
				time.Sleep(500 * time.Millisecond)
			}

			if err != nil {
				log.Printf("[ERROR] SNMP Get failed from %s: %v\n", server, err)
				return
			}

			if resultVal == nil {
				log.Printf("[WARNING] SNMP response is nil from %s\n", server)
				return
			}

			// Kiểm tra kiểu dữ liệu trả về
			var result int64
			switch v := resultVal.(type) {
			case int64:
				result = v
			case int:
				result = int64(v)
			case float64:
				result = int64(v)
			case string:
				log.Printf("[INFO] SNMP data from %s is string: %s\n", server, v)
				return
			case []byte:
				log.Printf("[INFO] SNMP data from %s is []byte: %s\n", server, string(v))
				return
			default:
				log.Printf("[ERROR] Invalid SNMP data type from %s: %T\n", server, v)
				return
			}

			results <- result // Đưa dữ liệu vào channel
		}(server)
	}

	// Đợi tất cả goroutines hoàn thành
	wg.Wait()
	close(results)

	// Lấy dữ liệu từ channel
	for res := range results {
		sum = append(sum, res)
	}

	return sum, svs
}

func GetMemBufferInfo(c *gin.Context) ([]int64, []string) {
	var sum []int64
	var svs []string
	var wg sync.WaitGroup
	results := make(chan int64, len(config.CMP_Servers))    // Channel để thu thập kết quả
	semaphore := make(chan struct{}, maxConcurrentSNMP) // Giới hạn 20 request SNMP chạy đồng thời

	for _, server := range config.CMP_Servers {
		wg.Add(1)
		semaphore <- struct{}{} // Giữ chỗ trong semaphore

		go func(server string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Giải phóng chỗ trong semaphore

			client, err := getSNMPClient(server)
			if err != nil {
				log.Printf("Failed to get SNMP client for %s: %v", server, err)
				return
			}

			var resultVal interface{}
			for i := 0; i < 3; i++ { // Thử lại tối đa 3 lần
				result, err := client.Get([]string{config.MemBufferReal})
				if err == nil && len(result.Variables) > 0 {
					svs = append(svs, server)
					resultVal = result.Variables[0].Value
					break
				}
				log.Printf("[WARNING] Retry %d: Failed SNMP from %s\n", i+1, server)
				time.Sleep(500 * time.Millisecond)
			}

			if err != nil {
				log.Printf("[ERROR] SNMP Get failed from %s: %v\n", server, err)
				return
			}

			if resultVal == nil {
				log.Printf("[WARNING] SNMP response is nil from %s\n", server)
				return
			}

			// Kiểm tra kiểu dữ liệu trả về
			var result int64
			switch v := resultVal.(type) {
			case int64:
				result = v
			case int:
				result = int64(v)
			case float64:
				result = int64(v)
			case string:
				log.Printf("[INFO] SNMP data from %s is string: %s\n", server, v)
				return
			case []byte:
				log.Printf("[INFO] SNMP data from %s is []byte: %s\n", server, string(v))
				return
			default:
				log.Printf("[ERROR] Invalid SNMP data type from %s: %T\n", server, v)
				return
			}

			results <- result // Đưa dữ liệu vào channel
		}(server)
	}

	// Đợi tất cả goroutines hoàn thành
	wg.Wait()
	close(results)

	// Lấy dữ liệu từ channel
	for res := range results {
		sum = append(sum, res)
	}

	return sum, svs
}

func GetCpuRawUser(c *gin.Context) ([]int64, []string) {
	var sum []int64
	var svs []string
	var wg sync.WaitGroup
	results := make(chan int64, len(config.CMP_Servers))    // Channel để thu thập kết quả
	semaphore := make(chan struct{}, maxConcurrentSNMP) // Giới hạn 20 request SNMP chạy đồng thời

	for _, server := range config.CMP_Servers {
		wg.Add(1)
		semaphore <- struct{}{} // Giữ chỗ trong semaphore

		go func(server string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Giải phóng chỗ trong semaphore

			client, err := getSNMPClient(server)
			if err != nil {
				log.Printf("Failed to get SNMP client for %s: %v", server, err)
				return
			}

			var resultVal interface{}
			for i := 0; i < 3; i++ { // Thử lại tối đa 3 lần
				result, err := client.Get([]string{config.SsCpuRawUser})
				if err == nil && len(result.Variables) > 0 {
					svs = append(svs, server)
					resultVal = result.Variables[0].Value
					break
				}
				log.Printf("[WARNING] Retry %d: Failed SNMP from %s\n", i+1, server)
				time.Sleep(500 * time.Millisecond)
			}

			if err != nil {
				log.Printf("[ERROR] SNMP Get failed from %s: %v\n", server, err)
				return
			}

			if resultVal == nil {
				log.Printf("[WARNING] SNMP response is nil from %s\n", server)
				return
			}

			// Kiểm tra kiểu dữ liệu trả về
			var result int64
			switch v := resultVal.(type) {
			case int64:
				result = v
			case int:
				result = int64(v)
			case float64:
				result = int64(v)
			case string:
				log.Printf("[INFO] SNMP data from %s is string: %s\n", server, v)
				return
			case []byte:
				log.Printf("[INFO] SNMP data from %s is []byte: %s\n", server, string(v))
				return
			default:
				log.Printf("[ERROR] Invalid SNMP data type from %s: %T\n", server, v)
				return
			}

			results <- result // Đưa dữ liệu vào channel
		}(server)
	}

	// Đợi tất cả goroutines hoàn thành
	wg.Wait()
	close(results)

	// Lấy dữ liệu từ channel
	for res := range results {
		sum = append(sum, res)
	}

	return sum, svs

}

func GetCpuRawNice(c *gin.Context) []int64 {
	var sum []int64
	var result int64
	var ok bool
	var val interface{}

	for _, server := range config.CMP_Servers {
		val = GetDataSNMP(c, server, config.SsCpuRawNice)
		if val == nil {
			log.Printf("[ERROR] Failed to fetch SNMP from %s\n", server)
		}
		result, ok = val.(int64)
		if !ok {
			log.Printf("[ERROR] Invalid SNMP data type from %s\n", server)
			continue
		}
		result = int64(val.(int))
		sum = append(sum, result)
	}
	//
	//val := GetDataSNMP(c, config.Server244, config.SsCpuRawNice)
	//fmt.Printf("102 gia tri %v\n", val)
	//result := int64(val.(uint))
	//sum = append(sum, result)
	//val = GetDataSNMP(c, config.Server245, config.SsCpuRawNice)
	//result = int64(val.(uint))
	//sum = append(sum, result)
	//val = GetDataSNMP(c, config.Server246, config.SsCpuRawNice)
	//result = int64(val.(uint))
	//sum = append(sum, result)
	//val = GetDataSNMP(c, config.Server250, config.SsCpuRawNice)
	//result = int64(val.(uint))
	//sum = append(sum, result)
	//val = GetDataSNMP(c, config.Server251, config.SsCpuRawNice)
	//result = int64(val.(uint))
	//sum = append(sum, result)
	//val = GetDataSNMP(c, config.Server252, config.SsCpuRawNice)
	//result = int64(val.(uint))
	//sum = append(sum, result)
	return sum
}

func GetCpuRawSystem(c *gin.Context) []int64 {
	var sum []int64
	var result int64
	var ok bool
	var val interface{}

	for _, server := range config.CMP_Servers {
		val = GetDataSNMP(c, server, config.SsCpuRawSystem)
		if val == nil {
			log.Printf("[ERROR] Failed to fetch SNMP from %s\n", server)
		}
		result, ok = val.(int64)
		if !ok {
			log.Printf("[ERROR] Invalid SNMP data type from %s\n", server)
			continue
		}
		result = int64(val.(int))
		sum = append(sum, result)
	}
	//val := GetDataSNMP(c, config.Server244, config.SsCpuRawSystem)
	//fmt.Printf("126 gia tri %v\n", val)
	//result := int64(val.(uint))
	//sum = append(sum, result)
	//val = GetDataSNMP(c, config.Server245, config.SsCpuRawSystem)
	//result = int64(val.(uint))
	//sum = append(sum, result)
	//val = GetDataSNMP(c, config.Server246, config.SsCpuRawSystem)
	//result = int64(val.(uint))
	//sum = append(sum, result)
	//val = GetDataSNMP(c, config.Server250, config.SsCpuRawSystem)
	//result = int64(val.(uint))
	//sum = append(sum, result)
	//val = GetDataSNMP(c, config.Server251, config.SsCpuRawSystem)
	//result = int64(val.(uint))
	//sum = append(sum, result)
	//val = GetDataSNMP(c, config.Server252, config.SsCpuRawSystem)
	//result = int64(val.(uint))
	//sum = append(sum, result)
	return sum
}

func GetCpuRawIdle(c *gin.Context) ([]int64, []string) {
	var sum []int64
	var svs []string
	var wg sync.WaitGroup
	results := make(chan int64, len(config.CMP_Servers))    // Channel để thu thập kết quả
	semaphore := make(chan struct{}, maxConcurrentSNMP) // Giới hạn 20 request SNMP chạy đồng thời

	for _, server := range config.CMP_Servers {
		wg.Add(1)
		semaphore <- struct{}{} // Giữ chỗ trong semaphore

		go func(server string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Giải phóng chỗ trong semaphore

			client, err := getSNMPClient(server)
			if err != nil {
				log.Printf("Failed to get SNMP client for %s: %v", server, err)
				return
			}

			var resultVal interface{}
			for i := 0; i < 3; i++ { // Thử lại tối đa 3 lần
				result, err := client.Get([]string{config.SsCpuRawIdle})
				if err == nil && len(result.Variables) > 0 {
					svs = append(svs, server)
					resultVal = result.Variables[0].Value
					break
				}
				log.Printf("[WARNING] Retry %d: Failed SNMP from %s\n", i+1, server)
				time.Sleep(500 * time.Millisecond)
			}

			if err != nil {
				log.Printf("[ERROR] SNMP Get failed from %s: %v\n", server, err)
				return
			}

			if resultVal == nil {
				log.Printf("[WARNING] SNMP response is nil from %s\n", server)
				return
			}

			// Kiểm tra kiểu dữ liệu trả về
			var result int64
			switch v := resultVal.(type) {
			case int64:
				result = v
			case int:
				result = int64(v)
			case float64:
				result = int64(v)
			case string:
				log.Printf("[INFO] SNMP data from %s is string: %s\n", server, v)
				return
			case []byte:
				log.Printf("[INFO] SNMP data from %s is []byte: %s\n", server, string(v))
				return
			default:
				log.Printf("[ERROR] Invalid SNMP data type from %s: %T\n", server, v)
				return
			}

			results <- result // Đưa dữ liệu vào channel
		}(server)
	}

	// Đợi tất cả goroutines hoàn thành
	wg.Wait()
	close(results)

	// Lấy dữ liệu từ channel
	for res := range results {
		sum = append(sum, res)
	}

	return sum, svs
}

func GetCpuRawWait(c *gin.Context) []int64 {
	var sum []int64
	var result int64
	var ok bool
	var val interface{}

	for _, server := range config.CMP_Servers {
		val = GetDataSNMP(c, server, config.SsCpuRawWait)
		if val == nil {
			log.Printf("[ERROR] Failed to fetch SNMP from %s\n", server)
		}
		result, ok = val.(int64)
		if !ok {
			log.Printf("[ERROR] Invalid SNMP data type from %s\n", server)
			continue
		}
		result = int64(val.(int))
		sum = append(sum, result)
	}
	//val := GetDataSNMP(c, config.Server244, config.SsCpuRawWait)
	//fmt.Printf("174 gia tri %v\n", val)
	//result := int64(val.(uint))
	//sum = append(sum, result)
	//val = GetDataSNMP(c, config.Server245, config.SsCpuRawWait)
	//result = int64(val.(uint))
	//sum = append(sum, result)
	//val = GetDataSNMP(c, config.Server246, config.SsCpuRawWait)
	//result = int64(val.(uint))
	//sum = append(sum, result)
	//val = GetDataSNMP(c, config.Server250, config.SsCpuRawWait)
	//result = int64(val.(uint))
	//sum = append(sum, result)
	//val = GetDataSNMP(c, config.Server251, config.SsCpuRawWait)
	//result = int64(val.(uint))
	//sum = append(sum, result)
	//val = GetDataSNMP(c, config.Server252, config.SsCpuRawWait)
	//result = int64(val.(uint))
	//sum = append(sum, result)
	return sum
}

func GetCpuRawKernel(c *gin.Context) ([]int64, []string) {
	var sum []int64
	var svs []string
	var wg sync.WaitGroup
	results := make(chan int64, len(config.CMP_Servers))    // Channel để thu thập kết quả
	semaphore := make(chan struct{}, maxConcurrentSNMP) // Giới hạn 20 request SNMP chạy đồng thời

	for _, server := range config.CMP_Servers {
		wg.Add(1)
		semaphore <- struct{}{} // Giữ chỗ trong semaphore

		go func(server string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Giải phóng chỗ trong semaphore

			client, err := getSNMPClient(server)
			if err != nil {
				log.Printf("Failed to get SNMP client for %s: %v", server, err)
				return
			}

			var resultVal interface{}
			for i := 0; i < 3; i++ { // Thử lại tối đa 3 lần
				result, err := client.Get([]string{config.SsCpuRawKernel})
				if err == nil && len(result.Variables) > 0 {
					svs = append(svs, server)
					resultVal = result.Variables[0].Value
					break
				}
				log.Printf("[WARNING] Retry %d: Failed SNMP from %s\n", i+1, server)
				time.Sleep(500 * time.Millisecond)
			}

			if err != nil {
				log.Printf("[ERROR] SNMP Get failed from %s: %v\n", server, err)
				return
			}

			if resultVal == nil {
				log.Printf("[WARNING] SNMP response is nil from %s\n", server)
				return
			}

			// Kiểm tra kiểu dữ liệu trả về
			var result int64
			switch v := resultVal.(type) {
			case int64:
				result = v
			case int:
				result = int64(v)
			case float64:
				result = int64(v)
			case string:
				log.Printf("[INFO] SNMP data from %s is string: %s\n", server, v)
				return
			case []byte:
				log.Printf("[INFO] SNMP data from %s is []byte: %s\n", server, string(v))
				return
			default:
				log.Printf("[ERROR] Invalid SNMP data type from %s: %T\n", server, v)
				return
			}

			results <- result // Đưa dữ liệu vào channel
		}(server)
	}

	// Đợi tất cả goroutines hoàn thành
	wg.Wait()
	close(results)

	// Lấy dữ liệu từ channel
	for res := range results {
		sum = append(sum, res)
	}

	return sum, svs
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

func GetDiskTotal(c *gin.Context) ([]int64, []string) {
    var sum []int64
    var svs []string
    var wg sync.WaitGroup
    var svsMutex sync.Mutex // Add mutex for thread-safe server list updates
    results := make(chan int64, len(config.CMP_Servers))    // Channel để thu thập kết quả
    semaphore := make(chan struct{}, maxConcurrentSNMP) // Giới hạn 20 request SNMP chạy đồng thời

    for _, server := range config.CMP_Servers {
        // Skip empty server addresses
        if server == "" {
            log.Printf("[WARNING] Empty server address found in config.Servers")
            continue
        }
        
        wg.Add(1)
        semaphore <- struct{}{} // Giữ chỗ trong semaphore

        go func(server string) {
            defer wg.Done()
            defer func() { <-semaphore }() // Giải phóng chỗ trong semaphore

            client, err := getSNMPClient(server)
            if err != nil {
                log.Printf("Failed to get SNMP client for %s: %v", server, err)
                return
            }

            var resultVal interface{}
            var fetchErr error
            
            for i := 0; i < 3; i++ { // Thử lại tối đa 3 lần
                result, err := client.Get([]string{config.DskTotalNew})
                if err == nil && len(result.Variables) > 0 && result.Variables[0].Value != nil {
                    // Use mutex to safely update the server list
                    svsMutex.Lock()
                    svs = append(svs, server)
                    svsMutex.Unlock()
                    
                    resultVal = result.Variables[0].Value
                    fetchErr = nil
                    break
                }
                
                fetchErr = err
                log.Printf("[WARNING] Retry %d: Failed SNMP from %s: %v", i+1, server, err)
                
                // Reconnect by getting a fresh client
                time.Sleep(500 * time.Millisecond)
                client, err = getSNMPClient(server)
                if err != nil {
                    log.Printf("Failed to reconnect SNMP client for %s: %v", server, err)
                }
            }

            if fetchErr != nil {
                log.Printf("[ERROR] SNMP Get failed from %s after 3 retries: %v", server, fetchErr)
                return
            }

            if resultVal == nil {
                log.Printf("[WARNING] SNMP response is nil from %s", server)
                return
            }

            // Kiểm tra kiểu dữ liệu trả về
            var result int64
            switch v := resultVal.(type) {
            case int64:
                result = v
            case int:
                result = int64(v)
            case uint:
                result = int64(v)
            case uint64:
                result = int64(v)
            case float64:
                result = int64(v)
            case string:
                log.Printf("[INFO] SNMP data from %s is string: %s", server, v)
                return
            case []byte:
                log.Printf("[INFO] SNMP data from %s is []byte: %s", server, string(v))
                return
            default:
                log.Printf("[ERROR] Invalid SNMP data type from %s: %T", server, v)
                return
            }

            results <- result // Đưa dữ liệu vào channel
        }(server)
    }

    // Đợi tất cả goroutines hoàn thành
    wg.Wait()
    close(results)

    // Lấy dữ liệu từ channel
    for res := range results {
        sum = append(sum, res)
    }

    return sum, svs
}

func GetDiskAvail(c *gin.Context) ([]int64, []string) {
    var sum []int64
    var svs []string
    var wg sync.WaitGroup
    var svsMutex sync.Mutex // Add mutex for thread-safe server list updates
    results := make(chan int64, len(config.CMP_Servers))    // Channel để thu thập kết quả
    semaphore := make(chan struct{}, maxConcurrentSNMP) // Giới hạn 20 request SNMP chạy đồng thời

    for _, server := range config.CMP_Servers {
        // Skip empty server addresses
        if server == "" {
            log.Printf("[WARNING] Empty server address found in config.Servers")
            continue
        }
        
        wg.Add(1)
        semaphore <- struct{}{} // Giữ chỗ trong semaphore

        go func(server string) {
            defer wg.Done()
            defer func() { <-semaphore }() // Giải phóng chỗ trong semaphore

            client, err := getSNMPClient(server)
            if err != nil {
                log.Printf("Failed to get SNMP client for %s: %v", server, err)
                return
            }

            var resultVal interface{}
            var fetchErr error
            
            for i := 0; i < 3; i++ { // Thử lại tối đa 3 lần
                // Fixed: Using the correct OID for available disk space
                result, err := client.Get([]string{config.DskAvailNew})
                if err == nil && len(result.Variables) > 0 && result.Variables[0].Value != nil {
                    // Use mutex to safely update the server list
                    svsMutex.Lock()
                    svs = append(svs, server)
                    svsMutex.Unlock()
                    
                    resultVal = result.Variables[0].Value
                    fetchErr = nil
                    break
                }
                
                fetchErr = err
                log.Printf("[WARNING] Retry %d: Failed SNMP from %s: %v", i+1, server, err)
                
                // Reconnect by getting a fresh client
                time.Sleep(500 * time.Millisecond)
                client, err = getSNMPClient(server)
                if err != nil {
                    log.Printf("Failed to reconnect SNMP client for %s: %v", server, err)
                }
            }

            if fetchErr != nil {
                log.Printf("[ERROR] SNMP Get failed from %s after 3 retries: %v", server, fetchErr)
                return
            }

            if resultVal == nil {
                log.Printf("[WARNING] SNMP response is nil from %s", server)
                return
            }

            // Kiểm tra kiểu dữ liệu trả về
            var result int64
            switch v := resultVal.(type) {
            case int64:
                result = v
            case int:
                result = int64(v)
            case uint:
                result = int64(v)
            case uint64:
                result = int64(v)
            case float64:
                result = int64(v)
            case string:
                log.Printf("[INFO] SNMP data from %s is string: %s", server, v)
                return
            case []byte:
                log.Printf("[INFO] SNMP data from %s is []byte: %s", server, string(v))
                return
            default:
                log.Printf("[ERROR] Invalid SNMP data type from %s: %T", server, v)
                return
            }

            results <- result // Đưa dữ liệu vào channel
        }(server)
    }

    // Đợi tất cả goroutines hoàn thành
    wg.Wait()
    close(results)

    // Lấy dữ liệu từ channel
    for res := range results {
        sum = append(sum, res)
    }

    return sum, svs
}

func GetIoReceiveData(c *gin.Context) ([]int64, []string) {
	var sum []int64
	var svs []string
	var wg sync.WaitGroup
	results := make(chan int64, len(config.CMP_Servers))    // Channel để thu thập kết quả
	semaphore := make(chan struct{}, maxConcurrentSNMP) // Giới hạn 20 request SNMP chạy đồng thời

	for _, server := range config.CMP_Servers {
		wg.Add(1)
		semaphore <- struct{}{} // Giữ chỗ trong semaphore

		go func(server string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Giải phóng chỗ trong semaphore

			client, err := getSNMPClient(server)
			if err != nil {
				log.Printf("Failed to get SNMP client for %s: %v", server, err)
				return
			}

			var resultVal interface{}
			for i := 0; i < 3; i++ { // Thử lại tối đa 3 lần
				result, err := client.Get([]string{config.IOReceive})
				if err == nil && len(result.Variables) > 0 {
					svs = append(svs, server)
					resultVal = result.Variables[0].Value
					break
				}
				log.Printf("[WARNING] Retry %d: Failed SNMP from %s\n", i+1, server)
				time.Sleep(500 * time.Millisecond)
			}

			if err != nil {
				log.Printf("[ERROR] SNMP Get failed from %s: %v\n", server, err)
				return
			}

			if resultVal == nil {
				log.Printf("[WARNING] SNMP response is nil from %s\n", server)
				return
			}

			// Kiểm tra kiểu dữ liệu trả về
			var result int64
			switch v := resultVal.(type) {
			case int64:
				result = v
			case int:
				result = int64(v)
			case float64:
				result = int64(v)
			case string:
				log.Printf("[INFO] SNMP data from %s is string: %s\n", server, v)
				return
			case []byte:
				log.Printf("[INFO] SNMP data from %s is []byte: %s\n", server, string(v))
				return
			default:
				log.Printf("[ERROR] Invalid SNMP data type from %s: %T\n", server, v)
				return
			}

			results <- result // Đưa dữ liệu vào channel
		}(server)
	}

	// Đợi tất cả goroutines hoàn thành
	wg.Wait()
	close(results)

	// Lấy dữ liệu từ channel
	for res := range results {
		sum = append(sum, res)
	}

	return sum, svs
}

func GetIoSentData(c *gin.Context) ([]int64, []string) {
	var sum []int64
	var svs []string
	var wg sync.WaitGroup
	results := make(chan int64, len(config.CMP_Servers))    // Channel để thu thập kết quả
	semaphore := make(chan struct{}, maxConcurrentSNMP) // Giới hạn 20 request SNMP chạy đồng thời

	for _, server := range config.CMP_Servers {
		wg.Add(1)
		semaphore <- struct{}{} // Giữ chỗ trong semaphore

		go func(server string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Giải phóng chỗ trong semaphore

			client, err := getSNMPClient(server)
			if err != nil {
				log.Printf("Failed to get SNMP client for %s: %v", server, err)
				return
			}

			var resultVal interface{}
			for i := 0; i < 3; i++ { // Thử lại tối đa 3 lần
				result, err := client.Get([]string{config.IOSent})
				if err == nil && len(result.Variables) > 0 {
					svs = append(svs, server)
					resultVal = result.Variables[0].Value
					break
				}
				log.Printf("[WARNING] Retry %d: Failed SNMP from %s\n", i+1, server)
				time.Sleep(500 * time.Millisecond)
			}

			if err != nil {
				log.Printf("[ERROR] SNMP Get failed from %s: %v\n", server, err)
				return
			}

			if resultVal == nil {
				log.Printf("[WARNING] SNMP response is nil from %s\n", server)
				return
			}

			// Kiểm tra kiểu dữ liệu trả về
			var result int64
			switch v := resultVal.(type) {
			case int64:
				result = v
			case int:
				result = int64(v)
			case float64:
				result = int64(v)
			case string:
				log.Printf("[INFO] SNMP data from %s is string: %s\n", server, v)
				return
			case []byte:
				log.Printf("[INFO] SNMP data from %s is []byte: %s\n", server, string(v))
				return
			default:
				log.Printf("[ERROR] Invalid SNMP data type from %s: %T\n", server, v)
				return
			}

			results <- result // Đưa dữ liệu vào channel
		}(server)
	}

	// Đợi tất cả goroutines hoàn thành
	wg.Wait()
	close(results)

	// Lấy dữ liệu từ channel
	for res := range results {
		sum = append(sum, res)
	}

	return sum, svs
}

func GetMemSwapTotal(c *gin.Context) ([]int64, []string) {
	var sum []int64
	var svs []string
	var wg sync.WaitGroup
	results := make(chan int64, len(config.CMP_Servers))    // Channel để thu thập kết quả
	semaphore := make(chan struct{}, maxConcurrentSNMP) // Giới hạn 20 request SNMP chạy đồng thời

	for _, server := range config.Servers {
		wg.Add(1)
		semaphore <- struct{}{} // Giữ chỗ trong semaphore

		go func(server string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Giải phóng chỗ trong semaphore

			client, err := getSNMPClient(server)
			if err != nil {
				log.Printf("Failed to get SNMP client for %s: %v", server, err)
				return
			}

			var resultVal interface{}
			for i := 0; i < 3; i++ { // Thử lại tối đa 3 lần
				result, err := client.Get([]string{config.MemSwapTotal})
				if err == nil && len(result.Variables) > 0 {
					svs = append(svs, server)
					resultVal = result.Variables[0].Value
					break
				}
				log.Printf("[WARNING] Retry %d: Failed SNMP from %s\n", i+1, server)
				time.Sleep(500 * time.Millisecond)
			}

			if err != nil {
				log.Printf("[ERROR] SNMP Get failed from %s: %v\n", server, err)
				return
			}

			if resultVal == nil {
				log.Printf("[WARNING] SNMP response is nil from %s\n", server)
				return
			}

			// Kiểm tra kiểu dữ liệu trả về
			var result int64
			switch v := resultVal.(type) {
			case int64:
				result = v
			case int:
				result = int64(v)
			case float64:
				result = int64(v)
			case string:
				log.Printf("[INFO] SNMP data from %s is string: %s\n", server, v)
				return
			case []byte:
				log.Printf("[INFO] SNMP data from %s is []byte: %s\n", server, string(v))
				return
			default:
				log.Printf("[ERROR] Invalid SNMP data type from %s: %T\n", server, v)
				return
			}

			results <- result // Đưa dữ liệu vào channel
		}(server)
	}

	// Đợi tất cả goroutines hoàn thành
	wg.Wait()
	close(results)

	// Lấy dữ liệu từ channel
	for res := range results {
		sum = append(sum, res)
	}

	return sum, svs
}

func GetMemSwapAvail(c *gin.Context) ([]int64, []string) {
	var sum []int64
	var svs []string
	var wg sync.WaitGroup
	results := make(chan int64, len(config.CMP_Servers))    // Channel để thu thập kết quả
	semaphore := make(chan struct{}, maxConcurrentSNMP) // Giới hạn 20 request SNMP chạy đồng thời

	for _, server := range config.Servers {
		wg.Add(1)
		semaphore <- struct{}{} // Giữ chỗ trong semaphore

		go func(server string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Giải phóng chỗ trong semaphore

			client, err := getSNMPClient(server)
			if err != nil {
				log.Printf("Failed to get SNMP client for %s: %v", server, err)
				return
			}

			var resultVal interface{}
			for i := 0; i < 3; i++ { // Thử lại tối đa 3 lần
				result, err := client.Get([]string{config.MemSwapAvail})
				if err == nil && len(result.Variables) > 0 {
					svs = append(svs, server)
					resultVal = result.Variables[0].Value
					break
				}
				log.Printf("[WARNING] Retry %d: Failed SNMP from %s\n", i+1, server)
				time.Sleep(500 * time.Millisecond)
			}

			if err != nil {
				log.Printf("[ERROR] SNMP Get failed from %s: %v\n", server, err)
				return
			}

			if resultVal == nil {
				log.Printf("[WARNING] SNMP response is nil from %s\n", server)
				return
			}

			// Kiểm tra kiểu dữ liệu trả về
			var result int64
			switch v := resultVal.(type) {
			case int64:
				result = v
			case int:
				result = int64(v)
			case float64:
				result = int64(v)
			case string:
				log.Printf("[INFO] SNMP data from %s is string: %s\n", server, v)
				return
			case []byte:
				log.Printf("[INFO] SNMP data from %s is []byte: %s\n", server, string(v))
				return
			default:
				log.Printf("[ERROR] Invalid SNMP data type from %s: %T\n", server, v)
				return
			}

			results <- result // Đưa dữ liệu vào channel
		}(server)
	}

	// Đợi tất cả goroutines hoàn thành
	wg.Wait()
	close(results)

	// Lấy dữ liệu từ channel
	for res := range results {
		sum = append(sum, res)
	}

	return sum, svs
}

func GetCpuPercentUser(c *gin.Context) ([]int64, []string) {
	var sum []int64
	var svs []string
	var wg sync.WaitGroup
	results := make(chan int64, len(config.CMP_Servers))    // Channel để thu thập kết quả
	semaphore := make(chan struct{}, maxConcurrentSNMP) // Giới hạn 20 request SNMP chạy đồng thời

	for _, server := range config.CMP_Servers {
		wg.Add(1)
		semaphore <- struct{}{} // Giữ chỗ trong semaphore

		go func(server string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Giải phóng chỗ trong semaphore

			client, err := getSNMPClient(server)
			if err != nil {
				log.Printf("Failed to get SNMP client for %s: %v", server, err)
				return
			}

			var resultVal interface{}
			for i := 0; i < 3; i++ { // Thử lại tối đa 3 lần
				result, err := client.Get([]string{config.SsCpuUser})
				if err == nil && len(result.Variables) > 0 {
					svs = append(svs, server)
					resultVal = result.Variables[0].Value
					break
				}
				log.Printf("[WARNING] Retry %d: Failed SNMP from %s\n", i+1, server)
				time.Sleep(500 * time.Millisecond)
			}

			if err != nil {
				log.Printf("[ERROR] SNMP Get failed from %s: %v\n", server, err)
				return
			}

			if resultVal == nil {
				log.Printf("[WARNING] SNMP response is nil from %s\n", server)
				return
			}

			// Kiểm tra kiểu dữ liệu trả về
			var result int64
			switch v := resultVal.(type) {
			case int64:
				result = v
			case int:
				result = int64(v)
			case float64:
				result = int64(v)
			case string:
				log.Printf("[INFO] SNMP data from %s is string: %s\n", server, v)
				return
			case []byte:
				log.Printf("[INFO] SNMP data from %s is []byte: %s\n", server, string(v))
				return
			default:
				log.Printf("[ERROR] Invalid SNMP data type from %s: %T\n", server, v)
				return
			}

			results <- result // Đưa dữ liệu vào channel
		}(server)
	}

	// Đợi tất cả goroutines hoàn thành
	wg.Wait()
	close(results)

	// Lấy dữ liệu từ channel
	for res := range results {
		sum = append(sum, res)
	}

	return sum, svs
}

func GetCpuPercentSystem(c *gin.Context) ([]int64, []string) {
	var sum []int64
	var svs []string
	var wg sync.WaitGroup
	results := make(chan int64, len(config.CMP_Servers))    // Channel để thu thập kết quả
	semaphore := make(chan struct{}, maxConcurrentSNMP) // Giới hạn 20 request SNMP chạy đồng thời

	for _, server := range config.CMP_Servers {
		wg.Add(1)
		semaphore <- struct{}{} // Giữ chỗ trong semaphore

		go func(server string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Giải phóng chỗ trong semaphore

			client, err := getSNMPClient(server)
			if err != nil {
				log.Printf("Failed to get SNMP client for %s: %v", server, err)
				return
			}

			var resultVal interface{}
			for i := 0; i < 3; i++ { // Thử lại tối đa 3 lần
				result, err := client.Get([]string{config.SsCpuSystem})
				if err == nil && len(result.Variables) > 0 {
					svs = append(svs, server)
					resultVal = result.Variables[0].Value
					break
				}
				log.Printf("[WARNING] Retry %d: Failed SNMP from %s\n", i+1, server)
				time.Sleep(500 * time.Millisecond)
			}

			if err != nil {
				log.Printf("[ERROR] SNMP Get failed from %s: %v\n", server, err)
				return
			}

			if resultVal == nil {
				log.Printf("[WARNING] SNMP response is nil from %s\n", server)
				return
			}

			// Kiểm tra kiểu dữ liệu trả về
			var result int64
			switch v := resultVal.(type) {
			case int64:
				result = v
			case int:
				result = int64(v)
			case float64:
				result = int64(v)
			case string:
				log.Printf("[INFO] SNMP data from %s is string: %s\n", server, v)
				return
			case []byte:
				log.Printf("[INFO] SNMP data from %s is []byte: %s\n", server, string(v))
				return
			default:
				log.Printf("[ERROR] Invalid SNMP data type from %s: %T\n", server, v)
				return
			}

			results <- result // Đưa dữ liệu vào channel
		}(server)
	}

	// Đợi tất cả goroutines hoàn thành
	wg.Wait()
	close(results)

	// Lấy dữ liệu từ channel
	for res := range results {
		sum = append(sum, res)
	}

	return sum, svs
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
		c.Error(err) // CHỈ gọi c.Error nếu err != nil
		return nil
	}
	defer snmp.Manager.Conn.Close()

	result, err2 := snmp.Manager.Get(oids) // Lấy dữ liệu từ SNMP agent
	if err2 != nil {
		log.Printf("Get() err: %v", err2)
		c.Error(err2) // CHỈ gọi c.Error nếu err2 != nil
		return nil
	}

	// Kiểm tra nếu không có dữ liệu trả về
	if len(result.Variables) == 0 {
		log.Printf("SNMP response empty from %s", target)
		return nil
	}

	return result.Variables[0].Value
}

func GetPreviousDataSNMP(c *gin.Context, target string, oid string) interface{} {
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

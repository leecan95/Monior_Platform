package model

import (
	"Monitor_Platform/config"
	"github.com/gin-gonic/gin"

	"fmt"
	_ "github.com/lib/pq"
)

func (c *PostgreDb) QueryListDataPackageLatency() (config.LatencyKpi, error) {
	fmt.Print("query data")
	var data config.LatencyKpi
	rows, err := c.db.Query("SELECT * FROM mv_listDataPackage_latency_one_day")
	//rows, err := c.db.Query("WITH Percentile AS (\n    SELECT PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY latency) AS percentile_95\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/report/overview%'\n),\nStats AS (\n    SELECT\n        COUNT(*) AS total_records,\n        COUNT(*) FILTER (WHERE latency < 5000) AS count_latency_below_5,\n        (SELECT percentile_95 FROM Percentile) AS latency_95th_percentile\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/report/overview%'\n)\nSELECT\n    count_latency_below_5,\n    latency_95th_percentile,\n    total_records,\n    CASE\n        WHEN count_latency_below_5 > (0.95 * total_records) THEN 'Yes'\n        ELSE 'No'\n    END AS exceeds_95_percent\nFROM Stats")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var count string
		var latency string
		var total string
		var result string
		if err := rows.Scan(&count, &latency, &total, &result); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.LatencyKpi{
			Api:        "listDataPackage",
			Total:      total,
			Count:      count,
			Percentile: latency,
			Result:     result,
		}
		fmt.Printf("count : %s, latency: %s, total: %s, result: %s\n", count, latency, total, result)
	}
	return data, nil
}

func GetListDataPackageLatency(c *gin.Context) (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryListDataPackageLatency()
	return data, nil
}

func (c *PostgreDb) QuerylistSubscribeLatency() (config.LatencyKpi, error) {
	fmt.Print("query data")
	var data config.LatencyKpi
	rows, err := c.db.Query("SELECT * FROM mv_listSubscribe_latency_one_day")
	//rows, err := c.db.Query("WITH Percentile AS (\n    SELECT PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY latency) AS percentile_95\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/report/overview%'\n),\nStats AS (\n    SELECT\n        COUNT(*) AS total_records,\n        COUNT(*) FILTER (WHERE latency < 5000) AS count_latency_below_5,\n        (SELECT percentile_95 FROM Percentile) AS latency_95th_percentile\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/report/overview%'\n)\nSELECT\n    count_latency_below_5,\n    latency_95th_percentile,\n    total_records,\n    CASE\n        WHEN count_latency_below_5 > (0.95 * total_records) THEN 'Yes'\n        ELSE 'No'\n    END AS exceeds_95_percent\nFROM Stats")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var count string
		var latency string
		var total string
		var result string
		if err := rows.Scan(&count, &latency, &total, &result); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.LatencyKpi{
			Api:        "listSubscribe",
			Total:      total,
			Count:      count,
			Percentile: latency,
			Result:     result,
		}
		fmt.Printf("count : %s, latency: %s, total: %s, result: %s\n", count, latency, total, result)
	}
	return data, nil
}

func GetlistSubscribeLatency(c *gin.Context) (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QuerylistSubscribeLatency()
	return data, nil
}

func (c *PostgreDb) QuerylistBalanceInfoLatency() (config.LatencyKpi, error) {
	fmt.Print("query data")
	var data config.LatencyKpi
	rows, err := c.db.Query("SELECT * FROM mv_listBalanceInfo_latency_one_day")
	//rows, err := c.db.Query("WITH Percentile AS (\n    SELECT PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY latency) AS percentile_95\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/report/overview%'\n),\nStats AS (\n    SELECT\n        COUNT(*) AS total_records,\n        COUNT(*) FILTER (WHERE latency < 5000) AS count_latency_below_5,\n        (SELECT percentile_95 FROM Percentile) AS latency_95th_percentile\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/report/overview%'\n)\nSELECT\n    count_latency_below_5,\n    latency_95th_percentile,\n    total_records,\n    CASE\n        WHEN count_latency_below_5 > (0.95 * total_records) THEN 'Yes'\n        ELSE 'No'\n    END AS exceeds_95_percent\nFROM Stats")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var count string
		var latency string
		var total string
		var result string
		if err := rows.Scan(&count, &latency, &total, &result); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.LatencyKpi{
			Api:        "listBalanceInfo",
			Total:      total,
			Count:      count,
			Percentile: latency,
			Result:     result,
		}
		fmt.Printf("count : %s, latency: %s, total: %s, result: %s\n", count, latency, total, result)
	}
	return data, nil
}

func GetlistBalanceInfoLatency(c *gin.Context) (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QuerylistBalanceInfoLatency()
	return data, nil
}

func (c *PostgreDb) QueryktmiLatency() (config.LatencyKpi, error) {
	fmt.Print("query data")
	var data config.LatencyKpi
	rows, err := c.db.Query("SELECT * FROM mv_ktmi_latency_one_day")
	//rows, err := c.db.Query("WITH Percentile AS (\n    SELECT PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY latency) AS percentile_95\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/report/overview%'\n),\nStats AS (\n    SELECT\n        COUNT(*) AS total_records,\n        COUNT(*) FILTER (WHERE latency < 5000) AS count_latency_below_5,\n        (SELECT percentile_95 FROM Percentile) AS latency_95th_percentile\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/report/overview%'\n)\nSELECT\n    count_latency_below_5,\n    latency_95th_percentile,\n    total_records,\n    CASE\n        WHEN count_latency_below_5 > (0.95 * total_records) THEN 'Yes'\n        ELSE 'No'\n    END AS exceeds_95_percent\nFROM Stats")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var count string
		var latency string
		var total string
		var result string
		if err := rows.Scan(&count, &latency, &total, &result); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.LatencyKpi{
			Api:        "Ktmi",
			Total:      total,
			Count:      count,
			Percentile: latency,
			Result:     result,
		}
		fmt.Printf("count : %s, latency: %s, total: %s, result: %s\n", count, latency, total, result)
	}
	return data, nil
}

func GetktmiLatency(c *gin.Context) (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryktmiLatency()
	return data, nil
}

func (c *PostgreDb) QueryisdnvalidateLatency() (config.LatencyKpi, error) {
	fmt.Print("query data")
	var data config.LatencyKpi
	rows, err := c.db.Query("SELECT * FROM mv_isdnvalidate_latency_one_day")
	//rows, err := c.db.Query("WITH Percentile AS (\n    SELECT PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY latency) AS percentile_95\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/report/overview%'\n),\nStats AS (\n    SELECT\n        COUNT(*) AS total_records,\n        COUNT(*) FILTER (WHERE latency < 5000) AS count_latency_below_5,\n        (SELECT percentile_95 FROM Percentile) AS latency_95th_percentile\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/report/overview%'\n)\nSELECT\n    count_latency_below_5,\n    latency_95th_percentile,\n    total_records,\n    CASE\n        WHEN count_latency_below_5 > (0.95 * total_records) THEN 'Yes'\n        ELSE 'No'\n    END AS exceeds_95_percent\nFROM Stats")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var count string
		var latency string
		var total string
		var result string
		if err := rows.Scan(&count, &latency, &total, &result); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.LatencyKpi{
			Api:        "isdnValidate",
			Total:      total,
			Count:      count,
			Percentile: latency,
			Result:     result,
		}
		fmt.Printf("count : %s, latency: %s, total: %s, result: %s\n", count, latency, total, result)
	}
	return data, nil
}

func GetisdnvalidateLatency(c *gin.Context) (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryisdnvalidateLatency()
	return data, nil
}


package model

import (
	"Monitor_Platform/config"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Helper function để xử lý giá trị NULL từ database
func getStringValue(nullStr sql.NullString) string {
	if nullStr.Valid {
		return nullStr.String
	}
	return "0" // Trả về giá trị mặc định khi NULL
}

var Schema = `
CREATE TABLE person (
    first_name text,
    last_name text,
    email text
);

CREATE TABLE place (
    country text,
    city text NULL,
    telcode integer
)`

var Schema2 = `
CREATE TABLE cpuUsage (
    id SERIAL PRIMARY KEY,
    cpu_id INT,
    value FLOAT,
    timestamp timestamptz DEFAULT current_timestamp
)`

var Schema3 = `
CREATE TABLE kpimodule (
    pod text,
    value FLOAT,
    timestamp timestamptz DEFAULT current_timestamp
)`

type PostgreDb struct {
	db *sql.DB
}

type PostgreDB struct {
	pool *pgxpool.Pool
}
type CpuUsage struct {
	ID        int     `db:"id"`
	CpuID     int     `db:"cpu_id"`
	Value     float64 `db:"value"`
	Timestamp string  `db:"timestamp"`
}

type KpiRequestModule struct {
	Pod       string  `db:"pod"`
	Method    string  `db:"method"`
	Value     float64 `db:"value"`
	Timestamp string  `db:"timestamp"`
}
type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

type Place struct {
	Country string
	City    sql.NullString
	TelCode int
}

type ApiLog struct {
	Url        string
	StatusCode int
	Latency    int64
	Timestamp  time.Time
}

func LoadDBConfig() config.DBConfig {
	return config.DBConfig{
		Host:              config.Env(config.EnvDBHost, config.DefDBHost),
		Port:              config.Env(config.EnvDBPort, config.DefDBPort),
		User:              config.Env(config.EnvDBUser, config.DefDBUser),
		Pass:              config.Env(config.EnvDBPass, config.DefDBPass),
		Name:              config.Env(config.EnvDB, config.DefDB),
		SSLMode:           config.Env(config.EnvDBSSLMode, config.DefDBSSLMode),
		SSLCert:           config.Env(config.EnvDBSSLCert, config.DefDBSSLCert),
		SSLKey:            config.Env(config.EnvDBSSLKey, config.DefDBSSLKey),
		SSLRootcert:       config.Env(config.EnvDBSSLRootCert, config.DefDBSSLRootCert),
		ConnectionTimeout: int(100),
		MinConn:           int32(30),
		MaxConn:           int32(100),
		MaxRetry:          uint64(100),
		RetryInterval:     uint64(10),
	}
}

func Connect(cfg config.DBConfig) (*sqlx.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootcert)

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}
func ConnectToDb() *sqlx.DB {
	var cfg config.DBConfig
	cfg = LoadDBConfig()
	db, err := Connect(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func ConnectNewDB(cfg config.DBConfig) (*PostgreDb, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootcert)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgreDb{db}, nil
}

func ConnectToNewDb() *PostgreDb {
	var cfg config.DBConfig
	cfg = LoadDBConfig()
	db, err := ConnectNewDB(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func ConnectTransDB(cfg config.DBConfig) (*PostgreDb, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, "transactions", cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootcert)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgreDb{db}, nil
}

func ConnectDevicesDB(cfg config.DBConfig) (*PostgreDb, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, "devices", cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootcert)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgreDb{db}, nil
}

func ConnectAttributesDB(cfg config.DBConfig) (*PostgreDb, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, "attributes", cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootcert)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgreDb{db}, nil
}

func ConnectPoolDB(cfg config.DBConfig) (*PostgreDB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s pool_max_conns=10", cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name)
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		fmt.Printf("loi ", err)
		return nil, err
	}
	dbPool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Printf("loi ", err)
		return nil, err
	}
	if err := dbPool.Ping(context.Background()); err != nil {
		fmt.Printf("loi ", err)
		return nil, err
	}
	return &PostgreDB{dbPool}, nil
}

func (db *PostgreDB) Close() {
	db.pool.Close()
}

func ConnectToTransDb() *PostgreDb {
	var cfg config.DBConfig
	cfg = LoadDBConfig()
	db, err := ConnectTransDB(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func (c *PostgreDb) QueryData() error {
	fmt.Print("query data")
	rows, err := c.db.Query("select * from kpimodule")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var pod string
		var value string
		var time string
		if err := rows.Scan(&pod, &value, &time); err != nil {
			fmt.Printf("Loi scan db ", err)
			return err
		}
		fmt.Printf("Pod : %s, Value: %s, Time: %s\n", pod, value, time)
	}
	return nil
}

func (c *PostgreDb) QueryLatency() (config.LatencyKpi, error) {
	fmt.Print("query data")
	var data config.LatencyKpi
	rows, err := c.db.Query("SELECT * FROM mv_transactions_overall_latency")
	//rows, err := c.db.Query("WITH Percentile AS (\n    SELECT PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY latency) AS percentile_95\n    FROM transactions\n),\nStats AS (\n    SELECT\n        COUNT(*) AS total_records,\n        COUNT(*) FILTER (WHERE latency < 5000) AS count_latency_below_5\n    FROM transactions\n)\nSELECT\n    s.count_latency_below_5,\n    p.percentile_95 AS latency_95th_percentile,\n    s.total_records,\n    CASE\n        WHEN s.count_latency_below_5 > (0.95 * s.total_records) THEN 'Yes'\n        ELSE 'No'\n    END AS exceeds_95_percent\nFROM Stats s, Percentile p")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var count sql.NullString
		var latency sql.NullString
		var total sql.NullString
		var result sql.NullString
		if err := rows.Scan(&count, &latency, &total, &result); err != nil {
			fmt.Printf("Loi scan db %v", err)
			return data, err
		}
		data = config.LatencyKpi{
			Api:        "Total",
			Total:      getStringValue(total),
			Count:      getStringValue(count),
			Percentile: getStringValue(latency),
			Result:     getStringValue(result),
		}
		fmt.Printf("count : %s, latency: %s, total: %s, result: %s\n", data.Count, data.Percentile, data.Total, data.Result)
	}
	return data, nil
}
func (p *PostgreDB) QueryLatencyPool() (config.LatencyKpi, error) {
	var data config.LatencyKpi
	conn, err := p.pool.Acquire(context.Background())
	if err != nil {
		return data, err
	}
	defer conn.Release()
	var count sql.NullString
	var latency sql.NullString
	var total sql.NullString
	var result sql.NullString
	err = conn.QueryRow(context.Background(), `WITH Percentile AS (
SELECT PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY latency) AS percentile_95
    FROM transactions
)
SELECT
    (SELECT COUNT(*) FROM transactions WHERE latency < 5000) AS count_latency_below_5,
    (SELECT percentile_95 FROM Percentile) AS latency_95th_percentile,
    (SELECT COUNT(*) FROM transactions) AS total_records,
    CASE
        WHEN (SELECT COUNT(*) FROM transactions WHERE latency < 5000) > ((SELECT percentile_95 FROM Percentile) * 0.95 * (SELECT COUNT(*) FROM transactions)) THEN 'Yes'
        ELSE 'No'
    END AS exceeds_95_percent
FROM Percentile`).Scan(&count, &latency, &total, &result)
	if err != nil {
		return data, err
	}
	data = config.LatencyKpi{
		Api:        "Total",
		Total:      getStringValue(total),
		Count:      getStringValue(count),
		Percentile: getStringValue(latency),
		Result:     getStringValue(result),
	}
	fmt.Printf("count : %s, latency: %s, total: %s, result: %s\n", data.Count, data.Percentile, data.Total, data.Result)

	return data, nil
}

func GetOverallLatency(c *gin.Context) (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryLatency()
	return data, nil
}

func ExportOverallLatency() (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryLatency()
	return data, nil
}

func GetOverallLatencyPool(c *gin.Context) (config.LatencyKpi, error) {
	var data config.LatencyKpi
	dbpool, exists := c.Get("dbpool")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database pool not found"})
	}
	pool := dbpool.(PostgreDB)
	data, _ = pool.QueryLatencyPool()
	return data, nil
}

// Get transactions data


func (c *PostgreDb) QueryTransactions(url string, offset, limit int) ([]config.TransactionQuery, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if url != "" {
		rows, err = c.db.Query(`
			SELECT
				method,
				url,
				status_code,
				response_body,
				service_name,
				latency,
				ts,
				user_id,
				project_id
			FROM transactions
			WHERE url LIKE $1
			ORDER BY ts DESC
			LIMIT $2 OFFSET $3
		`, "%"+url+"%", limit, offset)
	} else {
		rows, err = c.db.Query(`
			SELECT
				method,
				url,
				status_code,
				response_body,
				service_name,
				latency,
				ts,
				user_id,
				project_id
			FROM transactions
			ORDER BY ts DESC
			LIMIT $1 OFFSET $2
		`, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]config.TransactionQuery, 0, limit)

	for rows.Next() {
		var tx config.TransactionQuery
		if err := rows.Scan(
			&tx.Method,
			&tx.URL,
			&tx.StatusCode,
			&tx.ResponseBody,
			&tx.ServiceName,
			&tx.Latency,
			&tx.TS,
			&tx.UserID,
			&tx.ProjectID,
		); err != nil {
			fmt.Println("Error scanning row", err)
			return nil, err
		}
		results = append(results, tx)
		fmt.Printf("url : %s, response_body: %s, latency: %d, status_code: %d\n", tx.URL, tx.ResponseBody, tx.Latency, tx.StatusCode)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating rows", err)
		return nil, err
	}

	return results, nil
}


func GetTransactions(c *gin.Context, url string, offset int, limit int) ([]config.TransactionQuery, error) {
	var cfg config.DBConfig
	var data []config.TransactionQuery
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryTransactions(url, offset, limit)
	return data, nil
}

func (c *PostgreDb) QueryLatencyPercentileByURL(
	url string,
) (config.LatencyPercentileResult, error) {

	var result config.LatencyPercentileResult

	query := `
	WITH Percentile AS (
		SELECT 
			PERCENTILE_CONT(0.95) 
			WITHIN GROUP (ORDER BY latency) AS percentile_95
		FROM transactions
		WHERE url LIKE $1
		  AND ts >= EXTRACT(EPOCH FROM CURRENT_DATE) * 1000000
		  AND ts < EXTRACT(EPOCH FROM (CURRENT_DATE + INTERVAL '1 day')) * 1000000
	),
	Stats AS (
		SELECT
			COUNT(*) AS total_records,
			COUNT(*) FILTER (WHERE latency < 5000) AS count_latency_below_5
		FROM transactions
		WHERE url LIKE $1
		  AND ts >= EXTRACT(EPOCH FROM CURRENT_DATE) * 1000000
		  AND ts < EXTRACT(EPOCH FROM (CURRENT_DATE + INTERVAL '1 day')) * 1000000
	)
	SELECT
		s.count_latency_below_5,
		p.percentile_95,
		s.total_records,
		CASE
			WHEN s.count_latency_below_5 > (0.95 * s.total_records) 
			THEN 'Yes'
			ELSE 'No'
		END AS exceeds_95_percent
	FROM Stats s, Percentile p
	`

	row := c.db.QueryRow(query, "%"+url+"%")

	err := row.Scan(
		&result.CountBelow5,
		&result.Percentile95,
		&result.TotalRecords,
		&result.Exceeds95Pct,
	)
	if err != nil {
		return result, err
	}

	return result, nil
}


func (c *PostgreDb) QueryLoginLatency() (config.LatencyKpi, error) {
	fmt.Print("query data")
	var data config.LatencyKpi
	rows, err := c.db.Query("SELECT * FROM mv_transactions_login_latency")
	//rows, err := c.db.Query("WITH Percentile AS (\n    SELECT PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY latency) AS percentile_95\n    FROM transactions\n    WHERE url = '/app/vtracking/login'\n),\nStats AS (\n    SELECT\n        COUNT(*) AS total_records,\n        COUNT(*) FILTER (WHERE latency < 5000) AS count_latency_below_5,\n        (SELECT percentile_95 FROM Percentile) AS latency_95th_percentile\n    FROM transactions\n    WHERE url = '/app/vtracking/login'\n)\nSELECT\n    count_latency_below_5,\n    latency_95th_percentile,\n    total_records,\n    CASE\n        WHEN count_latency_below_5 > (0.95 * total_records) THEN 'Yes'\n        ELSE 'No'\n    END AS exceeds_95_percent\nFROM Stats")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var count sql.NullString
		var latency sql.NullString
		var total sql.NullString
		var result sql.NullString
		if err := rows.Scan(&count, &latency, &total, &result); err != nil {
			fmt.Printf("Loi scan db %v", err)
			return data, err
		}
		data = config.LatencyKpi{
			Api:        "Login",
			Total:      getStringValue(total),
			Count:      getStringValue(count),
			Percentile: getStringValue(latency),
			Result:     getStringValue(result),
		}
		fmt.Printf("count : %s, latency: %s, total: %s, result: %s\n", data.Count, data.Percentile, data.Total, data.Result)
	}
	return data, nil
}

func GetLoginLatency(c *gin.Context) (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryLoginLatency()
	return data, nil
}

func ExportLoginLatency() (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryLoginLatency()
	return data, nil
}

func (c *PostgreDb) QueryReportLatency() (config.LatencyKpi, error) {
	fmt.Print("query data")
	var data config.LatencyKpi
	rows, err := c.db.Query("SELECT * FROM mv_transactions_attributes_stats")
	//rows, err := c.db.Query("WITH Percentile AS (\n    SELECT PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY latency) AS percentile_95\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/report/overview%'\n),\nStats AS (\n    SELECT\n        COUNT(*) AS total_records,\n        COUNT(*) FILTER (WHERE latency < 5000) AS count_latency_below_5,\n        (SELECT percentile_95 FROM Percentile) AS latency_95th_percentile\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/report/overview%'\n)\nSELECT\n    count_latency_below_5,\n    latency_95th_percentile,\n    total_records,\n    CASE\n        WHEN count_latency_below_5 > (0.95 * total_records) THEN 'Yes'\n        ELSE 'No'\n    END AS exceeds_95_percent\nFROM Stats")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var count sql.NullString
		var latency sql.NullString
		var total sql.NullString
		var result sql.NullString
		if err := rows.Scan(&count, &latency, &total, &result); err != nil {
			fmt.Printf("Loi scan db %v", err)
			return data, err
		}
		data = config.LatencyKpi{
			Api:        "Report",
			Total:      getStringValue(total),
			Count:      getStringValue(count),
			Percentile: getStringValue(latency),
			Result:     getStringValue(result),
		}
		fmt.Printf("count : %s, latency: %s, total: %s, result: %s\n", data.Count, data.Percentile, data.Total, data.Result)
	}
	return data, nil
}

func GetReportLatency(c *gin.Context) (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryReportLatency()
	return data, nil
}

func ExportReportLatency() (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryReportLatency()
	return data, nil
}

func (c *PostgreDb) QueryGetimageLatency() (config.LatencyKpi, error) {
	fmt.Print("query data")
	var data config.LatencyKpi
	rows, err := c.db.Query("SELECT * FROM mv_transactions_images_latency")
	//rows, err := c.db.Query("WITH Percentile AS (\n    SELECT PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY latency) AS percentile_95\n    FROM transactions\n    WHERE url LIKE '/s3%'\n),\nStats AS (\n    SELECT\n        COUNT(*) AS total_records,\n        COUNT(*) FILTER (WHERE latency < 5000) AS count_latency_below_5,\n        (SELECT percentile_95 FROM Percentile) AS latency_95th_percentile\n    FROM transactions\n    WHERE url LIKE '/s3%'\n)\nSELECT\n    count_latency_below_5,\n    latency_95th_percentile,\n    total_records,\n    CASE\n        WHEN count_latency_below_5 > (0.95 * total_records) THEN 'Yes'\n        ELSE 'No'\n    END AS exceeds_95_percent\nFROM Stats;\n")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var count sql.NullString
		var latency sql.NullString
		var total sql.NullString
		var result sql.NullString
		if err := rows.Scan(&count, &latency, &total, &result); err != nil {
			fmt.Printf("Loi scan db %v", err)
			return data, err
		}
		data = config.LatencyKpi{
			Api:        "Get_images",
			Total:      getStringValue(total),
			Count:      getStringValue(count),
			Percentile: getStringValue(latency),
			Result:     getStringValue(result),
		}
		fmt.Printf("count : %s, latency: %s, total: %s, result: %s\n", data.Count, data.Percentile, data.Total, data.Result)
	}
	return data, nil
}

func GetImageLatency(c *gin.Context) (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryGetimageLatency()
	return data, nil
}

func ExportImageLatency() (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryGetimageLatency()
	return data, nil
}

func (c *PostgreDb) QueryTrackingLatency() (config.LatencyKpi, error) {
	fmt.Print("query data")
	var data config.LatencyKpi
	rows, err := c.db.Query("SELECT * FROM mv_transactions_tracking_latency")
	//rows, err := c.db.Query("WITH Percentile AS (\n    SELECT PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY latency) AS percentile_95\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/logged/VEHICLE/%'\n),\nStats AS (\n    SELECT\n        COUNT(*) AS total_records,\n        COUNT(*) FILTER (WHERE latency < 5000) AS count_latency_below_5,\n        (SELECT percentile_95 FROM Percentile) AS latency_95th_percentile\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/logged/VEHICLE/%'\n)\nSELECT\n    count_latency_below_5,\n    latency_95th_percentile,\n    total_records,\n    CASE\n        WHEN count_latency_below_5 > (0.95 * total_records) THEN 'Yes'\n        ELSE 'No'\n    END AS exceeds_95_percent\nFROM Stats")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var count sql.NullString
		var latency sql.NullString
		var total sql.NullString
		var result sql.NullString
		if err := rows.Scan(&count, &latency, &total, &result); err != nil {
			fmt.Printf("Loi scan db %v", err)
			return data, err
		}
		data = config.LatencyKpi{
			Api:        "Tracking",
			Total:      getStringValue(total),
			Count:      getStringValue(count),
			Percentile: getStringValue(latency),
			Result:     getStringValue(result),
		}
		fmt.Printf("count : %s, latency: %s, total: %s, result: %s\n", data.Count, data.Percentile, data.Total, data.Result)
	}
	return data, nil
}

func GetTrackingLatency(c *gin.Context) (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryTrackingLatency()
	return data, nil
}

// *** get latency one_day ***//
func (c *PostgreDb) QueryTrackingLatencyOneDay() (config.LatencyKpi, error) {
	fmt.Print("query data")
	var data config.LatencyKpi
	rows, err := c.db.Query("SELECT * FROM mv_tracking_latency_one_day")
	//rows, err := c.db.Query("WITH Percentile AS (\n    SELECT PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY latency) AS percentile_95\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/logged/VEHICLE/%'\n),\nStats AS (\n    SELECT\n        COUNT(*) AS total_records,\n        COUNT(*) FILTER (WHERE latency < 5000) AS count_latency_below_5,\n        (SELECT percentile_95 FROM Percentile) AS latency_95th_percentile\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/logged/VEHICLE/%'\n)\nSELECT\n    count_latency_below_5,\n    latency_95th_percentile,\n    total_records,\n    CASE\n        WHEN count_latency_below_5 > (0.95 * total_records) THEN 'Yes'\n        ELSE 'No'\n    END AS exceeds_95_percent\nFROM Stats")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var count sql.NullString
		var latency sql.NullString
		var total sql.NullString
		var result sql.NullString
		if err := rows.Scan(&count, &latency, &total, &result); err != nil {
			fmt.Printf("Loi scan db %v", err)
			return data, err
		}
		data = config.LatencyKpi{
			Api:        "Lấy dữ liệu hành trình",
			Total:      getStringValue(total),
			Count:      getStringValue(count),
			Percentile: getStringValue(latency),
			Result:     getStringValue(result),
		}
		fmt.Printf("count : %s, latency: %s, total: %s, result: %s\n", data.Count, data.Percentile, data.Total, data.Result)
	}
	return data, nil
}

func GetTrackingLatencyOneDay() (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryTrackingLatencyOneDay()
	return data, nil
}

func (c *PostgreDb) QueryLoginLatencyOneDay() (config.LatencyKpi, error) {
	fmt.Print("query data")
	var data config.LatencyKpi
	rows, err := c.db.Query("SELECT * FROM mv_login_latency_one_day")
	//rows, err := c.db.Query("WITH Percentile AS (\n    SELECT PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY latency) AS percentile_95\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/logged/VEHICLE/%'\n),\nStats AS (\n    SELECT\n        COUNT(*) AS total_records,\n        COUNT(*) FILTER (WHERE latency < 5000) AS count_latency_below_5,\n        (SELECT percentile_95 FROM Percentile) AS latency_95th_percentile\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/logged/VEHICLE/%'\n)\nSELECT\n    count_latency_below_5,\n    latency_95th_percentile,\n    total_records,\n    CASE\n        WHEN count_latency_below_5 > (0.95 * total_records) THEN 'Yes'\n        ELSE 'No'\n    END AS exceeds_95_percent\nFROM Stats")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var count sql.NullString
		var latency sql.NullString
		var total sql.NullString
		var result sql.NullString
		if err := rows.Scan(&count, &latency, &total, &result); err != nil {
			fmt.Printf("Loi scan db %v", err)
			return data, err
		}
		data = config.LatencyKpi{
			Api:        "Đăng nhập hệ thống",
			Total:      getStringValue(total),
			Count:      getStringValue(count),
			Percentile: getStringValue(latency),
			Result:     getStringValue(result),
		}
		fmt.Printf("count : %s, latency: %s, total: %s, result: %s\n", data.Count, data.Percentile, data.Total, data.Result)
	}
	return data, nil
}

func GetLoginLatencyOneDay() (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryLoginLatencyOneDay()
	return data, nil
}

func (c *PostgreDb) QueryImageLatencyOneDay() (config.LatencyKpi, error) {
	fmt.Print("query data")
	var data config.LatencyKpi
	rows, err := c.db.Query("SELECT * FROM mv_images_latency_one_day")
	//rows, err := c.db.Query("WITH Percentile AS (\n    SELECT PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY latency) AS percentile_95\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/logged/VEHICLE/%'\n),\nStats AS (\n    SELECT\n        COUNT(*) AS total_records,\n        COUNT(*) FILTER (WHERE latency < 5000) AS count_latency_below_5,\n        (SELECT percentile_95 FROM Percentile) AS latency_95th_percentile\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/logged/VEHICLE/%'\n)\nSELECT\n    count_latency_below_5,\n    latency_95th_percentile,\n    total_records,\n    CASE\n        WHEN count_latency_below_5 > (0.95 * total_records) THEN 'Yes'\n        ELSE 'No'\n    END AS exceeds_95_percent\nFROM Stats")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var count sql.NullString
		var latency sql.NullString
		var total sql.NullString
		var result sql.NullString
		if err := rows.Scan(&count, &latency, &total, &result); err != nil {
			fmt.Printf("Loi scan db %v", err)
			return data, err
		}
		data = config.LatencyKpi{
			Api:        "Lấy dữ liệu hình ảnh",
			Total:      getStringValue(total),
			Count:      getStringValue(count),
			Percentile: getStringValue(latency),
			Result:     getStringValue(result),
		}
		fmt.Printf("count : %s, latency: %s, total: %s, result: %s\n", data.Count, data.Percentile, data.Total, data.Result)
	}
	return data, nil
}

func GetImageLatencyOneDay() (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryImageLatencyOneDay()
	return data, nil
}

func (c *PostgreDb) QueryDashboardLatencyOneDay() (config.LatencyKpi, error) {
	fmt.Print("query data")
	var data config.LatencyKpi
	rows, err := c.db.Query("SELECT * FROM mv_dashboard_latency_one_day")
	//rows, err := c.db.Query("WITH Percentile AS (\n    SELECT PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY latency) AS percentile_95\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/logged/VEHICLE/%'\n),\nStats AS (\n    SELECT\n        COUNT(*) AS total_records,\n        COUNT(*) FILTER (WHERE latency < 5000) AS count_latency_below_5,\n        (SELECT percentile_95 FROM Percentile) AS latency_95th_percentile\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/logged/VEHICLE/%'\n)\nSELECT\n    count_latency_below_5,\n    latency_95th_percentile,\n    total_records,\n    CASE\n        WHEN count_latency_below_5 > (0.95 * total_records) THEN 'Yes'\n        ELSE 'No'\n    END AS exceeds_95_percent\nFROM Stats")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var count sql.NullString
		var latency sql.NullString
		var total sql.NullString
		var result sql.NullString
		if err := rows.Scan(&count, &latency, &total, &result); err != nil {
			fmt.Printf("Loi scan db %v", err)
			return data, err
		}
		data = config.LatencyKpi{
			Api:        "Tải bảng ứng dụng",
			Total:      getStringValue(total),
			Count:      getStringValue(count),
			Percentile: getStringValue(latency),
			Result:     getStringValue(result),
		}
		fmt.Printf("count : %s, latency: %s, total: %s, result: %s\n", data.Count, data.Percentile, data.Total, data.Result)
	}
	return data, nil
}

func GetDashboardLatencyOneDay() (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryDashboardLatencyOneDay()
	return data, nil
}

func (c *PostgreDb) QueryTotalLatencyOneDay() (config.LatencyKpi, error) {
	fmt.Print("query data")
	var data config.LatencyKpi
	rows, err := c.db.Query("SELECT * FROM mv_total_latency_one_day")
	//rows, err := c.db.Query("WITH Percentile AS (\n    SELECT PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY latency) AS percentile_95\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/logged/VEHICLE/%'\n),\nStats AS (\n    SELECT\n        COUNT(*) AS total_records,\n        COUNT(*) FILTER (WHERE latency < 5000) AS count_latency_below_5,\n        (SELECT percentile_95 FROM Percentile) AS latency_95th_percentile\n    FROM transactions\n    WHERE url LIKE '/attributes/vtracking/logged/VEHICLE/%'\n)\nSELECT\n    count_latency_below_5,\n    latency_95th_percentile,\n    total_records,\n    CASE\n        WHEN count_latency_below_5 > (0.95 * total_records) THEN 'Yes'\n        ELSE 'No'\n    END AS exceeds_95_percent\nFROM Stats")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var count sql.NullString
		var latency sql.NullString
		var total sql.NullString
		var result sql.NullString
		if err := rows.Scan(&count, &latency, &total, &result); err != nil {
			fmt.Printf("Loi scan db %v", err)
			return data, err
		}
		data = config.LatencyKpi{
			Api:        "Tổng quát",
			Total:      getStringValue(total),
			Count:      getStringValue(count),
			Percentile: getStringValue(latency),
			Result:     getStringValue(result),
		}
		fmt.Printf("count : %s, latency: %s, total: %s, result: %s\n", data.Count, data.Percentile, data.Total, data.Result)
	}
	return data, nil
}

func GetTotalLatencyOneDay() (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryTotalLatencyOneDay()
	return data, nil
}

func ExportTrackingLatency() (config.LatencyKpi, error) {
	var cfg config.DBConfig
	var data config.LatencyKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	data, _ = db2.QueryTrackingLatency()
	return data, nil
}

func (c *PostgreDb) QueryRequestCount() (config.SuccessKpi, error) {
	fmt.Print("query request data")
	var data config.SuccessKpi
	//rows, err := c.db.Query("SELECT COUNT(*)\nFROM transactions")
	rows, err := c.db.Query("SELECT * \nFROM mv_transactions_overall_request")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.SuccessKpi{
			Api:   "Total",
			Value: total,
		}
	}
	return data, nil
}

func (c *PostgreDb) QuerySuccessRequestCount() (config.SuccessKpi, error) {
	fmt.Print("query request data")
	var data config.SuccessKpi
	//rows, err := c.db.Query("SELECT COUNT(*)\nFROM transactions WHERE response_body NOT LIKE '\"code\"%'")
	rows, err := c.db.Query("SELECT * \nFROM mv_transactions_overall_request_success")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.SuccessKpi{
			Api:   "Total",
			Value: total,
		}
	}
	return data, nil
}

func GetRequestTotal(c *gin.Context) (config.SuccessKpi, error) {
	var cfg config.DBConfig
	var total, count, result config.SuccessKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	total, _ = db2.QueryRequestCount()
	count, _ = db2.QuerySuccessRequestCount()
	result = config.SuccessKpi{
		Api:   total.Api,
		Value: (count.Value / total.Value) * 100,
	}
	return result, nil
}

func ExportRequestTotal() (config.SuccessKpi, error) {
	var cfg config.DBConfig
	var total, count, result config.SuccessKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	total, _ = db2.QueryRequestCount()
	count, _ = db2.QuerySuccessRequestCount()
	result = config.SuccessKpi{
		Api:   total.Api,
		Value: (count.Value / total.Value) * 100,
	}
	return result, nil
}

func (c *PostgreDb) QueryRequestLoginCount() (config.SuccessKpi, error) {
	fmt.Print("query request data")
	var data config.SuccessKpi
	//rows, err := c.db.Query("SELECT COUNT(*)\nFROM transactions\nWHERE url = '/app/vtracking/login'")
	rows, err := c.db.Query("SELECT * \nFROM mv_transactions_login_request")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.SuccessKpi{
			Api:   "Login",
			Value: total,
		}
	}
	return data, nil
}

func (c *PostgreDb) QuerySuccessRequestLoginCount() (config.SuccessKpi, error) {
	fmt.Print("query request data")
	var data config.SuccessKpi
	//rows, err := c.db.Query("SELECT COUNT(*)\nFROM transactions\nWHERE response_body NOT LIKE '\"code\":0%'\n    AND url = '/app/vtracking/login'")
	rows, err := c.db.Query("SELECT * \nFROM mv_transactions_login_request_success")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.SuccessKpi{
			Api:   "Login",
			Value: total,
		}
	}
	return data, nil
}

func GetRequestLoginTotal(c *gin.Context) (config.SuccessKpi, error) {
	var cfg config.DBConfig
	var total, count, result config.SuccessKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	total, _ = db2.QueryRequestLoginCount()
	count, _ = db2.QuerySuccessRequestLoginCount()
	result = config.SuccessKpi{
		Api:   total.Api,
		Value: (count.Value / total.Value) * 100,
	}
	return result, nil
}

func ExportRequestLoginTotal() (config.SuccessKpi, error) {
	var cfg config.DBConfig
	var total, count, result config.SuccessKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	total, _ = db2.QueryRequestLoginCount()
	count, _ = db2.QuerySuccessRequestLoginCount()
	result = config.SuccessKpi{
		Api:   total.Api,
		Value: (count.Value / total.Value) * 100,
	}
	return result, nil
}

func (c *PostgreDb) QueryRequestReportCount() (config.SuccessKpi, error) {
	fmt.Print("query request data")
	var data config.SuccessKpi
	//rows, err := c.db.Query("SELECT COUNT(*)\nFROM transactions\n WHERE url LIKE '/attributes/vtracking/report/overview%'")
	rows, err := c.db.Query("SELECT * \nFROM mv_transactions_report_request")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.SuccessKpi{
			Api:   "Report",
			Value: total,
		}
	}
	return data, nil
}

func (c *PostgreDb) QuerySuccessRequestReportCount() (config.SuccessKpi, error) {
	fmt.Print("query request data")
	var data config.SuccessKpi
	//rows, err := c.db.Query("SELECT COUNT(*)\nFROM transactions WHERE url LIKE '/attributes/vtracking/report/overview%' AND response_body NOT LIKE '\"code\"%'")
	rows, err := c.db.Query("SELECT * \nFROM mv_transactions_report_request_success")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.SuccessKpi{
			Api:   "Report",
			Value: total,
		}
	}
	return data, nil
}

func GetRequestReportTotal(c *gin.Context) (config.SuccessKpi, error) {
	var cfg config.DBConfig
	var total, count, result config.SuccessKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	total, _ = db2.QueryRequestReportCount()
	count, _ = db2.QuerySuccessRequestReportCount()
	result = config.SuccessKpi{
		Api:   total.Api,
		Value: (count.Value / total.Value) * 100,
	}
	return result, nil
}

func ExportRequestReportTotal() (config.SuccessKpi, error) {
	var cfg config.DBConfig
	var total, count, result config.SuccessKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	total, _ = db2.QueryRequestReportCount()
	count, _ = db2.QuerySuccessRequestReportCount()
	result = config.SuccessKpi{
		Api:   total.Api,
		Value: (count.Value / total.Value) * 100,
	}
	return result, nil
}

func (c *PostgreDb) QueryRequestGetImageCount() (config.SuccessKpi, error) {
	fmt.Print("query request data")
	var data config.SuccessKpi
	//rows, err := c.db.Query("SELECT COUNT(*)\nFROM transactions\n WHERE url LIKE '/s3%'")
	rows, err := c.db.Query("SELECT * \nFROM mv_transactions_images_request")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.SuccessKpi{
			Api:   "Image",
			Value: total,
		}
	}
	return data, nil
}

func (c *PostgreDb) QuerySuccessRequestGetImageCount() (config.SuccessKpi, error) {
	fmt.Print("query request data")
	var data config.SuccessKpi
	//rows, err := c.db.Query("SELECT COUNT(*)\nFROM transactions WHERE response_body NOT LIKE '\"code\"%' AND url LIKE '/s3%'") //"code":0
	rows, err := c.db.Query("SELECT * \nFROM mv_transactions_images_request_success")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.SuccessKpi{
			Api:   "Image",
			Value: total,
		}
	}
	return data, nil
}

func GetRequestGetImageTotal(c *gin.Context) (config.SuccessKpi, error) {
	var cfg config.DBConfig
	var total, count, result config.SuccessKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	total, _ = db2.QueryRequestGetImageCount()
	count, _ = db2.QuerySuccessRequestGetImageCount()
	result = config.SuccessKpi{
		Api:   total.Api,
		Value: (count.Value / total.Value) * 100,
	}
	return result, nil
}

func ExportRequestGetImageTotal() (config.SuccessKpi, error) {
	var cfg config.DBConfig
	var total, count, result config.SuccessKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	total, _ = db2.QueryRequestGetImageCount()
	count, _ = db2.QuerySuccessRequestGetImageCount()
	result = config.SuccessKpi{
		Api:   total.Api,
		Value: (count.Value / total.Value) * 100,
	}
	return result, nil
}

func (c *PostgreDb) QueryRequestTrackingCount() (config.SuccessKpi, error) {
	fmt.Print("query request data")
	var data config.SuccessKpi
	//rows, err := c.db.Query("SELECT COUNT(*)\nFROM transactions\n WHERE url LIKE '/attributes/vtracking/logged/VEHICLE/%'")
	rows, err := c.db.Query("SELECT * \nFROM mv_transactions_tracking_request")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.SuccessKpi{
			Api:   "Tracking",
			Value: total,
		}
	}
	return data, nil
}

func (c *PostgreDb) QuerySuccessRequestTrackingCount() (config.SuccessKpi, error) {
	fmt.Print("query request data")
	var data config.SuccessKpi
	//rows, err := c.db.Query("SELECT COUNT(*)\nFROM transactions WHERE response_body NOT LIKE '\"code\"%' AND url LIKE '/attributes/vtracking/logged/VEHICLE/%'")
	rows, err := c.db.Query("SELECT * \nFROM mv_transactions_tracking_request_success")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.SuccessKpi{
			Api:   "Tracking",
			Value: total,
		}
	}
	return data, nil
}

func GetRequestTrackingTotal(c *gin.Context) (config.SuccessKpi, error) {
	var cfg config.DBConfig
	var total, count, result config.SuccessKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	total, _ = db2.QueryRequestTrackingCount()
	count, _ = db2.QuerySuccessRequestTrackingCount()
	result = config.SuccessKpi{
		Api:   total.Api,
		Value: (count.Value / total.Value) * 100,
	}
	return result, nil
}

func ExportRequestTrackingTotal() (config.SuccessKpi, error) {
	var cfg config.DBConfig
	var total, count, result config.SuccessKpi
	cfg = LoadDBConfig()
	db2, _ := ConnectTransDB(cfg)
	total, _ = db2.QueryRequestTrackingCount()
	count, _ = db2.QuerySuccessRequestTrackingCount()
	result = config.SuccessKpi{
		Api:   total.Api,
		Value: (count.Value / total.Value) * 100,
	}
	return result, nil
}

func (c *PostgreDb) QueryLicenseExpiredTodayCount() (config.DataLicense, error) {
	var data config.DataLicense
	rows, err := c.db.Query("SELECT * FROM mv_license_plate_expired_today")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.DataLicense{
			Api:   "License",
			Value: total,
		}
	}
	return data, nil
}

func GetLicenseExpiredTodayTotal() (config.DataLicense, error) {
	var cfg config.DBConfig
	var data config.DataLicense
	cfg = LoadDBConfig()
	db2, _ := ConnectDevicesDB(cfg)
	data, _ = db2.QueryLicenseExpiredTodayCount()
	return data, nil
}

func (c *PostgreDb) QueryLicenseExpiredThisMonthCount() (config.DataLicense, error) {
	var data config.DataLicense
	rows, err := c.db.Query("SELECT * FROM mv_license_plate_expired_this_month")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.DataLicense{
			Api:   "License",
			Value: total,
		}
	}
	return data, nil
}

func GetLicenseExpiredThisMonthTotal() (config.DataLicense, error) {
	var cfg config.DBConfig
	var data config.DataLicense
	cfg = LoadDBConfig()
	db2, _ := ConnectDevicesDB(cfg)
	data, _ = db2.QueryLicenseExpiredThisMonthCount()
	return data, nil
}

func (c *PostgreDb) QueryLicenseExpiredThisYearCount() (config.DataLicense, error) {
	var data config.DataLicense
	rows, err := c.db.Query("SELECT * FROM mv_license_plate_expired_this_year")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.DataLicense{
			Api:   "License",
			Value: total,
		}
	}
	return data, nil
}

func GetLicenseExpiredThisYearTotal() (config.DataLicense, error) {
	var cfg config.DBConfig
	var data config.DataLicense
	cfg = LoadDBConfig()
	db2, _ := ConnectDevicesDB(cfg)
	data, _ = db2.QueryLicenseExpiredThisYearCount()
	return data, nil
}
func (c *PostgreDb) QueryLicenseNewTodayCount() (config.DataLicense, error) {
	var data config.DataLicense
	rows, err := c.db.Query("SELECT * FROM mv_license_plate_new_today")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.DataLicense{
			Api:   "License",
			Value: total,
		}
	}
	return data, nil
}

func GetLicenseNewTodayTotal() (config.DataLicense, error) {
	var cfg config.DBConfig
	var data config.DataLicense
	cfg = LoadDBConfig()
	db2, _ := ConnectDevicesDB(cfg)
	data, _ = db2.QueryLicenseNewTodayCount()
	return data, nil
}

func (c *PostgreDb) QueryLicenseNewThisMonthCount() (config.DataLicense, error) {
	var data config.DataLicense
	rows, err := c.db.Query("SELECT * FROM mv_license_plate_new_this_month")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.DataLicense{
			Api:   "License",
			Value: total,
		}
	}
	return data, nil
}

func GetLicenseNewThisMonthTotal() (config.DataLicense, error) {
	var cfg config.DBConfig
	var data config.DataLicense
	cfg = LoadDBConfig()
	db2, _ := ConnectDevicesDB(cfg)
	data, _ = db2.QueryLicenseNewThisMonthCount()
	return data, nil
}

func (c *PostgreDb) QueryLicenseNewTotalCount() (config.DataLicense, error) {
	var data config.DataLicense
	rows, err := c.db.Query("SELECT * FROM mv_license_plate_new_total")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.DataLicense{
			Api:   "License",
			Value: total,
		}
	}
	return data, nil
}

func GetLicenseNewTotal() (config.DataLicense, error) {
	var cfg config.DBConfig
	var data config.DataLicense
	cfg = LoadDBConfig()
	db2, _ := ConnectDevicesDB(cfg)
	data, _ = db2.QueryLicenseNewTotalCount()
	return data, nil
}
func (c *PostgreDb) QueryLicenseValidThisMonthCount() (config.DataLicense, error) {
	var data config.DataLicense
	rows, err := c.db.Query("SELECT * FROM mv_license_plate_valid_this_month")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.DataLicense{
			Api:   "License",
			Value: total,
		}
	}
	return data, nil
}

func GetLicenseValidThisMonthTotal() (config.DataLicense, error) {
	var cfg config.DBConfig
	var data config.DataLicense
	cfg = LoadDBConfig()
	db2, _ := ConnectDevicesDB(cfg)
	data, _ = db2.QueryLicenseValidThisMonthCount()
	return data, nil
}

func (c *PostgreDb) QueryLicenseReNewTodayCount() (config.DataLicense, error) {
	var data config.DataLicense
	rows, err := c.db.Query("SELECT * FROM mv_license_plate_renew_today_1")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.DataLicense{
			Api:   "License",
			Value: total,
		}
	}
	return data, nil
}

func GetLicenseReNewToday() (config.DataLicense, error) {
	var cfg config.DBConfig
	var data config.DataLicense
	cfg = LoadDBConfig()
	db2, _ := ConnectAttributesDB(cfg)
	data, _ = db2.QueryLicenseReNewTodayCount()
	return data, nil
}

func (c *PostgreDb) QueryLicenseReNewThisMonthCount() (config.DataLicense, error) {
	var data config.DataLicense
	rows, err := c.db.Query("SELECT * FROM mv_license_plate_renew_this_month_1")
	if err != nil {
		fmt.Printf("Loi get db", err)
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var total float64
		if err := rows.Scan(&total); err != nil {
			fmt.Printf("Loi scan db ", err)
			return data, err
		}
		data = config.DataLicense{
			Api:   "License",
			Value: total,
		}
	}
	return data, nil
}

func GetLicenseReNewThisMonth() (config.DataLicense, error) {
	var cfg config.DBConfig
	var data config.DataLicense
	cfg = LoadDBConfig()
	db2, _ := ConnectAttributesDB(cfg)
	data, _ = db2.QueryLicenseReNewThisMonthCount()
	return data, nil
}

func (c *PostgreDb) BatchInsertApiLogs(logs []ApiLog) error {
	if len(logs) == 0 {
		return nil
	}

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO api_logs (url, status_code, latency, timestamp)
		VALUES ($1,$2,$3,$4)
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, l := range logs {
		_, err := stmt.Exec(
			l.Url,
			l.StatusCode,
			l.Latency,
			l.Timestamp,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (c *PostgreDb) SetDowntimeDB(pod string, downtime, total uint64) {

	created_time := time.Now()
	query := `INSERT INTO your_table_name (pod, downtime, total, created_time) VALUES ($1, $2, $3, $4)`
	_, err := c.db.Exec(query, pod, downtime, total, created_time)
	if err != nil {
		log.Printf("Failed to insert data into PostgreSQL: %v", err)
	} else {
		log.Println("Data inserted successfully into PostgreSQL")
	}
}

func Test(db *sqlx.DB) {
	// this Pings the database trying to connect
	// use sqlx.Open() for sql.Open() semantics
	//db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	//db := ConnectToDb()

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	//db.MustExec(Schema)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "John", "Doe", "johndoeDNE@gmail.net")
	tx.MustExec("INSERT INTO place (country, city, telcode) VALUES ($1, $2, $3)", "United States", "New York", "1")
	tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Hong Kong", "852")
	tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Singapore", "65")
	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	tx.NamedExec("INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &Person{"Jane", "Citizen", "jane.citzen@example.com"})
	tx.Commit()

	// Query the database, storing results in a []Person (wrapped in []interface{})
	people := []Person{}
	db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")
	jason, john := people[0], people[1]

	fmt.Printf("%#v\n%#v", jason, john)
	// Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}
	// Person{FirstName:"John", LastName:"Doe", Email:"johndoeDNE@gmail.net"}

	// You can also get a single result, a la QueryRow
	jason = Person{}
	err := db.Get(&jason, "SELECT * FROM person WHERE first_name=$1", "Jason")
	fmt.Printf("%#v\n", jason)
	// Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}

	// if you have null fields and use SELECT *, you must use sql.Null* in your struct
	places := []Place{}
	err = db.Select(&places, "SELECT * FROM place ORDER BY telcode ASC")
	if err != nil {
		fmt.Println(err)
		return
	}
	usa, singsing, honkers := places[0], places[1], places[2]

	fmt.Printf("%#v\n%#v\n%#v\n", usa, singsing, honkers)
	// Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
	// Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}
	// Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}

	// Loop through rows using only one struct
	place := Place{}
	rows, err := db.Queryx("SELECT * FROM place")
	for rows.Next() {
		err := rows.StructScan(&place)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%#v\n", place)
	}
	// Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
	// Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}
	// Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}

	// Named queries, using `:name` as the bindvar.  Automatic bindvar support
	// which takes into account the dbtype based on the driverName on sqlx.Open/Connect
	_, err = db.NamedExec(`INSERT INTO person (first_name,last_name,email) VALUES (:first,:last,:email)`,
		map[string]interface{}{
			"first": "Bin",
			"last":  "Smuth",
			"email": "bensmith@allblacks.nz",
		})

	// Selects Mr. Smith from the database
	rows, err = db.NamedQuery(`SELECT * FROM person WHERE first_name=:fn`, map[string]interface{}{"fn": "Bin"})

	// Named queries can also use structs.  Their bind names follow the same rules
	// as the name -> db mapping, so struct fields are lowercased and the `db` tag
	// is taken into consideration.
	rows, err = db.NamedQuery(`SELECT * FROM person WHERE first_name=:first_name`, jason)

	// batch insert

	// batch insert with structs
	personStructs := []Person{
		{FirstName: "Ardie", LastName: "Savea", Email: "asavea@ab.co.nz"},
		{FirstName: "Sonny Bill", LastName: "Williams", Email: "sbw@ab.co.nz"},
		{FirstName: "Ngani", LastName: "Laumape", Email: "nlaumape@ab.co.nz"},
	}

	_, err = db.NamedExec(`INSERT INTO person (first_name, last_name, email)
        VALUES (:first_name, :last_name, :email)`, personStructs)

	// batch insert with maps
	personMaps := []map[string]interface{}{
		{"first_name": "Ardie", "last_name": "Savea", "email": "asavea@ab.co.nz"},
		{"first_name": "Sonny Bill", "last_name": "Williams", "email": "sbw@ab.co.nz"},
		{"first_name": "Ngani", "last_name": "Laumape", "email": "nlaumape@ab.co.nz"},
	}

	_, err = db.NamedExec(`INSERT INTO person (first_name, last_name, email)
        VALUES (:first_name, :last_name, :email)`, personMaps)
}

func LoadCpu(db *sqlx.DB) {
	newData := CpuUsage{
		CpuID:     1,
		Value:     25.5,
		Timestamp: "",
	}
	tx := db.MustBegin()
	tx.NamedExec("INSERT INTO cpuUsage (cpu_id, value,timestamp) VALUES (:cpu_id, :value, current_timestamp)", newData)
	tx.Commit()

	data := []CpuUsage{}
	db.Select(&data, "SELECT * FROM cpuUsage")
	for _, result := range data {
		fmt.Printf("%#v\n", result)
	}
}

func KpiModule(db *sqlx.DB, data []config.KpiData) {
	if data == nil {
		var reponse []KpiRequestModule
		jsonData := `[
        {"pod": "user", "value": 0.662922382837676},
        {"pod": "device", "value": 5.682267922879199},
        {"pod": "organization", "value": 0.2607443523626743},
        {"pod": "attribute", "value": 9.793216861596788}
    ]`
		err := json.Unmarshal([]byte(jsonData), &reponse)
		if err != nil {
			fmt.Println("Lỗi giải mã JSON:", err)
			return
		}
		tx := db.MustBegin()
		for _, result := range reponse {
			tx.NamedExec("INSERT INTO kpimodule (pod, value,timestamp) VALUES (:pod, :value, current_timestamp)", result)
			// Xóa dữ liệu cũ nhất nếu có hơn 1000 dòng
			//_, err = tx.Exec("DELETE FROM kpimodule WHERE id = (SELECT id FROM kpimodule ORDER BY timestamp ASC LIMIT 1 OFFSET 999);")
			if err != nil {
				fmt.Println("Lôi xoá line", err)
			}
		}

		tx.Commit()
	}
	var reponse []KpiRequestModule

	for _, result := range data {
		req, _ := strconv.ParseFloat(result.Req, 64)
		num, _ := strconv.ParseFloat(result.Error, 64)
		cal := KpiRequestModule{
			Pod:       result.Pod,
			Method:    result.Method,
			Value:     (num / req) * 100,
			Timestamp: "",
		}
		reponse = append(reponse, cal)
	}
	tx := db.MustBegin()
	for _, result := range reponse {
		tx.NamedExec("INSERT INTO kpimodule (pod, value,timestamp) VALUES (:pod, :value, current_timestamp)", result)
		// Xóa dữ liệu cũ nhất nếu có hơn 1000 dòng
		_, err := tx.Exec("DELETE FROM kpimodule WHERE id = (SELECT id FROM kpimodule ORDER BY timestamp ASC LIMIT 1 OFFSET 999);")
		if err != nil {
			fmt.Println("Lôi xoá line", err)
		}
	}

	tx.Commit()

}

func (c *PostgreDb) SetCpuUsage(server string, value uint64) {

	created_time := time.Now()
	query := `WITH inserted AS (
		INSERT INTO cpuusage (server, value, timestamp)
		VALUES ($1, $2, $3)
		RETURNING ctid
	), delete_oldest AS (
		DELETE FROM cpuusage
		WHERE ctid IN (
			SELECT ctid
			FROM cpuusage
			ORDER BY timestamp
			LIMIT 1
		)
		AND (SELECT COUNT(*) FROM cpuusage) > 3000
		RETURNING ctid
	)
	SELECT 1`
	_, err := c.db.Exec(query, server, value, created_time)
	if err != nil {
		log.Printf("Failed to insert data into PostgreSQL: %v", err)
	} else {
		log.Println("Data inserted successfully into PostgreSQL")
	}
}

func SetCpu(server string, value uint64) {
	var cfg config.DBConfig
	cfg = LoadDBConfig()
	db, _ := ConnectNewDB(cfg)
	db.SetCpuUsage(server, value)
	defer db.db.Close()
}

func (c *PostgreDb) SetRamUsage(server string, total uint64, avail uint64) {

	created_time := time.Now()
	query := `	WITH inserted AS (
		INSERT INTO ramusage (server, total, avail, timestamp)
		VALUES ($1, $2, $3, $4)
		RETURNING ctid
	), delete_oldest AS (
		DELETE FROM ramusage
		WHERE ctid IN (
			SELECT ctid
			FROM ramusage
			ORDER BY timestamp
			LIMIT 1
		)
		AND (SELECT COUNT(*) FROM ramusage) > 3000
		RETURNING ctid
	)
	SELECT 1`
	_, err := c.db.Exec(query, server, total, avail, created_time)
	if err != nil {
		log.Printf("Failed to insert data into PostgreSQL: %v", err)
	} else {
		log.Println("Data inserted successfully into PostgreSQL")
	}
}

func SetRam(server string, total, avail uint64) {
	var cfg config.DBConfig
	cfg = LoadDBConfig()
	db, _ := ConnectNewDB(cfg)
	db.SetRamUsage(server, total, avail)
	defer db.db.Close()
}

func (c *PostgreDb) SetDownTimeKpi(data config.KpiData) {

}

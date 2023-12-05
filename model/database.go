package model

import (
	"Monitor_Platform/config"
	"strconv"

	//"Monitor_Platform/services"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

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

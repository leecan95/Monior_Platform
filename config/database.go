package config

import "time"
import "database/sql"

const (
	EnvDBHost        = ""
	EnvDBPort        = ""
	EnvDBUser        = "VT_USERS_DB_USER"
	EnvDBPass        = "VT_USERS_DB_PASS"
	EnvDB            = ""
	EnvDBSSLMode     = "VT_USERS_DB_SSL_MODE"
	EnvDBSSLCert     = "VT_USERS_DB_SSL_CERT"
	EnvDBSSLKey      = "VT_USERS_DB_SSL_KEY"
	EnvDBSSLRootCert = "VT_USERS_DB_SSL_ROOT_CERT"
	DefDBHost        = "10.207.189.122"
	DefDBPort        = "5000"
	DefDBUser        = "postgres"
	DefDBPass        = "newpassword"
	DefDB            = "monitor"
	DefDBSSLMode     = "disable"
	DefDBSSLCert     = ""
	DefDBSSLKey      = ""
	DefDBSSLRootCert = ""
)

type DBConfig struct {
	Host              string
	Port              string
	User              string
	Pass              string
	Name              string
	SSLMode           string
	SSLCert           string
	SSLKey            string
	SSLRootcert       string
	ConnectionTimeout int
	LogLevel          string
	MinConn           int32
	MaxConn           int32
	MaxConnLifeTime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
	MaxRetry          uint64
	RetryInterval     uint64
}

type TransactionDb struct {
	Method       string     `json:"method"`
	URL          string     `json:"url"`
	StatusCode   int        `json:"status_code"`
	ResponseBody *string `json:"response_body,omitempty"`
	ServiceName  string     `json:"service_name"`
	Latency      int64      `json:"latency"`
	TS           int64      `json:"ts"`
	UserID       string  `json:"user_id"`
	ProjectID    string  `json:"project_id"`
}

type TransactionQuery struct {
    Method       string
    URL          string
    StatusCode   sql.NullInt32
    ResponseBody sql.NullString
    ServiceName  sql.NullString
    Latency      int64
    TS           int64
    UserID       string
    ProjectID    string
}

var ActionURLMap = map[string]string{
	"Login":               "%/vtracking/login%",
	"ReceiveOTP":          "%/vtracking/users/verifycode%",
	// "ReceiveAlert":        "%alert%",
	"GetListVehicle":      "%/devices/vtracking/vehicle/filter%",
	"ConfigDevice":        "%config%",
	"GetLocation":         "%/devices/vtracking/vehicle%",
	"GetDetailVehicle":    "%/devices/vtracking/vehicle%",
	"GetListLocation":     "%/attributes/vtracking/logged/VEHICLE%",
	"GetImage":            "%/vtracking/s3%",
	// "GetVideoFromMemory":  "/vtracking/s3",
	// "GetStreamVideo":      "/vtracking/s3",
	"GetReport":           "%/report/attributes/vtracking%",
}

type LatencyPercentileResult struct {
	CountBelow5    int64   `json:"count_latency_below_5"`
	Percentile95   float64 `json:"latency_95th_percentile"`
	TotalRecords   int64   `json:"total_records"`
	Exceeds95Pct   string  `json:"exceeds_95_percent"`
}

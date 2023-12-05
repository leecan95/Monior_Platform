package config

import "time"

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
	DefDBHost        = "haproxy.viot-db"
	DefDBPort        = "5000"
	DefDBUser        = "viot"
	DefDBPass        = "viot"
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

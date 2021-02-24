package db

import (
	"database/sql"
	"sync"

	"github.com/micro/cli"
)

var (
	pgConfig *PGConfig
	pgDB     *sql.DB
	m        sync.RWMutex
	once     sync.Once
)

type PGConfig struct {
	DBName            string `json:"dbName"`            // The name of the database to connect to
	User              string `json:"user"`              // the user to sign in as
	Password          string `json:"password"`          // The user's password
	Host              string `json:"host"`              // The host to connect to.Values that start with / are for unix domain sockets. (default is localhost)
	Port              int    `json:"port"`              // The port to bind to.(default is 5432)
	SSLMode           string `json:"sslMode"`           // Whether or not to use SSL (default is require, this is not the default for libpq)
	ConnectTimeout    int    `json:"connectTimeout"`    // Maximum wait for connection, in seconds.Zero or not specified means wait indefinitely.
	SSLCert           string `json:"sslCert"`           // Cert file location. The file must contain PEM encoded data.
	SSLKey            string `json:"sslKey"`            // Key file location.The file must contain PEM encoded data.
	SSLRootCert       string `json:"sslRootCert"`       // The location of the root certificate file.The file must contain PEM encoded data.)
	MaxOpenConnection int    `json:"maxOpenConnection"` // use the default 0
	MaxIdleConnection int    `json:"maxIdleConnection"` // use the default 0
}

func Init(ctx *cli.Context) {
	m.Lock()
	defer m.Unlock()

	pgConfig = newPGConfig()

	if pgConfigFile := ctx.String("pg_config_file"); len(pgConfigFile) > 0 {
		if err := config.Get(pgConfigFile).Scan(&pgConfig); err != nil {
			panic(err)
		}
	}

	initDB()
}

func newPGConfig() *PGConfig {
	return &PGConfig{
		DBName:         "postgres",
		User:           "postgres",
		Password:       "password",
		Host:           "localhost",
		ConnectTimeout: 10,
		Port:           5432,
		SSLMode:        "disable",
	}
}

// initDB 初始化数据库
func initDB() {
	once.Do(func() {
		initPG()
	})
}

// GetPG get postgre db
func GetPG() *sql.DB {
	return pgDB
}

package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Port     int    `koanf:"port"`
	Host     string `koanf:"host"`
	DbName   string `koanf:"db_name"`
}

type MySQLDB struct {
	config Config
	db     *sql.DB
}

func (m *MySQLDB) Conn() *sql.DB {
	return m.db
}

func NewMYSQL(config Config) *MySQLDB {

	//* pardeTime=true changes the output type of DATE and DATETIME \
	//* values to time.time instead of []byte / string
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.Username,
		config.Password, config.Host, config.Port, config.DbName))
	if err != nil {
		panic(fmt.Errorf("can not open mysql database: %v", err))
	}

	//* options
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	//TODO migrate to mysql manually or ysung third party

	return &MySQLDB{config: config, db: db}
}

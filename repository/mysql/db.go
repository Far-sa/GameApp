package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLDB struct {
	Db *sql.DB
}

func NewMYSQL() *MySQLDB {

	dsn := "user:password@/dbname"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Errorf("can not open mysql database: %v", err))
	}

	//* options
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MySQLDB{Db: db}
}

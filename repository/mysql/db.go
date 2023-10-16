package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type MySQLDB struct {
	Db *sql.DB
}

func NewMYSQL() *MySQLDB {

	vErr := godotenv.Load(".env")
	if vErr != nil {
		log.Fatal("Error loading environment")
	}

	dsn := "gameapp:gameapp@tcp(localhost:3306)/gamedb"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Errorf("can not open mysql database: %v", err))
	}

	//* options
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	//TODO migrate to mysql manually or ysung third party

	return &MySQLDB{Db: db}
}

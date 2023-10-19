package migrator

import (
	"database/sql"
	"fmt"
	"game-app/repository/mysql"

	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	config      mysql.Config
	migratation *migrate.FileMigrationSource
}

func New(config mysql.Config) Migrator {

	migrations := &migrate.FileMigrationSource{
		Dir: "./repository/mysql/migrations",
	}

	return Migrator{config: config, migratation: migrations}
}

func (m Migrator) Up() {
	db, err := sql.Open("", "")
	if err != nil {
		// handle error
	}

	n, err := migrate.Exec(db, "", m.migratation, migrate.Up)
	if err != nil {
		// handle error
	}
	fmt.Println(n)
}

func (m Migrator) Down() {
	db, err := sql.Open("", "")
	if err != nil {
		// handle error
	}

	n, err := migrate.Exec(db, "", m.migratation, migrate.Down)
	if err != nil {
		// handle error
	}
	fmt.Println(n)
}

func (m Migrator) Status() {
	//TODO
}

package database

import (
	"log"
	"sync"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func MysQL() (*sql.DB, func()) {
	var db *sql.DB
	var err error
	var once sync.Once
	var cleanup func()

	once.Do(func() {
		db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/test")
		if err != nil {
			log.Fatalf("error in openning mysql connection: err %v", err)
		}
		err = db.Ping()
		if err != nil {
			log.Fatalf("error in mysql ping: err %v", err)
		}

		cleanup = func() {
			db.Close()
		}
	})

	return db, cleanup
}

package data

import (
	"database/sql"
	"log"
	"time"

	"github.com/spf13/viper"
)

func DbOpen() (db *sql.DB) {
	db, err := sql.Open("sqlite3", viper.GetString("DbLocation"))
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS record( id INTEGER PRIMARY KEY, timestamp TIMESTAMP, location TEXT, server TEXT, ping FLOAT, dl FLOAT, ul FLOAT)")
	if err != nil {
		log.Printf("Err: %#v\n", err)
	}
	stmt.Exec()
	return db
}

func DbInsert(db *sql.DB, timestamp time.Time, location string, server string, ping int64, dl float64, ul float64) {
	stmt, err := db.Prepare("INSERT INTO record(timestamp, location, server, ping, dl, ul) VALUES(?,?,?,?,?,?)")
	if err != nil {
		log.Printf("Err: %#v\n", err)
	}
	stmt.Exec(timestamp, location, server, ping, dl, ul)
}

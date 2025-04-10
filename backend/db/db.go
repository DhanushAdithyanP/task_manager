package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal("Error opening DB: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}

	log.Println("Connected to DB!")
}

package database_dir

import (
	"database/sql"
	"log"
)

func DBconnect() *sql.DB {
	connStr := "user=postgres dbname=postgres password= postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

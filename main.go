package main

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"todo/database_dir"
	"todo/router_dir"
)

func main() {

	db := database_dir.DBconnect()
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://database_dir/migration_dir", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	er := m.Up()
	if er == migrate.ErrNoChange {
		//
	}
	r := router_dir.Router()
	err1 := http.ListenAndServe(":8080", r)
	if err1 != nil {
		log.Fatalf("Error")
		return
	}

}

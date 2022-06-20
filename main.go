package main

import (
	"github.com/go-chi/chi"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"todo/database_dir"
	"todo/handler_dir"
	"todo/middleware_dir"
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
	r := chi.NewRouter()

	r.Post("/signup", handler_dir.SignUp)
	r.Post("/login", handler_dir.Login)
	r.HandleFunc("/home", middleware_dir.AuthMiddleware(handler_dir.Home))
	err1 := http.ListenAndServe(":8080", r)
	if err1 != nil {
		log.Fatalf("Error")
		return
	}

}

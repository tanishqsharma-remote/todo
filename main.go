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

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", handler_dir.SignUp)
		r.Post("/login", handler_dir.Login)
		r.Get("/logout", handler_dir.Logout)
	})

	r.Route("/todo", func(r chi.Router) {
		r.Post("/", handler_dir.CreateTask)
		r.Put("/", handler_dir.DoneTask)
		r.Delete("/", handler_dir.ArchiveTask)
		r.Group(func(r chi.Router) {
			r.Use(middleware_dir.AuthMiddleware)
			r.Get("/", handler_dir.GetTask)
		})
	})

	r.Route("/home", func(r chi.Router) {
		r.Use(middleware_dir.AuthMiddleware)
		r.Get("/", handler_dir.Home)
	})

	r.Route("/refresh", func(r chi.Router) {
		r.Use(middleware_dir.RefreshMiddleware)
		r.Get("/", handler_dir.Refresh)
	})

	err1 := http.ListenAndServe(":8080", r)
	if err1 != nil {
		log.Fatalf("Error")
		return
	}

}

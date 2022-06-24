package router_dir

import (
	"github.com/go-chi/chi"
	"todo/handler_dir"
	"todo/middleware_dir"
)

func Router() *chi.Mux {
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

	return r
}

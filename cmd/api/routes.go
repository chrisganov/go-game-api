package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() *chi.Mux {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Routes
	router.Get("/healthcheck", app.healthcheckHandler)

	router.Route("/users", func(r chi.Router) {
		r.Get("/", app.getUsers)
		r.Post("/", app.createUser)

		r.Get("/{id}", app.getUserById)
	})

	router.Route("/scores", func(r chi.Router) {
		r.Get("/", app.getAllScores)
		r.Post("/", app.createScore)
	})

	// Error
	router.NotFound(app.notFoundResponse)
	router.MethodNotAllowed(app.methodNotAllowedResponse)

	return router
}

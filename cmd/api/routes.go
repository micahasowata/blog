package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()
	app.stack(router)
	router.Post("/v1/users/register", app.registerUser)
	return router
}

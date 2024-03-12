package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hibiken/asynq"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()
	app.stack(router)
	router.Post("/v1/users/register", app.registerUser)
	router.Post("/v1/users/verify", app.verifyEmail)
	return router
}

func (app *application) jobs() asynq.Handler {
	mux := asynq.NewServeMux()
	mux.HandleFunc(typeWelcomeEmail, app.handleWelcomeEmailDelivery)
	return mux
}

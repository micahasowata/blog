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
	router.Post("/v1/tokens/login", app.createLoginToken)
	router.Post("/v1/users/login", app.loginUser)
	router.With(app.requireAccessToken).Post("/v1/users/logout", app.logoutUser)
	router.With(app.requireRefreshToken).Post("/v1/tokens/refresh", app.refreshToken)
	router.With(app.requireAccessToken).Get("/v1/users/me", app.getUserProfile)
	router.With(app.requireAccessToken).Patch("/v1/users/update", app.updateUserProfile)
	return router
}

func (app *application) jobs() asynq.Handler {
	mux := asynq.NewServeMux()
	mux.HandleFunc(typeOTPEmail, app.handleOTPEmailDelivery)
	mux.HandleFunc(typeLoginEmail, app.handleLoginEmailTask)
	return mux
}

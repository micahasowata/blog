package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func (app *application) setMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  zap.NewStdLog(app.logger),
		NoColor: true,
	}))
	router.Use(middleware.Recoverer)
	router.Use(middleware.CleanPath)
	router.Use(middleware.RequestID)
	router.Use(middleware.RequestSize(1_048_576))
}

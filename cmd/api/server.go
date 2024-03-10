package main

import (
	"net/http"
	"time"

	"github.com/pseidemann/finish"
	"go.uber.org/zap"
)

func (app *application) serve() {
	srv := &http.Server{
		Addr:         app.config.Address,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorLog:     zap.NewStdLog(app.logger),
	}

	manager := finish.New()
	manager.Log = app.logger.Sugar()

	manager.Add(srv)

	app.logger.Info("starting server", zap.String("address", srv.Addr))
	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			app.logger.Error(err.Error())
		}
	}()

	manager.Wait()
}

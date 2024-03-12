package main

import (
	"net/http"
	"time"

	"github.com/hibiken/asynq"
	"github.com/pseidemann/finish"
	"go.uber.org/zap"
)

func (app *application) serve() {
	server := &http.Server{
		Addr:         app.config.Address,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorLog:     zap.NewStdLog(app.logger),
	}

	rdc := asynq.RedisClientOpt{
		Addr: app.config.AsynqRDC,
	}
	cfg := asynq.Config{
		Logger:      app.logger.Sugar(),
		Concurrency: 10,
	}

	processor := asynq.NewServer(rdc, cfg)
	err := processor.Start(app.jobs())
	if err != nil {
		app.logger.Fatal("asynq server error", zap.Error(err))
	}

	manager := finish.New()
	manager.Log = app.logger.Sugar()

	manager.Add(server)

	app.logger.Info("starting server", zap.String("address", server.Addr))
	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			processor.Shutdown()
			app.logger.Error(err.Error())
		}
	}()

	manager.Wait()
}

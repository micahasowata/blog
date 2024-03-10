package main

import (
	"log"

	"github.com/micahasowata/blog/internal/config"
	"go.uber.org/zap"
)

type application struct {
	logger *zap.Logger
	config *config.Config
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err.Error())
	}

	config, err := config.New()
	if err != nil {
		log.Fatal(err.Error())
	}

	app := &application{
		logger: logger,
		config: config,
	}

	app.serve()
}

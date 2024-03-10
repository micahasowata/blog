package main

import (
	"log"

	"github.com/micahasowata/blog/internal/config"
	"github.com/micahasowata/jason"
	"go.uber.org/zap"
)

type application struct {
	*jason.Jason
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
		Jason:  jason.New(int64(config.MaxSize), false, true),
		logger: logger,
		config: config,
	}

	app.serve()
}

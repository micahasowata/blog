package main

import (
	"log"

	"github.com/micahasowata/blog/internal/config"
	"github.com/micahasowata/blog/internal/db"
	"github.com/micahasowata/blog/internal/models"
	"github.com/micahasowata/jason"
	"go.uber.org/zap"
)

type application struct {
	*jason.Jason

	logger *zap.Logger
	config *config.Config
	models *models.Models
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

	db, err := db.NewProduction(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	app := &application{
		Jason:  jason.New(int64(config.MaxSize), false, true),
		logger: logger,
		config: config,
		models: models.New(db),
	}

	app.serve()
}

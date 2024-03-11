package main

import (
	"errors"
	"log"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/micahasowata/blog/internal/config"
	"github.com/micahasowata/blog/internal/db"
	"github.com/micahasowata/blog/internal/models"
	"github.com/micahasowata/jason"
	"go.uber.org/zap"
)

type application struct {
	*jason.Jason

	logger     *zap.Logger
	config     *config.Config
	translator ut.Translator
	validate   *validator.Validate
	models     *models.Models
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

	localeEN := en.New()
	universal := ut.New(localeEN, localeEN)
	translator, ok := universal.GetTranslator("en")
	if !ok {
		log.Fatal(errors.New("unable to get validation translator"))
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	en_translations.RegisterDefaultTranslations(validate, translator)

	app := &application{
		Jason:      jason.New(int64(config.MaxSize), false, true),
		logger:     logger,
		config:     config,
		translator: translator,
		validate:   validate,
		models:     models.New(db),
	}

	app.serve()
}

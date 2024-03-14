package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/kataras/jwt"
	"github.com/micahasowata/blog/internal/config"
	"github.com/micahasowata/blog/internal/db"
	"github.com/micahasowata/blog/internal/models"
	"github.com/micahasowata/jason"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type application struct {
	*jason.Jason

	logger     *zap.Logger
	config     *config.Config
	translator ut.Translator
	validate   *validator.Validate
	models     *models.Models
	rclient    *redis.Client
	executor   *asynq.Client
	blocklist  *jwt.Blocklist
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = godotenv.Load()
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

	executor := asynq.NewClient(asynq.RedisClientOpt{
		Addr: config.RDB,
	})

	rclient := redis.NewClient(&redis.Options{
		Addr: config.RDB,
	})

	blocklist := jwt.NewBlocklistContext(context.Background(), 1*time.Hour)

	app := &application{
		Jason:      jason.New(int64(config.MaxSize), false, true),
		logger:     logger,
		config:     config,
		translator: translator,
		validate:   validate,
		models:     models.New(db),
		rclient:    rclient,
		executor:   executor,
		blocklist:  blocklist,
	}

	app.serve()
}

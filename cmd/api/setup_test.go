package main

import (
	"context"
	"testing"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kataras/jwt"
	"github.com/micahasowata/blog/internal/config"
	"github.com/micahasowata/blog/internal/db"
	"github.com/micahasowata/blog/internal/models"
	"github.com/micahasowata/jason"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	cfg, err := config.New()
	require.Nil(t, err)

	tdb, err := db.NewTest(cfg)
	require.Nil(t, err)

	return tdb
}

func setupApp(t *testing.T, db *pgxpool.Pool) *application {
	t.Helper()

	cfg, err := config.New()
	require.Nil(t, err)

	localeEN := en.New()
	universal := ut.New(localeEN, localeEN)

	translator, ok := universal.GetTranslator("en")
	require.NotEmpty(t, ok)

	validate := validator.New(validator.WithRequiredStructEnabled())
	en_translations.RegisterDefaultTranslations(validate, translator)

	executor := asynq.NewClient(asynq.RedisClientOpt{
		Addr: cfg.RDB,
	})

	rclient := redis.NewClient(&redis.Options{
		Addr: cfg.RDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	blocklist := jwt.NewBlocklistContext(ctx, 1*time.Hour)

	app := &application{
		Jason:      jason.New(int64(cfg.MaxSize), false, true),
		logger:     zap.NewExample(),
		config:     cfg,
		validate:   validate,
		translator: translator,
		models:     models.New(db),
		executor:   executor,
		rclient:    rclient,
		blocklist:  blocklist,
	}

	return app
}

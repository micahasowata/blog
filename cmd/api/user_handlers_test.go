package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/micahasowata/blog/internal/config"
	"github.com/micahasowata/blog/internal/db"
	"github.com/micahasowata/blog/internal/models"
	"github.com/micahasowata/jason"
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

	t.Log(cfg)

	localeEN := en.New()
	universal := ut.New(localeEN, localeEN)

	translator, ok := universal.GetTranslator("en")
	require.NotEmpty(t, ok)

	validate := validator.New(validator.WithRequiredStructEnabled())
	en_translations.RegisterDefaultTranslations(validate, translator)

	executor := asynq.NewClient(asynq.RedisClientOpt{
		Addr: cfg.AsynqRDC,
	})

	app := &application{
		Jason:      jason.New(int64(cfg.MaxSize), false, true),
		logger:     zap.NewExample(),
		config:     cfg,
		validate:   validate,
		translator: translator,
		models:     models.New(db),
		executor:   executor,
	}

	return app
}

func TestRegisterUser(t *testing.T) {
	t.Skip()

	tdb := setupDB(t)
	defer db.Clean(tdb)

	app := setupApp(t, tdb)

	server := httptest.NewServer(app.routes())
	defer server.Close()

	tests := []struct {
		name string
		body string
		code int
	}{
		{
			name: "valid",
			body: `{"name": "addam","username": "iamadam","email": "theaddam@gmail.com"}`,
			code: http.StatusOK,
		},
		{
			name: "bad body",
			body: `{"password":"9LdPaiw8B"}`,
			code: http.StatusBadRequest,
		},
		{
			name: "invalid body",
			body: `{"name": "addam","username": "§age","email": "theaddam@gmail.com"}`,
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "duplicate username",
			body: `{"name": "addam","username": "iamadam","email": "theaddam@gmail.com"}`,
			code: http.StatusConflict,
		},
		{
			name: "duplicate email",
			body: `{"name": "addam","username": "iamadam","email": "theaddam45@gmail.com"}`,
			code: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httpexpect.Default(t, server.URL)

			req.POST("/v1/users/register").
				WithHeader(jason.ContentType, jason.ContentTypeJSON).
				WithBytes([]byte(tt.body)).
				Expect().
				Status(tt.code)
		})
	}
}

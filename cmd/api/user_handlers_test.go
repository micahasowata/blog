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

	tdb, err := db.NewTest()
	require.Nil(t, err)

	return tdb
}

func setupApp(t *testing.T, db *pgxpool.Pool) *application {
	t.Helper()

	cfg := &config.Config{
		MaxSize: 1_048_576,
	}

	localeEN := en.New()
	universal := ut.New(localeEN, localeEN)
	translator, ok := universal.GetTranslator("en")
	require.NotEmpty(t, ok)
	validate := validator.New(validator.WithRequiredStructEnabled())
	en_translations.RegisterDefaultTranslations(validate, translator)

	app := &application{
		Jason:      jason.New(int64(cfg.MaxSize), false, true),
		logger:     zap.NewExample(),
		config:     cfg,
		validate:   validate,
		translator: translator,
		models:     models.New(db),
	}

	return app
}

func TestRegisterUser(t *testing.T) {
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
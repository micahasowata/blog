package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/micahasowata/blog/internal/db"
	"github.com/micahasowata/blog/internal/models"
	"github.com/micahasowata/jason"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func TestVerifyEmail(t *testing.T) {
	tdb := setupDB(t)
	defer db.Clean(tdb)

	app := setupApp(t, tdb)

	user := &models.Users{
		ID:       xid.New().String(),
		Name:     "addam",
		Username: "iamaddam",
		Email:    "addam@gmail.com",
	}

	_, err := app.models.Users.Insert(user)
	require.Nil(t, err)

	server := httptest.NewServer(app.routes())
	defer server.Close()

	token := setUpToken(t, app, user)

	body := fmt.Sprintf(`{"token":"%s"}`, token)

	tests := []struct {
		name string
		body string
		code int
	}{
		{
			name: "valid",
			body: body,
			code: http.StatusOK,
		},
		{
			name: "bad body",
			body: `{"password":"9LdPaiw8B"}`,
			code: http.StatusBadRequest,
		},
		{
			name: "invalid body",
			body: `{"token":"5674902"}`,
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "missing token",
			body: `{"token":"567490"}`,
			code: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httpexpect.Default(t, server.URL)

			req.POST("/v1/users/verify").
				WithHeader(jason.ContentType, jason.ContentTypeJSON).
				WithBytes([]byte(tt.body)).
				Expect().
				Status(tt.code)
		})
	}
}

func TestVerifyEmail_UserNotFound(t *testing.T) {
	tdb := setupDB(t)
	defer db.Clean(tdb)

	app := setupApp(t, tdb)

	user := &models.Users{
		ID:       xid.New().String(),
		Name:     "addam",
		Username: "iamaddam",
		Email:    "addam@gmail.com",
	}

	server := httptest.NewServer(app.routes())
	defer server.Close()

	token := setUpToken(t, app, user)

	body := fmt.Sprintf(`{"token":"%s"}`, token)

	t.Run("missing email", func(t *testing.T) {
		req := httpexpect.Default(t, server.URL)

		req.POST("/v1/users/verify").
			WithHeader(jason.ContentType, jason.ContentTypeJSON).
			WithBytes([]byte(body)).
			Expect().
			Status(http.StatusForbidden)
	})
}

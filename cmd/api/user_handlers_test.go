package main

import (
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

func TestGetUserProfile(t *testing.T) {
	tdb := setupDB(t)
	defer db.Clean(tdb)

	app := setupApp(t, tdb)

	user := &models.Users{
		ID:       xid.New().String(),
		Name:     "addam",
		Username: "iamaddam",
		Email:    "addam@gmail.com",
	}

	createdUser, err := app.models.Users.Insert(user)
	require.Nil(t, err)

	accessToken, err := app.newAccessToken(&tokenClaims{
		ID: createdUser.ID,
	})

	require.Nil(t, err)

	server := httptest.NewServer(app.routes())
	defer server.Close()

	req := httpexpect.Default(t, server.URL)

	req.GET("/v1/users/me").
		WithHeader("Authorization", "Bearer "+accessToken).
		Expect().
		Status(http.StatusOK)
}

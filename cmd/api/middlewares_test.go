package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/micahasowata/blog/internal/db"
	"github.com/micahasowata/blog/internal/models"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequireAccessToken(t *testing.T) {
	tdb := setupDB(t)
	defer db.Clean(tdb)

	app := setupApp(t, tdb)

	user := &models.Users{
		ID:       xid.New().String(),
		Name:     "Adam",
		Username: "iamadam",
		Email:    "adam45@gmail.com",
	}

	createdUser, err := app.models.Users.Insert(user)
	require.Nil(t, err)

	accessToken, err := app.newAccessToken(&tokenClaims{ID: createdUser.ID, StdClaims: nil})
	require.Nil(t, err)

	refreshToken, err := app.newRefreshToken(&tokenClaims{ID: createdUser.ID, StdClaims: nil})
	require.Nil(t, err)

	tests := []struct {
		name   string
		header string
		code   int
	}{
		{
			name:   "valid",
			header: "Bearer " + accessToken,
			code:   http.StatusOK,
		},
		{
			name:   "no token",
			header: "Bearer ",
			code:   http.StatusForbidden,
		},
		{
			name:   "refresh token",
			header: "Bearer " + refreshToken,
			code:   http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			r, err := http.NewRequest(http.MethodPut, "/", nil)
			require.Nil(t, err)

			r.Header.Set("Authorization", tt.header)

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte("OK"))
				require.Nil(t, err)
			})

			app.requireAccessToken(next).ServeHTTP(rr, r)

			rs := rr.Result()

			assert.Equal(t, tt.code, rs.StatusCode)
		})
	}
}

func TestRequireRefreshToken(t *testing.T) {
	tdb := setupDB(t)
	defer db.Clean(tdb)

	app := setupApp(t, tdb)

	user := &models.Users{
		ID:       xid.New().String(),
		Name:     "Adam",
		Username: "iamadam",
		Email:    "adam45@gmail.com",
	}

	createdUser, err := app.models.Users.Insert(user)
	require.Nil(t, err)

	accessToken, err := app.newAccessToken(&tokenClaims{ID: createdUser.ID, StdClaims: nil})
	require.Nil(t, err)

	refreshToken, err := app.newRefreshToken(&tokenClaims{ID: createdUser.ID, StdClaims: nil})
	require.Nil(t, err)

	tests := []struct {
		name   string
		header string
		code   int
	}{
		{
			name:   "valid",
			header: "Bearer " + refreshToken,
			code:   http.StatusOK,
		},
		{
			name:   "no token",
			header: "Bearer ",
			code:   http.StatusForbidden,
		},
		{
			name:   "access token",
			header: "Bearer " + accessToken,
			code:   http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			r, err := http.NewRequest(http.MethodPut, "/", nil)
			require.Nil(t, err)

			r.Header.Set("Authorization", tt.header)

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte("OK"))
				require.Nil(t, err)
			})

			app.requireRefreshToken(next).ServeHTTP(rr, r)

			rs := rr.Result()

			assert.Equal(t, tt.code, rs.StatusCode)
		})
	}
}

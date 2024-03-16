package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatValidationErr(t *testing.T) {
	app := setupApp(t, nil)

	t.Run("valid", func(t *testing.T) {
		email := "addam.go"
		err := app.validate.Var(email, "email")
		require.NotNil(t, err)

		mapOfErrs, err := app.formatValidationErr(err)
		require.Nil(t, err)

		require.NotEmpty(t, mapOfErrs)
	})

	t.Run("invalid", func(t *testing.T) {
		err := errors.New("just another error")

		mapOfErrs, err := app.formatValidationErr(err)
		require.NotNil(t, err)

		require.Empty(t, mapOfErrs)
	})
}

func TestNewToken(t *testing.T) {
	app := setupApp(t, nil)

	token := app.newToken()

	assert.Equal(t, len(token), 6)

	secondToken := app.newToken()

	assert.NotEqual(t, token, secondToken)
}

func TestUserIP(t *testing.T) {
	app := setupApp(t, nil)

	r := httptest.NewRequest(http.MethodPost, "/", nil)

	ip := app.userIP(r)

	assert.Equal(t, "86.44.17.109", ip)
}

func TestUserLocation(t *testing.T) {
	app := setupApp(t, nil)

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	ip := app.userIP(r)

	location, err := app.userLocation(ip)
	require.Nil(t, err)

	assert.NotEmpty(t, strings.Contains(location, "Ireland"))
}

func TestGetUserAgent(t *testing.T) {
	app := setupApp(t, nil)
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	ua := app.getUserAgent(r)

	assert.Equal(t, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36", ua)
}

func TestGetUserDeviceInfo(t *testing.T) {
	app := setupApp(t, nil)
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	ua := app.getUserAgent(r)

	info := app.getUserDeviceInfo(ua)

	assert.NotEmpty(t, strings.Contains(info, "Linux"))
}

package main

import (
	"errors"
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

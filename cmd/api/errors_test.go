package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/micahasowata/jason"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestMethodNotAllowed(t *testing.T) {
	r, err := http.NewRequest(http.MethodPut, "/v1/users/register", nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()

	app := setupApp(t, nil)

	app.methodNotAllowed(rr, r)

	rs := rr.Result()

	assert.Equal(t, http.StatusMethodNotAllowed, rs.StatusCode)
}

func TestNotFound(t *testing.T) {
	r, err := http.NewRequest(http.MethodPut, "/v2", nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()

	app := setupApp(t, nil)

	app.notFoundHandler(rr, r)

	rs := rr.Result()

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
}

func TestBadRequestHandler(t *testing.T) {
	app := setupApp(t, nil)

	tests := []struct {
		name string
		err  error
		code int
	}{
		{
			name: "valid",
			err:  &jason.Err{Msg: "invalid body"},
			code: http.StatusBadRequest,
		},
		{
			name: "invalid",
			err:  errors.New("invalid body"),
			code: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			app.badRequestHandler(rr, tt.err)

			rs := rr.Result()

			assert.Equal(t, tt.code, rs.StatusCode)
		})
	}
}

func TestServerErrorHandler(t *testing.T) {
	app := &application{
		Jason:  jason.New(100, false, true),
		logger: zap.NewExample(),
	}

	err := errors.New("just an error")

	rr := httptest.NewRecorder()

	app.badRequestHandler(rr, err)

	rs := rr.Result()

	assert.Equal(t, http.StatusInternalServerError, rs.StatusCode)
}

func TestValidationErrHandler(t *testing.T) {
	app := setupApp(t, nil)

	t.Run("valid", func(t *testing.T) {
		rr := httptest.NewRecorder()

		email := "addam.go"
		err := app.validate.Var(email, "email")
		require.NotNil(t, err)
		app.validationErrHandler(rr, err)

		rs := rr.Result()

		assert.Equal(t, http.StatusUnprocessableEntity, rs.StatusCode)
	})

	t.Run("invalid", func(t *testing.T) {
		rr := httptest.NewRecorder()

		err := errors.New("just another error")
		app.validationErrHandler(rr, err)

		rs := rr.Result()
		assert.Equal(t, http.StatusInternalServerError, rs.StatusCode)
	})
}

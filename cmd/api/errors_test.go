package main

import (
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

	app := &application{
		Jason: jason.New(100, false, true),
	}

	app.methodNotAllowed(rr, r)

	rs := rr.Result()

	assert.Equal(t, http.StatusMethodNotAllowed, rs.StatusCode)
}

func TestNotFound(t *testing.T) {
	r, err := http.NewRequest(http.MethodPut, "/v2", nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()

	app := &application{
		Jason: jason.New(100, false, true),
	}

	app.notFoundHandler(rr, r)

	rs := rr.Result()

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
}

func TestBadRequestHandler(t *testing.T) {
	rr := httptest.NewRecorder()

	err := &jason.Err{
		Msg: "invalid body",
	}

	app := &application{
		Jason:  jason.New(100, false, true),
		logger: zap.NewExample(),
	}

	app.badRequestHandler(rr, err)

	rs := rr.Result()

	assert.Equal(t, http.StatusBadRequest, rs.StatusCode)
}

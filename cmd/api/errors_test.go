package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/micahasowata/jason"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMethodNotAllowed(t *testing.T) {
	r, err := http.NewRequest(http.MethodPut, "/", nil)
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
	r, err := http.NewRequest(http.MethodPut, "/auth", nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()

	app := &application{
		Jason: jason.New(100, false, true),
	}

	app.notFoundHandler(rr, r)

	rs := rr.Result()

	assert.Equal(t, http.StatusNotFound, rs.StatusCode)
}

package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/micahasowata/blog/internal/models"
	"github.com/micahasowata/jason"
)

type errResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Details jason.Envelope `json:"details"`
	Cause   error          `json:"-"`
}

func (app *application) errorResponse(w http.ResponseWriter, e *errResponse) {
	if e.Code == http.StatusInternalServerError && e.Cause != nil {
		app.logger.Error(e.Cause.Error())
	}

	err := app.Write(w, e.Code, jason.Envelope{"error": e}, nil)
	if err != nil {
		w.WriteHeader(e.Code)
		return
	}
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := strings.ToLower(fmt.Sprintf("%s is not supported for %s", r.Method, r.URL.Path))

	e := &errResponse{
		Message: message,
		Code:    http.StatusMethodNotAllowed,
		Details: jason.Envelope{},
		Cause:   nil,
	}

	app.errorResponse(w, e)
}

func (app *application) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	message := strings.ToLower(fmt.Sprintf("%s does not exist", r.URL.Path))

	e := &errResponse{
		Message: message,
		Details: jason.Envelope{},
		Code:    http.StatusNotFound,
		Cause:   nil,
	}

	app.errorResponse(w, e)
}
func (app *application) serverErrorHandler(w http.ResponseWriter, err error) {
	e := &errResponse{
		Message: "request could no longer be process",
		Code:    http.StatusInternalServerError,
		Cause:   err,
	}
	app.errorResponse(w, e)
}

func (app *application) badRequestHandler(w http.ResponseWriter, err error) {
	bodyErr, ok := err.(*jason.Err)
	if !ok {
		app.serverErrorHandler(w, err)
		return
	}

	e := &errResponse{
		Message: bodyErr.Msg,
		Code:    http.StatusBadRequest,
		Cause:   err,
	}
	app.errorResponse(w, e)
}

func (app *application) validationErrHandler(w http.ResponseWriter, err error) {
	validationErrs, err := app.formatValidationErr(err)
	if err != nil {
		app.serverErrorHandler(w, err)
		return
	}

	e := &errResponse{
		Message: "invalid data in request data",
		Details: jason.Envelope{
			"errors": validationErrs,
		},
		Code:  http.StatusUnprocessableEntity,
		Cause: err,
	}

	app.errorResponse(w, e)
}

func (app *application) writeErrHandler(w http.ResponseWriter, err error) {
	app.logger.Error(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}

func (app *application) duplicateUserDataHandler(w http.ResponseWriter, err error) {
	e := &errResponse{
		Code:  http.StatusConflict,
		Cause: err,
	}

	switch {
	case errors.Is(err, models.ErrDuplicateUsername):
		e.Message = err.Error()
		app.errorResponse(w, e)
	case errors.Is(err, models.ErrDuplicateEmail):
		e.Message = err.Error()
		app.errorResponse(w, e)
	default:
		app.serverErrorHandler(w, err)
	}
}

func (app *application) invalidTokenHandler(w http.ResponseWriter, err error) {
	e := &errResponse{
		Code:    http.StatusForbidden,
		Message: "invalid token",
		Cause:   err,
	}

	app.errorResponse(w, e)
}

package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/micahasowata/jason"
)

type errHTTP struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type errResponse struct {
	Message     string         `json:"message"`
	Details     jason.Envelope `json:"details"`
	Description string         `json:"description"`
	Code        int            `json:"code"`
	Response    errHTTP        `json:"response"`
	Cause       error          `json:"-"`
}

func (app *application) errorResponse(w http.ResponseWriter, e *errResponse) {
	if e.Cause != nil {
		app.logger.Error(e.Cause.Error())
	}

	err := app.Write(w, e.Response.Code, jason.Envelope{"error": e}, nil)
	if err != nil {
		w.WriteHeader(e.Response.Code)
		return
	}
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := strings.ToLower(fmt.Sprintf("%s is not supported for %s", r.Method, r.URL.Path))
	description := fmt.Sprintf("%s. check the documentation for acceptable methods", message)

	e := &errResponse{
		Message:     message,
		Details:     jason.Envelope{},
		Description: description,
		Code:        0001,
		Response: errHTTP{
			Message: "invalid method for request URL",
			Code:    http.StatusMethodNotAllowed,
		},
		Cause: nil,
	}

	app.errorResponse(w, e)
}

func (app *application) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	message := strings.ToLower(fmt.Sprintf("%s does not exist", r.URL.Path))
	description := fmt.Sprintf("%s. check the documentation for available routes", message)

	e := &errResponse{
		Message:     message,
		Details:     jason.Envelope{},
		Description: description,
		Code:        0002,
		Response: errHTTP{
			Message: message,
			Code:    http.StatusNotFound,
		},
		Cause: nil,
	}

	app.errorResponse(w, e)
}

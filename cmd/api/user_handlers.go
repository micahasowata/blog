package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/micahasowata/blog/internal/models"
	"github.com/micahasowata/jason"
	"github.com/rs/xid"
)

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name" validate:"required,lte=150"`
		Username string `json:"username" validate:"required,gte=2,lte=25,ascii"`
		Email    string `json:"email" validate:"required,email,lte=150"`
	}

	err := app.Read(w, r, &input)
	if err != nil {
		app.badRequestHandler(w, err)
		return
	}

	err = app.validate.Struct(&input)
	if err != nil {
		app.validationErrHandler(w, err)
	}

	user := &models.Users{
		ID:       xid.New().String(),
		Name:     input.Name,
		Username: input.Username,
		Email:    input.Email,
	}

	token := app.newToken()

	err = app.rclient.Set(r.Context(), token, user.Email, 5*time.Hour).Err()
	if err != nil {
		app.serverErrorHandler(w, err)
		return
	}

	payload := otpEmailPayload{
		Subject: fmt.Sprintf("%s, welcome to Blog", strings.ToLower(user.Name)),
		Name:    user.Name,
		To:      user.Email,
		Token:   token,
		Kind:    "welcome",
	}

	task, err := app.newOTPEmailTask(payload)
	if err != nil {
		app.serverErrorHandler(w, err)
		return
	}

	_, err = app.executor.EnqueueContext(r.Context(), task)
	if err != nil {
		app.serverErrorHandler(w, err)
		return
	}

	user, err = app.models.Users.Insert(user)
	if err != nil {
		app.duplicateUserDataHandler(w, err)
		return
	}

	err = app.Write(w, http.StatusOK, jason.Envelope{"user": user}, nil)
	if err != nil {
		app.writeErrHandler(w, err)
		return
	}
}

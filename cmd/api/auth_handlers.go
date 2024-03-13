package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/micahasowata/blog/internal/models"
	"github.com/micahasowata/jason"
)

func (app *application) verifyEmail(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Token string `json:"token" validate:"required,len=6"`
	}

	err := app.Read(w, r, &input)
	if err != nil {
		app.badRequestHandler(w, err)
		return
	}

	err = app.validate.Struct(&input)
	if err != nil {
		app.validationErrHandler(w, err)
		return
	}

	email, err := app.rclient.Get(r.Context(), input.Token).Result()
	if err != nil || email == "" {
		app.invalidTokenHandler(w, err)
		return
	}

	user, err := app.models.Users.GetByEmail(email)
	if err != nil {
		app.invalidTokenHandler(w, err)
		return
	}

	if user.Verified {
		app.invalidTokenHandler(w, models.ErrUserNotFound)

		err := app.rclient.Del(r.Context(), input.Token).Err()
		if err != nil {
			app.serverErrorHandler(w, err)
			return
		}

		return
	}

	user, err = app.models.Users.VerifyEmail(email)
	if err != nil {
		app.invalidTokenHandler(w, err)
		return
	}

	err = app.rclient.Del(r.Context(), input.Token).Err()
	if err != nil {
		app.serverErrorHandler(w, err)
		return
	}

	err = app.Write(w, http.StatusOK, jason.Envelope{"user": user}, nil)
	if err != nil {
		app.writeErrHandler(w, err)
		return
	}
}

func (app *application) createLoginToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email" validate:"required,email"`
	}

	err := app.Read(w, r, &input)
	if err != nil {
		app.badRequestHandler(w, err)
		return
	}

	err = app.validate.Struct(&input)
	if err != nil {
		app.validationErrHandler(w, err)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			app.resourceNotFoundHandler(w, models.ErrUserNotFound)
		default:
			app.serverErrorHandler(w, err)
		}
		return
	}

	token := app.newToken()

	err = app.rclient.Set(r.Context(), token, user.Email, 5*time.Hour).Err()
	if err != nil {
		app.serverErrorHandler(w, err)
		return
	}

	payload := otpEmailPayload{
		Subject: fmt.Sprintf("%s login token", strings.ToLower(user.Name)),
		Name:    user.Name,
		To:      user.Email,
		Token:   token,
		Kind:    "login_token",
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

	err = app.Write(w, http.StatusOK, jason.Envelope{"user": user}, nil)
	if err != nil {
		app.writeErrHandler(w, err)
		return
	}
}

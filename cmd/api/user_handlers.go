package main

import (
	"errors"
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

func (app *application) getUserProfile(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(userID).(string)

	user, err := app.models.Users.GetByID(id)
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

func (app *application) updateUserProfile(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(userID).(string)

	var input struct {
		Name     *string `json:"name" validate:"omitempty,lte=150"`
		Username *string `json:"username" validate:"omitempty,gte=2,lte=25,ascii"`
		Email    *string `json:"email" validate:"omitempty,email,lte=150"`
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

	user, err := app.models.Users.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			app.resourceNotFoundHandler(w, models.ErrUserNotFound)
		default:
			app.serverErrorHandler(w, err)
		}
		return
	}

	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.Username != nil {
		user.Username = *input.Username
	}

	if input.Email != nil {
		user.Email = *input.Email
	}

	user, err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			app.resourceNotFoundHandler(w, models.ErrUserNotFound)
		default:
			app.duplicateUserDataHandler(w, err)
		}
		return
	}

	err = app.Write(w, http.StatusOK, jason.Envelope{"user": user}, nil)
	if err != nil {
		app.writeErrHandler(w, err)
		return
	}
}

func (app *application) deleteUserProfile(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(userID).(string)

	err := app.models.Users.Delete(id)
	if err != nil {
		app.serverErrorHandler(w, err)
		return
	}

	err = app.Write(w, http.StatusOK, jason.Envelope{"user": "user deleted successfully"}, nil)
	if err != nil {
		app.writeErrHandler(w, err)
		return
	}
}

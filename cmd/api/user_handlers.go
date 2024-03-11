package main

import (
	"net/http"

	"github.com/micahasowata/blog/internal/models"
	"github.com/micahasowata/jason"
	"github.com/rs/xid"
)

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name" validate:"required,lte=150"`
		Username string `json:"username" validate:"required,gte=2,lte=25"`
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

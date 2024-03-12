package main

import (
	"net/http"

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

	user, err := app.models.Users.VerifyEmail(email)
	if err != nil {
		app.invalidTokenHandler(w, err)
		return
	}

	err = app.Write(w, http.StatusOK, jason.Envelope{"user": user}, nil)
	if err != nil {
		app.writeErrHandler(w, err)
		return
	}
}

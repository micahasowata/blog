package main

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func (app *application) formatValidationErr(err error) (map[string]string, error) {
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string]string{}, err
	}

	errs := map[string]string{}

	for _, e := range validationErrs {
		errs[strings.ToLower(e.Field())] = e.Translate(app.translator)
	}

	return errs, nil
}

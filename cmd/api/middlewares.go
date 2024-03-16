package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func (app *application) stack(router *chi.Mux) {
	router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  zap.NewStdLog(app.logger),
		NoColor: true,
	}))
	router.Use(middleware.Recoverer)
	router.Use(middleware.CleanPath)
	router.Use(middleware.RequestID)
	router.Use(middleware.RequestSize(int64(app.config.MaxSize)))
	router.MethodNotAllowed(http.HandlerFunc(app.methodNotAllowed))
	router.NotFound(http.HandlerFunc(app.notFoundHandler))
}

func (app *application) requireAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		values := strings.Split(header, " ")

		if len(values) != 2 || values[0] != "Bearer" || values[1] == "" {
			app.invalidTokenHandler(w, errors.New("invalid auth token"))
			return
		}

		token := values[1]

		claims, err := app.verifyJWT(token)
		if err != nil {
			app.invalidTokenHandler(w, err)
			return
		}

		if claims.StdClaims.Subject != "access" || claims.StdClaims.Issuer != "blog-be" {
			app.invalidTokenHandler(w, fmt.Errorf("tampered token %+v", claims.StdClaims))
			return
		}

		user, err := app.models.Users.GetByID(claims.ID)
		if err != nil {
			app.invalidTokenHandler(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), userToken, token)
		ctx = context.WithValue(ctx, userID, user.ID)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		values := strings.Split(header, " ")

		if len(values) != 2 || values[0] != "Bearer" || values[1] == "" {
			app.invalidTokenHandler(w, errors.New("invalid auth token"))
			return
		}

		token := values[1]

		claims, err := app.verifyJWT(token)
		if err != nil {
			app.invalidTokenHandler(w, err)
			return
		}

		if claims.StdClaims.Subject != "refresh" || claims.StdClaims.Issuer != "blog-be" {
			app.invalidTokenHandler(w, fmt.Errorf("tampered token %+v", claims.StdClaims))
			return
		}

		user, err := app.models.Users.GetByID(claims.ID)
		if err != nil {
			app.invalidTokenHandler(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), userToken, token)
		ctx = context.WithValue(ctx, userID, user.ID)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

package main

import (
	"context"
	"testing"
	"time"

	"github.com/micahasowata/blog/internal/models"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func setUpToken(t *testing.T, app *application, user *models.Users) string {
	token := app.newToken()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	err := app.rclient.Set(ctx, token, user.Email, 5*time.Hour).Err()
	require.Nil(t, err)

	return token
}
func TestNewWelcomeEmailTask(t *testing.T) {
	app := setupApp(t, nil)

	user := &models.Users{
		ID:       xid.New().String(),
		Name:     "addam",
		Username: "iamadam",
		Email:    "addam@gmail.com",
	}

	token := setUpToken(t, app, user)

	task, err := app.newWelcomeEmailTask("iamadam", token, "adam45@gmail.com")
	require.Nil(t, err)
	require.NotNil(t, task)
}

func TestHandleWelcomeEmailDelivery(t *testing.T) {
	app := setupApp(t, nil)
	user := &models.Users{
		ID:       xid.New().String(),
		Name:     "addam",
		Username: "iamadam",
		Email:    "addam@gmail.com",
	}

	token := setUpToken(t, app, user)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	task, err := app.newWelcomeEmailTask("addam", token, "adam45@gmail.com")
	require.Nil(t, err)

	err = app.handleWelcomeEmailDelivery(ctx, task)
	require.Nil(t, err)

}

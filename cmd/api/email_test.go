package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewWelcomeEmailTask(t *testing.T) {
	app := setupApp(t, nil)

	task, err := app.newWelcomeEmailTask("iamadam", "adam45@gmail.com")
	require.Nil(t, err)
	require.NotNil(t, task)
}

func TestHandleWelcomeEmailDelivery(t *testing.T) {
	t.Skip()

	app := setupApp(t, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	task, err := app.newWelcomeEmailTask("addam", "adam45@gmail.com")
	require.Nil(t, err)

	err = app.handleWelcomeEmailDelivery(ctx, task)
	require.Nil(t, err)

}

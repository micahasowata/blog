package main

import (
	"context"

	"github.com/hibiken/asynq"
	jsoniter "github.com/json-iterator/go"
	"github.com/micahasowata/blog/internal/templates"
	"github.com/wneessen/go-mail"
)

const (
	typeWelcomeEmail = "email:welcome"
)

type welcomeEmailPayload struct {
	Name string
	To   string
	From string
}

func (app *application) newWelcomeEmailTask(name, to string) (*asynq.Task, error) {
	payload, err := jsoniter.Marshal(welcomeEmailPayload{Name: name, To: to, From: app.config.From})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(typeWelcomeEmail, payload, asynq.MaxRetry(3)), nil
}

func (app *application) handleWelcomeEmailDelivery(ctx context.Context, t *asynq.Task) error {
	payload := welcomeEmailPayload{}

	err := jsoniter.Unmarshal(t.Payload(), &payload)
	if err != nil {
		return err
	}

	message := mail.NewMsg()

	err = message.From(payload.From)
	if err != nil {
		return err
	}

	err = message.To(payload.To)
	if err != nil {
		return err
	}

	message.Subject("👋🏼 Welcome to Blog 👋🏼")

	err = message.SetBodyHTMLTemplate(templates.Parse("welcome"), &payload)
	if err != nil {
		return err
	}

	client, err := mail.NewClient(app.config.SMTPHost, mail.WithPort(app.config.SMTPPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(app.config.SMTPUsername),
		mail.WithPassword(app.config.SMTPPassword))

	if err != nil {
		return err
	}

	defer client.Close()

	err = client.DialAndSendWithContext(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

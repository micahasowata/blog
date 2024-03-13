package main

import (
	"context"

	"github.com/hibiken/asynq"
	jsoniter "github.com/json-iterator/go"
	"github.com/micahasowata/blog/internal/templates"
	"github.com/wneessen/go-mail"
)

const (
	typeOTPEmail = "email:otp"
)

type otpEmailPayload struct {
	Subject string
	Name    string
	To      string
	Token   string
	Kind    string
}

func (app *application) newOTPEmailTask(payload otpEmailPayload) (*asynq.Task, error) {
	p, err := jsoniter.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(typeOTPEmail, p, asynq.MaxRetry(3)), nil
}

func (app *application) handleOTPEmailDelivery(ctx context.Context, t *asynq.Task) error {
	payload := otpEmailPayload{}

	err := jsoniter.Unmarshal(t.Payload(), &payload)
	if err != nil {
		return err
	}

	message := mail.NewMsg()

	err = message.From(app.config.From)
	if err != nil {
		return err
	}

	err = message.To(payload.To)
	if err != nil {
		return err
	}

	message.Subject(payload.Subject)

	err = message.SetBodyHTMLTemplate(templates.Parse(payload.Kind), &payload)
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

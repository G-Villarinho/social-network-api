package client

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/social-network/config"
	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
	"github.com/mailersend/mailersend-go"
)

type MailerSendClient interface {
	SendEmail(ctx context.Context, payload domain.EmailPayload) error
}

type mailerSendClient struct {
	di     *internal.Di
	client *mailersend.Mailersend
}

func NewMailerSendClient(di *internal.Di) (MailerSendClient, error) {
	mailsender := mailersend.NewMailersend(config.Env.MaileSenderApiToken)

	return &mailerSendClient{
		di:     di,
		client: mailsender,
	}, nil
}

func (m *mailerSendClient) SendEmail(ctx context.Context, payload domain.EmailPayload) error {
	message := m.client.Email.NewMessage()

	message.SetFrom(mailersend.From{
		Name:  payload.FromName,
		Email: payload.From,
	})
	recipients := make([]mailersend.Recipient, len(payload.Recipients))
	for i, r := range payload.Recipients {
		recipients[i] = mailersend.Recipient{
			Name:  r.Name,
			Email: r.Email,
		}
	}
	message.SetRecipients(recipients)
	message.SetSubject(payload.Subject)
	message.SetHTML(payload.HTMLContent)
	message.SetText(payload.TextContent)

	_, err := m.client.Email.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

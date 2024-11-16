package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/G-Villarinho/social-network/client"
	"github.com/G-Villarinho/social-network/config"
	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
)

type emailService struct {
	di     *internal.Di
	client client.MailerSendClient
}

func NewEmailService(di *internal.Di) (domain.EmailService, error) {
	client, err := client.NewMailerSendClient(di)
	if err != nil {
		return nil, err
	}

	return &emailService{
		di:     di,
		client: client,
	}, nil
}

func (e *emailService) SendEmail(ctx context.Context, task domain.EmailPayloadTask) error {
	content, err := renderTemplate(task.Template, task.Params)
	if err != nil {
		return fmt.Errorf("render email template: %w", err)
	}

	payload := domain.EmailPayload{
		From:        config.Env.EmailSender,
		FromName:    "Social Network",
		Recipients:  []domain.Recipient{{Name: task.Recipient.Name, Email: task.Recipient.Email}},
		Subject:     task.Subject,
		HTMLContent: content,
	}

	if err := e.client.SendEmail(ctx, payload); err != nil {
		return fmt.Errorf("send email %s: %w", task.Template, err)
	}

	return nil
}

func renderTemplate(templateName domain.EmailTemplate, params map[string]string) (string, error) {
	content, err := os.ReadFile(filepath.Join("./templates", string(templateName)))
	if err != nil {
		return "", errors.New("read email template: " + err.Error())
	}

	template := string(content)
	for key, value := range params {
		placeholder := fmt.Sprintf("#%s#", key)
		template = strings.ReplaceAll(template, placeholder, value)
	}

	return template, nil
}

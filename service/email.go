package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/G-Villarinho/social-network/client"
	"github.com/G-Villarinho/social-network/config"
	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
)

type emailTemplate string

const (
	OTP                emailTemplate = "otp"
	SignInNotification emailTemplate = "sign-in-notification"
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

func (e *emailService) SendOTP(ctx context.Context, email string, name string, otp string) error {
	duration := strconv.Itoa(config.Env.Cache.CacheExp)

	params := map[string]string{
		"name":     name,
		"otp":      otp,
		"duration": duration,
	}

	content, err := renderTemplate(OTP, params)
	if err != nil {
		return fmt.Errorf("render email template: %w", err)
	}

	payload := domain.EmailPayload{
		From:        config.Env.EmailSender,
		FromName:    "Social Network",
		Recipients:  []domain.Recipient{{Name: name, Email: email}},
		Subject:     "Your OTP Code",
		HTMLContent: content,
	}

	if err := e.client.SendEmail(ctx, payload); err != nil {
		return fmt.Errorf("send OTP: %w", err)
	}

	return nil
}

func (e *emailService) SendSignInNotification(ctx context.Context, payload domain.SignInNotificationPayload) error {
	params := map[string]string{
		"name":      payload.Name,
		"device":    payload.Device,
		"location":  payload.Location,
		"loginTime": payload.LoginTime,
	}

	content, err := renderTemplate(SignInNotification, params)
	if err != nil {
		return fmt.Errorf("render email template: %w", err)
	}

	emailPayload := domain.EmailPayload{
		From:        config.Env.EmailSender,
		FromName:    "Social Network",
		Recipients:  []domain.Recipient{{Name: payload.Name, Email: payload.Email}},
		Subject:     "New Sign-In Detected",
		HTMLContent: content,
	}

	if err := e.client.SendEmail(ctx, emailPayload); err != nil {
		return fmt.Errorf("send sign-in notification: %w", err)
	}

	return nil
}

func renderTemplate(templateName emailTemplate, params map[string]string) (string, error) {
	content, err := os.ReadFile(filepath.Join("./templates", string(templateName)))
	if err != nil {
		return "", errors.New("read email template: " + err.Error())
	}

	template := string(content)
	for key, value := range params {
		placeholder := "#" + key + "#"
		template = strings.ReplaceAll(template, placeholder, value)
	}

	return template, nil
}

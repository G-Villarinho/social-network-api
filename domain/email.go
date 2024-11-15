package domain

import "context"

type EmailPayload struct {
	From        string
	FromName    string
	Recipients  []Recipient
	Subject     string
	HTMLContent string
	TextContent string
}

type Recipient struct {
	Name  string
	Email string
}

type SignInNotificationPayload struct {
	Email     string
	Name      string
	Device    string
	Location  string
	LoginTime string
}

type EmailService interface {
	SendOTP(ctx context.Context, email, name, otp string) error
	SendSignInNotification(ctx context.Context, payload SignInNotificationPayload) error
}

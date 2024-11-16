package domain

import "context"

type EmailTemplate string

const (
	OTP                EmailTemplate = "otp"
	SignInNotification EmailTemplate = "sign-in-notification"
)

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

type EmailPayloadTask struct {
	Template  EmailTemplate
	Subject   string
	Recipient Recipient
	Params    map[string]string
}

type EmailService interface {
	SendEmail(ctx context.Context, payload EmailPayloadTask) error
}

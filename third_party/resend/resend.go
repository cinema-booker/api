package resend

import (
	"fmt"

	resendGo "github.com/resend/resend-go/v2"
)

type Resend struct {
	Client    *resendGo.Client
	FromEmail string
}

func NewResendService(apiKey string, fromEmail string) *Resend {
	return &Resend{
		Client:    resendGo.NewClient(apiKey),
		FromEmail: fromEmail,
	}
}

func (t *Resend) SendEmail(to []string, subject string, content string) (*resendGo.SendEmailResponse, error) {
	params := &resendGo.SendEmailRequest{
		From:    t.FromEmail,
		To:      to,
		Subject: subject,
		Html:    content,
	}

	return t.Client.Emails.Send(params)
}

func (t *Resend) SendPasswordResetEmail(to []string, code string) (*resendGo.SendEmailResponse, error) {
	content := fmt.Sprintf("<p>Here is your temporary passowrd : <strong>%s</strong></p><p>Don't forget to update it once you logged in.</p>", code)

	return t.SendEmail(to, "Reset password", content)
}

package providers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mailgun/mailgun-go/v4"
)

type MailGun struct {
	mg mailgun.Mailgun
}

func (m *MailGun) Send(ctx context.Context, to string, from string, subject string, body string, htmlBody, replyTo *string) error {
	msg := m.mg.NewMessage(from, subject, body, to)
	if htmlBody != nil {
		msg.SetHtml(*htmlBody)
	}
	if replyTo != nil {
		msg.SetReplyTo(*replyTo)
	}

	_, _, err := m.mg.Send(ctx, msg)
	return err
}

func (m *MailGun) SendBatch(ctx context.Context, to []string, from string, subject string, body string, htmlBody, replyTo *string) error {
	msg := m.mg.NewMessage(from, subject, body)
	for _, address := range to {
		if err := msg.AddRecipient(address); err != nil {
			return err
		}
	}
	if htmlBody != nil {
		msg.SetHtml(*htmlBody)
	}
	if replyTo != nil {
		msg.SetReplyTo(*replyTo)
	}

	_, _, err := m.mg.Send(ctx, msg)
	return err
}

// NewMailgun creates a new MailGun email provider
func NewMailgun(id string) (Provider, error) {
	envId := strings.ToUpper(id)
	apiKey := os.Getenv(fmt.Sprintf("MAILER_PROVIDER_%s_API_KEY", envId))
	if len(apiKey) == 0 {
		return nil, fmt.Errorf("missing api key for mailgun provider %q", id)
	}
	domain := os.Getenv(fmt.Sprintf("MAILER_PROVIDER_%s_DOMAIN", envId))
	if len(domain) == 0 {
		return nil, fmt.Errorf("missing domain for mailgun provider %q", id)
	}

	return &MailGun{
		mg: mailgun.NewMailgun(domain, apiKey),
	}, nil
}

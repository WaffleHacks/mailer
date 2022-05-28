package providers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/cenkalti/backoff/v4"
	"github.com/mailgun/mailgun-go/v4"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/WaffleHacks/mailer/logging"
)

var mailgunTracer = otel.Tracer("github.com/WaffleHacks/mailer/providers/mailgun")

type MailGun struct {
	mg mailgun.Mailgun
}

func (m *MailGun) Send(ctx context.Context, l *logging.Logger, to, from, subject, body string, htmlBody, replyTo *string) error {
	_, span := mailgunTracer.Start(ctx, "send")
	defer span.End()

	msg := m.mg.NewMessage(from, subject, body, to)
	if htmlBody != nil {
		msg.SetHtml(*htmlBody)
	}
	if replyTo != nil {
		msg.SetReplyTo(*replyTo)
	}
	l.Debug("constructed message", zap.Bool("has-html", htmlBody != nil), zap.Bool("has-reply-to", replyTo != nil))

	return backoff.Retry(func() error {
		_, _, err := m.mg.Send(ctx, msg)
		if err, ok := err.(*mailgun.UnexpectedResponseError); ok && err.Actual == http.StatusTooManyRequests {
			l.Warn("rate limit encountered, backing off and retrying")
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		} else {
			return nil
		}
	}, backoff.NewExponentialBackOff())
}

func (m *MailGun) SendBatch(ctx context.Context, l *logging.Logger, to []string, from, subject, body string, htmlBody, replyTo *string) error {
	_, span := mailgunTracer.Start(ctx, "send-batch")
	defer span.End()

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
	l.Debug("constructed message", zap.Bool("has-html", htmlBody != nil), zap.Bool("has-reply-to", replyTo != nil), zap.Int("recipients", len(to)))

	return backoff.Retry(func() error {
		_, _, err := m.mg.Send(ctx, msg)
		if err, ok := err.(*mailgun.UnexpectedResponseError); ok && err.Actual == http.StatusTooManyRequests {
			l.Warn("rate limit encountered, backing off and retrying")
			return err
		} else if err != nil {
			return backoff.Permanent(err)
		} else {
			return nil
		}
	}, backoff.NewExponentialBackOff())
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

package providers

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"

	"github.com/WaffleHacks/mailer/logging"
)

type SMTP struct {
	dialer *gomail.Dialer
	sender gomail.SendCloser
	open   bool
}

func (s *SMTP) Send(ctx context.Context, l *logging.Logger, to, from, subject, body string, htmlBody, replyTo *string) error {
	// Reconnect if necessary
	if !s.open {
		sender, err := s.dialer.Dial()
		if err != nil {
			l.Error("failed to reconnect to smtp provider")
			return err
		}
		l.Info("reconnected to smtp provider")
		s.sender = sender
		s.open = true
	}

	// Construct the message
	msg := gomail.NewMessage(gomail.SetCharset("UTF-8"))
	msg.SetHeader("To", to)
	msg.SetHeader("From", from)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)
	if htmlBody != nil {
		msg.AddAlternative("text/html", *htmlBody)
	}
	if replyTo != nil {
		msg.SetHeader("Reply-To", *replyTo)
	}

	// Send the message
	if err := s.sender.Send(from, []string{to}, msg); err != nil {
		s.open = false
		return err
	}

	return nil
}

func NewSMTP(id string) (Provider, error) {
	envId := strings.ToUpper(id)
	host := os.Getenv(fmt.Sprintf("MAILER_PROVIDER_%s_HOST", envId))
	if len(host) == 0 {
		return nil, fmt.Errorf("missing host for smtp provider %q", id)
	}
	port, err := strconv.Atoi(os.Getenv(fmt.Sprintf("MAILER_PROVIDER_%s_PORT", envId)))
	if err != nil {
		return nil, fmt.Errorf("invalid port for smtp provider %q: %v", id, err)
	}
	ssl := strings.ToLower(os.Getenv(fmt.Sprintf("MAILER_PROVIDER_%s_SSL", envId)))
	username := os.Getenv(fmt.Sprintf("MAILER_PROVIDER_%s_USERNAME", envId))
	if len(username) == 0 {
		return nil, fmt.Errorf("missing username for smtp proivder %q", id)
	}
	password := os.Getenv(fmt.Sprintf("MAILER_PROVIDER_%s_PASSWORD", envId))
	if len(password) == 0 {
		return nil, fmt.Errorf("missing password for smtp provider %q", id)
	}

	// Create the dialer
	dialer := gomail.NewDialer(host, port, username, password)
	dialer.SSL = ssl == "y" || ssl == "yes" || ssl == "t" || ssl == "true"

	// Test the connection
	sender, err := dialer.Dial()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to smtp provider %q: %v", id, err)
	}

	return &SMTP{
		dialer: dialer,
		sender: sender,
		open:   true,
	}, nil
}

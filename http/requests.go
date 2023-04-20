package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/WaffleHacks/mailer/daemon"
	"github.com/WaffleHacks/mailer/logging"
)

// CommonRequest contains all the shared fields between the different sending types
type CommonRequest struct {
	From    string         `json:"from"`
	Subject string         `json:"subject"`
	Body    string         `json:"body"`
	Format  *daemon.Format `json:"format,omitempty"`
	ReplyTo *string        `json:"replyTo,omitempty"`
}

type SendRequest struct {
	// To is a single recipient email
	To string `json:"to"`
	CommonRequest
}

type SendBatchRequest struct {
	// To is a list of recipient emails
	To []string `json:"to"`
	CommonRequest
}

type SendTemplateRequest struct {
	// To is a map of all the recipient emails to their context data
	To map[string]map[string]string `json:"to"`
	CommonRequest
}

func parse[T any](r *http.Request) (T, error) {
	var body T
	err := json.NewDecoder(r.Body).Decode(&body)

	return body, err
}

func setCommonSpanAttributes(span trace.Span, common CommonRequest) {
	span.SetAttributes(fromAttr.String(common.From), subjectAttr.String(common.Subject))
	if common.ReplyTo != nil {
		span.SetAttributes(replyToAttr.String(*common.ReplyTo))
	}
	if common.Format != nil {
		span.SetAttributes(formatAttr.String(string(*common.Format)))
	}
}

func validate[T any](ctx context.Context, recipients T, common CommonRequest, recipientValidator func(T) error) error {
	_, span := tracer.Start(ctx, "validate")
	defer span.End()

	// Ensure inputs exist
	if len(common.From) == 0 {
		return errors.New("from is required")
	}
	if len(common.Subject) == 0 {
		return errors.New("subject is required")
	}
	if len(common.Body) == 0 {
		return errors.New("body is required")
	}

	// Ensure format is valid
	if common.Format != nil {
		if *common.Format != daemon.FormatPlain && *common.Format != daemon.FormatHTML {
			return errors.New(`invalid format, must be "PLAIN" or "HTML"`)
		}
	}

	// Validate email addresses are correct
	if err := recipientValidator(recipients); err != nil {
		return err
	}
	if err := isEmailValid(common.From, "from"); err != nil {
		return err
	}
	if common.ReplyTo != nil {
		if err := isEmailValid(*common.ReplyTo, "replyTo"); err != nil {
			return err
		}
	}

	return nil
}

func isEmailValid(address string, field string) error {
	if _, err := mail.ParseAddress(address); err != nil {
		logging.L().Named("http.validator").Warn("invalid email address format", zap.String("field", field), zap.String("email", address))
		return fmt.Errorf("invalid %s email address format", field)
	}

	return nil
}

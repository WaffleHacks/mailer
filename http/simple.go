package http

import (
	"context"
	"errors"
	"net/http"

	"go.opentelemetry.io/otel/trace"

	"github.com/WaffleHacks/mailer/daemon"
)

func (m *mailerServer) send(w http.ResponseWriter, r *http.Request) {
	body, err := parse[SendRequest](r)
	if err != nil {
		deserializationFailure(w)
		return
	}

	m.process(r.Context(), w, []string{body.To}, body.CommonRequest)
}

func (m *mailerServer) sendBatch(w http.ResponseWriter, r *http.Request) {
	body, err := parse[SendBatchRequest](r)
	if err != nil {
		deserializationFailure(w)
		return
	}

	m.process(r.Context(), w, body.To, body.CommonRequest)
}

func (m *mailerServer) process(ctx context.Context, w http.ResponseWriter, recipients []string, common CommonRequest) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(toAttr.StringSlice(recipients))
	setCommonSpanAttributes(span, common)

	// Validate the request body
	if err := validate(ctx, recipients, common, simpleRecipientsValidator); err != nil {
		failure(w, http.StatusBadRequest, err.Error())
		return
	}

	// Assume the format is plaintext if not specified
	format := daemon.FormatPlain
	if common.Format != nil {
		format = *common.Format
	}

	// Enqueue the message for sending
	m.queue <- daemon.Message{
		To:          recipients,
		From:        common.From,
		Subject:     common.Subject,
		Body:        common.Body,
		Type:        format,
		ReplyTo:     common.ReplyTo,
		SpanContext: span.SpanContext(),
	}

	success(w)
}

func simpleRecipientsValidator(recipients []string) error {
	if len(recipients) == 0 {
		return errors.New("at least one to address is required")
	}

	for _, recipient := range recipients {
		if err := isEmailValid(recipient, "to"); err != nil {
			return err
		}
	}

	return nil
}

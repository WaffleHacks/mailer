package http

import (
	"context"
	"errors"
	"net/http"

	"go.opentelemetry.io/otel/trace"

	"github.com/WaffleHacks/mailer/daemon"
	"github.com/WaffleHacks/mailer/template"
)

func (m *mailerServer) sendTemplate(w http.ResponseWriter, r *http.Request) {
	request, err := parse[SendTemplateRequest](r)
	if err != nil {
		deserializationFailure(w)
		return
	}

	current := trace.SpanFromContext(r.Context())
	current.SetAttributes(toAttr.Int(len(request.To)))
	setCommonSpanAttributes(current, request.CommonRequest)

	// Validate the request
	if err := validate(r.Context(), request.To, request.CommonRequest, templatedRecipientsValidator); err != nil {
		failure(w, http.StatusBadRequest, err.Error())
		return
	}

	// Assume the format is plaintext if not specified
	format := daemon.FormatPlain
	if request.Format != nil {
		format = *request.Format
	}
	isHtml := format == daemon.FormatHTML

	tmpl, err := buildTemplate(r.Context(), request.Body)

	// Template and send each message
	for email, ctx := range request.To {
		_, span := tracer.Start(r.Context(), "send")

		m.queue <- daemon.Message{
			To:          []string{email},
			From:        request.From,
			Subject:     request.Subject,
			Body:        tmpl.Render(ctx, isHtml),
			Type:        format,
			ReplyTo:     request.ReplyTo,
			SpanContext: span.SpanContext(),
		}

		span.End()
	}

	success(w)
}

func templatedRecipientsValidator(recipients map[string]map[string]string) error {
	if len(recipients) == 0 {
		return errors.New("at least one to address is required")
	}

	for email := range recipients {
		if err := isEmailValid(email, "to"); err != nil {
			return err
		}
	}

	return nil
}

func buildTemplate(ctx context.Context, body string) (*template.Template, error) {
	_, span := tracer.Start(ctx, "template")
	defer span.End()

	return template.New(body)
}

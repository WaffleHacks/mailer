package rpc

import (
	"context"

	"github.com/WaffleHacks/mailer/daemon"
	mailerv1 "github.com/WaffleHacks/mailer/gen/proto/go/mailer/v1"
	"github.com/WaffleHacks/mailer/logging"
	"github.com/WaffleHacks/mailer/template"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m *mailerServer) SendTemplate(ctx context.Context, in *mailerv1.SendTemplateRequest) (*mailerv1.SendTemplateResponse, error) {
	current := trace.SpanFromContext(ctx)
	current.SetAttributes(
		attribute.Int("recipients", len(in.To)),
		fromAttr.String(in.From),
		typeAttr.String(in.Type.String()),
		subjectAttr.String(in.Subject),
	)

	s := logging.GRPCRequest("SendTemplate", func(l *logging.Logger) *status.Status {
		// Validate the message
		replyTo, s := validateTemplated(ctx, l, in)
		if s != nil {
			return s
		}

		// Create the template
		_, tmplSpan := tracer.Start(ctx, "template")
		tmpl, err := template.New(in.Body)
		if err != nil {
			tmplSpan.End()
			return status.New(codes.InvalidArgument, err.Error())
		}
		tmplSpan.End()

		isHtml := false
		bodyType := mailerv1.BodyType_BODY_TYPE_PLAIN
		if in.Type == mailerv1.BodyType_BODY_TYPE_HTML {
			bodyType = mailerv1.BodyType_BODY_TYPE_HTML
			isHtml = true
		}

		// Template each message body
		for email, rawMap := range in.To {
			_, span := tracer.Start(ctx, "send")

			// Ensure the recipient is valid
			if s := isEmailValid(email, l, "invalid to email address format"); s != nil {
				span.End()
				return s
			}

			// Construct the template context map
			templateCtx, s := buildTemplateContext(ctx, rawMap)
			if s != nil {
				span.End()
				return s
			}

			// Send the message
			m.queue <- daemon.Message{
				To:          []string{email},
				From:        in.From,
				Subject:     in.Subject,
				Body:        tmpl.Render(templateCtx, isHtml),
				Type:        daemon.BodyType(bodyType),
				ReplyTo:     replyTo,
				SpanContext: span.SpanContext(),
			}
			span.End()
		}

		return nil
	})

	if s != nil {
		return nil, s.Err()
	}
	return &mailerv1.SendTemplateResponse{}, nil
}

// Validate the from address, subject, and body
func validateTemplated(ctx context.Context, logger *logging.Logger, in *mailerv1.SendTemplateRequest) (*string, *status.Status) {
	_, span := tracer.Start(ctx, "validate")
	defer span.End()

	// Ensure all inputs exist
	if len(in.To) == 0 {
		return nil, status.New(codes.InvalidArgument, "to is required")
	}
	if len(in.From) == 0 {
		return nil, status.New(codes.InvalidArgument, "from is required")
	}
	if len(in.Subject) == 0 {
		return nil, status.New(codes.InvalidArgument, "subject is required")
	}
	if len(in.Body) == 0 {
		return nil, status.New(codes.InvalidArgument, "body is required")
	}

	// Validate the email addresses have the proper format
	if s := isEmailValid(in.From, logger, "invalid from email address format"); s != nil {
		return nil, s
	}
	var replyTo *string
	if len(in.ReplyTo) != 0 {
		if s := isEmailValid(in.ReplyTo, logger, "invalid reply to email address format"); s != nil {
			return nil, s
		}
		replyTo = &in.ReplyTo
	}

	return replyTo, nil
}

func buildTemplateContext(ctx context.Context, raw *mailerv1.TemplateContext) (map[string]string, *status.Status) {
	_, span := tracer.Start(ctx, "template-context")
	defer span.End()

	templateCtx := make(map[string]string)

	if raw != nil {
		// Ensure it is valid
		if len(raw.Key) != len(raw.Value) {
			return nil, status.New(codes.InvalidArgument, "length of keys and values must match")
		}

		for i, key := range raw.Key {
			templateCtx[key] = raw.Value[i]
		}
	}

	return templateCtx, nil
}

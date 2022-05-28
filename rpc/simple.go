package rpc

import (
	"context"

	"github.com/WaffleHacks/mailer/daemon"
	mailerv1 "github.com/WaffleHacks/mailer/gen/proto/go/mailer/v1"
	"github.com/WaffleHacks/mailer/logging"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m *mailerServer) Send(ctx context.Context, in *mailerv1.SendRequest) (*mailerv1.SendResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		toAttr.String(in.To),
		fromAttr.String(in.From),
		typeAttr.String(in.Type.String()),
		subjectAttr.String(in.Subject),
	)

	s := logging.GRPCRequest("Send", func(l *logging.Logger) *status.Status {
		return m.process(ctx, l, &mailerv1.SendBatchRequest{
			To:      []string{in.To},
			From:    in.From,
			Subject: in.Subject,
			Body:    in.Body,
			Type:    in.Type,
			ReplyTo: in.ReplyTo,
		})
	})

	if s != nil {
		return nil, s.Err()
	}
	return &mailerv1.SendResponse{}, nil
}

func (m *mailerServer) SendBatch(ctx context.Context, in *mailerv1.SendBatchRequest) (*mailerv1.SendBatchResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		toAttr.StringSlice(in.To),
		fromAttr.String(in.From),
		typeAttr.String(in.Type.String()),
		subjectAttr.String(in.Subject),
	)

	s := logging.GRPCRequest("SendBatch", func(l *logging.Logger) *status.Status {
		return m.process(ctx, l, in)
	})

	if s != nil {
		return nil, s.Err()
	}
	return &mailerv1.SendBatchResponse{}, nil
}

// Do the work of processing the messages
func (m *mailerServer) process(ctx context.Context, logger *logging.Logger, in *mailerv1.SendBatchRequest) *status.Status {
	_, span := tracer.Start(ctx, "process")
	defer span.End()

	// Validate the message
	replyTo, err := validate(ctx, logger, in)
	if err != nil {
		return err
	}

	// Set defaults
	bodyType := in.Type
	if bodyType == mailerv1.BodyType_BODY_TYPE_UNSPECIFIED {
		bodyType = mailerv1.BodyType_BODY_TYPE_PLAIN
	}

	// Enqueue the message for sending
	m.queue <- daemon.Message{
		To:      in.To,
		From:    in.From,
		Subject: in.Subject,
		Body:    in.Body,
		Type:    daemon.BodyType(bodyType),
		ReplyTo: replyTo,
	}

	return nil
}

// Validate the message contents
func validate(ctx context.Context, logger *logging.Logger, in *mailerv1.SendBatchRequest) (*string, *status.Status) {
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
	for _, to := range in.To {
		if s := isEmailValid(to, logger, "invalid to email address format"); s != nil {
			return nil, s
		}
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

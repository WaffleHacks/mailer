package rpc

import (
	"context"
	"net/mail"

	"github.com/WaffleHacks/mailer/template"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/WaffleHacks/mailer/daemon"
	mailerv1 "github.com/WaffleHacks/mailer/gen/proto/go/mailer/v1"
	"github.com/WaffleHacks/mailer/logging"
)

type mailerServer struct {
	queue chan daemon.Message
}

// New creates a new gRPC interface to serve
func New(queue chan daemon.Message) *grpc.Server {
	server := grpc.NewServer()
	mailerv1.RegisterMailerServiceServer(server, &mailerServer{queue: queue})
	return server
}

func (m *mailerServer) Send(_ context.Context, in *mailerv1.SendRequest) (*mailerv1.SendResponse, error) {
	s := logging.GRPCRequest("Send", func(l *logging.Logger) *status.Status {
		return m.process(l, &mailerv1.SendBatchRequest{
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

func (m *mailerServer) SendBatch(_ context.Context, in *mailerv1.SendBatchRequest) (*mailerv1.SendBatchResponse, error) {
	s := logging.GRPCRequest("SendBatch", func(l *logging.Logger) *status.Status {
		return m.process(l, in)
	})

	if s != nil {
		return nil, s.Err()
	}
	return &mailerv1.SendBatchResponse{}, nil
}

func (m *mailerServer) SendTemplate(_ context.Context, in *mailerv1.SendTemplateRequest) (*mailerv1.SendTemplateResponse, error) {
	s := logging.GRPCRequest("SendTemplate", func(l *logging.Logger) *status.Status {
		// Create the template
		tmpl, err := template.New(in.Body)
		if err != nil {
			return status.New(codes.InvalidArgument, err.Error())
		}

		isHtml := false
		if in.Type == mailerv1.BodyType_BODY_TYPE_HTML {
			isHtml = true
		}

		// Template each message body
		for email, rawMap := range in.To {
			ctx := make(map[string]string)

			// Populate the context
			if rawMap != nil {
				// Ensure it is valid
				if len(rawMap.Key) != len(rawMap.Value) {
					return status.New(codes.InvalidArgument, "length of keys and values must match")
				}

				for i, key := range rawMap.Key {
					ctx[key] = rawMap.Value[i]
				}
			}

			// Send the message
			m.process(l, &mailerv1.SendBatchRequest{
				To:      []string{email},
				From:    in.From,
				Subject: in.Subject,
				Body:    tmpl.Render(ctx, isHtml),
				Type:    in.Type,
				ReplyTo: in.ReplyTo,
			})
		}

		return nil
	})

	if s != nil {
		return nil, s.Err()
	}
	return &mailerv1.SendTemplateResponse{}, nil
}

// Do the work of processing the messages
func (m *mailerServer) process(logger *logging.Logger, in *mailerv1.SendBatchRequest) *status.Status {
	// Ensure all inputs exist
	if len(in.To) == 0 {
		return status.New(codes.InvalidArgument, "to is required")
	}
	if len(in.From) == 0 {
		return status.New(codes.InvalidArgument, "from is required")
	}
	if len(in.Subject) == 0 {
		return status.New(codes.InvalidArgument, "subject is required")
	}
	if len(in.Body) == 0 {
		return status.New(codes.InvalidArgument, "body is required")
	}

	// Validate the email addresses have the proper format
	if _, err := mail.ParseAddress(in.From); err != nil {
		logger.Warn("invalid from email address format", zap.String("from", in.From))
		return status.New(codes.InvalidArgument, "invalid from email address format")
	}
	for _, to := range in.To {
		if _, err := mail.ParseAddress(to); err != nil {
			logger.Warn("invalid to email address format", zap.String("to", to))
			return status.New(codes.InvalidArgument, "invalid to email address format")
		}
	}
	var replyTo *string
	if len(in.ReplyTo) != 0 {
		if _, err := mail.ParseAddress(in.ReplyTo); err != nil {
			logger.Warn("invalid email address format", zap.String("reply-to", in.ReplyTo))
			return status.New(codes.InvalidArgument, "invalid reply to email address format")
		}
		replyTo = &in.ReplyTo
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

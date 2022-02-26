package rpc

import (
	"context"

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
	s := logging.GRPCRequest("Send", func(l *zap.Logger) *status.Status {
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
	s := logging.GRPCRequest("SendBatch", func(l *zap.Logger) *status.Status {
		return m.process(l, in)
	})

	if s != nil {
		return nil, s.Err()
	}
	return &mailerv1.SendBatchResponse{}, nil
}

// Do the work of processing the messages
func (m *mailerServer) process(logger *zap.Logger, in *mailerv1.SendBatchRequest) *status.Status {
	return status.New(codes.Unimplemented, "unimplemented")
}

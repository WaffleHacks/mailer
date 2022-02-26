package rpc

import (
	"context"

	"google.golang.org/grpc"

	mailerv1 "github.com/WaffleHacks/mailer/gen/proto/go/mailer/v1"
)

type mailerServer struct{}

// New creates a new gRPC interface to serve
func New() *grpc.Server {
	server := grpc.NewServer()
	mailerv1.RegisterMailerServiceServer(server, &mailerServer{})
	return server
}

func (m *mailerServer) Send(ctx context.Context, in *mailerv1.SendRequest) (*mailerv1.SendResponse, error) {
	return &mailerv1.SendResponse{}, nil
}

func (m *mailerServer) SendBatch(ctx context.Context, in *mailerv1.SendBatchRequest) (*mailerv1.SendBatchResponse, error) {
	return &mailerv1.SendBatchResponse{}, nil
}

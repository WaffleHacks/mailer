package rpc

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	mailerv1 "github.com/WaffleHacks/mailer/gen/proto/go/mailer/v1"
)

// NewGateway creates a new HTTP gateway mapping a HTTP request to a gRPC request
func NewGateway(grpcAddress, httpAddress string) (*http.Server, error) {
	// Connect to the gRPC API
	conn, err := grpc.DialContext(context.Background(), "dns:///"+grpcAddress, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Create the handler
	mux := runtime.NewServeMux()
	if err := mailerv1.RegisterMailerServiceHandler(context.Background(), mux, conn); err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:    httpAddress,
		Handler: mux,
	}
	return server, nil
}

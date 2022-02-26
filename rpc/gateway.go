package rpc

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/WaffleHacks/mailer/docs"
	mailerv1 "github.com/WaffleHacks/mailer/gen/proto/go/mailer/v1"
	"github.com/WaffleHacks/mailer/logging"
)

// NewGateway creates a new HTTP gateway mapping a HTTP request to a gRPC request
func NewGateway(grpcAddress, httpAddress string) (*http.Server, error) {
	// Connect to the gRPC API
	conn, err := grpc.DialContext(context.Background(), "dns:///"+grpcAddress, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Create the proxy handler
	mux := runtime.NewServeMux()
	if err := mailerv1.RegisterMailerServiceHandler(context.Background(), mux, conn); err != nil {
		return nil, err
	}

	// Create the docs handler
	documentation, err := docs.Handler()
	if err != nil {
		return nil, err
	}

	// Create a wrapping mux
	r := chi.NewRouter()
	r.Use(logging.Request(zap.L().Named("http")))
	r.Handle("/docs/*", documentation)
	r.Handle("/*", mux)

	server := &http.Server{
		Addr:    httpAddress,
		Handler: r,
	}
	return server, nil
}

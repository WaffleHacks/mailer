package rpc

import (
	"github.com/WaffleHacks/mailer/daemon"
	mailerv1 "github.com/WaffleHacks/mailer/gen/proto/go/mailer/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
)

var (
	tracer = otel.Tracer("github.com/WaffleHacks/mailer/rpc")

	toAttr      = attribute.Key("mailer.to")
	fromAttr    = attribute.Key("mailer.from")
	typeAttr    = attribute.Key("mailer.type")
	subjectAttr = attribute.Key("mailer.subject")
)

type mailerServer struct {
	queue chan daemon.Message
}

// New creates a new gRPC interface to serve
func New(queue chan daemon.Message) *grpc.Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	mailerv1.RegisterMailerServiceServer(server, &mailerServer{queue: queue})
	return server
}

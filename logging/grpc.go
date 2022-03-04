package logging

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	gonanoid "github.com/matoous/go-nanoid"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

// GRPCRequest acts as middleware to add logging to a gRPC request
func GRPCRequest(method string, action func(l *Logger) *status.Status) *status.Status {
	span := sentry.StartSpan(context.Background(), method)

	l := L().Named("grpc").With(zap.String("method", method), zap.String("id", gonanoid.MustID(8)))
	l.Info("started processing request")
	start := time.Now()

	s := action(l)

	latency := float64(time.Now().Sub(start).Nanoseconds()) / 1000000.0
	l.Info("finished processing request", zap.String("status", s.Code().String()), zap.Float64("latency", latency))

	span.Finish()
	return s
}

package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

// Config contains all the runtime configuration for OpenTelemetry tracing
type Config struct {
	Enable   bool
	Endpoint string
	Headers  map[string]string
}

func newExporter(ctx context.Context, development bool) (trace.SpanExporter, error) {
	if development {
		return jaeger.New(
			jaeger.WithAgentEndpoint(
				jaeger.WithAgentPort("6231"),
			),
		)
	} else {
		client := otlptracegrpc.NewClient(
			otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
		)
		return otlptrace.New(ctx, client)
	}
}

func newTraceProvider(exp trace.SpanExporter) (*trace.TracerProvider, error) {
	resouces, err := resource.New(context.Background(), resource.WithFromEnv())
	if err != nil {
		return nil, err
	}

	return trace.NewTracerProvider(trace.WithBatcher(exp), trace.WithResource(resouces)), nil
}

// Initialize sets up OpenTelemetry tracing if requested
func Initialize(ctx context.Context, enable, development bool) (*trace.TracerProvider, error) {
	if enable {
		exporter, err := newExporter(ctx, development)
		if err != nil {
			return nil, err
		}

		provider, err := newTraceProvider(exporter)
		if err != nil {
			return nil, err
		}

		otel.SetTracerProvider(provider)
		otel.SetTextMapPropagator(
			propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			),
		)

		return provider, nil
	}

	return nil, nil
}

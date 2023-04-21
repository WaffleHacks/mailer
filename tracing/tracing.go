package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"

	"github.com/WaffleHacks/mailer/version"
)

func newTraceProvider(exp trace.SpanExporter) (*trace.TracerProvider, error) {
	serviceVersion := version.Commit
	if version.Dirty {
		serviceVersion += "-dirty"
	}

	r, err := resource.New(context.Background(),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(semconv.ServiceVersion(serviceVersion)),
	)
	if err != nil {
		return nil, err
	}

	return trace.NewTracerProvider(trace.WithBatcher(exp), trace.WithResource(r)), nil
}

// Initialize sets up OpenTelemetry tracing if requested
func Initialize(ctx context.Context, enable bool) (*trace.TracerProvider, error) {
	if enable {
		exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient())
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

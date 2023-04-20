package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func newTraceProvider(exp trace.SpanExporter) (*trace.TracerProvider, error) {
	resouces, err := resource.New(context.Background(), resource.WithFromEnv())
	if err != nil {
		return nil, err
	}

	return trace.NewTracerProvider(trace.WithBatcher(exp), trace.WithResource(resouces)), nil
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

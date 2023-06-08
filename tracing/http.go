package tracing

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// AddHeaders is an HTTP middleware to enrich the current span with the request and response headers
func AddHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		span := trace.SpanFromContext(r.Context())
		addHeadersToSpan(span, r.Header)

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		addHeadersToSpan(span, ww.Header())
	})
}

func addHeadersToSpan(span trace.Span, headers http.Header) {
	for name, values := range headers {
		lowerName := strings.ToLower(name)
		key := attribute.Key("http.request.header." + strings.ReplaceAll(lowerName, "-", "_"))

		var attr attribute.KeyValue
		if lowerName == "cookie" || lowerName == "set-cookie" || lowerName == "authorization" {
			attr = key.StringSlice([]string{"[REDACTED]"})
		} else {
			attr = key.StringSlice(values)
		}

		span.SetAttributes(attr)
	}
}

package http

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"

	"github.com/WaffleHacks/mailer/daemon"
	"github.com/WaffleHacks/mailer/logging"
	"github.com/WaffleHacks/mailer/version"
)

var (
	ErrServerClosed = http.ErrServerClosed

	tracer = otel.Tracer("github.com/WaffleHacks/mailer/http")

	toAttr      = attribute.Key("mailer.to")
	fromAttr    = attribute.Key("mailer.from")
	formatAttr  = attribute.Key("mailer.format")
	subjectAttr = attribute.Key("mailer.subject")
	replyToAttr = attribute.Key("mailer.replyTo")
)

type mailerServer struct {
	queue chan daemon.Message
}

// New creates a new HTTP server
func New(address string, queue chan daemon.Message) *http.Server {
	router := chi.NewRouter()

	m := &mailerServer{queue: queue}

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logging.Request(zap.L().Named("http")))
	router.Use(middleware.Recoverer)

	router.Get("/health", healthcheck)
	router.Post("/send", m.send)
	router.Post("/send/batch", m.sendBatch)
	router.Post("/send/template", m.sendTemplate)

	return &http.Server{
		Addr: address,
		Handler: otelhttp.NewHandler(
			router,
			"request",
			otelhttp.WithFilter(routeFilter),
			otelhttp.WithServerName(serverName()),
		),
	}
}

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(version.Printable))
}

func serverName() string {
	hostname, err := os.Hostname()
	if err != nil {
		zap.L().Fatal("failed to retrieve hostname", zap.Error(err))
	}
	return hostname
}

// routeFilter checks if a route should be ignored by OpenTelemetry
func routeFilter(r *http.Request) bool {
	if r.Method == http.MethodGet && r.URL.Path == "/health" {
		return false
	}

	return true
}

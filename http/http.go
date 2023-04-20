package http

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"

	"github.com/WaffleHacks/mailer/daemon"
	"github.com/WaffleHacks/mailer/logging"
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
	r := chi.NewRouter()

	m := &mailerServer{queue: queue}

	tracingMiddleware := otelchi.Middleware(serverName(), otelchi.WithChiRoutes(r), otelchi.WithFilter(routeFilter))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(tracingMiddleware)
	r.Use(logging.Request(logging.L().Named("http")))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Post("/send", m.send)
	r.Post("/send/batch", m.sendBatch)
	r.Post("/send/template", m.sendTemplate)

	return &http.Server{
		Addr:    address,
		Handler: r,
	}
}

func serverName() string {
	hostname, err := os.Hostname()
	if err != nil {
		logging.L().Fatal("failed to retrieve hostname", zap.Error(err))
	}
	return hostname
}

// routeFilter checks if a route should be ignored by OpenTelemetry
func routeFilter(r *http.Request) bool {
	return r.Method == http.MethodGet && r.URL.Path == "/ping"
}

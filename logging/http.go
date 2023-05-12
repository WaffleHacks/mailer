package logging

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Request is an HTTP middleware to log requests and responses
func Request(logger *zap.Logger) func(next http.Handler) http.Handler {
	var f middleware.LogFormatter = &requestLogger{logger}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                        if r.Method == http.MethodGet && r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
                        }

			entry := f.NewLogEntry(r)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			buf := newLimitBuffer(512)
			ww.Tee(buf)

			start := time.Now()
			defer func() {
				var respBody []byte
				if ww.Status() >= 400 {
					respBody, _ = io.ReadAll(buf)
				}
				entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(start), respBody)
			}()

			next.ServeHTTP(ww, middleware.WithLogEntry(r, entry))
		})
	}
}

type requestLogger struct {
	Logger *zap.Logger
}

func (l *requestLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &RequestLoggerEntry{}
	entry.Logger = l.Logger.With(
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
		zap.String("remote", r.RemoteAddr),
		zap.String("id", middleware.GetReqID(r.Context())),
		zap.String("version", r.Proto),
	)
	entry.Logger.Info("started processing request")
	return entry
}

type RequestLoggerEntry struct {
	Logger *zap.Logger
}

func (l *RequestLoggerEntry) Write(status, bytes int, headers http.Header, elapsed time.Duration, _ interface{}) {
	latency := float64(elapsed.Nanoseconds()) / 1000000.0
	logger := l.Logger.With(zap.Int("status", status), zap.Int("bytes", bytes), zap.Float64("latency", latency))

	// Check if hijacked for websocket upgrade
	if status == 0 && bytes == 0 && len(headers) == 0 {
		logger.Info("upgraded to websocket connection")
	} else {
		logFunc := logLevelForStatus(logger, status)
		logFunc("finished processing request")
	}

}

func (l *RequestLoggerEntry) Panic(v interface{}, stack []byte) {
	stacktrace := "#"
	if !developmentOption {
		stacktrace = string(stack)
	}

	fields := []zap.Field{zap.String("stacktrace", stacktrace)}
	if developmentOption {
		fields = append(fields, zap.String("panic", fmt.Sprintf("%+v", v)))
	}

	l.Logger.Error("request handler panicked", fields...)
	if developmentOption {
		middleware.PrintPrettyStack(v)
	}
}

func logLevelForStatus(logger *zap.Logger, status int) func(msg string, fields ...zap.Field) {
	switch {
	case status <= 0:
		return logger.Warn
	case status < 400:
		return logger.Info
	case status >= 400 && status < 500:
		return logger.Warn
	case status >= 500:
		return logger.Error
	default:
		return logger.Info
	}
}

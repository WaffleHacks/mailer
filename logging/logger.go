package logging

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	developmentOption = false
	_globalL          *Logger
)

// L returns the global Logger
func L() *Logger {
	return _globalL
}

// Logger is a wrapper around a zap.Logger with support for Sentry
type Logger struct {
	core *zap.Logger
}

// New generates a new logger from a standard configuration
func New(level zap.AtomicLevel, development bool) (*Logger, error) {
	var encoding string
	if development {
		encoding = "console"
	} else {
		encoding = "json"
	}

	var sampling *zap.SamplingConfig
	if !development {
		sampling = &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		}
	}

	var encoder zapcore.EncoderConfig
	if development {
		encoder = zap.NewDevelopmentEncoderConfig()
	} else {
		encoder = zap.NewProductionEncoderConfig()
	}

	config := zap.Config{
		Level:             level,
		Development:       development,
		DisableCaller:     false,
		DisableStacktrace: true,
		Sampling:          sampling,
		Encoding:          encoding,
		EncoderConfig:     encoder,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stdout"},
	}
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	zap.ReplaceGlobals(logger)
	developmentOption = development

	_globalL = &Logger{
		core: logger,
	}

	return _globalL, nil
}

// Named is a wrapper around the Named method on zap.Logger
func (log *Logger) Named(s string) *Logger {
	if len(s) == 0 {
		return log
	}
	return &Logger{core: log.core.Named(s)}
}

// With is a wrapper around the With method on zap.Logger
func (log *Logger) With(fields ...zap.Field) *Logger {
	if len(fields) == 0 {
		return log
	}
	return &Logger{core: log.core.With(fields...)}
}

// Debug is a wrapper around the Debug method on zap.Logger
func (log *Logger) Debug(msg string, fields ...zap.Field) {
	log.core.Debug(msg, fields...)
}

// Info is a wrapper around the Info method on zap.Logger
func (log *Logger) Info(msg string, fields ...zap.Field) {
	log.core.Info(msg, fields...)
}

// Warn is a wrapper around the Warn method on zap.Logger
func (log *Logger) Warn(msg string, fields ...zap.Field) {
	log.core.Warn(msg, fields...)
}

// Error is a wrapper around the Error method on zap.Logger
func (log *Logger) Error(msg string, fields ...zap.Field) {
	log.core.Error(msg, fields...)

	log.attemptMessageCapture(msg, fields)
}

// Fatal is a wrapper around the Fatal method on zap.Logger
func (log *Logger) Fatal(msg string, fields ...zap.Field) {
	log.core.Fatal(msg, fields...)

	log.attemptMessageCapture(msg, fields)
}

// Attempt to capture the message and exception from the log
func (log *Logger) attemptMessageCapture(msg string, fields []zap.Field) {
	hub := sentry.CurrentHub()
	if hub == nil {
		return
	}

	// Capture the message
	hub.CaptureMessage(msg)

	// Try to find an error to report
	for _, field := range fields {
		if err, ok := field.Interface.(error); ok {
			hub.CaptureException(err)
		}
	}
}

// Sync is a wrapper around the Sync method on zap.Logger
func (log *Logger) Sync() error {
	return log.core.Sync()
}

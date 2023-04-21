package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/WaffleHacks/mailer/version"
)

var (
	developmentOption = false
)

// New generates a new logger from a standard configuration
func New(level zap.AtomicLevel, development bool) (*zap.Logger, error) {
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
	base, err := config.Build()
	if err != nil {
		return nil, err
	}

	logger := base.With(
		zap.String("build.version", version.Commit),
		zap.Bool("build.dirty", version.Dirty),
		zap.String("build.go", version.GoVersion),
	)
	zap.ReplaceGlobals(logger)
	developmentOption = development

	return logger, nil
}

package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	zap.ReplaceGlobals(logger)

	return logger, nil
}

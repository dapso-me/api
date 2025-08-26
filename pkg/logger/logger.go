package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	OutputPath      string
	ErrorOutputPath string
}

func New(cfg Config) (*zap.Logger, error) {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:      "timestamp",
			LevelKey:     "level",
			MessageKey:   "message",
			CallerKey:    "caller",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.FullCallerEncoder,
		},
		OutputPaths:      []string{cfg.OutputPath},
		ErrorOutputPaths: []string{cfg.ErrorOutputPath},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("logger: failed to build: %w", err)
	}

	return logger, nil
}

package utils

import (
	"fmt"

	"go.uber.org/zap"
)

// InitSugaredLogger is a helper function for initializing a *zap.SugaredLogger
func InitSugaredLogger(verbose bool) (*zap.SugaredLogger, error) {
	var zl *zap.Logger
	cfg := zap.Config{
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	if verbose {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	zl, err := cfg.Build()

	if err != nil {
		return nil, fmt.Errorf("error when initializing logger: %w", err)
	}

	sugar := zl.Sugar()
	sugar.Debug("logger initialization successful")
	return sugar, nil
}

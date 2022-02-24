package utils

import (
	"context"
	"time"
)

const (
	DefaultTimeout = "10s"
)

// InitContextWithTimeout creates a new context with a timeout. If
// timeout is invalid, it will use DefaultTimeout.
func InitContextWithTimeout(timeout string) (context.Context, context.CancelFunc, error) {
	t, err := time.ParseDuration(timeout)
	if err != nil {
		t, _ = time.ParseDuration(DefaultTimeout)
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	return ctx, cancel, err
}

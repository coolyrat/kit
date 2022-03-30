package logr

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.With("hello", "world",
		zap.Stack("stack"),
	).Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"attempt", 3,
		"backoff", time.Second,
	)
}

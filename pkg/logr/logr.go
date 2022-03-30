package logr

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config = zap.Config

func init() {
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

type logr struct {
	Logger
}

func InitLogger(conf *Config) *logr {
	if conf == nil {
		conf = &zap.Config{
			Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
			Development:      false,
			Encoding:         "console",
			EncoderConfig:    consoleEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}

	zLogr, err := conf.Build()
	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %v", err))
	}

	return &logr{zLogr.Sugar()}
}

func (l *logr) Reload(conf *Config) {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	l = InitLogger(conf)
}

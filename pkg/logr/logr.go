package logr

import (
	"fmt"
	"sync"

	"github.com/coolyrat/kit/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config = zap.Config

const (
	dataID     = "logger"
	configPath = "logger"
)

// TODO 如果一个包使用了logr，那将必定先初始化config，因为logr引用了config，也就是说，必定是使用读完配置
// 以后初始化的logr。只要config不引用logger就不会循环引用。
// 如要要用config，可以考虑在config里面创建一个子包再引用
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

	return &logr{
		SugaredLogger: zLogr.Sugar(),
	}
}

func reload() {
	var conf Config
	if err := config.Unmarshal(configPath, &conf); err != nil {
		logger.Error(err)
	}

	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	logger = InitLogger(&conf)
}

func init() {
	config.RegisterWatcher(dataID, reload)
	reload()
}

type logr struct {
	*zap.SugaredLogger
}

func (l *logr) Named(name string) Logger {
	return &logr{SugaredLogger: l.SugaredLogger.Named(name)}
}

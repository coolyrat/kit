package main

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

func main() {
	rawJSON := []byte(`level: info
encoding: json
outputPaths:
- stdout
errorOutputPaths:
- stderr
encoderConfig:
  messageKey: message
  levelKey: level
  levelEncoder: lowercase`)

	k := koanf.New(".")
	if err := k.Load(rawbytes.Provider(rawJSON), yaml.Parser()); err != nil {
		panic(err)
	}

	var cfg zap.Config
	// d := k.String("")
	// fmt.Println(d)
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{
		Tag:       "yaml",
		FlatPaths: false,
		DecoderConfig: &mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(), mapstructure.TextUnmarshallerHookFunc()),
			Metadata:         nil,
			Result:           &cfg,
			WeaklyTypedInput: true,
		},
	}); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	logger.Info("Hi, custom logger!")
	logger.Warn("Custom logger is warning you!")
	logger.Error("Let's do error instead.")
}

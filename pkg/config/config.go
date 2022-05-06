package config

import (
	"fmt"

	"github.com/coolyrat/kit/pkg/config/configcenter"
	"github.com/knadh/koanf"
	"github.com/mitchellh/mapstructure"
)

var conf = NewConfigFactory().Load()

type config struct {
	*koanf.Koanf
	configCenter configcenter.ConfigCenter
}

type File interface {
	DataID() string
	ConfigPath() string
}

type Watcher interface {
	File
	CallbackFn() func()
}

func RegisterWatcher(w Watcher) {
	load := func() {
		err := Unmarshal(w.ConfigPath(), w)
		if err != nil {
			panic(err)
		}
	}

	cb := w.CallbackFn()
	if cb == nil {
		cb = load
	}

	load()
	conf.configCenter.RegisterWatcher(w.DataID(), cb)
}

func InitConfigFile(f File) error {
	if err := Unmarshal(f.ConfigPath(), f); err != nil {
		return fmt.Errorf("init config file %s failed, %w", f.DataID(), err)
	}

	return nil
}

func MustInitConfigFile(f File) {
	if err := InitConfigFile(f); err != nil {
		panic(err)
	}
}

func Unmarshal(path string, v interface{}) error {
	return conf.UnmarshalWithConf(path, v, koanf.UnmarshalConf{
		Tag:       "yaml",
		FlatPaths: false,
		DecoderConfig: &mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(), mapstructure.TextUnmarshallerHookFunc()),
			Metadata:         nil,
			Result:           v,
			WeaklyTypedInput: true,
		},
	})
}

func Print() {
	conf.Print()
}

func GetString(path string) string {
	if s := conf.String(path); s != "" {
		return s
	}

	if s, ok := defaultConfig[path]; ok {
		return s.(string)
	}

	return ""
}

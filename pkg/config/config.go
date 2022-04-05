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

func RegisterWatcher(dataID string, cb func()) {
	conf.configCenter.RegisterWatcher(dataID, cb)
}

func Unmarshal(path string, v interface{}) error {
	fmt.Println(conf.String(path))
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

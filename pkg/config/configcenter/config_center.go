package configcenter

import (
	"fmt"

	"github.com/coolyrat/kit/pkg/koanf/providers/nacos"
	"github.com/knadh/koanf"
)

const NacosPath = "config-center.nacos"

type ConfigCenter interface {
	RegisterWatcher(dataID string, cb func())
	WatchConfig()
}

func Init(k *koanf.Koanf) ConfigCenter {
	if k.Exists(NacosPath) {
		conf := nacos.Config{Koanf: k}
		if err := k.UnmarshalWithConf(NacosPath, &conf, koanf.UnmarshalConf{Tag: "yaml"}); err != nil {
			panic(fmt.Errorf("failed to unmarshal nacos config: %w", err))
		}

		p, err := nacos.Provider(&conf)
		if err != nil {
			panic(fmt.Errorf("failed to create nacos provider: %w", err))
		}
		if err := k.Load(p, nil); err != nil {
			panic(fmt.Errorf("failed to load configs from nacos: %w", err))
		}
		p.WatchConfig()
		return &configCenter{p}
	}

	return &configCenter{&noop{}}
}

type configCenter struct {
	ConfigCenter
}

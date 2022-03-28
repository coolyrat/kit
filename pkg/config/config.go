package config

import (
	"fmt"

	"github.com/coolyrat/kit/pkg/koanf/providers/nacos"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
)

// var Config = NewConfigFactory().Load()

type config struct {
	*koanf.Koanf
	watchers map[string]*WatchedConfig
	changes  chan *nacos.Changes
}

func (c *config) WatchChange() {
	go func() {
		fmt.Println("start watching")
		for change := range c.changes {
			fmt.Println("WatchChange", change)
			// group, dataId, koanf, configPath
			k := koanf.New(".")
			k.Load(rawbytes.Provider([]byte(change.Data)), yaml.Parser())
			fmt.Println(k.Keys())
			fmt.Println(k.KeyMap())
			c.Koanf.Merge(k)
			c.Koanf.Print()
		}
		fmt.Println("stop watching")
	}()
}

func (c *config) RegisterWatcher(key string, fn func()) {
	// c.watchers[key] = append(c.watchers[key], fn)
}

type WatchedConfig struct {
	*koanf.Koanf
	key string
	fn  func()
}

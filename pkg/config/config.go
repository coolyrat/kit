package config

import (
	"github.com/coolyrat/kit/pkg/koanf/providers/nacos"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
)

var Config = NewConfigFactory().Load()

type config struct {
	*koanf.Koanf
	watchers watchers
	changes  chan *nacos.Changes
}

func (c *config) WatchChange() {
	go func() {
		for change := range c.changes {
			// group, dataId, koanf, configPath
			k := koanf.New(".")
			k.Load(rawbytes.Provider([]byte(change.Data)), yaml.Parser())
			c.Koanf.Merge(k)
			c.NotifyWatchers(change.DataID)
			c.Koanf.Print()
		}
	}()
}

func (c *config) RegisterWatcher(dataID string, cb func()) {
	if s, ok := c.watchers[dataID]; ok {
		c.watchers[dataID] = append(s, cb)
	} else {
		c.watchers[dataID] = []func(){cb}
	}
}

func (c *config) NotifyWatchers(dataID string) {
	if s, ok := c.watchers[dataID]; ok {
		for _, cb := range s {
			cb()
		}
	}
}

// watchers contains a list of callbacks for a DataID
type watchers map[string][]func()

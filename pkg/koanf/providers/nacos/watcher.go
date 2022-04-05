package nacos

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
)

// watchers contains a list of callbacks for a DataID
type watchers map[string][]func()

func (n *Nacos) RegisterWatcher(dataID string, cb func()) {
	if s, ok := n.watchers[dataID]; ok {
		n.watchers[dataID] = append(s, cb)
	} else {
		n.watchers[dataID] = []func(){cb}
	}
}

func (n *Nacos) NotifyWatchers(dataID string) {
	if s, ok := n.watchers[dataID]; ok {
		for _, cb := range s {
			cb()
		}
	}
}

func (n *Nacos) WatchConfig() {
	go func() {
		for change := range n.changes {
			k := koanf.New(".")
			// TODO add error handling
			k.Load(rawbytes.Provider([]byte(change.Data)), yaml.Parser())
			n.koanf.Merge(k)
			n.NotifyWatchers(change.DataID)
		}
	}()
}

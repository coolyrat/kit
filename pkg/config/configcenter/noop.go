package configcenter

type noop struct {
}

func (n *noop) RegisterWatcher(dataID string, cb func()) {
}

func (n *noop) WatchConfig() {
}

package svca

import (
	"github.com/coolyrat/kit/pkg/config"
	"github.com/coolyrat/kit/pkg/logr"
)

type service struct {
	config *Config
}

func NewService() *service {
	logr.Info("loading svc a")
	var c Config
	loadConfig := func() {
		config.Print()
		err := config.Unmarshal(configPath, &c)
		if err != nil {
			panic(err)
		}
	}
	loadConfig()
	config.RegisterWatcher(dataID, loadConfig)
	return &service{config: &c}
}

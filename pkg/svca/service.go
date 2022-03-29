package svca

import "github.com/coolyrat/kit/pkg/config"

type service struct {
	config *Config
}

func NewService() *service {
	var c Config
	loadConfig := func() {
		config.Config.Koanf.Print()
		err := config.Config.Unmarshal(configPath, &c)
		if err != nil {
			panic(err)
		}
	}
	loadConfig()
	config.Config.RegisterWatcher(dataID, loadConfig)
	return &service{config: &c}
}

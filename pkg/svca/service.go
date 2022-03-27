package svca

type service struct {
	config *Config
}

// func NewService() *service {
// 	var c Config
// 	err := config.Config.Unmarshal(configPath, &c)
// 	if err != nil {
// 		panic(err)
// 	}
// 	config.Config.RegisterWatcher(configPath, func(v interface{}) {
// 		err := config.Config.Unmarshal(configPath, &c)
// 		if err != nil {
// 			panic(err)
// 		}
// 	})
// 	return &service{config: &c}
// }

package svca

var configPath = "service.a"
var dataID = "svc_a"

type Config struct {
	Enable bool
	Age    int
}

func (c *Config) ConfigPath() string {
	return configPath
}

func (c *Config) Reload(v interface{}) {
	if newV, ok := v.(*Config); ok {
		c = newV
	}
}

func (c *Config) Watch() {

}

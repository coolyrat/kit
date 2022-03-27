package config

import (
	"os"
	"testing"
)

func Test_configFactory_Load(t *testing.T) {
	os.Setenv("KIT_CONFIG_FILE", "application.yml")
	os.Setenv("KIT_APPLICATION_ENV", "prod")
	NewConfigFactory().Load()
	select {}
}

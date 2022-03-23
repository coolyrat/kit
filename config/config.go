package config

import (
	"os"

	"github.com/coolyrat/kit/constant"
)

type configFactory struct {
	configFileEnv string
	appEnv        string
}

func (cf *configFactory) getConfigFile() string {
	if f := os.Getenv(constant.ConfigFileEnv); f != "" {
		return f
	}

	if cf.appEnv == "" {
		return constant.DefaultConfigFile
	}

	return ""
}

func loadConfigFile() {
	// configFactory := configFactory{
	// 	configFileEnv: os.Getenv(constant.ConfigFileEnv),
	// 	appEnv:        os.Getenv(constant.AppEnv),
	// }
	//
	// configFile := constant.DefaultConfigFile
	// if f := os.Getenv(constant.ConfigFileEnv); f != "" {
	// 	configFile = f
	// }
}

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/coolyrat/kit/constant"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

var EnvPrefix = "KIT_"

func GetEnv(key string) string {
	return os.Getenv(EnvPrefix + strings.ToUpper(key))
}

type configFactory struct {
	configFileEnv string
	appEnv        string
}

func NewConfigFactory() *configFactory {
	return &configFactory{
		configFileEnv: GetEnv(constant.ConfigFileEnv),
		appEnv:        GetEnv(constant.AppEnv),
	}
}

func (cf *configFactory) Load() {
	var k = koanf.New(".")

	confFile := cf.getConfigFile()
	if err := k.Load(file.Provider(confFile), yaml.Parser()); err != nil {
		panic(fmt.Errorf("failed to load config file %s: %w", confFile, err))
	}

	k.Load(env.Provider(EnvPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, EnvPrefix)), "_", ".", -1)
	}), nil)

	k.Print()
}

func (cf *configFactory) getConfigFile() string {
	if f := GetEnv(constant.ConfigFileEnv); f != "" {
		return f
	}

	if cf.appEnv == "" {
		return constant.DefaultConfigFile
	}

	return fmt.Sprintf("%s.%s.%s",
		filepath.Base(constant.DefaultConfigFile),
		cf.appEnv,
		filepath.Ext(constant.DefaultConfigFile))
}

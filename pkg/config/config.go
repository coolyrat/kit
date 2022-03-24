package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	constant2 "github.com/coolyrat/kit/pkg/constant"
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
		configFileEnv: GetEnv(constant2.ConfigFileEnv),
		appEnv:        GetEnv(constant2.AppEnv),
	}
}

func (cf *configFactory) Load() {
	var k = koanf.New(".")

	confFile := cf.getConfigFile()
	if err := k.Load(file.Provider(confFile), yaml.Parser()); err != nil {
		panic(fmt.Errorf("failed to load config from file %s: %w", confFile, err))
	}

	cb := func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, EnvPrefix)), "_", ".", -1)
	}
	if err := k.Load(env.Provider(EnvPrefix, ".", cb), nil); err != nil {
		panic(fmt.Errorf("failed to load config env with prefix %s: %w", EnvPrefix, err))
	}

	k.Print()
}

func (cf *configFactory) getConfigFile() string {
	if f := GetEnv(constant2.ConfigFileEnv); f != "" {
		return f
	}

	if cf.appEnv == "" {
		return constant2.DefaultConfigFile
	}

	return fmt.Sprintf("%s.%s.%s",
		filepath.Base(constant2.DefaultConfigFile),
		cf.appEnv,
		filepath.Ext(constant2.DefaultConfigFile))
}
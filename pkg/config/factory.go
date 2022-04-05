package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/coolyrat/kit/pkg/config/configcenter"
	"github.com/coolyrat/kit/pkg/config/env"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	koanfEnv "github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

type configFactory struct {
	configFileEnv string
	appEnv        string
}

func NewConfigFactory() *configFactory {
	return &configFactory{
		configFileEnv: env.GetEnv(env.ConfigFileEnv),
		appEnv:        env.GetEnv(env.AppEnv),
	}
}

func (cf *configFactory) Load() *config {
	var k = koanf.New(".")

	// Load config file
	confFile := cf.getConfigFile()
	if err := k.Load(file.Provider(confFile), yaml.Parser()); err != nil {
		panic(fmt.Errorf("failed to load config from file %s: %w", confFile, err))
	}

	// Load config from environment variables
	cb := func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, env.Prefix)), "_", ".", -1)
	}
	if err := k.Load(koanfEnv.Provider(env.Prefix, ".", cb), nil); err != nil {
		panic(fmt.Errorf("failed to load config env with prefix %s: %w", env.Prefix, err))
	}

	// Load config from config center
	configCenter := configcenter.Init(k)

	c := &config{
		configCenter: configCenter,
		Koanf:        k,
	}

	return c
}

func (cf *configFactory) getConfigFile() string {
	if cf.configFileEnv != "" {
		return cf.configFileEnv
	}

	if cf.appEnv == "" {
		return defaultConfigFile
	}

	return fmt.Sprintf("%s.%s.%s",
		filepath.Base(defaultConfigFile),
		cf.appEnv,
		filepath.Ext(defaultConfigFile))
}

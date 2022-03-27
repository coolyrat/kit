package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/coolyrat/kit/pkg/config/env"
	"github.com/coolyrat/kit/pkg/koanf/providers/nacos"
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
	ch := make(chan *nacos.Changes, 1)

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
	if k.Exists(CenterNacosPath) {
		var conf nacos.Config
		conf.Changes = ch
		if err := k.UnmarshalWithConf(CenterNacosPath, &conf, koanf.UnmarshalConf{Tag: "yaml"}); err != nil {
			panic(fmt.Errorf("failed to unmarshal nacos config: %w", err))
		}

		p, err := nacos.Provider(&conf)
		if err != nil {
			panic(fmt.Errorf("failed to create nacos provider: %w", err))
		}
		if err := k.Load(p, nil); err != nil {
			panic(fmt.Errorf("failed to load configs from nacos: %w", err))
		}
	}

	k.Print()

	c := &config{
		Koanf:    k,
		watchers: map[string]*WatchedConfig{},
		changes:  ch,
	}
	c.WatchChange()

	return c
}

func (cf *configFactory) getConfigFile() string {
	if f := env.GetEnv(env.ConfigFileEnv); f != "" {
		return f
	}

	if cf.appEnv == "" {
		return defaultConfigFile
	}

	return fmt.Sprintf("%s.%s.%s",
		filepath.Base(defaultConfigFile),
		cf.appEnv,
		filepath.Ext(defaultConfigFile))
}

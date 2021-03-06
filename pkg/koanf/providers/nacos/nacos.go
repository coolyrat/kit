package nacos

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/imdario/mergo"
	"github.com/knadh/koanf"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v3"
)

type Changes struct {
	Group, DataID, Data string
}

type Config struct {
	Koanf *koanf.Koanf

	Endpoint  string              `yaml:"endpoint" validate:"required,hostname_port"`
	Namespace string              `yaml:"namespace" validate:"required"`
	Username  string              `yaml:"username"`
	Password  string              `yaml:"password" validate:"required_with=Username"`
	Configs   map[string][]string `yaml:"configs" validate:"required"`
}

func (c *Config) build() (*Nacos, error) {
	v := validator.New()
	if err := v.Struct(c); err != nil {
		return nil, fmt.Errorf("nacos.Config invalid: %w", err)
	}

	hostPort := strings.Split(c.Endpoint, ":")
	host := hostPort[0]
	port, _ := strconv.ParseUint(hostPort[1], 10, 64)

	clientConf := *constant.NewClientConfig(
		constant.WithNamespaceId(c.Namespace),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithNotLoadCacheAtStart(true),
	)
	if c.Username != "" {
		clientConf.Username = c.Username
		clientConf.Password = c.Password
	}

	serverConf := []constant.ServerConfig{
		*constant.NewServerConfig(
			host,
			port,
		),
	}

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig":  clientConf,
		"serverConfigs": serverConf,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build nacos provider: %w", err)
	}

	return &Nacos{
		koanf:    c.Koanf,
		client:   client,
		configs:  c.Configs,
		result:   make(map[string]interface{}),
		watchers: make(watchers),
		changes:  make(chan *Changes, 1),
	}, nil
}

func Provider(config *Config) (*Nacos, error) {
	return config.build()
}

type Nacos struct {
	koanf   *koanf.Koanf
	client  config_client.IConfigClient
	configs map[string][]string
	result  map[string]interface{}

	watchers watchers
	changes  chan *Changes
}

func (n *Nacos) ReadBytes() ([]byte, error) {
	return nil, errors.New("nacos provider does not support this method")
}

func (n *Nacos) Read() (map[string]interface{}, error) {
	for group, dataIds := range n.configs {
		for _, dataId := range dataIds {
			data, err := n.client.GetConfig(vo.ConfigParam{
				Group:  group,
				DataId: dataId})
			if err != nil {
				return nil, fmt.Errorf("failed to get config from nacos group=%s dataId=%s: %w", group, dataId, err)
			}

			var result map[string]interface{}
			if err := yaml.Unmarshal([]byte(data), &result); err != nil {
				return nil, fmt.Errorf("failed to unmarshal from nacos group=%s dataId=%s: %s", group, dataId, err)
			}
			if err := mergo.Merge(&n.result, result); err != nil {
				return nil, fmt.Errorf("failed to merge config group=%s dataId=%s: %s", group, dataId, err)
			}

			err = n.client.ListenConfig(vo.ConfigParam{
				DataId: dataId,
				Group:  group,
				OnChange: func(namespace, group, dataId, data string) {
					n.changes <- &Changes{
						Group:  group,
						DataID: dataId,
						Data:   data,
					}
				},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to listen config from nacos group=%s dataId=%s: %w", group, dataId, err)
			}

		}
	}

	return n.result, nil
}

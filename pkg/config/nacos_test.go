package config

import (
	"fmt"
	"testing"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func TestNewConfigFactory(t *testing.T) {
	clientConfig := *constant.NewClientConfig(
		constant.WithUsername("kit-dev"),    // When namespace is public, fill in the blank string here.
		constant.WithPassword("kit-dev"),    // When namespace is public, fill in the blank string here.
		constant.WithNamespaceId("kit-dev"), // When namespace is public, fill in the blank string here.
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)

	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(
			"127.0.0.1",
			8848,
		),
	}

	c, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}

	content, err := c.GetConfig(vo.ConfigParam{
		DataId: "logger",
		Group:  "kit"})
	fmt.Println(content)
}

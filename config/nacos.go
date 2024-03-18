package config

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

var Client config_client.IConfigClient

func GetClient(Address string) error {
	err := ViperInit(Address)
	if err != nil {
		return err
	}

	ip := viper.GetString("Nacos.Ip")
	port := viper.GetInt("Nacos.Port")

	clientConfig := constant.ClientConfig{
		NamespaceId:         "", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: ip,
			Port:   uint64(port),
		},
	}

	// 创建动态配置客户端的另一种方式 (推荐)
	Client, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	return err
}

func GetConfig(DataId, Group string) (string, error) {

	config, err := Client.GetConfig(vo.ConfigParam{
		DataId: DataId,
		Group:  Group,
	})
	if err != nil {
		return "", err
	}

	return config, nil
}

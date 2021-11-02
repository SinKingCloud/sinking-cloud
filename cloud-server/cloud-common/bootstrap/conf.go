package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/spf13/viper"
)

// LoadConf 加载本地配置
func LoadConf(configPath string, configName string, configType string) {
	config := viper.New()
	config.AddConfigPath(configPath)
	config.SetConfigName(configName)
	config.SetConfigType(configType)
	if err := config.ReadInConfig(); err != nil {
		panic(err)
		return
	}
	util.Conf = config
}

package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/spf13/viper"
)

// LoadConf 加载本地配置
func LoadConf() {
	config := viper.New()
	config.AddConfigPath("./config")
	config.SetConfigName("application.yml")
	config.SetConfigType("yaml")
	if err := config.ReadInConfig(); err != nil {
		panic(err)
		return
	}
	util.Conf = config
}

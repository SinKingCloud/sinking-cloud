package bootstrap

import (
	config2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/spf13/viper"
)

// LoadConf 加载本地配置
func LoadConf() {
	config := viper.New()
	config.AutomaticEnv() //读取环境变量
	config.AddConfigPath("./config")
	config.SetConfigName("application.yml")
	config.SetConfigType("yaml")
	config.WatchConfig()
	if err := config.ReadInConfig(); err != nil {
		panic(err)
		return
	}
	//根据环境变量更新配置文件
	if config.GetString(config2.EnvApplicationAddr) != "" {
		config.Set(config2.ConfApplicationAddr, config.GetString(config2.EnvApplicationAddr))
	}
	if config.GetString(config2.EnvApplicationPort) != "" {
		config.Set(config2.ConfApplicationPort, config.GetString(config2.EnvApplicationPort))
	}
	if config.GetString(config2.EnvRegisterServer) != "" {
		config.Set(config2.ConfRegisterServer, config.GetString(config2.EnvRegisterServer))
	}
	if config.GetString(config2.EnvRegisterTokenKey) != "" {
		config.Set(config2.ConfRegisterTokenKey, config.GetString(config2.EnvRegisterTokenKey))
	}
	if config.GetString(config2.EnvRegisterTokenValue) != "" {
		config.Set(config2.ConfRegisterTokenValue, config.GetString(config2.EnvRegisterTokenValue))
	}
	//赋值到util
	util.Conf = config
}

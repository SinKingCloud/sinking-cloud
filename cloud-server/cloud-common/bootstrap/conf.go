package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/file"
	"github.com/spf13/viper"
	"strings"
)

// LoadConf 加载本地配置
func LoadConf() {
	path := constant.ConfPath
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	fileName := constant.ConfFile
	disk := file.NewDisk(path)
	_ = disk.AutoCreate(fileName)
	config := viper.New()
	config.AddConfigPath(path)
	config.SetConfigName(fileName)
	config.SetConfigType("yaml")
	config.WatchConfig()
	if err := config.ReadInConfig(); err != nil {
		panic(err)
		return
	}
	config.AutomaticEnv() //读取环境变量
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	keys := []string{
		constant.ServerPort,
		constant.ServerHost,
		constant.ServerMode,
		constant.ConsulAddress,
		constant.ConsulGroup,
		constant.ConsulName,
		constant.ConsulLocal,
		constant.ConsulToken,
	}
	for _, key := range keys {
		_ = config.BindEnv(key)
	}
	global.Config = config
}

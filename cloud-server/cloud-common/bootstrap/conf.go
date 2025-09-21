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
	config.AutomaticEnv() //读取环境变量
	config.AddConfigPath(path)
	config.SetConfigName(fileName)
	config.SetConfigType("yaml")
	config.WatchConfig()
	if err := config.ReadInConfig(); err != nil {
		panic(err)
		return
	}
	//赋值到util
	global.Conf = config
}

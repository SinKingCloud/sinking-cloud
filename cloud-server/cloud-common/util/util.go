package util

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/redis"
	"github.com/spf13/viper"
	"log"
)

var (
	Database *database.Database //数据库
	Conf     *viper.Viper       //配置文件
	Cache    *cache.Cache       //缓存系统
	Log      *log.Logger        //日志系统
	Redis    *redis.Client      //redis客户端
)

// IsDebug 是否debug模式
func IsDebug() bool {
	return Conf.GetString(constant.ServerMode) == "dev"
}

// ServerAddr 获取服务监听地址
func ServerAddr() (host string, port int) {
	host = Conf.GetString(constant.ServerHost)
	port = Conf.GetInt(constant.ServerPort)
	if host == "" {
		host = "0.0.0.0"
	}
	if port <= 0 {
		port = 5678
	}
	return host, port
}

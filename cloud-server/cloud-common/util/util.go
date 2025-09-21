package util

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/spf13/viper"
	"log"
)

var (
	Database *database.Database //数据库
	Conf     *viper.Viper       //文件配置
	Cache    cache.Cache
	Log      *log.Logger
)

// IsDebug 是否debug模式
func IsDebug() bool {
	return Conf.GetString(constant.ServerMode) == "dev"
}

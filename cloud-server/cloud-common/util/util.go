package util

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/logs"
	"github.com/SinKingCloud/sinking-go/sinking-sdk/sinking-sdk-go"
	"github.com/spf13/viper"
)

var (
	Database *database.Database       //数据库
	Cache    *cache.Cache             //缓存
	Log      *logs.Logs               //日志
	Conf     *viper.Viper             //文件配置
	Consul   *sinking_sdk_go.Register //服务注册发现与配置
)

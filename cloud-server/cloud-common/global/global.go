package global

import (
	"log"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/redis"
	consul "github.com/SinKingCloud/sinking-go/sinking-consul/client"
	"github.com/spf13/viper"
)

type Application struct {
	Database *database.Database //数据库
	Config   *viper.Viper       //配置文件
	Cache    cache.Interface    //缓存系统
	Log      *log.Logger        //日志系统
	Redis    *redis.Client      //redis
	Consul   *consul.Client     //consul
}

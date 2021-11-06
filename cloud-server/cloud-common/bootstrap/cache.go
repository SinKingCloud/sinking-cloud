package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// LoadCache 初始化缓存
func LoadCache() {
	conf := util.Consul.Config(group.CloudCommon).Name(config.RedisCommon).Viper()
	util.Cache = cache.NewRedis(conf.GetString("host"), conf.GetString("port"), conf.GetString("password"), conf.GetInt64("db"))
}

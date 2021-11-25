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
	util.Cache = cache.NewRedis(
		"127.0.0.1", //conf.GetString(config.RedisHost)
		conf.GetString(config.RedisPort),
		conf.GetString(config.RedisPassword),
		conf.GetInt64(config.RedisDb),
	)
}

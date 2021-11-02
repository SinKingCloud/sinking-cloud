package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// LoadCache 初始化缓存
func LoadCache() {
	conf := util.Consul.Config("default").Name("redis").Viper()
	util.Cache = cache.NewRedis(conf.GetString("host"), conf.GetString("port"), conf.GetString("pwd"), conf.GetInt64("db"))
}

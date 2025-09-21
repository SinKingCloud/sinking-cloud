package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// LoadCache 初始化缓存
func LoadCache() {
	redis := cache.NewRedis(
		util.Conf.GetString(constant.RedisHost),
		util.Conf.GetString(constant.RedisPort),
		util.Conf.GetString(constant.RedisPwd),
		util.Conf.GetInt64(constant.RedisDb),
	)
	// Redis 连接测试已在 NewRedis 中完成
	util.Cache = redis
	util.Log.Println("redis连接成功")
}

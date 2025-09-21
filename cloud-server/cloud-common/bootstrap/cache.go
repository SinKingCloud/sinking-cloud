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
	if _, err := redis.Pool.Dial(); err != nil {
		util.Log.Println("redis连接失败", err)
		panic("redis连接失败")
		return
	}
	util.Cache = redis
	util.Log.Println("redis连接成功")
}

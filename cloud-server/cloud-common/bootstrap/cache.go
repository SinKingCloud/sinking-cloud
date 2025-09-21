package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/redis"
)

// LoadCache 初始化缓存
func LoadCache() {
	client, err := redis.NewClient(
		global.Conf.GetString(constant.RedisHost),
		global.Conf.GetString(constant.RedisPort),
		global.Conf.GetString(constant.RedisPwd),
		global.Conf.GetInt64(constant.RedisDb),
	)
	if err != nil {
		global.Log.Println("redis连接失败", err)
		panic("redis连接失败: " + err.Error())
	}
	var cacheInterface cache.Interface = client
	global.Cache = &cacheInterface
	global.Redis = client
	global.Log.Println("redis连接成功")
}

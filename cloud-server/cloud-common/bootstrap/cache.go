package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/redis"
)

// LoadCache 初始化缓存
func LoadCache() {
	client, err := redis.NewClient(
		util.Conf.GetString(constant.RedisHost),
		util.Conf.GetString(constant.RedisPort),
		util.Conf.GetString(constant.RedisPwd),
		util.Conf.GetInt64(constant.RedisDb),
	)
	if err != nil {
		util.Log.Println("redis连接失败", err)
		panic("redis连接失败: " + err.Error())
	}
	var cacheInterface cache.Interface = client
	util.Cache = &cacheInterface
	util.Redis = client
	util.Log.Println("redis连接成功")
}

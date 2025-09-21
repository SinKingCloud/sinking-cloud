package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// LoadCache 初始化缓存
func LoadCache() {
	redis, err := cache.NewRedis(
		util.Conf.GetString(constant.RedisHost),
		util.Conf.GetString(constant.RedisPort),
		util.Conf.GetString(constant.RedisPwd),
		util.Conf.GetInt64(constant.RedisDb),
	)
	if err != nil {
		util.Log.Println("redis连接失败", err)
		panic("redis连接失败: " + err.Error())
	}
	util.Cache = redis
	util.Log.Println("redis连接成功")
}

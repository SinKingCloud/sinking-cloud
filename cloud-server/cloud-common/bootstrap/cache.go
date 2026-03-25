package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/redis"
)

// LoadCache 初始化缓存
func LoadCache() {
	if global.App.Consul == nil {
		panic("consul 未初始化")
		return
	}
	config := global.App.GetRedisConfig()
	if config == nil {
		panic("从consul获取redis配置失败")
		return
	}
	client, err := redis.NewClient(
		config[constant.RedisHost].(string),
		config[constant.RedisPort].(string),
		config[constant.RedisPwd].(string),
		config[constant.RedisDb].(int64),
	)
	if err != nil {
		global.App.Log.Println("redis 连接失败", err)
		panic("redis 连接失败: " + err.Error())
		return
	}
	global.App.SetCache(client)
	global.App.SetRedis(client)
	global.App.Log.Println("cache初始化成功")
}

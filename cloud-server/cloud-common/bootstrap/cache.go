package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/redis"
)

// LoadCache 初始化缓存
func LoadCache(name string) {
	if global.App.Consul == nil {
		panic("consul 未初始化")
		return
	}
	config, err := global.App.Consul.GetConfig(name)
	if err != nil || config == nil {
		panic("从consul获取配置失败: " + err.Error())
		return
	}
	client, err := redis.NewClient(
		config.GetString(constant.RedisHost),
		config.GetString(constant.RedisPort),
		config.GetString(constant.RedisPwd),
		config.GetInt64(constant.RedisDb),
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

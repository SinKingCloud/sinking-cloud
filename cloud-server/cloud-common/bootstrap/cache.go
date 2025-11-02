package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/redis"
)

// LoadCache 初始化缓存
func LoadCache() {
	if global.Consul == nil {
		panic("Consul未初始化")
		return
	}
	config, err := global.Consul.GetConfig(global.Config.GetString(constant.ConsulName))
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
		global.Log.Println("redis连接失败", err)
		panic("redis连接失败: " + err.Error())
	}
	global.Cache = client
	global.Redis = client
	global.Log.Println("redis连接成功")
}

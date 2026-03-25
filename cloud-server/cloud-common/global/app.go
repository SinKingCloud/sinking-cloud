package global

import (
	"fmt"
	"log"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/redis"
	consul "github.com/SinKingCloud/sinking-go/sinking-consul/client"
	"github.com/spf13/viper"
)

var App = &Application{}

// SetDataBase 设置数据库实例，支持链式调用
func (a *Application) SetDataBase(d *database.Database) *Application {
	a.Database = d
	return a
}

// SetConfig 设置配置文件实例，支持链式调用
func (a *Application) SetConfig(c *viper.Viper) *Application {
	a.Config = c
	return a
}

// SetCache 设置缓存系统实例，支持链式调用
func (a *Application) SetCache(c cache.Interface) *Application {
	a.Cache = c
	return a
}

// SetLog 设置日志实例，支持链式调用
func (a *Application) SetLog(l *log.Logger) *Application {
	a.Log = l
	return a
}

// SetRedis 设置Redis客户端实例，支持链式调用
func (a *Application) SetRedis(r *redis.Client) *Application {
	a.Redis = r
	return a
}

// SetConsul 设置Consul客户端实例，支持链式调用
// 注意：consul的Client实际是*api.Client（来自github.com/hashicorp/consul/api）
func (a *Application) SetConsul(c *consul.Client) *Application {
	a.Consul = c
	return a
}

// GetNatsUrl 获取nats链接地址
func (a *Application) GetNatsUrl() string {
	config, err := a.Consul.GetConfig(a.Consul.GetName())
	if err != nil {
		return ""
	}
	return fmt.Sprintf("nats://%s:%s@%s:%s",
		config.GetString(constant.NatsUser),
		config.GetString(constant.NatsPwd),
		config.GetString(constant.NatsHost),
		config.GetString(constant.NatsPort))
}

// GetRedisConfig 获取redis配置信息
func (a *Application) GetRedisConfig() map[string]interface{} {
	config, err := a.Consul.GetConfig(a.Consul.GetName())
	if err != nil {
		return nil
	}
	return map[string]interface{}{
		constant.RedisHost: config.GetString(constant.RedisHost),
		constant.RedisPort: config.GetString(constant.RedisPort),
		constant.RedisPwd:  config.GetString(constant.RedisPwd),
		constant.RedisDb:   config.GetInt64(constant.RedisDb),
	}
}

// GetPostgresConfig 获取postgres配置信息
func (a *Application) GetPostgresConfig() map[string]string {
	config, err := a.Consul.GetConfig(a.Consul.GetName())
	if err != nil {
		return nil
	}
	return map[string]string{
		constant.PostgresHost: config.GetString(constant.PostgresHost),
		constant.PostgresPort: config.GetString(constant.PostgresPort),
		constant.PostgresUser: config.GetString(constant.PostgresUser),
		constant.PostgresPwd:  config.GetString(constant.PostgresPwd),
		constant.PostgresName: config.GetString(constant.PostgresName),
	}
}

// IsDebug 是否debug模式
func (a *Application) IsDebug() bool {
	return a.Config.GetString(constant.ServerMode) == "dev"
}

// ServerAddr 获取服务监听地址
func (a *Application) ServerAddr() (host string, port int) {
	host = a.Config.GetString(constant.ServerHost)
	port = a.Config.GetInt(constant.ServerPort)
	if host == "" {
		host = "0.0.0.0"
	}
	if port <= 0 {
		port = 5678
	}
	return host, port
}

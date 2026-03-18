package global

import (
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

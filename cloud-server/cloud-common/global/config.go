package global

import (
	"fmt"
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
)

// GetNatsUrl 获取nats链接地址
func (a *Application) GetNatsUrl() string {
	config, err := a.Consul.GetConfig(a.Consul.GetName())
	if err != nil {
		return ""
	}
	user := config.GetString(constant.NatsUser)
	pwd := config.GetString(constant.NatsPwd)
	host := config.GetString(constant.NatsHost)
	port := config.GetString(constant.NatsPort)
	hostList := strings.Split(host, ",")
	var serverUrls []string
	for _, h := range hostList {
		h = strings.TrimSpace(h)
		if h != "" {
			url := fmt.Sprintf("nats://%s:%s@%s:%s", user, pwd, h, port)
			serverUrls = append(serverUrls, url)
		}
	}
	return strings.Join(serverUrls, ",")
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
		constant.PostgresMaster: config.GetString(constant.PostgresMaster),
		constant.PostgresSlave:  config.GetString(constant.PostgresSlave),
		constant.PostgresPort:   config.GetString(constant.PostgresPort),
		constant.PostgresUser:   config.GetString(constant.PostgresUser),
		constant.PostgresPwd:    config.GetString(constant.PostgresPwd),
		constant.PostgresName:   config.GetString(constant.PostgresName),
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

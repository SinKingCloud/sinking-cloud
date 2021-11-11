package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-go/sinking-sdk/sinking-sdk-go"
)

func LoadConsul() {
	//实例化一个server
	util.Consul = sinking_sdk_go.New(
		util.Conf.GetString(config.ConfRegisterServer),
		util.Conf.GetString(config.ConfRegisterTokenKey),
		util.Conf.GetString(config.ConfRegisterTokenValue),
		util.Conf.GetString(config.ConfProjectName),
		util.Conf.GetString(config.ConfProjectEnv),
	)
	//注册信息
	util.Consul.Register(
		util.Conf.GetString(config.ConfApplicationGroup),
		util.Conf.GetString(config.ConfApplicationName),
		util.Conf.GetString(config.ConfApplicationAddr)+":"+util.Conf.GetString(config.ConfApplicationPort),
	)
	//拉取服务
	list := util.Conf.Get(config.ConfApplicationUseService)
	service := map[string][]string{}
	for k, v := range list.(map[string]interface{}) {
		for _, v1 := range v.([]interface{}) {
			if service[k] == nil {
				service[k] = []string{}
			}
			service[k] = append(service[k], v1.(string))
		}
	}
	util.Consul.UseService(service)
	//监听
	util.Consul.Listen()
}

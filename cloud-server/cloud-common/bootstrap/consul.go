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
	util.Consul.UseService(util.Conf.GetStringMapString(config.ConfApplicationUseService))
	//监听
	util.Consul.Listen()
}

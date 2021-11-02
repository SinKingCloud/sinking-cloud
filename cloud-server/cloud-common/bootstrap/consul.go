package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-go/sinking-sdk/sinking-sdk-go"
)

func LoadConsul() {
	//实例化一个server
	util.Consul = sinking_sdk_go.New(
		util.Conf.GetString("register.server"),
		util.Conf.GetString("register.tokenKey"),
		util.Conf.GetString("register.tokenValue"),
		util.Conf.GetString("project.name"),
		util.Conf.GetString("project.env"),
	)
	//注册信息
	util.Consul.Register(
		util.Conf.GetString("application.group"),
		util.Conf.GetString("application.name"),
		util.Conf.GetString("application.addr")+":"+util.Conf.GetString("application.port"),
	)
	//拉取服务
	util.Consul.UseService(util.Conf.GetStringMapString("application.useService"))
	//监听
	util.Consul.Listen()
}

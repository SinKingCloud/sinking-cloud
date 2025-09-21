package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	consul "github.com/SinKingCloud/sinking-go/sinking-consul/client"
)

func LoadConsul() {
	cli := consul.NewClient(
		global.Config.GetStringSlice(constant.ConsulAddress),
		global.Config.GetString(constant.ConsulGroup),
		global.Config.GetString(constant.ConsulName),
		global.Config.GetString(constant.ConsulLocal),
		global.Config.GetString(constant.ConsulToken))
	err := cli.Connect()
	if err != nil {
		panic(err)
		return
	}
	global.Consul = cli
}

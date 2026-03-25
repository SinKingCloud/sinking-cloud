package bootstrap

import (
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	consul "github.com/SinKingCloud/sinking-go/sinking-consul/client"
)

func LoadConsul(name string) {
	var list []string
	str := global.App.Config.GetString(constant.ConsulAddress)
	if str == "" {
		list = global.App.Config.GetStringSlice(constant.ConsulAddress)
	} else {
		list = strings.Split(str, ",")
	}
	cli := consul.NewClient(
		list,
		global.App.Config.GetString(constant.ConsulGroup),
		strings.ToLower(name),
		global.App.Config.GetString(constant.ConsulLocal),
		global.App.Config.GetString(constant.ConsulToken))
	err := cli.Connect()
	if err != nil {
		global.App.Log.Println("consul连接失败", err)
		panic(err)
		return
	}
	global.App.SetConsul(cli)
	global.App.Log.Println("consul初始化成功")
}

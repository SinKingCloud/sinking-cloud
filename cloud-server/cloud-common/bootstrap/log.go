package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/logs"
)

func LoadLog() {
	conf := util.Consul.Config(group.CloudCommon).Name(config.LogStashCommon).Viper()
	util.Log = logs.NewLogstash(conf.GetString("host"), conf.GetInt("port"), conf.GetInt("timeout"))
}

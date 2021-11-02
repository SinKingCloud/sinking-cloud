package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/logs"
)

func LoadLog() {
	conf := util.Consul.Config("default").Name("logstash").Viper()
	util.Log = logs.NewLogstash(conf.GetString("host"), conf.GetInt("port"), conf.GetInt("timeout"))
}

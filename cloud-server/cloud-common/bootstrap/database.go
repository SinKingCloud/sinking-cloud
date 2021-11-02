package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// LoadDatabase 初始化数据库
func LoadDatabase() {
	conf := util.Consul.Config("default").Name("mysql").Viper()
	util.Database = database.NewMysql(
		conf.GetString("host"),
		conf.GetString("port"),
		conf.GetString("user"),
		conf.GetString("pwd"),
		conf.GetString("database"),
	)
}

package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// LoadDatabase 初始化数据库
func LoadDatabase() {
	conf := util.Consul.Config(group.CloudService).Name(config.OpenMysql).Viper()
	util.Database = database.NewMysql(
		conf.GetString(config.MysqlHost),
		conf.GetString(config.MysqlPort),
		conf.GetString(config.MysqlUser),
		conf.GetString(config.MysqlPassword),
		conf.GetString(config.MysqlDatabase),
	)
}

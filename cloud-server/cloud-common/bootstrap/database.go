package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// LoadDatabase 初始化数据库
func LoadDatabase() {
	if global.Consul == nil {
		panic("consul 未初始化")
		return
	}
	config, err := global.Consul.GetConfig(global.Config.GetString(constant.ConsulName))
	if err != nil || config == nil {
		panic("从consul获取配置失败: " + err.Error())
		return
	}
	global.Database = database.NewPostgresSql(
		config.GetString(constant.PostgresHost),
		config.GetString(constant.PostgresPort),
		config.GetString(constant.PostgresUser),
		config.GetString(constant.PostgresPwd),
		config.GetString(constant.PostgresName),
	)
	if global.Database.DbError != nil {
		global.Log.Println("数据库连接失败", global.Database.DbError)
		panic(global.Database.DbError)
		return
	}
}

package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// LoadDatabase 初始化数据库
func LoadDatabase() {
	if global.Consul == nil {
		panic("consul未初始化")
		return
	}
	config, err := global.Consul.GetConfig(global.Config.GetString(constant.ConsulName))
	if err != nil || config == nil {
		panic("从consul获取配置失败: " + err.Error())
		return
	}
	global.Database = database.NewPostgresSql(
		config.GetString(constant.MysqlHost),
		config.GetString(constant.MysqlPort),
		config.GetString(constant.MysqlUser),
		config.GetString(constant.MysqlPwd),
		config.GetString(constant.MysqlName),
	)
	if global.Database.DbError != nil {
		global.Log.Println("mysql连接失败", global.Database.DbError)
		panic("mysql连接失败")
		return
	}
	global.Log.Println("mysql连接成功")
}

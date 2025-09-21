package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// LoadDatabase 初始化数据库
func LoadDatabase() {
	global.Database = database.NewMysql(
		global.Conf.GetString(constant.MysqlHost),
		global.Conf.GetString(constant.MysqlPort),
		global.Conf.GetString(constant.MysqlUser),
		global.Conf.GetString(constant.MysqlPwd),
		global.Conf.GetString(constant.MysqlName),
	)
	if global.Database.DbError != nil {
		global.Log.Println("mysql连接失败", global.Database.DbError)
		panic("mysql连接失败")
		return
	}
	global.Log.Println("mysql连接成功")
}

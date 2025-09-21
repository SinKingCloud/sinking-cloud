package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// LoadDatabase 初始化数据库
func LoadDatabase() {
	util.Database = database.NewMysql(
		util.Conf.GetString(constant.MysqlHost),
		util.Conf.GetString(constant.MysqlPort),
		util.Conf.GetString(constant.MysqlUser),
		util.Conf.GetString(constant.MysqlPwd),
		util.Conf.GetString(constant.MysqlName),
	)
	if util.Database.DbError != nil {
		util.Log.Println("mysql连接失败", util.Database.DbError)
		panic("mysql连接失败")
		return
	}
	util.Log.Println("mysql连接成功")
}

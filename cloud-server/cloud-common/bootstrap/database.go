package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// LoadDatabase 初始化数据库
func LoadDatabase() {
	if global.App.Consul == nil {
		panic("consul 未初始化")
		return
	}
	config := global.App.GetPostgresConfig()
	if config == nil {
		panic("从consul获取postgress配置失败")
		return
	}
	global.App.SetDataBase(database.NewPostgresSql(
		config[constant.PostgresHost],
		config[constant.PostgresPort],
		config[constant.PostgresUser],
		config[constant.PostgresPwd],
		config[constant.PostgresName],
	))
	if global.App.Database.DbError != nil {
		global.App.Log.Println("数据库连接失败", global.App.Database.DbError)
		panic(global.App.Database.DbError)
		return
	}
	global.App.Log.Println("database初始化成功")
}

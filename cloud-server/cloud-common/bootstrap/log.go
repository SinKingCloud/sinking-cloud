package bootstrap

import (
	"log"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
)

func LoadLog() {
	global.App.SetLog(log.Default())
	global.App.Log.Println("log初始化成功")
}

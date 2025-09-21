package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"log"
)

func LoadLog() {
	global.Log = log.Default()
}

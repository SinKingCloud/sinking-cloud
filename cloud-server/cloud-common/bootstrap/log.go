package bootstrap

import (
	"log"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
)

func LoadLog() {
	global.Log = log.Default()
}

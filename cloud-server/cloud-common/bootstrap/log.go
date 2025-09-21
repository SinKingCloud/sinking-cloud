package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"log"
)

func LoadLog() {
	util.Log = log.Default()
}

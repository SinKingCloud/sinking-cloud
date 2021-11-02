package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/logs"
)

func LoadLog() {
	util.Log = logs.NewLogstash("", 5601, 5)
}

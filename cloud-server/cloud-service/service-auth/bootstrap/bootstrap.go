package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/bootstrap"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
)

func Load() {
	bootstrap.LoadLog()
	bootstrap.LoadConf()
	bootstrap.LoadConsul(constant.ServiceAuth)
	bootstrap.LoadCache()
}

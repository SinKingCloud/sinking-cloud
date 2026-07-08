package bootstrap

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/bootstrap"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
)

func Load() {
	bootstrap.Load(constant.ServiceSpace)
}

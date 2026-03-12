package bootstrap

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/bootstrap"

func Load() {
	bootstrap.LoadLog()
	bootstrap.LoadConf()
	bootstrap.LoadConsul()
}

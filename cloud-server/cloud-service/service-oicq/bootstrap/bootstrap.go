package bootstrap

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/bootstrap"

// Load 初始化所有功能
func Load() {
	bootstrap.LoadConf()
	bootstrap.LoadConsul()
	bootstrap.LoadLog()
	bootstrap.LoadCache()
	LoadDatabase()
}

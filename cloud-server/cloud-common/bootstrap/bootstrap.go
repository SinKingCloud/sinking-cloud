package bootstrap

// Load 初始化加载
// name 服务名称
func Load(name string) {
	LoadLog()
	LoadConf()
	LoadConsul(name)
	LoadDatabase(name)
	LoadCache(name)
}

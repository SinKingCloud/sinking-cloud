package bootstrap

func Load() {
	LoadLog()
	LoadConf()
	LoadConsul()
	LoadDatabase()
	LoadCache()
}

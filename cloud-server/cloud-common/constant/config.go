package constant

const (
	TempPath = "./temp" //缓存目录

	ConfPath = "./config"        //配置文件目录
	ConfFile = "application.yml" //配置文件

	ServerMode = "server.mode" //运行模式
	ServerHost = "server.host" //绑定地址
	ServerPort = "server.port" //绑定端口

	ConsulGroup   = "consul.group"   //服务分组
	ConsulAddress = "consul.address" //服务地址
	ConsulToken   = "consul.token"   //服务token
	ConsulLocal   = "consul.local"   //本机地址
)

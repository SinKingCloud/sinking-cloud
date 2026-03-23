package route

import (
	rpc2 "service-admin/app/http/rpc"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApp 加载app路由
func loadApp(s *sinking_web.Engine) {
	loadRpc(s)
}

// loadRpc 加载rpc
func loadRpc(s *sinking_web.Engine) {
	rpc.RegisterHandle(s)
	rpc.RegisterFunc(constant.RpcGetEnumDataFunName, rpc2.Enum.GetData) //获取枚举配置
	rpc.RegisterFunc("GetUserConfig", rpc2.Config.GetUserConfig)        //获取用户配置
	rpc.RegisterFunc("SetUserConfig", rpc2.Config.SetUserConfig)        //更新用户配置
	rpc.RegisterFunc("SetUserConfigs", rpc2.Config.SetUserConfigs)      //批量更新用户配置
	rpc.RegisterFunc("GetSystemConfig", rpc2.Config.GetSystemConfig)    //获取系统配置
	rpc.RegisterFunc("GetSystemConfigs", rpc2.Config.GetSystemConfigs)  //批量获取系统配置
	rpc.RegisterFunc("SetSystemConfig", rpc2.Config.SetSystemConfig)    //更新系统配置
	rpc.RegisterFunc("SetSystemConfigs", rpc2.Config.SetSystemConfigs)  //批量更新系统配置
	rpc.RegisterFunc("Login", rpc2.Auth.Login)                          //账号登录
	rpc.RegisterFunc("Logout", rpc2.Auth.Logout)                        //登陆注销
}

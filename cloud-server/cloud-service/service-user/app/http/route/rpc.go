package route

import (
	"service-user/app/enum"
	rpc2 "service-user/app/http/rpc"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRpc 加载rpc
func loadRpc(s *sinking_web.Engine) {
	rpc.RegisterHandle(s)
	rpc.RegisterFunc(constant.RpcGetEnumDataFunName, service2.Enum.GetDataHandler(enum.Data)) //获取枚举配置
	rpc.RegisterFunc("GetUserInfo", rpc2.User.GetUserInfo)                                    //获取用户详情
	rpc.RegisterFunc("GetUserConfig", rpc2.Config.GetUserConfig)                              //获取用户配置
	rpc.RegisterFunc("SetUserConfig", rpc2.Config.SetUserConfig)                              //更新用户配置
	rpc.RegisterFunc("SetUserConfigs", rpc2.Config.SetUserConfigs)                            //批量更新用户配置
	rpc.RegisterFunc("GenSocketToken", rpc2.Socket.GenToken)                                  //生成socket连接token
	rpc.RegisterFunc("Login", rpc2.Auth.Login)                                                //账号登录
	rpc.RegisterFunc("Logout", rpc2.Auth.Logout)                                              //登陆注销
	rpc.RegisterFunc("ClearExpireSessions", rpc2.Auth.ClearExpireSessions)                    //清理过期会话
}

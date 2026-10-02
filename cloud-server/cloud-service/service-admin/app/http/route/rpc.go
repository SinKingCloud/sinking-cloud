package route

import (
	"service-admin/app/enum"
	rpc2 "service-admin/app/http/rpc"

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
	rpc.RegisterFunc("CreateUserLog", rpc2.User.CreateUserLog)                                //创建用户操作日志
	rpc.RegisterFunc("GetUserConfig", rpc2.Config.GetUserConfig)                              //获取用户配置
	rpc.RegisterFunc("SetUserConfig", rpc2.Config.SetUserConfig)                              //更新用户配置
	rpc.RegisterFunc("SetUserConfigs", rpc2.Config.SetUserConfigs)                            //批量更新用户配置
	rpc.RegisterFunc("GetSystemConfig", rpc2.Config.GetSystemConfig)                          //获取系统配置
	rpc.RegisterFunc("GetSystemConfigs", rpc2.Config.GetSystemConfigs)                        //批量获取系统配置
	rpc.RegisterFunc("SetSystemConfig", rpc2.Config.SetSystemConfig)                          //更新系统配置
	rpc.RegisterFunc("SetSystemConfigs", rpc2.Config.SetSystemConfigs)                        //批量更新系统配置
	rpc.RegisterFunc("Login", rpc2.Auth.Login)                                                //账号登录
	rpc.RegisterFunc("Logout", rpc2.Auth.Logout)                                              //登陆注销
	rpc.RegisterFunc("ClearExpireSessions", rpc2.Auth.ClearExpireSessions)                    //清理过期会话
	rpc.RegisterFunc("RegisterPermissions", rpc2.Rbac.RegisterPermissions)                    //注册rbac权限
}

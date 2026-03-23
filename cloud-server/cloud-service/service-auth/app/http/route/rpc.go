package route

import (
	rpc2 "service-auth/app/http/rpc"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRpc 加载rpc
func loadRpc(s *sinking_web.Engine) {
	rpc.RegisterHandle(s)
	rpc.RegisterFunc("CheckCaptcha", rpc2.Auth.CheckCaptcha) //判断滑块验证码
	rpc.RegisterFunc("CheckSms", rpc2.Auth.CheckSms)         //判断短信验证码
	rpc.RegisterFunc("CheckMail", rpc2.Auth.CheckMail)       //判断邮件验证码
}

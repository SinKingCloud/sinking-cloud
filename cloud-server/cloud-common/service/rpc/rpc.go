package rpc

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/captcha"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/face"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/oauth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/sms"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/system_rbac"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/system_user"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/upload"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/user_app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/web_user"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/web_user_pay_log"
)

var (
	Face          = face.GetIns()             //人脸识别
	Log           = log.GetIns()              //用户日志
	Oauth         = oauth.GetIns()            //授权登录
	WebUser       = web_user.GetIns()         //网站用户
	WebUserPayLog = web_user_pay_log.GetIns() //余额日志
	SystemRbac    = system_rbac.GetIns()      //权限判断
	SystemUser    = system_user.GetIns()      //系统用户
	Upload        = upload.GetIns()           //上传文件
	UserApp       = user_app.GetIns()         //用户应用
	Sms           = sms.GetIns()              //短信服务
	Captcha       = captcha.GetIns()          //图形验证
)

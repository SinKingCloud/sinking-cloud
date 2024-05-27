package service

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_access"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_flie"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_log_request"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_menu"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_rbac"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_role"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_user"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/web_notice"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/web_user"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/web_user_order"
)

var (
	SystemAccess     = system_access.GetIns()      //系统权限
	SystemConfig     = system_config.GetIns()      //系统配置
	SystemFile       = system_file.GetIns()        //系统文件
	SystemLog        = system_log.GetIns()         //系统日志
	SystemLogRequest = system_log_request.GetIns() //请求日志
	SystemMenu       = system_menu.GetIns()        //系统菜单
	SystemRbac       = system_rbac.GetIns()        //权限信息
	SystemRole       = system_role.GetIns()        //系统角色
	SystemUser       = system_user.GetIns()        //系统用户
	WebNotice        = web_notice.GetIns()         //网站公告
	WebUser          = web_user.GetIns()           //网站用户
	WebUserOrder     = web_user_order.GetIns()     //网站用户订单
)

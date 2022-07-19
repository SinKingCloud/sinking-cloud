package auth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth/auth_app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth/auth_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth/auth_version"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth/user_auth"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"strconv"
)

// Check 授权检测
func (ControllerAuth) Check(c *sinking_server.Context) {
	type Form struct {
		Domain string `form:"domain" default:"" validate:"omitempty,max=50" label:"域名"`
		Code   string `form:"code" default:"" validate:"required,len=36" label:"授权码"`
		Ip     string `form:"ip" default:"" validate:"omitempty,ip" label:"绑定IP"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	if form.Domain != "" {
		where["domain"] = form.Domain
	}
	if form.Code != "" {
		where["code"] = form.Code
	}
	if form.Ip != "" {
		where["ip"] = form.Ip
	}
	where["user_id"] = strconv.Itoa(c.GetWebUserInfo().Id)
	userAuth, e := user_auth.GetIns().Find(where)
	if e != nil || userAuth == nil {
		c.Error("获取授权许可证失败，请前往官网下载正版产品")
		return
	}
	if userAuth.Status != user_auth.Normal {
		c.Error("您的许可证已被封禁,请前往官网查看详情")
		return
	}
	authApp, e := auth_app.GetIns().FindByIdWithCache(userAuth.AuthId, false)
	if e != nil || authApp == nil {
		c.Error("获取应用信息失败,请前往官网查看详情")
		return
	}
	ver := 0
	authVersion, e := auth_version.GetIns().FindByIdWithCache(authApp.VersionId, false)
	if e == nil && authVersion != nil {
		ver = authVersion.VersionNumber
	}
	c.SuccessWithData("检测成功", sinking_web.H{
		"name":                  authApp.Name,     //程序名称
		"describe":              authApp.Describe, //程序描述
		"price":                 authApp.Price,    //程序价格
		"latest_version_number": ver,              //最新版本号
	})
}

// Log 授权记录
func (ControllerAuth) Log(c *sinking_server.Context) {
	type Form struct {
		AuthId  int    `json:"auth_id" default:"" validate:"omitempty,numeric" label:"程序ID"`
		Content string `json:"content" validate:"required" label:"日志内容"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	err := auth_log.GetIns().Create(&model.AuthLog{
		Ip:      c.GetRequestIp(),
		AuthId:  form.AuthId,
		Content: form.Content,
	})
	if err == nil {
		c.Success("上报日志成功")
	} else {
		c.Error("上报日志失败")
	}
}

// Upgrade 获取升级信息
func (ControllerAuth) Upgrade(c *sinking_server.Context) {
	type Form struct {
		Code          string `form:"code" default:"" validate:"required,len=36" label:"授权码"`
		VersionNumber int    `form:"version_number" default:"" validate:"min=0" label:"当前版本号"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	userAuth, e := user_auth.GetIns().Find(map[string]string{"code": form.Code})
	if e != nil || userAuth == nil {
		c.Error("获取授权许可证失败，请前往官网下载正版产品")
		return
	}
	if userAuth.Status != user_auth.Normal {
		c.Error("您的许可证已被封禁,请前往官网查看详情")
		return
	}
	list, err := auth_version.GetIns().SelectUpgrade(userAuth.AuthId, form.VersionNumber)
	if err != nil || list == nil {
		c.Error("获取更新信息失败,请前往官网查看详情")
		return
	}
	c.SuccessWithData("获取更新信息成功", list)
}

// Package 获取程序下载地址
func (ControllerAuth) Package(c *sinking_server.Context) {
	type Form struct {
		Code          string `form:"code" default:"" validate:"required,len=36" label:"授权码"`
		VersionNumber int    `form:"version_number" default:"" validate:"min=0" label:"当前版本号"`
		Types         string `form:"type" default:"upgrade" validate:"required,oneof=upgrade release" label:"包类型"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	userAuth, e := user_auth.GetIns().Find(map[string]string{"code": form.Code})
	if e != nil || userAuth == nil {
		c.Error("获取授权许可证失败，请前往官网下载正版产品")
		return
	}
	if userAuth.Status != user_auth.Normal {
		c.Error("您的许可证已被封禁,请前往官网查看详情")
		return
	}
	authVersion, err := auth_version.GetIns().FindByAuthIdVerNum(userAuth.AuthId, form.VersionNumber)
	if err != nil || authVersion == nil {
		c.Error("获取文件失败,请前往官网查看详情")
		return
	}
	var url string
	var err2 error
	if form.Types == "release" {
		url, err2 = web_auth.Package.BuildRelease(authVersion.Id, userAuth)
	} else {
		url, err2 = web_auth.Package.BuildUpdate(authVersion.Id, userAuth)
	}
	if err2 != nil {
		c.Error(err2.Error())
		return
	}
	c.SuccessWithData("获取文件成功", url)
}

package auth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth/user_auth"
	"strconv"
)

// List 授权列表
func (ControllerAuth) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `json:"order_by_field" default:"id" validate:"oneof=id status create_time" label:"排序字段"`
		OrderByType     string `json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		AuthId          string `json:"auth_id" default:"" validate:"omitempty,numeric" label:"程序ID"`
		Domain          string `json:"domain" default:"" validate:"omitempty,max=50" label:"域名"`
		Code            string `json:"code" default:"" validate:"omitempty,len=36" label:"授权码"`
		Ip              string `json:"ip" default:"" validate:"omitempty,ip" label:"绑定IP"`
		Status          string `json:"status" default:"" validate:"omitempty,oneof=0 1" label:"状态"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"修改起始时间"`
		UpdateTimeEnd   string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"修改结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	where["user_id"] = strconv.Itoa(c.GetWebUserInfo().Id)
	if form.Domain != "" {
		where["domain"] = form.Domain
	}
	if form.Code != "" {
		where["code"] = form.Code
	}
	if form.Ip != "" {
		where["ip"] = form.Ip
	}
	if form.Status != "" {
		where["status"] = form.Status
	}
	if form.CreateTimeStart != "" {
		where["create_time_start"] = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where["create_time_end"] = form.CreateTimeEnd
	}
	if form.UpdateTimeStart != "" {
		where["update_time_start"] = form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where["update_time_end"] = form.UpdateTimeEnd
	}
	data, total, err := service.UserAuth.Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取列表失败")
	} else {
		rpc.Log.CreateUser(c.GetWebUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), log.UserEventShow, "查看授权列表", "查看程序授权列表")
		c.SuccessWithData("获取列表成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

// Download 获取程序下载地址
func (ControllerAuth) Download(c *sinking_server.Context) {
	type Form struct {
		Id int `json:"id" default:"" validate:"min=0" label:"用户程序ID"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	userAuth, e := service.UserAuth.FindById(form.Id)
	if e != nil || userAuth == nil || userAuth.UserId != c.GetWebUserInfo().Id {
		c.Error("获取授权许可证失败")
		return
	}
	if userAuth.Status != user_auth.Normal {
		c.Error("您的许可证已被封禁")
		return
	}
	authApp, e := service.AuthApp.FindByIdWithCache(userAuth.AuthId, false)
	if e != nil || authApp == nil {
		c.Error("获取应用信息失败")
		return
	}
	authVersion, err := service.AuthVersion.FindByIdWithCache(authApp.VersionId, false)
	if err != nil || authVersion == nil {
		c.Error("获取版本信息失败")
		return
	}
	url, err2 := service.WebAuthPackage.BuildRelease(authVersion.Id, userAuth)
	if err2 != nil {
		c.Error(err2.Error())
		return
	}
	c.SuccessWithData("获取成功", url)
}

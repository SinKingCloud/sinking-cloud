package uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/uin"
	"strconv"
)

type ControllerUin struct {
}

func (ControllerUin) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		AppId           string `form:"app_id" default:"" validate:"omitempty,numeric" label:"应用ID"`
		ServerKey       string `form:"server_key" default:"" validate:"omitempty,max=50" label:"地区标识"`
		UinStatus       string `form:"uin_status" default:"" validate:"omitempty,numeric" label:"挂机状态"`
		Uin             string `form:"uin" default:"" validate:"omitempty,numeric" label:"QQ账号"`
		Status          string `form:"status" default:"" validate:"omitempty,oneof=0 1" label:"QQ状态"`
		RunTimeStart    string `form:"run_time_start" default:"" validate:"-" label:"挂机起始时间"`
		RunTimeEnd      string `form:"run_time_end" default:"" validate:"-" label:"挂机结束时间"`
		CreateTimeStart string `form:"create_time_start" default:"" validate:"-" label:"创建起始时间"`
		CreateTimeEnd   string `form:"create_time_end" default:"" validate:"-" label:"创建结束时间"`
		UpdateTimeStart string `form:"update_time_start" default:"" validate:"-" label:"更新起始时间"`
		UpdateTimeEnd   string `form:"update_time_end" default:"" validate:"-" label:"更新结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	where["user_id"] = strconv.Itoa(c.GetWebUserInfo().Id)
	if form.AppId != "" {
		where["app_id"] = form.AppId
	}
	if form.ServerKey != "" {
		where["server_key"] = form.ServerKey
	}
	if form.UinStatus != "" {
		where["uin_status"] = form.UinStatus
	}
	if form.Uin != "" {
		where["uin"] = form.Uin
	}
	if form.Status != "" {
		where["status"] = form.Status
	}
	if form.RunTimeStart != "" {
		where["run_time_start"] = form.RunTimeStart
	}
	if form.RunTimeEnd != "" {
		where["run_time_end"] = form.RunTimeEnd
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
	data, total, err := uin.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取挂机列表失败")
	} else {
		c.SuccessWithData("获取挂机列表成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

// Self 账户信息
func (ControllerUin) Self(c *sinking_server.Context) {
	c.Success("success")
}

// Count 统计信息
func (ControllerUin) Count(c *sinking_server.Context) {
	c.Success("success")
}

// Buy 下单
func (ControllerUin) Buy(c *sinking_server.Context) {
	c.Success("success")
}

// Login 登陆QQ
func (ControllerUin) Login(c *sinking_server.Context) {
	c.Success("success")
}

// Setting 功能设置
func (ControllerUin) Setting(c *sinking_server.Context) {
	c.Success("success")
}

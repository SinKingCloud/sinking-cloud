package uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin/order"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin/uin"
	"strconv"
)

type ControllerOrder struct {
}

func (ControllerOrder) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `json:"order_by_field" default:"id" validate:"oneof=id" label:"排序字段"`
		OrderByType     string `json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		AppId           string `json:"app_id" default:"" validate:"omitempty,numeric" label:"应用ID"`
		Uin             string `json:"uin" default:"" json:"uin" validate:"required,numeric,min=5,max=10" label:"QQ账号"`
		Type            string `json:"type" default:"" validate:"omitempty,oneof=0 1 2" label:"订单类型"`
		Name            string `json:"name" default:"" validate:"omitempty,max=50" label:"标题"`
		TradeNo         string `json:"trade_no" default:"" validate:"omitempty,numeric" label:"订单号"`
		Status          string `json:"status" default:"" validate:"omitempty,oneof=0 1" label:"状态"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	where := make(map[string]string)
	where["user_id"] = strconv.Itoa(c.GetWebUserInfo().Id)
	if form.AppId != "" {
		where["app_id"] = form.AppId
	}
	where["uin_id"] = strconv.Itoa(uinInfo.Id)
	if form.Status != "" {
		where["status"] = form.Status
	}
	if form.TradeNo != "" {
		where["trade_no"] = form.TradeNo
	}
	if form.Name != "" {
		where["name"] = form.Name
	}
	if form.CreateTimeStart != "" {
		where["create_time_start"] = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where["create_time_end"] = form.CreateTimeEnd
	}
	data, total, err := order.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取订单列表失败")
	} else {
		c.SuccessWithData("获取订单列表成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

func (ControllerOrder) Refund(c *sinking_server.Context) {
	type Form struct {
		TradeNo string `json:"trade_no" default:"" validate:"required,numeric" label:"订单号"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	orderInfo, err := order.GetIns().FindByUserIdTradeNo(c.GetWebUserInfo().Id, form.TradeNo)
	if err != nil || orderInfo == nil {
		c.Error("未查询到有效订单")
		return
	}
	err = web_uin.Refund(form.TradeNo)
	if err == nil {
		c.Success("退单成功")
	} else {
		c.Error(err.Error())
	}
}

package web

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/system_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/web_user_order"
)

// List 列表
func (ControllerUserOrder) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id user_id pay_type order_type create_time update_time status" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Status          string `form:"status" default:"" validate:"omitempty,oneof=0 1" label:"状态"`
		Account         string `form:"account" default:"" validate:"omitempty,max=15" label:"用户账号"`
		NickName        string `form:"nick_name" default:"" validate:"omitempty,max=20" label:"用户昵称"`
		Email           string `form:"email" default:"" validate:"omitempty" label:"用户邮箱"`
		Phone           string `form:"phone" default:"" validate:"omitempty,numeric" label:"用户手机"`
		UserId          string `form:"user_id" default:"" validate:"omitempty,numeric" label:"用户ID"`
		PayType         string `form:"pay_type" default:"" validate:"omitempty,numeric" label:"支付类型"`
		OrderType       string `form:"order_type" default:"" validate:"omitempty,numeric" label:"订单类型"`
		TradeNo         string `form:"trade_no" default:"" validate:"omitempty,numeric" label:"订单号"`
		OutTradeNo      string `form:"out_trade_no" default:"" validate:"omitempty,numeric" label:"外部订单号"`
		Name            string `form:"name" default:"" validate:"omitempty,max=50" label:"订单名称"`
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
	if form.UserId != "" {
		where["user_id"] = form.UserId
	}
	if form.Status != "" {
		where["status"] = form.Status
	}
	if form.Account != "" {
		where["account"] = form.Account
	}
	if form.NickName != "" {
		where["nick_name"] = form.NickName
	}
	if form.Email != "" {
		where["email"] = form.Email
	}
	if form.Phone != "" {
		where["phone"] = form.Phone
	}
	if form.OrderType != "" {
		where["order_type"] = form.OrderType
	}
	if form.PayType != "" {
		where["pay_type"] = form.PayType
	}
	if form.TradeNo != "" {
		where["trade_no"] = form.TradeNo
	}
	if form.OutTradeNo != "" {
		where["out_trade_no"] = form.OutTradeNo
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
	if form.UpdateTimeStart != "" {
		where["update_time_start"] = form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where["update_time_end"] = form.UpdateTimeEnd
	}
	data, total, err := web_user_order.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err == nil {
		system_log.GetIns().CreateByQueue(c, 0, "查看订单", "查看订单列表")
		c.SuccessWithData("获取订单数据成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
		return
	}
	c.Error("获取数据失败")
}

// Delete 删除记录
func (ControllerUserOrder) Delete(c *sinking_server.Context) {
	type Form struct {
		CreateTimeStart string `json:"create_time_start" default:"" validate:"required" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"required" label:"创建结束时间"`
		Status          string `form:"status" default:"" validate:"omitempty,oneof=0 1" label:"订单状态"`
		UserId          string `form:"user_id" default:"" validate:"omitempty,numeric" label:"用户ID"`
		PayType         string `form:"pay_type" default:"" validate:"omitempty,numeric" label:"支付类型"`
		OrderType       string `form:"order_type" default:"" validate:"omitempty,numeric" label:"订单类型"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	if form.UserId != "" {
		where["user_id"] = form.UserId
	}
	if form.PayType != "" {
		where["pay_type"] = form.PayType
	}
	if form.OrderType != "" {
		where["order_type"] = form.OrderType
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
	err := web_user_order.GetIns().Delete(where)
	if err == nil {
		system_log.GetIns().CreateByQueue(c, 2, "清理订单", "清理用户订单")
		c.Success("删除订单成功")
	} else {
		c.Error("删除订单失败")
	}
}

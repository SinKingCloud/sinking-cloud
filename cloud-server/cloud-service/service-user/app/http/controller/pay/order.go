package pay

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/order"
	"strconv"
)

// OrderList 订单记录
func (ControllerPay) OrderList(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id status create_time" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		PayType         string `form:"pay_type" json:"pay_type" default:"" validate:"omitempty,numeric" label:"支付类型"`
		OrderType       string `form:"order_type" json:"order_type" default:"" validate:"omitempty,numeric" label:"订单类型"`
		TradeNo         string `form:"trade_no" json:"trade_no" default:"" validate:"omitempty,max=100" label:"系统订单号"`
		OutTradeNo      string `form:"out_trade_no" json:"out_trade_no" default:"" validate:"omitempty,max=100" label:"外部订单号"`
		Status          string `form:"status" json:"status" default:"" validate:"omitempty,oneof=0 1 2" label:"订单状态"`
		CreateTimeStart string `form:"create_time_start" json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `form:"create_time_end" json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `form:"update_time_start" json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `form:"update_time_end" json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	where["user_id"] = strconv.Itoa(c.GetWebUserInfo().Id)
	if form.PayType != "" {
		where["pay_type"] = form.PayType
	}
	if form.OrderType != "" {
		where["order_type"] = form.OrderType
	}
	if form.TradeNo != "" {
		where["trade_no"] = form.TradeNo
	}
	if form.OutTradeNo != "" {
		where["out_trade_no"] = form.OutTradeNo
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
	data, total, err := order.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取订单记录失败")
	} else {
		log.GetIns().CreateByQueue(c, log.EventShow, "查看订单", "查看订单记录")
		c.SuccessWithData("获取订单记录成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

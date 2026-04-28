package admin

import (
	"service-pay/app/repository/order"
	"service-pay/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

type orderController struct {
}

// List 订单列表
func (orderController) List(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("id", "desc", "id,user_id,channel_id,type,status,pay_status,pay_time,expired_time,complete_time,create_time,update_time")
	var form struct {
		UserId            string `json:"user_id" default:"" validate:"omitempty,numeric" label:"用户ID"`
		ChannelId         string `json:"channel_id" default:"" validate:"omitempty,numeric" label:"通道ID"`
		TradeNo           string `json:"trade_no" default:"" validate:"omitempty" label:"订单号"`
		Type              string `json:"type" default:"" validate:"omitempty,numeric" label:"订单类型"`
		Status            string `json:"status" default:"" validate:"omitempty,numeric" label:"订单状态"`
		OutTradeNo        string `json:"out_trade_no" default:"" validate:"omitempty" label:"外部订单号"`
		PayStatus         string `json:"pay_status" default:"" validate:"omitempty,numeric" label:"支付状态"`
		PayTimeStart      string `json:"pay_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"支付起始时间"`
		PayTimeEnd        string `json:"pay_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"支付结束时间"`
		CompleteTimeStart string `json:"complete_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"完成起始时间"`
		CompleteTimeEnd   string `json:"complete_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"完成结束时间"`
		CreateTimeStart   string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd     string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart   string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd     string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &order.SelectOrder{}
	if form.UserId != "" {
		where.UserId = form.UserId
	}
	if form.ChannelId != "" {
		where.ChannelId = form.ChannelId
	}
	if form.TradeNo != "" {
		where.TradeNo = form.TradeNo
	}
	if form.Type != "" {
		where.Type = form.Type
	}
	if form.Status != "" {
		where.Status = form.Status
	}
	if form.OutTradeNo != "" {
		where.OutTradeNo = form.OutTradeNo
	}
	if form.PayStatus != "" {
		where.PayStatus = form.PayStatus
	}
	if form.PayTimeStart != "" {
		where.PayTimeStart = form.PayTimeStart
	}
	if form.PayTimeEnd != "" {
		where.PayTimeEnd = form.PayTimeEnd
	}
	if form.CompleteTimeStart != "" {
		where.CompleteTimeStart = form.CompleteTimeStart
	}
	if form.CompleteTimeEnd != "" {
		where.CompleteTimeEnd = form.CompleteTimeEnd
	}
	if form.CreateTimeStart != "" {
		where.CreateTimeStart = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where.CreateTimeEnd = form.CreateTimeEnd
	}
	if form.UpdateTimeStart != "" {
		where.UpdateTimeStart = form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where.UpdateTimeEnd = form.UpdateTimeEnd
	}
	data, total, err := service.Order.Select(where, orderByField, orderByType, pageNum, pageSize)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", page.NewPage(total, pageNum, pageSize, data))
}

// Info 订单详情
func (orderController) Info(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	info, err := service.Order.FindById(form.Id)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", info)
}

// Deliver 订单发货
func (orderController) Deliver(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Order.Deliver(form.Id, c.GetUserInfo().Id)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("订单发货成功")
}

// Refound 订单退款
func (orderController) Refound(c *context.Context) {
	var form struct {
		Id    int64   `json:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
		Money float64 `json:"money" default:"0" validate:"omitempty,min=0" label:"退款金额"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Order.Refound(form.Id, form.Money, c.GetUserInfo().Id)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("订单退款成功")
}

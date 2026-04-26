package admin

import (
	"service-pay/app/repository/order"
	"service-pay/app/repository/order_log"
	"service-pay/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

type orderController struct {
}

// List 订单列表
func (orderController) List(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("id", "desc", "id,user_id,channel_id,type,status,pay_status,pay_type,pay_time,expired_time,complete_time,create_time,update_time")
	var form struct {
		UserId            string `json:"user_id" default:"" validate:"omitempty,numeric" label:"用户ID"`
		ChannelId         string `json:"channel_id" default:"" validate:"omitempty,numeric" label:"通道ID"`
		TradeNo           string `json:"trade_no" default:"" validate:"omitempty" label:"订单号"`
		Type              string `json:"type" default:"" validate:"omitempty,numeric" label:"订单类型"`
		Status            string `json:"status" default:"" validate:"omitempty,numeric" label:"订单状态"`
		OutTradeNo        string `json:"out_trade_no" default:"" validate:"omitempty" label:"外部订单号"`
		PayStatus         string `json:"pay_status" default:"" validate:"omitempty,numeric" label:"支付状态"`
		PayType           string `json:"pay_type" default:"" validate:"omitempty,numeric" label:"支付类型"`
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
	if form.PayType != "" {
		where.PayType = form.PayType
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

// Log 订单日志
func (orderController) Log(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("id", "desc", "id,order_id,operator_id,operator_type,type,create_time")
	var form struct {
		OrderId         string `json:"order_id" default:"" validate:"omitempty,numeric" label:"订单ID"`
		OperatorId      string `json:"operator_id" default:"" validate:"omitempty,numeric" label:"操作人ID"`
		OperatorType    string `json:"operator_type" default:"" validate:"omitempty,numeric" label:"操作人类型"`
		Type            string `json:"type" default:"" validate:"omitempty,numeric" label:"日志类型"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &order_log.SelectOrderLog{}
	if form.OrderId != "" {
		where.OrderId = form.OrderId
	}
	if form.OperatorId != "" {
		where.OperatorId = form.OperatorId
	}
	if form.OperatorType != "" {
		where.OperatorType = form.OperatorType
	}
	if form.Type != "" {
		where.Type = form.Type
	}
	if form.CreateTimeStart != "" {
		where.CreateTimeStart = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where.CreateTimeEnd = form.CreateTimeEnd
	}
	data, total, err := service.Order.SelectLog(where, orderByField, orderByType, pageNum, pageSize)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", page.NewPage(total, pageNum, pageSize, data))
}

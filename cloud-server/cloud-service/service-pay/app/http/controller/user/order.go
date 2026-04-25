package user

import (
	"service-pay/app/repository/order"
	"service-pay/app/repository/order_log"
	"service-pay/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

type orderController struct {
}

// List 订单列表
func (orderController) List(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("id", "desc", "id,type,status,pay_status,pay_type,pay_time,expired_time,complete_time,create_time,update_time")
	var form struct {
		TradeNo           string `json:"trade_no" default:"" validate:"omitempty" label:"订单号"`
		Type              string `json:"type" default:"" validate:"omitempty,numeric" label:"订单类型"`
		Status            string `json:"status" default:"" validate:"omitempty,numeric" label:"订单状态"`
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
	u := c.GetUserInfo()
	where := &order.SelectOrder{}
	if form.TradeNo != "" {
		where.TradeNo = form.TradeNo
	}
	if form.Type != "" {
		where.Type = form.Type
	}
	if form.Status != "" {
		where.Status = form.Status
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
	data, total, err := service.Order.SelectByUserId(u.Id, where, orderByField, orderByType, pageNum, pageSize)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", page.NewPage(total, pageNum, pageSize, data))
}

// Log 订单日志
func (orderController) Log(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("id", "desc", "id,type,create_time")
	var form struct {
		OrderId         int64  `json:"order_id" default:"0" validate:"required,numeric,min=1" label:"订单ID"`
		Type            string `json:"type" default:"" validate:"omitempty,numeric" label:"日志类型"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	u := c.GetUserInfo()
	info, err := service.Order.FindById(form.OrderId)
	if err != nil {
		c.Error(err.Error())
		return
	}
	if info.UserId != u.Id {
		c.Error("订单不存在")
		return
	}
	where := &order_log.SelectOrderLog{
		OrderId: strconv.FormatInt(form.OrderId, 10),
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

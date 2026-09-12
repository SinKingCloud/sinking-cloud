package log

import (
	"service-pay/app/repository/order_log"
	"service-pay/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// List 订单日志
func (Controller) List(c *context.Context) {
	query := c.ValidatePage("id", "desc", "id,create_time", "id,order_id,operator_id,operator_type,type,create_time")
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
		orderId, err := strconv.ParseInt(form.OrderId, 10, 64)
		if err != nil {
			c.Error("订单ID参数错误")
			return
		}
		where.OrderId = &orderId
	}
	if form.OperatorId != "" {
		operatorId, err := strconv.ParseInt(form.OperatorId, 10, 64)
		if err != nil {
			c.Error("操作人ID参数错误")
			return
		}
		where.OperatorId = &operatorId
	}
	if form.OperatorType != "" {
		operatorType, err := strconv.Atoi(form.OperatorType)
		if err != nil {
			c.Error("操作人类型参数错误")
			return
		}
		where.OperatorType = &operatorType
	}
	if form.Type != "" {
		types, err := strconv.Atoi(form.Type)
		if err != nil {
			c.Error("日志类型参数错误")
			return
		}
		where.Type = &types
	}
	if form.CreateTimeStart != "" {
		where.CreateTimeStart = &form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where.CreateTimeEnd = &form.CreateTimeEnd
	}
	data, err := service.Order.SelectLog(where, query)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", data)
}

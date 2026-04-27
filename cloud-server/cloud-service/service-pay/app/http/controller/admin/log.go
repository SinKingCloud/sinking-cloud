package admin

import (
	"service-pay/app/repository/order_log"
	"service-pay/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

type logController struct {
}

// List 订单日志
func (logController) List(c *context.Context) {
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

// Info 日志详情
func (logController) Info(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	info, err := service.Order.FindLogById(form.Id)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", info)
}

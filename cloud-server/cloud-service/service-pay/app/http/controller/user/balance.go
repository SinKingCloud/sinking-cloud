package user

import (
	"service-pay/app/repository/balance_log"
	"service-pay/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

type balanceController struct {
}

// Info 余额详情
func (balanceController) Info(c *context.Context) {
	u := c.GetUserInfo()
	info, err := service.Balance.FindByUserId(u.Id, false)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", info)
}

// Log 余额明细
func (balanceController) Log(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("id", "desc", "id,type,money,create_time")
	var form struct {
		Type            string `json:"type" default:"" validate:"omitempty,numeric" label:"类型"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	u := c.GetUserInfo()
	where := &balance_log.SelectBalanceLog{
		UserId: strconv.FormatInt(u.Id, 10),
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
	data, total, err := service.Balance.SelectLog(where, orderByField, orderByType, pageNum, pageSize)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", page.NewPage(total, pageNum, pageSize, data))
}

// Recharge 充值
func (balanceController) Recharge(c *context.Context) {
	var form struct {
		Money int    `json:"money" default:"" validate:"min=1" label:"充值金额"`
		Type  string `json:"type" default:"" label:"充值方式"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	u := c.GetUserInfo()
	tradeNo, err := service.Balance.CreateRechargeOrder(u.Id, form.Money, form.Type)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("创建订单成功", tradeNo)
}

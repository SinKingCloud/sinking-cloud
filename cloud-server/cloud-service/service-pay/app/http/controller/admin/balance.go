package admin

import (
	"service-pay/app/enum/balance_log_type"
	"service-pay/app/repository/balance_log"
	"service-pay/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

type balanceController struct {
}

// Info 用户余额详情
func (balanceController) Info(c *context.Context) {
	var form struct {
		UserId int64 `json:"user_id" default:"" validate:"numeric,min=1" label:"用户ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	info, err := service.Balance.FindByUserId(form.UserId, false)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", info)
}

// Operate 操作用户余额
func (balanceController) Operate(c *context.Context) {
	var form struct {
		UserId          int64   `json:"user_id" default:"0" validate:"required,numeric,min=1" label:"用户ID"`
		Money           float64 `json:"money" default:"0" validate:"required" label:"金额"`
		IsOperateFrozen bool    `json:"is_operate_frozen" default:"false" label:"是否操作冻结余额"`
		Describe        string  `json:"describe" default:"" validate:"required" label:"描述"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	logType := balance_log_type.Additional
	if form.Money < 0 {
		logType = balance_log_type.Deduct
	}
	if form.IsOperateFrozen {
		logType = balance_log_type.Freeze
		if form.Money < 0 {
			logType = balance_log_type.Defrost
		}
	}
	describe := "操作人ID:" + strconv.FormatInt(c.GetUserInfo().Id, 10) + "; " + form.Describe
	err := service.Balance.OperateMoney(form.UserId, form.Money, form.IsOperateFrozen, logType, describe)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("操作成功")
}

// Log 余额明细
func (balanceController) Log(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("id", "desc", "id,type,create_time")
	var form struct {
		UserId          int64  `json:"user_id" default:"" validate:"numeric,min=1" label:"用户ID"`
		Type            string `json:"type" default:"" validate:"omitempty,numeric" label:"类型"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &balance_log.SelectBalanceLog{
		UserId: strconv.FormatInt(form.UserId, 10),
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

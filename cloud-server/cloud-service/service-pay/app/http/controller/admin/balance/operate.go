package balance

import (
	"service-pay/app/enum/balance_log_type"
	"service-pay/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Operate 操作用户余额
func (Controller) Operate(c *context.Context) {
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

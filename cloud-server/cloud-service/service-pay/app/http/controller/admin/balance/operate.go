package balance

import (
	"fmt"
	"service-pay/app/enum/balance_log_type"
	"service-pay/app/service"
	"strconv"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Operate 操作用户余额
func (Controller) Operate(c *context.Context) {
	var form struct {
		UserId          int64   `json:"user_id" default:"0" validate:"numeric,min=1" label:"用户ID"`
		Money           float64 `json:"money" default:"0" validate:"numeric" label:"金额"`
		IsOperateFrozen bool    `json:"is_operate_frozen" default:"false" label:"是否操作冻结余额"`
		Describe        string  `json:"describe" default:"" validate:"omitempty" label:"描述"`
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
	describe := "操作人ID:" + strconv.FormatInt(c.GetUserInfo().Id, 10)
	if form.Describe != "" {
		describe += " 备注:" + form.Describe
	}
	err := service.Balance.OperateMoney(form.UserId, form.Money, form.IsOperateFrozen, logType, describe)
	if err != nil {
		c.Error(err.Error())
		return
	}
	_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventEdit, balance_log_type.Map()[logType], fmt.Sprintf("%s：用户[%d]，金额[%.4f]", balance_log_type.Map()[logType], form.UserId, form.Money))
	c.Success("操作成功")
}

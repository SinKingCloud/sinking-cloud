package balance

import (
	"fmt"
	"service-pay/app/repository/balance_log"
	"service-pay/app/service"
	"strconv"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Log 余额明细
func (Controller) Log(c *context.Context) {
	query := c.ValidatePage("id", "desc", "id,create_time", "id,type,create_time")
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
		UserId: &form.UserId,
	}
	if form.Type != "" {
		types, err := strconv.Atoi(form.Type)
		if err != nil {
			c.Error("类型参数错误")
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
	data, err := service.Balance.SelectLog(where, query)
	if err != nil {
		c.Error(err.Error())
		return
	}
	_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventShow, "查看用户余额明细", fmt.Sprintf("查看用户[%d]的余额明细", form.UserId))
	c.SuccessWithData("获取数据成功", data)
}

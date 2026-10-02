package balance

import (
	"fmt"
	"service-pay/app/service"
	"strconv"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// Info 用户余额详情
func (Controller) Info(c *context.Context) {
	var form struct {
		UserId int64 `json:"user_id" default:"" validate:"numeric,min=1" label:"用户ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	info, err := service.Balance.GetByUserId(form.UserId, false)
	if err != nil {
		c.Error(err.Error())
		return
	}
	_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventShow, "查看用户余额", fmt.Sprintf("查看用户[%d]的余额", form.UserId))
	c.SuccessWithData("获取数据成功", sinking_web.H{
		"money":             strconv.FormatFloat(info.Money, 'f', 2, 64),
		"frozen":            strconv.FormatFloat(info.Frozen, 'f', 2, 64),
		"total_recharge":    strconv.FormatFloat(info.TotalRecharge, 'f', 2, 64),
		"total_consume":     strconv.FormatFloat(info.TotalConsume, 'f', 2, 64),
		"last_operate_time": info.UpdateTime.String(),
	})
}

package balance

import (
	"service-pay/app/service"
	"strconv"

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
	info, err := service.Balance.FindByUserId(form.UserId, false)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", sinking_web.H{
		"money":             strconv.FormatFloat(info.Money, 'f', 2, 64),
		"frozen":            strconv.FormatFloat(info.Frozen, 'f', 2, 64),
		"total_recharge":    strconv.FormatFloat(info.TotalRecharge, 'f', 2, 64),
		"total_consume":     strconv.FormatFloat(info.TotalConsume, 'f', 2, 64),
		"last_operate_time": info.UpdateTime.String(),
	})
}

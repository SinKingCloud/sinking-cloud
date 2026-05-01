package balance

import (
	"service-pay/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Info 余额详情
func (Controller) Info(c *context.Context) {
	u := c.GetUserInfo()
	info, err := service.Balance.FindByUserId(u.Id, false)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", info)
}

package share

import (
	"service-space/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Info 分享导图详情
func (Controller) Info(c *context.Context) {
	var form struct {
		Code string `json:"code" default:"" validate:"required,len=8" label:"分享码"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	info, err := service.Mind.GetDataByShareCode(form.Code)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", info)
}

package mind

import (
	"service-space/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Share 获取公开导图分享
func (Controller) Share(c *context.Context) {
	var form struct {
		Code string `json:"code" default:"" validate:"required,len=8" label:"分享码"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data, err := service.Mind.GetDataByShareCode(form.Code)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", data)
}

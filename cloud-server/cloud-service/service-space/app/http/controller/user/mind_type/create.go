package mind_type

import (
	"service-space/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Create 创建导图分类
func (Controller) Create(c *context.Context) {
	var form struct {
		Name string `json:"name" default:"" validate:"required,max=24" label:"分类名称"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data, err := service.Mind.CreateType(c.GetUserInfo().Id, form.Name)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("创建数据成功", data)
}

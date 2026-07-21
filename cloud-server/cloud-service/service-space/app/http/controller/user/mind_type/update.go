package mind_type

import (
	"service-space/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Update 更新导图分类
func (Controller) Update(c *context.Context) {
	var form struct {
		Id   int64  `json:"id" validate:"required,min=1" label:"分类ID"`
		Name string `json:"name" default:"" validate:"required,max=24" label:"分类名称"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	if err := service.Mind.UpdateType(c.GetUserInfo().Id, form.Id, form.Name); err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("更新数据成功")
}

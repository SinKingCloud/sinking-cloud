package mind_data

import (
	"service-space/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Create 创建导图数据
func (Controller) Create(c *context.Context) {
	var form struct {
		TypeId   int64  `json:"type_id" default:"0" validate:"numeric,gte=0" label:"导图分类ID"`
		Title    string `json:"title" default:"" validate:"required,max=60" label:"导图标题"`
		Describe string `json:"describe" default:"" validate:"omitempty,max=500" label:"导图描述"`
		Value    string `json:"value" default:"" validate:"required" label:"导图内容"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data, err := service.Mind.CreateData(c.GetUserInfo().Id, form.TypeId, form.Title, form.Describe, form.Value)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("创建数据成功", data)
}

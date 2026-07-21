package mind_data

import (
	"service-space/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Save 保存导图数据
func (Controller) Save(c *context.Context) {
	var form struct {
		Id     int64  `json:"id" validate:"required,min=1" label:"导图ID"`
		Value  string `json:"value" default:"" validate:"required" label:"导图内容"`
		Type   int    `json:"type" default:"0" validate:"numeric" label:"保存类型"`
		Remark string `json:"remark" default:"" validate:"omitempty,max=60" label:"保存备注"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data, err := service.Mind.SaveData(c.GetUserInfo().Id, form.Id, form.Type, form.Remark, form.Value)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("保存成功", data)
}

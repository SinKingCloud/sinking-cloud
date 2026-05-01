package scene

import (
	"service-driver/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Delete 删除数据
func (Controller) Delete(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Scene.DeleteById(form.Id)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("删除数据成功")
}

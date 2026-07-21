package mind_type

import (
	"service-space/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Delete 删除导图分类
func (Controller) Delete(c *context.Context) {
	var form struct {
		Id int64 `json:"id" validate:"required,min=1" label:"分类ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	if err := service.Mind.DeleteType(c.GetUserInfo().Id, form.Id); err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("删除数据成功")
}

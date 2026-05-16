package master

import (
	"errors"
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Delete 删除用户
func (Controller) Delete(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"required,min=1" label:"用户ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	if form.Id == c.GetUserInfo().Id {
		c.Error("不能删除自身账户")
		return
	}
	err := service.User.DeleteById(form.Id)
	if err != nil {
		c.Error(errors.Join(errors.New("删除数据失败"), err).Error())
		return
	}
	c.Success("删除数据成功")
}

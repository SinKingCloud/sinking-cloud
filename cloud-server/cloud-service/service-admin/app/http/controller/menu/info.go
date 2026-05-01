package menu

import (
	"errors"
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Info 页面信息
func (Controller) Info(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	info, err := service.Rbac.GetMenuOrControlInfoById(form.Id)
	if err != nil {
		c.Error(errors.Join(errors.New("获取数据失败"), err))
		return
	}
	c.SuccessWithData("获取数据成功", info)
}

package file

import (
	"service-driver/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Delete 删除数据
func (Controller) Delete(c *context.Context) {
	var form struct {
		Ids   []int64 `json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
		Clear bool    `json:"clear" default:"false" label:"是否清理文件"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.File.DeleteByIds(form.Ids, form.Clear)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("更新数据成功")
}

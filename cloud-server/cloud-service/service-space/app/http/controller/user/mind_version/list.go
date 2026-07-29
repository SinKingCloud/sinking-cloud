package mind_version

import (
	"service-space/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// List 导图版本列表
func (Controller) List(c *context.Context) {
	query := c.ValidatePage("id", "desc", "id,create_time,update_time")
	var form struct {
		DataId int64  `json:"data_id" validate:"required,min=1" label:"导图ID"`
		Type   string `json:"type" default:"" validate:"required,numeric" label:"版本类型"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	versionType, err := strconv.Atoi(form.Type)
	if err != nil {
		c.Error("版本类型参数错误")
		return
	}
	data, err := service.Mind.SelectVersions(c.GetUserInfo().Id, form.DataId, versionType, query)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", data)
}

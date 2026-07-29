package mind_type

import (
	"service-space/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// List 自定义导图分类列表
func (Controller) List(c *context.Context) {
	query := c.ValidatePage("id", "desc", "id,create_time,update_time")
	data, err := service.Mind.SelectTypes(c.GetUserInfo().Id, query)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", data)
}

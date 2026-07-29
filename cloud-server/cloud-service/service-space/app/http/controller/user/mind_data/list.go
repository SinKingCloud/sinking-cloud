package mind_data

import (
	mindDataRepository "service-space/app/repository/mind_data"
	"service-space/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// List 导图数据列表
func (Controller) List(c *context.Context) {
	query := c.ValidatePage("id", "desc", "id,create_time,update_time", "id,create_time,update_time")
	var form struct {
		TypeId  int64  `json:"type_id" default:"0" validate:"omitempty,min=1" label:"导图分类ID"`
		Keyword string `json:"keyword" default:"" validate:"omitempty,max=100" label:"搜索关键词"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &mindDataRepository.SelectMindData{}
	if form.TypeId > 0 {
		where.TypeId = form.TypeId
	}
	if form.Keyword != "" {
		where.Keyword = form.Keyword
	}
	data, err := service.Mind.SelectDataByUserId(c.GetUserInfo().Id, where, query)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", data)
}

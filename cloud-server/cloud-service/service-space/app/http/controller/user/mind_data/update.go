package mind_data

import (
	mindDataRepository "service-space/app/repository/mind_data"
	"service-space/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Update 更新导图数据
func (Controller) Update(c *context.Context) {
	var form struct {
		Id       int64   `json:"id" validate:"required,min=1" label:"导图ID"`
		TypeId   string  `json:"type_id" default:"" validate:"omitempty,numeric" label:"导图分类ID"`
		Title    string  `json:"title" default:"" validate:"omitempty,max=60" label:"导图标题"`
		Describe *string `json:"describe" validate:"omitempty,max=500" label:"导图描述"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &mindDataRepository.UpdateMindData{}
	if form.TypeId != "" {
		data.TypeId, _ = strconv.ParseInt(form.TypeId, 10, 64)
	}
	if form.Title != "" {
		data.Title = form.Title
	}
	if form.Describe != nil {
		data.Describe = *form.Describe
	}
	if err := service.Mind.UpdateData(c.GetUserInfo().Id, form.Id, data); err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("更新数据成功")
}

package message

import (
	repositoryMessage "service-user/app/repository/message"
	"service-user/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// List 消息列表
func (Controller) List(c *context.Context) {
	query := c.ValidatePage("id", "desc", "id,status,read_time,create_time,update_time")
	var form struct {
		Title           string `json:"title" default:"" validate:"omitempty,max=100" label:"消息标题"`
		Status          string `json:"status" default:"" validate:"omitempty,numeric" label:"消息状态"`
		ReadTimeStart   string `json:"read_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"已读起始时间"`
		ReadTimeEnd     string `json:"read_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"已读结束时间"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &repositoryMessage.SelectMessage{}
	if form.Title != "" {
		where.Title = form.Title
	}
	if form.Status != "" {
		where.Status = form.Status
	}
	if form.ReadTimeStart != "" {
		where.ReadTimeStart = form.ReadTimeStart
	}
	if form.ReadTimeEnd != "" {
		where.ReadTimeEnd = form.ReadTimeEnd
	}
	if form.CreateTimeStart != "" {
		where.CreateTimeStart = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where.CreateTimeEnd = form.CreateTimeEnd
	}
	if form.UpdateTimeStart != "" {
		where.UpdateTimeStart = form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where.UpdateTimeEnd = form.UpdateTimeEnd
	}
	data, total, err := service.Message.SelectByUserId(c.GetUserInfo().Id, where, query)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", page.NewPage(total, query.Page, query.PageSize, data))
}

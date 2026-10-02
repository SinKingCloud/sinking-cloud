package master

import (
	"service-admin/app/repository/user_log"
	"service-admin/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Log 用户操作日志
func (Controller) Log(c *context.Context) {
	query := c.ValidatePage("id", "desc", "id,create_time", "id,type,create_time")
	var form struct {
		Id              int64  `json:"id" validate:"required,min=1" label:"用户ID"`
		Keyword         string `json:"keyword" default:"" validate:"omitempty,max=100" label:"关键词"`
		RequestId       string `json:"request_id" default:"" validate:"omitempty" label:"请求ID"`
		RequestIp       string `json:"request_ip" default:"" validate:"omitempty,max=50" label:"请求IP"`
		Location        string `json:"location" default:"" validate:"omitempty,max=100" label:"请求地点"`
		Type            string `json:"type" default:"" validate:"omitempty,numeric" label:"类型"`
		Title           string `json:"title" default:"" validate:"omitempty" label:"日志标题"`
		Content         string `json:"content" default:"" validate:"omitempty" label:"日志内容"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &user_log.SelectUserLog{
		UserId: &form.Id,
	}
	if form.Keyword != "" {
		where.Keyword = &form.Keyword
	}
	if form.RequestId != "" {
		where.RequestId = &form.RequestId
	}
	if form.RequestIp != "" {
		where.RequestIp = &form.RequestIp
	}
	if form.Location != "" {
		where.Location = &form.Location
	}
	if form.Type != "" {
		types, err := strconv.Atoi(form.Type)
		if err != nil {
			c.Error("类型参数错误")
			return
		}
		where.Type = &types
	}
	if form.Title != "" {
		where.Title = &form.Title
	}
	if form.Content != "" {
		where.Content = &form.Content
	}
	if form.CreateTimeStart != "" {
		where.CreateTimeStart = &form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where.CreateTimeEnd = &form.CreateTimeEnd
	}
	data, err := service.User.SelectLog(where, query)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", data)
}

package security

import (
	sessionRepository "service-user/app/repository/session"
	"service-user/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Session 会话列表
func (Controller) Session(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("create_time", "desc", "create_time,update_time,expire_time")
	var form struct {
		Device          string `json:"device" default:"" validate:"omitempty,max=100" label:"设备"`
		Ip              string `json:"ip" default:"" validate:"omitempty,max=50" label:"IP"`
		Location        string `json:"location" default:"" validate:"omitempty,max=100" label:"登录地点"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"登录起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"登录结束时间"`
		UpdateTimeStart string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
		ExpireTimeStart string `json:"expire_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"过期起始时间"`
		ExpireTimeEnd   string `json:"expire_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"过期结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &sessionRepository.SelectSession{}
	if form.Device != "" {
		where.Device = form.Device
	}
	if form.Ip != "" {
		where.Ip = form.Ip
	}
	if form.Location != "" {
		where.Location = form.Location
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
	if form.ExpireTimeStart != "" {
		where.ExpireTimeStart = form.ExpireTimeStart
	}
	if form.ExpireTimeEnd != "" {
		where.ExpireTimeEnd = form.ExpireTimeEnd
	}
	data, total, err := service.Auth.SelectSessionByUserId(c.GetUserInfo().Id, where, orderByField, orderByType, pageNum, pageSize)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", page.NewPage(total, pageNum, pageSize, data))
}

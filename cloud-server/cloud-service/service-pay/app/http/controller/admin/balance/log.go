package balance

import (
	"service-pay/app/repository/balance_log"
	"service-pay/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Log 余额明细
func (Controller) Log(c *context.Context) {
	query := c.ValidatePage("id", "desc", "id,type,create_time")
	var form struct {
		UserId          int64  `json:"user_id" default:"" validate:"numeric,min=1" label:"用户ID"`
		Type            string `json:"type" default:"" validate:"omitempty,numeric" label:"类型"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &balance_log.SelectBalanceLog{
		UserId: strconv.FormatInt(form.UserId, 10),
	}
	if form.Type != "" {
		where.Type = form.Type
	}
	if form.CreateTimeStart != "" {
		where.CreateTimeStart = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where.CreateTimeEnd = form.CreateTimeEnd
	}
	data, total, err := service.Balance.SelectLog(where, query)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", page.NewPage(total, query.Page, query.PageSize, data))
}

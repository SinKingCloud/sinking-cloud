package pay

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/log"
	"strconv"
)

func (ControllerPay) LogList(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id create_time" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Type            string `form:"type" json:"type" default:"" validate:"omitempty,oneof=0 1" label:"记录类型"`
		LogType         string `form:"log_type" json:"log_type" default:"" validate:"omitempty,numeric" label:"日志类型"`
		CreateTimeStart string `form:"create_time_start" json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `form:"create_time_end" json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	where["user_id"] = strconv.Itoa(c.GetWebUserInfo().Id)
	if form.Type != "" {
		where["type"] = form.Type
	}
	if form.LogType != "" {
		where["log_type"] = form.LogType
	}
	if form.CreateTimeStart != "" {
		where["create_time_start"] = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where["create_time_end"] = form.CreateTimeEnd
	}
	data, total, err := service.PayLog.Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取消费记录失败")
	} else {
		service.Log.CreateByQueue(c, log.EventShow, "查看消费", "查看消费记录")
		c.SuccessWithData("获取消费记录成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

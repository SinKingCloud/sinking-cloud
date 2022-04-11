package notice

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/service/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/service/notice"
)

// Info 获取详情
func (ControllerNotice) Info(c *sinking_server.Context) {
	type Form struct {
		Id int `form:"id" default:"" validate:"omitempty,max=36" label:"记录ID"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	data, err := notice.GetIns().FindById(form.Id)
	if err == nil && data != nil {
		log.GetIns().CreateByQueue(c, log.EventShow, "查看公告", "查看公告"+data.Title)
		c.SuccessWithData("获取详情数据成功", data)
	} else {
		c.Error("获取详情数据失败")
	}
}

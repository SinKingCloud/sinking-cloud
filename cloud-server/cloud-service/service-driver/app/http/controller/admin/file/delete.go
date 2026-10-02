package file

import (
	"fmt"
	"service-driver/app/service"

	commonService "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Delete 删除数据
func (Controller) Delete(c *context.Context) {
	var form struct {
		Ids   []int64 `json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
		Clear bool    `json:"clear" default:"false" label:"是否清理文件"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.File.DeleteByIds(form.Ids, form.Clear)
	if err != nil {
		c.Error(err.Error())
		return
	}
	content := fmt.Sprintf("删除文件记录%v", form.Ids)
	if form.Clear {
		content += "，请求清理存储文件"
	}
	_ = commonService.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventDelete, "删除文件记录", content)
	c.Success("更新数据成功")
}

package oauth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_oauth/oauth"
)

// Delete 删除
func (ControllerOauth) Delete(c *sinking_server.Context) {
	type Form struct {
		Id int `form:"id" json:"id" default:"" validate:"min=1" label:"ID"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	info, err := service.Oauth.FindByIdWithCache(form.Id, false)
	if err != nil || info == nil || info.UserId != c.GetWebUserInfo().Id {
		c.Error("接入网站不存在")
		return
	}
	if info.Status == oauth.Error {
		c.Error("该接入网站已被封禁，无法操作")
		return
	}
	err = service.Oauth.DeleteById(form.Id)
	if err == nil {
		c.Success("删除接入网站成功")
	} else {
		c.Error("删除接入网站失败")
	}
}

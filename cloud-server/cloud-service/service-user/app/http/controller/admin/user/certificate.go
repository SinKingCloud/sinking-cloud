package user

import (
	"service-user/app/enum/real_auth_status"
	"service-user/app/enum/real_auth_type"
	"service-user/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// Certificate 用户实名信息
func (Controller) Certificate(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"required,min=1" label:"用户ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	info, err := service.User.GetRealAuthInfo(form.Id, false)
	if err != nil {
		c.Error(err.Error())
		return
	}
	if info == nil {
		c.Error("未获取到实名信息")
		return
	}
	data := sinking_web.H{
		"status": info.Status == real_auth_status.Authed,
	}
	if info.Status == real_auth_status.Authed {
		data["auth_time"] = info.AuthTime.String()
		data["card_name"] = info.CardName
		data["card_no"] = info.CardNo
		if info.Files != nil {
			data["bucket"] = info.Files.Bucket
			data["pics"] = info.Files.PicPaths
			data["video"] = info.Files.VideoPath
		}
		if info.Type == real_auth_type.Company {
			data["company_name"] = info.CompanyName
			data["company_no"] = info.CompanyNo
		}
	}
	c.SuccessWithData("获取数据成功", data)
}

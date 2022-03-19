package verify

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/aliyun"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/service/captcha"
	"strconv"
)

func (ControllerVerify) Sms(c *sinking_server.Context) {
	type Form struct {
		Phone   string `form:"phone" json:"phone" default:"" validate:"-" label:"手机号"`
		Token   string `form:"token" json:"token" default:"" validate:"required,len=16" label:"验证码标识"`
		Captcha string `form:"captcha" json:"captcha" default:"" validate:"required,numeric,len=4" label:"验证码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	if !captcha.Check(form.Token, form.Captcha) {
		c.Error("验证码验证失败请重试")
		return
	}
	configs := service.GetSystemConfigServiceIns().GetByCache(false)
	code := strconv.FormatInt(str.NetStrTool().RandInt64(100000, 999999), 10)
	if aliyun.Sms.Send(
		configs[config.AliyunAccessKey],
		configs[config.AliyunAccessSecret],
		form.Phone,
		configs[config.AliyunSmsSignName],
		configs[config.AliyunSmsTemplateCode],
		"{\"code\":"+code+"}",
	) {
		util.Cache.SetWithExpire(constant.CacheNameWithSmsPhone+form.Phone, code, constant.CacheTimeWithSmsPhone)
		c.Success("验证码已发送")
	} else {
		c.Error("验证码发送失败")
	}
}

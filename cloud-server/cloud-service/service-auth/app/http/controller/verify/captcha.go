package verify

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service/captcha"
	"image/png"
)

// Code 生成验证码
func (ControllerVerify) Code(c *sinking_server.Context) {
	type Form struct {
		Token string `form:"token" default:"" validate:"required,len=16" label:"验证码标识"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	img, _ := captcha.Get(form.Token)
	err := png.Encode(c.Writer, img)
	if err != nil {
		return
	}
}

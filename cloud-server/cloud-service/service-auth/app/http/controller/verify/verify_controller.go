package verify

import (
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/captcha"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"image/png"
)

func (ControllerVerify) Code(c *sinking_server.Context) {
	type Form struct {
		Token string `form:"token" default:"" validate:"required" label:"验证码标识"`
	}
	form := &Form{}
	ok, msg := c.ValidatorGet(form)
	if !ok {
		c.Error(msg)
		return
	}
	img, code := captcha.NewCaptchaRandCode(4)
	err := png.Encode(c.Writer, img)
	if err != nil {
		return
	}
	//验证码写入redis
	fmt.Println(code)
}

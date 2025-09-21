package captcha

import (
	"github.com/afocus/captcha"
	"image/color"
)

var c *captcha.Captcha

func init() {
	c = captcha.New()
	if err := c.AddFontFromBytes(NewTtf()); err != nil {
		panic(err.Error())
	}
	c.SetSize(256, 128)
	c.SetDisturbance(64)
	c.SetFrontColor(color.RGBA{R: 0, G: 81, B: 235, A: 255})
}

// GenCode 生成随机验证码图片
// num 验证码个数
func GenCode(num int) (*captcha.Image, string) {
	img, code := c.Create(num, captcha.ALL)
	return img, code
}

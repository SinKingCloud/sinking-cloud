package spider

import (
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-spider/app/util/spider/ocr_captcha"
)

type ControllerSpider struct {
}

// OcrCaptcha 获取ocr cookie
func (ControllerSpider) OcrCaptcha(c *sinking_server.Context) {
	data := c.GetRpcData()
	url := data["captcha_url"].(string)
	cookie, e := ocr_captcha.Captcha(url)
	fmt.Println(cookie, e)
	if e == nil {
		c.SuccessWithData("success", cookie)
	} else {
		c.Error(e.Error())
	}
}

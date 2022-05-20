package ocr_captcha

import (
	"errors"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"math/rand"
	"time"
)

func Captcha(securityUrl string) (string, error) {
	url := "https://duguang.aliyun.com/" + securityUrl
	browser := rod.New().MustConnect()
	//实例化
	var page *rod.Page
	c := func(browser *rod.Browser, page *rod.Page) {
		_ = page.Close()
		_ = browser.Close()
	}
	cookies := ""
	var e error
	err := rod.Try(func() {
		page = browser.Timeout(30 * time.Second).MustPage(url)
		//设置设备
		err1 := page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
			Width:  390,
			Height: 844,
			Mobile: true,
		})
		if err1 != nil {
			e = errors.New("设置视口错误")
			return
		}
		//等待页面加载完毕
		page.MustWaitLoad()
		page.MustElement("#nc_1_n1z")
		//获取拖动对象
		button := page.MustElement("#nc_1_n1z").MustShape().OnePointInside()
		//鼠标拖动
		mouse := page.Mouse
		mouse.MustMove(button.X+10, button.Y)
		mouse.MustDown(proto.InputMouseButtonLeft)
		err := mouse.Move(button.X+340, button.Y+30, 180+rand.Intn(3))
		if err == nil {
			mouse.MustUp("left")
			net, e1 := page.Cookies([]string{"https://duguang.aliyun.com"})
			if e1 != nil {
				e = errors.New("获取cookie失败")
				return
			}
			fail := true
			for _, v := range net {
				if v.Name == "x5sec" {
					fail = false
					cookies = v.Value
				}
			}
			if fail {
				e = errors.New("获取cookie失败")
			}
		} else {
			e = errors.New("滑动验证失败")
		}
	})
	c(browser, page)
	if err != nil {
		e = errors.New("验证超时")
	}
	return cookies, e
}

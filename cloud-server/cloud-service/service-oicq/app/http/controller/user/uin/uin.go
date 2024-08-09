package uin

import (
	"bytes"
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc"
	userLog "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/command/queues/producer"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/app_config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin/log"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin/server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin/uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_user"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/util/oicq"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/util/server"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"github.com/SinKingCloud/sinking-oicq/binary-mobile/client"
	client2 "github.com/SinKingCloud/sinking-oicq/binary-pc/client"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"gorm.io/gorm"
	"image/png"
	log2 "log"
	"math/rand"
	"strconv"
	"time"
)

type ControllerUin struct {
}

// List 列表
func (ControllerUin) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		AppId           string `form:"app_id" json:"app_id" default:"" validate:"omitempty,numeric" label:"应用ID"`
		ServerKey       string `form:"server_key" json:"server_key" default:"" validate:"omitempty,max=50" label:"地区标识"`
		UinStatus       string `form:"uin_status" json:"uin_status" default:"" validate:"omitempty,numeric" label:"挂机状态"`
		Uin             string `form:"uin" json:"uin" default:"" validate:"omitempty,numeric" label:"QQ账号"`
		Status          string `form:"status" json:"status" default:"" validate:"omitempty,oneof=0 1 2" label:"QQ状态"`
		RunTimeStart    string `form:"run_time_start" json:"run_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"挂机起始时间"`
		RunTimeEnd      string `form:"run_time_end" json:"run_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"挂机结束时间"`
		CreateTimeStart string `form:"create_time_start" json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `form:"create_time_end" json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `form:"update_time_start" json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `form:"update_time_end" json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	where["user_id"] = strconv.Itoa(c.GetWebUserInfo().Id)
	if form.AppId != "" {
		where["app_id"] = form.AppId
	}
	if form.ServerKey != "" {
		where["server_key"] = form.ServerKey
	}
	if form.UinStatus != "" {
		where["uin_status"] = form.UinStatus
	}
	if form.Uin != "" {
		where["uin"] = form.Uin
	}
	if form.Status != "" {
		where["status"] = form.Status
	}
	if form.RunTimeStart != "" {
		where["run_time_start"] = form.RunTimeStart
	}
	if form.RunTimeEnd != "" {
		where["run_time_end"] = form.RunTimeEnd
	}
	if form.CreateTimeStart != "" {
		where["create_time_start"] = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where["create_time_end"] = form.CreateTimeEnd
	}
	if form.UpdateTimeStart != "" {
		where["update_time_start"] = form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where["update_time_end"] = form.UpdateTimeEnd
	}
	data, total, err := uin.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取挂机列表失败")
	} else {
		rpc.Log.CreateUser(c.GetWebUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), userLog.UserEventShow, "查看挂机", "查看挂机列表")
		c.SuccessWithData("获取挂机列表成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

// Test 测试加速
func (ControllerUin) Test(c *sinking_server.Context) {
	type Form struct {
		ServerKey string `form:"server_key" json:"server_key" validate:"required,max=10" label:"挂机地区"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	s, err := web_uin.GetServer(form.ServerKey)
	if err != nil || s == nil {
		c.Error("未获取到可用节点,请更换挂机节点")
		return
	}
	appConfig := app_config.GetIns().GetAppConfigs(constant.AppUinId, false)
	if appConfig[constant.AppUinSpeedTest] != "1" {
		c.Error("系统已关闭测试功能")
		return
	}
	qps := 0
	if appConfig[constant.AppUinSpeedQps] != "" {
		num, e := strconv.Atoi(appConfig[constant.AppUinSpeedQps])
		if e == nil {
			qps = num
		}
	}
	num := 0
	if appConfig[constant.AppUinSpeedNum] != "" {
		n, e := strconv.Atoi(appConfig[constant.AppUinSpeedNum])
		if e == nil {
			num = n
		}
	}
	var proxy *oicq.Proxy
	_ = json.Unmarshal([]byte(s.Proxy), &proxy)
	if proxy == nil {
		c.Error("该挂机节点不可用")
		return
	}
	key := "SpeedTestU_" + strconv.Itoa(c.GetWebUserInfo().Id)
	if qps > 0 {
		n, exists := util.Cache.Get(key)
		if exists {
			num, e := strconv.Atoi(n)
			if e == nil && num >= qps {
				c.Error("并发超过系统最大限制")
				return
			}
		}
	}
	util.Cache.Inc(key)
	util.Cache.SetExpire(key, int64(120*1000))
	cli := client2.NewClientEmpty().ChooseServerMode(true).UseTcp(true).UseProxy(proxy.Addr, proxy.Port, proxy.User, proxy.Pwd, true)
	e, u := cli.FetchQrCode()
	if e != nil || u == "" {
		c.Error("获取二维码失败")
		return
	}
	go func(cli *client2.QQClient, key string, num int) {
		defer func(cli *client2.QQClient) {
			cli.Release()
			cli.Disconnect()
			util.Cache.Dec(key)
			cli = nil
		}(cli)
		for i := 0; i < 60; i++ {
			e1, q := cli.FetchQrCodeStatus()
			if e1 != nil {
				break
			}
			if q == client2.ScanCodeSuccess {
				//上线登陆
				e2, l := cli.Login(true)
				if e2 != nil {
					break
				}
				if l.Status == client2.LoginSuccess {
					cli.Release()
					cli.Disconnect()
					key2 := "SpeedTestQ_" + strconv.FormatInt(cli.Uin, 10)
					util.Cache.Inc(key2)
					util.Cache.SetExpire(key2, (str.NewTimeTool().GetZeroTime(time.Now()).Unix()+86400-time.Now().Unix())*1000)
					if num > 0 {
						n, exists := util.Cache.Get(key2)
						if exists {
							n1, ne := strconv.Atoi(n)
							if ne == nil && n1 > num {
								break
							}
						}
					}
					//获取加速项目
					e4, status := cli.Speed.Info()
					if e4 == nil && status != nil {
						//执行加速操作
						if status.IContinueLogin != "1" {
							_ = cli.Speed.ContinueLogin()
						}
						if status.IDailySign != "1" {
							_, _ = cli.Speed.DailySign()
						}
						if status.IQzoneState != "1" {
							_, _ = cli.Speed.QZone()
						}
						if status.WeishiVideoview != "1" {
							_, _, _ = cli.Speed.WeiShi("")
						}
						if status.IAddFriend != "1" {
							_, _ = cli.Speed.Friend()
						}
						if status.IQQMusic != "1" {
							_, _, _, _ = cli.Speed.Music("", "")
						}
					}
				}
				break
			}
			if q == client2.ScanCodeExpire || q == client2.ScanCodeError {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}(cli, key, num)
	c.SuccessWithData("获取二维码成功,请在60s内扫描二维码", u)
}

// Experience 领取体验
func (ControllerUin) Experience(c *sinking_server.Context) {
	type Form struct {
		Uin       string `form:"uin" json:"uin" validate:"required,numeric,min=5,max=10" label:"QQ账号"`
		ServerKey string `form:"server_key" json:"server_key" validate:"omitempty,max=10" label:"挂机地区"`
		AppId     int    `form:"app_id" json:"app_id" validate:"omitempty,numeric,gte=0" label:"应用ID"`
		Pwd       string `form:"pwd" json:"pwd" validate:"omitempty,len=32" label:"QQ密码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	u := c.GetWebUserInfo()
	//判断是否新添加
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(u.Id, form.Uin)
	if err != nil || uinInfo == nil {
		//新开通
		if form.Pwd == "" || len(form.Pwd) != 32 {
			c.Error("QQ密码格式不正确")
			return
		}
		if form.ServerKey == "" {
			c.Error("挂机地区不能为空")
			return
		}
	}
	//判断是否携带APPID
	if form.AppId == 0 {
		appInfo := server.GetAppInfo(c)
		if appInfo != nil {
			form.AppId = appInfo.Id
		}
	}
	//领取
	err = web_uin.Experience(u.Id, form.AppId, form.ServerKey, form.Uin, form.Pwd)
	if err != nil {
		c.Error(err.Error())
	} else {
		c.Success("领取体验套装成功")
	}
}

// Buy 下单
func (ControllerUin) Buy(c *sinking_server.Context) {
	type Form struct {
		Uin       string `form:"uin" json:"uin" validate:"required,numeric,min=5,max=10" label:"QQ账号"`
		ServerKey string `form:"server_key" json:"server_key" validate:"omitempty,max=10" label:"挂机地区"`
		AppId     int    `form:"app_id" json:"app_id" validate:"omitempty,numeric,gte=0" label:"应用ID"`
		Pwd       string `form:"pwd" json:"pwd" validate:"omitempty,len=32" label:"QQ密码"`
		Type      int    `form:"type" json:"type" validate:"oneof=0 1 2" label:"下单类型"`
		Month     int    `form:"month" json:"month" validate:"oneof=1 3 6 12 999" label:"下单月数"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	//判断是否新添加
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		//新开通
		if form.Pwd == "" || len(form.Pwd) != 32 {
			c.Error("QQ密码格式不正确")
			return
		}
		if form.ServerKey == "" {
			c.Error("挂机地区不能为空")
			return
		}
	}
	//判断是否携带APPID
	if form.AppId == 0 {
		appInfo := server.GetAppInfo(c)
		if appInfo != nil {
			form.AppId = appInfo.Id
		}
	}
	u := c.GetWebUserInfo()
	//下单
	err = web_uin.Buy(u.Id, form.AppId, form.ServerKey, form.Uin, form.Pwd, form.Type, form.Month, true)
	if err != nil {
		c.Error(err.Error())
	} else {
		c.Success("下单成功")
	}
}

// ChangePassword 更改挂机QQ密码
func (ControllerUin) ChangePassword(c *sinking_server.Context) {
	type Form struct {
		Uin string `form:"uin" json:"uin" validate:"required,numeric" label:"QQ账号"`
		Pwd string `form:"pwd" json:"pwd" validate:"required,len=32" label:"QQ密码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	if uinInfo.Status == uin.Error {
		c.Error("该QQ已被系统锁定不可操作,请咨询管理员")
		return
	}
	if uinInfo.PlayerDay <= 0 && uinInfo.SignDay <= 0 {
		c.Error("该QQ未下单或已过期,无法操作")
		return
	}
	if !util.Cache.Lock(constant.LockUserUinOperate+uinInfo.Uin, constant.LockTimeUserUinOperate) {
		c.Error("操作过于频繁,请等待5s后再试")
		return
	}
	err = web_uin.ChangePassword(uinInfo.Id, form.Pwd)
	if err == nil {
		c.Success("更改挂机密码成功,请重新登录上线")
	} else {
		c.Error("更改挂机密码失败,请稍后重试")
	}
}

// ChangeServer 更改挂机地区
func (ControllerUin) ChangeServer(c *sinking_server.Context) {
	type Form struct {
		Uin       string `form:"uin" json:"uin" validate:"required,numeric" label:"QQ账号"`
		ServerKey string `form:"server_key" json:"server_key" validate:"required,max=32" label:"挂机地区"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	if uinInfo.Status == uin.Error {
		c.Error("该QQ已被系统锁定不可操作,请咨询管理员")
		return
	}
	if uinInfo.PlayerDay <= 0 && uinInfo.SignDay <= 0 {
		c.Error("该QQ未下单或已过期,无法操作")
		return
	}
	err = web_uin.ChangeServer(uinInfo.Id, form.ServerKey)
	if err == nil {
		c.Success("更改挂机地区成功")
	} else {
		c.Error(err.Error())
	}
}

// ChangeTiming 修改挂机时间
func (ControllerUin) ChangeTiming(c *sinking_server.Context) {
	type Form struct {
		Uin    string `form:"uin" json:"uin" validate:"required,numeric" label:"QQ账号"`
		Timing int    `form:"timing" json:"timing" validate:"numeric,gt=0,lt=1440" label:"挂机时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	if uinInfo.Status == uin.Error {
		c.Error("该QQ已被系统锁定不可操作,请咨询管理员")
		return
	}
	if uinInfo.PlayerDay <= 0 && uinInfo.SignDay <= 0 {
		c.Error("该QQ未下单或已过期,无法操作")
		return
	}
	err = web_uin.ChangeTiming(uinInfo.Id, form.Timing)
	if err == nil {
		c.Success("更改挂机时间成功")
	} else {
		c.Error(err.Error())
	}
}

// ReSetDevice 重置挂机设备
func (ControllerUin) ReSetDevice(c *sinking_server.Context) {
	type Form struct {
		Uin string `form:"uin" json:"uin" validate:"required,numeric" label:"QQ账号"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	if uinInfo.Status == uin.Error {
		c.Error("该QQ已被系统锁定不可操作,请咨询管理员")
		return
	}
	err = web_uin.ReSetDeviceId(uinInfo.Id)
	if err == nil {
		c.Success("重置挂机设备成功")
	} else {
		c.Error(err.Error())
	}
}

// ChangeStatus 修改挂机状态
func (ControllerUin) ChangeStatus(c *sinking_server.Context) {
	type Form struct {
		Uin    string `form:"uin" json:"uin" validate:"required,numeric" label:"QQ账号"`
		Status int    `form:"status" json:"status" validate:"oneof=0 1" label:"挂机状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	if uinInfo.Status == uin.Error {
		c.Error("该QQ已被系统锁定不可操作,请咨询管理员")
		return
	}
	err = web_uin.ChangeStatus(uinInfo.Id, form.Status)
	if err == nil {
		c.Success("更改挂机状态成功")
	} else {
		c.Error(err.Error())
	}
}

// Login 登陆QQ
func (ControllerUin) Login(c *sinking_server.Context) {
	type Form struct {
		Uin    string `form:"uin" json:"uin" validate:"required,numeric" label:"QQ账号"`
		Step   string `form:"step" json:"step" validate:"oneof=normal captcha get_sms sms" label:"登陆步骤"`
		Rand   string `form:"rand" json:"rand" validate:"omitempty,max=500" label:"登陆rand"`
		Ticket string `form:"ticket" json:"ticket" validate:"omitempty,max=500" label:"登陆ticket"`
		Code   string `form:"code" json:"code" validate:"omitempty,numeric,len=6" label:"短信验证码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if uin.GetIns().Count(map[string]string{"uin": form.Uin, "status": strconv.Itoa(uin.Error)}) > 0 {
		c.Error("该QQ已被系统锁定不可操作,请咨询管理员")
		return
	}
	user := web_user.GetIns().GetCurrentUserInfo(c.GetWebUserInfo().Id, false)
	if user.RealAuth != web_user.RealAuthSuccess {
		c.Error("该对接账户未完成实名认证,请完成实名认证")
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	if uinInfo.PlayerDay <= 0 && uinInfo.SignDay <= 0 {
		c.Error("该QQ未下单或已过期,无法操作")
		return
	}
	if form.Step == "normal" {
		appConfig := app_config.GetIns().GetAppConfigs(constant.AppUinId, false)
		if appConfig[constant.AppUinLoginPwd] == "1" {
			c.Error("系统已关闭密码更新")
			return
		}
		checkMin := true
		if appConfig[constant.AppUinUpdateDate] != "" {
			t1, e1 := time.Parse("2006-01-02 15:04:05", appConfig[constant.AppUinUpdateDate])
			t2, e2 := time.Parse("2006-01-02 15:04:05", uinInfo.CreateTime.String())
			if e1 == nil && e2 == nil && t2.Before(t1) {
				checkMin = false
			}
		}
		if checkMin {
			minLevel := 0
			if appConfig[constant.AppUinUpdateLevel] != "" {
				num, e := strconv.Atoi(appConfig[constant.AppUinUpdateLevel])
				if e == nil {
					minLevel = num
				}
			}
			if minLevel > 0 {
				minDays := 0
				if appConfig[constant.AppUinUpdateDays] != "" {
					num, e := strconv.Atoi(appConfig[constant.AppUinUpdateDays])
					if e == nil {
						minDays = num
					}
				}
				u, _ := strconv.ParseInt(uinInfo.Uin, 10, 64)
				levelInfo := web_uin.GetGuestUinInfo(u)
				if levelInfo == nil {
					levelInfo = web_uin.GetGuestUinInfoBak(u)
				}
				if levelInfo != nil && levelInfo.Level <= minLevel && uinInfo.PlayerDay < minDays && uinInfo.SignDay < minDays {
					c.Error("您的剩余挂机天数低于" + strconv.Itoa(minDays) + "天。请续费至高于" + strconv.Itoa(minDays) + "天再更新状态")
					return
				}
			}
		}
	}
	srv, err := server2.GetIns().FindById(uinInfo.ServerId)
	if err != nil || srv == nil || srv.Status != server2.Normal {
		c.Error("未查询到挂机地区信息,请更换挂机地区")
		return
	}
	changeProto := func(protos client.ClientProtocol) {
		//判断是否切换协议
		proto := oicq.GetOicqProtoCode(protos)
		if uinInfo.Proto != proto {
			err = uin.GetIns().UpdateByUin(uinInfo.Uin, map[string]interface{}{"proto": proto})
			if err != nil {
				return
			}
			uinInfo.Proto = proto
		}
	}
	//更换协议
	changeProto(client.AndroidPad)
	//获取单例
	cli, err := web_uin.GetOicqInsById(uinInfo.Id, true)
	if err != nil {
		c.Error("连接服务器失败")
		return
	}
	cli.LoginTime = time.Now().Unix()
	ins := cli.Mobile
	if form.Step == "get_sms" {
		if ins.RequestSMS() {
			c.Return(200, "获取验证码成功,请提交验证码", nil)
		} else {
			c.Return(500, "获取验证码失败,请稍后重试", nil)
		}
		return
	}
	var resp *client.LoginResponse
	request := func(uinInfo *uin.AppUinJob, fun func() (*client.LoginResponse, error)) (resp *client.LoginResponse, err error) {
		resp, err = fun()
		log2.Println(uinInfo.Uin, resp, err)
		if err != nil {
			log2.Println(uinInfo.Uin, "登陆错误-1", resp, err)
			ins.Release()
			ins.Disconnect()
			time.Sleep(100 * time.Millisecond)
			_ = ins.Connect()
			time.Sleep(100 * time.Millisecond)
			resp, err = fun()
			if err != nil {
				log2.Println(uinInfo.Uin, "登陆错误-2", resp, err)
			}
		}
		//设备异常,生成新的设备信息
		if err == nil && resp.Error == client.DeviceVersionError {
			field := "device"
			if uin.GetIns().UpdateById(uinInfo.Id, map[string]interface{}{field: 100000 + rand.Intn(899999)}) == nil {
				q, e := strconv.ParseInt(uinInfo.Uin, 10, 64)
				if e == nil {
					oicq.DelOicqInstance(q) //删除实例
					oicq.DelTempOicqInstance(q)
					cli, err = web_uin.GetOicqInsById(uinInfo.Id, true)
					if err != nil {
						return nil, err
					}
					ins = cli.Mobile
					resp, err = fun()
					return resp, err
				}
			}
		}
		return resp, err
	}
	switch form.Step {
	case "normal":
		resp, err = request(uinInfo, func() (*client.LoginResponse, error) {
			return ins.Login()
		})
		break
	case "captcha":
		if form.Ticket == "" {
			c.Error("ticket不能为空")
			return
		}
		resp, err = request(uinInfo, func() (*client.LoginResponse, error) {
			return ins.SubmitTicket(form.Ticket)
		})
		break
	case "sms":
		if form.Code == "" {
			c.Error("验证码不能为空")
			return
		}
		resp, err = request(uinInfo, func() (*client.LoginResponse, error) {
			return ins.SubmitSMS(form.Code)
		})
		break
	}
	if err != nil || resp == nil {
		log2.Println(uinInfo.Uin, "登陆错误", resp, err)
		c.Error("登陆发生错误,请使用扫码登陆重试")
		return
	}
	response := func(uinInfo *uin.AppUinJob, resp *client.LoginResponse) (code int32, message interface{}, data interface{}) {
		if resp.Success {
			return 200, "更新状态成功", nil
		} else if resp.Error == client.SliderNeededError {
			return 400, "请验证滑动验证码", resp.VerifyUrl //需要滑动验证码
		} else if resp.Error == client.PhoneVerifyNeededError {
			return 399, "当前登录存在安全风险,请验证手机号后登陆", nil //需要发送验证码
		} else if resp.Error == client.SMSOrVerifyNeededError || resp.Error == client.UnsafeDeviceError || (resp.Error == client.SMSNeededError && resp.SMSPhone != "") {
			if resp.VerifyUrl != "" {
				resp.VerifyUrl, err = ins.GetQrCodeAuthToken(resp.VerifyUrl)
				if err == nil {
					resp.VerifyUrl = "https://accounts.qq.com/safe/scanresult?str_url=" + resp.VerifyUrl + "&envfrom=&verify_id=&verify_scene=&uin=" + uinInfo.Uin
				}
			}
			return 398, "请选择一种方式验证设备锁", sinking_web.H{
				"phone":  resp.SMSPhone,
				"qrcode": resp.VerifyUrl,
			}
		} else if resp.Error == client.UnknownLoginError {
			msg := "账号被冻结,请解冻后在更新状态"
			if resp.ErrorMessage != "" {
				msg = resp.ErrorMessage
			}
			return 500, msg, nil
		} else if resp.Error == client.TooManySMSRequestError {
			return 500, "获取验证码频繁,请稍后重试", nil
		} else if resp.Error == client.SMSNeededError {
			return 500, "登陆环境异常,请稍后重试", nil
		} else if resp.Error == client.OtherLoginError {
			log2.Println(uinInfo.Uin, "登陆异常", resp, err)
			if resp.ErrorMessage == "" {
				resp.ErrorMessage = "登陆异常,请重试"
			}
			return 500, resp.ErrorMessage, nil
		} else {
			log2.Println(uinInfo.Uin, "登陆失败", resp, err)
			if resp.ErrorMessage == "" {
				resp.ErrorMessage = "登陆失败,请重试"
			}
			return 500, resp.ErrorMessage, nil
		}
	}
	delIns := func(uin string) {
		q, e := strconv.ParseInt(uin, 10, 64)
		if e == nil {
			oicq.DelOicqInstance(q) //删除实例
			oicq.DelTempOicqInstance(q)
		}
	}
	code, message, data := response(uinInfo, resp)
	if code == 200 {
		token := ins.GenToken()
		if web_uin.SaveLoginToken(uinInfo, token) == nil {
			go func(uinInfo *uin.AppUinJob) {
				//写入QQ操作记录
				_ = log.GetIns().Create(&model.AppUinLog{
					UserId:  uinInfo.UserId,
					AppId:   uinInfo.AppId,
					UinId:   uinInfo.Id,
					Uin:     uinInfo.Uin,
					Type:    log.Login,
					Title:   "更新状态",
					Content: "密码更新账号状态",
				})
				_ = ins.Init(true, false)
				delIns(uinInfo.Uin)
			}(uinInfo)
			code, message, data = 200, "更新状态成功", nil
		}
	}
	c.Return(code, message, data)
}

// QrLogin 扫码登陆
func (ControllerUin) QrLogin(c *sinking_server.Context) {
	type Form struct {
		Uin  string `form:"uin" json:"uin" validate:"required,numeric" label:"QQ账号"`
		Step string `form:"step" json:"step" validate:"required,oneof=code query" label:"验证步骤"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if uin.GetIns().Count(map[string]string{"uin": form.Uin, "status": strconv.Itoa(uin.Error)}) > 0 {
		c.Error("该QQ已被系统锁定不可操作,请咨询管理员")
		return
	}
	user := web_user.GetIns().GetCurrentUserInfo(c.GetWebUserInfo().Id, false)
	if user.RealAuth != web_user.RealAuthSuccess {
		c.Error("该对接账户未完成实名认证,请完成实名认证")
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	if uinInfo.PlayerDay <= 0 && uinInfo.SignDay <= 0 {
		c.Error("该QQ未下单或已过期,无法操作")
		return
	}
	srv, err := server2.GetIns().FindById(uinInfo.ServerId)
	if err != nil || srv == nil || srv.Status != server2.Normal {
		c.Error("未查询到挂机地区信息,请更换挂机地区")
		return
	}
	key := "QrLogin_" + strconv.Itoa(user.Id) + "_" + form.Uin
	n, exists := util.Cache.Get(key)
	if exists {
		num, e := strconv.Atoi(n)
		if e == nil && num >= 1 {
			c.Error("请求频繁.请稍后重试")
			return
		}
	}
	util.Cache.Inc(key)
	util.Cache.SetExpire(key, int64(1000))
	if form.Step == "code" {
		changeProto := func(protos client.ClientProtocol) {
			//判断是否切换协议
			proto := oicq.GetOicqProtoCode(protos)
			if uinInfo.Proto != proto {
				err = uin.GetIns().UpdateByUin(uinInfo.Uin, map[string]interface{}{"proto": proto})
				if err != nil {
					return
				}
				uinInfo.Proto = proto
			}
		}
		//更换协议
		changeProto(client.AndroidWatch)
	}
	//获取单例
	cli, err := web_uin.GetOicqInsById(uinInfo.Id, true)
	cli.LoginTime = time.Now().Unix()
	if err != nil {
		c.Error("连接服务器失败")
		return
	}
	if form.Step == "code" {
		appConfig := app_config.GetIns().GetAppConfigs(constant.AppUinId, false)
		if appConfig[constant.AppUinLoginQrCode] == "1" {
			c.Error("系统已关闭扫码更新")
			return
		}
		checkMin := true
		if appConfig[constant.AppUinUpdateDate] != "" {
			t1, e1 := time.Parse("2006-01-02 15:04:05", appConfig[constant.AppUinUpdateDate])
			t2, e2 := time.Parse("2006-01-02 15:04:05", uinInfo.CreateTime.String())
			if e1 == nil && e2 == nil && t2.Before(t1) {
				checkMin = false
			}
		}
		if checkMin {
			minLevel := 0
			if appConfig[constant.AppUinUpdateLevel] != "" {
				num, e := strconv.Atoi(appConfig[constant.AppUinUpdateLevel])
				if e == nil {
					minLevel = num
				}
			}
			if minLevel > 0 {
				minDays := 0
				if appConfig[constant.AppUinUpdateDays] != "" {
					num, e := strconv.Atoi(appConfig[constant.AppUinUpdateDays])
					if e == nil {
						minDays = num
					}
				}
				u, _ := strconv.ParseInt(uinInfo.Uin, 10, 64)
				levelInfo := web_uin.GetGuestUinInfo(u)
				if levelInfo == nil {
					levelInfo = web_uin.GetGuestUinInfoBak(u)
				}
				if levelInfo != nil && levelInfo.Level <= minLevel && uinInfo.PlayerDay < minDays && uinInfo.SignDay < minDays {
					c.Error("您的剩余挂机天数低于" + strconv.Itoa(minDays) + "天。请续费至高于" + strconv.Itoa(minDays) + "天再更新状态")
					return
				}
			}
		}
		rsp, e := cli.Mobile.FetchQRCode()
		if e != nil || rsp.ImageData == nil {
			c.Error("获取二维码失败(-1)")
			return
		}
		imgObj, err := png.Decode(bytes.NewReader(rsp.ImageData))
		if err != nil {
			c.Error("解析二维码失败(-2)")
			return
		}
		bmp, err := gozxing.NewBinaryBitmapFromImage(imgObj)
		if err != nil {
			c.Error("解析二维码失败(-3)")
			return
		}
		qrReader := qrcode.NewQRCodeReader()
		result, err := qrReader.Decode(bmp, nil)
		if err != nil || result == nil || result.GetText() == "" {
			c.Error("解析二维码失败(-4)")
			return
		}
		cli.QrLoginSig = rsp.Sig
		c.SuccessWithData("获取二维码成功", result.GetText())
		return
	} else {
		if cli.QrLoginSig == nil {
			c.Error("获取状态失败,请先拉取二维码")
			return
		}
		s, e := cli.Mobile.QueryQRCodeStatus(cli.QrLoginSig)
		if e != nil {
			c.Error("获取状态失败,请重新拉取二维码")
			return
		}
		switch s.State {
		case client.QRCodeCanceled:
			c.SuccessWithData("获取状态成功", sinking_web.H{"status": 500, "message": "扫码被用户取消"})
			return
		case client.QRCodeTimeout:
			c.SuccessWithData("获取状态成功", sinking_web.H{"status": 501, "message": "二维码过期"})
			return
		case client.QRCodeWaitingForConfirm:
			c.SuccessWithData("获取状态成功", sinking_web.H{"status": 502, "message": "扫码成功, 请在手机端确认登录"})
			return
		case client.QRCodeImageFetch, client.QRCodeWaitingForScan:
			c.SuccessWithData("获取状态成功", sinking_web.H{"status": 503, "message": "等待扫码, 请在手机端确认登录"})
			return
		case client.QRCodeConfirmed:
			res, e := cli.Mobile.QRCodeLogin(s.LoginInfo)
			if e != nil {
				c.SuccessWithData("获取状态成功", sinking_web.H{"status": 504, "message": "扫码成功,发起登陆请求失败"})
				return
			}
			if !res.Success {
				//设备异常,生成新的设备信息
				if res.Error == client.DeviceVersionError {
					if uin.GetIns().UpdateById(uinInfo.Id, map[string]interface{}{"device": 100000 + rand.Intn(899999)}) == nil {
						q, e := strconv.ParseInt(uinInfo.Uin, 10, 64)
						if e == nil {
							oicq.DelOicqInstance(q) //删除实例
							oicq.DelTempOicqInstance(q)
						}
					}
					c.SuccessWithData("获取状态成功", sinking_web.H{"status": 505, "message": "登陆失败,请获取二维码重试"})
				} else {
					c.SuccessWithData("获取状态成功", sinking_web.H{"status": 505, "message": res.ErrorMessage})
				}
				return
			}
			if strconv.FormatInt(cli.Mobile.Uin, 10) != uinInfo.Uin {
				cli.QrLoginSig = nil
				c.SuccessWithData("获取状态成功", sinking_web.H{"status": 506, "message": "登陆账号不匹配"})
				return
			}
			token := cli.Mobile.GenToken()
			if web_uin.SaveLoginToken(uinInfo, token) == nil {
				go func(uinInfo *uin.AppUinJob) {
					//写入QQ操作记录
					_ = log.GetIns().Create(&model.AppUinLog{
						UserId:  uinInfo.UserId,
						AppId:   uinInfo.AppId,
						UinId:   uinInfo.Id,
						Uin:     uinInfo.Uin,
						Type:    log.Login,
						Title:   "更新状态",
						Content: "扫码更新账号状态",
					})
					_ = cli.Mobile.Init(true, false)
					q, e := strconv.ParseInt(uinInfo.Uin, 10, 64)
					if e == nil {
						cli.QrLoginSig = nil
						oicq.DelOicqInstance(q)
						oicq.DelTempOicqInstance(q)
					}
				}(uinInfo)
			}
			c.SuccessWithData("获取状态成功", sinking_web.H{"status": 200, "message": "更新状态成功"})
			return
		}
		c.Error("获取状态失败,请重新拉取二维码")
		return
	}
}

// Sms 短信验证
func (ControllerUin) Sms(c *sinking_server.Context) {
	type Form struct {
		Uin    string `form:"uin" json:"uin" validate:"required,numeric" label:"QQ账号"`
		Step   string `form:"step" json:"step" validate:"required,oneof=1 2 3" label:"验证步骤"`
		Rand   string `form:"rand" json:"rand" validate:"required,max=500" label:"登陆rand"`
		Ticket string `form:"ticket" json:"ticket" validate:"required,max=500" label:"登陆ticket"`
		Phone  string `form:"phone" json:"phone" validate:"omitempty,numeric" label:"手机号码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	if uinInfo.Status == uin.Error {
		c.Error("该QQ已被系统锁定不可操作,请咨询管理员")
		return
	}
	if uinInfo.PlayerDay <= 0 && uinInfo.SignDay <= 0 {
		c.Error("该QQ未下单或已过期,无法操作")
		return
	}
	srv, err := server2.GetIns().FindById(uinInfo.ServerId)
	if err != nil || srv == nil || srv.Status != server2.Normal {
		c.Error("未查询到挂机地区信息,请更换挂机地区")
		return
	}
	cli, err := web_uin.GetOicqInsById(uinInfo.Id, true)
	if err != nil {
		c.Error("连接服务器失败")
		return
	}
	ins := cli.Mobile
	key := "UinLoginSms" + uinInfo.Uin
	d := make(map[string]string)
	stars, exists := util.Cache.Get(key)
	if exists && stars != "" {
		d = str.NewStrTool().ToMap(stars)
	}
	data := make(map[string]interface{})
	switch form.Step {
	case "1":
		e := ins.SubmitEnvVerifyStepCheckTicket(form.Ticket, form.Rand)
		if e == nil {
			e = ins.SubmitEnvVerifyStepCheckRisk(form.Rand)
			if e == nil {
				e = ins.SubmitEnvVerifyStepQueryVerifyMethod(form.Rand)
				if e == nil {
					e2, area, phone := ins.SubmitEnvVerifyStepQueryPhone(form.Rand)
					e = e2
					if e == nil {
						data["area"] = area
						data["phone"] = phone
						util.Cache.SetWithExpire(key, str.NewStrTool().ToJson(map[string]string{
							"phone": phone,
							"area":  area,
						}), 600*1000) //写入缓存
					}
				}
			}
		}
		if e != nil {
			c.Error("初始化验证失败")
		} else {
			c.SuccessWithData("滑动验证码验证成功,请输入完整手机号", data)
		}
		break
	case "2":
		if d == nil || d["phone"] == "" || d["area"] == "" {
			c.Error("缺失参数,请重新验证")
			return
		}
		if form.Phone == "" {
			c.Error("手机号不能为空")
			return
		} else {
			e := ins.SubmitEnvVerifyStepCheckPhone(form.Rand, form.Phone, d["area"])
			if e == nil {
				e2, sendTo, sendSms := ins.SubmitEnvVerifyStepGetSms(form.Rand, d["phone"], d["area"])
				e = e2
				if e2 == nil {
					data["to"] = sendTo
					data["sms"] = sendSms
				}
			}
			if e != nil {
				c.Error("手机号验证失败")
			} else {
				c.SuccessWithData("手机号验证成功,请发送短信继续验证", data)
			}
		}
		break
	case "3":
		if d == nil || d["phone"] == "" || d["area"] == "" {
			c.Error("缺失参数,请重新验证")
			return
		}
		e, k := ins.SubmitEnvVerifyStepCheckSms(form.Rand, d["phone"], d["area"])
		if e != nil {
			c.Error("验证失败,请检查是否已发送短信")
			return
		}
		e = ins.SubmitEnvVerifyStepAuthPassWord(form.Rand, d["phone"], k)
		if e != nil {
			c.Error("上报结果失败,请联系管理员")
			return
		}
		c.Success("验证成功,请重新登录")
		break
	}
}

// Delete 删除QQ
func (ControllerUin) Delete(c *sinking_server.Context) {
	type Form struct {
		Uin string `form:"uin" json:"uin" validate:"required,numeric" label:"QQ账号"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	if uinInfo.Status == uin.Error {
		c.Error("该QQ已被系统锁定不可操作,请咨询管理员")
		return
	}
	err = uin.GetIns().DeleteById(uinInfo.Id)
	if err != nil {
		c.Error("删除QQ失败")
		return
	}
	go func(uinInfo *uin.AppUinJob) {
		if uinInfo.UinStatus == uin.Normal {
			_ = web_uin.ServerUinDelete(uinInfo.Id)
		}
		_ = server2.GetIns().UpdateById(uinInfo.ServerId, map[string]interface{}{
			"uin_num": gorm.Expr("`uin_num`-1"),
		})
	}(uinInfo)
	c.Success("删除QQ成功")
}

// Setting 功能设置
func (ControllerUin) Setting(c *sinking_server.Context) {
	type Form struct {
		Uin     string                     `json:"uin" validate:"required,numeric" label:"QQ账号"`
		Setting map[web_uin.JobName]string `json:"setting" validate:"required,min=1" label:"设置项"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	if uinInfo.PlayerDay <= 0 && uinInfo.SignDay <= 0 {
		c.Error("该QQ未下单或已过期,无法操作")
		return
	}
	err = web_uin.ChangeSetting(uinInfo.Id, form.Setting)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("操作成功")
}

// Job 执行任务
func (ControllerUin) Job(c *sinking_server.Context) {
	type Form struct {
		Uin string `json:"uin" validate:"required,numeric" label:"QQ账号"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	if uinInfo.Status == uin.Error {
		c.Error("该QQ已被系统锁定不可操作,请咨询管理员")
		return
	}
	if uinInfo.UinStatus != uin.Normal {
		c.Error("该QQ状态异常,请更新状态后重新执行")
		return
	}
	if uinInfo.PlayerDay <= 0 && uinInfo.SignDay <= 0 {
		c.Error("该QQ未下单或已过期,无法操作")
		return
	}
	key := constant.LockUserUinJob + uinInfo.Uin
	if util.Cache.Lock(key, constant.LockTimeUserUinJob) {
		go func(uinInfo *uin.AppUinJob) {
			producer.NewUinJobProducer().Build(uinInfo).Send()
			// 写入记录
			_ = log.GetIns().Create(&model.AppUinLog{
				UserId:  uinInfo.UserId,
				AppId:   uinInfo.AppId,
				UinId:   uinInfo.Id,
				Uin:     uinInfo.Uin,
				Type:    log.Run,
				Title:   "执行挂机",
				Content: "手动执行挂机",
			})
		}(&uin.AppUinJob{
			Id:         uinInfo.Id,
			UserId:     uinInfo.UserId,
			AppId:      uinInfo.AppId,
			ServerId:   uinInfo.ServerId,
			ServerKey:  uinInfo.ServerKey,
			Uin:        uinInfo.Uin,
			UinStatus:  uinInfo.UinStatus,
			SignDay:    uinInfo.SignDay,
			PlayerDay:  uinInfo.PlayerDay,
			RunDay:     uinInfo.RunDay,
			Timing:     uinInfo.Timing,
			Status:     uinInfo.Status,
			RunTime:    uinInfo.RunTime,
			CreateTime: uinInfo.CreateTime,
			UpdateTime: uinInfo.UpdateTime,
		})
		c.Success("提交执行任务指令成功,即将开始挂机")
	} else {
		c.Error("提交过于频繁,请等待60s后重试")
	}
}

// Info 账号信息
func (ControllerUin) Info(c *sinking_server.Context) {
	type Form struct {
		Uin string `json:"uin" validate:"required,numeric" label:"QQ账号"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	// QQ信息
	tempUinInfo := map[string]interface{}{
		"app_id":      uinInfo.AppId,
		"server_key":  uinInfo.ServerKey,
		"uin":         uinInfo.Uin,
		"uin_status":  uinInfo.UinStatus,
		"sign_day":    uinInfo.SignDay,
		"player_day":  uinInfo.PlayerDay,
		"run_day":     uinInfo.RunDay,
		"timing":      uinInfo.Timing,
		"status":      uinInfo.Status,
		"run_time":    uinInfo.RunTime,
		"create_time": uinInfo.CreateTime,
		"update_time": uinInfo.UpdateTime,
	}
	uinLevelInfo, e := web_uin.GetUinLevelInfo(uinInfo.Id)
	if e != nil {
		tempUinInfo["nick_name"] = uinInfo.Uin
		tempUinInfo["level"] = 0
		tempUinInfo["activate_day"] = "0"
		tempUinInfo["today_day"] = "0"
	} else {
		tempUinInfo["nick_name"] = uinLevelInfo.NickName
		tempUinInfo["level"] = uinLevelInfo.Level
		tempUinInfo["activate_day"] = uinLevelInfo.ActivateDay
		tempUinInfo["today_day"] = uinLevelInfo.TodayDay
	}
	// 任务状态
	jobStatus := make(map[web_uin.JobName]string)
	_ = json.Unmarshal([]byte(uinLevelInfo.SpeedInfo), &jobStatus)
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", uinLevelInfo.UpdateTime.String(), time.Local)
	todayUnix := str.NewTimeTool().GetZeroTime(time.Now()).Unix()
	newJobStatus := make(map[web_uin.JobName]string)
	sinking_web.Try(func() {
		for k, v := range web_uin.JobMap {
			for k1 := range v {
				days := 0
				if k == web_uin.Speed {
					days = uinInfo.PlayerDay
				} else {
					days = uinInfo.SignDay
				}
				if days <= 0 {
					newJobStatus[k1] = "3"
				} else {
					if jobStatus != nil {
						if todayUnix > stamp.Unix() || jobStatus[k1] == "" {
							newJobStatus[k1] = "0"
							continue
						} else {
							newJobStatus[k1] = jobStatus[k1]
						}
					}
				}
			}
		}
	}, nil)
	// 设置信息
	var tempSetting map[web_uin.JobName]string
	newSetting := make(map[web_uin.JobName]string)
	err = json.Unmarshal([]byte(uinInfo.Setting), &tempSetting)
	if err != nil {
		newSetting = web_uin.Setting
	} else {
		for k, v := range web_uin.Setting {
			if tempSetting[k] != "" {
				newSetting[k] = tempSetting[k]
			} else {
				newSetting[k] = v
			}
		}
	}
	c.SuccessWithData("获取成功", sinking_web.H{
		"uin":        tempUinInfo,
		"setting":    newSetting,
		"job_status": newJobStatus,
	})
}

// Heart 发送心跳
func (ControllerUin) Heart(c *sinking_server.Context) {
	type Form struct {
		Uin  string `json:"uin" validate:"required,numeric" label:"QQ账号"`
		Wait bool   `json:"wait" validate:"omitempty" label:"是否等待返回结果"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	if uinInfo.Status == uin.Error {
		c.Error("该QQ已被系统锁定不可操作,请咨询管理员")
		return
	}
	if uinInfo.PlayerDay <= 0 && uinInfo.SignDay <= 0 {
		c.Error("该QQ未下单或已过期,无法操作")
		return
	}
	if uinInfo.UinStatus != uin.Normal {
		c.Error("该QQ状态异常,请更新状态后重新执行")
		return
	}
	userConfigs := app_config.GetIns().GetUserConfigs(uinInfo.UserId, false)
	if userConfigs[constant.UserAuthCookies] != "1" && userConfigs[constant.UserAuthQrLogin] != "1" && userConfigs[constant.UserAuthDeviceLock] != "1" {
		c.Error("您的账户暂未开通此权限")
		return
	}
	//1.获取实例
	cli, err := web_uin.GetOicqInsById(uinInfo.Id, true)
	cli.LoginTime = time.Now().Unix()
	if err != nil {
		c.Error("连接服务器失败")
		return
	}
	//2.判断是否在线
	if !cli.Mobile.Online.Load() {
		//不在线则登陆账号
		e, l := cli.AutoLogin()
		if e != nil {
			c.Error("账号状态失效,请先更新状态")
			return
		}
		if !l.Success {
			if l.ErrorMessage != "" {
				c.Error(l.ErrorMessage)
			} else {
				c.Error("登录账号失败")
			}
			return
		}
	}
	//3.发送心跳
	if e := cli.Mobile.DoHeartbeat(form.Wait); e != nil {
		c.Error("发送心跳失败")
	} else {
		c.Success("发送心跳成功")
	}
}

// GetCookies 获取cookies信息
func (ControllerUin) GetCookies(c *sinking_server.Context) {
	type Form struct {
		Uin string `json:"uin" validate:"required,numeric" label:"QQ账号"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	if uinInfo.Status == uin.Error {
		c.Error("该QQ已被系统锁定不可操作,请咨询管理员")
		return
	}
	if uinInfo.PlayerDay <= 0 && uinInfo.SignDay <= 0 {
		c.Error("该QQ未下单或已过期,无法操作")
		return
	}
	if uinInfo.UinStatus != uin.Normal {
		c.Error("该QQ状态异常,请更新状态后重新执行")
		return
	}
	userConfigs := app_config.GetIns().GetUserConfigs(uinInfo.UserId, false)
	if userConfigs[constant.UserAuthCookies] != "1" {
		c.Error("您的账户暂未开通此权限")
		return
	}
	//1.获取实例
	cli, err := web_uin.GetOicqInsById(uinInfo.Id, true)
	cli.LoginTime = time.Now().Unix()
	if err != nil {
		c.Error("连接服务器失败")
		return
	}
	//2.判断是否在线
	if !cli.Mobile.Online.Load() {
		//不在线则登陆账号
		e, l := cli.AutoLogin()
		if e != nil {
			c.Error("账号状态失效,请先更新状态")
			return
		}
		if !l.Success {
			if l.ErrorMessage != "" {
				c.Error(l.ErrorMessage)
			} else {
				c.Error("登录账号失败")
			}
			return
		}
	}
	//3.获取cookies
	clientKey := cli.Mobile.GetClientKey()
	cookies := cli.Mobile.GetAllCookies()
	c.SuccessWithData("获取cookies成功", sinking_web.H{
		"client_key": clientKey,
		"cookies":    cookies,
	})
}

// AllowQrCode 授权登陆二维码
func (ControllerUin) AllowQrCode(c *sinking_server.Context) {
	type Form struct {
		Uin string `json:"uin" validate:"required,numeric" label:"QQ账号"`
		Pk  string `json:"pk" validate:"required,min=10,max=64" label:"二维码PK"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	if uinInfo.Status == uin.Error {
		c.Error("该QQ已被系统锁定不可操作,请咨询管理员")
		return
	}
	if uinInfo.PlayerDay <= 0 && uinInfo.SignDay <= 0 {
		c.Error("该QQ未下单或已过期,无法操作")
		return
	}
	if uinInfo.UinStatus != uin.Normal {
		c.Error("该QQ状态异常,请更新状态后重新执行")
		return
	}
	userConfigs := app_config.GetIns().GetUserConfigs(uinInfo.UserId, false)
	if userConfigs[constant.UserAuthQrLogin] != "1" {
		c.Error("您的账户暂未开通此权限")
		return
	}
	//1.获取实例
	cli, err := web_uin.GetOicqInsById(uinInfo.Id, true)
	cli.LoginTime = time.Now().Unix()
	if err != nil {
		c.Error("连接服务器失败")
		return
	}
	//2.判断是否在线
	if !cli.Mobile.Online.Load() {
		//不在线则登陆账号
		e, l := cli.AutoLogin()
		if e != nil {
			c.Error("账号状态失效,请先更新状态")
			return
		}
		if !l.Success {
			if l.ErrorMessage != "" {
				c.Error(l.ErrorMessage)
			} else {
				c.Error("登录账号失败")
			}
			return
		}
	}
	//3.扫码授权
	e := cli.QrLogin(form.Pk)
	if e != nil {
		c.Error("授权登录失败")
		return
	}
	c.Success("授权登录成功")
}

// AllowDeviceLock 授权登陆设备锁
func (ControllerUin) AllowDeviceLock(c *sinking_server.Context) {
	type Form struct {
		Uin       string `json:"uin" validate:"required,numeric" label:"QQ账号"`
		VerifyUrl string `json:"verify_url" validate:"required,url" label:"验证链接"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	if uinInfo.Status == uin.Error {
		c.Error("该QQ已被系统锁定不可操作,请咨询管理员")
		return
	}
	if uinInfo.PlayerDay <= 0 && uinInfo.SignDay <= 0 {
		c.Error("该QQ未下单或已过期,无法操作")
		return
	}
	if uinInfo.UinStatus != uin.Normal {
		c.Error("该QQ状态异常,请更新状态后重新执行")
		return
	}
	userConfigs := app_config.GetIns().GetUserConfigs(uinInfo.UserId, false)
	if userConfigs[constant.UserAuthDeviceLock] != "1" {
		c.Error("您的账户暂未开通此权限")
		return
	}
	//1.获取实例
	cli, err := web_uin.GetOicqInsById(uinInfo.Id, true)
	cli.LoginTime = time.Now().Unix()
	if err != nil {
		c.Error("连接服务器失败")
		return
	}
	//2.判断是否在线
	if !cli.Mobile.Online.Load() {
		//不在线则登陆账号
		e, l := cli.AutoLogin()
		if e != nil {
			c.Error("账号状态失效,请先更新状态")
			return
		}
		if !l.Success {
			if l.ErrorMessage != "" {
				c.Error(l.ErrorMessage)
			} else {
				c.Error("登录账号失败")
			}
			return
		}
	}
	//3.授权登陆
	e := cli.AllowDeviceLock(form.VerifyUrl)
	if e != nil {
		c.Error("授权登录失败")
		return
	}
	c.Success("授权登录成功")
}

// LevelInfo 账号等级信息
func (ControllerUin) LevelInfo(c *sinking_server.Context) {
	type Form struct {
		Uin string `json:"uin" validate:"required,numeric" label:"QQ账号"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	u, err := strconv.ParseInt(form.Uin, 10, 64)
	if err != nil || u <= 10000 || u >= 9999999999 {
		c.Error("获取信息失败")
		return
	}
	var info *web_uin.UinInfo
	key := "LevelInfo_" + form.Uin
	a, exists := util.Cache.Get(key)
	if exists {
		_ = json.Unmarshal([]byte(a), &info)
	}
	if info == nil {
		info = web_uin.GetGuestUinInfo(u)
		if info == nil {
			info = web_uin.GetGuestUinInfoBak(u)
		}
		if info != nil {
			bytes1, e := json.Marshal(info)
			if e == nil {
				util.Cache.SetWithExpire(key, string(bytes1), 3600*1000)
			}
		}
	}
	if info == nil {
		c.Error("获取信息失败")
		return
	}
	c.SuccessWithData("获取信息成功", info)
}

// Activity 活动列表
func (ControllerUin) Activity(c *sinking_server.Context) {
	appConfig := app_config.GetIns().GetAppConfigs(constant.AppUinId, false)
	// 首次购买活动
	firstBuy := 0
	firstBuyOpen := 0
	if appConfig[constant.AppUinFirstBuyGift] != "" {
		n, e := strconv.Atoi(appConfig[constant.AppUinFirstBuyGift])
		if e == nil && n > 0 {
			firstBuy = n
			firstBuyOpen = 1
		}
	}
	// 开通赠送活动
	buyOpen := 0
	if appConfig[constant.AppUinBuyActivityOpen] == "1" {
		buyOpen = 1
	}
	fun := func(value string, d int) int {
		if value == "" {
			return d
		}
		n, e := strconv.Atoi(value)
		if e == nil && n > 0 {
			return n
		}
		return d
	}
	//扫码加速活动
	testOpen := 0
	if appConfig[constant.AppUinSpeedTest] == "1" {
		testOpen = 1
	}
	//领取体验套装活动
	experienceOpen := 0
	if appConfig[constant.AppUinExperienceOpen] == "1" {
		experienceOpen = 1
	}
	c.SuccessWithData("获取成功", sinking_web.H{
		"test_speed":         testOpen,                                         //测试扫码加速
		"experience_open":    experienceOpen,                                   //领取体验套装活动
		"experience_days":    fun(appConfig[constant.AppUinExperienceDays], 0), //领取体验套装活动天数
		"first_buy_open":     firstBuyOpen,                                     //首次购买活动开关
		"first_buy_multiple": firstBuy,                                         //首次购买到账倍数
		"buy_open":           buyOpen,                                          //购买赠送活动开关
		"buy_month_gift":     fun(appConfig[constant.AppUinBuyMonthGift], 0),   //购买1月赠送天数
		"buy_quarter_gift":   fun(appConfig[constant.AppUinBuyQuarterGift], 0), //购买3月赠送天数
		"buy_half_gift":      fun(appConfig[constant.AppUinBuyHalfGift], 0),    //购买6月赠送天数
		"buy_year_gift":      fun(appConfig[constant.AppUinBuyYearGift], 0),    //购买12月赠送天数
	})
}

// Project 功能列表
func (ControllerUin) Project(c *sinking_server.Context) {
	c.SuccessWithData("获取成功", web_uin.Project)
}

// Self 账户信息
func (ControllerUin) Self(c *sinking_server.Context) {
	user := web_user.GetIns().GetCurrentUserInfo(c.GetWebUserInfo().Id, false)
	level, _ := web_uin.GetUserLevelWithCache(user.Id, false)
	userConfigs := app_config.GetIns().GetUserConfigs(user.Id, false)
	c.SuccessWithData("获取成功", sinking_web.H{
		"money":  user.Money,
		"level":  level,
		"expire": userConfigs[constant.UserUinLevelExpire],
	})
}

// Count 统计信息
func (ControllerUin) Count(c *sinking_server.Context) {
	u := c.GetWebUserInfo()
	data := make(map[string]int64)
	data["uin_num_all_normal"] = web_uin.GetAllNum(u.Id, 0)
	data["uin_num_all_un_normal"] = web_uin.GetAllNum(u.Id, 1)
	data["uin_num_today_normal"] = web_uin.GetTodayAddNum(u.Id, 0)
	data["uin_num_today_un_normal"] = web_uin.GetTodayAddNum(u.Id, 1)
	c.SuccessWithData("获取成功", data)
}

// Auth 权限信息
func (ControllerUin) Auth(c *sinking_server.Context) {
	user := web_user.GetIns().GetCurrentUserInfo(c.GetWebUserInfo().Id, false)
	userConfigs := app_config.GetIns().GetUserConfigs(user.Id, false)
	c.SuccessWithData("获取成功", sinking_web.H{
		"qr_login":    userConfigs[constant.UserAuthQrLogin] == "1",
		"device_lock": userConfigs[constant.UserAuthDeviceLock] == "1",
		"cookies":     userConfigs[constant.UserAuthCookies] == "1",
	})
}

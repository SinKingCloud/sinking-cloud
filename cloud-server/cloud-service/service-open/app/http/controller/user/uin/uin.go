package uin

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/command/queues/producer"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/app_config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/rpc/web_user_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/log"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_user"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/util/server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/util/uin_tool"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"gorm.io/gorm"
	"strconv"
)

type ControllerUin struct {
}

// List 列表
func (ControllerUin) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		AppId           string `form:"app_id" default:"" validate:"omitempty,numeric" label:"应用ID"`
		ServerKey       string `form:"server_key" default:"" validate:"omitempty,max=50" label:"地区标识"`
		UinStatus       string `form:"uin_status" default:"" validate:"omitempty,numeric" label:"挂机状态"`
		Uin             string `form:"uin" default:"" validate:"omitempty,numeric" label:"QQ账号"`
		Status          string `form:"status" default:"" validate:"omitempty,oneof=0 1 2" label:"QQ状态"`
		RunTimeStart    string `form:"run_time_start" default:"" validate:"-" label:"挂机起始时间"`
		RunTimeEnd      string `form:"run_time_end" default:"" validate:"-" label:"挂机结束时间"`
		CreateTimeStart string `form:"create_time_start" default:"" validate:"-" label:"创建起始时间"`
		CreateTimeEnd   string `form:"create_time_end" default:"" validate:"-" label:"创建结束时间"`
		UpdateTimeStart string `form:"update_time_start" default:"" validate:"-" label:"更新起始时间"`
		UpdateTimeEnd   string `form:"update_time_end" default:"" validate:"-" label:"更新结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
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
		rpc.WebUserLog.Create(c.GetWebUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), web_user_log.EventShow, "查看挂机", "查看挂机列表")
		c.SuccessWithData("获取挂机列表成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

// Buy 下单
func (ControllerUin) Buy(c *sinking_server.Context) {
	type Form struct {
		Uin       string `json:"uin" validate:"required,numeric" label:"QQ账号"`
		ServerKey string `json:"server_key" validate:"omitempty,max=10" label:"挂机地区"`
		AppId     int    `json:"app_id" validate:"omitempty,numeric,gte=0" label:"应用ID"`
		Pwd       string `json:"pwd" validate:"omitempty,len=32" label:"QQ密码"`
		Type      int    `json:"type" validate:"oneof=0 1 2" label:"下单类型"`
		Month     int    `json:"month" validate:"oneof=1 3 6 12" label:"下单月数"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
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
	//下单
	err = web_uin.Buy(c.GetWebUserInfo().Id, form.AppId, form.ServerKey, form.Uin, form.Pwd, form.Type, form.Month, true)
	if err != nil {
		c.Error(err.Error())
	} else {
		c.Success("下单成功")
	}
}

// ChangePassword 更改挂机QQ密码
func (ControllerUin) ChangePassword(c *sinking_server.Context) {
	type Form struct {
		Uin string `json:"uin" validate:"required,numeric" label:"QQ账号"`
		Pwd string `json:"pwd" validate:"required,len=32" label:"QQ密码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
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
	err = web_uin.ChangePassword(uinInfo.Id, form.Pwd)
	if err == nil {
		c.Success("更改挂机密码成功,请重新登录上线")
	} else {
		c.Error("更改挂机密码失败")
	}
}

// ChangePhone 更改挂机手机号
func (ControllerUin) ChangePhone(c *sinking_server.Context) {
	type Form struct {
		Uin   string `json:"uin" validate:"required,numeric" label:"QQ账号"`
		Phone string `json:"phone" validate:"required,len=11" label:"绑定手机"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
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
	err = web_uin.ChangePhone(uinInfo.Id, form.Phone)
	if err == nil {
		c.Success("更改挂机手机成功,请重新登录上线")
	} else {
		c.Error("更改挂机手机失败")
	}
}

// ChangeServer 更改挂机地区
func (ControllerUin) ChangeServer(c *sinking_server.Context) {
	type Form struct {
		Uin       string `json:"uin" validate:"required,numeric" label:"QQ账号"`
		ServerKey string `json:"server_key" validate:"required,max=32" label:"挂机地区"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
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
	err, change := web_uin.ChangeServer(uinInfo.Id, form.ServerKey)
	if err == nil {
		t := ""
		if change {
			t = "请重新登录上线"
		}
		c.Success("更改挂机地区成功" + t)
	} else {
		c.Error(err.Error())
	}
}

// ChangeTiming 修改挂机时间
func (ControllerUin) ChangeTiming(c *sinking_server.Context) {
	type Form struct {
		Uin    string `json:"uin" validate:"required,numeric" label:"QQ账号"`
		Timing int    `json:"timing" validate:"numeric,gt=0,lt=1440" label:"挂机时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
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

// ChangeStatus 修改挂机状态
func (ControllerUin) ChangeStatus(c *sinking_server.Context) {
	type Form struct {
		Uin    string `json:"uin" validate:"required,numeric" label:"QQ账号"`
		Status int    `json:"status" validate:"oneof=0 1" label:"挂机状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
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
		Uin       string `json:"uin" validate:"required,numeric" label:"QQ账号"`
		LoginType string `json:"login_type" validate:"oneof=pwd phone" label:"登陆方式"`
		Step      string `json:"step" validate:"oneof=normal captcha sms" label:"登陆步骤"`
		Rand      string `json:"rand" validate:"omitempty,max=500" label:"登陆rand"`
		Ticket    string `json:"ticket" validate:"omitempty,max=500" label:"登陆ticket"`
		Code      string `json:"code" validate:"omitempty,numeric,len=6" label:"短信验证码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
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
	q := uin_tool.NewUinTool(srv.Addr, srv.Proxy)
	var code int
	var e error
	var data string
	var verify string
	var loginUin string
	switch form.Step {
	case "normal":
		captcha := ""
		if form.LoginType == "phone" {
			//手机登录
			if uinInfo.LoginType != 1 {
				if uinInfo.Phone == "" {
					c.Error("登陆错误")
					return
				}
				err := web_uin.ChangePhone(uinInfo.Id, uinInfo.Phone)
				if err != nil {
					c.Error("登陆错误")
					return
				} else {
					uinInfo.LoginType = 1
				}
			}
		} else {
			//密码登陆
			if uinInfo.LoginType != 0 {
				err := web_uin.ChangePassword(uinInfo.Id, uinInfo.Password)
				if err != nil {
					c.Error("登陆错误")
					return
				} else {
					uinInfo.LoginType = 0
				}
			}
		}
		if uinInfo.LoginType == 1 {
			if form.Ticket == "" {
				c.Error("ticket不能为空")
				return
			}
			if form.Rand == "" {
				c.Error("rand不能为空")
				return
			}
			captcha = form.Rand + "/" + form.Ticket
		}
		code, e, data, loginUin = q.Login(uinInfo.Uin, uinInfo.Password, uinInfo.Phone, uinInfo.LoginType, captcha)
		if e != nil || code != uin_tool.Success {
			if code == uin_tool.NeedCaptcha {
				c.Return(400, e.Error(), data)
			} else if code == uin_tool.NeedCode {
				c.Return(401, e.Error(), nil)
			} else {
				c.Error(e.Error())
			}
			return
		}
		break
	case "captcha":
		if form.Ticket == "" {
			c.Error("ticket不能为空")
			return
		}
		verify = form.Ticket
		break
	case "sms":
		if form.Code == "" {
			c.Error("验证码不能为空")
			return
		}
		verify = form.Code
		break
	}
	if form.Step != "normal" {
		code, e, data, loginUin = q.VerifyLogin(uinInfo.Uin, verify)
		if e != nil || code != uin_tool.Success {
			if code == uin_tool.NeedCode {
				c.Return(401, e.Error(), nil)
			} else if code == uin_tool.NeedCaptcha {
				c.Return(400, e.Error(), data)
			} else {
				c.Error(e.Error())
			}
			return
		}
	}
	if loginUin != uinInfo.Uin {
		c.Error("该手机号绑定的QQ与当前挂机账号不一致。")
		return
	}
	go func(uinInfo *model.AppUin, status int) {
		//更改QQ状态
		_ = uin.GetIns().UpdateById(uinInfo.Id, map[string]interface{}{"uin_status": status})
		//写入记录
		_ = log.GetIns().Create(&model.AppUinLog{
			UserId:  uinInfo.UserId,
			AppId:   uinInfo.AppId,
			UinId:   uinInfo.Id,
			Uin:     uinInfo.Uin,
			Type:    log.Login,
			Title:   "登陆账号",
			Content: "更新账号状态",
		})
	}(uinInfo, 0)
	c.Success("登陆成功")
}

// Delete 删除QQ
func (ControllerUin) Delete(c *sinking_server.Context) {
	type Form struct {
		Uin string `json:"uin" validate:"required,numeric" label:"QQ账号"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
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
	go func(uinInfo *model.AppUin) {
		_ = server2.GetIns().UpdateById(uinInfo.ServerId, map[string]interface{}{
			"uin_num": gorm.Expr("`uin_num`-1"),
		})
	}(uinInfo)
	c.Success("删除QQ成功")
}

// Setting 功能设置
func (ControllerUin) Setting(c *sinking_server.Context) {
	type Form struct {
		Uin     string            `json:"uin" validate:"required,numeric" label:"QQ账号"`
		Setting map[string]string `json:"setting" validate:"required,min=1" label:"设置项"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
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
	if ok, msg := c.ValidatorJson(form); !ok {
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
		go func(uinInfo *model.AppUin) {
			producer.NewUinJobProducer().Build(uinInfo).Send()
			//写入记录
			_ = log.GetIns().Create(&model.AppUinLog{
				UserId:  uinInfo.UserId,
				AppId:   uinInfo.AppId,
				UinId:   uinInfo.Id,
				Uin:     uinInfo.Uin,
				Type:    log.Run,
				Title:   "执行挂机",
				Content: "手动执行挂机",
			})
		}(uinInfo)
		c.Success("提交执行任务指令成功,即将开始挂机")
	} else {
		c.Error("提交过于频繁,请五分钟后在重试")
	}
}

// Info 账号信息
func (ControllerUin) Info(c *sinking_server.Context) {
	type Form struct {
		Uin string `json:"uin" validate:"required,numeric" label:"QQ账号"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	tempInfo := &uin.AppUin{
		AppId:      uinInfo.AppId,
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
	}
	var tempSetting map[string]string
	err = json.Unmarshal([]byte(uinInfo.Setting), &tempSetting)
	if err != nil {
		tempSetting = web_uin.Setting
	} else {
		for k, v := range web_uin.Setting {
			if tempSetting[k] == "" {
				tempSetting[k] = v
			}
		}
	}
	c.SuccessWithData("获取成功", sinking_web.H{
		"uin":     tempInfo,
		"setting": tempSetting,
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

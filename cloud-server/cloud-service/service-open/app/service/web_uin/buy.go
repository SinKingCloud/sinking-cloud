package web_uin

import (
	"errors"
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/app_config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/order"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/server"
	uin2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_user"
	"gorm.io/gorm"
	"strconv"
)

// Buy 业务下单
func Buy(userId int, appId int, serverKey string, uin string, pwd string, types int, month int, cost bool) error {
	key := constant.LockUinBuy + strconv.Itoa(userId) + "_" + uin
	//上锁，避免同时下单太多
	if !util.Cache.Lock(key, constant.LockTimeUinBuy) {
		return errors.New("系统繁忙,请稍后下单")
	}
	defer util.Cache.UnLock(key) //释放并发锁
	user := web_user.GetIns().GetCurrentUserInfo(userId, true)
	if user == nil {
		return errors.New("获取用户信息失败")
	}
	if appId != 0 {
		app, err := web_app.GetIns().FindByIdWithCache(appId, false)
		if err != nil || app == nil || app.UserId != user.Id {
			return errors.New("未获取到应用信息")
		}
	}
	level, err := GetUserLevelWithCache(user.Id, false)
	if err != nil || level == nil {
		return errors.New("未获取到折扣价格信息，请联系管理员")
	}
	//获取价格
	price := GetPrice(level, types, month)
	if cost {
		//判断用户余额是否足够
		if user.Money < price {
			return errors.New("账户余额不足,请充值")
		}
		//扣款
		if price > 0 && !rpc.WebUser.Money(userId, 0, fmt.Sprintf("%.2f", price)) {
			return errors.New("扣款失败，请联系管理员")
		}
	}
	uinInfo, err := uin2.GetIns().FindByUserIdWithUin(user.Id, uin)
	if err != nil || uinInfo == nil {
		server, err2 := GetServer(serverKey)
		if err2 != nil || server == nil {
			err = errors.New("未获取到可用节点,请更换挂机节点")
		} else {
			//写入qq记录
			uinLog := &model.AppUin{
				UserId:    user.Id,
				AppId:     appId,
				ServerId:  server.Id,
				ServerKey: server.Key,
				Uin:       uin,
				Password:  pwd,
				UinStatus: 1,
				Device:    0,
				SignDay:   0,
				PlayerDay: 0,
				RunDay:    0,
				Setting:   str.NewJsonTool(Setting).ToString(),
				Timing:    30,
				Status:    0,
			}
			err = uin2.GetIns().Create(uinLog)
			if err == nil {
				_ = server2.GetIns().UpdateById(server.Id, map[string]interface{}{
					"uin_num": gorm.Expr("`uin_num`+1"),
				})
				uinInfo, err = uin2.GetIns().FindByUserIdWithUin(user.Id, uin)
			}
		}
	}
	//下单错误异常
	var orderErr error
	//是否首次下单(1:否/2:是)
	firstBuy := 1
	//赠送天数
	giftDays := 0
	//写入qq记录失败,退回余额
	if err != nil || uinInfo == nil {
		orderErr = err
	} else {
		data := make(map[string]interface{})
		//应用ID变动
		if appId != 0 && appId != uinInfo.AppId {
			data["app_id"] = appId
		}
		//判断是否首次下单,首次下单赠送天数
		if log.GetIns().Count(map[string]string{"user_id": strconv.Itoa(userId), "uin": uinInfo.Uin}) <= 0 {
			appConfig := app_config.GetIns().GetAppConfigs(constant.AppUinId, false)
			if appConfig[constant.AppUinFirstBuyGift] != "" {
				n, e := strconv.Atoi(appConfig[constant.AppUinFirstBuyGift])
				if e == nil && n > 0 {
					firstBuy = n
				}
			}
		}
		//赠送天数
		giftDays = GetGiftDays(month)
		//实际到账天数
		days := month*30*firstBuy + giftDays
		//更改到期时间
		if types == 0 || types == 1 {
			data["player_day"] = gorm.Expr("`player_day`+" + strconv.Itoa(days))
		}
		if types == 0 || types == 2 {
			data["sign_day"] = gorm.Expr("`sign_day`+" + strconv.Itoa(days))
		}
		//密码变动
		if pwd != "" && pwd != uinInfo.Password {
			_ = ChangePassword(uinInfo.Id, pwd)
		}
		//挂机地点变动
		if serverKey != "" && serverKey != uinInfo.ServerKey {
			_ = ChangeServer(uinInfo.Id, serverKey)
		}
		err = uin2.GetIns().UpdateById(uinInfo.Id, data)
		if err != nil {
			orderErr = errors.New("下单失败，请联系管理员")
		}
	}
	if orderErr != nil {
		if cost && price > 0 {
			rpc.WebUser.Money(userId, 1, fmt.Sprintf("%.2f", price))
		}
		return orderErr
	}
	//写入消费及QQ下单记录
	go func() {
		typeName := "全套套装"
		if types == 1 {
			typeName = "加速套装"
		} else if types == 2 {
			typeName = "签到套装"
		}
		t := "在线"
		if !cost {
			t = "后台"
		}
		content := t + "下单" + strconv.Itoa(month) + "月" + typeName + "。"
		if firstBuy > 1 {
			content += "首次下单" + strconv.Itoa(firstBuy) + "倍到账。"
		}
		if giftDays >= 1 {
			content += "活动优惠赠送" + strconv.Itoa(giftDays) + "天时长。"
		}
		_ = log.GetIns().Create(&model.AppUinLog{
			UserId:  user.Id,
			AppId:   appId,
			UinId:   uinInfo.Id,
			Uin:     uinInfo.Uin,
			Type:    log.Buy,
			Title:   "在线下单",
			Content: content,
		})
		_ = order.GetIns().Create(&model.AppUinOrder{
			UserId: user.Id,
			AppId:  appId,
			UinId:  uinInfo.Id,
			Type:   types,
			Money:  price,
			Name:   content,
			Status: 0,
		})
		if cost {
			rpc.WebUserPayLog.Create(user.Id, 1, 1, price, "在线消费", "为QQ"+uinInfo.Uin+"在线下单"+strconv.Itoa(month)+"月"+typeName)
		}
		//清理用户缓存
		web_user.GetIns().Clear(user.Id)
	}()
	return nil
}

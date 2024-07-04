package web_uin

import (
	"errors"
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/app_config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin/order"
	uin2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin/uin"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// Refund 退单
func Refund(tradeNo string) error {
	// todo uin_order新增trade_no refund_money days等字段
	if tradeNo == "" {
		return errors.New("订单号不能为空")
	}
	//设置并发锁
	key := constant.LockUinRefund + tradeNo
	if !util.Cache.Lock(key, constant.LockTimeUinRefund) {
		return errors.New("系统繁忙,请稍后操作")
	}
	defer util.Cache.UnLock(key)
	//获取退单设置
	appConfig := app_config.GetIns().GetAppConfigs(constant.AppUinId, false)
	if appConfig[constant.AppUinRefundOpen] != "1" {
		return errors.New("退单功能未开启")
	}
	day, err := strconv.Atoi(appConfig[constant.AppUinRefundDay])
	if err != nil {
		day = 0
	}
	freeDay, err := strconv.Atoi(appConfig[constant.AppUinRefundFree])
	if err != nil {
		freeDay = 0
	}
	//根据订单号查找订单信息
	orderInfo, err := order.GetIns().FindByTradeNo(tradeNo)
	if err != nil || orderInfo == nil {
		return errors.New("未查找到订单记录")
	}
	//判断是否已退单过
	if orderInfo.Status != int(order.Success) {
		return errors.New("该订单已退订过,不可重复退订")
	}
	//判断是否已过退单日期
	if day > 0 {
		stamp, e := time.ParseInLocation("2006-01-02 15:04:05", orderInfo.CreateTime.String(), time.Local)
		if e != nil || (stamp.Unix()+int64(day)*86400) < time.Now().Unix() {
			return errors.New("改订单已超过最大可退单时效(自下单日起" + strconv.Itoa(day) + "天可退单)")
		}
	}
	isFree := false
	if freeDay > 0 {
		stamp, e := time.ParseInLocation("2006-01-02 15:04:05", orderInfo.CreateTime.String(), time.Local)
		if e != nil || (stamp.Unix()+int64(freeDay)*86400) > time.Now().Unix() {
			isFree = true
		}
	}
	//计算价格
	money := orderInfo.Money
	if !isFree {
		point, err := strconv.ParseFloat(appConfig[constant.AppUinRefundPoint], 64)
		if err != nil {
			point = 0
		}
		min, err := strconv.ParseFloat(appConfig[constant.AppUinRefundMin], 64)
		if err != nil {
			min = 0
		}
		max, err := strconv.ParseFloat(appConfig[constant.AppUinRefundMax], 64)
		if err != nil {
			max = 0
		}
		if point > 0 {
			costMoney := orderInfo.Money / 100 * point
			if min > 0 && costMoney < min {
				costMoney = min
			}
			if max > 0 && costMoney > max {
				costMoney = max
			}
			money = orderInfo.Money - costMoney
		}
	}
	//退款至账户
	if money > 0 {
		if !rpc.WebUser.Money(orderInfo.UserId, 1, fmt.Sprintf("%.2f", money)) {
			return errors.New("退单失败,请联系管理员(-1)")
		}
		rpc.WebUserPayLog.Create(orderInfo.UserId, 0, 2, money, "在线退单", "订单号"+orderInfo.TradeNo+"在线退单,订单"+fmt.Sprintf("%.2f", orderInfo.Money)+"元,手续费"+fmt.Sprintf("%.2f", orderInfo.Money-money)+"元")
	}
	//更改订单状态
	err = order.GetIns().UpdateById(orderInfo.Id, map[string]interface{}{"status": order.Refund, "refund_money": money})
	if err != nil {
		return errors.New("退单失败,请联系管理员(-2)")
	}
	//扣除天数
	if orderInfo.Days > 0 {
		uinInfo, e := uin2.GetIns().FindByIdWithAppUin(orderInfo.UinId)
		if e == nil && uinInfo != nil {
			data := make(map[string]interface{})
			if orderInfo.Type == int(order.All) || orderInfo.Type == int(order.Player) {
				if uinInfo.PlayerDay-orderInfo.Days <= 0 {
					data["player_day"] = 0
				} else {
					data["player_day"] = gorm.Expr("`player_day`-" + strconv.Itoa(orderInfo.Days))
				}
			}
			if orderInfo.Type == int(order.All) || orderInfo.Type == int(order.Sign) {
				if uinInfo.SignDay-orderInfo.Days <= 0 {
					data["sign_day"] = 0
				} else {
					data["sign_day"] = gorm.Expr("`sign_day`-" + strconv.Itoa(orderInfo.Days))
				}
			}
			_ = uin2.GetIns().UpdateById(orderInfo.UinId, data)
			_ = log.GetIns().Create(&model.AppUinLog{
				UserId:  orderInfo.UserId,
				AppId:   orderInfo.AppId,
				UinId:   orderInfo.UinId,
				Uin:     uinInfo.Uin,
				Type:    log.Refund,
				Title:   "在线退单",
				Content: "在线申请订单号" + orderInfo.TradeNo + "退单",
			})
		}
	}
	return nil
}

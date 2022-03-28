package web_uin

import (
	"errors"
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/order"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/server"
	uin2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_user"
	"gorm.io/gorm"
	"strconv"
)

// Buy 业务下单
func Buy(userId int, appId int, serverKey string, uin string, pwd string, types int, month int) error {
	user := web_user.GetIns().GetCurrentUserInfo(userId, false)
	if user == nil {
		return errors.New("获取用户资料失败")
	}
	if appId != 0 {
		app, err := web_app.GetIns().FindByIdWithCache(appId, false)
		if err != nil || app == nil || app.UserId != appId {
			return errors.New("未获取到应用信息")
		}
	}
	level, err := GetUserLevelWithCache(user.Id, false)
	if err != nil || level == nil {
		return errors.New("未获取到折扣价格信息，请联系管理员")
	}
	//获取价格
	price := GetPrice(level, types, month)
	//判断用户余额是否足够
	if user.Money < price {
		return errors.New("账户余额不足,请充值")
	}
	//扣款
	if price > 0 && !rpc.WebUser.Money(userId, 0, fmt.Sprintf("%.2f", price)) {
		return errors.New("扣款失败，请联系管理员")
	}
	uinInfo, err := uin2.GetIns().FindByAppIdWithUin(user.Id, uin)
	if err != nil || uinInfo == nil {
		server, err2 := GetServer(serverKey)
		if err2 != nil || server == nil {
			return errors.New("未获取到可用节点,请更换挂机节点")
		}
		//写入qq记录
		uinLog := &model.AppUin{
			UserId:    user.Id,
			AppId:     appId,
			ServerId:  server.Id,
			ServerKey: server.Key,
			Uin:       uin,
			Password:  pwd,
			UinStatus: 1,
			SignDay:   0,
			PlayerDay: 0,
			RunDay:    0,
			Setting:   str.NewJsonTool(Setting).ToString(),
			Timing:    0,
			Status:    0,
		}
		err3 := uin2.GetIns().Create(uinLog)
		if err3 == nil {
			go func() {
				_ = server2.GetIns().UpdateById(server.Id, map[string]interface{}{
					"uin_num": gorm.Expr("`uin_num`+1"),
				})
			}()
			uinInfo, err = uin2.GetIns().FindByAppIdWithUin(user.Id, uin)
		}
	}
	var orderErr error
	//写入qq记录失败,退回余额
	if err != nil || uinInfo == nil {
		orderErr = errors.New("写入QQ记录失败，请联系管理员")
	} else {
		//更新到期时间
		data := make(map[string]interface{})
		if types == 0 || types == 1 {
			data["player_day"] = gorm.Expr("`player_day`+" + strconv.Itoa(month*30))
		}
		if types == 0 || types == 2 {
			data["sign_day"] = gorm.Expr("`sign_day`+" + strconv.Itoa(month*30))
		}
		if pwd != uinInfo.Password {
			data["password"] = pwd
			data["uin_status"] = 1
			go func() {
				_ = log.GetIns().Create(&model.AppUinLog{
					UserId:  user.Id,
					AppId:   appId,
					UinId:   uinInfo.Id,
					Title:   "修改密码",
					Content: "修改挂机密码",
				})
			}()
		}
		if serverKey != uinInfo.ServerKey {
			server, err := GetServer(serverKey)
			if err != nil || server == nil {
				return errors.New("未获取到可用节点,请更换挂机节点")
			}
			data["server_key"] = server.Key
			data["server_id"] = server.Id
			data["uin_status"] = 1
			go func() {
				//写入QQ操作记录
				_ = log.GetIns().Create(&model.AppUinLog{
					UserId:  user.Id,
					AppId:   appId,
					UinId:   uinInfo.Id,
					Title:   "挂机地点变动",
					Content: "挂机地点变动为" + server.Name,
				})
				_ = server2.GetIns().UpdateById(uinInfo.ServerId, map[string]interface{}{
					"uin_num": gorm.Expr("`uin_num`-1"),
				})
				_ = server2.GetIns().UpdateById(server.Id, map[string]interface{}{
					"uin_num": gorm.Expr("`uin_num`+1"),
				})
			}()
		}
		err = uin2.GetIns().UpdateById(uinInfo.Id, data)
		if err != nil {
			orderErr = errors.New("下单失败，请联系管理员")
		}
	}
	if orderErr != nil {
		if price > 0 {
			rpc.WebUser.Money(userId, 1, fmt.Sprintf("%.2f", price))
		}
		return orderErr
	}
	//写入消费及QQ下单记录
	go func() {
		typeName := "全套代挂"
		if types == 1 {
			typeName = "加速套装"
		} else if types == 2 {
			typeName = "签到套装"
		}
		_ = log.GetIns().Create(&model.AppUinLog{
			UserId:  user.Id,
			AppId:   appId,
			UinId:   uinInfo.Id,
			Title:   "在线下单",
			Content: "在线下单" + strconv.Itoa(month) + "月" + typeName,
		})
		_ = order.GetIns().Create(&model.AppUinOrder{
			UserId: user.Id,
			AppId:  appId,
			UinId:  uinInfo.Id,
			Type:   types,
			Money:  price,
			Name:   "在线下单" + strconv.Itoa(month) + "月" + typeName,
			Status: 0,
		})
		rpc.WebUserPayLog.Create(user.Id, 1, 1, price, "在线消费", "为QQ"+uinInfo.Uin+"在线下单"+strconv.Itoa(month)+"月"+typeName)
		//清理用户缓存
		web_user.GetIns().Clear(user.Id)
	}()
	return nil
}

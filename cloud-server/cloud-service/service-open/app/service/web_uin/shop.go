package web_uin

import (
	"errors"
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/constant"
	model2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/app_config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/rpc"
	level2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/level"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_user"
	"strconv"
	"time"
)

// BuyLevel 开通折扣等级
func BuyLevel(userId int, levelId int, month int, cost bool) error {
	level, err := level2.GetIns().FindById(levelId)
	if err != nil || level.Status != 0 {
		return errors.New("获取折扣等级信息失败")
	}
	if level.Price <= 0 {
		return errors.New("免费等级无需开通,将会自动结转")
	}
	user := web_user.GetIns().GetCurrentUserInfo(userId, true)
	if user == nil {
		return errors.New("获取用户信息失败")
	}
	//获取价格
	price := float64(month) * level.Price
	if cost && price > 0 {
		//判断账户余额
		if user.Money < price {
			return errors.New("账户余额不足")
		}
		//扣款
		if price > 0 && !rpc.WebUser.Money(user.Id, 0, fmt.Sprintf("%.2f", price)) {
			return errors.New("扣款失败，请联系管理员")
		}
	}
	//写入用户等级
	isNew := false
	userLevel, err := GetUserLevelWithCache(user.Id, false)
	if err != nil || userLevel == nil {
		//直接写入等级
		isNew = true
	} else {
		if userLevel.Id != level.Id {
			isNew = true
		}
	}
	//判断是否续费等级
	expireTimeStr := time.Now().Format("2006-01-02 15:04:05")
	if !isNew {
		userConfigs := app_config.GetIns().GetUserConfigs(user.Id, false)
		if userConfigs[constant.UserUinLevelExpire] != "" {
			expireTimeStr = userConfigs[constant.UserUinLevelExpire]
		}
	}
	//计算到期时间
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", expireTimeStr, time.Local)
	//写入到期时间
	app_config.GetIns().SetUserConfigs(user.Id, map[string]string{
		constant.UserUinLevelExpire: stamp.AddDate(0, month, 0).Format("2006-01-02 15:04:05"),
		constant.UserUinLevel:       strconv.Itoa(level.Id),
	})
	//刷新用户配置
	app_config.GetIns().ClearUserConfig(user.Id)
	go func(user *model.WebUser, level *model2.AppUinLevel, cost bool, price float64) {
		name := "自主"
		if !cost {
			price = 0
			name = "管理员"
		}
		//写入用户余额记录
		rpc.WebUserPayLog.Create(user.Id, 1, 1, price, "在线消费", name+"开通"+strconv.Itoa(month)+"月"+level.Name+"折扣等级")
		//刷新用户折扣等级
		_, _ = GetUserLevelWithCache(user.Id, true)
	}(user, level, cost, price)
	return nil
}

package web_uin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/app_config"
	"strconv"
	"time"
)

// GetUserLevel 获取用户折扣等级
func GetUserLevel(userId int) (data *model.AppUinLevel) {
	//判断用户是否开通了等级
	userConfigs := app_config.GetIns().GetUserConfigs(userId, true)
	var levelId int
	var err error
	isExpire := true //是否到期
	if userConfigs[constant.UserUinLevel] != "" {
		//已设定等级
		levelId, err = strconv.Atoi(userConfigs["uin.level"])
		if err != nil {
			levelId = 0
		} else {
			//判断等级是否到期
			if userConfigs[constant.UserUinLevelExpire] != "" {
				stamp, err := time.ParseInLocation("2006-01-02 15:04:05", userConfigs["uin.level.expire"], time.Local)
				if err != nil {
					levelId = 0
				} else {
					if time.Now().Unix() > stamp.Unix() {
						//已过期
						levelId = 0
					} else {
						isExpire = false
					}
				}
			} else {
				levelId = 0
			}
		}
	}
	if levelId != 0 {
		//直接返回设置等级
		err = util.Database.Db.Where("`id` = ?", levelId).First(&data).Error
		if err != nil {
			data = nil
		}
	} else {
		//根据消费金额计算等级
		err = util.Database.Db.Where("`month_spend`<=?", fmt.Sprintf("%.2f", GetUserMonthSpend(userId))).Order("id desc").First(&data).Error
		if err != nil {
			data = nil
		}
	}
	//写入用户配置
	if data != nil {
		d := map[string]string{
			constant.UserUinLevel: strconv.Itoa(data.Id),
		}
		if isExpire {
			d[constant.UserUinLevelExpire] = str.NewTimeTool().GetFirstDateOfMonth(time.Now().AddDate(0, +1, 0)).Format("2006-01-02 15:04:05")
		}
		app_config.GetIns().SetUserConfigs(userId, d)
		//刷新用户
		app_config.GetIns().ClearUserConfig(userId)
	}
	if data == nil {
		err = util.Database.Db.Order("id asc").First(&data).Error
		if err != nil {
			data = nil
		}
	}
	return
}

// GetUserMonthSpend 获取用户上月消费
func GetUserMonthSpend(userId int) float64 {
	var money float64
	err := util.Database.Db.Model(&model.AppUinOrder{}).
		Where("`user_id`=? and `status`=0 and `create_time` between ? and ?", userId,
			str.NewTimeTool().GetFirstDateOfMonth(time.Now().AddDate(0, -1, 0)).Format("2006-01-02 15:04:05"),
			str.NewTimeTool().GetLastDateOfMonth(time.Now().AddDate(0, -1, 0)).Format("2006-01-02 15:04:05")).
		Pluck("COALESCE(SUM(`money`),0) as money", &money).Error
	if err != nil {
		money = 0
	}
	return money
}

// GetUserLevelWithCache 获取用户折扣等级(缓存)
func GetUserLevelWithCache(userId int, refresh bool) (data *model.AppUinLevel, err error) {
	key := "UinUserLevel_" + time.Now().Format("2006-01-02") + strconv.Itoa(userId)
	if refresh {
		data = GetUserLevel(userId)
		if data == nil {
			err = errors.New("not found")
			return
		}
		bytes, err := json.Marshal(data)
		if err == nil {
			util.Cache.SetWithExpire(key, string(bytes), 86400*1000)
		}
	} else {
		get, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(get), &data)
		} else {
			data = GetUserLevel(userId)
			if data == nil {
				err = errors.New("not found")
				return
			}
			bytes, err := json.Marshal(data)
			if err == nil {
				util.Cache.SetWithExpire(key, string(bytes), 86400*1000)
			}
		}
	}
	return
}

// GetPrice 获取价格
func GetPrice(data *model.AppUinLevel, types int, month int) (price float64) {
	type Price struct {
		Month   float64
		Quarter float64
		Half    float64
		Year    float64
		Forever float64
	}
	f := func(data *model.AppUinLevel, types int) *Price {
		var price *Price
		switch types {
		case 0:
			//全套加速价格
			price = &Price{
				Month:   data.MonthPrice,
				Quarter: data.QuarterPrice,
				Half:    data.HalfPrice,
				Year:    data.YearPrice,
				Forever: data.ForeverPrice,
			}
			break
		case 1:
			//等价加速价格
			price = &Price{
				Month:   data.PlayerMonthPrice,
				Quarter: data.PlayerQuarterPrice,
				Half:    data.PlayerHalfPrice,
				Year:    data.PlayerYearPrice,
				Forever: data.PlayerForeverPrice,
			}
			break
		default:
			//扩展加速价格
			price = &Price{
				Month:   data.SignMonthPrice,
				Quarter: data.SignQuarterPrice,
				Half:    data.SignHalfPrice,
				Year:    data.SignYearPrice,
				Forever: data.SignForeverPrice,
			}
			break
		}
		return price
	}
	switch month {
	case 1:
		price = f(data, types).Month
		break
	case 3:
		price = f(data, types).Quarter
		break
	case 6:
		price = f(data, types).Half
		break
	case 12:
		price = f(data, types).Year
		break
	case 999:
		price = f(data, types).Forever
		break
	}
	return
}

// GetGiftDays 获取赠送天数
func GetGiftDays(month int) int {
	appConfig := app_config.GetIns().GetAppConfigs(constant.AppUinId, false)
	//判断是否开启活动
	if appConfig[constant.AppUinBuyActivityOpen] != "1" {
		return 0
	}
	fun := func(appConfig map[string]string, key string) int {
		days := 0
		if appConfig[key] != "" {
			n, e := strconv.Atoi(appConfig[key])
			if e == nil && n > 0 {
				days = n
			}
		}
		return days
	}
	key := ""
	switch month {
	case 1:
		key = constant.AppUinBuyMonthGift
		break
	case 3:
		key = constant.AppUinBuyQuarterGift
		break
	case 6:
		key = constant.AppUinBuyHalfGift
		break
	case 12:
		key = constant.AppUinBuyYearGift
		break
	}
	if key == "" {
		return 0
	}
	return fun(appConfig, key)
}

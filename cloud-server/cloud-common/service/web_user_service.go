package service

import (
	"encoding/json"
	"errors"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/web_user_config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

// webUserService 单例对象
type webUserService struct {
}

// userObj 单例对象
var (
	webUserObj     *webUserService
	onceWebUserObj sync.Once
)

// GetWebUserServiceIns 获取单例
func GetWebUserServiceIns() *webUserService {
	onceWebUserObj.Do(func() {
		webUserObj = &webUserService{}
	})
	return webUserObj
}

// ClearCache 清理缓存
func (webUserService) ClearCache(id int) {
	util.Cache.Delete(util2.CacheNameWithWebUser + strconv.Itoa(id))
}

// FindById 查询用户信息(根据用户ID)
func (webUserService) FindById(id int) (user *model.WebUser, err error) {
	err = util.Database.Db.Where("id = ?", id).First(&user).Error
	return
}

// UpdateById 通过ID更新
func (w webUserService) UpdateById(id int, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.WebUser{}).Where("`id` = ? ", id).Updates(data).Error
	if err == nil {
		w.ClearCache(id)
	}
	return
}

// FindByIdWithCache 查询缓存用户信息(根据用户ID)
func (w webUserService) FindByIdWithCache(id int, refresh bool) (user *model.WebUser, err error) {
	key := util2.CacheNameWithWebUser + strconv.Itoa(id)
	if refresh {
		user, err = w.FindById(id)
		if err != nil {
			return
		}
		data, err := json.Marshal(user)
		if err == nil {
			util.Cache.SetWithExpire(key, string(data), util2.CacheTimeWithWebUser)
		}
	} else {
		data, exists := util.Cache.Get(key)
		if exists {
			err = json.Unmarshal([]byte(data), &user)
		} else {
			user, err = w.FindById(id)
			if err != nil {
				return
			}
			data, err := json.Marshal(user)
			if err == nil {
				util.Cache.SetWithExpire(key, string(data), util2.CacheTimeWithWebUser)
			}
		}
	}
	return
}

// Money 用户余额调整(根据用户ID)
func (webUserService) Money(userId int, money string, mod int) (err error) {
	l := "+"
	if mod == 0 {
		l = "-"
	}
	err = util.Database.Db.Model(&model.WebUser{Id: userId}).Updates(map[string]interface{}{
		"money":       gorm.Expr("`money`" + l + money),
		"update_time": str.DateTime(time.Now()),
	}).Error
	return
}

// SetRealAuthStatu 设置实名状态(根据用户ID)
func (w webUserService) SetRealAuthStatu(userId int, status bool, cardNo string, cardName string) (err error) {
	if userId <= 0 {
		return errors.New("用户不存在")
	}
	if cardNo != "" {
		GetWebUserConfSrvIns().SetConfig(userId, web_user_config.PersonCardNo, cardNo)
	}
	if cardName != "" {
		GetWebUserConfSrvIns().SetConfig(userId, web_user_config.PersonCardName, cardName)
	}
	s := "0"
	if status {
		s = "1"
	}
	err = w.UpdateById(userId, map[string]interface{}{
		"real_auth": s,
	})
	return
}

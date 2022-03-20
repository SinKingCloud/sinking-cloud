package service

import (
	"encoding/json"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// webUserService 单例对象
type webUserService struct {
}

// userObj 单例对象
var webUserObj *webUserService

// GetWebUserServiceIns 获取单例
func GetWebUserServiceIns() *webUserService {
	if webUserObj == nil {
		webUserObj = &webUserService{}
	}
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
			_ = json.Unmarshal([]byte(data), &user)
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
func (webUserService) Money(id int, money string, mod int) (err error) {
	l := "+"
	if mod == 0 {
		l = "-"
	}
	err = util.Database.Db.Model(&model.WebUser{Id: id}).Updates(map[string]interface{}{
		"money":       gorm.Expr("`money`" + l + money),
		"update_time": str.DateTime(time.Now()),
	}).Error
	return
}

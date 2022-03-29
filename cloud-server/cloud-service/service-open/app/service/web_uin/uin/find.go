package uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
)

// FindById 通过ID查询
func (s *service) FindById(id int) (data *model.AppUin, err error) {
	err = util.Database.Db.Where("`id` = ? ", id).First(&data).Error
	return
}

// FindByUserIdWithUin 通过用户id uin查询
func (s *service) FindByUserIdWithUin(userId int, uin string) (data *model.AppUin, err error) {
	err = util.Database.Db.Where("`user_id` = ? and `uin` = ?", userId, uin).First(&data).Error
	return
}

// FindByAppIdWithUin 通过应用id uin查询
func (s *service) FindByAppIdWithUin(appId int, uin string) (data *model.AppUin, err error) {
	err = util.Database.Db.Where("`app_id` = ? and `uin` = ?", appId, uin).First(&data).Error
	return
}

package uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/model"
)

// FindById 通过ID查询
func (s *Service) FindById(id int) (data *AppUinJob, err error) {
	err = util.Database.Db.Model(&model.AppUin{}).Where("`id` = ? ", id).First(&data).Error
	return
}

// FindByIdAll 通过ID查询
func (s *Service) FindByIdAll(id int) (data *model.AppUin, err error) {
	err = util.Database.Db.Model(&model.AppUin{}).Where("`id` = ? ", id).First(&data).Error
	return
}

// FindByUin 通过uin查询
func (s *Service) FindByUin(uin string) (data *AppUinJob, err error) {
	err = util.Database.Db.Model(&model.AppUin{}).Where("`uin` = ? ", uin).First(&data).Error
	return
}

// FindByUinAll 通过uin查询
func (s *Service) FindByUinAll(uin string) (data *model.AppUin, err error) {
	err = util.Database.Db.Model(&model.AppUin{}).Where("`uin` = ? ", uin).First(&data).Error
	return
}

// FindByUserIdWithUinAll 通过用户id uin查询
func (s *Service) FindByUserIdWithUinAll(userId int, uin string) (data *model.AppUin, err error) {
	err = util.Database.Db.Model(&model.AppUin{}).Where("`user_id` = ? and `uin` = ?", userId, uin).First(&data).Error
	return
}

// FindByUserIdWithUin 通过用户id uin查询
func (s *Service) FindByUserIdWithUin(userId int, uin string) (data *AppUinJob, err error) {
	err = util.Database.Db.Model(&model.AppUin{}).Where("`user_id` = ? and `uin` = ?", userId, uin).First(&data).Error
	return
}

// FindByAppIdWithUin 通过应用id uin查询
func (s *Service) FindByAppIdWithUin(appId int, uin string) (data *model.AppUin, err error) {
	err = util.Database.Db.Model(&model.AppUin{}).Where("`app_id` = ? and `uin` = ?", appId, uin).First(&data).Error
	return
}

// FindByIdWithAppUin 通过ID查询
func (s *Service) FindByIdWithAppUin(id int) (data *AppUin, err error) {
	err = util.Database.Db.Model(&model.AppUin{}).Where("`id` = ? ", id).First(&data).Error
	return
}

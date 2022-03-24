package app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
)

// FindByUserIdWithId 通过用户ID及ID查询
func (s *service) FindByUserIdWithId(userId int, id int) (data *model.WebUserApp, err error) {
	err = util.Database.Db.Where("`id` = ? and `user_id` = ?", id, userId).First(&data).Error
	return
}

// FindById 通过ID查询
func (s *service) FindById(id int) (data *model.WebUserApp, err error) {
	err = util.Database.Db.Where("`id` = ? ", id).First(&data).Error
	return
}

package oauth_user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// FindById 通过ID查询
func (s *Service) FindById(id int) (data *model.AppOauthUser, err error) {
	err = util.Database.Db.Where("`id` = ? ", id).First(&data).Error
	return
}

// FindByOauthIdAndUserId 通过oauthID和userId查询
func (s *Service) FindByOauthIdAndUserId(oauthId int, userId int) (data *model.AppOauthUser, err error) {
	err = util.Database.Db.Where("`oauth_id` = ? and `user_id` = ?", oauthId, userId).First(&data).Error
	return
}

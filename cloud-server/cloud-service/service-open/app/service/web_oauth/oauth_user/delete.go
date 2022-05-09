package oauth_user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// DeleteById 通过ID删除
func (s *service) DeleteById(id int) (err error) {
	err = util.Database.Db.Where("`id` = ? ", id).Delete(&model.AppOauthUser{}).Error
	return
}

// DeleteByOauthId 通过OauthID删除
func (s *service) DeleteByOauthId(id int) (err error) {
	err = util.Database.Db.Where("`oauth_id` = ? ", id).Delete(&model.AppOauthUser{}).Error
	return
}

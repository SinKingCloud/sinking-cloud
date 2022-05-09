package oauth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_oauth/oauth_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_oauth/oauth_user"
)

// DeleteById 通过ID删除
func (s *service) DeleteById(id int) (err error) {
	go func(id int) {
		_ = oauth_log.GetIns().DeleteByOauthId(id)
		_ = oauth_user.GetIns().DeleteByOauthId(id)
	}(id)
	err = util.Database.Db.Where("`id` = ? ", id).Delete(&model.AppOauth{}).Error
	if err == nil {
		s.Clear(id)
	}
	return
}

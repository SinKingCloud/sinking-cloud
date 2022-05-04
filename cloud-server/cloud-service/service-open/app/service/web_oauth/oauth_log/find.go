package oauth_log

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
)

// FindById 通过ID查询
func (s *service) FindById(id int) (data *model.AppOauthLog, err error) {
	err = util.Database.Db.Where("`id` = ? ", id).First(&data).Error
	return
}

// FindByTicket 通过Ticket查询
func (s *service) FindByTicket(ticket string) (data *model.AppOauthLog, err error) {
	err = util.Database.Db.Where("`ticket` = ? ", ticket).First(&data).Error
	return
}

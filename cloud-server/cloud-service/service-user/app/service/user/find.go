package user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// FindByPhone 查询信息
func (*Service) FindByPhone(phone string) (user *model.WebUser, err error) {
	err = util.Database.Db.Where("`phone` = ?", phone).First(&user).Error
	return
}

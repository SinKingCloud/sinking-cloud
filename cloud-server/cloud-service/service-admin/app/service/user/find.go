package user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// FindById 查询用户信息
func (*service) FindById(id int) (user *model.SystemUser, err error) {
	err = util.Database.Db.Where("`id` = ?", id).First(&user).Error
	return
}

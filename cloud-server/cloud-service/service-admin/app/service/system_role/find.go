package system_role

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// FindById 查询角色信息(根据角色ID)
func (*Service) FindById(id int) (role *model.SystemRole, err error) {
	err = util.Database.Db.Where("`id` = ?", id).First(&role).Error
	return
}

package notice

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// FindById 查询信息
func (*service) FindById(id int) (user *model.WebNotice, err error) {
	err = util.Database.Db.Where("`id` = ?", id).First(&user).Error
	return
}

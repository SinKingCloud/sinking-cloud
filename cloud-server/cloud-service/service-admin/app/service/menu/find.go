package menu

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// FindByKey 通过key查找页面数据
func (s *service) FindByKey(key string) (menu *Menu, err error) {
	err = util.Database.Db.Model(&model.SystemMenu{}).Where("`key` = ? and `status`=0", key).First(&menu).Error
	return
}

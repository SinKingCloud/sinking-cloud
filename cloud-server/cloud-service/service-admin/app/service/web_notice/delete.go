package web_notice

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// DeleteByIds 通过ID列表删除
func (s *service) DeleteByIds(ids []int) (err error) {
	err = util.Database.Db.Where("`id` in ? ", ids).Delete(&model.WebNotice{}).Error
	return
}

// DeleteById 通过ID删除
func (s *service) DeleteById(id []int) (err error) {
	err = util.Database.Db.Where("`id` = ? ", id).Delete(&model.WebNotice{}).Error
	return
}

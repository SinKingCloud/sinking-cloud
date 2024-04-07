package notice

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// UpdateById 通过ID更新
func (s *Service) UpdateById(id int, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.WebNotice{}).Where("`id` = ? ", id).Updates(data).Error
	return
}

// UpdateByLookNumInc 查看次数自增
func (s *Service) UpdateByLookNumInc(id int, num int) (err error) {
	err = util.Database.Db.Model(&model.WebNotice{}).Where("`id` = ? ", id).Updates(map[string]interface{}{
		"update_time": str.DateTime(time.Now()),
		"look_num":    gorm.Expr("`look_num`+" + strconv.Itoa(num)),
	}).Error
	return
}

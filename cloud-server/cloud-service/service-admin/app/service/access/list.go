package access

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

func (s *service) GetAccessList(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (err error, total int64, list []*model.SystemAccess) {
	offset := pageSize * (page - 1)
	err = util.Database.Db.Model(&model.SystemAccess{}).
		Where(where).
		Count(&total).
		Limit(pageSize).
		Offset(offset).
		Order(orderByField + " " + orderByType).
		Find(&list).Error
	return
}

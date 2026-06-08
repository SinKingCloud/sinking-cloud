package session

import (
	"service-admin/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// Select 查询会话列表
func (r *Repository) Select(where *SelectSession, orderByField string, orderByType string, page int, pageSize int) (list []*model.Session, total int64, err error) {
	offset := pageSize * (page - 1)
	query := r.Database.Db.Model(&model.Session{})
	if where != nil {
		if where.UserId > 0 {
			query = query.Where("user_id = ?", where.UserId)
		}
		if where.Device != "" {
			query = query.Where("device LIKE ?", "%"+where.Device+"%")
		}
		if where.Ip != "" {
			query = query.Where("ip LIKE ?", "%"+where.Ip+"%")
		}
		if where.Location != "" {
			query = query.Where("location LIKE ?", "%"+where.Location+"%")
		}
		if where.CreateTimeStart != "" {
			query = query.Where("create_time >= ?", where.CreateTimeStart)
		}
		if where.CreateTimeEnd != "" {
			query = query.Where("create_time <= ?", where.CreateTimeEnd)
		}
		if where.UpdateTimeStart != "" {
			query = query.Where("update_time >= ?", where.UpdateTimeStart)
		}
		if where.UpdateTimeEnd != "" {
			query = query.Where("update_time <= ?", where.UpdateTimeEnd)
		}
		if where.ExpireTimeStart != "" {
			query = query.Where("expire_time >= ?", where.ExpireTimeStart)
		}
		if where.ExpireTimeEnd != "" {
			query = query.Where("expire_time <= ?", where.ExpireTimeEnd)
		}
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(r.GetOrderBy(orderByField, orderByType)).Find(&list).Error
	return
}

// SelectExpireSessions 查询过期会话
func (r *Repository) SelectExpireSessions(tx ...*gorm.DB) (map[int64]string, error) {
	type expireSession struct {
		Id    int64  `gorm:"column:id"`
		Token string `gorm:"column:token"`
	}

	list := make([]*expireSession, 0)
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	err := db.
		Table("cloud_sessions").
		Select("id, token").
		Where("expire_time < ?", str.DateTime(time.Now())).
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	result := make(map[int64]string, len(list))
	for _, item := range list {
		if item == nil || item.Id <= 0 {
			continue
		}
		result[item.Id] = item.Token
	}
	return result, nil
}

package session

import (
	"service-user/app/model"
	"time"
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
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}

// SelectUnExpireTokensByUserIds 根据userid列表查询所有未过期的tokens
func (r *Repository) SelectUnExpireTokensByUserIds(userIds []int64) (tokens []string, err error) {
	if len(userIds) == 0 {
		return nil, nil
	}
	tokenMap := make(map[string]struct{}, len(userIds))
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	err = r.Database.BatchExecute(userIds, 1000, func(batch interface{}) error {
		var batchTokens []string
		err2 := r.Database.Db.
			Model(&model.Session{}).
			Distinct("token").
			Where("user_id IN ?", batch.([]int64)).
			Where("expire_time > ?", nowStr).
			Pluck("token", &batchTokens).Error
		if err2 != nil {
			return err2
		}
		for _, token := range batchTokens {
			if _, ok := tokenMap[token]; ok {
				continue
			}
			tokenMap[token] = struct{}{}
			tokens = append(tokens, token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

package session

import (
	"service-user/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// Select 查询会话列表
func (r *Repository) Select(where *SelectSession, queryPage *page.Query) (*page.Result[*model.Session], error) {
	query := r.Database.Db.Model(&model.Session{})
	if where != nil {
		if where.UserId != nil {
			query = query.Where("user_id = ?", *where.UserId)
		}
		if where.Keyword != nil {
			keyword := "%" + *where.Keyword + "%"
			query = query.Where("(device LIKE ? OR ip LIKE ? OR location LIKE ?)", keyword, keyword, keyword)
		}
		if where.Device != nil {
			query = query.Where("device LIKE ?", "%"+*where.Device+"%")
		}
		if where.Ip != nil {
			query = query.Where("ip LIKE ?", "%"+*where.Ip+"%")
		}
		if where.Location != nil {
			query = query.Where("location LIKE ?", "%"+*where.Location+"%")
		}
		if where.CreateTimeStart != nil {
			query = query.Where("create_time >= ?", *where.CreateTimeStart)
		}
		if where.CreateTimeEnd != nil {
			query = query.Where("create_time <= ?", *where.CreateTimeEnd)
		}
		if where.UpdateTimeStart != nil {
			query = query.Where("update_time >= ?", *where.UpdateTimeStart)
		}
		if where.UpdateTimeEnd != nil {
			query = query.Where("update_time <= ?", *where.UpdateTimeEnd)
		}
		if where.ExpireTimeStart != nil {
			query = query.Where("expire_time >= ?", *where.ExpireTimeStart)
		}
		if where.ExpireTimeEnd != nil {
			query = query.Where("expire_time <= ?", *where.ExpireTimeEnd)
		}
	}
	return r.Repository.SelectPage(query, queryPage, "token")
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

package scene

import (
	"service-driver/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// SelectIdNameMap 查询id name 映射
func (r *Repository) SelectIdNameMap() (map[int64]string, error) {
	var scene []struct {
		Id   int64
		Name string
	}
	err := r.Database.Db.
		Model(&model.Scene{}).
		Select("id", "name").
		Order("id ASC").
		Find(&scene).Error
	if err != nil {
		return nil, err
	}
	result := make(map[int64]string, len(scene))
	for _, c := range scene {
		result[c.Id] = c.Name
	}
	return result, nil
}

// Select 查询数据
func (r *Repository) Select(where *SelectScene, queryPage *page.Query) (*page.Result[*model.Scene], error) {
	query := r.Database.Db.Model(&model.Scene{})
	if where != nil {
		if where.ChannelId != nil {
			query = query.Where("channel_id = ?", *where.ChannelId)
		}
		if where.Status != nil {
			query = query.Where("status = ?", *where.Status)
		}
		if where.Login != nil {
			query = query.Where("login = ?", *where.Login)
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
		if where.Name != nil {
			query = query.Where("name LIKE ?", "%"+*where.Name+"%")
		}
		if where.Remark != nil {
			query = query.Where("remark LIKE ?", "%"+*where.Remark+"%")
		}
	}
	return r.Repository.SelectPage(query, queryPage)
}

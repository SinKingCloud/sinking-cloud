package channel

import (
	"service-pay/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// SelectIdNameMap 查询id name 映射
func (r *Repository) SelectIdNameMap() (map[int64]string, error) {
	var data []struct {
		Id   int64  `gorm:"column:id" json:"id"`
		Name string `gorm:"column:name" json:"name"`
	}
	err := r.Database.Db.
		Model(&model.Channel{}).
		Select("id", "name").
		Order("id ASC").
		Find(&data).Error
	if err != nil {
		return nil, err
	}
	result := make(map[int64]string, len(data))
	for _, c := range data {
		result[c.Id] = c.Name
	}
	return result, nil
}

// Select 查询数据
func (r *Repository) Select(where *SelectChannel, queryPage *page.Query) (*page.Result[*Channel], error) {
	query := r.Database.Db.Model(&model.Channel{})
	if where != nil {
		if where.Type != "" {
			query = query.Where("type = ?", where.Type)
		}
		if where.Plugin != "" {
			query = query.Where("plugin = ?", where.Plugin)
		}
		if where.Status != "" {
			query = query.Where("status = ?", where.Status)
		}
		if where.Name != "" {
			query = query.Where("name LIKE ?", "%"+where.Name+"%")
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
	}
	return r.Repository.SelectPage(query, queryPage)
}

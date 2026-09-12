package role

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// SelectByIds 根据ID集合查询角色
func (r *Repository) SelectByIds(ids []int64) ([]*model.Role, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var roleList []*model.Role
	err := r.Database.BatchExecute(ids, 1000, func(batch interface{}) error {
		batchIds := batch.([]int64)
		var roles []*model.Role
		err := r.Database.Db.
			Model(&model.Role{}).
			Where("id IN ?", batchIds).
			Find(&roles).Error
		if err != nil {
			return err
		}
		roleList = append(roleList, roles...)
		return nil
	})
	return roleList, err
}

// SelectIdNameMap 查询id name 映射
func (r *Repository) SelectIdNameMap() (map[int64]string, error) {
	var data []struct {
		Id   int64  `gorm:"column:id" json:"id"`
		Name string `gorm:"column:name" json:"name"`
	}
	err := r.Database.Db.
		Model(&model.Role{}).
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
func (r *Repository) Select(where *SelectRole, queryPage *page.Query) (*page.Result[*model.Role], error) {
	query := r.Database.Db.Model(&model.Role{})
	if where != nil {
		if where.Name != nil {
			query = query.Where("name LIKE ?", "%"+*where.Name+"%")
		}
		if where.Status != nil {
			query = query.Where("status = ?", *where.Status)
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
	}
	return r.Repository.SelectPage(query, queryPage)
}

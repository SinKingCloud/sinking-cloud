package file

import (
	"service-driver/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// SelectChannelIdSceneIdHashPathByIds 根据 ids 查询清理所需的通道、场景、哈希和路径。
func (r *Repository) SelectChannelIdSceneIdHashPathByIds(ids []int64) ([]*model.File, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	list := make([]*model.File, 0, len(ids))
	err := r.Database.BatchExecute(ids, 1000, func(batch interface{}) error {
		var rows []*model.File
		err := r.Database.Db.
			Model(&model.File{}).
			Select("id", "channel_id", "scene_id", "hash", "path").
			Where("id IN ?", batch.([]int64)).
			Find(&rows).Error
		if err != nil {
			return err
		}
		list = append(list, rows...)
		return nil
	})
	return list, err
}

// Select 查询数据
func (r *Repository) Select(where *SelectFile, queryPage *page.Query) (*page.Result[*model.File], error) {
	query := r.Database.Db.Model(&model.File{})
	if where != nil {
		if where.ChannelId != nil {
			query = query.Where("channel_id = ?", *where.ChannelId)
		}
		if where.SceneId != nil {
			query = query.Where("scene_id = ?", *where.SceneId)
		}
		if where.Ext != nil {
			query = query.Where("ext = ?", *where.Ext)
		}
		if where.Hash != nil {
			query = query.Where("hash = ?", *where.Hash)
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
		if where.Name != nil {
			query = query.Where("name LIKE ?", "%"+*where.Name+"%")
		}
		if where.Path != nil {
			query = query.Where("path LIKE ?", *where.Path+"%")
		}
	}
	return r.Repository.SelectPage(query, queryPage)
}

package mind_share

import (
	"service-space/app/enum/mind_share_status"
	"service-space/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// SelectByDataId 根据导图数据ID查询分享列表
func (r *Repository) SelectByDataId(dataId int64, queryPage *page.Query) (list []*model.MindShare, total int64, err error) {
	if dataId <= 0 || queryPage == nil {
		return []*model.MindShare{}, 0, nil
	}
	offset := queryPage.PageSize * (queryPage.Page - 1)
	query := r.Database.Db.Model(&model.MindShare{}).Where("data_id = ?", dataId)
	err = query.Count(&total).Limit(queryPage.PageSize).Offset(offset).Order(r.GetOrderBy(queryPage.OrderByField, queryPage.OrderByType)).Find(&list).Error
	return
}

// SelectValidDataIds 查询存在有效分享的数据ID集合
func (r *Repository) SelectValidDataIds(dataIds []int64) ([]int64, error) {
	if len(dataIds) == 0 {
		return []int64{}, nil
	}
	validDataIds := make([]int64, 0, len(dataIds))
	dataIdMap := make(map[int64]struct{}, len(dataIds))
	now := str.DateTime(time.Now())
	err := r.Database.BatchExecute(dataIds, 1000, func(batch interface{}) error {
		var batchDataIds []int64
		if err := r.Database.Db.
			Model(&model.MindShare{}).
			Distinct("data_id").
			Where("data_id IN ?", batch.([]int64)).
			Where("status = ?", mind_share_status.Enable).
			Where("expire_time IS NULL OR expire_time > ?", now).
			Pluck("data_id", &batchDataIds).Error; err != nil {
			return err
		}
		for _, dataId := range batchDataIds {
			if _, ok := dataIdMap[dataId]; ok {
				continue
			}
			dataIdMap[dataId] = struct{}{}
			validDataIds = append(validDataIds, dataId)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return validDataIds, nil
}

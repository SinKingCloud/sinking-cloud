package uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// Select 查询
func (s *service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*AppUin, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.AppUin{})
	if where["user_id"] != "" {
		query.Where("`user_id` = ?", where["user_id"])
	}
	if where["app_id"] != "" {
		query.Where("`app_id` = ?", where["app_id"])
	}
	if where["server_id"] != "" {
		query.Where("`server_id` = ?", where["server_id"])
	}
	if where["server_key"] != "" {
		query.Where("`server_key` = ?", where["server_key"])
	}
	if where["uin_status"] != "" {
		query.Where("`uin_status` = ?", where["uin_status"])
	}
	if where["uin"] != "" {
		query.Where("`uin` = ?", where["uin"])
	}
	if where["status"] != "" {
		query.Where("`status` = ?", where["status"])
	}
	if where["run_time_start"] != "" {
		query.Where("`run_time` >= ?", where["run_time_start"])
	}
	if where["run_time_end"] != "" {
		query.Where("`run_time` <= ?", where["run_time_end"])
	}
	if where["create_time_start"] != "" {
		query.Where("`create_time` >= ?", where["create_time_start"])
	}
	if where["create_time_end"] != "" {
		query.Where("`create_time` <= ?", where["create_time_end"])
	}
	if where["update_time_start"] != "" {
		query.Where("`update_time` >= ?", where["update_time_start"])
	}
	if where["update_time_end"] != "" {
		query.Where("`update_time` <= ?", where["update_time_end"])
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}

// SelectUinJob 查询全部代挂任务
func (s *service) SelectUinJob(limit int) (list []*model.AppUin, err error) {
	query := util.Database.Db.Model(&model.AppUin{})
	query.Where("`status`=0")
	err = query.Limit(limit).Find(&list).Error
	return
}

// SelectAllUinJob 查询全部代挂任务
func (s *service) SelectAllUinJob(limit int, minute int64, runTime string) (list []*model.AppUin, err error) {
	query := util.Database.Db.Model(&model.AppUin{})
	query.Where("`timing`=? and `run_time`<=? and (`player_day`>0 or `sign_day`>0)", minute, runTime)
	err = query.Limit(limit).Find(&list).Error
	return
}

// SelectUinCheckJob 查询挂机检测任务
func (s *service) SelectUinCheckJob(limit int, startTime string, endTime string) (list []*model.AppUin, err error) {
	query := util.Database.Db.Model(&model.AppUin{})
	query.Where("`uin_status`=0 and `status` =0 and `check_time` between ? and ?", startTime, endTime)
	err = query.Limit(limit).Find(&list).Error
	return
}

// SelectByUin 查询QQ列表
func (s *service) SelectByUin(uin string) (list []*model.AppUin, err error) {
	err = util.Database.Db.Model(&model.AppUin{}).
		Where("`uin`=?", uin).
		Find(&list).Error
	return
}

// SelectByServerId 查询QQ列表
func (s *service) SelectByServerId(limit int, serverId int) (list []*model.AppUin, err error) {
	err = util.Database.Db.Model(&model.AppUin{}).
		Where("`server_id`=?", serverId).
		Limit(limit).
		Find(&list).Error
	return
}

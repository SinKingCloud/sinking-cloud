package uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
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

// AppUinJob QQ记录
type AppUinJob struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId     int          `gorm:"column:user_id" json:"user_id"`
	AppId      int          `gorm:"column:app_id" json:"app_id"`
	ServerId   int          `gorm:"column:server_id" json:"server_id"`
	ServerKey  string       `gorm:"column:server_key" json:"server_key"`
	Uin        string       `gorm:"column:uin" json:"uin"`
	UinStatus  int          `gorm:"column:uin_status" json:"uin_status"`
	Device     int          `gorm:"column:device" json:"device"`
	SignDay    int          `gorm:"column:sign_day" json:"sign_day"`
	PlayerDay  int          `gorm:"column:player_day" json:"player_day"`
	RunDay     int          `gorm:"column:run_day" json:"run_day"`
	Setting    string       `gorm:"column:setting" json:"setting"`
	Timing     int          `gorm:"column:timing" json:"timing"`
	Status     int          `gorm:"column:status" json:"status"`
	RunTime    str.DateTime `gorm:"column:run_time" json:"run_time"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// SelectAllUinJob 查询全部代挂任务
func (s *service) SelectAllUinJob(limit int, minute int64, runTime string) (list []*AppUinJob, err error) {
	query := util.Database.Db.Model(&model.AppUin{})
	query.Where("`timing`=? and `run_time`<=? and (`player_day`>0 or `sign_day`>0)", minute, runTime)
	err = query.Limit(limit).Find(&list).Error
	return
}

// SelectUinCheckJob 查询挂机检测任务
func (s *service) SelectUinCheckJob(limit int, startTime string, endTime string) (list []*AppUinJob, err error) {
	query := util.Database.Db.Model(&model.AppUin{})
	query.Where("`uin_status`=0 and `status` =0 and `check_time` between ? and ?", startTime, endTime)
	err = query.Limit(limit).Find(&list).Error
	return
}

// SelectByUin 查询QQ列表
func (s *service) SelectByUin(uin string) (list []*AppUinJob, err error) {
	err = util.Database.Db.Model(&model.AppUin{}).
		Where("`uin`=?", uin).
		Find(&list).Error
	return
}

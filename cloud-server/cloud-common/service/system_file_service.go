package service

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"time"
)

// systemFileService config_service
type systemFileService struct {
}

// systemFileObj systemConfigObj对象
var systemFileObj *systemFileService

// GetSystemFileServiceIns 获取单例
func GetSystemFileServiceIns() *systemFileService {
	if systemFileObj == nil {
		systemFileObj = &systemFileService{}
	}
	return systemFileObj
}

func (*systemFileService) Create(file *model.SystemFile) (err error) {
	var count int64
	err = util.Database.Db.Model(&model.SystemFile{}).Where("`hash`=?", file.Hash).Count(&count).Error
	if err != nil {
		return
	}
	if count <= 0 {
		//新增
		err = util.Database.Db.Create(&file).Error
	} else {
		//修改
		file.UpdateTime = str.DateTime(time.Now())
		err = util.Database.Db.Where("`hash`=?", file.Hash).Updates(&file).Error
	}
	return
}

package system_file

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/aliyun"
)

// DeleteByIds 通过ID列表删除
func (s *Service) DeleteByIds(ids []int) (err error) {
	err = util.Database.Db.Where("`id` in ? ", ids).Delete(&model.SystemFile{}).Error
	return
}

// DeleteById 通过ID删除
func (s *Service) DeleteById(id []int) (err error) {
	err = util.Database.Db.Where("`id` = ? ", id).Delete(&model.SystemFile{}).Error
	return
}

// DeleteFiles 删除文件
func (s *Service) DeleteFiles(filePaths []string) {
	configs := service2.GetSystemConfigServiceIns().GetGroupByCache(config.AliyunGroup, false)
	for _, v := range filePaths {
		aliyun.Oss.Delete(
			configs[config.AliyunAccessKey],
			configs[config.AliyunAccessSecret],
			configs[config.AliyunOssEndPoint],
			configs[config.AliyunOssBuckName],
			v,
		)
	}
}

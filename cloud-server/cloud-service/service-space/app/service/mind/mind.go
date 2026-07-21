package mind

import (
	"service-space/app/model"
	mindDataRepository "service-space/app/repository/mind_data"
	mindShareRepository "service-space/app/repository/mind_share"
	mindTypeRepository "service-space/app/repository/mind_type"
	mindVersionRepository "service-space/app/repository/mind_version"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Service 导图服务接口
type Service interface {
	CreateData(userId int64, typeId int64, title string, describe string, value string) (*Data, error)                   //创建导图数据
	CopyData(userId int64, id int64) (*Data, error)                                                                      //复制导图数据
	DeleteData(userId int64, id int64) error                                                                             //删除导图数据
	GetDataInfoByUserId(userId int64, id int64) (*Data, error)                                                           //查询用户导图数据详情
	SaveData(userId int64, id int64, versionType int, remark string, value string) (*Version, error)                     //保存导图数据
	SelectDataByUserId(userId int64, where *mindDataRepository.SelectMindData, page *page.Query) ([]*Data, int64, error) //查询用户导图数据列表
	UpdateData(userId int64, id int64, data *mindDataRepository.UpdateMindData) error                                    //更新导图数据信息
	CreateType(userId int64, name string) (*Type, error)                                                                 //创建自定义分类
	DeleteType(userId int64, id int64) error                                                                             //删除自定义分类
	SelectTypes(userId int64, page *page.Query) ([]*Type, int64, error)                                                  //自定义分类列表
	UpdateType(userId int64, id int64, name string) error                                                                //更新自定义分类
	GetVersion(userId int64, dataId int64, id int64) (*Version, error)                                                   //版本详情
	RestoreVersion(userId int64, dataId int64, id int64) (*Data, error)                                                  //恢复版本
	SelectVersions(userId int64, dataId int64, versionType int, page *page.Query) ([]*Version, int64, error)             //版本列表
	CreateShare(userId int64, dataId int64, status int, expireTime string, remark string) (*model.MindShare, error)      //创建分享
	DeleteShare(userId int64, dataId int64, id int64) error                                                              //删除分享
	SelectSharesByUserId(userId int64, dataId int64, page *page.Query) ([]*model.MindShare, int64, error)                //分享列表
	UpdateShare(userId int64, dataId int64, id int64, data *mindShareRepository.UpdateMindShare) error                   //更新分享
	GetDataByShareCode(code string) (*SharedData, error)                                                                 //通过分享码获取导图
}

// 导图服务
type service struct {
	mindDataRepository    mindDataRepository.Interface
	mindTypeRepository    mindTypeRepository.Interface
	mindVersionRepository mindVersionRepository.Interface
	mindShareRepository   mindShareRepository.Interface
	database              *database.Database
	cache                 cache.Interface
	configService         *config.Service
}

// NewService 实例化service
func NewService(mindDataRepository mindDataRepository.Interface, mindTypeRepository mindTypeRepository.Interface, mindVersionRepository mindVersionRepository.Interface, mindShareRepository mindShareRepository.Interface, database *database.Database, cache cache.Interface, configService *config.Service) *service {
	return &service{
		mindDataRepository:    mindDataRepository,
		mindTypeRepository:    mindTypeRepository,
		mindVersionRepository: mindVersionRepository,
		mindShareRepository:   mindShareRepository,
		database:              database,
		cache:                 cache,
		configService:         configService,
	}
}

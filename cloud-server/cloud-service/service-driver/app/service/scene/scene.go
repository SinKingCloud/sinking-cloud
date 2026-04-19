package scene

import (
	"service-driver/app/repository/file"
	"service-driver/app/repository/scene"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Service service接口
type Service interface {
}

// Session服务
type service struct {
	repositoryScene scene.Interface
	repositoryFile  file.Interface
	database        *database.Database
	cache           cache.Interface
}

// NewService 实例化service
func NewService(repositoryScene scene.Interface, repositoryFile file.Interface, database *database.Database, cache cache.Interface) *service {
	return &service{
		repositoryScene: repositoryScene,
		repositoryFile:  repositoryFile,
		database:        database,
		cache:           cache,
	}
}

package menu

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"

// Interface 菜单仓储接口
type Interface interface {
}

// Repository 用户仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建User仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}

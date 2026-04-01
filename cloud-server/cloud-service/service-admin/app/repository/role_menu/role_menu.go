package role_menu

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"

// Interface 角色菜单仓储接口
type Interface interface {
	SelectMenuIdsByRoleId(roleId int64) ([]int64, error)
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}

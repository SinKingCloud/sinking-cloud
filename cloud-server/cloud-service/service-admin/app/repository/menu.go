package repository

import (
	"service-admin/app/repository/menu"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Menu 菜单仓储接口
type Menu interface {
}

// NewMenu 创建Menu仓储
func NewMenu(db *database.Database) *menu.Repository {
	return &menu.Repository{Database: db}
}

package order_log

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Interface 仓储接口
type Interface interface {
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}

package config

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Repository 配置仓储
type Repository struct {
	Database *database.Database
}

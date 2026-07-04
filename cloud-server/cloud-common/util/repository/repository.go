package repository

import (
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository 基础仓储
type Repository struct {
	Database *database.Database
	Model    interface{}
}

// NewRepository 创建基础仓储
func NewRepository(db *database.Database, model interface{}) *Repository {
	return &Repository{Database: db, Model: model}
}

// GetOrderBy 获取排序条件
func (r *Repository) GetOrderBy(field string, orderType string) clause.OrderByColumn {
	if field == "" {
		return clause.OrderByColumn{}
	}
	orderType = strings.ToLower(orderType)
	if orderType != "asc" && orderType != "desc" {
		orderType = "desc"
	}
	stmt := &gorm.Statement{DB: r.Database.Db}
	if err := stmt.Parse(r.Model); err != nil {
		return clause.OrderByColumn{}
	}
	if modelField := stmt.Schema.LookUpField(field); modelField != nil {
		field = modelField.DBName
	} else {
		if stmt.Schema.PrioritizedPrimaryField != nil {
			field = stmt.Schema.PrioritizedPrimaryField.DBName
		} else {
			return clause.OrderByColumn{}
		}
	}
	return clause.OrderByColumn{
		Column: clause.Column{Name: field},
		Desc:   orderType == "desc",
	}
}

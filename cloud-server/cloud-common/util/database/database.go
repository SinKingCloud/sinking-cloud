package database

import (
	"gorm.io/gorm"
)

// Database 数据库连接
type Database struct {
	Db      *gorm.DB
	DbError error
}

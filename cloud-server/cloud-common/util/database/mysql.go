package database

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewMysql 实例化一个mysql连接
func NewMysql(host string, port string, user string, pwd string, database string) *Database {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pwd, host, port, database)
	Db, DbError := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger(),
	})
	sqlDB, err := Db.DB()
	if err != nil {
		return &Database{Db: Db, DbError: errors.New("set sql run config failed")}
	}
	sqlDB.SetMaxIdleConns(1000)
	sqlDB.SetMaxOpenConns(10000)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	return &Database{Db: Db, DbError: DbError}
}

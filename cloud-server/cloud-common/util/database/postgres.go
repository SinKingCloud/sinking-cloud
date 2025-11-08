package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

// NewPostgresSql 实例化一个postgres连接
func NewPostgresSql(host string, port string, user string, pwd string, database string) *Database {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		host, port, user, pwd, database)
	Db, DbError := gorm.Open(postgres.Open(dsn), &gorm.Config{
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

package database

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// NewPostgresSql 实例化一个postgres连接
func NewPostgresSql(master []string, slave []string, port string, user string, pwd string, database string) *Database {
	dsnTpl := "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai"
	var sources []gorm.Dialector
	for _, h := range master {
		h = strings.ReplaceAll(h, " ", "")
		if h != "" {
			sources = append(sources, postgres.Open(fmt.Sprintf(dsnTpl, h, port, user, pwd, database)))
		}
	}
	var replicas []gorm.Dialector
	for _, h := range slave {
		h = strings.ReplaceAll(h, " ", "")
		if h != "" {
			replicas = append(replicas, postgres.Open(fmt.Sprintf(dsnTpl, h, port, user, pwd, database)))
		}
	}
	if len(sources) == 0 {
		return &Database{DbError: errors.New("master is empty")}
	}
	Db, DbError := gorm.Open(sources[0], &gorm.Config{
		Logger: newLogger(),
	})
	if DbError == nil && (len(sources) > 1 || len(replicas) > 0) {
		DbError = Db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  sources,
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}))
	}
	sqlDB, err := Db.DB()
	if err != nil {
		return &Database{Db: Db, DbError: errors.New("set sql run config failed")}
	}
	sqlDB.SetMaxIdleConns(1000)
	sqlDB.SetMaxOpenConns(10000)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	return &Database{Db: Db, DbError: DbError}
}

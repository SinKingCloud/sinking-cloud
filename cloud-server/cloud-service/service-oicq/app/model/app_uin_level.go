package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// AppUinLevel QQ折扣等级
type AppUinLevel struct {
	Id                 int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Name               string       `gorm:"column:name" json:"name"`
	Describe           string       `gorm:"column:describe" json:"describe"`
	Price              float64      `gorm:"column:price" json:"price"`
	MonthSpend         float64      `gorm:"column:month_spend" json:"month_spend"`
	MonthPrice         float64      `gorm:"column:month_price" json:"month_price"`
	QuarterPrice       float64      `gorm:"column:quarter_price" json:"quarter_price"`
	HalfPrice          float64      `gorm:"column:half_price" json:"half_price"`
	YearPrice          float64      `gorm:"column:year_price" json:"year_price"`
	ForeverPrice       float64      `gorm:"column:forever_price" json:"forever_price"`
	PlayerMonthPrice   float64      `gorm:"column:player_month_price" json:"player_month_price"`
	PlayerQuarterPrice float64      `gorm:"column:player_quarter_price" json:"player_quarter_price"`
	PlayerHalfPrice    float64      `gorm:"column:player_half_price" json:"player_half_price"`
	PlayerYearPrice    float64      `gorm:"column:player_year_price" json:"player_year_price"`
	PlayerForeverPrice float64      `gorm:"column:player_forever_price" json:"player_forever_price"`
	SignMonthPrice     float64      `gorm:"column:sign_month_price" json:"sign_month_price"`
	SignQuarterPrice   float64      `gorm:"column:sign_quarter_price" json:"sign_quarter_price"`
	SignHalfPrice      float64      `gorm:"column:sign_half_price" json:"sign_half_price"`
	SignYearPrice      float64      `gorm:"column:sign_year_price" json:"sign_year_price"`
	SignForeverPrice   float64      `gorm:"column:sign_forever_price" json:"sign_forever_price"`
	Status             int          `gorm:"column:status" json:"status"`
	CreateTime         str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime         str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*AppUinLevel) TableName() string {
	return "app_uin_levels"
}

// BeforeCreate 创建前
func (t *AppUinLevel) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *AppUinLevel) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}

package user

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// SelectUser 用户查询
type SelectUser struct {
	Account         *string
	Phone           *string
	Email           *string
	NickName        *string
	RealAuthType    *int
	RealAuthStatus  *int
	Status          *int
	CreateTimeStart *string
	CreateTimeEnd   *string
	UpdateTimeStart *string
	UpdateTimeEnd   *string
}

// UpdateUser 用户更新
type UpdateUser struct {
	Account        *string
	Phone          *string
	Email          *string
	Password       *string
	Avatar         *string
	NickName       *string
	RealAuthType   *int
	RealAuthStatus *int
	Status         *int
}

// User 用户列表展示数据
type User struct {
	Id             int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Account        string       `gorm:"column:account" json:"account"`
	Phone          string       `gorm:"column:phone" json:"phone"`
	Email          string       `gorm:"column:email" json:"email"`
	Avatar         string       `gorm:"column:avatar" json:"avatar"`
	NickName       string       `gorm:"column:nick_name" json:"nick_name"`
	RealAuthType   int          `gorm:"column:real_auth_type" json:"real_auth_type"`
	RealAuthStatus int          `gorm:"column:real_auth_status" json:"real_auth_status"`
	Status         int          `gorm:"column:status" json:"status"`
	CreateTime     str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime     str.DateTime `gorm:"column:update_time" json:"update_time"`
}

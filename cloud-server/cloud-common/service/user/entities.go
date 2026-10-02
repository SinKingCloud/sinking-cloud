package user

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// WebUser 网站用户
type WebUser struct {
	Id             int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Account        string       `gorm:"column:account" json:"account"`
	Phone          string       `gorm:"column:phone" json:"phone"`
	Email          string       `gorm:"column:email" json:"email"`
	Password       string       `gorm:"column:password" json:"password"`
	Avatar         string       `gorm:"column:avatar" json:"avatar"`
	NickName       string       `gorm:"column:nick_name" json:"nick_name"`
	RealAuthStatus int          `gorm:"column:real_auth_status" json:"real_auth_status"`
	Status         int          `gorm:"column:status" json:"status"`
	CreateTime     str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime     str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// SystemUser 系统用户
type SystemUser struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Account    string       `gorm:"column:account" json:"account"`
	Phone      string       `gorm:"column:phone" json:"phone"`
	Email      string       `gorm:"column:email" json:"email"`
	Password   string       `gorm:"column:password" json:"password"`
	Avatar     string       `gorm:"column:avatar" json:"avatar"`
	NickName   string       `gorm:"column:nick_name" json:"nick_name"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// GetWebUserInfoParam 获取网站用户详情请求参数
type GetWebUserInfoParam struct {
	UserId  int64 // 用户ID
	Refresh bool  // 是否强制刷新缓存
}

// GetSystemUserInfoParam 获取系统用户详情请求参数
type GetSystemUserInfoParam struct {
	UserId  int64 // 用户ID
	Refresh bool  // 是否强制刷新缓存
}

// CreateSystemUserLogParam 创建系统用户操作日志请求参数
type CreateSystemUserLogParam struct {
	Id        int64
	UserId    int64
	RequestId string
	RequestIp string
	Type      int
	Title     string
	Content   string
}

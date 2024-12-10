package page

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
)

type ValidatePage struct {
	Page     int `json:"page" default:"1" validate:"omitempty,gte=1" label:"页码编号"`
	PageSize int `json:"page_size" default:"20" validate:"omitempty,default=10,lte=1000" label:"每页数量"`
}

// ValidatePageDefault 默认分页验证
func ValidatePageDefault(c *sinking_server.Context) *ValidatePage {
	pageInfo := &ValidatePage{}
	if ok, _ := c.ValidatorAll(pageInfo); !ok {
		pageInfo.Page = 1
		pageInfo.PageSize = 20
	}
	if pageInfo.Page <= 0 {
		pageInfo.Page = 1
	}
	if pageInfo.PageSize <= 0 || pageInfo.PageSize > 1000 {
		pageInfo.PageSize = 20
	}
	return pageInfo
}

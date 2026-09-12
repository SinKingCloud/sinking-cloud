package real_auth

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// UpdateRealAuth 实名记录更新数据
type UpdateRealAuth struct {
	Type        *int
	CompanyName *string
	CompanyNo   *string
	CardName    *string
	CardNo      *string
	Files       *string
	Status      *int
	AuthTime    *str.DateTime
}

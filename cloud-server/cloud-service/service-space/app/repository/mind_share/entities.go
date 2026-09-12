package mind_share

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// UpdateMindShare 导图分享更新
type UpdateMindShare struct {
	Status     *int
	ExpireTime *str.DateTime
	Remark     *string
}

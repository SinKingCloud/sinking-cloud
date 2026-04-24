package balance

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// Log 余额明细展示数据
type Log struct {
	Id           int64        `json:"id"`
	Type         int          `json:"type"`
	Money        float64      `json:"money"`
	MoneyBefore  float64      `json:"money_before"`
	MoneyAfter   float64      `json:"money_after"`
	FrozenBefore float64      `json:"frozen_before"`
	FrozenAfter  float64      `json:"frozen_after"`
	Describe     string       `json:"describe"`
	CreateTime   str.DateTime `json:"create_time"`
}

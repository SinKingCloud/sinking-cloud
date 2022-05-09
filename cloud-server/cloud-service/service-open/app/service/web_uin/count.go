package web_uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/order"
	uin2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/uin"
	"strconv"
	"time"
)

// GetAllNum 获取QQ个数
func GetAllNum(userId int, status int) int64 {
	return uin2.GetIns().Count(map[string]string{
		"user_id":    strconv.Itoa(userId),
		"uin_status": strconv.Itoa(status),
	})
}

// GetTodayAddNum 获取今日新增个数
func GetTodayAddNum(userId int, status int) int64 {
	return uin2.GetIns().Count(map[string]string{
		"user_id":           strconv.Itoa(userId),
		"uin_status":        strconv.Itoa(status),
		"create_time_start": str.NewTimeTool().GetZeroTime(time.Now()).Format("2006-01-02 15:04:05"),
	})
}

// GetAllCost 获取总消费
func GetAllCost(userId int, status int) int64 {
	return order.GetIns().Count(map[string]string{
		"user_id": strconv.Itoa(userId),
		"status":  strconv.Itoa(status),
	})
}

// GetTodayCost 获取今日消费
func GetTodayCost(userId int, status int) int64 {
	return order.GetIns().Count(map[string]string{
		"user_id":           strconv.Itoa(userId),
		"status":            strconv.Itoa(status),
		"create_time_start": str.NewTimeTool().GetZeroTime(time.Now()).Format("2006-01-02 15:04:05"),
	})
}

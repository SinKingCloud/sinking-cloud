package cron

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
)

func Init() {
	//注册清理任务
	_ = service.Cron.RegisterRpcTask(1, "清理系统过期订单", "0 0 * * * *", constant.ServicePay, "CleanExpiredUnPayOrder", false, true)
}

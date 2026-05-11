package cron

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
)

func Init() {
	// 注册清理任务
	_ = service.Cron.RegisterRpcTask(3, "清理后台过期会话", "0 * * * * *", constant.ServiceAdmin, "ClearExpireSessions", false, true)
}

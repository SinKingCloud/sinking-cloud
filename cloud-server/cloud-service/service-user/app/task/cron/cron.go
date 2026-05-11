package cron

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
)

func Init() {
	// 注册清理任务
	_ = service.Cron.RegisterRpcTask(2, "清理用户过期会话", "0 * * * * *", constant.ServiceUser, "ClearExpireSessions", false, true)
}

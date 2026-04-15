package route

import (
	"service-cron/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApi 加载API
func loadAdmin(s *sinking_web.Engine) {
	g := service.Auth.GetAdminRouteGroup(s)
	service.Auth.LoadAdminRouteWithRegister(g, func(ANY func(string, string, func(*context.Context))) {
		ANY("/list", "任务列表", controller.List)
		ANY("/update", "更新任务", controller.Update)
		ANY("/delete", "删除任务", controller.Delete)
		ANY("/create", "创建任务", controller.Create)
		ANY("/info", "任务详情", controller.Info)
		ANY("/run", "执行任务", controller.Run)
		ANY("/log", "任务日志", controller.Log)
	})
}

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
		ANY("/task/list", "任务列表", controller.Admin.Task.List)
		ANY("/task/update", "更新任务", controller.Admin.Task.Update)
		ANY("/task/delete", "删除任务", controller.Admin.Task.Delete)
		ANY("/task/create", "创建任务", controller.Admin.Task.Create)
		ANY("/task/info", "任务详情", controller.Admin.Task.Info)
		ANY("/task/run", "执行任务", controller.Admin.Task.Run)
		ANY("/task/log", "任务日志", controller.Admin.Task.Log)
	})
}

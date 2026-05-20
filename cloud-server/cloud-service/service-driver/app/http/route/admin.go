package route

import (
	"service-driver/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApi 加载API
func loadAdmin(s *sinking_web.Engine) {
	g := service.Auth.GetAdminRouteGroup(s)
	service.Auth.LoadAdminRouteWithRegister(g, func(ANY func(string, string, func(*context.Context))) {
		ANY("/channel/list", "存储列表", controller.Admin.Channel.List)
		ANY("/channel/info", "存储详情", controller.Admin.Channel.Info)
		ANY("/channel/create", "创建存储", controller.Admin.Channel.Create)
		ANY("/channel/delete", "删除存储", controller.Admin.Channel.Delete)
		ANY("/channel/update", "更新存储", controller.Admin.Channel.Update)

		ANY("/scene/list", "场景列表", controller.Admin.Scene.List)
		ANY("/scene/create", "创建场景", controller.Admin.Scene.Create)
		ANY("/scene/delete", "删除场景", controller.Admin.Scene.Delete)
		ANY("/scene/update", "更新场景", controller.Admin.Scene.Update)

		ANY("/file/list", "文件列表", controller.Admin.File.List)
		ANY("/file/delete", "删除文件", controller.Admin.File.Delete)
		ANY("/file/update", "修改文件", controller.Admin.File.Update)
	})
}

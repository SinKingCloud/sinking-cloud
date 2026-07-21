package route

import (
	"service-space/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadUser 加载user控制器
func loadUser(s *sinking_web.Engine) {
	g := service.Auth.GetUserRouteGroup(s)
	mind := g.Group("/mind")
	{
		mindData := mind.Group("/data")
		{
			mindData.ANY("/list", context.HandleFunc(controller.User.MindData.List))     //导图数据列表
			mindData.ANY("/info", context.HandleFunc(controller.User.MindData.Info))     //导图数据详情
			mindData.ANY("/create", context.HandleFunc(controller.User.MindData.Create)) //创建导图数据
			mindData.ANY("/update", context.HandleFunc(controller.User.MindData.Update)) //更新导图数据
			mindData.ANY("/delete", context.HandleFunc(controller.User.MindData.Delete)) //删除导图数据
			mindData.ANY("/copy", context.HandleFunc(controller.User.MindData.Copy))     //复制导图数据
			mindData.ANY("/save", context.HandleFunc(controller.User.MindData.Save))     //保存导图数据
		}

		mindType := mind.Group("/type")
		{
			mindType.ANY("/list", context.HandleFunc(controller.User.MindType.List))     //导图分类列表
			mindType.ANY("/create", context.HandleFunc(controller.User.MindType.Create)) //创建导图分类
			mindType.ANY("/update", context.HandleFunc(controller.User.MindType.Update)) //更新导图分类
			mindType.ANY("/delete", context.HandleFunc(controller.User.MindType.Delete)) //删除导图分类
		}

		version := mind.Group("/version")
		{
			version.ANY("/list", context.HandleFunc(controller.User.MindVersion.List))       //导图版本列表
			version.ANY("/info", context.HandleFunc(controller.User.MindVersion.Info))       //导图版本详情
			version.ANY("/restore", context.HandleFunc(controller.User.MindVersion.Restore)) //恢复导图版本
		}

		share := mind.Group("/share")
		{
			share.ANY("/list", context.HandleFunc(controller.User.MindShare.List))     //导图分享列表
			share.ANY("/create", context.HandleFunc(controller.User.MindShare.Create)) //创建导图分享
			share.ANY("/update", context.HandleFunc(controller.User.MindShare.Update)) //更新导图分享
			share.ANY("/delete", context.HandleFunc(controller.User.MindShare.Delete)) //删除导图分享
		}
	}
}

package admin

import (
	"service-driver/app/http/controller/admin/channel"
	"service-driver/app/http/controller/admin/file"
	"service-driver/app/http/controller/admin/scene"
)

type Controller struct {
	Channel channel.Controller
	File    file.Controller
	Scene   scene.Controller
}

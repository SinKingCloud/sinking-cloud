package user_app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service"
)

type ControllerUserApp struct {
}

// FindById 通过id查询
func (ControllerUserApp) FindById(c *sinking_server.Context) {
	data := c.GetRpcData()
	id := int(data["id"].(float64))
	app, err := service.WebApp.FindByIdWithCache(id, true)
	if err != nil || app == nil {
		c.Error("fail")
	} else {
		c.SuccessWithData("success", app)
	}
}

package rbac

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_rbac"
)

type ControllerRbac struct {
}

// Check 判断rbac
func (ControllerRbac) Check(c *sinking_server.Context) {
	data := c.GetRpcData()
	userId := int(data["user_id"].(float64))
	roleId := int(data["role_id"].(float64))
	r := system_rbac.GetIns().Check(userId, roleId, data["path"].(string), data["method"].(string))
	c.SuccessWithData("success", r)
}

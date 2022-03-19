package sinking_server

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/validator"

// GetRpcData rpc参数获取
func (c *Context) GetRpcData() map[string]interface{} {
	type Rpc struct {
		Group string                 `json:"group" validate:"required" label:"服务分组"`
		Name  string                 `json:"name" validate:"required" label:"服务名"`
		Data  map[string]interface{} `json:"data" validate:"-" label:"数据"`
	}
	rpc2 := &Rpc{}
	if c.BindJson(rpc2) != nil {
		c.Error("rpc参数绑定失败")
		c.Abort()
		return nil
	}
	a, b := validator.Check(rpc2)
	if !a {
		c.Error(b)
		c.Abort()
		return nil
	}
	return rpc2.Data
}

package test

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"

func Index(c *sinking_server.Context) {
	c.Success("请求成功")
}

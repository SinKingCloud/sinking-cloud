package upload

import (
	"service-driver/app/service"
	"strconv"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Simple 简单上传
func (Controller) Simple(c *context.Context) {
	sceneIdText := c.Request.FormValue("scene_id")
	if sceneIdText == "" {
		c.Error("场景ID不能为空")
		return
	}
	sceneId, err := strconv.ParseInt(sceneIdText, 10, 64)
	if err != nil || sceneId <= 0 {
		c.Error("场景ID不合法")
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.Error("上传文件不能为空")
		return
	}
	defer func() {
		_ = file.Close()
	}()
	var isLogin bool
	u, err := service2.User.GetJwtUserInfo(c)
	if err == nil && u != nil && u.Id > 0 {
		isLogin = true
	}
	url, err := service.File.SimpleUpload(c.Request.Context(), isLogin, sceneId, file, header.Filename, header.Size)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("上传成功", url)
}

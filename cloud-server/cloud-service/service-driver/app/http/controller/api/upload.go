package api

import (
	"service-driver/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

type uploadController struct {
}

// Simple 简单上传
func (uploadController) Simple(c *context.Context) {
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
	url, err := service.File.SimpleUpload(c.Request.Context(), sceneId, file, header.Filename, header.Size)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("上传成功", url)
}

// Chunk 分片上传
func (uploadController) Chunk(c *context.Context) {

}

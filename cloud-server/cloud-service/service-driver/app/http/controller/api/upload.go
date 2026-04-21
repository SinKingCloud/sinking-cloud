package api

import (
	"service-driver/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

type uploadController struct {
}

// Simple 简单上传
func (uploadController) Simple(c *context.Context) {
	var form struct {
		SceneId int64 `json:"scene_id" default:"0" validate:"required,min=1" label:"场景ID"`
	}
	if ok, msg := c.ValidatorPost(&form); !ok {
		c.Error(msg)
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
	url, err := service.File.SimpleUpload(c.Request.Context(), form.SceneId, file, header.Filename, header.Size)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("上传成功", url)
}

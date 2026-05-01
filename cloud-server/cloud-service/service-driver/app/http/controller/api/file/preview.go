package file

import (
	"net/http"
	"service-driver/app/service"
	"strconv"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Preview 文件预览
func (Controller) Preview(c *context.Context) {
	var form struct {
		SceneId int64  `json:"scene_id" default:"" validate:"required,min=1" label:"场景ID"`
		Hash    string `json:"hash" default:"" validate:"required,len=32" label:"文件哈希"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	url, err := service.File.GetUrlByHash(form.SceneId, form.Hash)
	if err != nil {
		c.Error(err.Error())
		return
	}
	second := 7 * 24 * 3600
	expire := time.Now().Add(time.Duration(second) * time.Second).UTC()
	c.SetHeader("Cache-Control", "public, max-age="+strconv.Itoa(second))
	c.SetHeader("Expires", expire.Format(http.TimeFormat))
	c.Redirect(http.StatusFound, url)
}

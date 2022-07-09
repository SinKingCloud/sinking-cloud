package upload

import (
	"encoding/base64"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/aliyun"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"time"
)

type ControllerUpload struct {
}

// Temp 上传临时文件
func (ControllerUpload) Temp(c *sinking_server.Context) {
	data := c.GetRpcData()
	configs := service.GetSystemConfigServiceIns().GetByCache(false)
	fileName := data["name"].(string)
	fileData := data["data"].(string)
	fileByte, err := base64.StdEncoding.DecodeString(fileData)
	if err != nil {
		c.Error("error")
		return
	}
	savePath := configs[config.AliyunOssTempSavePath] + "/" + time.Now().Format("20060102") + "/" + fileName
	if aliyun.Oss.Upload(
		configs[config.AliyunAccessKey],
		configs[config.AliyunAccessSecret],
		configs[config.AliyunOssEndPoint],
		configs[config.AliyunOssBuckName],
		fileByte,
		savePath,
	) {
		c.SuccessWithData("success", configs[config.AliyunOssDomain]+savePath)
	} else {
		c.Error("error")
	}
}

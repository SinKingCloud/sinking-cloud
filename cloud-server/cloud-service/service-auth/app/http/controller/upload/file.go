package upload

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/aliyun"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"path"
	"strconv"
	"strings"
	"time"
)

// File 上传文件
func (ControllerUpload) File(c *sinking_server.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.Error("请上传文件")
		return
	}
	ext := path.Ext(file.Filename)
	configs := service.GetSystemConfigServiceIns().GetByCache(false)
	if !strings.Contains(configs[config.AliyunOssAllowExt], ext) {
		c.Error("该文件类型不允许上传")
		return
	}
	f, err := file.Open()
	if err != nil {
		c.Error("打开文件失败")
		return
	}
	data := make([]byte, file.Size)
	count, err := f.Read(data)
	if err != nil {
		c.Error("读取文件失败")
		return
	}
	maxSize, err := strconv.Atoi(configs[config.AliyunOssMaxSize])
	if err != nil || count/1000 > maxSize {
		c.Error("超过最大上传文件大小")
		return
	}
	hash := str.NetByteTool(data).Md5()
	savePath := configs[config.AliyunOssSavePath] + "/" + time.Now().Format("20060102") + "/" + hash + ext
	if aliyun.Oss.Upload(
		configs[config.AliyunAccessKey],
		configs[config.AliyunAccessSecret],
		configs[config.AliyunOssEndPoint],
		configs[config.AliyunOssBuckName],
		data,
		savePath,
	) {
		type File struct {
			Name       string `json:"name"`
			Url        string `json:"url"`
			Hash       string `json:"hash"`
			Extensions string `json:"extensions"`
			Size       int    `json:"size"`
		}
		c.SuccessWithData("上传文件成功", &File{
			Name:       file.Filename,
			Url:        configs[config.AliyunOssDomain] + savePath,
			Hash:       hash,
			Extensions: ext,
			Size:       count,
		})
	} else {
		c.Error("上传文件失败")
	}
}

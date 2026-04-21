package api

import (
	"io"
	"service-driver/app/service"
	fileService "service-driver/app/service/file"
	"strconv"
	"strings"

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
	step := strings.TrimSpace(strings.ToLower(c.Request.FormValue("step")))
	if step == "" {
		c.Error("分片上传步骤不能为空")
		return
	}
	hash := c.Request.FormValue("hash")
	if len(hash) != 32 {
		c.Error("文件哈希不合法")
		return
	}
	filename := c.Request.FormValue("name")
	uploadId := c.Request.FormValue("upload_id")
	totalSize := int64(0)
	chunkSize := int64(0)
	partNumber := 0
	var (
		body       io.Reader
		headerSize int64
	)
	switch step {
	case fileService.ChunkUploadStepUpload:
		totalSizeText := c.Request.FormValue("total_size")
		if totalSizeText == "" {
			c.Error("文件总大小不能为空")
			return
		}
		totalSize, err = strconv.ParseInt(totalSizeText, 10, 64)
		if err != nil || totalSize <= 0 {
			c.Error("文件总大小不合法")
			return
		}
		partNumberText := c.Request.FormValue("part_number")
		if partNumberText == "" {
			c.Error("分片序号不能为空")
			return
		}
		partNumber, err = strconv.Atoi(partNumberText)
		if err != nil || partNumber <= 0 {
			c.Error("分片序号不合法")
			return
		}
		chunkSizeText := c.Request.FormValue("chunk_size")
		if chunkSizeText != "" {
			chunkSize, err = strconv.ParseInt(chunkSizeText, 10, 64)
			if err != nil || chunkSize <= 0 {
				c.Error("分片大小不合法")
				return
			}
		}
		formFile, header, err := c.Request.FormFile("file")
		if err != nil {
			c.Error("上传文件不能为空")
			return
		}
		defer func() {
			_ = formFile.Close()
		}()
		body = formFile
		headerSize = header.Size
		if filename == "" {
			filename = header.Filename
		}
	case fileService.ChunkUploadStepStatus:
		if strings.TrimSpace(uploadId) == "" {
			c.Error("上传任务ID不能为空")
			return
		}
	case fileService.ChunkUploadStepFinish:
		if strings.TrimSpace(uploadId) == "" {
			c.Error("上传任务ID不能为空")
			return
		}
		if strings.TrimSpace(filename) == "" {
			c.Error("上传文件名称不能为空")
			return
		}
	default:
		c.Error("分片上传步骤不合法")
		return
	}
	data, err := service.File.ChunkUpload(c.Request.Context(), step, sceneId, uploadId, hash, filename, totalSize, chunkSize, partNumber, body, headerSize)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("上传成功", data)
}

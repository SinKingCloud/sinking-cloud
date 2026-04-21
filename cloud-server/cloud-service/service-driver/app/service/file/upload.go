package file

import (
	"bytes"
	"context"
	"errors"
	"io"
	"path/filepath"
	"service-driver/app/enum/channel_status"
	"service-driver/app/enum/file_status"
	"service-driver/app/enum/scene_status"
	"service-driver/app/model"
	"strconv"
	"strings"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// SimpleUpload 简单上传文件。
func (s *service) SimpleUpload(ctx context.Context, sceneId int64, body io.Reader, filename string, size int64) (string, error) {
	if sceneId <= 0 {
		return "", errors.New("上传场景不能为空")
	}
	if body == nil {
		return "", errors.New("上传文件不能为空")
	}
	filename = strings.TrimSpace(filepath.Base(filename))
	if filename == "" {
		return "", errors.New("上传文件名称不能为空")
	}
	sceneInfo, err := s.sceneService.FindById(sceneId, false)
	if err != nil {
		return "", errors.Join(errors.New("获取上传场景失败"), err)
	}
	if sceneInfo == nil {
		return "", errors.New("上传场景不存在")
	}
	if sceneInfo.Status != scene_status.Enable {
		return "", errors.New("上传场景不可用")
	}
	channelInfo, err := s.channelService.FindById(sceneInfo.ChannelId, false)
	if err != nil {
		return "", errors.Join(errors.New("获取存储通道失败"), err)
	}
	if channelInfo == nil {
		return "", errors.New("存储通道不存在")
	}
	if channelInfo.Status != channel_status.Enable {
		return "", errors.New("存储通道不可用")
	}
	if sceneInfo.MaxSize > 0 && size > 0 && size > sceneInfo.MaxSize {
		return "", errors.New("上传文件大小超过场景限制")
	}
	content, err := io.ReadAll(body)
	if err != nil {
		return "", errors.Join(errors.New("读取上传文件失败"), err)
	}
	fileSize := int64(len(content))
	if sceneInfo.MaxSize > 0 && fileSize > sceneInfo.MaxSize {
		return "", errors.New("上传文件大小超过场景限制")
	}
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")
	if !s.allowUploadExt(sceneInfo.AllowExt, ext) {
		return "", errors.New("上传文件类型不允许")
	}
	hash, err := str.NewByteTool(bytes.NewReader(content)).Md5()
	if err != nil {
		return "", errors.New("生成文件哈希失败")
	}
	exists, err := s.findByHashWithSceneId(hash, sceneInfo.Id, false)
	if err != nil {
		return "", errors.Join(errors.New("获取文件信息失败"), err)
	}
	if exists != nil {
		return s.buildPreviewUrl(sceneInfo.Id, hash), nil
	}
	client, err := s.channelService.GetStorage(sceneInfo.ChannelId)
	if err != nil {
		return "", errors.Join(errors.New("获取存储信息失败"), err)
	}
	path := s.buildUploadPath(hash, ext)
	_, err = client.Put(ctx, path, bytes.NewReader(content))
	if err != nil {
		return "", errors.Join(errors.New("上传文件失败"), err)
	}
	data := &model.File{
		ChannelId: sceneInfo.ChannelId,
		SceneId:   sceneInfo.Id,
		Name:      filename,
		Ext:       ext,
		Hash:      hash,
		Size:      fileSize,
		Path:      path,
		Status:    file_status.Enable,
	}
	err = s.repositoryFile.Create(data)
	if err != nil {
		if err2 := client.Delete(ctx, path); err2 != nil {
			return "", errors.Join(errors.New("创建文件记录失败"), err, errors.New("清理已上传文件失败"), err2)
		}
		return "", errors.Join(errors.New("创建文件记录失败"), err)
	}
	s.clearCache([]*model.File{data})
	return s.buildPreviewUrl(sceneInfo.Id, hash), nil
}

// allowUploadExt 判断场景是否允许当前扩展名。
func (s *service) allowUploadExt(allowExt string, ext string) bool {
	allowExt = strings.TrimSpace(strings.ToLower(allowExt))
	if allowExt == "" || allowExt == "*" {
		return true
	}
	if ext == "" {
		return false
	}
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	fields := strings.FieldsFunc(allowExt, func(r rune) bool {
		switch r {
		case ',', '，', ';', '；', '|', '、', ' ', '\n', '\r', '\t':
			return true
		default:
			return false
		}
	})
	for _, item := range fields {
		item = strings.TrimSpace(strings.TrimPrefix(strings.ToLower(item), "."))
		if item == "*" || item == ext {
			return true
		}
	}
	return false
}

// buildUploadPath 生成文件上传路径。
func (s *service) buildUploadPath(hash string, ext string) string {
	name := hash
	if ext != "" {
		name += "." + ext
	}
	return strings.Join([]string{
		time.Now().Format("20060102"),
		name,
	}, "/")
}

// buildPreviewUrl 生成文件预览地址。
func (s *service) buildPreviewUrl(sceneId int64, hash string) string {
	return "/preview/" + strconv.FormatInt(sceneId, 10) + "/" + hash
}

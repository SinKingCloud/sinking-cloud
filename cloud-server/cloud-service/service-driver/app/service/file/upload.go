package file

import (
	"bytes"
	"context"
	"errors"
	"io"
	"path/filepath"
	"service-driver/app/enum/channel_status"
	"service-driver/app/enum/file_status"
	"service-driver/app/enum/scene_login"
	"service-driver/app/enum/scene_status"
	"service-driver/app/model"
	"strconv"
	"strings"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/storage"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// SimpleUpload 简单上传文件。
func (s *service) SimpleUpload(ctx context.Context, isLogin bool, sceneId int64, body io.Reader, filename string, size int64) (string, error) {
	if sceneId <= 0 {
		return "", errors.New("上传场景不能为空")
	}
	if body == nil {
		return "", errors.New("上传文件不能为空")
	}
	sceneInfo, client, err := s.prepareUpload(sceneId, isLogin)
	if err != nil {
		return "", err
	}
	maxSizeBytes := s.sceneMaxSizeBytes(sceneInfo)
	if maxSizeBytes > 0 && size > 0 && size > maxSizeBytes {
		return "", errors.New("上传文件大小超过场景限制")
	}
	filename, ext, err := s.prepareUploadFilename(sceneInfo, filename)
	if err != nil {
		return "", err
	}
	content, err := io.ReadAll(body)
	if err != nil {
		return "", errors.Join(errors.New("读取上传文件失败"), err)
	}
	fileSize := int64(len(content))
	if maxSizeBytes > 0 && fileSize > maxSizeBytes {
		return "", errors.New("上传文件大小超过场景限制")
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

// sceneMaxSizeBytes 获取场景限制大小，场景配置单位为 KB。
func (s *service) sceneMaxSizeBytes(sceneInfo *model.Scene) int64 {
	if sceneInfo == nil || sceneInfo.MaxSize <= 0 {
		return 0
	}
	return sceneInfo.MaxSize * 1024
}

// ChunkUpload 分片上传文件。
func (s *service) ChunkUpload(ctx context.Context, isLogin bool, step string, sceneId int64, uploadId string, hash string, filename string, totalSize int64, chunkSize int64, partNumber int, body io.Reader, size int64) (*ChunkUploadResult, error) {
	step = strings.TrimSpace(strings.ToLower(step))
	if step == "" {
		return nil, errors.New("分片上传步骤不能为空")
	}
	if sceneId <= 0 {
		return nil, errors.New("上传场景不能为空")
	}
	hash = strings.TrimSpace(strings.ToLower(hash))
	if len(hash) != 32 {
		return nil, errors.New("文件哈希不合法")
	}
	sceneInfo, client, err := s.prepareUpload(sceneId, isLogin)
	if err != nil {
		return nil, err
	}
	switch step {
	case ChunkUploadStepInit:
		return s.chunkUploadInit(ctx, client, sceneInfo, hash, filename, totalSize, chunkSize)
	case ChunkUploadStepUpload:
		return s.chunkUploadPart(ctx, client, uploadId, partNumber, body, size)
	case ChunkUploadStepStatus:
		return s.chunkUploadStatus(ctx, client, uploadId)
	case ChunkUploadStepFinish:
		return s.chunkUploadFinish(ctx, client, sceneInfo, uploadId, hash, filename)
	default:
		return nil, errors.New("分片上传步骤不合法")
	}
}

// chunkUploadInit 初始化上传任务。
func (s *service) chunkUploadInit(ctx context.Context, client storage.Storage, sceneInfo *model.Scene, hash string, filename string, totalSize int64, chunkSize int64) (*ChunkUploadResult, error) {
	_, ext, err := s.prepareUploadFilename(sceneInfo, filename)
	if err != nil {
		return nil, err
	}
	if totalSize <= 0 {
		return nil, errors.New("文件总大小不合法")
	}
	maxSizeBytes := s.sceneMaxSizeBytes(sceneInfo)
	if maxSizeBytes > 0 && totalSize > maxSizeBytes {
		return nil, errors.New("上传文件大小超过场景限制")
	}
	exists, err := s.findByHashWithSceneId(hash, sceneInfo.Id, false)
	if err != nil {
		return nil, errors.Join(errors.New("获取文件信息失败"), err)
	}
	if exists != nil {
		return &ChunkUploadResult{
			Done:           true,
			Url:            s.buildPreviewUrl(sceneInfo.Id, hash),
			CommittedBytes: exists.Size,
		}, nil
	}
	uploadId, err := client.StartUpload(ctx, s.buildUploadPath(hash, ext), totalSize, chunkSize)
	if err != nil {
		return nil, errors.Join(errors.New("创建上传任务失败"), err)
	}
	return &ChunkUploadResult{
		UploadId: uploadId,
	}, nil
}

// prepareUpload 加载上传场景和存储客户端。
func (s *service) prepareUpload(sceneId int64, isLogin bool) (*model.Scene, storage.Storage, error) {
	sceneInfo, err := s.sceneService.FindById(sceneId, false)
	if err != nil {
		return nil, nil, errors.Join(errors.New("获取上传场景失败"), err)
	}
	if sceneInfo == nil {
		return nil, nil, errors.New("上传场景不存在")
	}
	if sceneInfo.Status != scene_status.Enable {
		return nil, nil, errors.New("上传场景不可用")
	}
	channelInfo, err := s.channelService.FindById(sceneInfo.ChannelId, false)
	if err != nil {
		return nil, nil, errors.Join(errors.New("获取存储通道失败"), err)
	}
	if channelInfo == nil {
		return nil, nil, errors.New("存储通道不存在")
	}
	if channelInfo.Status != channel_status.Enable {
		return nil, nil, errors.New("存储通道不可用")
	}
	if sceneInfo.Login == scene_login.Yes && !isLogin {
		return nil, nil, errors.New("您还未登录")
	}
	client, err := s.channelService.GetStorage(sceneInfo.ChannelId)
	if err != nil {
		return nil, nil, errors.Join(errors.New("获取存储信息失败"), err)
	}
	return sceneInfo, client, nil
}

// chunkUploadPart 上传单个分片。
func (s *service) chunkUploadPart(ctx context.Context, client storage.Storage, uploadId string, partNumber int, body io.Reader, size int64) (*ChunkUploadResult, error) {
	if body == nil {
		return nil, errors.New("上传文件不能为空")
	}
	if strings.TrimSpace(uploadId) == "" {
		return nil, errors.New("上传任务ID不能为空")
	}
	if size <= 0 {
		return nil, errors.New("分片大小不合法")
	}
	if partNumber <= 0 {
		return nil, errors.New("分片序号不合法")
	}
	if err := client.UploadPart(ctx, uploadId, partNumber, body, size); err != nil {
		return nil, errors.Join(errors.New("上传分片失败"), err)
	}
	return &ChunkUploadResult{
		UploadId: uploadId,
		Done:     false,
	}, nil
}

// chunkUploadStatus 查询上传状态。
func (s *service) chunkUploadStatus(ctx context.Context, client storage.Storage, uploadId string) (*ChunkUploadResult, error) {
	if strings.TrimSpace(uploadId) == "" {
		return nil, errors.New("上传任务ID不能为空")
	}
	state, err := client.GetUpload(ctx, uploadId)
	if err != nil {
		return nil, errors.Join(errors.New("获取上传状态失败"), err)
	}
	if state == nil {
		return nil, errors.New("获取上传状态失败")
	}
	return &ChunkUploadResult{
		UploadId:       uploadId,
		Done:           state.Done,
		UploadedParts:  state.UploadedParts,
		CommittedBytes: state.CommittedBytes,
	}, nil
}

// chunkUploadFinish 完成上传并创建文件记录。
func (s *service) chunkUploadFinish(ctx context.Context, client storage.Storage, sceneInfo *model.Scene, uploadId string, hash string, filename string) (*ChunkUploadResult, error) {
	if strings.TrimSpace(uploadId) == "" {
		return nil, errors.New("上传任务ID不能为空")
	}
	filename, ext, err := s.prepareUploadFilename(sceneInfo, filename)
	if err != nil {
		return nil, err
	}
	info, err := client.FinishUpload(ctx, uploadId)
	if err != nil {
		return nil, errors.Join(errors.New("完成上传失败"), err)
	}
	if info == nil {
		return nil, errors.New("完成上传失败")
	}
	exists, err := s.findByHashWithSceneId(hash, sceneInfo.Id, true)
	if err != nil {
		return nil, errors.Join(errors.New("获取文件信息失败"), err)
	}
	if exists == nil {
		data := &model.File{
			ChannelId: sceneInfo.ChannelId,
			SceneId:   sceneInfo.Id,
			Name:      filename,
			Ext:       ext,
			Hash:      hash,
			Size:      info.Size,
			Path:      info.Key,
			Status:    file_status.Enable,
		}
		err = s.repositoryFile.Create(data)
		if err != nil {
			return nil, errors.Join(errors.New("创建文件记录失败"), err)
		}
		s.clearCache([]*model.File{data})
	} else {
		s.clearCache([]*model.File{exists})
	}
	return &ChunkUploadResult{
		UploadId:       uploadId,
		Done:           true,
		Url:            s.buildPreviewUrl(sceneInfo.Id, hash),
		CommittedBytes: info.Size,
	}, nil
}

// prepareUploadFilename 规范化文件名并校验扩展名。
func (s *service) prepareUploadFilename(sceneInfo *model.Scene, filename string) (string, string, error) {
	filename = strings.TrimSpace(filepath.Base(filename))
	if filename == "" {
		return "", "", errors.New("上传文件名称不能为空")
	}
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")
	if !s.allowUploadExt(sceneInfo.AllowExt, ext) {
		return "", "", errors.New("上传文件类型不允许")
	}
	return filename, ext, nil
}

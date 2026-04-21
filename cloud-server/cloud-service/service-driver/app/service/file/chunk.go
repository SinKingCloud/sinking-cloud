package file

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"service-driver/app/enum/channel_status"
	"service-driver/app/enum/file_status"
	"service-driver/app/enum/scene_status"
	"service-driver/app/model"
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/storage"
)

// ChunkUpload 分片上传文件。
func (s *service) ChunkUpload(ctx context.Context, step string, sceneId int64, uploadId string, hash string, filename string, totalSize int64, chunkSize int64, partNumber int, body io.Reader, size int64) (*ChunkUploadResult, error) {
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
	sceneInfo, client, err := s.prepareChunkUpload(sceneId)
	if err != nil {
		return nil, err
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
	switch step {
	case ChunkUploadStepUpload:
		return s.chunkUploadPart(ctx, client, sceneInfo, uploadId, hash, filename, totalSize, chunkSize, partNumber, body, size)
	case ChunkUploadStepStatus:
		return s.chunkUploadStatus(ctx, client, uploadId)
	case ChunkUploadStepFinish:
		return s.chunkUploadFinish(ctx, client, sceneInfo, uploadId, hash, filename)
	default:
		return nil, errors.New("分片上传步骤不合法")
	}
}

// prepareChunkUpload 加载分片上传场景和存储客户端。
func (s *service) prepareChunkUpload(sceneId int64) (*model.Scene, storage.Storage, error) {
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
	client, err := s.channelService.GetStorage(sceneInfo.ChannelId)
	if err != nil {
		return nil, nil, errors.Join(errors.New("获取存储信息失败"), err)
	}
	return sceneInfo, client, nil
}

// chunkUploadPart 上传单个分片。
func (s *service) chunkUploadPart(ctx context.Context, client storage.Storage, sceneInfo *model.Scene, uploadId string, hash string, filename string, totalSize int64, chunkSize int64, partNumber int, body io.Reader, size int64) (*ChunkUploadResult, error) {
	if body == nil {
		return nil, errors.New("上传文件不能为空")
	}
	_, ext, err := s.prepareChunkFilename(sceneInfo, filename)
	if err != nil {
		return nil, err
	}
	if totalSize <= 0 {
		return nil, errors.New("文件总大小不合法")
	}
	if sceneInfo.MaxSize > 0 && totalSize > sceneInfo.MaxSize {
		return nil, errors.New("上传文件大小超过场景限制")
	}
	if size <= 0 {
		return nil, errors.New("分片大小不合法")
	}
	if partNumber <= 0 {
		return nil, errors.New("分片序号不合法")
	}
	path := s.buildUploadPath(hash, ext)
	if chunkSize <= 0 {
		chunkSize = size
	}
	if uploadId == "" {
		uploadId, err = client.StartUpload(ctx, path, totalSize, chunkSize)
		if err != nil {
			return nil, errors.Join(errors.New("创建上传任务失败"), err)
		}
	}
	if err = client.UploadPart(ctx, uploadId, partNumber, body, size); err != nil {
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
	filename, ext, err := s.prepareChunkFilename(sceneInfo, filename)
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

// prepareChunkFilename 规范化文件名并校验扩展名。
func (s *service) prepareChunkFilename(sceneInfo *model.Scene, filename string) (string, string, error) {
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

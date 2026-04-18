package storage

import (
	"context"
	"io"
	"time"
)

// Storage 定义对业务层暴露的极简存储接口。
type Storage interface {
	// Put 以普通方式上传一个对象。
	Put(ctx context.Context, key string, body io.Reader) (*ObjectInfo, error)
	// UploadLarge 以自动分片方式上传一个大对象。
	UploadLarge(ctx context.Context, key string, body io.Reader, size, chunkSize int64) (*ObjectInfo, error)
	// StartUpload 创建一个可续传的分片上传任务。
	StartUpload(ctx context.Context, key string, size, chunkSize int64) (string, error)
	// UploadPart 上传指定上传任务的单个分片。
	UploadPart(ctx context.Context, uploadID string, partNumber int, body io.Reader, size int64) error
	// GetUpload 查询上传任务当前状态。
	GetUpload(ctx context.Context, uploadID string) (*UploadState, error)
	// FinishUpload 完成一个已经上传完分片的任务。
	FinishUpload(ctx context.Context, uploadID string) (*ObjectInfo, error)
	// CancelUpload 取消一个未完成的上传任务。
	CancelUpload(ctx context.Context, uploadID string) error
	// Delete 删除单个对象。
	Delete(ctx context.Context, key string) error
	// DeleteMany 批量删除多个对象。
	DeleteMany(ctx context.Context, keys []string) error
	// Stat 获取对象元信息。
	Stat(ctx context.Context, key string) (*ObjectInfo, error)
	// Exists 判断对象是否存在。
	Exists(ctx context.Context, key string) (bool, error)
	// List 按前缀分页列出对象。
	List(ctx context.Context, request ListRequest) (*ListResult, error)
	// Copy 在存储侧复制一个对象。
	Copy(ctx context.Context, sourceKey, destinationKey string) (*CopyResult, error)
	// Move 移动对象到新的逻辑路径。
	Move(ctx context.Context, sourceKey, destinationKey string) (*CopyResult, error)
	// PublicURL 生成对象的公开访问地址。
	PublicURL(key string) (string, error)
	// SignedURL 生成对象的临时签名访问地址。
	SignedURL(ctx context.Context, key string, expire time.Duration) (string, error)
}

// Put 以最简方式上传一个对象，只需要 key 和内容流。
func (c *Client) Put(ctx context.Context, key string, body io.Reader) (*ObjectInfo, error) {
	return c.upload(ctx, uploadRequest{
		key:  key,
		body: body,
	})
}

// UploadLarge 以极简方式上传大对象，内部自动完成分片上传和合并。
//
// size 传实际文件大小，未知时可传 0。
// chunkSize 传单个分片大小，传 0 时自动使用驱动默认策略。
func (c *Client) UploadLarge(ctx context.Context, key string, body io.Reader, size, chunkSize int64) (*ObjectInfo, error) {
	return c.uploadLarge(ctx, largeUploadRequest{
		uploadRequest: uploadRequest{
			key:  key,
			body: body,
			size: size,
		},
		chunkSize: chunkSize,
	})
}

// StartUpload 初始化一个分片上传任务，并返回可持久化的字符串 uploadID。
func (c *Client) StartUpload(ctx context.Context, key string, size, chunkSize int64) (string, error) {
	session, err := c.startUploadSession(ctx, startUploadSessionRequest{
		key:       key,
		size:      size,
		chunkSize: chunkSize,
	})
	if err != nil {
		return "", err
	}

	return c.encodeUploadID(*session)
}

// UploadPart 上传一个分片，只需要 uploadID、分片序号和内容。
func (c *Client) UploadPart(ctx context.Context, uploadID string, partNumber int, body io.Reader, size int64) error {
	session, err := c.decodeUploadID(uploadID)
	if err != nil {
		return err
	}
	if partNumber <= 0 || size <= 0 {
		return errInvalidUploadChunk
	}

	offset := int64(partNumber-1) * session.chunkSize
	isFinal := false
	if session.size > 0 {
		isFinal = offset+size >= session.size
	}

	_, err = c.uploadChunk(ctx, uploadChunkRequest{
		session: session,
		index:   partNumber,
		offset:  offset,
		size:    size,
		chunkID: buildChunkID(partNumber),
		isFinal: isFinal,
		body:    body,
	})
	return err
}

// GetUpload 查询上传任务状态。
func (c *Client) GetUpload(ctx context.Context, uploadID string) (*UploadState, error) {
	session, err := c.decodeUploadID(uploadID)
	if err != nil {
		return nil, err
	}

	status, err := c.getUploadStatus(ctx, session)
	if err != nil {
		return nil, err
	}

	uploadedParts := make([]int, 0, len(status.uploadedChunks))
	for _, chunk := range status.uploadedChunks {
		if chunk.index > 0 {
			uploadedParts = append(uploadedParts, chunk.index)
		}
	}

	return &UploadState{
		ID:             uploadID,
		Key:            session.key,
		Size:           session.size,
		ChunkSize:      session.chunkSize,
		UploadedParts:  uploadedParts,
		CommittedBytes: status.committedBytes,
		Done:           status.finalized,
	}, nil
}

// FinishUpload 完成一个分片上传任务。
func (c *Client) FinishUpload(ctx context.Context, uploadID string) (*ObjectInfo, error) {
	session, err := c.decodeUploadID(uploadID)
	if err != nil {
		return nil, err
	}

	request := completeUploadRequest{
		session: session,
	}
	if c.capabilities().uploadStatusQuery {
		status, err := c.getUploadStatus(ctx, session)
		if err != nil {
			return nil, err
		}
		request.chunks = status.uploadedChunks
	}

	return c.completeUpload(ctx, request)
}

// CancelUpload 取消一个未完成的分片上传任务。
func (c *Client) CancelUpload(ctx context.Context, uploadID string) error {
	session, err := c.decodeUploadID(uploadID)
	if err != nil {
		return err
	}

	return c.abortUpload(ctx, session)
}

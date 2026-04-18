package storage

import (
	"context"
	"io"
	"time"
)

// Storage 定义对业务层暴露的极简存储接口。
type Storage interface {
	Put(ctx context.Context, key string, body io.Reader) (*ObjectInfo, error)
	StartUpload(ctx context.Context, key string, size, chunkSize int64) (string, error)
	UploadPart(ctx context.Context, uploadID string, partNumber int, body io.Reader, size int64) error
	GetUpload(ctx context.Context, uploadID string) (*UploadState, error)
	FinishUpload(ctx context.Context, uploadID string) (*ObjectInfo, error)
	CancelUpload(ctx context.Context, uploadID string) error
	Delete(ctx context.Context, key string) error
	DeleteMany(ctx context.Context, keys []string) error
	Stat(ctx context.Context, key string) (*ObjectInfo, error)
	Exists(ctx context.Context, key string) (bool, error)
	List(ctx context.Context, request ListRequest) (*ListResult, error)
	Copy(ctx context.Context, sourceKey, destinationKey string) (*CopyResult, error)
	Move(ctx context.Context, sourceKey, destinationKey string) (*CopyResult, error)
	SignedURL(ctx context.Context, key string, expire time.Duration) (string, error)
}

// Put 以最简方式上传一个对象，只需要 key 和内容流。
func (c *Client) Put(ctx context.Context, key string, body io.Reader) (*ObjectInfo, error) {
	return c.upload(ctx, uploadRequest{
		key:  key,
		body: body,
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

	return encodeUploadID(*session)
}

// UploadPart 上传一个分片，只需要 uploadID、分片序号和内容。
func (c *Client) UploadPart(ctx context.Context, uploadID string, partNumber int, body io.Reader, size int64) error {
	session, err := decodeUploadID(uploadID)
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
	session, err := decodeUploadID(uploadID)
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
	session, err := decodeUploadID(uploadID)
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
	session, err := decodeUploadID(uploadID)
	if err != nil {
		return err
	}

	return c.abortUpload(ctx, session)
}

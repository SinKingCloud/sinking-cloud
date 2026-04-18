package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
)

// upload 执行一次普通对象上传。
func (c *Client) upload(ctx context.Context, request uploadRequest) (*ObjectInfo, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if request.body == nil {
		return nil, errUploadBodyRequired
	}

	return c.driver.upload(ctx, request)
}

// uploadLarge 以服务端托管方式完成一次自动分片上传。
func (c *Client) uploadLarge(ctx context.Context, request largeUploadRequest) (_ *ObjectInfo, err error) {
	if err = checkContext(ctx); err != nil {
		return nil, err
	}
	if request.body == nil {
		return nil, errUploadBodyRequired
	}

	capabilities := c.capabilities()
	if !capabilities.largeUpload {
		return nil, errLargeUploadNotSupported
	}

	chunkSize := normalizeChunkSize(request.chunkSize, capabilities.chunkPolicy, request.size)
	if request.size > 0 && request.size <= chunkSize {
		return c.upload(ctx, request.uploadRequest)
	}

	session, err := c.startUploadSession(ctx, startUploadSessionRequest{
		key:                request.key,
		size:               request.size,
		chunkSize:          chunkSize,
		contentType:        request.contentType,
		cacheControl:       request.cacheControl,
		contentDisposition: request.contentDisposition,
		metadata:           request.metadata,
	})
	if err != nil {
		return nil, err
	}

	completed := false
	defer func() {
		if !completed {
			_ = c.abortUpload(context.Background(), *session)
		}
	}()

	buffer := make([]byte, int(chunkSize))
	uploadedChunks := make([]uploadedChunk, 0, int(maxInt64(1, capacityBySize(request.size, chunkSize))))
	index := 1
	var offset int64

	for {
		if err = checkContext(ctx); err != nil {
			return nil, err
		}

		n, readErr := io.ReadFull(request.body, buffer)
		if readErr == io.EOF && n == 0 {
			break
		}
		if readErr != nil && !errors.Is(readErr, io.ErrUnexpectedEOF) {
			return nil, fmt.Errorf("read upload chunk failed: %w", readErr)
		}
		if n == 0 {
			break
		}

		isFinal := false
		if request.size > 0 {
			isFinal = offset+int64(n) >= request.size
		} else if errors.Is(readErr, io.ErrUnexpectedEOF) {
			isFinal = true
		}

		chunk, uploadErr := c.uploadChunk(ctx, uploadChunkRequest{
			session: *session,
			index:   index,
			offset:  offset,
			size:    int64(n),
			isFinal: isFinal,
			body:    bytes.NewReader(buffer[:n]),
		})
		if uploadErr != nil {
			return nil, uploadErr
		}

		uploadedChunks = append(uploadedChunks, *chunk)
		index++
		offset += int64(n)

		if errors.Is(readErr, io.ErrUnexpectedEOF) {
			break
		}
	}

	if len(uploadedChunks) == 0 {
		return c.upload(ctx, request.uploadRequest)
	}

	result, err := c.completeUpload(ctx, completeUploadRequest{
		session: *session,
		chunks:  uploadedChunks,
	})
	if err != nil {
		return nil, err
	}

	completed = true
	return result, nil
}

// startUploadSession 初始化一个内部上传会话。
func (c *Client) startUploadSession(ctx context.Context, request startUploadSessionRequest) (*uploadSession, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if !c.capabilities().largeUpload {
		return nil, errLargeUploadNotSupported
	}

	request.chunkSize = normalizeChunkSize(request.chunkSize, c.capabilities().chunkPolicy, request.size)
	return c.driver.startUploadSession(ctx, request)
}

// uploadChunk 上传单个分片到当前驱动。
func (c *Client) uploadChunk(ctx context.Context, request uploadChunkRequest) (*uploadedChunk, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if !c.capabilities().largeUpload {
		return nil, errLargeUploadNotSupported
	}
	if request.body == nil {
		return nil, errUploadBodyRequired
	}

	return c.driver.uploadChunk(ctx, request)
}

// getUploadStatus 查询上传会话当前已提交的分片状态。
func (c *Client) getUploadStatus(ctx context.Context, session uploadSession) (*uploadStatus, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if !c.capabilities().uploadStatusQuery {
		return nil, errUploadStatusNotSupported
	}

	status, err := c.driver.getUploadStatus(ctx, session)
	if err != nil {
		return nil, err
	}

	status.uploadedChunks = sortUploadedChunks(status.uploadedChunks)
	return status, nil
}

// completeUpload 根据分片清单完成整个上传任务。
func (c *Client) completeUpload(ctx context.Context, request completeUploadRequest) (*ObjectInfo, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if !c.capabilities().largeUpload {
		return nil, errLargeUploadNotSupported
	}

	request.chunks = sortUploadedChunks(request.chunks)
	return c.driver.completeUpload(ctx, request)
}

// abortUpload 终止一个尚未完成的上传会话。
func (c *Client) abortUpload(ctx context.Context, session uploadSession) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if !c.capabilities().largeUpload {
		return errLargeUploadNotSupported
	}

	return c.driver.abortUpload(ctx, session)
}

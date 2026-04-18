package storage

import (
	"context"
	"time"
)

type driver interface {
	provider() Provider
	capabilities() capabilities
	upload(ctx context.Context, request uploadRequest) (*ObjectInfo, error)
	startUploadSession(ctx context.Context, request startUploadSessionRequest) (*uploadSession, error)
	uploadChunk(ctx context.Context, request uploadChunkRequest) (*uploadedChunk, error)
	getUploadStatus(ctx context.Context, session uploadSession) (*uploadStatus, error)
	completeUpload(ctx context.Context, request completeUploadRequest) (*ObjectInfo, error)
	abortUpload(ctx context.Context, session uploadSession) error
	delete(ctx context.Context, key string) error
	deleteMany(ctx context.Context, keys []string) error
	stat(ctx context.Context, key string) (*ObjectInfo, error)
	exists(ctx context.Context, key string) (bool, error)
	list(ctx context.Context, request ListRequest) (*ListResult, error)
	copy(ctx context.Context, sourceKey, destinationKey string) (*CopyResult, error)
	signedURL(ctx context.Context, key string, expire time.Duration) (string, error)
}

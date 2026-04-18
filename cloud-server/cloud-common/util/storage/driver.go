package storage

import (
	"context"
	"time"
)

// driver 定义底层厂商驱动需要实现的统一能力。
type driver interface {
	// provider 返回当前驱动所属的存储厂商。
	provider() Provider
	// capabilities 返回当前驱动支持的能力集合。
	capabilities() capabilities
	// upload 执行一次普通对象上传。
	upload(ctx context.Context, request uploadRequest) (*ObjectInfo, error)
	// startUploadSession 初始化一个分片上传会话。
	startUploadSession(ctx context.Context, request startUploadSessionRequest) (*uploadSession, error)
	// uploadChunk 上传分片数据。
	uploadChunk(ctx context.Context, request uploadChunkRequest) (*uploadedChunk, error)
	// getUploadStatus 查询当前上传会话的分片状态。
	getUploadStatus(ctx context.Context, session uploadSession) (*uploadStatus, error)
	// completeUpload 使用分片清单完成整个上传任务。
	completeUpload(ctx context.Context, request completeUploadRequest) (*ObjectInfo, error)
	// abortUpload 终止一个尚未完成的上传会话。
	abortUpload(ctx context.Context, session uploadSession) error
	// delete 删除单个对象。
	delete(ctx context.Context, key string) error
	// deleteMany 批量删除多个对象。
	deleteMany(ctx context.Context, keys []string) error
	// stat 获取对象元信息。
	stat(ctx context.Context, key string) (*ObjectInfo, error)
	// exists 判断对象是否存在。
	exists(ctx context.Context, key string) (bool, error)
	// list 按条件列出对象。
	list(ctx context.Context, request ListRequest) (*ListResult, error)
	// copy 在存储侧复制对象。
	copy(ctx context.Context, sourceKey, destinationKey string) (*CopyResult, error)
	// publicURL 生成对象的公开访问地址。
	publicURL(key string) (string, error)
	// signedURL 生成对象的临时签名访问地址。
	signedURL(ctx context.Context, key string, expire time.Duration) (string, error)
}

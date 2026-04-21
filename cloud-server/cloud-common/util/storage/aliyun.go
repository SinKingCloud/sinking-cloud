package storage

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// aliyunOSSDriver 是阿里云 OSS 的具体驱动实现。
type aliyunOSSDriver struct {
	client *oss.Client // 是阿里云 OSS SDK 客户端。
	bucket *oss.Bucket // 是已经绑定好的 Bucket 句柄。
	config Config      // 保存当前驱动初始化配置。
}

// provider 返回当前驱动对应的存储厂商。
func (d *aliyunOSSDriver) provider() Provider {
	return ProviderAliyunOSS
}

// capabilities 返回阿里云 OSS 驱动支持的能力集合。
func (d *aliyunOSSDriver) capabilities() capabilities {
	return capabilities{
		smallUpload:       true,
		largeUpload:       true,
		resumeUpload:      true,
		uploadStatusQuery: true,
		signedURL:         true,
		copyObject:        true,
		batchDelete:       true,
		protocols:         []uploadProtocol{uploadProtocolMultipart},
		chunkPolicy: chunkPolicy{
			defaultSize: defaultChunkSize,
			minSize:     minCompatibleChunkSize,
			maxChunks:   10000,
		},
	}
}

// upload 执行一次阿里云 OSS 普通对象上传。
func (d *aliyunOSSDriver) upload(ctx context.Context, request uploadRequest) (*ObjectInfo, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if request.body == nil {
		return nil, errUploadBodyRequired
	}

	providerKey, err := buildProviderKey(d.config.Prefix, request.key)
	if err != nil {
		return nil, err
	}

	options := d.buildObjectOptions(
		detectContentType(request.contentType, request.key),
		request.cacheControl,
		request.contentDisposition,
		request.metadata,
	)

	if err := d.bucket.PutObject(providerKey, request.body, options...); err != nil {
		return nil, fmt.Errorf("upload object failed: %w", err)
	}

	return d.statByProviderKey(ctx, providerKey)
}

// startUploadSession 初始化一个阿里云 OSS multipart 上传会话。
func (d *aliyunOSSDriver) startUploadSession(ctx context.Context, request startUploadSessionRequest) (*uploadSession, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	providerKey, err := buildProviderKey(d.config.Prefix, request.key)
	if err != nil {
		return nil, err
	}

	chunkSize := normalizeChunkSize(request.chunkSize, d.capabilities().chunkPolicy, request.size)
	contentType := detectContentType(request.contentType, request.key)
	options := d.buildObjectOptions(contentType, request.cacheControl, request.contentDisposition, request.metadata)

	result, err := d.bucket.InitiateMultipartUpload(providerKey, options...)
	if err != nil {
		return nil, fmt.Errorf("init upload session failed: %w", err)
	}

	return &uploadSession{
		provider:           ProviderAliyunOSS,
		protocol:           uploadProtocolMultipart,
		bucket:             d.config.Bucket,
		key:                request.key,
		providerKey:        providerKey,
		sessionID:          result.UploadID,
		size:               request.size,
		chunkSize:          chunkSize,
		contentType:        contentType,
		cacheControl:       request.cacheControl,
		contentDisposition: request.contentDisposition,
		metadata:           cloneMetadata(request.metadata),
		createdAt:          time.Now(),
	}, nil
}

// uploadChunk 上传 multipart 中的单个分片。
func (d *aliyunOSSDriver) uploadChunk(ctx context.Context, request uploadChunkRequest) (*uploadedChunk, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if err := d.validateUploadChunkRequest(request); err != nil {
		return nil, err
	}

	part, err := d.bucket.UploadPart(d.uploadSessionHandle(request.session), request.body, request.size, request.index)
	if err != nil {
		return nil, fmt.Errorf("upload chunk failed: %w", err)
	}

	return &uploadedChunk{
		index:        request.index,
		offset:       request.offset,
		size:         request.size,
		etag:         strings.Trim(part.ETag, "\""),
		lastModified: time.Now(),
	}, nil
}

// getUploadStatus 查询当前 multipart 会话下已上传的全部分片。
func (d *aliyunOSSDriver) getUploadStatus(ctx context.Context, session uploadSession) (*uploadStatus, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if err := d.validateUploadSession(session); err != nil {
		return nil, err
	}

	chunks, committedBytes, err := collectUploadedChunks(func(partNumberMarker int) (oss.ListUploadedPartsResult, error) {
		options := []oss.Option{
			oss.MaxParts(1000),
		}
		if partNumberMarker > 0 {
			options = append(options, oss.PartNumberMarker(partNumberMarker))
		}

		result, err := d.bucket.ListUploadedParts(d.uploadSessionHandle(session), options...)
		if err != nil {
			return oss.ListUploadedPartsResult{}, fmt.Errorf("get upload status failed: %w", err)
		}
		return result, nil
	})
	if err != nil {
		return nil, err
	}

	return &uploadStatus{
		session:        session,
		uploadedChunks: chunks,
		committedBytes: committedBytes,
	}, nil
}

// completeUpload 使用阿里云返回的分片清单完成 multipart 上传。
func (d *aliyunOSSDriver) completeUpload(ctx context.Context, request completeUploadRequest) (*ObjectInfo, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if err := d.validateUploadSession(request.session); err != nil {
		return nil, err
	}
	if len(request.chunks) == 0 {
		return nil, errInvalidCompleteUpload
	}

	parts := make([]oss.UploadPart, 0, len(request.chunks))
	for _, chunk := range request.chunks {
		if chunk.index <= 0 || strings.TrimSpace(chunk.etag) == "" {
			return nil, errInvalidCompleteUpload
		}
		parts = append(parts, oss.UploadPart{
			PartNumber: chunk.index,
			ETag:       chunk.etag,
		})
	}

	sort.Sort(oss.UploadParts(parts))
	if _, err := d.bucket.CompleteMultipartUpload(d.uploadSessionHandle(request.session), parts); err != nil {
		return nil, fmt.Errorf("complete upload failed: %w", err)
	}

	return d.statByProviderKey(ctx, request.session.providerKey)
}

// abortUpload 终止一个尚未完成的 multipart 上传会话。
func (d *aliyunOSSDriver) abortUpload(ctx context.Context, session uploadSession) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if err := d.validateUploadSession(session); err != nil {
		return err
	}
	if err := d.bucket.AbortMultipartUpload(d.uploadSessionHandle(session)); err != nil {
		return fmt.Errorf("abort upload failed: %w", err)
	}
	return nil
}

// delete 删除单个阿里云 OSS 对象。
func (d *aliyunOSSDriver) delete(ctx context.Context, key string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}

	providerKey, err := buildProviderKey(d.config.Prefix, key)
	if err != nil {
		return err
	}
	if err := d.bucket.DeleteObject(providerKey); err != nil {
		return fmt.Errorf("delete object failed: %w", err)
	}
	return nil
}

// deleteMany 批量删除多个阿里云 OSS 对象。
func (d *aliyunOSSDriver) deleteMany(ctx context.Context, keys []string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	const aliyunDeleteManyBatchSize = 1000
	providerKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		providerKey, err := buildProviderKey(d.config.Prefix, key)
		if err != nil {
			return err
		}
		providerKeys = append(providerKeys, providerKey)
	}
	num := len(providerKeys)
	for start := 0; start < num; start += aliyunDeleteManyBatchSize {
		if err := checkContext(ctx); err != nil {
			return err
		}
		end := start + aliyunDeleteManyBatchSize
		if end > num {
			end = num
		}
		if _, err := d.bucket.DeleteObjects(providerKeys[start:end]); err != nil {
			return fmt.Errorf("batch delete objects failed: %w", err)
		}
	}
	return nil
}

// stat 获取单个对象的元信息。
func (d *aliyunOSSDriver) stat(ctx context.Context, key string) (*ObjectInfo, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	providerKey, err := buildProviderKey(d.config.Prefix, key)
	if err != nil {
		return nil, err
	}

	return d.statByProviderKey(ctx, providerKey)
}

// exists 判断指定对象在阿里云 OSS 中是否存在。
func (d *aliyunOSSDriver) exists(ctx context.Context, key string) (bool, error) {
	if err := checkContext(ctx); err != nil {
		return false, err
	}

	providerKey, err := buildProviderKey(d.config.Prefix, key)
	if err != nil {
		return false, err
	}

	exists, err := d.bucket.IsObjectExist(providerKey)
	if err != nil {
		return false, fmt.Errorf("check object exists failed: %w", err)
	}
	return exists, nil
}

// list 按条件列出阿里云 OSS 中的对象。
func (d *aliyunOSSDriver) list(ctx context.Context, request ListRequest) (*ListResult, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	limit := request.Limit
	if limit <= 0 {
		limit = 100
	}

	options := []oss.Option{
		oss.Prefix(buildProviderPrefix(d.config.Prefix, request.Prefix)),
		oss.MaxKeys(limit),
	}

	if marker := normalizeListPrefix(request.Marker); marker != "" {
		options = append(options, oss.Marker(buildProviderPrefix(d.config.Prefix, marker)))
	}

	result, err := d.bucket.ListObjects(options...)
	if err != nil {
		return nil, fmt.Errorf("list objects failed: %w", err)
	}

	objects := make([]ObjectInfo, 0, len(result.Objects))
	for _, object := range result.Objects {
		objects = append(objects, d.objectFromListObject(object))
	}

	return &ListResult{
		Objects:    objects,
		NextMarker: stripProviderPrefix(d.config.Prefix, result.NextMarker),
		HasMore:    result.IsTruncated,
	}, nil
}

// copy 在阿里云 OSS 侧复制一个对象。
func (d *aliyunOSSDriver) copy(ctx context.Context, sourceKey, destinationKey string) (*CopyResult, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	srcProviderKey, err := buildProviderKey(d.config.Prefix, sourceKey)
	if err != nil {
		return nil, err
	}
	dstProviderKey, err := buildProviderKey(d.config.Prefix, destinationKey)
	if err != nil {
		return nil, err
	}

	if _, err := d.bucket.CopyObject(srcProviderKey, dstProviderKey); err != nil {
		return nil, fmt.Errorf("copy object failed: %w", err)
	}

	info, err := d.statByProviderKey(ctx, dstProviderKey)
	if err != nil {
		return nil, err
	}

	return &CopyResult{
		SourceKey:      sourceKey,
		DestinationKey: destinationKey,
		ETag:           info.ETag,
		CopiedAt:       time.Now(),
	}, nil
}

// publicURL 生成对象的公开访问地址。
func (d *aliyunOSSDriver) publicURL(key string) (string, error) {
	providerKey, err := buildProviderKey(d.config.Prefix, key)
	if err != nil {
		return "", err
	}

	return d.objectURL(providerKey), nil
}

// signedURL 生成对象的临时签名访问地址。
func (d *aliyunOSSDriver) signedURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	if err := checkContext(ctx); err != nil {
		return "", err
	}

	providerKey, err := buildProviderKey(d.config.Prefix, key)
	if err != nil {
		return "", err
	}

	signedURL, err := d.bucket.SignURL(providerKey, oss.HTTPGet, int64(expire.Seconds()))
	if err != nil {
		return "", fmt.Errorf("generate signed url failed: %w", err)
	}
	return signedURL, nil
}

// buildObjectOptions 将通用对象属性转换成阿里云 OSS SDK 选项。
func (d *aliyunOSSDriver) buildObjectOptions(contentType, cacheControl, contentDisposition string, metadata map[string]string) []oss.Option {
	options := make([]oss.Option, 0, len(metadata)+3)

	if strings.TrimSpace(contentType) != "" {
		options = append(options, oss.ContentType(contentType))
	}
	if strings.TrimSpace(cacheControl) != "" {
		options = append(options, oss.CacheControl(cacheControl))
	}
	if strings.TrimSpace(contentDisposition) != "" {
		options = append(options, oss.ContentDisposition(contentDisposition))
	}
	for key, value := range metadata {
		if strings.TrimSpace(key) == "" {
			continue
		}
		options = append(options, oss.Meta(key, value))
	}

	return options
}

// uploadSessionHandle 将内部上传会话转换成 OSS SDK 需要的句柄结构。
func (d *aliyunOSSDriver) uploadSessionHandle(session uploadSession) oss.InitiateMultipartUploadResult {
	return oss.InitiateMultipartUploadResult{
		Bucket:   d.config.Bucket,
		Key:      session.providerKey,
		UploadID: session.sessionID,
	}
}

// validateUploadSession 校验内部上传会话是否满足阿里云 OSS 调用要求。
func (d *aliyunOSSDriver) validateUploadSession(session uploadSession) error {
	if strings.TrimSpace(session.providerKey) == "" || strings.TrimSpace(session.sessionID) == "" {
		return errInvalidUploadSession
	}
	if session.protocol != "" && session.protocol != uploadProtocolMultipart {
		return errInvalidUploadSession
	}
	return nil
}

// validateUploadChunkRequest 校验单个分片上传参数是否合法。
func (d *aliyunOSSDriver) validateUploadChunkRequest(request uploadChunkRequest) error {
	if err := d.validateUploadSession(request.session); err != nil {
		return err
	}
	if request.body == nil || request.size <= 0 || request.index <= 0 {
		return errInvalidUploadChunk
	}
	return nil
}

// statByProviderKey 根据底层厂商对象路径获取对象元信息。
func (d *aliyunOSSDriver) statByProviderKey(ctx context.Context, providerKey string) (*ObjectInfo, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	props, err := d.bucket.GetObjectDetailedMeta(providerKey)
	if err != nil {
		return nil, fmt.Errorf("get object meta failed: %w", err)
	}

	var size int64
	if sizeText := props.Get("Content-Length"); sizeText != "" {
		_, _ = fmt.Sscanf(sizeText, "%d", &size)
	}

	var lastModified time.Time
	if modifiedText := props.Get("Last-Modified"); modifiedText != "" {
		lastModified, _ = time.Parse(time.RFC1123, modifiedText)
	}

	return &ObjectInfo{
		Key:          stripProviderPrefix(d.config.Prefix, providerKey),
		ProviderKey:  providerKey,
		Size:         size,
		LastModified: lastModified,
		ETag:         strings.Trim(props.Get("ETag"), "\""),
		ContentType:  props.Get("Content-Type"),
		URL:          d.objectURL(providerKey),
		Metadata:     metadataFromHeaders(props, "x-oss-meta-"),
	}, nil
}

// objectFromListObject 将 OSS 列表结果转换成统一对象信息。
func (d *aliyunOSSDriver) objectFromListObject(object oss.ObjectProperties) ObjectInfo {
	return ObjectInfo{
		Key:          stripProviderPrefix(d.config.Prefix, object.Key),
		ProviderKey:  object.Key,
		Size:         object.Size,
		LastModified: object.LastModified,
		ETag:         strings.Trim(object.ETag, "\""),
		URL:          d.objectURL(object.Key),
	}
}

// objectURL 根据当前配置生成对象访问地址。
func (d *aliyunOSSDriver) objectURL(providerKey string) string {
	if strings.TrimSpace(d.config.Domain) != "" {
		return strings.TrimRight(ensureScheme(d.config.Domain, d.config.UseHTTPS), "/") + "/" + providerKey
	}

	endpoint := strings.TrimRight(d.config.Endpoint, "/")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	scheme := "https"
	if !d.config.UseHTTPS {
		scheme = "http"
	}

	return fmt.Sprintf("%s://%s.%s/%s", scheme, d.config.Bucket, endpoint, providerKey)
}

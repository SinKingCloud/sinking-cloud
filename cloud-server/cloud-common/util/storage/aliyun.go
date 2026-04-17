package storage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// aliyunOSSDriver 是阿里云 OSS 的对象存储驱动实现。
type aliyunOSSDriver struct {
	client *oss.Client
	bucket *oss.Bucket
	config Config
}

// newAliyunOSSDriver 创建阿里云 OSS 驱动。
func newAliyunOSSDriver(config Config) (*aliyunOSSDriver, error) {
	if strings.TrimSpace(config.Endpoint) == "" ||
		strings.TrimSpace(config.Bucket) == "" ||
		strings.TrimSpace(config.AccessKeyID) == "" ||
		strings.TrimSpace(config.AccessKeySecret) == "" {
		return nil, fmt.Errorf("aliyun oss config is incomplete")
	}

	config.UseHTTPS = resolveHTTPSPreference(config.UseHTTPS, config.Endpoint, config.Domain)
	config.Prefix = normalizePrefix(config.Prefix)

	clientOptions := make([]oss.ClientOption, 0, 4)
	if strings.TrimSpace(config.SecurityToken) != "" {
		clientOptions = append(clientOptions, oss.SecurityToken(config.SecurityToken))
	}
	if strings.TrimSpace(config.Region) != "" {
		clientOptions = append(clientOptions, oss.Region(config.Region), oss.AuthVersion(oss.AuthV4))
	}

	client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret, clientOptions...)
	if err != nil {
		return nil, fmt.Errorf("create aliyun oss client failed: %w", err)
	}

	bucket, err := client.Bucket(config.Bucket)
	if err != nil {
		return nil, fmt.Errorf("open aliyun oss bucket failed: %w", err)
	}

	return &aliyunOSSDriver{
		client: client,
		bucket: bucket,
		config: config,
	}, nil
}

// provider 返回当前驱动的 provider 类型。
func (d *aliyunOSSDriver) provider() Provider {
	return ProviderAliyunOSS
}

// capabilities 返回阿里云 OSS 驱动支持的能力集合。
func (d *aliyunOSSDriver) capabilities() Capabilities {
	return Capabilities{
		SmallUpload:       true,
		LargeUpload:       true,
		ResumeUpload:      true,
		UploadStatusQuery: true,
		SignedURL:         true,
		CopyObject:        true,
		BatchDelete:       true,
		Protocols:         []UploadProtocol{UploadProtocolMultipart},
		ChunkPolicy: ChunkPolicy{
			DefaultSize: DefaultChunkSize,
			MinSize:     MinCompatibleChunkSize,
			MaxChunks:   10000,
		},
	}
}

// upload 执行简单对象上传。
func (d *aliyunOSSDriver) upload(ctx context.Context, request UploadRequest) (*ObjectInfo, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if request.Body == nil {
		return nil, ErrUploadBodyRequired
	}

	providerKey, err := buildProviderKey(d.config.Prefix, request.Key)
	if err != nil {
		return nil, err
	}

	options := d.buildObjectOptions(
		detectContentType(request.ContentType, request.Key),
		request.CacheControl,
		request.ContentDisposition,
		request.Metadata,
	)

	if err := d.bucket.PutObject(providerKey, request.Body, options...); err != nil {
		return nil, fmt.Errorf("upload object failed: %w", err)
	}

	return d.statByProviderKey(ctx, providerKey)
}

// startUploadSession 初始化阿里云 OSS multipart 上传会话。
func (d *aliyunOSSDriver) startUploadSession(ctx context.Context, request StartUploadSessionRequest) (*UploadSession, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	providerKey, err := buildProviderKey(d.config.Prefix, request.Key)
	if err != nil {
		return nil, err
	}

	chunkSize := normalizeChunkSize(request.ChunkSize, d.capabilities().ChunkPolicy, request.Size)
	contentType := detectContentType(request.ContentType, request.Key)
	options := d.buildObjectOptions(
		contentType,
		request.CacheControl,
		request.ContentDisposition,
		request.Metadata,
	)

	result, err := d.bucket.InitiateMultipartUpload(providerKey, options...)
	if err != nil {
		return nil, fmt.Errorf("init upload session failed: %w", err)
	}

	return &UploadSession{
		Provider:           ProviderAliyunOSS,
		Protocol:           UploadProtocolMultipart,
		Bucket:             d.config.Bucket,
		Key:                request.Key,
		ProviderKey:        providerKey,
		SessionID:          result.UploadID,
		Size:               request.Size,
		ChunkSize:          chunkSize,
		ContentType:        contentType,
		CacheControl:       request.CacheControl,
		ContentDisposition: request.ContentDisposition,
		Metadata:           cloneMetadata(request.Metadata),
		CreatedAt:          time.Now(),
	}, nil
}

// uploadChunk 上传单个 multipart 分片。
func (d *aliyunOSSDriver) uploadChunk(ctx context.Context, request UploadChunkRequest) (*UploadedChunk, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if err := d.validateUploadChunkRequest(request); err != nil {
		return nil, err
	}

	part, err := d.bucket.UploadPart(d.uploadSessionHandle(request.Session), request.Body, request.Size, request.Index)
	if err != nil {
		return nil, fmt.Errorf("upload chunk failed: %w", err)
	}

	return &UploadedChunk{
		Index:        request.Index,
		Offset:       request.Offset,
		Size:         request.Size,
		ETag:         strings.Trim(part.ETag, "\""),
		LastModified: time.Now(),
	}, nil
}

// getUploadStatus 查询 multipart 上传会话状态。
func (d *aliyunOSSDriver) getUploadStatus(ctx context.Context, session UploadSession) (*UploadStatus, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if err := d.validateUploadSession(session); err != nil {
		return nil, err
	}

	result, err := d.bucket.ListUploadedParts(d.uploadSessionHandle(session))
	if err != nil {
		return nil, fmt.Errorf("get upload status failed: %w", err)
	}

	chunks := make([]UploadedChunk, 0, len(result.UploadedParts))
	var committedBytes int64
	for _, part := range result.UploadedParts {
		size := int64(part.Size)
		chunks = append(chunks, UploadedChunk{
			Index:        part.PartNumber,
			Size:         size,
			ETag:         strings.Trim(part.ETag, "\""),
			LastModified: part.LastModified,
		})
		committedBytes += size
	}

	return &UploadStatus{
		Session:        session,
		UploadedChunks: chunks,
		CommittedBytes: committedBytes,
	}, nil
}

// completeUpload 完成 multipart 上传。
func (d *aliyunOSSDriver) completeUpload(ctx context.Context, request CompleteUploadRequest) (*ObjectInfo, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if err := d.validateUploadSession(request.Session); err != nil {
		return nil, err
	}
	if len(request.Chunks) == 0 {
		return nil, ErrInvalidCompleteUpload
	}

	parts := make([]oss.UploadPart, 0, len(request.Chunks))
	for _, chunk := range request.Chunks {
		if chunk.Index <= 0 || strings.TrimSpace(chunk.ETag) == "" {
			return nil, ErrInvalidCompleteUpload
		}
		parts = append(parts, oss.UploadPart{
			PartNumber: chunk.Index,
			ETag:       chunk.ETag,
		})
	}

	sort.Sort(oss.UploadParts(parts))
	if _, err := d.bucket.CompleteMultipartUpload(d.uploadSessionHandle(request.Session), parts); err != nil {
		return nil, fmt.Errorf("complete upload failed: %w", err)
	}

	return d.statByProviderKey(ctx, request.Session.ProviderKey)
}

// abortUpload 终止 multipart 上传。
func (d *aliyunOSSDriver) abortUpload(ctx context.Context, session UploadSession) error {
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

// delete 删除单个对象。
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

// deleteMany 批量删除对象。
func (d *aliyunOSSDriver) deleteMany(ctx context.Context, keys []string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}

	providerKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		providerKey, err := buildProviderKey(d.config.Prefix, key)
		if err != nil {
			return err
		}
		providerKeys = append(providerKeys, providerKey)
	}

	if _, err := d.bucket.DeleteObjects(providerKeys); err != nil {
		return fmt.Errorf("batch delete objects failed: %w", err)
	}
	return nil
}

// stat 获取单个对象元信息。
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

// exists 判断对象是否存在。
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

// list 分页列举对象。
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

// copy 复制对象。
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

// signedURL 生成对象的临时访问地址。
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

// buildObjectOptions 构建 OSS 上传相关的可选参数。
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

// uploadSessionHandle 将通用上传会话转换成 OSS SDK 所需的上传句柄。
func (d *aliyunOSSDriver) uploadSessionHandle(session UploadSession) oss.InitiateMultipartUploadResult {
	return oss.InitiateMultipartUploadResult{
		Bucket:   d.config.Bucket,
		Key:      session.ProviderKey,
		UploadID: session.SessionID,
	}
}

// validateUploadSession 校验阿里云 OSS 上传会话。
func (d *aliyunOSSDriver) validateUploadSession(session UploadSession) error {
	if strings.TrimSpace(session.ProviderKey) == "" || strings.TrimSpace(session.SessionID) == "" {
		return ErrInvalidUploadSession
	}
	if session.Protocol != "" && session.Protocol != UploadProtocolMultipart {
		return ErrInvalidUploadSession
	}
	return nil
}

// validateUploadChunkRequest 校验阿里云 OSS 分块上传请求。
func (d *aliyunOSSDriver) validateUploadChunkRequest(request UploadChunkRequest) error {
	if err := d.validateUploadSession(request.Session); err != nil {
		return err
	}
	if request.Body == nil || request.Size <= 0 || request.Index <= 0 {
		return ErrInvalidUploadChunk
	}
	return nil
}

// statByProviderKey 使用真实 provider key 获取对象元信息。
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
		Metadata:     metadataFromHeaders(props),
	}, nil
}

// objectFromListObject 将 OSS 列表条目转换为统一对象模型。
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

// objectURL 构建对象访问地址。
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

// metadataFromHeaders 提取 OSS 对象元数据。
func metadataFromHeaders(headers http.Header) map[string]string {
	if len(headers) == 0 {
		return nil
	}

	metadata := make(map[string]string)
	for key, values := range headers {
		if len(values) == 0 {
			continue
		}

		lowerKey := strings.ToLower(key)
		if !strings.HasPrefix(lowerKey, "x-oss-meta-") {
			continue
		}

		metadata[strings.TrimPrefix(lowerKey, "x-oss-meta-")] = values[0]
	}

	if len(metadata) == 0 {
		return nil
	}
	return metadata
}

// resolveHTTPSPreference 根据域名和 endpoint 推断是否默认使用 HTTPS。
func resolveHTTPSPreference(useHTTPS bool, values ...string) bool {
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		switch {
		case strings.HasPrefix(value, "http://"):
			return false
		case strings.HasPrefix(value, "https://"):
			return true
		}
	}

	if useHTTPS {
		return true
	}
	return true
}

// ensureScheme 确保域名包含协议头。
func ensureScheme(raw string, useHTTPS bool) string {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}

	scheme := "https://"
	if !useHTTPS {
		scheme = "http://"
	}

	if parsed, err := url.Parse(raw); err == nil && parsed.Scheme != "" {
		return raw
	}

	return scheme + raw
}
